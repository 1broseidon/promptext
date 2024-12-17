package filter

import (
	"path/filepath"
	"strings"
)

// UnifiedFilter provides backward compatibility with the old API
type UnifiedFilter struct {
	chain *FilterChain
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

	return &UnifiedFilter{chain: chain}
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
	for _, dir := range uf.defaultIgnores {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") || path == dir {
			return true
		}
	}
	return false
}

// matchesExcludePatterns checks if the path matches any exclude patterns
func (uf *UnifiedFilter) matchesExcludePatterns(path string) bool {
	for _, exclude := range uf.configExcludes {
		// Try exact match first
		if exclude == path {
			return true
		}

		// Try glob pattern match
		if matched, err := filepath.Match(exclude, filepath.Base(path)); err == nil && matched {
			return true
		}

		// Try path contains pattern
		if strings.Contains(path, exclude) {
			return true
		}
	}
	return false
}

// hasAllowedExtension checks if the file has an allowed extension
func (uf *UnifiedFilter) hasAllowedExtension(path string) bool {
	ext := filepath.Ext(path)

	// Check against default ignored extensions first
	for _, ignoreExt := range uf.defaultIgnoreExts {
		if strings.EqualFold(ignoreExt, ext) {
			return false
		}
	}

	// If no allowed extensions specified, include all non-excluded files
	if len(uf.allowedExtensions) == 0 {
		return true
	}

	// If allowed extensions specified, only include matching files
	for _, allowedExt := range uf.allowedExtensions {
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

// Helper function to parse gitignore patterns
func parseGitIgnorePatterns(patterns []string) []*Pattern {
	var result []*Pattern
	for _, p := range patterns {
		if pattern, err := NewPattern(p); err == nil {
			result = append(result, pattern)
		}
	}
	return result
}
