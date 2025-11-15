package helpers

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// SaveManifestDirect saves a manifest directly to the test repository
// without requiring git validation. This is for testing purposes only.
func SaveManifestDirect(m *manifest.Manifest, sha string, repoDir string) error {
	manifestDir := filepath.Join(repoDir, ".gip", "manifest")
	
	// Ensure directory exists
	if err := os.MkdirAll(manifestDir, 0755); err != nil {
		return fmt.Errorf("failed to create manifest directory: %w", err)
	}

	// Set schema version if not set
	if m.SchemaVersion == "" {
		m.SchemaVersion = "2.0"
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal manifest: %w", err)
	}

	// Write to file
	manifestPath := filepath.Join(manifestDir, sha+".json")
	if err := os.WriteFile(manifestPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// LoadManifestDirect loads a manifest directly from the test repository
// without requiring git validation. This is for testing purposes only.
func LoadManifestDirect(sha string, repoDir string) (*manifest.Manifest, error) {
	manifestPath := filepath.Join(repoDir, ".gip", "manifest", sha+".json")
	
	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, fmt.Errorf("manifest not found for commit %s: %w", sha, err)
	}

	var m manifest.Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Apply migration if needed (same logic as production)
	if m.SchemaVersion == "" || m.SchemaVersion == "1.0" {
		m = migrateV1ToV2Test(m)
	}

	return &m, nil
}

// migrateV1ToV2Test converts a v1.0 manifest to v2.0 format for testing
func migrateV1ToV2Test(v1 manifest.Manifest) manifest.Manifest {
	v1.SchemaVersion = "2.0"
	
	for i := range v1.Entries {
		// Set default hunk ID if missing
		if v1.Entries[i].Anchor.HunkID == "" {
			v1.Entries[i].Anchor.HunkID = "H#0"
		}
		
		// Set default change type if missing
		if v1.Entries[i].ChangeType == "" {
			v1.Entries[i].ChangeType = "modify"
		}
		
		// Migrate old v1.0 compatibility flags to v2.0 structure
		if v1.Entries[i].Compatibility.BinaryBreaking || v1.Entries[i].Compatibility.SourceBreaking {
			v1.Entries[i].Compatibility.Breaking = true
		}
	}
	
	return v1
}

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
