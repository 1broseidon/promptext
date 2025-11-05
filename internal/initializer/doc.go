// Package initializer provides automatic configuration file generation for promptext.
//
// The initializer package implements the --init flag functionality, which creates
// .promptext.yml configuration files with intelligent defaults based on detected
// project types.
//
// # Key Features
//
// - Smart project type detection for 15+ frameworks and languages
// - Framework-specific file extensions and exclusion patterns
// - Multi-language project support (e.g., Go + Node.js)
// - Interactive prompts for user preferences
// - Safe overwrite protection with --force flag support
//
// # Architecture
//
// The package consists of three main components:
//
// 1. Detector (detector.go): Scans the project directory for framework indicators
// and returns detected project types sorted by priority.
//
// 2. Template Generator (templates.go): Creates configuration templates with
// framework-specific extensions and exclusion patterns. Templates are merged
// when multiple project types are detected.
//
// 3. Initializer (initializer.go): Orchestrates the initialization process,
// handles user interaction, and writes the final configuration file.
//
// # Template Merging Strategy
//
// When multiple project types are detected (e.g., Go + Node.js), templates are
// merged using the following strategy:
//
// - Extensions: All unique extensions from all detected types are included
// - Excludes: All unique exclude patterns from all detected types are included
// - Deduplication: Both extensions and excludes are deduplicated to prevent redundancy
// - Priority: Project types are sorted by priority before template generation
//
// For example, a Go + Node.js project will include:
// - Extensions: .go, .mod, .sum, .js, .ts, .json, .md
// - Excludes: vendor/, bin/, node_modules/, dist/, *_test.go, *.test.js, etc.
//
// # Usage
//
// Basic initialization:
//
//	init := initializer.NewInitializer("/path/to/project", false, false)
//	err := init.Run() // Interactive mode with prompts
//
// Quick initialization (no prompts):
//
//	init := initializer.NewInitializer("/path/to/project", false, true)
//	err := init.RunQuick() // Uses defaults, excludes tests
//
// Force overwrite:
//
//	init := initializer.NewInitializer("/path/to/project", true, false)
//	err := init.Run() // Overwrites existing config without asking
//
// # Supported Frameworks
//
// JavaScript/TypeScript:
// - Next.js (next.config.js)
// - Nuxt.js (nuxt.config.js)
// - Vite (vite.config.js)
// - Vue.js (vue.config.js)
// - Angular (angular.json)
// - Svelte (svelte.config.js)
// - Node.js (package.json)
//
// Backend:
// - Go (go.mod)
// - Django (manage.py)
// - Flask (app.py, wsgi.py)
// - Laravel (artisan)
// - Ruby/Rails (Gemfile)
// - PHP (composer.json)
//
// Systems:
// - Rust (Cargo.toml)
// - Java/Maven (pom.xml)
// - Java/Gradle (build.gradle)
// - .NET (*.csproj, *.fsproj, *.vbproj)
//
// # Detection Priority Levels
//
// Project types are prioritized to ensure framework-specific configurations
// take precedence over generic language configurations:
//
// - PriorityFrameworkSpecific (100): Framework-specific (Next.js, Django, etc.)
// - PriorityBuildTool (90): Build tools and configs (Vite, Flask, etc.)
// - PriorityLanguage (80): Language-specific (Go, Rust, Java, .NET)
// - PriorityGeneric (70): Generic language (Python, PHP)
// - PriorityBasic (60): Basic/generic (Node.js)
//
// # Security Considerations
//
// - Path validation: Verifies target directory exists and is a directory
// - Safe overwriting: Requires --force flag or user confirmation
// - Input validation: Prompts validate user input length and format
// - No arbitrary code execution: Only reads filesystem metadata
//
// # Examples
//
// See the test files for comprehensive usage examples:
// - detector_test.go: Project detection examples
// - templates_test.go: Template generation examples
// - initializer_test.go: Full initialization flow examples
package initializer
