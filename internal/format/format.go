package format

import (
	"encoding/xml"
	"fmt"
	"strings"
)

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatXML      OutputFormat = "xml"
)

// DirectoryNode represents a node in the directory tree
type DirectoryNode struct {
	Name     string           `xml:"name,attr"`
	Type     string           `xml:"type,attr"` // "file" or "dir"
	Children []*DirectoryNode `xml:"node,omitempty"`
}

type ProjectOutput struct {
	XMLName       xml.Name         `xml:"project"`
	DirectoryTree *DirectoryNode   `xml:"directoryTree"`
	GitInfo       *GitInfo         `xml:"gitInfo,omitempty"`
	Metadata      *Metadata        `xml:"metadata,omitempty"`
	Files         []FileInfo       `xml:"files>file,omitempty"`
	Overview      *ProjectOverview `xml:"overview,omitempty"`
	FileStats     *FileStatistics  `xml:"fileStats,omitempty"`
	Dependencies  *DependencyInfo  `xml:"dependencies,omitempty"`
	Analysis      *ProjectAnalysis `xml:"analysis,omitempty"`
}

type ProjectOverview struct {
	Description string   `xml:"description"`
	Purpose     string   `xml:"purpose"`
	Features    []string `xml:"features>feature,omitempty"`
}

type FileStatistics struct {
	TotalFiles   int            `xml:"totalFiles"`
	FilesByType  map[string]int `xml:"-"` // Exclude from direct XML marshaling
	TotalLines   int            `xml:"totalLines"`
	PackageCount int            `xml:"packageCount"`
}

type DependencyInfo struct {
	Imports   map[string][]string `xml:"imports>file"`
	Packages  []string            `xml:"packages>package"`
	CoreFiles []string            `xml:"coreFiles>file"`
}

type ProjectAnalysis struct {
	EntryPoints   map[string]string `xml:"entryPoints,omitempty"`
	ConfigFiles   map[string]string `xml:"configFiles,omitempty"`
	CoreFiles     map[string]string `xml:"coreFiles,omitempty"`
	TestFiles     map[string]string `xml:"testFiles,omitempty"`
	Documentation map[string]string `xml:"documentation,omitempty"`
}

// Helper function to convert DirectoryNode to markdown string
func (d *DirectoryNode) ToMarkdown(level int) string {
	var sb strings.Builder

	// Skip root node name but include its children
	if level > 0 {
		indent := strings.Repeat("  ", level-1)
		prefix := "└── "
		if d.Type == "dir" {
			sb.WriteString(fmt.Sprintf("%s%s%s/\n", indent, prefix, d.Name))
		} else {
			sb.WriteString(fmt.Sprintf("%s%s%s\n", indent, prefix, d.Name))
		}
	}

	// Process children
	if d.Children != nil {
		for _, child := range d.Children {
			// For root level, don't increment the level
			nextLevel := level
			if level > 0 {
				nextLevel++
			}
			sb.WriteString(child.ToMarkdown(nextLevel))
		}
	}

	return sb.String()
}

type GitInfo struct {
	Branch        string `xml:"branch"`
	CommitHash    string `xml:"commitHash"`
	CommitMessage string `xml:"commitMessage"`
}

type Metadata struct {
	Language     string   `xml:"language"`
	Version      string   `xml:"version"`
	Dependencies []string `xml:"dependencies>dependency,omitempty"`
}

type FileInfo struct {
	Path    string `xml:"path,attr"`
	Content string `xml:"content"`
}

// Formatter interface for different output formats
type Formatter interface {
	Format(project *ProjectOutput) (string, error)
}

// Get appropriate formatter based on format string
func GetFormatter(format string) (Formatter, error) {
	switch OutputFormat(format) {
	case FormatMarkdown:
		return &MarkdownFormatter{}, nil
	case FormatXML:
		return &XMLFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s (supported: markdown, xml)", format)
	}
}
