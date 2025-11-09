package main

import (
	"fmt"
	"log"

	"github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
	fmt.Println("=== Token Budget and Relevance Filtering Example ===")

	// Example 1: Basic token budget
	fmt.Println("Example 1: Extract with 5000 token budget")
	fmt.Println("---")
	result, err := promptext.Extract("../..",
		promptext.WithExtensions(".go"),
		promptext.WithTokenBudget(5000),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Files included: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("Files excluded: %d\n", result.ExcludedFiles)
	fmt.Printf("Token count: %d\n", result.TokenCount)
	fmt.Printf("Total tokens (if all included): %d\n", result.TotalTokens)

	if result.ProjectOutput.Budget != nil {
		fmt.Printf("Budget max: %d\n", result.ProjectOutput.Budget.MaxTokens)
		fmt.Printf("Budget used: %d\n", result.ProjectOutput.Budget.EstimatedTokens)
	}

	fmt.Println()

	// Example 2: Relevance filtering
	fmt.Println("Example 2: Filter by relevance (processor-related files)")
	fmt.Println("---")
	result, err = promptext.Extract("../..",
		promptext.WithExtensions(".go"),
		promptext.WithRelevance("processor", "extract"),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Relevant files found: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("Files excluded (not relevant): %d\n", result.ExcludedFiles)
	fmt.Println("Included files:")
	for _, file := range result.ProjectOutput.Files {
		fmt.Printf("  - %s (%d tokens)\n", file.Path, file.Tokens)
	}

	fmt.Println()

	// Example 3: Combining relevance and token budget
	fmt.Println("Example 3: Relevance filtering + token budget")
	fmt.Println("---")
	result, err = promptext.Extract("../..",
		promptext.WithExtensions(".go"),
		promptext.WithRelevance("format", "formatter"),
		promptext.WithTokenBudget(8000),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Files included: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("Files excluded: %d\n", result.ExcludedFiles)
	fmt.Printf("Token count: %d / %d\n", result.TokenCount, 8000)

	if len(result.ExcludedFileList) > 0 {
		fmt.Println("\nExcluded files (sample):")
		for i, excluded := range result.ExcludedFileList {
			if i >= 5 {
				fmt.Printf("  ... and %d more\n", len(result.ExcludedFileList)-5)
				break
			}
			fmt.Printf("  - %s (%d tokens)\n", excluded.Path, excluded.Tokens)
		}
	}

	fmt.Println()

	// Example 4: Smart extraction for AI context
	fmt.Println("Example 4: AI-optimized extraction (CLI-related code)")
	fmt.Println("---")
	result, err = promptext.Extract("../..",
		promptext.WithExtensions(".go"),
		promptext.WithRelevance("main", "cli", "command"),
		promptext.WithTokenBudget(10000),
		promptext.WithFormat(promptext.FormatPTX),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Optimized for AI context:\n")
	fmt.Printf("  Files: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("  Tokens: %d (within AI model limits)\n", result.TokenCount)
	fmt.Printf("  Format: PTX (AI-optimized)\n")
	fmt.Printf("  Output size: %d chars\n", len(result.FormattedOutput))

	// Show token efficiency
	if result.TotalTokens > 0 {
		efficiency := float64(result.TokenCount) / float64(result.TotalTokens) * 100
		fmt.Printf("  Efficiency: %.1f%% of total project\n", efficiency)
	}

	fmt.Println()

	// Example 5: Different strategies for different model limits
	fmt.Println("Example 5: Adapt to different AI model token limits")
	fmt.Println("---")

	limits := map[string]int{
		"GPT-3.5 (4k context)":  2000, // Leave room for response
		"GPT-4 (8k context)":    6000,
		"Claude (100k context)": 80000,
		"GPT-4-Turbo (128k)":    100000,
	}

	for model, limit := range limits {
		result, err := promptext.Extract("../..",
			promptext.WithExtensions(".go"),
			promptext.WithTokenBudget(limit),
		)
		if err != nil {
			continue
		}

		fmt.Printf("%s:\n", model)
		fmt.Printf("  Included: %d files, %d tokens\n",
			len(result.ProjectOutput.Files), result.TokenCount)
		fmt.Printf("  Excluded: %d files\n", result.ExcludedFiles)
	}

	fmt.Println("\n✓ All token budget examples completed!")
}
