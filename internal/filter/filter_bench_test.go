package filter

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// createTempFilterTest creates a comprehensive directory structure for filter testing
func createTempFilterTest(b *testing.B, numFiles int, complexity string) string {
	b.Helper()

	tempDir, err := os.MkdirTemp("", "filter_bench")
	if err != nil {
		b.Fatal(err)
	}

	// Create directory structure based on complexity
	var dirs []string
	var patterns []string

	switch complexity {
	case "simple":
		dirs = []string{"src", "test", "docs"}
		patterns = []string{"*.tmp", "*.log"}
	case "medium":
		dirs = []string{
			"src/main", "src/utils", "src/handlers",
			"test/unit", "test/integration",
			"docs/api", "docs/guides",
			"build", "dist", "node_modules/pkg1", "node_modules/pkg2",
		}
		patterns = []string{"*.tmp", "*.log", "node_modules/", "build/", "*.test"}
	case "complex":
		dirs = []string{
			"src/main/java/com/example", "src/main/resources",
			"src/test/java/com/example", "src/test/resources",
			"target/classes", "target/test-classes", "target/surefire-reports",
			"node_modules/pkg1/dist", "node_modules/pkg2/lib", "node_modules/pkg3/build",
			"vendor/github.com/pkg1", "vendor/github.com/pkg2",
			".git/objects", ".git/refs", ".git/hooks",
			"coverage/html", "coverage/lcov",
			"docs/api/v1", "docs/api/v2", "docs/guides/getting-started",
			"build/debug", "build/release", "build/tmp",
			"cache/build", "cache/test", "tmp/uploads",
		}
		patterns = []string{
			"node_modules/", "vendor/", ".git/", "target/", "build/", "dist/",
			"coverage/", "tmp/", "cache/", "*.tmp", "*.log", "*.test",
			"*.class", "*.o", "*.so", "*.exe", "*.dll",
		}
	default:
		dirs = []string{"src"}
		patterns = []string{"*.tmp"}
	}

	// Create directories
	for _, dir := range dirs {
		if err := os.MkdirAll(filepath.Join(tempDir, dir), 0755); err != nil {
			b.Fatal(err)
		}
	}

	// Create files distributed across directories
	filesPerDir := numFiles / len(dirs)
	remainder := numFiles % len(dirs)

	extensions := []string{".go", ".js", ".py", ".java", ".md", ".json", ".yaml", ".tmp", ".log", ".test", ".class"}
	fileCount := 0

	for i, dir := range dirs {
		dirFiles := filesPerDir
		if i < remainder {
			dirFiles++
		}

		for j := 0; j < dirFiles && fileCount < numFiles; j++ {
			ext := extensions[fileCount%len(extensions)]
			filename := filepath.Join(tempDir, dir, fmt.Sprintf("file_%d%s", j, ext))
			content := fmt.Sprintf("// File content for %s\n", filename)

			if err := os.WriteFile(filename, []byte(content), 0644); err != nil {
				b.Fatal(err)
			}
			fileCount++
		}
	}

	// Create .gitignore file if complex
	if complexity == "complex" {
		gitignoreContent := strings.Join(patterns, "\n")
		gitignorePath := filepath.Join(tempDir, ".gitignore")
		if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
			b.Fatal(err)
		}
	}

	return tempDir
}

// Benchmark filter creation with different complexities
func BenchmarkFilter_Creation_SimpleRules(b *testing.B) {
	benchmarkFilterCreation(b, "simple")
}

func BenchmarkFilter_Creation_MediumRules(b *testing.B) {
	benchmarkFilterCreation(b, "medium")
}

func BenchmarkFilter_Creation_ComplexRules(b *testing.B) {
	benchmarkFilterCreation(b, "complex")
}

func benchmarkFilterCreation(b *testing.B, complexity string) {
	var opts Options

	switch complexity {
	case "simple":
		opts = Options{
			Includes:        []string{".go", ".js"},
			Excludes:        []string{"*.tmp", "*.log"},
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}
	case "medium":
		opts = Options{
			Includes: []string{".go", ".js", ".py", ".java", ".md"},
			Excludes: []string{
				"*.tmp", "*.log", "node_modules/", "build/", "dist/",
				"*.test", "coverage/", "tmp/",
			},
			UseDefaultRules: true,
			UseGitIgnore:    true,
		}
	case "complex":
		opts = Options{
			Includes: []string{".go", ".js", ".py", ".java", ".rb", ".php", ".rs", ".md", ".yaml", ".json"},
			Excludes: []string{
				"node_modules/", "vendor/", ".git/", "target/", "build/", "dist/",
				"coverage/", "tmp/", "cache/", "*.tmp", "*.log", "*.test",
				"*.class", "*.o", "*.so", "*.exe", "*.dll", "*.jar", "*.war",
				"*.pyc", "*.pyo", "__pycache__/", ".pytest_cache/",
				".vscode/", ".idea/", "*.swp", "*.swo", "*~",
			},
			UseDefaultRules: true,
			UseGitIgnore:    true,
		}
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		filter := New(opts)
		if filter == nil {
			b.Fatal("Failed to create filter")
		}
	}
}

// Benchmark filtering performance on different file counts
func BenchmarkFilter_ShouldProcess_100Files(b *testing.B) {
	benchmarkFilterShouldProcess(b, 100, "simple")
}

func BenchmarkFilter_ShouldProcess_1000Files(b *testing.B) {
	benchmarkFilterShouldProcess(b, 1000, "medium")
}

func BenchmarkFilter_ShouldProcess_10000Files(b *testing.B) {
	benchmarkFilterShouldProcess(b, 10000, "complex")
}

func benchmarkFilterShouldProcess(b *testing.B, numFiles int, complexity string) {
	tempDir := createTempFilterTest(b, numFiles, complexity)
	defer os.RemoveAll(tempDir)

	// Collect all file paths
	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}

	// Create filter based on complexity
	var opts Options
	switch complexity {
	case "simple":
		opts = Options{
			Includes:        []string{".go", ".js"},
			Excludes:        []string{"*.tmp", "*.log"},
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}
	case "medium":
		opts = Options{
			Includes:        []string{".go", ".js", ".py", ".md"},
			Excludes:        []string{"*.tmp", "*.log", "node_modules/", "build/"},
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}
	case "complex":
		opts = Options{
			Includes: []string{".go", ".js", ".py", ".java", ".md", ".yaml"},
			Excludes: []string{
				"node_modules/", "vendor/", ".git/", "target/", "build/",
				"*.tmp", "*.log", "*.test", "*.class", "coverage/",
			},
			UseDefaultRules: true,
			UseGitIgnore:    false,
		}
	}

	filter := New(opts)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processedCount := 0
		for _, file := range files {
			if filter.ShouldProcess(file) {
				processedCount++
			}
		}
		if processedCount == 0 {
			b.Log("Warning: No files were processed")
		}
	}

	b.ReportMetric(float64(len(files)), "files_checked")
}

// Benchmark pattern matching performance
func BenchmarkFilter_PatternMatching_Simple(b *testing.B) {
	benchmarkPatternMatching(b, []string{"*.tmp", "*.log"}, 1000)
}

func BenchmarkFilter_PatternMatching_Complex(b *testing.B) {
	patterns := []string{
		"node_modules/*", "vendor/*", ".git/*", "target/*", "build/*",
		"*.tmp", "*.log", "*.test", "*.class", "*.o", "*.so",
		"coverage/*", "tmp/*", "cache/*", "__pycache__/*",
	}
	benchmarkPatternMatching(b, patterns, 1000)
}

func benchmarkPatternMatching(b *testing.B, patterns []string, numFiles int) {
	// Create test file paths
	testPaths := make([]string, numFiles)
	pathTemplates := []string{
		"src/main.go", "src/handler.js", "test/main_test.go", "docs/readme.md",
		"node_modules/pkg/index.js", "vendor/lib/file.go", ".git/config",
		"build/output.tmp", "tmp/temp.log", "cache/data.class",
		"coverage/report.html", "target/classes/App.class",
	}

	for i := 0; i < numFiles; i++ {
		template := pathTemplates[i%len(pathTemplates)]
		testPaths[i] = fmt.Sprintf("%s_%d", template, i)
	}

	opts := Options{
		Excludes:        patterns,
		UseDefaultRules: false,
		UseGitIgnore:    false,
	}

	filter := New(opts)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			_ = filter.ShouldProcess(path)
		}
	}
}

// Benchmark gitignore parsing performance
func BenchmarkFilter_GitIgnore_Parsing(b *testing.B) {
	// Create temporary .gitignore with varying complexity
	tempDir, err := os.MkdirTemp("", "gitignore_bench")
	if err != nil {
		b.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	gitignoreContent := `# Dependencies
node_modules/
vendor/
bower_components/

# Build outputs  
build/
dist/
target/
out/
bin/

# IDE files
.vscode/
.idea/
*.swp
*.swo
*~

# OS files
.DS_Store
Thumbs.db
*.tmp

# Logs
*.log
logs/

# Coverage
coverage/
.nyc_output/
*.lcov

# Cache
.cache/
.tmp/
tmp/

# Environment
.env
.env.local
.env.development.local
.env.test.local
.env.production.local

# Package manager
npm-debug.log*
yarn-debug.log*
yarn-error.log*
package-lock.json
yarn.lock

# Runtime
*.pid
*.seed
*.pid.lock
`

	gitignorePath := filepath.Join(tempDir, ".gitignore")
	if err := os.WriteFile(gitignorePath, []byte(gitignoreContent), 0644); err != nil {
		b.Fatal(err)
	}

	// Change to temp directory for relative path testing
	oldWd, err := os.Getwd()
	if err != nil {
		b.Fatal(err)
	}
	defer os.Chdir(oldWd)
	os.Chdir(tempDir)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		patterns, err := ParseGitIgnore(".")
		if err != nil {
			b.Fatal(err)
		}
		if len(patterns) == 0 {
			b.Fatal("No patterns parsed")
		}
	}
}

// Benchmark concurrent filtering
func BenchmarkFilter_Concurrent(b *testing.B) {
	tempDir := createTempFilterTest(b, 1000, "medium")
	defer os.RemoveAll(tempDir)

	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}

	opts := Options{
		Includes:        []string{".go", ".js", ".py"},
		Excludes:        []string{"*.tmp", "*.log", "node_modules/"},
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}

	filter := New(opts)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for _, file := range files {
				_ = filter.ShouldProcess(file)
			}
		}
	})
}

// Benchmark memory usage during filtering
func BenchmarkFilter_MemoryUsage(b *testing.B) {
	tempDir := createTempFilterTest(b, 5000, "complex")
	defer os.RemoveAll(tempDir)

	var files []string
	err := filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			rel, err := filepath.Rel(tempDir, path)
			if err != nil {
				return err
			}
			files = append(files, rel)
		}
		return nil
	})
	if err != nil {
		b.Fatal(err)
	}

	opts := Options{
		Includes: []string{".go", ".js", ".py", ".java", ".md", ".yaml"},
		Excludes: []string{
			"node_modules/", "vendor/", ".git/", "target/",
			"*.tmp", "*.log", "*.test", "*.class",
		},
		UseDefaultRules: true,
		UseGitIgnore:    true,
	}

	filter := New(opts)

	// Measure memory usage
	var m1, m2 runtime.MemStats
	runtime.GC()
	runtime.ReadMemStats(&m1)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processedCount := 0
		for _, file := range files {
			if filter.ShouldProcess(file) {
				processedCount++
			}
		}
	}

	runtime.ReadMemStats(&m2)

	// Report memory metrics
	b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/1024/1024, "MB_allocated")
	b.ReportMetric(float64(len(files)), "files_processed")
}

// Benchmark filter rule evaluation order
func BenchmarkFilter_RuleOrder_IncludeFirst(b *testing.B) {
	opts := Options{
		Includes:        []string{".go", ".js", ".py"},
		Excludes:        []string{"*.tmp", "*.log", "test/"},
		UseDefaultRules: false,
		UseGitIgnore:    false,
	}
	benchmarkRuleOrder(b, opts, "include_first")
}

func BenchmarkFilter_RuleOrder_ExcludeFirst(b *testing.B) {
	opts := Options{
		Excludes:        []string{"*.tmp", "*.log", "test/"},
		Includes:        []string{".go", ".js", ".py"},
		UseDefaultRules: false,
		UseGitIgnore:    false,
	}
	benchmarkRuleOrder(b, opts, "exclude_first")
}

func benchmarkRuleOrder(b *testing.B, opts Options, scenario string) {
	testPaths := []string{
		"main.go", "handler.js", "utils.py", "test.tmp", "debug.log",
		"test/main_test.go", "src/app.js", "lib/util.py",
		"build/output.tmp", "logs/error.log",
	}

	filter := New(opts)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			_ = filter.ShouldProcess(path)
		}
	}
}

// Benchmark file type detection
func BenchmarkFilter_GetFileType(b *testing.B) {
	testFiles := []string{
		"main.go", "handler.js", "utils.py", "README.md", "config.yaml",
		"main_test.go", "spec.js", "test_utils.py", "package.json",
		"Dockerfile", "docker-compose.yml", "Makefile", "go.mod",
		"index.html", "style.css", "app.ts", "component.jsx",
	}

	opts := Options{
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}
	filter := New(opts)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, file := range testFiles {
			fileType := GetFileType(file, filter)
			_ = fileType
		}
	}
}

// Benchmark realistic project filtering scenarios
func BenchmarkFilter_RealisticGoProject(b *testing.B) {
	projectFiles := []string{
		"main.go", "go.mod", "go.sum", "README.md", ".gitignore",
		"cmd/server/main.go", "cmd/client/main.go",
		"internal/handler/user.go", "internal/handler/auth.go",
		"internal/service/user.go", "internal/repository/user.go",
		"pkg/config/config.go", "pkg/logger/logger.go",
		"test/user_test.go", "test/auth_test.go", "test/integration_test.go",
		"config/app.yaml", "config/db.yaml",
		"scripts/build.sh", "scripts/deploy.sh",
		"docs/api.md", "docs/deployment.md",
		"vendor/github.com/pkg1/file.go", "vendor/github.com/pkg2/file.go",
		"build/main", "build/server", "tmp/test.log", "tmp/debug.tmp",
	}

	opts := Options{
		Includes:        []string{".go", ".md", ".yaml", ".sh"},
		Excludes:        []string{"vendor/", "build/", "tmp/", "*.tmp", "*.log"},
		UseDefaultRules: true,
		UseGitIgnore:    false,
	}

	filter := New(opts)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		processedCount := 0
		for _, file := range projectFiles {
			if filter.ShouldProcess(file) {
				processedCount++
			}
		}

		if i == 0 {
			b.ReportMetric(float64(processedCount), "files_processed")
			b.ReportMetric(float64(len(projectFiles)), "total_files")
		}
	}
}
