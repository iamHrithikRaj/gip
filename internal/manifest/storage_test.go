package manifest

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestSaveAndLoad(t *testing.T) {
	t.Parallel()
	
	// Create temporary directory for test
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Initialize git repo (required by EnsureGipDir)
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	
	// Create test manifest
	testCommit := "abc123def456"
	m := &Manifest{
		Commit: testCommit,
		Entries: []Entry{
			{
				Anchor: Anchor{
					File:   "test.go",
					Symbol: "TestFunc",
					HunkID: "1",
				},
				ChangeType: "modify",
				Contract: Contract{
					Preconditions:  []string{"input is valid"},
					Postconditions: []string{"returns result"},
					ErrorModel:     []string{"none"},
				},
				BehaviorClass: []string{"feature"},
				SideEffects:   []string{"none"},
				Rationale:     "Added new feature",
			},
		},
	}
	
	// Test Save
	if err := Save(m, testCommit); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	
	// Verify file exists
	manifestPath := filepath.Join(".gip", "manifest", testCommit+".json")
	if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
		t.Fatalf("Manifest file was not created at %s", manifestPath)
	}
	
	// Test Load
	loaded, err := Load(testCommit)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	// Verify loaded data
	if loaded.Commit != m.Commit {
		t.Errorf("Expected commit %s, got %s", m.Commit, loaded.Commit)
	}
	
	if len(loaded.Entries) != len(m.Entries) {
		t.Errorf("Expected %d entries, got %d", len(m.Entries), len(loaded.Entries))
	}
	
	if len(loaded.Entries) > 0 {
		entry := loaded.Entries[0]
		expectedEntry := m.Entries[0]
		
		if entry.Anchor.File != expectedEntry.Anchor.File {
			t.Errorf("Expected file %s, got %s", expectedEntry.Anchor.File, entry.Anchor.File)
		}
		
		if entry.Anchor.Symbol != expectedEntry.Anchor.Symbol {
			t.Errorf("Expected symbol %s, got %s", expectedEntry.Anchor.Symbol, entry.Anchor.Symbol)
		}
		
		if entry.Rationale != expectedEntry.Rationale {
			t.Errorf("Expected rationale %s, got %s", expectedEntry.Rationale, entry.Rationale)
		}
		
		if len(entry.BehaviorClass) != len(expectedEntry.BehaviorClass) {
			t.Errorf("Expected %d behavior classes, got %d", len(expectedEntry.BehaviorClass), len(entry.BehaviorClass))
		}
	}
}

func TestLoadNonExistent(t *testing.T) {
	t.Parallel()
	
	// Create temporary directory
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Initialize git repo
	exec.Command("git", "init").Run()
	os.MkdirAll(".gip/manifest", 0755)
	
	// Try to load non-existent manifest
	_, err := Load("nonexistent123")
	if err == nil {
		t.Error("Expected error loading non-existent manifest, got nil")
	}
}

func TestSaveWithoutGipDir(t *testing.T) {
	t.Parallel()
	
	// Create temporary directory without .gip
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Initialize git repo
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	
	// Should create .gip directory automatically
	m := &Manifest{
		Commit: "test123",
		Entries: []Entry{},
	}
	
	if err := Save(m, "test123"); err != nil {
		t.Fatalf("Save should create .gip directory, but failed: %v", err)
	}
	
	// Verify .gip directory was created
	if _, err := os.Stat(".gip/manifest"); os.IsNotExist(err) {
		t.Error(".gip/manifest directory was not created")
	}
}

func TestSaveMultipleEntries(t *testing.T) {
	t.Parallel()
	
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Initialize git repo
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	
	// Create manifest with multiple entries
	testCommit := "multi123"
	m := &Manifest{
		Commit: testCommit,
		Entries: []Entry{
			{
				Anchor:        Anchor{File: "file1.go", Symbol: "Func1"},
				BehaviorClass: []string{"feature"},
				Rationale:     "First change",
			},
			{
				Anchor:        Anchor{File: "file2.go", Symbol: "Func2"},
				BehaviorClass: []string{"bugfix"},
				Rationale:     "Second change",
			},
			{
				Anchor:        Anchor{File: "file3.go", Symbol: "Func3"},
				BehaviorClass: []string{"refactor"},
				Rationale:     "Third change",
			},
		},
	}
	
	if err := Save(m, testCommit); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	
	loaded, err := Load(testCommit)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	if len(loaded.Entries) != 3 {
		t.Errorf("Expected 3 entries, got %d", len(loaded.Entries))
	}
	
	// Verify all entries
	for i, entry := range loaded.Entries {
		expected := m.Entries[i]
		if entry.Anchor.File != expected.Anchor.File {
			t.Errorf("Entry %d: expected file %s, got %s", i, expected.Anchor.File, entry.Anchor.File)
		}
		if entry.Rationale != expected.Rationale {
			t.Errorf("Entry %d: expected rationale %s, got %s", i, expected.Rationale, entry.Rationale)
		}
	}
}

func TestEmptyManifest(t *testing.T) {
	t.Parallel()
	
	tmpDir := t.TempDir()
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	os.Chdir(tmpDir)
	
	// Initialize git repo
	exec.Command("git", "init").Run()
	exec.Command("git", "config", "user.name", "Test").Run()
	exec.Command("git", "config", "user.email", "test@test.com").Run()
	
	// Create empty manifest (no entries)
	testCommit := "empty123"
	m := &Manifest{
		Commit:  testCommit,
		Entries: []Entry{},
	}
	
	if err := Save(m, testCommit); err != nil {
		t.Fatalf("Save failed for empty manifest: %v", err)
	}
	
	loaded, err := Load(testCommit)
	if err != nil {
		t.Fatalf("Load failed for empty manifest: %v", err)
	}
	
	if loaded.Commit != testCommit {
		t.Errorf("Expected commit %s, got %s", testCommit, loaded.Commit)
	}
	
	if len(loaded.Entries) != 0 {
		t.Errorf("Expected 0 entries, got %d", len(loaded.Entries))
	}
}
