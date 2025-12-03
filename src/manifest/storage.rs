//! Manifest storage operations - saving and loading manifests using Git Notes
//!
//! Manifests are stored as JSON in the custom git ref `refs/notes/gip`.
//! This allows manifests to be shared across the team when pushing/pulling.

use crate::git;
use crate::manifest::types::*;
use anyhow::{Context, Result};
use std::fs;
use std::path::Path;

/// Save writes a manifest to Git Notes
pub fn save(manifest: &Manifest, commit_sha: &str, repo_path: Option<&Path>) -> Result<()> {
    // Serialize as JSON
    let json = serde_json::to_string_pretty(manifest).context("Failed to serialize manifest")?;

    // Write to Git Notes
    git::add_note(commit_sha, &json, repo_path).context("Failed to save manifest to git notes")?;

    Ok(())
}

/// Load reads a manifest from Git Notes
pub fn load(commit_sha: &str, repo_path: Option<&Path>) -> Result<Manifest> {
    // Read from Git Notes
    let data =
        git::get_note(commit_sha, repo_path).context("Failed to read manifest from git notes")?;

    // Parse JSON
    let mut manifest: Manifest =
        serde_json::from_str(&data).context("Failed to parse manifest JSON")?;

    // Migrate v1.0 → v2.0 if needed
    if manifest.schema_version.is_empty() || manifest.schema_version == SCHEMA_VERSION_1_0 {
        manifest = migrate_v1_to_v2(manifest);
    }

    Ok(manifest)
}

/// SavePending saves a manifest as pending (before commit)
pub fn save_pending(manifest: &Manifest, gip_dir: &Path) -> Result<()> {
    // Ensure .gip directory exists
    fs::create_dir_all(gip_dir).context("Failed to create .gip directory")?;

    let path = gip_dir.join("pending.json");

    // Serialize as JSON
    let json =
        serde_json::to_string_pretty(manifest).context("Failed to serialize pending manifest")?;

    // Write to file
    fs::write(&path, json)
        .with_context(|| format!("Failed to write pending manifest to {:?}", path))?;

    Ok(())
}

/// LoadPending loads the pending manifest
pub fn load_pending(gip_dir: &Path) -> Result<Manifest> {
    let path = gip_dir.join("pending.json");

    let data = fs::read_to_string(&path)
        .with_context(|| format!("Failed to read pending manifest from {:?}", path))?;

    let manifest: Manifest =
        serde_json::from_str(&data).context("Failed to parse pending manifest")?;

    Ok(manifest)
}

/// Migrate v1.0 manifest to v2.0 format
pub fn migrate_v1_to_v2(mut manifest: Manifest) -> Manifest {
    // Update schema version
    manifest.schema_version = SCHEMA_VERSION_2_0.to_string();

    // Migrate compatibility flags if needed
    for entry in &mut manifest.entries {
        if let Some(ref mut compat) = entry.compatibility {
            // If old v1.0 fields are set but new fields aren't, migrate them
            if compat.deprecations.is_none() {
                compat.deprecations = Some(Vec::new());
            }
            if compat.migrations.is_none() {
                compat.migrations = Some(Vec::new());
            }

            // Map old breaking change flags to new format
            if compat.binary_breaking == Some(true) || compat.source_breaking == Some(true) {
                compat.breaking = true;
            }
        }

        // Ensure hunk_id exists (v2.0 feature)
        if entry.anchor.hunk_id.is_empty() {
            entry.anchor.hunk_id = "H#0".to_string();
        }
    }

    manifest
}

#[cfg(test)]
mod tests {
    use super::*;
    use tempfile::TempDir;

    fn create_test_manifest() -> Manifest {
        Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "abc123def456".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/main.rs".to_string(),
                    symbol: "main".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: CHANGE_MODIFY.to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec![],
                    postconditions: vec![],
                    error_model: vec![],
                },
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                side_effects: vec![],
                compatibility: None,
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                rationale: "Test rationale".to_string(),
                inherits_global_intent: None,
            }],
        }
    }

    #[test]
    fn test_save_and_load_pending() {
        let temp_dir = TempDir::new().unwrap();
        let gip_dir = temp_dir.path();

        let manifest = create_test_manifest();

        // Save pending
        save_pending(&manifest, gip_dir).unwrap();

        // Verify file exists
        assert!(gip_dir.join("pending.json").exists());

        // Load pending
        let loaded = load_pending(gip_dir).unwrap();
        assert_eq!(loaded, manifest);
    }

    #[test]
    fn test_migrate_v1_to_v2() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_1_0.to_string(),
            commit: "old123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "old.rs".to_string(),
                    symbol: "old_fn".to_string(),
                    hunk_id: "".to_string(), // Empty in v1.0
                },
                change_type: CHANGE_MODIFY.to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec![],
                    postconditions: vec![],
                    error_model: vec![],
                },
                behavior_class: vec![BEHAVIOR_BUGFIX.to_string()],
                side_effects: vec![],
                compatibility: Some(Compatibility {
                    breaking: false,
                    deprecations: None,
                    migrations: None,
                    binary_breaking: Some(true),
                    source_breaking: Some(false),
                    data_model_migration: None,
                }),
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                rationale: "".to_string(),
                inherits_global_intent: None,
            }],
        };

        let migrated = migrate_v1_to_v2(manifest);

        // Check schema version updated
        assert_eq!(migrated.schema_version, SCHEMA_VERSION_2_0);

        // Check hunk_id was set
        assert_eq!(migrated.entries[0].anchor.hunk_id, "H#0");

        // Check compatibility fields migrated
        let compat = migrated.entries[0].compatibility.as_ref().unwrap();
        assert!(compat.breaking); // binary_breaking=true should set breaking=true
        assert!(compat.deprecations.is_some());
        assert!(compat.migrations.is_some());
    }
}
