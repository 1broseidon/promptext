package rules

type Rule interface {
	Match(path string) bool
	Action() RuleAction
}

type RuleAction int

const (
	Include RuleAction = iota
	Exclude
	Skip
)

// BaseRule provides common functionality
type BaseRule struct {
	pattern string
	action  RuleAction
}

func (r *BaseRule) Action() RuleAction {
	return r.action
}
