package types

import "testing"

func TestNewBaseRule(t *testing.T) {
	rule := NewBaseRule("*.go", Include)
	if rule.Pattern != "*.go" {
		t.Fatalf("expected pattern to be preserved, got %s", rule.Pattern)
	}
	if rule.ActionType != Include {
		t.Fatalf("expected ActionType Include, got %v", rule.ActionType)
	}
	if rule.Action() != Include {
		t.Fatalf("expected Action() to return Include")
	}
}

func TestBaseRuleActionValues(t *testing.T) {
	excludeRule := NewBaseRule("*.txt", Exclude)
	if excludeRule.Action() != Exclude {
		t.Fatalf("expected exclude action")
	}
	skipRule := NewBaseRule("*.tmp", Skip)
	if skipRule.Action() != Skip {
		t.Fatalf("expected skip action")
	}
}
