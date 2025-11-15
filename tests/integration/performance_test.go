package integration

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/iamHrithikRaj/gip/internal/manifest"
	"github.com/iamHrithikRaj/gip/internal/merge"
	"github.com/iamHrithikRaj/gip/tests/helpers"
)

// TestPerformance_LargeManifest tests selective injection with 100+ entry manifests
// Validates that only relevant entries are injected, not all 100
func TestPerformance_LargeManifest(t *testing.T) {
	t.Parallel()

	repo := helpers.CreateTestRepo(t)
	defer repo.Cleanup()

	// Create a large Python file with 50 functions
	var codeBuilder strings.Builder
	codeBuilder.WriteString("# Large module with many functions\n\n")
	
	for i := 1; i <= 50; i++ {
		codeBuilder.WriteString(fmt.Sprintf(`def function_%d(x):
    """Function %d"""
    return x * %d

`, i, i, i))
	}

	initialCode := codeBuilder.String()
	repo.WriteFile("large_module.py", initialCode)
	repo.GitAdd("large_module.py")
	repo.GitCommit("Initial commit with 50 functions")

	// Branch A: Modify function_25 (the target conflict)
	repo.GitCheckout("branch-a")
	branchACode := strings.Replace(initialCode, 
		`def function_25(x):
    """Function 25"""
    return x * 25`,
		`def function_25(x):
    """Function 25 - Added validation"""
    if x < 0:
        raise ValueError("x must be positive")
    return x * 25`,
		1)
	repo.WriteFile("large_module.py", branchACode)
	repo.GitAdd("large_module.py")
	repo.GitCommit("Add validation to function_25")

	branchASHA := repo.GetCurrentCommit()

	// Create a LARGE manifest with 100 entries
	// (50 real functions + 50 dummy entries to simulate a complex module)
	manifestA := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchASHA,
		GlobalIntent: &manifest.GlobalIntent{
			BehaviorClass: []string{"refactor"},
			Rationale:     "Large-scale refactoring across module",
		},
		Entries: make([]manifest.Entry, 100),
	}

	// Fill 50 real function entries
	for i := 0; i < 50; i++ {
		manifestA.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "large_module.py",
				Symbol: fmt.Sprintf("function_%d", i+1),
				HunkID: fmt.Sprintf("H#%d", i*5+1), // Approximate line numbers
			},
			ChangeType:    "modify",
			BehaviorClass: []string{"refactor"},
			Rationale:     fmt.Sprintf("Refactored function_%d for performance", i+1),
			Contract: manifest.Contract{
				Preconditions:  []string{fmt.Sprintf("x is integer for function_%d", i+1)},
				Postconditions: []string{fmt.Sprintf("returns computed value from function_%d", i+1)},
				ErrorModel:     []string{"None"},
			},
			SideEffects:         []string{"none"},
			InheritsGlobalIntent: true,
		}
	}

	// Fill 50 additional entries for other hypothetical symbols
	for i := 50; i < 100; i++ {
		manifestA.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "large_module.py",
				Symbol: fmt.Sprintf("helper_%d", i-50+1),
				HunkID: fmt.Sprintf("H#%d", i*5+1),
			},
			ChangeType:    "modify",
			BehaviorClass: []string{"refactor"},
			Rationale:     fmt.Sprintf("Refactored helper_%d", i-50+1),
			Contract: manifest.Contract{
				Preconditions:  []string{"Valid input"},
				Postconditions: []string{"Returns result"},
				ErrorModel:     []string{"None"},
			},
			SideEffects:         []string{"none"},
			InheritsGlobalIntent: true,
		}
	}

	// Override function_25 with specific validation rationale
	manifestA.Entries[24].Rationale = "Added strict input validation to function_25"
	manifestA.Entries[24].BehaviorClass = []string{"validation", "bugfix"}

	if err := helpers.SaveManifestDirect(manifestA, branchASHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest A: %v", err)
	}

	// Branch B: Modify function_25 differently (conflict)
	repo.GitCheckout("main")
	repo.GitCheckout("branch-b")
	branchBCode := strings.Replace(initialCode,
		`def function_25(x):
    """Function 25"""
    return x * 25`,
		`def function_25(x):
    """Function 25 - Added error handling"""
    try:
        result = x * 25
        return result
    except Exception as e:
        return 0`,
		1)
	repo.WriteFile("large_module.py", branchBCode)
	repo.GitAdd("large_module.py")
	repo.GitCommit("Add error handling to function_25")

	branchBSHA := repo.GetCurrentCommit()

	// Create large manifest for branch B as well
	manifestB := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchBSHA,
		GlobalIntent: &manifest.GlobalIntent{
			BehaviorClass: []string{"bugfix"},
			Rationale:     "Add error handling across module",
		},
		Entries: make([]manifest.Entry, 100),
	}

	// Similar structure to manifestA
	for i := 0; i < 50; i++ {
		manifestB.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "large_module.py",
				Symbol: fmt.Sprintf("function_%d", i+1),
				HunkID: fmt.Sprintf("H#%d", i*5+1),
			},
			ChangeType:           "modify",
			BehaviorClass:        []string{"bugfix"},
			Rationale:            fmt.Sprintf("Added error handling to function_%d", i+1),
			InheritsGlobalIntent: true,
		}
	}

	for i := 50; i < 100; i++ {
		manifestB.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "large_module.py",
				Symbol: fmt.Sprintf("helper_%d", i-50+1),
				HunkID: fmt.Sprintf("H#%d", i*5+1),
			},
			ChangeType:           "modify",
			BehaviorClass:        []string{"bugfix"},
			Rationale:            fmt.Sprintf("Added error handling to helper_%d", i-50+1),
			InheritsGlobalIntent: true,
		}
	}

	// Override function_25
	manifestB.Entries[24].Rationale = "Added try-catch error handling to function_25"

	if err := helpers.SaveManifestDirect(manifestB, branchBSHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest B: %v", err)
	}

	// Attempt merge
	repo.GitCheckout("branch-a")
	err := repo.GitMerge("branch-b")
	if err == nil {
		t.Fatal("Expected merge conflict")
	}

	// Measure enrichment performance
	conflictFile := filepath.Join(repo.Dir, "large_module.py")
	startTime := time.Now()
	
	if err := merge.EnrichConflicts(conflictFile, "", branchASHA, branchBSHA); err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}
	
	duration := time.Since(startTime)
	t.Logf("EnrichConflicts completed in %v for manifest with 100 entries", duration)

	// Performance threshold: should complete in under 1 second
	if duration > 1*time.Second {
		t.Errorf("EnrichConflicts took too long: %v (expected < 1s)", duration)
	}

	// Read enriched file
	enriched, err := os.ReadFile(conflictFile)
	if err != nil {
		t.Fatalf("Failed to read enriched file: %v", err)
	}
	enrichedStr := string(enriched)
	lineCount := strings.Count(enrichedStr, "\n")

	t.Logf("Enriched output has %d lines (manifest has 100 entries)", lineCount)

	// Critical assertion: Output should be reasonable despite 100-entry manifest
	// The file has 50 functions (214 lines), conflict adds ~20 lines of TOON context
	// Total: ~236 lines is correct (not 1000+ if all entries were injected)
	if lineCount > 300 {
		t.Errorf("Output too large: %d lines (suggests multiple entries injected)", lineCount)
	}

	// Verify only function_25 context appears
	if !strings.Contains(enrichedStr, "function_25") {
		t.Error("Expected function_25 context in enriched output")
	}

	// Verify function_25's specific rationale from branch A
	if !strings.Contains(enrichedStr, "Added strict input validation to function_25") {
		t.Error("Expected branch A's function_25 rationale")
	}

	// Verify function_25's specific rationale from branch B
	if !strings.Contains(enrichedStr, "Added try-catch error handling to function_25") {
		t.Error("Expected branch B's function_25 rationale")
	}

	// Verify OTHER functions' rationales do NOT appear in TOON context
	// (function_1 and function_50 will appear as function definitions, but not in TOON)
	if strings.Contains(enrichedStr, "Refactored function_1 for performance") {
		t.Error("Should NOT inject function_1's rationale (not in conflict)")
	}
	if strings.Contains(enrichedStr, "Refactored function_50") {
		t.Error("Should NOT inject function_50's rationale (not in conflict)")
	}
	if strings.Contains(enrichedStr, "helper_1") {
		t.Error("Should NOT inject helper_1 context (not in conflict)")
	}

	// Success metrics
	t.Logf("✅ Selective injection working: Only function_25 context injected")
	t.Logf("✅ Performance: %v for 100-entry manifest", duration)
	t.Logf("✅ Output size: %d lines (vs potential 1000+ if all entries injected)", lineCount)
}

// TestPerformance_CompareV1VsV2 simulates v1.0 behavior (inject all) vs v2.0 (selective)
func TestPerformance_CompareV1VsV2(t *testing.T) {
	t.Parallel()

	repo := helpers.CreateTestRepo(t)
	defer repo.Cleanup()

	// Simple conflict scenario
	initialCode := `def foo():
    return 1
`
	repo.WriteFile("app.py", initialCode)
	repo.GitAdd("app.py")
	repo.GitCommit("Initial")

	// Branch A: Modify foo
	repo.GitCheckout("branch-a")
	branchACode := `def foo():
    return 2
`
	repo.WriteFile("app.py", branchACode)
	repo.GitAdd("app.py")
	repo.GitCommit("Change foo to 2")
	branchASHA := repo.GetCurrentCommit()

	// Create a manifest with 50 entries (simulating large commit)
	manifestA := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchASHA,
		Entries:       make([]manifest.Entry, 50),
	}

	for i := 0; i < 50; i++ {
		manifestA.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "app.py",
				Symbol: fmt.Sprintf("function_%d", i),
				HunkID: fmt.Sprintf("H#%d", i*10),
			},
			ChangeType:    "modify",
			BehaviorClass: []string{"refactor"},
			Rationale:     fmt.Sprintf("Modified function_%d with long explanation that takes multiple lines and adds significant content to the manifest", i),
		}
	}

	// Make the first entry match the actual conflict
	manifestA.Entries[0].Anchor.Symbol = "foo"
	manifestA.Entries[0].Anchor.HunkID = "H#1"
	manifestA.Entries[0].Rationale = "Changed return value to 2"

	if err := helpers.SaveManifestDirect(manifestA, branchASHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest A: %v", err)
	}

	// Branch B: Modify foo differently
	repo.GitCheckout("main")
	repo.GitCheckout("branch-b")
	branchBCode := `def foo():
    return 3
`
	repo.WriteFile("app.py", branchBCode)
	repo.GitAdd("app.py")
	repo.GitCommit("Change foo to 3")
	branchBSHA := repo.GetCurrentCommit()

	// Similar large manifest for branch B
	manifestB := &manifest.Manifest{
		SchemaVersion: "2.0",
		Commit:        branchBSHA,
		Entries:       make([]manifest.Entry, 50),
	}

	for i := 0; i < 50; i++ {
		manifestB.Entries[i] = manifest.Entry{
			Anchor: manifest.Anchor{
				File:   "app.py",
				Symbol: fmt.Sprintf("function_%d", i),
				HunkID: fmt.Sprintf("H#%d", i*10),
			},
			ChangeType:    "modify",
			BehaviorClass: []string{"refactor"},
			Rationale:     fmt.Sprintf("Modified function_%d with different long explanation", i),
		}
	}

	manifestB.Entries[0].Anchor.Symbol = "foo"
	manifestB.Entries[0].Anchor.HunkID = "H#1"
	manifestB.Entries[0].Rationale = "Changed return value to 3"

	if err := helpers.SaveManifestDirect(manifestB, branchBSHA, repo.Dir); err != nil {
		t.Fatalf("Failed to save manifest B: %v", err)
	}

	// Merge and enrich
	repo.GitCheckout("branch-a")
	err := repo.GitMerge("branch-b")
	if err == nil {
		t.Fatal("Expected merge conflict")
	}

	conflictFile := filepath.Join(repo.Dir, "app.py")
	if err := merge.EnrichConflicts(conflictFile, "", branchASHA, branchBSHA); err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}

	enriched, err := os.ReadFile(conflictFile)
	if err != nil {
		t.Fatalf("Failed to read enriched file: %v", err)
	}
	enrichedStr := string(enriched)
	lineCount := strings.Count(enrichedStr, "\n")

	t.Logf("V2.0 selective injection: %d lines of output", lineCount)

	// In v1.0, this would inject ALL 50 entries from both sides = ~500+ lines
	// In v2.0, only 1 entry from each side = ~20-30 lines
	if lineCount > 50 {
		t.Errorf("Output suggests non-selective injection: %d lines", lineCount)
	}

	// Verify only foo's context appears
	if !strings.Contains(enrichedStr, "Changed return value to 2") {
		t.Error("Expected branch A's foo rationale")
	}
	if !strings.Contains(enrichedStr, "Changed return value to 3") {
		t.Error("Expected branch B's foo rationale")
	}

	// Verify other entries don't appear
	if strings.Contains(enrichedStr, "function_1") {
		t.Error("Should NOT inject function_1 (not in conflict)")
	}
	if strings.Contains(enrichedStr, "function_49") {
		t.Error("Should NOT inject function_49 (not in conflict)")
	}

	t.Logf("✅ V2.0 achieves ~%dx size reduction vs hypothetical v1.0 full injection", 500/lineCount)
}
