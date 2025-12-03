use crate::git;
use crate::manifest::{self, Manifest};
use anyhow::Result;
use colored::*;

pub fn run(commit: Option<String>, export: bool, json: bool) -> Result<()> {
    let commit_sha = match commit {
        Some(c) => c,
        None => git::get_current_commit()?,
    };

    let manifest = match manifest::load(&commit_sha, None) {
        Ok(m) => m,
        Err(_) => {
            println!(
                "{}",
                format!("No context found for commit {}", commit_sha).yellow()
            );
            return Ok(());
        }
    };

    if json {
        let output = serde_json::to_string_pretty(&manifest)?;
        println!("{}", output);
        return Ok(());
    }

    if export {
        let output = manifest::serialize_manifest_toon(&manifest)?;
        println!("{}", output);
        return Ok(());
    }

    // Pretty print for terminal
    print_manifest(&manifest);

    Ok(())
}

fn print_manifest(manifest: &Manifest) {
    println!(
        "┌─ Commit {} (schema v{})",
        manifest.commit.cyan(),
        manifest.schema_version
    );

    if let Some(ref gi) = manifest.global_intent {
        println!("│");
        println!("│  Global Intent:");
        println!("│  Behavior: {}", gi.behavior_class.join(", ").blue());
        println!("│  Rationale: {}", gi.rationale);
    }

    for entry in &manifest.entries {
        println!("│");
        println!("│  File: {}", entry.anchor.file.yellow());
        println!("│  Symbol: {}", entry.anchor.symbol.yellow());
        println!("│  Change: {}", entry.change_type.green());
        println!("│  Rationale: {}", entry.rationale);

        if !entry.behavior_class.is_empty() {
            println!("│  Behavior: {}", entry.behavior_class.join(", ").blue());
        }

        if !entry.contract.preconditions.is_empty() {
            println!("│  Preconditions: {:?}", entry.contract.preconditions);
        }
    }
    println!("└───────────────────────────────────────────────────────────────");
}
