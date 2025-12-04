use crate::git;
use crate::manifest::{self, Manifest};
use anyhow::{Context, Result};
use colored::*;
use std::fs;
use toon_format::{decode, DecodeOptions};

const TEMPLATE: &str = r#"; Gip Manifest Template
; This file describes the semantic intent of your changes.
; It is used to enrich merge conflicts with context.
;
; INSTRUCTIONS FOR LLM/AGENTS:
; 1. Analyze the code changes in the current commit.
; 2. Update the fields below to reflect the actual changes.
; 3. 'rationale' should explain WHY the change was made.
; 4. 'behaviorClass' options: feature, bugfix, refactor, perf, security, config.
; 5. 'changeType' options: add, modify, delete, rename.
; 6. Remove these instruction comments if desired, but keep the structure.

schemaVersion: "2.0"
commit: HEAD
entries[1]:
  - anchor:
      file: src/main.rs
      symbol: main
      hunkId: H#1
    changeType: modify
    rationale: Describe your changes here
    behaviorClass[1]: feature
    contract:
      preconditions[1]: none
      postconditions[1]: program_runs
      errorModel[1]: panic_on_error
"#;

pub fn run(message: Option<String>, force: bool, args: &[String]) -> Result<()> {
    // 1. Check for manifest.toon
    let root = git::get_repo_root()?;
    let manifest_path = root.join(".gip").join("manifest.toon");

    let manifest_content = if manifest_path.exists() {
        Some(fs::read_to_string(&manifest_path).context("Failed to read manifest.toon")?)
    } else {
        None
    };

    // Validation Logic
    if !force {
        let mut reject = false;
        let mut reason = String::new();

        if manifest_content.is_none() {
            // Create template
            let gip_dir = root.join(".gip");
            if !gip_dir.exists() {
                fs::create_dir_all(&gip_dir)?;
            }
            fs::write(&manifest_path, TEMPLATE)?;

            reject = true;
            reason = format!(
                "Manifest file was missing. Created new template at {}",
                manifest_path.display()
            );
        } else if let Some(ref content) = manifest_content {
            // Normalize line endings for comparison
            let normalized_content = content.replace("\r\n", "\n");
            let normalized_template = TEMPLATE.replace("\r\n", "\n");

            if normalized_content.trim() == normalized_template.trim() {
                reject = true;
                reason = "Manifest file is unchanged from template".to_string();
            } else if content.contains("Describe your changes here") {
                reject = true;
                reason =
                    "Manifest contains placeholder text 'Describe your changes here'".to_string();
            }
        }

        if reject {
            // Print LLM friendly error
            eprintln!(
                "{}",
                "ERROR: Commit rejected due to missing or incomplete manifest."
                    .red()
                    .bold()
            );
            eprintln!("Reason: {}", reason);
            eprintln!(
                "\nTo commit, you must fill out the manifest file at: {}",
                manifest_path.display()
            );
            eprintln!("This file describes the intent of your changes in TOON format.");
            eprintln!("\nHere is the template structure you need to follow:");
            eprintln!("---------------------------------------------------");
            eprintln!("{}", TEMPLATE);
            eprintln!("---------------------------------------------------");
            eprintln!("\nINSTRUCTIONS FOR AGENT/LLM:");
            eprintln!("1. Read the file at: {}", manifest_path.display());
            eprintln!("2. Understand the code changes you are committing.");
            eprintln!("3. Fill out the 'rationale', 'changeType', and 'behaviorClass' fields in the manifest file.");
            eprintln!("4. Save the file.");
            eprintln!("5. Retry the commit command.");
            eprintln!("\nIf you really want to commit without a manifest, use the --force flag.");
            anyhow::bail!("Commit rejected. See output for details.");
        }
    }

    let manifest: Option<Manifest> = if let Some(content) = manifest_content {
        let opts = DecodeOptions::new().with_strict(false);
        Some(decode(&content, &opts).context("Failed to parse manifest.toon")?)
    } else {
        if !force {
            // Should be caught above, but just in case
            anyhow::bail!("Manifest missing and force not set");
        }
        println!(
            "{}",
            "No manifest.toon found. Committing without context (FORCE).".yellow()
        );
        None
    };

    if manifest.is_some() {
        println!("{}", "✓ Manifest validated".green());
    }

    // 3. Commit using git
    let mut git_args = vec!["commit".to_string()];
    if let Some(msg) = message {
        git_args.push("-m".to_string());
        git_args.push(msg);
    }
    git_args.extend_from_slice(args);

    // Run git commit
    crate::commands::passthrough::run(&git_args)?;

    // 4. Attach manifest as git note if it exists
    if let Some(manifest) = manifest {
        let commit_sha = git::get_current_commit()?;

        // Update manifest with actual commit SHA
        let mut final_manifest = manifest.clone();
        final_manifest.commit = commit_sha.clone();

        manifest::save(&final_manifest, &commit_sha, None)?;

        println!("{}", "✓ Changes committed with context".green());
        println!("{}", "✓ Manifest attached as git note".green());
    }

    Ok(())
}
