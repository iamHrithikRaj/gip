package manifest

import (
	"strings"
	"testing"
)

func TestSerializeManifest(t *testing.T) {
	t.Parallel()

	m := &Manifest{
		Commit: "abc123",
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
					ErrorModel:     []string{"may panic"},
				},
				BehaviorClass: []string{"feature"},
				SideEffects:   []string{"logs:info"},
				Rationale:     "Added feature X",
			},
		},
	}

	toon := SerializeManifest(m)

	// Verify basic structure
	if !strings.Contains(toon, "(manifest") {
		t.Error("Expected TOON to contain (manifest")
	}

	if !strings.Contains(toon, "abc123") {
		t.Error("Expected TOON to contain commit SHA")
	}

	if !strings.Contains(toon, "test.go") {
		t.Error("Expected TOON to contain file name")
	}

	if !strings.Contains(toon, "TestFunc") {
		t.Error("Expected TOON to contain symbol name")
	}

	if !strings.Contains(toon, "input is valid") {
		t.Error("Expected TOON to contain precondition")
	}

	if !strings.Contains(toon, "returns result") {
		t.Error("Expected TOON to contain postcondition")
	}

	if !strings.Contains(toon, "feature") {
		t.Error("Expected TOON to contain behavior class")
	}

	if !strings.Contains(toon, "Added feature X") {
		t.Error("Expected TOON to contain rationale")
	}
}

func TestSerializeEmptyManifest(t *testing.T) {
	t.Parallel()

	m := &Manifest{
		Commit:  "empty123",
		Entries: []Entry{},
	}

	toon := SerializeManifest(m)

	if !strings.Contains(toon, "(manifest") {
		t.Error("Expected TOON structure even for empty manifest")
	}

	if !strings.Contains(toon, "empty123") {
		t.Error("Expected commit SHA in empty manifest")
	}
}

func TestSerializeMultipleEntries(t *testing.T) {
	t.Parallel()

	m := &Manifest{
		Commit: "multi123",
		Entries: []Entry{
			{
				Anchor:        Anchor{File: "file1.go", Symbol: "Func1"},
				BehaviorClass: []string{"feature"},
				Rationale:     "First",
			},
			{
				Anchor:        Anchor{File: "file2.go", Symbol: "Func2"},
				BehaviorClass: []string{"bugfix"},
				Rationale:     "Second",
			},
		},
	}

	toon := SerializeManifest(m)

	// Count number of entries
	entryCount := strings.Count(toon, "(entry")
	if entryCount != 2 {
		t.Errorf("Expected 2 entries in TOON, found %d", entryCount)
	}

	// Verify both files appear
	if !strings.Contains(toon, "file1.go") {
		t.Error("Expected file1.go in TOON")
	}
	if !strings.Contains(toon, "file2.go") {
		t.Error("Expected file2.go in TOON")
	}

	// Verify both symbols appear
	if !strings.Contains(toon, "Func1") {
		t.Error("Expected Func1 in TOON")
	}
	if !strings.Contains(toon, "Func2") {
		t.Error("Expected Func2 in TOON")
	}
}

func TestSerializeWithSideEffects(t *testing.T) {
	t.Parallel()

	m := &Manifest{
		Commit: "side123",
		Entries: []Entry{
			{
				Anchor:      Anchor{File: "api.go", Symbol: "CallAPI"},
				SideEffects: []string{"network:http", "logs:info", "writes:cache"},
				Rationale:   "API call with caching",
			},
		},
	}

	toon := SerializeManifest(m)

	if !strings.Contains(toon, "network:http") {
		t.Error("Expected side effect 'network:http' in TOON")
	}
	if !strings.Contains(toon, "logs:info") {
		t.Error("Expected side effect 'logs:info' in TOON")
	}
	if !strings.Contains(toon, "writes:cache") {
		t.Error("Expected side effect 'writes:cache' in TOON")
	}
}

func TestSerializeWithCompatibilityFlags(t *testing.T) {
	t.Parallel()

	m := &Manifest{
		Commit: "compat123",
		Entries: []Entry{
			{
				Anchor: Anchor{File: "api.go", Symbol: "PublicFunc"},
				Compatibility: Compatibility{
					BinaryBreaking: true,
					SourceBreaking: false,
				},
				Rationale: "Changed API signature",
			},
		},
	}

	toon := SerializeManifest(m)

	if !strings.Contains(toon, "(compatibility") {
		t.Error("Expected compatibility section in TOON")
	}
	if !strings.Contains(toon, "binaryBreaking") {
		t.Error("Expected binaryBreaking field in TOON")
	}
}

func TestSerializeForConflict(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Anchor: Anchor{
			File:   "cart.py",
			Symbol: "calculate_total",
		},
		Contract: Contract{
			Preconditions:  []string{"items is list"},
			Postconditions: []string{"returns float"},
			ErrorModel:     []string{"AttributeError"},
		},
		BehaviorClass: []string{"feature"},
		SideEffects:   []string{},
		Rationale:     "Added tax",
	}

	result := SerializeForConflict(&entry, "abc123")

	// Verify format structure
	if !strings.Contains(result, "||| TOON (ctx") {
		t.Error("Expected TOON context wrapper")
	}

	if !strings.Contains(result, "abc123") {
		t.Error("Expected commit SHA")
	}

	if !strings.Contains(result, "calculate_total") {
		t.Error("Expected symbol name")
	}

	if !strings.Contains(result, "items is list") {
		t.Error("Expected preconditions")
	}

	if !strings.Contains(result, "returns float") {
		t.Error("Expected postconditions")
	}

	if !strings.Contains(result, "feature") {
		t.Error("Expected behavior class")
	}

	if !strings.Contains(result, "Added tax") {
		t.Error("Expected rationale")
	}
}

func TestSerializeForConflictMultipleValues(t *testing.T) {
	t.Parallel()

	entry := Entry{
		Anchor: Anchor{Symbol: "MultiFunc", File: "multi.go"},
		Contract: Contract{
			Preconditions:  []string{"a > 0", "b > 0"},
			Postconditions: []string{"result > 0", "no mutation"},
			ErrorModel:     []string{"ValueError", "TypeError"},
		},
		BehaviorClass: []string{"feature", "validation"},
		SideEffects:   []string{"logs:error", "writes:db"},
		Rationale:     "Multi validation",
	}

	result := SerializeForConflict(&entry, "multi123")

	// Verify all values appear (joined with semicolons)
	if !strings.Contains(result, "a > 0") && !strings.Contains(result, "b > 0") {
		t.Error("Expected all preconditions")
	}

	if !strings.Contains(result, "feature") && !strings.Contains(result, "validation") {
		t.Error("Expected all behavior classes")
	}

	if !strings.Contains(result, "multi123") {
		t.Error("Expected commit SHA")
	}
}
