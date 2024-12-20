package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"path/filepath"
	"strings"
)

type PatternRule struct {
	types.BaseRule
	patterns []string
}

func (r *PatternRule) Patterns() []string {
	return r.patterns
}

func NewPatternRule(patterns []string, action types.RuleAction) types.Rule {
	return &PatternRule{
		BaseRule: types.NewBaseRule("", action),
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

		// Handle wildcard patterns (e.g., .aider*)
		if strings.Contains(pattern, "*") {
			matched, _ := filepath.Match(pattern, filepath.Base(normalizedPath))
			if matched {
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
