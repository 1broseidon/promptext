package promptext

import (
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/processor"
)

// Result contains the output of a code extraction operation.
// It includes both the raw structured data and formatted output.
type Result struct {
	// ProjectOutput contains the structured project data
	ProjectOutput *ProjectOutput

	// FormattedOutput contains the output formatted according to the selected format
	FormattedOutput string

	// TokenCount is the estimated token count for the included files
	TokenCount int

	// TotalTokens is the total estimated tokens if all files were included
	TotalTokens int

	// ExcludedFiles is the number of files excluded due to token budget or relevance
	ExcludedFiles int

	// ExcludedFileList contains details about excluded files
	ExcludedFileList []ExcludedFileInfo
}

// ExcludedFileInfo contains information about an excluded file.
type ExcludedFileInfo struct {
	Path   string
	Tokens int
}

// ProjectOutput represents the complete structured output of a project extraction.
// This is the main data structure that contains all project information.
type ProjectOutput struct {
	// DirectoryTree is the hierarchical directory structure
	DirectoryTree *DirectoryNode

	// GitInfo contains git repository information (if available)
	GitInfo *GitInfo

	// Metadata contains project metadata (language, version, dependencies)
	Metadata *Metadata

	// Files contains the actual file contents and metadata
	Files []FileInfo

	// FileStats contains statistics about the processed files
	FileStats *FileStatistics

	// Budget contains token budget and truncation information
	Budget *BudgetInfo

	// FilterConfig describes the filter configuration used
	FilterConfig *FilterConfig
}

// DirectoryNode represents a node in the directory tree hierarchy.
type DirectoryNode struct {
	Name     string
	Type     string // "file" or "dir"
	Children []*DirectoryNode
}

// GitInfo contains git repository information.
type GitInfo struct {
	Branch        string
	CommitHash    string
	CommitMessage string
}

// Metadata contains project metadata information.
type Metadata struct {
	Language     string
	Version      string
	Dependencies []string
}

// FileInfo represents a single file and its contents.
type FileInfo struct {
	Path       string
	Content    string
	Tokens     int
	Truncation *TruncationInfo
}

// TruncationInfo describes how a file was truncated.
type TruncationInfo struct {
	Mode           string
	OriginalTokens int
}

// FileStatistics contains statistics about the processed files.
type FileStatistics struct {
	TotalFiles   int
	TotalLines   int
	PackageCount int
}

// BudgetInfo tracks token budget and truncation statistics.
type BudgetInfo struct {
	MaxTokens       int
	EstimatedTokens int
	FileTruncations int
}

// FilterConfig describes the filter configuration used to generate the output.
type FilterConfig struct {
	Includes []string
	Excludes []string
}

// As converts the result to a different output format.
// This is useful when you want to convert already-extracted data to a different format
// without re-processing the files.
//
// Example:
//
//	result, _ := promptext.Extract(".", WithFormat(promptext.FormatPTX))
//	markdownOutput, _ := result.As(promptext.FormatMarkdown)
//	jsonlOutput, _ := result.As(promptext.FormatJSONL)
func (r *Result) As(format Format) (string, error) {
	formatter, err := GetFormatter(string(format))
	if err != nil {
		return "", err
	}
	return formatter.Format(r.ProjectOutput)
}

// fromInternalProcessResult converts internal processor.ProcessResult to public Result
func fromInternalProcessResult(internal *processor.ProcessResult, formattedOutput string) *Result {
	if internal == nil {
		return nil
	}

	result := &Result{
		ProjectOutput:    fromInternalProjectOutput(internal.ProjectOutput),
		FormattedOutput:  formattedOutput,
		TokenCount:       internal.TokenCount,
		TotalTokens:      internal.TotalTokens,
		ExcludedFiles:    internal.ExcludedFiles,
		ExcludedFileList: make([]ExcludedFileInfo, len(internal.ExcludedFileList)),
	}

	for i, excluded := range internal.ExcludedFileList {
		result.ExcludedFileList[i] = ExcludedFileInfo{
			Path:   excluded.Path,
			Tokens: excluded.Tokens,
		}
	}

	return result
}

// fromInternalProjectOutput converts internal format.ProjectOutput to public ProjectOutput
func fromInternalProjectOutput(internal *format.ProjectOutput) *ProjectOutput {
	if internal == nil {
		return nil
	}

	output := &ProjectOutput{}

	// Convert DirectoryTree
	if internal.DirectoryTree != nil {
		output.DirectoryTree = fromInternalDirectoryNode(internal.DirectoryTree)
	}

	// Convert GitInfo
	if internal.GitInfo != nil {
		output.GitInfo = &GitInfo{
			Branch:        internal.GitInfo.Branch,
			CommitHash:    internal.GitInfo.CommitHash,
			CommitMessage: internal.GitInfo.CommitMessage,
		}
	}

	// Convert Metadata
	if internal.Metadata != nil {
		output.Metadata = &Metadata{
			Language:     internal.Metadata.Language,
			Version:      internal.Metadata.Version,
			Dependencies: internal.Metadata.Dependencies,
		}
	}

	// Convert Files
	output.Files = make([]FileInfo, len(internal.Files))
	for i, file := range internal.Files {
		output.Files[i] = FileInfo{
			Path:    file.Path,
			Content: file.Content,
			Tokens:  file.Tokens,
		}
		if file.Truncation != nil {
			output.Files[i].Truncation = &TruncationInfo{
				Mode:           file.Truncation.Mode,
				OriginalTokens: file.Truncation.OriginalTokens,
			}
		}
	}

	// Convert FileStats
	if internal.FileStats != nil {
		output.FileStats = &FileStatistics{
			TotalFiles:   internal.FileStats.TotalFiles,
			TotalLines:   internal.FileStats.TotalLines,
			PackageCount: internal.FileStats.PackageCount,
		}
	}

	// Convert Budget
	if internal.Budget != nil {
		output.Budget = &BudgetInfo{
			MaxTokens:       internal.Budget.MaxTokens,
			EstimatedTokens: internal.Budget.EstimatedTokens,
			FileTruncations: internal.Budget.FileTruncations,
		}
	}

	// Convert FilterConfig
	if internal.FilterConfig != nil {
		output.FilterConfig = &FilterConfig{
			Includes: internal.FilterConfig.Includes,
			Excludes: internal.FilterConfig.Excludes,
		}
	}

	return output
}

// fromInternalDirectoryNode converts internal format.DirectoryNode to public DirectoryNode
func fromInternalDirectoryNode(internal *format.DirectoryNode) *DirectoryNode {
	if internal == nil {
		return nil
	}

	node := &DirectoryNode{
		Name: internal.Name,
		Type: internal.Type,
	}

	if len(internal.Children) > 0 {
		node.Children = make([]*DirectoryNode, len(internal.Children))
		for i, child := range internal.Children {
			node.Children[i] = fromInternalDirectoryNode(child)
		}
	}

	return node
}
