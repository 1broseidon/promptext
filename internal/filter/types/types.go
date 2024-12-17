package types

type RuleAction int

const (
    Include RuleAction = iota
    Exclude
    Skip
)

// Rule defines a single filtering rule
type Rule interface {
    Match(path string) bool
    Action() RuleAction
}

// BaseRule provides common functionality
type BaseRule struct {
    pattern string
    action  RuleAction
}

func (r *BaseRule) Action() RuleAction {
    return r.action
}
