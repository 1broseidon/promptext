package filter

import "github.com/1broseidon/promptext/internal/gitignore"

// DefaultIgnoreDirs contains common directories that should be ignored
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

// ShouldProcessFile is maintained for backward compatibility
// Use UnifiedFilter.ShouldProcess instead for new code
func ShouldProcessFile(path string, extensions, excludes []string, gitIgnore *gitignore.GitIgnore) bool {
	filter := NewUnifiedFilter(gitIgnore, extensions, excludes)
	return filter.ShouldProcess(path)
}
