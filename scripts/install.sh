#!/bin/bash
# Install GIP locally for development

set -e

echo "Building GIP..."
go build -o gip ./cmd/gip

echo "Installing GIP..."
go install ./cmd/gip

echo ""
echo "✅ GIP installed successfully!"
echo "Run 'gip --version' to verify"
