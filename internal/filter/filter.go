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
	if strings.Contains(path, "_test.go") || strings.Contains(path, "test_") ||
		strings.HasPrefix(path, "test_") || strings.HasSuffix(base, ".test.js") ||
		strings.HasSuffix(base, ".spec.js") || strings.HasSuffix(base, "_test.py") {
		info.IsTest = true
		info.Type = "test"
		info.Category = "test:" + strings.TrimPrefix(ext, ".")
		return info
	}

	// Check for entry points
	if base == "main.go" || base == "index.js" || base == "app.py" ||
		base == "index.ts" || base == "app.js" || base == "server.js" {
		info.IsEntryPoint = true
		info.Type = "source"
		info.Category = "entry:" + strings.TrimPrefix(ext, ".")
		return info
	}

	// Check for config files
	switch ext {
	case ".yml", ".yaml":
		info.Type = "config"
		info.Category = "config:yaml"
	case ".json":
		info.Type = "config"
		info.Category = "config:json"
	case ".toml":
		info.Type = "config"
		info.Category = "config:toml"
	case ".ini", ".conf", ".config":
		info.Type = "config"
		info.Category = "config:ini"
	}

	// Check for documentation
	switch ext {
	case ".md":
		info.Type = "doc"
		info.Category = "doc:markdown"
	case ".txt":
		info.Type = "doc"
		info.Category = "doc:text"
	case ".rst":
		info.Type = "doc"
		info.Category = "doc:rst"
	case ".adoc":
		info.Type = "doc"
		info.Category = "doc:asciidoc"
	}

	// Check for source code files
	switch ext {
	case ".go":
		info.Type = "source"
		info.Category = "source:go"
	case ".js":
		info.Type = "source"
		info.Category = "source:javascript"
	case ".ts":
		info.Type = "source"
		info.Category = "source:typescript"
	case ".jsx", ".tsx":
		info.Type = "source"
		info.Category = "source:react"
	case ".py":
		info.Type = "source"
		info.Category = "source:python"
	case ".java":
		info.Type = "source"
		info.Category = "source:java"
	case ".rb":
		info.Type = "source"
		info.Category = "source:ruby"
	case ".php":
		info.Type = "source"
		info.Category = "source:php"
	case ".rs":
		info.Type = "source"
		info.Category = "source:rust"
	}

	// Check for build and dependency files
	switch base {
	case "package.json", "package-lock.json", "yarn.lock":
		info.Type = "dependency"
		info.Category = "dep:node"
	case "go.mod", "go.sum":
		info.Type = "dependency"
		info.Category = "dep:go"
	case "requirements.txt", "Pipfile", "pyproject.toml":
		info.Type = "dependency"
		info.Category = "dep:python"
	case "Gemfile", "Gemfile.lock":
		info.Type = "dependency"
		info.Category = "dep:ruby"
	case "composer.json", "composer.lock":
		info.Type = "dependency"
		info.Category = "dep:php"
	case "Cargo.toml", "Cargo.lock":
		info.Type = "dependency"
		info.Category = "dep:rust"
	}

	// If no specific type was set but we have an extension, mark as source
	if info.Type == "" && ext != "" {
		info.Type = "source"
		info.Category = "source:other"
	}

	// If still no type, mark as unknown
	if info.Type == "" {
		info.Type = "unknown"
		info.Category = "unknown"
	}

	return info
}
