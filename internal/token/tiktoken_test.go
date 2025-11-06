package token

import (
	"fmt"
	"math"
	"os"
	"strings"
	"testing"
)

func TestTokenCounter_EstimateTokens(t *testing.T) {
	tc := NewTokenCounter()

	tests := []struct {
		name string
		text string
	}{
		{"Simple prose", "The quick brown fox jumps over the lazy dog."},
		{"JSON structure", `{"name": "test", "version": "1.0.0", "dependencies": {}}`},
		{"Code with braces", `func main() { fmt.Println("hello") }`},
		{"Very long word", strings.Repeat("a", 1000)},
		{"Code block with imports", "import (\n\t\"fmt\"\n\t\"os\"\n)\n\nfunc main() {\n\tfmt.Println(\"hello\")\n}"},
		{"Markdown with code", "# Title\n\nSome text with `code` and **bold**.\n\n```go\nfunc main() {}\n```"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := tc.EstimateTokens(tt.text)
			if tokens == 0 {
				t.Fatalf("EstimateTokens returned 0")
			}

			if tc.IsFallbackMode() {
				expected := tc.approximateTokens(tt.text)
				if tokens != expected {
					t.Fatalf("fallback estimate mismatch: got %d, want %d", tokens, expected)
				}
			} else {
				ratio := float64(len(tt.text)) / float64(tokens)
				if ratio < 1.0 || ratio > 6.0 {
					t.Fatalf("unreasonable chars/token ratio %.2f for %q", ratio, tt.text)
				}
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

	if tc.IsFallbackMode() {
		expected := tc.approximateTokens(lockContent)
		if tokens != expected {
			t.Fatalf("fallback large-file mismatch: got %d, want %d", tokens, expected)
		}
		return
	}

	ratio := float64(chars) / float64(tokens)
	if ratio < 2.0 || ratio > 5.0 {
		t.Fatalf("Large JSON ratio %.2f outside reasonable bounds", ratio)
	}

	expectedTokens := float64(chars) / 3.0
	tolerance := expectedTokens * 0.4
	if math.Abs(float64(tokens)-expectedTokens) > tolerance {
		t.Fatalf("Token count %d not within tolerance of %.0f", tokens, expectedTokens)
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
	counter := &TokenCounter{fallbackMode: true, encodingName: "approximation"}

	cases := []struct {
		name string
		text string
		want int
	}{
		{"empty", "", 0},
		{"short prose", "Hello, world!", counter.approximateTokens("Hello, world!")},
		{"code", "func main() { println(\"hello\") }", counter.approximateTokens("func main() { println(\"hello\") }")},
		{"json", `{"key": "value", "number": 123}`, counter.approximateTokens(`{"key": "value", "number": 123}`)},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := counter.EstimateTokens(c.text)
			if got != c.want {
				t.Fatalf("EstimateTokens(%q) = %d, want %d", c.text, got, c.want)
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

func TestIsLikelyCodeThreshold(t *testing.T) {
	lowCode := strings.Repeat("{", 9) + strings.Repeat("a", 991)
	highCode := strings.Repeat("{", 12) + strings.Repeat("a", 988)

	if isLikelyCode(lowCode) {
		t.Fatalf("expected low density to be prose")
	}
	if !isLikelyCode(highCode) {
		t.Fatalf("expected high density to be code")
	}
}

func TestTokenCounterCacheInitialization(t *testing.T) {
	original := os.Getenv("TIKTOKEN_CACHE_DIR")
	t.Setenv("TIKTOKEN_CACHE_DIR", "")

	os.Unsetenv("TIKTOKEN_CACHE_DIR")
	ensureCacheDir()
	if os.Getenv("TIKTOKEN_CACHE_DIR") == "" {
		t.Fatalf("expected cache dir to be set by package init")
	}

	if original != "" {
		os.Setenv("TIKTOKEN_CACHE_DIR", original)
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
