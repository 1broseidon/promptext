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
}

// NewUnifiedFilter creates a new UnifiedFilter with all exclusion patterns
func NewUnifiedFilter(gitIgnore *gitignore.GitIgnore, extensions, excludes []string) *UnifiedFilter {
	return &UnifiedFilter{
		gitIgnore:         gitIgnore,
		configExcludes:    excludes,
		allowedExtensions: extensions,
		defaultIgnores:    defaultIgnoreDirs,
	}
}

// ShouldProcess determines if a file should be processed based on all rules
func (uf *UnifiedFilter) ShouldProcess(path string) bool {
	// 1. Check default ignore directories first
	for _, dir := range uf.defaultIgnores {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") {
			return false
		}
	}

	// 2. Check gitignore patterns
	if uf.gitIgnore != nil && uf.gitIgnore.ShouldIgnore(path) {
		return false
	}

	// 3. Check exclude patterns from config
	for _, exclude := range uf.configExcludes {
		if matched, err := filepath.Match(exclude, filepath.Base(path)); err == nil && matched {
			return false
		}
		if strings.Contains(path, exclude) {
			return false
		}
	}

	// 4. If no extensions specified, include all non-excluded files
	if len(uf.allowedExtensions) == 0 {
		return true
	}

	// 5. Only include files with matching extensions
	ext := filepath.Ext(path)
	for _, allowedExt := range uf.allowedExtensions {
		// Normalize extensions for comparison
		if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
			return true
		}
	}

	// File extension doesn't match any allowed extensions
	return false
}
