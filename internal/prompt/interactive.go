// Package prompt provides interactive CLI prompts for collecting manifest metadata
// from developers during the commit process.
//
// It uses survey/v2 to present user-friendly prompts for change type, scope,
// description, impact, and related changes, with intelligent defaults based on
// git diff analysis.
package prompt

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/iamHrithikRaj/gip/internal/diff"
	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// InteractiveCommit runs the interactive commit flow
func InteractiveCommit(batchMode bool) error {
	color.Cyan("Analyzing staged changes...")

	// Analyze git diff to extract symbols
	changes, err := diff.AnalyzeStagedChanges()
	if err != nil {
		return fmt.Errorf("failed to analyze changes: %w", err)
	}

	if len(changes) == 0 {
		color.Yellow("No staged changes found. Stage changes first with: git add")
		return nil
	}

	fmt.Printf("\nFound %d modified symbols:\n", len(changes))
	for i, change := range changes {
		fmt.Printf("  %d. %s::%s (%d lines)\n", i+1, change.File, change.Symbol, change.LinesChanged)
	}
	fmt.Println()

	// Check if this is a multi-function commit
	var globalIntent *manifest.GlobalIntent
	isMultiFunction := len(changes) > 1
	
	if isMultiFunction {
		useGlobalIntent, err := promptUseGlobalIntent(len(changes))
		if err != nil {
			return err
		}
		
		if useGlobalIntent {
			globalIntent, err = promptGlobalIntent()
			if err != nil {
				return err
			}
		}
	}

	// Ask if batch mode (if not already set)
	if !batchMode {
		batchMode, err := promptBatchMode()
		if err != nil {
			return err
		}
		// Note: batchMode is intentionally not used after this point
		// It was kept for future batch processing features
		_ = batchMode
	}

	// If no global intent, prompt for single-function behavior class and rationale
	var behaviorClasses []string
	var rationale string
	
	if globalIntent == nil {
		behaviorClasses, err = promptBehaviorClass()
		if err != nil {
			return err
		}
		
		rationale, err = promptRationale()
		if err != nil {
			return err
		}
	}

	// Collect contract information (shared across all functions)
	contract, err := promptContract()
	if err != nil {
		return err
	}

	// Collect side effects (shared)
	sideEffects, err := promptSideEffects()
	if err != nil {
		return err
	}

	// Collect compatibility info (shared)
	compat, err := promptCompatibility()
	if err != nil {
		return err
	}

	// Feature flags (shared)
	featureFlags, err := promptFeatureFlags()
	if err != nil {
		return err
	}

	// Build manifest entries from changes
	entries := make([]manifest.Entry, len(changes))
	for i, change := range changes {
		entries[i] = diff.ToManifestEntry(change)
		
		// Ask if this entry should inherit global intent
		if globalIntent != nil {
			inheritGlobal, err := promptInheritGlobalIntent(change.Symbol)
			if err != nil {
				return err
			}
			
			entries[i].InheritsGlobalIntent = inheritGlobal
			
			if inheritGlobal {
				// Inherit behavior class from global intent
				entries[i].BehaviorClass = globalIntent.BehaviorClass
				entries[i].Rationale = "" // Empty rationale means use global
			} else {
				// Prompt for unique behavior class and rationale
				uniqueClasses, err := promptBehaviorClass()
				if err != nil {
					return err
				}
				uniqueRationale, err := promptRationale()
				if err != nil {
					return err
				}
				entries[i].BehaviorClass = uniqueClasses
				entries[i].Rationale = uniqueRationale
			}
		} else {
			// No global intent, use prompts from earlier
			entries[i].BehaviorClass = behaviorClasses
			entries[i].Rationale = rationale
		}
		
		entries[i].Contract = *contract
		entries[i].SideEffects = sideEffects
		entries[i].Compatibility = *compat
		entries[i].FeatureFlags = append(entries[i].FeatureFlags, featureFlags...)
	}

	m := &manifest.Manifest{
		Commit:       "pending",
		GlobalIntent: globalIntent,
		Entries:      entries,
	}

	// Save pending manifest
	if err := manifest.SavePending(m); err != nil {
		return fmt.Errorf("failed to save manifest: %w", err)
	}

	color.Green("✓ Manifest created successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review the manifest in .gip/pending-manifest.toon")
	fmt.Println("  2. Run: git commit")

	return nil
}

func promptBehaviorClass() ([]string, error) {
	color.Yellow("\n[Behavior Class]")
	fmt.Println("What kind of change is this?")

	options := []string{
		"bugfix     - Fixes incorrect behavior",
		"feature    - Adds new capability",
		"refactor   - Code restructuring",
		"perf       - Performance optimization",
		"security   - Security hardening",
		"validation - Input validation",
		"docs       - Documentation",
		"config     - Configuration change",
		"migration  - Data/API migration",
	}

	var selected []int
	prompt := &survey.MultiSelect{
		Message: "Select one or more:",
		Options: options,
	}

	if err := survey.AskOne(prompt, &selected); err != nil {
		return nil, err
	}

	// Extract behavior class names
	classes := make([]string, len(selected))
	allClasses := manifest.AllBehaviorClasses()
	for i, idx := range selected {
		classes[i] = allClasses[idx]
	}

	return classes, nil
}

func promptBatchMode() (bool, error) {
	var batch bool
	prompt := &survey.Confirm{
		Message: "Apply same contract to all symbols?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &batch); err != nil {
		return false, err
	}
	return batch, nil
}

func promptContract() (*manifest.Contract, error) {
	contract := &manifest.Contract{}

	color.Yellow("\n[Contract - Preconditions]")
	fmt.Println("What must be true BEFORE this function is called?")
	fmt.Println("Examples:")
	fmt.Println("  - \"items must be a non-empty list\"")
	fmt.Println("  - \"user must be authenticated\"")
	fmt.Println()

	contract.Preconditions = promptMultiLine("Enter preconditions (empty line to finish):")

	color.Yellow("\n[Contract - Postconditions]")
	fmt.Println("What is GUARANTEED after successful execution?")
	fmt.Println("Examples:")
	fmt.Println("  - \"returns numeric total >= 0\"")
	fmt.Println("  - \"User object has non-null id\"")
	fmt.Println()

	contract.Postconditions = promptMultiLine("Enter postconditions:")

	color.Yellow("\n[Contract - Error Model]")
	fmt.Println("When do errors occur? (format: \"ExceptionType when condition\")")
	fmt.Println("Examples:")
	fmt.Println("  - \"ValueError when discount < 0\"")
	fmt.Println("  - \"TypeError when items not iterable\"")
	fmt.Println()

	contract.ErrorModel = promptMultiLine("Enter error conditions:")

	return contract, nil
}

func promptMultiLine(message string) []string {
	var results []string

	for {
		var input string
		prompt := &survey.Input{
			Message: ">",
		}
		if err := survey.AskOne(prompt, &input); err != nil {
			break
		}

		input = strings.TrimSpace(input)
		if input == "" {
			break
		}
		results = append(results, input)
	}

	return results
}

func promptSideEffects() ([]string, error) {
	color.Yellow("\n[Side Effects]")
	fmt.Println("Any external effects? (format: action:target)")
	fmt.Println("Examples: logs:channel, reads:file, writes:db, env:VAR, net:host")

	return promptMultiLine("Enter side effects (empty to skip):"), nil
}

func promptCompatibility() (*manifest.Compatibility, error) {
	color.Yellow("\n[Compatibility]")

	compat := &manifest.Compatibility{}

	prompt1 := &survey.Confirm{
		Message: "Binary breaking (old compiled code breaks)?",
		Default: false,
	}
	survey.AskOne(prompt1, &compat.BinaryBreaking)

	prompt2 := &survey.Confirm{
		Message: "Source breaking (old source won't compile)?",
		Default: false,
	}
	survey.AskOne(prompt2, &compat.SourceBreaking)

	prompt3 := &survey.Confirm{
		Message: "Data migration needed?",
		Default: false,
	}
	survey.AskOne(prompt3, &compat.DataModelMigration)

	return compat, nil
}

func promptFeatureFlags() ([]string, error) {
	color.Yellow("\n[Feature Flags]")
	fmt.Println("Any feature flags? (comma-separated, or empty)")

	var input string
	prompt := &survey.Input{
		Message: "Flags:",
	}
	if err := survey.AskOne(prompt, &input); err != nil {
		return nil, err
	}

	if input == "" {
		return []string{}, nil
	}

	flags := strings.Split(input, ",")
	for i := range flags {
		flags[i] = strings.TrimSpace(flags[i])
	}

	return flags, nil
}

func promptRationale() (string, error) {
	color.Yellow("\n[Rationale]")
	fmt.Println("One-line description of what this guarantees:")

	var rationale string
	prompt := &survey.Input{
		Message: "Rationale:",
	}
	if err := survey.AskOne(prompt, &rationale); err != nil {
		return "", err
	}

	return rationale, nil
}

func promptUseGlobalIntent(numSymbols int) (bool, error) {
	color.Cyan(fmt.Sprintf("\n[Multi-Function Commit Detected: %d symbols]", numSymbols))
	fmt.Println("Do all these functions share the same intent?")
	fmt.Println("(If yes, you'll describe the intent once and each function can inherit it)")

	var useGlobal bool
	prompt := &survey.Confirm{
		Message: "Use global intent?",
		Default: true,
	}
	if err := survey.AskOne(prompt, &useGlobal); err != nil {
		return false, err
	}
	return useGlobal, nil
}

func promptGlobalIntent() (*manifest.GlobalIntent, error) {
	color.Yellow("\n[Global Intent - Commit Level]")
	fmt.Println("Describe the overall intent of this commit (shared across all functions)")

	// Prompt for behavior classes
	classes, err := promptBehaviorClass()
	if err != nil {
		return nil, err
	}

	// Prompt for rationale
	color.Yellow("\n[Global Rationale]")
	fmt.Println("One-line description of why this commit exists:")

	var rationale string
	rationalePrompt := &survey.Input{
		Message: "Global rationale:",
	}
	if err := survey.AskOne(rationalePrompt, &rationale); err != nil {
		return nil, err
	}

	return &manifest.GlobalIntent{
		BehaviorClass: classes,
		Rationale:     rationale,
	}, nil
}

func promptInheritGlobalIntent(symbolName string) (bool, error) {
	color.Yellow(fmt.Sprintf("\n[Function: %s]", symbolName))
	
	var inherit bool
	prompt := &survey.Confirm{
		Message: "Inherit global intent? (or provide unique intent for this function)",
		Default: true,
	}
	if err := survey.AskOne(prompt, &inherit); err != nil {
		return false, err
	}
	return inherit, nil
}
