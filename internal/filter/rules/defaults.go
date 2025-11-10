package rules

import "github.com/1broseidon/promptext/internal/filter/types"

func DefaultExcludes() []types.Rule {
	return []types.Rule{
		// Pattern-based exclusions (fastest - checked first)
		NewPatternRule([]string{
			".DS_Store",
			"Thumbs.db",

			// Version control
			".git/",
			".git*",
			".svn/",
			".hg/",
			".bzr/",

			// ========== JavaScript/TypeScript/Node.js ==========
			// Dependencies
			"node_modules/",
			"bower_components/",
			"jspm_packages/",

			// Build outputs
			"dist/",
			"build/",
			"out/",
			".output/", // Nuxt 3

			// Framework-specific
			".next/",           // Next.js
			".nuxt/",           // Nuxt.js
			".svelte-kit/",     // SvelteKit
			".remix/",          // Remix
			".astro/",          // Astro
			".docusaurus/",     // Docusaurus
			".vuepress/",       // VuePress
			"_site/",           // Jekyll

			// Bundler caches
			".parcel-cache/",
			".turbo/",
			".rollup.cache/",

			// Package manager
			".npm/",
			".yarn/",
			".pnp.*",
			".pnpm-debug.log",

			// Deployment
			".vercel/",
			".netlify/",

			// ========== Python ==========
			// Virtual environments
			".venv/",
			"venv/",
			"env/",
			".env/",
			"virtualenv/",
			".virtualenv/",

			// Cache and build
			"__pycache__/",
			".pytest_cache/",
			".mypy_cache/",
			".ruff_cache/",
			".tox/",
			".nox/",
			"*.egg-info/",
			".eggs/",
			".Python",
			"pip-wheel-metadata/",
			".ipynb_checkpoints/",
			"htmlcov/",
			".coverage",

			// ========== Ruby ==========
			".bundle/",
			"vendor/bundle/",
			".gem/",

			// ========== PHP ==========
			"vendor/", // Composer (also used by Go)

			// ========== Go ==========
			"vendor/", // Go modules

			// ========== Rust ==========
			"target/", // Cargo

			// ========== Java/Kotlin/Scala ==========
			"target/",          // Maven
			".gradle/",         // Gradle
			".mvn/",            // Maven wrapper
			"out/",             // IntelliJ

			// ========== C#/.NET ==========
			"bin/",
			"obj/",
			"packages/",
			"*.nupkg",

			// ========== Swift/iOS ==========
			".build/",
			".swiftpm/",
			"DerivedData/",
			"Pods/",
			"xcuserdata/",

			// ========== Dart/Flutter ==========
			".dart_tool/",
			".flutter-plugins",
			".flutter-plugins-dependencies",

			// ========== Elixir ==========
			"_build/",
			"deps/",
			".elixir_ls/",

			// ========== Android ==========
			".externalNativeBuild/",
			".cxx/",
			"local.properties",

			// ========== IDE and Editor ==========
			".idea/",
			".vscode/",
			".vs/",
			"*.sublime-*",
			"*.swp",
			"*.swo",
			"*~",

			// ========== Test Coverage ==========
			"coverage/",
			".nyc_output/",
			"test-results/",
			"junit/",
			".phpunit.result.cache",

			// ========== Infrastructure/DevOps ==========
			".terraform/",
			".vagrant/",
			".docker/",

			// ========== General Cache ==========
			".cache/",
			".temp/",
			".tmp/",
			".sass-cache/",

			// ========== Logs and Temp ==========
			"logs/",
			"*.log",
			"tmp/",
			"temp/",

			// ========== Database ==========
			"*.db-shm",
			"*.db-wal",
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
