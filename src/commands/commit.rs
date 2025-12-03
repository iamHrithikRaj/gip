use crate::git;
use crate::manifest::{self, Manifest};
use anyhow::{Context, Result};
use colored::*;
use std::fs;
use toon_format::{decode, DecodeOptions};

pub fn run(message: Option<String>, args: &[String]) -> Result<()> {
    // 1. Check for manifest.toon
    let root = git::get_repo_root()?;
    let manifest_path = root.join("manifest.toon");

    let manifest: Manifest = if manifest_path.exists() {
        let content = fs::read_to_string(&manifest_path).context("Failed to read manifest.toon")?;
        let opts = DecodeOptions::new().with_strict(false);
        decode(&content, &opts).context("Failed to parse manifest.toon")?
    } else {
        // If no manifest, maybe allow commit without it?
        // The docs say "Gip works without manifests. Commits without manifests are just regular git commits."
        // But also "Gip attaches a manifest...".
        // If manifest.toon exists, use it. If not, proceed as normal git commit?
        // But the user might want to enforce it.
        // For now, if no manifest.toon, we just commit.
        println!("{}", "No manifest.toon found. Committing without context.".yellow());
        
        // Pass through to git commit
        let mut git_args = vec!["commit".to_string()];
        if let Some(msg) = message {
            git_args.push("-m".to_string());
            git_args.push(msg);
        }
        git_args.extend_from_slice(args);
        
        return crate::commands::passthrough::run(&git_args);
    };

    // 2. Validate manifest (basic validation is done by decode)
    // TODO: Add more semantic validation

    println!("{}", "✓ Manifest validated".green());

    // 3. Commit using git
    let mut git_args = vec!["commit".to_string()];
    if let Some(msg) = message {
        git_args.push("-m".to_string());
        git_args.push(msg);
    }
    git_args.extend_from_slice(args);

    // Run git commit
    crate::commands::passthrough::run(&git_args)?;

    // 4. Attach manifest as git note
    let commit_sha = git::get_current_commit()?;
    
    // Update manifest with actual commit SHA
    let mut final_manifest = manifest.clone();
    final_manifest.commit = commit_sha.clone();

    manifest::save(&final_manifest, &commit_sha, None)?;

    println!("{}", "✓ Changes committed with context".green());
    println!("{}", "✓ Manifest attached as git note".green());

    Ok(())
}
