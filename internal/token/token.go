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
		symbolPattern:  regexp.MustCompile(`[^\w\s]+`), // Changed to capture sequences
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

	// Split on whitespace first
	parts := strings.Fields(line)
	
	for _, part := range parts {
		// Handle markdown syntax
		if strings.HasPrefix(part, "[") && strings.Contains(part, "](") {
			// Count link text and URL separately
			count += 2
			continue
		}
		if strings.HasPrefix(part, "*") || strings.HasPrefix(part, "**") {
			// Count emphasized text as one token
			count++
			continue
		}
		
		// Count remaining words and symbol sequences
		words := tc.wordPattern.FindAllString(part, -1)
		symbols := tc.symbolPattern.FindAllString(part, -1)
		count += len(words) + len(symbols)
	}

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

	// Handle comments separately
	if idx := strings.Index(line, "//"); idx >= 0 {
		beforeComment := strings.TrimSpace(line[:idx])
		// Count the code before comment plus 1 for the comment itself
		return tc.countCodeTokens(beforeComment) + 1
	}

	// Split on whitespace first
	parts := strings.FieldsFunc(line, unicode.IsSpace)

	for _, part := range parts {
		if part == "" {
			continue
		}

		// Handle string literals
		if strings.HasPrefix(part, "\"") && strings.HasSuffix(part, "\"") {
			count++ // Count entire string as one token
			continue
		}

		// Handle operators and symbols as sequences
		symbols := tc.symbolPattern.FindAllString(part, -1)
		words := tc.wordPattern.FindAllString(part, -1)
		
		count += len(symbols) + len(words)
	}

	return count
}
