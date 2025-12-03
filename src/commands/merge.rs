use crate::git;
use crate::merge;
use anyhow::Result;
use colored::*;

pub fn run(args: &[String]) -> Result<()> {
    println!("{}", "Merging with Gip...".cyan());

    // 1. Run git merge
    let mut git_args = vec!["merge".to_string()];
    git_args.extend_from_slice(args);

    // We allow git merge to fail (conflict), so we don't use passthrough::run which exits on error
    let status = std::process::Command::new("git").args(&git_args).status()?;

    if status.success() {
        println!("{}", "Merge successful".green());
        return Ok(());
    }

    // 2. If failed, check for conflicts
    println!(
        "{}",
        "Merge conflict detected. Enriching markers...".yellow()
    );

    // We need to know the SHAs being merged.
    // HEAD is ours.
    // The other side is in MERGE_HEAD usually, or we can parse args.
    // But args might be a branch name.
    // `git rev-parse MERGE_HEAD` should work if merge is in progress.

    let ours_sha = git::get_current_commit()?;
    let theirs_sha = match git::run_git_cmd(&["rev-parse", "MERGE_HEAD"], None) {
        Ok(sha) => sha,
        Err(_) => {
            // Maybe rebase? Or just failed merge without starting?
            println!(
                "{}",
                "Could not determine MERGE_HEAD. Skipping enrichment.".red()
            );
            std::process::exit(status.code().unwrap_or(1));
        }
    };

    let count = merge::enrich_all_conflicts(&ours_sha, &theirs_sha)?;

    if count > 0 {
        println!(
            "{}",
            format!("✓ Enriched {} conflicted files with context", count).green()
        );
    } else {
        println!("{}", "No context available for conflicts".yellow());
    }

    // Exit with the original status code
    std::process::exit(status.code().unwrap_or(1));
}
