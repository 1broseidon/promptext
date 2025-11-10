// Package main implements a semantic code search tool using promptext.
//
// This example demonstrates how to build a natural language code search tool
// that finds relevant code across large codebases using keyword-based relevance
// scoring and AI context extraction.
//
// Usage:
//   go run main.go "Where is user authentication handled?"
//   go run main.go "How does database connection pooling work?"
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go \"your search query\"")
		fmt.Println("\nExamples:")
		fmt.Println("  go run main.go \"Where is user authentication handled?\"")
		fmt.Println("  go run main.go \"How does database connection pooling work?\"")
		fmt.Println("  go run main.go \"Find all API endpoint definitions\"")
		os.Exit(1)
	}

	query := strings.Join(os.Args[1:], " ")
	fmt.Printf("🔍 Searching for: %s\n\n", query)

	// Extract keywords from the natural language query
	// In a production tool, you might use NLP or call an AI API to extract keywords
	keywords := extractKeywords(query)
	fmt.Printf("📋 Keywords extracted: %v\n\n", keywords)

	// Get the current directory (or accept as argument)
	searchDir := "."
	if dir := os.Getenv("SEARCH_DIR"); dir != "" {
		searchDir = dir
	}

	absPath, err := filepath.Abs(searchDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📂 Searching in: %s\n", filepath.Base(absPath))

	// Use promptext to find relevant files
	// The relevance scoring will prioritize files that match our keywords
	result, err := promptext.Extract(searchDir,
		// Filter by programming language extensions (customize as needed)
		promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java"),

		// Use relevance filtering to find files matching our keywords
		promptext.WithRelevance(keywords...),

		// Limit to a reasonable token budget (enough for context but not overwhelming)
		promptext.WithTokenBudget(5000),

		// Use PTX format for AI-friendly output
		promptext.WithFormat(promptext.FormatPTX),
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error during search: %v\n", err)
		os.Exit(1)
	}

	// Display search results
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Printf("✨ Found %d relevant files (%d tokens)\n",
		len(result.ProjectOutput.Files), result.TokenCount)

	if result.ExcludedFiles > 0 {
		fmt.Printf("   ℹ️  %d additional files excluded due to token budget\n", result.ExcludedFiles)
	}
	fmt.Println(strings.Repeat("=", 60))

	// List the relevant files found
	fmt.Println("\n📄 Relevant Files:")
	for i, file := range result.ProjectOutput.Files {
		fmt.Printf("   %d. %s (%d tokens)\n", i+1, file.Path, file.Tokens)
	}

	// In a real application, you would:
	// 1. Send result.FormattedOutput to an AI API (Claude, GPT, etc.)
	// 2. Ask the AI to answer the query based on the extracted code
	// 3. Display the AI's response to the user

	fmt.Println("\n💡 Next Steps:")
	fmt.Println("   The extracted code context is ready to send to an AI assistant.")
	fmt.Println("   You can paste the output below into ChatGPT/Claude to get answers:")
	fmt.Println()
	fmt.Println("   Example prompt:")
	fmt.Printf("   \"Based on this code: %s\"\n", query)
	fmt.Println()

	// Optionally save to file for manual review
	outputFile := "search-results.ptx"
	if err := os.WriteFile(outputFile, []byte(result.FormattedOutput), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Could not save results to file: %v\n", err)
	} else {
		fmt.Printf("💾 Full context saved to: %s\n", outputFile)
		fmt.Printf("   (%d characters, %d tokens)\n", len(result.FormattedOutput), result.TokenCount)
	}

	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🎯 Search complete!")
}

// extractKeywords performs simple keyword extraction from a natural language query.
// In a production system, you might use:
// - NLP libraries (like spacy, nltk)
// - AI APIs to extract semantically relevant terms
// - Custom domain-specific keyword dictionaries
func extractKeywords(query string) []string {
	// Convert to lowercase for matching
	lower := strings.ToLower(query)

	// Remove common question words and filler words
	stopWords := map[string]bool{
		"where": true, "is": true, "are": true, "how": true, "does": true,
		"do": true, "the": true, "a": true, "an": true, "in": true,
		"on": true, "at": true, "to": true, "for": true, "of": true,
		"with": true, "by": true, "from": true, "what": true, "which": true,
		"who": true, "when": true, "why": true, "can": true, "could": true,
		"would": true, "should": true, "all": true, "any": true, "some": true,
		"this": true, "that": true, "these": true, "those": true, "i": true,
		"you": true, "we": true, "they": true, "it": true, "be": true,
		"been": true, "being": true, "have": true, "has": true, "had": true,
		"find": true, "show": true, "get": true, "tell": true, "me": true,
	}

	// Split into words and filter
	words := strings.FieldsFunc(lower, func(r rune) bool {
		return !((r >= 'a' && r <= 'z') || (r >= '0' && r <= '9'))
	})

	keywords := make([]string, 0)
	seen := make(map[string]bool)

	for _, word := range words {
		// Skip stop words and very short words
		if len(word) < 3 || stopWords[word] {
			continue
		}

		// Avoid duplicates
		if !seen[word] {
			keywords = append(keywords, word)
			seen[word] = true
		}
	}

	// If we didn't find any keywords, use the original query
	if len(keywords) == 0 {
		keywords = words
	}

	return keywords
}
