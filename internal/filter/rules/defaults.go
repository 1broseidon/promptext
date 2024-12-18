package rules

import "github.com/1broseidon/promptext/internal/filter/types"

func DefaultExcludes() []types.Rule {
	return []types.Rule{
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
		NewBinaryRule(),
	}
}
