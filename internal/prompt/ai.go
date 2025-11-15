package prompt

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/iamHrithikRaj/gip/internal/diff"
	"github.com/iamHrithikRaj/gip/internal/llm"
	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// AICommit generates a manifest using LLM and commits
func AICommit(userIntent string) error {
	color.Cyan("Analyzing staged changes...")

	// Get staged diff
	diffContent, err := diff.GetStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged diff: %w", err)
	}

	if diffContent == "" {
		color.Yellow("No staged changes found. Stage changes first with: git add")
		return nil
	}

	// Get API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable not set\nSet it with: export OPENAI_API_KEY=sk-...")
	}

	// Create LLM config
	config := llm.DefaultConfig()
	config.APIKey = apiKey

	// Create generator
	gen, err := llm.NewGenerator(config)
	if err != nil {
		return fmt.Errorf("failed to create LLM generator: %w", err)
	}

	color.Yellow("Calling OpenAI API to generate manifest...")
	fmt.Printf("Intent: %s\n", userIntent)
	fmt.Printf("Model: %s\n\n", config.Model)

	// Generate manifest
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	m, err := gen.GenerateManifest(ctx, diffContent, userIntent)
	if err != nil {
		return fmt.Errorf("failed to generate manifest: %w", err)
	}

	// Show summary
	color.Green("✓ Manifest generated successfully!")
	fmt.Printf("\nSummary:\n")
	fmt.Printf("  Symbols: %d\n", len(m.Entries))
	
	if m.GlobalIntent != nil {
		fmt.Printf("  Global Intent: %s\n", m.GlobalIntent.Rationale)
		fmt.Printf("  Behavior: %v\n", m.GlobalIntent.BehaviorClass)
	}
	
	for i, entry := range m.Entries {
		fmt.Printf("  %d. %s::%s\n", i+1, entry.Anchor.File, entry.Anchor.Symbol)
		if entry.Rationale != "" {
			fmt.Printf("     Rationale: %s\n", entry.Rationale)
		}
	}
	fmt.Println()

	// Save pending manifest
	if err := manifest.SavePending(m); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	color.Green("✓ Manifest saved!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the manifest in .gip/pending-manifest.toon")
	fmt.Println("  2. Run: git commit")

	return nil
}
