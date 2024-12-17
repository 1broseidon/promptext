package filter

import (
    "path/filepath"
    "strings"
)

// Filter represents a path filtering system
type Filter struct {
    includes []string // Patterns to explicitly include
    excludes []string // Patterns to exclude
    ignoreDefault bool // Whether to use default ignore patterns
}

// Options configures the Filter behavior
type Options struct {
    Includes []string
    Excludes []string
    IgnoreDefault bool
}

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

// New creates a new Filter with the given options
func New(opts Options) *Filter {
    return &Filter{
        includes: opts.Includes,
        excludes: append([]string{}, opts.Excludes...),
        ignoreDefault: opts.IgnoreDefault,
    }
}

// ShouldProcess determines if a path should be processed
func (f *Filter) ShouldProcess(path string) bool {
    // Check excludes first (excludes take precedence)
    if f.isExcluded(path) {
        return false
    }

    // If no includes specified, include everything not excluded
    if len(f.includes) == 0 {
        return true
    }

    // Check includes
    return f.isIncluded(path)
}

// isExcluded checks if a path matches any exclude patterns
func (f *Filter) isExcluded(path string) bool {
    normalizedPath := filepath.ToSlash(path)

    // Check default ignores if enabled
    if f.ignoreDefault {
        for _, pattern := range DefaultIgnoreDirs {
            if strings.Contains(normalizedPath, pattern) {
                return true
            }
        }
        
        ext := filepath.Ext(normalizedPath)
        for _, ignoreExt := range DefaultIgnoreExtensions {
            if strings.EqualFold(ignoreExt, ext) {
                return true
            }
        }
    }

    // Check custom excludes
    for _, pattern := range f.excludes {
        if f.matchPattern(pattern, normalizedPath) {
            return true
        }
    }

    return false
}

// isIncluded checks if a path matches any include patterns
func (f *Filter) isIncluded(path string) bool {
    normalizedPath := filepath.ToSlash(path)
    
    for _, pattern := range f.includes {
        if f.matchPattern(pattern, normalizedPath) {
            return true
        }
    }
    
    return false
}

// matchPattern checks if a path matches a pattern
func (f *Filter) matchPattern(pattern, path string) bool {
    // Handle directory patterns (ending with /)
    if strings.HasSuffix(pattern, "/") {
        return strings.Contains(path, pattern)
    }

    // Handle extension patterns (starting with .)
    if strings.HasPrefix(pattern, ".") {
        return strings.HasSuffix(path, pattern)
    }

    // Handle glob patterns
    if strings.Contains(pattern, "*") {
        matched, _ := filepath.Match(pattern, filepath.Base(path))
        return matched
    }

    // Exact match
    return path == pattern
}
