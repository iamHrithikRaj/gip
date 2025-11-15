// Package manifest provides functionality for creating, storing, and loading
// Gip manifests that capture structured context about code changes.
//
// Manifests are stored as JSON files in .gip/manifest/<sha>.json and can be
// serialized to TOON format for display in conflict markers during merge operations.
package manifest

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/iamHrithikRaj/gip/internal/git"
)

// Save writes a manifest to disk as JSON (easy to parse)
func Save(m *Manifest, commitSHA string) error {
	// Ensure .gip/manifest directory exists
	if err := git.EnsureGipDir(); err != nil {
		return fmt.Errorf("failed to create .gip directory: %w", err)
	}

	// Set schema version if not already set (v2.0)
	if m.SchemaVersion == "" {
		m.SchemaVersion = SchemaVersionCurrent
	}

	// Get manifest path (use .json extension)
	manifestDir, err := git.GetManifestDir()
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/%s.json", manifestDir, commitSHA)

	// Serialize as JSON for easy parsing
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write manifest: %w", err)
	}

	return nil
}

// Load reads a manifest from disk and automatically migrates v1.0 → v2.0
func Load(commitSHA string) (*Manifest, error) {
	// Try .json first (new format)
	manifestDir, err := git.GetManifestDir()
	if err != nil {
		return nil, err
	}

	jsonPath := fmt.Sprintf("%s/%s.json", manifestDir, commitSHA)
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		// Fallback to .toon for backward compatibility
		toonPath := fmt.Sprintf("%s/%s.toon", manifestDir, commitSHA)
		data, err = os.ReadFile(toonPath)
		if err != nil {
			return nil, fmt.Errorf("manifest not found for commit %s: %w", commitSHA, err)
		}
		// For .toon files, parse as JSON (they should be JSON now)
	}

	// Parse JSON
	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return nil, fmt.Errorf("failed to parse manifest: %w", err)
	}

	// Migrate v1.0 → v2.0 if needed
	if manifest.SchemaVersion == "" || manifest.SchemaVersion == SchemaVersion10 {
		manifest = migrateV1ToV2(manifest)
	}

	return &manifest, nil
}

// SavePending saves a manifest as pending (before commit)
func SavePending(m *Manifest) error {
	if err := git.EnsureGipDir(); err != nil {
		return err
	}

	// Set schema version if not already set (v2.0)
	if m.SchemaVersion == "" {
		m.SchemaVersion = SchemaVersionCurrent
	}

	gipDir, err := git.GetGipDir()
	if err != nil {
		return err
	}

	// Serialize as JSON
	jsonData, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize manifest: %w", err)
	}

	path := fmt.Sprintf("%s/pending-manifest.json", gipDir)
	return os.WriteFile(path, jsonData, 0644)
}

// Exists checks if a manifest exists for a commit
func Exists(commitSHA string) bool {
	path, err := git.GetManifestPath(commitSHA)
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}

// migrateV1ToV2 upgrades a v1.0 manifest to v2.0 schema
func migrateV1ToV2(v1 Manifest) Manifest {
	// Set schema version
	v1.SchemaVersion = SchemaVersion20

	// Migrate each entry
	for i := range v1.Entries {
		entry := &v1.Entries[i]

		// v1.0 used HunkID, v2.0 uses Anchor.HunkID (already compatible)
		// If hunk is empty, set default
		if entry.Anchor.HunkID == "" {
			entry.Anchor.HunkID = "H#0" // Unknown hunk
		}

		// If changeType is empty, assume "modify"
		if entry.ChangeType == "" {
			entry.ChangeType = ChangeModify
		}

		// Migrate old Compatibility fields to new format
		if entry.Compatibility.BinaryBreaking || entry.Compatibility.SourceBreaking {
			entry.Compatibility.Breaking = true
		}

		// No globalIntent in v1.0, so inheritsGlobalIntent stays false
		entry.InheritsGlobalIntent = false
	}

	return v1
}
