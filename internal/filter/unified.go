package filter

import (
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/gitignore"
)

// UnifiedFilter combines all filtering rules into a single structure
type UnifiedFilter struct {
	gitIgnore         *gitignore.GitIgnore
	configExcludes    []string
	allowedExtensions []string
	defaultIgnores    []string
	defaultIgnoreExts []string
}

// NewUnifiedFilter creates a new UnifiedFilter with all exclusion patterns
func NewUnifiedFilter(gitIgnore *gitignore.GitIgnore, extensions, excludes []string) *UnifiedFilter {
	return &UnifiedFilter{
		gitIgnore:         gitIgnore,
		configExcludes:    excludes,
		allowedExtensions: extensions,
		defaultIgnores:    DefaultIgnoreDirs,
		defaultIgnoreExts: DefaultIgnoreExtensions,
	}
}

// GetFileType determines the type of file based on its path and patterns
func (uf *UnifiedFilter) GetFileType(path string) string {
	ext := strings.ToLower(filepath.Ext(path))
	base := strings.ToLower(filepath.Base(path))

	// Check for tests
	if strings.Contains(path, "_test.") || strings.Contains(path, "test_") {
		return "test"
	}

	// Check for entry points
	for lang, patterns := range entryPointPatterns {
		for _, pattern := range patterns {
			if matched, _ := filepath.Match(pattern, base); matched {
				return "entry:" + lang
			}
		}
	}

	// Check for configs
	for _, pattern := range configPatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return "config"
		}
	}

	// Check for documentation
	for _, pattern := range docPatterns {
		if matched, _ := filepath.Match(pattern, base); matched {
			return "doc"
		}
	}

	return "source"
}

// ShouldProcess determines if a file should be processed based on all rules
func (uf *UnifiedFilter) ShouldProcess(path string) bool {
	// 1. Check all exclusion patterns first

	// Check default ignore directories
	for _, dir := range uf.defaultIgnores {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") || path == dir {
			return false
		}
	}

	// Check gitignore patterns
	if uf.gitIgnore != nil && uf.gitIgnore.ShouldIgnore(path) {
		return false
	}

	// Check exclude patterns from config
	for _, exclude := range uf.configExcludes {
		// Try exact match first
		if exclude == path {
			return false
		}

		// Try glob pattern match
		if matched, err := filepath.Match(exclude, filepath.Base(path)); err == nil && matched {
			return false
		}

		// Try path contains pattern
		if strings.Contains(path, exclude) {
			return false
		}
	}

	// 2. After exclusions, check extensions

	// Get file extension and check against default ignored extensions first
	ext := filepath.Ext(path)
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
