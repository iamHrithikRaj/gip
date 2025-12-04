<div align="center">

# Gip (Git++)

### Git with Intent Preservation

**The missing link between your code changes and the context that explains them.**

[![CI](https://github.com/iamHrithikRaj/gip/actions/workflows/ci.yml/badge.svg)](https://github.com/iamHrithikRaj/gip/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/iamHrithikRaj/gip?style=flat&logo=github)](https://github.com/iamHrithikRaj/gip/releases)
[![Downloads](https://img.shields.io/github/downloads/iamHrithikRaj/gip/total?style=flat&logo=github)](https://github.com/iamHrithikRaj/gip/releases)
[![Rust](https://img.shields.io/badge/Rust-1.75+-orange.svg?style=flat&logo=rust)](https://www.rust-lang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
<br>
[![Windows](https://img.shields.io/badge/Windows-0078D6?style=flat&logo=windows)](https://microsoft.com/windows)
[![Linux](https://img.shields.io/badge/Linux-FCC624?style=flat&logo=linux&logoColor=black)](https://kernel.org)
[![macOS](https://img.shields.io/badge/macOS-000000?style=flat&logo=apple)](https://apple.com/macos)

[Installation](#installation) • [Quick Start](#quick-start) • [Examples](#examples) • [Architecture](#architecture) • [Contributing](#contributing)

</div>

---

## The Problem

Every day, millions of developers write commit messages like:

```
fix: resolve issue with user authentication
```

Then six months later, someone asks: *"Why did we change the authentication flow?"*

The commit message says **what** changed. The diff shows **how** it changed. But the **why**—the intent, the trade-offs, the alternatives considered—is lost forever.

This matters because:

- **Code reviews** become archaeology expeditions
- **AI Agents** lack the confidence to resolve conflicts autonomously because git markers lack semantic data
- **Merge conflicts** are resolved by guessing intent
- **Onboarding** requires tribal knowledge transfer
- **Refactoring** risks breaking implicit assumptions

## The Solution

**Gip is the protocol for Agentic Version Control.**

It solves the context problem by enforcing a **structured handshake** between the coder (human or AI) and the repository.

Instead of accepting a simple text message, Gip expects a **Manifest** (`manifest.toon`)—a machine-readable declaration of *Intent*, *Behavior*, and *Contracts*.

**The Agentic Workflow:**
1.  **The Request**: You (or an agent) attempt to commit changes.
2.  **The Protocol**: Gip demands a `manifest.toon` describing the semantic impact of those changes.
3.  **The Generation**: Since the manifest is structured (TOON), **LLMs can automatically generate it** by analyzing the diff.
4.  **The Preservation**: Gip attaches this structured context to the commit history, creating a permanent, queryable record.

This turns your git history from a flat list of changes into a **semantic knowledge base**, enabling agents to understand the *evolution* of your software, not just its current state.

## Why Gip?

### For Developers

| Without Gip | With Gip |
|-------------|----------|
| "Why did someone change this?" | Intent documented in manifest |
| Agents generate code without context | Agents understand design decisions |
| Merge conflicts = guesswork | Semantic hints for resolution |
| Onboarding takes months | Self-documenting codebase |

### For Teams

- **Consistent Context**: Every commit follows the same structured format
- **Searchable Intent**: Find commits by behavioral impact, not just text
- **Review Efficiency**: Reviewers see intent before diving into diffs
- **Knowledge Preservation**: Decisions outlive team members

### For Agentic Development

Gip was designed for the Agentic future, enabling **true agentic workflows**:

1.  **Context Injection**: Agents can query `gip context` to understand the codebase before making changes.
2.  **Structured Commits**: Gip forces agents to document their intent via the manifest rejection loop.
3.  **Autonomous Resolution**: Enriched conflict markers provide the ground truth agents need to resolve merges without human help.

```bash
# Export context for your Agent
gip context --export > context.toon
```

The TOON format used by Gip is **49% smaller than JSON** and **23% smaller than YAML**, making it ideal for LLM context windows.

---

## Installation

### One-Line Install

You can install Gip with a single command. This will download the latest pre-built binary for your platform and install it to your user folder.

**Windows (PowerShell):**
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; irm https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.ps1 | iex
```

**Linux / macOS:**
```bash
curl -sL https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.sh | bash
```

### Manual Install (Build from Source)

If you prefer to inspect the code first or build from source:

**Windows (PowerShell):**
```powershell
git clone https://github.com/iamHrithikRaj/gip.git
cd gip
.\scripts\install.ps1
```

**Linux / macOS:**
```bash
git clone https://github.com/iamHrithikRaj/gip.git
cd gip
./scripts/install.sh
```

### Pre-built Binaries (Manual Download)

Download the latest release for your platform:

| Platform | Architecture | Download |
|----------|--------------|----------|
| Windows  | x64          | [gip-windows-x64.zip](https://github.com/iamHrithikRaj/gip/releases/latest) |
| Linux    | x64          | [gip-linux-x64.tar.gz](https://github.com/iamHrithikRaj/gip/releases/latest) |
| macOS    | x64/ARM64    | [gip-macos-universal.tar.gz](https://github.com/iamHrithikRaj/gip/releases/latest) |

```powershell
# Windows (PowerShell)
Expand-Archive gip-windows-x64.zip -DestinationPath $env:LOCALAPPDATA\Programs\gip
$env:PATH += ";$env:LOCALAPPDATA\Programs\gip"
```

```bash
# Linux/macOS
tar -xzf gip-*.tar.gz
sudo mv gip /usr/local/bin/
```

### Build from Source

**Requirements:**
- Rust 1.75+ (stable)
- Git 2.28+ (for worktree support)

```bash
git clone https://github.com/iamHrithikRaj/gip.git
cd gip

# 1. Build the project (Release mode)
cargo build --release

# 2. Run Tests
cargo test

# 3. Install
cargo install --path .
```

### Verify Installation

```bash
gip --version
# gip 1.0.0

gip --help
# Git with Intent Preservation - Context-aware git wrapper
```

### Troubleshooting (Windows)

**Conflict with `Get-NetIPConfiguration`**

PowerShell has a default alias `gip` for `Get-NetIPConfiguration`. The installation script will attempt to remove this alias for you. If you still see network information when running `gip`, you can manually add this to your PowerShell profile:

```powershell
Remove-Item alias:gip -ErrorAction SilentlyContinue
```

---

## Using Gip

### 1. The Commit Workflow (Capturing Intent)

This is the primary loop for developers and agents.

1.  **Initialize**: Run `gip init` once to set up the repo. This creates `.gip/manifest.toon` and adds `.gip` to `.gitignore`.
2.  **Generate Manifest**:
    *   **Humans**: Edit `.gip/manifest.toon` manually.
    *   **Agents**: Use the diff to automatically generate the manifest. Gip provides the prompt structure in the template.
    ```yaml
    # .gip/manifest.toon (Generated by Agent)
    schemaVersion: "2.0"
    globalIntent:
      behaviorClass: ["perf"]
      rationale: "Add user preference caching to reduce API calls"
    entries:
      - anchor: { file: "src/cache.rs", symbol: "UserCache" }
        changeType: "add"
        rationale: "Implemented LRU cache for user preferences"
        contract:
          postconditions: ["Preferences cached for 5 minutes"]
        compatibility: { breaking: false }
    ```
3.  **Commit**: Run `gip commit -m "perf: add preference caching"`.
    *   Gip validates the manifest.
    *   Gip commits the code.
    *   Gip attaches the manifest as a Git Note.
    *   *Note*: If the manifest is missing or incomplete, the commit is rejected with instructions for the Agent/LLM. Use `--force` to bypass.

### 2. The Conflict Resolution Workflow (Enriched Markers)

This is where Gip shines for autonomous agents. When a merge conflict occurs, Gip injects semantic context directly into the conflict markers.

```bash
gip merge feature-branch
# or
gip rebase main
```

**Standard Git Conflict:**
```rust
<<<<<<< HEAD
fn process_payment(amount: f32, currency: Currency) {
=======
fn process_payment(amount: f32) {
>>>>>>> feature-branch
```
*Result: Agent guesses which one is correct.*

**Gip Enriched Conflict:**
```rust
<<<<<<< HEAD
fn process_payment(amount: f32, currency: Currency) {
||| Gip CONTEXT (HEAD - Your changes)
||| behaviorClass: refactor
||| rationale: Added currency support for internationalization
||| breaking: true
||| migrations[0]: Update all callsites to pass Currency::USD by default
||| symbol: process_payment
=======
fn process_payment(amount: f32) {
||| Gip CONTEXT (feature-branch - Their changes)
||| behaviorClass: feature
||| rationale: Simplified payment flow for guest checkout
||| breaking: false
||| symbol: process_payment
>>>>>>> feature-branch
```
*Result: Agent sees `breaking: true` and `migrations` instructions, allowing it to correctly update the feature branch code to match the new signature.*

### 3. The Context Workflow (Querying Knowledge)

Turn your git history into a RAG-ready knowledge base.

*   **For Humans**: `gip context <sha>` shows a readable summary of *why* a change happened.
*   **For Agents**: `gip context --export` dumps the full semantic history in TOON format, perfect for context injection before starting a task.

```bash
# "Why do we use Redis here?"
gip context src/auth_service.rs
```
**Output:**
```text
Intent: security
Rationale: Prevent brute force attacks observed in logs
Side Effects: Writes to 'rate_limit_keys' in Redis
```

---

## Manifest Schema

The manifest is a structured document in TOON format (v2.0):

```yaml
schemaVersion: "2.0"              # Schema version (required)

globalIntent:                     # Commit-level intent
  behaviorClass: ["feature"]      # [feature, bugfix, refactor, perf, security...]
  rationale: "string"             # High-level explanation

entries:                          # List of changes
  - anchor:
      file: "string"              # File path
      symbol: "string"            # Function/Class name
    changeType: "modify"          # [add, modify, delete, rename]
    rationale: "string"           # Why this specific change?
    
    contract:                     # Behavioral contract
      preconditions: ["string"]
      postconditions: ["string"]
      errorModel: ["string"]
      
    compatibility:                # Compatibility flags
      breaking: bool
      migrations: ["string"]      # Migration instructions
      
    decisions:                    # Design decisions
      - rationale: "string"
        alternatives: ["string"]
```

### Validation

Gip validates manifests before committing:

```bash
gip commit -m "feat: add feature"
# ✗ Manifest validation failed:
#   - Missing required field: 'intent'
#   - 'behavior.breaking' must be boolean
```

---

## CLI Reference

### Command Overview

| Command | Purpose | Usage Example | Key Features |
| :--- | :--- | :--- | :--- |
| **`init`** | Initialize Gip | `gip init` | Creates `.gip/` and `.gip/manifest.toon` template. |
| **`commit`** | Commit with Context | `gip commit -m "msg"` | Validates manifest, commits changes, and attaches Git Note. Enforces manifest presence. |
| **`push`** | Push Code + Notes | `gip push` | Pushes commits AND `refs/notes/gip` to remote. |
| **`context`** | Query Intent | `gip context <sha>` | Retrieves semantic context for humans or agents. |
| **`merge`** | Smart Merge | `gip merge feature` | Injects intent & contracts into conflict markers. |
| **`rebase`** | Smart Rebase | `gip rebase main` | Preserves intent during rebase conflicts. |

### Merge & Rebase (Enriched Conflicts)

Gip wraps `git merge` and `git rebase` to provide **enriched conflict markers** when conflicts occur. This is critical for agents to resolve conflicts autonomously.

```bash
gip merge feature-branch
gip rebase main
```

When a conflict occurs, Gip automatically:
1. Detects conflicted files
2. Looks up manifests for both sides of the conflict
3. Injects structured context into conflict markers (see [Conflict Resolution Workflow](#2-the-conflict-resolution-workflow-enriched-markers))

### The `context` Command

The `context` command is the bridge between your git history and AI agents.

| Usage | Description |
| :--- | :--- |
| `gip context` | Show the human-readable manifest for the current `HEAD`. |
| `gip context <sha>` | Show the manifest for a specific commit. |
| `gip context --export` | Export raw TOON format (optimized for LLM context windows). |

### Git Passthrough

Any command not listed above is passed directly to git, so you can use `gip` as your daily driver:

```bash
gip status
gip log --oneline --graph
gip checkout -b feature/new-ui
```

---



---

## Architecture

```
┌──────────────────────────────────────────────────────────────┐
│                         Gip CLI                              │
│  ┌─────────┐  ┌─────────┐  ┌─────────┐  ┌─────────────────┐  │
│  │  init   │  │ commit  │  │  push   │  │    context      │  │
│  └────┬────┘  └────┬────┘  └────┬────┘  └────────┬────────┘  │
│       │            │            │                │           │
│       v            v            v                v           │
│  ┌────────────────────────────────────────────────────────┐  │
│  │                    Command Router                      │  │
│  └────────────────────────────────────────────────────────┘  │
│       │            │            │                │           │
│       v            v            v                v           │
│  ┌─────────────────────┐  ┌──────────────────────────────┐  │
│  │    Git Adapter      │  │      Manifest Parser         │  │
│  │  (git operations)   │  │       (TOON support)         │  │
│  └─────────────────────┘  └──────────────────────────────┘  │
│                                        │                     │
│                                        v                     │
│                              ┌──────────────────┐            │
│                              │   toon-format    │            │
│                              │  TOON serializer │            │
│                              └──────────────────┘            │
└──────────────────────────────────────────────────────────────┘
                    │
                    v
         ┌──────────────────┐
         │   Git Repository │
         │  ┌────────────┐  │
         │  │  commits   │  │
         │  │  + notes   │◄─┼── Manifests stored as Git notes
         │  └────────────┘  │
         └──────────────────┘
```

### Key Components

| Component | Purpose | Location |
|-----------|---------|----------|
| **Git Adapter** | Wraps git CLI with typed Rust interface | `src/git.rs` |
| **Manifest Parser** | Parses and validates manifest files | `src/manifest/mod.rs` |
| **Command Handlers** | Individual command implementations | `src/commands/` |
| **toon** | TOON format serialization | `src/toon/` |

### Storage Model

Gip uses [Git notes](https://git-scm.com/docs/git-notes) to store manifests in TOON format:

```bash
# Manifests stored in refs/notes/gip
git notes --ref=gip show <commit>
```

This approach:
- ✅ Travels with repository clones (with `--notes` fetch)
- ✅ Invisible to normal git workflows
- ✅ Doesn't pollute commit history
- ✅ Can be synced independently of commits
- ✅ Works with existing Git infrastructure

---

## Comparison

| Feature | Gip | Conventional Commits | GitHub PR Templates |
|---------|-----|---------------------|---------------------|
| Structured format | ✅ TOON | ❌ Free-form | ⚠️ Markdown |
| Machine-readable | ✅ Yes | ⚠️ Parseable | ❌ No |
| Travels with repo | ✅ Git notes | ✅ Message | ❌ Platform-specific |
| Design decisions | ✅ First-class | ❌ No | ⚠️ Template-dependent |
| Agent-Ready | ✅ Yes | ❌ No | ❌ No |
| Breaking change flag | ✅ Boolean field | ⚠️ `!` convention | ❌ No |
| Validation | ✅ Schema-based | ❌ No | ❌ No |

---

## FAQ

### Does Gip replace commit messages?

No. Gip enhances commit messages with structured context. You still write normal commit messages—GIP adds the manifest as supplementary data.

### What if I forget to create a manifest?

Gip enforces manifest creation by default. If you try to commit without filling out `.gip/manifest.toon`, the commit will be rejected with explicit instructions for you (or your Agent).

**This enables a fully autonomous loop:**
1. Agent tries to commit.
2. Gip rejects commit and prints instructions to `stdout`.
3. Agent reads instructions, fills out the manifest, and retries.
4. **No human intervention is required.**

**Example Rejection Output:**
```text
ERROR: Commit rejected due to missing or incomplete manifest.
Reason: Manifest file was missing. Created new template at .gip/manifest.toon

INSTRUCTIONS FOR AGENT/LLM:
1. Read the file at: .gip/manifest.toon
2. Understand the code changes you are committing.
3. Fill out the 'rationale', 'changeType', and 'behaviorClass' fields.
4. Save the file.
5. Retry the commit command.
```

If you absolutely must commit without context, use the force flag:
```bash
gip commit --force -m "emergency fix"
```

### How do I sync notes with remotes?

```bash
# Fetch notes from remote
git fetch origin refs/notes/gip:refs/notes/gip

# Push notes to remote
gip push --with-notes
# or: git push origin refs/notes/gip
```

### Is Gip compatible with my Git workflow?

Yes. Gip is a transparent wrapper around git. It doesn't modify git's behavior—only adds the notes mechanism. Works with:
- GitHub, GitLab, Azure DevOps, Bitbucket
- Git Flow, GitHub Flow, Trunk-based development
- Pre-commit hooks, CI/CD pipelines
- Any git GUI that supports notes

### What Git versions are supported?

Gip requires Git 2.28+ for full functionality (worktree support). Core features work with Git 2.17+.

### Can I use Gip with existing repositories?

Absolutely. Run `gip init` in any git repository. Gip won't modify existing history—it only adds context to future commits.

### What's the performance impact?

Minimal. Gip adds ~50ms per commit (manifest parsing + note attachment). The binary is a single ~2MB executable with no runtime dependencies.

---

## Contributing

We welcome contributions! Gip is built with enterprise-grade standards in mind.

### Development Setup

```bash
# Clone the repository
git clone https://github.com/iamHrithikRaj/gip.git
cd gip

# Build with tests
cargo build --tests

# Run tests
cargo test
```

### Code Standards

- **Rust 2021** with modern idioms
- **Enterprise naming**: Full descriptive names (`is_repository()` not `is_repo()`)
- **Error handling**: Result types over panics for expected failures

See [CODING_STANDARDS.md](docs/CODING_STANDARDS.md) for detailed guidelines.

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feat/amazing-feature`)
3. **Write a manifest** for your changes (`manifest.toon`)
4. Commit with Gip (`gip commit -m "feat: add amazing feature"`)
5. Push and open a PR

### Issue Templates

- **Bug Report**: Include version, OS, reproduction steps
- **Feature Request**: Describe use case and proposed solution
- **Documentation**: Point to specific gaps or errors

---

## Roadmap

- [x] **Merge & Rebase**: Intent-aware conflict resolution with enriched markers ✓
- [ ] **VS Code Extension**: Manifest editor with IntelliSense
- [ ] **GitHub Action**: CI validation of manifests
- [ ] **Language Server**: Real-time manifest validation
- [ ] **libgit2 Backend**: Native git operations (no CLI dependency)
- [ ] **Cherry-pick Support**: Enriched conflicts for cherry-pick operations

---

## Related Projects

- [toon-format](https://crates.io/crates/toon-format) - TOON format serialization library
- [gip-legacy](https://github.com/AZBucky/gip-legacy) - Original Rust implementation

---

## License

MIT License - See [LICENSE](LICENSE) for details.

---

<div align="center">

**Built for the Agentic developer.**

[Report Bug](https://github.com/iamHrithikRaj/gip/issues) · [Request Feature](https://github.com/iamHrithikRaj/gip/issues) · [Documentation](docs/)

</div>

