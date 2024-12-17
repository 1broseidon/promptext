package token

import (
	"regexp"
	"strings"
)

// TokenCounter provides functionality to estimate token counts for LLM input
type TokenCounter struct {
	// Pre-compiled patterns for token splitting
	wordPattern     *regexp.Regexp
	numberPattern   *regexp.Regexp
	symbolPattern   *regexp.Regexp
	spacingPattern  *regexp.Regexp
	markdownPattern *regexp.Regexp
}

// NewTokenCounter creates a new TokenCounter instance
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{
		wordPattern:     regexp.MustCompile(`[a-zA-Z_]\w*`),
		numberPattern:   regexp.MustCompile(`\d+`),
		symbolPattern:   regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'",.<>/?\\|` + "`" + `~]`),
		spacingPattern:  regexp.MustCompile(`\s+`),
		markdownPattern: regexp.MustCompile(`(\*\*[^*]+\*\*|\*[^*]+\*|\[[^\]]+\]\([^)]+\))`),
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

	// Handle markdown formatting first
	line = tc.markdownPattern.ReplaceAllStringFunc(line, func(match string) string {
		if strings.HasPrefix(match, "[") {
			// For links, count both text and URL
			count += 2
		} else {
			// For bold/italic, count as one plus the words inside
			inner := strings.Trim(match, "*")
			words := tc.wordPattern.FindAllString(inner, -1)
			count += len(words)
		}
		return ""
	})

	// Split remaining text into tokens
	tokens := strings.Fields(line)
	for _, token := range tokens {
		// Count words
		words := tc.wordPattern.FindAllString(token, -1)
		count += len(words)

		// Count numbers
		numbers := tc.numberPattern.FindAllString(token, -1)
		count += len(numbers)

		// Count symbols individually
		symbols := tc.symbolPattern.FindAllString(token, -1)
		count += len(symbols)
	}

	return count
}

// countCodeTokens handles code with more granular tokenization
func (tc *TokenCounter) countCodeTokens(line string) int {
	count := 0

	// Handle comments
	if idx := strings.Index(line, "//"); idx >= 0 {
		beforeComment := line[:idx]
		afterComment := line[idx:]
		count += tc.countCodePart(beforeComment)
		count += 2 // Count comment marker and content separately
		return count
	}

	return tc.countCodePart(line)
}

func (tc *TokenCounter) countCodePart(code string) int {
	count := 0

	// Handle string literals
	stringPattern := regexp.MustCompile(`"[^"]*"`)
	code = stringPattern.ReplaceAllStringFunc(code, func(match string) string {
		count += 2 // Count quotes and content separately
		return ""
	})

	// Split remaining code into tokens
	tokens := strings.Fields(code)
	for _, token := range tokens {
		// Count identifiers
		words := tc.wordPattern.FindAllString(token, -1)
		count += len(words)

		// Count numbers separately
		numbers := tc.numberPattern.FindAllString(token, -1)
		count += len(numbers)

		// Count symbols individually
		symbols := tc.symbolPattern.FindAllString(token, -1)
		count += len(symbols)
	}

	// Add extra token for newlines in code blocks
	if strings.Contains(code, "\n") {
		count++
	}

	return count
}
