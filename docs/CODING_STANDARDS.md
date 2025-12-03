# Gip Coding Standards

## Naming Conventions

### General Rules

| Element | Convention | Example |
|---------|------------|---------|
| Files | `snake_case.rs` | `git_adapter.rs` |
| Structs/Enums | `PascalCase` | `GitAdapter`, `ManifestEntry` |
| Functions/Methods | `snake_case` | `get_file_history()`, `parse_manifest()` |
| Variables | `snake_case` | `commit_sha`, `file_history` |
| Constants | `UPPER_SNAKE_CASE` | `DEFAULT_TIMEOUT`, `MAX_RETRIES` |
| Modules | `snake_case` | `commands`, `git` |
| Traits | `PascalCase` | `Display`, `Debug` |

## Code Style

### Modules

```rust
//! Brief description of the module
//!
//! Detailed description if needed.

use std::path::Path;

use crate::error::Result;

// ...
```

### Structs

```rust
/// Brief description of the struct
/// 
/// Detailed description if needed.
#[derive(Debug, Clone)]
pub struct GitAdapter {
    repo_path: String,
    is_initialized: bool,
}

impl GitAdapter {
    /// Creates a new GitAdapter
    pub fn new(repo_path: &str) -> Self {
        Self {
            repo_path: repo_path.to_string(),
            is_initialized: false,
        }
    }
}
```

### Functions

```rust
/// Calculate the total price with tax
///
/// # Arguments
///
/// * `base_price` - The base price (must be >= 0)
/// * `tax_rate` - The tax rate as a decimal (e.g., 0.08 for 8%)
///
/// # Returns
///
/// The total price including tax
pub fn calculate_total_price(base_price: f64, tax_rate: f64) -> Result<f64> {
    if base_price < 0.0 {
        return Err(Error::InvalidInput("Base price cannot be negative".into()));
    }
    Ok(base_price * (1.0 + tax_rate))
}
```

### Error Handling

```rust
// Use Result types for operations that can fail
pub fn parse_manifest(input: &str) -> Result<Manifest> {
    // ...
}

// Use custom Error types with thiserror
#[derive(Error, Debug)]
pub enum Error {
    #[error("IO error: {0}")]
    Io(#[from] std::io::Error),
    #[error("Invalid manifest: {0}")]
    InvalidManifest(String),
}
```

### Modern Features

```rust
// Use pattern matching
match value {
    Some(v) => process(v),
    None => handle_none(),
}

// Use if let
if let Some(v) = value {
    process(v);
}

// Use iterators
let results: Vec<_> = items.iter()
    .map(|x| process(x))
    .collect();
```

## Best Practices

### Do

- ✅ Use `clippy` to catch common mistakes
- ✅ Use `rustfmt` for consistent formatting
- ✅ Write unit tests for all public APIs
- ✅ Use `Result` for error handling instead of panics
- ✅ Document public APIs with doc comments (`///`)

### Don't

- ❌ Use `unwrap()` in production code (use `expect()` or `?`)
- ❌ Use `unsafe` unless absolutely necessary
- ❌ Ignore compiler warnings

## Static Analysis

Run before committing:

```bash
# Format code
cargo fmt

# Static analysis
cargo clippy -- -D warnings

# Run tests
cargo test
```
