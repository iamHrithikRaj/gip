package llm

import (
	"strings"
)

// ConstructPrompt builds the full LLM prompt with schema and examples
func ConstructPrompt(diff string, userIntent string) string {
	var sb strings.Builder

	// Header
	sb.WriteString("Generate a Gip manifest for this Git commit.\n\n")

	// User intent
	sb.WriteString("## User Intent\n")
	sb.WriteString(userIntent)
	sb.WriteString("\n\n")

	// Schema
	sb.WriteString("## Schema (JSON)\n")
	sb.WriteString(GetSchemaTemplate())
	sb.WriteString("\n\n")

	// Examples
	sb.WriteString("## Examples\n")
	sb.WriteString(GetExamples())
	sb.WriteString("\n\n")

	// Diff
	sb.WriteString("## Git Diff\n")
	sb.WriteString("```diff\n")
	sb.WriteString(diff)
	sb.WriteString("\n```\n\n")

	// Instructions
	sb.WriteString("## Instructions\n")
	sb.WriteString(GetInstructions())

	return sb.String()
}

// GetSchemaTemplate returns the JSON schema template
func GetSchemaTemplate() string {
	return `{
  "schemaVersion": "2.0",
  "commit": "<will-be-filled>",
  "globalIntent": {
    "behaviorClass": ["feature"|"bugfix"|"refactor"|"perf"|"security"|"validation"|"docs"|"config"|"migration"],
    "rationale": "Why this commit exists (commit-level description)"
  },
  "entries": [
    {
      "anchor": {
        "file": "path/to/file.ext",
        "symbol": "functionName",
        "hunk": "H#<line_number>"
      },
      "changeType": "add"|"modify"|"delete"|"rename",
      "signatureDelta": {
        "before": "old signature (if changed)",
        "after": "new signature (if changed)"
      },
      "contract": {
        "preconditions": [
          "What must be true BEFORE execution",
          "Example: 'x must be non-negative'",
          "Example: 'user must be authenticated'"
        ],
        "postconditions": [
          "What is guaranteed AFTER successful execution",
          "Example: 'returns non-null User object'",
          "Example: 'database transaction is committed'"
        ],
        "errorModel": [
          "When errors occur (format: 'ExceptionType when condition')",
          "Example: 'ValueError when x < 0'",
          "Example: 'IOError when file not found'"
        ]
      },
      "behaviorClass": ["feature"|"bugfix"|...],
      "sideEffects": [
        "none (pure function)",
        "logs:info (writes logs)",
        "reads:db (reads database)",
        "writes:cache (modifies cache)",
        "network:http (makes HTTP calls)"
      ],
      "rationale": "Function-specific rationale (if different from global)",
      "compatibility": {
        "breaking": false,
        "deprecations": ["List of deprecated features"],
        "migrations": ["Migration steps if needed"]
      },
      "featureFlags": ["FLAG_NAME"],
      "inheritsGlobalIntent": true|false
    }
  ]
}`
}

// GetExamples returns few-shot examples
func GetExamples() string {
	return `### Example 1: Bugfix (Single Function)
` + "```json" + `
{
  "schemaVersion": "2.0",
  "commit": "fix123",
  "entries": [
    {
      "anchor": {"file": "cart.py", "symbol": "calculate_total", "hunk": "H#42"},
      "changeType": "modify",
      "contract": {
        "preconditions": ["items is non-empty list", "all prices are numeric"],
        "postconditions": ["returns float >= 0"],
        "errorModel": ["ValueError when items empty", "TypeError when price not numeric"]
      },
      "behaviorClass": ["bugfix"],
      "sideEffects": ["none"],
      "rationale": "Fixed edge case where empty cart returned None instead of 0.0"
    }
  ]
}
` + "```" + `

### Example 2: Feature (Multi-Function with Global Intent)
` + "```json" + `
{
  "schemaVersion": "2.0",
  "commit": "feat456",
  "globalIntent": {
    "behaviorClass": ["feature", "validation"],
    "rationale": "Add strict validation across parsing layer"
  },
  "entries": [
    {
      "anchor": {"file": "user.py", "symbol": "parseUser", "hunk": "H#10"},
      "changeType": "modify",
      "signatureDelta": {
        "before": "def parseUser(raw: str) -> User",
        "after": "def parseUser(raw: str, strict: bool = False) -> User"
      },
      "contract": {
        "preconditions": ["raw is non-empty string"],
        "postconditions": ["returns User object with non-null id"],
        "errorModel": ["ValueError when strict=True and id missing"]
      },
      "behaviorClass": ["feature", "validation"],
      "sideEffects": ["none"],
      "rationale": "",
      "inheritsGlobalIntent": true
    },
    {
      "anchor": {"file": "order.py", "symbol": "parseOrder", "hunk": "H#22"},
      "changeType": "modify",
      "contract": {
        "preconditions": ["raw is valid JSON string"],
        "postconditions": ["returns Order object"],
        "errorModel": ["ValueError when required fields missing in strict mode"]
      },
      "behaviorClass": ["feature", "validation"],
      "sideEffects": ["none"],
      "rationale": "",
      "inheritsGlobalIntent": true
    }
  ]
}
` + "```" + `

### Example 3: Refactor
` + "```json" + `
{
  "schemaVersion": "2.0",
  "commit": "refactor789",
  "entries": [
    {
      "anchor": {"file": "utils.py", "symbol": "formatDate", "hunk": "H#15"},
      "changeType": "modify",
      "contract": {
        "preconditions": ["date is datetime object"],
        "postconditions": ["returns ISO 8601 string"],
        "errorModel": ["TypeError when date is not datetime"]
      },
      "behaviorClass": ["refactor", "perf"],
      "sideEffects": ["none"],
      "rationale": "Simplified date formatting logic, 2x performance improvement"
    }
  ]
}
` + "```" + ``
}

// GetInstructions returns generation instructions
func GetInstructions() string {
	return `1. Analyze the git diff to identify all changed functions/symbols
2. For each symbol, extract the file path and approximate line number
3. If multiple functions share the same intent, use globalIntent with inheritsGlobalIntent=true
4. Fill contract fields:
   - preconditions: What must be true BEFORE the function runs
   - postconditions: What is guaranteed AFTER successful execution
   - errorModel: When exceptions occur (format: "ExceptionType when condition")
5. Determine behaviorClass: feature, bugfix, refactor, perf, security, validation, docs, config, migration
6. List sideEffects: none, logs:level, reads:target, writes:target, network:protocol
7. If function signature changed, fill signatureDelta.before and signatureDelta.after
8. Output ONLY valid JSON matching the schema above
9. Do NOT include markdown code blocks (no ` + "```json" + ` markers)
10. Ensure all required fields are present`
}
