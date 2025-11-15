package llm

import (
	"context"
	"testing"
)

// MockClient is a mock LLM client for testing
type MockClient struct {
	Response string
	Error    error
}

func (m *MockClient) Generate(ctx context.Context, prompt string) (string, error) {
	if m.Error != nil {
		return "", m.Error
	}
	return m.Response, nil
}

func TestValidateManifest(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name: "valid manifest",
			input: `{
				"schemaVersion": "2.0",
				"commit": "abc123",
				"entries": [
					{
						"anchor": {"file": "test.py", "symbol": "foo", "hunk": "H#10"},
						"changeType": "modify",
						"contract": {
							"preconditions": ["x > 0"],
							"postconditions": ["returns int"],
							"errorModel": ["ValueError when x < 0"]
						},
						"behaviorClass": ["feature"],
						"sideEffects": ["none"],
						"rationale": "test"
					}
				]
			}`,
			wantErr: false,
		},
		{
			name: "manifest with markdown code block",
			input: "```json\n" + `{
				"schemaVersion": "2.0",
				"commit": "abc123",
				"entries": [
					{
						"anchor": {"file": "test.py", "symbol": "foo", "hunk": "H#10"},
						"changeType": "modify",
						"contract": {
							"preconditions": ["x > 0"],
							"postconditions": ["returns int"],
							"errorModel": ["ValueError when x < 0"]
						},
						"behaviorClass": ["feature"],
						"sideEffects": ["none"],
						"rationale": "test"
					}
				]
			}` + "\n```",
			wantErr: false,
		},
		{
			name: "manifest with global intent",
			input: `{
				"schemaVersion": "2.0",
				"commit": "abc123",
				"globalIntent": {
					"behaviorClass": ["feature"],
					"rationale": "Add validation"
				},
				"entries": [
					{
						"anchor": {"file": "test.py", "symbol": "foo", "hunk": "H#10"},
						"changeType": "modify",
						"contract": {
							"preconditions": [],
							"postconditions": [],
							"errorModel": []
						},
						"behaviorClass": ["feature"],
						"sideEffects": ["none"],
						"inheritsGlobalIntent": true
					}
				]
			}`,
			wantErr: false,
		},
		{
			name:    "invalid JSON",
			input:   `{invalid json}`,
			wantErr: true,
		},
		{
			name: "missing commit",
			input: `{
				"entries": []
			}`,
			wantErr: true,
		},
		{
			name: "empty entries",
			input: `{
				"commit": "abc123",
				"entries": []
			}`,
			wantErr: true,
		},
		{
			name: "missing file in anchor",
			input: `{
				"commit": "abc123",
				"entries": [
					{
						"anchor": {"symbol": "foo"},
						"contract": {
							"preconditions": [],
							"postconditions": [],
							"errorModel": []
						},
						"behaviorClass": ["feature"],
						"sideEffects": ["none"]
					}
				]
			}`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ValidateManifest(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateManifest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStripMarkdown(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "json code block",
			input: "```json\n{\"test\": true}\n```",
			want:  "\n{\"test\": true}\n",
		},
		{
			name:  "generic code block",
			input: "```\n{\"test\": true}\n```",
			want:  "\n{\"test\": true}\n",
		},
		{
			name:  "no code block",
			input: "{\"test\": true}",
			want:  "{\"test\": true}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripMarkdown(tt.input)
			if got != tt.want {
				t.Errorf("stripMarkdown() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestGeneratorWithMock(t *testing.T) {
	mockResponse := `{
		"schemaVersion": "2.0",
		"commit": "test123",
		"entries": [
			{
				"anchor": {"file": "test.py", "symbol": "foo", "hunk": "H#10"},
				"changeType": "modify",
				"contract": {
					"preconditions": ["x > 0"],
					"postconditions": ["returns int"],
					"errorModel": ["ValueError when x < 0"]
				},
				"behaviorClass": ["feature"],
				"sideEffects": ["none"],
				"rationale": "test function"
			}
		]
	}`

	gen := &Generator{
		config: DefaultConfig(),
		client: &MockClient{Response: mockResponse},
	}

	ctx := context.Background()
	manifest, err := gen.GenerateManifest(ctx, "diff content", "Add feature")
	if err != nil {
		t.Fatalf("GenerateManifest() error = %v", err)
	}

	if manifest.Commit != "test123" {
		t.Errorf("got commit %q, want 'test123'", manifest.Commit)
	}

	if len(manifest.Entries) != 1 {
		t.Errorf("got %d entries, want 1", len(manifest.Entries))
	}
}

func TestConstructPrompt(t *testing.T) {
	diff := "@@ -1,3 +1,4 @@\n def foo():\n+    print('hello')\n     pass"
	intent := "Add logging"

	prompt := ConstructPrompt(diff, intent)

	// Verify key sections are present
	if !contains(prompt, "User Intent") {
		t.Error("prompt missing 'User Intent' section")
	}
	if !contains(prompt, "Schema") {
		t.Error("prompt missing 'Schema' section")
	}
	if !contains(prompt, "Examples") {
		t.Error("prompt missing 'Examples' section")
	}
	if !contains(prompt, "Git Diff") {
		t.Error("prompt missing 'Git Diff' section")
	}
	if !contains(prompt, "Instructions") {
		t.Error("prompt missing 'Instructions' section")
	}
	if !contains(prompt, intent) {
		t.Errorf("prompt missing user intent: %q", intent)
	}
	if !contains(prompt, diff) {
		t.Error("prompt missing diff content")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && findSubstring(s, substr))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
