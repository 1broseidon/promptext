package processor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	"github.com/1broseidon/promptext/internal/filter"
)

// createTempCodebase creates a temporary directory with specified number of files
func createTempCodebase(b *testing.B, numFiles int, avgFileSize int) string {
	b.Helper()
	
	tempDir, err := os.MkdirTemp("", "promptext_bench")
	if err != nil {
		b.Fatal(err)
	}
	
	// Create a realistic directory structure
	dirs := []string{
		"cmd/main",
		"internal/service",
		"internal/handler",
		"internal/model",
		"internal/repository",
		"pkg/utils",
		"pkg/config",
		"test/unit",
		"test/integration",
		"docs",
		"scripts",
	}
	
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tempDir, dir), 0755); err != nil {
			b.Fatal(err)
		}
	}
	
	// Generate sample Go code content
	sampleCode := generateSampleGoCode(avgFileSize)
	sampleTest := generateSampleTestCode(avgFileSize / 2)
	sampleConfig := generateSampleConfig(avgFileSize / 4)
	
	filesPerDir := numFiles / len(dirs)
	remainder := numFiles % len(dirs)
	
	fileCount := 0
	for i, dir := range dirs {
		dirFiles := filesPerDir
		if i < remainder {
			dirFiles++
		}
		
		for j := 0; j < dirFiles && fileCount < numFiles; j++ {
			var content string
			var ext string
			
			// Create different types of files based on directory
			switch {
			case strings.Contains(dir, "test"):
				content = sampleTest
				ext = ".go"
			case strings.Contains(dir, "config") || dir == "docs":
				content = sampleConfig
				if strings.Contains(dir, "config") {
					ext = ".yaml"
				} else {
					ext = ".md"
				}
			default:
				content = sampleCode
				ext = ".go"
			}
			
			filename := filepath.Join(tempDir, dir, fmt.Sprintf("file_%d%s", j, ext))
			if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
				b.Fatal(err)
			}
			fileCount++
		}
	}
	
	// Create some binary files to test filtering
	binaryFiles := []string{
		"binary.exe",
		"image.png",
		"archive.zip",
		"font.ttf",
	}
	
	for _, binFile := range binaryFiles {
		binaryContent := make([]byte, 1024)
		// Add some binary markers
		binaryContent[0] = 0x00
		binaryContent[1] = 0xFF
		binaryContent[100] = 0x00
		
		binPath := filepath.Join(tempDir, binFile)
		if err := os.WriteFile(binPath, binaryContent, 0644); err != nil {
			b.Fatal(err)
		}
	}
	
	return tempDir
}

func generateSampleGoCode(size int) string {
	template := `package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

type Server struct {
	addr string
	port int
}

func NewServer(addr string, port int) *Server {
	return &Server{
		addr: addr,
		port: port,
	}
}

func (s *Server) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("/api/users", s.usersHandler)
	
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.addr, s.port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	
	log.Printf("Starting server on %s:%d", s.addr, s.port)
	return server.ListenAndServe()
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "OK")
}

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getUsers(w, r)
	case http.MethodPost:
		s.createUser(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) getUsers(w http.ResponseWriter, r *http.Request) {
	// Implementation here
	fmt.Fprint(w, "[]")
}

func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	// Implementation here
	w.WriteHeader(http.StatusCreated)
}

func main() {
	server := NewServer("localhost", 8080)
	if err := server.Start(); err != nil {
		log.Fatal(err)
	}
}
`
	
	// Pad or trim to approximate target size
	if len(template) < size {
		padding := strings.Repeat("// Additional comment line for size padding\n", (size-len(template))/50)
		return template + padding
	}
	
	if len(template) > size && size > 100 {
		return template[:size-10] + "\n}\n"
	}
	
	return template
}

func generateSampleTestCode(size int) string {
	template := `package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestServer_NewServer(t *testing.T) {
	server := NewServer("localhost", 8080)
	assert.NotNil(t, server)
	assert.Equal(t, "localhost", server.addr)
	assert.Equal(t, 8080, server.port)
}

func TestServer_healthHandler(t *testing.T) {
	server := NewServer("localhost", 8080)
	// Test implementation here
	assert.NotNil(t, server)
}

func BenchmarkServer_Start(b *testing.B) {
	server := NewServer("localhost", 8080)
	for i := 0; i < b.N; i++ {
		// Benchmark implementation
		_ = server
	}
}
`
	
	if len(template) < size {
		padding := strings.Repeat("// Test comment padding\n", (size-len(template))/25)
		return template + padding
	}
	
	return template
}

func generateSampleConfig(size int) string {
	template := `# Application Configuration
server:
  host: localhost
  port: 8080
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: myapp
  user: postgres
  password: secret
  max_connections: 10

logging:
  level: info
  format: json
  output: stdout

features:
  enable_metrics: true
  enable_tracing: true
  cache_size: 1000
`
	
	if len(template) < size {
		padding := strings.Repeat("# Additional config comment\n", (size-len(template))/30)
		return template + padding
	}
	
	return template
}

// Benchmark file processing pipeline with different codebase sizes
func BenchmarkProcessDirectory_100Files(b *testing.B) {
	benchmarkProcessDirectory(b, 100, 2000)
}

func BenchmarkProcessDirectory_1000Files(b *testing.B) {
	benchmarkProcessDirectory(b, 1000, 2000)
}

func BenchmarkProcessDirectory_10000Files(b *testing.B) {
	benchmarkProcessDirectory(b, 10000, 2000)
}

func benchmarkProcessDirectory(b *testing.B, numFiles, avgFileSize int) {
	tempDir := createTempCodebase(b, numFiles, avgFileSize)
	defer os.RemoveAll(tempDir)
	
	// Create realistic filter configuration
	filterOpts := filter.Options{
		Includes:        []string{".go", ".yaml", ".md"},
		Excludes:        []string{"*.exe", "*.png", "*.zip", "*.ttf"},
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}
	
	f := filter.New(filterOpts)
	config := Config{
		DirPath:    tempDir,
		Extensions: []string{".go", ".yaml", ".md"},
		Excludes:   []string{"*.exe", "*.png", "*.zip", "*.ttf"},
		GitIgnore:  false,
		Filter:     f,
	}
	
	// Reset timer after setup
	b.ResetTimer()
	
	// Track memory usage
	var m1, m2 runtime.MemStats
	
	for i := 0; i < b.N; i++ {
		runtime.GC()
		runtime.ReadMemStats(&m1)
		
		result, err := ProcessDirectory(config, false)
		if err != nil {
			b.Fatal(err)
		}
		
		runtime.ReadMemStats(&m2)
		
		// Verify some processing happened
		if result.ProjectOutput == nil || len(result.ProjectOutput.Files) == 0 {
			b.Fatal("No files processed")
		}
		
		// Report memory usage for first iteration
		if i == 0 {
			b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024/1024, "MB/op")
			b.ReportMetric(float64(len(result.ProjectOutput.Files)), "files_processed")
			b.ReportMetric(float64(result.TokenCount), "total_tokens")
		}
	}
}

// Benchmark file processing with different file size distributions
func BenchmarkProcessDirectory_ManySmallFiles(b *testing.B) {
	benchmarkProcessDirectory(b, 5000, 500) // Many small files
}

func BenchmarkProcessDirectory_FewLargeFiles(b *testing.B) {
	benchmarkProcessDirectory(b, 100, 50000) // Few large files
}

// Benchmark individual components
func BenchmarkProcessFile_SingleFile(b *testing.B) {
	tempDir := createTempCodebase(b, 1, 2000)
	defer os.RemoveAll(tempDir)
	
	// Find the created file
	var testFile string
	err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && filepath.Ext(path) == ".go" {
			testFile = path
			return filepath.SkipAll
		}
		return nil
	})
	if err != nil || testFile == "" {
		b.Fatal("No test file found")
	}
	
	filterOpts := filter.Options{
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}
	f := filter.New(filterOpts)
	
	config := Config{
		DirPath: tempDir,
		Filter:  f,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fileInfo, err := processFile(testFile, config)
		if err != nil {
			b.Fatal(err)
		}
		if fileInfo == nil {
			b.Fatal("File was not processed")
		}
	}
}

// Benchmark walkdir performance vs processing
func BenchmarkWalkDir_DirectoryTraversal(b *testing.B) {
	tempDir := createTempCodebase(b, 1000, 2000)
	defer os.RemoveAll(tempDir)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		fileCount := 0
		err := filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				fileCount++
			}
			return nil
		})
		if err != nil {
			b.Fatal(err)
		}
		if fileCount == 0 {
			b.Fatal("No files found during traversal")
		}
	}
}

// Benchmark memory allocation patterns
func BenchmarkProcessDirectory_MemoryProfile(b *testing.B) {
	tempDir := createTempCodebase(b, 500, 2000)
	defer os.RemoveAll(tempDir)
	
	filterOpts := filter.Options{
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}
	f := filter.New(filterOpts)
	
	config := Config{
		DirPath: tempDir,
		Filter:  f,
	}
	
	// Force garbage collection before test
	runtime.GC()
	
	var m1, m2 runtime.MemStats
	runtime.ReadMemStats(&m1)
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result, err := ProcessDirectory(config, false)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
	
	runtime.ReadMemStats(&m2)
	
	// Report detailed memory metrics
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024/1024, "total_MB")
	b.ReportMetric(float64(m2.Mallocs-m1.Mallocs), "malloc_ops")
	b.ReportMetric(float64(m2.Sys-m1.Sys)/1024/1024, "sys_MB")
	b.ReportMetric(float64(m2.HeapAlloc-m1.HeapAlloc)/1024/1024, "heap_MB")
}

// Benchmark realistic codebase scenarios
func BenchmarkProcessDirectory_GoProject(b *testing.B) {
	// Simulate a typical Go project structure
	tempDir := createRealisticGoProject(b)
	defer os.RemoveAll(tempDir)
	
	filterOpts := filter.Options{
		Includes:        []string{".go", ".yaml", ".yml", ".md"},
		UseDefaultRules: true,
		UseGitIgnore:    true,
	}
	f := filter.New(filterOpts)
	
	config := Config{
		DirPath: tempDir,
		Filter:  f,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result, err := ProcessDirectory(config, false)
		if err != nil {
			b.Fatal(err)
		}
		
		if i == 0 {
			b.ReportMetric(float64(len(result.ProjectOutput.Files)), "files_processed")
			b.ReportMetric(float64(result.TokenCount), "total_tokens")
		}
	}
}

func createRealisticGoProject(b *testing.B) string {
	b.Helper()
	
	tempDir, err := os.MkdirTemp("", "go_project_bench")
	if err != nil {
		b.Fatal(err)
	}
	
	// Create realistic Go project structure
	structure := map[string]string{
		"main.go":                     generateSampleGoCode(3000),
		"go.mod":                      "module example.com/myapp\n\ngo 1.21\n",
		"README.md":                   "# My App\n\nThis is a sample application.",
		"Dockerfile":                  "FROM golang:1.21\nWORKDIR /app\nCOPY . .\n",
		".gitignore":                  "*.exe\n*.log\nvendor/\n.env\n",
		"cmd/server/main.go":          generateSampleGoCode(2000),
		"internal/handler/user.go":    generateSampleGoCode(4000),
		"internal/handler/auth.go":    generateSampleGoCode(3500),
		"internal/service/user.go":    generateSampleGoCode(5000),
		"internal/repository/user.go": generateSampleGoCode(3000),
		"pkg/config/config.go":        generateSampleGoCode(2000),
		"pkg/logger/logger.go":        generateSampleGoCode(1500),
		"test/user_test.go":           generateSampleTestCode(2000),
		"test/auth_test.go":           generateSampleTestCode(1800),
		"config/app.yaml":             generateSampleConfig(800),
		"docs/api.md":                 "# API Documentation\n\n## Endpoints\n\n### Users\n",
		"scripts/build.sh":            "#!/bin/bash\ngo build -o app ./cmd/server\n",
	}
	
	for filePath, content := range structure {
		fullPath := filepath.Join(tempDir, filePath)
		dir := filepath.Dir(fullPath)
		
		if err := os.MkdirAll(dir, 0755); err != nil {
			b.Fatal(err)
		}
		
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			b.Fatal(err)
		}
	}
	
	return tempDir
}

// Test performance against the existing sample projects
func BenchmarkProcessDirectory_SampleGoService(b *testing.B) {
	// Use the existing sample if available
	sampleDir := "../../../samples/go-service"
	if _, err := os.Stat(sampleDir); os.IsNotExist(err) {
		b.Skip("Sample go-service not found")
	}
	
	filterOpts := filter.Options{
		UseDefaultRules: true,
		UseGitIgnore:    true,
	}
	f := filter.New(filterOpts)
	
	config := Config{
		DirPath: sampleDir,
		Filter:  f,
	}
	
	b.ResetTimer()
	
	for i := 0; i < b.N; i++ {
		result, err := ProcessDirectory(config, false)
		if err != nil {
			b.Fatal(err)
		}
		_ = result
	}
}