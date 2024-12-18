#!/bin/bash

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

echo "Installing promptext..."
echo "OS: $OS"
echo "Architecture: $ARCH"
echo "Installation directory: $INSTALL_DIR"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

echo "Downloading latest promptext release..."
if ! curl -#L "$RELEASE_URL" -o promptext.tar.gz; then
	echo "Failed to download release"
	exit 1
fi

echo "Extracting archive..."
tar xzf promptext.tar.gz

echo "Installing promptext..."
# Remove existing installation if present
if [ -f "$INSTALL_DIR/promptext" ]; then
	echo "Removing existing installation..."
	sudo rm "$INSTALL_DIR/promptext"
fi

echo "Moving binary to $INSTALL_DIR..."
sudo mv promptext "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/promptext"

echo "Cleaning up..."
cd - >/dev/null
rm -rf "$TMP_DIR"

# Function to check if alias exists
check_alias_exists() {
    if [ -f "$HOME/.bashrc" ] && grep -q "alias prx=" "$HOME/.bashrc"; then
        return 0
    fi
    if [ -f "$HOME/.zshrc" ] && grep -q "alias prx=" "$HOME/.zshrc"; then
        return 0
    fi
    if command -v prx >/dev/null 2>&1; then
        return 0
    fi
    return 1
}

# Try to create alias if it doesn't exist
if check_alias_exists; then
    echo "⚠️  Note: 'prx' alias already exists on your system."
    echo "✨ Installation complete! You can use the 'promptext' command."
else
    # Determine the appropriate shell config file
    SHELL_CONFIG=""
    if [ -f "$HOME/.zshrc" ]; then
        SHELL_CONFIG="$HOME/.zshrc"
    elif [ -f "$HOME/.bashrc" ]; then
        SHELL_CONFIG="$HOME/.bashrc"
    fi

    if [ -n "$SHELL_CONFIG" ]; then
        echo "alias prx=promptext" >> "$SHELL_CONFIG"
        echo "✨ Installation complete! You can use either 'promptext' or 'prx' command."
        echo "Note: Please restart your terminal or run 'source $SHELL_CONFIG' to use the 'prx' alias."
    else
        echo "⚠️  Note: Could not create 'prx' alias (shell config file not found)."
        echo "✨ Installation complete! You can use the 'promptext' command."
    fi
fi
