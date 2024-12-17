package rules

import (
    "path/filepath"
    "github.com/1broseidon/promptext/internal/filter/types"
)

type ExtensionRule struct {
    types.BaseRule
    extensions map[string]bool
}

func NewExtensionRule(exts []string, action types.RuleAction) types.Rule {
    extMap := make(map[string]bool)
    for _, ext := range exts {
        extMap[ext] = true
    }
    return &ExtensionRule{
        BaseRule: types.NewBaseRule("", action),
        extensions: extMap,
    }
}

func (r *ExtensionRule) Match(path string) bool {
    ext := filepath.Ext(path)
    return r.extensions[ext]
}
