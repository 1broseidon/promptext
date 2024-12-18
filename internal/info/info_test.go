package info

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetProjectInfo(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "project-info-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"go.mod":          "module test\n\ngo 1.17\n\nrequire github.com/stretchr/testify v1.8.0",
		"main.go":         "package main\n\nfunc main() {}\n",
		"README.md":       "# Test Project",
		".gitignore":      "*.tmp\n",
		"internal/foo.go": "package internal\n",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Initialize test filter
	f := filter.New(filter.Options{
		Includes: []string{".go"},
		Excludes: []string{},
	})

	// Test GetProjectInfo
	t.Run("basic project structure", func(t *testing.T) {
		info, err := GetProjectInfo(tmpDir, f)
		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.NotNil(t, info.DirectoryTree)
		// Verify basic structure instead of specific temp dir name
		assert.NotEmpty(t, info.DirectoryTree.Name)
		assert.Equal(t, "dir", info.DirectoryTree.Type)
	})
}

func TestGenerateDirectoryTree(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "directory-tree-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directory structure
	files := []string{
		"main.go",
		"internal/pkg1/file1.go",
		"internal/pkg2/file2.go",
		"docs/README.md",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	f := filter.New(filter.Options{
		Includes: []string{".go", ".md"},
		Excludes: []string{},
	})

	t.Run("directory tree generation", func(t *testing.T) {
		tree, err := generateDirectoryTree(tmpDir, f)
		assert.NoError(t, err)
		assert.NotNil(t, tree)

		// Verify root node
		assert.Equal(t, filepath.Base(tmpDir), tree.Name)
		assert.Equal(t, "dir", tree.Type)

		// Verify directory structure
		foundMain := false
		foundInternal := false
		foundDocs := false

		for _, child := range tree.Children {
			switch child.Name {
			case "main.go":
				foundMain = true
			case "internal":
				foundInternal = true
			case "docs":
				foundDocs = true
			}
		}

		assert.True(t, foundMain, "main.go not found")
		assert.True(t, foundInternal, "internal/ not found")
		assert.True(t, foundDocs, "docs/ not found")
	})
}

func TestGetProjectMetadata(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "metadata-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Go project", func(t *testing.T) {
		// Create go.mod file
		goMod := `module test
go 1.17
require github.com/stretchr/testify v1.8.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
		assert.NoError(t, err)

		metadata, err := getProjectMetadata(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "Go", metadata.Language)
		assert.Equal(t, "1.17", metadata.Version)
		assert.Contains(t, metadata.Dependencies, "github.com/stretchr/testify")
	})

	t.Run("Node.js project", func(t *testing.T) {
		// Clean up previous files
		os.RemoveAll(tmpDir)
		os.MkdirAll(tmpDir, 0755)

		// Create package.json file
		packageJSON := `{
			"name": "test",
			"version": "1.0.0",
			"dependencies": {
				"express": "^4.17.1"
			}
		}`
		err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644)
		assert.NoError(t, err)

		metadata, err := getProjectMetadata(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "JavaScript/Node.js", metadata.Language)
		assert.Contains(t, metadata.Dependencies, "express")
	})
}

func TestAnalyzeProject(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "analysis-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"main.go":           "package main\n\nfunc main() {}\n",
		"internal/core.go":  "package internal\n",
		"config.yaml":       "key: value\n",
		"README.md":         "# Test Project",
		"test/main_test.go": "package test\n",
		".gitignore":        "*.tmp\n",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("project analysis", func(t *testing.T) {
		analysis := AnalyzeProject(tmpDir)
		assert.NotNil(t, analysis)

		// Check entry points
		assert.Contains(t, analysis.EntryPoints, "main.go")

		// Check core files
		assert.Contains(t, analysis.CoreFiles, "internal/core.go")

		// Check config files
		assert.Contains(t, analysis.ConfigFiles, "config.yaml")

		// Check documentation
		assert.Contains(t, analysis.Documentation, "README.md")

		// Check test files
		assert.Contains(t, analysis.TestFiles, "test/main_test.go")
	})
}
