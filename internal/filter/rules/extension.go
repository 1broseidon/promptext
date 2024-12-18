package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"path/filepath"
	"strings"
)

type ExtensionRule struct {
	types.BaseRule
	extensions map[string]bool
}

func NewExtensionRule(exts []string, action types.RuleAction) types.Rule {
	extMap := make(map[string]bool)
	for _, ext := range exts {
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		extMap[ext] = true
	}
	return &ExtensionRule{
		BaseRule:   types.NewBaseRule("", action),
		extensions: extMap,
	}
}

func (r *ExtensionRule) Match(path string) bool {
	ext := filepath.Ext(path)
	if ext == "" {
		return false
	}
	matches := r.extensions[ext]
	// For excludes, we want to return true if it matches (to trigger exclusion)
	// For includes, we want to return true if it matches (to trigger inclusion)
	return matches
}
