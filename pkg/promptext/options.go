package promptext

// Option is a functional option for configuring the extraction process.
type Option func(*config)

// config holds the internal configuration for extraction.
// This is kept private to maintain API stability.
type config struct {
	extensions        []string
	excludes          []string
	gitignore         bool
	useDefaultRules   bool
	relevanceKeywords string
	tokenBudget       int
	format            Format
	verbose           bool
	debug             bool
}

// newDefaultConfig creates a config with sensible defaults.
func newDefaultConfig() *config {
	return &config{
		extensions:      nil,       // nil means all supported extensions
		excludes:        nil,       // nil means no custom excludes
		gitignore:       true,      // respect .gitignore by default
		useDefaultRules: true,      // use built-in filtering rules by default
		tokenBudget:     0,         // 0 means unlimited
		format:          FormatPTX, // PTX is the default format
		verbose:         false,     // quiet by default
		debug:           false,     // no debug logging by default
	}
}

// WithExtensions specifies file extensions to include in the extraction.
// Extensions should include the dot (e.g., ".go", ".js", ".py").
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithExtensions(".go", ".mod"))
func WithExtensions(extensions ...string) Option {
	return func(c *config) {
		c.extensions = extensions
	}
}

// WithExcludes specifies patterns to exclude from the extraction.
// Patterns can be file names, directory names, or glob patterns.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithExcludes("vendor/", "*.test.go", "node_modules/"))
func WithExcludes(patterns ...string) Option {
	return func(c *config) {
		c.excludes = patterns
	}
}

// WithGitIgnore controls whether .gitignore patterns should be respected.
// By default, .gitignore patterns are used.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithGitIgnore(false))
func WithGitIgnore(enabled bool) Option {
	return func(c *config) {
		c.gitignore = enabled
	}
}

// WithDefaultRules controls whether built-in filtering rules should be used.
// Built-in rules filter out common files like binaries, lockfiles, and generated files.
// By default, these rules are enabled.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithDefaultRules(false))
func WithDefaultRules(enabled bool) Option {
	return func(c *config) {
		c.useDefaultRules = enabled
	}
}

// WithRelevance filters and prioritizes files based on keyword relevance.
// Files are scored based on keyword matches in filenames, directories, imports, and content.
// Only files with keyword matches will be included.
//
// Scoring weights:
//   - Filename match: 10x
//   - Directory match: 5x
//   - Import/package match: 3x
//   - Content match: 1x
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithRelevance("auth", "login", "OAuth"))
func WithRelevance(keywords ...string) Option {
	return func(c *config) {
		// Join keywords with spaces for internal processing
		keywordStr := ""
		for i, kw := range keywords {
			if i > 0 {
				keywordStr += " "
			}
			keywordStr += kw
		}
		c.relevanceKeywords = keywordStr
	}
}

// WithTokenBudget sets a maximum token budget for the extraction.
// Files are prioritized by relevance and entry point status, and lower-priority
// files are excluded when the budget would be exceeded.
//
// The budget applies to file contents only; metadata, git info, and directory
// tree are counted separately.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithTokenBudget(8000))
func WithTokenBudget(maxTokens int) Option {
	return func(c *config) {
		c.tokenBudget = maxTokens
	}
}

// WithFormat specifies the output format for the extraction.
// Available formats: FormatPTX, FormatTOON, FormatJSONL, FormatTOONStrict, FormatMarkdown, FormatXML.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatJSONL))
func WithFormat(format Format) Option {
	return func(c *config) {
		c.format = format
	}
}

// WithVerbose enables verbose output logging during extraction.
// This is useful for debugging or understanding what files are being processed.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithVerbose(true))
func WithVerbose(enabled bool) Option {
	return func(c *config) {
		c.verbose = enabled
	}
}

// WithDebug enables debug logging with detailed timing information.
// This is useful for performance analysis and troubleshooting.
//
// Example:
//
//	result, _ := promptext.Extract(".", promptext.WithDebug(true))
func WithDebug(enabled bool) Option {
	return func(c *config) {
		c.debug = enabled
		if enabled {
			c.verbose = true // Debug implies verbose
		}
	}
}
