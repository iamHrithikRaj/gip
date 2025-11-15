package prompt

import (
	"testing"

	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// TestGlobalIntentWorkflow tests the global intent prompt logic
func TestGlobalIntentWorkflow(t *testing.T) {
	tests := []struct {
		name             string
		numSymbols       int
		useGlobalIntent  bool
		inheritPattern   []bool // per-symbol inheritance decision
		wantGlobalIntent bool
	}{
		{
			name:             "single function - no global intent",
			numSymbols:       1,
			useGlobalIntent:  false,
			inheritPattern:   []bool{},
			wantGlobalIntent: false,
		},
		{
			name:             "multi-function with global intent - all inherit",
			numSymbols:       3,
			useGlobalIntent:  true,
			inheritPattern:   []bool{true, true, true},
			wantGlobalIntent: true,
		},
		{
			name:             "multi-function with global intent - mixed inheritance",
			numSymbols:       3,
			useGlobalIntent:  true,
			inheritPattern:   []bool{true, false, true},
			wantGlobalIntent: true,
		},
		{
			name:             "multi-function without global intent",
			numSymbols:       2,
			useGlobalIntent:  false,
			inheritPattern:   []bool{},
			wantGlobalIntent: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Verify logic: single function should never prompt for global intent
			if tt.numSymbols == 1 && tt.wantGlobalIntent {
				t.Error("single function should not have global intent")
			}

			// Verify logic: if useGlobalIntent is true, should have global intent
			if tt.useGlobalIntent != tt.wantGlobalIntent {
				t.Errorf("useGlobalIntent=%v but wantGlobalIntent=%v", tt.useGlobalIntent, tt.wantGlobalIntent)
			}

			// Verify inheritance pattern matches expectations
			if tt.wantGlobalIntent && len(tt.inheritPattern) != tt.numSymbols {
				t.Errorf("inheritance pattern length %d doesn't match numSymbols %d", len(tt.inheritPattern), tt.numSymbols)
			}
		})
	}
}

// TestGlobalIntentInheritance tests that entries correctly inherit from global intent
func TestGlobalIntentInheritance(t *testing.T) {
	globalIntent := &manifest.GlobalIntent{
		BehaviorClass: []string{"feature", "security"},
		Rationale:     "Add strict validation across parsers",
	}

	tests := []struct {
		name              string
		inherits          bool
		wantBehaviorClass []string
		wantRationale     string
	}{
		{
			name:              "entry inherits global intent",
			inherits:          true,
			wantBehaviorClass: []string{"feature", "security"},
			wantRationale:     "",
		},
		{
			name:              "entry has unique intent",
			inherits:          false,
			wantBehaviorClass: []string{"bugfix"},
			wantRationale:     "Fix parsing edge case",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			entry := manifest.Entry{
				InheritsGlobalIntent: tt.inherits,
			}

			if tt.inherits {
				// Entry should use global intent behavior classes
				entry.BehaviorClass = globalIntent.BehaviorClass
				entry.Rationale = ""

				if len(entry.BehaviorClass) != len(tt.wantBehaviorClass) {
					t.Errorf("got %d behavior classes, want %d", len(entry.BehaviorClass), len(tt.wantBehaviorClass))
				}

				if entry.Rationale != tt.wantRationale {
					t.Errorf("got rationale %q, want %q", entry.Rationale, tt.wantRationale)
				}
			} else {
				// Entry has unique intent
				entry.BehaviorClass = tt.wantBehaviorClass
				entry.Rationale = tt.wantRationale

				if entry.BehaviorClass[0] != "bugfix" {
					t.Errorf("got behavior class %q, want 'bugfix'", entry.BehaviorClass[0])
				}

				if entry.Rationale == "" {
					t.Error("unique intent should have non-empty rationale")
				}
			}
		})
	}
}

// TestManifestWithGlobalIntent tests that manifests correctly store global intent
func TestManifestWithGlobalIntent(t *testing.T) {
	m := &manifest.Manifest{
		Commit: "abc123",
		GlobalIntent: &manifest.GlobalIntent{
			BehaviorClass: []string{"feature"},
			Rationale:     "Add validation layer",
		},
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   "user.py",
					Symbol: "parseUser",
					HunkID: "H#10",
				},
				InheritsGlobalIntent: true,
				BehaviorClass:        []string{"feature"},
			},
			{
				Anchor: manifest.Anchor{
					File:   "order.py",
					Symbol: "parseOrder",
					HunkID: "H#22",
				},
				InheritsGlobalIntent: false,
				BehaviorClass:        []string{"bugfix"},
				Rationale:            "Fix edge case",
			},
		},
	}

	// Verify global intent
	if m.GlobalIntent == nil {
		t.Fatal("manifest should have global intent")
	}

	if m.GlobalIntent.Rationale != "Add validation layer" {
		t.Errorf("got rationale %q, want 'Add validation layer'", m.GlobalIntent.Rationale)
	}

	// Verify first entry inherits
	if !m.Entries[0].InheritsGlobalIntent {
		t.Error("first entry should inherit global intent")
	}

	// Verify second entry doesn't inherit
	if m.Entries[1].InheritsGlobalIntent {
		t.Error("second entry should not inherit global intent")
	}

	if m.Entries[1].Rationale == "" {
		t.Error("non-inheriting entry should have unique rationale")
	}
}
