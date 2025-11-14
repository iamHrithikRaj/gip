#!/bin/bash
# Install Gip locally for development

set -e

echo "Building Gip..."
go build -o gip ./cmd/gip

echo "Installing Gip..."
go install ./cmd/gip

echo ""
echo "✅ Gip installed successfully!"
echo "Run 'gip --version' to verify"
