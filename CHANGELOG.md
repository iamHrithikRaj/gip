# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Comprehensive test suite with 40+ tests (36 unit, 4 integration)
- Test helpers for git operations and manifest creation
- Package-level documentation for all internal packages
- Linter configuration with golangci-lint
- Parallel test execution with `t.Parallel()`
- Complete contribution guidelines (CONTRIBUTING.md)
- Code of Conduct (CODE_OF_CONDUCT.md)
- GitHub issue templates (bug, feature, question)
- Pull request template with checklist
- Security policy (SECURITY.md)
- Best practices documentation (BEST_PRACTICES.md)
- MIT License file

### Changed
- Tests now run in parallel for faster execution
- CI/CD workflow includes linter checks
- Improved error handling in manifest storage
- Enhanced test coverage to >80%

### Fixed
- SHA slicing panic when SHA length < 8 characters
- Git repository initialization defaults to `main` branch
- Proper error handling when no commits exist

## [0.1.0] - Initial Release

### Added
- Core GIP functionality for manifest creation
- Interactive prompts for collecting change metadata
- Custom Git merge driver with context enrichment
- TOON format serialization for conflict markers
- Multi-language diff analysis (Python, Go, JS, TS, Java, C/C++, Ruby, Rust, PHP)
- CLI commands:
  - `gip init` - Initialize GIP in repository ✅
  - `gip commit -c` - Interactive commit with manifest ✅
  - `gip version` - Display version ✅
  - `gip status` - Show GIP status (planned)
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
  - Multi-Go-version testing (1.21, 1.22, 1.23)
  - Coverage reporting
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
