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
	File           string
	Symbol         string
	HunkID         string
	ChangeType     string
	LinesChanged   int
	SignatureBefore string
	SignatureAfter  string
	TestFiles      []string
	FeatureFlags   []string
}

// AnalyzeStagedChanges analyzes git staged changes and extracts symbols
func AnalyzeStagedChanges() ([]SymbolChange, error) {
	// TODO: Execute `git diff --cached` and parse output
	// For now, return mock data
	
	changes := []SymbolChange{
		{
			File:          "src/cart.py",
			Symbol:        "calculateTotal",
			HunkID:        "H#88",
			ChangeType:    manifest.ChangeModify,
			LinesChanged:  12,
			SignatureBefore: "calculateTotal(items)",
			SignatureAfter:  "calculateTotal(items, discount=0)",
			TestFiles:     []string{"tests/test_cart.py"},
			FeatureFlags:  []string{"DISCOUNT_ENABLED"},
		},
	}
	
	return changes, nil
}

// ExtractSymbolFromDiff extracts symbol name from diff hunk
func ExtractSymbolFromDiff(diffHunk string) string {
	// Try to detect function/method/class definitions
	patterns := []string{
		`def\s+(\w+)\s*\(`,           // Python function
		`func\s+(\w+)\s*\(`,          // Go function
		`function\s+(\w+)\s*\(`,      // JavaScript function
		`class\s+(\w+)`,              // Class definition
		`(\w+)\s*=\s*function`,       // Variable function assignment
		`const\s+(\w+)\s*=`,          // Const declaration
		`let\s+(\w+)\s*=`,            // Let declaration
		`var\s+(\w+)\s*=`,            // Var declaration
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
	return manifest.Entry{
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
	}
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
