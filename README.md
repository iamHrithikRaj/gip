# Gip (Git++) 🚀

[![CI](https://github.com/iamHrithikRaj/gip/actions/workflows/test.yml/badge.svg)](https://github.com/iamHrithikRaj/gip/actions/workflows/test.yml)
[![Go Version](https://img.shields.io/github/go-mod/go-version/iamHrithikRaj/gip)](https://go.dev/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/iamHrithikRaj/gip)](https://goreportcard.com/report/github.com/iamHrithikRaj/gip)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**A lightweight Git wrapper that enriches merge conflicts with structured context for humans and LLMs.**

Gip transforms cryptic merge conflicts into self-documenting resolution tasks by injecting structured manifests directly into conflict markers—no repo scanning required.

## ✨ v2.0 Features

- 🤖 **AI-Generated Manifests** - Use OpenAI to auto-generate contract documentation from diffs
- 🎯 **Global Intent** - Document commit-level intent for multi-function changes
- 📦 **Selective Injection** - Only relevant context appears in conflicts (20x size reduction)
- 🔍 **Breaking Change Detection** - Automatic detection of signature changes
- 🚩 **Feature Flag Tracking** - Scan and document feature flags in code
- 🔄 **Backward Compatible** - Seamlessly reads v1.0 manifests

---

## The Problem

Traditional Git conflict markers provide **zero semantic context**:

```python
<<<<<<< HEAD
total += item.price * 1.08
=======
total += item.price + 5.99
>>>>>>> feature-branch
```

**Questions left unanswered:**
- What was the intent behind each change?
- What are the preconditions and postconditions?
- Are there side effects or error handling differences?
- Which change is a feature vs a bugfix?

Both humans and LLMs waste time scanning the entire codebase to understand conflicts.

---

## The Solution

Gip adds **TOON manifests** to commits, capturing structured context:

```python
<<<<<<< HEAD
total += item.price * 1.08
||| Gip CONTEXT (HEAD - Your changes)
||| Commit: 3c7d3422
||| behaviorClass[1]: feature
||| preconditions[1]: items is list with .price
||| postconditions[1]: returns total with 8% sales tax
||| errorModel[1]: AttributeError if no .price
||| rationale: Added 8% sales tax to comply with state law
||| sideEffects[0]:
||| symbol: calculate_total
=======
total += item.price + 5.99
||| Gip CONTEXT (MERGE_HEAD - Their changes)
||| Commit: 1825fc7c
||| behaviorClass[1]: feature
||| preconditions[1]: items is list with .price
||| postconditions[1]: returns total plus flat shipping
||| errorModel[1]: AttributeError if no .price
||| rationale: Added $5.99 flat shipping fee for all orders
||| sideEffects[0]:
||| symbol: calculate_total
>>>>>>> feature-branch
```

**Now you know:**
- ✅ One adds tax (8%), the other adds shipping ($5.99)
- ✅ Both are features, not bugfixes
- ✅ Both have same preconditions and error handling
- ✅ Resolution is clear: keep both or choose based on business logic

---

## Features

- 🎯 **Zero Scanning**: Context injected directly into conflicts—no repo traversal needed
- 📦 **Single Binary**: No runtime dependencies, just copy and run (6.6MB self-contained)
- 🔄 **Git-Compatible**: Drop-in replacement—use `gip` instead of `git`
- 🤖 **LLM-Friendly**: TOON format optimized for LLM token efficiency (30-60% smaller than JSON)
- 🔗 **Non-Invasive**: Works alongside standard Git workflow
- 📝 **Structured Context**: Contracts, behavior classes, side effects, compatibility flags

---

## Installation

### Option 1: Download Binary (Recommended)

**No dependencies required!** Download the latest release for your platform:

**Windows (PowerShell)**:
```powershell
# Download from GitHub releases
Invoke-WebRequest -Uri "https://github.com/iamHrithikRaj/gip/releases/latest/download/gip-windows-amd64.exe" -OutFile "gip.exe"

# Move to a directory in your PATH (optional)
Move-Item gip.exe $env:LOCALAPPDATA\Programs\gip.exe

# Verify
gip --version
```

**macOS/Linux**:
```bash
# Download
curl -LO https://github.com/iamHrithikRaj/gip/releases/latest/download/gip-$(uname -s)-$(uname -m)

# Make executable
chmod +x gip-*

# Move to PATH
sudo mv gip-* /usr/local/bin/gip

# Verify
gip --version
```

**Single binary**: The executable is **6.6MB** and has **zero runtime dependencies**. No Go, no libraries—just copy and run!

### Option 2: Build from Source

**Prerequisites**: Go 1.21+

```bash
git clone https://github.com/iamHrithikRaj/gip
cd gip
go build -o gip ./cmd/gip
```

**Dependencies** (automatically fetched by Go):
- `github.com/spf13/cobra` - CLI framework
- `github.com/AlecAivazis/survey/v2` - Interactive prompts
- `github.com/alpkeskin/gotoon` - TOON encoding
- `github.com/fatih/color` - Terminal colors

---

## Quick Start

### 1. Initialize Gip

```bash
cd your-git-repo
gip init
```

This installs Git hooks and configures the repository. You only need to do this once per repo.

### 2. Commit with AI (v2.0 - Recommended)

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=sk-...

# Stage your changes
git add modified_file.py

# Let AI generate the manifest
gip commit --ai --intent "Add strict validation to user parser"
```

The AI will:
- Analyze your git diff
- Extract function signatures and changes
- Generate structured contracts (preconditions, postconditions, error models)
- Detect breaking changes automatically
- Identify feature flags in your code

### 2. Alternative: Interactive Mode

For manual control or when API keys aren't available:

```bash
# Stage your changes
git add modified_file.py

# Interactive commit with prompts
gip commit -c
```

**v2.0 Features in Interactive Mode:**
- **Multi-function detection**: Prompts for global intent if multiple functions changed
- **Inheritance**: Each function can inherit global intent or provide unique rationale
- **Breaking change hints**: Warns if signature changes detected

You'll be prompted to fill in:
- **Global Intent** (if multi-function): Shared rationale across all changes
- **Behavior Class**: Feature, bugfix, refactor, perf, security, etc.
- **Preconditions**: What must be true before the change
- **Postconditions**: What's guaranteed after the change
- **Error Model**: How errors are handled
- **Side Effects**: Logs, file I/O, network calls, etc.
- **Rationale**: Why this change was made

### 3. Merge with Enriched Conflicts

```bash
# Use gip instead of git for merging
gip merge feature-branch

# If conflicts occur, they'll be enriched with TOON context
# Resolve manually, then:
gip add resolved_file.py
gip commit
```

---

## How It Works

### 1. Commit Phase

When you run `gip commit -c`:
1. Gip analyzes staged changes (git diff)
2. Extracts symbols (functions, classes) that changed
3. Prompts for structured context (contracts, behavior, rationale)
4. Stores manifest as `.gip/manifest/<commit-sha>.json`

### 2. Merge Phase

When you run `gip merge <branch>`:
1. Gip calls `git merge` internally
2. If conflicts occur, Gip post-processes conflicted files
3. For each conflict block:
   - Loads manifests for both commits (HEAD and MERGE_HEAD)
   - Converts JSON manifests to TOON format
   - Injects context with `|||` prefix between conflict markers

### 3. Storage Format

Manifests are stored as **JSON** (easy parsing) but displayed as **TOON** (token-efficient):

**Stored** (`.gip/manifest/<sha>.json`):
```json
{
  "commit": "abc123",
  "entries": [{
    "anchor": {"file": "cart.py", "symbol": "calculate_total"},
    "contract": {
      "preconditions": ["items is list"],
      "postconditions": ["returns float"],
      "errorModel": ["AttributeError"]
    },
    "behaviorClass": ["feature"],
    "rationale": "Added tax"
  }]
}
```

**Displayed** (in conflict markers):
```
||| behaviorClass[1]: feature
||| preconditions[1]: items is list
||| postconditions[1]: returns float
||| errorModel[1]: AttributeError
||| rationale: Added tax
||| symbol: calculate_total
```

---

## Manifest Schema

### Entry Structure

| Field | Type | Description | Example |
|-------|------|-------------|---------|
| `anchor.file` | string | File path | `"src/cart.py"` |
| `anchor.symbol` | string | Function/class name | `"calculate_total"` |
| `contract.preconditions` | string[] | Input guarantees | `["items is non-empty"]` |
| `contract.postconditions` | string[] | Output guarantees | `["returns float >= 0"]` |
| `contract.errorModel` | string[] | Error handling | `["ValueError on invalid"]` |
| `behaviorClass` | string[] | Change category | `["feature", "bugfix"]` |
| `sideEffects` | string[] | Side effects | `["logs:info", "writes:cache"]` |
| `rationale` | string | Why this change | `"Added tax calculation"` |

### Behavior Classes

- `feature` - New functionality
- `bugfix` - Bug correction
- `refactor` - Code restructuring
- `perf` - Performance optimization
- `security` - Security improvement
- `validation` - Input validation
- `docs` - Documentation
- `config` - Configuration change

### Side Effects Format

- `logs:<level>` - Logging (info, error, debug)
- `reads:<resource>` - Reads (file, db, api)
- `writes:<resource>` - Writes (file, db, cache)
- `network:<protocol>` - Network calls (http, grpc)
- `none` - Pure function, no side effects

---

## CLI Reference

### Core Commands

```bash
gip init                    # Initialize Gip in current repo
gip commit -c              # Interactive commit with context
gip merge <branch>         # Merge with enriched conflicts
gip status                 # Show Gip manifest status
```

### Git Passthrough

Gip transparently passes through **any Git command** it doesn't recognize:

```bash
gip log --oneline          # Same as: git log --oneline
gip branch -a              # Same as: git branch -a
gip diff                   # Same as: git diff
gip add .                  # Same as: git add .
gip push origin main       # Same as: git push origin main
gip pull                   # Same as: git pull
gip checkout -b new-feat   # Same as: git checkout -b new-feat
```

**You can use `gip` as a complete drop-in replacement for `git`!**

---

## Examples

### Example 1: Feature Conflict

**Scenario**: Two developers add different features to the same function.

**Branch A** (tax calculation):
```python
def calculate_total(items):
    return sum(item.price * 1.08 for item in items)
```

**Branch B** (shipping fee):
```python
def calculate_total(items):
    return sum(item.price for item in items) + 5.99
```

**Gip-enriched conflict**:
```python
<<<<<<< HEAD
def calculate_total(items):
    return sum(item.price * 1.08 for item in items)
||| Gip CONTEXT (HEAD)
||| rationale: Added 8% sales tax per state law
||| behaviorClass[1]: feature
||| postconditions[1]: returns total with 8% tax applied
||| sideEffects[0]:
=======
def calculate_total(items):
    return sum(item.price for item in items) + 5.99
||| Gip CONTEXT (MERGE_HEAD)
||| rationale: Added flat shipping fee for all orders
||| behaviorClass[1]: feature  
||| postconditions[1]: returns total plus $5.99 shipping
||| sideEffects[0]:
>>>>>>> feature-shipping
```

**Resolution**: Combine both changes:
```python
def calculate_total(items):
    subtotal = sum(item.price * 1.08 for item in items)  # Tax
    return subtotal + 5.99  # Shipping
```

### Example 2: Bugfix vs Feature

**Scenario**: One developer fixes a crash, another adds a new feature to the same function.

**Branch A** (bugfix - empty input validation):
```python
def apply_discount(items, discount):
    if not items:
        raise ValueError("items cannot be empty")
    return sum(item.price for item in items) * (1 - discount)
```

**Branch B** (feature - discount parameter with default):
```python
def apply_discount(items, discount=0.1):
    return sum(item.price for item in items) * (1 - discount)
```

**Gip-enriched conflict**:
```python
<<<<<<< HEAD
def apply_discount(items, discount):
    if not items:
        raise ValueError("items cannot be empty")
    return sum(item.price for item in items) * (1 - discount)
||| Gip CONTEXT (HEAD)
||| behaviorClass[1]: bugfix
||| rationale: Fixed crash on empty input
||| errorModel[1]: ValueError on empty list
||| preconditions[1]: items must be non-empty list
||| postconditions[1]: returns discounted total
||| sideEffects[0]:
=======
def apply_discount(items, discount=0.1):
    return sum(item.price for item in items) * (1 - discount)
||| Gip CONTEXT (MERGE_HEAD)
||| behaviorClass[1]: feature
||| rationale: Added discount parameter with 10% default
||| errorModel[0]:
||| preconditions[1]: items is list, discount is float
||| postconditions[1]: returns discounted total
||| sideEffects[0]:
>>>>>>> feature-discount
```

**Resolution**: Combine both - keep feature's default parameter AND bugfix's validation:
```python
def apply_discount(items, discount=0.1):
    if not items:
        raise ValueError("items cannot be empty")
    return sum(item.price for item in items) * (1 - discount)
```

**Why this resolution?**
- ✅ Feature adds convenience (default 10% discount)
- ✅ Bugfix prevents production crash (empty list validation)
- ✅ Both improvements are complementary, not conflicting

---

## Architecture

```
gip/
├── cmd/
│   └── gip/
│       └── main.go          # CLI entry point with cobra
│                            # Passthrough for unknown commands
├── internal/
│   ├── manifest/
│   │   ├── types.go         # Manifest schema (Go structs)
│   │   ├── storage.go       # JSON save/load operations
│   │   └── toon.go          # TOON serialization (gotoon.Encode)
│   ├── diff/
│   │   └── analyzer.go      # Git diff parsing, symbol extraction
│   ├── merge/
│   │   └── driver.go        # Conflict detection & enrichment
│   ├── git/
│   │   └── integration.go   # Git operations (hooks, config)
│   └── prompt/
│       └── interactive.go   # Interactive CLI prompts (survey)
└── .gip/
    └── manifest/
        └── <commit-sha>.json # JSON manifest storage
```

**Design Decisions**:
- **JSON storage**: Easy parsing with Go's native `encoding/json`
- **TOON display**: Token-efficient format for conflict markers (gotoon encoder)
- **No Git modification**: Manifests stored separately in `.gip/` directory
- **Post-processing**: Merge driver runs after Git merge, enriches existing conflicts
- **Transparent wrapper**: Unknown commands passed to `git` via `exec.Command`

---

## Comparison with Alternatives

| Feature | Gip | Plain Git | Git Notes | GitHub PRs |
|---------|-----|-----------|-----------|------------|
| Context in conflicts | ✅ Inline | ❌ None | 🟡 Separate | 🟡 Separate |
| Token efficiency | ✅ TOON | N/A | ❌ JSON | ❌ Markdown |
| Offline support | ✅ Yes | ✅ Yes | ✅ Yes | ❌ No |
| Structured contracts | ✅ Yes | ❌ No | 🟡 Freeform | 🟡 Freeform |
| LLM-friendly | ✅ Yes | ❌ No | 🟡 Partial | 🟡 Partial |
| Setup complexity | 🟡 One command | ✅ None | 🟡 Manual | ❌ Platform lock-in |
| Binary size | ✅ 6.6MB | ✅ ~2MB | ✅ ~2MB | ❌ Cloud |
| Dependencies | ✅ Zero | ✅ Zero | ✅ Zero | ❌ Network |

---

## Design Principles

- **Objective, not subjective**: No "priority" fields, only verifiable contracts
- **Local resolution**: All context embedded in conflict markers—no external scanning
- **Language-agnostic**: Works with any text-based file format
- **Lightweight**: Single binary, minimal overhead
- **Deterministic**: Same conflicts always produce same enriched output
- **Non-invasive**: Standard Git workflow remains unchanged

---

## FAQ

**Q: Does Gip modify Git's internal storage?**  
A: No. Gip stores manifests in `.gip/manifest/` as separate JSON files. Your `.git/` directory remains untouched.

**Q: Can I use Gip with existing Git repos?**  
A: Yes! Run `gip init` in any Git repo to start using Gip. Existing commits without manifests will still merge normally.

**Q: What if I forget to use `gip commit -c`?**  
A: You can still use regular `git commit`. The pre-commit hook will suggest using `gip commit -c`, but won't block the commit.

**Q: Do collaborators need Gip installed?**  
A: No for basic functionality—manifests are just JSON files in `.gip/`. Yes to get enriched conflicts when merging.

**Q: Does Gip work with GUI Git clients?**  
A: Partially. You can use GUI clients for staging/committing, but enriched conflicts only appear when using `gip merge`.

**Q: How much overhead does Gip add?**  
A: Minimal. Manifest files are ~1KB each (JSON format). The binary is 6.6MB. No runtime dependencies.

**Q: Can I customize the manifest schema?**  
A: Yes, edit `internal/manifest/types.go` and rebuild with `go build ./cmd/gip`.

**Q: Does Gip support rebasing?**  
A: Currently only `gip merge` enriches conflicts. Rebase support is planned.

**Q: What about binary files?**  
A: Gip only processes text files. Binary conflicts are left unchanged.

---

## Contributing

Contributions welcome! This is an experimental project exploring structured context for version control.

### Development Setup

```bash
# Clone
git clone https://github.com/iamHrithikRaj/gip
cd gip

# Install dependencies
go mod download

# Build
go build -o gip ./cmd/gip

# Run tests
go test ./...
```

### Testing

Gip has comprehensive unit tests for all core packages:

```bash
# Run all unit tests
make test

# Run tests in short mode (faster)
make test-unit

# Run tests with coverage
make coverage

# Run all tests with verbose output
make test-all
```

**Test Coverage:**
- ✅ `internal/manifest` - 12 tests (storage, TOON serialization)
- ✅ `internal/merge` - 6 tests (conflict detection, enrichment)
- ✅ `internal/diff` - 18 tests (symbol extraction, change detection)

**CI/CD:**
- Tests run automatically on push/PR via GitHub Actions
- Multi-OS testing (Windows, Linux, macOS)
- Multi-Go-version testing (1.21, 1.22, 1.23)

See [TEST_STRUCTURE_PROPOSAL.md](TEST_STRUCTURE_PROPOSAL.md) for full test architecture.

---

## Current Status: v0.1.0 (Alpha)

Gip is in active development. Here's what works today:

### ✅ Working Features
- **`gip init`** - Initialize Gip in a repository
- **`gip commit -c`** - Interactive commit with manifest creation
- **`gip version`** - Display version information
- **Custom merge driver** - Enriches conflicts with TOON context
- **Multi-language support** - Python, Go, JS, TS, Java, C/C++, Ruby, Rust, PHP
- **Test suite** - 40+ tests with CI/CD

### 🚧 In Development
- **`gip status`** - View manifest information
- **`gip commit -t`** - Generate manifest templates
- **`gip commit -m`** - Use pre-written manifests
- **Rebase support** - Enriched conflicts during rebase

## Roadmap

Future planned features:

- [ ] **Gip MCP** - Model Context Protocol server for seamless LLM integration
- [ ] VS Code extension for inline manifest editing
- [ ] GitHub Actions for automated manifest generation
- [ ] LLM-assisted conflict resolution suggestions
- [ ] Manifest diff visualization tool
- [ ] Team manifest templates
- [ ] Semantic merge strategies based on contracts
- [ ] Integration with popular merge tools (Meld, KDiff3, etc.)

**Want to help?** Check out [open issues](https://github.com/iamHrithikRaj/gip/issues) or [contribute](CONTRIBUTING.md)!

---

## Contributing

We welcome contributions! Gip is an open-source project that benefits from community involvement.

### Quick Start

1. **Fork and clone** the repository
2. **Install Go 1.21+** and dependencies
3. **Make your changes** following our coding standards
4. **Add tests** for new functionality
5. **Submit a pull request**

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:
- Setting up your development environment
- Coding standards and best practices
- Testing requirements
- Commit message conventions
- Pull request process

**Good First Issues**: Look for issues labeled `good first issue` to get started!

### Code of Conduct

This project adheres to the [Contributor Covenant Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code.

---

## License

MIT License - see [LICENSE](LICENSE) for details.

By contributing to Gip, you agree that your contributions will be licensed under the MIT License.

---

## Credits

- **TOON Format**: [github.com/johannschopplich/toon](https://github.com/johannschopplich/toon)
- **GoTOON Library**: [github.com/alpkeskin/gotoon](https://github.com/alpkeskin/gotoon)
- **Inspired by**: Git Notes, semantic merge tools, and the need for better conflict resolution UX

---

## Support

- 📖 **Documentation**: See this README and [CONTRIBUTING.md](CONTRIBUTING.md)
- 🐛 **Bug Reports**: [Open an issue](https://github.com/iamHrithikRaj/gip/issues/new?template=bug_report.md)
- 💡 **Feature Requests**: [Suggest a feature](https://github.com/iamHrithikRaj/gip/issues/new?template=feature_request.md)
- ❓ **Questions**: [Ask a question](https://github.com/iamHrithikRaj/gip/issues/new?template=question.md)

---

**Made with ❤️ by a developer tired of cryptic merge conflicts.**

**Gip = Git++** 🎯 Better merges through structured context.
