package promptext

import (
	"errors"
	"os"
	"path/filepath"
	"testing"
)

func TestExtract_SimpleCase(t *testing.T) {
	// Create a temporary directory with test files
	tmpDir := t.TempDir()

	// Create test files
	testFile := filepath.Join(tmpDir, "test.go")
	err := os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Extract with default options
	result, err := Extract(tmpDir)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify result
	if result == nil {
		t.Fatal("Result is nil")
	}

	if result.FormattedOutput == "" {
		t.Error("FormattedOutput is empty")
	}

	if result.ProjectOutput == nil {
		t.Fatal("ProjectOutput is nil")
	}

	if len(result.ProjectOutput.Files) == 0 {
		t.Error("No files in result")
	}

	// Check that our test file is included
	found := false
	for _, file := range result.ProjectOutput.Files {
		if filepath.Base(file.Path) == "test.go" {
			found = true
			if file.Content == "" {
				t.Error("File content is empty")
			}
			break
		}
	}
	if !found {
		t.Error("test.go not found in results")
	}
}

func TestExtract_WithExtensions(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files with different extensions
	goFile := filepath.Join(tmpDir, "test.go")
	jsFile := filepath.Join(tmpDir, "test.js")
	txtFile := filepath.Join(tmpDir, "test.txt")

	os.WriteFile(goFile, []byte("package main"), 0644)
	os.WriteFile(jsFile, []byte("console.log()"), 0644)
	os.WriteFile(txtFile, []byte("text file"), 0644)

	// Extract only .go files
	result, err := Extract(tmpDir, WithExtensions(".go"))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify only .go file is included
	if len(result.ProjectOutput.Files) == 0 {
		t.Fatal("No files in result")
	}

	for _, file := range result.ProjectOutput.Files {
		if filepath.Ext(file.Path) != ".go" {
			t.Errorf("Unexpected file extension: %s", file.Path)
		}
	}
}

func TestExtract_WithExcludes(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	mainFile := filepath.Join(tmpDir, "main.go")
	testFile := filepath.Join(tmpDir, "main_test.go")

	os.WriteFile(mainFile, []byte("package main"), 0644)
	os.WriteFile(testFile, []byte("package main"), 0644)

	// Extract excluding test files
	result, err := Extract(tmpDir, WithExcludes("*_test.go"))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify test file is excluded
	for _, file := range result.ProjectOutput.Files {
		if filepath.Base(file.Path) == "main_test.go" {
			t.Error("Test file should be excluded")
		}
	}
}

func TestExtract_WithFormat(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	formats := []Format{
		FormatPTX,
		FormatMarkdown,
		FormatJSONL,
		FormatXML,
	}

	for _, format := range formats {
		t.Run(string(format), func(t *testing.T) {
			result, err := Extract(tmpDir, WithFormat(format))
			if err != nil {
				t.Fatalf("Extract failed for format %s: %v", format, err)
			}

			if result.FormattedOutput == "" {
				t.Errorf("FormattedOutput is empty for format %s", format)
			}
		})
	}
}

func TestExtract_WithTokenBudget(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple files with content
	for i := 0; i < 10; i++ {
		filename := filepath.Join(tmpDir, filepath.FromSlash("test"+string(rune('0'+i))+".go"))
		content := "package main\n\n// This is a test file with some content\n"
		os.WriteFile(filename, []byte(content), 0644)
	}

	// Extract with small token budget
	result, err := Extract(tmpDir, WithTokenBudget(100))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Should have excluded some files
	if result.ExcludedFiles == 0 {
		t.Log("Warning: Expected some files to be excluded due to token budget")
	}

	// Token count should respect budget (with some overhead)
	if result.TokenCount > 200 { // Allow some overhead for metadata
		t.Errorf("Token count %d exceeds budget significantly", result.TokenCount)
	}
}

func TestExtract_InvalidDirectory(t *testing.T) {
	_, err := Extract("/nonexistent/directory/path")
	if err == nil {
		t.Fatal("Expected error for invalid directory")
	}

	var dirErr *DirectoryError
	if !errors.As(err, &dirErr) {
		t.Errorf("Expected DirectoryError, got %T", err)
	}
}

func TestExtract_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	// Try to extract from empty directory
	_, err := Extract(tmpDir)
	if err == nil {
		t.Fatal("Expected error for empty directory")
	}

	if !errors.Is(err, ErrNoFilesMatched) {
		t.Errorf("Expected ErrNoFilesMatched, got %v", err)
	}
}

func TestExtractor_Reusability(t *testing.T) {
	tmpDir1 := t.TempDir()
	tmpDir2 := t.TempDir()

	os.WriteFile(filepath.Join(tmpDir1, "test1.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir2, "test2.go"), []byte("package main"), 0644)

	// Create extractor once
	extractor := NewExtractor(WithFormat(FormatPTX))

	// Use it for multiple directories
	result1, err := extractor.Extract(tmpDir1)
	if err != nil {
		t.Fatalf("First extract failed: %v", err)
	}

	result2, err := extractor.Extract(tmpDir2)
	if err != nil {
		t.Fatalf("Second extract failed: %v", err)
	}

	if result1 == nil || result2 == nil {
		t.Fatal("Results are nil")
	}

	if len(result1.ProjectOutput.Files) == 0 || len(result2.ProjectOutput.Files) == 0 {
		t.Error("Expected files in both results")
	}
}

func TestExtractor_BuilderPattern(t *testing.T) {
	tmpDir := t.TempDir()
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "test_test.go"), []byte("package main"), 0644)

	// Test builder pattern
	result, err := NewExtractor().
		WithExtensions(".go").
		WithExcludes("*_test.go").
		WithFormat(FormatMarkdown).
		Extract(tmpDir)

	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	// Verify test file is excluded
	for _, file := range result.ProjectOutput.Files {
		if filepath.Base(file.Path) == "test_test.go" {
			t.Error("Test file should be excluded")
		}
	}
}

func TestResult_As(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.go")
	os.WriteFile(testFile, []byte("package main"), 0644)

	// Extract with PTX format
	result, err := Extract(tmpDir, WithFormat(FormatPTX))
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Convert to different formats
	markdownOutput, err := result.As(FormatMarkdown)
	if err != nil {
		t.Fatalf("Conversion to Markdown failed: %v", err)
	}
	if markdownOutput == "" {
		t.Error("Markdown output is empty")
	}

	jsonlOutput, err := result.As(FormatJSONL)
	if err != nil {
		t.Fatalf("Conversion to JSONL failed: %v", err)
	}
	if jsonlOutput == "" {
		t.Error("JSONL output is empty")
	}
}

func TestVersion(t *testing.T) {
	if Version == "" {
		t.Error("Version should not be empty")
	}
}

func TestOptions_Combination(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte("package main\n\nfunc main() {}"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "auth.go"), []byte("package main\n\n// Authentication logic"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte("package main\n\n// Some other code"), 0644)

	// Combine multiple options
	result, err := Extract(tmpDir,
		WithExtensions(".go"),
		WithExcludes("test.go"),
		WithFormat(FormatJSONL),
		WithTokenBudget(5000),
	)

	if err != nil {
		t.Fatalf("Extract with combined options failed: %v", err)
	}

	if result == nil {
		t.Fatal("Result is nil")
	}

	// Verify test.go is excluded
	for _, file := range result.ProjectOutput.Files {
		if filepath.Base(file.Path) == "test.go" {
			t.Error("test.go should be excluded")
		}
	}
}

func TestExtract_CurrentDirectory(t *testing.T) {
	// Test with "." and "" (should use current directory)
	originalDir, _ := os.Getwd()
	tmpDir := t.TempDir()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	os.WriteFile("test.go", []byte("package main"), 0644)

	// Test with "."
	result1, err := Extract(".")
	if err != nil {
		t.Fatalf("Extract with '.' failed: %v", err)
	}
	if result1 == nil {
		t.Fatal("Result is nil")
	}

	// Test with empty string
	result2, err := Extract("")
	if err != nil {
		t.Fatalf("Extract with '' failed: %v", err)
	}
	if result2 == nil {
		t.Fatal("Result is nil")
	}
}
