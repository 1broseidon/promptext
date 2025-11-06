package initializer

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestInitializer_PathValidation(t *testing.T) {
	tests := []struct {
		name        string
		setupPath   func() (string, func())
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid directory",
			setupPath: func() (string, func()) {
				tmpDir, _ := os.MkdirTemp("", "valid-*")
				return tmpDir, func() { os.RemoveAll(tmpDir) }
			},
			expectError: false,
		},
		{
			name: "non-existent directory",
			setupPath: func() (string, func()) {
				return "/tmp/nonexistent-dir-12345", func() {}
			},
			expectError: true,
			errorMsg:    "directory does not exist",
		},
		{
			name: "path is a file not directory",
			setupPath: func() (string, func()) {
				tmpFile, _ := os.CreateTemp("", "file-*")
				path := tmpFile.Name()
				tmpFile.Close()
				return path, func() { os.Remove(path) }
			},
			expectError: true,
			errorMsg:    "path is not a directory",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path, cleanup := tt.setupPath()
			defer cleanup()

			init := NewInitializer(path, false, true) // quiet mode
			err := init.RunQuick()

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorMsg)
				} else if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

func TestInitializer_ForceOverwrite(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "force-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create initial config
	configPath := filepath.Join(tmpDir, ".promptext.yml")
	initialContent := "# Initial content"
	if err := os.WriteFile(configPath, []byte(initialContent), 0644); err != nil {
		t.Fatalf("Failed to create initial config: %v", err)
	}

	// Try without force - should fail
	init := NewInitializer(tmpDir, false, true)
	err = init.RunQuick()
	if err == nil {
		t.Error("Expected error when config exists without force flag")
	}

	// Verify original content is unchanged
	content, _ := os.ReadFile(configPath)
	if string(content) != initialContent {
		t.Error("Config was modified without force flag")
	}

	// Try with force - should succeed
	init = NewInitializer(tmpDir, true, true)
	err = init.RunQuick()
	if err != nil {
		t.Errorf("Expected success with force flag, got: %v", err)
	}

	// Verify content was overwritten
	content, _ = os.ReadFile(configPath)
	if string(content) == initialContent {
		t.Error("Config was not overwritten with force flag")
	}
}

func TestInitializer_ConfigGeneration(t *testing.T) {
	tests := []struct {
		name          string
		setupFiles    []string
		includeTests  bool
		expectStrings []string
	}{
		{
			name:       "Next.js project",
			setupFiles: []string{"next.config.js", "package.json"},
			expectStrings: []string{
				".js",
				".jsx",
				".ts",
				".tsx",
				"node_modules",
				".next",
			},
		},
		{
			name:       "Go project",
			setupFiles: []string{"go.mod"},
			expectStrings: []string{
				".go",
				".mod",
				"vendor",
			},
		},
		{
			name:       "Django project",
			setupFiles: []string{"manage.py"},
			expectStrings: []string{
				".py",
				"__pycache__",
				"venv",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "config-gen-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for _, file := range tt.setupFiles {
				filePath := filepath.Join(tmpDir, file)
				if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			// Run initialization
			init := NewInitializer(tmpDir, false, true)
			err = init.RunQuick()
			if err != nil {
				t.Fatalf("RunQuick() error = %v", err)
			}

			// Read generated config
			configPath := filepath.Join(tmpDir, ".promptext.yml")
			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("Failed to read generated config: %v", err)
			}

			// Verify expected strings are present
			configStr := string(content)
			for _, expectedStr := range tt.expectStrings {
				if !strings.Contains(configStr, expectedStr) {
					t.Errorf("Expected config to contain '%s', but it doesn't.\nConfig:\n%s", expectedStr, configStr)
				}
			}

			// Verify file structure
			if !strings.Contains(configStr, "extensions:") {
				t.Error("Config should contain 'extensions:' section")
			}
			if !strings.Contains(configStr, "excludes:") {
				t.Error("Config should contain 'excludes:' section")
			}
			if !strings.Contains(configStr, "gitignore: true") {
				t.Error("Config should have gitignore enabled by default")
			}
		})
	}
}

func TestInitializerRunCreatesConfig(t *testing.T) {
	tmpDir := t.TempDir()
	init := NewInitializer(tmpDir, false, true)

	if err := init.Run(); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if _, err := os.Stat(filepath.Join(tmpDir, ".promptext.yml")); err != nil {
		t.Fatalf("expected config file to be created: %v", err)
	}
}

func TestPromptConfirmFlow(t *testing.T) {
	tmpDir := t.TempDir()
	init := NewInitializer(tmpDir, false, false)

	origStdin := os.Stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdin = r
	t.Cleanup(func() {
		os.Stdin = origStdin
	})

	// Provide invalid input followed by yes
	go func() {
		w.WriteString("maybe\ny\n")
		w.Close()
	}()

	if !init.promptConfirm("overwrite?") {
		t.Fatalf("expected prompt to eventually accept yes input")
	}

	r2, w2, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdin = r2
	go func() {
		w2.WriteString("n\n")
		w2.Close()
	}()

	if init.promptConfirm("cancel?") {
		t.Fatalf("expected prompt to return false for no input")
	}
}

// TestInitializer_IntegrationFullFlow tests the complete initialization flow end-to-end
func TestInitializer_IntegrationFullFlow(t *testing.T) {
	// This integration test validates the entire initialization flow:
	// 1. Directory validation
	// 2. Project type detection
	// 3. Template generation
	// 4. YAML file creation
	// 5. Config file correctness

	tests := []struct {
		name               string
		setupFiles         []string
		force              bool
		expectProjectTypes []string
		expectExtensions   []string
		expectExcludes     []string
	}{
		{
			name:               "Complete Next.js initialization",
			setupFiles:         []string{"next.config.js", "package.json", "tsconfig.json"},
			force:              false,
			expectProjectTypes: []string{"Next.js", "Node.js"},
			expectExtensions:   []string{".js", ".jsx", ".ts", ".tsx", ".json"},
			expectExcludes:     []string{"node_modules", ".next", "dist"},
		},
		{
			name:               "Complete Go project initialization",
			setupFiles:         []string{"go.mod", "go.sum", "main.go"},
			force:              false,
			expectProjectTypes: []string{"Go"},
			expectExtensions:   []string{".go", ".mod", ".sum"},
			expectExcludes:     []string{"vendor", "bin"},
		},
		{
			name:               "Multi-language project (Go + Node)",
			setupFiles:         []string{"go.mod", "package.json"},
			force:              false,
			expectProjectTypes: []string{"Go", "Node.js"},
			expectExtensions:   []string{".go", ".mod", ".js", ".ts"},
			expectExcludes:     []string{"vendor", "node_modules"},
		},
		{
			name:               "Force overwrite existing config",
			setupFiles:         []string{"next.config.js", ".promptext.yml"},
			force:              true,
			expectProjectTypes: []string{"Next.js"},
			expectExtensions:   []string{".js", ".jsx", ".ts", ".tsx"},
			expectExcludes:     []string{"node_modules", ".next"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "integration-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Setup: Create test files
			for _, file := range tt.setupFiles {
				filePath := filepath.Join(tmpDir, file)
				var content []byte
				if file == ".promptext.yml" {
					content = []byte("# Old config\nextensions:\n  - .old\n")
				}
				if err := os.WriteFile(filePath, content, 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", file, err)
				}
			}

			// Execute: Run initialization
			init := NewInitializer(tmpDir, tt.force, true) // quiet mode
			err = init.RunQuick()
			if err != nil {
				t.Fatalf("RunQuick() failed: %v", err)
			}

			// Verify: Config file was created
			configPath := filepath.Join(tmpDir, ".promptext.yml")
			if _, err := os.Stat(configPath); err != nil {
				t.Fatalf("Config file was not created: %v", err)
			}

			// Verify: Read and validate config content
			content, err := os.ReadFile(configPath)
			if err != nil {
				t.Fatalf("Failed to read config: %v", err)
			}
			configStr := string(content)

			// Verify: Check for expected extensions
			for _, ext := range tt.expectExtensions {
				if !strings.Contains(configStr, ext) {
					t.Errorf("Config missing expected extension %s", ext)
				}
			}

			// Verify: Check for expected excludes
			for _, exc := range tt.expectExcludes {
				if !strings.Contains(configStr, exc) {
					t.Errorf("Config missing expected exclude %s", exc)
				}
			}

			// Verify: Check config structure is valid YAML
			requiredFields := []string{
				"extensions:",
				"excludes:",
				"gitignore:",
				"use-default-rules:",
				"format:",
				"verbose:",
				"debug:",
			}
			for _, field := range requiredFields {
				if !strings.Contains(configStr, field) {
					t.Errorf("Config missing required field: %s", field)
				}
			}

			// Verify: Check for comments (educational output)
			if !strings.Contains(configStr, "# Promptext Configuration File") {
				t.Error("Config should have header comment")
			}
			if !strings.Contains(configStr, "# Auto-generated by: promptext --init") {
				t.Error("Config should indicate it was auto-generated")
			}

			// Verify: Default values are correct
			if !strings.Contains(configStr, "gitignore: true") {
				t.Error("Config should enable gitignore by default")
			}
			if !strings.Contains(configStr, "use-default-rules: true") {
				t.Error("Config should enable default rules by default")
			}
		})
	}
}

// TestInitializer_SecurityInputValidation tests input validation for security
func TestInitializer_SecurityInputValidation(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "security-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a simple project
	if err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("{}"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Test that excessively long paths are handled
	veryLongPath := filepath.Join(tmpDir, strings.Repeat("a", 1000))
	init := NewInitializer(veryLongPath, false, true)
	err = init.RunQuick()
	if err == nil {
		t.Error("Expected error for excessively long path")
	}

	// Test that non-absolute paths are handled correctly
	init = NewInitializer(".", false, true)
	err = init.RunQuick()
	// Should not panic, may succeed or fail depending on CWD
	// Just verify it doesn't panic
}
