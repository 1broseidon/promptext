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

	for _, pattern := range gi.patterns {
		matched, err := filepath.Match(pattern, filepath.Base(path))
		if err == nil && matched {
			return true
		}

		// Handle directory patterns
		if strings.HasSuffix(pattern, "/") {
			dirPattern := strings.TrimSuffix(pattern, "/")
			if strings.Contains(path, dirPattern) {
				return true
			}
		}
	}
	return false
}
