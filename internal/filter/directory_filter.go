package filter

import (
    "path/filepath"
    "strings"
)

// DirectoryFilter filters files based on directory patterns
type DirectoryFilter struct {
    patterns []string
    exclude  bool
}

// NewDirectoryFilter creates a new directory filter
func NewDirectoryFilter(patterns []string, exclude bool) *DirectoryFilter {
    return &DirectoryFilter{
        patterns: patterns,
        exclude:  exclude,
    }
}

// Match checks if the path matches any directory patterns
func (df *DirectoryFilter) Match(path string) (bool, error) {
    normalizedPath := filepath.ToSlash(path)
    
    for _, pattern := range df.patterns {
        // Exact match
        if pattern == normalizedPath {
            return true, nil
        }

        // Directory prefix match
        if strings.HasSuffix(pattern, "/") {
            dirPattern := strings.TrimSuffix(pattern, "/")
            if strings.HasPrefix(normalizedPath, dirPattern+"/") || normalizedPath == dirPattern {
                return true, nil
            }
        }

        // Path component match
        components := strings.Split(normalizedPath, "/")
        for _, component := range components {
            if component == pattern {
                return true, nil
            }
        }
    }
    
    return false, nil
}

// Priority returns the filter priority
func (df *DirectoryFilter) Priority() int {
    if df.exclude {
        return 90 // Higher priority for exclusion
    }
    return 30
}

// ShouldInclude determines if the path should be included
func (df *DirectoryFilter) ShouldInclude(path string) bool {
    matched, _ := df.Match(path)
    return matched != df.exclude
}
