package promptext

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/log"
	"github.com/1broseidon/promptext/internal/processor"
)

// Version is the current version of the promptext library.
const Version = "0.7.0-alpha"

// Extract is the main entry point for extracting code context from a directory.
// It processes the directory according to the provided options and returns a Result
// containing both structured data and formatted output.
//
// The dir parameter can be an absolute or relative path. If empty or ".", the current
// working directory is used.
//
// Extract uses sensible defaults that work out of the box:
//   - All supported file types are included
//   - .gitignore patterns are respected
//   - Built-in filtering rules exclude binaries, lockfiles, etc.
//   - Output format is PTX (recommended for AI assistants)
//
// Example - Simple usage:
//
//	result, err := promptext.Extract(".")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.FormattedOutput)
//
// Example - With options:
//
//	result, err := promptext.Extract("/path/to/project",
//	    promptext.WithExtensions(".go", ".mod"),
//	    promptext.WithTokenBudget(8000),
//	    promptext.WithFormat(promptext.FormatJSONL),
//	)
//
// Example - Relevance filtering:
//
//	result, err := promptext.Extract(".",
//	    promptext.WithRelevance("auth", "login"),
//	    promptext.WithTokenBudget(5000),
//	)
func Extract(dir string, opts ...Option) (*Result, error) {
	// Create extractor with options
	extractor := NewExtractor(opts...)

	// Extract from directory
	return extractor.Extract(dir)
}

// Extractor provides a reusable extractor that can process multiple directories
// with the same configuration. This is useful when you need to extract code
// from multiple projects with consistent settings.
//
// Example:
//
//	extractor := promptext.NewExtractor(
//	    promptext.WithExtensions(".go"),
//	    promptext.WithTokenBudget(8000),
//	)
//	result1, _ := extractor.Extract("/path/to/project1")
//	result2, _ := extractor.Extract("/path/to/project2")
type Extractor struct {
	config *config
}

// NewExtractor creates a new Extractor with the given options.
// The extractor can be reused to process multiple directories with the same configuration.
//
// Example:
//
//	extractor := promptext.NewExtractor(
//	    promptext.WithExtensions(".go", ".mod"),
//	    promptext.WithExcludes("vendor/", "*.test.go"),
//	)
//	result, err := extractor.Extract(".")
func NewExtractor(opts ...Option) *Extractor {
	cfg := newDefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}
	return &Extractor{config: cfg}
}

// Extract processes the specified directory and returns the extraction result.
// The directory path can be absolute or relative.
//
// Example:
//
//	extractor := promptext.NewExtractor(promptext.WithFormat(promptext.FormatPTX))
//	result, err := extractor.Extract("/path/to/project")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(result.FormattedOutput)
func (e *Extractor) Extract(dir string) (*Result, error) {
	// Validate and resolve directory path
	absPath, err := resolvePath(dir)
	if err != nil {
		return nil, &DirectoryError{
			Path: dir,
			Err:  err,
		}
	}

	// Check if directory exists and is accessible
	if err := validateDirectory(absPath); err != nil {
		return nil, &DirectoryError{
			Path: absPath,
			Err:  err,
		}
	}

	// Configure logging
	if e.config.debug {
		log.Enable()
		log.SetColorEnabled(true)
	}

	// Create filter options
	filterOpts := filter.Options{
		Includes:        e.config.extensions,
		Excludes:        e.config.excludes,
		UseDefaultRules: e.config.useDefaultRules,
		UseGitIgnore:    e.config.gitignore,
	}

	// Create filter
	f := filter.New(filterOpts)

	// Create processor configuration
	procConfig := processor.Config{
		DirPath:           absPath,
		Extensions:        e.config.extensions,
		Excludes:          e.config.excludes,
		GitIgnore:         e.config.gitignore,
		Filter:            f,
		RelevanceKeywords: e.config.relevanceKeywords,
		MaxTokens:         e.config.tokenBudget,
	}

	// Process directory
	procResult, err := processor.ProcessDirectory(procConfig, e.config.verbose)
	if err != nil {
		return nil, fmt.Errorf("error processing directory: %w", err)
	}

	// Check if any files were processed
	if len(procResult.ProjectOutput.Files) == 0 {
		return nil, ErrNoFilesMatched
	}

	// Get formatter
	formatter, err := GetFormatter(string(e.config.format))
	if err != nil {
		return nil, err
	}

	// Format output
	formattedOutput, err := formatter.Format(fromInternalProjectOutput(procResult.ProjectOutput))
	if err != nil {
		return nil, &FormatError{
			Format: string(e.config.format),
			Err:    err,
		}
	}

	// Convert to public Result type
	result := fromInternalProcessResult(procResult, formattedOutput)

	return result, nil
}

// WithExtensions is a convenience method to add extensions to the extractor.
// It returns a new Extractor with the updated configuration.
//
// Example:
//
//	extractor := promptext.NewExtractor().WithExtensions(".go", ".mod")
func (e *Extractor) WithExtensions(extensions ...string) *Extractor {
	e.config.extensions = extensions
	return e
}

// WithExcludes is a convenience method to add exclude patterns to the extractor.
// It returns a new Extractor with the updated configuration.
//
// Example:
//
//	extractor := promptext.NewExtractor().WithExcludes("vendor/", "*.test.go")
func (e *Extractor) WithExcludes(patterns ...string) *Extractor {
	e.config.excludes = patterns
	return e
}

// WithFormat is a convenience method to set the output format.
// It returns a new Extractor with the updated configuration.
//
// Example:
//
//	extractor := promptext.NewExtractor().WithFormat(promptext.FormatJSONL)
func (e *Extractor) WithFormat(format Format) *Extractor {
	e.config.format = format
	return e
}

// resolvePath resolves a directory path to an absolute path.
// It handles special cases like "." and empty string (current directory).
func resolvePath(dir string) (string, error) {
	// Handle empty or current directory
	if dir == "" || dir == "." {
		cwd, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %w", err)
		}
		return cwd, nil
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(dir)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return absPath, nil
}

// validateDirectory checks if a directory exists and is accessible.
func validateDirectory(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("%w: directory does not exist", ErrInvalidDirectory)
		}
		return fmt.Errorf("%w: %v", ErrInvalidDirectory, err)
	}

	if !info.IsDir() {
		return fmt.Errorf("%w: path is not a directory", ErrInvalidDirectory)
	}

	return nil
}

// Helper function to join extensions for internal use
func joinExtensions(extensions []string) string {
	return strings.Join(extensions, ",")
}

// Helper function to join excludes for internal use
func joinExcludes(excludes []string) string {
	return strings.Join(excludes, ",")
}
