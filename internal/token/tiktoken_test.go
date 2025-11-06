package token

import (
	"fmt"
	"strings"
	"testing"
)

func TestTokenCounter_Accuracy(t *testing.T) {
	tc := NewTokenCounter()

	tests := []struct {
		name     string
		text     string
		minRatio float64 // min chars per token (higher = more efficient)
		maxRatio float64 // max chars per token
	}{
		{
			name:     "Simple prose",
			text:     "The quick brown fox jumps over the lazy dog.",
			minRatio: 3.5, // Prose: ~4 chars/token
			maxRatio: 5.0,
		},
		{
			name:     "JSON structure",
			text:     `{"name": "test", "version": "1.0.0", "dependencies": {}}`,
			minRatio: 2.5, // JSON: ~3 chars/token (more punctuation)
			maxRatio: 4.0,
		},
		{
			name:     "Code with braces",
			text:     `func main() { fmt.Println("hello") }`,
			minRatio: 2.5, // Code: ~3.5 chars/token
			maxRatio: 4.5,
		},
		{
			name:     "Very long word",
			text:     strings.Repeat("a", 1000),
			minRatio: 1.0,  // Long words break into many tokens
			maxRatio: 10.0, // Single char repetition has higher ratio
		},
		{
			name:     "Code block with imports",
			text:     "import (\n\t\"fmt\"\n\t\"os\"\n)\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}",
			minRatio: 2.0,
			maxRatio: 4.5,
		},
		{
			name:     "Markdown with code",
			text:     "# Title\n\nSome text with `code` and **bold**.\n\n```go\nfunc main() {}\n```",
			minRatio: 2.5,
			maxRatio: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tc.EstimateTokens(tt.text)
			if tokens == 0 {
				t.Errorf("EstimateTokens returned 0")
				return
			}

			ratio := float64(len(tt.text)) / float64(tokens)
			t.Logf("Text: %d chars, %d tokens, %.2f chars/token (mode: %s)",
				len(tt.text), tokens, ratio, tc.encodingName)

			if ratio < tt.minRatio || ratio > tt.maxRatio {
				t.Errorf("Chars/token ratio %.2f outside expected range [%.2f, %.2f]",
					ratio, tt.minRatio, tt.maxRatio)
			}
		})
	}
}

func TestTokenCounter_LargeFile(t *testing.T) {
	tc := NewTokenCounter()

	// Simulate a package-lock.json structure
	lockContent := `{
  "name": "test-project",
  "version": "1.0.0",
  "lockfileVersion": 2,
  "requires": true,
  "packages": {`

	// Add 1000 dependency entries
	for i := 0; i < 1000; i++ {
		lockContent += fmt.Sprintf(`
    "node_modules/package-%d": {
      "version": "1.2.3",
      "resolved": "https://registry.npmjs.org/package-%d/-/package-%d-1.2.3.tgz",
      "integrity": "sha512-abcdef1234567890abcdef1234567890abcdef1234567890",
      "dependencies": {}
    },`, i, i, i)
	}
	lockContent += "\n  }\n}"

	tokens := tc.EstimateTokens(lockContent)
	chars := len(lockContent)

	t.Logf("Large file: %d chars, %d tokens (%.2f chars/token)",
		chars, tokens, float64(chars)/float64(tokens))

	// JSON should be ~2.8-4.0 chars per token
	ratio := float64(chars) / float64(tokens)
	if ratio < 2.5 || ratio > 4.0 {
		t.Errorf("Large JSON file has unexpected ratio: %.2f (expected 2.5-4.0)", ratio)
	}

	// Sanity check: chars/tokens ratio should be reasonable for JSON
	// Typical JSON: ~3 chars/token, so 100KB ≈ 33K tokens, 250KB ≈ 83K tokens
	expectedTokens := float64(chars) / 3.0
	tolerance := expectedTokens * 0.3 // Allow 30% variance
	if float64(tokens) < (expectedTokens-tolerance) || float64(tokens) > (expectedTokens+tolerance) {
		t.Errorf("Token count %d not within 30%% of expected ~%.0f tokens for %d chars",
			tokens, expectedTokens, chars)
	}
}

func TestTokenCounter_EmptyString(t *testing.T) {
	tc := NewTokenCounter()
	tokens := tc.EstimateTokens("")
	if tokens != 0 {
		t.Errorf("Empty string should be 0 tokens, got %d", tokens)
	}
}

func TestTokenCounter_FallbackMode(t *testing.T) {
	counter := NewTokenCounter()

	t.Logf("Token counter mode: %s (fallback: %v)", counter.GetEncodingName(), counter.IsFallbackMode())

	// Test that we always get reasonable estimates
	testCases := []struct {
		name string
		text string
	}{
		{"empty", ""},
		{"short prose", "Hello, world!"},
		{"code", "func main() { println(\"hello\") }"},
		{"json", `{"key": "value", "number": 123}`},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			tokens := counter.EstimateTokens(testCase.text)

			if testCase.text == "" {
				if tokens != 0 {
					t.Errorf("Empty string should have 0 tokens, got %d", tokens)
				}
			} else {
				if tokens == 0 {
					t.Errorf("Non-empty string should have >0 tokens")
				}

				// Check ratio is reasonable
				ratio := float64(len(testCase.text)) / float64(tokens)
				if ratio < 1.0 || ratio > 6.0 {
					t.Errorf("Unreasonable chars/token ratio: %.2f (text: %d chars, %d tokens)",
						ratio, len(testCase.text), tokens)
				}
			}
		})
	}
}

func TestIsLikelyCode(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected bool
	}{
		{
			name:     "Plain prose",
			text:     "This is a simple sentence with no code.",
			expected: false,
		},
		{
			name:     "Code with braces",
			text:     "func main() { return nil }",
			expected: true,
		},
		{
			name:     "JSON",
			text:     `{"key": "value", "array": [1, 2, 3]}`,
			expected: true,
		},
		{
			name:     "Go code",
			text:     "if x == y {\n\treturn true\n} else {\n\treturn false\n}",
			expected: true,
		},
		{
			name:     "Markdown",
			text:     "# Title\n\nSome text without code symbols",
			expected: false,
		},
		{
			name:     "Empty string",
			text:     "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLikelyCode(tt.text)
			if got != tt.expected {
				t.Errorf("isLikelyCode() = %v, want %v for text: %q", got, tt.expected, tt.text)
			}
		})
	}
}

func TestApproximateTokens(t *testing.T) {
	tc := &TokenCounter{
		encoding:     nil,
		fallbackMode: true,
		encodingName: "approximation",
	}

	tests := []struct {
		name     string
		text     string
		minRatio float64
		maxRatio float64
	}{
		{
			name:     "Prose text",
			text:     "The quick brown fox jumps over the lazy dog. This is a longer piece of text.",
			minRatio: 3.5,
			maxRatio: 5.0,
		},
		{
			name:     "Code",
			text:     "func main() {\n\tfmt.Println(\"Hello, World!\")\n\treturn nil\n}",
			minRatio: 2.5,
			maxRatio: 4.5,
		},
		{
			name:     "Mixed content",
			text:     "# Documentation\n\nThis function does something:\n\n```go\nfunc doThing() {}\n```",
			minRatio: 2.5,
			maxRatio: 5.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tc.approximateTokens(tt.text)
			if tokens == 0 {
				t.Errorf("approximateTokens returned 0")
				return
			}

			ratio := float64(len(tt.text)) / float64(tokens)
			t.Logf("Approximation: %d chars, %d tokens, %.2f chars/token",
				len(tt.text), tokens, ratio)

			if ratio < tt.minRatio || ratio > tt.maxRatio {
				t.Errorf("Chars/token ratio %.2f outside expected range [%.2f, %.2f]",
					ratio, tt.minRatio, tt.maxRatio)
			}
		})
	}
}

func TestDebugTokenCount(t *testing.T) {
	tc := NewTokenCounter()

	// This test just verifies that DebugTokenCount runs without error
	// Actual debug output is only shown when debug logging is enabled
	text := "func main() { fmt.Println(\"hello\") }"
	tokens := tc.DebugTokenCount(text, "test-function")

	if tokens == 0 {
		t.Errorf("DebugTokenCount returned 0 tokens")
	}

	// Test with empty string
	emptyTokens := tc.DebugTokenCount("", "empty-test")
	if emptyTokens != 0 {
		t.Errorf("DebugTokenCount(\"\") returned %d tokens, want 0", emptyTokens)
	}

	// Test with longer text that will be previewed
	longText := strings.Repeat("This is a test sentence that will be used to create a very long text. ", 10)
	longTokens := tc.DebugTokenCount(longText, "long-test")
	if longTokens == 0 {
		t.Errorf("DebugTokenCount with long text returned 0 tokens")
	}
}

func TestNewTokenCounter_EdgeCases(t *testing.T) {
	// Just verify NewTokenCounter returns a usable counter
	tc := NewTokenCounter()

	if tc == nil {
		t.Fatal("NewTokenCounter returned nil")
	}

	if tc.encodingName == "" {
		t.Error("NewTokenCounter returned counter with empty encoding name")
	}

	// Verify it works regardless of fallback mode
	tokens := tc.EstimateTokens("test")
	if tokens == 0 {
		t.Error("Counter should estimate >0 tokens for non-empty text")
	}
}

func TestTokenCounter_GetEncodingName(t *testing.T) {
	tc := NewTokenCounter()
	name := tc.GetEncodingName()

	if name == "" {
		t.Error("GetEncodingName returned empty string")
	}

	// Should be either "cl100k_base" or "approximation"
	if name != "cl100k_base" && name != "approximation" {
		t.Errorf("GetEncodingName returned unexpected value: %q", name)
	}
}

func TestTokenCounter_IsFallbackMode(t *testing.T) {
	tc := NewTokenCounter()

	// Just verify the method works
	fallback := tc.IsFallbackMode()

	// Should match encoding name
	if fallback && tc.GetEncodingName() != "approximation" {
		t.Error("IsFallbackMode is true but encoding name is not 'approximation'")
	}
	if !fallback && tc.GetEncodingName() != "cl100k_base" {
		t.Error("IsFallbackMode is false but encoding name is not 'cl100k_base'")
	}
}
