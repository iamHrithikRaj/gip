#!/bin/bash
# scripts/bootstrap.sh
# One-liner installer for Gip
# Usage: curl -sL https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.sh | bash

set -e

REPO_URL="https://github.com/iamHrithikRaj/gip.git"
TEMP_DIR=$(mktemp -d)

echo "Downloading Gip source..."
git clone --depth 1 "$REPO_URL" "$TEMP_DIR/gip"

echo "Running installer..."
bash "$TEMP_DIR/gip/scripts/install.sh"

echo "Cleaning up..."
rm -rf "$TEMP_DIR"
