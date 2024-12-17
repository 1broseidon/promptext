package filter

import (
    "path/filepath"
    "strings"
)

// GitIgnoreFilter implements filtering based on gitignore patterns
type GitIgnoreFilter struct {
    gitIgnore *GitIgnore
}

// NewGitIgnoreFilter creates a new GitIgnore filter
func NewGitIgnoreFilter(gitIgnore *GitIgnore) *GitIgnoreFilter {
    return &GitIgnoreFilter{
        gitIgnore: gitIgnore,
    }
}

// Match checks if the path matches any gitignore patterns
func (gf *GitIgnoreFilter) Match(path string) (bool, error) {
    if gf.gitIgnore == nil {
        return false, nil
    }
    return gf.gitIgnore.ShouldIgnore(path), nil
}

// Priority returns the filter priority
func (gf *GitIgnoreFilter) Priority() int {
    return 100 // Highest priority
}

// ShouldInclude determines if the path should be included
func (gf *GitIgnoreFilter) ShouldInclude(path string) bool {
    if gf.gitIgnore == nil {
        return true
    }
    return !gf.gitIgnore.ShouldIgnore(path)
}
