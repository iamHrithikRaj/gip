use crate::git;
use anyhow::Result;
use colored::*;

pub fn run(args: &[String]) -> Result<()> {
    // 1. Push code
    println!("{}", "Pushing code...".cyan());
    let mut git_args = vec!["push".to_string()];
    git_args.extend_from_slice(args);

    crate::commands::passthrough::run(&git_args)?;

    // 2. Push notes
    println!("{}", "Pushing context notes...".cyan());
    // Assuming 'origin' for now, or parse from args if provided
    // Ideally we should detect the remote being pushed to.
    // For simplicity, we'll try to push to origin.
    // TODO: Parse remote from args
    let remote = "origin";

    match git::push_notes(remote) {
        Ok(_) => println!("{}", "✓ Context notes pushed".green()),
        Err(e) => println!(
            "{}",
            format!("Warning: Failed to push notes: {}", e).yellow()
        ),
    }

    Ok(())
}
