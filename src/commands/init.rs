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
    let manifest_path = root.join("manifest.toon");

    if !manifest_path.exists() {
        let template = r#"schemaVersion: "2.0"
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
        println!("Created manifest.toon template");
    }

    println!("{}", "✓ Gip initialized successfully".green());
    println!("Created: .gip/");
    if manifest_path.exists() {
        println!("Created: manifest.toon (template)");
    }

    Ok(())
}
