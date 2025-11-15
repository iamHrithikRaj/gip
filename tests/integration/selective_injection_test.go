package integration

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/iamHrithikRaj/gip/internal/manifest"
	"github.com/iamHrithikRaj/gip/internal/merge"
	"github.com/iamHrithikRaj/gip/tests/helpers"
)

// TestSelectiveInjection_SingleFunction tests that only the relevant manifest
// entry is injected for a single-function conflict
func TestSelectiveInjection_SingleFunction(t *testing.T) {
	// Setup git repo
	repo := helpers.CreateTestRepo(t)
	defer repo.Cleanup()

	// Create initial file with function
	initialCode := `def calculate(x):
    return x * 2
`
	repo.WriteFile("math.py", initialCode)
	repo.GitAdd("math.py")
	repo.GitCommit("Initial commit")

	// Branch A: Change to multiply by 3
	repo.GitCheckout("branch-a")
	branchACode := `def calculate(x):
    return x * 3
`
	repo.WriteFile("math.py", branchACode)
	repo.GitAdd("math.py")
	repo.GitCommit("Multiply by 3")

	// Create manifest for branch A
	branchASHA := repo.GetCurrentCommit()
	manifestA := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchASHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "math.py",
					Symbol: "calculate",
					HunkID: "H#1",
				},
				ChangeType: "modify",
				Contract: manifest.Contract{
					Preconditions:  []string{"x must be number"},
					Postconditions: []string{"returns x * 3"},
					ErrorModel:     []string{"none"},
				},
				BehaviorClass: []string{"feature"},
				Rationale:     "Changed multiplier to 3",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestA, branchASHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest A: %v", err)
	}

	// Branch B: Change to multiply by 4
	repo.GitCheckout("main")
	repo.GitCheckout("branch-b")
	branchBCode := `def calculate(x):
    return x * 4
`
	repo.WriteFile("math.py", branchBCode)
	repo.GitAdd("math.py")
	repo.GitCommit("Multiply by 4")

	// Create manifest for branch B
	branchBSHA := repo.GetCurrentCommit()
	manifestB := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchBSHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "math.py",
					Symbol: "calculate",
					HunkID: "H#1",
				},
				ChangeType: "modify",
				Contract: manifest.Contract{
					Preconditions:  []string{"x must be number"},
					Postconditions: []string{"returns x * 4"},
					ErrorModel:     []string{"none"},
				},
				BehaviorClass: []string{"performance"},
				Rationale:     "Changed multiplier to 4",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestB, branchBSHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest B: %v", err)
	}

	// Attempt merge (should conflict)
	repo.GitCheckout("branch-a")
	err := repo.GitMerge("branch-b")
	if err == nil {
		t.Fatal("Expected merge conflict, but merge succeeded")
	}

	// Enrich conflicts
	conflictFile := filepath.Join(repo.Dir, "math.py")
	if err := merge.EnrichConflicts(conflictFile, "", branchASHA, branchBSHA); err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}

	// Read enriched content
	enriched, err := os.ReadFile(conflictFile)
	if err != nil {
		t.Fatalf("Failed to read enriched file: %v", err)
	}
	enrichedStr := string(enriched)

	// Verify selective injection: should contain ONLY calculate's context
	if !strings.Contains(enrichedStr, "calculate") {
		t.Error("Expected 'calculate' symbol in enriched output")
	}

	// Verify branch A's rationale appears
	if !strings.Contains(enrichedStr, "Changed multiplier to 3") {
		t.Error("Expected branch A rationale in enriched output")
	}

	// Verify branch B's rationale appears
	if !strings.Contains(enrichedStr, "Changed multiplier to 4") {
		t.Error("Expected branch B rationale in enriched output")
	}

	// Verify context markers
	if !strings.Contains(enrichedStr, "||| Gip CONTEXT") {
		t.Error("Expected Gip context markers")
	}

	// Count lines in enriched output (should be compact, not 100+ lines)
	// With selective injection, we expect ~20-40 lines (original code + conflict markers + 2 manifest entries)
	lines := strings.Split(enrichedStr, "\n")
	if len(lines) > 50 {
		t.Errorf("Enriched output too verbose: %d lines (expected < 50 for selective injection)", len(lines))
	}
	
	t.Logf("Enriched output has %d lines (selective injection working correctly)", len(lines))
}

// TestSelectiveInjection_MultiFunction tests that with multiple functions changed,
// only the relevant entry is injected for each conflict
func TestSelectiveInjection_MultiFunction(t *testing.T) {
	repo := helpers.CreateTestRepo(t)
	defer repo.Cleanup()

	// Create initial file with two functions
	initialCode := `def add(x, y):
    return x + y

def multiply(x, y):
    return x * y
`
	repo.WriteFile("math.py", initialCode)
	repo.GitAdd("math.py")
	repo.GitCommit("Initial commit")

	// Branch A: Modify both functions
	repo.GitCheckout("branch-a")
	branchACode := `def add(x, y):
    # Added validation
    return x + y

def multiply(x, y):
    # Added logging
    return x * y
`
	repo.WriteFile("math.py", branchACode)
	repo.GitAdd("math.py")
	repo.GitCommit("Add comments to both functions")

	branchASHA := repo.GetCurrentCommit()
	manifestA := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchASHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "math.py",
					Symbol: "add",
					HunkID: "H#1",
				},
				ChangeType:    "modify",
				BehaviorClass: []string{"docs"},
				Rationale:     "Added validation comment to add",
			},
			{
				Anchor: manifest.Anchor{
					File:   "math.py",
					Symbol: "multiply",
					HunkID: "H#4",
				},
				ChangeType:    "modify",
				BehaviorClass: []string{"docs"},
				Rationale:     "Added logging comment to multiply",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestA, branchASHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest A: %v", err)
	}

	// Branch B: Modify multiply differently (conflict only in multiply)
	repo.GitCheckout("main")
	repo.GitCheckout("branch-b")
	branchBCode := `def add(x, y):
    return x + y

def multiply(x, y):
    # Added error handling
    return x * y
`
	repo.WriteFile("math.py", branchBCode)
	repo.GitAdd("math.py")
	repo.GitCommit("Add error handling to multiply")

	branchBSHA := repo.GetCurrentCommit()
	manifestB := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchBSHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "math.py",
					Symbol: "multiply",
					HunkID: "H#4",
				},
				ChangeType:    "modify",
				BehaviorClass: []string{"bugfix"},
				Rationale:     "Added error handling to multiply",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestB, branchBSHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest B: %v", err)
	}

	// Attempt merge
	repo.GitCheckout("branch-a")
	err := repo.GitMerge("branch-b")
	if err == nil {
		t.Fatal("Expected merge conflict")
	}

	// Enrich conflicts
	conflictFile := filepath.Join(repo.Dir, "math.py")
	if err := merge.EnrichConflicts(conflictFile, "", branchASHA, branchBSHA); err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}

	enriched, err := os.ReadFile(conflictFile)
	if err != nil {
		t.Fatalf("Failed to read enriched file: %v", err)
	}
	enrichedStr := string(enriched)

	// Verify ONLY multiply's context appears (not add's)
	if !strings.Contains(enrichedStr, "multiply") {
		t.Error("Expected 'multiply' context in enriched output")
	}

	// Should contain branch A's multiply rationale
	if !strings.Contains(enrichedStr, "Added logging comment to multiply") {
		t.Error("Expected branch A's multiply rationale")
	}

	// Should contain branch B's multiply rationale
	if !strings.Contains(enrichedStr, "Added error handling to multiply") {
		t.Error("Expected branch B's multiply rationale")
	}

	// Should NOT contain add's rationale (no conflict in add)
	if strings.Contains(enrichedStr, "Added validation comment to add") {
		t.Error("Should NOT inject add's context (no conflict there)")
	}
}

// TestSelectiveInjection_V1Fallback tests backward compatibility with v1.0 manifests
func TestSelectiveInjection_V1Fallback(t *testing.T) {
	repo := helpers.CreateTestRepo(t)
	defer repo.Cleanup()

	initialCode := `def process():
    return True
`
	repo.WriteFile("app.py", initialCode)
	repo.GitAdd("app.py")
	repo.GitCommit("Initial")

	repo.GitCheckout("branch-a")
	branchACode := `def process():
    return False
`
	repo.WriteFile("app.py", branchACode)
	repo.GitAdd("app.py")
	repo.GitCommit("Change to False")

	branchASHA := repo.GetCurrentCommit()

	// Create v1.0 manifest (no hunk ID, no schema version)
	manifestV1 := &manifest.Manifest{
		Commit: branchASHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "app.py",
					Symbol: "process",
					// No HunkID - v1.0 format
				},
				ChangeType:    "modify",
				BehaviorClass: []string{"bugfix"},
				Rationale:     "Fixed return value",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestV1, branchASHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save v1 manifest: %v", err)
	}

	repo.GitCheckout("main")
	repo.GitCheckout("branch-b")
	branchBCode := `def process():
    return None
`
	repo.WriteFile("app.py", branchBCode)
	repo.GitAdd("app.py")
	repo.GitCommit("Change to None")

	branchBSHA := repo.GetCurrentCommit()
	manifestV1B := &manifest.Manifest{
		Commit: branchBSHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "app.py",
					Symbol: "process",
				},
				ChangeType:    "modify",
				BehaviorClass: []string{"refactor"},
				Rationale:     "Use None instead",
			},
		},
	}
	if err := helpers.SaveManifestDirect(manifestV1B, branchBSHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save v1 manifest B: %v", err)
	}

	// Merge and enrich
	repo.GitCheckout("branch-a")
	repo.GitMerge("branch-b") // Will conflict

	conflictFile := filepath.Join(repo.Dir, "app.py")
	if err := merge.EnrichConflicts(conflictFile, "", branchASHA, branchBSHA); err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}

	enriched, err := os.ReadFile(conflictFile)
	if err != nil {
		t.Fatalf("Failed to read enriched file: %v", err)
	}
	enrichedStr := string(enriched)

	// Should still inject context (fallback to file+symbol matching without hunk)
	if !strings.Contains(enrichedStr, "Fixed return value") {
		t.Error("V1 manifest context not injected (fallback should work)")
	}

	if !strings.Contains(enrichedStr, "Use None instead") {
		t.Error("V1 manifest B context not injected")
	}
}
