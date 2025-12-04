use crate::git;
use anyhow::Result;
use colored::*;
use std::fs;

pub fn run() -> Result<()> {
    println!("{}", "Initializing Gip...".cyan());

    if !git::is_git_repo() {
        anyhow::bail!("Not a git repository. Run 'git init' first.");
    }

    git::ensure_gip_dir()?;

    // Create a template manifest if it doesn't exist
    let root = git::get_repo_root()?;
    let gip_dir = root.join(".gip");
    let manifest_path = gip_dir.join("manifest.toon");

    if !manifest_path.exists() {
        let template = r#"; Gip Manifest Template
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
        fs::write(&manifest_path, template)?;
        println!("Created .gip/manifest.toon template");
    }

    // Add .gip to .gitignore
    let gitignore_path = root.join(".gitignore");
    let mut gitignore_content = if gitignore_path.exists() {
        fs::read_to_string(&gitignore_path)?
    } else {
        String::new()
    };

    if !gitignore_content.contains(".gip") {
        if !gitignore_content.is_empty() && !gitignore_content.ends_with('\n') {
            gitignore_content.push('\n');
        }
        gitignore_content.push_str(".gip\n");
        fs::write(&gitignore_path, gitignore_content)?;
        println!("Added .gip to .gitignore");
    }

    println!("{}", "✓ Gip initialized successfully".green());
    println!("Created: .gip/");
    if manifest_path.exists() {
        println!("Created: .gip/manifest.toon (template)");
    }

    Ok(())
}
