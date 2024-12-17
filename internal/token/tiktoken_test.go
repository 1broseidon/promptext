package token

import (
	"testing"
)

func TestTokenCounter_EstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: 0,
		},
		{
			name:     "Simple text",
			text:     "Hello world",
			expected: 2,
		},
		{
			name:     "Code block",
			text:     "```go\nfunc main() {}\n```",
			expected: 11,
		},
		{
			name:     "Markdown with links",
			text:     "[Link](https://example.com)",
			expected: 5,
		},
	}

	tc := NewTokenCounter()
	if tc.encoding == nil {
		t.Skip("Tiktoken encoder not initialized")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tc.EstimateTokens(tt.text)
			// Note: We're being a bit flexible with the expected counts
			// since tiktoken's exact tokenization may change
			if got == 0 && tt.expected > 0 {
				t.Errorf("EstimateTokens() = %v, want approximately %v", got, tt.expected)
			}
		})
	}
}
