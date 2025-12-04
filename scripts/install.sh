#!/bin/bash
# scripts/install.sh
# Installs Gip from source

set -e

echo -e "\033[0;36mInstalling Gip...\033[0m"
echo "Performing user-space installation (no sudo required)."

# Check prerequisites
if ! command -v cargo &> /dev/null; then
    echo "Error: Cargo is required but not found. Please install the toolchain."
    exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
INSTALL_DIR="$HOME/.local/bin"

# Build
echo -e "\033[0;33mBuilding Release...\033[0m"
cd "$REPO_ROOT"
cargo build --release

# Install
echo -e "\033[0;33mInstalling to $INSTALL_DIR...\033[0m"
mkdir -p "$INSTALL_DIR"
cp "$REPO_ROOT/target/release/gip" "$INSTALL_DIR/gip"

# Check for conflict and offer rename
if command -v gip >/dev/null 2>&1; then
    EXISTING_GIP=$(command -v gip)
    # Don't warn if it's the one we just installed
    if [ "$EXISTING_GIP" != "$INSTALL_DIR/gip" ]; then
        echo -e "\033[0;33m\n[!] Detected existing 'gip' command at $EXISTING_GIP\033[0m"
        
        read -p "Do you want to override/keep 'gip'? (Y/N) - N will rename installed binary to 'git++' [Default: Y] " choice
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

echo -e "\033[0;32m\nGip installed successfully!\033[0m"

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "\033[0;33mWarning: $INSTALL_DIR is not in your PATH.\033[0m"
    echo "Add the following to your shell config (.bashrc, .zshrc, etc.):"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR\""
fi

echo "Run 'gip --version' to verify."
