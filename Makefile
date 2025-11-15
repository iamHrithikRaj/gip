# Gip Makefile

.PHONY: help test test-unit test-integration test-all coverage build clean install format lint

# Default target
help:
	@echo "Gip - Git++ Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make test             - Run all tests"
	@echo "  make test-unit        - Run library tests only"
	@echo "  make test-integration - Run integration tests"
	@echo "  make test-all         - Run all tests with output"
	@echo "  make coverage         - Generate coverage report"
	@echo "  make build            - Build gip binary (debug)"
	@echo "  make release          - Build optimized release binary"
	@echo "  make install          - Install gip to ~/.cargo/bin"
	@echo "  make clean            - Clean build artifacts"
	@echo "  make format           - Format code with rustfmt"
	@echo "  make lint             - Run clippy linter"

# Run all tests
test:
	@echo "Running tests..."
	cargo test

# Run library tests only
test-unit:
	@echo "Running library tests..."
	cargo test --lib

# Run integration tests
test-integration:
	@echo "Running integration tests..."
	cargo test --test '*'

# Run all tests with output
test-all:
	@echo "Running all tests (verbose)..."
	cargo test -- --nocapture

# Generate coverage report (requires cargo-tarpaulin)
coverage:
	@echo "Generating coverage report..."
	cargo tarpaulin --out Html --output-dir target/coverage
	@echo "Coverage report generated: target/coverage/index.html"
	go tool cover -func=coverage.out

# Build debug binary
build:
	@echo "Building gip (debug)..."
	cargo build
	@echo "Build complete: target/debug/gip"

# Build release binary
release:
	@echo "Building gip (release)..."
	cargo build --release
	@echo "Build complete: target/release/gip"

# Install to ~/.cargo/bin
install:
	@echo "Installing gip..."
	cargo install --path .

# Format code
format:
	@echo "Formatting code..."
	cargo fmt

# Run linter
lint:
	@echo "Running clippy..."
	cargo clippy -- -D warnings

# Clean build artifacts
clean:
	@echo "Cleaning..."
	cargo clean
	@echo "Clean complete"

# Run tests and build
all: test build

# Quick check - format, lint, test, then build
check: format lint test
	@echo ""
	@echo "All checks passed! Building..."
	@$(MAKE) build
