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
        // Handle patterns with wildcards
        if strings.Contains(pattern, "*") || strings.Contains(pattern, "?") || strings.Contains(pattern, "[") {
            // Try matching against base name first for *.ext patterns
            if strings.HasPrefix(pattern, "*") {
                if matched, err := filepath.Match(pattern, filepath.Base(normalizedPath)); err == nil && matched {
                    return true, nil
                }
            }
            
            // Try matching against full path
            matched, err := filepath.Match(pattern, normalizedPath)
            if err != nil {
                return false, err
            }
            if matched {
                return true, nil
            }
        } else {
            // Exact match or directory prefix match
            if strings.HasSuffix(pattern, "/") {
                // For directory patterns, check if it's at the start or after a slash
                dirPattern := strings.TrimSuffix(pattern, "/")
                parts := strings.Split(normalizedPath, "/")
                for _, part := range parts {
                    if part == dirPattern {
                        return true, nil
                    }
                }
            } else if pattern == normalizedPath {
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
