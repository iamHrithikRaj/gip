//! Gip - Git++ CLI
//!
//! A lightweight Git wrapper that enriches merge conflicts with structured context

use anyhow::Result;
use clap::{Parser, Subcommand};
use gip::commands;

#[derive(Parser)]
#[command(name = "gip")]
#[command(version, about = "Git with Intent Preservation - Context-aware git wrapper", long_about = None)]
#[command(disable_help_subcommand = true)]
struct Cli {
    #[command(subcommand)]
    command: Option<Commands>,
}

#[derive(Subcommand)]
enum Commands {
    /// Initialize Gip in the current repository
    Init,

    /// Commit with manifest attachment
    Commit {
        /// Commit message
        #[arg(short, long)]
        message: Option<String>,

        /// Force commit without manifest
        #[arg(short, long)]
        force: bool,

        /// Additional git arguments
        #[arg(allow_hyphen_values = true, trailing_var_arg = true)]
        args: Vec<String>,
    },

    /// Push code AND context notes to remote
    Push {
        /// Additional git arguments
        #[arg(allow_hyphen_values = true, trailing_var_arg = true)]
        args: Vec<String>,
    },

    /// Merge with enriched conflict markers
    Merge {
        /// Additional git arguments (e.g. branch name)
        #[arg(allow_hyphen_values = true, trailing_var_arg = true)]
        args: Vec<String>,
    },

    /// Rebase with enriched conflict markers
    Rebase {
        /// Additional git arguments (e.g. branch name)
        #[arg(allow_hyphen_values = true, trailing_var_arg = true)]
        args: Vec<String>,
    },

    /// Show semantic history/context
    Context {
        /// Commit SHA or file path (optional)
        target: Option<String>,

        /// Export context to TOON format
        #[arg(long)]
        export: bool,
    },

    #[command(external_subcommand)]
    External(Vec<String>),
}

fn main() -> Result<()> {
    let cli = Cli::parse();

    match cli.command {
        Some(Commands::Init) => commands::init::run(),
        Some(Commands::Commit {
            message,
            force,
            args,
        }) => commands::commit::run(message, force, &args),
        Some(Commands::Push { args }) => commands::push::run(&args),
        Some(Commands::Merge { args }) => commands::merge::run(&args),
        Some(Commands::Rebase { args }) => commands::rebase::run(&args),
        Some(Commands::Context { target, export }) => commands::context::run(target, export),
        Some(Commands::External(args)) => commands::passthrough::run(&args),
        None => {
            // Show help if no args
            use clap::CommandFactory;
            Cli::command().print_help()?;
            Ok(())
        }
    }
}
