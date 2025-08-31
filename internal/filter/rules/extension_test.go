package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewExtensionRule(t *testing.T) {
	tests := []struct {
		name        string
		extensions  []string
		action      types.RuleAction
		expectCount int
		desc        string
	}{
		{
			name:        "single extension with dot",
			extensions:  []string{".go"},
			action:      types.Include,
			expectCount: 1,
			desc:        "single extension already with dot",
		},
		{
			name:        "single extension without dot",
			extensions:  []string{"js"},
			action:      types.Exclude,
			expectCount: 1,
			desc:        "single extension without dot should get dot added",
		},
		{
			name:        "multiple extensions mixed",
			extensions:  []string{".go", "js", ".py", "ts"},
			action:      types.Include,
			expectCount: 4,
			desc:        "mixed extensions with and without dots",
		},
		{
			name:        "empty extension list",
			extensions:  []string{},
			action:      types.Skip,
			expectCount: 0,
			desc:        "empty extension list",
		},
		{
			name:        "duplicate extensions",
			extensions:  []string{"go", ".go", "js", ".js"},
			action:      types.Exclude,
			expectCount: 2,
			desc:        "duplicate extensions should be deduplicated",
		},
		{
			name:        "case sensitive extensions",
			extensions:  []string{".GO", ".go", ".Js", ".js"},
			action:      types.Include,
			expectCount: 4,
			desc:        "case sensitive extension handling",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewExtensionRule(tt.extensions, tt.action)
			
			require.NotNil(t, rule)
			assert.Equal(t, tt.action, rule.Action())
			
			// Verify rule type
			extRule, ok := rule.(*ExtensionRule)
			require.True(t, ok, "Expected ExtensionRule type")
			assert.Equal(t, tt.expectCount, len(extRule.extensions), tt.desc)
		})
	}
}

func TestExtensionRule_Match_BasicExtensions(t *testing.T) {
	rule := NewExtensionRule([]string{"go", ".js", ".py", "ts"}, types.Include)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Matching extensions
		{"go file", "main.go", true, "go extension should match"},
		{"js file", "app.js", true, "js extension should match"},
		{"python file", "script.py", true, "python extension should match"},
		{"typescript file", "component.ts", true, "typescript extension should match"},

		// Nested paths with matching extensions
		{"nested go", "src/utils/helper.go", true, "nested go file"},
		{"nested js", "frontend/components/App.js", true, "nested js file"},
		{"deeply nested py", "backend/api/v1/handlers/auth.py", true, "deeply nested python file"},
		{"nested ts", "types/models/user.ts", true, "nested typescript file"},

		// Non-matching extensions
		{"txt file", "README.txt", false, "txt extension should not match"},
		{"md file", "docs.md", false, "md extension should not match"},
		{"json file", "config.json", false, "json extension should not match"},
		{"no extension", "Makefile", false, "file without extension"},
		{"hidden file no ext", ".gitignore", false, "hidden file without extension"},

		// Case sensitivity
		{"uppercase GO", "main.GO", false, "case sensitive - GO vs go"},
		{"mixed case Js", "app.Js", false, "case sensitive - Js vs js"},
		{"uppercase PY", "script.PY", false, "case sensitive - PY vs py"},

		// Edge cases
		{"multiple dots", "file.test.js", true, "multiple dots - should match .js"},
		{"dot in filename", "jquery-3.6.0.js", true, "dots in filename should work"},
		{"extension-like in name", "gofile.txt", false, "extension in filename but wrong actual extension"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestExtensionRule_Match_EdgeCaseExtensions(t *testing.T) {
	rule := NewExtensionRule([]string{".gitignore", ".env", ".config"}, types.Exclude)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Files that look like extensions but are exact matches
		{"gitignore file", ".gitignore", true, ".gitignore matches .gitignore extension"},
		{"env file", ".env", true, ".env matches .env extension"},
		{"actual gitignore ext", "project.gitignore", true, "file with .gitignore extension"},
		{"actual env ext", "development.env", true, "file with .env extension"},
		{"actual config ext", "app.config", true, "file with .config extension"},

		// Hidden files with extensions
		{"hidden js", ".hidden.js", false, "hidden file with .js extension (not in rule)"},
		{"hidden env ext", ".secrets.env", true, "hidden file with .env extension"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestExtensionRule_Match_EmptyAndSpecialExtensions(t *testing.T) {
	tests := []struct {
		name       string
		extensions []string
		testCases  []struct {
			path     string
			expected bool
			desc     string
		}
	}{
		{
			name:       "empty extension rule",
			extensions: []string{},
			testCases: []struct {
				path     string
				expected bool
				desc     string
			}{
				{"any.file", false, "empty rule should not match anything"},
				{"no-extension", false, "empty rule should not match files without extension"},
				{".hidden", false, "empty rule should not match hidden files"},
			},
		},
		{
			name:       "very long extensions",
			extensions: []string{".verylongextensionnamehere", ".anotherverylongext"},
			testCases: []struct {
				path     string
				expected bool
				desc     string
			}{
				{"file.verylongextensionnamehere", true, "very long extension should match"},
				{"file.anotherverylongext", true, "another long extension should match"},
				{"file.verylongextension", false, "partial long extension should not match"},
			},
		},
		{
			name:       "numeric extensions",
			extensions: []string{".123", ".v2", ".001"},
			testCases: []struct {
				path     string
				expected bool
				desc     string
			}{
				{"backup.123", true, "numeric extension should match"},
				{"file.v2", true, "alphanumeric extension should match"},
				{"data.001", true, "zero-padded numeric extension should match"},
				{"file.124", false, "different numeric extension should not match"},
			},
		},
		{
			name:       "special character extensions",
			extensions: []string{".c++", ".~tmp", ".#bak"},
			testCases: []struct {
				path     string
				expected bool
				desc     string
			}{
				{"source.c++", true, "c++ extension should match"},
				{"temp.~tmp", true, "tilde extension should match"},
				{"backup.#bak", true, "hash extension should match"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewExtensionRule(tt.extensions, types.Include)
			
			for _, tc := range tt.testCases {
				t.Run(tc.desc, func(t *testing.T) {
					result := rule.Match(tc.path)
					assert.Equal(t, tc.expected, result, "Path: %s - %s", tc.path, tc.desc)
				})
			}
		})
	}
}

func TestExtensionRule_Match_NoExtension(t *testing.T) {
	rule := NewExtensionRule([]string{".txt", ".md", ".go"}, types.Include)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		{"no extension", "Makefile", false, "file without extension"},
		{"dockerfile", "Dockerfile", false, "dockerfile without extension"},
		{"license", "LICENSE", false, "license file without extension"},
		{"readme", "README", false, "readme without extension"},
		{"hidden no ext", ".gitignore", false, "hidden file without extension"},
		{"hidden with dot", ".bashrc", false, "hidden file starting with dot but no extension"},
		{"empty filename", "", false, "empty filename"},
		{"just dot", ".", false, "just a dot"},
		{"double dot", "..", false, "double dot"},
		{"directory-like", "src/", false, "directory-like path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestExtensionRule_Match_CaseSensitivity(t *testing.T) {
	// Test case sensitivity explicitly
	rule := NewExtensionRule([]string{".GO", ".Js", ".PY"}, types.Include)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Exact case matches
		{"exact GO", "main.GO", true, "exact case match for GO"},
		{"exact Js", "app.Js", true, "exact case match for Js"},
		{"exact PY", "script.PY", true, "exact case match for PY"},

		// Different case - should not match
		{"lowercase go", "main.go", false, "lowercase should not match uppercase rule"},
		{"lowercase js", "app.js", false, "lowercase should not match mixed case rule"},
		{"lowercase py", "script.py", false, "lowercase should not match uppercase rule"},
		{"all caps JS", "app.JS", false, "all caps should not match mixed case rule"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestExtensionRule_Action(t *testing.T) {
	tests := []struct {
		name   string
		action types.RuleAction
	}{
		{"include action", types.Include},
		{"exclude action", types.Exclude},
		{"skip action", types.Skip},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewExtensionRule([]string{".test"}, tt.action)
			assert.Equal(t, tt.action, rule.Action())
		})
	}
}

func TestExtensionRule_DotHandling(t *testing.T) {
	tests := []struct {
		name       string
		input      []string
		testPath   string
		expected   bool
		desc       string
	}{
		{
			name:     "dot added automatically",
			input:    []string{"go", "js"},
			testPath: "main.go",
			expected: true,
			desc:     "dot should be added automatically to extensions without dot",
		},
		{
			name:     "dot preserved",
			input:    []string{".go", ".js"},
			testPath: "main.go",
			expected: true,
			desc:     "dot should be preserved for extensions with dot",
		},
		{
			name:     "mixed dot handling",
			input:    []string{".go", "js", ".py", "ts"},
			testPath: "app.js",
			expected: true,
			desc:     "mixed dot handling should work correctly",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewExtensionRule(tt.input, types.Include)
			result := rule.Match(tt.testPath)
			assert.Equal(t, tt.expected, result, tt.desc)
		})
	}
}

func TestExtensionRule_PathNormalization(t *testing.T) {
	rule := NewExtensionRule([]string{".go", ".js"}, types.Include)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Different path separators should not affect extension matching
		{"unix path", "src/main.go", true, "unix path separator"},
		{"windows path", "src\\main.go", true, "windows path separator"},
		{"mixed separators", "src/utils\\helper.go", true, "mixed path separators"},
		
		// Complex paths
		{"relative path", "./src/main.go", true, "relative path with dot slash"},
		{"parent path", "../utils/helper.go", true, "parent directory reference"},
		{"absolute path", "/home/user/project/main.go", true, "absolute path"},
		
		// URLs or network paths (should still work for extension)
		{"url-like", "http://example.com/script.js", true, "url-like path"},
		{"network path", "\\\\server\\share\\file.go", true, "network path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func BenchmarkExtensionRule_SingleExtension(b *testing.B) {
	rule := NewExtensionRule([]string{".go"}, types.Include)
	path := "src/main.go"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkExtensionRule_MultipleExtensions(b *testing.B) {
	rule := NewExtensionRule([]string{".go", ".js", ".py", ".ts", ".java", ".cpp", ".c", ".h"}, types.Include)
	path := "src/component.ts"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkExtensionRule_NoExtension(b *testing.B) {
	rule := NewExtensionRule([]string{".go", ".js", ".py"}, types.Include)
	path := "Makefile"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkExtensionRule_NoMatch(b *testing.B) {
	rule := NewExtensionRule([]string{".go", ".js", ".py"}, types.Include)
	path := "config.json"
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}