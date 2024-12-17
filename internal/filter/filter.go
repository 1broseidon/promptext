package filter


// DefaultIgnoreExtensions contains file extensions that should be ignored by default
var DefaultIgnoreExtensions = []string{
	// Images
	".jpg", ".jpeg", ".png", ".gif", ".bmp", ".tiff", ".webp",
	".ico", ".icns", ".svg", ".eps", ".raw", ".cr2", ".nef",
	// Binary/Executable
	".exe", ".dll", ".so", ".dylib", ".bin", ".obj",
	// Archives
	".zip", ".tar", ".gz", ".7z", ".rar", ".iso",
	// Other binary formats
	".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx",
	".class", ".pyc", ".pyo", ".pyd", ".o", ".a",
}

// DefaultIgnoreDirs contains common directories that should be ignored
var DefaultIgnoreDirs = []string{
	".git/",
	"node_modules/",
	"vendor/",
	".idea/",
	".vscode/",
	"__pycache__/",
	".pytest_cache/",
	"dist/",
	"build/",
	"coverage/",
	"bin/",
	".terraform/",
}

// ShouldProcessFile is maintained for backward compatibility
// Use UnifiedFilter.ShouldProcess instead for new code
func ShouldProcessFile(path string, extensions, excludes []string, gitIgnore *GitIgnore) bool {
	filter := NewUnifiedFilter(gitIgnore, extensions, excludes)
	return filter.ShouldProcess(path)
}
