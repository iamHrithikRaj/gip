//! Manifest module - handles manifest storage, loading, and migration
//!
//! Provides functionality for creating, storing, and loading Gip manifests that
//! capture structured context about code changes.

pub mod types;
pub mod storage;
pub mod toon;

pub use types::*;
pub use storage::{save, load, save_pending, load_pending, migrate_v1_to_v2};
pub use toon::{serialize_manifest, serialize_manifest_toon};
