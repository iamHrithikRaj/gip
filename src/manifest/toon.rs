//! TOON format serialization for manifests
//!
//! Converts manifests to TOON (Token-Oriented Object Notation) format for displaying
//! in Git conflict markers. Uses the official toon-format crate for spec-compliant
//! serialization with optimal token efficiency.

use crate::manifest::types::*;
use anyhow::Result;

/// Serialize a Manifest to TOON format using the official toon-format library
pub fn serialize_manifest_toon(manifest: &Manifest) -> Result<String> {
    use toon_format::encode_default;

    // Use default encoding with key folding for token efficiency
    let toon = encode_default(manifest)?;
    Ok(toon)
}

/// Serialize a Manifest to legacy custom TOON format (for backward compatibility)
///
/// This function maintains the original Gip TOON format for existing manifests.
/// Consider using `serialize_manifest_toon()` for new implementations.
pub fn serialize_manifest(manifest: &Manifest) -> String {
    let mut output = String::new();

    output.push_str("; Gip Manifest\n");
    output.push_str("(manifest\n");
    output.push_str(&format!("  (schemaVersion {})\n", manifest.schema_version));
    output.push_str(&format!("  (commit #{})\n", manifest.commit));

    // Global intent (v2.0)
    if let Some(ref gi) = manifest.global_intent {
        output.push_str("  (globalIntent\n");
        if !gi.behavior_class.is_empty() {
            output.push_str(&format!(
                "    (behaviorClass [ {} ])\n",
                gi.behavior_class.join(" ")
            ));
        }
        if !gi.rationale.is_empty() {
            output.push_str(&format!("    (rationale \"\"\"{}\"\"\")\n", gi.rationale));
        }
        output.push_str("  )\n");
    }

    output.push_str("  (entries\n");

    for entry in &manifest.entries {
        output.push_str("    (entry\n");

        // Anchor
        output.push_str("      (anchor\n");
        output.push_str(&format!("        (file {})\n", entry.anchor.file));
        output.push_str(&format!("        (symbol {})\n", entry.anchor.symbol));
        output.push_str(&format!("        (hunk {}))\n", entry.anchor.hunk_id));

        // Change type
        output.push_str(&format!("      (changeType {})\n", entry.change_type));

        // Signature delta
        if let Some(ref delta) = entry.signature_delta {
            output.push_str("      (signatureDelta\n");
            output.push_str(&format!("        (before {})\n", delta.before));
            output.push_str(&format!("        (after {}))\n", delta.after));
        }

        // Contract
        output.push_str("      (contract\n");
        if !entry.contract.preconditions.is_empty() {
            output.push_str("        (preconditions\n");
            for pre in &entry.contract.preconditions {
                output.push_str(&format!("          [ \"\"\"{}\"\"\" ]\n", pre));
            }
            output.push_str("        )\n");
        }
        if !entry.contract.postconditions.is_empty() {
            output.push_str("        (postconditions\n");
            for post in &entry.contract.postconditions {
                output.push_str(&format!("          [ \"\"\"{}\"\"\" ]\n", post));
            }
            output.push_str("        )\n");
        }
        if !entry.contract.error_model.is_empty() {
            output.push_str("        (errorModel\n");
            for err in &entry.contract.error_model {
                output.push_str(&format!("          [ \"\"\"{}\"\"\" ]\n", err));
            }
            output.push_str("        )\n");
        }
        output.push_str("      )\n");

        // Behavior class
        if !entry.behavior_class.is_empty() {
            output.push_str(&format!(
                "      (behaviorClass [ {} ])\n",
                entry.behavior_class.join(" ")
            ));
        }

        // Side effects
        if !entry.side_effects.is_empty() {
            output.push_str(&format!(
                "      (sideEffects [ {} ])\n",
                entry.side_effects.join(" ")
            ));
        }

        // Compatibility
        if let Some(ref compat) = entry.compatibility {
            output.push_str("      (compatibility\n");
            output.push_str(&format!(
                "        (breaking {})\n",
                if compat.breaking { "true" } else { "false" }
            ));
            if let Some(ref deps) = compat.deprecations {
                if !deps.is_empty() {
                    output.push_str("        (deprecations\n");
                    for dep in deps {
                        output.push_str(&format!("          [ \"\"\"{}\"\"\" ]\n", dep));
                    }
                    output.push_str("        )\n");
                }
            }
            if let Some(ref migs) = compat.migrations {
                if !migs.is_empty() {
                    output.push_str("        (migrations\n");
                    for mig in migs {
                        output.push_str(&format!("          [ \"\"\"{}\"\"\" ]\n", mig));
                    }
                    output.push_str("        )\n");
                }
            }
            output.push_str("      )\n");
        }

        // Tests touched
        if let Some(ref tests) = entry.tests_touched {
            if !tests.is_empty() {
                output.push_str(&format!("      (testsTouched [ {} ])\n", tests.join(" ")));
            }
        }

        // Feature flags
        if let Some(ref flags) = entry.feature_flags {
            if !flags.is_empty() {
                output.push_str(&format!("      (featureFlags [ {} ])\n", flags.join(" ")));
            }
        }

        // Rationale
        if !entry.rationale.is_empty() {
            output.push_str(&format!(
                "      (rationale \"\"\"{}\"\"\")\n",
                entry.rationale
            ));
        }

        // Inherits global intent
        if let Some(inherits) = entry.inherits_global_intent {
            output.push_str(&format!(
                "      (inheritsGlobalIntent {})\n",
                if inherits { "true" } else { "false" }
            ));
        }

        output.push_str("    )\n");
    }

    output.push_str("  )\n");
    output.push_str(")\n");

    output
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_serialize_simple_manifest() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "abc123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/main.rs".to_string(),
                    symbol: "main".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: CHANGE_ADD.to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec!["none".to_string()],
                    postconditions: vec!["program runs".to_string()],
                    error_model: vec![],
                },
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                side_effects: vec![],
                compatibility: None,
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                rationale: "Initial implementation".to_string(),
                inherits_global_intent: None,
            }],
        };

        let toon = serialize_manifest(&manifest);

        // Verify key elements are present
        assert!(toon.contains("; Gip Manifest"));
        assert!(toon.contains("(manifest"));
        assert!(toon.contains("(schemaVersion 2.0)"));
        assert!(toon.contains("(commit #abc123)"));
        assert!(toon.contains("(file src/main.rs)"));
        assert!(toon.contains("(symbol main)"));
        assert!(toon.contains("(hunk H#1)"));
        assert!(toon.contains("(changeType add)"));
        assert!(toon.contains("(behaviorClass [ feature ])"));
        assert!(toon.contains("(rationale \"\"\"Initial implementation\"\"\")"));
    }

    #[test]
    fn test_serialize_with_global_intent() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "xyz789".to_string(),
            global_intent: Some(GlobalIntent {
                behavior_class: vec![BEHAVIOR_REFACTOR.to_string()],
                rationale: "Complete module refactor".to_string(),
            }),
            entries: vec![],
        };

        let toon = serialize_manifest(&manifest);

        assert!(toon.contains("(globalIntent"));
        assert!(toon.contains("(behaviorClass [ refactor ])"));
        assert!(toon.contains("(rationale \"\"\"Complete module refactor\"\"\")"));
    }

    #[test]
    fn test_serialize_with_signature_delta() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "sig123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "lib.rs".to_string(),
                    symbol: "process".to_string(),
                    hunk_id: "H#10".to_string(),
                },
                change_type: CHANGE_MODIFY.to_string(),
                signature_delta: Some(SignatureDelta {
                    before: "fn process(x: i32)".to_string(),
                    after: "fn process(x: i32, y: i32)".to_string(),
                }),
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
                rationale: "".to_string(),
                inherits_global_intent: None,
            }],
        };

        let toon = serialize_manifest(&manifest);

        assert!(toon.contains("(signatureDelta"));
        assert!(toon.contains("(before fn process(x: i32))"));
        assert!(toon.contains("(after fn process(x: i32, y: i32))"));
    }

    #[test]
    fn test_serialize_with_compatibility() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "compat123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "api.rs".to_string(),
                    symbol: "old_api".to_string(),
                    hunk_id: "H#5".to_string(),
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
                compatibility: Some(Compatibility {
                    breaking: true,
                    deprecations: Some(vec!["old parameter removed".to_string()]),
                    migrations: Some(vec!["use new_api instead".to_string()]),
                    binary_breaking: None,
                    source_breaking: None,
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

        let toon = serialize_manifest(&manifest);

        assert!(toon.contains("(compatibility"));
        assert!(toon.contains("(breaking true)"));
        assert!(toon.contains("(deprecations"));
        assert!(toon.contains("old parameter removed"));
        assert!(toon.contains("(migrations"));
        assert!(toon.contains("use new_api instead"));
    }

    #[test]
    fn test_serialize_manifest_toon_format() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "abc123".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/main.rs".to_string(),
                    symbol: "main".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: CHANGE_ADD.to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec!["none".to_string()],
                    postconditions: vec!["program runs".to_string()],
                    error_model: vec![],
                },
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                side_effects: vec![],
                compatibility: None,
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                rationale: "Initial implementation".to_string(),
                inherits_global_intent: None,
            }],
        };

        // Test official TOON format
        let toon = serialize_manifest_toon(&manifest).unwrap();

        // Should contain key fields
        assert!(toon.contains("schemaVersion") || toon.contains("schema_version"));
        assert!(toon.contains("abc123"));

        // TOON format should be valid and parseable
        assert!(!toon.is_empty());

        // Compare with pretty-printed JSON for realistic comparison
        let json_pretty = serde_json::to_string_pretty(&manifest).unwrap();
        println!("TOON size: {} bytes", toon.len());
        println!("JSON (pretty) size: {} bytes", json_pretty.len());

        // TOON should be more compact than pretty-printed JSON
        assert!(
            toon.len() < json_pretty.len(),
            "TOON ({} bytes) should be more compact than pretty JSON ({} bytes)",
            toon.len(),
            json_pretty.len()
        );
    }

    #[test]
    fn test_toon_round_trip() {
        use toon_format::{decode_default, encode_default};

        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "test456".to_string(),
            global_intent: Some(GlobalIntent {
                behavior_class: vec![BEHAVIOR_REFACTOR.to_string()],
                rationale: "Test refactor".to_string(),
            }),
            entries: vec![],
        };

        // Encode to TOON
        let toon = encode_default(&manifest).unwrap();

        // Decode back
        let decoded: Manifest = decode_default(&toon).unwrap();

        // Should match original
        assert_eq!(decoded.commit, manifest.commit);
        assert_eq!(decoded.schema_version, manifest.schema_version);
        assert!(decoded.global_intent.is_some());
    }

    #[test]
    fn test_toon_round_trip_full_entry() {
        use toon_format::{decode_default, encode_default};

        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "HEAD".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/main.rs".to_string(),
                    symbol: "main".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: "modify".to_string(),
                rationale: "Describe your changes here".to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: None,
                    outputs: None,
                    preconditions: vec!["none".to_string()],
                    postconditions: vec!["program_runs".to_string()],
                    error_model: vec!["panic_on_error".to_string()],
                },
                behavior_class: vec!["feature".to_string()],
                side_effects: vec![],
                compatibility: None,
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                inherits_global_intent: None,
            }],
        };

        // Encode to TOON
        let toon = encode_default(&manifest).unwrap();
        println!("TOON:\n{}", toon);

        // Decode back
        let decoded: Manifest = decode_default(&toon).unwrap();

        // Should match original
        assert_eq!(decoded, manifest);
    }

    #[test]
    fn test_serialize_with_all_optional_fields() {
        let manifest = Manifest {
            schema_version: SCHEMA_VERSION_2_0.to_string(),
            commit: "full123".to_string(),
            global_intent: Some(GlobalIntent {
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                rationale: "Global change".to_string(),
            }),
            entries: vec![Entry {
                anchor: Anchor {
                    file: "full.rs".to_string(),
                    symbol: "full_fn".to_string(),
                    hunk_id: "H#99".to_string(),
                },
                change_type: CHANGE_ADD.to_string(),
                signature_delta: Some(SignatureDelta {
                    before: "".to_string(),
                    after: "fn full_fn()".to_string(),
                }),
                contract: Contract {
                    inputs: Some(vec!["a".to_string()]),
                    outputs: Some("b".to_string()),
                    preconditions: vec!["a > 0".to_string()],
                    postconditions: vec!["b > a".to_string()],
                    error_model: vec!["panic".to_string()],
                },
                behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
                side_effects: vec!["logs:stdout".to_string()],
                compatibility: Some(Compatibility {
                    breaking: false,
                    deprecations: None,
                    migrations: None,
                    binary_breaking: None,
                    source_breaking: None,
                    data_model_migration: None,
                }),
                tests_touched: Some(vec!["test.rs".to_string()]),
                perf_budget: None,
                security_notes: None,
                feature_flags: Some(vec!["FLAG_A".to_string()]),
                rationale: "Full entry".to_string(),
                inherits_global_intent: Some(true),
            }],
        };

        let toon = serialize_manifest(&manifest);

        // Check all sections are present
        assert!(toon.contains("(globalIntent"));
        assert!(toon.contains("(signatureDelta"));
        assert!(toon.contains("(compatibility"));
        assert!(toon.contains("(testsTouched"));
        assert!(toon.contains("(featureFlags"));
        assert!(toon.contains("(inheritsGlobalIntent true)"));
    }
}
