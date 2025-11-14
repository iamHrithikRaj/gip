# Contributing to Gip

First off, thank you for considering contributing to Gip! It's people like you that make Gip such a great tool.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [How Can I Contribute?](#how-can-i-contribute)
- [Development Workflow](#development-workflow)
- [Coding Standards](#coding-standards)
- [Testing Guidelines](#testing-guidelines)
- [Commit Message Guidelines](#commit-message-guidelines)
- [Pull Request Process](#pull-request-process)

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- **Go 1.21 or higher** - [Install Go](https://golang.org/doc/install)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Make** (optional but recommended) - For running build tasks

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/gip.git
   cd gip
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/iamHrithikRaj/gip.git
   ```

## Development Setup

### Install Dependencies

```bash
# Download Go module dependencies
go mod download

# Install development tools
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
```

### Build the Project

```bash
# Using Make
make build

# Or directly with Go
go build -o gip ./cmd/gip
```

### Run Tests

```bash
# Run all tests
make test

# Run only unit tests
make test-unit

# Run with coverage
make coverage

# Run linter
golangci-lint run ./...
```

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the problem
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, Go version, Git version)
- **Code samples** or **test cases** if applicable

Use the bug report template when creating issues.

### Suggesting Enhancements

Enhancement suggestions are tracked as GitHub issues. When creating an enhancement suggestion:

- **Use a clear and descriptive title**
- **Provide detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **Include examples** of how it would work

Use the feature request template when creating enhancement issues.

### Your First Code Contribution

Unsure where to begin? Look for issues labeled:

- `good first issue` - Good for newcomers
- `help wanted` - Need assistance from the community
- `documentation` - Documentation improvements

### Pull Requests

We actively welcome your pull requests! See the [Pull Request Process](#pull-request-process) below.

## Development Workflow

### 1. Create a Branch

```bash
# Update your fork
git checkout main
git pull upstream main

# Create a feature branch
git checkout -b feature/your-feature-name
# Or for bug fixes
git checkout -b fix/bug-description
```

### 2. Make Your Changes

- Write clear, idiomatic Go code
- Follow the [Coding Standards](#coding-standards)
- Add tests for new functionality
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
make test

# Run linter
golangci-lint run ./...

# Check test coverage
make coverage
```

### 4. Commit Your Changes

Follow the [Commit Message Guidelines](#commit-message-guidelines):

```bash
git add .
git commit -m "feat: add new feature description"
```

### 5. Push and Create PR

```bash
git push origin feature/your-feature-name
```

Then create a pull request on GitHub.

## Coding Standards

### Go Style Guide

We follow standard Go conventions:

- Use `gofmt` for formatting (automatic with most editors)
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)

### Specific Guidelines

#### Package Documentation

Every package must have a package-level comment:

```go
// Package merge implements Gip's custom merge driver that enriches Git conflict
// markers with structured context from Gip manifests.
package merge
```

#### Function Documentation

Exported functions must have godoc comments:

```go
// EnrichConflicts reads a file with Git conflict markers and enriches them
// with context from Gip manifests for each commit SHA.
//
// Returns an error if the file cannot be read or manifests cannot be loaded.
func EnrichConflicts(filePath, ancestorSHA, currentSHA, otherSHA string) error {
```

#### Error Handling

- Always check and handle errors
- Wrap errors with context: `fmt.Errorf("failed to load manifest: %w", err)`
- Use sentinel errors for common cases: `var ErrManifestNotFound = errors.New("manifest not found")`
- Never panic in production code

#### Naming Conventions

- Use **camelCase** for variables: `manifestPath`
- Use **PascalCase** for exported types: `ManifestStorage`
- Use **descriptive names**: prefer `manifestPath` over `mp`
- Acronyms should be uppercase: `HTTPServer`, `URLParser`

#### Code Organization

- Keep functions small and focused (< 50 lines when possible)
- Group related functionality in the same file
- Put tests in `*_test.go` files alongside the code
- Use `internal/` for private packages

### Linter Compliance

Your code must pass `golangci-lint`:

```bash
golangci-lint run ./...
```

Our configuration enables these linters:
- `gofmt`, `govet`, `errcheck`, `staticcheck`
- `gosimple`, `unused`, `ineffassign`
- `misspell`, `gocyclo`, `dupl`, `goconst`

## Testing Guidelines

### Test Coverage

- Aim for **>80% test coverage** for new code
- All exported functions should have tests
- Write both positive and negative test cases

### Test Structure

Use table-driven tests for multiple scenarios:

```go
func TestExtractSymbol(t *testing.T) {
    t.Parallel()
    
    tests := []struct {
        name     string
        input    string
        expected string
    }{
        {
            name:     "Python function",
            input:    "def calculate():",
            expected: "calculate",
        },
        // ... more test cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := ExtractSymbol(tt.input)
            if result != tt.expected {
                t.Errorf("expected %s, got %s", tt.expected, result)
            }
        })
    }
}
```

### Test Helpers

- Use `t.Helper()` in test helper functions
- Use `t.TempDir()` for temporary directories
- Use `t.Cleanup()` for cleanup operations
- Add `t.Parallel()` to independent tests

### Integration Tests

- Place integration tests in `tests/integration/`
- Use `testing.Short()` to skip in short mode:
  ```go
  if testing.Short() {
      t.Skip("Skipping integration test in short mode")
  }
  ```

### Test Commands

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run only short tests (skip integration)
go test -short ./...

# Run specific package
go test ./internal/manifest/...

# Verbose output
go test -v ./...
```

## Commit Message Guidelines

We follow [Conventional Commits](https://www.conventionalcommits.org/):

### Format

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Types

- `feat`: New feature
- `fix`: Bug fix
- `docs`: Documentation changes
- `style`: Formatting, missing semicolons, etc.
- `refactor`: Code restructuring without functionality change
- `perf`: Performance improvements
- `test`: Adding or updating tests
- `chore`: Maintenance tasks, dependency updates

### Examples

```bash
feat(merge): add support for multi-hunk conflicts

Add ability to enrich conflicts that span multiple hunks
with appropriate manifest context for each section.

Closes #42

---

fix(manifest): handle short SHA values

Previously crashed when SHA was less than 8 characters.
Now checks length before slicing.

Fixes #38

---

docs(readme): add installation instructions

Add detailed installation steps for Linux, macOS, and Windows.

---

test(diff): add edge case tests for symbol extraction

Add tests for empty diffs, malformed code, and unusual
language constructs.
```

### Rules

- Use imperative mood: "add feature" not "added feature"
- Don't capitalize first letter of subject
- No period at the end of subject
- Limit subject line to 50 characters
- Wrap body at 72 characters
- Use body to explain *what* and *why*, not *how*

## Pull Request Process

### Before Submitting

Ensure your PR meets these requirements:

- [ ] All tests pass: `make test`
- [ ] Linter passes: `golangci-lint run ./...`
- [ ] Code is formatted: `gofmt -w .`
- [ ] New code has tests
- [ ] Documentation is updated
- [ ] Commit messages follow guidelines
- [ ] Branch is up to date with `main`

### PR Description

Use the pull request template and include:

- **Summary** of changes
- **Motivation** for the changes
- **Testing** performed
- **Related issues** (use `Fixes #123` or `Closes #123`)

### Review Process

1. **Automated checks** must pass (CI/CD)
2. **Code review** by at least one maintainer
3. **Address feedback** and update PR
4. **Squash commits** if requested
5. **Merge** by maintainer once approved

### After Your PR is Merged

- Delete your branch
- Pull the latest `main`:
  ```bash
  git checkout main
  git pull upstream main
  ```
- Celebrate! 🎉

## Development Tips

### Project Structure

```
gip/
├── cmd/gip/              # Main application entry point
├── internal/             # Private packages
│   ├── manifest/         # Manifest storage and serialization
│   ├── merge/            # Merge driver implementation
│   ├── diff/             # Diff analysis
│   ├── git/              # Git integration
│   ├── prompt/           # Interactive prompts
│   └── toon/             # TOON serialization
├── tests/                # Integration tests
│   ├── fixtures/         # Test data
│   ├── helpers/          # Test utilities
│   └── integration/      # Integration test suites
└── scripts/              # Development scripts
```

### Running Gip Locally

```bash
# Build and install locally
go install ./cmd/gip

# Or run directly
go run ./cmd/gip --version

# Test in a repository
cd /path/to/test/repo
gip init
gip commit
```

### Debugging

```bash
# Enable verbose output (if available)
export GIP_DEBUG=1

# Use Go's built-in tools
go run -race ./cmd/gip  # Race detector
go test -race ./...     # Race detection in tests

# Use delve debugger
dlv debug ./cmd/gip
```

### Performance Profiling

```bash
# CPU profiling
go test -cpuprofile=cpu.prof ./internal/merge/
go tool pprof cpu.prof

# Memory profiling
go test -memprofile=mem.prof ./internal/merge/
go tool pprof mem.prof
```

## Resources

- [Go Documentation](https://golang.org/doc/)
- [Effective Go](https://golang.org/doc/effective_go.html)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Conventional Commits](https://www.conventionalcommits.org/)
- [GitHub Flow](https://guides.github.com/introduction/flow/)

## Questions?

- Open an issue with the `question` label
- Check existing documentation
- Review closed issues and PRs

## License

By contributing, you agree that your contributions will be licensed under the MIT License.

---

**Thank you for contributing to Gip!** Your efforts help make conflict resolution better for everyone. 🚀
