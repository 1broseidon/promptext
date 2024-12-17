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

// FilterChain manages multiple filters
type FilterChain struct {
    filters []Filter
}

// NewFilterChain creates a new filter chain
func NewFilterChain() *FilterChain {
    return &FilterChain{
        filters: make([]Filter, 0),
    }
}

// Add appends a new filter to the chain
func (fc *FilterChain) Add(filter Filter) {
    fc.filters = append(fc.filters, filter)
    // Sort filters by priority (higher priority first)
    sortFilters(fc.filters)
}

// ShouldProcess checks if a path should be processed through all filters
func (fc *FilterChain) ShouldProcess(path string) bool {
    for _, filter := range fc.filters {
        if !filter.ShouldInclude(path) {
            return false
        }
    }
    return true
}

// sortFilters sorts filters by priority (higher first)
func sortFilters(filters []Filter) {
    for i := 0; i < len(filters)-1; i++ {
        for j := 0; j < len(filters)-i-1; j++ {
            if filters[j].Priority() < filters[j+1].Priority() {
                filters[j], filters[j+1] = filters[j+1], filters[j]
            }
        }
    }
}
