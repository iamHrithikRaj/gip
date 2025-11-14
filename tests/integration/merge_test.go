package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/iamHrithikRaj/gip/internal/manifest"
	"github.com/iamHrithikRaj/gip/internal/merge"
	"github.com/iamHrithikRaj/gip/tests/helpers"
)

func TestMergeWorkflow(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test repository
	repoDir := helpers.SetupGitRepo(t)
	t.Logf("Test repo: %s", repoDir)

	// Create initial file
	initialContent := `def calculate_total(items):
    total = 0
    for item in items:
        total += item.price
    return total
`
	sha1 := helpers.CreateCommit(t, "cart.py", initialContent, "Initial commit")
	t.Logf("Initial commit: %s", sha1)

	// Create manifest for initial commit
	m1 := helpers.NewTestManifest(sha1, "cart.py", "calculate_total", "Initial implementation")
	manifest.Save(m1, sha1)

	// Branch A: Add tax calculation
	helpers.CreateBranch(t, "feature-tax")
	taxContent := `def calculate_total(items):
    total = 0
    for item in items:
        total += item.price * 1.08  # Add 8% tax
    return total
`
	sha2 := helpers.CreateCommit(t, "cart.py", taxContent, "Add tax calculation")

	// Create manifest for tax feature
	m2 := helpers.NewFeatureManifest(sha2, "cart.py", "calculate_total", []string{"none"})
	m2.Entries[0].Rationale = "Added 8% sales tax"
	m2.Entries[0].Contract.Postconditions = []string{"returns total with 8% tax"}
	manifest.Save(m2, sha2)

	// Go back to main and create Branch B
	helpers.CheckoutBranch(t, "main")
	helpers.CreateBranch(t, "feature-shipping")
	shippingContent := `def calculate_total(items):
    total = 0
    for item in items:
        total += item.price
    return total + 5.99  # Add flat shipping
`
	sha3 := helpers.CreateCommit(t, "cart.py", shippingContent, "Add shipping fee")

	// Create manifest for shipping feature
	m3 := helpers.NewFeatureManifest(sha3, "cart.py", "calculate_total", []string{"none"})
	m3.Entries[0].Rationale = "Added $5.99 flat shipping fee"
	m3.Entries[0].Contract.Postconditions = []string{"returns total plus shipping"}
	manifest.Save(m3, sha3)

	// Try to merge feature-tax (should conflict)
	success, err := helpers.MergeBranch(t, "feature-tax")
	if err != nil {
		t.Fatalf("Merge failed with error: %v", err)
	}

	if success {
		t.Fatal("Expected merge conflict, but merge succeeded")
	}

	// Verify conflict exists
	conflictedFiles := helpers.GetConflictedFiles(t)
	if len(conflictedFiles) != 1 || conflictedFiles[0] != "cart.py" {
		t.Fatalf("Expected cart.py to have conflicts, got: %v", conflictedFiles)
	}

	// Enrich conflicts
	err = merge.EnrichConflicts("cart.py", sha1, sha3, sha2)
	if err != nil {
		t.Fatalf("Failed to enrich conflicts: %v", err)
	}

	// Verify enrichment
	content := helpers.ReadFile(t, "cart.py")

	// Check for Gip context markers
	helpers.AssertFileContains(t, "cart.py", "||| Gip CONTEXT")
	helpers.AssertFileContains(t, "cart.py", "HEAD")
	helpers.AssertFileContains(t, "cart.py", "MERGE_HEAD")

	// Check for commit SHAs
	if len(sha2) >= 8 && len(sha3) >= 8 {
		helpers.AssertFileContains(t, "cart.py", sha2[:8])
		helpers.AssertFileContains(t, "cart.py", sha3[:8])
	}

	t.Logf("Enriched conflict content:\n%s", content)
}

func TestCommitWorkflow(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Setup test repository
	repoDir := helpers.SetupGitRepo(t)
	t.Logf("Test repo: %s", repoDir)

	// Create a file and commit it first
	content := `package main

func HelloWorld() string {
    return "Hello, World!"
}
`
	sha := helpers.CreateCommit(t, "hello.go", content, "Add HelloWorld")

	// Create a manifest manually
	m := helpers.NewTestManifest(sha, "hello.go", "HelloWorld", "Initial greeting function")

	// Save manifest
	err := manifest.Save(m, sha)
	helpers.AssertNoError(t, err, "Failed to save manifest")

	// Load manifest back
	loaded, err := manifest.Load(sha)
	helpers.AssertNoError(t, err, "Failed to load manifest")

	// Verify loaded manifest
	if loaded.Commit != sha {
		t.Errorf("Expected commit %s, got %s", sha, loaded.Commit)
	}

	if len(loaded.Entries) != 1 {
		t.Fatalf("Expected 1 entry, got %d", len(loaded.Entries))
	}

	entry := loaded.Entries[0]
	if entry.Anchor.Symbol != "HelloWorld" {
		t.Errorf("Expected symbol HelloWorld, got %s", entry.Anchor.Symbol)
	}

	if entry.Rationale != "Initial greeting function" {
		t.Errorf("Expected correct rationale, got %s", entry.Rationale)
	}
}

func TestMultipleFileConflicts(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	repoDir := helpers.SetupGitRepo(t)
	t.Logf("Test repo: %s", repoDir)

	// Create two files
	file1 := "func1.py"
	file2 := "func2.py"

	helpers.WriteFile(t, file1, "def func1():\n    return 1\n")
	helpers.WriteFile(t, file2, "def func2():\n    return 2\n")
	helpers.StageFile(t, file1)
	helpers.StageFile(t, file2)
	sha1 := helpers.CreateCommit(t, file1, "def func1():\n    return 1\n", "Initial commit")

	// Branch A: Modify both files
	helpers.CreateBranch(t, "branch-a")
	helpers.WriteFile(t, file1, "def func1():\n    return 10  # Branch A\n")
	helpers.WriteFile(t, file2, "def func2():\n    return 20  # Branch A\n")
	helpers.StageFile(t, file1)
	helpers.StageFile(t, file2)
	sha2 := helpers.CreateCommit(t, file1, "def func1():\n    return 10  # Branch A\n", "Branch A changes")

	// Branch B: Modify both files differently
	helpers.CheckoutBranch(t, "main")
	helpers.CreateBranch(t, "branch-b")
	helpers.WriteFile(t, file1, "def func1():\n    return 100  # Branch B\n")
	helpers.WriteFile(t, file2, "def func2():\n    return 200  # Branch B\n")
	helpers.StageFile(t, file1)
	helpers.StageFile(t, file2)
	sha3 := helpers.CreateCommit(t, file1, "def func1():\n    return 100  # Branch B\n", "Branch B changes")

	// Create manifests
	m2 := helpers.NewManifestWithMultipleEntries(sha2, []helpers.ManifestEntrySpec{
		helpers.DefaultEntrySpec(file1, "func1"),
		helpers.DefaultEntrySpec(file2, "func2"),
	})
	m2.Entries[0].Rationale = "Branch A changes to func1"
	m2.Entries[1].Rationale = "Branch A changes to func2"
	manifest.Save(m2, sha2)

	m3 := helpers.NewManifestWithMultipleEntries(sha3, []helpers.ManifestEntrySpec{
		helpers.DefaultEntrySpec(file1, "func1"),
		helpers.DefaultEntrySpec(file2, "func2"),
	})
	m3.Entries[0].Rationale = "Branch B changes to func1"
	m3.Entries[1].Rationale = "Branch B changes to func2"
	manifest.Save(m3, sha3)

	// Merge branch-a (should conflict)
	success, err := helpers.MergeBranch(t, "branch-a")
	helpers.AssertNoError(t, err, "Merge should not error")

	if success {
		t.Fatal("Expected conflicts")
	}

	// Get conflicted files
	conflicted := helpers.GetConflictedFiles(t)
	if len(conflicted) != 2 {
		t.Fatalf("Expected 2 conflicted files, got %d: %v", len(conflicted), conflicted)
	}

	// Enrich both files
	for _, file := range conflicted {
		merge.EnrichConflicts(file, sha1, sha3, sha2)
		helpers.AssertFileContains(t, file, "||| Gip CONTEXT")
	}
}

func TestNoConflictMerge(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	repoDir := helpers.SetupGitRepo(t)
	t.Logf("Test repo: %s", repoDir)

	// Create file on main
	helpers.WriteFile(t, "file1.txt", "Line 1\n")
	_ = helpers.CreateCommit(t, "file1.txt", "Line 1\n", "Initial")

	// Branch A: Add to end
	helpers.CreateBranch(t, "branch-a")
	helpers.WriteFile(t, "file1.txt", "Line 1\nLine 2 from A\n")
	_ = helpers.CreateCommit(t, "file1.txt", "Line 1\nLine 2 from A\n", "Add line 2")

	// Branch B: Add different file
	helpers.CheckoutBranch(t, "main")
	helpers.CreateBranch(t, "branch-b")
	helpers.WriteFile(t, "file2.txt", "New file\n")
	_ = helpers.CreateCommit(t, "file2.txt", "New file\n", "Add file2")

	// Merge should succeed (no conflicts)
	success, err := helpers.MergeBranch(t, "branch-a")
	helpers.AssertNoError(t, err, "Merge should succeed")

	if !success {
		t.Fatal("Expected successful merge")
	}

	// Verify both changes are present
	if _, err := os.Stat(filepath.Join(repoDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt should exist")
	}

	if _, err := os.Stat(filepath.Join(repoDir, "file2.txt")); os.IsNotExist(err) {
		t.Error("file2.txt should exist")
	}

	content1 := helpers.ReadFile(t, "file1.txt")
	if content1 != "Line 1\nLine 2 from A\n" {
		t.Errorf("Unexpected content in file1.txt: %s", content1)
	}
}
