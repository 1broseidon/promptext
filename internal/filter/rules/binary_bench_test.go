package rules

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

// createTempBinaryFiles creates various types of test files for benchmarking
func createTempBinaryFiles(b *testing.B, numFiles int) string {
	b.Helper()
	
	tempDir, err := os.MkdirTemp("", "binary_bench")
	if err != nil {
		b.Fatal(err)
	}
	
	// Create different types of files for comprehensive testing
	fileTypes := []struct {
		ext     string
		binary  bool
		content func() []byte
	}{
		{".go", false, func() []byte { return []byte("package main\n\nfunc main() {\n\tprintln(\"Hello, World!\")\n}") }},
		{".js", false, func() []byte { return []byte("console.log('Hello, World!');") }},
		{".txt", false, func() []byte { return []byte("This is a plain text file with some content.") }},
		{".md", false, func() []byte { return []byte("# Markdown File\n\nThis is **bold** text.") }},
		{".yaml", false, func() []byte { return []byte("key: value\nlist:\n  - item1\n  - item2") }},
		{".exe", true, func() []byte { 
			content := make([]byte, 1024)
			content[0] = 0x4D  // 'M'
			content[1] = 0x5A  // 'Z' - PE header
			content[100] = 0x00 // null byte
			return content
		}},
		{".png", true, func() []byte {
			content := make([]byte, 512)
			content[0] = 0x89  // PNG signature
			content[1] = 0x50  // 'P'
			content[2] = 0x4E  // 'N'
			content[3] = 0x47  // 'G'
			content[50] = 0x00 // null byte
			return content
		}},
		{".zip", true, func() []byte {
			content := make([]byte, 256)
			content[0] = 0x50  // 'P' - ZIP signature
			content[1] = 0x4B  // 'K'
			content[2] = 0x03
			content[3] = 0x04
			content[25] = 0x00 // null byte
			return content
		}},
		{".pdf", true, func() []byte {
			content := make([]byte, 1024)
			copy(content, "%PDF-1.4") // PDF header
			content[100] = 0x00       // null byte
			return content
		}},
		{".so", true, func() []byte {
			content := make([]byte, 512)
			content[0] = 0x7F  // ELF signature
			content[1] = 0x45  // 'E'
			content[2] = 0x4C  // 'L'
			content[3] = 0x46  // 'F'
			content[50] = 0x00 // null byte
			return content
		}},
	}
	
	filesPerType := numFiles / len(fileTypes)
	remainder := numFiles % len(fileTypes)
	
	fileCount := 0
	for i, fileType := range fileTypes {
		typeFiles := filesPerType
		if i < remainder {
			typeFiles++
		}
		
		for j := 0; j < typeFiles && fileCount < numFiles; j++ {
			filename := filepath.Join(tempDir, fmt.Sprintf("file_%d_%d%s", i, j, fileType.ext))
			content := fileType.content()
			
			if err := os.WriteFile(filename, content, 0644); err != nil {
				b.Fatal(err)
			}
			fileCount++
		}
	}
	
	return tempDir
}

// createLargeFiles creates files of different sizes for size-based testing
func createLargeFiles(b *testing.B) string {
	b.Helper()
	
	tempDir, err := os.MkdirTemp("", "large_files_bench")
	if err != nil {
		b.Fatal(err)
	}
	
	// Create files of different sizes
	sizes := []struct {
		name string
		size int64
		ext  string
	}{
		{"small.txt", 1024, ".txt"},                    // 1KB text
		{"medium.log", 100 * 1024, ".log"},            // 100KB text  
		{"large.json", 1024 * 1024, ".json"},          // 1MB text
		{"huge_text.sql", 15 * 1024 * 1024, ".sql"},   // 15MB text (larger than threshold)
		{"small_binary.exe", 512, ".exe"},             // Small binary
		{"medium_binary.dll", 50 * 1024, ".dll"},      // Medium binary
		{"large_binary.so", 5 * 1024 * 1024, ".so"},   // Large binary
		{"huge_binary.bin", 50 * 1024 * 1024, ".bin"}, // Huge binary
	}
	
	for _, file := range sizes {
		path := filepath.Join(tempDir, file.name)
		
		var content []byte
		if filepath.Ext(file.name) == ".txt" || filepath.Ext(file.name) == ".log" ||
		   filepath.Ext(file.name) == ".json" || filepath.Ext(file.name) == ".sql" {
			// Text content - repeating pattern
			pattern := "This is sample text content for benchmarking binary detection performance. "
			needed := int(file.size) / len(pattern) + 1
			textContent := ""
			for i := 0; i < needed; i++ {
				textContent += pattern
			}
			content = []byte(textContent[:file.size])
		} else {
			// Binary content with null bytes and binary markers
			content = make([]byte, file.size)
			// Add binary signatures based on extension
			switch filepath.Ext(file.name) {
			case ".exe":
				content[0] = 0x4D // PE header
				content[1] = 0x5A
			case ".dll":
				content[0] = 0x4D // PE header
				content[1] = 0x5A
			case ".so":
				content[0] = 0x7F // ELF header
				content[1] = 0x45
				content[2] = 0x4C
				content[3] = 0x46
			}
			// Add null bytes throughout
			for i := int64(10); i < file.size; i += 100 {
				content[i] = 0x00
			}
		}
		
		if err := os.WriteFile(path, content, 0644); err != nil {
			b.Fatal(err)
		}
	}
	
	return tempDir
}

// Benchmark binary detection for different file counts
func BenchmarkBinaryRule_Match_100Files(b *testing.B) {
	benchmarkBinaryDetection(b, 100)
}

func BenchmarkBinaryRule_Match_1000Files(b *testing.B) {
	benchmarkBinaryDetection(b, 1000)
}

func BenchmarkBinaryRule_Match_10000Files(b *testing.B) {
	benchmarkBinaryDetection(b, 10000)
}

func benchmarkBinaryDetection(b *testing.B, numFiles int) {
	tempDir := createTempBinaryFiles(b, numFiles)
	defer os.RemoveAll(tempDir)
	
	// Collect all file paths
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	rule := NewBinaryRule()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		detectedBinary := 0
		for _, file := range files {
			if rule.Match(file) {
				detectedBinary++
			}
		}
		// Ensure some binary files were detected
		if detectedBinary == 0 {
			b.Fatal("No binary files detected")
		}
	}
	
	// Report files processed per operation
	b.ReportMetric(float64(len(files)), "files_per_op")
}

// Benchmark extension-only vs full binary detection
func BenchmarkBinaryRule_ExtensionCheck(b *testing.B) {
	tempDir := createTempBinaryFiles(b, 1000)
	defer os.RemoveAll(tempDir)
	
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		detectedBinary := 0
		for _, file := range files {
			ext := filepath.Ext(file)
			if binaryExtensions[ext] {
				detectedBinary++
			}
		}
	}
}

// Benchmark content analysis performance
func BenchmarkBinaryRule_ContentCheck(b *testing.B) {
	tempDir := createTempBinaryFiles(b, 100)
	defer os.RemoveAll(tempDir)
	
	// Get only files that don't have binary extensions (to force content analysis)
	var textFiles []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			ext := filepath.Ext(path)
			if !binaryExtensions[ext] {
				textFiles = append(textFiles, path)
			}
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	rule := &BinaryRule{}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, file := range textFiles {
			_ = rule.isBinaryContent(file)
		}
	}
}

// Benchmark large file handling
func BenchmarkBinaryRule_LargeFiles(b *testing.B) {
	tempDir := createLargeFiles(b)
	defer os.RemoveAll(tempDir)
	
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	rule := NewBinaryRule()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, file := range files {
			_ = rule.Match(file)
		}
	}
	
	runtime.ReadMemStats(&m2)
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024/1024, "MB_allocated")
}

// Compare performance with different buffer sizes for content analysis
func BenchmarkBinaryRule_BufferSize512(b *testing.B) {
	benchmarkContentBufferSize(b, 512)
}

func BenchmarkBinaryRule_BufferSize1024(b *testing.B) {
	benchmarkContentBufferSize(b, 1024)
}

func BenchmarkBinaryRule_BufferSize2048(b *testing.B) {
	benchmarkContentBufferSize(b, 2048)
}

func benchmarkContentBufferSize(b *testing.B, bufferSize int) {
	// Create a temporary file with mixed content
	tempDir, err := os.MkdirTemp("", "buffer_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create a file that will require content analysis
	testFile := filepath.Join(tempDir, "test.unknown")
	content := make([]byte, 4096) // 4KB file
	// Mix of text and some non-printable characters
	copy(content, "This is text content ")
	for i := 100; i < 4096; i += 200 {
		if i < len(content) {
			content[i] = byte(i % 128) // Various ASCII values
		}
	}
	
	if err := os.WriteFile(testFile, content, 0644); err != nil {
		b.Fatal(err)
	}
	
	// Custom rule with different buffer size
	rule := &BinaryRule{}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		_ = rule.isBinaryContent(testFile)
	}
}

// Benchmark edge cases
func BenchmarkBinaryRule_EmptyFiles(b *testing.B) {
	tempDir, err := os.MkdirTemp("", "empty_files_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)
	
	// Create multiple empty files
	var files []string
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("empty_%d.txt", i))
		if err := os.WriteFile(filename, []byte{}, 0644); err != nil {
			b.Fatal(err)
		}
		files = append(files, filename)
	}
	
	rule := NewBinaryRule()
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		for _, file := range files {
			_ = rule.Match(file)
		}
	}
}

// Benchmark concurrent binary detection
func BenchmarkBinaryRule_Concurrent(b *testing.B) {
	tempDir := createTempBinaryFiles(b, 1000)
	defer os.RemoveAll(tempDir)
	
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	rule := NewBinaryRule()
	
	b.ResetTimer()
	
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, file := range files {
				_ = rule.Match(file)
			}
		}
	})
}

// Benchmark binary detection vs old method (if there was one)
func BenchmarkBinaryRule_OptimizedVsNaive(b *testing.B) {
	tempDir := createTempBinaryFiles(b, 500)
	defer os.RemoveAll(tempDir)
	
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}
	
	optimizedRule := NewBinaryRule()
	
	// Naive implementation - always reads content
	naiveBinaryCheck := func(path string) bool {
		file, err := os.Open(path)
		if err != nil {
			return false
		}
		defer file.Close()
		
		buf := make([]byte, 1024)
		n, err := file.Read(buf)
		if err != nil {
			return false
		}
		
		// Check for null bytes
		for i := 0; i < n; i++ {
			if buf[i] == 0 {
				return true
			}
		}
		return false
	}
	
	b.Run("Optimized", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, file := range files {
				_ = optimizedRule.Match(file)
			}
		}
	})
	
	b.Run("Naive", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, file := range files {
				_ = naiveBinaryCheck(file)
			}
		}
	})
}