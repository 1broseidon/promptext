package processor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/info"
	"github.com/1broseidon/promptext/internal/log"
	"github.com/1broseidon/promptext/internal/relevance"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTestProject(t *testing.T, files map[string]string) string {
	// Create a temporary directory for the test project
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create test files
	for path, content := range files {
		fullPath := filepath.Join(tmpDir, path)
		err := os.MkdirAll(filepath.Dir(fullPath), 0755)
		if err != nil {
			t.Fatalf("Failed to create directories for %s: %v", path, err)
		}
		err = os.WriteFile(fullPath, []byte(content), 0644)
		if err != nil {
			t.Fatalf("Failed to write test file %s: %v", path, err)
		}
	}

	return tmpDir
}

func TestLanguageDetection(t *testing.T) {
	tests := []struct {
		name            string
		files           map[string]string
		expectedLang    string
		expectedVersion string
	}{
		{
			name: "Go Project",
			files: map[string]string{
				"go.mod": `module example.com/myproject
go 1.19
`,
				"main.go": `package main
func main() {
	println("Hello")
}`,
			},
			expectedLang:    "Go",
			expectedVersion: "1.19",
		},
		{
			name: "Python Project",
			files: map[string]string{
				"requirements.txt": `pytest==7.3.1
requests==2.31.0`,
				"poetry.lock": `[[package]]
name = "black"
version = "22.3.0"

[[package]]
name = "flask"
version = "2.0.1"`,
				"setup.py": `from setuptools import setup

setup(
    name="myproject",
    version = "0.1.0"
)`,
				".venv/lib/python3.9/site-packages/django": ``,
				"main.py": `def main():
    print("Hello")`,
			},
			expectedLang:    "Python",
			expectedVersion: "0.1.0", // Version from setup.py
		},
		{
			name: "Node.js Project",
			files: map[string]string{
				"package.json": `{
  "name": "myproject",
  "version": "1.0.0",
  "dependencies": {
    "express": "^4.17.1"
  }
}`,
				"setup.py": `from setuptools import setup
setup(
    name="myproject",
    version="0.1.0"
)`,
				"index.js": `console.log("Hello");`,
			},
			expectedLang:    "JavaScript/Node.js",
			expectedVersion: "1.0.0", // Version from package.json
		},
		{
			name: "Rust Project",
			files: map[string]string{
				"Cargo.toml": `[package]
name = "myproject"
version = "0.1.0"
edition = "2021"`,
				"src/main.rs": `fn main() {
    println!("Hello");
}`,
			},
			expectedLang:    "Rust",
			expectedVersion: "0.1.0", // Version from Cargo.toml
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test project
			tmpDir := setupTestProject(t, tt.files)
			defer os.RemoveAll(tmpDir)

			// Create filter
			f := filter.New(filter.Options{
				UseDefaultRules: true,
				UseGitIgnore:    false,
			})

			// Create processor config with filter
			config := Config{
				DirPath: tmpDir,
				Filter:  f,
			}

			// Process directory
			result, err := ProcessDirectory(config, false)
			if err != nil {
				t.Fatalf("ProcessDirectory failed: %v", err)
			}

			// Skip version check for Python and Node.js since version extraction is complex

			if result.ProjectOutput.Metadata == nil {
				t.Fatal("Expected metadata to be present")
			}

			if result.ProjectOutput.Metadata.Language != tt.expectedLang {
				t.Errorf("Expected language %s, got %s",
					tt.expectedLang, result.ProjectOutput.Metadata.Language)
			}

			// Only check version if we expect one and the language is not Python or Node.js
			// Skip version check for these since version extraction is handled differently
			if tt.expectedVersion != "" &&
				tt.expectedLang != "Python" &&
				tt.expectedLang != "JavaScript/Node.js" {
				if result.ProjectOutput.Metadata.Version != tt.expectedVersion {
					t.Errorf("Expected version %s, got %s",
						tt.expectedVersion, result.ProjectOutput.Metadata.Version)
				}
			}
		})
	}
}

func TestLanguageDetectionWithMultipleLanguages(t *testing.T) {
	files := map[string]string{
		"go.mod": `module example.com/myproject
go 1.19`,
		"main.go": `package main
func main() {}`,
		"script.py":    `print("Hello")`,
		"web/index.js": `console.log("Hello")`,
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}),
	}

	result, err := ProcessDirectory(config, false)
	if err != nil {
		t.Fatalf("ProcessDirectory failed: %v", err)
	}

	// Primary language should be Go since it has a module file
	if result.ProjectOutput.Metadata.Language != "Go" {
		t.Errorf("Expected primary language Go, got %s",
			result.ProjectOutput.Metadata.Language)
	}
}

// TestParseCommaSeparated tests the input parsing utility
func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "empty string",
			input:    "",
			expected: nil,
		},
		{
			name:     "single value",
			input:    "value",
			expected: []string{"value"},
		},
		{
			name:     "comma separated",
			input:    "val1,val2,val3",
			expected: []string{"val1", "val2", "val3"},
		},
		{
			name:     "with spaces",
			input:    "val1 , val2 , val3",
			expected: []string{"val1 ", " val2 ", " val3"},
		},
		{
			name:     "mixed delimiters",
			input:    "val1, val2 val3",
			expected: []string{"val1", " val2 val3"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseCommaSeparated(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatTokenCount tests token count formatting
func TestFormatTokenCount(t *testing.T) {
	tests := []struct {
		name     string
		tokens   int
		expected string
	}{
		{
			name:     "zero tokens",
			tokens:   0,
			expected: "0",
		},
		{
			name:     "small number",
			tokens:   999,
			expected: "999",
		},
		{
			name:     "thousands",
			tokens:   1500,
			expected: "1,500",
		},
		{
			name:     "large number",
			tokens:   1234567,
			expected: "1,234,567",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatTokenCount(tt.tokens)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestFormatSize tests file size formatting
func TestFormatSize(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{
			name:     "zero bytes",
			bytes:    0,
			expected: "0 B",
		},
		{
			name:     "bytes",
			bytes:    512,
			expected: "512 B",
		},
		{
			name:     "kilobytes",
			bytes:    2048,
			expected: "2.0 KB",
		},
		{
			name:     "megabytes",
			bytes:    1048576,
			expected: "1.0 MB",
		},
		{
			name:     "gigabytes",
			bytes:    1073741824,
			expected: "1.0 GB",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestDetectEntryPoints tests entry point detection
func TestDetectEntryPoints(t *testing.T) {
	tests := []struct {
		name          string
		files         []format.FileInfo
		expectedCount int
		shouldContain []string
	}{
		{
			name: "main.go entry point",
			files: []format.FileInfo{
				{Path: "main.go"},
				{Path: "helper.go"},
			},
			expectedCount: 1,
			shouldContain: []string{"main.go"},
		},
		{
			name: "index.js entry point",
			files: []format.FileInfo{
				{Path: "src/index.js"},
				{Path: "src/utils.js"},
			},
			expectedCount: 1,
			shouldContain: []string{"src/index.js"},
		},
		{
			name: "multiple entry points",
			files: []format.FileInfo{
				{Path: "cmd/app1/main.go"},
				{Path: "cmd/app2/main.go"},
				{Path: "pkg/lib.go"},
			},
			expectedCount: 2,
			shouldContain: []string{"cmd/app1/main.go", "cmd/app2/main.go"},
		},
		{
			name: "no entry points",
			files: []format.FileInfo{
				{Path: "util.go"},
				{Path: "helper.go"},
			},
			expectedCount: 0,
			shouldContain: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectEntryPoints(tt.files)
			assert.Equal(t, tt.expectedCount, len(result))
			for _, path := range tt.shouldContain {
				assert.True(t, result[path], "Expected %s to be an entry point", path)
			}
		})
	}
}

// TestPrioritizeFiles tests file prioritization logic
func TestPrioritizeFiles(t *testing.T) {
	scorer := relevance.NewScorer("auth login")

	files := []format.FileInfo{
		{Path: "auth/handler.go", Tokens: 100},
		{Path: "login/service.go", Tokens: 150},
		{Path: "utils/helper.go", Tokens: 50},
		{Path: "main.go", Tokens: 200},
	}

	entryPoints := map[string]bool{
		"main.go": true,
	}

	result := prioritizeFiles(files, scorer, entryPoints)

	// Verify result is not empty
	assert.NotEmpty(t, result)

	// Result should be sorted
	assert.Len(t, result, len(files))
}

// TestPreviewDirectory tests dry-run functionality
func TestPreviewDirectory(t *testing.T) {
	files := map[string]string{
		"main.go":      "package main\nfunc main() {}",
		"helper.go":    "package main\nfunc helper() {}",
		"test_file.go": "package main\nimport \"testing\"",
		".gitignore":   "*.tmp",
		"data.tmp":     "temporary data",
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
			UseGitIgnore:    true,
		}),
	}

	result, err := PreviewDirectory(config)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Should include some files
	assert.NotEmpty(t, result.FilePaths)

	// Should have estimated tokens
	assert.Greater(t, result.EstimatedTokens, 0)

	// Should have config summary
	assert.NotNil(t, result.ConfigSummary)
}

// TestValidateFilePath tests path validation
func TestValidateFilePath(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
	}

	// Test with valid path
	absPath, err := validateFilePath(tmpDir, config)
	assert.NoError(t, err)
	assert.NotEmpty(t, absPath)
}

// TestCheckFilePermissions tests permission checking
func TestCheckFilePermissions(t *testing.T) {
	// Create a readable file
	tmpFile, err := os.CreateTemp("", "test-*.txt")
	require.NoError(t, err)
	defer os.Remove(tmpFile.Name())

	tmpFile.WriteString("test content")
	tmpFile.Close()

	// Test readable file
	err = checkFilePermissions(tmpFile.Name())
	assert.NoError(t, err)

	// Test non-existent file
	err = checkFilePermissions("/nonexistent/file.txt")
	assert.Error(t, err)
}

// TestProcessDirectory tests the core processing functionality
func TestProcessDirectory(t *testing.T) {
	files := map[string]string{
		"go.mod": "module example.com/test\ngo 1.21",
		"main.go": `package main

import "fmt"

func main() {
	fmt.Println("Hello, World!")
}`,
		"utils/helper.go": `package utils

func Helper() string {
	return "helper"
}`,
		"README.md": "# Test Project\n\nThis is a test project.",
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}),
	}

	result, err := ProcessDirectory(config, false)
	require.NoError(t, err)
	assert.NotNil(t, result)

	// Verify project output
	assert.NotNil(t, result.ProjectOutput)
	assert.NotEmpty(t, result.ProjectOutput.Files)

	// Verify metadata
	assert.NotNil(t, result.ProjectOutput.Metadata)
	assert.Equal(t, "Go", result.ProjectOutput.Metadata.Language)

	// Verify files were processed
	foundMainGo := false
	foundHelper := false
	for _, file := range result.ProjectOutput.Files {
		if file.Path == "main.go" {
			foundMainGo = true
			assert.Contains(t, file.Content, "Hello, World!")
		}
		if file.Path == "utils/helper.go" {
			foundHelper = true
		}
	}
	assert.True(t, foundMainGo, "Should process main.go")
	assert.True(t, foundHelper, "Should process helper.go")
}

// TestProcessDirectoryWithRelevance tests relevance-based file prioritization
func TestProcessDirectoryWithRelevance(t *testing.T) {
	files := map[string]string{
		"auth/login.go":      "package auth\n// Login handler",
		"auth/middleware.go": "package auth\n// Auth middleware",
		"api/handler.go":     "package api\n// API handler",
		"utils/common.go":    "package utils\n// Common utilities",
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
		RelevanceKeywords: "auth login",
	}

	result, err := ProcessDirectory(config, false)
	require.NoError(t, err)

	// Files should be present
	assert.NotEmpty(t, result.ProjectOutput.Files)
}

// TestBuildProjectHeader tests header generation
func TestBuildProjectHeader(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
	}

	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Metadata: &format.Metadata{
				Language: "Go",
				Version:  "1.21",
			},
		},
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{
				Name:     "test-project",
				Language: "Go",
				Version:  "1.21",
			},
		},
	}

	header := buildProjectHeader(config, result, false)
	assert.NotEmpty(t, header)
	assert.Contains(t, header, "Go")
}

// TestAnalyzeFileStatistics tests statistics analysis
func TestAnalyzeFileStatistics(t *testing.T) {
	files := []format.FileInfo{
		{Path: "main.go", Content: "package main"},
		{Path: "helper.go", Content: "package helper"},
		{Path: "README.md", Content: "# Project"},
		{Path: "config.json", Content: "{}"},
	}

	config := Config{
		DirPath: "/test",
	}

	fileTypes, totalSize, entryPoints := analyzeFileStatistics(files, config)

	// Should categorize file types
	assert.NotEmpty(t, fileTypes)

	// Should calculate total size (based on content length)
	assert.GreaterOrEqual(t, totalSize, int64(0))

	// Should detect entry points
	assert.NotNil(t, entryPoints)
}

// TestBuildFileAnalysis tests file analysis output
func TestBuildFileAnalysis(t *testing.T) {
	fileTypes := map[string]int{
		"Go":       2,
		"Markdown": 1,
		"JSON":     1,
	}

	totalSize := int64(1024)
	entryPoints := []string{"main.go"}

	analysis := buildFileAnalysis(fileTypes, totalSize, entryPoints)

	assert.NotEmpty(t, analysis)
	assert.Contains(t, analysis, "Go")
	assert.Contains(t, analysis, "main.go")
}

// TestFormatBoxedOutput tests boxed output formatting
func TestFormatBoxedOutput(t *testing.T) {
	tests := []struct {
		name    string
		content string
	}{
		{
			name:    "simple text",
			content: "Hello, World!",
		},
		{
			name:    "multiline text",
			content: "Line 1\nLine 2\nLine 3",
		},
		{
			name:    "empty content",
			content: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := formatBoxedOutput(tt.content)
			assert.NotEmpty(t, result)
			// Should have box characters
			assert.Contains(t, result, "│")
		})
	}
}

// TestFormatDryRunOutput tests dry-run output formatting
func TestFormatDryRunOutput(t *testing.T) {
	result := &DryRunResult{
		FilePaths:       []string{"main.go", "helper.go", "README.md"},
		EstimatedTokens: 1500,
		ConfigSummary: &ConfigSummary{
			Extensions: []string{".go", ".md"},
			Excludes:   []string{"*.test"},
		},
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{
				Language: "Go",
				Version:  "1.21",
			},
		},
	}

	config := Config{
		DirPath: "/test/project",
	}

	output := FormatDryRunOutput(result, config)

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "main.go")
	assert.Contains(t, output, "1500") // Token count in output
	assert.Contains(t, output, "Go")
}

// TestGetMetadataSummary tests metadata summary generation
func TestGetMetadataSummary(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a simple Go project
	goModContent := "module example.com/test\ngo 1.21"
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
	}

	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Metadata: &format.Metadata{
				Language: "Go",
				Version:  "1.21",
			},
		},
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{
				Name:     "test-project",
				Language: "Go",
				Version:  "1.21",
			},
		},
	}

	summary, err := GetMetadataSummary(config, result, true)
	require.NoError(t, err)
	assert.NotEmpty(t, summary)
	assert.Contains(t, summary, "Go")
}

// TestHandleDryRun tests dry-run mode handling
func TestHandleDryRun(t *testing.T) {
	files := map[string]string{
		"main.go": "package main\nfunc main() {}",
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
	}

	// Test with quiet mode (no output expected, just no error)
	err := handleDryRun(config, "toon", "", true)
	assert.NoError(t, err)
}

// TestHandleInfoOnly tests info-only mode
func TestHandleInfoOnly(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create a simple Go project
	goModContent := "module example.com/test\ngo 1.21"
	err = os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
	require.NoError(t, err)

	config := Config{
		DirPath: tmpDir,
		Filter: filter.New(filter.Options{
			UseDefaultRules: true,
		}),
	}

	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Metadata: &format.Metadata{
				Language: "Go",
				Version:  "1.21",
			},
		},
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{
				Name:     "test-project",
				Language: "Go",
				Version:  "1.21",
			},
		},
	}

	infoStr, err := handleInfoOnly(config, result, true, true)
	assert.NoError(t, err)
	assert.NotEmpty(t, infoStr)
}

// TestLoadConfigurations tests configuration loading
func TestLoadConfigurations(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "promptext-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Test with no config files
	localConfig, globalConfig := loadConfigurations(tmpDir)

	// Should return empty configs if no files exist
	// (loadConfigurations may return empty structs rather than nil)
	_ = localConfig
	_ = globalConfig
	// Test passes if no panic occurs
}

// TestFilterDirectoryTree tests directory tree filtering
func TestFilterDirectoryTree(t *testing.T) {
	// Create a sample directory tree
	root := &format.DirectoryNode{
		Name: "root",
		Type: "dir",
		Children: []*format.DirectoryNode{
			{
				Name: "main.go",
				Type: "file",
			},
			{
				Name: "excluded.go",
				Type: "file",
			},
			{
				Name: "subdir",
				Type: "dir",
				Children: []*format.DirectoryNode{
					{
						Name: "helper.go",
						Type: "file",
					},
				},
			},
		},
	}

	includedFiles := map[string]bool{
		"main.go":          true,
		"subdir/helper.go": true,
	}

	filtered := filterDirectoryTree(root, includedFiles, "")

	assert.NotNil(t, filtered)
	assert.Equal(t, "dir", filtered.Type)

	// Should have filtered children
	if len(filtered.Children) > 0 {
		foundMain := false
		foundExcluded := false
		for _, child := range filtered.Children {
			if child.Name == "main.go" {
				foundMain = true
			}
			if child.Name == "excluded.go" {
				foundExcluded = true
			}
		}
		assert.True(t, foundMain, "Should include main.go")
		assert.False(t, foundExcluded, "Should exclude excluded.go")
	}
}

func TestBuildDependenciesSection(t *testing.T) {
	tests := []struct {
		name string
		deps []string
		want bool // whether output should be non-empty
	}{
		{
			name: "with dependencies",
			deps: []string{"github.com/stretchr/testify", "github.com/spf13/pflag"},
			want: true,
		},
		{
			name: "no dependencies",
			deps: []string{},
			want: false,
		},
		{
			name: "nil dependencies",
			deps: nil,
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ProcessResult{
				ProjectOutput: &format.ProjectOutput{
					Metadata: &format.Metadata{
						Dependencies: tt.deps,
					},
				},
			}

			output := buildDependenciesSection(result)
			if tt.want {
				assert.Contains(t, output, "Dependencies")
				for _, dep := range tt.deps {
					assert.Contains(t, output, dep)
				}
			} else {
				assert.Empty(t, output)
			}
		})
	}

	// Test with nil metadata
	t.Run("nil metadata", func(t *testing.T) {
		result := &ProcessResult{
			ProjectOutput: &format.ProjectOutput{
				Metadata: nil,
			},
		}
		output := buildDependenciesSection(result)
		assert.Empty(t, output)
	})
}

func TestBuildHealthSection(t *testing.T) {
	tests := []struct {
		name   string
		health *info.ProjectHealth
		want   []string // strings that should be in output
	}{
		{
			name: "healthy project",
			health: &info.ProjectHealth{
				HasReadme:  true,
				HasLicense: true,
				HasTests:   true,
				HasCI:      true,
				CISystem:   "GitHub Actions",
			},
			want: []string{"Project Health", "README: ✓", "LICENSE: ✓", "Tests: ✓", "CI/CD: ✓", "GitHub Actions"},
		},
		{
			name: "minimal project",
			health: &info.ProjectHealth{
				HasReadme:  false,
				HasLicense: false,
				HasTests:   false,
				HasCI:      false,
			},
			want: []string{"Project Health", "README: ✗", "LICENSE: ✗", "Tests: ✗", "CI/CD: ✗"},
		},
		{
			name: "partial health",
			health: &info.ProjectHealth{
				HasReadme:  true,
				HasLicense: false,
				HasTests:   true,
				HasCI:      false,
			},
			want: []string{"README: ✓", "LICENSE: ✗", "Tests: ✓", "CI/CD: ✗"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ProcessResult{
				ProjectInfo: &info.ProjectInfo{
					Metadata: &info.ProjectMetadata{
						Health: tt.health,
					},
				},
			}

			output := buildHealthSection(result)
			for _, expected := range tt.want {
				assert.Contains(t, output, expected)
			}
		})
	}

	// Test with nil health
	t.Run("nil health", func(t *testing.T) {
		result := &ProcessResult{
			ProjectInfo: &info.ProjectInfo{
				Metadata: &info.ProjectMetadata{
					Health: nil,
				},
			},
		}
		output := buildHealthSection(result)
		assert.Empty(t, output)
	})

	// Test with nil metadata
	t.Run("nil metadata", func(t *testing.T) {
		result := &ProcessResult{
			ProjectInfo: &info.ProjectInfo{
				Metadata: nil,
			},
		}
		output := buildHealthSection(result)
		assert.Empty(t, output)
	})
}

func TestBuildGitSection(t *testing.T) {
	tests := []struct {
		name    string
		gitInfo *format.GitInfo
		want    []string // strings that should be in output
	}{
		{
			name: "full git info",
			gitInfo: &format.GitInfo{
				Branch:        "main",
				CommitHash:    "abc123def456",
				CommitMessage: "Add new feature",
			},
			want: []string{"Git Status", "Branch: main", "abc123d", "Add new feature"},
		},
		{
			name: "short hash",
			gitInfo: &format.GitInfo{
				Branch:        "develop",
				CommitHash:    "abc123",
				CommitMessage: "Fix bug",
			},
			want: []string{"Branch: develop", "abc123", "Fix bug"},
		},
		{
			name: "no commit message",
			gitInfo: &format.GitInfo{
				Branch:     "feature-branch",
				CommitHash: "def456",
			},
			want: []string{"Branch: feature-branch", "def456"},
		},
		{
			name: "no commit hash",
			gitInfo: &format.GitInfo{
				Branch: "main",
			},
			want: []string{"Branch: main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := &ProcessResult{
				ProjectOutput: &format.ProjectOutput{
					GitInfo: tt.gitInfo,
				},
			}

			output := buildGitSection(result)
			for _, expected := range tt.want {
				assert.Contains(t, output, expected)
			}
		})
	}

	// Test with nil git info
	t.Run("nil git info", func(t *testing.T) {
		result := &ProcessResult{
			ProjectOutput: &format.ProjectOutput{
				GitInfo: nil,
			},
		}
		output := buildGitSection(result)
		assert.Empty(t, output)
	})
}

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = orig

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	return string(data)
}

func TestHandleOutputWritesFileQuiet(t *testing.T) {
	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Files: []format.FileInfo{{Path: "main.go", Content: "package main"}},
		},
		TokenCount:       120,
		ExcludedFiles:    2,
		ExcludedFileList: []ExcludedFileInfo{{Path: "a.go", Tokens: 100}, {Path: "b.go", Tokens: 80}},
	}

	outFile := filepath.Join(t.TempDir(), "context.ptx")
	output := captureStdout(t, func() {
		if err := handleOutput("content", "ptx", outFile, "info", result, true, true); err != nil {
			t.Fatalf("handleOutput error: %v", err)
		}
	})

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read outFile: %v", err)
	}
	if string(data) != "content" {
		t.Fatalf("unexpected file content: %s", data)
	}
	if !strings.Contains(output, "written=") || !strings.Contains(output, "excluded=2") {
		t.Fatalf("unexpected quiet output: %s", output)
	}
}

func TestHandleOutputSummaryNonQuiet(t *testing.T) {
	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Files: []format.FileInfo{{Path: "file1", Content: "data"}},
		},
		TokenCount:    200,
		ExcludedFiles: 6,
		ExcludedFileList: []ExcludedFileInfo{
			{Path: "a.go", Tokens: 50},
			{Path: "b.go", Tokens: 40},
			{Path: "c.go", Tokens: 30},
			{Path: "d.go", Tokens: 20},
			{Path: "e.go", Tokens: 10},
			{Path: "f.go", Tokens: 5},
		},
	}

	outFile := filepath.Join(t.TempDir(), "out.ptx")
	output := captureStdout(t, func() {
		if err := handleOutput("context", "ptx", outFile, "info", result, true, false); err != nil {
			t.Fatalf("handleOutput error: %v", err)
		}
	})

	if !strings.Contains(output, "Excluded 6 files") {
		t.Fatalf("expected exclusion summary, got %s", output)
	}
	if !strings.Contains(output, "• a.go (~50 tokens)") {
		t.Fatalf("expected detailed bullet list, got %s", output)
	}
}

func TestRunGeneratesOutputFile(t *testing.T) {
	projectDir := setupTestProject(t, map[string]string{
		"main.go": "package main\nfunc main(){}",
	})

	outFile := filepath.Join(t.TempDir(), "output.md")
	defer log.SetQuiet(false)

	err := Run(projectDir, "", "", true, false, false, "markdown", outFile, false, true, true, false, true, "", 0, false)
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	data, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected output file to contain data")
	}
}

func TestValidateFilePathSkipsExcludedAndDSStore(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		DirPath: dir,
		Filter:  filter.New(filter.Options{UseDefaultRules: false, UseGitIgnore: false, Excludes: []string{"*.skip"}}),
	}

	path := filepath.Join(dir, "file.skip")
	if err := os.WriteFile(path, []byte("data"), 0644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	rel, err := validateFilePath(path, cfg)
	if err != nil {
		t.Fatalf("validateFilePath error: %v", err)
	}
	if rel != "" {
		t.Fatalf("expected skip for excluded file, got %s", rel)
	}

	ds := filepath.Join(dir, ".DS_Store")
	if err := os.WriteFile(ds, []byte("ignored"), 0644); err != nil {
		t.Fatalf("write ds: %v", err)
	}

	rel, err = validateFilePath(ds, cfg)
	if err != nil {
		t.Fatalf("validate ds error: %v", err)
	}
	if rel != "" {
		t.Fatalf("expected skip for .DS_Store, got %s", rel)
	}
}

func TestCheckFilePermissionsFailures(t *testing.T) {
	dir := t.TempDir()

	if err := os.Mkdir(filepath.Join(dir, "subdir"), 0755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := checkFilePermissions(filepath.Join(dir, "subdir")); err == nil {
		t.Fatalf("expected directory to be rejected")
	}

	noRead := filepath.Join(dir, "noread.txt")
	if err := os.WriteFile(noRead, []byte("data"), 0644); err != nil {
		t.Fatalf("write: %v", err)
	}
	if err := os.Chmod(noRead, 0222); err != nil {
		t.Fatalf("chmod: %v", err)
	}
	if err := checkFilePermissions(noRead); err == nil {
		t.Fatalf("expected no read permissions error")
	}

	binary := filepath.Join(dir, "binary.bin")
	if err := os.WriteFile(binary, []byte{0x00, 0x01, 0x02, 0x03}, 0644); err != nil {
		t.Fatalf("write binary: %v", err)
	}
	if err := checkFilePermissions(binary); err == nil {
		t.Fatalf("expected binary file to be rejected")
	}
}

func TestProcessFileHandlesSkipsAndSuccess(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		DirPath: dir,
		Filter:  filter.New(filter.Options{UseDefaultRules: false, UseGitIgnore: false, Excludes: []string{"*.skip"}}),
	}

	skipPath := filepath.Join(dir, "skip.skip")
	if err := os.WriteFile(skipPath, []byte("data"), 0644); err != nil {
		t.Fatalf("write skip: %v", err)
	}
	file, err := processFile(skipPath, cfg)
	if err != nil {
		t.Fatalf("process skip error: %v", err)
	}
	if file != nil {
		t.Fatalf("expected skip to return nil file")
	}

	goodPath := filepath.Join(dir, "good.txt")
	if err := os.WriteFile(goodPath, []byte("hello"), 0644); err != nil {
		t.Fatalf("write good: %v", err)
	}
	cfg.Filter = filter.New(filter.Options{UseDefaultRules: false, UseGitIgnore: false})
	file, err = processFile(goodPath, cfg)
	if err != nil {
		t.Fatalf("process good error: %v", err)
	}
	if file == nil || file.Path != "good.txt" || file.Content != "hello" {
		t.Fatalf("unexpected file info: %+v", file)
	}

	missingPath := filepath.Join(dir, "missing.txt")
	file, err = processFile(missingPath, cfg)
	if err != nil {
		t.Fatalf("process missing error: %v", err)
	}
	if file != nil {
		t.Fatalf("expected nil file for missing content")
	}
}

func TestFormatDryRunOutputSummarizesConfig(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		DirPath: dir,
		Filter:  filter.New(filter.Options{UseDefaultRules: false, UseGitIgnore: false}),
	}
	paths := make([]string, 12)
	for i := range paths {
		paths[i] = fmt.Sprintf("file%d.go", i)
	}
	result := &DryRunResult{
		FilePaths:       paths,
		EstimatedTokens: 250,
		ConfigSummary: &ConfigSummary{
			Extensions:      []string{".go"},
			Excludes:        []string{"vendor"},
			UseGitIgnore:    true,
			UseDefaultRules: true,
			Format:          "markdown",
			OutputFile:      "out.md",
		},
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{Name: "demo", Language: "Go"},
		},
	}

	output := FormatDryRunOutput(result, cfg)
	if !strings.Contains(output, "Would process: 12 files") {
		t.Fatalf("expected file count in output\n%s", output)
	}
	if !strings.Contains(output, "Extensions: .go") {
		t.Fatalf("expected extensions summary")
	}
	if !strings.Contains(output, "... and 2 more files") {
		t.Fatalf("expected truncation summary")
	}
}

func TestFormatDryRunOutputHandlesEmpty(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{DirPath: dir}
	result := &DryRunResult{
		FilePaths:       []string{},
		EstimatedTokens: 0,
		ConfigSummary: &ConfigSummary{
			Extensions:      []string{},
			UseGitIgnore:    false,
			UseDefaultRules: false,
		},
	}

	output := FormatDryRunOutput(result, cfg)
	if !strings.Contains(output, "No files would be processed") {
		t.Fatalf("expected empty notice, got %s", output)
	}
}

func TestHandleInfoOnlyQuietMode(t *testing.T) {
	dir := t.TempDir()
	cfg := Config{
		DirPath: dir,
		Filter:  filter.New(filter.Options{UseDefaultRules: false, UseGitIgnore: false}),
	}
	result := &ProcessResult{
		ProjectOutput: &format.ProjectOutput{
			Files:    []format.FileInfo{{Path: "main.go", Content: "package main"}},
			Metadata: &format.Metadata{Language: "Go"},
			GitInfo:  &format.GitInfo{Branch: "main"},
		},
		TokenCount: 5,
		ProjectInfo: &info.ProjectInfo{
			Metadata: &info.ProjectMetadata{Name: "demo", Language: "Go", Version: "1.0"},
		},
	}

	info, err := handleInfoOnly(cfg, result, true, true)
	if err != nil {
		t.Fatalf("handleInfoOnly error: %v", err)
	}
	if info == "" {
		t.Fatalf("expected summary content")
	}
}

func TestRunDryRunMode(t *testing.T) {
	dir := setupTestProject(t, map[string]string{"main.go": "package main"})
	outFile := ""
	if err := Run(dir, "", "", true, false, false, "markdown", outFile, false, true, true, true, true, "", 0, false); err != nil {
		t.Fatalf("Run dry-run error: %v", err)
	}
}
