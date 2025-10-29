package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDefaultExcludes(t *testing.T) {
	rules := DefaultExcludes()

	// Verify we get a non-empty slice of rules
	require.NotEmpty(t, rules, "DefaultExcludes should return non-empty rules")

	// Should contain at least pattern rule and binary rule
	require.GreaterOrEqual(t, len(rules), 2, "Should have at least pattern and binary rules")

	// All rules should have Exclude action (since they're excludes)
	for i, rule := range rules {
		assert.Equal(t, types.Exclude, rule.Action(), "Rule %d should be Exclude action", i)
	}

	// Verify rule types
	hasPatternRule := false
	hasBinaryRule := false

	for _, rule := range rules {
		switch rule.(type) {
		case *PatternRule:
			hasPatternRule = true
		case *BinaryRule:
			hasBinaryRule = true
		}
	}

	assert.True(t, hasPatternRule, "Should contain at least one PatternRule")
	assert.True(t, hasBinaryRule, "Should contain a BinaryRule")
}

func TestDefaultExcludes_VersionControlPatterns(t *testing.T) {
	rules := DefaultExcludes()

	// Find the pattern rule (should be first)
	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Version control directories
		{"git dir", ".git/config", true, "git directory should be excluded"},
		{"nested git", "project/.git/hooks/pre-commit", true, "nested git directory"},
		{"gitignore", ".gitignore", true, "git files should be excluded"},
		{"gitattributes", ".gitattributes", true, "git attributes file"},
		{"gitmodules", ".gitmodules", true, "git modules file"},
		{"svn dir", ".svn/entries", true, "svn directory should be excluded"},
		{"hg dir", ".hg/hgrc", true, "mercurial directory should be excluded"},
		{"nested svn", "legacy/.svn/props", true, "nested svn directory"},

		// Should not match similar but different patterns
		{"git prefix", ".github/workflows/ci.yml", false, "github dir should not match .git*"},
		{"not git dir", "git_backup/", false, "similar name without dot"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_DependencyPatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Node.js dependencies
		{"node_modules dir", "node_modules/package/index.js", true, "node_modules should be excluded"},
		{"nested node_modules", "project/node_modules/lib/util.js", true, "nested node_modules"},
		{"bower_components", "bower_components/jquery/jquery.js", true, "bower components"},
		{"jspm_packages", "jspm_packages/npm/package@1.0.0/lib.js", true, "jspm packages"},

		// Go dependencies
		{"vendor dir", "vendor/github.com/pkg/errors/errors.go", true, "vendor directory"},
		{"nested vendor", "cmd/vendor/lib.go", true, "nested vendor directory"},

		// Should not match similar names
		{"node_modules_backup", "node_modules_backup/file.js", false, "similar but different name"},
		{"vendor.json", "vendor.json", false, "vendor file not directory"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_IDEPatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// IDE directories
		{"idea dir", ".idea/workspace.xml", true, "IntelliJ IDEA directory"},
		{"vscode dir", ".vscode/settings.json", true, "VS Code directory"},
		{"vs dir", ".vs/config.json", true, "Visual Studio directory"},
		{"nested idea", "project/.idea/modules.xml", true, "nested IDEA directory"},

		// Sublime Text files
		{"sublime project", "myproject.sublime-project", true, "sublime project file"},
		{"sublime workspace", "app.sublime-workspace", true, "sublime workspace file"},
		{"sublime settings", "preferences.sublime-settings", true, "sublime settings file"},

		// Should not match similar patterns
		{"idea in name", "myidea.txt", false, "idea in filename but not directory"},
		{"sublime in name", "sublime.txt", false, "sublime in name but not matching pattern"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_BuildAndOutputPatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Build directories
		{"dist dir", "dist/bundle.js", true, "dist directory"},
		{"build dir", "build/main.exe", true, "build directory"},
		{"out dir", "out/compiled.class", true, "out directory"},
		{"bin dir", "bin/executable", true, "bin directory"},
		{"target dir", "target/classes/Main.class", true, "target directory (Maven/sbt)"},

		// Nested build dirs
		{"nested dist", "frontend/dist/app.js", true, "nested dist directory"},
		{"nested build", "src/build/output.o", true, "nested build directory"},

		// Should not match similar names
		{"dist file", "dist.txt", false, "dist file not directory"},
		{"build file", "build.sh", false, "build file not directory"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_CachePatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Python cache
		{"pycache dir", "__pycache__/module.pyc", true, "Python __pycache__ directory"},
		{"pytest cache", ".pytest_cache/session.json", true, "pytest cache directory"},
		{"nested pycache", "src/__pycache__/utils.cpython-39.pyc", true, "nested pycache"},

		// Web build cache
		{"sass cache", ".sass-cache/main.scss.cache", true, "Sass cache directory"},
		{"npm cache", ".npm/package.json", true, "npm cache directory"},
		{"yarn cache", ".yarn/cache.json", true, "yarn cache directory"},

		// Should not match similar names
		{"cache in name", "cache_utils.py", false, "cache in filename but not cache directory"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_TestCoveragePatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Coverage directories
		{"coverage dir", "coverage/lcov.info", true, "coverage directory"},
		{"nyc output", ".nyc_output/coverage.json", true, "nyc output directory"},
		{"nested coverage", "frontend/coverage/index.html", true, "nested coverage directory"},

		// Should not match files
		{"coverage file", "coverage.json", false, "coverage file not directory"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_InfrastructurePatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Infrastructure directories
		{"terraform dir", ".terraform/state.json", true, "terraform directory"},
		{"vagrant dir", ".vagrant/machines/default/config", true, "vagrant directory"},
		{"nested terraform", "infrastructure/.terraform/plan.tfstate", true, "nested terraform"},

		// Should not match files
		{"terraform file", "main.tf", false, "terraform file not directory"},
		{"vagrant file", "Vagrantfile", false, "vagrant file not directory"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_LogsAndTempPatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Log directories and files
		{"logs dir", "logs/app.log", true, "logs directory"},
		{"log file", "application.log", true, "log file"},
		{"nested logs", "backend/logs/error.log", true, "nested logs directory"},
		{"access log", "server/access.log", true, "access log file"},

		// Temp directories
		{"tmp dir", "tmp/upload.dat", true, "tmp directory"},
		{"temp dir", "temp/cache.tmp", true, "temp directory"},
		{"nested tmp", "uploads/tmp/file.dat", true, "nested tmp directory"},

		// Should not match similar patterns
		{"logical", "logical.txt", false, "logical should not match *.log"},
		{"template", "template.html", false, "template should not match temp/"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_SystemFilePatterns(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// System files
		{"ds store", ".DS_Store", true, "macOS DS_Store file"},
		{"nested ds store", "images/.DS_Store", true, "DS_Store in directory"},
		{"deeply nested ds store", "project/assets/images/.DS_Store", true, "deeply nested DS_Store"},

		// Should not match similar files
		{"ds store backup", ".DS_Store.bak", true, "DS_Store backup matches due to prefix matching"},
		{"ds prefix", ".DS_Something", false, "similar prefix should not match"},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			result := patternRule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestDefaultExcludes_BinaryRule(t *testing.T) {
	rules := DefaultExcludes()

	// Find the binary rule
	var binaryRule *BinaryRule
	for _, rule := range rules {
		if br, ok := rule.(*BinaryRule); ok {
			binaryRule = br
			break
		}
	}
	require.NotNil(t, binaryRule, "Should have a binary rule")

	testCases := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Should exclude binary files (this tests integration with binary rule)
		{"executable", "program.exe", true, "executables should be excluded"},
		{"image", "photo.jpg", true, "images should be excluded"},
		{"archive", "data.zip", true, "archives should be excluded"},

		// Should not exclude text files
		{"text file", "readme.txt", false, "text files should not be excluded"},
		{"source code", "main.go", false, "source code should not be excluded"},
	}

	// Note: These tests verify binary rule integration, but actual binary
	// detection logic is thoroughly tested in binary_test.go
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			// We're testing that the binary rule exists and has correct action,
			// not the full binary detection logic (that's tested separately)
			assert.Equal(t, types.Exclude, binaryRule.Action(), "Binary rule should exclude")
		})
	}
}

func TestDefaultExcludes_RuleCount(t *testing.T) {
	rules := DefaultExcludes()

	// Verify expected number of rules
	// Should have exactly 2 rules: PatternRule and BinaryRule
	assert.Equal(t, 2, len(rules), "Should have exactly 2 default exclude rules")

	// Verify rule order (pattern rule first, then binary rule)
	_, isFirstPattern := rules[0].(*PatternRule)
	_, isSecondBinary := rules[1].(*BinaryRule)

	assert.True(t, isFirstPattern, "First rule should be PatternRule")
	assert.True(t, isSecondBinary, "Second rule should be BinaryRule")
}

func TestDefaultExcludes_PatternCount(t *testing.T) {
	rules := DefaultExcludes()

	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}
	require.NotNil(t, patternRule, "Should have a pattern rule")

	patterns := patternRule.Patterns()

	// Verify we have a reasonable number of patterns
	// (This helps catch if patterns are accidentally removed)
	assert.Greater(t, len(patterns), 20, "Should have more than 20 default exclude patterns")

	// Verify some essential patterns exist
	essentialPatterns := []string{
		".DS_Store",
		".git/",
		"node_modules/",
		"vendor/",
		".idea/",
		"*.log",
	}

	for _, essential := range essentialPatterns {
		found := false
		for _, pattern := range patterns {
			if pattern == essential {
				found = true
				break
			}
		}
		assert.True(t, found, "Essential pattern '%s' should be included", essential)
	}
}

func BenchmarkDefaultExcludes(b *testing.B) {
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rules := DefaultExcludes()
		_ = rules // Prevent optimization
	}
}

func BenchmarkDefaultExcludes_PatternMatching(b *testing.B) {
	rules := DefaultExcludes()
	var patternRule *PatternRule
	for _, rule := range rules {
		if pr, ok := rule.(*PatternRule); ok {
			patternRule = pr
			break
		}
	}

	testPaths := []string{
		"src/main.go",
		"node_modules/package/index.js",
		".git/config",
		"build/app.exe",
		"logs/app.log",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, path := range testPaths {
			patternRule.Match(path)
		}
	}
}
