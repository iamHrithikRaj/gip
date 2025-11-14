package helpers

import (
	"github.com/hrithikraj/gip/internal/manifest"
)

// NewTestManifest creates a basic manifest for testing
func NewTestManifest(commit, file, symbol, rationale string) *manifest.Manifest {
	return &manifest.Manifest{
		Commit: commit,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   file,
					Symbol: symbol,
					HunkID: "1",
				},
				ChangeType: manifest.ChangeModify,
				Contract: manifest.Contract{
					Preconditions:  []string{"valid input"},
					Postconditions: []string{"returns result"},
					ErrorModel:     []string{"none"},
				},
				BehaviorClass: []string{"feature"},
				SideEffects:   []string{"none"},
				Rationale:     rationale,
			},
		},
	}
}

// NewBugfixManifest creates a bugfix manifest
func NewBugfixManifest(commit, file, symbol string) *manifest.Manifest {
	return &manifest.Manifest{
		Commit: commit,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   file,
					Symbol: symbol,
					HunkID: "1",
				},
				ChangeType: manifest.ChangeModify,
				Contract: manifest.Contract{
					Preconditions:  []string{"input may be invalid"},
					Postconditions: []string{"handles edge cases"},
					ErrorModel:     []string{"returns error on invalid input"},
				},
				BehaviorClass: []string{"bugfix"},
				SideEffects:   []string{"none"},
				Rationale:     "Fixed edge case bug",
			},
		},
	}
}

// NewFeatureManifest creates a feature manifest with side effects
func NewFeatureManifest(commit, file, symbol string, sideEffects []string) *manifest.Manifest {
	return &manifest.Manifest{
		Commit: commit,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   file,
					Symbol: symbol,
					HunkID: "1",
				},
				ChangeType: manifest.ChangeAdd,
				Contract: manifest.Contract{
					Preconditions:  []string{"none"},
					Postconditions: []string{"new functionality available"},
					ErrorModel:     []string{"may throw exceptions"},
				},
				BehaviorClass: []string{"feature"},
				SideEffects:   sideEffects,
				Rationale:     "Added new feature",
			},
		},
	}
}

// NewRefactorManifest creates a refactor manifest
func NewRefactorManifest(commit, file, symbol string) *manifest.Manifest {
	return &manifest.Manifest{
		Commit: commit,
		Entries: []manifest.Entry{
			{
				Anchor: manifest.Anchor{
					File:   file,
					Symbol: symbol,
					HunkID: "1",
				},
				ChangeType: manifest.ChangeModify,
				SignatureDelta: manifest.SignatureDelta{
					Before: "oldSignature()",
					After:  "newSignature()",
				},
				Contract: manifest.Contract{
					Preconditions:  []string{"same as before"},
					Postconditions: []string{"same behavior, cleaner code"},
					ErrorModel:     []string{"same as before"},
				},
				BehaviorClass: []string{"refactor"},
				SideEffects:   []string{"none"},
				Compatibility: manifest.Compatibility{
					BinaryBreaking: true,
					SourceBreaking: false,
				},
				Rationale: "Improved code structure",
			},
		},
	}
}

// NewManifestWithMultipleEntries creates a manifest with multiple entries
func NewManifestWithMultipleEntries(commit string, entries []ManifestEntrySpec) *manifest.Manifest {
	m := &manifest.Manifest{
		Commit:  commit,
		Entries: []manifest.Entry{},
	}
	
	for _, spec := range entries {
		entry := manifest.Entry{
			Anchor: manifest.Anchor{
				File:   spec.File,
				Symbol: spec.Symbol,
				HunkID: spec.HunkID,
			},
			ChangeType: spec.ChangeType,
			Contract: manifest.Contract{
				Preconditions:  spec.Preconditions,
				Postconditions: spec.Postconditions,
				ErrorModel:     spec.ErrorModel,
			},
			BehaviorClass: spec.BehaviorClass,
			SideEffects:   spec.SideEffects,
			Rationale:     spec.Rationale,
		}
		m.Entries = append(m.Entries, entry)
	}
	
	return m
}

// ManifestEntrySpec specifies a manifest entry for testing
type ManifestEntrySpec struct {
	File           string
	Symbol         string
	HunkID         string
	ChangeType     string
	Preconditions  []string
	Postconditions []string
	ErrorModel     []string
	BehaviorClass  []string
	SideEffects    []string
	Rationale      string
}

// DefaultEntrySpec returns a default entry spec
func DefaultEntrySpec(file, symbol string) ManifestEntrySpec {
	return ManifestEntrySpec{
		File:           file,
		Symbol:         symbol,
		HunkID:         "1",
		ChangeType:     manifest.ChangeModify,
		Preconditions:  []string{"valid input"},
		Postconditions: []string{"returns result"},
		ErrorModel:     []string{"none"},
		BehaviorClass:  []string{"feature"},
		SideEffects:    []string{"none"},
		Rationale:      "Test change",
	}
}
