package gitignore

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

type GitIgnore struct {
	patterns []string
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

func (gi *GitIgnore) ShouldIgnore(path string) bool {
	if len(gi.patterns) == 0 {
		return false
	}

	// Get both the full path and base name for matching
	baseName := filepath.Base(path)

	for _, pattern := range gi.patterns {
		// Handle directory patterns
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.Contains(path, dirPattern) {
				return true
			}
			continue
		}

		// Try exact matches first
		if pattern == baseName || pattern == path {
			return true
		}

		// Handle glob patterns
		if strings.Contains(pattern, "*") {
			// For patterns starting with *, try matching against base name first
			if strings.HasPrefix(pattern, "*") {
				matched, err := filepath.Match(pattern, baseName)
				if err == nil && matched {
					return true
				}
			}

			// Try matching against full path
			matched, err := filepath.Match(pattern, path)
			if err == nil && matched {
				return true
			}

			// Try matching against each path segment
			segments := strings.Split(path, string(filepath.Separator))
			for _, segment := range segments {
				matched, err := filepath.Match(pattern, segment)
				if err == nil && matched {
					return true
				}
			}

			// For patterns like .aider*, try matching with the pattern as a suffix
			if strings.HasSuffix(pattern, "*") {
				prefix := strings.TrimSuffix(pattern, "*")
				if strings.HasPrefix(path, prefix) || strings.HasPrefix(baseName, prefix) {
					return true
				}
			}
		}
	}
	return false
}
