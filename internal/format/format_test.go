package format

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/1broseidon/promptext/internal/references"
)

func TestGetFormatter(t *testing.T) {
	tests := []struct {
		name        string
		format      string
		wantType    Formatter
		wantErr     bool
		errContains string
	}{
		{
			name:     "markdown formatter",
			format:   "markdown",
			wantType: &MarkdownFormatter{},
			wantErr:  false,
		},
		{
			name:     "xml formatter",
			format:   "xml",
			wantType: &XMLFormatter{},
			wantErr:  false,
		},
		{
			name:        "invalid formatter",
			format:      "invalid",
			wantType:    nil,
			wantErr:     true,
			errContains: "unsupported format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetFormatter(tt.format)
			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				} else if !strings.Contains(err.Error(), tt.errContains) {
					t.Errorf("error %q should contain %q", err.Error(), tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got == nil {
				t.Fatal("got nil formatter")
			}
			switch tt.wantType.(type) {
			case *MarkdownFormatter:
				if _, ok := got.(*MarkdownFormatter); !ok {
					t.Error("expected MarkdownFormatter")
				}
			case *XMLFormatter:
				if _, ok := got.(*XMLFormatter); !ok {
					t.Error("expected XMLFormatter")
				}
			}
		})
	}
}

func TestMarkdownFormatter_Format(t *testing.T) {
	formatter := &MarkdownFormatter{}
	
	tests := []struct {
		name    string
		input   *ProjectOutput
		want    []string // Strings that should be present in output
		unwant  []string // Strings that should not be present in output
		wantErr bool
	}{
		{
			name: "basic project",
			input: &ProjectOutput{
				Overview: &ProjectOverview{
					Description: "Test Project",
					Features:    []string{"Feature 1", "Feature 2"},
				},
				FileStats: &FileStatistics{
					TotalFiles: 10,
					TotalLines: 500,
					FilesByType: map[string]int{
						".go": 5,
						".md": 2,
					},
				},
				GitInfo: &GitInfo{
					Branch:        "main",
					CommitHash:    "abc123",
					CommitMessage: "test commit",
				},
			},
			want: []string{
				"# Project Overview",
				"Test Project",
				"Feature 1",
				"Feature 2",
				"Total Files: 10",
				"Total Lines: 500",
				".go: 5 files",
				"Branch: main",
				"Commit: abc123",
			},
			unwant: []string{
				"<project>",
				"</project>",
			},
		},
		{
			name: "with source files",
			input: &ProjectOutput{
				Files: []FileInfo{
					{
						Path:    "test.go",
						Content: "package main\n\nfunc main() {}\n",
						References: &references.ReferenceMap{
							Internal: map[string][]string{
								"test.go": {"other.go"},
							},
						},
					},
				},
			},
			want: []string{
				"## Source Files",
				"### test.go",
				"```go",
				"package main",
				"func main()",
				"```",
				"References:",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatter.Format(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("Format() output should contain %q", want)
				}
			}

			for _, unwant := range tt.unwant {
				if strings.Contains(got, unwant) {
					t.Errorf("Format() output should not contain %q", unwant)
				}
			}
		})
	}
}

func TestXMLFormatter_Format(t *testing.T) {
	formatter := &XMLFormatter{}
	
	tests := []struct {
		name    string
		input   *ProjectOutput
		want    []string
		wantErr bool
	}{
		{
			name: "basic project",
			input: &ProjectOutput{
				Overview: &ProjectOverview{
					Description: "Test Project",
					Features:    []string{"Feature 1", "Feature 2"},
				},
				FileStats: &FileStatistics{
					TotalFiles: 10,
					TotalLines: 500,
					FilesByType: map[string]int{
						".go": 5,
						".md": 2,
					},
				},
			},
			want: []string{
				`<?xml version="1.0" encoding="UTF-8"?>`,
				"<project>",
				"<overview>",
				"<description><![CDATA[Test Project]]></description>",
				"<feature>Feature 1</feature>",
				"<feature>Feature 2</feature>",
				"<fileStats>",
				"<totalFiles>10</totalFiles>",
				"<totalLines>500</totalLines>",
				"</project>",
			},
		},
		{
			name: "with files",
			input: &ProjectOutput{
				Files: []FileInfo{
					{
						Path:    "test.go",
						Content: "package main\n\nfunc main() {}\n",
					},
				},
			},
			want: []string{
				"<files>",
				`<file path="test.go"`,
				"<content><![CDATA[",
				"package main",
				"func main()",
				"]]></content>",
				"</file>",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := formatter.Format(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Format() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Verify it's valid XML
			if err := xml.Unmarshal([]byte(got), &ProjectOutput{}); err != nil {
				t.Errorf("Format() produced invalid XML: %v", err)
			}

			for _, want := range tt.want {
				if !strings.Contains(got, want) {
					t.Errorf("Format() output should contain %q", want)
				}
			}
		})
	}
}

func TestDirectoryNode_ToMarkdown(t *testing.T) {
	tests := []struct {
		name string
		node *DirectoryNode
		want string
	}{
		{
			name: "simple directory",
			node: &DirectoryNode{
				Name: "root",
				Type: "dir",
				Children: []*DirectoryNode{
					{
						Name: "file1.txt",
						Type: "file",
					},
					{
						Name: "dir1",
						Type: "dir",
						Children: []*DirectoryNode{
							{
								Name: "file2.txt",
								Type: "file",
							},
						},
					},
				},
			},
			want: "└── file1.txt\n└── dir1/\n  └── file2.txt\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.node.ToMarkdown(1)
			if got != tt.want {
				t.Errorf("ToMarkdown() = %q, want %q", got, tt.want)
			}
		})
	}
}
