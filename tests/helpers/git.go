package helpers

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// TestRepo represents a temporary git repository for testing
type TestRepo struct {
	Dir         string
	t           *testing.T
	originalDir string
}

// CreateTestRepo creates a new temporary git repository
func CreateTestRepo(t *testing.T) *TestRepo {
	t.Helper()

	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	repo := &TestRepo{
		Dir:         tmpDir,
		t:           t,
		originalDir: originalDir,
	}

	// Change to temp directory
	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo
	repo.runGit("init", "-b", "main")
	repo.runGit("config", "user.name", "Test User")
	repo.runGit("config", "user.email", "test@example.com")

	// Create .gip directory structure
	if err := os.MkdirAll(filepath.Join(tmpDir, ".gip", "manifest"), 0755); err != nil {
		t.Fatalf("Failed to create .gip directory: %v", err)
	}

	return repo
}

// Cleanup restores the original directory
func (r *TestRepo) Cleanup() {
	os.Chdir(r.originalDir)
}

// WriteFile writes a file to the repository
func (r *TestRepo) WriteFile(name, content string) {
	r.t.Helper()
	path := filepath.Join(r.Dir, name)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		r.t.Fatalf("Failed to write file %s: %v", name, err)
	}
}

// GitAdd stages files
func (r *TestRepo) GitAdd(files ...string) {
	r.t.Helper()
	r.runGit(append([]string{"add"}, files...)...)
}

// GitCommit creates a commit
func (r *TestRepo) GitCommit(message string) {
	r.t.Helper()
	r.runGit("commit", "-m", message)
}

// GitCheckout checks out a branch (creates if doesn't exist with -b flag)
func (r *TestRepo) GitCheckout(branch string) {
	r.t.Helper()
	// Try checkout first, if fails, create new branch
	cmd := exec.Command("git", "checkout", branch)
	cmd.Dir = r.Dir
	if err := cmd.Run(); err != nil {
		// Branch doesn't exist, create it
		r.runGit("checkout", "-b", branch)
	}
}

// GitMerge attempts to merge a branch
func (r *TestRepo) GitMerge(branch string) error {
	cmd := exec.Command("git", "merge", branch)
	cmd.Dir = r.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		// Merge failed (possibly conflict), but that's expected in tests
		if strings.Contains(string(output), "CONFLICT") {
			return fmt.Errorf("merge conflict: %s", string(output))
		}
		return err
	}
	return nil
}

// GetCurrentCommit returns the current commit SHA
func (r *TestRepo) GetCurrentCommit() string {
	r.t.Helper()
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = r.Dir
	output, err := cmd.Output()
	if err != nil {
		r.t.Fatalf("Failed to get current commit: %v", err)
	}
	return strings.TrimSpace(string(output))
}

// runGit runs a git command in the repository
func (r *TestRepo) runGit(args ...string) {
	r.t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = r.Dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		r.t.Fatalf("Git command failed: git %s\nOutput: %s\nError: %v",
			strings.Join(args, " "), string(output), err)
	}
}

// SetupGitRepo creates a temporary Git repository for testing
func SetupGitRepo(t *testing.T) string {
	t.Helper()

	tmpDir := t.TempDir()

	// Change to temp directory
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get working directory: %v", err)
	}

	// Register cleanup to restore directory
	t.Cleanup(func() {
		os.Chdir(originalDir)
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change to temp directory: %v", err)
	}

	// Initialize git repo
	runCommand(t, "git", "init", "-b", "main")
	runCommand(t, "git", "config", "user.name", "Test User")
	runCommand(t, "git", "config", "user.email", "test@example.com")

	return tmpDir
}

// CreateCommit creates a file and commits it
func CreateCommit(t *testing.T, filename, content, message string) string {
	t.Helper()

	// Write file
	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", filename, err)
	}

	// Stage and commit
	runCommand(t, "git", "add", filename)
	runCommand(t, "git", "commit", "-m", message)

	// Get commit SHA
	output := runCommand(t, "git", "rev-parse", "HEAD")
	return strings.TrimSpace(output)
}

// CreateBranch creates and checks out a new branch
func CreateBranch(t *testing.T, name string) {
	t.Helper()
	runCommand(t, "git", "checkout", "-b", name)
}

// CheckoutBranch checks out an existing branch
func CheckoutBranch(t *testing.T, name string) {
	t.Helper()
	runCommand(t, "git", "checkout", name)
}

// MergeBranch attempts to merge a branch (may conflict)
func MergeBranch(t *testing.T, branch string) (bool, error) {
	t.Helper()

	cmd := exec.Command("git", "merge", branch)
	output, err := cmd.CombinedOutput()

	if err != nil {
		// Check if it's a merge conflict (expected in some tests)
		if strings.Contains(string(output), "CONFLICT") {
			return false, nil // No error, but conflicts exist
		}
		return false, err
	}

	return true, nil // Success, no conflicts
}

// GetConflictedFiles returns list of files with merge conflicts
func GetConflictedFiles(t *testing.T) []string {
	t.Helper()

	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	output, err := cmd.Output()
	if err != nil {
		return []string{}
	}

	files := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(files) == 1 && files[0] == "" {
		return []string{}
	}

	return files
}

// HasConflictMarkers checks if a file contains conflict markers
func HasConflictMarkers(t *testing.T, filename string) bool {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filename, err)
	}

	str := string(content)
	return strings.Contains(str, "<<<<<<<") &&
		strings.Contains(str, "=======") &&
		strings.Contains(str, ">>>>>>>")
}

// ReadFile reads a file and returns its content
func ReadFile(t *testing.T, filename string) string {
	t.Helper()

	content, err := os.ReadFile(filename)
	if err != nil {
		t.Fatalf("Failed to read file %s: %v", filename, err)
	}

	return string(content)
}

// WriteFile writes content to a file
func WriteFile(t *testing.T, filename, content string) {
	t.Helper()

	dir := filepath.Dir(filename)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("Failed to create directory %s: %v", dir, err)
		}
	}

	if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
		t.Fatalf("Failed to write file %s: %v", filename, err)
	}
}

// StageFile stages a file for commit
func StageFile(t *testing.T, filename string) {
	t.Helper()
	runCommand(t, "git", "add", filename)
}

// GetCurrentCommit returns the current commit SHA
func GetCurrentCommit(t *testing.T) string {
	t.Helper()
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		// No commits yet
		return ""
	}
	return strings.TrimSpace(string(output))
}

// GetCurrentBranch returns the current branch name
func GetCurrentBranch(t *testing.T) string {
	t.Helper()
	output := runCommand(t, "git", "branch", "--show-current")
	return strings.TrimSpace(output)
}

// InitializeGip initializes Gip in the current directory
func InitializeGip(t *testing.T, gipBinary string) {
	t.Helper()

	// Copy gip binary to current directory if path provided
	if gipBinary != "" && gipBinary != "gip" {
		content, err := os.ReadFile(gipBinary)
		if err != nil {
			t.Fatalf("Failed to read gip binary: %v", err)
		}
		if err := os.WriteFile("gip.exe", content, 0755); err != nil {
			t.Fatalf("Failed to copy gip binary: %v", err)
		}
	}

	// Run gip init
	runCommand(t, "gip", "init")
}

// runCommand executes a command and returns output
func runCommand(t *testing.T, command string, args ...string) string {
	t.Helper()

	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Logf("Command failed: %s %v", command, args)
		t.Logf("Output: %s", string(output))
		t.Logf("Error: %v", err)
		// Don't fail immediately - let caller handle errors
	}

	return string(output)
}

// AssertFileContains checks if a file contains expected text
func AssertFileContains(t *testing.T, filename, expected string) {
	t.Helper()

	content := ReadFile(t, filename)
	if !strings.Contains(content, expected) {
		t.Errorf("File %s does not contain expected text: %s", filename, expected)
		t.Logf("Actual content:\n%s", content)
	}
}

// AssertFileNotContains checks if a file does NOT contain text
func AssertFileNotContains(t *testing.T, filename, unexpected string) {
	t.Helper()

	content := ReadFile(t, filename)
	if strings.Contains(content, unexpected) {
		t.Errorf("File %s should not contain text: %s", filename, unexpected)
	}
}

// AssertNoError checks that error is nil
func AssertNoError(t *testing.T, err error, message string) {
	t.Helper()

	if err != nil {
		t.Fatalf("%s: %v", message, err)
	}
}

// AssertError checks that error is not nil
func AssertError(t *testing.T, err error, message string) {
	t.Helper()

	if err == nil {
		t.Fatalf("%s: expected error but got nil", message)
	}
}
