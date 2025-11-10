// Package main implements an automated documentation generator using promptext.
//
// This tool extracts code from your project and uses AI to generate or update
// documentation automatically. It helps keep docs in sync with code changes.
//
// Features:
// - Extract public API surfaces
// - Generate API reference docs
// - Update README sections
// - Create usage examples
// - Maintain changelogs
//
// Usage:
//   # Generate API reference
//   go run main.go --type api --output docs/api-reference.md
//
//   # Update README
//   go run main.go --type readme --output README.md
//
//   # Generate all docs
//   go run main.go --all
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/pkg/promptext"
)

// DocType represents the type of documentation to generate
type DocType string

const (
	DocTypeAPI     DocType = "api"
	DocTypeREADME  DocType = "readme"
	DocTypeGuide   DocType = "guide"
	DocTypeExample DocType = "example"
)

// Config holds the documentation generation configuration
type Config struct {
	DocType    DocType
	OutputFile string
	SourceDir  string
	GenerateAll bool
}

func main() {
	config := parseFlags()

	fmt.Println("📚 Documentation Generator")
	fmt.Println(strings.Repeat("=", 60))

	if config.GenerateAll {
		fmt.Println("Generating all documentation...")
		generateAllDocs(config)
	} else {
		fmt.Printf("Generating %s documentation...\n\n", config.DocType)
		if err := generateDoc(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n🎯 Documentation generation complete!")
}

func parseFlags() *Config {
	config := &Config{}

	docType := flag.String("type", "api", "Type of documentation to generate (api, readme, guide, example)")
	flag.StringVar(&config.OutputFile, "output", "", "Output file path (default: docs/{type}.md)")
	flag.StringVar(&config.SourceDir, "source", ".", "Source directory to analyze")
	flag.BoolVar(&config.GenerateAll, "all", false, "Generate all documentation types")

	flag.Parse()

	config.DocType = DocType(*docType)

	// Set default output if not specified
	if config.OutputFile == "" {
		config.OutputFile = fmt.Sprintf("docs/%s.md", config.DocType)
	}

	return config
}

func generateAllDocs(config *Config) {
	types := []DocType{DocTypeAPI, DocTypeREADME, DocTypeGuide, DocTypeExample}

	for _, docType := range types {
		c := *config
		c.DocType = docType
		c.OutputFile = fmt.Sprintf("docs/%s.md", docType)

		fmt.Printf("\n📄 Generating %s documentation...\n", docType)
		if err := generateDoc(&c); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: Failed to generate %s: %v\n", docType, err)
		}
	}
}

func generateDoc(config *Config) error {
	// Validate output path to prevent directory traversal
	if err := validateOutputPath(config.OutputFile); err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Extract code based on documentation type
	result, err := extractCodeForDocs(config)
	if err != nil {
		return fmt.Errorf("failed to extract code: %w", err)
	}

	fmt.Printf("✅ Extracted %d files (%d tokens)\n",
		len(result.ProjectOutput.Files), result.TokenCount)

	// Generate documentation prompt
	prompt := generateDocPrompt(config, result)

	// Save prompt for review/manual use
	promptFile := strings.TrimSuffix(config.OutputFile, filepath.Ext(config.OutputFile)) + "-prompt.txt"
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		return fmt.Errorf("failed to save prompt: %w", err)
	}

	// Save extracted context
	contextFile := strings.TrimSuffix(config.OutputFile, filepath.Ext(config.OutputFile)) + "-context.ptx"
	if err := os.WriteFile(contextFile, []byte(result.FormattedOutput), 0644); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}

	fmt.Printf("\n💾 Files saved:\n")
	fmt.Printf("   - Prompt: %s\n", promptFile)
	fmt.Printf("   - Context: %s\n", contextFile)
	fmt.Printf("   - Output template: %s\n", config.OutputFile)

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(config.OutputFile), 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate sample output structure
	if err := generateDocTemplate(config); err != nil {
		return fmt.Errorf("failed to generate template: %w", err)
	}

	fmt.Println("\n💡 Next Steps:")
	fmt.Println("   1. Review the generated prompt and context")
	fmt.Println("   2. Send to your AI assistant (Claude, GPT, etc.)")
	fmt.Println("   3. AI will generate documentation based on the code")
	fmt.Println("   4. Review and commit the generated documentation")

	return nil
}

func extractCodeForDocs(config *Config) (*promptext.Result, error) {
	var opts []promptext.Option

	switch config.DocType {
	case DocTypeAPI:
		// For API docs, extract public packages (exclude tests and internal)
		opts = []promptext.Option{
			promptext.WithExtensions(".go"),
			promptext.WithExcludes("*_test.go", "internal/", "testdata/", "vendor/"),
			promptext.WithTokenBudget(20000), // Larger budget for comprehensive API docs
		}

	case DocTypeREADME:
		// For README, focus on entry points and main packages
		opts = []promptext.Option{
			promptext.WithExtensions(".go", ".md"),
			promptext.WithRelevance("main", "cmd", "api", "example"),
			promptext.WithTokenBudget(10000),
		}

	case DocTypeGuide:
		// For guides, extract relevant features
		opts = []promptext.Option{
			promptext.WithExtensions(".go"),
			promptext.WithExcludes("*_test.go", "vendor/"),
			promptext.WithTokenBudget(15000),
		}

	case DocTypeExample:
		// For examples, extract example code and tests
		opts = []promptext.Option{
			promptext.WithExtensions(".go"),
			promptext.WithRelevance("example", "demo", "sample"),
			promptext.WithTokenBudget(8000),
		}
	}

	opts = append(opts, promptext.WithFormat(promptext.FormatMarkdown))

	return promptext.Extract(config.SourceDir, opts...)
}

func generateDocPrompt(config *Config, result *promptext.Result) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("# Documentation Generation Request: %s\n\n", config.DocType))

	switch config.DocType {
	case DocTypeAPI:
		prompt.WriteString("## Task\n")
		prompt.WriteString("Generate comprehensive API reference documentation from the provided code.\n\n")
		prompt.WriteString("## Requirements\n")
		prompt.WriteString("1. Document all public types, functions, and methods\n")
		prompt.WriteString("2. Include parameter descriptions and return values\n")
		prompt.WriteString("3. Provide usage examples for each major component\n")
		prompt.WriteString("4. Organize by package/module\n")
		prompt.WriteString("5. Use clear, beginner-friendly language\n\n")

	case DocTypeREADME:
		prompt.WriteString("## Task\n")
		prompt.WriteString("Generate or update the README.md file based on the current codebase.\n\n")
		prompt.WriteString("## Requirements\n")
		prompt.WriteString("1. **Overview**: Brief description of what the project does\n")
		prompt.WriteString("2. **Features**: List key features and capabilities\n")
		prompt.WriteString("3. **Installation**: Step-by-step installation instructions\n")
		prompt.WriteString("4. **Quick Start**: Simple example to get started\n")
		prompt.WriteString("5. **Usage**: Common usage patterns\n")
		prompt.WriteString("6. **API Reference**: Link to detailed API docs\n")
		prompt.WriteString("7. **Contributing**: How to contribute\n")
		prompt.WriteString("8. **License**: License information\n\n")

	case DocTypeGuide:
		prompt.WriteString("## Task\n")
		prompt.WriteString("Generate a comprehensive user guide based on the codebase.\n\n")
		prompt.WriteString("## Requirements\n")
		prompt.WriteString("1. **Introduction**: What the library does and why\n")
		prompt.WriteString("2. **Core Concepts**: Key abstractions and patterns\n")
		prompt.WriteString("3. **Getting Started**: Step-by-step tutorial\n")
		prompt.WriteString("4. **Common Tasks**: How to accomplish typical goals\n")
		prompt.WriteString("5. **Advanced Usage**: Power user features\n")
		prompt.WriteString("6. **Best Practices**: Recommended patterns\n")
		prompt.WriteString("7. **Troubleshooting**: Common issues and solutions\n\n")

	case DocTypeExample:
		prompt.WriteString("## Task\n")
		prompt.WriteString("Generate practical, runnable examples based on the code.\n\n")
		prompt.WriteString("## Requirements\n")
		prompt.WriteString("1. Create 3-5 complete, working examples\n")
		prompt.WriteString("2. Each example should be copy-paste ready\n")
		prompt.WriteString("3. Include comments explaining what each part does\n")
		prompt.WriteString("4. Show different use cases and patterns\n")
		prompt.WriteString("5. Add expected output for each example\n\n")
	}

	prompt.WriteString("## Source Code\n\n")
	prompt.WriteString("The following code has been extracted from the project:\n\n")
	prompt.WriteString(fmt.Sprintf("- **Files**: %d\n", len(result.ProjectOutput.Files)))
	prompt.WriteString(fmt.Sprintf("- **Tokens**: %d\n\n", result.TokenCount))

	if result.ProjectOutput.Metadata != nil {
		prompt.WriteString("## Project Information\n\n")
		prompt.WriteString(fmt.Sprintf("- **Language**: %s\n", result.ProjectOutput.Metadata.Language))
		if result.ProjectOutput.Metadata.Version != "" {
			prompt.WriteString(fmt.Sprintf("- **Version**: %s\n", result.ProjectOutput.Metadata.Version))
		}
		prompt.WriteString("\n")
	}

	prompt.WriteString("## Output Format\n\n")
	prompt.WriteString("Generate well-structured Markdown documentation with:\n")
	prompt.WriteString("- Clear headings and sections\n")
	prompt.WriteString("- Code blocks with syntax highlighting\n")
	prompt.WriteString("- Links between related sections\n")
	prompt.WriteString("- Table of contents for longer docs\n")
	prompt.WriteString("- Emoji for visual appeal (sparingly)\n\n")

	prompt.WriteString("---\n\n")
	prompt.WriteString("_The extracted code context is provided in the accompanying file._\n")

	return prompt.String()
}

func generateDocTemplate(config *Config) error {
	var template string

	switch config.DocType {
	case DocTypeAPI:
		template = `# API Reference

> This documentation is auto-generated. To regenerate: ` + "`go run examples/doc-generator/main.go --type api`" + `

## Table of Contents

- [Overview](#overview)
- [Packages](#packages)
- [Types](#types)
- [Functions](#functions)

## Overview

[AI will generate overview here based on the code]

## Packages

### Package: main

[AI will document package details here]

## Types

### Type: Example

[AI will document type details here]

## Functions

### Function: Example()

[AI will document function details here]

---

_Generated on: ` + "{{DATE}}" + `_
_Last updated: ` + "{{DATE}}" + `_
`

	case DocTypeREADME:
		template = `# Project Name

> Brief description here

## Features

- Feature 1
- Feature 2
- Feature 3

## Installation

` + "```bash" + `
go get github.com/example/project
` + "```" + `

## Quick Start

` + "```go" + `
// Example code here
` + "```" + `

## Usage

[Detailed usage instructions]

## Documentation

- [API Reference](docs/api.md)
- [User Guide](docs/guide.md)
- [Examples](docs/example.md)

## Contributing

[Contributing guidelines]

## License

[License information]
`

	case DocTypeGuide:
		template = `# User Guide

## Introduction

[Introduction to the project]

## Core Concepts

[Key concepts and patterns]

## Getting Started

[Step-by-step tutorial]

## Common Tasks

[How to accomplish typical goals]

## Advanced Usage

[Power user features]

## Best Practices

[Recommended patterns]

## Troubleshooting

[Common issues and solutions]
`

	case DocTypeExample:
		template = `# Examples

## Example 1: Basic Usage

` + "```go" + `
// Code here
` + "```" + `

**Output:**
` + "```" + `
// Expected output
` + "```" + `

## Example 2: Advanced Usage

[More examples...]
`
	}

	return os.WriteFile(config.OutputFile, []byte(template), 0644)
}

// validateOutputPath ensures the output path is safe and doesn't escape the current directory
func validateOutputPath(path string) error {
	// Clean the path to resolve .. and .
	cleaned := filepath.Clean(path)

	// Get absolute path for both current dir and output path
	cwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	// Make output path absolute
	absPath := cleaned
	if !filepath.IsAbs(cleaned) {
		absPath = filepath.Join(cwd, cleaned)
	}

	// Ensure the output path is under current directory or in typical doc locations
	relPath, err := filepath.Rel(cwd, absPath)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	// Prevent paths that escape current directory
	if strings.HasPrefix(relPath, "..") {
		return fmt.Errorf("output path cannot be outside current directory: %s", path)
	}

	return nil
}
