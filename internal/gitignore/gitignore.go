package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	patterns []string
	Patterns []string // Exported for testing
}

func New(path string) (*GitIgnore, error) {
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

	return &GitIgnore{patterns: patterns}, nil
}

// matchExact checks if the pattern exactly matches either the base name or full path
func (gi *GitIgnore) matchExact(pattern, path, baseName string) bool {
	return pattern == baseName || pattern == path
}

// matchDirectory checks if a directory pattern matches the path
func (gi *GitIgnore) matchDirectory(pattern, path string) bool {
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
func (gi *GitIgnore) matchGlobPattern(pattern, path, baseName string) bool {
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
	if len(gi.patterns) == 0 {
		return false
	}

	baseName := filepath.Base(path)

	for _, pattern := range gi.patterns {
		// Try exact matches first
		if gi.matchExact(pattern, path, baseName) {
			return true
		}

		// Check directory patterns
		if gi.matchDirectory(pattern, path) {
			return true
		}

		// Handle glob patterns
		if strings.Contains(pattern, "*") && gi.matchGlobPattern(pattern, path, baseName) {
			return true
		}
	}
	return false
}
