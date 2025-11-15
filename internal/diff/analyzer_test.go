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

func TestParseGitDiff(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		diffOutput    string
		expectedCount int
		expectedFile  string
		expectedSymbol string
		expectedHunk  string
	}{
		{
			name: "Single Python function modification",
			diffOutput: `diff --git a/src/cart.py b/src/cart.py
index 1234567..abcdefg 100644
--- a/src/cart.py
+++ b/src/cart.py
@@ -42,7 +42,10 @@ def calculateTotal(items):
-def calculateTotal(items):
-    return sum(items)
+def calculateTotal(items, discount=0):
+    total = sum(items)
+    return total * (1 - discount)
`,
			expectedCount:  1,
			expectedFile:   "src/cart.py",
			expectedSymbol: "calculateTotal",
			expectedHunk:   "H#42",
		},
		{
			name: "Multiple hunks in same file",
			diffOutput: `diff --git a/src/user.py b/src/user.py
index aaa..bbb 100644
--- a/src/user.py
+++ b/src/user.py
@@ -10,3 +10,5 @@ def parseUser(data):
+def parseUser(data, strict=False):
+    if strict:
+        validate(data)
     return User(data)
@@ -50,2 +52,4 @@ def validateUser(user):
+def validateUser(user, full=True):
+    if full:
+        checkAll(user)
     return True
`,
			expectedCount: 2,
			expectedFile:  "src/user.py",
		},
		{
			name: "Go function with context",
			diffOutput: `diff --git a/internal/api/handler.go b/internal/api/handler.go
index xxx..yyy 100644
--- a/internal/api/handler.go
+++ b/internal/api/handler.go
@@ -100,5 +100,8 @@ func HandleRequest(w http.ResponseWriter, r *http.Request) {
+	if err := validateRequest(r); err != nil {
+		http.Error(w, err.Error(), http.StatusBadRequest)
+		return
+	}
 	processRequest(r)
`,
			expectedCount:  1,
			expectedFile:   "internal/api/handler.go",
			expectedSymbol: "HandleRequest",
			expectedHunk:   "H#100",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			changes, err := ParseGitDiff(tt.diffOutput)
			if err != nil {
				t.Fatalf("ParseGitDiff failed: %v", err)
			}

			if len(changes) != tt.expectedCount {
				t.Errorf("Expected %d changes, got %d", tt.expectedCount, len(changes))
			}

			if len(changes) > 0 {
				if tt.expectedFile != "" && changes[0].File != tt.expectedFile {
					t.Errorf("Expected file '%s', got '%s'", tt.expectedFile, changes[0].File)
				}
				if tt.expectedSymbol != "" && changes[0].Symbol != tt.expectedSymbol {
					t.Errorf("Expected symbol '%s', got '%s'", tt.expectedSymbol, changes[0].Symbol)
				}
				if tt.expectedHunk != "" && changes[0].HunkID != tt.expectedHunk {
					t.Errorf("Expected hunk '%s', got '%s'", tt.expectedHunk, changes[0].HunkID)
				}
			}
		})
	}
}

func TestExtractSymbolFromContext(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		context  string
		expected string
	}{
		{
			name:     "Python function",
			context:  "def calculateTotal(items)",
			expected: "calculateTotal",
		},
		{
			name:     "Go function",
			context:  "func HandleRequest",
			expected: "HandleRequest",
		},
		{
			name:     "JavaScript function",
			context:  "function processData()",
			expected: "processData",
		},
		{
			name:     "Class method",
			context:  "Calculator::add",
			expected: "add",
		},
		{
			name:     "Simple name",
			context:  "myFunction",
			expected: "myFunction",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractSymbolFromContext(tt.context)
			if result != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}
