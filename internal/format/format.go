package format

import (
	"encoding/xml"
	"fmt"
)

type OutputFormat string

const (
	FormatMarkdown OutputFormat = "markdown"
	FormatXML      OutputFormat = "xml"
	FormatJSON     OutputFormat = "json"
)

type ProjectOutput struct {
	XMLName       xml.Name   `xml:"project" json:"-"`
	DirectoryTree string     `xml:"directoryTree" json:"directoryTree"`
	GitInfo       *GitInfo   `xml:"gitInfo,omitempty" json:"gitInfo,omitempty"`
	Metadata      *Metadata  `xml:"metadata,omitempty" json:"metadata,omitempty"`
	Files         []FileInfo `xml:"files>file,omitempty" json:"files,omitempty"`
}

type GitInfo struct {
	Branch        string `xml:"branch" json:"branch"`
	CommitHash    string `xml:"commitHash" json:"commitHash"`
	CommitMessage string `xml:"commitMessage" json:"commitMessage"`
}

type Metadata struct {
	Language     string   `xml:"language" json:"language"`
	Version      string   `xml:"version" json:"version"`
	Dependencies []string `xml:"dependencies>dependency,omitempty" json:"dependencies,omitempty"`
}

type FileInfo struct {
	Path    string `xml:"path,attr" json:"path"`
	Content string `xml:"content" json:"content"`
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
	case FormatJSON:
		return &JSONFormatter{}, nil
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}
