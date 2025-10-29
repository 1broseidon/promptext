package rules

import "github.com/1broseidon/promptext/internal/filter/types"

func DefaultExcludes() []types.Rule {
	return []types.Rule{
		// Pattern-based exclusions (fastest - checked first)
		NewPatternRule([]string{
			".DS_Store",

			// Version control
			".git/",
			".git*",
			".svn/",
			".hg/",

			// Dependencies and packages
			"node_modules/",
			"vendor/",
			"bower_components/",
			"jspm_packages/",

			// IDE and editor
			".idea/",
			".vscode/",
			".vs/",
			"*.sublime-*",

			// Build and output
			"dist/",
			"build/",
			"out/",
			"bin/",
			"target/",

			// Cache directories
			"__pycache__/",
			".pytest_cache/",
			".sass-cache/",
			".npm/",
			".yarn/",
			".pnp.*",

			// Test coverage
			"coverage/",
			".nyc_output/",

			// Infrastructure
			".terraform/",
			".vagrant/",

			// Logs and temp
			"logs/",
			"*.log",
			"tmp/",
			"temp/",
		}, types.Exclude),

		// Lock file detection (multi-layered, ordered by confidence)
		// Layer 1: Signature-based (99% confidence - most reliable)
		NewLockFileRule(),

		// Layer 2: Ecosystem-aware (95% confidence - context-aware)
		NewEcosystemRule("."),

		// Layer 3: Generated file detection (85% confidence - heuristic)
		NewGeneratedFileRule(1), // 1MB threshold

		// Binary file detection (always reliable)
		NewBinaryRule(),
	}
}
