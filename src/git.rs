//! Git integration module - utilities for interacting with Git repositories
//!
//! Provides functions for executing Git commands, retrieving commit information,
//! and configuring Gip's custom merge driver.

use anyhow::{Context, Result};
use std::path::PathBuf;
use std::process::Command;

/// Check if current directory is a Git repository
pub fn is_git_repo() -> bool {
    Command::new("git")
        .args(["rev-parse", "--git-dir"])
        .output()
        .map(|output| output.status.success())
        .unwrap_or(false)
}

/// Get the root directory of the Git repository
pub fn get_repo_root() -> Result<PathBuf> {
    let output = Command::new("git")
        .args(["rev-parse", "--show-toplevel"])
        .output()
        .context("Failed to execute git command")?;

    if !output.status.success() {
        anyhow::bail!("Not a git repository");
    }

    let path = String::from_utf8(output.stdout)
        .context("Invalid UTF-8 in git output")?
        .trim()
        .to_string();

    Ok(PathBuf::from(path))
}

/// Get the current commit SHA
pub fn get_current_commit() -> Result<String> {
    let output = Command::new("git")
        .args(["rev-parse", "HEAD"])
        .output()
        .context("Failed to execute git command")?;

    if !output.status.success() {
        anyhow::bail!("Failed to get current commit");
    }

    Ok(String::from_utf8(output.stdout)
        .context("Invalid UTF-8 in git output")?
        .trim()
        .to_string())
}

/// Get the diff of staged changes
pub fn get_staged_diff() -> Result<String> {
    let output = Command::new("git")
        .args(["diff", "--cached"])
        .output()
        .context("Failed to execute git command")?;

    String::from_utf8(output.stdout)
        .context("Invalid UTF-8 in git output")
}

/// Check if there are staged changes
pub fn has_staged_changes() -> bool {
    Command::new("git")
        .args(["diff", "--cached", "--quiet"])
        .status()
        .map(|status| !status.success())
        .unwrap_or(false)
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
    std::fs::create_dir_all(manifest_dir)
        .context("Failed to create .gip directory structure")?;
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
