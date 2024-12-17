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

            // Try to resolve the reference
            resolved := resolveReference(ref, currentDir, rootDir, allFiles)
            if resolved != "" {
                // Convert to relative path for display
                if rel, err := filepath.Rel(rootDir, filepath.Join(rootDir, resolved)); err == nil {
                    refs.Internal[currentDir] = append(refs.Internal[currentDir], rel)
                }
            } else if isExternalReference(ref) {
                // Known external reference
                refs.External[currentDir] = append(refs.External[currentDir], ref)
            } else {
                // Attempt to resolve relative path as internal
                candidate := filepath.Join(currentDir, ref)
                if matchFile(candidate, rootDir, allFiles) {
                    if rel, err := filepath.Rel(rootDir, filepath.Join(rootDir, candidate)); err == nil {
                        refs.Internal[currentDir] = append(refs.Internal[currentDir], rel)
                        continue
                    }
                }
                // If still not found, it's external by default
                refs.External[currentDir] = append(refs.External[currentDir], ref)
            }
        }
    }

    return refs
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
    
    // Check for package names without paths or with '@' prefix (npm packages)
    if !strings.Contains(ref, "/") && !strings.Contains(ref, ".") || strings.HasPrefix(ref, "@") {
        return true
    }
    
    return false
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
