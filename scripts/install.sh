#!/bin/bash
# scripts/install.sh
# Installs Gip (C++ Edition) from source

set -e

echo -e "\033[0;36mInstalling Gip (C++ Edition)...\033[0m"
echo "Performing user-space installation (no sudo required)."

# Check prerequisites
if ! command -v cmake &> /dev/null; then
    echo "Error: CMake is required but not found."
    exit 1
fi

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
BUILD_DIR="$REPO_ROOT/build"
INSTALL_DIR="$HOME/.local"

# Clean build dir
if [ -d "$BUILD_DIR" ]; then
    echo "Cleaning build directory..."
    rm -rf "$BUILD_DIR"
fi

# Configure
echo -e "\033[0;33mConfiguring CMake...\033[0m"
cmake -S "$REPO_ROOT" -B "$BUILD_DIR" -DCMAKE_BUILD_TYPE=Release -DCMAKE_INSTALL_PREFIX="$INSTALL_DIR"

# Build
echo -e "\033[0;33mBuilding Release...\033[0m"
cmake --build "$BUILD_DIR" --config Release

# Install
echo -e "\033[0;33mInstalling to $INSTALL_DIR/bin...\033[0m"
cmake --install "$BUILD_DIR"

echo -e "\033[0;32m\nGip installed successfully!\033[0m"

# Check PATH
if [[ ":$PATH:" != *":$INSTALL_DIR/bin:"* ]]; then
    echo -e "\033[0;33mWarning: $INSTALL_DIR/bin is not in your PATH.\033[0m"
    echo "Add the following to your shell config (.bashrc, .zshrc, etc.):"
    echo "  export PATH=\"\$PATH:$INSTALL_DIR/bin\""
fi

echo "Run 'gitp --version' to verify."
