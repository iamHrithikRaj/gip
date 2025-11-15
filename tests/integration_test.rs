//! Integration tests for Gip
//!
//! Tests the complete workflow: init, commit, merge with conflicts

use std::fs;
use std::path::PathBuf;
use std::process::Command;
use tempfile::TempDir;

/// Helper to run gip commands
fn run_gip(args: &[&str], repo_dir: &PathBuf) -> Result<String, String> {
    let output = Command::new("cargo")
        .arg("run")
        .arg("--")
        .args(args)
        .current_dir(repo_dir)
        .output()
        .map_err(|e| format!("Failed to execute gip: {}", e))?;

    if output.status.success() {
        Ok(String::from_utf8_lossy(&output.stdout).to_string())
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

/// Helper to run git commands (for setup only)
fn run_git(args: &[&str], repo_dir: &PathBuf) -> Result<String, String> {
    let output = Command::new("git")
        .args(args)
        .current_dir(repo_dir)
        .output()
        .map_err(|e| format!("Failed to execute git: {}", e))?;

    if output.status.success() {
        Ok(String::from_utf8_lossy(&output.stdout).to_string())
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

#[test]
#[ignore] // Run with: cargo test --test integration_test -- --ignored
fn test_gip_init() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path().to_path_buf();

    // Initialize git repo first
    run_git(&["init"], &repo_path).unwrap();
    run_git(&["config", "user.name", "Test User"], &repo_path).unwrap();
    run_git(&["config", "user.email", "test@example.com"], &repo_path).unwrap();

    // Initialize gip
    let output = run_gip(&["init"], &repo_path).unwrap();
    assert!(output.contains("Gip initialized") || output.contains("initialized"));

    // Check .gip directory exists
    let gip_dir = repo_path.join(".gip");
    assert!(gip_dir.exists(), ".gip directory should exist");

    let manifest_dir = gip_dir.join("manifest");
    assert!(manifest_dir.exists(), ".gip/manifest directory should exist");
}

#[test]
#[ignore]
fn test_manifest_storage_and_retrieval() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path().to_path_buf();

    // Setup repo
    run_git(&["init"], &repo_path).unwrap();
    run_git(&["config", "user.name", "Test User"], &repo_path).unwrap();
    run_git(&["config", "user.email", "test@example.com"], &repo_path).unwrap();

    // Initialize gip
    run_gip(&["init"], &repo_path).unwrap();

    // Create a test manifest manually
    use gip::manifest::types::*;
    use gip::manifest::storage;

    let manifest = Manifest {
        schema_version: SCHEMA_VERSION_2_0.to_string(),
        commit: "test123".to_string(),
        global_intent: Some(GlobalIntent {
            behavior_class: vec!["feature".to_string()],
            rationale: "Test feature".to_string(),
        }),
        entries: vec![Entry {
            anchor: Anchor {
                file: "src/lib.rs".to_string(),
                symbol: "test_function".to_string(),
                hunk_id: "H#1".to_string(),
            },
            change_type: CHANGE_ADD.to_string(),
            signature_delta: None,
            contract: Contract {
                inputs: None,
                outputs: None,
                preconditions: vec!["none".to_string()],
                postconditions: vec!["returns success".to_string()],
                error_model: vec!["none".to_string()],
            },
            behavior_class: vec!["feature".to_string()],
            side_effects: vec![],
            compatibility: None,
            tests_touched: None,
            perf_budget: None,
            security_notes: None,
            feature_flags: None,
            rationale: "Added test function".to_string(),
            inherits_global_intent: Some(true),
        }],
    };

    let manifest_dir = repo_path.join(".gip").join("manifest");
    storage::save(&manifest, "test123", &manifest_dir).unwrap();

    // Load it back
    let loaded = storage::load("test123", &manifest_dir).unwrap();
    assert_eq!(loaded.commit, "test123");
    assert_eq!(loaded.schema_version, SCHEMA_VERSION_2_0);
    assert!(loaded.global_intent.is_some());
    assert_eq!(loaded.entries.len(), 1);
}

#[test]
#[ignore]
fn test_toon_serialization() {
    use gip::manifest::types::*;
    use gip::manifest::toon;

    let manifest = Manifest {
        schema_version: SCHEMA_VERSION_2_0.to_string(),
        commit: "abc123".to_string(),
        global_intent: Some(GlobalIntent {
            behavior_class: vec!["bugfix".to_string()],
            rationale: "Fixed critical bug".to_string(),
        }),
        entries: vec![Entry {
            anchor: Anchor {
                file: "src/payment.rs".to_string(),
                symbol: "process_payment".to_string(),
                hunk_id: "H#1".to_string(),
            },
            change_type: CHANGE_MODIFY.to_string(),
            signature_delta: Some(SignatureDelta {
                before: "process_payment(amount: f64)".to_string(),
                after: "process_payment(amount: f64, currency: Currency)".to_string(),
            }),
            contract: Contract {
                inputs: None,
                outputs: None,
                preconditions: vec!["amount > 0".to_string()],
                postconditions: vec!["payment created".to_string()],
                error_model: vec!["InvalidAmount".to_string()],
            },
            behavior_class: vec!["bugfix".to_string()],
            side_effects: vec!["writes:database".to_string()],
            compatibility: Some(Compatibility {
                breaking: true,
                deprecations: Some(vec!["old signature".to_string()]),
                migrations: Some(vec!["add Currency::USD default".to_string()]),
                binary_breaking: None,
                source_breaking: None,
                data_model_migration: None,
            }),
            tests_touched: Some(vec!["tests/payment_test.rs".to_string()]),
            perf_budget: None,
            security_notes: None,
            feature_flags: Some(vec!["MULTI_CURRENCY".to_string()]),
            rationale: "Support multiple currencies".to_string(),
            inherits_global_intent: Some(false),
        }],
    };

    // Test official TOON format
    let toon_official = toon::serialize_manifest_toon(&manifest).unwrap();
    assert!(!toon_official.is_empty());
    println!("Official TOON format:\n{}", toon_official);

    // Test legacy TOON format
    let toon_legacy = toon::serialize_manifest(&manifest);
    assert!(toon_legacy.contains("(manifest"));
    assert!(toon_legacy.contains("(globalIntent"));
    assert!(toon_legacy.contains("(signatureDelta"));
    assert!(toon_legacy.contains("(compatibility"));
    assert!(toon_legacy.contains("(featureFlags"));
    println!("\nLegacy TOON format:\n{}", toon_legacy);

    // Verify token efficiency
    let json = serde_json::to_string_pretty(&manifest).unwrap();
    println!("\nSize comparison:");
    println!("  JSON (pretty): {} bytes", json.len());
    println!("  TOON (legacy): {} bytes", toon_legacy.len());
    println!("  TOON (official): {} bytes", toon_official.len());
}

#[test]
#[ignore]
fn test_conflict_scenario() {
    // This test simulates a real conflict scenario
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path().to_path_buf();

    // Setup repo
    run_git(&["init"], &repo_path).unwrap();
    run_git(&["config", "user.name", "Test User"], &repo_path).unwrap();
    run_git(&["config", "user.email", "test@example.com"], &repo_path).unwrap();

    // Initialize gip
    run_gip(&["init"], &repo_path).unwrap();

    // Create initial file
    let file_path = repo_path.join("calculator.py");
    fs::write(
        &file_path,
        r#"def calculate_total(items):
    total = 0
    for item in items:
        total += item.price
    return total
"#,
    )
    .unwrap();

    run_git(&["add", "calculator.py"], &repo_path).unwrap();
    run_git(&["commit", "-m", "Initial commit"], &repo_path).unwrap();

    // Create and switch to feature branch
    run_git(&["checkout", "-b", "add-tax"], &repo_path).unwrap();

    // Modify file on feature branch (add tax)
    fs::write(
        &file_path,
        r#"def calculate_total(items):
    total = 0
    for item in items:
        total += item.price * 1.08  # Add 8% tax
    return total
"#,
    )
    .unwrap();

    run_git(&["add", "calculator.py"], &repo_path).unwrap();
    run_git(&["commit", "-m", "Add 8% sales tax"], &repo_path).unwrap();
    
    let tax_commit = run_git(&["rev-parse", "HEAD"], &repo_path).unwrap().trim().to_string();

    // Create manifest for tax commit
    use gip::manifest::types::*;
    use gip::manifest::storage;

    let tax_manifest = Manifest {
        schema_version: SCHEMA_VERSION_2_0.to_string(),
        commit: tax_commit.clone(),
        global_intent: None,
        entries: vec![Entry {
            anchor: Anchor {
                file: "calculator.py".to_string(),
                symbol: "calculate_total".to_string(),
                hunk_id: "H#1".to_string(),
            },
            change_type: CHANGE_MODIFY.to_string(),
            signature_delta: None,
            contract: Contract {
                inputs: None,
                outputs: None,
                preconditions: vec!["items is list with .price".to_string()],
                postconditions: vec!["returns total with 8% sales tax".to_string()],
                error_model: vec!["AttributeError if no .price".to_string()],
            },
            behavior_class: vec!["feature".to_string()],
            side_effects: vec![],
            compatibility: None,
            tests_touched: None,
            perf_budget: None,
            security_notes: None,
            feature_flags: None,
            rationale: "Added 8% sales tax to comply with state law".to_string(),
            inherits_global_intent: None,
        }],
    };

    let manifest_dir = repo_path.join(".gip").join("manifest");
    storage::save(&tax_manifest, &tax_commit, &manifest_dir).unwrap();

    // Switch back to main and make conflicting change
    run_git(&["checkout", "main"], &repo_path).unwrap();

    fs::write(
        &file_path,
        r#"def calculate_total(items):
    total = 0
    for item in items:
        total += item.price + 5.99  # Add $5.99 shipping
    return total
"#,
    )
    .unwrap();

    run_git(&["add", "calculator.py"], &repo_path).unwrap();
    run_git(&["commit", "-m", "Add flat shipping fee"], &repo_path).unwrap();
    
    let shipping_commit = run_git(&["rev-parse", "HEAD"], &repo_path).unwrap().trim().to_string();

    // Create manifest for shipping commit
    let shipping_manifest = Manifest {
        schema_version: SCHEMA_VERSION_2_0.to_string(),
        commit: shipping_commit.clone(),
        global_intent: None,
        entries: vec![Entry {
            anchor: Anchor {
                file: "calculator.py".to_string(),
                symbol: "calculate_total".to_string(),
                hunk_id: "H#1".to_string(),
            },
            change_type: CHANGE_MODIFY.to_string(),
            signature_delta: None,
            contract: Contract {
                inputs: None,
                outputs: None,
                preconditions: vec!["items is list with .price".to_string()],
                postconditions: vec!["returns total plus flat shipping".to_string()],
                error_model: vec!["AttributeError if no .price".to_string()],
            },
            behavior_class: vec!["feature".to_string()],
            side_effects: vec![],
            compatibility: None,
            tests_touched: None,
            perf_budget: None,
            security_notes: None,
            feature_flags: None,
            rationale: "Added $5.99 flat shipping fee for all orders".to_string(),
            inherits_global_intent: None,
        }],
    };

    storage::save(&shipping_manifest, &shipping_commit, &manifest_dir).unwrap();

    // Attempt merge (will conflict)
    let merge_result = run_git(&["merge", "add-tax"], &repo_path);
    assert!(merge_result.is_err() || merge_result.unwrap().contains("CONFLICT"));

    // Read the conflicted file
    let conflicted_content = fs::read_to_string(&file_path).unwrap();
    println!("\n=== CONFLICTED FILE (WITHOUT GIP ENRICHMENT) ===\n{}", conflicted_content);

    // Verify both manifests can be loaded
    let tax_loaded = storage::load(&tax_commit, &manifest_dir).unwrap();
    let shipping_loaded = storage::load(&shipping_commit, &manifest_dir).unwrap();

    println!("\n=== TAX COMMIT MANIFEST ===");
    println!("{}", serde_json::to_string_pretty(&tax_loaded).unwrap());

    println!("\n=== SHIPPING COMMIT MANIFEST ===");
    println!("{}", serde_json::to_string_pretty(&shipping_loaded).unwrap());

    // Test TOON serialization for conflict markers
    use gip::manifest::toon;
    
    println!("\n=== TAX TOON (for conflict marker) ===");
    println!("{}", toon::serialize_manifest(&tax_loaded));

    println!("\n=== SHIPPING TOON (for conflict marker) ===");
    println!("{}", toon::serialize_manifest(&shipping_loaded));

    // Verify manifests have the expected data
    assert_eq!(tax_loaded.entries[0].rationale, "Added 8% sales tax to comply with state law");
    assert_eq!(shipping_loaded.entries[0].rationale, "Added $5.99 flat shipping fee for all orders");
}
