package token

import (
	"regexp"
	"strings"
)

// TokenCounter provides functionality to estimate token counts for LLM input
type TokenCounter struct {
	// Pre-compiled patterns for token splitting
	// wordPattern     *regexp.Regexp // Not used anymore
	// numberPattern   *regexp.Regexp // Not used anymore
	// symbolPattern   *regexp.Regexp // Not used anymore
	// spacingPattern  *regexp.Regexp // Not used anymore
	// markdownPattern *regexp.Regexp // Not used anymore
}

// NewTokenCounter creates a new TokenCounter instance
func NewTokenCounter() *TokenCounter {
	return &TokenCounter{
		// wordPattern:     regexp.MustCompile(`[a-zA-Z_]\w*`), // Not used anymore
		// numberPattern:   regexp.MustCompile(`\d+`), // Not used anymore
		// symbolPattern:   regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};:'",.<>/?\\|` + "`" + `~]`), // Not used anymore
		// spacingPattern:  regexp.MustCompile(`\s+`), // Not used anymore
		// markdownPattern: regexp.MustCompile(`(\*\*[^*]+\*\*|\*[^*]+\*|\[[^\]]+\]\([^)]+\))`), // Not used anymore
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
			tokenCount += len(tokenizeCode(line))
		} else {
			line = stripMarkdown(line)
			tokenCount += len(tokenizeText(line))
		}
	}

	return tokenCount
}

func stripMarkdown(line string) string {
    // Replace links like [text](url) with "text url"
    line = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`).ReplaceAllString(line, "$1 $2")
    // Replace bold/italic markers with just spaces
    // For example, "**bold**" -> "bold", "*italic*" -> "italic"
    line = regexp.MustCompile(`\*{1,2}([^\*]+)\*{1,2}`).ReplaceAllString(line, "$1")
    return line
}

func tokenizeText(line string) []string {
    // Split by any non-alphanumeric character, preserving apostrophes and underscores in words.
    // For simplicity, let's consider tokens as sequences of letters/numbers or single punctuation.
    return regexp.MustCompile(`[A-Za-z0-9]+|[^\sA-Za-z0-9]`).FindAllString(line, -1)
}

func tokenizeCode(line string) []string {
    tokens := []string{}

    // Handle comments first: Split by "//"
    parts := strings.SplitN(line, "//", 2)
    codePart := parts[0]
    var commentPart string
    if len(parts) > 1 {
        commentPart = parts[1]
    }

    // Extract string literals in code
    stringPattern := regexp.MustCompile(`"[^"]*"`)
    codePart = stringPattern.ReplaceAllStringFunc(codePart, func(match string) string {
        // We'll treat the entire string literal as one token.
        tokens = append(tokens, match)
        return " "
    })

    // Now tokenize the remaining code (identifiers, numbers, symbols)
    codeTokens := regexp.MustCompile(`[A-Za-z0-9]+|[^\sA-Za-z0-9]`).FindAllString(codePart, -1)
    tokens = append(tokens, codeTokens...)

    // Tokenize the comment part if exists
    if commentPart != "" {
        commentTokens := regexp.MustCompile(`[A-Za-z0-9]+|[^\sA-Za-z0-9]`).FindAllString(commentPart, -1)
        tokens = append(tokens, commentTokens...)
    }

    return tokens
}
