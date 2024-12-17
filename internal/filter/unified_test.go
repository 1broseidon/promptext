package filter

import (
	"testing"

)

func TestUnifiedFilter_ShouldProcess(t *testing.T) {
	// Helper function to create a gitignore for testing
	createGitIgnore := func(patterns []string) *gitignore.GitIgnore {
		return &gitignore.GitIgnore{Patterns: patterns}
	}

	tests := []struct {
		name            string
		path            string
		allowedExts     []string
		excludePatterns []string
		gitIgnoreRules  []string
		want            bool
	}{
		{
			name: "default ignored directory",
			path: "node_modules/package.json",
			want: false,
		},
		{
			name: "default ignored extension",
			path: "images/photo.jpg",
			want: false,
		},
		{
			name:        "allowed extension",
			path:        "src/main.go",
			allowedExts: []string{".go"},
			want:        true,
		},
		{
			name:        "not allowed extension when filter active",
			path:        "src/main.js",
			allowedExts: []string{".go"},
			want:        false,
		},
		{
			name:            "excluded by pattern",
			path:            "test/test_file.go",
			excludePatterns: []string{"test/"},
			want:            false,
		},
		{
			name:           "excluded by gitignore",
			path:           "build/output.txt",
			gitIgnoreRules: []string{"build/"},
			want:           false,
		},
		{
			name:        "case insensitive extension match",
			path:        "doc.MD",
			allowedExts: []string{".md"},
			want:        true,
		},
		{
			name:        "no extension filters",
			path:        "README.md",
			allowedExts: []string{},
			want:        true,
		},
		{
			name:        "multiple allowed extensions",
			path:        "src/main.go",
			allowedExts: []string{".js", ".go", ".py"},
			want:        true,
		},
		{
			name:            "multiple exclusion patterns",
			path:            "vendor/lib/package.go",
			excludePatterns: []string{"vendor/", "test/"},
			want:            false,
		},
		{
			name: "binary file extension",
			path: "bin/program.exe",
			want: false,
		},
		{
			name: "image file extension",
			path: "assets/logo.png",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gitIgnore := createGitIgnore(tt.gitIgnoreRules)
			filter := NewUnifiedFilter(gitIgnore, tt.allowedExts, tt.excludePatterns)

			got := filter.ShouldProcess(tt.path)
			if got != tt.want {
				t.Errorf("UnifiedFilter.ShouldProcess(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNewUnifiedFilter(t *testing.T) {
	gitIgnore := &gitignore.GitIgnore{}
	extensions := []string{".go", ".md"}
	excludes := []string{"vendor/"}

	filter := NewUnifiedFilter(gitIgnore, extensions, excludes)

	if filter.gitIgnore != gitIgnore {
		t.Error("NewUnifiedFilter did not set gitIgnore correctly")
	}

	if len(filter.allowedExtensions) != len(extensions) {
		t.Errorf("NewUnifiedFilter got %d extensions, want %d", len(filter.allowedExtensions), len(extensions))
	}

	if len(filter.configExcludes) != len(excludes) {
		t.Errorf("NewUnifiedFilter got %d excludes, want %d", len(filter.configExcludes), len(excludes))
	}

	if len(filter.defaultIgnores) != len(DefaultIgnoreDirs) {
		t.Errorf("NewUnifiedFilter got %d default ignores, want %d", len(filter.defaultIgnores), len(DefaultIgnoreDirs))
	}

	if len(filter.defaultIgnoreExts) != len(DefaultIgnoreExtensions) {
		t.Errorf("NewUnifiedFilter got %d default ignore extensions, want %d", len(filter.defaultIgnoreExts), len(DefaultIgnoreExtensions))
	}
}

func TestUnifiedFilter_PathMatching(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		excludes []string
		want     bool
	}{
		{
			name:     "exact path match",
			path:     "vendor/package.json",
			excludes: []string{"vendor/package.json"},
			want:     false,
		},
		{
			name:     "glob pattern match",
			path:     "test/file.test.js",
			excludes: []string{"*.test.js"},
			want:     false,
		},
		{
			name:     "directory prefix match",
			path:     "src/vendor/lib.go",
			excludes: []string{"vendor/"},
			want:     false,
		},
		{
			name:     "no match",
			path:     "src/main.go",
			excludes: []string{"vendor/", "test/"},
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filter := NewUnifiedFilter(nil, nil, tt.excludes)
			got := filter.ShouldProcess(tt.path)
			if got != tt.want {
				t.Errorf("PathMatching(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
