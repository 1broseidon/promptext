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
        // Go imports: captures both "utils" and "./utils" style
        regexp.MustCompile(`(?m)^\s*import\s+(?:"([^"]+)"|'([^']+)')`),
        // Node.js require: matches require('utils') or require('./utils')
        regexp.MustCompile(`(?m)require\(['"]([^'"]+)['"]\)`),
        // Markdown links: [text](utils) or [text](./utils)
        regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),
        // HTML includes: src="utils.js" or src="./utils.js"
        regexp.MustCompile(`(?m)(?:src|href)=["']([^'"]+)["']`),
    }
)

// ExtractFileReferences finds references to other files within the given content
func ExtractFileReferences(content, currentDir, rootDir string, allFiles []string) []string {
    var refs []string
    refMap := make(map[string]bool)

    for _, pattern := range referencePatterns {
        matches := pattern.FindAllStringSubmatch(content, -1)
        for _, match := range matches {
            if len(match) < 2 {
                continue
            }
            ref := match[1]
            
            // The captured group might be in match[1] or match[2] depending on pattern
            ref := ""
            for i := 1; i < len(match); i++ {
                if match[i] != "" {
                    ref = match[i]
                    break
                }
            }
            if ref == "" {
                continue
            }

            // Try to resolve the reference
            resolved := resolveReference(ref, currentDir, rootDir, allFiles)
            if resolved != "" && !refMap[resolved] {
                refs = append(refs, resolved)
                refMap[resolved] = true
            }
        }
    }

    return refs
}

func isNonLocalReference(ref string) bool {
    if strings.HasPrefix(ref, "#") {
        return true
    }
    for _, prefix := range nonLocalPrefixes {
        if strings.HasPrefix(ref, prefix) {
            return true
        }
    }
    if strings.Contains(ref, "://") {
        return true
    }
    return false
}

func matchExact(absCandidate, rootDir string, allFiles []string) string {
    for _, file := range allFiles {
        absFile, _ := filepath.Abs(filepath.Join(rootDir, file))
        if absFile == absCandidate {
            return file
        }
    }
    return ""
}

func resolveReference(ref, currentDir, rootDir string, allFiles []string) string {
    if isNonLocalReference(ref) {
        return ""
    }

    // Clean and normalize the reference path
    ref = filepath.Clean(ref)
    
    // Form absolute path from current directory
    absCandidate, err := filepath.Abs(filepath.Join(rootDir, currentDir, ref))
    if err != nil {
        return ""
    }

    // Try exact match first
    if file := matchExact(absCandidate, rootDir, allFiles); file != "" {
        return file
    }

    // If no extension present, try with common extensions
    if filepath.Ext(ref) == "" {
        for _, ext := range []string{".go", ".md", ".yml", ".yaml", ".json"} {
            candidateWithExt := absCandidate + ext
            if file := matchExact(candidateWithExt, rootDir, allFiles); file != "" {
                return file
            }
        }
    }

    return ""
}
