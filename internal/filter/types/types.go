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
	Pattern    string
	ActionType RuleAction
}

func NewBaseRule(pattern string, action RuleAction) BaseRule {
	return BaseRule{
		Pattern:    pattern,
		ActionType: action,
	}
}

func (r *BaseRule) Action() RuleAction {
	return r.ActionType
}
