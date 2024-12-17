package filter

import "path/filepath"

type Filter struct {
    rules []Rule
}

func New(rules ...Rule) *Filter {
    return &Filter{rules: rules}
}

// ShouldProcess determines if a path should be processed
func (f *Filter) ShouldProcess(path string) bool {
    path = filepath.Clean(path)
    
    for _, rule := range f.rules {
        if rule.Match(path) {
            return rule.Action() == Include
        }
    }
    
    return true
}

// IsExcluded checks if a path is explicitly excluded
func (f *Filter) IsExcluded(path string) bool {
    path = filepath.Clean(path)
    
    for _, rule := range f.rules {
        if rule.Match(path) && rule.Action() == Exclude {
            return true
        }
    }
    
    return false
}

// GetFileType determines the type of file based on its path
func GetFileType(path string, f *Filter) string {
    // First check if the path should be excluded
    if f != nil && f.IsExcluded(path) {
        return ""
    }

    // Check for test files
    if strings.Contains(path, "_test.go") || strings.Contains(path, "test_") || strings.HasPrefix(path, "test_") {
        return "test"
    }

    // Check for entry points
    base := filepath.Base(path)
    if base == "main.go" || base == "index.js" || base == "app.py" {
        return "entry:main"
    }

    // Check for config files
    switch filepath.Ext(path) {
    case ".yml", ".yaml", ".json", ".toml", ".ini", ".conf", ".config":
        return "config"
    }
    
    // Check for documentation
    switch filepath.Ext(path) {
    case ".md", ".txt", ".rst", ".adoc":
        return "doc"
    }
    
    // Default to empty string for other files
    return ""
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
    // Normalize both pattern and path
    pattern = filepath.ToSlash(pattern)
    path = filepath.ToSlash(path)

    // Handle directory patterns (ending with /)
    if strings.HasSuffix(pattern, "/") {
        pattern = strings.TrimSuffix(pattern, "/")
        return strings.HasPrefix(path, pattern+"/") || path == pattern
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
