package filter

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	Patterns []string // Exported for testing and external use
}

func NewGitIgnore(path string) (*GitIgnore, error) {
	file, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &GitIgnore{}, nil
		}
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pattern := strings.TrimSpace(scanner.Text())
		if pattern != "" && !strings.HasPrefix(pattern, "#") {
			patterns = append(patterns, pattern)
		}
	}

	return &GitIgnore{Patterns: patterns}, nil
}

// matchExact checks if the pattern exactly matches either the base name or full path,
// or if the pattern matches the start of the path for directory-like patterns
func (gi *GitIgnore) MatchExact(pattern, path, baseName string) bool {
	if pattern == baseName || pattern == path {
		return true
	}
	// Handle patterns that should match at the start of the path
	if !strings.HasSuffix(pattern, "/") {
		return strings.HasPrefix(path, pattern+"/")
	}
	return false
}

// matchDirectory checks if a directory pattern matches the path
func (gi *GitIgnore) MatchDirectory(pattern, path string) bool {
	if strings.HasSuffix(pattern, "/") {
		dirPattern := strings.TrimSuffix(pattern, "/")
		parts := strings.Split(path, string(filepath.Separator))
		for _, part := range parts {
			if part == dirPattern {
				return true
			}
		}
	}
	return false
}

// matchGlobPattern checks if a glob pattern matches any part of the path
func (gi *GitIgnore) MatchGlobPattern(pattern, path, baseName string) bool {
	// For patterns starting with *, try matching against base name first
	if strings.HasPrefix(pattern, "*") {
		if matched, err := filepath.Match(pattern, baseName); err == nil && matched {
			return true
		}
	}

	// Try matching against full path
	if matched, err := filepath.Match(pattern, path); err == nil && matched {
		return true
	}

	// Try matching against each path segment
	segments := strings.Split(path, string(filepath.Separator))
	for _, segment := range segments {
		if matched, err := filepath.Match(pattern, segment); err == nil && matched {
			return true
		}
	}

	// For patterns like .name*, try matching with the pattern as a suffix
	if strings.HasSuffix(pattern, "*") {
		prefix := strings.TrimSuffix(pattern, "*")
		if strings.HasPrefix(path, prefix) || strings.HasPrefix(baseName, prefix) {
			return true
		}
	}

	return false
}

func (gi *GitIgnore) ShouldIgnore(path string) bool {
	if len(gi.Patterns) == 0 {
		return false
	}

	baseName := filepath.Base(path)

	for _, pattern := range gi.Patterns {
		// Try exact matches first
		if gi.MatchExact(pattern, path, baseName) {
			return true
		}

		// Check directory patterns
		if gi.MatchDirectory(pattern, path) {
			return true
		}

		// Handle glob patterns
		if strings.Contains(pattern, "*") && gi.MatchGlobPattern(pattern, path, baseName) {
			return true
		}
	}
	return false
}
