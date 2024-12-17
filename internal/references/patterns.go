
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
