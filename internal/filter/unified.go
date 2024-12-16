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
	// 1. Check default ignore directories
	for _, dir := range uf.defaultIgnores {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") {
			return false
		}
	}

	// 2. Check gitignore patterns
	if uf.gitIgnore != nil && uf.gitIgnore.ShouldIgnore(path) {
		return false
	}

	// 3. Check exclude patterns
	for _, exclude := range uf.configExcludes {
		if matched, err := filepath.Match(exclude, filepath.Base(path)); err == nil && matched {
			return false
		}
		if strings.Contains(path, exclude) {
			return false
		}
	}

	// 4. Check extensions if specified
	if len(uf.allowedExtensions) > 0 {
		ext := filepath.Ext(path)
		for _, allowedExt := range uf.allowedExtensions {
			if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
				return true
			}
		}
		return false
	}

	return true
}
