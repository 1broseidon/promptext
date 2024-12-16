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
	XMLName       xml.Name       `xml:"project"`
	DirectoryTree *DirectoryNode `xml:"directoryTree"`
	GitInfo       *GitInfo       `xml:"gitInfo,omitempty"`
	Metadata      *Metadata      `xml:"metadata,omitempty"`
	Files         []FileInfo     `xml:"files>file,omitempty"`
}

// Helper function to convert DirectoryNode to markdown string
func (d *DirectoryNode) ToMarkdown(level int) string {
	var sb strings.Builder
	indent := strings.Repeat("  ", level)
	prefix := "├──"
	if level > 0 {
		prefix = "└──"
	}

	if level > 0 {
		sb.WriteString(fmt.Sprintf("%s%s %s", indent, prefix, d.Name))
		if d.Type == "dir" {
			sb.WriteString("/")
		}
		sb.WriteString("\n")
	}

	if d.Children != nil {
		for i, child := range d.Children {
			if i == len(d.Children)-1 {
				sb.WriteString(child.ToMarkdown(level + 1))
			} else {
				sb.WriteString(child.ToMarkdown(level + 1))
			}
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
