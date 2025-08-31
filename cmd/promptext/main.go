package main

import (
	"fmt"
	"log"
	"os"

	"github.com/1broseidon/promptext/internal/processor"
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
    -f, --format FORMAT       Output format: markdown, md, xml (default: markdown)
    -o, --output FILE         Write output to file instead of clipboard
    -n, --no-copy            Don't copy output to clipboard
    -i, --info               Show only project summary (no file contents)
        --verbose            Display full content in terminal

PROCESSING OPTIONS:
        --dry-run            Preview files that would be processed without reading content
    -q, --quiet              Suppress non-essential output for scripting

DEBUG OPTIONS:
    -D, --debug              Enable debug logging and timing information
    -h, --help               Show this help message
    -v, --version            Show version information

EXAMPLES:
    # Basic usage - process current directory, copy to clipboard
    prx

    # Process specific project with Go files only
    prx -d /path/to/project -e .go

    # Quick project overview without file contents
    prx -i

    # Export specific file types to XML with debug info
    prx -e .js,.ts,.json -f xml -o project.xml -D

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

CONFIGURATION:
    Create a .promptext.yml file in your project root for persistent settings:

    extensions:
      - .go
      - .js
      - .py
    excludes:
      - vendor/
      - node_modules/
    format: markdown
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

	// Input options
	dirPath := pflag.StringP("directory", "d", ".", "Directory to process (default: current directory)")
	extension := pflag.StringP("extension", "e", "", "File extensions to include (comma-separated, e.g., .go,.js,.py)")
	gitignore := pflag.BoolP("gitignore", "g", true, "Use .gitignore patterns for filtering")
	useDefaultRules := pflag.BoolP("use-default-rules", "u", true, "Use built-in filtering rules for common files")

	// Filtering options
	exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude (comma-separated, e.g., vendor/,*.test.go)")

	// Output options
	format := pflag.StringP("format", "f", "markdown", "Output format: markdown, md, or xml")
	outFile := pflag.StringP("output", "o", "", "Write output to file instead of clipboard")
	noCopy := pflag.BoolP("no-copy", "n", false, "Don't copy output to clipboard")
	infoOnly := pflag.BoolP("info", "i", false, "Show only project summary without file contents")
	verbose := pflag.Bool("verbose", false, "Display full content in terminal while processing")

	// Processing options
	dryRun := pflag.Bool("dry-run", false, "Preview files that would be processed without reading content")
	quiet := pflag.BoolP("quiet", "q", false, "Suppress non-essential output for scripting")

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

	// Handle positional argument for directory
	args := pflag.Args()
	if len(args) > 0 {
		*dirPath = args[0]
	}

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore, *useDefaultRules, *dryRun, *quiet); err != nil {
		log.Fatal(err)
	}
}
