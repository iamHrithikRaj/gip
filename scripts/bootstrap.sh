#!/bin/bash
# scripts/bootstrap.sh
# One-liner installer for Gip (Binary Release)
# Usage: curl -sL https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.sh | bash

set -e

REPO_OWNER="iamHrithikRaj"
REPO_NAME="gip"
INSTALL_DIR="$HOME/.local/bin"

# Detect OS and Architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

case "${OS}" in
    Linux*)
        if [ "$ARCH" = "x86_64" ]; then
            ASSET_NAME="gip-linux-x64.tar.gz"
        else
            echo "Unsupported Linux architecture: ${ARCH}"
            exit 1
        fi
        ;;
    Darwin*)
        if [ "$ARCH" = "arm64" ]; then
            ASSET_NAME="gip-macos-arm64.tar.gz"
        elif [ "$ARCH" = "x86_64" ]; then
            ASSET_NAME="gip-macos-x64.tar.gz"
        else
            echo "Unsupported macOS architecture: ${ARCH}"
            exit 1
        fi
        ;;
    *)
        echo "Unsupported OS: ${OS}"
        exit 1
        ;;
esac

echo -e "\033[0;36mFetching latest release info...\033[0m"
LATEST_RELEASE_URL="https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
DOWNLOAD_URL=$(curl -s "$LATEST_RELEASE_URL" | grep "browser_download_url.*$ASSET_NAME" | cut -d : -f 2,3 | tr -d \" | xargs)

if [ -z "$DOWNLOAD_URL" ]; then
    echo "Error: Could not find $ASSET_NAME in the latest release."
    exit 1
fi

TEMP_DIR=$(mktemp -d)
TEMP_FILE="$TEMP_DIR/$ASSET_NAME"

echo -e "\033[0;36mDownloading Gip...\033[0m"
curl -sL "$DOWNLOAD_URL" -o "$TEMP_FILE"

echo -e "\033[0;36mInstalling to $INSTALL_DIR...\033[0m"
mkdir -p "$INSTALL_DIR"
tar -xzf "$TEMP_FILE" -C "$INSTALL_DIR"

# Check for conflict and offer rename
if command -v gip >/dev/null 2>&1; then
    EXISTING_GIP=$(command -v gip)
    # Don't warn if it's the one we just installed
    if [ "$EXISTING_GIP" != "$INSTALL_DIR/gip" ]; then
        echo -e "\033[0;33m\n[!] Detected existing 'gip' command at $EXISTING_GIP\033[0m"
        
        # Read from /dev/tty to allow interaction even when piped
        if [ -t 0 ]; then
            read -p "Do you want to override/keep 'gip'? (Y/N) - N will rename installed binary to 'git++' [Default: Y] " choice
        else
            # Fallback if no TTY available (e.g. non-interactive CI)
            echo "Non-interactive session detected. Keeping 'gip'."
            choice="Y"
        fi

        case "$choice" in 
          n|N ) 
            echo -e "\033[0;36mRenaming installed binary to 'git++'...\033[0m"
            mv "$INSTALL_DIR/gip" "$INSTALL_DIR/git++"
            echo -e "\033[0;32mInstalled as 'git++'. You can run it using: git++\033[0m"
            ;;
          * ) 
            echo "Keeping 'gip'."
            ;;
        esac
    fi
fi

# Cleanup
rm -rf "$TEMP_DIR"

echo -e "\033[0;32m\nGip installed successfully!\033[0m"

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "\033[0;33mWarning: $INSTALL_DIR is not in your PATH.\033[0m"
    echo "Add the following to your shell config (.bashrc, .zshrc, etc.):"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo "Run 'gip --version' to verify."
