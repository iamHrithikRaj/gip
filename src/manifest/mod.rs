//! Manifest module - handles manifest storage, loading, and migration
//!
//! Provides functionality for creating, storing, and loading Gip manifests that
//! capture structured context about code changes.

pub mod storage;
pub mod toon;
pub mod types;

pub use storage::{load, load_pending, migrate_v1_to_v2, save, save_pending};
pub use toon::{serialize_manifest, serialize_manifest_toon};
pub use types::*;
