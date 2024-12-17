package filter

import (
    "path/filepath"
    "strings"
)

// ExtensionFilter filters files based on their extensions
type ExtensionFilter struct {
    extensions []string
    exclude    bool
}

// NewExtensionFilter creates a new extension filter
func NewExtensionFilter(extensions []string, exclude bool) *ExtensionFilter {
    return &ExtensionFilter{
        extensions: extensions,
        exclude:    exclude,
    }
}

// Match checks if the file extension matches any of the configured extensions
func (ef *ExtensionFilter) Match(path string) (bool, error) {
    if len(ef.extensions) == 0 {
        return !ef.exclude, nil
    }

    ext := filepath.Ext(path)
    for _, allowedExt := range ef.extensions {
        if strings.EqualFold(strings.TrimPrefix(allowedExt, "."), strings.TrimPrefix(ext, ".")) {
            return true, nil
        }
    }
    return false, nil
}

// Priority returns the filter priority
func (ef *ExtensionFilter) Priority() int {
    if ef.exclude {
        return 80 // Higher priority for exclusion
    }
    return 20
}

// ShouldInclude determines if the file should be included based on its extension
func (ef *ExtensionFilter) ShouldInclude(path string) bool {
    matched, _ := ef.Match(path)
    return matched != ef.exclude
}
