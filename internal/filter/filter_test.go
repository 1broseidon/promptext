package filter

import (
	"os"
	"path/filepath"
	"testing"
)

func createTestFile(t *testing.T, path string, content []byte) {
	t.Helper()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, content, 0644); err != nil {
		t.Fatal(err)
	}
}

func TestFilter_ShouldProcess(t *testing.T) {
	// Create temporary test directory
	tmpDir := t.TempDir()

	// Create test files
	textFile := filepath.Join(tmpDir, "text.txt")
	createTestFile(t, textFile, []byte("This is a text file\nwith multiple lines\n"))

	binaryFile := filepath.Join(tmpDir, "binary.dat")
	binaryContent := []byte{0x00, 0x01, 0x02, 0x03} // Binary content with null bytes
	createTestFile(t, binaryFile, binaryContent)

	tests := []struct {
		name string
		opts Options
		path string
		want bool
	}{
		{
			name: "text file",
			opts: Options{},
			path: textFile,
			want: true,
		},
		{
			name: "binary file",
			opts: Options{
				UseDefaultRules: true, // Enable default rules including binary detection
			},
			path: binaryFile,
			want: false,
		},
		{
			name: "extension include",
			opts: Options{
				Includes: []string{".go"},
			},
			path: "src/main.go",
			want: true,
		},
		{
			name: "directory exclude",
			opts: Options{
				Excludes: []string{"vendor/"},
			},
			path: "vendor/module/file.go",
			want: false,
		},
		{
			name: "default ignores",
			opts: Options{
				UseDefaultRules: true,
			},
			path: "node_modules/package.json",
			want: false,
		},
		{
			name: "mixed patterns",
			opts: Options{
				Includes:        []string{".go", ".md"},
				Excludes:        []string{"vendor/", "test/"},
				UseDefaultRules: true,
			},
			path: "src/main.go",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := New(tt.opts)
			if got := f.ShouldProcess(tt.path); got != tt.want {
				t.Errorf("ShouldProcess(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}

	// Cleanup
	os.RemoveAll(tmpDir)
}

func TestParseGitIgnore(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    []string
		wantErr bool
	}{
		{
			name:    "basic patterns",
			content: "node_modules/\n*.log\n.env\n",
			want:    []string{"node_modules/", "*.log", ".env"},
			wantErr: false,
		},
		{
			name:    "with comments",
			content: "# Build output\ndist/\n# Dependencies\nnode_modules/\n",
			want:    []string{"dist/", "node_modules/"},
			wantErr: false,
		},
		{
			name:    "with empty lines",
			content: "dist/\n\nnode_modules/\n\n*.log\n",
			want:    []string{"dist/", "node_modules/", "*.log"},
			wantErr: false,
		},
		{
			name:    "empty file",
			content: "",
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "only comments",
			content: "# Comment 1\n# Comment 2\n",
			want:    []string{},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			gitignorePath := filepath.Join(tmpDir, ".gitignore")
			if err := os.WriteFile(gitignorePath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to create .gitignore: %v", err)
			}

			got, err := ParseGitIgnore(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseGitIgnore() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if len(got) != len(tt.want) {
				t.Errorf("ParseGitIgnore() got %d patterns, want %d", len(got), len(tt.want))
				return
			}

			for i, pattern := range got {
				if pattern != tt.want[i] {
					t.Errorf("ParseGitIgnore()[%d] = %q, want %q", i, pattern, tt.want[i])
				}
			}
		})
	}
}

func TestParseGitIgnore_NoFile(t *testing.T) {
	tmpDir := t.TempDir()
	patterns, err := ParseGitIgnore(tmpDir)
	if err != nil {
		t.Errorf("ParseGitIgnore() with missing file should return nil error, got: %v", err)
	}
	if patterns != nil {
		t.Errorf("ParseGitIgnore() with missing file should return nil patterns, got: %v", patterns)
	}
}

func TestMergeAndDedupePatterns(t *testing.T) {
	tests := []struct {
		name     string
		patterns [][]string
		want     []string
	}{
		{
			name:     "single set",
			patterns: [][]string{{"*.log", "dist/"}},
			want:     []string{"*.log", "dist/"},
		},
		{
			name:     "multiple sets with duplicates",
			patterns: [][]string{{"*.log", "dist/"}, {"*.log", "node_modules/"}, {"dist/", "vendor/"}},
			want:     []string{"*.log", "dist/", "node_modules/", "vendor/"},
		},
		{
			name:     "empty sets",
			patterns: [][]string{{}, {}},
			want:     []string{},
		},
		{
			name:     "mixed empty and filled",
			patterns: [][]string{{}, {"*.log"}, {}},
			want:     []string{"*.log"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MergeAndDedupePatterns(tt.patterns...)
			if len(got) != len(tt.want) {
				t.Errorf("MergeAndDedupePatterns() got %d patterns, want %d", len(got), len(tt.want))
				return
			}
			for i, pattern := range got {
				if pattern != tt.want[i] {
					t.Errorf("MergeAndDedupePatterns()[%d] = %q, want %q", i, pattern, tt.want[i])
				}
			}
		})
	}
}

func TestGetFileType(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		wantType     string
		wantCategory string
		wantTest     bool
		wantEntry    bool
	}{
		// Test files
		{
			name:         "go test file",
			path:         "internal/processor/processor_test.go",
			wantType:     "test",
			wantCategory: "test:go",
			wantTest:     true,
			wantEntry:    false,
		},
		{
			name:         "js spec file",
			path:         "src/components/Button.spec.js",
			wantType:     "test",
			wantCategory: "test:js",
			wantTest:     true,
			wantEntry:    false,
		},
		{
			name:         "python test file",
			path:         "tests/test_utils.py",
			wantType:     "test",
			wantCategory: "test:py",
			wantTest:     true,
			wantEntry:    false,
		},

		// Entry points
		{
			name:         "main.go",
			path:         "cmd/promptext/main.go",
			wantType:     "source",
			wantCategory: "entry:go",
			wantTest:     false,
			wantEntry:    true,
		},
		{
			name:         "index.js",
			path:         "src/index.js",
			wantType:     "source",
			wantCategory: "entry:js",
			wantTest:     false,
			wantEntry:    true,
		},
		{
			name:         "app.py",
			path:         "app.py",
			wantType:     "source",
			wantCategory: "entry:py",
			wantTest:     false,
			wantEntry:    true,
		},

		// Config files
		{
			name:         "yaml config",
			path:         ".github/workflows/build.yml",
			wantType:     "config",
			wantCategory: "config:yaml",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "json config",
			path:         ".eslintrc.json",
			wantType:     "config",
			wantCategory: "config:json",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "toml config",
			path:         "pyproject.toml",
			wantType:     "config",
			wantCategory: "config:toml",
			wantTest:     false,
			wantEntry:    false,
		},

		// Documentation
		{
			name:         "markdown doc",
			path:         "README.md",
			wantType:     "doc",
			wantCategory: "doc:markdown",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "text doc",
			path:         "LICENSE.txt",
			wantType:     "doc",
			wantCategory: "doc:text",
			wantTest:     false,
			wantEntry:    false,
		},

		// Source files
		{
			name:         "go source",
			path:         "internal/processor/processor.go",
			wantType:     "source",
			wantCategory: "source:go",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "javascript source",
			path:         "src/utils/helper.js",
			wantType:     "source",
			wantCategory: "source:javascript",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "typescript source",
			path:         "src/components/Button.ts",
			wantType:     "source",
			wantCategory: "source:typescript",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "python source",
			path:         "utils/helper.py",
			wantType:     "source",
			wantCategory: "source:python",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "rust source",
			path:         "src/main.rs",
			wantType:     "source",
			wantCategory: "source:rust",
			wantTest:     false,
			wantEntry:    false,
		},

		// Dependency files (some categorized by extension first)
		{
			name:         "package.json",
			path:         "package.json",
			wantType:     "config", // .json extension matches config before dependency
			wantCategory: "config:json",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "go.mod",
			path:         "go.mod",
			wantType:     "dependency",
			wantCategory: "dep:go",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "requirements.txt",
			path:         "requirements.txt",
			wantType:     "doc", // .txt extension matches doc before dependency
			wantCategory: "doc:text",
			wantTest:     false,
			wantEntry:    false,
		},
		{
			name:         "Cargo.toml",
			path:         "Cargo.toml",
			wantType:     "config", // .toml extension matches config before dependency
			wantCategory: "config:toml",
			wantTest:     false,
			wantEntry:    false,
		},

		// Unknown/other
		{
			name:         "unknown extension",
			path:         "data.bin",
			wantType:     "source",
			wantCategory: "source:other",
			wantTest:     false,
			wantEntry:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileType(tt.path, nil)
			if got.Type != tt.wantType {
				t.Errorf("GetFileType(%q).Type = %q, want %q", tt.path, got.Type, tt.wantType)
			}
			if got.Category != tt.wantCategory {
				t.Errorf("GetFileType(%q).Category = %q, want %q", tt.path, got.Category, tt.wantCategory)
			}
			if got.IsTest != tt.wantTest {
				t.Errorf("GetFileType(%q).IsTest = %v, want %v", tt.path, got.IsTest, tt.wantTest)
			}
			if got.IsEntryPoint != tt.wantEntry {
				t.Errorf("GetFileType(%q).IsEntryPoint = %v, want %v", tt.path, got.IsEntryPoint, tt.wantEntry)
			}
		})
	}
}

func TestGetFileType_WithFilter(t *testing.T) {
	// Test with filter that excludes certain paths
	f := New(Options{
		Excludes:        []string{"vendor/"},
		UseDefaultRules: false,
	})

	tests := []struct {
		name     string
		path     string
		wantType string
	}{
		{
			name:     "excluded path returns empty",
			path:     "vendor/module/file.go",
			wantType: "",
		},
		{
			name:     "included path works",
			path:     "src/main.go",
			wantType: "source",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetFileType(tt.path, f)
			if got.Type != tt.wantType {
				t.Errorf("GetFileType(%q).Type = %q, want %q", tt.path, got.Type, tt.wantType)
			}
		})
	}
}

func TestIsTestFile(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"processor_test.go", true},
		{"test_utils.py", true},
		{"utils.test.js", true},
		{"button.spec.js", true},
		{"main.go", false},
		{"utils.py", false},
		{"component.js", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			base := filepath.Base(tt.path)
			got := isTestFile(tt.path, base)
			if got != tt.want {
				t.Errorf("isTestFile(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestIsEntryPoint(t *testing.T) {
	tests := []struct {
		base string
		want bool
	}{
		{"main.go", true},
		{"index.js", true},
		{"app.py", true},
		{"index.ts", true},
		{"server.js", true},
		{"helper.go", false},
		{"utils.js", false},
		{"config.py", false},
	}

	for _, tt := range tests {
		t.Run(tt.base, func(t *testing.T) {
			got := isEntryPoint(tt.base)
			if got != tt.want {
				t.Errorf("isEntryPoint(%q) = %v, want %v", tt.base, got, tt.want)
			}
		})
	}
}

func TestGetConfigType(t *testing.T) {
	tests := []struct {
		ext          string
		wantType     string
		wantCategory string
	}{
		{".yml", "config", "config:yaml"},
		{".yaml", "config", "config:yaml"},
		{".json", "config", "config:json"},
		{".toml", "config", "config:toml"},
		{".ini", "config", "config:ini"},
		{".conf", "config", "config:ini"},
		{".go", "", ""},
		{".js", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			gotType, gotCategory := getConfigType(tt.ext)
			if gotType != tt.wantType {
				t.Errorf("getConfigType(%q) type = %q, want %q", tt.ext, gotType, tt.wantType)
			}
			if gotCategory != tt.wantCategory {
				t.Errorf("getConfigType(%q) category = %q, want %q", tt.ext, gotCategory, tt.wantCategory)
			}
		})
	}
}

func TestGetDocType(t *testing.T) {
	tests := []struct {
		ext          string
		wantType     string
		wantCategory string
	}{
		{".md", "doc", "doc:markdown"},
		{".txt", "doc", "doc:text"},
		{".rst", "doc", "doc:rst"},
		{".adoc", "doc", "doc:asciidoc"},
		{".go", "", ""},
		{".js", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			gotType, gotCategory := getDocType(tt.ext)
			if gotType != tt.wantType {
				t.Errorf("getDocType(%q) type = %q, want %q", tt.ext, gotType, tt.wantType)
			}
			if gotCategory != tt.wantCategory {
				t.Errorf("getDocType(%q) category = %q, want %q", tt.ext, gotCategory, tt.wantCategory)
			}
		})
	}
}

func TestGetSourceType(t *testing.T) {
	tests := []struct {
		ext          string
		wantType     string
		wantCategory string
	}{
		{".go", "source", "source:go"},
		{".js", "source", "source:javascript"},
		{".ts", "source", "source:typescript"},
		{".jsx", "source", "source:react"},
		{".tsx", "source", "source:react"},
		{".py", "source", "source:python"},
		{".rs", "source", "source:rust"},
		{".java", "source", "source:java"},
		{".txt", "", ""},
		{".md", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			gotType, gotCategory := getSourceType(tt.ext)
			if gotType != tt.wantType {
				t.Errorf("getSourceType(%q) type = %q, want %q", tt.ext, gotType, tt.wantType)
			}
			if gotCategory != tt.wantCategory {
				t.Errorf("getSourceType(%q) category = %q, want %q", tt.ext, gotCategory, tt.wantCategory)
			}
		})
	}
}

func TestGetDependencyType(t *testing.T) {
	tests := []struct {
		base         string
		wantType     string
		wantCategory string
	}{
		{"package.json", "dependency", "dep:node"},
		{"package-lock.json", "dependency", "dep:node"},
		{"yarn.lock", "dependency", "dep:node"},
		{"go.mod", "dependency", "dep:go"},
		{"go.sum", "dependency", "dep:go"},
		{"requirements.txt", "dependency", "dep:python"},
		{"Pipfile", "dependency", "dep:python"},
		{"pyproject.toml", "dependency", "dep:python"},
		{"Cargo.toml", "dependency", "dep:rust"},
		{"Cargo.lock", "dependency", "dep:rust"},
		{"Gemfile", "dependency", "dep:ruby"},
		{"composer.json", "dependency", "dep:php"},
		{"main.go", "", ""},
		{"README.md", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.base, func(t *testing.T) {
			gotType, gotCategory := getDependencyType(tt.base)
			if gotType != tt.wantType {
				t.Errorf("getDependencyType(%q) type = %q, want %q", tt.base, gotType, tt.wantType)
			}
			if gotCategory != tt.wantCategory {
				t.Errorf("getDependencyType(%q) category = %q, want %q", tt.base, gotCategory, tt.wantCategory)
			}
		})
	}
}

func TestFilter_IsExcluded(t *testing.T) {
	f := New(Options{
		Excludes:        []string{"vendor/", "*.log"},
		UseDefaultRules: false,
	})

	tests := []struct {
		path string
		want bool
	}{
		{"vendor/module/file.go", true},
		{"src/main.go", false},
		{"debug.log", true},
		{"README.md", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := f.IsExcluded(tt.path)
			if got != tt.want {
				t.Errorf("IsExcluded(%q) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}

func TestNew_GitIgnoreIntegration(t *testing.T) {
	// Create temporary directory with .gitignore
	tmpDir := t.TempDir()
	gitignoreContent := "node_modules/\n*.log\ndist/\n"
	gitignorePath := filepath.Join(tmpDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		t.Fatalf("Failed to create .gitignore: %v", err)
	}

	// Change to temp directory
	oldWd, _ := os.Getwd()
	defer os.Chdir(oldWd)
	os.Chdir(tmpDir)

	// Create filter with gitignore enabled
	f := New(Options{
		UseGitIgnore:    true,
		UseDefaultRules: false,
	})

	tests := []struct {
		path string
		want bool
	}{
		{"node_modules/package/file.js", false},
		{"dist/bundle.js", false},
		{"debug.log", false},
		{"src/main.go", true},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			got := f.ShouldProcess(tt.path)
			if got != tt.want {
				t.Errorf("ShouldProcess(%q) with gitignore = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
