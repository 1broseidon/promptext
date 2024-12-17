package processor

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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
				"setup.py": `from setuptools import setup
setup(name="myproject", version="0.1.0")`,
				"main.py": `def main():
    print("Hello")`,
			},
			expectedLang:    "Python",
			expectedVersion: "0.1.0",
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
				"index.js": `console.log("Hello");`,
			},
			expectedLang:    "JavaScript",
			expectedVersion: "1.0.0",
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
			expectedVersion: "0.1.0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test project
			tmpDir := setupTestProject(t, tt.files)
			defer os.RemoveAll(tmpDir)

			// Create processor config
			config := Config{
				DirPath: tmpDir,
			}

			// Process directory
			result, err := ProcessDirectory(config, false)
			if err != nil {
				t.Fatalf("ProcessDirectory failed: %v", err)
			}

			if result.ProjectOutput.Metadata == nil {
				t.Fatal("Expected metadata to be present")
			}

			if result.ProjectOutput.Metadata.Language != tt.expectedLang {
				t.Errorf("Expected language %s, got %s", 
					tt.expectedLang, result.ProjectOutput.Metadata.Language)
			}

			if result.ProjectOutput.Metadata.Version != tt.expectedVersion {
				t.Errorf("Expected version %s, got %s", 
					tt.expectedVersion, result.ProjectOutput.Metadata.Version)
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
		"script.py": `print("Hello")`,
		"web/index.js": `console.log("Hello")`,
	}

	tmpDir := setupTestProject(t, files)
	defer os.RemoveAll(tmpDir)

	config := Config{
		DirPath: tmpDir,
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
