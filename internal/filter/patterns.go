package filter

import (
    "path/filepath"
    "strings"
)

// Pattern represents a parsed path pattern
type Pattern struct {
    Original  string
    Negated   bool
    IsGlob    bool
    IsDir     bool
    Segments  []string
}

// NewPattern creates and validates a new pattern
func NewPattern(pattern string) (*Pattern, error) {
    p := &Pattern{
        Original: pattern,
        Negated:  strings.HasPrefix(pattern, "!"),
        IsGlob:   strings.Contains(pattern, "*") || strings.Contains(pattern, "?"),
        IsDir:    strings.HasSuffix(pattern, "/"),
    }
    
    // Remove negation prefix if present
    if p.Negated {
        pattern = pattern[1:]
    }
    
    // Remove trailing slash for directory patterns
    if p.IsDir {
        pattern = strings.TrimSuffix(pattern, "/")
    }
    
    // Split pattern into segments
    p.Segments = strings.Split(pattern, string(filepath.Separator))
    
    return p, nil
}

func (p *Pattern) Match(path string) bool {
    if p == nil {
        return false
    }
    
    // Handle directory-only patterns
    if p.IsDir && !strings.HasSuffix(path, "/") {
        path += "/"
    }
    
    // Split path into segments
    pathSegments := strings.Split(path, string(filepath.Separator))
    
    // Simple prefix/suffix matching for non-glob patterns
    if !p.IsGlob {
        patternPath := strings.Join(p.Segments, string(filepath.Separator))
        return strings.HasPrefix(path, patternPath) || strings.HasSuffix(path, patternPath)
    }
    
    // For glob patterns, try matching each segment
    return p.matchSegments(pathSegments)
}

func (p *Pattern) matchSegments(pathSegments []string) bool {
    if len(p.Segments) == 0 {
        return true
    }
    
    for i := 0; i <= len(pathSegments)-len(p.Segments); i++ {
        matched := true
        for j, segment := range p.Segments {
            if !matchSegment(segment, pathSegments[i+j]) {
                matched = false
                break
            }
        }
        if matched {
            return true
        }
    }
    return false
}

func matchSegment(pattern, segment string) bool {
    // Handle special glob characters
    if pattern == "*" {
        return true
    }
    
    matched, err := filepath.Match(pattern, segment)
    return err == nil && matched
}
