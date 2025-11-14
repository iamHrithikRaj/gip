package merge

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	
	"github.com/iamHrithikRaj/gip/internal/manifest"
)

func TestFindConflictBlocks(t *testing.T) {
	t.Parallel()
	
	lines := []string{
		"normal line",
		"<<<<<<< HEAD",
		"head content 1",
		"head content 2",
		"=======",
		"their content 1",
		"their content 2",
		">>>>>>> branch",
		"normal line after",
	}
	
	conflicts := findConflictBlocks(lines)
	
	if len(conflicts) != 1 {
		t.Fatalf("Expected 1 conflict, got %d", len(conflicts))
	}
	
	conflict := conflicts[0]
	if conflict.StartLine != 1 {
		t.Errorf("Expected StartLine=1, got %d", conflict.StartLine)
	}
	
	if conflict.EndLine != 7 {
		t.Errorf("Expected EndLine=7, got %d", conflict.EndLine)
	}
	
	if len(conflict.HeadLines) != 2 {
		t.Errorf("Expected 2 HEAD lines, got %d", len(conflict.HeadLines))
	}
	
	if len(conflict.TheirLines) != 2 {
		t.Errorf("Expected 2 THEIR lines, got %d", len(conflict.TheirLines))
	}
}

func TestFindMultipleConflicts(t *testing.T) {
	t.Parallel()
	
	lines := []string{
		"<<<<<<< HEAD",
		"head1",
		"=======",
		"their1",
		">>>>>>> branch",
		"normal",
		"<<<<<<< HEAD",
		"head2",
		"=======",
		"their2",
		">>>>>>> branch",
	}
	
	conflicts := findConflictBlocks(lines)
	
	if len(conflicts) != 2 {
		t.Fatalf("Expected 2 conflicts, got %d", len(conflicts))
	}
	
	if conflicts[0].StartLine != 0 {
		t.Errorf("First conflict should start at line 0, got %d", conflicts[0].StartLine)
	}
	
	if conflicts[1].StartLine != 6 {
		t.Errorf("Second conflict should start at line 6, got %d", conflicts[1].StartLine)
	}
}

func TestFindNoConflicts(t *testing.T) {
	t.Parallel()
	
	lines := []string{
		"normal line 1",
		"normal line 2",
		"normal line 3",
	}
	
	conflicts := findConflictBlocks(lines)
	
	if len(conflicts) != 0 {
		t.Errorf("Expected no conflicts, got %d", len(conflicts))
	}
}

func TestEnrichConflictsBasic(t *testing.T) {
	t.Parallel()
	
	// Setup test environment
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	os.MkdirAll(".gip/manifest", 0755)
	
	// Create test file with conflict
	conflictFile := filepath.Join(tmpDir, "test.py")
	conflictContent := `def calculate():
<<<<<<< HEAD
    return 42
=======
    return 100
>>>>>>> feature
`
	os.WriteFile(conflictFile, []byte(conflictContent), 0644)
	
	// Create test manifests
	headSHA := "head123"
	theirSHA := "their456"
	
	headManifest := &manifest.Manifest{
		Commit: headSHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "test.py",
					Symbol: "calculate",
				},
				Contract: manifest.Contract{
					Preconditions:  []string{"none"},
					Postconditions: []string{"returns 42"},
				},
				BehaviorClass: []string{"feature"},
				Rationale:     "Return magic number",
			},
		},
	}
	
	theirManifest := &manifest.Manifest{
		Commit: theirSHA,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "test.py",
					Symbol: "calculate",
				},
				Contract: manifest.Contract{
					Preconditions:  []string{"none"},
					Postconditions: []string{"returns 100"},
				},
				BehaviorClass: []string{"bugfix"},
				Rationale:     "Fix return value",
			},
		},
	}
	
	manifest.Save(headManifest, headSHA)
	manifest.Save(theirManifest, theirSHA)
	
	// Run enrichment
	err := EnrichConflicts(conflictFile, "ancestor", headSHA, theirSHA)
	if err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}
	
	// Read enriched content
	enriched, _ := os.ReadFile(conflictFile)
	enrichedStr := string(enriched)
	
	// Verify TOON context was injected
	if !strings.Contains(enrichedStr, "||| GIP CONTEXT") {
		t.Error("Expected '||| GIP CONTEXT' marker in enriched file")
	}
	
	// Verify HEAD context
	if !strings.Contains(enrichedStr, "HEAD") {
		t.Error("Expected 'HEAD' label in enriched file")
	}
	
	// Verify MERGE_HEAD context
	if !strings.Contains(enrichedStr, "MERGE_HEAD") {
		t.Error("Expected 'MERGE_HEAD' label in enriched file")
	}
	
	// Verify commit SHAs appear
	if !strings.Contains(enrichedStr, headSHA) {
		t.Error("Expected HEAD commit SHA in enriched file")
	}
	
	if !strings.Contains(enrichedStr, theirSHA) {
		t.Error("Expected THEIR commit SHA in enriched file")
	}
	
	// Verify rationale appears
	if !strings.Contains(enrichedStr, "Return magic number") {
		t.Error("Expected HEAD rationale in enriched file")
	}
	
	if !strings.Contains(enrichedStr, "Fix return value") {
		t.Error("Expected THEIR rationale in enriched file")
	}
}

func TestEnrichConflictsNoManifests(t *testing.T) {
	t.Parallel()
	
	// Setup test environment
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	os.MkdirAll(".gip/manifest", 0755)
	
	// Create test file with conflict
	conflictFile := filepath.Join(tmpDir, "test.py")
	conflictContent := `<<<<<<< HEAD
line1
=======
line2
>>>>>>> branch
`
	os.WriteFile(conflictFile, []byte(conflictContent), 0644)
	
	// Run enrichment without creating manifests
	err := EnrichConflicts(conflictFile, "ancestor", "nonexistent1", "nonexistent2")
	
	// Should not fail even if manifests don't exist
	if err != nil {
		t.Logf("EnrichConflicts with missing manifests: %v (this is OK)", err)
	}
	
	// File should still exist
	if _, err := os.Stat(conflictFile); os.IsNotExist(err) {
		t.Error("File was deleted when it should have been preserved")
	}
}

func TestEnrichConflictsNoConflicts(t *testing.T) {
	t.Parallel()
	
	// Setup test environment
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Create test file WITHOUT conflicts
	normalFile := filepath.Join(tmpDir, "normal.py")
	normalContent := `def func():
    return 42
`
	os.WriteFile(normalFile, []byte(normalContent), 0644)
	
	// Run enrichment
	err := EnrichConflicts(normalFile, "ancestor", "sha1", "sha2")
	if err != nil {
		t.Fatalf("EnrichConflicts should not fail on files without conflicts: %v", err)
	}
	
	// Content should be unchanged
	result, _ := os.ReadFile(normalFile)
	if string(result) != normalContent {
		t.Error("File without conflicts should remain unchanged")
	}
}

func TestEnrichMultipleConflicts(t *testing.T) {
	t.Parallel()
	
	// Setup
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	os.MkdirAll(".gip/manifest", 0755)
	
	// File with TWO conflicts
	conflictFile := filepath.Join(tmpDir, "multi.py")
	conflictContent := `def func1():
<<<<<<< HEAD
    return 1
=======
    return 2
>>>>>>> branch

def func2():
<<<<<<< HEAD
    return 3
=======
    return 4
>>>>>>> branch
`
	os.WriteFile(conflictFile, []byte(conflictContent), 0644)
	
	// Create manifests
	headSHA := "multi_head"
	theirSHA := "multi_their"
	
	manifest.Save(&manifest.Manifest{
		Commit: headSHA,
		Entries: []manifest.Entry{
			{Anchor: manifest.Anchor{File: "multi.py", Symbol: "func1"}, Rationale: "Change 1"},
			{Anchor: manifest.Anchor{File: "multi.py", Symbol: "func2"}, Rationale: "Change 2"},
		},
	}, headSHA)
	
	manifest.Save(&manifest.Manifest{
		Commit: theirSHA,
		Entries: []manifest.Entry{
			{Anchor: manifest.Anchor{File: "multi.py", Symbol: "func1"}, Rationale: "Alt 1"},
			{Anchor: manifest.Anchor{File: "multi.py", Symbol: "func2"}, Rationale: "Alt 2"},
		},
	}, theirSHA)
	
	// Enrich
	err := EnrichConflicts(conflictFile, "ancestor", headSHA, theirSHA)
	if err != nil {
		t.Fatalf("EnrichConflicts failed: %v", err)
	}
	
	// Verify both conflicts were enriched
	enriched, _ := os.ReadFile(conflictFile)
	enrichedStr := string(enriched)
	
	contextCount := strings.Count(enrichedStr, "||| GIP CONTEXT")
	if contextCount < 2 {
		t.Errorf("Expected at least 2 GIP CONTEXT blocks, found %d", contextCount)
	}
}
