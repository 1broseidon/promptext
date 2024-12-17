package filter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1broseidon/promptext/internal/filter"
)

func TestNewGitIgnore(t *testing.T) {
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

			got, err := filter.NewGitIgnore(gitignorePath)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewGitIgnore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(got.patterns) != len(tt.patterns) {
					t.Errorf("NewGitIgnore() got %v patterns, want %v patterns", len(got.patterns), len(tt.patterns))
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
		got, err := filter.NewGitIgnore(filepath.Join(tmpDir, "nonexistent"))
		if err != nil {
			t.Errorf("NewGitIgnore() error = %v, want nil for non-existent file", err)
			return
		}
		if len(got.patterns) != 0 {
			t.Errorf("NewGitIgnore() got %v patterns, want 0 for non-existent file", len(got.patterns))
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
			name: "common IDE and editor files",
			patterns: []string{
				".idea/",
				".vscode/",
				"*.swp",
				"*.swo",
				".DS_Store",
				"Thumbs.db",
			},
			paths: map[string]bool{
				".idea/workspace.xml":     true,
				".vscode/settings.json":   true,
				"file.swp":               true,
				".DS_Store":              true,
				"src/.DS_Store":          true,
				"normal.txt":             false,
			},
		},
		{
			name: "build and dependency directories",
			patterns: []string{
				"node_modules/",
				"dist/",
				"build/",
				"vendor/",
				"target/",
			},
			paths: map[string]bool{
				"node_modules/package.json": true,
				"dist/bundle.js":           true,
				"build/output":             true,
				"vendor/module/file.go":    true,
				"target/debug":             true,
				"src/file.js":             false,
			},
		},
		{
			name: "common config and env files",
			patterns: []string{
				".env",
				".env.local",
				"config.json",
				"*.config.js",
				"settings.json",
			},
			paths: map[string]bool{
				".env":              true,
				".env.local":        true,
				"config.json":       true,
				"webpack.config.js": true,
				"src/.env":          true,
				"normal.json":       false,
			},
		},
		{
			name: "common log and cache files",
			patterns: []string{
				"*.log",
				"npm-debug.log*",
				".npm",
				".cache/",
				"*.cache",
			},
			paths: map[string]bool{
				"error.log":          true,
				"npm-debug.log":      true,
				"npm-debug.log.1":    true,
				".npm/package.json":  true,
				".cache/files":       true,
				"style.css.cache":    true,
				"important.txt":      false,
			},
		},
		{
			name: "common test and coverage files",
			patterns: []string{
				"coverage/",
				"*.test.js",
				"__tests__/",
				"*.spec.ts",
				"junit.xml",
			},
			paths: map[string]bool{
				"coverage/lcov.info":     true,
				"src/module.test.js":     true,
				"__tests__/utils.js":     true,
				"component.spec.ts":      true,
				"src/component.spec.ts":  true,
				"junit.xml":             true,
				"src/main.js":           false,
			},
		},
		{
			name: "common backup and temp files",
			patterns: []string{
				"*~",
				"*.bak",
				"*.tmp",
				"temp/",
				"tmp/",
			},
			paths: map[string]bool{
				"document~":        true,
				"backup.bak":       true,
				"data.tmp":         true,
				"temp/file":        true,
				"tmp/cache":        true,
				"src/file.txt":     false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gi := &filter.GitIgnore{Patterns: tt.patterns}

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
	gi := &filter.GitIgnore{}

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
