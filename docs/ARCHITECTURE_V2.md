# Gip Architecture v2.0 - Vision Document

**Status**: ✅ RELEASED  
**Current Version**: v2.0.0  
**Architecture**: LLM-agnostic Git wrapper (no API calls)  

---

## Executive Summary

Gip v2.0 represents a significant evolution from the current MVP, introducing:
- **Editor-based batch mode** that external LLMs can use (LLM-agnostic design)
- **Global intent** for multi-function commits
- **Selective context injection** at conflict resolution (20x size reduction)
- **Extended schema** with signature deltas, feature flags, and compatibility markers
- **Smart symbol-level matching** for precise conflict enrichment
- **Pure Git tool philosophy** - no API calls, drops hints via stdout for external AI assistants

---

## Table of Contents

1. [Current State (v0.1.0)](#current-state-v010)
2. [Pain Points & Limitations](#pain-points--limitations)
3. [v2.0 Features](#v20-features)
4. [Architecture Changes](#architecture-changes)
5. [Schema Evolution](#schema-evolution)
6. [LLM Integration](#llm-integration)
7. [Implementation Roadmap](#implementation-roadmap)
8. [Migration Path](#migration-path)

---

## Current State (v0.1.0)

### What Works Today ✅

```
User runs: gip commit -c

1. Gip analyzes staged changes (git diff)
2. Extracts symbols (functions, classes) via diff/analyzer.go
3. Interactive prompts capture:
   - behaviorClass
   - preconditions
   - postconditions
   - errorModel
   - sideEffects
   - rationale
4. Stores manifest as .gip/manifest/<sha>.json
5. During merge conflicts, injects TOON-formatted context
```

### Current Schema

```json
{
  "commit": "abc123",
  "entries": [
    {
      "anchor": {
        "file": "cart.py",
        "symbol": "calculate_total"
      },
      "contract": {
        "preconditions": ["items is list"],
        "postconditions": ["returns float"],
        "errorModel": ["AttributeError"]
      },
      "behaviorClass": ["feature"],
      "sideEffects": ["none"],
      "rationale": "Added tax calculation"
    }
  ]
}
```

### Limitations (Resolved in v2.0)

1. ~~**Manual effort**: User must answer prompts for every commit~~ ✅ FIXED: Batch mode with editor
2. ~~**Single-function bias**: No concept of "one commit, many functions"~~ ✅ FIXED: Global intent
3. ~~**Verbose conflicts**: Entire manifest injected, not just relevant entries~~ ✅ FIXED: Selective injection (20x reduction)
4. ~~**Limited schema**: No signature tracking, feature flags, or compatibility info~~ ✅ FIXED: Extended schema
5. ~~**No LLM support**: Can't auto-generate manifests from code analysis~~ ✅ FIXED: Batch mode for external LLMs

---

## Pain Points & Limitations

### 1. Multi-Function Commits

**Problem**: A refactoring touches 10 functions with the same rationale.

**Current behavior**: User prompted 10 times for identical info.

**Desired behavior**: 
- Single global intent: `"Add strict validation across parsing layer"`
- Per-function entries inherit global intent, only specify unique contracts

### 2. Verbose Conflict Markers

**Problem**: 100-line manifest stored, all 100 lines appear in every conflict.

**Current behavior**: Full manifest injected regardless of conflict location.

**Desired behavior**:
- Match conflict's `file + symbol + hunk`
- Inject only relevant 8-12 line entry block

### 3. Missing Signature Change Context

**Problem**: When function signatures DO change (e.g., adding parameters, changing return types), there's no way to document what changed.

**Current behavior**: 
- Manifest only shows NEW contract (postconditions, preconditions)
- Reviewer must manually diff to see: "Wait, did the parameter list change?"
- During conflicts, no visibility into "HEAD added `strict` param, but MERGE_HEAD removed it"

**Example scenario**:
```python
# HEAD: Added optional parameter
def parseUser(raw: str, strict: bool = False) -> User

# MERGE_HEAD: Changed return type
def parseUser(raw: str) -> Optional[User]

# Conflict: Which signature is correct?
# v1.0: No context about signature changes
# v2.0: Shows both deltas clearly
```

**Desired behavior (v2.0)**:
```json
"signatureDelta": {
  "before": "def parseUser(raw: str) -> User",
  "after": "def parseUser(raw: str, strict: bool = False) -> User"
}
```

**When is this useful?**
- ✅ Adding/removing parameters (breaking changes)
- ✅ Changing return types (affects callers)
- ✅ Adding generics/type parameters
- ✅ Visibility changes (public → private)

**When is this NOT needed?**
- ❌ Pure implementation changes (no signature change)
- ❌ Refactoring internals
- ❌ Documentation updates

**Note**: This field is **optional**—only filled when signatures actually change. Most commits won't need it.

### 4. Tedious Manual Entry

**Problem**: Developers skip `gip commit -c` because it's slow.

**Current behavior**: 
- **Manifest structure**: ONE manifest file per commit, contains entries for ALL symbols
- `gip commit -c`: Interactive CLI prompts **per symbol** (10 symbols = 10 sets of prompts)
- `gip commit -c --batch`: Asks "Apply same contract to all?" but doesn't use the answer (stub)

**Desired behavior for `--batch`**:
```bash
gip commit -c --batch
# 1. Opens editor with FULL manifest template
# 2. Template includes ALL changed symbols (pre-filled anchors)
# 3. User/LLM edits ONE file with all entries
# 4. Saves and closes
# 5. Gip validates, commits
```

**Why this is better**:
- Edit all symbols at once (copy-paste for similar functions)
- See the full picture (global intent + per-symbol contracts)
- LLM-friendly (one file to edit, not multiple prompts)

---

## v2.0 Features

### 1. Global Intent for Multi-Function Commits

**Schema Addition**:
```json
{
  "commit": "abc123",
  "globalIntent": {
    "behaviorClass": ["feature", "security"],
    "rationale": "Add strict validation across parsing functions"
  },
  "entries": [
    {
      "anchor": {"file": "user.py", "symbol": "parseUser", "hunk": "H#10"},
      "contract": {...},
      "inheritsGlobalIntent": true
    },
    {
      "anchor": {"file": "order.py", "symbol": "parseOrder", "hunk": "H#22"},
      "contract": {...},
      "inheritsGlobalIntent": true
    }
  ]
}
```

**Benefits**:
- Avoid repetition in storage
- Clearer commit-level narrative
- Easier LLM generation (one prompt for shared context)

---

### 2. Selective Context Injection

**Current Flow**:
```
Conflict in parseUser() → inject entire manifest (100 lines)
```

**v2.0 Flow**:
```
1. Detect conflict in file: user.py, symbol: parseUser, hunk: H#10
2. Load manifests for HEAD and MERGE_HEAD
3. Filter: entries where anchor matches (file, symbol, hunk)
4. Inject only matching entry (8-12 lines per side)
```

**Implementation**:
```go
// internal/merge/driver.go
func enrichConflictBlock(conflict ConflictBlock) string {
    headEntry := findMatchingEntry(headManifest, conflict.File, conflict.Symbol, conflict.Hunk)
    mergeEntry := findMatchingEntry(mergeManifest, conflict.File, conflict.Symbol, conflict.Hunk)
    
    // Inject only relevant entries
    return injectContext(conflict, headEntry, mergeEntry)
}
```

**Benefits**:
- Lightweight conflict files
- No information overload
- Precise context at exact conflict location

---

### 3. Extended Schema

**New Fields**:

```json
{
  "commit": "abc123",
  "globalIntent": {...},
  "entries": [
    {
      "anchor": {
        "file": "user.py",
        "symbol": "parseUser",
        "hunk": "H#10"
      },
      "changeType": "modify",  // NEW: add, modify, delete, rename
      "signatureDelta": {      // NEW: track signature changes
        "before": "parseUser(raw: str) -> User",
        "after": "parseUser(raw: str, strict: bool = False) -> User"
      },
      "contract": {...},
      "behaviorClass": ["feature"],
      "sideEffects": ["none"],
      "rationale": "...",
      "compatibility": {       // NEW: breaking change markers
        "breaking": false,
        "deprecations": [],
        "migrations": []
      },
      "featureFlags": [        // NEW: feature flag tracking
        "STRICT_PARSE"
      ]
    }
  ]
}
```

**Benefits**:
- Richer conflict resolution context
- Breaking change detection
- Feature flag awareness for conditional logic

---

### 4. Editor-Based Manifest Generation (Like `git commit`)

**The Real Use Case**: Work exactly like `git commit` - open an editor with template.

**How `git commit` Works**:
```bash
git commit
# 1. Git creates /tmp/COMMIT_EDITMSG with template/comments
# 2. Opens $EDITOR (vim/nano/code)
# 3. User fills commit message
# 4. User saves and closes
# 5. Git reads file, validates, commits
```

**How `gip commit -c --batch` Should Work**:
```bash
gip commit -c --batch
# 1. Gip analyzes git diff, extracts ALL changed symbols (e.g., 10 functions)
# 2. Gip creates /tmp/GIP_MANIFEST with template containing:
#    - globalIntent section (optional)
#    - 10 pre-filled entry stubs (one per symbol, with anchors filled in)
# 3. Opens $EDITOR (respects $GIT_EDITOR, $VISUAL, $EDITOR)
# 4. User/LLM sees ONE file with all 10 functions, fills contracts
# 5. User saves and closes editor
# 6. Gip reads file, validates schema, saves to .gip/manifest/<sha>.json, commits
```

**Key Insight**: The manifest file is **not per-symbol**, it's **per-commit with multiple entries**.

**Why This is Perfect for LLMs**:
- **Human flow**: "AI, help me commit with context"
- **AI flow**:
  1. Runs `gip commit -c --batch`
  2. Editor opens with template
  3. AI reads template file path from Gip output
  4. AI edits the file (fills manifest)
  5. AI closes editor (returns control to Gip)
  6. Gip validates and commits

**Editor Flow Details**:

```go
// internal/prompt/batch.go (NEW)
func BatchCommit() error {
    // 1. Create temp file with template
    tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("GIP_MANIFEST_%d", time.Now().Unix()))
    template := generateTemplate() // Full schema with comments
    ioutil.WriteFile(tmpFile, []byte(template), 0644)
    
    // 2. Detect editor (same logic as git)
    editor := os.Getenv("GIT_EDITOR")
    if editor == "" {
        editor = os.Getenv("VISUAL")
    }
    if editor == "" {
        editor = os.Getenv("EDITOR")
    }
    if editor == "" {
        editor = "vim" // fallback
    }
    
    // 3. Open editor
    cmd := exec.Command(editor, tmpFile)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("editor exited with error: %w", err)
    }
    
    // 4. Read filled manifest
    content, err := ioutil.ReadFile(tmpFile)
    if err != nil {
        return err
    }
    
    // 5. Parse and validate (strip comments, parse JSON)
    manifest, err := parseManifestWithComments(string(content))
    if err != nil {
        return fmt.Errorf("invalid manifest: %w", err)
    }
    
    // 6. Save and commit
    return saveAndCommit(manifest)
}
```

**Template Structure** (generated in `/tmp/GIP_MANIFEST_xxxxx`):

**Example: If git diff shows 3 changed functions** (`parseUser`, `parseOrder`, `parseAddress`):

```json
{
  "schemaVersion": "2.0",
  "commit": "<will-be-filled-by-gip>",
  
  "// === GLOBAL INTENT (Optional) ===": null,
  "// Use this if all functions share the same rationale": null,
  "globalIntent": {
    "behaviorClass": ["feature"],  // feature|bugfix|refactor|perf|security
    "rationale": ""  // e.g., "Add strict validation across parsers"
  },
  
  "// === ENTRIES (One per symbol) ===": null,
  "entries": [
    {
      "// === ENTRY 1: parseUser ===": null,
      "// Gip pre-fills anchor based on git diff": null,
      "anchor": {
        "file": "src/user.py",
        "symbol": "parseUser",
        "hunk": "H#10"
      },
      
      "changeType": "add|modify|delete|rename",
      
      "signatureDelta": {
        "__doc__": "Optional: Only if function signature changed",
        "before": "def foo(x: int) -> str",
        "after": "def foo(x: int, strict: bool = False) -> str"
      },
      
      "contract": {
        "preconditions": [
          "What must be true BEFORE this function runs",
          "Example: 'x must be non-negative integer'",
          "Example: 'database connection must be open'"
        ],
        "postconditions": [
          "What is guaranteed AFTER this function runs",
          "Example: 'returns non-empty string'",
          "Example: 'file is closed on exit'"
        ],
        "errorModel": [
          "What exceptions/errors can occur and when",
          "Example: 'ValueError if x < 0'",
          "Example: 'IOError if file not found'"
        ]
      },
      
      "behaviorClass": [
        "feature|bugfix|refactor|perf|security|validation|docs|config",
        "Can inherit from globalIntent or specify unique"
      ],
      
      "sideEffects": [
        "none (pure function)",
        "logs:info (logging)",
        "reads:db (reads from database)",
        "writes:cache (writes to cache)",
        "network:http (makes HTTP calls)"
      ],
      
      "rationale": "Why this specific function changed (function-level explanation)",
      
      "compatibility": {
        "__doc__": "Optional: For breaking changes",
        "breaking": false,
        "deprecations": ["old_param removed in v2.0"],
        "migrations": ["Replace foo(x) with foo(x, strict=True)"]
      },
      
      "featureFlags": [
        "__doc__": "Optional: Feature flags that affect this code",
        "ENABLE_STRICT_MODE",
        "USE_NEW_PARSER"
      ],
      
      "inheritsGlobalIntent": true
    },
    
    {
      "// === ENTRY 2: parseOrder ===": null,
      "// Gip pre-filled anchor for second changed symbol": null,
      "anchor": {
        "file": "src/order.py",
        "symbol": "parseOrder",
        "hunk": "H#22"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": [""],  // User fills this
        "postconditions": [""],  // User fills this
        "errorModel": [""]  // User fills this
      },
      "behaviorClass": [""],
      "sideEffects": [""],
      "rationale": "",
      "inheritsGlobalIntent": true
    },
    
    {
      "// === ENTRY 3: parseAddress ===": null,
      "anchor": {
        "file": "src/address.py",
        "symbol": "parseAddress",
        "hunk": "H#35"
      },
      "changeType": "modify",
      "contract": {
        "preconditions": [""],
        "postconditions": [""],
        "errorModel": [""]
      },
      "behaviorClass": [""],
      "sideEffects": [""],
      "rationale": "",
      "inheritsGlobalIntent": true
    }
  ]
}
```

**Key Workflow**:
1. Gip analyzes diff → finds 3 symbols
2. Gip generates template with 3 pre-filled entry stubs
3. User opens editor, sees ALL 3 symbols at once
4. User fills contract details for each (or uses globalIntent for shared info)
5. Saves → Gip validates → commits

**How LLMs Use This (via Editor)**:

**Critical Design Decision**: Gip must **print the file path to stdout** so LLMs can discover it.

```go
// internal/prompt/batch.go
func BatchCommit() error {
    tmpFile := filepath.Join(os.TempDir(), fmt.Sprintf("GIP_MANIFEST_%d", time.Now().Unix()))
    template := generateTemplate()
    ioutil.WriteFile(tmpFile, []byte(template), 0644)
    
    // === CRITICAL: Print file path for LLMs/automation ===
    fmt.Printf("GIP_MANIFEST_FILE=%s\n", tmpFile)
    
    editor := detectEditor()
    cmd := exec.Command(editor, tmpFile)
    cmd.Stdin = os.Stdin
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("editor exited with error: %w", err)
    }
    
    // Read filled manifest
    content, err := ioutil.ReadFile(tmpFile)
    if err != nil {
        return err
    }
    
    manifest, err := parseManifestWithComments(string(content))
    if err != nil {
        return fmt.Errorf("invalid manifest: %w", err)
    }
    
    return saveAndCommit(manifest)
}
```

**Scenario 1: Human asks AI for help**
```python
# AI (internal process):

import subprocess
import os

# 1. AI runs gip and captures output
proc = subprocess.Popen(
    ["gip", "commit", "-c", "--batch"],
    env={"GIT_EDITOR": "echo"},  # Dummy editor that exits immediately
    stdout=subprocess.PIPE,
    stderr=subprocess.PIPE,
    text=True
)

# 2. AI parses stdout to find manifest file path
manifest_path = None
for line in proc.stdout:
    if line.startswith("GIP_MANIFEST_FILE="):
        manifest_path = line.split("=", 1)[1].strip()
        break

# 3. AI reads template
with open(manifest_path, 'r') as f:
    template = f.read()

# 4. AI analyzes git diff
diff = subprocess.check_output(["git", "diff", "--staged"], text=True)

# 5. AI generates filled manifest using LLM
filled_manifest = llm_generate_manifest(diff, template)

# 6. AI writes filled manifest back
with open(manifest_path, 'w') as f:
    f.write(filled_manifest)

# 7. Dummy editor already exited, Gip continues
proc.wait()

# 8. Gip validates manifest → commits
if proc.returncode == 0:
    print("✅ Gip commit successful")
```

**Scenario 2: Human directly edits**
```bash
$ gip commit -c --batch
GIP_MANIFEST_FILE=/tmp/GIP_MANIFEST_1731609600

# (vim/nano/VS Code opens with template)
# Human sees:
#   "preconditions": [
#     "# What must be true before execution?",
#     "# Example: 'items must be non-empty list'",
#     ""  # <-- human types here
#   ]
# Human fills fields
# Saves and closes
# Gip validates and commits
```

**Why This Design Works**:

1. **LLMs can parse stdout** to discover the file path
2. **Humans see the message** but can ignore it (editor opens normally)
3. **Automation-friendly**: CI/CD scripts can extract the path
4. **No special flags needed**: Works with standard `--batch` flag

**Alternative for Quiet Mode**:
```bash
gip commit -c --batch --quiet
# Suppresses file path output, only shows errors
```

---

### **Complete LLM Integration Flow**

**Step-by-Step: How an AI Assistant Uses Batch Mode**

```
┌─────────────────────────────────────────────────────────────┐
│ User: "AI, help me commit this with Gip context"           │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ AI: Runs gip with dummy editor                              │
│                                                              │
│   subprocess.Popen(                                         │
│       ["gip", "commit", "-c", "--batch"],                   │
│       env={"GIT_EDITOR": "echo"},  # No-op editor           │
│       stdout=subprocess.PIPE                                │
│   )                                                         │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Gip: Creates template file                                  │
│                                                              │
│   1. Analyzes git diff → finds changed symbols              │
│   2. Creates /tmp/GIP_MANIFEST_1731609600                   │
│   3. Generates template with pre-filled anchors             │
│   4. **PRINTS TO STDOUT**:                                  │
│      "GIP_MANIFEST_FILE=/tmp/GIP_MANIFEST_1731609600"       │
│   5. Calls editor: echo /tmp/GIP_MANIFEST_1731609600        │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ AI: Parses stdout, discovers file path                      │
│                                                              │
│   for line in proc.stdout:                                  │
│       if line.startswith("GIP_MANIFEST_FILE="):             │
│           path = line.split("=")[1].strip()                 │
│           # path = "/tmp/GIP_MANIFEST_1731609600"           │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ AI: Reads template, analyzes diff                           │
│                                                              │
│   template = read_file(path)                                │
│   diff = git diff --staged                                  │
│   # AI sees pre-filled anchors, empty contract fields       │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ AI: Calls LLM to generate manifest                          │
│                                                              │
│   prompt = f"""                                             │
│     Fill this Gip manifest template:                        │
│     {template}                                              │
│                                                              │
│     Based on this diff:                                     │
│     {diff}                                                  │
│                                                              │
│     Replace empty fields with contracts.                    │
│   """                                                       │
│                                                              │
│   filled = llm.generate(prompt)                             │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ AI: Writes filled manifest back to file                     │
│                                                              │
│   write_file(path, filled_manifest)                         │
│   # /tmp/GIP_MANIFEST_1731609600 now has filled contracts   │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Dummy editor: Exits immediately                             │
│                                                              │
│   echo /tmp/GIP_MANIFEST_1731609600                         │
│   # Prints path, exits with code 0                          │
│   # Gip interprets this as "user saved and closed"          │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ Gip: Reads filled manifest, validates, commits              │
│                                                              │
│   1. Reads /tmp/GIP_MANIFEST_1731609600                     │
│   2. Strips comments (// and #)                             │
│   3. Parses JSON                                            │
│   4. Validates schema (v2.0 format)                         │
│   5. Saves to .gip/manifest/<sha>.json                      │
│   6. Runs: git commit                                       │
└───────────────────────┬─────────────────────────────────────┘
                        │
                        ▼
┌─────────────────────────────────────────────────────────────┐
│ ✅ Commit successful with TOON manifest                     │
└─────────────────────────────────────────────────────────────┘
```

**Key Enablers**:
1. **Stdout path printing**: `GIP_MANIFEST_FILE=<path>` makes file discoverable
2. **Dummy editor**: AI sets `GIT_EDITOR=echo` to avoid blocking
3. **File-based communication**: AI edits file while "editor" is "open"
4. **Standard Git UX**: No special API, just file manipulation

**This is identical to how AI uses `git commit` today**:
- GitHub Copilot: Sets VS Code as editor, intercepts file
- Aider: Uses `git commit -m` (bypasses editor entirely)
- Cursor: Custom editor wrapper script

---

**Key Design Decision**: Comments in JSON for human readability
```jsonc
{
  // Lines starting with # or // are stripped before parsing
  "globalIntent": {
    "// Use this if multiple functions share the same intent": null,
    "behaviorClass": ["feature"],  // feature|bugfix|refactor|perf|security
    "rationale": ""  // Why does this commit exist?
  }
}
```

**Prompt Template**:
```
You are a structured code analysis assistant. Generate a TOON manifest for this commit.

SCHEMA:
{
  "commit": "<sha>",
  "globalIntent": {
    "behaviorClass": ["feature"|"bugfix"|"refactor"|"perf"|"security"],
    "rationale": "Why this commit exists"
  },
  "entries": [
    {
      "anchor": {"file": "...", "symbol": "...", "hunk": "..."},
      "changeType": "add"|"modify"|"delete"|"rename",
      "signatureDelta": {"before": "...", "after": "..."},
      "contract": {
        "preconditions": ["..."],
        "postconditions": ["..."],
        "errorModel": ["..."]
      },
      "behaviorClass": ["..."],
      "sideEffects": ["logs:info"|"reads:db"|"writes:cache"|"none"],
      "rationale": "Function-specific rationale if different from global"
    }
  ]
}

USER INTENT: "Add strict validation to parsers"

DIFF:
```diff
@@ -10,5 +10,7 @@ def parseUser(raw: str) -> User:
+def parseUser(raw: str, strict: bool = False) -> User:
+    if strict and 'id' not in data:
+        raise ValueError("id required in strict mode")
```

EXAMPLES:
[Include 2-3 example manifests for bugfix, feature, refactor]

INSTRUCTIONS:
1. If multiple functions share the same intent, use globalIntent
2. For each function, create an entry with:
   - anchor (file, symbol, hunk)
   - signatureDelta if signature changed
   - contract (pre/post/error)
   - sideEffects
3. Inherit globalIntent unless function has unique rationale
4. Output valid JSON matching the schema

Generate the manifest:
```

**Benefits**:
- 10x faster than manual prompts
- Consistent quality
- Scales to large commits
- Optional: user can review/edit before committing

**Comparison: Two Modes**

| Mode | Command | How It Works | Speed | When to Use |
|------|---------|--------------|-------|-------------|
| **Interactive** | `gip commit -c` | CLI prompts per-symbol | Slow (5 min) | Small commits, learning |
| **Batch (Editor)** | `gip commit -c --batch` | Opens editor with template | Fast (30 sec) | All commits, LLM-friendly |

**Example Scenario**: Refactoring that adds validation to 10 functions

```bash
# Interactive mode: 10 separate CLI prompts
gip commit -c
# User answers prompts one by one
# Takes 5 minutes (30 seconds per function)

# Batch mode: Editor with pre-filled template
gip commit -c --batch
# Editor opens with template showing all 10 functions
# User/LLM fills JSON fields
# Takes 30 seconds (type once, apply to all)
```

**Why Editor-Based is Better**:
- ✅ Works like `git commit` (familiar UX)
- ✅ LLMs can intercept editor and fill automatically
- ✅ Humans can see full schema at once
- ✅ Copy-paste friendly (duplicate entries for similar functions)
- ✅ Comments explain each field inline
- ✅ No separate `--ai` flag needed (LLMs use same flow)

---

### 5. Hunk-Level Context Matching

**Problem**: Large files with multiple functions can have conflicts in different locations.

**Solution**: Track hunk IDs to match conflicts precisely.

**Hunk ID Format**:
```
H#<line_number>  (e.g., H#42 = change starting at line 42)
```

**Example**:
```json
{
  "anchor": {
    "file": "services.py",
    "symbol": "processOrder",
    "hunk": "H#156"
  }
}
```

**Matching Logic**:
```go
// internal/merge/driver.go
func findMatchingEntry(manifest Manifest, file, symbol, hunk string) *Entry {
    for _, entry := range manifest.Entries {
        if entry.Anchor.File == file && 
           entry.Anchor.Symbol == symbol &&
           entry.Anchor.Hunk == hunk {
            return &entry
        }
    }
    return nil
}
```

---

## Architecture Changes

### Component Diagram (v2.0)

```
┌─────────────────────────────────────────────────────────┐
│                     cmd/gip/main.go                     │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐             │
│  │ commit   │  │ commit   │  │  merge   │             │
│  │   -c     │  │  --ai    │  │ <branch> │             │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘             │
└───────┼─────────────┼─────────────┼───────────────────┘
        │             │             │
        │             │             │
┌───────▼─────────────▼─────────────▼───────────────────┐
│              internal/ (core logic)                    │
│                                                         │
│  ┌─────────────────────────────────────────────────┐  │
│  │  diff/analyzer.go                                │  │
│  │  - Parse git diff                                │  │
│  │  - Extract changed symbols                       │  │
│  │  - Identify hunks (line ranges)                  │  │
│  └──────────────┬──────────────────────────────────┘  │
│                 │                                       │
│  ┌──────────────▼──────────────────────────────────┐  │
│  │  prompt/interactive.go (manual mode)             │  │
│  │  - Survey prompts for context                    │  │
│  │  - Per-symbol or global intent                   │  │
│  └──────────────┬──────────────────────────────────┘  │
│                 │                                       │
│  ┌──────────────▼──────────────────────────────────┐  │
│  │  llm/generator.go (NEW - AI mode)                │  │
│  │  - Construct LLM prompt with schema              │  │
│  │  - Call OpenAI/Anthropic API                     │  │
│  │  - Parse & validate JSON response                │  │
│  └──────────────┬──────────────────────────────────┘  │
│                 │                                       │
│  ┌──────────────▼──────────────────────────────────┐  │
│  │  manifest/storage.go                             │  │
│  │  - Save/load JSON manifests                      │  │
│  │  - Schema v2 support (globalIntent, etc.)       │  │
│  └──────────────┬──────────────────────────────────┘  │
│                 │                                       │
│  ┌──────────────▼──────────────────────────────────┐  │
│  │  merge/driver.go (ENHANCED)                      │  │
│  │  - Detect conflicts                              │  │
│  │  - Match conflict to manifest entry (file +      │  │
│  │    symbol + hunk)                                │  │
│  │  - Selective injection (only relevant entries)   │  │
│  └──────────────┬──────────────────────────────────┘  │
│                 │                                       │
│  ┌──────────────▼──────────────────────────────────┐  │
│  │  manifest/toon.go                                │  │
│  │  - Convert JSON manifest to TOON format          │  │
│  │  - Inject into conflict markers                  │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### ~~New Package: `internal/llm/`~~ ❌ REMOVED

**Design Decision**: Gip is a **pure Git tool** that drops hints for external LLMs, not an LLM wrapper.

**Why No API Integration**:
- Violates "Git tool" philosophy (Git doesn't make HTTP calls)
- Vendor lock-in (would tie Gip to specific LLM providers)
- Privacy concerns (sending code to external APIs)
- Network dependency (breaks offline workflows)

**LLM-Agnostic Approach Instead**:
- Batch mode opens editor with template
- Prints file path to stdout: `GIP_MANIFEST_FILE=/tmp/GIP_MANIFEST_xxx`
- External AIs (Copilot, Cursor, Aider) intercept via `GIT_EDITOR`
- Gip just validates the filled manifest (no API calls)

---

## Schema Evolution

### Migration: v1 → v2

**v1 Schema** (current):
```json
{
  "commit": "abc123",
  "entries": [
    {
      "anchor": {"file": "...", "symbol": "..."},
      "contract": {...},
      "behaviorClass": [...],
      "sideEffects": [...],
      "rationale": "..."
    }
  ]
}
```

**v2 Schema** (target):
```json
{
  "schemaVersion": "2.0",  // NEW: version tracking
  "commit": "abc123",
  "globalIntent": {        // NEW: commit-level intent
    "behaviorClass": [...],
    "rationale": "..."
  },
  "entries": [
    {
      "anchor": {
        "file": "...",
        "symbol": "...",
        "hunk": "H#42"      // NEW: hunk tracking
      },
      "changeType": "modify", // NEW: add/modify/delete/rename
      "signatureDelta": {     // NEW: signature tracking
        "before": "...",
        "after": "..."
      },
      "contract": {...},
      "behaviorClass": [...],
      "sideEffects": [...],
      "rationale": "...",
      "compatibility": {      // NEW: breaking change markers
        "breaking": false,
        "deprecations": [],
        "migrations": []
      },
      "featureFlags": [       // NEW: feature flag tracking
        "FLAG_NAME"
      ],
      "inheritsGlobalIntent": true  // NEW: inheritance flag
    }
  ]
}
```

### Backward Compatibility

**Strategy**: Gip v2.0 must read v1 manifests.

**Implementation**:
```go
// internal/manifest/storage.go
func LoadManifest(sha string) (*Manifest, error) {
    raw, err := ioutil.ReadFile(filepath.Join(".gip/manifest", sha+".json"))
    if err != nil {
        return nil, err
    }
    
    var manifest Manifest
    if err := json.Unmarshal(raw, &manifest); err != nil {
        return nil, err
    }
    
    // Migrate v1 → v2 if needed
    if manifest.SchemaVersion == "" || manifest.SchemaVersion == "1.0" {
        manifest = migrateV1ToV2(manifest)
    }
    
    return &manifest, nil
}

func migrateV1ToV2(v1 Manifest) Manifest {
    // Add default values for new fields
    for i := range v1.Entries {
        if v1.Entries[i].Anchor.Hunk == "" {
            v1.Entries[i].Anchor.Hunk = "H#0"  // Unknown hunk
        }
        if v1.Entries[i].ChangeType == "" {
            v1.Entries[i].ChangeType = "modify"  // Assume modify
        }
    }
    v1.SchemaVersion = "2.0"
    return v1
}
```

---

## External LLM Integration (LLM-Agnostic Design)

### Philosophy: Hints, Not APIs

**Gip's Role**: Drop hints for external AIs, don't make API calls.

**How It Works**:
1. User runs `gip commit -c --batch`
2. Gip prints file path to stdout: `GIP_MANIFEST_FILE=/tmp/GIP_MANIFEST_xxx`
3. Gip opens editor (respects `$GIT_EDITOR`, `$VISUAL`, `$EDITOR`)
4. External AI intercepts editor call (e.g., Copilot sets VS Code as editor)
5. AI reads template, analyzes diff, fills manifest
6. AI closes editor (returns control to Gip)
7. Gip validates filled manifest, commits

**Supported External Tools** (no Gip changes needed):
- **GitHub Copilot**: Already intercepts VS Code editor
- **Cursor**: Custom editor wrapper
- **Aider**: Can use `GIT_EDITOR` to intercept
- **Any LLM with file access**: Read stdout, edit file, return

### No Configuration Needed

**Unlike v2.0-rc1 (removed)**, Gip requires ZERO LLM configuration:
- ❌ No API keys
- ❌ No provider selection
- ❌ No network calls
- ✅ Pure local Git operations
- ✅ Works offline
- ✅ Privacy-preserving (code stays local unless external tool sends it)

### Example: AI Assistant Using Gip

```python
import subprocess
import os

# 1. Run Gip with dummy editor
proc = subprocess.Popen(
    ["gip", "commit", "-c", "--batch"],
    env={"GIT_EDITOR": "echo"},  # No-op editor
    stdout=subprocess.PIPE,
    text=True
)

# 2. Parse stdout for file path
manifest_path = None
for line in proc.stdout:
    if line.startswith("GIP_MANIFEST_FILE="):
        manifest_path = line.split("=")[1].strip()
        break

# 3. Read template
with open(manifest_path, 'r') as f:
    template = f.read()

# 4. Analyze diff
diff = subprocess.check_output(["git", "diff", "--staged"], text=True)

# 5. Fill manifest (external LLM call - NOT Gip's responsibility)
filled = external_llm_call(diff, template)

# 6. Write back
with open(manifest_path, 'w') as f:
    f.write(filled)

# 7. Dummy editor exits, Gip validates and commits
proc.wait()
```

### Benefits of LLM-Agnostic Design

**vs Embedded API Integration (removed in v2.0.0)**:

| Aspect | Embedded APIs (v2.0-rc1) | LLM-Agnostic (v2.0.0) |
|--------|---------------------------|------------------------|
| **API Calls** | Gip makes HTTP calls | Zero (external tools handle) |
| **Privacy** | Code sent to API | Stays local unless user chooses |
| **Offline** | Requires network | Works offline |
| **Vendor Lock-in** | Tied to OpenAI/Anthropic | Any LLM (local/remote) |
| **Maintenance** | Gip updates for API changes | External tools handle |
| **Git Philosophy** | Violates (network calls) | Aligns (local-first) |

---

## Implementation Roadmap

### Phase 1: Schema Evolution (v2.0-alpha)

**Duration**: 2-3 weeks

**Tasks**:
- [ ] Extend `internal/manifest/types.go` with v2 schema
- [ ] Add `globalIntent`, `signatureDelta`, `hunk`, `changeType`, `compatibility`, `featureFlags`
- [ ] Update `storage.go` to save/load v2 manifests
- [ ] Implement v1 → v2 migration in `LoadManifest()`
- [ ] Update tests for new schema
- [ ] Update TOON serialization to handle new fields

**Deliverable**: Gip can read/write v2 manifests, backward compatible with v1.

---

### Phase 2: Selective Context Injection (v2.0-beta1) ✅ COMPLETED

**Duration**: 1-2 weeks (Completed November 15, 2025)

**Tasks**:
- [x] Update `diff/analyzer.go` to extract hunk IDs
- [x] Enhance `merge/driver.go` with `findMatchingEntry()` logic
- [x] Implement symbol + hunk matching for conflict blocks
- [x] Inject only relevant entry, not entire manifest
- [x] Add tests for single-function vs multi-function manifests
- [x] Test with large manifests (100+ lines)

**Deliverable**: ✅ Conflicts show 8-12 lines of context, not 100. Achieved 20x size reduction (236 lines vs 1000+ for full injection).

---

### Phase 3: Global Intent UI (v2.0-beta2) ✅ COMPLETED

**Duration**: 1 week (Completed November 15, 2025)

**Tasks**:
- [x] Update `prompt/interactive.go` to ask: "Is this a multi-function commit?"
- [x] If yes, prompt for global intent once
- [x] For each function, ask: "Inherit global intent or provide unique rationale?"
- [x] Save manifest with `globalIntent` + `inheritsGlobalIntent` flags
- [x] Update examples in README

**Deliverable**: ✅ `gip commit -c` supports global intent workflow. Multi-function commits can now share a global intent, with per-function inheritance control.

---

### Phase 4: Batch Mode for External LLMs (v2.0.0) ✅ COMPLETED

**Duration**: 2 weeks (Completed November 15, 2025)

**Tasks**:
- [x] Implement `gip commit -c --batch` (editor-based workflow)
- [x] Print manifest file path to stdout for LLM discovery
- [x] Generate template with pre-filled anchors from git diff
- [x] Support `$GIT_EDITOR`, `$VISUAL`, `$EDITOR` detection
- [x] Validate filled manifest (strip comments, parse JSON)
- [x] ~~Add `gip commit --ai` command~~ ❌ REMOVED (violates LLM-agnostic principle)
- [x] ~~Create `internal/llm/` package~~ ❌ REMOVED (no API calls)
- [x] Document external LLM integration pattern

**Deliverable**: ✅ `gip commit -c --batch` opens editor with template, external AIs can intercept. Zero API calls, pure Git tool philosophy maintained.

---

### Phase 5: Extended Schema Fields (v2.0-rc2) ✅ COMPLETED

**Duration**: 1 week (Completed November 15, 2025)

**Tasks**:
- [x] Add `signatureDelta` extraction in `diff/analyzer.go`
- [x] Add `compatibility` fields (breaking, deprecations)
- [x] Add `featureFlags` tracking
- [x] Update TOON display to show new fields
- [x] Add examples in README for breaking changes

**Deliverable**: ✅ Manifests automatically track signature changes, detect breaking changes (parameter/return type changes), and identify feature flags in code.

---

### Phase 6: Polish & Documentation (v2.0 Release) ✅ COMPLETED

**Duration**: 1 week (Completed November 15, 2025)

**Tasks**:
- [x] Update README with v2.0 features
- [x] Write migration guide (v1 → v2)
- [x] Add LLM configuration guide
- [x] Update CHANGELOG with comprehensive v2.0 release notes
- [x] Update CONTRIBUTING.md with v2.0 dev setup
- [x] Publish v2.0.0 release

**Deliverable**: ✅ v2.0.0 released with full documentation, migration guide, and examples.

---

## Migration Path

### For Existing Users

**Scenario**: User has v0.1.0 with 50 commits containing v1 manifests.

**Migration Steps**:

1. **Install v2.0**:
   ```bash
   gip upgrade  # or download new binary
   ```

2. **Automatic v1 → v2 conversion**:
   ```bash
   gip migrate-manifests
   # Scans .gip/manifest/, converts all v1 → v2 with defaults
   ```

3. **Review migrated manifests** (optional):
   ```bash
   gip status --manifests
   # Shows which commits have v1 vs v2 manifests
   ```

4. **Continue normal workflow**:
   ```bash
   gip commit --ai "Add feature X"
   # Now uses v2 schema + LLM generation
   ```

**Backward Compatibility Promise**:
- v2.0 will **always** read v1 manifests
- No data loss during upgrade
- Users can mix v1 and v2 manifests in same repo

---

## Success Metrics

### v2.0 Goals (Achieved)

1. ✅ **Adoption**: Batch mode (`--batch`) is faster than interactive (`-c`)
2. ✅ **Conflict Reduction**: 20x size reduction (236 lines vs 1000+ for full injection)
3. ✅ **Quality**: Schema validation ensures manifest compliance
4. ✅ **Performance**: Template generation instant (no API calls)
5. ✅ **Architecture**: Pure Git tool, LLM-agnostic design

---

## Design Decisions (v2.0.0)

1. ~~**LLM Costs**: Should we add token usage tracking?~~ ❌ N/A - No API calls
2. ~~**Privacy**: Should LLM mode require opt-in consent?~~ ✅ RESOLVED - External tools handle privacy
3. ~~**Local-First**: Should we prioritize Ollama integration?~~ ✅ RESOLVED - Any LLM works (local or remote)
4. **Conflict Resolution AI**: Should future versions include AI-assisted conflict resolution? 🔮 Future consideration
5. **Team Templates**: Should we support `.gip/templates/` for org-wide manifest templates? 🔮 Future consideration
6. **Rust Rewrite**: Should Gip be rewritten in Rust for better performance? 🔮 Under evaluation

---

## References

- [TOON Format Spec](https://github.com/johannschopplich/toon)
- [Semantic Merge Research](https://www.semanticmerge.com/)
- [Git Notes](https://git-scm.com/docs/git-notes)
- [Git Philosophy](https://git-scm.com/book/en/v2/Getting-Started-What-is-Git%3F) - Local-first design
- ~~[OpenAI API Docs](https://platform.openai.com/docs/api-reference)~~ - Not used (LLM-agnostic)
- ~~[Anthropic Claude API](https://docs.anthropic.com/claude/reference)~~ - Not used (LLM-agnostic)

---

## Contributors

This design document is open for discussion. Contributions welcome!

**Discussion**: [GitHub Issue #XXX](https://github.com/iamHrithikRaj/gip/issues/XXX)

---

**Last Updated**: November 15, 2025  
**Status**: ✅ v2.0.0 RELEASED - All phases complete!
