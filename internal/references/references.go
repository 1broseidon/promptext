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
                    ref = match[i]
                    break
                }
            }
            
            if ref == "" {
                continue
            }

            // Check if it's an external reference
            if isExternalReference(ref) {
                refs.External[currentDir] = append(
                    refs.External[currentDir], 
                    ref,
                )
                continue
            }

            // Try to resolve the reference
            resolved := resolveReference(ref, currentDir, rootDir, allFiles)
            if resolved != "" {
                refs.Internal[currentDir] = append(
                    refs.Internal[currentDir], 
                    resolved,
                )
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
    
    // Check for package names without paths
    if !strings.Contains(ref, "/") && !strings.Contains(ref, ".") {
        return true
    }
    
    return false
}

func resolveReference(ref, currentDir, rootDir string, allFiles []string) string {
    // Clean and normalize the reference path
    ref = filepath.Clean(ref)
    
    // Try direct match first
    candidate := filepath.Join(currentDir, ref)
    if matchFile(candidate, rootDir, allFiles) {
        return candidate
    }
    
    // If no extension provided, try common extensions
    if filepath.Ext(ref) == "" {
        for _, ext := range commonExtensions {
            candidateWithExt := candidate + ext
            if matchFile(candidateWithExt, rootDir, allFiles) {
                return candidateWithExt
            }
        }
    }
    
    // Try absolute path within project
    if matchFile(ref, rootDir, allFiles) {
        return ref
    }
    
    return ""
}

func matchFile(path, rootDir string, allFiles []string) bool {
    // Normalize path for comparison
    path = filepath.Clean(path)
    
    // Check if file exists in project
    for _, file := range allFiles {
        if filepath.Clean(file) == path {
            return true
        }
    }
    
    return false
}
