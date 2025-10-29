package token

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

// generateTextContent creates text content of specified size and type
func generateTextContent(size int, contentType string) string {
	var pattern string

	switch contentType {
	case "code":
		pattern = `package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})
	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

type User struct {
	ID       int    ` + "`json:\"id\"`" + `
	Name     string ` + "`json:\"name\"`" + `
	Email    string ` + "`json:\"email\"`" + `
	Password string ` + "`json:\"-\"`" + `
}

func (u *User) Validate() error {
	if u.Name == "" {
		return fmt.Errorf("name is required")
	}
	if u.Email == "" {
		return fmt.Errorf("email is required")
	}
	return nil
}
`
	case "markdown":
		pattern = `# Project Documentation

## Overview

This project implements a high-performance web service using Go. The service provides REST API endpoints for user management and authentication.

### Features

- **User Management**: Create, read, update, and delete users
- **Authentication**: JWT-based authentication system
- **Rate Limiting**: Request rate limiting per user
- **Monitoring**: Health checks and metrics endpoints

### Installation

` + "```bash" + `
go mod download
go build -o server ./cmd/server
./server
` + "```" + `

### API Endpoints

#### Users

- ` + "`GET /api/users`" + ` - List all users
- ` + "`POST /api/users`" + ` - Create a new user
- ` + "`GET /api/users/{id}`" + ` - Get user by ID
- ` + "`PUT /api/users/{id}`" + ` - Update user
- ` + "`DELETE /api/users/{id}`" + ` - Delete user

#### Authentication

- ` + "`POST /auth/login`" + ` - User login
- ` + "`POST /auth/logout`" + ` - User logout
- ` + "`POST /auth/refresh`" + ` - Refresh token

### Configuration

The service can be configured using environment variables or a YAML configuration file.
`
	case "json":
		pattern = `{
  "name": "example-service",
  "version": "1.0.0",
  "description": "A high-performance web service built with Go",
  "author": "Development Team",
  "license": "MIT",
  "dependencies": {
    "gin-gonic/gin": "v1.9.1",
    "golang-jwt/jwt": "v5.0.0",
    "gorm.io/gorm": "v1.25.4",
    "gorm.io/driver/postgres": "v1.5.2",
    "go-redis/redis": "v9.2.0",
    "stretchr/testify": "v1.8.4"
  },
  "scripts": {
    "build": "go build -o bin/server ./cmd/server",
    "test": "go test ./...",
    "run": "./bin/server",
    "docker:build": "docker build -t example-service .",
    "docker:run": "docker run -p 8080:8080 example-service"
  },
  "configuration": {
    "server": {
      "host": "localhost",
      "port": 8080,
      "timeout": 30
    },
    "database": {
      "host": "localhost",
      "port": 5432,
      "name": "service_db",
      "user": "postgres",
      "password": "secret"
    },
    "redis": {
      "host": "localhost",
      "port": 6379,
      "db": 0
    }
  }
}`
	case "plain":
		pattern = `This is plain text content for testing token counting performance.
The content contains regular English sentences with various punctuation marks.
We include numbers like 123, 456, and 789 to test numeric token handling.
Special characters are also included: !@#$%^&*()_+-=[]{}|;:,.<>?
This helps ensure comprehensive testing of the tokenization algorithm.

Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor 
incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis 
nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.

Technical terms and abbreviations: HTTP, REST, API, JWT, SQL, NoSQL, CPU, RAM, 
SSD, GPU, TCP/IP, DNS, URL, JSON, XML, YAML, CSV, PDF, HTML, CSS, JavaScript.
`
	case "mixed":
		// Combination of code, markdown, and plain text
		pattern = `# API Handler Implementation

` + "```go" + `
func (h *Handler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
	}
	
	if err := h.userService.Create(user); err != nil {
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	
	c.JSON(201, user)
}
` + "```" + `

This function handles user creation requests with proper validation and error handling.

**Key Features:**
- JSON request binding
- Input validation
- Service layer integration
- Appropriate HTTP status codes

Configuration example:
` + "```yaml" + `
server:
  port: 8080
  timeout: 30s
database:
  driver: postgres
  host: localhost
  port: 5432
` + "```" + `
`
	default:
		pattern = "Default text content for benchmarking token counting performance. "
	}

	// Repeat pattern to reach desired size
	if len(pattern) >= size {
		return pattern[:size]
	}

	repeats := size/len(pattern) + 1
	result := strings.Repeat(pattern, repeats)
	return result[:size]
}

// Benchmark token counting for different text sizes
func BenchmarkTokenCounter_SmallText(b *testing.B) {
	benchmarkTokenCounting(b, 1000, "code") // 1KB
}

func BenchmarkTokenCounter_MediumText(b *testing.B) {
	benchmarkTokenCounting(b, 10*1024, "code") // 10KB
}

func BenchmarkTokenCounter_LargeText(b *testing.B) {
	benchmarkTokenCounting(b, 100*1024, "code") // 100KB
}

func BenchmarkTokenCounter_HugeText(b *testing.B) {
	benchmarkTokenCounting(b, 1024*1024, "code") // 1MB
}

func benchmarkTokenCounting(b *testing.B, size int, contentType string) {
	content := generateTextContent(size, contentType)
	tokenCounter := NewTokenCounter()

	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	// Track memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	var totalTokens int
	for i := 0; i < b.N; i++ {
		tokens := tokenCounter.EstimateTokens(content)
		totalTokens = tokens // Prevent optimization
	}

	runtime.ReadMemStats(&m2)

	// Report metrics
	b.ReportMetric(float64(totalTokens), "tokens_counted")
	b.ReportMetric(float64(len(content))/1024, "content_kb")
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024/1024, "MB_allocated")
	b.ReportMetric(float64(totalTokens)/float64(len(content))*1000, "tokens_per_1k_chars")
}

// Benchmark different content types
func BenchmarkTokenCounter_CodeContent(b *testing.B) {
	benchmarkContentType(b, "code")
}

func BenchmarkTokenCounter_MarkdownContent(b *testing.B) {
	benchmarkContentType(b, "markdown")
}

func BenchmarkTokenCounter_JSONContent(b *testing.B) {
	benchmarkContentType(b, "json")
}

func BenchmarkTokenCounter_PlainContent(b *testing.B) {
	benchmarkContentType(b, "plain")
}

func BenchmarkTokenCounter_MixedContent(b *testing.B) {
	benchmarkContentType(b, "mixed")
}

func benchmarkContentType(b *testing.B, contentType string) {
	const size = 10 * 1024 // 10KB for each content type
	content := generateTextContent(size, contentType)
	tokenCounter := NewTokenCounter()

	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.ResetTimer()

	var totalTokens int
	for i := 0; i < b.N; i++ {
		totalTokens = tokenCounter.EstimateTokens(content)
	}

	b.ReportMetric(float64(totalTokens), "tokens_counted")
	b.ReportMetric(float64(totalTokens)/float64(len(content))*1000, "tokens_per_1k_chars")
}

// Benchmark many small files vs few large files
func BenchmarkTokenCounter_ManySmallTexts(b *testing.B) {
	const numTexts = 1000
	const textSize = 500 // 500 bytes each

	texts := make([]string, numTexts)
	for i := 0; i < numTexts; i++ {
		texts[i] = generateTextContent(textSize, "code")
	}

	tokenCounter := NewTokenCounter()
	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		totalTokens := 0
		for _, text := range texts {
			totalTokens += tokenCounter.EstimateTokens(text)
		}
		if totalTokens == 0 {
			b.Fatal("No tokens counted")
		}
	}

	b.ReportMetric(float64(numTexts), "texts_processed")
}

func BenchmarkTokenCounter_FewLargeTexts(b *testing.B) {
	const numTexts = 10
	const textSize = 50 * 1024 // 50KB each

	texts := make([]string, numTexts)
	for i := 0; i < numTexts; i++ {
		texts[i] = generateTextContent(textSize, "code")
	}

	tokenCounter := NewTokenCounter()
	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		totalTokens := 0
		for _, text := range texts {
			totalTokens += tokenCounter.EstimateTokens(text)
		}
		if totalTokens == 0 {
			b.Fatal("No tokens counted")
		}
	}

	b.ReportMetric(float64(numTexts), "texts_processed")
}

// Benchmark token counter initialization
func BenchmarkTokenCounter_Initialization(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		tc := NewTokenCounter()
		if tc == nil {
			b.Fatal("Failed to create token counter")
		}
	}
}

// Benchmark concurrent token counting
func BenchmarkTokenCounter_Concurrent(b *testing.B) {
	content := generateTextContent(10*1024, "mixed")
	tokenCounter := NewTokenCounter()

	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			tokens := tokenCounter.EstimateTokens(content)
			if tokens == 0 {
				b.Fatal("No tokens counted")
			}
		}
	})
}

// Benchmark empty and edge case content
func BenchmarkTokenCounter_EdgeCases(b *testing.B) {
	tokenCounter := NewTokenCounter()
	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	testCases := []struct {
		name    string
		content string
	}{
		{"empty", ""},
		{"single_char", "a"},
		{"whitespace", "   \n\t  "},
		{"special_chars", "!@#$%^&*()_+-=[]{}|;':\",./<>?"},
		{"unicode", "こんにちは 世界 🌍 emoji test"},
		{"numbers", "1234567890 42 3.14159 -123.45"},
		{"mixed_newlines", "line1\nline2\r\nline3\rline4"},
	}

	b.ResetTimer()

	for _, tc := range testCases {
		b.Run(tc.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				tokens := tokenCounter.EstimateTokens(tc.content)
				_ = tokens
			}
		})
	}
}

// Benchmark realistic codebase token counting scenario
func BenchmarkTokenCounter_RealisticCodebase(b *testing.B) {
	// Simulate a realistic codebase with different file types and sizes
	files := []struct {
		name    string
		content string
		size    int
	}{
		{"main.go", generateTextContent(3000, "code"), 3000},
		{"handler.go", generateTextContent(5000, "code"), 5000},
		{"model.go", generateTextContent(2000, "code"), 2000},
		{"README.md", generateTextContent(2500, "markdown"), 2500},
		{"config.json", generateTextContent(1000, "json"), 1000},
		{"test.go", generateTextContent(4000, "code"), 4000},
		{"api.md", generateTextContent(6000, "markdown"), 6000},
		{"package.json", generateTextContent(800, "json"), 800},
	}

	tokenCounter := NewTokenCounter()
	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		totalTokens := 0
		for _, file := range files {
			tokens := tokenCounter.EstimateTokens(file.content)
			totalTokens += tokens
		}

		if i == 0 {
			b.ReportMetric(float64(totalTokens), "total_tokens")
			b.ReportMetric(float64(len(files)), "files_processed")
		}
	}
}

// Benchmark token counting with memory profiling
func BenchmarkTokenCounter_MemoryProfile(b *testing.B) {
	sizes := []int{1024, 10 * 1024, 100 * 1024, 1024 * 1024} // 1KB to 1MB
	tokenCounter := NewTokenCounter()

	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	for _, size := range sizes {
		content := generateTextContent(size, "mixed")

		b.Run(fmt.Sprintf("size_%dKB", size/1024), func(b *testing.B) {
			var m1, m2 runtime.MemStats
			runtime.GC()
			runtime.ReadMemStats(&m1)

			b.ResetTimer()

			for i := 0; i < b.N; i++ {
				tokens := tokenCounter.EstimateTokens(content)
				_ = tokens
			}

			runtime.ReadMemStats(&m2)

			// Report memory metrics
			b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024, "KB_allocated")
			b.ReportMetric(float64(m2.Mallocs-m1.Mallocs), "malloc_count")
		})
	}
}

// Benchmark token counting performance vs simple character count
func BenchmarkTokenCounter_VsCharCount(b *testing.B) {
	content := generateTextContent(50*1024, "mixed") // 50KB mixed content
	tokenCounter := NewTokenCounter()

	if tokenCounter.encoding == nil {
		b.Skip("tiktoken encoding not available")
	}

	b.Run("TokenCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			tokens := tokenCounter.EstimateTokens(content)
			_ = tokens
		}
	})

	b.Run("CharCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			chars := len(content)
			_ = chars
		}
	})

	b.Run("RuneCount", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			runes := len([]rune(content))
			_ = runes
		}
	})
}
