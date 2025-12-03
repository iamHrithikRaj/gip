use crate::git;
use crate::merge;
use anyhow::Result;
use colored::*;

pub fn run(args: &[String]) -> Result<()> {
    println!("{}", "Rebasing with Gip...".cyan());

    // 1. Run git rebase
    let mut git_args = vec!["rebase".to_string()];
    git_args.extend_from_slice(args);

    let status = std::process::Command::new("git").args(&git_args).status()?;

    if status.success() {
        println!("{}", "Rebase successful".green());
        return Ok(());
    }

    // 2. If failed, check for conflicts
    println!(
        "{}",
        "Rebase conflict detected. Enriching markers...".yellow()
    );

    // In rebase:
    // HEAD is the commit being replayed (theirs in merge terms, but ours in rebase terms?)
    // REBASE_HEAD is the commit we are rebasing onto?
    // Actually:
    // "ours" (HEAD) is the upstream we are rebasing ONTO.
    // "theirs" is the commit being applied.
    // But git conflict markers might be swapped depending on rebase type.

    // Let's try to get REBASE_HEAD and HEAD.
    // During rebase, HEAD is detached at the commit we are rebasing onto (upstream).
    // REBASE_HEAD is the commit being applied.

    let ours_sha = git::get_current_commit()?; // Upstream
    let theirs_sha = match git::run_git_cmd(&["rev-parse", "REBASE_HEAD"], None) {
        Ok(sha) => sha,
        Err(_) => {
            // Maybe interactive rebase or something else?
            // Try to find stopped commit.
            println!(
                "{}",
                "Could not determine REBASE_HEAD. Skipping enrichment.".red()
            );
            std::process::exit(status.code().unwrap_or(1));
        }
    };

    // Note: In rebase, "ours" is upstream, "theirs" is the patch.
    // But conflict markers usually show HEAD as upstream.
    let count = merge::enrich_all_conflicts(&ours_sha, &theirs_sha)?;

    if count > 0 {
        println!(
            "{}",
            format!("✓ Enriched {} conflicted files with context", count).green()
        );
    } else {
        println!("{}", "No context available for conflicts".yellow());
    }

    std::process::exit(status.code().unwrap_or(1));
}
