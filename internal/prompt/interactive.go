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
	"github.com/hrithikraj/gip/internal/diff"
	"github.com/hrithikraj/gip/internal/manifest"
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
	
	// Select behavior classes
	behaviorClasses, err := promptBehaviorClass()
	if err != nil {
		return err
	}
	
	// Ask if batch mode (if not already set)
	if !batchMode {
		batchMode, err = promptBatchMode()
		if err != nil {
			return err
		}
	}
	
	// Collect contract information
	contract, err := promptContract()
	if err != nil {
		return err
	}
	
	// Collect side effects
	sideEffects, err := promptSideEffects()
	if err != nil {
		return err
	}
	
	// Collect compatibility info
	compat, err := promptCompatibility()
	if err != nil {
		return err
	}
	
	// Feature flags
	featureFlags, err := promptFeatureFlags()
	if err != nil {
		return err
	}
	
	// Rationale
	rationale, err := promptRationale()
	if err != nil {
		return err
	}
	
	// Build manifest entries from changes
	entries := make([]manifest.Entry, len(changes))
	for i, change := range changes {
		entries[i] = diff.ToManifestEntry(change)
		entries[i].BehaviorClass = behaviorClasses
		entries[i].Contract = *contract
		entries[i].SideEffects = sideEffects
		entries[i].Compatibility = *compat
		entries[i].FeatureFlags = append(entries[i].FeatureFlags, featureFlags...)
		entries[i].Rationale = rationale
	}
	
	m := &manifest.Manifest{
		Commit:  "pending",
		Entries: entries,
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
