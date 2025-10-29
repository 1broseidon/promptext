package token

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/log"
	"github.com/pkoukk/tiktoken-go"
)

func init() {
	// Set default cache directory if TIKTOKEN_CACHE_DIR is not set
	if os.Getenv("TIKTOKEN_CACHE_DIR") == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Debug("Warning: Could not get user home directory: %v", err)
			return
		}

		cacheDir := filepath.Join(homeDir, ".promptext", "cache")
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			log.Debug("Warning: Could not create cache directory: %v", err)
			return
		}

		os.Setenv("TIKTOKEN_CACHE_DIR", cacheDir)
		log.Debug("Set tiktoken cache to: %s", cacheDir)
	}
}

type TokenCounter struct {
	encoding     *tiktoken.Tiktoken
	fallbackMode bool
	encodingName string
}

// NewTokenCounter creates a token counter with proper fallback
func NewTokenCounter() *TokenCounter {
	// Try cl100k_base (GPT-4, GPT-3.5-turbo)
	enc, err := tiktoken.GetEncoding("cl100k_base")
	if err != nil {
		log.Debug("Failed to load cl100k_base encoding: %v", err)
		log.Info("Token counting using approximation (tiktoken unavailable)")
		return &TokenCounter{
			encoding:     nil,
			fallbackMode: true,
			encodingName: "approximation",
		}
	}

	log.Debug("Initialized tiktoken with cl100k_base encoding")
	return &TokenCounter{
		encoding:     enc,
		fallbackMode: false,
		encodingName: "cl100k_base",
	}
}

// EstimateTokens counts tokens using tiktoken or falls back to approximation
func (tc *TokenCounter) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	if tc.fallbackMode || tc.encoding == nil {
		return tc.approximateTokens(text)
	}

	// Use tiktoken for accurate count
	tokens := tc.encoding.Encode(text, nil, nil)
	return len(tokens)
}

// approximateTokens provides a fallback when tiktoken is unavailable
// Uses a more sophisticated approximation than simple char/4
func (tc *TokenCounter) approximateTokens(text string) int {
	// More accurate approximation based on empirical GPT-4 tokenization:
	// - Whitespace-separated words: ~1.3 tokens per word on average
	// - Punctuation and special chars add tokens
	// - Code has higher token density than prose

	wordCount := len(strings.Fields(text))
	charCount := len(text)

	// Base estimate: 1.3 tokens per word
	estimate := float64(wordCount) * 1.3

	// Adjust for code characteristics
	if isLikelyCode(text) {
		// Code has more special chars, brackets, dots - higher token density
		// Use ~3.5 chars per token for code
		codeEstimate := float64(charCount) / 3.5

		// Take average of word-based and char-based estimates
		estimate = (estimate + codeEstimate) / 2
	} else {
		// Prose: ~4 chars per token
		proseEstimate := float64(charCount) / 4.0
		estimate = (estimate + proseEstimate) / 2
	}

	return int(estimate)
}

// isLikelyCode determines if text is likely code vs prose
func isLikelyCode(text string) bool {
	// Count code indicators
	indicators := 0

	codeChars := []string{"{", "}", "(", ")", "[", "]", ";", "=>", "->", "==", "!="}
	for _, char := range codeChars {
		count := strings.Count(text, char)
		indicators += count
	}

	// If we see 10+ code chars per 1000 characters, it's likely code
	if len(text) == 0 {
		return false
	}
	codeCharDensity := float64(indicators) / float64(len(text)) * 1000
	return codeCharDensity > 10
}

// GetEncodingName returns the name of the encoding being used
func (tc *TokenCounter) GetEncodingName() string {
	return tc.encodingName
}

// IsFallbackMode returns true if using approximation instead of tiktoken
func (tc *TokenCounter) IsFallbackMode() bool {
	return tc.fallbackMode
}

// DebugTokenCount provides detailed breakdown of token estimation
func (tc *TokenCounter) DebugTokenCount(text string, label string) int {
	tokens := tc.EstimateTokens(text)

	if log.IsDebugEnabled() {
		chars := len(text)
		words := len(strings.Fields(text))
		lines := strings.Count(text, "\n") + 1

		log.Debug("Token estimation for %s:", label)
		log.Debug("  Characters: %d", chars)
		log.Debug("  Words: %d", words)
		log.Debug("  Lines: %d", lines)
		log.Debug("  Tokens: %d", tokens)
		if tokens > 0 {
			log.Debug("  Chars/token: %.2f", float64(chars)/float64(tokens))
		}
		log.Debug("  Method: %s", tc.encodingName)

		// Show first 100 chars for context
		preview := text
		if len(preview) > 100 {
			preview = preview[:100] + "..."
		}
		log.Debug("  Preview: %s", strings.ReplaceAll(preview, "\n", "\\n"))
	}

	return tokens
}
