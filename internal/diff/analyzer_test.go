package diff

import (
	"testing"
)

func TestExtractSymbolFromDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		diffHunk string
		expected string
	}{
		{
			name: "Python function",
			diffHunk: `+def calculate_total(items):
+    return sum(items)`,
			expected: "calculate_total",
		},
		{
			name: "Go function",
			diffHunk: `+func ProcessData(input string) error {
+    return nil
+}`,
			expected: "ProcessData",
		},
		{
			name: "JavaScript function",
			diffHunk: `+function handleClick() {
+    console.log('clicked');
+}`,
			expected: "handleClick",
		},
		{
			name: "Python class",
			diffHunk: `+class Calculator:
+    def __init__(self):
+        pass`,
			expected: "Calculator",
		},
		{
			name:     "JavaScript const",
			diffHunk: `+const API_URL = 'https://api.example.com';`,
			expected: "API_URL",
		},
		{
			name: "JavaScript arrow function",
			diffHunk: `+const fetchData = () => {
+    return fetch(API_URL);
+}`,
			expected: "fetchData",
		},
		{
			name: "Go method",
			diffHunk: `+func (c *Calculator) Add(a, b int) int {
+    return a + b
+}`,
			expected: "Add",
		},
		{
			name: "No recognizable symbol",
			diffHunk: `+    x = 42
+    y = 100`,
			expected: "unknown",
		},
		{
			name:     "Empty diff",
			diffHunk: "",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSymbolFromDiff(tt.diffHunk)
			if result != tt.expected {
				t.Errorf("Expected symbol '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDetectChangeType(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		diffLines []string
		expected  string
	}{
		{
			name: "Added lines only",
			diffLines: []string{
				"diff --git a/file.go b/file.go",
				"+func NewFunc() {",
				"+    return nil",
				"+}",
			},
			expected: "add",
		},
		{
			name: "Deleted lines only",
			diffLines: []string{
				"diff --git a/file.go b/file.go",
				"-func OldFunc() {",
				"-    return nil",
				"-}",
			},
			expected: "delete",
		},
		{
			name: "Modified - both added and deleted",
			diffLines: []string{
				"diff --git a/file.go b/file.go",
				"-func Calculate() int {",
				"-    return 42",
				"+func Calculate() int {",
				"+    return 100",
				"}",
			},
			expected: "modify",
		},
		{
			name: "Mostly additions",
			diffLines: []string{
				"-old line",
				"+new line 1",
				"+new line 2",
				"+new line 3",
				"+new line 4",
				"+new line 5",
			},
			expected: "add",
		},
		{
			name: "Mostly deletions",
			diffLines: []string{
				"-old line 1",
				"-old line 2",
				"-old line 3",
				"-old line 4",
				"-old line 5",
				"+new line",
			},
			expected: "delete",
		},
		{
			name:      "Empty diff",
			diffLines: []string{},
			expected:  "modify",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectChangeType(tt.diffLines)
			if result != tt.expected {
				t.Errorf("Expected change type '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestAnalyzeStagedChanges(t *testing.T) {
	t.Parallel()

	// This is a basic smoke test since the actual implementation
	// depends on git environment
	changes, err := AnalyzeStagedChanges()

	if err != nil {
		t.Logf("AnalyzeStagedChanges returned error (may be expected if not in git repo): %v", err)
	}

	// Should return data (mock or real)
	if changes == nil {
		t.Error("Expected non-nil changes slice")
	}

	// If we got changes, verify structure
	if len(changes) > 0 {
		change := changes[0]

		if change.File == "" {
			t.Error("Expected File to be populated")
		}

		if change.Symbol == "" {
			t.Error("Expected Symbol to be populated")
		}

		if change.ChangeType == "" {
			t.Error("Expected ChangeType to be populated")
		}
	}
}

func TestExtractSymbolMultiplePatterns(t *testing.T) {
	t.Parallel()

	// Test that we correctly extract from various language patterns
	diffHunks := []struct {
		code     string
		language string
		expected string
	}{
		{
			code:     "def process_data(input):\n    pass",
			language: "Python",
			expected: "process_data",
		},
		{
			code:     "func ValidateInput(s string) bool {\n    return true\n}",
			language: "Go",
			expected: "ValidateInput",
		},
		{
			code:     "function fetchUsers() {\n    return [];\n}",
			language: "JavaScript",
			expected: "fetchUsers",
		},
		{
			code:     "class UserManager {\n    constructor() {}\n}",
			language: "Class",
			expected: "UserManager",
		},
		{
			code:     "const apiEndpoint = 'https://api.example.com'",
			language: "Constant",
			expected: "apiEndpoint",
		},
		{
			code:     "let counter = 0",
			language: "Variable",
			expected: "counter",
		},
	}

	for _, tc := range diffHunks {
		t.Run(tc.language, func(t *testing.T) {
			result := ExtractSymbolFromDiff(tc.code)
			if result != tc.expected {
				t.Errorf("%s: expected '%s', got '%s'", tc.language, tc.expected, result)
			}
		})
	}
}

func TestExtractSymbolEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		diffHunk string
		expected string
	}{
		{
			name:     "Multiple functions - take first",
			diffHunk: "def first():\n    pass\ndef second():\n    pass",
			expected: "first",
		},
		{
			name:     "Function with complex signature",
			diffHunk: "def complex_func(a, b, c=None, *args, **kwargs):\n    pass",
			expected: "complex_func",
		},
		{
			name:     "Indented function (method)",
			diffHunk: "    def method_name(self):\n        pass",
			expected: "method_name",
		},
		{
			name:     "Function with type hints",
			diffHunk: "def typed_func(x: int) -> str:\n    return str(x)",
			expected: "typed_func",
		},
		{
			name:     "Anonymous function",
			diffHunk: "map(function(x) { return x * 2; }, arr)",
			expected: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractSymbolFromDiff(tt.diffHunk)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestDetectChangeTypeEdgeCases(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			name: "Ignore diff headers",
			lines: []string{
				"--- a/file.go",
				"+++ b/file.go",
				"+added line",
			},
			expected: "add",
		},
		{
			name: "Equal additions and deletions",
			lines: []string{
				"-old1",
				"-old2",
				"+new1",
				"+new2",
			},
			expected: "modify",
		},
		{
			name: "One more addition than deletion",
			lines: []string{
				"-old1",
				"+new1",
				"+new2",
			},
			expected: "modify",
		},
		{
			name: "One more deletion than addition",
			lines: []string{
				"-old1",
				"-old2",
				"+new1",
			},
			expected: "modify",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectChangeType(tt.lines)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
