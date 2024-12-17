package filter

import (
	"path/filepath"
	"strings"
)

var (
	entryPointPatterns = map[string][]string{
		"go":     {"main.go"},
		"python": {"__main__.py", "main.py"},
		"node":   {"index.js", "server.js", "app.js"},
	}

	configPatterns = []string{
		"*.json", "*.yaml", "*.yml", "*.toml", "*.ini",
		".env*", "config.*", "*.config.*",
	}

	docPatterns = []string{
		"README*", "*.md", "docs/*", "*.txt",
		"LICENSE*", "CHANGELOG*", "CONTRIBUTING*",
	}
)

// UnifiedFilter provides backward compatibility with the old API
type UnifiedFilter struct {
	chain             *FilterChain
	DefaultIgnores    []string
	DefaultIgnoreExts []string
	AllowedExtensions []string
	ConfigExcludes    []string
	GitIgnore         *GitIgnore
}

// NewUnifiedFilter creates a new UnifiedFilter with all exclusion patterns
func NewUnifiedFilter(gitIgnore *GitIgnore, extensions, excludes []string) *UnifiedFilter {
	chain := NewFilterChain()

	// Add gitignore filter if provided
	if gitIgnore != nil {
		chain.Add(NewGitIgnoreFilter(gitIgnore))
	}

	// Add directory filter for default ignores
	chain.Add(NewDirectoryFilter(DefaultIgnoreDirs, true))

	// Add extension filters
	if len(extensions) > 0 {
		chain.Add(NewExtensionFilter(extensions, false)) // Include only these extensions
	}
	chain.Add(NewExtensionFilter(DefaultIgnoreExtensions, true)) // Exclude these extensions

	// Add directory filter for custom excludes
	if len(excludes) > 0 {
		chain.Add(NewDirectoryFilter(excludes, true))
	}

	return &UnifiedFilter{
		chain:             chain,
		DefaultIgnores:    DefaultIgnoreDirs,
		DefaultIgnoreExts: DefaultIgnoreExtensions,
		AllowedExtensions: extensions,
		ConfigExcludes:    excludes,
		GitIgnore:         gitIgnore,
	}
}

// GetFileType determines the type of file based on its path and patterns
func (uf *UnifiedFilter) GetFileType(path string) string {
	// Quick check for node_modules first
	if strings.Contains(path, "node_modules/") {
		return "dependency"
	}

	// Check for tests
	if strings.Contains(path, "_test.") || strings.Contains(path, "test_") {
		return "test"
	}

	// Check for entry points with full path support
	for lang, patterns := range entryPointPatterns {
		for _, pattern := range patterns {
			// Try matching against full path first
			if matched, _ := filepath.Match(pattern, path); matched {
				return "entry:" + lang
			}
			// Fall back to base name for simple patterns
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return "entry:" + lang
			}
		}
	}

	// Check for configs
	for _, pattern := range configPatterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return "config"
		}
	}

	// Check for documentation
	for _, pattern := range docPatterns {
		// Try full path first
		if matched, _ := filepath.Match(pattern, path); matched {
			return "doc"
		}
		// Then try base name
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return "doc"
		}
	}

	return "source"
}

// isInNodeModules checks if the path is within node_modules
func (uf *UnifiedFilter) isInNodeModules(path string) bool {
	return strings.Contains(path, "node_modules/")
}

// isInDefaultIgnoreDir checks if the path is in a default ignore directory
func (uf *UnifiedFilter) isInDefaultIgnoreDir(path string) bool {
	for _, dir := range uf.DefaultIgnores {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") || path == dir {
			return true
		}
	}
	return false
}

// matchesExcludePatterns checks if the path matches any exclude patterns
func (uf *UnifiedFilter) matchesExcludePatterns(path string) bool {
	for _, exclude := range uf.ConfigExcludes {
		// Try exact match first
		if exclude == path {
			return true
		}

		// Handle directory prefix patterns (ending with /)
		if strings.HasSuffix(exclude, "/") && strings.Contains(path, exclude) {
			return true
		}

		// Handle glob patterns
		if strings.Contains(exclude, "*") {
			baseName := filepath.Base(path)
			
			// For patterns like *.ext, prioritize matching against the base name
			if strings.HasPrefix(exclude, "*") {
				if matched, err := filepath.Match(exclude, baseName); err == nil && matched {
					return true
				}
			}
			
			// Try matching against full path
			if matched, err := filepath.Match(exclude, path); err == nil && matched {
				return true
			}
			
			// Try matching against each path component
			parts := strings.Split(path, "/")
			for _, part := range parts {
				if matched, err := filepath.Match(exclude, part); err == nil && matched {
					return true
				}
			}
		} else {
			// For non-glob patterns, check if the path contains the pattern
			if strings.Contains(path, exclude) {
				return true
			}
		}
	}
	return false
}

// hasAllowedExtension checks if the file has an allowed extension
func (uf *UnifiedFilter) hasAllowedExtension(path string) bool {
	ext := filepath.Ext(path)

	// Check against default ignored extensions first
	for _, ignoreExt := range uf.DefaultIgnoreExts {
		if strings.EqualFold(ignoreExt, ext) {
			return false
		}
	}

	// If no allowed extensions specified, include all non-excluded files
	if len(uf.AllowedExtensions) == 0 {
		return true
	}

	// If allowed extensions specified, only include matching files
	for _, allowedExt := range uf.AllowedExtensions {
		if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
			return true
		}
	}

	return false
}

// ShouldProcess determines if a file should be processed based on all rules
func (uf *UnifiedFilter) ShouldProcess(path string) bool {
	return uf.chain.ShouldProcess(path)
}

