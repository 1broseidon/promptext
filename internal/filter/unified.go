package filter

import (
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/gitignore"
)

// UnifiedFilter combines all filtering rules into a single structure
type UnifiedFilter struct {
	gitIgnore           *gitignore.GitIgnore
	configExcludes      []string
	allowedExtensions   []string
	defaultIgnores      []string
	defaultIgnoreExts   []string
}

// NewUnifiedFilter creates a new UnifiedFilter with all exclusion patterns
func NewUnifiedFilter(gitIgnore *gitignore.GitIgnore, extensions, excludes []string) *UnifiedFilter {
	return &UnifiedFilter{
		gitIgnore:           gitIgnore,
		configExcludes:      excludes,
		allowedExtensions:   extensions,
		defaultIgnores:      DefaultIgnoreDirs,
		defaultIgnoreExts:   DefaultIgnoreExtensions,
	}
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

	// If no extensions specified, include all non-excluded files
	if len(uf.allowedExtensions) == 0 {
		return true
	}

	// Check default ignored extensions
	ext := filepath.Ext(path)
	for _, ignoreExt := range uf.defaultIgnoreExts {
		if strings.EqualFold(ignoreExt, ext) {
			return false
		}
	}

	// If allowed extensions specified, only include files with matching extensions
	if len(uf.allowedExtensions) > 0 {
		for _, allowedExt := range uf.allowedExtensions {
			// Normalize extensions for comparison
			if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
				return true
			}
		}
		return false
	}

	return false
}
