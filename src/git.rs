//! Git integration module - utilities for interacting with Git repositories
//!
//! Provides functions for executing Git commands, retrieving commit information,
//! and configuring Gip's custom merge driver.

use anyhow::{Context, Result};
use std::path::{Path, PathBuf};
use std::process::Command;

/// Helper to run git command with optional CWD
pub fn run_git_cmd(args: &[&str], cwd: Option<&Path>) -> Result<String> {
    let mut cmd = Command::new("git");
    cmd.args(args);
    if let Some(dir) = cwd {
        cmd.current_dir(dir);
    }

    let output = cmd.output().context("Failed to execute git command")?;

    if !output.status.success() {
        let stderr = String::from_utf8_lossy(&output.stderr);
        anyhow::bail!("Git command failed: {}", stderr);
    }

    Ok(String::from_utf8(output.stdout)
        .context("Invalid UTF-8 in git output")?
        .trim()
        .to_string())
}

/// Check if current directory is a Git repository
pub fn is_git_repo() -> bool {
    run_git_cmd(&["rev-parse", "--git-dir"], None).is_ok()
}

/// Get the root directory of the Git repository
pub fn get_repo_root() -> Result<PathBuf> {
    let path = run_git_cmd(&["rev-parse", "--show-toplevel"], None)?;
    Ok(PathBuf::from(path))
}

/// Get the current commit SHA
pub fn get_current_commit() -> Result<String> {
    run_git_cmd(&["rev-parse", "HEAD"], None)
}

/// Get the diff of staged changes
pub fn get_staged_diff() -> Result<String> {
    run_git_cmd(&["diff", "--cached"], None)
}

/// Check if there are staged changes
pub fn has_staged_changes() -> bool {
    run_git_cmd(&["diff", "--cached", "--quiet"], None).is_err()
}

/// Add a note to a commit using the custom gip ref
pub fn add_note(commit_sha: &str, content: &str, cwd: Option<&Path>) -> Result<()> {
    run_git_cmd(
        &[
            "notes",
            "--ref=gip",
            "add",
            "-f",
            "-m",
            content,
            commit_sha,
        ],
        cwd,
    )?;
    Ok(())
}

/// Get a note from a commit using the custom gip ref
pub fn get_note(commit_sha: &str, cwd: Option<&Path>) -> Result<String> {
    run_git_cmd(&["notes", "--ref=gip", "show", commit_sha], cwd)
}

/// Push gip notes to remote
pub fn push_notes(remote: &str) -> Result<()> {
    run_git_cmd(&["push", remote, "refs/notes/gip"], None)?;
    Ok(())
}

/// Fetch gip notes from remote
pub fn fetch_notes(remote: &str) -> Result<()> {
    run_git_cmd(&["fetch", remote, "refs/notes/gip:refs/notes/gip"], None)?;
    Ok(())
}


/// Get the .gip directory path
pub fn get_gip_dir() -> Result<PathBuf> {
    let root = get_repo_root()?;
    Ok(root.join(".gip"))
}

/// Get the manifest storage directory
pub fn get_manifest_dir() -> Result<PathBuf> {
    let gip_dir = get_gip_dir()?;
    Ok(gip_dir.join("manifest"))
}

/// Ensure .gip directory structure exists
pub fn ensure_gip_dir() -> Result<()> {
    let manifest_dir = get_manifest_dir()?;
    std::fs::create_dir_all(manifest_dir).context("Failed to create .gip directory structure")?;
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_is_git_repo() {
        // This will depend on test environment
        // For now, just check it doesn't panic
        let _ = is_git_repo();
    }

    #[test]
    fn test_get_gip_dir_structure() {
        // Test the path construction logic
        let result = get_repo_root();
        if let Ok(root) = result {
            let gip_dir = root.join(".gip");
            let manifest_dir = gip_dir.join("manifest");
            assert_eq!(manifest_dir.file_name().unwrap(), "manifest");
        }
    }
}
