#!/bin/bash

# Function to detect OS
detect_os() {
    if [[ "$OSTYPE" == "darwin"* ]]; then
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
        aarch64|arm64)
            echo "arm64"
            ;;
        *)
            echo "Unsupported architecture: $arch" >&2
            exit 1
            ;;
    esac
}

# Function to get latest release URL
get_latest_release() {
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

# Detect system information
OS=$(detect_os)
ARCH=$(detect_arch)
RELEASE_URL=$(get_latest_release $OS $ARCH)
INSTALL_DIR="/usr/local/bin"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

echo "Downloading latest promptext release..."
if ! curl -sL "$RELEASE_URL" -o promptext.tar.gz; then
    echo "Failed to download release"
    exit 1
fi

echo "Extracting archive..."
tar xzf promptext.tar.gz

echo "Installing promptext..."
sudo mv promptext "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/promptext"

echo "Cleaning up..."
cd - > /dev/null
rm -rf "$TMP_DIR"

echo "Installation complete! You can now use 'promptext' command."
