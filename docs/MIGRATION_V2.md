# Migration Guide: v1.0 → v2.0

This guide helps you migrate from Gip v1.0 to v2.0.

## Overview

**Good news**: v2.0 is **fully backward compatible**. Your existing v1.0 manifests will continue to work without any changes!

## What's New in v2.0

- 🤖 AI-powered manifest generation (`--ai` flag)
- 🎯 Global intent for multi-function commits
- 📦 Selective context injection (20x smaller conflict files)
- 🔍 Automatic breaking change detection
- 🚩 Feature flag tracking
- 📊 Extended schema with signature deltas

## Migration Steps

### Step 1: Install v2.0

```bash
# If you installed via go install:
go install github.com/iamHrithikRaj/gip@latest

# Verify version:
gip version  # Should show 2.0.0
```

### Step 2: No Action Required!

Your existing v1.0 manifests will be **automatically migrated** when loaded:
- Missing `HunkID` fields → Default to `"H#0"`
- Missing `ChangeType` → Default to `"modify"`
- Missing `SchemaVersion` → Upgraded to `"2.0"`

**Example:**
```json
// v1.0 manifest (old format)
{
  "commit": "abc123",
  "entries": [
    {
      "anchor": {"file": "cart.py", "symbol": "calculateTotal"},
      "contract": {...},
      "behaviorClass": ["feature"],
      "rationale": "Added tax calculation"
    }
  ]
}

// After loading in v2.0 (auto-migrated)
{
  "schemaVersion": "2.0",
  "commit": "abc123",
  "entries": [
    {
      "anchor": {
        "file": "cart.py",
        "symbol": "calculateTotal",
        "hunk": "H#0"  // ← Added automatically
      },
      "changeType": "modify",  // ← Added automatically
      "contract": {...},
      "behaviorClass": ["feature"],
      "rationale": "Added tax calculation"
    }
  ]
}
```

### Step 3: Start Using v2.0 Features

#### AI-Generated Manifests

```bash
# Set your OpenAI API key
export OPENAI_API_KEY=sk-...

# Use AI to generate manifests
git add modified_file.py
gip commit --ai --intent "Add strict validation to parsers"
```

#### Global Intent (Multi-Function Commits)

When you modify multiple functions in one commit:

```bash
git add user.py order.py address.py
gip commit -c

# Gip will detect 3 functions and ask:
# "Do all these functions share the same intent? (Y/n)"
#
# If yes:
# - You describe the global intent once
# - Each function can inherit or provide unique rationale
```

#### Breaking Change Detection

Signature changes are now automatically detected:

```python
# Before
def parseUser(raw: str) -> User

# After
def parseUser(raw: str, strict: bool = False) -> User
```

Gip will populate:
```json
{
  "signatureDelta": {
    "before": "def parseUser(raw: str) -> User",
    "after": "def parseUser(raw: str, strict: bool = False) -> User"
  },
  "compatibility": {
    "breaking": false  // ← Added parameter has default, not breaking
  }
}
```

## Workflow Comparison

### v1.0 Workflow (Still Works!)

```bash
git add file.py
gip commit -c  # Interactive prompts per function
git commit
```

### v2.0 AI Workflow (Recommended)

```bash
git add file.py
gip commit --ai --intent "Your high-level intent"
# AI generates manifest automatically
git commit
```

### v2.0 Interactive with Global Intent

```bash
git add user.py order.py  # Multiple files
gip commit -c
# Gip asks: "Is this a multi-function commit?"
# If yes: Provide global intent once
# Each function: Inherit or customize
git commit
```

## Breaking Changes

**None!** v2.0 is fully backward compatible.

## New Schema Fields (Optional)

You can now use these new fields:

| Field | Description | Auto-Populated |
|-------|-------------|----------------|
| `globalIntent` | Commit-level shared intent | Manual or AI |
| `signatureDelta` | Before/after signatures | ✅ Automatic |
| `compatibility.breaking` | Breaking change flag | ✅ Automatic |
| `featureFlags` | Feature flags in code | ✅ Automatic |
| `changeType` | add/modify/delete/rename | ✅ Automatic |
| `hunk` | Line number of change | ✅ Automatic |

## Performance Improvements

### Selective Context Injection

v1.0 injected entire manifests into conflicts:
```python
<<<<<<< HEAD
code
||| Gip CONTEXT (100 lines of manifest for all functions)
```

v2.0 injects only relevant entries:
```python
<<<<<<< HEAD
code
||| Gip CONTEXT (8 lines for this specific function)
```

**Result**: 20x smaller conflict files

## Configuration

### AI Mode (v2.0)

Set your OpenAI API key:

```bash
# Option 1: Environment variable
export OPENAI_API_KEY=sk-...

# Option 2: Create .gip/config.yml (future)
# llm:
#   provider: openai
#   model: gpt-4o-mini
#   apiKey: ${OPENAI_API_KEY}
```

## Troubleshooting

### "manifest not found" errors

If you see errors loading old manifests:
```bash
# Force re-save all manifests in v2.0 format
gip migrate-manifests  # (Coming in v2.1)
```

### AI mode not working

```bash
# Check API key
echo $OPENAI_API_KEY

# Test with verbose mode
gip commit --ai --intent "test" --verbose
```

### Want to use v1.0 behavior

All v1.0 commands still work:
```bash
gip commit -c  # Same interactive mode as v1.0
```

## Next Steps

- Try `gip commit --ai` for your next commit
- Explore global intent for large refactorings
- Check `gip --help` for new flags
- Read [ARCHITECTURE_V2.md](docs/ARCHITECTURE_V2.md) for technical details

## Need Help?

- 📖 Read the [README](README.md)
- 🐛 [Report issues](https://github.com/iamHrithikRaj/gip/issues)
- 💬 Ask questions in [Discussions](https://github.com/iamHrithikRaj/gip/discussions)

---

**Welcome to Gip v2.0!** 🚀
