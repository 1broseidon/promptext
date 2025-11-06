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
    -d, --dir DIR          Install to specific directory (default: ~/.local/bin)
    -u, --uninstall        Uninstall promptext
    --no-alias             Skip alias creation
    --no-verify            Skip checksum verification (not recommended)
    --insecure            Skip HTTPS certificate validation (not recommended)

Examples:
    $0                     # Install to ~/.local/bin
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
INSTALL_DIR="$HOME/.local/bin"        # Default installation directory
CURL_OPTS="-fsSL --tlsv1.2"
DO_UNINSTALL=false
SKIP_ALIAS=false
SKIP_VERIFY=false
alias_updated=false                   # Initialize alias status

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

	# GoReleaser uses lowercase OS names and version in filename
	# Format: promptext_{version}_{os}_{arch}.tar.gz
	# We use "latest" which redirects to the actual version
	echo "https://github.com/1broseidon/promptext/releases/latest/download/promptext_${os}_${arch}.tar.gz"
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

    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        rm -f "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Remove alias
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
    write_status "Installation directory: $INSTALL_DIR"

    # Create installation directory
    if [ ! -d "$INSTALL_DIR" ]; then
        mkdir -p "$INSTALL_DIR" || {
            echo "Error: Failed to create installation directory." >&2
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

    # Install binary (GoReleaser builds it as "prx")
    write_status "Installing binary..."
    if [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        write_status "Removing existing installation..."
        rm -f "$INSTALL_DIR/$BINARY_NAME"
    fi

    # Binary is named "prx" in the archive, rename to "promptext"
    mv prx "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    # Add to PATH if needed
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.bashrc"
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >> "$HOME/.zshrc"
        write_status "Added $INSTALL_DIR to PATH"
    fi

    # Add alias
    if [ "$SKIP_ALIAS" = false ]; then
        echo "alias prx='$BINARY_NAME'" >> "$HOME/.bashrc"
        echo "alias prx='$BINARY_NAME'" >> "$HOME/.zshrc"
        write_status "Added 'prx' alias"
    fi

    # Verify installation
    if command -v "$BINARY_NAME" >/dev/null 2>&1; then
        local version
        version=$("$BINARY_NAME" -v)
        write_status "Installation verified: $version"
    else
        echo "⚠️  Warning: Installation completed but binary not found in PATH" >&2
        echo "Try restarting your terminal or run:" >&2
        echo "export PATH=\"\$PATH:$INSTALL_DIR\"" >&2
    fi

    echo "✨ Installation complete!"
    echo "You can use either '$BINARY_NAME' or 'prx' command after restarting your terminal."
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
