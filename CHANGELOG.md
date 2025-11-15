# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2025-11-15

### 🚀 Major Features

#### AI-Powered Manifest Generation
- **`gip commit --ai --intent "description"`** - Automatically generate manifests using OpenAI GPT-4o-mini
- Structured prompts with schema and few-shot examples
- Retry logic with error feedback (up to 3 attempts)
- Schema validation for LLM responses
- Markdown code block stripping for robust parsing

#### Global Intent for Multi-Function Commits
- Detect when multiple functions are modified in a single commit
- Prompt for commit-level global intent (shared rationale)
- Per-function inheritance control (inherit or provide unique intent)
- Reduces repetition in manifests for large refactorings

#### Selective Context Injection
- **20x size reduction** in conflict files (236 lines vs 1000+)
- 3-tier matching: file+symbol+hunk → file+symbol → file
- Bottom-up symbol extraction (finds closest function definition)
- Performance: ~150ms for 100-entry manifests

#### Extended Schema (v2.0)
- **Schema Versioning**: `schemaVersion: "2.0"` with automatic v1→v2 migration
- **Global Intent**: Commit-level `behaviorClass` and `rationale`
- **Signature Delta**: Before/after signatures for breaking change tracking
- **Breaking Change Detection**: Automatic detection of:
  - Parameter additions/removals
  - Return type changes
  - Required vs optional parameters
- **Feature Flag Tracking**: Scans code for `FLAG_*`, `ENABLE_*`, `FEATURE_*` patterns
- **Compatibility Markers**: `breaking`, `deprecations`, `migrations` fields
- **Change Type**: `add`, `modify`, `delete`, `rename`
- **Hunk IDs**: Precise line-level tracking (e.g., `H#42`)

### Added

- `internal/llm/` package for LLM integration
  - `generator.go` - Main LLM client interface
  - `openai.go` - OpenAI API client
  - `prompts.go` - Prompt engineering with templates and examples
  - Mock client for testing
- `internal/prompt/ai.go` - AI commit workflow
- `internal/prompt/global_intent_test.go` - Global intent unit tests
- Breaking change detection in `internal/diff/analyzer.go`
- Signature extraction and comparison
- Feature flag detection with regex patterns
- Test helpers: `SaveManifestDirect()`, `LoadManifestDirect()`
- Integration tests:
  - `selective_injection_test.go` - 3 scenarios (single, multi, v1 fallback)
  - `performance_test.go` - Large manifest tests (100 entries)
  - `global_intent_test.go` - Inheritance patterns
- `internal/diff/breaking_test.go` - Breaking change detection tests

### Changed

- **Manifest Schema**: Upgraded from v1.0 to v2.0 (backward compatible)
- `InteractiveCommit()` now supports multi-function detection
- `ToManifestEntry()` populates `SignatureDelta` and `Compatibility` fields
- TOON serialization includes all v2.0 fields
- Merge driver uses selective injection by default
- CLI flags:
  - Added `--ai` and `--intent` flags to commit command
  - Enhanced help text with v2.0 examples

### Fixed

- Unused import in `prompts.go`
- Test isolation issues with parallel execution
- Manifest loading in test environments (no full git validation)

### Backward Compatibility

- ✅ Reads v1.0 manifests seamlessly
- ✅ Automatic migration: adds default values (HunkID="H#0", ChangeType="modify")
- ✅ No breaking changes to existing workflows
- ✅ v1.0 commands still work (`gip commit -c`)

### Performance

- Selective injection: 20x reduction in conflict file size
- LLM generation: < 5 seconds for typical commits
- Manifest enrichment: ~150ms for 100-entry manifests

### Documentation

- Updated ARCHITECTURE_V2.md with all phase completions
- Enhanced README with v2.0 features and examples
- Added migration guide in CHANGELOG
- Inline code documentation for new packages

## [Unreleased]

### Changed - Rust Rewrite
- **Complete rewrite from Go to Rust** for memory safety and performance
- Binary size reduced from 6.6MB to 3.6MB (45% smaller)
- Comprehensive test suite with 22+ unit tests using TDD approach
- Module-level documentation for all public APIs
- Linter configuration with cargo clippy
- Memory-safe implementation with zero unsafe code
- Complete migration of all functionality to Rust
- Updated all documentation to remove Go references
- Build system now uses Cargo instead of Go modules
- CI/CD workflow updated for Rust toolchain

### Added
- Rust-based implementation with modern error handling (anyhow, thiserror)
- Type-safe serialization with serde
- Test coverage tracking with cargo-tarpaulin support
- Complete contribution guidelines updated for Rust (CONTRIBUTING.md)
- Code of Conduct (CODE_OF_CONDUCT.md)
- Security policy (SECURITY.md)
- MIT License file

### Fixed
- SHA slicing panic when SHA length < 8 characters
- Git repository initialization defaults to `main` branch
- Proper error handling when no commits exist

## [0.1.0] - Initial Release

### Added
- Core Gip functionality for manifest creation
- Interactive prompts for collecting change metadata
- Custom Git merge driver with context enrichment
- TOON format serialization for conflict markers
- Multi-language diff analysis (Python, Go, JS, TS, Java, C/C++, Ruby, Rust, PHP)
- CLI commands:
  - `gip init` - Initialize Gip in repository ✅
  - `gip commit -c` - Interactive commit with manifest ✅
  - `gip version` - Display version ✅
  - `gip status` - Show Gip status (planned)
  - `gip commit -t` - Generate manifest template (planned)
  - `gip commit -m <file>` - Use pre-written manifest (planned)
- Manifest storage in `.gip/manifest/<sha>.json`
- Support for:
  - Change type detection (add/modify/delete/refactor)
  - Contract specification (preconditions, postconditions, error model)
  - Behavior classification
  - Compatibility flags (breaking changes)
  - Side effects tracking
  - Related changes linking
- GitHub Actions CI/CD workflow
  - Multi-OS testing (Ubuntu, Windows, macOS)
  - Multi-Rust-version testing (stable, beta)
  - Coverage reporting with cargo-tarpaulin
  - Build artifact generation
- Makefile with common development tasks
- Comprehensive README with examples

### Documentation
- Installation instructions for multiple platforms
- Usage examples and workflow guide
- Architecture overview
- TOON format explanation
- Troubleshooting guide
- Roadmap and future features

[Unreleased]: https://github.com/iamHrithikRaj/gip/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/iamHrithikRaj/gip/releases/tag/v0.1.0
