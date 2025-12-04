//! Merge driver module - handles conflict enrichment
//!
//! Provides functionality for detecting Git conflict markers and injecting
//! structured context from Gip manifests into them.

use crate::git;
use crate::manifest::{self, Manifest};
use anyhow::{Context, Result};
use std::fs;
use std::path::Path;

const CONFLICT_START: &str = "<<<<<<<";
const CONFLICT_MIDDLE: &str = "=======";
const CONFLICT_END: &str = ">>>>>>>";
// const CONFLICT_BASE: &str = "|||||||";

/// Enrich all conflicted files with context
pub fn enrich_all_conflicts(ours_sha: &str, theirs_sha: &str) -> Result<usize> {
    let conflicted_files = get_conflicted_files()?;
    let mut enriched_count = 0;

    for file in conflicted_files {
        if enrich_conflict_markers(&file, ours_sha, theirs_sha)? {
            enriched_count += 1;
        }
    }

    Ok(enriched_count)
}

/// Get list of conflicted files
fn get_conflicted_files() -> Result<Vec<String>> {
    // git diff --name-only --diff-filter=U
    let output = git::run_git_cmd(&["diff", "--name-only", "--diff-filter=U"], None)?;

    Ok(output.lines().map(|s| s.trim().to_string()).collect())
}

/// Enrich conflict markers in a single file
fn enrich_conflict_markers(file_path: &str, ours_sha: &str, theirs_sha: &str) -> Result<bool> {
    let path = Path::new(file_path);
    if !path.exists() {
        return Ok(false);
    }

    let content = fs::read_to_string(path).context("Failed to read conflicted file")?;

    if !content.contains(CONFLICT_START) {
        return Ok(false);
    }

    // Load manifests
    let ours_manifest = manifest::load(ours_sha, None).ok();
    let theirs_manifest = manifest::load(theirs_sha, None).ok();

    if ours_manifest.is_none() && theirs_manifest.is_none() {
        return Ok(false);
    }

    let mut output = String::new();
    let lines: Vec<&str> = content.lines().collect();
    let mut current_line_idx = 0;

    while current_line_idx < lines.len() {
        let line = lines[current_line_idx];

        if line.starts_with(CONFLICT_START) {
            output.push_str(line);
            output.push('\n');

            // Get context before this marker for symbol detection
            let context_start = if current_line_idx > 50 { current_line_idx - 50 } else { 0 };
            let context = &lines[context_start..current_line_idx];

            if let Some(ref m) = ours_manifest {
                let context = format_enriched_marker("HEAD", "Your changes", m, file_path, Some(context));
                output.push_str(&context);
            }
        } else if line.starts_with(CONFLICT_MIDDLE) {
            output.push_str(line);
            output.push('\n');
        } else if line.starts_with(CONFLICT_END) {
            // Extract branch name from marker if possible
            let branch = line.trim_start_matches(CONFLICT_END).trim();

            // Get context before this marker (including the conflict body)
            // We search further back to find the symbol definition
            let context_start = if current_line_idx > 100 { current_line_idx - 100 } else { 0 };
            let context = &lines[context_start..current_line_idx];

            if let Some(ref m) = theirs_manifest {
                let context = format_enriched_marker(branch, "Their changes", m, file_path, Some(context));
                output.push_str(&context);
            }

            output.push_str(line);
            output.push('\n');
        } else {
            output.push_str(line);
            output.push('\n');
        }
        current_line_idx += 1;
    }

    fs::write(path, output).context("Failed to write enriched file")?;
    Ok(true)
}

fn format_enriched_marker(
    side: &str,
    description: &str,
    manifest: &Manifest,
    file_path: &str,
    context: Option<&[&str]>,
) -> String {
    let mut output = String::new();

    output.push_str(&format!("||| Gip CONTEXT ({} - {})\n", side, description));
    output.push_str(&format!("||| Commit: {}\n", manifest.commit));

    // Find relevant entry
    let entry = find_entry(manifest, file_path, context);

    if let Some(e) = entry {
        if !e.behavior_class.is_empty() {
            output.push_str(&format!(
                "||| behaviorClass: {}\n",
                e.behavior_class.join(", ")
            ));
        }

        if !e.rationale.is_empty() {
            output.push_str(&format!("||| rationale: {}\n", e.rationale));
        }

        if let Some(ref compat) = e.compatibility {
            output.push_str(&format!("||| breaking: {}\n", compat.breaking));

            if let Some(ref migs) = compat.migrations {
                for (i, mig) in migs.iter().enumerate() {
                    output.push_str(&format!("||| migrations[{}]: {}\n", i, mig));
                }
            }
        }

        if let Some(ref inputs) = e.contract.inputs {
            for (i, input) in inputs.iter().enumerate() {
                output.push_str(&format!("||| inputs[{}]: {}\n", i, input));
            }
        }

        if let Some(ref outputs) = e.contract.outputs {
            output.push_str(&format!("||| outputs: {}\n", outputs));
        }

        if !e.contract.preconditions.is_empty() {
            for (i, pre) in e.contract.preconditions.iter().enumerate() {
                output.push_str(&format!("||| preconditions[{}]: {}\n", i, pre));
            }
        }

        if !e.contract.postconditions.is_empty() {
            for (i, post) in e.contract.postconditions.iter().enumerate() {
                output.push_str(&format!("||| postconditions[{}]: {}\n", i, post));
            }
        }

        if !e.contract.error_model.is_empty() {
            for (i, err) in e.contract.error_model.iter().enumerate() {
                output.push_str(&format!("||| errorModel[{}]: {}\n", i, err));
            }
        }

        if !e.side_effects.is_empty() {
            for (i, side) in e.side_effects.iter().enumerate() {
                output.push_str(&format!("||| sideEffects[{}]: {}\n", i, side));
            }
        }

        output.push_str(&format!("||| symbol: {}\n", e.anchor.symbol));
    } else {
        // Fallback to global intent if no specific entry found
        if let Some(ref gi) = manifest.global_intent {
            output.push_str(&format!(
                "||| behaviorClass: {}\n",
                gi.behavior_class.join(", ")
            ));
            output.push_str(&format!("||| rationale: {}\n", gi.rationale));
        }
    }

    output
}

fn find_entry<'a>(
    manifest: &'a Manifest,
    file_path: &str,
    context: Option<&[&str]>,
) -> Option<&'a crate::manifest::Entry> {
    // 1. Filter entries by file path
    let filename = Path::new(file_path).file_name()?.to_str()?;
    
    let file_entries: Vec<&crate::manifest::Entry> = manifest.entries.iter().filter(|e| {
        e.anchor.file == file_path || 
        Path::new(&e.anchor.file)
            .file_name()
            .map(|n| n.to_str().unwrap_or(""))
            == Some(filename)
    }).collect();

    if file_entries.is_empty() {
        return None;
    }

    // 2. If context is available, try to match symbol
    if let Some(lines) = context {
        let mut best_entry: Option<&crate::manifest::Entry> = None;
        let mut min_indent = usize::MAX;

        // We search backwards from the conflict
        for line in lines.iter().rev() {
            // Calculate indentation (spaces/tabs)
            let indent = line.chars().take_while(|c| c.is_whitespace()).count();
            
            for entry in &file_entries {
                if line.contains(&entry.anchor.symbol) {
                    // Found a match.
                    // Heuristic: The enclosing function definition usually has 
                    // lower indentation than the code inside it (including calls).
                    // We prefer the match with the lowest indentation found so far.
                    if indent < min_indent {
                        best_entry = Some(entry);
                        min_indent = indent;
                    }
                }
            }
        }

        if let Some(entry) = best_entry {
            return Some(entry);
        }
    }

    // 3. Fallback: return the first entry for this file
    Some(file_entries[0])
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::manifest::types::*;

    #[test]
    fn test_format_enriched_marker_full() {
        let manifest = Manifest {
            schema_version: "2.0".to_string(),
            commit: "abc1234".to_string(),
            global_intent: None,
            entries: vec![Entry {
                anchor: Anchor {
                    file: "src/payment.rs".to_string(),
                    symbol: "processPayment".to_string(),
                    hunk_id: "H#1".to_string(),
                },
                change_type: "modify".to_string(),
                signature_delta: None,
                contract: Contract {
                    inputs: Some(vec![
                        "amount: float".to_string(),
                        "currency: string".to_string(),
                    ]),
                    outputs: Some("bool success".to_string()),
                    preconditions: vec![],
                    postconditions: vec![],
                    error_model: vec!["throws PaymentException".to_string()],
                },
                behavior_class: vec!["feature".to_string()],
                side_effects: vec![],
                compatibility: Some(Compatibility {
                    breaking: true,
                    deprecations: None,
                    migrations: Some(vec!["Update payment config".to_string()]),
                    binary_breaking: None,
                    source_breaking: None,
                    data_model_migration: None,
                }),
                tests_touched: None,
                perf_budget: None,
                security_notes: None,
                feature_flags: None,
                rationale: "Added new payment method".to_string(),
                inherits_global_intent: None,
            }],
        };

        let marker = format_enriched_marker("HEAD", "Your changes", &manifest, "src/payment.rs", None);

        assert!(marker.contains("||| Gip CONTEXT (HEAD - Your changes)"));
        assert!(marker.contains("||| Commit: abc1234"));
        assert!(marker.contains("||| behaviorClass: feature"));
        assert!(marker.contains("||| rationale: Added new payment method"));
        assert!(marker.contains("||| breaking: true"));
        assert!(marker.contains("||| migrations[0]: Update payment config"));
        assert!(marker.contains("||| inputs[0]: amount: float"));
        assert!(marker.contains("||| outputs: bool success"));
        assert!(marker.contains("||| symbol: processPayment"));
        assert!(marker.contains("||| errorModel[0]: throws PaymentException"));
    }
    
    #[test]
    fn test_find_entry_with_symbol_context() {
        let manifest = Manifest {
            schema_version: "2.0".to_string(),
            commit: "abc".to_string(),
            global_intent: None,
            entries: vec![
                Entry {
                    anchor: Anchor { file: "src/main.rs".to_string(), symbol: "main".to_string(), hunk_id: "1".to_string() },
                    change_type: "mod".to_string(), rationale: "main logic".to_string(),
                    behavior_class: vec![], contract: Contract { inputs: None, outputs: None, preconditions: vec![], postconditions: vec![], error_model: vec![] },
                    side_effects: vec![], compatibility: None, tests_touched: None, perf_budget: None, security_notes: None, feature_flags: None, inherits_global_intent: None, signature_delta: None
                },
                Entry {
                    anchor: Anchor { file: "src/main.rs".to_string(), symbol: "helper".to_string(), hunk_id: "2".to_string() },
                    change_type: "mod".to_string(), rationale: "helper logic".to_string(),
                    behavior_class: vec![], contract: Contract { inputs: None, outputs: None, preconditions: vec![], postconditions: vec![], error_model: vec![] },
                    side_effects: vec![], compatibility: None, tests_touched: None, perf_budget: None, security_notes: None, feature_flags: None, inherits_global_intent: None, signature_delta: None
                }
            ]
        };
        
        let context = vec![
            "fn helper() {",
            "    // some code",
        ];
        
        let entry = find_entry(&manifest, "src/main.rs", Some(&context));
        assert_eq!(entry.unwrap().anchor.symbol, "helper");
        
        let context_main = vec![
            "fn main() {",
            "    helper();",
        ];
        let entry_main = find_entry(&manifest, "src/main.rs", Some(&context_main));
        assert_eq!(entry_main.unwrap().anchor.symbol, "main");
    }
}
