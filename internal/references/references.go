package references

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Common patterns for finding references in code
var referencePatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?m)^import\s+["']([^"']+)["']`),                    // Go imports
	regexp.MustCompile(`(?m)^from\s+([^\s]+)\s+import`),                     // Python imports
	regexp.MustCompile(`require\s*\(?\s*["']([^"']+)["']`),                  // Node.js requires
	regexp.MustCompile(`import\s+.*?from\s+["']([^"']+)["']`),               // ES6 imports
	regexp.MustCompile(`@import\s+["']([^"']+)["']`),                        // CSS imports
	regexp.MustCompile(`#include\s+["<]([^>"']+)[>"]`),                      // C/C++ includes
	regexp.MustCompile(`source\s+["']([^"']+)["']`),                         // Shell source
	regexp.MustCompile(`href\s*=\s*["']([^"']+)["']`),                       // HTML links
	regexp.MustCompile(`src\s*=\s*["']([^"']+)["']`),                        // HTML sources
	regexp.MustCompile(`url\s*\(\s*["']?([^"'\)]+)["']?\s*\)`),             // CSS urls
	regexp.MustCompile(`\[.*?\]\(([^)\s]+)\)`),                              // Markdown links
}

// Common non-local prefixes that indicate external references
var nonLocalPrefixes = []string{
	"http://", "https://",
	"git://", "git+",
	"npm:", "pip:",
	"gem:", "mvn:",
}

// Common file extensions to try when resolving references
var commonExtensions = []string{
	".go", ".mod",
	".js", ".jsx", ".ts", ".tsx",
	".py", ".rb", ".php",
	".java", ".scala", ".kt",
	".c", ".cpp", ".h", ".hpp",
	".css", ".scss", ".less",
	".html", ".htm",
	".md", ".rst", ".txt",
	".json", ".yaml", ".yml", ".toml",
}

// ReferenceMap stores file reference information
type ReferenceMap struct {
	// Internal references within the project
	Internal map[string][]string
	// External references (packages, URLs, etc.)
	External map[string][]string
}

// NewReferenceMap creates a new ReferenceMap
func NewReferenceMap() *ReferenceMap {
	return &ReferenceMap{
		Internal: make(map[string][]string),
		External: make(map[string][]string),
	}
}


// ExtractFileReferences finds references to other files within the given content
func ExtractFileReferences(content, currentDir, rootDir string, allFiles []string) *ReferenceMap {
	refs := NewReferenceMap()

	for _, pattern := range referencePatterns {
		matches := pattern.FindAllStringSubmatch(content, -1)
		for _, match := range matches {
			if len(match) < 2 {
				continue
			}

			// Get the first non-empty capture group
			var ref string
			for i := 1; i < len(match); i++ {
				if match[i] != "" {
					ref = strings.TrimSpace(match[i])
					break
				}
			}

			if ref == "" || ref == "." || ref == ".." {
				continue
			}

			// Remove any query parameters or fragments
			if idx := strings.IndexAny(ref, "?#"); idx != -1 {
				ref = ref[:idx]
			}

			// Initialize maps if needed
			if refs.Internal[currentDir] == nil {
				refs.Internal[currentDir] = []string{}
			}
			if refs.External[currentDir] == nil {
				refs.External[currentDir] = []string{}
			}

			addReference(refs, ref, currentDir, rootDir, allFiles)
		}
	}

	return refs
}

func addReference(refs *ReferenceMap, ref, currentDir, rootDir string, allFiles []string) {
	// Check if it's external
	if isExternalReference(ref) {
		if _, ok := refs.External[currentDir]; !ok {
			refs.External[currentDir] = []string{}
		}
		refs.External[currentDir] = append(refs.External[currentDir], ref)
		return
	}

	// Try to resolve as internal reference
	resolved := resolveReference(ref, currentDir, rootDir, allFiles)
	if resolved != "" {
		if _, ok := refs.Internal[currentDir]; !ok {
			refs.Internal[currentDir] = []string{}
		}
		refs.Internal[currentDir] = append(refs.Internal[currentDir], resolved)
	} else {
		// Only add as external if it's not a relative path
		if !strings.HasPrefix(ref, "./") && !strings.HasPrefix(ref, "../") {
			if _, ok := refs.External[currentDir]; !ok {
				refs.External[currentDir] = []string{}
			}
			refs.External[currentDir] = append(refs.External[currentDir], ref)
		}
	}
}

func isExternalReference(ref string) bool {
	// Check against non-local prefixes
	for _, prefix := range nonLocalPrefixes {
		if strings.HasPrefix(ref, prefix) {
			return true
		}
	}

	// Check for URLs
	if strings.Contains(ref, "://") {
		return true
	}

	if strings.HasPrefix(ref, "@") || strings.HasPrefix(ref, "github.com/") || strings.HasPrefix(ref, "golang.org/") || strings.HasPrefix(ref, "gopkg.in/") {
		return true
	}

	return false
}

func resolveGoDirectory(ref, rootDir string, allFiles []string) string {
	// Try config.go in the directory
	candidate := filepath.Join(ref, "config.go")
	if matchFile(candidate, rootDir, allFiles) {
		if rel, err := filepath.Rel(rootDir, candidate); err == nil {
			return rel
		}
		return candidate
	}
	return ""
}

func fallbackUpDirectories(ref, currentDir, rootDir string, allFiles []string) string {
	dir := currentDir
	for {
		dir = filepath.Dir(dir)
		candidate := filepath.Join(dir, ref)
		if matchFile(candidate, rootDir, allFiles) {
			if rel, err := filepath.Rel(rootDir, candidate); err == nil {
				return rel
			}
			return candidate
		}
		if dir == "." || dir == rootDir {
			break
		}
	}
	return ""
}

func resolveReference(ref, currentDir, rootDir string, allFiles []string) string {
	// Clean and normalize the reference path
	ref = filepath.Clean(ref)

	// Handle relative paths
	if strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../") {
		ref = filepath.Join(currentDir, ref)
	}

	candidates := []string{
		// Try as-is first
		ref,
		// Try relative to current directory
		filepath.Join(currentDir, ref),
		// Try absolute path within project
		filepath.Join(rootDir, ref),
	}

	// If no extension provided, try common extensions
	if filepath.Ext(ref) == "" {
		withExtensions := []string{}
		for _, candidate := range candidates {
			for _, ext := range commonExtensions {
				withExtensions = append(withExtensions, candidate+ext)
			}
		}
		candidates = append(candidates, withExtensions...)

		// Try Go directory logic for imports without extension
		if strings.Contains(ref, "/") {
			if resolved := resolveGoDirectory(ref, rootDir, allFiles); resolved != "" {
				return resolved
			}
		}
	}

	// Try all candidates
	for _, candidate := range candidates {
		if matchFile(candidate, rootDir, allFiles) {
			// Return path relative to root
			if rel, err := filepath.Rel(rootDir, candidate); err == nil {
				return rel
			}
			return candidate
		}
	}

	// Try fallback up directories
	if resolved := fallbackUpDirectories(ref, currentDir, rootDir, allFiles); resolved != "" {
		return resolved
	}

	// Final fallback: try the file directly at the project root
	rootCandidate := filepath.Join(rootDir, filepath.Base(ref))
	if matchFile(rootCandidate, rootDir, allFiles) {
		if rel, err := filepath.Rel(rootDir, rootCandidate); err == nil {
			return rel
		}
		return rootCandidate
	}

	return ""
}

func matchFile(path, rootDir string, allFiles []string) bool {
	// Normalize path for comparison
	path = filepath.Clean(path)

	// Convert candidate path to a relative path before comparing
	rel, err := filepath.Rel(rootDir, path)
	if err == nil {
		path = rel
	}

	// Check if file exists in project
	for _, file := range allFiles {
		if filepath.Clean(file) == path {
			return true
		}
	}

	return false
}
