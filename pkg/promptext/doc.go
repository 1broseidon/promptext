// Package promptext provides a Go library for extracting code context from projects.
//
// Promptext analyzes codebases, filters relevant files, estimates token usage,
// and provides formatted output suitable for AI assistants. It's designed to be
// simple by default while offering powerful configuration options.
//
// # Quick Start
//
// The simplest usage extracts the current directory with sensible defaults:
//
//	result, err := promptext.Extract(".")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.FormattedOutput)
//
// # Filtering Files
//
// Use options to control which files are included:
//
//	result, err := promptext.Extract(".",
//	    promptext.WithExtensions(".go", ".mod", ".sum"),
//	    promptext.WithExcludes("vendor/", "*.test.go"),
//	)
//
// # Relevance and Token Budgets
//
// Filter files by keyword relevance and limit token usage:
//
//	result, err := promptext.Extract(".",
//	    promptext.WithRelevance("auth", "login", "OAuth"),
//	    promptext.WithTokenBudget(8000),
//	)
//
// This is particularly useful for large codebases where you want to focus on
// specific functionality while staying within AI model token limits.
//
// # Output Formats
//
// Choose from multiple output formats:
//
//	// PTX v2.0 (recommended for AI assistants)
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatPTX))
//
//	// JSONL (machine-friendly, one JSON object per line)
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatJSONL))
//
//	// Markdown (human-readable)
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatMarkdown))
//
//	// XML (machine-parseable)
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatXML))
//
// # Format Conversion
//
// Convert results to different formats without re-processing:
//
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatPTX))
//	markdownOutput, _ := result.As(promptext.FormatMarkdown)
//	jsonlOutput, _ := result.As(promptext.FormatJSONL)
//
// # Reusable Extractors
//
// Create an extractor to process multiple directories with the same configuration:
//
//	extractor := promptext.NewExtractor(
//	    promptext.WithExtensions(".go"),
//	    promptext.WithTokenBudget(8000),
//	)
//	result1, _ := extractor.Extract("/path/to/project1")
//	result2, _ := extractor.Extract("/path/to/project2")
//
// # Custom Formatters
//
// Register custom formatters for specialized output needs:
//
//	type MyFormatter struct{}
//
//	func (f *MyFormatter) Format(output *promptext.ProjectOutput) (string, error) {
//	    // Custom formatting logic
//	    return "custom format", nil
//	}
//
//	promptext.RegisterFormatter("myformat", &MyFormatter{})
//	result, _ := promptext.Extract(".", promptext.WithFormat("myformat"))
//
// # Error Handling
//
// The library provides typed errors for common cases:
//
//	result, err := promptext.Extract("/nonexistent")
//	if errors.Is(err, promptext.ErrInvalidDirectory) {
//	    // Handle invalid directory
//	}
//
//	result, err := promptext.Extract(".",
//	    promptext.WithExtensions(".xyz"),
//	)
//	if errors.Is(err, promptext.ErrNoFilesMatched) {
//	    // Handle no matching files
//	}
//
// # Accessing Structured Data
//
// The Result contains both formatted output and structured data:
//
//	result, _ := promptext.Extract(".")
//
//	// Formatted output ready for AI
//	fmt.Println(result.FormattedOutput)
//
//	// Structured data for programmatic access
//	fmt.Printf("Files: %d\n", len(result.ProjectOutput.Files))
//	fmt.Printf("Token count: %d\n", result.TokenCount)
//	fmt.Printf("Language: %s\n", result.ProjectOutput.Metadata.Language)
//
//	// Iterate over files
//	for _, file := range result.ProjectOutput.Files {
//	    fmt.Printf("%s: %d tokens\n", file.Path, file.Tokens)
//	}
//
// # Configuration Options
//
// Available options:
//
//   - WithExtensions(extensions ...string) - Include specific file extensions
//   - WithExcludes(patterns ...string) - Exclude files matching patterns
//   - WithGitIgnore(enabled bool) - Respect .gitignore patterns (default: true)
//   - WithDefaultRules(enabled bool) - Use built-in filtering rules (default: true)
//   - WithRelevance(keywords ...string) - Filter by keyword relevance
//   - WithTokenBudget(maxTokens int) - Limit output to token budget
//   - WithFormat(format Format) - Set output format
//   - WithVerbose(enabled bool) - Enable verbose logging
//   - WithDebug(enabled bool) - Enable debug logging with timing
//
// # Design Principles
//
// 1. Simple by Default: Works with zero configuration
// 2. Composable: Options can be combined naturally
// 3. Discoverable: IDE autocomplete reveals all options
// 4. Safe: Errors are typed and checkable with errors.Is()
// 5. Extensible: Custom formatters can be registered
//
// # Examples
//
// See the examples/ directory in the repository for complete working examples:
//
//   - examples/basic/ - Simple extraction
//   - examples/filtered/ - With extensions and excludes
//   - examples/token-budget/ - AI-focused extraction with token limits
//   - examples/custom-formatter/ - Custom output format
//
// # Version
//
// This is version 0.7.0 (Phase 1) of the library API.
// The API may evolve during the 0.x releases. Version 1.0.0 will provide
// API stability guarantees and backward compatibility.
//
// For more information, visit https://github.com/1broseidon/promptext
package promptext
