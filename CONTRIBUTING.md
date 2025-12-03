# Contributing to Gip

Thank you for your interest in contributing to Gip! This document provides guidelines and instructions for contributing.

## Code of Conduct

Please be respectful and constructive in all interactions.

## Getting Started

### Prerequisites

- C++17 compatible compiler (GCC 9+, Clang 10+, MSVC 2019+)
- CMake 3.15+
- Git

### Building

```bash
git clone https://github.com/iamHrithikRaj/gip.git
cd gip
mkdir build && cd build
cmake .. -DGIP_BUILD_TESTS=ON
cmake --build .
```

### Running Tests

```bash
cd build
ctest --output-on-failure
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

We use clang-format. Run before committing:

```bash
find src include -name '*.cpp' -o -name '*.h' | xargs clang-format -i
```

### Guidelines

- Use `snake_case` for file names
- Use `CamelCase` for types
- Use `camelCase` for functions and variables
- Prefix private members with underscore suffix (`member_`)
- Use `[[nodiscard]]` for functions with return values that shouldn't be ignored
- Document public APIs with Doxygen comments

### Example

```cpp
/// @brief Calculate tax for a given price
/// @param price The base price (must be >= 0)
/// @return Price with tax applied
/// @throws std::invalid_argument if price is negative
[[nodiscard]] double calculateTax(double price);
```

## Testing

### Unit Tests

Located in `tests/unit/`. Use Catch2 framework:

```cpp
TEST_CASE("Description", "[tag]") {
    SECTION("specific case") {
        REQUIRE(result == expected);
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
