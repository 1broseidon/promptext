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
				IgnoreDefault: true, // Enable default rules including binary detection
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
				IgnoreDefault: true,
			},
			path: "node_modules/package.json",
			want: false,
		},
		{
			name: "mixed patterns",
			opts: Options{
				Includes:      []string{".go", ".md"},
				Excludes:      []string{"vendor/", "test/"},
				IgnoreDefault: true,
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
