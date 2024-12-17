package gitignore

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNew(t *testing.T) {
	// Create temp directory for test files
	tmpDir, err := os.MkdirTemp("", "gitignore_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		patterns []string
	}{
		{
			name: "basic patterns",
			content: `node_modules/
*.log
.env
# comment
.DS_Store

`,
			wantErr:  false,
			patterns: []string{"node_modules/", "*.log", ".env", ".DS_Store"},
		},
		{
			name:     "empty file",
			content:  "",
			wantErr:  false,
			patterns: nil,
		},
		{
			name: "only comments",
			content: `# ignore this
# and this`,
			wantErr:  false,
			patterns: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create gitignore file
			gitignorePath := filepath.Join(tmpDir, ".gitignore")
			if err := os.WriteFile(gitignorePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write .gitignore file: %v", err)
			}

			got, err := New(gitignorePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.patterns) != len(tt.patterns) {
					t.Errorf("New() got %v patterns, want %v patterns", len(got.patterns), len(tt.patterns))
					return
				}

				for i, pattern := range got.patterns {
					if pattern != tt.patterns[i] {
						t.Errorf("Pattern %d: got %v, want %v", i, pattern, tt.patterns[i])
					}
				}
			}
		})
	}

	// Test non-existent file
	t.Run("non-existent file", func(t *testing.T) {
		got, err := New(filepath.Join(tmpDir, "nonexistent"))
		if err != nil {
			t.Errorf("New() error = %v, want nil for non-existent file", err)
			return
		}
		if len(got.patterns) != 0 {
			t.Errorf("New() got %v patterns, want 0 for non-existent file", len(got.patterns))
		}
	})
}

func TestGitIgnore_ShouldIgnore(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		paths    map[string]bool // path -> should ignore
	}{
		{
			name:     "empty patterns",
			patterns: []string{},
			paths: map[string]bool{
				"file.txt": false,
			},
		},
		{
			name: "exact matches",
			patterns: []string{
				".env",
				"config.json",
			},
			paths: map[string]bool{
				".env":         true,
				"config.json":  true,
				"other.json":   false,
				"src/.env":     true,
				"config.json2": false,
			},
		},
		{
			name: "directory patterns",
			patterns: []string{
				"node_modules/",
				"vendor/",
			},
			paths: map[string]bool{
				"node_modules/package.json": true,
				"vendor/module/file.go":     true,
				"src/vendor/":               true,
				"myvendor/file":            false,
			},
		},
		{
			name: "glob patterns",
			patterns: []string{
				"*.log",
				".git*",
				"temp*",
			},
			paths: map[string]bool{
				"error.log":       true,
				".gitignore":      true,
				"logs/error.log":  true,
				"temp":           true,
				"temp.txt":       true,
				"nottemp.txt":    false,
				".gitkeep":       true,
				"logs/.gitkeep":  true,
			},
		},
		{
			name: "complex patterns",
			patterns: []string{
				"**/*.test.js",
				"build/",
				".env*",
				"*.bak",
			},
			paths: map[string]bool{
				"src/file.test.js":    true,
				"test/unit.test.js":   true,
				"build/output.js":     true,
				".env":                true,
				".env.local":          true,
				"backup.bak":          true,
				"src/backup.bak":      true,
				"normal.js":           false,
				"src/component.js":    false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &GitIgnore{patterns: tt.patterns}

			for path, want := range tt.paths {
				got := gi.ShouldIgnore(path)
				if got != want {
					t.Errorf("ShouldIgnore(%q) = %v, want %v", path, got, want)
				}
			}
		})
	}
}

func TestGitIgnore_MatchFunctions(t *testing.T) {
	gi := &GitIgnore{}

	t.Run("matchExact", func(t *testing.T) {
		tests := []struct {
			pattern  string
			path     string
			baseName string
			want     bool
		}{
			{".env", ".env", ".env", true},
			{".env", "src/.env", ".env", true},
			{".env", "other", "other", false},
		}

		for _, tt := range tests {
			if got := gi.matchExact(tt.pattern, tt.path, tt.baseName); got != tt.want {
				t.Errorf("matchExact(%q, %q, %q) = %v, want %v", 
					tt.pattern, tt.path, tt.baseName, got, tt.want)
			}
		}
	})

	t.Run("matchDirectory", func(t *testing.T) {
		tests := []struct {
			pattern string
			path    string
			want    bool
		}{
			{"node_modules/", "node_modules/package.json", true},
			{"vendor/", "src/vendor/module.go", true},
			{"dist/", "something/else", false},
			{"temp", "temp/file", false}, // Not a directory pattern
		}

		for _, tt := range tests {
			if got := gi.matchDirectory(tt.pattern, tt.path); got != tt.want {
				t.Errorf("matchDirectory(%q, %q) = %v, want %v", 
					tt.pattern, tt.path, got, tt.want)
			}
		}
	})

	t.Run("matchGlobPattern", func(t *testing.T) {
		tests := []struct {
			pattern  string
			path     string
			baseName string
			want     bool
		}{
			{"*.log", "error.log", "error.log", true},
			{"*.log", "logs/error.log", "error.log", true},
			{".git*", ".gitignore", ".gitignore", true},
			{"temp*", "temporary", "temporary", true},
			{"*.txt", "not-match.log", "not-match.log", false},
		}

		for _, tt := range tests {
			if got := gi.matchGlobPattern(tt.pattern, tt.path, tt.baseName); got != tt.want {
				t.Errorf("matchGlobPattern(%q, %q, %q) = %v, want %v", 
					tt.pattern, tt.path, tt.baseName, got, tt.want)
			}
		}
	})
}
