package filter

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter/rules"
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/1broseidon/promptext/internal/log"
)

type Options struct {
	Includes        []string
	Excludes        []string
	UseDefaultRules bool // Controls whether to apply default filtering rules
	UseGitIgnore    bool
}

// ParseGitIgnore reads .gitignore file and returns patterns
func ParseGitIgnore(rootDir string) ([]string, error) {
	gitignorePath := filepath.Join(rootDir, ".gitignore")
	file, err := os.Open(gitignorePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}

	return patterns, scanner.Err()
}

// MergeAndDedupePatterns combines and deduplicates patterns
func MergeAndDedupePatterns(patterns ...[]string) []string {
	seen := make(map[string]bool)
	var result []string

	for _, patternSet := range patterns {
		for _, pattern := range patternSet {
			if !seen[pattern] {
				seen[pattern] = true
				result = append(result, pattern)
			}
		}
	}

	return result
}

type Filter struct {
	rules []types.Rule
}

func New(opts Options) *Filter {
	var filterRules []types.Rule
	var excludePatterns []string

	var defaultPatterns, gitPatterns, configPatterns []string

	log.Phase("Filter Configuration")

	// Collect patterns from all sources
	if opts.UseDefaultRules {
		defaultRules := rules.DefaultExcludes()
		for _, rule := range defaultRules {
			if patternRule, ok := rule.(*rules.PatternRule); ok {
				defaultPatterns = append(defaultPatterns, patternRule.Patterns()...)
			}
		}
		log.Debug("Default exclude patterns: %d", len(defaultPatterns))
	}

	if opts.UseGitIgnore {
		if patterns, err := ParseGitIgnore("."); err == nil && len(patterns) > 0 {
			gitPatterns = patterns
			log.Debug("Gitignore patterns: %d", len(gitPatterns))
		}
	}

	if len(opts.Excludes) > 0 {
		configPatterns = opts.Excludes
		log.Debug("Config exclude patterns: %d", len(configPatterns))
	}

	// Merge all patterns
	excludePatterns = MergeAndDedupePatterns([][]string{defaultPatterns, gitPatterns, configPatterns}...)

	// Log final consolidated patterns in array style
	if len(excludePatterns) > 0 {
		log.Debug("Final consolidated exclude patterns (%d): [%s]", len(excludePatterns), strings.Join(excludePatterns, ", "))
	}

	// Add default rules first if enabled
	if opts.UseDefaultRules {
		filterRules = append(filterRules, rules.DefaultExcludes()...)
	}

	// Add pattern-based rules
	if len(excludePatterns) > 0 {
		filterRules = append(filterRules,
			rules.NewPatternRule(excludePatterns, types.Exclude),
			rules.NewExtensionRule(excludePatterns, types.Exclude))
	}

	// Add include rules
	if len(opts.Includes) > 0 {
		filterRules = append(filterRules, rules.NewExtensionRule(opts.Includes, types.Include))
	}

	return &Filter{rules: filterRules}
}

// ShouldProcess determines if a path should be processed
func (f *Filter) ShouldProcess(path string) bool {
	path = filepath.Clean(path)

	// First check excludes silently
	if f.IsExcluded(path) {
		return false
	}

	// Check for binary files early
	for _, rule := range f.rules {
		if br, ok := rule.(*rules.BinaryRule); ok {
			if br.Match(path) {
				log.Debug("Skipping binary file: %s", path)
				return false
			}
		}
	}

	// Then check includes
	for _, rule := range f.rules {
		if rule.Match(path) && rule.Action() == types.Include {
			return true
		}
	}

	// If there are include rules but none matched, exclude silently
	for _, rule := range f.rules {
		if rule.Action() == types.Include {
			return false
		}
	}

	// No rules matched, default to include without logging
	return true
}

// IsExcluded checks if a path is explicitly excluded
func (f *Filter) IsExcluded(path string) bool {
	path = filepath.Clean(path)

	for _, rule := range f.rules {
		if rule.Match(path) && rule.Action() == types.Exclude {
			return true
		}
	}

	return false
}

// FileTypeInfo contains detailed information about a file's type and category
type FileTypeInfo struct {
	Type         string // Primary type (e.g., source, config, doc)
	Category     string // Specific category (e.g., go-source, yaml-config)
	IsTest       bool   // Whether it's a test file
	IsEntryPoint bool   // Whether it's an entry point
	Size         int64  // File size in bytes
}

// Helper functions to reduce cyclomatic complexity

func isTestFile(path, base string) bool {
	return strings.Contains(path, "_test.go") || 
		strings.Contains(path, "test_") ||
		strings.HasPrefix(path, "test_") || 
		strings.HasSuffix(base, ".test.js") ||
		strings.HasSuffix(base, ".spec.js") || 
		strings.HasSuffix(base, "_test.py")
}

func isEntryPoint(base string) bool {
	entryPoints := map[string]bool{
		"main.go":   true,
		"index.js":  true,
		"app.py":    true,
		"index.ts":  true,
		"app.js":    true,
		"server.js": true,
	}
	return entryPoints[base]
}

func getConfigType(ext string) (string, string) {
	configTypes := map[string]string{
		".yml":    "yaml",
		".yaml":   "yaml",
		".json":   "json",
		".toml":   "toml",
		".ini":    "ini",
		".conf":   "ini",
		".config": "ini",
	}
	if configType, ok := configTypes[ext]; ok {
		return "config", "config:" + configType
	}
	return "", ""
}

func getDocType(ext string) (string, string) {
	docTypes := map[string]string{
		".md":   "markdown",
		".txt":  "text",
		".rst":  "rst",
		".adoc": "asciidoc",
	}
	if docType, ok := docTypes[ext]; ok {
		return "doc", "doc:" + docType
	}
	return "", ""
}

func getSourceType(ext string) (string, string) {
	sourceTypes := map[string]string{
		".go":   "go",
		".js":   "javascript",
		".ts":   "typescript",
		".jsx":  "react",
		".tsx":  "react",
		".py":   "python",
		".java": "java",
		".rb":   "ruby",
		".php":  "php",
		".rs":   "rust",
	}
	if sourceType, ok := sourceTypes[ext]; ok {
		return "source", "source:" + sourceType
	}
	return "", ""
}

func getDependencyType(base string) (string, string) {
	depTypes := map[string]string{
		"package.json":      "node",
		"package-lock.json": "node",
		"yarn.lock":         "node",
		"go.mod":            "go",
		"go.sum":            "go",
		"requirements.txt":  "python",
		"Pipfile":           "python",
		"pyproject.toml":    "python",
		"Gemfile":           "ruby",
		"Gemfile.lock":      "ruby",
		"composer.json":     "php",
		"composer.lock":     "php",
		"Cargo.toml":        "rust",
		"Cargo.lock":        "rust",
	}
	if depType, ok := depTypes[base]; ok {
		return "dependency", "dep:" + depType
	}
	return "", ""
}

// GetFileType determines detailed type information for a file
func GetFileType(path string, f *Filter) FileTypeInfo {
	info := FileTypeInfo{}

	// First check if the path should be excluded
	if f != nil && f.IsExcluded(path) {
		return info
	}

	// Get file size if possible
	if stat, err := os.Stat(path); err == nil {
		info.Size = stat.Size()
	}

	base := filepath.Base(path)
	ext := filepath.Ext(path)

	// Check for test files
	if isTestFile(path, base) {
		info.IsTest = true
		info.Type = "test"
		info.Category = "test:" + strings.TrimPrefix(ext, ".")
		return info
	}

	// Check for entry points
	if isEntryPoint(base) {
		info.IsEntryPoint = true
		info.Type = "source"
		info.Category = "entry:" + strings.TrimPrefix(ext, ".")
		return info
	}

	// Check for config files
	if fileType, category := getConfigType(ext); fileType != "" {
		info.Type = fileType
		info.Category = category
		return info
	}

	// Check for documentation
	if fileType, category := getDocType(ext); fileType != "" {
		info.Type = fileType
		info.Category = category
		return info
	}

	// Check for source code files
	if fileType, category := getSourceType(ext); fileType != "" {
		info.Type = fileType
		info.Category = category
		return info
	}

	// Check for build and dependency files
	if fileType, category := getDependencyType(base); fileType != "" {
		info.Type = fileType
		info.Category = category
		return info
	}

	// If no specific type was set but we have an extension, mark as source
	if ext != "" {
		info.Type = "source"
		info.Category = "source:other"
	} else {
		// If still no type, mark as unknown
		info.Type = "unknown"
		info.Category = "unknown"
	}

	return info
}
