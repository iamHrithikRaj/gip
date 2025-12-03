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
- **AI Agents** fail to resolve conflicts because git markers lack semantic data
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
3.  **The Generation**: Since the manifest is structured (TOON/JSON), **LLMs can automatically generate it** by analyzing the diff.
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
git clone https://github.com/AZBucky/gip.git
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

## Quick Start

### 1. Initialize in Your Repository

```bash
cd your-project
gip init
# ✓ Gip initialized successfully
# Created: .gip/
# Created: manifest.toon (template)
```

### 2. Create a Manifest

Edit `manifest.toon` to describe your changes:

```yaml
# manifest.toon
schemaVersion: "2.0"
globalIntent:
  behaviorClass: ["perf"]
  rationale: "Add user preference caching to reduce API calls"
entries:
  - anchor:
      file: "src/cache.rs"
      symbol: "UserCache"
    changeType: "add"
    rationale: "Implemented LRU cache for user preferences"
    contract:
      postconditions: ["Preferences cached for 5 minutes"]
    compatibility:
      breaking: false
```

### 3. Commit with Context

```bash
gip commit -m "perf: add preference caching"
# ✓ Manifest validated
# ✓ Changes committed with context
# ✓ Manifest attached as git note
```

### 4. Access Context Later

```bash
# View manifest for any commit
gip context abc1234

# Export for Agents
gip context --export --json > context.json

# View all behavioral changes since last release
gip context --since v2.0.0 --behavior
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
| **`init`** | Initialize Gip | `gip init` | Creates `.gip/` and `manifest.toon` template. |
| **`commit`** | Commit with Context | `gip commit -m "msg"` | Validates manifest, commits changes, and attaches Git Note. |
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

**Example of enriched conflict markers:**

```cpp
<<<<<<< HEAD
total += item.price * 1.08
||| Gip CONTEXT (HEAD - Your changes)
||| behaviorClass: feature
||| rationale: Added 8% sales tax to comply with state law
||| postconditions[0]: returns total with 8% sales tax
||| symbol: calculate_total
=======
total += item.price + 5.99
||| Gip CONTEXT (feature-branch - Their changes)
||| behaviorClass: feature
||| rationale: Added $5.99 flat shipping fee for all orders
||| postconditions[0]: returns total plus flat shipping
||| symbol: calculate_total
>>>>>>> feature-branch
```

### The `context` Command

The `context` command is the bridge between your git history and AI agents.

| Usage | Description |
| :--- | :--- |
| `gip context` | Show the human-readable manifest for the current `HEAD`. |
| `gip context <sha>` | Show the manifest for a specific commit. |
| `gip context --json` | Output in JSON format (for CI/CD pipelines or scripts). |
| `gip context --export` | Export raw TOON format (optimized for LLM context windows). |

### Git Passthrough

Any command not listed above is passed directly to git, so you can use `gip` as your daily driver:

```bash
gip status
gip log --oneline --graph
gip checkout -b feature/new-ui
```

---

## Examples

### 1. The Agentic Workflow (Ask LLM on Push)

Gip is designed to be the "missing link" for AI agents. When an agent (or human) tries to commit without context, Gip intervenes to ensure knowledge is captured.

**Step 1: Agent tries to commit**
```bash
gip commit -m "refactor: update user schema"
# ✗ Commit Rejected: Missing Context Manifest
# 
# Please retry with this block appended to your commit message:
# gip:
# {
#   schemaVersion: "2.0",
#   entries: [
#     {
#       anchor: { file: "src/user.rs", symbol: "User" },
#       behaviorClass: ["refactor"],
#       rationale: "<why?>",
#       compatibility: { breaking: false }
#     }
#   ]
# }
```

**Step 2: Agent analyzes code & fills manifest**
The agent reads the diff, understands the semantic impact, and fills out the manifest—including critical safety info like breaking changes.

```bash
gip commit -m "refactor: update user schema

gip:
{
  schemaVersion: "2.0",
  entries: [
    {
      anchor: { file: "src/user.rs", symbol: "User::validate" },
      behaviorClass: ["refactor"],
      rationale: "Enforce strict password complexity rules per security audit",
      compatibility: {
        breaking: true,
        migrations: ["Run db_migrate_v2.sql to update existing hashes"]
      },
      contract: {
        inputs: ["password: string", "options: ValidationOptions"],
        errorModel: ["throws WeakPasswordException"]
      }
    }
  ]
}"
# ✓ Committed with manifest: 7f2a9d1
```

### 2. Safe Conflict Resolution

Traditional git markers offer no semantic context, making it impossible for LLMs to reliably resolve logical conflicts. Gip solves this by injecting **contracts** and **breaking changes** directly into the markers, providing the explicit ground truth agents need for safe, autonomous resolution.

**Scenario:**
*   **HEAD**: Adds a new parameter to `processPayment` (Breaking Change).
*   **Feature Branch**: Calls `processPayment` with the old signature.

**Standard Git Conflict:**
```cpp
<<<<<<< HEAD
void processPayment(float amount, Currency currency) {
=======
void processPayment(float amount) {
>>>>>>> feature-branch
```
*Result: Developer blindly keeps HEAD, but Feature Branch code elsewhere breaks.*

**Gip Enriched Conflict:**
```cpp
<<<<<<< HEAD
void processPayment(float amount, Currency currency) {
||| Gip CONTEXT (HEAD - Your changes)
||| behaviorClass: refactor
||| rationale: Added currency support for internationalization
||| breaking: true
||| migrations[0]: Update all callsites to pass Currency::USD by default
||| inputs[0]: amount: float
||| inputs[1]: currency: Currency
||| symbol: processPayment
=======
void processPayment(float amount) {
||| Gip CONTEXT (feature-branch - Their changes)
||| behaviorClass: feature
||| rationale: Simplified payment flow for guest checkout
||| breaking: false
||| symbol: processPayment
>>>>>>> feature-branch
```

**Resolution:** The LLM (or human) sees `breaking: true` and the `migrations` instruction. It knows it **must** update the feature branch's logic to match the new signature, rather than just picking one side.

### 3. Semantic Knowledge Base

Stop guessing why code exists. Query the repository for intent.

```bash
gip context src/auth_service.cpp
```

**Output:**
```text
┌─ Commit 8b12f9a (2025-11-15 by hrithik)
│
│  feat: implement rate limiting
│
│  Intent: security
│  Rationale: Prevent brute force attacks observed in logs
│  Preconditions: Redis connection must be active
│  Side Effects: Writes to 'rate_limit_keys' in Redis
│  Error Model: Returns 429 Too Many Requests if quota exceeded
└───────────────────────────────────────────────────────────────
```

This acts as a **RAG-ready knowledge base** for your LLM, allowing it to answer questions like *"Why do we use Redis for rate limiting?"* or *"What are the side effects of the auth module?"* without hallucinating.

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
│  │  (git operations)   │  │   (TOON/JSON/YAML support)   │  │
│  └─────────────────────┘  └──────────────────────────────┘  │
│                                        │                     │
│                                        v                     │
│                              ┌──────────────────┐            │
│                              │  ctoon (static)  │            │
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

Gip uses [Git notes](https://git-scm.com/docs/git-notes) to store manifests:

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
| Structured format | ✅ TOON/JSON | ❌ Free-form | ⚠️ Markdown |
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

Gip works without manifests. Commits without manifests are just regular git commits. You can add context later:

```bash
gip context --add <sha>
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
git clone https://github.com/AZBucky/gip.git
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

- [ctoon](ctoon/) - TOON format serialization library
- [gip-legacy](https://github.com/AZBucky/gip-legacy) - Original Rust implementation

---

## License

MIT License - See [LICENSE](LICENSE) for details.

---

<div align="center">

**Built for the Agentic developer.**

[Report Bug](https://github.com/AZBucky/gip/issues) · [Request Feature](https://github.com/AZBucky/gip/issues) · [Documentation](docs/)

</div>

