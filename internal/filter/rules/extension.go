package rules

import (
    "path/filepath"
    "github.com/1broseidon/promptext/internal/filter"
)

type ExtensionRule struct {
    BaseRule
    extensions map[string]bool
}

func NewExtensionRule(exts []string, action filter.RuleAction) Rule {
    extMap := make(map[string]bool)
    for _, ext := range exts {
        extMap[ext] = true
    }
    return &ExtensionRule{
        BaseRule: filter.BaseRule{action: action},
        extensions: extMap,
    }
}

func (r *ExtensionRule) Match(path string) bool {
    ext := filepath.Ext(path)
    return r.extensions[ext]
}
