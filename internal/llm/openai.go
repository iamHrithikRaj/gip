package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAIClient handles OpenAI API calls
type OpenAIClient struct {
	apiKey      string
	model       string
	temperature float64
	maxTokens   int
	timeout     time.Duration
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config Config) (*OpenAIClient, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("OpenAI API key is required (set OPENAI_API_KEY env var)")
	}

	return &OpenAIClient{
		apiKey:      config.APIKey,
		model:       config.Model,
		temperature: config.Temperature,
		maxTokens:   config.MaxTokens,
		timeout:     config.Timeout,
	}, nil
}

// Generate calls the OpenAI API
func (c *OpenAIClient) Generate(ctx context.Context, prompt string) (string, error) {
	reqBody := map[string]interface{}{
		"model": c.model,
		"messages": []map[string]string{
			{
				"role":    "system",
				"content": "You are a code analysis expert. Generate valid JSON manifests for Git commits. Output only valid JSON, no markdown formatting.",
			},
			{
				"role":    "user",
				"content": prompt,
			},
		},
		"temperature": c.temperature,
		"max_tokens":  c.maxTokens,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	client := &http.Client{Timeout: c.timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Error struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	if result.Error.Message != "" {
		return "", fmt.Errorf("API error: %s", result.Error.Message)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from API")
	}

	return result.Choices[0].Message.Content, nil
}
