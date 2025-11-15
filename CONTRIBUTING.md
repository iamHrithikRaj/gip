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

- **Rust 1.70 or higher** - [Install Rust](https://rustup.rs/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- **Cargo** (comes with Rust) - Package manager and build tool

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
# Dependencies are managed by Cargo automatically
# Just make sure you have Rust installed
rustc --version
cargo --version

# Install development tools (optional)
cargo install cargo-tarpaulin  # For code coverage
cargo install cargo-watch      # For auto-rebuild
rustup component add rustfmt   # For code formatting
rustup component add clippy    # For linting
```

### Build the Project

```bash
# Debug build (faster compilation, slower runtime)
cargo build

# Release build (optimized)
cargo build --release

# The binary will be in target/debug/gip or target/release/gip
```

### Run Tests

```bash
# Run all tests
cargo test

# Run tests with output
cargo test -- --nocapture

# Run specific module tests
cargo test manifest

# Run tests with coverage
cargo tarpaulin --out Html

# Run linter (clippy)
cargo clippy -- -D warnings

# Format code
cargo fmt
```

## How Can I Contribute?

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Clear title and description**
- **Steps to reproduce** the problem
- **Expected behavior** vs **actual behavior**
- **Environment details** (OS, Rust version, Git version)
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

- Write clear, idiomatic Rust code
- Follow the [Coding Standards](#coding-standards)
- Add tests for new functionality (TDD approach preferred)
- Update documentation as needed

### 3. Test Your Changes

```bash
# Run all tests
cargo test

# Run linter
cargo clippy -- -D warnings

# Check formatting
cargo fmt --check

# Check test coverage
cargo tarpaulin --out Html
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

### Rust Style Guide

We follow standard Rust conventions:

- Use `cargo fmt` for formatting (automatic with most editors)
- Follow [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)
- Use `cargo clippy` for additional linting
- Write documentation comments for public APIs

### Specific Guidelines

#### Module Documentation

Every module must have module-level documentation:

```

#### Function Documentation

Public functions must have documentation comments:

```rust
/// Enrich conflicts reads a file with Git conflict markers and enriches them
/// with context from Gip manifests for each commit SHA.
///
/// # Arguments
/// * `file_path` - Path to the file with conflicts
/// * `ancestor_sha` - SHA of the common ancestor
/// * `current_sha` - SHA of the current branch
/// * `other_sha` - SHA of the branch being merged
///
/// # Errors
/// Returns an error if the file cannot be read or manifests cannot be loaded.
pub fn enrich_conflicts(
    file_path: &str,
    ancestor_sha: &str,
    current_sha: &str,
    other_sha: &str,
) -> Result<()> {
```

#### Error Handling

- Use `Result<T>` for functions that can fail
- Use `anyhow::Result` for application errors
- Use `thiserror` for custom error types
- Add context with `.context()` or `.with_context()`
- Avoid `unwrap()` and `expect()` in production code

#### Naming Conventions

- Use **snake_case** for variables and functions: `manifest_path`
- Use **PascalCase** for types: `ManifestStorage`
- Use **descriptive names**: prefer `manifest_path` over `mp`
- Acronyms follow case convention: `HttpServer`, `url_parser`

#### Code Organization

- Keep functions small and focused (< 50 lines when possible)
- Group related functionality in the same module
- Put tests in `#[cfg(test)]` modules in the same file
- Use private modules for internal implementation details

### Linter Compliance

Your code must pass `clippy`:

```bash
cargo clippy -- -D warnings
```

Common clippy lints to follow:
- Avoid unnecessary allocations
- Use idiomatic Rust patterns
- Prefer iterators over manual loops
- Use `?` operator for error propagation
- Avoid redundant clones

## Testing Guidelines

### Test Coverage

- Aim for **>80% test coverage** for new code
- All exported functions should have tests
- Write both positive and negative test cases

### Test Structure

Use Rust's built-in testing framework:

```rust
#[cfg(test)]
mod tests {
    use super::*;
    use pretty_assertions::assert_eq;
    
    #[test]
    fn test_extract_symbol() {
        let test_cases = vec![
            ("def calculate():", "calculate"),
            ("fn process() {", "process"),
            // ... more test cases
        ];
        
        for (input, expected) in test_cases {
            let result = extract_symbol(input);
            assert_eq!(result, expected, "Failed for input: {}", input);
        }
    }
    
    #[test]
    fn test_with_tempdir() {
        use tempfile::TempDir;
        let temp_dir = TempDir::new().unwrap();
        // test code using temp_dir
        // TempDir is automatically cleaned up when dropped
    }
}
```

### Test Helpers

- Use `tempfile` crate for temporary directories
- Use `pretty_assertions` for better assertion output
- Use `assert_cmd` for testing CLI commands
- Use `#[should_panic]` for tests that expect panics

### Integration Tests

- Place integration tests in `tests/`
- Mark them with `#[ignore]` to skip by default:
  ```rust
  #[test]
  #[ignore]  // Run with: cargo test -- --ignored
  fn test_integration_scenario() {
      // Test code
  }
  ```

### Test Commands

```bash
# Run all tests
cargo test

# Run with output
cargo test -- --nocapture

# Run specific module tests
cargo test manifest

# Run with coverage
cargo tarpaulin --out Html

# Run doc tests
cargo test --doc
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

- [ ] All tests pass: `cargo test`
- [ ] Linter passes: `cargo clippy -- -D warnings`
- [ ] Code is formatted: `cargo fmt`
- [ ] New code has tests (TDD approach)
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
├── src/
│   ├── main.rs           # CLI entry point
│   ├── lib.rs            # Library root
│   ├── manifest/         # Manifest storage and serialization
│   │   ├── mod.rs        # Module exports
│   │   ├── types.rs      # Data structures
│   │   ├── storage.rs    # File I/O
│   │   └── toon.rs       # TOON serialization
│   ├── merge/            # Merge driver implementation
│   ├── diff/             # Diff analysis
│   ├── git.rs            # Git integration
│   ├── prompt/           # Interactive prompts
│   └── toon/             # TOON utilities
├── tests/                # Integration tests
│   └── integration_tests.rs
├── Cargo.toml            # Dependencies and project config
└── scripts/              # Development scripts
```

### Running Gip Locally

```bash
# Build and install locally
cargo install --path .

# Or run directly
cargo run -- --version

# Test in a repository
cd /path/to/test/repo
cargo run -- init
cargo run -- commit
```

### Debugging

```bash
# Enable debug logging (if available)
export RUST_LOG=debug
export RUST_BACKTRACE=1

# Run with debug symbols
cargo build
# Use your preferred debugger (gdb, lldb, VS Code, etc.)

# Watch mode for development
cargo watch -x test -x run
```

### Performance Profiling

```bash
# Build with profiling
cargo build --release

# Use flamegraph for visualization (requires cargo-flamegraph)
cargo flamegraph

# Benchmark tests
cargo bench
```

## Resources

- [Rust Documentation](https://doc.rust-lang.org/)
- [Rust Book](https://doc.rust-lang.org/book/)
- [Rust API Guidelines](https://rust-lang.github.io/api-guidelines/)
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
