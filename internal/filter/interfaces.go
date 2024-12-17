package filter

// Matcher defines the base interface for pattern matching
type Matcher interface {
    // Match returns true if the path matches the pattern
    Match(path string) (bool, error)
}

// Filter extends Matcher with priority and filtering logic
type Filter interface {
    Matcher
    // Priority returns the filter's priority (higher runs first)
    Priority() int
    // ShouldInclude returns true if the path should be included
    ShouldInclude(path string) bool
}
