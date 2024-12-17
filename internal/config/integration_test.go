package config

import (
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/1broseidon/promptext/internal/processor"
)

func TestFilterIntegration(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "filter_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"main.go":           "package main",
		"lib/helper.go":     "package lib",
		"vendor/dep.go":     "package vendor",
		"node_modules/x.js": "module.exports = {};",
		"test.txt":          "test file",
		".env":              "SECRET=123",
	}

	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Create .gitignore
	gitignoreContent := []string{
		"vendor/",
		"node_modules/",
		".env",
	}
	if err := os.WriteFile(filepath.Join(tmpDir, ".gitignore"), []byte(strings.Join(gitignoreContent, "\n")), 0644); err != nil {
		t.Fatal(err)
	}

	// Create .promptext.yml
	promptextContent := `
extensions:
  - .go
excludes:
  - test.txt
`
	if err := os.WriteFile(filepath.Join(tmpDir, ".promptext.yml"), []byte(promptextContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Process directory
	config := processor.Config{
		DirPath: tmpDir,
	}

	result, err := processor.ProcessDirectory(config, false)
	if err != nil {
		t.Fatal(err)
	}

	// Verify results
	expectedFiles := []string{
		"main.go",
		"lib/helper.go",
	}

	foundFiles := make([]string, 0)
	for _, file := range result.ProjectOutput.Files {
		foundFiles = append(foundFiles, file.Path)
	}

	// Sort both slices for comparison
	sort.Strings(expectedFiles)
	sort.Strings(foundFiles)

	if !reflect.DeepEqual(expectedFiles, foundFiles) {
		t.Errorf("Filtered files mismatch:\nwant: %v\ngot:  %v", expectedFiles, foundFiles)
	}

	// Verify excluded files are not present
	excludedFiles := []string{
		"vendor/dep.go",
		"node_modules/x.js",
		".env",
		"test.txt",
	}

	for _, excluded := range excludedFiles {
		for _, found := range foundFiles {
			if found == excluded {
				t.Errorf("Found excluded file in output: %s", excluded)
			}
		}
	}
}
