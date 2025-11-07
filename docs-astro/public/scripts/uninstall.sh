#!/usr/bin/env bash
set -euo pipefail

# Promptext Uninstall Script
# Removes promptext binary, aliases, and PATH modifications

VERSION="1.0.0"

# Configuration
BINARY_NAME="promptext"
DEFAULT_INSTALL_DIR="$HOME/.local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    echo -e "${GREEN}→${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}⚠${NC} $1"
}

print_error() {
    echo -e "${RED}✗${NC} $1" >&2
}

print_success() {
    echo -e "${GREEN}✓${NC} $1"
}

# Help text
show_help() {
    cat << EOF
promptext uninstall script v${VERSION}

Usage: $0 [options]

Options:
    -h, --help              Show this help message
    -d, --dir DIR          Custom installation directory to check
    -y, --yes              Skip confirmation prompt

Examples:
    $0                      # Uninstall with confirmation
    $0 -y                   # Uninstall without confirmation
    $0 --dir ~/bin         # Uninstall from custom directory
EOF
    exit 0
}

# Parse command line arguments
INSTALL_DIR="$DEFAULT_INSTALL_DIR"
SKIP_CONFIRM=false

while [[ $# -gt 0 ]]; do
    case $1 in
        -h|--help)
            show_help
            ;;
        -d|--dir)
            INSTALL_DIR="$2"
            shift
            ;;
        -y|--yes)
            SKIP_CONFIRM=true
            ;;
        *)
            print_error "Unknown option: $1"
            show_help
            exit 1
            ;;
    esac
    shift
done

# Find promptext installations
find_installations() {
    local found=()

    # Check default location
    if [ -f "$DEFAULT_INSTALL_DIR/$BINARY_NAME" ]; then
        found+=("$DEFAULT_INSTALL_DIR/$BINARY_NAME")
    fi

    # Check custom location if different
    if [ "$INSTALL_DIR" != "$DEFAULT_INSTALL_DIR" ] && [ -f "$INSTALL_DIR/$BINARY_NAME" ]; then
        found+=("$INSTALL_DIR/$BINARY_NAME")
    fi

    # Check if in PATH
    if command -v $BINARY_NAME >/dev/null 2>&1; then
        local path_location=$(command -v $BINARY_NAME)
        # Add if not already in found list
        if [[ ! " ${found[@]} " =~ " ${path_location} " ]]; then
            found+=("$path_location")
        fi
    fi

    echo "${found[@]}"
}

# Remove shell aliases
remove_aliases() {
    local removed=false
    local shell_configs=(
        "$HOME/.bashrc"
        "$HOME/.zshrc"
        "$HOME/.config/fish/config.fish"
    )

    print_info "Removing shell aliases..."

    for config in "${shell_configs[@]}"; do
        if [ -f "$config" ]; then
            if grep -q "alias prx=" "$config" 2>/dev/null; then
                # Use different sed syntax based on OS
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    sed -i '' '/alias prx=/d' "$config"
                else
                    sed -i '/alias prx=/d' "$config"
                fi
                print_success "Removed alias from $(basename $config)"
                removed=true
            fi
        fi
    done

    if [ "$removed" = false ]; then
        print_info "No aliases found in shell configs"
    fi
}

# Remove PATH modifications
remove_path_entries() {
    local removed=false
    local shell_configs=(
        "$HOME/.bashrc"
        "$HOME/.zshrc"
    )

    print_info "Checking for PATH modifications..."

    for config in "${shell_configs[@]}"; do
        if [ -f "$config" ]; then
            if grep -q "PATH.*$INSTALL_DIR" "$config" 2>/dev/null; then
                # Use different sed syntax based on OS
                if [[ "$OSTYPE" == "darwin"* ]]; then
                    sed -i '' "\|export PATH.*$INSTALL_DIR|d" "$config"
                else
                    sed -i "\|export PATH.*$INSTALL_DIR|d" "$config"
                fi
                print_success "Removed PATH entry from $(basename $config)"
                removed=true
            fi
        fi
    done

    if [ "$removed" = false ]; then
        print_info "No PATH modifications found"
    fi
}

# Remove Homebrew symlink (if exists)
remove_homebrew_symlink() {
    local homebrew_bin
    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if [ -d "/opt/homebrew/bin" ]; then
            homebrew_bin="/opt/homebrew/bin"
        else
            homebrew_bin="/usr/local/bin"
        fi
    else
        # Linux
        homebrew_bin="/home/linuxbrew/.linuxbrew/bin"
    fi

    local prx_link="$homebrew_bin/prx"
    if [ -L "$prx_link" ]; then
        rm "$prx_link"
        print_success "Removed Homebrew prx symlink"
    fi
}

# Main uninstall process
main() {
    echo "Promptext Uninstaller v${VERSION}"
    echo ""

    # Find installations
    local installations=($(find_installations))

    if [ ${#installations[@]} -eq 0 ]; then
        print_warning "No promptext installation found"
        echo ""
        print_info "Checked locations:"
        echo "  - $DEFAULT_INSTALL_DIR/$BINARY_NAME"
        echo "  - System PATH"
        exit 0
    fi

    # Show what will be removed
    print_info "Found promptext installation(s):"
    for install in "${installations[@]}"; do
        echo "  - $install"
    done
    echo ""

    # Confirmation prompt
    if [ "$SKIP_CONFIRM" = false ]; then
        read -p "Remove promptext? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            print_info "Uninstall cancelled"
            exit 0
        fi
    fi

    # Remove binaries
    print_info "Removing binaries..."
    for install in "${installations[@]}"; do
        if [ -f "$install" ]; then
            rm -f "$install"
            print_success "Removed $install"
        fi
    done

    # Remove aliases and PATH entries
    remove_aliases
    remove_path_entries
    remove_homebrew_symlink

    echo ""
    print_success "Promptext has been uninstalled!"
    echo ""
    print_info "Note: You may need to restart your terminal for changes to take effect"
    print_info "To reinstall: curl -sSL promptext.sh/install | bash"
}

# Run main function
main
