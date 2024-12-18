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

// GetFileType determines the type of file based on its path
func GetFileType(path string, f *Filter) string {
	// First check if the path should be excluded
	if f != nil && f.IsExcluded(path) {
		return ""
	}

	// Check for test files
	if strings.Contains(path, "_test.go") || strings.Contains(path, "test_") || strings.HasPrefix(path, "test_") {
		return "test"
	}

	// Check for entry points
	base := filepath.Base(path)
	if base == "main.go" || base == "index.js" || base == "app.py" {
		return "entry:main"
	}

	// Check for config files
	switch filepath.Ext(path) {
	case ".yml", ".yaml", ".json", ".toml", ".ini", ".conf", ".config":
		return "config"
	}

	// Check for documentation
	switch filepath.Ext(path) {
	case ".md", ".txt", ".rst", ".adoc":
		return "doc"
	}

	// Default to empty string for other files
	return ""
}
