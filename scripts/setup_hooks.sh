#!/bin/bash

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$( dirname "$SCRIPT_DIR" )"
HOOKS_DIR="$ROOT_DIR/.git/hooks"

echo "Setting up git hooks..."

mkdir -p "$HOOKS_DIR"

cp "$SCRIPT_DIR/git-hooks/pre-commit" "$HOOKS_DIR/pre-commit"
chmod +x "$HOOKS_DIR/pre-commit"

echo "Successfully installed pre-commit hook to $HOOKS_DIR/pre-commit"
