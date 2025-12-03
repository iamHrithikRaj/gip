# Contributing to Gip

Thank you for your interest in contributing to Gip! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- Rust 1.75+ (stable)
- Git

### Building

```bash
git clone https://github.com/iamHrithikRaj/gip.git
cd gip
cargo build
```

### Running Tests

```bash
cargo test
```

## Development Workflow

### Branch Naming

- `feature/description` - New features
- `bugfix/description` - Bug fixes
- `refactor/description` - Code refactoring
- `docs/description` - Documentation updates

### Commit Messages

We eat our own dog food! Use Gip for commits:

```
feat: add new feature

gip:
{
  schemaVersion: "2.0",
  entries: [{
    file: "src/new_feature.cpp",
    behavior: "feature",
    rationale: "Why this feature was added"
  }]
}
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Ensure tests pass
5. Submit a pull request

## Code Style

### Formatting

We use rustfmt. Run before committing:

```bash
cargo fmt
```

### Guidelines

- Use `snake_case` for file names, functions, and variables
- Use `PascalCase` for types and traits
- Use `UPPER_SNAKE_CASE` for constants
- Document public APIs with doc comments (`///`)

### Example

```rust
/// Calculate tax for a given price
///
/// # Arguments
///
/// * `price` - The base price (must be >= 0)
///
/// # Returns
///
/// Price with tax applied, or an error if price is negative
pub fn calculate_tax(price: f64) -> Result<f64, Error> {
    if price < 0.0 {
        return Err(Error::InvalidInput);
    }
    Ok(price * 1.08)
}
```

## Testing

### Unit Tests

Located in `src/` alongside code or in `tests/unit/`. Use standard test framework:

```rust
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_calculate_tax() {
        assert_eq!(calculate_tax(100.0).unwrap(), 108.0);
    }
}
```

### Integration Tests

Located in `tests/integration/`. Test actual git operations.

## Documentation

- Update README.md for user-facing changes
- Update docs/ARCHITECTURE.md for structural changes
- Add Doxygen comments for new APIs

## Questions?

Open an issue for discussion before starting major work.
