#!/bin/bash
# Install Gip locally for development

set -e

echo "Building Gip..."
cargo build --release

echo "Installing Gip..."
cargo install --path .

echo ""
echo "✅ Gip installed successfully!"
echo "Binary location: ~/.cargo/bin/gip"
echo "Run 'gip --version' to verify"
