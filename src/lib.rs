//! Gip - Git++ 
//! 
//! A lightweight Git wrapper that enriches merge conflicts with structured context
//! for humans and LLMs.

pub mod manifest;
pub mod diff;
pub mod git;
pub mod merge;
pub mod prompt;
pub mod toon;

// Re-export commonly used types
pub use manifest::{Manifest, Entry, Contract};

/// Result type alias using anyhow::Error
pub type Result<T> = anyhow::Result<T>;
