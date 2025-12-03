# Gip C++ Coding Standards

## Naming Conventions

### General Rules

| Element | Convention | Example |
|---------|------------|---------|
| Files | `snake_case.cpp`, `snake_case.h` | `git_adapter.cpp` |
| Classes/Structs | `PascalCase` | `GitAdapter`, `ManifestEntry` |
| Functions/Methods | `camelCase` | `getFileHistory()`, `parseManifest()` |
| Variables | `camelCase` | `commitSha`, `fileHistory` |
| Constants | `kPascalCase` | `kDefaultTimeout`, `kMaxRetries` |
| Macros | `UPPER_SNAKE_CASE` | `GIP_VERSION`, `GIP_USE_LIBGIT2` |
| Namespaces | `lowercase` | `gip`, `gip::commands` |
| Enums | `PascalCase` with `k` prefix values | `enum class LogLevel { kDebug, kError }` |
| Private members | `camelCase_` (trailing underscore) | `repoPath_`, `isInitialized_` |

### Prefixes and Suffixes

| Prefix/Suffix | Usage | Example |
|---------------|-------|---------|
| `k` | Constants and enum values | `kMaxBufferSize` |
| `g_` | Global variables (avoid!) | `g_instance` |
| `_` (suffix) | Private member variables | `buffer_` |
| `p` | Pointer parameters (optional) | `pBuffer` |
| `Is/Has/Can` | Boolean methods | `isValid()`, `hasManifest()` |
| `Get/Set` | Accessors/mutators | `getPath()`, `setEnabled()` |

## Code Style

### Headers

```cpp
/// @file filename.h
/// @brief Brief description
/// @author Author Name
/// @copyright MIT License

#pragma once

#include <standard_headers>

#include "project/headers.h"

namespace gip {
// ...
}  // namespace gip
```

### Include Order

1. Corresponding header (for .cpp files)
2. C system headers
3. C++ standard library headers
4. Third-party library headers
5. Project headers

```cpp
#include "git_adapter.h"  // Corresponding header

#include <cstdio>         // C system headers

#include <string>         // C++ standard headers
#include <vector>

#include <libgit2/git2.h> // Third-party

#include "gip/types.h"    // Project headers
```

### Classes

```cpp
/// @brief Brief description of the class
/// 
/// Detailed description if needed.
class GitAdapter {
public:
    // Types and aliases first
    using PathList = std::vector<std::string>;
    
    // Static constants
    static constexpr int kDefaultTimeout = 30;
    
    // Constructors and destructor
    GitAdapter();
    explicit GitAdapter(std::string_view repoPath);
    ~GitAdapter();
    
    // Deleted/defaulted special members
    GitAdapter(const GitAdapter&) = delete;
    GitAdapter& operator=(const GitAdapter&) = delete;
    GitAdapter(GitAdapter&&) noexcept = default;
    GitAdapter& operator=(GitAdapter&&) noexcept = default;
    
    // Public methods (alphabetical or logical grouping)
    [[nodiscard]] bool isRepository() const;
    [[nodiscard]] std::string getRepositoryRoot() const;
    
private:
    // Private methods
    void initialize();
    
    // Private members (with trailing underscore)
    std::string repoPath_;
    bool isInitialized_ = false;
};
```

### Functions

```cpp
/// @brief Calculate the total price with tax
/// @param basePrice The base price (must be >= 0)
/// @param taxRate The tax rate as a decimal (e.g., 0.08 for 8%)
/// @return The total price including tax
/// @throws std::invalid_argument if basePrice < 0 or taxRate < 0
[[nodiscard]] double calculateTotalPrice(double basePrice, double taxRate);
```

### Error Handling

```cpp
// Prefer std::optional for values that may not exist
[[nodiscard]] std::optional<Manifest> parseManifest(std::string_view input);

// Use Result types for operations that can fail with details
struct Result {
    bool success;
    std::string error;
    
    [[nodiscard]] explicit operator bool() const noexcept { return success; }
};

// Exceptions only for truly exceptional cases
void criticalOperation() {
    if (catastrophicFailure) {
        throw std::runtime_error("Critical system failure");
    }
}
```

### Modern C++ Features

```cpp
// Use auto for complex types, be explicit for simple ones
auto it = container.find(key);  // OK
int count = 0;                   // Prefer explicit for primitives

// Use range-based for loops
for (const auto& entry : manifest.entries) {
    // ...
}

// Use structured bindings
for (const auto& [key, value] : map) {
    // ...
}

// Use [[nodiscard]] for functions where ignoring return is a bug
[[nodiscard]] bool isValid() const;

// Use constexpr where possible
constexpr int kMaxEntries = 100;

// Use string_view for read-only string parameters
void process(std::string_view input);

// Use smart pointers
std::unique_ptr<Resource> resource = std::make_unique<Resource>();
```

### Comments

```cpp
// Single-line comments for brief explanations
int count = 0;  // Track number of processed items

/*
 * Multi-line comments for longer explanations
 * that span multiple lines.
 */

/// @brief Doxygen comments for API documentation
/// @param input The input to process
/// @return The processed result
std::string process(std::string_view input);

// TODO(username): Description of what needs to be done
// FIXME(username): Description of the bug to fix
// HACK(username): Description of temporary workaround
```

### Formatting

- **Indentation**: 4 spaces (no tabs)
- **Line length**: 100 characters max
- **Braces**: Same line for control structures, new line for functions (configurable)
- **Spacing**: Space after keywords, around operators

```cpp
// Control structures
if (condition) {
    doSomething();
} else {
    doSomethingElse();
}

for (int i = 0; i < count; ++i) {
    process(i);
}

// Switch statements
switch (value) {
    case kFirst:
        handleFirst();
        break;
    case kSecond:
        handleSecond();
        break;
    default:
        handleDefault();
        break;
}
```

## Best Practices

### Do

- ✅ Use `const` wherever possible
- ✅ Use `[[nodiscard]]` for functions with important return values
- ✅ Use `noexcept` where appropriate
- ✅ Prefer composition over inheritance
- ✅ Use RAII for resource management
- ✅ Write unit tests for all public APIs

### Don't

- ❌ Use raw `new`/`delete` (use smart pointers)
- ❌ Use C-style casts (use `static_cast`, `dynamic_cast`, etc.)
- ❌ Use `using namespace` in headers
- ❌ Use global mutable state
- ❌ Catch exceptions by value (catch by const reference)
- ❌ Leave magic numbers in code (use named constants)

## Static Analysis

Run before committing:

```bash
# Format code
find src include -name '*.cpp' -o -name '*.h' | xargs clang-format -i

# Static analysis
clang-tidy src/*.cpp -- -I include -std=c++17

# Check for issues
cppcheck --enable=all --suppress=missingIncludeSystem -I include src/
```
