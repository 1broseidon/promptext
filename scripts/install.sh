#!/usr/bin/env bash
set -euo pipefail

# Script version
VERSION="0.2.0"

# Help text
show_help() {
    cat << EOF
promptext installer v${VERSION}

Usage: $0 [options]

Options:
    -h, --help              Show this help message
    -d, --dir DIR          Install to specific directory (default: /usr/local/bin)
    -u, --uninstall        Uninstall promptext
    --user                 Install for current user only (no sudo required)
    --no-alias             Skip alias creation
    --no-verify            Skip checksum verification (not recommended)
    --insecure            Skip HTTPS certificate validation (not recommended)

Examples:
    $0                     # Install system-wide (requires sudo)
    $0 --user             # Install for current user only
    $0 --dir ~/bin        # Install to custom directory
    $0 --uninstall        # Uninstall promptext
EOF
    exit 0
}

# Parse command line arguments
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                ;;
            -d|--dir)
                INSTALL_DIR="$2"
                shift
                ;;
            -u|--uninstall)
                DO_UNINSTALL=true
                ;;
            --user)
                install_user_level=true
                INSTALL_PATH="$USER_INSTALL_DIR"
                ;;
            --no-alias)
                SKIP_ALIAS=true
                ;;
            --no-verify)
                SKIP_VERIFY=true
                ;;
            --insecure)
                CURL_OPTS="$CURL_OPTS -k"
                echo "⚠️  Warning: HTTPS certificate validation disabled"
                ;;
            *)
                echo "Error: Unknown option $1" >&2
                show_help
                exit 1
                ;;
        esac
        shift
    done
}

# Configuration and defaults
OWNER="1broseidon"
REPO="promptext"
BINARY_NAME="promptext"
INSTALL_DIR="/usr/local/bin"          # Default system-wide install
USER_INSTALL_DIR="$HOME/.local/bin"   # Default user-level install
CURL_OPTS="-fsSL --tlsv1.2"
DO_UNINSTALL=false
SKIP_ALIAS=false
SKIP_VERIFY=false
install_user_level=false
INSTALL_PATH="$INSTALL_DIR"

# Shell configuration files to check (ordered by preference)
SHELL_CONFIGS=(
    "$HOME/.zshrc"
    "$HOME/.bashrc"
    "$HOME/.config/fish/config.fish"
    "$HOME/.kshrc"
    "$HOME/.cshrc"
)

# Cleanup function
cleanup() {
    if [ -n "${TMP_DIR:-}" ] && [ -d "$TMP_DIR" ]; then
        rm -rf "$TMP_DIR"
    fi
}

# Error handler
error_handler() {
    echo "Error: Installation failed on line $1" >&2
    cleanup
    exit 1
}

# Set error handler
trap 'error_handler ${LINENO}' ERR

# Function to write status messages with optional progress indicator
write_status() {
    local msg="$1"
    local progress="${2:-}"
    if [ -n "$progress" ]; then
        echo -ne "→ $msg... $progress\r"
    else
        echo "→ $msg"
    fi
}

# Function to check if we have sudo access
check_sudo_access() {
    if sudo -n true 2>/dev/null; then
        return 0
    else
        return 1
    fi
}

# Function to confirm sudo usage
confirm_sudo() {
    local action="$1"
    echo "⚠️  This operation requires sudo to $action"
    
    if ! check_sudo_access; then
        echo "Error: This operation requires sudo privileges." >&2
        echo "Please run with sudo or use --user flag for user-level installation." >&2
        exit 1
    fi
    
    read -p "Continue? [y/N] " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        echo "Operation cancelled by user" >&2
        exit 1
    fi
}

# Check for required dependencies
check_dependencies() {
    local missing_deps=()

    if ! command -v curl >/dev/null 2>&1; then
        missing_deps+=("curl")
    fi

    if ! command -v tar >/dev/null 2>&1; then
        missing_deps+=("tar")
    fi

    if ! command -v sha256sum >/dev/null 2>&1; then
        missing_deps+=("sha256sum")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        echo "Error: Missing required dependencies: ${missing_deps[*]}" >&2
        echo "Please install them using your package manager." >&2
        exit 1
    fi
}

# Function to detect OS
detect_os() {
	if [[ $OSTYPE == "darwin"* ]]; then
		echo "darwin"
	else
		echo "linux"
	fi
}

# Function to detect architecture
detect_arch() {
	local arch=$(uname -m)
	case $arch in
	x86_64)
		echo "x86_64"
		;;
	aarch64 | arm64)
		echo "arm64"
		;;
	*)
		echo "Unsupported architecture: $arch" >&2
		exit 1
		;;
	esac
}

# Function to get latest release URL
get_latest_release_url() {
	local os=$1
	local arch=$2
	local os_name
	if [ "$os" == "darwin" ]; then
		os_name="Darwin"
	else
		os_name="Linux"
	fi

	echo "https://github.com/1broseidon/promptext/releases/latest/download/promptext_${os_name}_${arch}.tar.gz"
}

# Function to verify checksum
verify_checksum() {
    local file=$1
    local release_url="$2"
    local checksum_files=("checksums.txt" "SHA256SUMS" "sha256sums.txt")
    local found=false

    write_status "Verifying checksum..."
    local actual_checksum=$(sha256sum "$file" | awk '{print $1}')
    local filename=$(basename "$release_url")

    # Try different checksum files
    for checksum_file in "${checksum_files[@]}"; do
        local checksum_url="${release_url%/*}/$checksum_file"
        local response=$(curl -fsL "$checksum_url" 2>/dev/null)
        if [ $? -eq 0 ] && [ -n "$response" ]; then
            write_status "Found checksum file: $checksum_file"
            local expected_checksum=$(echo "$response" | grep -F "$filename" | awk '{print $1}')
            if [ -n "$expected_checksum" ]; then
                found=true
                if [ "$actual_checksum" != "$expected_checksum" ]; then
                    echo "Error: Checksum verification failed." >&2
                    echo "Expected: $expected_checksum" >&2
                    echo "Got: $actual_checksum" >&2
                    exit 1
                fi
                write_status "Checksum verification successful"
                return 0
            fi
        fi
    done

    if [ "$found" = false ]; then
        if [ "$SKIP_VERIFY" = true ]; then
            write_status "Skipping checksum verification (--no-verify)"
            return 0
        fi
        echo "⚠️  Warning: No checksum file found. Use --no-verify to skip verification." >&2
        exit 1
    fi
}

# Function to uninstall promptext
uninstall_promptext() {
    write_status "Uninstalling promptext..."

    local install_path
    if [ "$install_user_level" = true ]; then
        install_path="$USER_INSTALL_DIR/$BINARY_NAME"
    else
        install_path="$INSTALL_DIR/$BINARY_NAME"
    fi

    if [ -f "$install_path" ]; then
        if [ "$install_user_level" = false ]; then
            if [ ! -w "$install_path" ]; then
                confirm_sudo "uninstall from $install_path"
            fi
            sudo rm -f "$install_path"
        else
            rm -f "$install_path"
        fi
    fi

    # Remove alias (simplified, can be improved)
    sed -i '/alias prx=promptext/d' ~/.bashrc 2>/dev/null
    sed -i '/alias prx=promptext/d' ~/.zshrc 2>/dev/null

    write_status "Promptext uninstalled."
    exit 0
}

# Main installation function
install_promptext() {
    write_status "Installing promptext..."
    write_status "OS: $OS"
    write_status "Architecture: $ARCH"
    write_status "Installation directory: $INSTALL_PATH"

    # Early check for sudo if needed
    if [ "$install_user_level" = false ]; then
        if [ ! -w "$INSTALL_PATH" ]; then
            confirm_sudo "install to $INSTALL_PATH"
        fi
    fi

    # Check installation directory permissions
    if [ ! -w "$(dirname "$INSTALL_PATH")" ] && [ "$install_user_level" = false ]; then
        confirm_sudo "install to $INSTALL_PATH"
    fi

    # Create installation directory
    if [ ! -d "$INSTALL_PATH" ]; then
        mkdir -p "$INSTALL_PATH" || {
            echo "Error: Failed to create installation directory. Try using --user or sudo." >&2
            exit 1
        }
    fi

    # Download and verify
    write_status "Downloading latest release..."
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR" || exit 1
    trap cleanup EXIT

    RELEASE_URL=$(get_latest_release_url "$OS" "$ARCH")
    if ! curl $CURL_OPTS "$RELEASE_URL" -o promptext.tar.gz; then
        echo "Error: Failed to download release. Try --insecure if using a proxy." >&2
        exit 1
    fi

    if [ "$SKIP_VERIFY" = false ]; then
        verify_checksum "promptext.tar.gz" "$RELEASE_URL"
    else
        write_status "Skipping checksum verification (--no-verify)"
    fi

    write_status "Extracting archive..."
    if ! tar xzf promptext.tar.gz; then
        echo "Error: Failed to extract archive. Try downloading again." >&2
        exit 1
    fi

    # Install binary
    write_status "Installing binary..."
    if [ -f "$INSTALL_PATH/$BINARY_NAME" ]; then
        write_status "Removing existing installation..."
        if [ "$install_user_level" = false ]; then
            confirm_sudo "remove existing installation"
            sudo rm -f "$INSTALL_PATH/$BINARY_NAME"
        else
            rm -f "$INSTALL_PATH/$BINARY_NAME"
        fi
    fi

    if [ "$install_user_level" = false ]; then
        sudo mv promptext "$INSTALL_PATH/"
        sudo chmod +x "$INSTALL_PATH/$BINARY_NAME"
    else
        mv promptext "$INSTALL_PATH/"
        chmod +x "$INSTALL_PATH/$BINARY_NAME"
    fi

    # Update PATH if needed
    if [ "$install_user_level" = true ]; then
        write_status "Updating PATH..."
        local path_updated=false
        for config in "${SHELL_CONFIGS[@]}"; do
            if [ -f "$config" ]; then
                if ! grep -q "export PATH=.*:$INSTALL_PATH" "$config"; then
                    echo "export PATH=\"\$PATH:$INSTALL_PATH\"" >> "$config"
                    path_updated=true
                    break
                fi
            fi
        done
        if [ "$path_updated" = false ]; then
            echo "⚠️  Warning: Could not update PATH. Add '$INSTALL_PATH' to your PATH manually." >&2
        fi
    fi

    # Create alias if requested
    if [ "$SKIP_ALIAS" = false ]; then
        write_status "Creating alias..."
        local alias_updated=false
        for config in "${SHELL_CONFIGS[@]}"; do
            if [ -f "$config" ]; then
                if ! grep -q "alias prx=$BINARY_NAME" "$config"; then
                    echo "alias prx=$BINARY_NAME" >> "$config"
                    alias_updated=true
                    break
                fi
            fi
        done
        if [ "$alias_updated" = false ]; then
            echo "⚠️  Warning: Could not create alias. You can still use '$BINARY_NAME' directly." >&2
        fi
    fi

    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" -v)
        write_status "Installation verified: $version"
    else
        echo "⚠️  Warning: Installation completed but binary not found in PATH" >&2
        echo "Try restarting your terminal or adding '$INSTALL_PATH' to your PATH" >&2
    fi

    echo "✨ Installation complete!"
    if [ "$SKIP_ALIAS" = false ] && [ "$alias_updated" = true ]; then
        echo "You can use either '$BINARY_NAME' or 'prx' command after restarting your terminal."
    else
        echo "You can use the '$BINARY_NAME' command after restarting your terminal."
    fi
    echo "To uninstall, run: $0 --uninstall"
}

# Main script
check_dependencies
OS=$(detect_os)
ARCH=$(detect_arch)

# Parse command line arguments
parse_args "$@"

# Handle uninstall
if [ "$DO_UNINSTALL" = true ]; then
    uninstall_promptext
fi

# Proceed with installation
install_promptext
