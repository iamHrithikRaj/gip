# Gip Architecture

## Overview

Gip (Git with Intent Protocol) is a drop-in replacement for Git that enforces semantic context for LLM-native development workflows.

## Design Principles

1. **Lightweight** - Minimal dependencies, fast execution
2. **Non-invasive** - Works alongside existing Git workflows
3. **LLM-first** - Designed to be used by AI assistants
4. **Portable** - Cross-platform (Linux, macOS, Windows)

## Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                          Gip CLI                                 │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                      main.rs                                ││
│  │  - Command parsing (clap)                                   ││
│  │  - Route to handlers                                        ││
│  └─────────────────────────────────────────────────────────────┘│
│                              │                                   │
│  ┌───────────┬───────────┬───────────┬───────────┬───────────┐ │
│  │  commit   │   init    │  context  │   push    │ passthru  │ │
│  │  handler  │  handler  │  handler  │  handler  │  handler  │ │
│  └───────────┴───────────┴───────────┴───────────┴───────────┘ │
│                              │                                   │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    Core Modules                              ││
│  │  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  ││
│  │  │   git.rs    │  │  manifest/  │  │      merge.rs       │  ││
│  │  │             │  │             │  │                     │  ││
│  │  └─────────────┘  └─────────────┘  └─────────────────────┘  ││
│  └─────────────────────────────────────────────────────────────┘│
│                              │                                   │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    External Dependencies                     ││
│  │  ┌─────────────┐  ┌─────────────┐                           ││
│  │  │ toon-format │  │    git2     │                           ││
│  │  │             │  │             │                           ││
│  │  └─────────────┘  └─────────────┘                           ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
                              │
                              ▼
┌─────────────────────────────────────────────────────────────────┐
│                         Git Repository                           │
│  ┌─────────────────────────────────────────────────────────────┐│
│  │                    refs/notes/gip                            ││
│  │  - Manifest storage                                          ││
│  │  - Travels with repo                                         ││
│  └─────────────────────────────────────────────────────────────┘│
└─────────────────────────────────────────────────────────────────┘
```

## Component Details

### CLI Layer (`src/main.rs`)

Entry point that:
- Parses command line arguments
- Routes to appropriate command handlers
- Provides help and version info

### Command Handlers (`src/commands/`)

| Handler | File | Responsibility |
|---------|------|----------------|
| `commit` | `commit.rs` | Intercepts commits, enforces manifests |
| `init` | `init.rs` | Initializes repo with AI instructions |
| `context` | `context.rs` | Queries semantic history |
| `push` | `push.rs` | Pushes code and notes |
| `passthrough` | `passthrough.rs` | Forwards unknown commands to git |

### Core Libraries

#### Git Integration (`src/git.rs`)

Abstracts git operations:
- `is_git_repo()` - Check if in a git repo
- `get_repo_root()` - Get repo root path
- `get_staged_diff()` - Get diff content
- `add_note()` / `get_note()` - Manage git notes
- `push_notes()` - Push code and notes
- `run_git_cmd()` - Run raw git commands

#### Manifest Module (`src/manifest/`)

Handles manifest operations:
- `load()` / `save()` - Storage in Git Notes
- `serialize_manifest_toon()` - TOON serialization
- `types.rs` - Data structures (Manifest, Entry, etc.)

#### Merge Driver (`src/merge.rs`)

Handles conflict enrichment:
- `enrich_all_conflicts()` - Detects and enriches conflicts
- `enrich_conflict_markers()` - Injects context into markers

### Data Storage

Manifests are stored in Git Notes under `refs/notes/gip`:
- Travels with repository
- Pushed/pulled with code
- Queryable via standard git

## Data Flow

### Commit Flow

```
1. User runs: gip commit -m "message with manifest"
2. commit handler parses message
3. ManifestParser extracts manifest block
4. If valid:
   a. Strip manifest from message
   b. Create commit with clean message
   c. Store manifest in notes
5. If missing/invalid:
   a. Reject commit
   b. Generate template
   c. Show error with instructions
```

### Context Query Flow

```
1. User runs: gip context src/file.rs
2. GitAdapter queries file history
3. For each commit, retrieve note
4. Parse manifests from notes
5. Apply filters (--behavior, --since)
6. Format output (terminal/JSON)
7. Optionally export to file
```

## Extension Points

### Adding New Commands

1. Create handler in `src/commands/`
2. Add to routing in `main.rs`
3. Update help text
4. Add tests

### Adding New Output Formats

1. Add serializer method to `src/manifest/toon.rs`
2. Add format option to context command
3. Update documentation

## Security Considerations

- No network calls (except git operations)
- No credential storage
- No code execution from manifests
- Manifests are text-only metadata
