package rules

import (
    "path/filepath"
    "strings"
    "github.com/1broseidon/promptext/internal/filter"
)

type PatternRule struct {
    filter.BaseRule
    patterns []string
}

func NewPatternRule(patterns []string, action filter.RuleAction) *PatternRule {
    return &PatternRule{
        BaseRule: filter.BaseRule{action: action},
        patterns: patterns,
    }
}

func (r *PatternRule) Match(path string) bool {
    normalizedPath := filepath.ToSlash(path)
    for _, pattern := range r.patterns {
        pattern = filepath.ToSlash(pattern)
        
        // Handle directory patterns
        if strings.HasSuffix(pattern, "/") {
            if strings.HasPrefix(normalizedPath, pattern) || 
               strings.Contains(normalizedPath, "/"+pattern) {
                return true
            }
            continue
        }
        
        // Handle exact matches and path-based patterns
        if strings.HasPrefix(normalizedPath, pattern) || 
           strings.Contains(normalizedPath, "/"+pattern) ||
           normalizedPath == pattern {
            return true
        }
    }
    return false
}
