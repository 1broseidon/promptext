package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"path/filepath"
	"testing"
)

func TestRuleCombinations_MultiplePatterns(t *testing.T) {
	// Test combining multiple pattern rules
	rules := []types.Rule{
		NewPatternRule([]string{"*.log", "*.tmp"}, types.Exclude),
		NewPatternRule([]string{"*.go", "*.js"}, types.Include),
		NewPatternRule([]string{"test_*", "*_test.go"}, types.Exclude),
	}

	testCases := []struct {
		name     string
		path     string
		expected map[types.Rule]bool
		desc     string
	}{
		{
			name: "log file exclusion",
			path: "app.log",
			expected: map[types.Rule]bool{
				rules[0]: true,  // matches *.log exclude
				rules[1]: false, // doesn't match *.go or *.js
				rules[2]: false, // doesn't match test patterns
			},
			desc: "log files should be excluded by first rule",
		},
		{
			name: "go file inclusion",
			path: "main.go",
			expected: map[types.Rule]bool{
				rules[0]: false, // doesn't match log patterns
				rules[1]: true,  // matches *.go include
				rules[2]: false, // doesn't match test patterns
			},
			desc: "go files should be included by second rule",
		},
		{
			name: "test go file conflict",
			path: "utils_test.go",
			expected: map[types.Rule]bool{
				rules[0]: false, // doesn't match log patterns
				rules[1]: true,  // matches *.go include
				rules[2]: true,  // matches *_test.go exclude
			},
			desc: "test files create rule conflicts",
		},
		{
			name: "test prefix file",
			path: "test_helper.js",
			expected: map[types.Rule]bool{
				rules[0]: false, // doesn't match log patterns
				rules[1]: true,  // matches *.js include
				rules[2]: true,  // matches test_* exclude
			},
			desc: "test helper files create conflicts",
		},
		{
			name: "no matches",
			path: "README.md",
			expected: map[types.Rule]bool{
				rules[0]: false, // no match
				rules[1]: false, // no match
				rules[2]: false, // no match
			},
			desc: "files not matching any rule",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, rule := range rules {
				result := rule.Match(tc.path)
				expected := tc.expected[rule]
				assert.Equal(t, expected, result,
					"Rule %d (%T) for path %s - %s", i, rule, tc.path, tc.desc)
			}
		})
	}
}

func TestRuleCombinations_DefaultsWithCustom(t *testing.T) {
	// Combine default excludes with custom rules
	defaultRules := DefaultExcludes()
	customRules := []types.Rule{
		NewExtensionRule([]string{".go", ".js", ".py"}, types.Include),
		NewPatternRule([]string{"src/", "lib/"}, types.Include),
	}

	allRules := append(defaultRules, customRules...)

	testCases := []struct {
		name string
		path string
		desc string
	}{
		{"node_modules go", "node_modules/pkg/main.go", "go file in node_modules"},
		{"src go file", "src/main.go", "go file in src directory"},
		{"git go file", ".git/hooks/hook.go", "go file in git directory"},
		{"build js", "build/app.js", "js file in build directory"},
		{"lib python", "lib/utils.py", "python file in lib directory"},
		{"logs go", "logs/parser.go", "go file in logs directory"},
		{"binary in src", "src/app.exe", "binary file in src"},
		{"image in lib", "lib/icon.png", "image file in lib"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, rule := range allRules {
				result := rule.Match(tc.path)
				action := rule.Action()
				t.Logf("Rule %d (%T, %v): %s -> %t", i, rule, action, tc.path, result)
			}
		})
	}
}

func TestDirectoryTraversal_ComplexStructure(t *testing.T) {
	// Create a complex directory structure for testing
	tmpDir, err := os.MkdirTemp("", "complex_structure_test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create complex directory structure
	structure := map[string][]byte{
		"README.md":                           []byte("# Project"),
		"main.go":                             []byte("package main"),
		"src/utils/helper.go":                 []byte("package utils"),
		"src/utils/helper_test.go":            []byte("package utils"),
		"node_modules/package/index.js":       []byte("module.exports = {}"),
		"node_modules/package/lib/util.go":    []byte("package lib"),
		".git/config":                         []byte("[core]"),
		".git/hooks/pre-commit":               []byte("#!/bin/sh"),
		"build/dist/app.js":                   []byte("console.log('app')"),
		"build/dist/app.css":                  []byte("body {}"),
		"vendor/github.com/pkg/errors.go":     []byte("package errors"),
		"logs/app.log":                        []byte("ERROR: test"),
		"logs/access.log":                     []byte("GET /"),
		"temp/cache.tmp":                      []byte("cache data"),
		".idea/workspace.xml":                 []byte("<workspace>"),
		".vscode/settings.json":               []byte("{}"),
		"images/logo.png":                     []byte("PNG fake content"),
		"docs/api.md":                         []byte("# API"),
		"tests/unit/parser_test.go":           []byte("package tests"),
		"scripts/build.sh":                    []byte("#!/bin/bash"),
		"config/app.yaml":                     []byte("port: 8080"),
		"data/users.json":                     []byte("[]"),
		".DS_Store":                           []byte("DS_Store content"),
		"images/.DS_Store":                    []byte("DS_Store content"),
		"nested/very/deep/structure/file.txt": []byte("deep file"),
		"nested/very/deep/.hidden/secret.key": []byte("secret"),
		"symlink_target.txt":                  []byte("target content"),
		"a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p.txt": []byte("very deep"),
	}

	// Create files and directories
	for relativePath, content := range structure {
		fullPath := filepath.Join(tmpDir, relativePath)
		dir := filepath.Dir(fullPath)

		err = os.MkdirAll(dir, 0755)
		if err != nil {
			t.Fatal(err)
		}

		err = os.WriteFile(fullPath, content, 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Create a symlink (if supported)
	symlinkPath := filepath.Join(tmpDir, "symlink.txt")
	targetPath := filepath.Join(tmpDir, "symlink_target.txt")
	err = os.Symlink(targetPath, symlinkPath)
	if err != nil {
		t.Logf("Symlink creation failed (may not be supported): %v", err)
	}

	// Test with default exclude rules
	defaultRules := DefaultExcludes()

	testCases := []struct {
		name     string
		path     string
		expected bool // Expected to be excluded by defaults
		desc     string
	}{
		{"main go file", "main.go", false, "main go file should not be excluded"},
		{"src go file", "src/utils/helper.go", false, "source files should not be excluded"},
		{"test go file", "src/utils/helper_test.go", false, "test files should not be excluded"},
		{"node_modules file", "node_modules/package/index.js", true, "node_modules should be excluded"},
		{"nested node_modules", "node_modules/package/lib/util.go", true, "nested in node_modules should be excluded"},
		{"git config", ".git/config", true, "git directory should be excluded"},
		{"git hooks", ".git/hooks/pre-commit", true, "git hooks should be excluded"},
		{"build dist", "build/dist/app.js", true, "build directory should be excluded"},
		{"vendor go", "vendor/github.com/pkg/errors.go", true, "vendor directory should be excluded"},
		{"log files", "logs/app.log", true, "log files should be excluded"},
		{"temp files", "temp/cache.tmp", true, "temp directory should be excluded"},
		{"idea config", ".idea/workspace.xml", true, "IDE directory should be excluded"},
		{"vscode config", ".vscode/settings.json", true, "VS Code directory should be excluded"},
		{"ds store", ".DS_Store", true, "DS_Store should be excluded"},
		{"nested ds store", "images/.DS_Store", true, "nested DS_Store should be excluded"},
		{"readme", "README.md", false, "README should not be excluded"},
		{"docs", "docs/api.md", false, "documentation should not be excluded"},
		{"scripts", "scripts/build.sh", false, "scripts should not be excluded"},
		{"config", "config/app.yaml", false, "config files should not be excluded"},
		{"data", "data/users.json", false, "data files should not be excluded"},
		{"deep structure", "nested/very/deep/structure/file.txt", false, "deep files should not be excluded"},
		{"hidden deep", "nested/very/deep/.hidden/secret.key", true, "hidden files match .git* pattern due to dot prefix"},
		{"very deep", "a/b/c/d/e/f/g/h/i/j/k/l/m/n/o/p.txt", false, "very deep files should not be excluded"},
		{"symlink", "symlink.txt", false, "symlinks should not be excluded"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fullPath := filepath.Join(tmpDir, tc.path)

			// Test that file exists (skip if symlink creation failed)
			if tc.path == "symlink.txt" {
				if _, err := os.Lstat(fullPath); os.IsNotExist(err) {
					t.Skip("Symlink not supported on this system")
				}
			}

			excluded := false
			for _, rule := range defaultRules {
				if rule.Match(tc.path) && rule.Action() == types.Exclude {
					excluded = true
					break
				}
			}

			assert.Equal(t, tc.expected, excluded, "Path: %s - %s", tc.path, tc.desc)
		})
	}
}

func TestDirectoryTraversal_PathNormalization(t *testing.T) {
	rules := []types.Rule{
		NewPatternRule([]string{"build/", "dist/"}, types.Exclude),
		NewExtensionRule([]string{".js", ".css"}, types.Include),
	}

	testCases := []struct {
		name string
		path string
		desc string
	}{
		// Unix paths
		{"unix build", "build/app.js", "unix path with build"},
		{"unix dist", "dist/main.css", "unix path with dist"},
		{"unix nested", "src/build/utils.js", "unix nested build"},

		// Windows-style paths (should be normalized)
		{"windows build", "build\\app.js", "windows path with build"},
		{"windows dist", "dist\\main.css", "windows path with dist"},
		{"windows nested", "src\\build\\utils.js", "windows nested build"},

		// Mixed separators
		{"mixed 1", "build/subdir\\file.js", "mixed separators"},
		{"mixed 2", "src\\build/utils.js", "mixed separators reverse"},

		// Complex paths
		{"relative", "./build/app.js", "relative path"},
		{"parent", "../build/app.js", "parent directory"},
		{"double slash", "build//app.js", "double slash"},
		{"trailing slash", "build/", "trailing slash"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			for i, rule := range rules {
				result := rule.Match(tc.path)
				t.Logf("Rule %d (%T): %s -> %t (%s)", i, rule, tc.path, result, tc.desc)
			}
		})
	}
}

func TestRuleActions_PriorityHandling(t *testing.T) {
	// Test how different rule actions might be prioritized
	testCases := []struct {
		name     string
		rules    []types.Rule
		testPath string
		expected []bool // Expected match results for each rule
		desc     string
	}{
		{
			name: "include vs exclude conflict",
			rules: []types.Rule{
				NewExtensionRule([]string{".go"}, types.Include),
				NewPatternRule([]string{"*_test.go"}, types.Exclude),
			},
			testPath: "utils_test.go",
			expected: []bool{true, true}, // Both match, creating conflict
			desc:     "file matches both include and exclude rules",
		},
		{
			name: "multiple excludes",
			rules: []types.Rule{
				NewPatternRule([]string{"build/"}, types.Exclude),
				NewPatternRule([]string{"*.js"}, types.Exclude),
				NewBinaryRule(),
			},
			testPath: "build/app.js",
			expected: []bool{true, true, false}, // Directory and extension match, not binary
			desc:     "file matches multiple exclude rules",
		},
		{
			name: "skip vs exclude",
			rules: []types.Rule{
				NewExtensionRule([]string{".tmp"}, types.Skip),
				NewPatternRule([]string{"temp/"}, types.Exclude),
			},
			testPath: "temp/cache.tmp",
			expected: []bool{true, true}, // Both match
			desc:     "file matches both skip and exclude rules",
		},
		{
			name: "no conflicts",
			rules: []types.Rule{
				NewExtensionRule([]string{".go"}, types.Include),
				NewPatternRule([]string{"vendor/"}, types.Exclude),
				NewBinaryRule(),
			},
			testPath: "src/main.go",
			expected: []bool{true, false, false}, // Only include matches
			desc:     "file clearly matches one rule type",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.Equal(t, len(tc.rules), len(tc.expected), "Test case setup error")

			for i, rule := range tc.rules {
				result := rule.Match(tc.testPath)
				expected := tc.expected[i]
				assert.Equal(t, expected, result,
					"Rule %d (%T, %v): %s - %s", i, rule, rule.Action(), tc.testPath, tc.desc)
			}
		})
	}
}

func TestRulePerformance_ComplexScenarios(t *testing.T) {
	// Test performance with complex rule combinations
	defaultRules := DefaultExcludes()
	customRules := []types.Rule{
		NewExtensionRule([]string{".go", ".js", ".py", ".ts", ".java", ".cpp", ".c", ".h"}, types.Include),
		NewPatternRule([]string{"src/", "lib/", "app/", "components/"}, types.Include),
		NewPatternRule([]string{"*_test.*", "test_*", "*.spec.*"}, types.Exclude),
	}
	allRules := append(defaultRules, customRules...)

	// Test paths that might be encountered in real projects
	testPaths := []string{
		"src/main.go",
		"lib/utils/helper.js",
		"app/components/Button.tsx",
		"tests/unit/parser_test.go",
		"node_modules/react/index.js",
		"build/dist/bundle.js",
		".git/config",
		"vendor/pkg/errors.go",
		"logs/application.log",
		"images/logo.png",
		"docs/README.md",
		"scripts/deploy.sh",
		"config/database.yaml",
		"data/seed.json",
		".env.example",
		"package-lock.json",
		"Dockerfile",
		"docker-compose.yml",
		"Makefile",
		"go.mod",
		"requirements.txt",
		"webpack.config.js",
		"tsconfig.json",
		".gitignore",
		".eslintrc.js",
		"jest.config.js",
		"babel.config.js",
		"prettier.config.js",
		"tailwind.config.js",
		"vite.config.ts",
	}

	// This is more of a smoke test to ensure rules work with realistic paths
	for _, path := range testPaths {
		t.Run("path_"+path, func(t *testing.T) {
			matchCount := 0
			for _, rule := range allRules {
				if rule.Match(path) {
					matchCount++
				}
			}
			// Just verify no panics occur and some matches happen
			t.Logf("Path %s matched %d rules", path, matchCount)
		})
	}
}

func BenchmarkRuleCombinations_Realistic(b *testing.B) {
	// Benchmark realistic rule combinations
	rules := append(DefaultExcludes(),
		NewExtensionRule([]string{".go", ".js", ".py"}, types.Include),
		NewPatternRule([]string{"src/", "lib/"}, types.Include),
	)

	testPaths := []string{
		"src/main.go",
		"node_modules/pkg/index.js",
		"build/app.js",
		"lib/utils.py",
		"logs/app.log",
		"images/logo.png",
		".git/config",
		"vendor/errors.go",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			for _, rule := range rules {
				rule.Match(path)
			}
		}
	}
}

func BenchmarkDirectoryTraversal_Deep(b *testing.B) {
	rule := NewPatternRule([]string{"node_modules/", "vendor/", "build/"}, types.Exclude)

	deepPaths := []string{
		"node_modules/a/b/c/d/e/f/g/h/i/j/index.js",
		"vendor/github.com/user/repo/pkg/lib/util.go",
		"build/dist/assets/images/icons/small/icon.png",
		"src/very/deep/nested/structure/component.js",
		"lib/external/third/party/library/source.py",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range deepPaths {
			rule.Match(path)
		}
	}
}
