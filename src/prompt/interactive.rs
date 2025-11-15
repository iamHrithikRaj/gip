//! Interactive prompts for gathering manifest information

use anyhow::Result;

/// Prompt user for manifest information
pub fn prompt_for_manifest() -> Result<()> {
    // TODO: Implement interactive prompts
    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_prompt_for_manifest() {
        // TODO: Implement interactive prompt tests
        let result = prompt_for_manifest();
        assert!(result.is_ok());
    }
}
