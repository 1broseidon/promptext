package token

import (
	"regexp"
	"strings"
	"unicode"
)

// TokenCounter provides functionality to estimate token counts for LLM input
type TokenCounter struct {
	// Pre-compiled patterns for token splitting
	wordPattern    *regexp.Regexp
	numberPattern  *regexp.Regexp
	symbolPattern  *regexp.Regexp
	spacingPattern *regexp.Regexp
}

// NewTokenCounter creates a new TokenCounter instance
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{
		wordPattern:    regexp.MustCompile(`\w+`),
		numberPattern:  regexp.MustCompile(`\d+`),
		symbolPattern:  regexp.MustCompile(`[^\w\s]`),
		spacingPattern: regexp.MustCompile(`\s+`),
	}
}

// EstimateTokens provides a reasonable estimation of token count for LLM input
// This is an approximation based on common tokenization patterns
func (tc *TokenCounter) EstimateTokens(text string) int {
	if text == "" {
		return 0
	}

	// Track total tokens
	tokenCount := 0

	// Split text into lines to handle code blocks and markdown separately
	lines := strings.Split(text, "\n")
	inCodeBlock := false

	for _, line := range lines {
		// Check for code block delimiters
		if strings.HasPrefix(line, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}

		if inCodeBlock {
			// Code tends to be tokenized more granularly
			tokenCount += tc.countCodeTokens(line)
		} else {
			// Regular text/markdown
			tokenCount += tc.countTextTokens(line)
		}
	}

	return tokenCount
}

// countTextTokens handles regular text and markdown
func (tc *TokenCounter) countTextTokens(line string) int {
	count := 0

	// Count words
	words := tc.wordPattern.FindAllString(line, -1)
	count += len(words)

	// Count symbols (punctuation, markdown characters, etc.)
	symbols := tc.symbolPattern.FindAllString(line, -1)
	count += len(symbols)

	return count
}

// countCodeTokens handles code with more granular tokenization
func (tc *TokenCounter) countCodeTokens(line string) int {
	count := 0

	// Trim the line
	line = strings.TrimSpace(line)
	if line == "" {
		return 0
	}

	// Split on whitespace first
	parts := strings.FieldsFunc(line, unicode.IsSpace)

	for _, part := range parts {
		// Count each character of operators and symbols
		symbols := tc.symbolPattern.FindAllString(part, -1)
		count += len(symbols)

		// Count words (identifiers, keywords)
		words := tc.wordPattern.FindAllString(part, -1)
		count += len(words)

	}

	return count
}
