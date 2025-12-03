#!/bin/bash
# scripts/bootstrap.sh
# One-liner installer for Gip (Binary Release)
# Usage: curl -sL https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.sh | bash

set -e

REPO_OWNER="iamHrithikRaj"
REPO_NAME="gip"
INSTALL_DIR="$HOME/.local/bin"

# Detect OS
OS="$(uname -s)"
case "${OS}" in
    Linux*)     ASSET_NAME="gip-linux-x64.tar.gz";;
    Darwin*)    ASSET_NAME="gip-macos-universal.tar.gz";;
    *)          echo "Unsupported OS: ${OS}"; exit 1;;
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
