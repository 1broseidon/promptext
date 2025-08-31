package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Mock rule implementation for testing the base rule functionality
type MockRule struct {
	types.BaseRule
	matchResult bool
}

func NewMockRule(action types.RuleAction, matchResult bool) *MockRule {
	return &MockRule{
		BaseRule:    types.NewBaseRule("test-pattern", action),
		matchResult: matchResult,
	}
}

func (r *MockRule) Match(path string) bool {
	return r.matchResult
}

func TestBaseRule_Action(t *testing.T) {
	testCases := []struct {
		name   string
		action types.RuleAction
		desc   string
	}{
		{
			name:   "include action",
			action: types.Include,
			desc:   "rule with include action",
		},
		{
			name:   "exclude action", 
			action: types.Exclude,
			desc:   "rule with exclude action",
		},
		{
			name:   "skip action",
			action: types.Skip,
			desc:   "rule with skip action",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rule := NewMockRule(tc.action, true)
			
			result := rule.Action()
			assert.Equal(t, tc.action, result, tc.desc)
		})
	}
}

func TestRuleAction_Constants(t *testing.T) {
	// Test that the constants are properly defined
	assert.Equal(t, types.RuleAction(0), types.Include, "Include should be 0")
	assert.Equal(t, types.RuleAction(1), types.Exclude, "Exclude should be 1") 
	assert.Equal(t, types.RuleAction(2), types.Skip, "Skip should be 2")
}

func TestRule_Interface(t *testing.T) {
	// Test that our rules implement the Rule interface properly
	var rule types.Rule
	
	// Test with different rule implementations
	rule = NewPatternRule([]string{"*.test"}, types.Include)
	assert.NotNil(t, rule.Match, "PatternRule should implement Match")
	assert.NotNil(t, rule.Action, "PatternRule should implement Action")
	
	rule = NewExtensionRule([]string{".go"}, types.Exclude)
	assert.NotNil(t, rule.Match, "ExtensionRule should implement Match")
	assert.NotNil(t, rule.Action, "ExtensionRule should implement Action")
	
	rule = NewBinaryRule()
	assert.NotNil(t, rule.Match, "BinaryRule should implement Match")
	assert.NotNil(t, rule.Action, "BinaryRule should implement Action")
}

func TestMockRule_Implementation(t *testing.T) {
	// Test the mock rule we created for testing
	tests := []struct {
		name        string
		action      types.RuleAction
		matchResult bool
		testPath    string
		desc        string
	}{
		{
			name:        "always match include",
			action:      types.Include,
			matchResult: true,
			testPath:    "any/path.txt",
			desc:        "mock rule that always matches with include action",
		},
		{
			name:        "never match exclude",
			action:      types.Exclude,
			matchResult: false,
			testPath:    "any/path.txt",
			desc:        "mock rule that never matches with exclude action",
		},
		{
			name:        "always match skip",
			action:      types.Skip,
			matchResult: true,
			testPath:    "different/path.go",
			desc:        "mock rule that always matches with skip action",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewMockRule(tt.action, tt.matchResult)
			
			// Test Action method
			assert.Equal(t, tt.action, rule.Action(), "Action should return correct value")
			
			// Test Match method
			assert.Equal(t, tt.matchResult, rule.Match(tt.testPath), "Match should return configured result")
		})
	}
}

func BenchmarkBaseRule_Action(b *testing.B) {
	rule := NewMockRule(types.Include, true)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Action()
	}
}