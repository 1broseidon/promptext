package references

import (
    "path/filepath"
    "regexp"
    "strings"
)

var referencePatterns = []*regexp.Regexp{
    // Go imports
    regexp.MustCompile(`(?m)^import\s+(?:\([^)]+\)|["']([^"']+)["'])`),
    // require('X'), Node.js style
    regexp.MustCompile(`(?m)require\(['"]([^'"]+)['"]\)`),
    // Markdown links to files: [text](/path/to/file)
    regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),
    // HTML includes
    regexp.MustCompile(`(?m)(?:src|href)=["']([^'"]+)["']`),
}

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
    // Try different path combinations
    candidates := []string{
        ref,
        filepath.Join(currentDir, ref),
        filepath.Clean(ref),
    }

    // Add common extensions if no extension present
    if filepath.Ext(ref) == "" {
        exts := []string{".go", ".md", ".yml", ".yaml", ".json"}
        baseCandidates := candidates
        for _, candidate := range baseCandidates {
            for _, ext := range exts {
                candidates = append(candidates, candidate+ext)
            }
        }
    }

    // Check each candidate against allFiles
    for _, candidate := range candidates {
        for _, file := range allFiles {
            if strings.HasSuffix(file, candidate) {
                return file
            }
        }
    }

    return ""
}
