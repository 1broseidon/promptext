package main

import (
	"fmt"
	"log"
	"os"

	"github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
	// Example 1: Simple extraction with defaults
	fmt.Println("=== Example 1: Simple Extraction ===")
	result, err := promptext.Extract(".")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Files processed: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("Token count: %d\n", result.TokenCount)
	if result.ProjectOutput.Metadata != nil {
		fmt.Printf("Language: %s\n", result.ProjectOutput.Metadata.Language)
	}

	// Example 2: Extract with specific extensions
	fmt.Println("\n=== Example 2: Extract Go Files Only ===")
	result, err = promptext.Extract(".",
		promptext.WithExtensions(".go", ".mod"),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Go files processed: %d\n", len(result.ProjectOutput.Files))
	for _, file := range result.ProjectOutput.Files {
		fmt.Printf("  - %s (%d tokens)\n", file.Path, file.Tokens)
	}

	// Example 3: Extract with exclusions
	fmt.Println("\n=== Example 3: Extract with Exclusions ===")
	result, err = promptext.Extract(".",
		promptext.WithExtensions(".go"),
		promptext.WithExcludes("*_test.go", "vendor/"),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Files (excluding tests): %d\n", len(result.ProjectOutput.Files))

	// Example 4: Extract with token budget
	fmt.Println("\n=== Example 4: Extract with Token Budget ===")
	result, err = promptext.Extract(".",
		promptext.WithTokenBudget(5000),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Printf("Files included: %d\n", len(result.ProjectOutput.Files))
	fmt.Printf("Files excluded: %d\n", result.ExcludedFiles)
	fmt.Printf("Token count: %d\n", result.TokenCount)

	// Example 5: Different output formats
	fmt.Println("\n=== Example 5: Different Output Formats ===")

	// PTX format (default)
	result, err = promptext.Extract(".",
		promptext.WithFormat(promptext.FormatPTX),
		promptext.WithExtensions(".go"),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("PTX format output length: %d chars\n", len(result.FormattedOutput))

	// Convert to other formats
	markdownOutput, err := result.As(promptext.FormatMarkdown)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Markdown format output length: %d chars\n", len(markdownOutput))

	jsonlOutput, err := result.As(promptext.FormatJSONL)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("JSONL format output length: %d chars\n", len(jsonlOutput))

	// Example 6: Save to file
	fmt.Println("\n=== Example 6: Save to File ===")
	result, err = promptext.Extract(".",
		promptext.WithExtensions(".go", ".mod"),
		promptext.WithFormat(promptext.FormatPTX),
	)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	outputFile := "output.ptx"
	err = os.WriteFile(outputFile, []byte(result.FormattedOutput), 0644)
	if err != nil {
		log.Fatalf("Error writing file: %v", err)
	}
	fmt.Printf("Output saved to %s\n", outputFile)

	// Example 7: Reusable extractor
	fmt.Println("\n=== Example 7: Reusable Extractor ===")
	extractor := promptext.NewExtractor(
		promptext.WithExtensions(".go"),
		promptext.WithFormat(promptext.FormatJSONL),
	)

	// Use the same extractor for different directories
	result1, err := extractor.Extract(".")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Current directory: %d files\n", len(result1.ProjectOutput.Files))

	// Example 8: Builder pattern
	fmt.Println("\n=== Example 8: Builder Pattern ===")
	result, err = promptext.NewExtractor().
		WithExtensions(".go").
		WithExcludes("*_test.go").
		WithFormat(promptext.FormatMarkdown).
		Extract(".")

	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fmt.Printf("Builder pattern result: %d files\n", len(result.ProjectOutput.Files))

	fmt.Println("\n✓ All examples completed successfully!")
}
