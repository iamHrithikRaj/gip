package diff

import (
	"testing"
)

func TestDetectBreakingChange(t *testing.T) {
	tests := []struct {
		name     string
		before   string
		after    string
		breaking bool
	}{
		{
			name:     "no change",
			before:   "def foo(x: int) -> str",
			after:    "def foo(x: int) -> str",
			breaking: false,
		},
		{
			name:     "add optional parameter",
			before:   "def foo(x: int) -> str",
			after:    "def foo(x: int, y: int = 0) -> str",
			breaking: false,
		},
		{
			name:     "add required parameter",
			before:   "def foo(x: int) -> str",
			after:    "def foo(x: int, y: int) -> str",
			breaking: true,
		},
		{
			name:     "remove parameter",
			before:   "def foo(x: int, y: int) -> str",
			after:    "def foo(x: int) -> str",
			breaking: true,
		},
		{
			name:     "change return type",
			before:   "def foo(x: int) -> str",
			after:    "def foo(x: int) -> Optional[str]",
			breaking: true,
		},
		{
			name:     "go function - add parameter with default",
			before:   "func Foo(x int) string",
			after:    "func Foo(x int, y int) string",
			breaking: true, // Go doesn't have default params, so this is breaking
		},
		{
			name:     "typescript - add optional parameter",
			before:   "function foo(x: number): string",
			after:    "function foo(x: number, y?: number): string",
			breaking: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectBreakingChange(tt.before, tt.after)
			if result != tt.breaking {
				t.Errorf("DetectBreakingChange() = %v, want %v", result, tt.breaking)
			}
		})
	}
}

func TestExtractParameters(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		want      []string
	}{
		{
			name:      "python function",
			signature: "def foo(x: int, y: str) -> bool",
			want:      []string{"x: int", "y: str"},
		},
		{
			name:      "with self",
			signature: "def foo(self, x: int) -> bool",
			want:      []string{"x: int"},
		},
		{
			name:      "no parameters",
			signature: "def foo() -> bool",
			want:      []string{},
		},
		{
			name:      "go function",
			signature: "func Foo(x int, y string) bool",
			want:      []string{"x int", "y string"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractParameters(tt.signature)
			if len(result) != len(tt.want) {
				t.Errorf("extractParameters() = %v, want %v", result, tt.want)
			}
		})
	}
}

func TestExtractReturnType(t *testing.T) {
	tests := []struct {
		name      string
		signature string
		want      string
	}{
		{
			name:      "python with arrow",
			signature: "def foo(x: int) -> str",
			want:      "str",
		},
		{
			name:      "python optional",
			signature: "def foo(x: int) -> Optional[str]",
			want:      "Optional[str]",
		},
		{
			name:      "go function",
			signature: "func Foo(x int) string",
			want:      "string",
		},
		{
			name:      "typescript",
			signature: "function foo(x: number): string",
			want:      "string",
		},
		{
			name:      "no return type",
			signature: "def foo(x: int)",
			want:      "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractReturnType(tt.signature)
			if result != tt.want {
				t.Errorf("extractReturnType() = %q, want %q", result, tt.want)
			}
		})
	}
}

func TestDetectFeatureFlags(t *testing.T) {
	tests := []struct {
		name string
		code string
		want []string
	}{
		{
			name: "enable flags",
			code: "if ENABLE_NEW_PARSER:\n    use_new_parser()",
			want: []string{"NEW_PARSER"},
		},
		{
			name: "flag prefix",
			code: "if FLAG_STRICT_MODE:\n    validate()",
			want: []string{"STRICT_MODE"},
		},
		{
			name: "feature flags",
			code: "if FEATURE_V2_API:\n    call_v2()",
			want: []string{"V2_API"},
		},
		{
			name: "no flags",
			code: "def foo():\n    return bar()",
			want: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := DetectFeatureFlags(tt.code)
			if len(result) != len(tt.want) {
				t.Errorf("DetectFeatureFlags() = %v, want %v", result, tt.want)
			}
		})
	}
}
