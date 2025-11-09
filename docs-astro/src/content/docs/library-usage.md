---
title: Library Usage
description: Complete guide to using promptext as a Go library in your applications
---

## Overview

As of v0.7.0, promptext can be used as a Go library in your applications. This allows you to programmatically extract code context, analyze projects, and integrate AI-optimized context generation into your own tools.

## Installation

Add promptext to your Go project:

```bash
go get github.com/1broseidon/promptext/pkg/promptext@latest
```

## Quick Start

The simplest way to use promptext is with the `Extract()` function:

```go
package main

import (
    "fmt"
    "log"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    // Extract code context from current directory
    result, err := promptext.Extract(".")
    if err != nil {
        log.Fatal(err)
    }

    // Output is ready to paste into AI assistant
    fmt.Println(result.FormattedOutput)

    // Access metadata
    fmt.Printf("Processed %d files (~%d tokens)\n",
        len(result.ProjectOutput.Files),
        result.TokenCount)
}
```

## Configuration Options

Use functional options to customize extraction behavior:

### File Filtering

```go
// Filter by file extensions
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go", ".mod", ".sum"),
)

// Exclude patterns
result, err := promptext.Extract(".",
    promptext.WithExcludes("vendor/", "*.test.go", "testdata/"),
)

// Combine filters
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go"),
    promptext.WithExcludes("vendor/", "internal/test/"),
)
```

### Relevance Filtering

Prioritize files by keywords (useful for focused context):

```go
// Find authentication-related code
result, err := promptext.Extract(".",
    promptext.WithRelevance("auth", "login", "session", "OAuth"),
)

// Focus on database operations
result, err := promptext.Extract(".",
    promptext.WithRelevance("database", "sql", "query", "migration"),
)
```

Files are scored based on:
- **Filename matches** (10x weight)
- **Directory path matches** (5x weight)
- **Import statements** (3x weight)
- **Content matches** (1x weight)

### Token Budget Management

Limit output to fit within AI model token limits:

```go
// Limit to 8K tokens for smaller context windows
result, err := promptext.Extract(".",
    promptext.WithTokenBudget(8000),
)

// Combine with relevance for smart prioritization
result, err := promptext.Extract(".",
    promptext.WithRelevance("api", "handler"),
    promptext.WithTokenBudget(10000),
)

// Check what was excluded
if result.ExcludedFiles > 0 {
    fmt.Printf("Excluded %d files to fit budget\n", result.ExcludedFiles)
    for _, excluded := range result.ExcludedFileList {
        fmt.Printf("  - %s (~%d tokens)\n", excluded.Path, excluded.Tokens)
    }
}
```

### Output Format Selection

Choose the format that best suits your needs:

```go
// PTX format (default, token-optimized)
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatPTX),
)

// Markdown (readable, good for documentation)
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatMarkdown),
)

// JSONL (streaming, programmatic processing)
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatJSONL),
)

// XML (structured, enterprise systems)
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatXML),
)
```

### GitIgnore and Default Rules

Control filtering behavior:

```go
// Disable .gitignore patterns (include all files)
result, err := promptext.Extract(".",
    promptext.WithGitIgnore(false),
)

// Disable built-in filtering rules
result, err := promptext.Extract(".",
    promptext.WithDefaultRules(false),
)
```

### Debug and Verbose Logging

Enable logging for troubleshooting:

```go
// Verbose output
result, err := promptext.Extract(".",
    promptext.WithVerbose(true),
)

// Debug mode with detailed timing
result, err := promptext.Extract(".",
    promptext.WithDebug(true),
)
```

## Reusable Extractor

For processing multiple directories with the same configuration:

```go
// Create extractor with configuration
extractor := promptext.NewExtractor(
    promptext.WithExtensions(".go", ".mod"),
    promptext.WithTokenBudget(8000),
    promptext.WithFormat(promptext.FormatPTX),
)

// Reuse for multiple projects
result1, err := extractor.Extract("/path/to/project1")
result2, err := extractor.Extract("/path/to/project2")
result3, err := extractor.Extract("/path/to/project3")
```

You can also chain configuration methods:

```go
extractor := promptext.NewExtractor().
    WithExtensions(".go").
    WithExcludes("vendor/").
    WithFormat(promptext.FormatMarkdown)

result, err := extractor.Extract(".")
```

## Format Conversion

Convert results to different formats without re-processing:

```go
// Extract once
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatPTX),
)

// Convert to other formats as needed
markdownOutput, err := result.As(promptext.FormatMarkdown)
jsonlOutput, err := result.As(promptext.FormatJSONL)
xmlOutput, err := result.As(promptext.FormatXML)
```

## Accessing Structured Data

The `Result` type provides complete access to extracted data:

```go
result, err := promptext.Extract(".")

// Project metadata
metadata := result.ProjectOutput.Metadata
fmt.Printf("Language: %s\n", metadata.Language)
fmt.Printf("Version: %s\n", metadata.Version)
fmt.Printf("Dependencies: %v\n", metadata.Dependencies)

// Git information
if git := result.ProjectOutput.GitInfo; git != nil {
    fmt.Printf("Branch: %s\n", git.Branch)
    fmt.Printf("Commit: %s\n", git.CommitHash)
}

// File statistics
stats := result.ProjectOutput.FileStats
fmt.Printf("Total Files: %d\n", stats.TotalFiles)
fmt.Printf("Total Lines: %d\n", stats.TotalLines)

// Token budget information
fmt.Printf("Token Count: %d\n", result.TokenCount)
fmt.Printf("Total Tokens: %d\n", result.TotalTokens)
fmt.Printf("Excluded Files: %d\n", result.ExcludedFiles)

// Individual file access
for _, file := range result.ProjectOutput.Files {
    fmt.Printf("File: %s (%d tokens)\n", file.Path, file.Tokens)
    if file.Truncation != nil {
        fmt.Printf("  Truncated: %s (original: %d tokens)\n",
            file.Truncation.Mode,
            file.Truncation.OriginalTokens)
    }
}

// Directory tree
printTree(result.ProjectOutput.DirectoryTree, 0)
```

## Error Handling

The library provides well-typed errors:

```go
result, err := promptext.Extract("/invalid/path")
if err != nil {
    switch {
    case errors.Is(err, promptext.ErrInvalidDirectory):
        log.Println("Directory does not exist or is not accessible")
    case errors.Is(err, promptext.ErrNoFilesMatched):
        log.Println("No files matched the filter criteria")
    case errors.Is(err, promptext.ErrTokenBudgetTooLow):
        log.Println("Token budget too low to include any files")
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

Check for specific error types:

```go
var dirErr *promptext.DirectoryError
if errors.As(err, &dirErr) {
    log.Printf("Directory error at %s: %v", dirErr.Path, dirErr.Err)
}

var formatErr *promptext.FormatError
if errors.As(err, &formatErr) {
    log.Printf("Format error for %s: %v", formatErr.Format, formatErr.Err)
}
```

## Common Use Cases

### AI Code Review Tool

```go
func analyzeForReview(projectPath string) error {
    // Extract with focus on recent changes
    extractor := promptext.NewExtractor(
        promptext.WithExtensions(".go", ".js", ".ts", ".py"),
        promptext.WithExcludes("vendor/", "node_modules/", "test/"),
        promptext.WithTokenBudget(15000), // Fit in context window
        promptext.WithFormat(promptext.FormatPTX),
    )

    result, err := extractor.Extract(projectPath)
    if err != nil {
        return fmt.Errorf("extraction failed: %w", err)
    }

    // Send to AI for review
    review, err := sendToAI(result.FormattedOutput)
    if err != nil {
        return err
    }

    fmt.Println(review)
    return nil
}
```

### Codebase Documentation Generator

```go
func generateDocs(projectPath string) error {
    // Extract full codebase context
    result, err := promptext.Extract(projectPath,
        promptext.WithExtensions(".go", ".md"),
        promptext.WithFormat(promptext.FormatMarkdown),
    )
    if err != nil {
        return err
    }

    // Save to documentation file
    return os.WriteFile("docs/codebase-overview.md",
        []byte(result.FormattedOutput), 0644)
}
```

### Smart Context Search

```go
func findRelevantCode(projectPath string, query string) (*promptext.Result, error) {
    // Split query into keywords
    keywords := strings.Fields(query)

    // Extract with relevance filtering
    return promptext.Extract(projectPath,
        promptext.WithRelevance(keywords...),
        promptext.WithTokenBudget(8000),
    )
}

// Usage
result, err := findRelevantCode("/my/project", "authentication jwt token")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d relevant files:\n", len(result.ProjectOutput.Files))
for _, file := range result.ProjectOutput.Files {
    fmt.Printf("  - %s\n", file.Path)
}
```

### Multi-Project Analyzer

```go
func analyzeProjects(projectPaths []string) error {
    // Create shared configuration
    extractor := promptext.NewExtractor(
        promptext.WithExtensions(".go"),
        promptext.WithExcludes("vendor/", "test/"),
        promptext.WithFormat(promptext.FormatJSONL),
    )

    for _, path := range projectPaths {
        result, err := extractor.Extract(path)
        if err != nil {
            log.Printf("Failed to analyze %s: %v", path, err)
            continue
        }

        // Process results
        fmt.Printf("Project: %s\n", path)
        fmt.Printf("  Files: %d\n", len(result.ProjectOutput.Files))
        fmt.Printf("  Tokens: %d\n", result.TokenCount)
        fmt.Printf("  Language: %s\n", result.ProjectOutput.Metadata.Language)
        fmt.Println()
    }

    return nil
}
```

### Token-Aware Context Builder

```go
func buildContextWithBudget(projectPath string, maxTokens int) (string, error) {
    // Start with high-priority files
    result, err := promptext.Extract(projectPath,
        promptext.WithRelevance("main", "api", "handler"),
        promptext.WithTokenBudget(maxTokens),
    )
    if err != nil {
        return "", err
    }

    // Report what was included
    included := len(result.ProjectOutput.Files)
    total := included + result.ExcludedFiles

    fmt.Printf("Context built: %d/%d files (~%d tokens)\n",
        included, total, result.TokenCount)

    if result.ExcludedFiles > 0 {
        fmt.Printf("Excluded %d files to stay within %d token budget\n",
            result.ExcludedFiles, maxTokens)
    }

    return result.FormattedOutput, nil
}
```

## Custom Formatters

Register custom output formatters:

```go
// Implement the Formatter interface
type MyCustomFormatter struct{}

func (f *MyCustomFormatter) Format(output *promptext.ProjectOutput) (string, error) {
    var buf strings.Builder

    // Custom formatting logic
    buf.WriteString("=== Custom Format ===\n")
    for _, file := range output.Files {
        buf.WriteString(fmt.Sprintf("File: %s\n%s\n\n", file.Path, file.Content))
    }

    return buf.String(), nil
}

// Register the formatter
promptext.RegisterFormatter("custom", &MyCustomFormatter{})

// Use it
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.Format("custom")),
)
```

## Best Practices

### 1. Use Relevance Filtering for Large Codebases

```go
// Instead of extracting everything
result, err := promptext.Extract("/large/project")

// Focus on relevant code
result, err := promptext.Extract("/large/project",
    promptext.WithRelevance("feature", "bug", "fix"),
    promptext.WithTokenBudget(10000),
)
```

### 2. Reuse Extractors for Consistent Processing

```go
// Create once
extractor := promptext.NewExtractor(
    promptext.WithExtensions(".go"),
    promptext.WithExcludes("vendor/", "*_test.go"),
)

// Use multiple times
for _, project := range projects {
    result, err := extractor.Extract(project)
    // ... process result
}
```

### 3. Handle Errors Appropriately

```go
result, err := promptext.Extract(path, opts...)
if err != nil {
    // Check for specific errors
    if errors.Is(err, promptext.ErrNoFilesMatched) {
        return nil, fmt.Errorf("no files found with current filters")
    }
    return nil, fmt.Errorf("extraction failed: %w", err)
}
```

### 4. Monitor Token Budgets

```go
result, err := promptext.Extract(".", promptext.WithTokenBudget(8000))
if err != nil {
    return err
}

percentage := float64(result.TokenCount) / float64(result.TotalTokens) * 100
fmt.Printf("Using %.1f%% of codebase (%d/%d tokens)\n",
    percentage, result.TokenCount, result.TotalTokens)
```

## API Reference

For complete API documentation, see [pkg.go.dev](https://pkg.go.dev/github.com/1broseidon/promptext/pkg/promptext).

### Core Functions

- `Extract(dir string, opts ...Option) (*Result, error)` - Extract code context from directory
- `NewExtractor(opts ...Option) *Extractor` - Create reusable extractor
- `RegisterFormatter(name string, formatter Formatter)` - Register custom formatter
- `GetFormatter(name string) (Formatter, error)` - Get registered formatter

### Options

- `WithExtensions(...string)` - Filter by file extensions
- `WithExcludes(...string)` - Exclude file patterns
- `WithRelevance(...string)` - Keyword-based relevance filtering
- `WithTokenBudget(int)` - Set token budget limit
- `WithFormat(Format)` - Set output format
- `WithGitIgnore(bool)` - Control .gitignore respect
- `WithDefaultRules(bool)` - Control built-in filtering rules
- `WithVerbose(bool)` - Enable verbose logging
- `WithDebug(bool)` - Enable debug logging

### Result Types

- `Result` - Extraction result with formatted output and metadata
- `ProjectOutput` - Structured project data
- `FileInfo` - Individual file information
- `ExcludedFileInfo` - Information about excluded files
- `DirectoryNode` - Directory tree structure
- `GitInfo` - Git repository information
- `Metadata` - Project metadata

### Error Types

- `ErrInvalidDirectory` - Invalid or inaccessible directory
- `ErrNoFilesMatched` - No files matched filter criteria
- `ErrTokenBudgetTooLow` - Token budget too low
- `ErrInvalidFormat` - Unsupported output format
- `DirectoryError` - Directory access error with path
- `FormatError` - Format conversion error

## Examples

The library includes working examples in the `examples/` directory:

- [`examples/basic/`](https://github.com/1broseidon/promptext/tree/main/examples/basic) - Fundamental usage patterns
- [`examples/token-budget/`](https://github.com/1broseidon/promptext/tree/main/examples/token-budget) - Token budget management

Run examples:

```bash
cd examples/basic && go run main.go
cd examples/token-budget && go run main.go
```

## Migration from CLI to Library

If you're currently using the CLI and want to integrate into your Go application:

**CLI:**
```bash
prx -e .go -r "api handler" --max-tokens 8000 -o context.ptx
```

**Library equivalent:**
```go
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go"),
    promptext.WithRelevance("api", "handler"),
    promptext.WithTokenBudget(8000),
)
if err != nil {
    log.Fatal(err)
}

os.WriteFile("context.ptx", []byte(result.FormattedOutput), 0644)
```

## Support

- [GitHub Issues](https://github.com/1broseidon/promptext/issues) - Bug reports and feature requests
- [API Documentation](https://pkg.go.dev/github.com/1broseidon/promptext/pkg/promptext) - Complete API reference
- [Examples](https://github.com/1broseidon/promptext/tree/main/examples) - Working code examples
