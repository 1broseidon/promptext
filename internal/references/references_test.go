package references

import (
	"testing"
)

func TestExtractFileReferences(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		currentDir string
		rootDir   string
		allFiles  []string
		want      *ReferenceMap
	}{
		{
			name: "Markdown links",
			content: `See [Configuration](docs/config.md)
Check [external](https://example.com)
[Relative](../README.md)`,
			currentDir: "docs/guide",
			rootDir:    "/project",
			allFiles:   []string{"docs/config.md", "README.md"},
			want: &ReferenceMap{
				Internal: map[string][]string{
					"docs/guide": {"docs/config.md", "README.md"},
				},
				External: map[string][]string{
					"docs/guide": {"https://example.com"},
				},
			},
		},
		{
			name: "Comment references",
			content: `// ref: utils/helper.go
// see: https://pkg.go.dev/fmt
# reference: ../config.yaml`,
			currentDir: "internal/app",
			rootDir:    "/project",
			allFiles:   []string{"utils/helper.go", "config.yaml"},
			want: &ReferenceMap{
				Internal: map[string][]string{
					"internal/app": {"utils/helper.go", "config.yaml"},
				},
				External: map[string][]string{
					"internal/app": {"https://pkg.go.dev/fmt"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractFileReferences(tt.content, tt.currentDir, tt.rootDir, tt.allFiles)
			
			// Compare internal references
			for dir, refs := range tt.want.Internal {
				gotRefs, ok := got.Internal[dir]
				if !ok || !stringSliceEqual(gotRefs, refs) {
					t.Errorf("Internal references mismatch for %s:\ngot  %v\nwant %v", 
						dir, gotRefs, refs)
				}
			}

			// Compare external references
			for dir, refs := range tt.want.External {
				gotRefs, ok := got.External[dir]
				if !ok || !stringSliceEqual(gotRefs, refs) {
					t.Errorf("External references mismatch for %s:\ngot  %v\nwant %v",
						dir, gotRefs, refs)
				}
			}
		})
	}
}

func TestIsExternalReference(t *testing.T) {
	tests := []struct {
		name string
		ref  string
		want bool
	}{
		{"HTTP URL", "https://example.com", true},
		{"HTTPS URL", "https://github.com", true},
		{"GitHub package", "github.com/user/pkg", true},
		{"NPM package", "@angular/core", true},
		{"Local file", "./utils/helper.go", false},
		{"Relative path", "../config.yaml", false},
		{"Absolute path", "/etc/config", false},
		{"Simple filename", "main.go", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isExternalReference(tt.ref); got != tt.want {
				t.Errorf("isExternalReference(%q) = %v, want %v", tt.ref, got, tt.want)
			}
		})
	}
}

func TestResolveReference(t *testing.T) {
	tests := []struct {
		name       string
		ref        string
		currentDir string
		rootDir    string
		allFiles   []string
		want       string
	}{
		{
			name:       "Exact match",
			ref:        "config.go",
			currentDir: "internal",
			rootDir:    "/project",
			allFiles:   []string{"internal/config.go"},
			want:       "internal/config.go",
		},
		{
			name:       "Relative path",
			ref:        "../utils/helper",
			currentDir: "internal/app",
			rootDir:    "/project",
			allFiles:   []string{"internal/utils/helper.go"},
			want:       "internal/utils/helper.go",
		},
		{
			name:       "With extension",
			ref:        "./config.yaml",
			currentDir: "config",
			rootDir:    "/project",
			allFiles:   []string{"config/config.yaml"},
			want:       "config/config.yaml",
		},
		{
			name:       "No match",
			ref:        "nonexistent.go",
			currentDir: "src",
			rootDir:    "/project",
			allFiles:   []string{"src/main.go"},
			want:       "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := resolveReference(tt.ref, tt.currentDir, tt.rootDir, tt.allFiles)
			if got != tt.want {
				t.Errorf("resolveReference() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func stringSliceEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	seen := make(map[string]bool)
	for _, str := range a {
		seen[str] = true
	}
	for _, str := range b {
		if !seen[str] {
			return false
		}
	}
	return true
}
