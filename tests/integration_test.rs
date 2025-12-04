//! Integration tests for Gip
//!
//! Tests the complete workflow: init, commit, merge with conflicts

use assert_cmd::Command;
use predicates::prelude::*;
use std::fs;
use std::path::Path;
use tempfile::TempDir;

/// Helper to run git commands
fn run_git(args: &[&str], repo_dir: &Path) {
    let output = std::process::Command::new("git")
        .args(args)
        .current_dir(repo_dir)
        .output()
        .expect("Failed to execute git");

    if !output.status.success() {
        panic!(
            "Git command failed: git {}\nStderr: {}",
            args.join(" "),
            String::from_utf8_lossy(&output.stderr)
        );
    }
}

#[test]
fn test_gip_init() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path();

    // Initialize git repo first
    run_git(&["init"], repo_path);
    run_git(&["config", "user.name", "Test User"], repo_path);
    run_git(&["config", "user.email", "test@example.com"], repo_path);

    // Initialize gip
    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("init")
        .assert()
        .success()
        .stdout(predicate::str::contains("Gip initialized successfully"));

    // Check .gip directory exists
    let gip_dir = repo_path.join(".gip");
    assert!(gip_dir.exists(), ".gip directory should exist");

    let manifest_path = gip_dir.join("manifest.toon");
    assert!(
        manifest_path.exists(),
        "manifest.toon template should exist in .gip"
    );

    // Check .gitignore
    let gitignore_path = repo_path.join(".gitignore");
    assert!(gitignore_path.exists(), ".gitignore should exist");
    let gitignore_content = fs::read_to_string(gitignore_path).unwrap();
    assert!(gitignore_content.contains(".gip"), ".gitignore should contain .gip");
}

#[test]
fn test_gip_commit_and_context() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path();

    // Setup
    run_git(&["init"], repo_path);
    run_git(&["config", "user.name", "Test User"], repo_path);
    run_git(&["config", "user.email", "test@example.com"], repo_path);

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path).arg("init").assert().success();

    // Create a file
    fs::write(repo_path.join("test.txt"), "hello world").unwrap();
    run_git(&["add", "test.txt"], repo_path);

    // Modify manifest to be valid
    let manifest_path = repo_path.join(".gip").join("manifest.toon");
    let manifest_content = r#"schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: test.txt
      symbol: main
      hunkId: H#1
    changeType: add
    rationale: Initial test file
    behaviorClass[1]: feature
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
"#;
    fs::write(&manifest_path, manifest_content).unwrap();

    // Commit with gip
    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("commit")
        .arg("-m")
        .arg("feat: initial commit")
        .assert()
        .success()
        .stdout(predicate::str::contains("Changes committed with context"));

    // Verify git commit happened
    let output = std::process::Command::new("git")
        .args(&["log", "-1", "--pretty=%B"])
        .current_dir(repo_path)
        .output()
        .unwrap();
    let commit_msg = String::from_utf8_lossy(&output.stdout);
    assert!(commit_msg.contains("feat: initial commit"));

    // Verify context exists
    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("context")
        .assert()
        .success()
        .stdout(predicate::str::contains("Rationale"));
}

#[test]
fn test_gip_merge_enrichment() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path();

    // Setup
    run_git(&["init"], repo_path);
    run_git(&["checkout", "-b", "main"], repo_path);
    run_git(&["config", "user.name", "Test User"], repo_path);
    run_git(&["config", "user.email", "test@example.com"], repo_path);

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path).arg("init").assert().success();

    // Initial commit
    fs::write(repo_path.join("file.txt"), "base content").unwrap();
    run_git(&["add", "file.txt"], repo_path);

    // Modify manifest for initial commit
    let manifest_path = repo_path.join(".gip").join("manifest.toon");
    let manifest_content_init = r#"schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: file.txt
      symbol: main
      hunkId: H#1
    changeType: add
    rationale: Initial file
    behaviorClass[1]: feature
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
"#;
    fs::write(&manifest_path, manifest_content_init).unwrap();

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("commit")
        .arg("-m")
        .arg("initial")
        .assert()
        .success();

    // Create branch
    run_git(&["checkout", "-b", "feature"], repo_path);

    // Modify file on feature
    fs::write(repo_path.join("file.txt"), "feature content").unwrap();
    run_git(&["add", "file.txt"], repo_path);

    // Update manifest for feature
    let manifest_content = r#"schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: file.txt
      symbol: main
      hunkId: H#1
    changeType: modify
    rationale: Feature change rationale
    behaviorClass[1]: feature
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
"#;
    fs::write(&manifest_path, manifest_content).unwrap();

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("commit")
        .arg("-m")
        .arg("feature commit")
        .assert()
        .success();

    // Switch back to main
    run_git(&["checkout", "main"], repo_path);

    // Modify file on main (conflict)
    fs::write(repo_path.join("file.txt"), "main content").unwrap();
    run_git(&["add", "file.txt"], repo_path);

    // Update manifest for main
    let manifest_content_main = r#"schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: file.txt
      symbol: main
      hunkId: H#1
    changeType: modify
    rationale: Main change rationale
    behaviorClass[1]: refactor
    contract:
      preconditions[0]:
      postconditions[0]:
      errorModel[0]:
"#;
    fs::write(&manifest_path, manifest_content_main).unwrap();

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("commit")
        .arg("-m")
        .arg("main commit")
        .assert()
        .success();

    // Merge feature into main using gip
    // This should fail with conflict, but enrich markers
    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("merge")
        .arg("feature")
        .assert()
        .failure(); // Expect failure due to conflict

    // Check file content for enriched markers
    let content = fs::read_to_string(repo_path.join("file.txt")).unwrap();

    // Should contain standard git markers
    assert!(content.contains("<<<<<<< HEAD"));
    assert!(content.contains("======="));
    assert!(content.contains(">>>>>>> feature"));

    // Should contain Gip enrichment
    assert!(content.contains("||| Gip CONTEXT (HEAD - Your changes)"));
    assert!(content.contains("||| behaviorClass: refactor"));
    assert!(content.contains("||| rationale: Main change rationale"));

    assert!(content.contains("||| Gip CONTEXT (feature - Their changes)"));
    assert!(content.contains("||| behaviorClass: feature"));
    assert!(content.contains("||| rationale: Feature change rationale"));
}

#[test]
fn test_gip_commit_rejects_incomplete_manifest() {
    let temp_dir = TempDir::new().unwrap();
    let repo_path = temp_dir.path();

    // Setup
    run_git(&["init"], repo_path);
    run_git(&["config", "user.name", "Test User"], repo_path);
    run_git(&["config", "user.email", "test@example.com"], repo_path);

    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path).arg("init").assert().success();

    // Create a file
    fs::write(repo_path.join("test.txt"), "hello world").unwrap();
    run_git(&["add", "test.txt"], repo_path);

    // Try to commit without modifying manifest
    // This should fail and print instructions
    let mut cmd = Command::cargo_bin("gip").unwrap();
    cmd.current_dir(repo_path)
        .arg("commit")
        .arg("-m")
        .arg("feat: should fail")
        .assert()
        .failure()
        .stderr(predicate::str::contains("INSTRUCTIONS FOR AGENT/LLM"))
        .stderr(predicate::str::contains("Fill out the 'rationale'"))
        .stderr(predicate::str::contains("manifest.toon"));
}
