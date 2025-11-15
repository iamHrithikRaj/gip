// Package merge implements Gip's custom merge driver that enriches Git conflict
// markers with structured context from Gip manifests.
//
// The merge driver detects conflict markers, loads manifests for each version,
// and converts them to TOON format for inline display in the conflicted file.
package merge

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// ConflictBlock represents a Git conflict region
type ConflictBlock struct {
	StartLine  int
	EndLine    int
	HeadLines  []string
	TheirLines []string
}

// EnrichConflicts processes a file with conflicts and injects custom context
func EnrichConflicts(filePath, ancestorSHA, currentSHA, otherSHA string) error {
	// Read the conflicted file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	lines := strings.Split(string(content), "\n")

	// Find all conflict blocks
	conflicts := findConflictBlocks(lines)

	if len(conflicts) == 0 {
		return nil // No conflicts
	}

	fmt.Printf("Found %d conflict(s) in %s\n", len(conflicts), filePath)

	// Enrich each conflict with custom context
	enrichedLines := enrichConflictBlocks(lines, conflicts, currentSHA, otherSHA, filePath)

	// Write back the enriched file
	enrichedContent := strings.Join(enrichedLines, "\n")
	if err := os.WriteFile(filePath, []byte(enrichedContent), 0644); err != nil {
		return fmt.Errorf("failed to write enriched file: %w", err)
	}

	return nil
}

// findConflictBlocks locates all Git conflict markers in the file
func findConflictBlocks(lines []string) []ConflictBlock {
	var conflicts []ConflictBlock
	var current *ConflictBlock
	inHead := false

	conflictStart := regexp.MustCompile(`^<<<<<<<\s+`)
	conflictMiddle := regexp.MustCompile(`^=======\s*$`)
	conflictEnd := regexp.MustCompile(`^>>>>>>>\s+`)

	for i, line := range lines {
		if conflictStart.MatchString(line) {
			current = &ConflictBlock{
				StartLine:  i,
				HeadLines:  []string{},
				TheirLines: []string{},
			}
			inHead = true
		} else if current != nil && conflictMiddle.MatchString(line) {
			inHead = false
		} else if current != nil && conflictEnd.MatchString(line) {
			current.EndLine = i
			conflicts = append(conflicts, *current)
			current = nil
			inHead = false
		} else if current != nil {
			if inHead {
				current.HeadLines = append(current.HeadLines, line)
			} else {
				current.TheirLines = append(current.TheirLines, line)
			}
		}
	}

	return conflicts
}

// extractSymbolFromConflict attempts to detect the function/class being modified
// by analyzing the conflict lines. This enables selective context injection.
// Searches from bottom to top to find the closest (most relevant) symbol.
func extractSymbolFromConflict(conflictLines []string) string {
	// Look for common function/method/class patterns in conflict lines
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`^\s*def\s+(\w+)\s*\(`),           // Python: def functionName(
		regexp.MustCompile(`^\s*func\s+(\w+)\s*\(`),          // Go: func functionName(
		regexp.MustCompile(`^\s*function\s+(\w+)\s*\(`),      // JS: function functionName(
		regexp.MustCompile(`^\s*(\w+)\s*:\s*function\s*\(`),  // JS: name: function(
		regexp.MustCompile(`^\s*(public|private|protected)\s+\w+\s+(\w+)\s*\(`), // Java/C#
		regexp.MustCompile(`^\s*class\s+(\w+)`),              // class ClassName
	}

	// Search from bottom to top to find the closest/most relevant symbol
	for i := len(conflictLines) - 1; i >= 0; i-- {
		line := conflictLines[i]
		for _, pattern := range patterns {
			if matches := pattern.FindStringSubmatch(line); matches != nil {
				if len(matches) > 1 {
					// Return last captured group (symbol name)
					return matches[len(matches)-1]
				}
			}
		}
	}

	return "" // No symbol detected
}

// enrichConflictBlocks injects custom Gip context into conflict regions
// Phase 2: Enhanced with symbol extraction for selective context injection
func enrichConflictBlocks(lines []string, conflicts []ConflictBlock, currentSHA, otherSHA string, filePath string) []string {
	result := make([]string, 0)
	lastEnd := 0

	for _, conflict := range conflicts {
		// Phase 2: Extract symbol from conflict for selective injection
		// Look at conflict lines AND some context before the conflict to find the enclosing function
		contextStart := conflict.StartLine - 10 // Look up to 10 lines back
		if contextStart < 0 {
			contextStart = 0
		}
		contextLines := append(lines[contextStart:conflict.StartLine], conflict.HeadLines...)
		contextLines = append(contextLines, conflict.TheirLines...)
		symbol := extractSymbolFromConflict(contextLines)
		// TODO: Extract hunk ID from conflict line numbers once we integrate with diff analyzer

		// Add lines before conflict (including <<<<<<< marker)
		result = append(result, lines[lastEnd:conflict.StartLine+1]...)

		// Add HEAD section
		result = append(result, conflict.HeadLines...)

		// Add Gip context for HEAD
		result = append(result, "||| Gip CONTEXT (HEAD - Your changes)")
		shortHead := currentSHA
		if len(currentSHA) > 8 {
			shortHead = currentSHA[:8]
		}
		result = append(result, fmt.Sprintf("||| Commit: %s", shortHead))

		// Phase 2: Use selective injection if symbol detected
		headToon := getManifestToonForConflict(currentSHA, filePath, symbol, "")
		if headToon != "" {
			for _, line := range strings.Split(headToon, "\n") {
				result = append(result, "||| "+line)
			}
		}

		// Add separator
		result = append(result, "=======")

		// Add THEIR section
		result = append(result, conflict.TheirLines...)

		// Add Gip context for THEIRS
		result = append(result, "||| Gip CONTEXT (MERGE_HEAD - Their changes)")
		shortOther := otherSHA
		if len(otherSHA) > 8 {
			shortOther = otherSHA[:8]
		}
		result = append(result, fmt.Sprintf("||| Commit: %s", shortOther))

		// Phase 2: Use selective injection if symbol detected
		theirToon := getManifestToonForConflict(otherSHA, filePath, symbol, "")
		if theirToon != "" {
			for _, line := range strings.Split(theirToon, "\n") {
				result = append(result, "||| "+line)
			}
		}

		// Add end marker (>>>>>>> line)
		result = append(result, lines[conflict.EndLine])

		lastEnd = conflict.EndLine + 1
	}

	// Add remaining lines after last conflict
	if lastEnd < len(lines) {
		result = append(result, lines[lastEnd:]...)
	}

	return result
}

// getCommitSHA resolves a Git ref to a commit SHA
func getCommitSHA(ref string) (string, error) {
	cmd := exec.Command("git", "rev-parse", ref)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// findMatchingEntry searches for a manifest entry matching file + symbol + hunk.
// Phase 2: Selective Context Injection - returns only the entry relevant to a specific conflict.
// Matching strategy:
//   1. Match by File + Symbol + HunkID (precise, v2.0 manifests)
//   2. Match by File + Symbol (fallback for v1.0 or missing hunk)
//   3. Match by File only (last resort)
//
// Returns the best matching entry, or nil if no suitable match found.
func findMatchingEntry(m *manifest.Manifest, filePath, symbol, hunkID string) *manifest.Entry {
	if m == nil || len(m.Entries) == 0 {
		return nil
	}

	// Best match: File + Symbol + HunkID (v2.0 selective injection)
	if symbol != "" && hunkID != "" {
		for i := range m.Entries {
			e := &m.Entries[i]
			if e.Anchor.File == filePath &&
				e.Anchor.Symbol == symbol &&
				e.Anchor.HunkID == hunkID {
				return e
			}
		}
	}

	// Good match: File + Symbol (fallback for v1.0 or if hunk doesn't match)
	if symbol != "" {
		for i := range m.Entries {
			e := &m.Entries[i]
			if e.Anchor.File == filePath && e.Anchor.Symbol == symbol {
				return e
			}
		}
	}

	// Last resort: File only (return first entry for this file)
	for i := range m.Entries {
		e := &m.Entries[i]
		if e.Anchor.File == filePath {
			return e
		}
	}

	return nil
}

// getManifestToonForConflict loads a manifest and uses selective context injection.
// Phase 2: Returns ONLY the entry matching the conflict (file + symbol + hunk).
// This replaces the old behavior of injecting entire manifests (100 lines → 8-12 lines).
func getManifestToonForConflict(commitSHA, filePath, symbol, hunkID string) string {
	// Load manifest - try production load
	m, err := manifest.Load(commitSHA)
	if err != nil {
		// If in a test environment without proper git setup, try direct load
		// This allows tests to work without full git repo initialization
		if m, err = loadManifestDirect(commitSHA); err != nil {
			return "" // No manifest found
		}
	}

	// Normalize file path - try to match with both absolute and relative paths
	// This handles cases where manifest uses relative path but conflict has absolute path
	entry := findMatchingEntry(m, filePath, symbol, hunkID)
	if entry == nil {
		// Try with just the base file name (last component of path)
		baseFile := filePath
		if idx := strings.LastIndexAny(filePath, "/\\"); idx != -1 {
			baseFile = filePath[idx+1:]
		}
		entry = findMatchingEntry(m, baseFile, symbol, hunkID)
	}

	if entry == nil {
		return "" // No matching entry
	}

	// Use the manifest's TOON serialization function for consistency
	return manifest.SerializeForConflict(entry, commitSHA)
}

// getManifestToon loads a manifest and converts it to TOON format (legacy v1.0 behavior)
// Deprecated: Use getManifestToonForConflict for Phase 2 selective injection
func getManifestToon(commitSHA, filePath string) string {
	return getManifestToonForConflict(commitSHA, filePath, "", "")
}

// EnrichAllConflicts finds all conflicted files and enriches them
func EnrichAllConflicts() error {
	// Get list of conflicted files
	cmd := exec.Command("git", "diff", "--name-only", "--diff-filter=U")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get conflicted files: %w", err)
	}

	conflictedFiles := strings.Split(strings.TrimSpace(string(output)), "\n")
	if len(conflictedFiles) == 0 || conflictedFiles[0] == "" {
		return nil // No conflicts
	}

	// Get commit SHAs
	currentSHA, err := getCommitSHA("HEAD")
	if err != nil {
		return fmt.Errorf("failed to get HEAD SHA: %w", err)
	}

	otherSHA, err := getCommitSHA("MERGE_HEAD")
	if err != nil {
		return fmt.Errorf("failed to get MERGE_HEAD SHA: %w", err)
	}

	// Enrich each conflicted file
	for _, filePath := range conflictedFiles {
		if filePath == "" {
			continue
		}

		if err := EnrichConflicts(filePath, "", currentSHA, otherSHA); err != nil {
			fmt.Printf("Warning: Failed to enrich %s: %v\n", filePath, err)
		}
	}

	return nil
}

// loadManifestDirect attempts to load a manifest directly from .gip/manifest
// without git validation. This is used as a fallback for test environments.
func loadManifestDirect(commitSHA string) (*manifest.Manifest, error) {
	// Try current directory first
	manifestPath := filepath.Join(".gip", "manifest", commitSHA+".json")
	
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("manifest not found: %w", err)
	}

	var m manifest.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	return &m, nil
}

// DetectConflicts checks if a file has Git conflict markers
func DetectConflicts(filePath string) (bool, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	conflictMarker := regexp.MustCompile(`^<<<<<<<\s+`)

	for scanner.Scan() {
		if conflictMarker.MatchString(scanner.Text()) {
			return true, nil
		}
	}

	return false, scanner.Err()
}
