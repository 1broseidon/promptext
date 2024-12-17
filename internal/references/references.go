package references

import (
	"path/filepath"
	"strings"
)

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

func isGoStdlibPackage(pkg string) bool {
	return !strings.Contains(pkg, ".") && !strings.Contains(pkg, "/") && pkg == strings.ToLower(pkg)
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

			// Handle Go import blocks
			if pattern == referencePatterns[1] && len(match) > 2 && match[2] != "" {
				importBlock := match[2]
				importLines := strings.Split(importBlock, "\n")
				for _, line := range importLines {
					line = strings.TrimSpace(line)
					if line == "" {
						continue
					}
					if strings.HasPrefix(line, "\"") && strings.HasSuffix(line, "\"") {
						ref := strings.Trim(line, "\"")
						if ref != "" && !isGoStdlibPackage(ref) {
							addReference(refs, ref, currentDir, rootDir, allFiles)
						}
					}
				}
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

			if ref == "" || ref == "." || ref == ".." || isGoStdlibPackage(ref) {
				continue
			}

			// Remove any query parameters or fragments
			if idx := strings.IndexAny(ref, "?#"); idx != -1 {
				ref = ref[:idx]
			}

			// Handle Go single-line imports
			if pattern == referencePatterns[0] && ref != "" {
				if isGoStdlibPackage(ref) {
					continue
				}
				addReference(refs, ref, currentDir, rootDir, allFiles)
				continue
			}

			// Handle Python relative imports
			if strings.HasPrefix(ref, ".") && !strings.Contains(ref, "/") {
				parts := strings.Split(ref, " ")
				modPath := strings.TrimSpace(parts[0])

				// Convert relative import path to filesystem path
				levels := strings.Count(modPath, ".") - 1
				targetDir := currentDir
				for i := 0; i < levels; i++ {
					targetDir = filepath.Dir(targetDir)
				}
				modPath = strings.TrimLeft(modPath, ".")
				if modPath != "" {
					modPath = filepath.Join(targetDir, strings.Replace(modPath, ".", "/", -1))
				} else {
					modPath = targetDir
				}

				// Try to resolve the module path
				resolved := resolveReference(modPath, currentDir, rootDir, allFiles)
				if resolved != "" {
					if _, ok := refs.Internal[currentDir]; !ok {
						refs.Internal[currentDir] = []string{}
					}
					refs.Internal[currentDir] = append(refs.Internal[currentDir], resolved)
					continue
				}
			}

			// Handle Python from ... import ...
			if pattern == referencePatterns[3] && len(match) > 2 {
				baseModule := match[1]
				importedNames := strings.TrimSpace(match[2])
				names := strings.Split(importedNames, ",")

				targetDir := currentDir
				if !strings.HasPrefix(baseModule, ".") {
					// For non-relative imports, use parent directory as base
					targetDir = filepath.Dir(currentDir)
				} else if strings.HasPrefix(baseModule, ".") {
					// For relative imports, calculate target directory
					levels := strings.Count(baseModule, ".")
					for i := 0; i < levels; i++ {
						targetDir = filepath.Dir(targetDir)
					}
					baseModule = strings.TrimLeft(baseModule, ".")
				}

				for _, name := range names {
					name = strings.TrimSpace(name)
					if name == "" || name == "*" {
						continue
					}

					modPath := filepath.Join(targetDir, strings.Replace(baseModule, ".", "/", -1))
					if baseModule != "" {
						modPath = filepath.Join(modPath, name)
					} else {
						modPath = filepath.Join(targetDir, name)
					}

					resolved := resolveReference(modPath, targetDir, rootDir, allFiles)
					if resolved != "" {
						if _, ok := refs.Internal[currentDir]; !ok {
							refs.Internal[currentDir] = []string{}
						}
						refs.Internal[currentDir] = append(refs.Internal[currentDir], resolved)
					} else {
						// Combine baseModule and name for external reference
						addReference(refs, baseModule, currentDir, rootDir, allFiles)
					}
				}
				continue
			}

			// Handle Python import ...
			if pattern == referencePatterns[2] {
				resolved := resolveReference(ref, currentDir, rootDir, allFiles)
				if resolved != "" {
					if _, ok := refs.Internal[currentDir]; !ok {
						refs.Internal[currentDir] = []string{}
					}
					refs.Internal[currentDir] = append(refs.Internal[currentDir], resolved)
					continue
				}
				if resolved == "" {
					addReference(refs, ref, currentDir, rootDir, allFiles)					
				}
				continue
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

	// Check for standard library packages and other external packages
	if !strings.Contains(ref, "/") && !strings.Contains(ref, ".") {
		// Standard library packages
		return true
	}
	if strings.HasPrefix(ref, "@") || strings.HasPrefix(ref, "github.com/") || strings.HasPrefix(ref, "golang.org/") {
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
