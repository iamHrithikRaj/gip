// Package diff provides functionality for analyzing Git diffs and extracting
// relevant code symbols (functions, classes, methods) from changed code.
//
// The analyzer supports multiple programming languages including Python, Go,
// JavaScript, TypeScript, Java, C/C++, Ruby, Rust, and PHP.
package diff

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// SymbolChange represents a detected change to a symbol
type SymbolChange struct {
	File            string
	Symbol          string
	HunkID          string
	ChangeType      string
	LinesChanged    int
	SignatureBefore string
	SignatureAfter  string
	TestFiles       []string
	FeatureFlags    []string
}

// GetStagedDiff returns the raw git diff output for staged changes
func GetStagedDiff() (string, error) {
	// TODO: Execute `git diff --cached` and return output
	// For now, return empty string (will be implemented with git integration)
	return "", nil
}

// AnalyzeStagedChanges analyzes git staged changes and extracts symbols
func AnalyzeStagedChanges() ([]SymbolChange, error) {
	// TODO: Execute `git diff --cached` and parse output
	// For now, return mock data with proper hunk tracking

	changes := []SymbolChange{
		{
			File:            "src/cart.py",
			Symbol:          "calculateTotal",
			HunkID:          "H#88",  // Line 88 in file
			ChangeType:      manifest.ChangeModify,
			LinesChanged:    12,
			SignatureBefore: "calculateTotal(items)",
			SignatureAfter:  "calculateTotal(items, discount=0)",
			TestFiles:       []string{"tests/test_cart.py"},
			FeatureFlags:    []string{"DISCOUNT_ENABLED"},
		},
	}

	return changes, nil
}

// ParseGitDiff parses git diff output and extracts hunks with line numbers
// Diff format: @@ -oldStart,oldLines +newStart,newLines @@ context
func ParseGitDiff(diffOutput string) ([]SymbolChange, error) {
	changes := []SymbolChange{}
	lines := strings.Split(diffOutput, "\n")
	
	var currentFile string
	var currentHunkStart int
	var hunkLines []string
	
	// Regex for diff header: @@ -10,5 +10,7 @@ function_context
	hunkHeaderRe := regexp.MustCompile(`^@@\s+-\d+,\d+\s+\+(\d+),\d+\s+@@(.*)`)
	fileHeaderRe := regexp.MustCompile(`^\+\+\+\s+b/(.+)`)
	
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		// Detect file being modified
		if matches := fileHeaderRe.FindStringSubmatch(line); len(matches) > 1 {
			currentFile = matches[1]
			continue
		}
		
		// Detect hunk header
		if matches := hunkHeaderRe.FindStringSubmatch(line); len(matches) > 1 {
			// Save previous hunk if exists
			if currentFile != "" && len(hunkLines) > 0 {
				change := extractSymbolFromHunk(currentFile, currentHunkStart, hunkLines)
				if change != nil {
					changes = append(changes, *change)
				}
			}
			
			// Parse new hunk start line
			fmt.Sscanf(matches[1], "%d", &currentHunkStart)
			hunkLines = []string{}
			
			// Extract context hint (function name after @@)
			if len(matches) > 2 && strings.TrimSpace(matches[2]) != "" {
				hunkLines = append(hunkLines, "CONTEXT:"+strings.TrimSpace(matches[2]))
			}
			continue
		}
		
		// Collect hunk lines (skip file headers)
		if currentFile != "" && !strings.HasPrefix(line, "diff --git") && 
		   !strings.HasPrefix(line, "index ") && !strings.HasPrefix(line, "---") {
			hunkLines = append(hunkLines, line)
		}
	}
	
	// Process last hunk
	if currentFile != "" && len(hunkLines) > 0 {
		change := extractSymbolFromHunk(currentFile, currentHunkStart, hunkLines)
		if change != nil {
			changes = append(changes, *change)
		}
	}
	
	return changes, nil
}

// extractSymbolFromHunk extracts symbol information from a diff hunk
func extractSymbolFromHunk(file string, startLine int, hunkLines []string) *SymbolChange {
	if len(hunkLines) == 0 {
		return nil
	}
	
	// Check for context hint (function name from @@ line)
	var contextHint string
	if strings.HasPrefix(hunkLines[0], "CONTEXT:") {
		contextHint = strings.TrimPrefix(hunkLines[0], "CONTEXT:")
		hunkLines = hunkLines[1:]
	}
	
	// Extract symbol from hunk content or context
	symbol := ExtractSymbolFromDiff(strings.Join(hunkLines, "\n"))
	if symbol == "unknown" && contextHint != "" {
		symbol = extractSymbolFromContext(contextHint)
	}
	
	// Detect change type
	changeType := DetectChangeType(hunkLines)
	
	// Extract signatures for modified functions
	var sigBefore, sigAfter string
	if changeType == manifest.ChangeModify {
		for _, line := range hunkLines {
			if strings.HasPrefix(line, "-") && strings.Contains(line, "(") {
				sigBefore = ExtractSignature(strings.TrimPrefix(line, "-"))
			}
			if strings.HasPrefix(line, "+") && strings.Contains(line, "(") {
				sigAfter = ExtractSignature(strings.TrimPrefix(line, "+"))
			}
		}
	}
	
	// Detect feature flags
	fullHunk := strings.Join(hunkLines, "\n")
	flags := DetectFeatureFlags(fullHunk)
	
	// Count changed lines
	linesChanged := 0
	for _, line := range hunkLines {
		if strings.HasPrefix(line, "+") || strings.HasPrefix(line, "-") {
			linesChanged++
		}
	}
	
	return &SymbolChange{
		File:            file,
		Symbol:          symbol,
		HunkID:          fmt.Sprintf("H#%d", startLine),
		ChangeType:      changeType,
		LinesChanged:    linesChanged,
		SignatureBefore: sigBefore,
		SignatureAfter:  sigAfter,
		TestFiles:       DetectTestFiles(file),
		FeatureFlags:    flags,
	}
}

// extractSymbolFromContext extracts symbol name from git diff context hint
// Context format: "function_name" or "ClassName::methodName"
func extractSymbolFromContext(context string) string {
	context = strings.TrimSpace(context)
	
	// Remove common prefixes
	context = strings.TrimPrefix(context, "def ")
	context = strings.TrimPrefix(context, "func ")
	context = strings.TrimPrefix(context, "function ")
	context = strings.TrimPrefix(context, "class ")
	
	// Extract symbol before parenthesis or colon
	if idx := strings.Index(context, "("); idx != -1 {
		return strings.TrimSpace(context[:idx])
	}
	if idx := strings.Index(context, ":"); idx != -1 {
		parts := strings.Split(context, ":")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	
	// Return first word
	parts := strings.Fields(context)
	if len(parts) > 0 {
		return parts[0]
	}
	
	return context
}

// ExtractSymbolFromDiff extracts symbol name from diff hunk
func ExtractSymbolFromDiff(diffHunk string) string {
	// Try to detect function/method/class definitions
	patterns := []string{
		`def\s+(\w+)\s*\(`,      // Python function
		`func\s+(\w+)\s*\(`,     // Go function
		`function\s+(\w+)\s*\(`, // JavaScript function
		`class\s+(\w+)`,         // Class definition
		`(\w+)\s*=\s*function`,  // Variable function assignment
		`const\s+(\w+)\s*=`,     // Const declaration
		`let\s+(\w+)\s*=`,       // Let declaration
		`var\s+(\w+)\s*=`,       // Var declaration
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		if matches := re.FindStringSubmatch(diffHunk); len(matches) > 1 {
			return matches[1]
		}
	}

	return "unknown"
}

// DetectChangeType determines if symbol was added, modified, deleted, or renamed
func DetectChangeType(diffLines []string) string {
	addedLines := 0
	deletedLines := 0

	for _, line := range diffLines {
		if strings.HasPrefix(line, "+") && !strings.HasPrefix(line, "+++") {
			addedLines++
		} else if strings.HasPrefix(line, "-") && !strings.HasPrefix(line, "---") {
			deletedLines++
		}
	}

	if deletedLines == 0 && addedLines > 0 {
		return manifest.ChangeAdd
	} else if addedLines == 0 && deletedLines > 0 {
		return manifest.ChangeDelete
	}

	return manifest.ChangeModify
}

// ExtractSignature attempts to extract function/method signature
func ExtractSignature(code string) string {
	// Normalize whitespace
	signature := strings.TrimSpace(code)

	// Find first line that looks like a signature
	lines := strings.Split(signature, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Check for function-like patterns
		if strings.Contains(line, "(") && strings.Contains(line, ")") {
			// Extract signature up to opening brace or colon
			if idx := strings.Index(line, "{"); idx != -1 {
				return strings.TrimSpace(line[:idx])
			}
			if idx := strings.Index(line, ":"); idx != -1 {
				return strings.TrimSpace(line[:idx])
			}
			return line
		}
	}

	return signature
}

// DetectFeatureFlags scans code for feature flag patterns
func DetectFeatureFlags(code string) []string {
	flags := []string{}

	// Common feature flag patterns
	patterns := []string{
		`FLAG_(\w+)`,
		`ENABLE_(\w+)`,
		`FEATURE_(\w+)`,
		`if\s+(\w+_FLAG)`,
		`featureFlags\[["'](\w+)["']\]`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringSubmatch(code, -1)
		for _, match := range matches {
			if len(match) > 1 {
				flags = append(flags, match[1])
			}
		}
	}

	return uniqueStrings(flags)
}

// DetectTestFiles finds related test files based on filename
func DetectTestFiles(filename string) []string {
	tests := []string{}

	// Extract base name
	base := strings.TrimSuffix(filename, filepath.Ext(filename))
	dir := filepath.Dir(filename)

	// Common test patterns
	testPatterns := []string{
		fmt.Sprintf("test_%s", filepath.Base(base)),
		fmt.Sprintf("%s_test", filepath.Base(base)),
		fmt.Sprintf("Test%s", strings.Title(filepath.Base(base))),
	}

	for _, pattern := range testPatterns {
		// Look in common test directories
		testDirs := []string{
			filepath.Join(dir, "tests"),
			filepath.Join(dir, "test"),
			filepath.Join(dir, "__tests__"),
			"tests",
			"test",
		}

		for _, testDir := range testDirs {
			testPath := filepath.Join(testDir, pattern)
			tests = append(tests, testPath)
		}
	}

	return tests
}

// GroupChangesBySymbol groups multiple hunks that belong to the same symbol
func GroupChangesBySymbol(changes []SymbolChange) map[string][]SymbolChange {
	grouped := make(map[string][]SymbolChange)

	for _, change := range changes {
		key := fmt.Sprintf("%s::%s", change.File, change.Symbol)
		grouped[key] = append(grouped[key], change)
	}

	return grouped
}

// ToManifestEntry converts a SymbolChange to a manifest Entry
func ToManifestEntry(change SymbolChange) manifest.Entry {
	// Detect if change is breaking
	isBreaking := DetectBreakingChange(change.SignatureBefore, change.SignatureAfter)
	
	entry := manifest.Entry{
		Anchor: manifest.Anchor{
			File:   change.File,
			Symbol: change.Symbol,
			HunkID: change.HunkID,
		},
		ChangeType: change.ChangeType,
		SignatureDelta: manifest.SignatureDelta{
			Before: change.SignatureBefore,
			After:  change.SignatureAfter,
		},
		TestsTouched: change.TestFiles,
		FeatureFlags: change.FeatureFlags,
		Compatibility: manifest.Compatibility{
			Breaking: isBreaking,
		},
	}
	
	return entry
}

// DetectBreakingChange determines if signature change is breaking
func DetectBreakingChange(before, after string) bool {
	if before == "" || after == "" {
		return false
	}
	
	// If signatures are identical, not breaking
	if before == after {
		return false
	}
	
	// Check for parameter removals (breaking)
	beforeParams := extractParameters(before)
	afterParams := extractParameters(after)
	
	// If parameters were removed, it's breaking
	if len(beforeParams) > len(afterParams) {
		return true
	}
	
	// Check if required parameters were added (breaking)
	// New params without defaults are breaking
	for i, param := range afterParams {
		if i >= len(beforeParams) {
			// New parameter
			if !hasDefaultValue(param) {
				return true
			}
		}
	}
	
	// Check for return type changes (potentially breaking)
	beforeReturn := extractReturnType(before)
	afterReturn := extractReturnType(after)
	
	if beforeReturn != "" && afterReturn != "" && beforeReturn != afterReturn {
		// Return type changed - consider breaking
		return true
	}
	
	return false
}

// extractParameters extracts parameter list from signature
func extractParameters(signature string) []string {
	// Find content between first ( and last )
	start := strings.Index(signature, "(")
	end := strings.LastIndex(signature, ")")
	
	if start == -1 || end == -1 || start >= end {
		return []string{}
	}
	
	paramStr := signature[start+1 : end]
	if strings.TrimSpace(paramStr) == "" {
		return []string{}
	}
	
	// Split by comma
	params := strings.Split(paramStr, ",")
	result := []string{}
	for _, p := range params {
		p = strings.TrimSpace(p)
		if p != "" && p != "self" && p != "this" {
			result = append(result, p)
		}
	}
	
	return result
}

// hasDefaultValue checks if parameter has a default value
func hasDefaultValue(param string) bool {
	// Check for assignment (e.g., x = 5, y: int = 10)
	if strings.Contains(param, "=") {
		return true
	}
	
	// Check for optional marker (e.g., y?: number)
	if strings.Contains(param, "?:") {
		return true
	}
	
	// Type annotation alone doesn't mean default value
	// "x: int" has no default, but "x: int = 5" does
	return false
}

// extractReturnType extracts return type from signature
func extractReturnType(signature string) string {
	// Python: -> Type
	if idx := strings.Index(signature, "->"); idx != -1 {
		returnType := strings.TrimSpace(signature[idx+2:])
		// Remove trailing colon if present
		returnType = strings.TrimSuffix(returnType, ":")
		return returnType
	}
	
	// Go: func name(params) ReturnType
	// TypeScript: function(params): ReturnType
	if strings.Contains(signature, ")") {
		parts := strings.Split(signature, ")")
		if len(parts) > 1 {
			returnPart := strings.TrimSpace(parts[1])
			returnPart = strings.TrimPrefix(returnPart, ":")
			returnPart = strings.TrimSpace(returnPart)
			return returnPart
		}
	}
	
	return ""
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]bool)
	result := []string{}

	for _, item := range slice {
		if !seen[item] {
			seen[item] = true
			result = append(result, item)
		}
	}

	return result
}
