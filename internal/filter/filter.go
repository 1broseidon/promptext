package filter

import (
    "path/filepath"
    "strings"
)

func ShouldProcessFile(path string, extensions, excludes []string) bool {
    // Check if file should be excluded
    for _, exclude := range excludes {
        if strings.Contains(path, exclude) {
            return false
        }
    }

    // If no extensions specified, process all files
    if len(extensions) == 0 {
        return true
    }

    // Check if file extension matches
    ext := filepath.Ext(path)
    for _, allowedExt := range extensions {
        if strings.TrimPrefix(allowedExt, ".") == strings.TrimPrefix(ext, ".") {
            return true
        }
    }

    return false
}
