//! Manifest types and core data structures for Gip.
//!
//! This module defines the schema for storing structured context about code changes,
//! including contracts, behavior classifications, and compatibility information.

use serde::{Deserialize, Serialize};

/// Schema version constants
pub const SCHEMA_VERSION_1_0: &str = "1.0";
pub const SCHEMA_VERSION_2_0: &str = "2.0";
pub const SCHEMA_VERSION_CURRENT: &str = SCHEMA_VERSION_2_0;

/// Behavior class constants
pub const BEHAVIOR_BUGFIX: &str = "bugfix";
pub const BEHAVIOR_FEATURE: &str = "feature";
pub const BEHAVIOR_REFACTOR: &str = "refactor";
pub const BEHAVIOR_PERF: &str = "perf";
pub const BEHAVIOR_SECURITY: &str = "security";
pub const BEHAVIOR_VALIDATION: &str = "validation";
pub const BEHAVIOR_DOCS: &str = "docs";
pub const BEHAVIOR_CONFIG: &str = "config";
pub const BEHAVIOR_MIGRATION: &str = "migration";

/// Change type constants
pub const CHANGE_ADD: &str = "add";
pub const CHANGE_MODIFY: &str = "modify";
pub const CHANGE_DELETE: &str = "delete";
pub const CHANGE_RENAME: &str = "rename";

/// Manifest represents a Gip change manifest for a commit
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct Manifest {
    pub schema_version: String,
    pub commit: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub global_intent: Option<GlobalIntent>,
    pub entries: Vec<Entry>,
}

/// GlobalIntent represents commit-level rationale for multi-function changes (v2.0)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct GlobalIntent {
    pub behavior_class: Vec<String>,
    pub rationale: String,
}

/// Entry represents a single symbol/hunk modification
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct Entry {
    pub anchor: Anchor,
    pub change_type: String,
    pub rationale: String,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub signature_delta: Option<SignatureDelta>,
    pub behavior_class: Vec<String>,
    pub contract: Contract,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub side_effects: Vec<String>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub compatibility: Option<Compatibility>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub tests_touched: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub perf_budget: Option<PerfBudget>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub security_notes: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub feature_flags: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub inherits_global_intent: Option<bool>,
}

/// Anchor identifies the location of the change
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct Anchor {
    pub file: String,
    pub symbol: String,
    pub hunk_id: String,
}

/// SignatureDelta captures API surface changes
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct SignatureDelta {
    pub before: String,
    pub after: String,
}

/// Contract defines the behavioral contract
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct Contract {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub inputs: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub outputs: Option<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub preconditions: Vec<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub postconditions: Vec<String>,
    #[serde(default, skip_serializing_if = "Vec::is_empty")]
    pub error_model: Vec<String>,
}

/// Compatibility flags (v2.0 enhanced)
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct Compatibility {
    pub breaking: bool,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub deprecations: Option<Vec<String>>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub migrations: Option<Vec<String>>,
    // v1.0 compatibility (deprecated but kept for backward compat)
    #[serde(skip_serializing_if = "Option::is_none")]
    pub binary_breaking: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub source_breaking: Option<bool>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub data_model_migration: Option<bool>,
}

/// PerfBudget captures performance expectations
#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
#[serde(rename_all = "camelCase")]
pub struct PerfBudget {
    #[serde(skip_serializing_if = "Option::is_none")]
    pub expected_max_latency_ms: Option<i32>,
    #[serde(skip_serializing_if = "Option::is_none")]
    pub cpu_delta_pct: Option<i32>,
}

impl Manifest {
    /// Creates a new Manifest with the current schema version
    pub fn new(commit: String) -> Self {
        Self {
            schema_version: SCHEMA_VERSION_CURRENT.to_string(),
            commit,
            global_intent: None,
            entries: Vec::new(),
        }
    }

    /// Returns all valid behavior class options
    pub fn all_behavior_classes() -> Vec<&'static str> {
        vec![
            BEHAVIOR_BUGFIX,
            BEHAVIOR_FEATURE,
            BEHAVIOR_REFACTOR,
            BEHAVIOR_PERF,
            BEHAVIOR_SECURITY,
            BEHAVIOR_VALIDATION,
            BEHAVIOR_DOCS,
            BEHAVIOR_CONFIG,
            BEHAVIOR_MIGRATION,
        ]
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_manifest_new() {
        let manifest = Manifest::new("abc123".to_string());
        assert_eq!(manifest.schema_version, SCHEMA_VERSION_CURRENT);
        assert_eq!(manifest.commit, "abc123");
        assert!(manifest.global_intent.is_none());
        assert_eq!(manifest.entries.len(), 0);
    }

    #[test]
    fn test_manifest_serialization() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "test123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/main.rs".to_string(),
                    symbol: "main".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: CHANGE_ADD.to_string(),
                rationale: "Initial implementation".to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec!["none".to_string()],
                    postconditions: vec!["program runs".to_string()],
                    error_model: vec!["panic on error".to_string()],
                },
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                side_effects: vec![],
                compatibility: None,
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                inherits_global_intent: None,
            }],
        };

        // Serialize to JSON
        let json = serde_json::to_string_pretty(&manifest).unwrap();
        assert!(json.contains("\"schemaVersion\": \"2.0\""));
        assert!(json.contains("\"commit\": \"test123\""));

        // Deserialize back
        let deserialized: Manifest = serde_json::from_str(&json).unwrap();
        assert_eq!(deserialized, manifest);
    }

    #[test]
    fn test_manifest_with_global_intent() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "commit789".to_string(),
            global_intent: Some(GlobalIntent {
                behavior_class: vec![BEHAVIOR_REFACTOR.to_string()],
                rationale: "Refactor entire module".to_string(),
            }),
            entries: vec![],
        };

        let json = serde_json::to_string_pretty(&manifest).unwrap();
        let deserialized: Manifest = serde_json::from_str(&json).unwrap();
        assert_eq!(deserialized, manifest);
        assert!(deserialized.global_intent.is_some());
    }

    #[test]
    fn test_entry_with_signature_delta() {
        let entry = Entry {
            anchor: Anchor {
                file: "lib.rs".to_string(),
                symbol: "process".to_string(),
                hunk_id: "H#42".to_string(),
            },
            change_type: CHANGE_MODIFY.to_string(),
            rationale: "Add support for two parameters".to_string(),
            signature_delta: Some(SignatureDelta {
                before: "fn process(x: i32)".to_string(),
                after: "fn process(x: i32, y: i32)".to_string(),
            }),
            contract: Contract {
                inputs: Some(vec!["x: i32".to_string(), "y: i32".to_string()]),
                outputs: Some("i32".to_string()),
                preconditions: vec!["x > 0".to_string()],
                postconditions: vec!["returns sum".to_string()],
                error_model: vec!["none".to_string()],
            },
            behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
            side_effects: vec!["none".to_string()],
            compatibility: Some(Compatibility {
                breaking: true,
                deprecations: Some(vec!["old signature deprecated".to_string()]),
                migrations: Some(vec!["add second parameter".to_string()]),
                binary_breaking: None,
                source_breaking: None,
                data_model_migration: None,
            }),
            tests_touched: Some(vec!["tests/process_test.rs".to_string()]),
            perf_budget: None,
            security_notes: None,
            feature_flags: None,
            inherits_global_intent: Some(false),
        };

        let json = serde_json::to_string_pretty(&entry).unwrap();
        let deserialized: Entry = serde_json::from_str(&json).unwrap();
        assert_eq!(deserialized, entry);
        assert!(deserialized.signature_delta.is_some());
        assert!(deserialized.compatibility.as_ref().unwrap().breaking);
    }

    #[test]
    fn test_all_behavior_classes() {
        let classes = Manifest::all_behavior_classes();
        assert_eq!(classes.len(), 9);
        assert!(classes.contains(&BEHAVIOR_BUGFIX));
        assert!(classes.contains(&BEHAVIOR_FEATURE));
        assert!(classes.contains(&BEHAVIOR_SECURITY));
    }

    #[test]
    fn test_compatibility_backward_compat() {
        // Test v1.0 compatibility fields
        let compat = Compatibility {
            breaking: false,
            deprecations: None,
            migrations: None,
            binary_breaking: Some(true),
            source_breaking: Some(false),
            data_model_migration: Some(true),
        };

        let json = serde_json::to_string(&compat).unwrap();
        let deserialized: Compatibility = serde_json::from_str(&json).unwrap();
        assert_eq!(deserialized.binary_breaking, Some(true));
        assert_eq!(deserialized.source_breaking, Some(false));
    }
}
