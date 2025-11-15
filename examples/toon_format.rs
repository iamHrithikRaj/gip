//! Example demonstrating TOON format usage in Gip
//!
//! This example shows how manifests are serialized to TOON format
//! for efficient token usage with LLMs.

use gip::manifest::*;

fn main() -> anyhow::Result<()> {
    // Create a sample manifest
    let manifest = Manifest {
        schema_version: SCHEMA_VERSION_2_0.to_string(),
        commit: "abc123def456".to_string(),
        global_intent: Some(GlobalIntent {
            behavior_class: vec![BEHAVIOR_REFACTOR.to_string()],
            rationale: "Refactor payment processing for better maintainability".to_string(),
        }),
        entries: vec![Entry {
            anchor: Anchor {
                file: "src/payment.rs".to_string(),
                symbol: "process_payment".to_string(),
                hunk_id: "H#42".to_string(),
            },
            change_type: CHANGE_MODIFY.to_string(),
            signature_delta: Some(SignatureDelta {
                before: "fn process_payment(amount: f64)".to_string(),
                after: "fn process_payment(amount: f64, currency: Currency)".to_string(),
            }),
            contract: Contract {
                inputs: Some(vec![
                    "amount: f64".to_string(),
                    "currency: Currency".to_string(),
                ]),
                outputs: Some("Result<PaymentId>".to_string()),
                preconditions: vec!["amount > 0".to_string(), "currency is valid".to_string()],
                postconditions: vec![
                    "payment created".to_string(),
                    "returns unique ID".to_string(),
                ],
                error_model: vec!["InvalidAmount".to_string(), "InvalidCurrency".to_string()],
            },
            behavior_class: vec![BEHAVIOR_FEATURE.to_string()],
            side_effects: vec!["writes:database".to_string(), "logs:audit".to_string()],
            compatibility: Some(Compatibility {
                breaking: true,
                deprecations: Some(vec!["old signature removed".to_string()]),
                migrations: Some(vec!["add Currency::USD as default".to_string()]),
                binary_breaking: None,
                source_breaking: None,
                data_model_migration: None,
            }),
            tests_touched: Some(vec!["tests/payment_tests.rs".to_string()]),
            perf_budget: Some(PerfBudget {
                expected_max_latency_ms: Some(100),
                cpu_delta_pct: Some(5),
            }),
            security_notes: Some(vec!["validates currency code".to_string()]),
            feature_flags: Some(vec!["MULTI_CURRENCY".to_string()]),
            rationale: "Add multi-currency support for international payments".to_string(),
            inherits_global_intent: Some(false),
        }],
    };

    println!("=== JSON Format ===");
    let json = serde_json::to_string_pretty(&manifest)?;
    println!("{}", json);
    println!("\nJSON size: {} bytes", json.len());
    println!("Estimated tokens: ~{}", json.len() / 4);

    println!("\n=== TOON Format (Official) ===");
    let toon = serialize_manifest_toon(&manifest)?;
    println!("{}", toon);
    println!("\nTOON size: {} bytes", toon.len());
    println!("Estimated tokens: ~{}", toon.len() / 4);

    let savings = ((json.len() - toon.len()) as f64 / json.len() as f64) * 100.0;
    println!("\n✨ Token savings: {:.2}%", savings);

    println!("\n=== Legacy TOON Format ===");
    let legacy_toon = serialize_manifest(&manifest);
    println!("{}", legacy_toon);
    println!("\nLegacy TOON size: {} bytes", legacy_toon.len());

    Ok(())
}
