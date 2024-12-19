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

# Function to check if alias exists and what it points to
check_alias_exists() {
    local alias_target=""
    
    # Check in .zshrc
    if [ -f "$HOME/.zshrc" ]; then
        alias_target=$(grep "alias prx=" "$HOME/.zshrc" | cut -d= -f2- 2>/dev/null)
    fi
    
    # Check in .bashrc if not found in .zshrc
    if [ -z "$alias_target" ] && [ -f "$HOME/.bashrc" ]; then
        alias_target=$(grep "alias prx=" "$HOME/.bashrc" | cut -d= -f2- 2>/dev/null)
    fi
    
    # Check if prx exists as a command
    if [ -z "$alias_target" ] && command -v prx >/dev/null 2>&1; then
        alias_target=$(command -v prx)
    fi
    
    if [ -n "$alias_target" ]; then
        # Remove quotes and whitespace
        alias_target=$(echo "$alias_target" | tr -d '"' | tr -d "'" | xargs)
        if [ "$alias_target" = "promptext" ]; then
            echo "exists_and_matches"
        else
            echo "exists_different:$alias_target"
        fi
    else
        echo "not_exists"
    fi
}

# Try to create alias if it doesn't exist
ALIAS_STATUS=$(check_alias_exists)
if [ "$ALIAS_STATUS" = "exists_and_matches" ]; then
    echo "✨ Installation complete! The 'prx' alias is already correctly set to promptext."
elif [[ "$ALIAS_STATUS" == exists_different:* ]]; then
    CURRENT_TARGET="${ALIAS_STATUS#exists_different:}"
    echo "⚠️  Note: 'prx' alias exists but points to: $CURRENT_TARGET"
    echo "✨ Installation complete! You can use the 'promptext' command."
elif [ "$ALIAS_STATUS" = "not_exists" ]; then
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
