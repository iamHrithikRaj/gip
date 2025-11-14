// Package git provides utilities for interacting with Git repositories,
// including executing Git commands, retrieving commit information, and
// configuring GIP's custom merge driver.
//
// This package handles all Git command execution and repository operations
// needed by GIP to capture and store manifest metadata.
package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// IsGitRepo checks if current directory is a Git repository
func IsGitRepo() bool {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	err := cmd.Run()
	return err == nil
}

// GetRepoRoot returns the root directory of the Git repository
func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("not a git repository")
	}
	return strings.TrimSpace(string(output)), nil
}

// GetCurrentCommit returns the current commit SHA
func GetCurrentCommit() (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetStagedDiff returns the diff of staged changes
func GetStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--cached")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

// HasStagedChanges checks if there are staged changes
func HasStagedChanges() bool {
	cmd := exec.Command("git", "diff", "--cached", "--quiet")
	err := cmd.Run()
	return err != nil // Returns true if there are changes (exit code != 0)
}

// GetGipDir returns the .gip directory path
func GetGipDir() (string, error) {
	root, err := GetRepoRoot()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, ".gip"), nil
}

// GetManifestDir returns the manifest storage directory
func GetManifestDir() (string, error) {
	gipDir, err := GetGipDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(gipDir, "manifest"), nil
}

// EnsureGipDir creates .gip directory structure if it doesn't exist
func EnsureGipDir() error {
	manifestDir, err := GetManifestDir()
	if err != nil {
		return err
	}
	return os.MkdirAll(manifestDir, 0755)
}

// GetManifestPath returns the path for a commit's manifest
func GetManifestPath(commitSHA string) (string, error) {
	manifestDir, err := GetManifestDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(manifestDir, fmt.Sprintf("%s.toon", commitSHA)), nil
}

// InstallGitHooks installs GIP git hooks
func InstallGitHooks() error {
	root, err := GetRepoRoot()
	if err != nil {
		return err
	}

	hooksDir := filepath.Join(root, ".git", "hooks")

	// Pre-commit hook
	preCommitPath := filepath.Join(hooksDir, "pre-commit")
	preCommitScript := `#!/bin/sh
# GIP pre-commit hook
# Validates that manifest exists for staged changes

if [ -f .gip/pending-manifest.toon ]; then
    echo "✓ GIP manifest found"
    exit 0
fi

echo "⚠ Warning: No GIP manifest found"
echo "Run: gip commit -c"
exit 0
`

	if err := os.WriteFile(preCommitPath, []byte(preCommitScript), 0755); err != nil {
		return fmt.Errorf("failed to create pre-commit hook: %w", err)
	}

	// Post-commit hook
	postCommitPath := filepath.Join(hooksDir, "post-commit")
	postCommitScript := `#!/bin/sh
# GIP post-commit hook
# Moves pending manifest to commit-specific file

COMMIT=$(git rev-parse HEAD)

if [ -f .gip/pending-manifest.toon ]; then
    mv .gip/pending-manifest.toon .gip/manifest/${COMMIT}.toon
    echo "✓ GIP manifest saved: .gip/manifest/${COMMIT}.toon"
fi
`

	if err := os.WriteFile(postCommitPath, []byte(postCommitScript), 0755); err != nil {
		return fmt.Errorf("failed to create post-commit hook: %w", err)
	}

	return nil
}

// SetupMergeDriver is no longer needed - GIP wraps git merge directly
// Kept for backward compatibility but does nothing
func SetupMergeDriver() error {
	return nil
}
