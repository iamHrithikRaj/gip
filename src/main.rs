//! Gip - Git++ CLI
//!
//! A lightweight Git wrapper that enriches merge conflicts with structured context

use anyhow::Result;
use clap::{Parser, Subcommand};

#[derive(Parser)]
#[command(name = "gip")]
#[command(version, about = "A lightweight Git wrapper that enriches merge conflicts with structured context", long_about = None)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize Gip in the current repository
    Init,

    /// Commit with manifest creation
    Commit {
        /// Commit message
        #[arg(short, long)]
        message: Option<String>,

        /// Create manifest interactively
        #[arg(short = 'c', long)]
        create_manifest: bool,
    },

    /// Show version information
    Version,
}

fn main() -> Result<()> {
    let cli = Cli::parse();

    match cli.command {
        Some(Commands::Init) => {
            println!("Initializing Gip...");
            gip::git::ensure_gip_dir()?;
            println!("Gip initialized successfully!");
        }
        Some(Commands::Commit {
            message: _,
            create_manifest: _,
        }) => {
            println!("Commit functionality not yet implemented");
        }
        Some(Commands::Version) => {
            println!("gip version {}", env!("CARGO_PKG_VERSION"));
        }
        None => {
            println!("No command specified. Use --help for usage information.");
        }
    }

    Ok(())
}
