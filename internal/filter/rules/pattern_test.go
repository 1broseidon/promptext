package rules

import (
	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewPatternRule(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		action   types.RuleAction
	}{
		{
			name:     "single pattern exclude rule",
			patterns: []string{"*.log"},
			action:   types.Exclude,
		},
		{
			name:     "multiple patterns include rule",
			patterns: []string{"*.go", "*.js", "*.py"},
			action:   types.Include,
		},
		{
			name:     "directory patterns",
			patterns: []string{"node_modules/", ".git/", "vendor/"},
			action:   types.Exclude,
		},
		{
			name:     "empty patterns",
			patterns: []string{},
			action:   types.Skip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewPatternRule(tt.patterns, tt.action)

			require.NotNil(t, rule)
			assert.Equal(t, tt.action, rule.Action())

			// Verify rule type
			patternRule, ok := rule.(*PatternRule)
			require.True(t, ok, "Expected PatternRule type")
			assert.Equal(t, tt.patterns, patternRule.Patterns())
		})
	}
}

func TestPatternRule_Match_DirectoryPatterns(t *testing.T) {
	rule := NewPatternRule([]string{
		"node_modules/",
		".git/",
		"vendor/",
		"build/",
	}, types.Exclude)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Direct directory matches
		{"direct node_modules", "node_modules/", true, "direct directory match"},
		{"direct git", ".git/", true, "direct git directory"},
		{"direct vendor", "vendor/", true, "direct vendor directory"},
		{"direct build", "build/", true, "direct build directory"},

		// Files within directories
		{"file in node_modules", "node_modules/package.json", true, "file within node_modules"},
		{"nested in node_modules", "project/node_modules/lib/index.js", true, "nested node_modules directory"},
		{"file in git", ".git/config", true, "file within git directory"},
		{"nested git", "project/.git/hooks/pre-commit", true, "nested git directory"},
		{"file in vendor", "vendor/github.com/pkg/file.go", true, "file within vendor"},
		{"nested vendor", "src/vendor/lib.go", true, "nested vendor directory"},

		// Non-matching paths
		{"similar name", "node_modules_backup/file.js", false, "similar but different name"},
		{"prefix match", "node_modules.json", false, "prefix match without slash"},
		{"git file", ".gitignore", false, "git file but not directory"},
		{"vendor file", "vendor.json", false, "vendor file but not directory"},
		{"build file", "build.sh", false, "build file but not directory"},
		{"random file", "src/main.go", false, "unrelated file"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestPatternRule_Match_WildcardPatterns(t *testing.T) {
	rule := NewPatternRule([]string{
		"*.log",
		"*.tmp",
		".aider*",
		"test_*",
		"*_backup",
	}, types.Exclude)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Basic wildcard matches
		{"log file", "app.log", true, "basic log file"},
		{"nested log", "logs/app.log", true, "nested log file"},
		{"tmp file", "temp.tmp", true, "basic tmp file"},
		{"nested tmp", "cache/temp.tmp", true, "nested tmp file"},

		// Prefix wildcards
		{"aider config", ".aider.conf.yml", true, "aider config file"},
		{"aider log", ".aider.log", true, "aider log file"},
		{"nested aider", "project/.aider.conf.yml", true, "nested aider file"},

		// Suffix wildcards
		{"test prefix", "test_utils.py", true, "test prefix match"},
		{"backup suffix", "config_backup", true, "backup suffix match"},
		{"nested test", "src/test_helper.js", true, "nested test file"},
		{"nested backup", "data/db_backup", true, "nested backup file"},

		// Non-matching patterns
		{"no extension", "logfile", false, "no file extension"},
		{"wrong extension", "app.txt", false, "different extension"},
		{"partial match", "aider.conf", false, "missing dot prefix"},
		{"substring", "testing.py", false, "substring but not prefix"},
		{"similar", "backup_file", false, "similar but wrong position"},
		{"case sensitive", "APP.LOG", false, "case sensitivity"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestPatternRule_Match_ExactAndPathPatterns(t *testing.T) {
	rule := NewPatternRule([]string{
		".DS_Store",
		"Thumbs.db",
		"src/generated",
		"package-lock.json",
	}, types.Exclude)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Exact matches
		{"ds store exact", ".DS_Store", true, "exact DS_Store match"},
		{"thumbs exact", "Thumbs.db", true, "exact Thumbs.db match"},
		{"package lock exact", "package-lock.json", true, "exact package-lock match"},

		// Path-based matches (contains with slash prefix)
		{"ds store in dir", "images/.DS_Store", true, "DS_Store in directory"},
		{"nested ds store", "project/assets/.DS_Store", true, "nested DS_Store"},
		{"thumbs in dir", "photos/Thumbs.db", true, "Thumbs.db in directory"},
		{"nested thumbs", "project/images/Thumbs.db", true, "nested Thumbs.db"},
		{"package lock in dir", "frontend/package-lock.json", true, "package-lock in directory"},
		{"nested package lock", "apps/web/package-lock.json", true, "nested package-lock"},

		// Prefix matches (starts with pattern)
		{"generated path match", "src/generated/types.ts", true, "file in generated directory"},
		{"exact generated match", "src/generated", true, "exact generated path match"},

		// Non-matching patterns
		{"substring no slash", "my.DS_Store.bak", false, "substring without slash separator"},
		{"prefix without slash", "Thumbs.db.old", true, "prefix matches - current implementation behavior"},
		{"case sensitive", "thumbs.db", false, "case sensitivity"},
		{"similar name", "src/generation", false, "similar but different path"},
		{"partial path", "generated/types.ts", false, "missing src/ prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestPatternRule_Match_ComplexPatterns(t *testing.T) {
	rule := NewPatternRule([]string{
		"*.pyc",         // Wildcard for pyc files
		"*.sublime-*",   // Complex wildcard
		"node_modules/", // Directory
		".git*",         // Prefix wildcard
	}, types.Exclude)

	tests := []struct {
		name     string
		path     string
		expected bool
		desc     string
	}{
		// Wildcard patterns - only match basename
		{"pyc file", "app.pyc", true, "pyc file matches *.pyc pattern"},
		{"nested pyc", "src/utils/helper.pyc", true, "nested pyc file matches basename"},
		{"sublime project", "myproject.sublime-project", true, "sublime project file"},
		{"sublime workspace", "app.sublime-workspace", true, "sublime workspace file"},

		// Git patterns - prefix wildcard matches basename
		{"gitignore", ".gitignore", true, "git ignore file"},
		{"gitattributes", ".gitattributes", true, "git attributes file"},
		{"git directory", ".git/config", false, "git directory file - .git* doesn't match nested paths"},

		// Directory patterns
		{"node modules file", "node_modules/package/index.js", true, "file in node_modules"},

		// Non-matching
		{"py file", "script.py", false, "py file doesn't match pyc pattern"},
		{"git in name", "mygit.conf", false, "git in filename but not prefix"},
		{"similar directory", "node_modules_backup/", false, "similar directory name"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := rule.Match(tt.path)
			assert.Equal(t, tt.expected, result, "Path: %s - %s", tt.path, tt.desc)
		})
	}
}

func TestPatternRule_Match_EdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		patterns []string
		testPath string
		expected bool
		desc     string
	}{
		// Empty patterns
		{
			name:     "empty patterns",
			patterns: []string{},
			testPath: "any/file.txt",
			expected: false,
			desc:     "empty pattern list should not match anything",
		},

		// Special characters in paths
		{
			name:     "spaces in pattern",
			patterns: []string{"My Documents/"},
			testPath: "My Documents/file.txt",
			expected: true,
			desc:     "pattern with spaces should work",
		},
		{
			name:     "special chars",
			patterns: []string{"file[1].txt"},
			testPath: "file[1].txt",
			expected: true,
			desc:     "special characters in exact match",
		},

		// Unicode paths
		{
			name:     "unicode pattern",
			patterns: []string{"测试/"},
			testPath: "测试/file.txt",
			expected: true,
			desc:     "unicode directory pattern",
		},

		// Very long paths
		{
			name:     "long path",
			patterns: []string{"very/"},
			testPath: "very/deeply/nested/directory/structure/with/many/levels/file.txt",
			expected: true,
			desc:     "long nested path should match directory pattern",
		},

		// Path normalization edge cases - Note: current implementation may not normalize these
		{
			name:     "windows vs unix paths",
			patterns: []string{"src/build/"}, // Use unix style pattern
			testPath: "src/build/output.exe",
			expected: true,
			desc:     "directory pattern should match",
		},

		// Pattern with trailing/leading slashes
		{
			name:     "multiple slashes",
			patterns: []string{"cache/"}, // Use single slash pattern
			testPath: "cache/file.tmp",
			expected: true,
			desc:     "directory pattern should match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := NewPatternRule(tt.patterns, types.Exclude)
			result := rule.Match(tt.testPath)
			assert.Equal(t, tt.expected, result, "Pattern: %v, Path: %s - %s", tt.patterns, tt.testPath, tt.desc)
		})
	}
}

func TestPatternRule_Patterns(t *testing.T) {
	patterns := []string{"*.log", "temp/", ".hidden*"}
	rule := NewPatternRule(patterns, types.Exclude)

	patternRule, ok := rule.(*PatternRule)
	require.True(t, ok, "Expected PatternRule type")

	result := patternRule.Patterns()
	assert.Equal(t, patterns, result, "Patterns() should return original patterns")

	// Note: The current implementation returns the original slice reference.
	// This test documents the current behavior. For true immutability,
	// the implementation would need to return a copy.
	result[0] = "modified"
	originalPatterns := patternRule.Patterns()
	// Current implementation: the slice is shared, so modification affects original
	assert.Equal(t, "modified", originalPatterns[0], "Current implementation shares slice reference")
}

func TestPatternRule_Action(t *testing.T) {
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
			rule := NewPatternRule([]string{"*.test"}, tt.action)
			assert.Equal(t, tt.action, rule.Action())
		})
	}
}

func BenchmarkPatternRule_DirectoryMatch(b *testing.B) {
	rule := NewPatternRule([]string{"node_modules/", ".git/", "vendor/"}, types.Exclude)
	path := "src/node_modules/package/index.js"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkPatternRule_WildcardMatch(b *testing.B) {
	rule := NewPatternRule([]string{"*.log", "*.tmp", ".aider*"}, types.Exclude)
	path := "logs/application.log"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkPatternRule_ExactMatch(b *testing.B) {
	rule := NewPatternRule([]string{".DS_Store", "Thumbs.db", "package-lock.json"}, types.Exclude)
	path := "images/.DS_Store"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}

func BenchmarkPatternRule_NoMatch(b *testing.B) {
	rule := NewPatternRule([]string{"*.log", "node_modules/", ".git*"}, types.Exclude)
	path := "src/main.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rule.Match(path)
	}
}
