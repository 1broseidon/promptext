package rules

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBinaryRule_ExtensionDetection(t *testing.T) {
	rule := NewBinaryRule()

	testCases := []struct {
		filename string
		expected bool
		desc     string
	}{
		// Should be detected as binary by extension (no file I/O)
		{"test.exe", true, "executable file"},
		{"image.jpg", true, "JPEG image"},
		{"archive.zip", true, "ZIP archive"},
		{"document.pdf", true, "PDF document"},
		{"library.so", true, "shared object"},
		{"font.ttf", true, "TrueType font"},

		// Should not be detected as binary by extension
		{"script.py", false, "Python script"},
		{"config.json", false, "JSON config"},
		{"readme.md", false, "Markdown file"},
		{"source.go", false, "Go source"},
		{"style.css", false, "CSS file"},
		{"data.txt", false, "text file"},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			// Create a temporary file with the test name
			tmpDir, err := os.MkdirTemp("", "binary_test")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tmpDir)

			testFile := filepath.Join(tmpDir, tc.filename)

			// Create a small text file - if extension detection works,
			// content won't matter for binary extensions
			err = os.WriteFile(testFile, []byte("hello world"), 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.filename, result)
			}
		})
	}
}

func TestBinaryRule_SizeDetection(t *testing.T) {
	rule := NewBinaryRule()

	tmpDir, err := os.MkdirTemp("", "binary_size_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test large file detection (should be binary due to size)
	largeFile := filepath.Join(tmpDir, "large.unknown")
	largeData := make([]byte, 11*1024*1024) // 11MB
	for i := range largeData {
		largeData[i] = 'A' // Fill with text to ensure it's detected by size, not content
	}
	err = os.WriteFile(largeFile, largeData, 0644)
	if err != nil {
		t.Fatal(err)
	}

	if !rule.Match(largeFile) {
		t.Error("Large file should be detected as binary")
	}

	// Test empty file (should not be binary)
	emptyFile := filepath.Join(tmpDir, "empty.unknown")
	err = os.WriteFile(emptyFile, []byte{}, 0644)
	if err != nil {
		t.Fatal(err)
	}

	if rule.Match(emptyFile) {
		t.Error("Empty file should not be detected as binary")
	}
}

func TestBinaryRule_ContentDetection(t *testing.T) {
	rule := NewBinaryRule()

	tmpDir, err := os.MkdirTemp("", "binary_content_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name     string
		content  []byte
		expected bool
		desc     string
	}{
		{
			name:     "text.unknown",
			content:  []byte("This is plain text content\nwith multiple lines\n"),
			expected: false,
			desc:     "plain text should not be binary",
		},
		{
			name:     "null_bytes.unknown",
			content:  []byte("text with\x00null byte"),
			expected: true,
			desc:     "content with null bytes should be binary",
		},
		{
			name: "high_non_printable.unknown",
			content: func() []byte {
				// Create content with high ratio of non-printable chars
				data := make([]byte, 100)
				for i := 0; i < 60; i++ { // 60% non-printable
					data[i] = byte(i + 128) // Non-ASCII
				}
				for i := 60; i < 100; i++ {
					data[i] = 'A' // Printable
				}
				return data
			}(),
			expected: true,
			desc:     "high ratio of non-printable chars should be binary",
		},
		{
			name: "low_non_printable.unknown",
			content: func() []byte {
				// Create content with low ratio of non-printable chars
				data := make([]byte, 100)
				for i := 0; i < 20; i++ { // 20% non-printable
					data[i] = byte(i + 128) // Non-ASCII
				}
				for i := 20; i < 100; i++ {
					data[i] = 'A' // Printable
				}
				return data
			}(),
			expected: false,
			desc:     "low ratio of non-printable chars should be text",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.name)
			err = os.WriteFile(testFile, tc.content, 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.name, result)
			}
		})
	}
}

func BenchmarkBinaryRule_ExtensionOnly(b *testing.B) {
	rule := NewBinaryRule()
	tmpDir, err := os.MkdirTemp("", "bench_binary")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file with binary extension
	testFile := filepath.Join(tmpDir, "test.jpg")
	err = os.WriteFile(testFile, []byte("fake jpeg content"), 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(testFile)
	}
}

func BenchmarkBinaryRule_ContentAnalysis(b *testing.B) {
	rule := NewBinaryRule()
	tmpDir, err := os.MkdirTemp("", "bench_binary_content")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file without binary extension - forces content analysis
	testFile := filepath.Join(tmpDir, "test.unknown")
	content := make([]byte, 1024)
	for i := range content {
		content[i] = 'A' // Text content
	}
	err = os.WriteFile(testFile, content, 0644)
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(testFile)
	}
}

func TestBinaryRule_ComprehensiveExtensions(t *testing.T) {
	rule := NewBinaryRule()

	// Test comprehensive binary extensions
	binaryFiles := []struct {
		filename string
		desc     string
	}{
		// Executables and libraries
		{"program.exe", "Windows executable"},
		{"library.dll", "Windows dynamic library"},
		{"libtest.so", "Linux shared object"},
		{"framework.dylib", "macOS dynamic library"},
		{"binary.bin", "generic binary"},
		{"object.obj", "object file"},
		{"archive.a", "static library"},
		{"java.class", "Java class file"},
		{"application.jar", "Java archive"},
		{"webapp.war", "Web application archive"},

		// Archives and compressed files
		{"data.zip", "ZIP archive"},
		{"backup.tar", "TAR archive"},
		{"compressed.gz", "gzip file"},
		{"archive.7z", "7-Zip archive"},
		{"old.rar", "RAR archive"},
		{"data.bz2", "bzip2 file"},
		{"file.xz", "XZ compressed"},
		{"package.tgz", "compressed tar"},
		{"backup.tbz", "bzip2 tar"},
		{"archive.tbz2", "bzip2 tar"},
		{"data.lz", "lzip file"},
		{"compressed.lzma", "LZMA file"},

		// Documents
		{"document.pdf", "PDF document"},
		{"report.doc", "Word document"},
		{"report.docx", "Word 2007+ document"},
		{"text.odt", "OpenDocument text"},
		{"spreadsheet.xls", "Excel spreadsheet"},
		{"data.xlsx", "Excel 2007+ spreadsheet"},
		{"calc.ods", "OpenDocument spreadsheet"},
		{"slides.ppt", "PowerPoint presentation"},
		{"presentation.pptx", "PowerPoint 2007+ presentation"},
		{"slides.odp", "OpenDocument presentation"},
		{"document.rtf", "Rich Text Format"},
		{"document.pages", "Apple Pages"},
		{"spreadsheet.numbers", "Apple Numbers"},
		{"presentation.key", "Apple Keynote"},

		// Images
		{"photo.jpg", "JPEG image"},
		{"picture.jpeg", "JPEG image"},
		{"image.png", "PNG image"},
		{"animation.gif", "GIF image"},
		{"bitmap.bmp", "Bitmap image"},
		{"image.tiff", "TIFF image"},
		{"photo.tif", "TIFF image"},
		{"web.webp", "WebP image"},
		{"icon.ico", "Icon file"},
		{"vector.svg", "SVG image"},
		{"design.psd", "Photoshop file"},
		{"vector.ai", "Adobe Illustrator"},
		{"print.eps", "Encapsulated PostScript"},
		{"photo.raw", "RAW image"},
		{"canon.cr2", "Canon RAW"},
		{"nikon.nef", "Nikon RAW"},

		// Audio and Video
		{"song.mp3", "MP3 audio"},
		{"audio.wav", "WAV audio"},
		{"music.flac", "FLAC audio"},
		{"audio.aac", "AAC audio"},
		{"sound.ogg", "OGG audio"},
		{"audio.wma", "WMA audio"},
		{"music.m4a", "M4A audio"},
		{"video.mp4", "MP4 video"},
		{"movie.avi", "AVI video"},
		{"video.mkv", "Matroska video"},
		{"clip.mov", "QuickTime video"},
		{"video.wmv", "Windows Media Video"},
		{"flash.flv", "Flash video"},
		{"web.webm", "WebM video"},
		{"video.m4v", "M4V video"},

		// Databases
		{"data.db", "database file"},
		{"app.sqlite", "SQLite database"},
		{"cache.sqlite3", "SQLite 3 database"},
		{"legacy.mdb", "Access database"},
		{"new.accdb", "Access 2007+ database"},
		{"table.dbf", "dBASE file"},

		// Fonts
		{"font.ttf", "TrueType font"},
		{"font.otf", "OpenType font"},
		{"web.woff", "Web font"},
		{"web.woff2", "Web font 2"},
		{"legacy.eot", "Embedded OpenType"},

		// System and other binary formats
		{"disk.iso", "ISO image"},
		{"installer.dmg", "macOS disk image"},
		{"floppy.img", "disk image"},
		{"package.deb", "Debian package"},
		{"package.rpm", "RPM package"},
		{"installer.msi", "Windows installer"},
		{"package.pkg", "macOS package"},
		{"application.app", "macOS application"},
		{"module.pyc", "Python bytecode"},
		{"optimized.pyo", "Python optimized"},
		{"extension.pyd", "Python extension"},
	}

	tmpDir, err := os.MkdirTemp("", "binary_comprehensive_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tc := range binaryFiles {
		t.Run(tc.desc, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.filename)
			err = os.WriteFile(testFile, []byte("test content"), 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if !result {
				t.Errorf("File %s (%s) should be detected as binary", tc.filename, tc.desc)
			}
		})
	}
}

func TestBinaryRule_TextFiles(t *testing.T) {
	rule := NewBinaryRule()

	// Test files that should NOT be detected as binary
	textFiles := []struct {
		filename string
		content  []byte
		desc     string
	}{
		{"script.py", []byte("#!/usr/bin/env python\nprint('hello')"), "Python script"},
		{"source.go", []byte("package main\n\nfunc main() {\n\tfmt.Println(\"Hello\")\n}"), "Go source"},
		{"app.js", []byte("console.log('Hello, World!');"), "JavaScript"},
		{"style.css", []byte("body { margin: 0; padding: 0; }"), "CSS stylesheet"},
		{"page.html", []byte("<!DOCTYPE html><html><head></head><body></body></html>"), "HTML document"},
		{"data.json", []byte("{\"name\": \"test\", \"value\": 42}"), "JSON data"},
		{"config.yaml", []byte("database:\n  host: localhost\n  port: 5432"), "YAML config"},
		{"README.md", []byte("# Project\n\nThis is a readme file."), "Markdown"},
		{"data.txt", []byte("Plain text file content\nwith multiple lines"), "Text file"},
		{"script.sh", []byte("#!/bin/bash\necho \"Hello World\""), "Shell script"},
		{"Makefile", []byte("all:\n\techo \"Building...\""), "Makefile"},
		{"Dockerfile", []byte("FROM alpine:latest\nRUN apk add --no-cache curl"), "Dockerfile"},
		{"config.xml", []byte("<?xml version=\"1.0\"?>\n<config><item>value</item></config>"), "XML config"},
		{"data.csv", []byte("name,age,city\nJohn,30,NYC\nJane,25,LA"), "CSV data"},
		{"requirements.txt", []byte("requests>=2.25.1\nflask==2.0.1"), "Python requirements"},
		{"package.json", []byte("{\"name\": \"app\", \"version\": \"1.0.0\"}"), "Node.js package"},
		{"go.mod", []byte("module example.com/app\n\ngo 1.19"), "Go module"},
		{"LICENSE", []byte("MIT License\n\nCopyright (c) 2023"), "License file"},
		{".gitignore", []byte("*.log\nnode_modules/\n.env"), "Git ignore"},
		{"changelog.rst", []byte("Changelog\n=========\n\nVersion 1.0\n-----------"), "reStructuredText"},
	}

	tmpDir, err := os.MkdirTemp("", "text_files_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	for _, tc := range textFiles {
		t.Run(tc.desc, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.filename)
			err = os.WriteFile(testFile, tc.content, 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if result {
				t.Errorf("File %s (%s) should NOT be detected as binary", tc.filename, tc.desc)
			}
		})
	}
}

func TestBinaryRule_EdgeCaseSizes(t *testing.T) {
	rule := NewBinaryRule()
	tmpDir, err := os.MkdirTemp("", "size_edge_cases")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name     string
		size     int
		content  func(size int) []byte
		expected bool
		desc     string
	}{
		{
			name: "exactly_10mb.unknown",
			size: 10 * 1024 * 1024, // Exactly 10MB
			content: func(size int) []byte {
				data := make([]byte, size)
				for i := range data {
					data[i] = 'A'
				}
				return data
			},
			expected: false,
			desc:     "exactly 10MB text should not be binary",
		},
		{
			name: "just_over_10mb.unknown",
			size: 10*1024*1024 + 1, // Just over 10MB
			content: func(size int) []byte {
				data := make([]byte, size)
				for i := range data {
					data[i] = 'A'
				}
				return data
			},
			expected: true,
			desc:     "just over 10MB should be binary due to size",
		},
		{
			name: "very_large.unknown",
			size: 50 * 1024 * 1024, // 50MB
			content: func(size int) []byte {
				data := make([]byte, size)
				for i := range data {
					data[i] = 'A'
				}
				return data
			},
			expected: true,
			desc:     "very large file should be binary",
		},
		{
			name: "small_binary_content.unknown",
			size: 1024,
			content: func(size int) []byte {
				data := make([]byte, size)
				// Fill with binary content (null bytes)
				for i := 0; i < size/2; i++ {
					data[i] = 0 // null byte
				}
				for i := size / 2; i < size; i++ {
					data[i] = 'A'
				}
				return data
			},
			expected: true,
			desc:     "small file with binary content should be binary",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.name)
			content := tc.content(tc.size)

			err = os.WriteFile(testFile, content, 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.name, result)
			}
		})
	}
}

func TestBinaryRule_ContentEdgeCases(t *testing.T) {
	rule := NewBinaryRule()
	tmpDir, err := os.MkdirTemp("", "content_edge_cases")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name     string
		content  []byte
		expected bool
		desc     string
	}{
		{
			name:     "utf8_text.unknown",
			content:  []byte("Hello 世界 🌍 Ñiño café résumé"),
			expected: true,
			desc:     "UTF-8 text with international chars detected as binary due to high non-ASCII content",
		},
		{
			name:     "control_chars.unknown",
			content:  []byte("Text with\ttab and\nnewline\rand\fform feed"),
			expected: false,
			desc:     "text with control characters should not be binary",
		},
		{
			name:     "single_null_byte.unknown",
			content:  []byte("normal text\x00more text"),
			expected: true,
			desc:     "single null byte should make it binary",
		},
		{
			name:     "multiple_null_bytes.unknown",
			content:  []byte("text\x00with\x00multiple\x00null\x00bytes"),
			expected: true,
			desc:     "multiple null bytes should make it binary",
		},
		{
			name: "high_ascii.unknown",
			content: func() []byte {
				data := make([]byte, 200)
				for i := 0; i < 50; i++ {
					data[i] = byte(200 + i%50) // High ASCII values
				}
				for i := 50; i < 200; i++ {
					data[i] = 'A' // Normal text
				}
				return data
			}(),
			expected: false,
			desc:     "25% high ASCII should not be binary (below 30% threshold)",
		},
		{
			name: "very_high_non_printable.unknown",
			content: func() []byte {
				data := make([]byte, 100)
				for i := 0; i < 70; i++ {
					data[i] = byte(200) // Non-printable
				}
				for i := 70; i < 100; i++ {
					data[i] = 'A' // Printable
				}
				return data
			}(),
			expected: true,
			desc:     "70% non-printable should be binary (above 30% threshold)",
		},
		{
			name:     "binary_signature.unknown",
			content:  []byte{0xFF, 0xD8, 0xFF, 0xE0}, // JPEG signature
			expected: true,
			desc:     "binary signature should be detected as binary",
		},
		{
			name:     "pdf_signature.unknown",
			content:  []byte("%PDF-1.4\n%âãÏÓ"),
			expected: true,
			desc:     "PDF with binary content should be detected",
		},
		{
			name:     "zip_signature.unknown",
			content:  []byte{'P', 'K', 0x03, 0x04}, // ZIP signature
			expected: true,
			desc:     "ZIP signature should be detected as binary",
		},
		{
			name: "just_under_threshold.unknown",
			content: func() []byte {
				data := make([]byte, 1000)
				// 29% non-printable (just under 30% threshold)
				for i := 0; i < 290; i++ {
					data[i] = byte(128 + i%100) // Non-printable
				}
				for i := 290; i < 1000; i++ {
					data[i] = 'A' // Printable
				}
				return data
			}(),
			expected: true,
			desc:     "29% non-printable still triggers binary detection",
		},
		{
			name: "just_over_threshold.unknown",
			content: func() []byte {
				data := make([]byte, 1000)
				// 31% non-printable (just over 30% threshold)
				for i := 0; i < 310; i++ {
					data[i] = byte(128 + i%100) // Non-printable
				}
				for i := 310; i < 1000; i++ {
					data[i] = 'A' // Printable
				}
				return data
			}(),
			expected: true,
			desc:     "just over 30% threshold should be binary",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			testFile := filepath.Join(tmpDir, tc.name)
			err = os.WriteFile(testFile, tc.content, 0644)
			if err != nil {
				t.Fatal(err)
			}

			result := rule.Match(testFile)
			if result != tc.expected {
				t.Errorf("Expected %v for %s, got %v", tc.expected, tc.name, result)
			}
		})
	}
}

func TestBinaryRule_FileSystemErrors(t *testing.T) {
	rule := NewBinaryRule()

	// Test non-existent file
	result := rule.Match("/path/that/does/not/exist")
	if result {
		t.Error("Non-existent file should not be detected as binary")
	}

	// Test directory (will cause read error)
	tmpDir, err := os.MkdirTemp("", "fs_errors_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	result = rule.Match(tmpDir)
	if result {
		t.Error("Directory should not be detected as binary")
	}
}
