use anyhow::{Context, Result};
use std::process::Command;

pub fn run(args: &[String]) -> Result<()> {
    let mut cmd = Command::new("git");
    cmd.args(args);

    let status = cmd.status().context("Failed to execute git command")?;

    if !status.success() {
        std::process::exit(status.code().unwrap_or(1));
    }

    Ok(())
}
