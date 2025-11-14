# Gip Makefile

.PHONY: help test test-unit test-integration test-all coverage build clean install

# Default target
help:
	@echo "Gip - Git++ Makefile"
	@echo ""
	@echo "Available targets:"
	@echo "  make test             - Run all unit tests"
	@echo "  make test-unit        - Run unit tests only (short mode)"
	@echo "  make test-integration - Run integration tests (requires test setup)"
	@echo "  make test-all         - Run all tests with verbose output"
	@echo "  make coverage         - Generate coverage report"
	@echo "  make build            - Build gip binary"
	@echo "  make install          - Install gip to GOPATH/bin"
	@echo "  make clean            - Clean build artifacts and test directories"

# Run all unit tests
test:
	@echo "Running unit tests..."
	go test ./internal/...

# Run unit tests with short flag (skip slow tests)
test-unit:
	@echo "Running unit tests (short mode)..."
	go test -short ./internal/...

# Run integration tests (if they exist in tests/integration/)
test-integration:
	@echo "Running integration tests..."
	@if [ -d "tests/integration" ]; then \
		go test -v ./tests/integration/...; \
	else \
		echo "No integration tests found in tests/integration/"; \
	fi

# Run all tests with verbose output
test-all:
	@echo "Running all tests (verbose)..."
	go test -v ./...

# Generate coverage report
coverage:
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"
	@echo ""
	@echo "Coverage summary:"
	go tool cover -func=coverage.out

# Build the binary
build:
	@echo "Building gip..."
	go build -o gip.exe ./cmd/gip
	@echo "Build complete: gip.exe"

# Install to GOPATH/bin
install:
	@echo "Installing gip..."
	go install ./cmd/gip

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@if exist gip.exe del /F /Q gip.exe 2>nul || echo ""
	@if exist gip-merge-driver.exe del /F /Q gip-merge-driver.exe 2>nul || echo ""
	@if exist coverage.out del /F /Q coverage.out 2>nul || echo ""
	@if exist coverage.html del /F /Q coverage.html 2>nul || echo ""
	@if exist test-merge-repo rmdir /S /Q test-merge-repo 2>nul || echo ""
	@if exist test-repo rmdir /S /Q test-repo 2>nul || echo ""
	@echo "Clean complete"

# Run tests and build
all: test build

# Quick check - run tests and build if they pass
check: test
	@echo ""
	@echo "Tests passed! Building..."
	@$(MAKE) build
