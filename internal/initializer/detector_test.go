package initializer

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFileDetector_Detect(t *testing.T) {
	tests := []struct {
		name          string
		files         []string
		expectedTypes []string
	}{
		{
			name:          "Next.js project",
			files:         []string{"next.config.js", "package.json"},
			expectedTypes: []string{"nextjs", "node"},
		},
		{
			name:          "Go project",
			files:         []string{"go.mod", "main.go"},
			expectedTypes: []string{"go"},
		},
		{
			name:          "Django project",
			files:         []string{"manage.py", "requirements.txt"},
			expectedTypes: []string{"django", "python"},
		},
		{
			name:          "Rust project",
			files:         []string{"Cargo.toml", "src/main.rs"},
			expectedTypes: []string{"rust"},
		},
		{
			name:          "Mixed Go + Node project",
			files:         []string{"go.mod", "package.json"},
			expectedTypes: []string{"go", "node"},
		},
		{
			name:          "Angular project",
			files:         []string{"angular.json", "package.json"},
			expectedTypes: []string{"angular", "node"},
		},
		{
			name:          "Laravel project",
			files:         []string{"artisan", "composer.json"},
			expectedTypes: []string{"laravel", "php"},
		},
		{
			name:          "Empty project",
			files:         []string{},
			expectedTypes: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "detector-test-*")
			if err != nil {
				t.Fatalf("Failed to create temp dir: %v", err)
			}
			defer os.RemoveAll(tmpDir)

			// Create test files
			for _, file := range tt.files {
				filePath := filepath.Join(tmpDir, file)
				dir := filepath.Dir(filePath)

				// Create directory if needed
				if dir != tmpDir {
					if err := os.MkdirAll(dir, 0755); err != nil {
						t.Fatalf("Failed to create directory %s: %v", dir, err)
					}
				}

				// Create file
				if err := os.WriteFile(filePath, []byte{}, 0644); err != nil {
					t.Fatalf("Failed to create file %s: %v", filePath, err)
				}
			}

			// Run detection
			detector := NewFileDetector()
			detected, err := detector.Detect(tmpDir)
			if err != nil {
				t.Fatalf("Detect() error = %v", err)
			}

			// Check results
			if len(detected) != len(tt.expectedTypes) {
				t.Errorf("Expected %d project types, got %d", len(tt.expectedTypes), len(detected))
				t.Errorf("Expected: %v", tt.expectedTypes)
				t.Errorf("Got: %v", getTypeNames(detected))
				return
			}

			// Verify each expected type is present
			detectedMap := make(map[string]bool)
			for _, pt := range detected {
				detectedMap[pt.Name] = true
			}

			for _, expectedType := range tt.expectedTypes {
				if !detectedMap[expectedType] {
					t.Errorf("Expected project type %s not found", expectedType)
				}
			}
		})
	}
}

func TestFileDetector_Priority(t *testing.T) {
	// Create temporary directory with multiple project indicators
	tmpDir, err := os.MkdirTemp("", "detector-priority-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create files for Next.js (priority 100) and generic Node (priority 60)
	files := []string{"next.config.js", "package.json"}
	for _, file := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, file), []byte{}, 0644); err != nil {
			t.Fatalf("Failed to create file %s: %v", file, err)
		}
	}

	// Run detection
	detector := NewFileDetector()
	detected, err := detector.Detect(tmpDir)
	if err != nil {
		t.Fatalf("Detect() error = %v", err)
	}

	// Verify Next.js is listed before Node.js (due to higher priority)
	if len(detected) < 2 {
		t.Fatalf("Expected at least 2 project types, got %d", len(detected))
	}

	if detected[0].Name != "nextjs" {
		t.Errorf("Expected nextjs to be first (highest priority), got %s", detected[0].Name)
	}
}

// Helper function to extract type names from ProjectType slice
func getTypeNames(types []ProjectType) []string {
	names := make([]string, len(types))
	for i, pt := range types {
		names[i] = pt.Name
	}
	return names
}
