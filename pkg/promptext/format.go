package promptext

import (
	"github.com/1broseidon/promptext/internal/format"
)

// Format represents the output format type for code extraction.
type Format string

// Supported output formats.
const (
	// FormatPTX is the PTX v2.0 format (TOON-based with enhanced manifest).
	// This is the recommended format for AI assistants.
	FormatPTX Format = "ptx"

	// FormatTOON is an alias for FormatPTX (backward compatibility).
	FormatTOON Format = "toon"

	// FormatTOONStrict is the TOON v1.3 strict compliance format with escaped strings.
	FormatTOONStrict Format = "toon-strict"

	// FormatJSONL is a machine-friendly JSONL format (one JSON object per line).
	FormatJSONL Format = "jsonl"

	// FormatMarkdown is a human-readable markdown format.
	FormatMarkdown Format = "markdown"

	// FormatXML is a machine-parseable XML format.
	FormatXML Format = "xml"
)

// Formatter is the interface that all output formatters must implement.
// This interface allows developers to create custom formatters for their specific needs.
//
// Example custom formatter:
//
//	type MyFormatter struct{}
//
//	func (f *MyFormatter) Format(output *ProjectOutput) (string, error) {
//	    // Custom formatting logic here
//	    return formatted, nil
//	}
//
//	RegisterFormatter("myformat", &MyFormatter{})
type Formatter interface {
	Format(output *ProjectOutput) (string, error)
}

var customFormatters = make(map[string]Formatter)

// RegisterFormatter registers a custom formatter that can be used with the library.
// This allows developers to extend the library with their own output formats.
//
// Example:
//
//	type MyFormatter struct{}
//	func (f *MyFormatter) Format(output *ProjectOutput) (string, error) {
//	    return "custom format", nil
//	}
//	promptext.RegisterFormatter("myformat", &MyFormatter{})
//	result, _ := promptext.Extract(".", WithFormat("myformat"))
func RegisterFormatter(name string, formatter Formatter) {
	customFormatters[name] = formatter
}

// GetFormatter returns the appropriate formatter for the given format string.
// It first checks custom formatters, then falls back to built-in formatters.
func GetFormatter(formatStr string) (Formatter, error) {
	// Check custom formatters first
	if customFormatter, ok := customFormatters[formatStr]; ok {
		return customFormatter, nil
	}

	// Fall back to built-in formatters
	internalFormatter, err := format.GetFormatter(formatStr)
	if err != nil {
		return nil, &FormatError{
			Format: formatStr,
			Err:    ErrInvalidFormat,
		}
	}

	// Wrap the internal formatter to implement our public Formatter interface
	return &formatterAdapter{internal: internalFormatter}, nil
}

// formatterAdapter adapts internal formatters to the public Formatter interface
type formatterAdapter struct {
	internal format.Formatter
}

func (a *formatterAdapter) Format(output *ProjectOutput) (string, error) {
	// Convert public ProjectOutput to internal format.ProjectOutput
	internalOutput := toInternalProjectOutput(output)
	return a.internal.Format(internalOutput)
}

// toInternalProjectOutput converts public ProjectOutput to internal format.ProjectOutput
func toInternalProjectOutput(output *ProjectOutput) *format.ProjectOutput {
	if output == nil {
		return nil
	}

	internal := &format.ProjectOutput{}

	// Convert DirectoryTree
	if output.DirectoryTree != nil {
		internal.DirectoryTree = toInternalDirectoryNode(output.DirectoryTree)
	}

	// Convert GitInfo
	if output.GitInfo != nil {
		internal.GitInfo = &format.GitInfo{
			Branch:        output.GitInfo.Branch,
			CommitHash:    output.GitInfo.CommitHash,
			CommitMessage: output.GitInfo.CommitMessage,
		}
	}

	// Convert Metadata
	if output.Metadata != nil {
		internal.Metadata = &format.Metadata{
			Language:     output.Metadata.Language,
			Version:      output.Metadata.Version,
			Dependencies: output.Metadata.Dependencies,
		}
	}

	// Convert Files
	internal.Files = make([]format.FileInfo, len(output.Files))
	for i, file := range output.Files {
		internal.Files[i] = format.FileInfo{
			Path:    file.Path,
			Content: file.Content,
			Tokens:  file.Tokens,
		}
		if file.Truncation != nil {
			internal.Files[i].Truncation = &format.TruncationInfo{
				Mode:           file.Truncation.Mode,
				OriginalTokens: file.Truncation.OriginalTokens,
			}
		}
	}

	// Convert FileStats
	if output.FileStats != nil {
		internal.FileStats = &format.FileStatistics{
			TotalFiles:   output.FileStats.TotalFiles,
			TotalLines:   output.FileStats.TotalLines,
			PackageCount: output.FileStats.PackageCount,
		}
	}

	// Convert Budget
	if output.Budget != nil {
		internal.Budget = &format.BudgetInfo{
			MaxTokens:       output.Budget.MaxTokens,
			EstimatedTokens: output.Budget.EstimatedTokens,
			FileTruncations: output.Budget.FileTruncations,
		}
	}

	// Convert FilterConfig
	if output.FilterConfig != nil {
		internal.FilterConfig = &format.FilterConfig{
			Includes: output.FilterConfig.Includes,
			Excludes: output.FilterConfig.Excludes,
		}
	}

	return internal
}

// toInternalDirectoryNode converts public DirectoryNode to internal format.DirectoryNode
func toInternalDirectoryNode(node *DirectoryNode) *format.DirectoryNode {
	if node == nil {
		return nil
	}

	internal := &format.DirectoryNode{
		Name: node.Name,
		Type: node.Type,
	}

	if len(node.Children) > 0 {
		internal.Children = make([]*format.DirectoryNode, len(node.Children))
		for i, child := range node.Children {
			internal.Children[i] = toInternalDirectoryNode(child)
		}
	}

	return internal
}
