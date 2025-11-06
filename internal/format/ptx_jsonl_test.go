package format

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestPTXFormatterFormatIncludesManifestAndStructure(t *testing.T) {
	formatter := &PTXFormatter{}
	project := &ProjectOutput{
		Metadata: &Metadata{
			Language:     "Go",
			Version:      "1.21",
			Dependencies: []string{"dep1", "dep2"},
		},
		GitInfo: &GitInfo{
			Branch:        "main",
			CommitHash:    "abc123",
			CommitMessage: "fix bug",
		},
		Budget: &BudgetInfo{
			MaxTokens:       8000,
			EstimatedTokens: 6400,
			FileTruncations: 2,
		},
		FilterConfig: &FilterConfig{
			Includes: []string{"*.go"},
			Excludes: []string{"vendor"},
		},
		FileStats: &FileStatistics{
			TotalFiles:   2,
			TotalLines:   42,
			PackageCount: 1,
			FilesByType: map[string]int{
				".go": 2,
			},
		},
		DirectoryTree: &DirectoryNode{
			Name: "root",
			Type: "dir",
			Children: []*DirectoryNode{
				{Name: "main.go", Type: "file"},
				{Name: "pkg", Type: "dir", Children: []*DirectoryNode{{Name: "pkg.go", Type: "file"}}},
			},
		},
		Analysis: &ProjectAnalysis{
			EntryPoints: map[string]string{"cmd/main.go": "main"},
			ConfigFiles: map[string]string{"config.yaml": "app config"},
			CoreFiles:   map[string]string{"pkg/pkg.go": "core"},
			TestFiles:   map[string]string{"pkg/pkg_test.go": "tests"},
		},
		Files: []FileInfo{
			{Path: "pkg/pkg.go", Content: "package pkg\n", Tokens: 100},
			{Path: "main.go", Content: "package main\n", Truncation: &TruncationInfo{Mode: "head", OriginalTokens: 200}},
		},
	}

	out, err := formatter.Format(project)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	checks := []string{
		"schema: ptx/v2.0",
		"language: Go",
		"dependencies[2]: dep1,dep2",
		"branch: main",
		"max_tokens: 8000",
		"includes[1]: *.go",
		"totalFiles: 2",
		"structure:",
		"entryPoints[1]{desc,path}:",
		"configFiles[1]{desc,path}:",
		"path: pkg/pkg.go",
		"tokens: 100",
		"mode: head",
	}
	for _, want := range checks {
		if !strings.Contains(out, want) {
			t.Fatalf("expected output to contain %q, got:\n%s", want, out)
		}
	}
}

func TestPTXFormatterTreeToDirectoryMap(t *testing.T) {
	formatter := &PTXFormatter{}
	tree := &DirectoryNode{
		Name: "root",
		Type: "dir",
		Children: []*DirectoryNode{
			{Name: "main.go", Type: "file"},
			{Name: "pkg", Type: "dir", Children: []*DirectoryNode{{Name: "pkg.go", Type: "file"}}},
		},
	}

	structure := formatter.treeToDirectoryMap(tree)
	if len(structure) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(structure))
	}
	filesVal, ok := structure[""]
	if !ok {
		t.Fatalf("missing root directory entry: %#v", structure)
	}
	files, ok := filesVal.([]string)
	if !ok || len(files) != 1 || files[0] != "main.go" {
		t.Fatalf("unexpected root files: %#v", filesVal)
	}
	pkgVal, ok := structure["pkg"]
	if !ok {
		t.Fatalf("missing pkg directory entry: %#v", structure)
	}
	pkgFiles, ok := pkgVal.([]string)
	if !ok || len(pkgFiles) != 1 || pkgFiles[0] != "pkg.go" {
		t.Fatalf("unexpected pkg files: %#v", pkgVal)
	}
}

func TestPTXFormatterMapToList(t *testing.T) {
	formatter := &PTXFormatter{}
	input := map[string]string{
		"a.go": "main",
		"b.go": "helper",
	}
	list := formatter.mapToList(input)
	if len(list) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(list))
	}
	paths := make(map[string]bool)
	for _, entry := range list {
		path, _ := entry["path"].(string)
		paths[path] = true
	}
	if !paths["a.go"] || !paths["b.go"] {
		t.Fatalf("missing paths in list: %#v", paths)
	}
}

func TestTOONStrictFormatterEscapesContent(t *testing.T) {
	formatter := &TOONStrictFormatter{}
	project := &ProjectOutput{
		GitInfo: &GitInfo{CommitMessage: "line1\nline2"},
		Files:   []FileInfo{{Path: "main.go", Content: "func main() {\n}\n"}},
	}

	out, err := formatter.Format(project)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}
	normalized := strings.ReplaceAll(out, "\n", "\\n")
	if !strings.Contains(normalized, "line1\\\\nline2") {
		t.Fatalf("expected escaped commit message, got %s", normalized)
	}
	if !strings.Contains(normalized, "func main() {\\\\n}\\\\n\",main.go") {
		t.Fatalf("expected escaped file content, got %s", normalized)
	}
}

func TestJSONLFormatterProducesDeterministicLines(t *testing.T) {
	formatter := &JSONLFormatter{}
	project := &ProjectOutput{
		Metadata: &Metadata{Language: "Go"},
		GitInfo:  &GitInfo{Branch: "main", CommitHash: "abc"},
		Budget:   &BudgetInfo{MaxTokens: 5000, EstimatedTokens: 4000, FileTruncations: 1},
		FilterConfig: &FilterConfig{
			Includes: []string{"*.go"},
		},
		Files: []FileInfo{
			{Path: "b.go", Content: "package b", Tokens: 10},
			{Path: "a.go", Content: "package a", Truncation: &TruncationInfo{Mode: "tail", OriginalTokens: 20}},
		},
	}

	out, err := formatter.Format(project)
	if err != nil {
		t.Fatalf("Format error: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 6 {
		t.Fatalf("expected 6 lines, got %d: %q", len(lines), lines)
	}

	typeLine := func(idx int) string {
		var payload map[string]interface{}
		if err := json.Unmarshal([]byte(lines[idx]), &payload); err != nil {
			t.Fatalf("line %d not valid json: %v", idx, err)
		}
		t.Helper()
		return payload["type"].(string)
	}

	if got := typeLine(0); got != "metadata" {
		t.Fatalf("line 0 type = %s, want metadata", got)
	}
	if got := typeLine(1); got != "git" {
		t.Fatalf("line 1 type = %s, want git", got)
	}
	if got := typeLine(2); got != "budget" {
		t.Fatalf("line 2 type = %s, want budget", got)
	}
	if got := typeLine(3); got != "filters" {
		t.Fatalf("line 3 type = %s, want filters", got)
	}
	if got := typeLine(4); got != "file" || !strings.Contains(lines[4], "\"path\":\"a.go\"") {
		t.Fatalf("line 4 should be file for a.go, got %s", lines[4])
	}
	if !strings.Contains(lines[4], "\"truncation\"") {
		t.Fatalf("expected truncation metadata in first file line")
	}
	if got := typeLine(5); got != "file" || !strings.Contains(lines[5], "\"path\":\"b.go\"") {
		t.Fatalf("line 5 should be file for b.go, got %s", lines[5])
	}
	if !strings.Contains(lines[5], "\"tokens\":10") {
		t.Fatalf("expected token count in second file line")
	}
}
