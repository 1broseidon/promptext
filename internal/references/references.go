package references

import (
    "path/filepath"
    "regexp"
    "strings"
)

var (
    // Common prefixes that indicate non-local references
    nonLocalPrefixes = []string{
        "http://", "https://", "mailto:", "tel:", "ftp://",
        "github.com/", "golang.org/", "gopkg.in/",
    }

    referencePatterns = []*regexp.Regexp{
        // Go imports (excluding standard library and external packages)
        regexp.MustCompile(`(?m)^import\s+(?:\([^)]+\)|["']\.\/([^"']+)["'])`),
        // require('./X'), Node.js style local imports
        regexp.MustCompile(`(?m)require\(['"]\./([^'"]+)['"]\)`),
        // Markdown links to local files: [text](./path/to/file)
        regexp.MustCompile(`(?m)\[[^\]]*\]\(((?!http|mailto|#)[^)]+)\)`),
        // HTML includes with relative paths
        regexp.MustCompile(`(?m)(?:src|href)=["'](\.\/[^'"]+)["']`),
    }
)

// ExtractFileReferences finds references to other files within the given content
func ExtractFileReferences(content, currentDir string, allFiles []string) []string {
    var refs []string
    refMap := make(map[string]bool)

    for _, pattern := range referencePatterns {
        matches := pattern.FindAllStringSubmatch(content, -1)
        for _, match := range matches {
            if len(match) < 2 {
                continue
            }
            ref := match[1]
            
            // Skip URLs and package imports
            if strings.HasPrefix(ref, "http") || strings.HasPrefix(ref, "github.com/") {
                continue
            }

            // Try to resolve the reference
            resolved := resolveReference(ref, currentDir, allFiles)
            if resolved != "" && !refMap[resolved] {
                refs = append(refs, resolved)
                refMap[resolved] = true
            }
        }
    }

    return refs
}

func resolveReference(ref, currentDir string, allFiles []string) string {
    // Skip non-local references
    for _, prefix := range nonLocalPrefixes {
        if strings.HasPrefix(ref, prefix) {
            return ""
        }
    }

    // Skip fragment-only references
    if strings.HasPrefix(ref, "#") {
        return ""
    }

    // Clean and normalize the reference path
    ref = filepath.Clean(ref)
    
    // Handle absolute vs relative paths
    var candidates []string
    if strings.HasPrefix(ref, "./") || strings.HasPrefix(ref, "../") {
        // Relative path - try resolving from current directory
        candidates = []string{filepath.Join(currentDir, ref)}
    } else if strings.Contains(ref, "/") {
        // Path-like reference - try as-is and relative to current dir
        candidates = []string{
            ref,
            filepath.Join(currentDir, ref),
        }
    } else {
        // Bare reference - only look in current directory
        candidates = []string{filepath.Join(currentDir, ref)}
    }

    // First try exact matches
    for _, candidate := range candidates {
        cleanCandidate := filepath.Clean(candidate)
        for _, file := range allFiles {
            if filepath.Clean(file) == cleanCandidate {
                return file
            }
        }
    }

    // If no extension present, try with common extensions
    if filepath.Ext(ref) == "" {
        exts := []string{".go", ".md", ".yml", ".yaml", ".json"}
        baseCandidates := candidates
        for _, candidate := range baseCandidates {
            for _, ext := range exts {
                withExt := candidate + ext
                cleanWithExt := filepath.Clean(withExt)
                for _, file := range allFiles {
                    if filepath.Clean(file) == cleanWithExt {
                        return file
                    }
                }
            }
        }
    }

    return ""
}
