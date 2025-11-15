// Package llm provides LLM-based manifest generation capabilities.
// It handles API calls to OpenAI, Anthropic, and other providers to automatically
// generate structured TOON manifests from git diffs and user intent.
package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/iamHrithikRaj/gip/internal/manifest"
)

// Provider represents an LLM provider
type Provider string

const (
	ProviderOpenAI    Provider = "openai"
	ProviderAnthropic Provider = "anthropic"
	ProviderOllama    Provider = "ollama"
)

// Config holds LLM configuration
type Config struct {
	Provider    Provider
	Model       string
	APIKey      string
	Temperature float64
	MaxTokens   int
	Timeout     time.Duration
}

// DefaultConfig returns sensible defaults
func DefaultConfig() Config {
	return Config{
		Provider:    ProviderOpenAI,
		Model:       "gpt-4o-mini",
		Temperature: 0.2,
		MaxTokens:   2000,
		Timeout:     30 * time.Second,
	}
}

// Generator handles LLM-based manifest generation
type Generator struct {
	config Config
	client Client // Interface for different LLM providers
}

// Client is the interface for LLM API calls
type Client interface {
	Generate(ctx context.Context, prompt string) (string, error)
}

// NewGenerator creates a new LLM generator
func NewGenerator(config Config) (*Generator, error) {
	var client Client
	var err error

	switch config.Provider {
	case ProviderOpenAI:
		client, err = NewOpenAIClient(config)
	case ProviderAnthropic:
		return nil, fmt.Errorf("anthropic provider not yet implemented")
	case ProviderOllama:
		return nil, fmt.Errorf("ollama provider not yet implemented")
	default:
		return nil, fmt.Errorf("unknown provider: %s", config.Provider)
	}

	if err != nil {
		return nil, err
	}

	return &Generator{
		config: config,
		client: client,
	}, nil
}

// GenerateManifest generates a manifest from a git diff and user intent
func (g *Generator) GenerateManifest(ctx context.Context, diff string, userIntent string) (*manifest.Manifest, error) {
	const maxRetries = 3

	for attempt := 1; attempt <= maxRetries; attempt++ {
		prompt := ConstructPrompt(diff, userIntent)

		response, err := g.client.Generate(ctx, prompt)
		if err != nil {
			return nil, fmt.Errorf("LLM API error (attempt %d/%d): %w", attempt, maxRetries, err)
		}

		// Try to parse and validate the response
		m, err := ValidateManifest(response)
		if err == nil {
			return m, nil
		}

		// If validation failed and we have retries left, add error context
		if attempt < maxRetries {
			userIntent = fmt.Sprintf("%s\n\nPrevious attempt failed with error: %v\nPlease fix the JSON and try again.", userIntent, err)
		}
	}

	return nil, fmt.Errorf("failed to generate valid manifest after %d retries", maxRetries)
}

// ValidateManifest parses and validates an LLM-generated manifest
func ValidateManifest(raw string) (*manifest.Manifest, error) {
	// Strip markdown code blocks if present
	raw = stripMarkdown(raw)

	var m manifest.Manifest
	if err := json.Unmarshal([]byte(raw), &m); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	// Basic validation
	if m.Commit == "" {
		return nil, fmt.Errorf("missing commit field")
	}

	if len(m.Entries) == 0 {
		return nil, fmt.Errorf("manifest must have at least one entry")
	}

	// Validate each entry
	for i, entry := range m.Entries {
		if entry.Anchor.File == "" {
			return nil, fmt.Errorf("entry %d: missing file in anchor", i)
		}
		if entry.Anchor.Symbol == "" {
			return nil, fmt.Errorf("entry %d: missing symbol in anchor", i)
		}
		if len(entry.BehaviorClass) == 0 && (m.GlobalIntent == nil || !entry.InheritsGlobalIntent) {
			return nil, fmt.Errorf("entry %d: must have behaviorClass or inherit from globalIntent", i)
		}
	}

	return &m, nil
}

// stripMarkdown removes markdown code blocks from LLM responses
func stripMarkdown(raw string) string {
	// Remove ```json ... ``` blocks
	if len(raw) > 7 && raw[:7] == "```json" {
		end := len(raw)
		if raw[end-3:] == "```" {
			raw = raw[7 : end-3]
		}
	}
	// Remove ``` ... ``` blocks
	if len(raw) > 3 && raw[:3] == "```" {
		end := len(raw)
		if raw[end-3:] == "```" {
			raw = raw[3 : end-3]
		}
	}
	return raw
}
