# Gip Test Suite

This directory contains the comprehensive test suite for Gip.

## Structure

```
tests/
├── integration/          # End-to-end integration tests
│   └── merge_test.go    # Merge workflow tests
├── fixtures/            # Test data and example files
│   ├── manifests/       # Sample manifest JSON files
│   │   ├── feature.json
│   │   ├── bugfix.json
│   │   └── refactor.json
│   └── conflicts/       # Sample conflict files
│       ├── simple.py
│       └── multi-hunk.go
└── helpers/             # Test utility functions
    ├── git.go          # Git repo setup, commit, merge helpers
    └── manifest.go     # Manifest creation helpers
```

## Running Tests

### Unit Tests

Unit tests are located alongside the code in `internal/*/` directories:

```bash
# Run all unit tests
go test ./internal/...

# Run with coverage
go test -cover ./internal/...

# Run specific package
go test ./internal/manifest/...
go test ./internal/merge/...
go test ./internal/diff/...
```

### Integration Tests

Integration tests are in `tests/integration/`:

```bash
# Run all integration tests
go test ./tests/integration/...

# Run with verbose output
go test -v ./tests/integration/...

# Skip integration tests (use -short flag)
go test -short ./...
```

### Using Makefile

```bash
# Run all unit tests
make test

# Run tests in short mode
make test-unit

# Run integration tests
make test-integration

# Generate coverage report
make coverage
```

## Test Helpers

### Git Helpers (`tests/helpers/git.go`)

```go
import "github.com/iamHrithikRaj/gip/tests/helpers"

func TestExample(t *testing.T) {
    // Setup test repository
    repoDir := helpers.SetupGitRepo(t)
    
    // Create and commit a file
    sha := helpers.CreateCommit(t, "file.txt", "content", "Initial commit")
    
    // Create a branch
    helpers.CreateBranch(t, "feature")
    
    // Make changes and merge
    helpers.CheckoutBranch(t, "main")
    success, err := helpers.MergeBranch(t, "feature")
    
    // Check for conflicts
    files := helpers.GetConflictedFiles(t)
    
    // Assertions
    helpers.AssertFileContains(t, "file.txt", "expected text")
    helpers.AssertNoError(t, err, "Operation should succeed")
}
```

### Manifest Helpers (`tests/helpers/manifest.go`)

```go
import (
    "github.com/iamHrithikRaj/gip/internal/manifest"
    "github.com/iamHrithikRaj/gip/tests/helpers"
)

func TestManifest(t *testing.T) {
    // Create a basic test manifest
    m := helpers.NewTestManifest("abc123", "file.go", "MyFunc", "Added feature")
    
    // Create a bugfix manifest
    m := helpers.NewBugfixManifest("def456", "file.go", "MyFunc")
    
    // Create a feature manifest with side effects
    m := helpers.NewFeatureManifest("ghi789", "api.go", "CallAPI", []string{"network:http"})
    
    // Create manifest with multiple entries
    specs := []helpers.ManifestEntrySpec{
        helpers.DefaultEntrySpec("file1.go", "Func1"),
        helpers.DefaultEntrySpec("file2.go", "Func2"),
    }
    m := helpers.NewManifestWithMultipleEntries("commit123", specs)
    
    // Save and load
    manifest.Save(m, m.Commit)
    loaded, err := manifest.Load(m.Commit)
}
```

## Test Fixtures

### Manifest Fixtures

Example manifest files for testing:

- `fixtures/manifests/feature.json` - Feature addition
- `fixtures/manifests/bugfix.json` - Bug fix
- `fixtures/manifests/refactor.json` - Refactoring with breaking changes

Load in tests:

```go
content, _ := os.ReadFile("tests/fixtures/manifests/feature.json")
var m manifest.Manifest
json.Unmarshal(content, &m)
```

### Conflict Fixtures

Example files with Git conflict markers:

- `fixtures/conflicts/simple.py` - Single conflict block
- `fixtures/conflicts/multi-hunk.go` - Multiple conflict blocks

## Writing Tests

### Unit Test Template

```go
package mypackage

import "testing"

func TestMyFunction(t *testing.T) {
    // Setup
    input := "test"
    
    // Execute
    result := MyFunction(input)
    
    // Verify
    if result != "expected" {
        t.Errorf("Expected 'expected', got '%s'", result)
    }
}
```

### Integration Test Template

```go
package integration

import (
    "testing"
    "github.com/iamHrithikRaj/gip/tests/helpers"
)

func TestWorkflow(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipping integration test in short mode")
    }
    
    // Setup test repo
    repoDir := helpers.SetupGitRepo(t)
    
    // Test workflow
    // ...
    
    // Assertions
    helpers.AssertFileContains(t, "file.txt", "expected")
}
```

## Coverage Goals

| Package | Current | Target |
|---------|---------|--------|
| `manifest` | ~90% | 90%+ |
| `merge` | ~75% | 85%+ |
| `diff` | ~70% | 80%+ |
| `git` | 0% | 60%+ |
| `prompt` | 0% | 50%+ |

## CI/CD

Tests run automatically via GitHub Actions on:
- Push to `main` or `develop`
- Pull requests to `main`

Test matrix:
- **OS**: Ubuntu, Windows, macOS
- **Go versions**: 1.21, 1.22, 1.23

See `.github/workflows/test.yml` for configuration.

## Manual Testing

For manual/exploratory testing, see PowerShell scripts in `scripts/`:

```powershell
# Test full merge workflow
.\scripts\test-gip-merge.ps1

# Test merge driver directly
.\scripts\test-merge-driver.ps1

# Test simple merge scenario
.\scripts\test-merge-simple.ps1
```

## Troubleshooting

### Tests fail with "not a git repository"

Some tests require a Git repository. Use `helpers.SetupGitRepo(t)` at the start of your test.

### Integration tests take too long

Use `-short` flag to skip integration tests:
```bash
go test -short ./...
```

### Coverage report not generating

Ensure you have `go tool cover` available:
```bash
go tool cover -html=coverage.out -o coverage.html
```

## Contributing

When adding new features:

1. ✅ Write unit tests first (TDD)
2. ✅ Add integration tests for workflows
3. ✅ Update fixtures if needed
4. ✅ Run `make test-all` before committing
5. ✅ Ensure coverage doesn't decrease

See [TEST_STRUCTURE_PROPOSAL.md](../TEST_STRUCTURE_PROPOSAL.md) for architectural details.
