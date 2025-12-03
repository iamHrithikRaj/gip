<div align="center">

# Gip (Git++)

### Git with Intent Preservation

**The missing link between your code changes and the context that explains them.**

[![C++17](https://img.shields.io/badge/C++-17-blue.svg?style=flat&logo=cplusplus)](https://en.cppreference.com/w/cpp/17)
[![CMake](https://img.shields.io/badge/CMake-3.15+-064F8C.svg?style=flat&logo=cmake)](https://cmake.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
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

**Gip captures structured semantic context alongside your commits.**

```bash
# Instead of just committing...
git commit -m "refactor: extract payment service"

# Commit with intent preservation
gip commit -m "refactor: extract payment service"
```

Gip attaches a **manifest**—a machine-readable document describing the *intent* behind changes—to your commits via Git notes. This context travels with your repository, invisible to normal workflows but available when you need it.

```yaml
# What Gip captures (manifest.toon)
intent: "Extract payment logic to enable multi-provider support"
behavior: "No functional changes - pure refactoring"
testing: "All existing payment tests pass unchanged"
reviewers: ["@payments-team"]
decisions:
  - "Chose strategy pattern over factory for provider extensibility"
  - "Kept synchronous API despite async provider requirements"
```

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

### One-Line Install (Build from Source)

You can install Gip with a single command. This will clone the repository to a temporary location, build it, and install it to your user folder.

**Windows (PowerShell):**
```powershell
Set-ExecutionPolicy Bypass -Scope Process -Force; irm https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.ps1 | iex
```

**Linux / macOS:**
```bash
curl -sL https://raw.githubusercontent.com/iamHrithikRaj/gip/main/scripts/bootstrap.sh | bash
```

### Manual Install

If you prefer to inspect the code first:

**Windows (PowerShell):**
```powershell
git clone --recursive https://github.com/iamHrithikRaj/gip.git
cd gip
.\scripts\install.ps1
```

**Linux / macOS:**
```bash
git clone --recursive https://github.com/iamHrithikRaj/gip.git
cd gip
./scripts/install.sh
```

### Pre-built Binaries

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
- C++17 compatible compiler (MSVC 19.14+, GCC 8+, Clang 7+)
- CMake 3.15+
- Git 2.28+ (for worktree support)

```bash
git clone --recursive https://github.com/AZBucky/gip.git
cd gip

# 1. Configure the project
cmake -S . -B build

# 2. Build the project (Release mode)
cmake --build build --config Release

# 3. Run Tests
ctest --test-dir build --output-on-failure

# 4. Install (optional)
cmake --install build
```

**Build Options:**

| Option | Default | Description |
|--------|---------|-------------|
| `GIP_BUILD_TESTS` | `ON` | Build unit tests |
| `GIP_USE_LIBGIT2` | `OFF` | Use libgit2 instead of git CLI |
| `GIP_ENABLE_SANITIZERS` | `OFF` | Enable ASan/UBSan |

### Verify Installation

```bash
gip --version
# gip 1.0.0 (c++ edition)

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
version: "1.0"
intent: "Add user preference caching to reduce API calls"
behavior:
  changes: "Preferences now cached for 5 minutes"
  breaking: false
testing:
  required: ["unit", "integration"]
  coverage: "New cache layer fully tested"
decisions:
  - rationale: "Chose 5-minute TTL based on user session analytics"
    alternatives: ["1-minute TTL", "No caching", "Persistent cache"]
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

The manifest is a structured document in TOON format:

```yaml
version: "1.0"                    # Schema version (required)

intent: "string"                  # Why are you making this change? (required)

behavior:                         # How does this change behavior?
  changes: "string"               # Description of behavioral changes
  breaking: bool                  # Is this a breaking change?
  
testing:                          # Testing requirements
  required: ["unit", "e2e"]       # Required test types
  coverage: "string"              # Coverage notes
  
dependencies:                     # Dependency changes
  added: ["package@version"]
  removed: ["package"]
  updated: ["package@old->new"]
  
decisions:                        # Design decisions (array)
  - rationale: "string"           # Why this decision?
    alternatives: ["string"]      # What else was considered?
    tradeoffs: "string"           # What are the tradeoffs?

reviewers: ["@username"]          # Suggested reviewers
context: "string"                 # Additional context or links
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

### Core Commands

```bash
gip init                          # Initialize Gip in repository
gip commit [-m "msg"]             # Commit with manifest attachment
gip push [--with-notes]           # Push commits and notes to remote
gip context [<commit>]            # Display context for commit(s)
gip merge <branch>                # Merge with enriched conflict markers
gip rebase <branch>               # Rebase with enriched conflict markers
```

### Merge & Rebase Commands

Gip wraps `git merge` and `git rebase` to provide **enriched conflict markers** when conflicts occur:

```bash
gip merge feature-branch          # Merge with intent-enriched conflicts
gip rebase main                   # Rebase with intent-enriched conflicts
gip rebase -i HEAD~3              # Interactive rebase (conflicts enriched)
```

When a conflict occurs, Gip automatically:
1. Detects conflicted files
2. Looks up manifests for both sides of the conflict
3. Injects structured context into conflict markers

**Example of enriched conflict markers:**

```
<<<<<<< HEAD
total += item.price * 1.08
||| Gip CONTEXT (HEAD - Your changes)
||| Commit: 3c7d3422
||| behaviorClass: feature
||| rationale: Added 8% sales tax to comply with state law
||| postconditions[0]: returns total with 8% sales tax
||| sideEffects: none
||| symbol: calculate_total
=======
total += item.price + 5.99
||| Gip CONTEXT (feature-branch - Their changes)
||| Commit: 1825fc7c
||| behaviorClass: feature
||| rationale: Added $5.99 flat shipping fee for all orders
||| postconditions[0]: returns total plus flat shipping
||| sideEffects: none
||| symbol: calculate_total
>>>>>>> feature-branch
```

**Now LLMs can understand:**
- **behaviorClass**: Both are features (not bugfixes)
- **rationale**: Tax compliance vs. shipping fee
- **postconditions**: One applies 8% tax, the other adds $5.99
- **Resolution**: These are additive changes—combine both!

### Context Command Options

```bash
gip context                       # Show context for HEAD
gip context <sha>                 # Show context for specific commit
gip context --all                 # Show context for all commits
gip context --since <ref>         # Show context since reference
gip context --behavior            # Show only behavioral changes
gip context --export              # Export full context document
gip context --json                # Output in JSON format
gip context --toon                # Output in TOON format (default)
```

### Git Passthrough

Any unrecognized command passes through to Git:

```bash
gip status                        # → git status
gip log --oneline                 # → git log --oneline
gip branch -a                     # → git branch -a
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
#       file: "src/user.cpp",
#       behavior: "<feature|bugfix|refactor>",
#       rationale: "<why?>",
#       breaking: false,
#       ...
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
  entries: [
    {
      file: "src/user.cpp",
      symbol: "User::validate",
      behavior: "refactor",
      rationale: "Enforce strict password complexity rules per security audit",
      breaking: true,
      migrations: ["Run db_migrate_v2.sql to update existing hashes"],
      inputs: ["password: string", "options: ValidationOptions"],
      errorModel: ["throws WeakPasswordException"]
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
| **Git Adapter** | Wraps git CLI with typed C++ interface | `src/git_adapter.cpp` |
| **Manifest Parser** | Parses and validates manifest files | `src/manifest.cpp` |
| **Command Handlers** | Individual command implementations | `src/commands/` |
| **ctoon** | TOON format serialization (static library) | `ctoon/` |

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
# Clone with submodules
git clone --recursive https://github.com/AZBucky/gip.git
cd gip

# Build with tests
cmake -B build -DGIP_BUILD_TESTS=ON
cmake --build build

# Run tests
ctest --test-dir build
```

### Code Standards

- **C++17** with modern idioms
- **Enterprise naming**: Full descriptive names (`isRepository()` not `isRepo()`)
- **Const correctness**: Use `const` wherever possible
- **RAII**: Resource management via constructors/destructors
- **Error handling**: Return types over exceptions for expected failures

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

