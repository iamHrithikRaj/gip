# Gip (Git++) 🚀

[![CI](https://github.com/iamHrithikRaj/gip/actions/workflows/test.yml/badge.svg)](https://github.com/iamHrithikRaj/gip/actions/workflows/test.yml)
[![Rust Version](https://img.shields.io/badge/rust-1.70%2B-blue.svg)](https://www.rust-lang.org/)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![Crates.io](https://img.shields.io/crates/v/gip.svg)](https://crates.io/crates/gip)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](CONTRIBUTING.md)

**An intelligent Git extension that enriches merge conflicts with structured context for autonomous LLM resolution and human clarity.**

Gip transforms merge conflicts from cryptic puzzles into self-documenting resolution tasks by injecting structured parameters directly into conflict markers—enabling accurate autonomous resolution without human intervention.

---

## The Problem

Traditional Git conflict markers show **raw text differences without semantic context**:

```python
<<<<<<< HEAD
total += item.price * 1.08
=======
total += item.price + 5.99
>>>>>>> feature-branch
```

**What's missing?**
- Is this a bugfix, feature, or side effect?
- What was the rationale behind each change?
- What are the preconditions and postconditions?
- Are there side effects or error handling differences?

**Why this matters:**

Without structured context, **LLMs can't reliably resolve conflicts on their own**. They have to guess intent from commit messages (which are inconsistent and vague) or scan the entire codebase for clues. This makes autonomous conflict resolution unreliable and still requires human intervention.

Humans face the same problem—wasting time context-switching between editors, GitHub, and chat tools to understand what each side of the conflict actually does.

---

## The Solution

Gip introduces **structured manifests** that define clear parameters for conflict resolution—making it possible for LLMs to resolve conflicts autonomously and accurately.

Instead of relying on inconsistent commit messages, Gip stores **structured metadata** alongside each change:

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

**Now you have structured parameters for resolution:**
- ✅ **Behavior Class**: Both are features (not bugfixes)
- ✅ **Rationale**: Tax compliance vs. shipping fee
- ✅ **Postconditions**: One applies 8% tax, the other adds $5.99
- ✅ **Side Effects**: None (both are pure calculations)
- ✅ **Resolution Strategy**: These are additive changes—combine both

**For LLMs:** These structured parameters enable autonomous, reliable conflict resolution without human intervention.

**For Humans:** No need to hunt through commit history or ask teammates—the context is right there.

---

## Why Gip?

### The Motivation

**The real problem isn't that merge conflicts are cryptic—it's that they're not informative.**

Traditional Git only shows raw text differences (`<<<` and `>>>`). Neither humans nor LLMs get the semantic context needed to resolve conflicts intelligently:
- Is this a bugfix, feature, or side effect?
- What was the rationale?
- What are the preconditions and postconditions?
- Are there side effects or breaking changes?

**For LLMs:**
Without structured parameters, AI can't reliably resolve conflicts autonomously. It has to guess from vague commit messages or scan entire codebases, making resolution unreliable and still requiring human intervention.

**For Humans:**
You waste time:
- 🔍 Hunting through commit messages (often vague: "fix bug")
- 📖 Reading through entire PRs and code reviews
- 💬 Messaging teammates asking "what was this change about?"
- 🧠 Context-switching between your editor, GitHub, Slack, and JIRA

### The Solution: Structured Parameters for Intelligent Resolution

Gip introduces **manifests**—structured metadata that defines clear parameters for conflict resolution. Instead of relying on inconsistent commit messages, Gip provides a machine-readable format that LLMs can reliably parse and act on.

**What Gip Captures:**
- **Behavior Class**: Is this a feature, bugfix, refactor, or performance optimization?
- **Rationale**: Why was this change made?
- **Contracts**: Preconditions, postconditions, and error models
- **Side Effects**: Logs, database writes, network calls, etc.
- **Breaking Changes**: Did function signatures change?

**How It Works:**

**For Humans (Optional Context):**
- By default, Gip works exactly like Git
- Use `gip commit -c` to add structured context when it matters
- The extra context is optional—only add it when you want

**For LLMs (Automatic Context):**
- Gip can generate manifests automatically using AI
- Context comes "for free" without manual effort
- When conflicts occur, manifests are injected right into the conflict markers
- LLMs can now resolve conflicts autonomously with high accuracy

**The Result:**
Instead of merging blindly, Gip **enriches conflicts with parameters** that enable intelligent, autonomous resolution—dramatically improving LLM accuracy and eliminating unnecessary human intervention.

### Key Features

- 🤖 **Autonomous LLM Resolution** - Structured parameters enable LLMs to resolve conflicts accurately without human intervention
- 🎯 **Zero Context Switching** - Humans get everything they need, directly in the conflict
- 🔄 **Works Like Git** - By default, Gip is a drop-in replacement for Git
- 📝 **Optional Human Context** - Add manifests with `-c` flag only when you need them
- 🪄 **Automatic AI Context** - Generate manifests automatically for LLM workflows
- 📦 **Single Binary** - 0.4MB, zero dependencies, just copy and run
- 🪶 **Token Efficient** - TOON format is 49% smaller than JSON for LLM contexts
- ⚡ **Lightning Fast** - Built in Rust for maximum performance

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

**Single binary**: The executable is **3.6MB** and has **zero runtime dependencies**. Just copy and run!

### Option 2: Build from Source

**Prerequisites**: Rust 1.70+ and Cargo

```bash
git clone https://github.com/iamHrithikRaj/gip
cd gip
cargo build --release
```

The binary will be available at `target/release/gip` (or `gip.exe` on Windows).

**Dependencies** (automatically managed by Cargo):
- `clap` - CLI framework
- `dialoguer` - Interactive prompts
- `git2` - Git integration
- `serde` & `serde_json` - Serialization
- `anyhow` - Error handling

---

## Quick Start

### 1. Initialize Gip

```bash
cd your-git-repo
gip init
```

This installs Git hooks and configures the repository. You only need to do this once per repo.

### 2. Commit with Context

**Option A: AI-Powered (Recommended)**

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=sk-...

# Stage your changes
git add modified_file.py

# Let AI generate the manifest
gip commit --ai --intent "Add strict validation to user parser"
```

The AI automatically:
- Analyzes your git diff
- Extracts function signatures and changes
- Generates structured contracts (preconditions, postconditions, error models)
- Detects breaking changes
- Identifies feature flags in your code

**Option B: Interactive Mode**

For manual control or when you prefer not to use AI:

```bash
# Stage your changes
git add modified_file.py

# Interactive commit with prompts
gip commit -c
```

Gip intelligently:
- Detects when multiple functions changed (prompts for global intent)
- Allows each function to inherit global intent or provide unique rationale
- Warns you if function signatures changed (potential breaking changes)

You'll be prompted to provide:
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

**Resolution** (LLM or Human):
```python
def calculate_total(items):
    subtotal = sum(item.price * 1.08 for item in items)  # Tax
    return subtotal + 5.99  # Shipping
```

**Why this works:** The structured parameters (behavior class: feature, no side effects, additive postconditions) clearly indicate both changes should be combined.

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

**Resolution** (LLM or Human):
```python
def apply_discount(items, discount=0.1):
    if not items:
        raise ValueError("items cannot be empty")
    return sum(item.price for item in items) * (1 - discount)
```

**Why this works:**
- ✅ **Behavior classes** indicate: bugfix (critical) + feature (enhancement)
- ✅ **Error models** show: HEAD adds validation, MERGE_HEAD has none
- ✅ **Resolution strategy**: Keep both—bugfix prevents crashes, feature adds convenience
- ✅ **LLM confidence**: High—parameters clearly indicate complementary changes

The structured manifest enables autonomous, intelligent resolution instead of blind merging.

---

## Architecture

```
gip/
├── src/
│   ├── main.rs              # CLI entry point with clap
│   ├── lib.rs               # Public API
│   ├── manifest/
│   │   ├── mod.rs           # Module exports
│   │   ├── types.rs         # Manifest schema (structs with serde)
│   │   ├── storage.rs       # JSON save/load operations
│   │   └── toon.rs          # TOON serialization
│   ├── diff/
│   │   ├── mod.rs
│   │   └── analyzer.rs      # Git diff parsing, symbol extraction
│   ├── merge/
│   │   ├── mod.rs
│   │   └── driver.rs        # Conflict detection & enrichment
│   ├── git.rs               # Git operations (hooks, config)
│   ├── prompt/
│   │   ├── mod.rs
│   │   └── interactive.rs   # Interactive CLI prompts (dialoguer)
│   └── toon/
│       ├── mod.rs
│       └── serializer.rs    # TOON format utilities
└── .gip/
    └── manifest/
        └── <commit-sha>.json # JSON manifest storage
```

**Design Decisions**:
- **Rust implementation**: Memory safety, zero-cost abstractions, excellent performance
- **JSON storage**: Easy parsing with serde's serialization
- **TOON display**: Token-efficient format for conflict markers
- **No Git modification**: Manifests stored separately in `.gip/` directory
- **Post-processing**: Merge driver runs after Git merge, enriches existing conflicts
- **Transparent wrapper**: Unknown commands passed to `git` via subprocess

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
| Binary size | ✅ 3.6MB | ✅ ~2MB | ✅ ~2MB | ❌ Cloud |
| Dependencies | ✅ Zero | ✅ Zero | ✅ Zero | ❌ Network |
| Memory Safety | ✅ Rust | 🟡 C | 🟡 C | N/A |

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
A: Minimal. Manifest files are ~1KB each (JSON format). The binary is 3.6MB. No runtime dependencies.

**Q: Can I customize the manifest schema?**  
A: Yes, edit `src/manifest/types.rs` and rebuild with `cargo build --release`.

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

# Build
cargo build --release

# Run tests
cargo test
```

### Testing

Gip has comprehensive unit tests for all core modules:

```bash
# Run all tests
cargo test

# Run tests with output
cargo test -- --nocapture

# Run specific module tests
cargo test manifest

# Generate coverage report (requires cargo-tarpaulin)
cargo tarpaulin --out Html
```

**Test Coverage:**
- ✅ `manifest` - 22+ tests (types, storage, TOON serialization, migration)
- ✅ `merge` - Conflict detection and enrichment tests
- ✅ `diff` - Symbol extraction and change detection tests
- ✅ `git` - Repository operations tests

**CI/CD:**
- Tests run automatically on push/PR via GitHub Actions
- Multi-OS testing (Windows, Linux, macOS)
- Rust stable and beta versions tested

See [TEST_STRUCTURE_PROPOSAL.md](TEST_STRUCTURE_PROPOSAL.md) for full test architecture.

---

## What's Next?

Gip is in active development with an ambitious roadmap:

### Currently Available
- ✅ Core manifest creation and storage
- ✅ Structured contract documentation
- ✅ TOON format serialization (49% token savings)
- ✅ Multi-language symbol extraction
- ✅ Comprehensive test suite (24+ passing tests)

### Roadmap
- **Gip MCP** - Model Context Protocol server for seamless LLM integration
- **VS Code Extension** - Inline manifest editing and visualization
- **GitHub Actions** - Automated manifest generation in CI/CD
- **LLM-Assisted Resolution** - Get AI suggestions for resolving conflicts
- **Team Templates** - Share manifest templates across your organization
- **Semantic Merge** - Auto-resolve conflicts based on contract compatibility
- **Tool Integration** - Work with Meld, KDiff3, and other merge tools

**Want to help shape the future?** Check out [open issues](https://github.com/iamHrithikRaj/gip/issues) or [contribute](CONTRIBUTING.md)!

---

## Contributing

We welcome contributions! Gip is an open-source project that benefits from community involvement.

### Quick Start

1. **Fork and clone** the repository
2. **Install Rust 1.70+** (via [rustup](https://rustup.rs/))
3. **Make your changes** following our coding standards
4. **Add tests** for new functionality (TDD approach)
5. **Submit a pull request**

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines on:
- Setting up your development environment
- Rust coding standards and best practices
- Testing requirements (test-driven development)
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

- **TOON Format**: [toon-format crate](https://crates.io/crates/toon-format) - Official spec-compliant Rust implementation
  - Specification: [TOON v2.0](https://github.com/toon-format/spec)
  - Token-efficient format for LLM-friendly data serialization
- **Inspired by**: Git Notes, semantic merge tools, and the need for better conflict resolution UX
- **Built with**: Rust, for memory safety and performance

---

## Support

- 📖 **Documentation**: See this README and [CONTRIBUTING.md](CONTRIBUTING.md)
- 🐛 **Bug Reports**: [Open an issue](https://github.com/iamHrithikRaj/gip/issues/new?template=bug_report.md)
- 💡 **Feature Requests**: [Suggest a feature](https://github.com/iamHrithikRaj/gip/issues/new?template=feature_request.md)
- ❓ **Questions**: [Ask a question](https://github.com/iamHrithikRaj/gip/issues/new?template=question.md)