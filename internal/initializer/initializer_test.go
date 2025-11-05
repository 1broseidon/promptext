package initializer

import (
	"os"
	"path/filepath"
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
				} else if tt.errorMsg != "" && !contains(err.Error(), tt.errorMsg) {
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
				if !contains(configStr, expectedStr) {
					t.Errorf("Expected config to contain '%s', but it doesn't.\nConfig:\n%s", expectedStr, configStr)
				}
			}

			// Verify file structure
			if !contains(configStr, "extensions:") {
				t.Error("Config should contain 'extensions:' section")
			}
			if !contains(configStr, "excludes:") {
				t.Error("Config should contain 'excludes:' section")
			}
			if !contains(configStr, "gitignore: true") {
				t.Error("Config should have gitignore enabled by default")
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		func() bool {
			for i := 0; i <= len(s)-len(substr); i++ {
				if s[i:i+len(substr)] == substr {
					return true
				}
			}
			return false
		}())
}
