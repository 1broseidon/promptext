# Promptext Library Examples

This directory contains example programs demonstrating how to use the promptext library in your Go applications.

## Running the Examples

Each example can be run independently:

```bash
# Basic usage examples
cd basic
go run main.go

# Token budget and relevance filtering
cd token-budget
go run main.go
```

## Examples

### 1. Basic (`basic/`)

Demonstrates fundamental library usage:
- Simple extraction with defaults
- Filtering by file extensions
- Excluding patterns
- Different output formats
- Format conversion
- Saving to files
- Reusable extractors
- Builder pattern

**Run it:**
```bash
cd basic && go run main.go
```

### 2. Token Budget (`token-budget/`)

Shows how to work with token budgets and relevance filtering:
- Setting token budgets for AI model limits
- Relevance-based file filtering
- Combining relevance and token budgets
- Optimizing for different AI models
- Understanding token efficiency

**Run it:**
```bash
cd token-budget && go run main.go
```

## Common Patterns

### Simple Extraction

```go
result, err := promptext.Extract(".")
if err != nil {
    log.Fatal(err)
}
fmt.Println(result.FormattedOutput)
```

### With Options

```go
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go", ".mod"),
    promptext.WithExcludes("vendor/", "*_test.go"),
    promptext.WithFormat(promptext.FormatPTX),
)
```

### Relevance Filtering

```go
result, err := promptext.Extract(".",
    promptext.WithRelevance("auth", "login"),
    promptext.WithTokenBudget(8000),
)
```

### Reusable Extractor

```go
extractor := promptext.NewExtractor(
    promptext.WithExtensions(".go"),
    promptext.WithTokenBudget(5000),
)

result1, _ := extractor.Extract("/project1")
result2, _ := extractor.Extract("/project2")
```

### Format Conversion

```go
result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatPTX))

// Convert to different formats
markdown, _ := result.As(promptext.FormatMarkdown)
jsonl, _ := result.As(promptext.FormatJSONL)
```

## Use Cases

### 1. AI Context Generation

Generate optimized code context for AI assistants:

```go
result, err := promptext.Extract(".",
    promptext.WithRelevance("authentication"),
    promptext.WithTokenBudget(8000),
    promptext.WithFormat(promptext.FormatPTX),
)
// Send result.FormattedOutput to AI
```

### 2. Code Documentation

Extract code for documentation purposes:

```go
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go"),
    promptext.WithExcludes("*_test.go", "vendor/"),
    promptext.WithFormat(promptext.FormatMarkdown),
)
os.WriteFile("docs/codebase.md", []byte(result.FormattedOutput), 0644)
```

### 3. CI/CD Integration

Analyze code in CI pipelines:

```go
result, err := promptext.Extract(".",
    promptext.WithFormat(promptext.FormatJSONL),
)
// Parse result.FormattedOutput as JSONL for analysis
```

### 4. Code Search and Analysis

Find relevant code across large codebases:

```go
result, err := promptext.Extract(".",
    promptext.WithRelevance("database", "query", "migration"),
)

for _, file := range result.ProjectOutput.Files {
    fmt.Printf("%s: %d tokens\n", file.Path, file.Tokens)
}
```

## Available Options

- `WithExtensions(extensions ...string)` - Filter by file extensions
- `WithExcludes(patterns ...string)` - Exclude file patterns
- `WithGitIgnore(enabled bool)` - Respect .gitignore (default: true)
- `WithDefaultRules(enabled bool)` - Use built-in rules (default: true)
- `WithRelevance(keywords ...string)` - Filter by keyword relevance
- `WithTokenBudget(maxTokens int)` - Limit output tokens
- `WithFormat(format Format)` - Set output format
- `WithVerbose(enabled bool)` - Enable verbose logging
- `WithDebug(enabled bool)` - Enable debug logging

## Output Formats

- `FormatPTX` - PTX v2.0 (recommended for AI, TOON-based)
- `FormatTOON` - Alias for PTX (backward compatibility)
- `FormatJSONL` - Machine-friendly JSONL
- `FormatTOONStrict` - TOON v1.3 strict compliance
- `FormatMarkdown` - Human-readable markdown
- `FormatXML` - Machine-parseable XML

## Error Handling

```go
result, err := promptext.Extract("/invalid/path")
if err != nil {
    if errors.Is(err, promptext.ErrInvalidDirectory) {
        // Handle invalid directory
    }
    if errors.Is(err, promptext.ErrNoFilesMatched) {
        // Handle no matching files
    }
    // Handle other errors
}
```

## Further Reading

- [Library Documentation](../pkg/promptext/doc.go)
- [Main README](../README.md)
- [API Reference](https://pkg.go.dev/github.com/1broseidon/promptext/pkg/promptext)
