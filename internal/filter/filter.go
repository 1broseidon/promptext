package filter

import (
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/gitignore"
)

var defaultIgnoreDirs = []string{
	".git",
	"node_modules",
	"vendor",
	".idea",
	".vscode",
	"__pycache__",
	".pytest_cache",
	"dist",
	"build",
	"coverage",
	"bin",
	".terraform",
}

func ShouldProcessFile(path string, extensions, excludes []string, gitIgnore *gitignore.GitIgnore) bool {
	// 1. Check default ignore directories first
	for _, dir := range defaultIgnoreDirs {
		if strings.Contains(path, "/"+dir+"/") || strings.HasPrefix(path, dir+"/") {
			return false
		}
	}

	// 2. Check gitignore patterns - this takes precedence over everything
	if gitIgnore != nil && gitIgnore.ShouldIgnore(path) {
		return false
	}

	// 3. Check exclude patterns - these also take precedence
	for _, exclude := range excludes {
		// Support both glob patterns and direct path contains
		if matched, err := filepath.Match(exclude, filepath.Base(path)); err == nil && matched {
			return false
		}
		if strings.Contains(path, exclude) {
			return false
		}
	}

	// 4. Only check extensions if file hasn't been excluded
	if len(extensions) == 0 {
		return true
	}

	ext := filepath.Ext(path)
	for _, allowedExt := range extensions {
		// Normalize extensions by trimming dots and comparing
		if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
			return true
		}
	}

	return false
}
