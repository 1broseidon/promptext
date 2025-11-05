package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/initializer"
	"github.com/1broseidon/promptext/internal/processor"
	"github.com/1broseidon/promptext/internal/update"
	"github.com/spf13/pflag"
)

// Build information. Populated at build time using -ldflags:
//
//	go build -ldflags "-X main.version=v0.2.4 -X main.commit=`git rev-parse HEAD` -X main.date=`date -u +%Y-%m-%d`"
var (
	version = "dev"     // version from git tag
	date    = "unknown" // build date in YYYY-MM-DD format
)

// customUsage provides a modern, well-organized help text for the CLI
func customUsage() {
	fmt.Printf(`promptext %s - Smart code context extractor for AI assistants

USAGE:
    prx [OPTIONS] [DIRECTORY]
    promptext [OPTIONS] [DIRECTORY]

DESCRIPTION:
    promptext analyzes your codebase, filters relevant files, estimates token 
    usage using tiktoken (GPT-3.5/4 compatible), and provides formatted output 
    suitable for AI prompts. Output is automatically copied to clipboard unless 
    disabled with --no-copy.

INPUT OPTIONS:
    -d, --directory DIR        Directory to process (default: current directory)
    -e, --extension LIST       File extensions to include, comma-separated
                               Examples: .go  or  .go,.js,.ts,.py
    -g, --gitignore           Use .gitignore patterns for filtering (default: true)
    -u, --use-default-rules   Use built-in filtering rules for common files (default: true)

FILTERING OPTIONS:
    -x, --exclude LIST        Patterns to exclude, comma-separated
                              Examples: vendor/,node_modules/  or  *.test.go,dist/

OUTPUT OPTIONS:
    -f, --format FORMAT       Output format (default: ptx)
                              • ptx, toon: PTX v2.0 format with enhanced manifest (TOON-based) [default]
                              • jsonl: Machine-friendly JSONL (one JSON object per line)
                              • toon-strict: TOON v1.3 strict compliance (escaped strings)
                              • markdown, md: Human-readable markdown
                              • xml: Machine-parseable XML
    -o, --output FILE         Write output to file instead of clipboard
    -n, --no-copy            Don't copy output to clipboard
    -i, --info               Show only project summary (no file contents)
        --verbose            Display full content in terminal

PROCESSING OPTIONS:
        --dry-run            Preview files that would be processed without reading content
    -q, --quiet              Suppress non-essential output for scripting

RELEVANCE & TOKEN BUDGET:
    -r, --relevant KEYWORDS  Filter and prioritize files by keyword relevance (comma or space separated)
                             Automatically excludes files with no keyword matches
                             Scoring weights: filename (10x), directory (5x), imports (3x), content (1x)
        --max-tokens NUMBER  Maximum token budget for output (excludes lower-priority files when exceeded)
                             Combines with --relevant to include highest-scoring files within budget

DEBUG OPTIONS:
    -D, --debug              Enable debug logging and timing information
    -h, --help               Show this help message
    -v, --version            Show version information

UPDATE OPTIONS:
        --update             Update promptext to the latest version from GitHub
        --check-update       Check if a new version is available without updating

INITIALIZATION OPTIONS:
        --init               Initialize a new .promptext.yml config file with smart defaults
                             Detects project type and suggests framework-specific settings
        --force              Force overwrite of existing config (use with --init)

EXAMPLES:
    # Basic usage - process current directory, copy to clipboard
    prx

    # Process specific project with Go files only
    prx -d /path/to/project -e .go

    # Quick project overview without file contents
    prx -i

    # Export specific file types to XML with debug info
    prx -e .js,.ts,.json -f xml -o project.xml -D

    # Use PTX v2.0 format for AI-optimized structure with enhanced manifest
    prx -f ptx -o project.ptx

    # Use JSONL for machine-friendly processing and pipelines
    prx -f jsonl -o project.jsonl

    # Use strict TOON v1.3 for maximum token compression
    prx -f toon-strict -o project.toon

    # Process with custom exclusions and see output in terminal
    prx -x "vendor/,*.test.go,dist/" -v

    # Analyze without using .gitignore patterns
    prx -g=false -x "node_modules/,target/,build/"

    # Full analysis with debug logging for performance tuning
    prx -D -v -x "test/,spec/,__tests__/"

    # Preview files that would be processed without reading them
    prx --dry-run -e .go,.js

    # Quiet mode for use in scripts (minimal output)
    prx -q -f xml -o output.xml

    # Auto-detect format from output file extension
    prx -o context.ptx                     # Automatically uses PTX format
    prx -o context.toon                    # Automatically uses PTX format (backward compat)
    prx -o context.md                      # Automatically uses markdown format

    # Filter to only authentication-related files
    prx --relevant "auth login OAuth"

    # Filter to database files, limit to 8000 tokens
    prx --relevant "database" --max-tokens 8000

    # Filter to API files, limit to top 5000 tokens worth
    prx -r "api routes handlers" --max-tokens 5000 -o api-context.toon

    # Check for updates and install latest version
    prx --check-update                         # Check only
    prx --update                               # Update to latest version

    # Initialize config file with smart defaults based on project type
    prx --init                                 # Interactive mode
    prx --init --force                         # Overwrite existing config

CONFIGURATION:
    Create a .promptext.yml file in your project root for persistent settings:

    extensions:
      - .go
      - .js
      - .py
    excludes:
      - vendor/
      - node_modules/
    format: toon
    verbose: false

    CLI flags override configuration file settings.

TOKEN ESTIMATION:
    Token counts are estimated using tiktoken (GPT-3.5/GPT-4 compatible) to help
    you understand context window usage. Use --info to see token estimates without
    full file contents.

VERSION: %s (%s)
HOME:    https://github.com/1broseidon/promptext
DOCS:    https://1broseidon.github.io/promptext/

`, version, version, date)
}

func main() {
	// Set custom usage function
	pflag.Usage = customUsage

	// Define command line flags with improved descriptions
	help := pflag.BoolP("help", "h", false, "Show this help message")
	showVersion := pflag.BoolP("version", "v", false, "Show version information and exit")

	// Update options
	checkUpdate := pflag.Bool("check-update", false, "Check if a new version is available")
	doUpdate := pflag.Bool("update", false, "Update to the latest version from GitHub")

	// Initialization options
	initConfig := pflag.Bool("init", false, "Initialize a new .promptext.yml config file with smart defaults")
	forceInit := pflag.Bool("force", false, "Force overwrite of existing config (use with --init)")

	// Input options
	dirPath := pflag.StringP("directory", "d", ".", "Directory to process (default: current directory)")
	extension := pflag.StringP("extension", "e", "", "File extensions to include (comma-separated, e.g., .go,.js,.py)")
	gitignore := pflag.BoolP("gitignore", "g", true, "Use .gitignore patterns for filtering")
	useDefaultRules := pflag.BoolP("use-default-rules", "u", true, "Use built-in filtering rules for common files")

	// Filtering options
	exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude (comma-separated, e.g., vendor/,*.test.go)")

	// Output options
	format := pflag.StringP("format", "f", "ptx", "Output format: ptx, toon, jsonl, toon-strict, markdown, md, or xml (default: ptx)")
	outFile := pflag.StringP("output", "o", "", "Write output to file instead of clipboard")
	noCopy := pflag.BoolP("no-copy", "n", false, "Don't copy output to clipboard")
	infoOnly := pflag.BoolP("info", "i", false, "Show only project summary without file contents")
	verbose := pflag.Bool("verbose", false, "Display full content in terminal while processing")

	// Processing options
	dryRun := pflag.Bool("dry-run", false, "Preview files that would be processed without reading content")
	quiet := pflag.BoolP("quiet", "q", false, "Suppress non-essential output for scripting")

	// Relevance and token budget options
	relevant := pflag.StringP("relevant", "r", "", "Keywords to prioritize files (comma or space separated, multi-factor scoring)")
	maxTokens := pflag.Int("max-tokens", 0, "Maximum token budget for output (excludes lower-priority files when exceeded)")
	explainSelection := pflag.Bool("explain-selection", false, "Show detailed priority scoring breakdown for file selection")

	// Debug options
	debug := pflag.BoolP("debug", "D", false, "Enable debug logging and timing information")

	pflag.Parse()

	// Handle help and version flags
	if *help {
		customUsage()
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("promptext version %s (%s)\n", version, date)
		os.Exit(0)
	}

	// Handle update flags
	if *checkUpdate {
		available, latestVersion, err := update.CheckForUpdate(version)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error checking for updates: %v\n", err)
			os.Exit(1)
		}
		if available {
			fmt.Printf("A new version is available: %s (current: %s)\n", latestVersion, version)
			fmt.Println("Run 'promptext --update' to update to the latest version")
		} else {
			fmt.Printf("You are running the latest version (%s)\n", version)
		}
		os.Exit(0)
	}

	if *doUpdate {
		if err := update.Update(version, true); err != nil {
			fmt.Fprintf(os.Stderr, "Error updating: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Handle initialization flag
	if *initConfig {
		// Get absolute path
		absPath, err := filepath.Abs(*dirPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving directory path: %v\n", err)
			os.Exit(1)
		}

		// Create and run initializer
		init := initializer.NewInitializer(absPath, *forceInit, *quiet)
		if err := init.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Error initializing config: %v\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	// Automatic update check (non-blocking, silently fails on network issues)
	// Only runs during normal operation, not for update/version/help commands
	go update.CheckAndNotifyUpdate(version)

	// Handle positional argument for directory
	args := pflag.Args()
	if len(args) > 0 {
		*dirPath = args[0]
	}

	// Format auto-detection from output file extension
	if *outFile != "" {
		ext := strings.ToLower(filepath.Ext(*outFile))
		detectedFormat := ""
		switch ext {
		case ".ptx":
			detectedFormat = "ptx"
		case ".toon":
			detectedFormat = "toon" // Maps to PTX for backward compatibility
		case ".md", ".markdown":
			detectedFormat = "markdown"
		case ".xml":
			detectedFormat = "xml"
		}

		// Check for format conflict and warn
		if detectedFormat != "" && *format != detectedFormat {
			// User explicitly set format flag
			formatFlag := pflag.Lookup("format")
			if formatFlag.Changed {
				// Warn about conflict
				fmt.Fprintf(os.Stderr, "⚠️  Warning: format flag '%s' conflicts with output extension '%s' - using '%s' (flag takes precedence)\n", *format, ext, *format)
			} else {
				// Auto-detect format from extension since flag wasn't explicitly set
				*format = detectedFormat
			}
		}
	}

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore, *useDefaultRules, *dryRun, *quiet, *relevant, *maxTokens, *explainSelection); err != nil {
		log.Fatal(err)
	}
}
