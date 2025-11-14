// Package merge implements GIP's custom merge driver that enriches Git conflict
// markers with structured context from GIP manifests.
//
// The merge driver detects conflict markers, loads manifests for each version,
// and converts them to TOON format for inline display in the conflicted file.
package merge

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
	
	"github.com/alpkeskin/gotoon"
	"github.com/hrithikraj/gip/internal/manifest"
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

// enrichConflictBlocks injects custom GIP context into conflict regions
func enrichConflictBlocks(lines []string, conflicts []ConflictBlock, currentSHA, otherSHA string, filePath string) []string {
	result := make([]string, 0)
	lastEnd := 0
	
	for _, conflict := range conflicts {
		// Add lines before conflict (including <<<<<<< marker)
		result = append(result, lines[lastEnd:conflict.StartLine+1]...)
		
		// Add HEAD section
		result = append(result, conflict.HeadLines...)
		
		// Add GIP context for HEAD
		result = append(result, "||| GIP CONTEXT (HEAD - Your changes)")
		shortHead := currentSHA
		if len(currentSHA) > 8 {
			shortHead = currentSHA[:8]
		}
		result = append(result, fmt.Sprintf("||| Commit: %s", shortHead))
		
		// Try to add TOON manifest for HEAD
		headToon := getManifestToon(currentSHA, filePath)
		if headToon != "" {
			for _, line := range strings.Split(headToon, "\n") {
				result = append(result, "||| "+line)
			}
		}
		
		// Add separator
		result = append(result, "=======")
		
		// Add THEIR section
		result = append(result, conflict.TheirLines...)
		
		// Add GIP context for THEIRS
		result = append(result, "||| GIP CONTEXT (MERGE_HEAD - Their changes)")
		shortOther := otherSHA
		if len(otherSHA) > 8 {
			shortOther = otherSHA[:8]
		}
		result = append(result, fmt.Sprintf("||| Commit: %s", shortOther))
		
		// Try to add TOON manifest for THEIRS
		theirToon := getManifestToon(otherSHA, filePath)
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

// getManifestToon loads a manifest and converts it to TOON format
func getManifestToon(commitSHA, filePath string) string {
	// Load manifest
	m, err := manifest.Load(commitSHA)
	if err != nil {
		return "" // No manifest found
	}
	
	// Find entry for this file
	var entry *manifest.Entry
	for i := range m.Entries {
		if m.Entries[i].Anchor.File == filePath {
			entry = &m.Entries[i]
			break
		}
	}
	
	if entry == nil {
		return "" // No entry for this file
	}
	
	// Convert entry to map for gotoon
	data := map[string]interface{}{
		"symbol":         entry.Anchor.Symbol,
		"file":           entry.Anchor.File,
		"changeType":     entry.ChangeType,
		"preconditions":  entry.Contract.Preconditions,
		"postconditions": entry.Contract.Postconditions,
		"errorModel":     entry.Contract.ErrorModel,
		"behaviorClass":  entry.BehaviorClass,
		"sideEffects":    entry.SideEffects,
		"rationale":      entry.Rationale,
	}
	
	// Encode to TOON
	toon, err := gotoon.Encode(data, gotoon.WithIndent(2))
	if err != nil {
		return "" // Encoding failed
	}
	
	return toon
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
