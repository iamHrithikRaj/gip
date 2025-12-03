//! Gip - Git++
//!
//! A lightweight Git wrapper that enriches merge conflicts with structured context
//! for humans and LLMs.

pub mod commands;
pub mod git;
pub mod manifest;
pub mod merge;

// Re-export commonly used types
pub use manifest::{Contract, Entry, Manifest};

/// Result type alias using anyhow::Error
pub type Result<T> = anyhow::Result<T>;
