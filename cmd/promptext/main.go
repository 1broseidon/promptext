package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/initializer"
	"github.com/1broseidon/promptext/internal/processor"
	"github.com/1broseidon/promptext/internal/update"
	"github.com/1broseidon/promptext/pkg/promptext"
	"github.com/atotto/clipboard"
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
	customUsageWithWriter(os.Stdout)
}

func customUsageWithWriter(w io.Writer) {
	fmt.Fprintf(w, `promptext %s - Smart code context extractor for AI assistants

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

type initializerRunner interface {
	Run() error
}

type initializerFactory func(root string, force bool, quiet bool) initializerRunner

type processorFunc func(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool, useDefaultRules bool, dryRun bool, quiet bool, relevanceKeywords string, maxTokens int, explainSelection bool) error

// runWithLibrary uses the promptext library for extraction instead of calling processor.Run() directly.
// This provides a thin CLI wrapper around the library while maintaining backward compatibility.
func runWithLibrary(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool, useDefaultRules bool, dryRun bool, quiet bool, relevanceKeywords string, maxTokens int, explainSelection bool) error {
	// For dry-run and explain-selection modes, fall back to processor.Run() as they use internal-only features
	if dryRun || explainSelection {
		return processor.Run(dirPath, extension, exclude, noCopy, infoOnly, verbose, outputFormat, outFile, debug, gitignore, useDefaultRules, dryRun, quiet, relevanceKeywords, maxTokens, explainSelection)
	}

	// Build library options from CLI flags
	opts := []promptext.Option{}

	// Extensions
	if extension != "" {
		exts := strings.Split(extension, ",")
		opts = append(opts, promptext.WithExtensions(exts...))
	}

	// Excludes
	if exclude != "" {
		excludes := strings.Split(exclude, ",")
		opts = append(opts, promptext.WithExcludes(excludes...))
	}

	// GitIgnore
	opts = append(opts, promptext.WithGitIgnore(gitignore))

	// Default rules
	opts = append(opts, promptext.WithDefaultRules(useDefaultRules))

	// Relevance keywords
	if relevanceKeywords != "" {
		keywords := strings.FieldsFunc(relevanceKeywords, func(r rune) bool {
			return r == ',' || r == ' '
		})
		opts = append(opts, promptext.WithRelevance(keywords...))
	}

	// Token budget
	if maxTokens > 0 {
		opts = append(opts, promptext.WithTokenBudget(maxTokens))
	}

	// Format
	opts = append(opts, promptext.WithFormat(promptext.Format(outputFormat)))

	// Verbose and debug
	if debug {
		opts = append(opts, promptext.WithDebug(true))
	} else if verbose {
		opts = append(opts, promptext.WithVerbose(true))
	}

	// Extract using the library
	result, err := promptext.Extract(dirPath, opts...)
	if err != nil {
		return err
	}

	// Handle info-only mode
	if infoOnly {
		if quiet {
			fmt.Printf("files=%d tokens=%d\n", len(result.ProjectOutput.Files), result.TokenCount)
		} else {
			// Format project info display
			var info strings.Builder

			// Project header
			if result.ProjectOutput.Metadata != nil && result.ProjectOutput.Metadata.Language != "" {
				info.WriteString(fmt.Sprintf("📦 %s", getProjectDisplayName(dirPath)))
				info.WriteString(fmt.Sprintf(" (%s)", result.ProjectOutput.Metadata.Language))
			} else {
				info.WriteString(fmt.Sprintf("📦 %s", getProjectDisplayName(dirPath)))
			}

			// File and token count
			fileCount := len(result.ProjectOutput.Files)
			totalFileCount := fileCount + result.ExcludedFiles

			if result.ExcludedFiles > 0 {
				info.WriteString(fmt.Sprintf("\nIncluded: %d/%d files • ~%s tokens",
					fileCount, totalFileCount, formatTokenCount(result.TokenCount)))
				if result.TotalTokens > result.TokenCount {
					info.WriteString(fmt.Sprintf("\nFull project: %d files • ~%s tokens",
						totalFileCount, formatTokenCount(result.TotalTokens)))
				}
			} else {
				info.WriteString(fmt.Sprintf("\nIncluded: %d files • ~%s tokens",
					fileCount, formatTokenCount(result.TokenCount)))
			}

			fmt.Printf("\033[32m%s\033[0m\n", info.String())
		}
		return nil
	}

	// Build exclusion message if files were excluded
	exclusionMsg := ""
	if result.ExcludedFiles > 0 {
		if quiet {
			exclusionMsg = fmt.Sprintf(" excluded=%d", result.ExcludedFiles)
		} else {
			// Build detailed exclusion summary
			var summary strings.Builder
			summary.WriteString(fmt.Sprintf("\n\n⚠️ Excluded %d files due to token budget:\n", result.ExcludedFiles))

			// Show first 5 excluded files with token counts
			displayCount := 5
			if len(result.ExcludedFileList) < displayCount {
				displayCount = len(result.ExcludedFileList)
			}

			totalExcludedTokens := 0
			for i := 0; i < displayCount; i++ {
				excluded := result.ExcludedFileList[i]
				summary.WriteString(fmt.Sprintf("• %s (~%s tokens)\n", excluded.Path, formatTokenCount(excluded.Tokens)))
				totalExcludedTokens += excluded.Tokens
			}

			// Add summary for remaining files
			if len(result.ExcludedFileList) > displayCount {
				remaining := len(result.ExcludedFileList) - displayCount
				remainingTokens := 0
				for i := displayCount; i < len(result.ExcludedFileList); i++ {
					remainingTokens += result.ExcludedFileList[i].Tokens
				}
				totalExcludedTokens += remainingTokens
				summary.WriteString(fmt.Sprintf("... and %d more files (~%s tokens)\n", remaining, formatTokenCount(remainingTokens)))
			}

			summary.WriteString(fmt.Sprintf("Total excluded: ~%s tokens", formatTokenCount(totalExcludedTokens)))
			exclusionMsg = summary.String()
		}
	}

	// Format basic project info for display
	var info strings.Builder
	if result.ProjectOutput.Metadata != nil && result.ProjectOutput.Metadata.Language != "" {
		info.WriteString(fmt.Sprintf("📦 %s", getProjectDisplayName(dirPath)))
		info.WriteString(fmt.Sprintf(" (%s)", result.ProjectOutput.Metadata.Language))
	} else {
		info.WriteString(fmt.Sprintf("📦 %s", getProjectDisplayName(dirPath)))
	}

	fileCount := len(result.ProjectOutput.Files)
	totalFileCount := fileCount + result.ExcludedFiles

	if result.ExcludedFiles > 0 {
		info.WriteString(fmt.Sprintf("\nIncluded: %d/%d files • ~%s tokens",
			fileCount, totalFileCount, formatTokenCount(result.TokenCount)))
		if result.TotalTokens > result.TokenCount {
			info.WriteString(fmt.Sprintf("\nFull project: %d files • ~%s tokens",
				totalFileCount, formatTokenCount(result.TotalTokens)))
		}
	} else {
		info.WriteString(fmt.Sprintf("\nIncluded: %d files • ~%s tokens",
			fileCount, formatTokenCount(result.TokenCount)))
	}

	infoFormatted := info.String()

	// Handle output
	if outFile != "" {
		if err := os.WriteFile(outFile, []byte(result.FormattedOutput), 0644); err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
		if quiet {
			fmt.Printf("written=%s format=%s files=%d tokens=%d%s\n", outFile, outputFormat, len(result.ProjectOutput.Files), result.TokenCount, exclusionMsg)
		} else {
			fmt.Printf("\033[32m%s%s\n\n✓ Code context written to %s (%s format)\033[0m\n", infoFormatted, exclusionMsg, outFile, outputFormat)
		}
	} else if !noCopy {
		if err := clipboard.WriteAll(result.FormattedOutput); err != nil {
			if !quiet {
				fmt.Printf("Warning: Failed to copy to clipboard: %v\n", err)
			}
			if quiet {
				return fmt.Errorf("clipboard copy failed")
			}
		} else {
			if quiet {
				fmt.Printf("clipboard=ok format=%s files=%d tokens=%d%s\n", outputFormat, len(result.ProjectOutput.Files), result.TokenCount, exclusionMsg)
			} else {
				fmt.Printf("\033[32m%s%s\n\n✓ Copied to clipboard!\033[0m\n", infoFormatted, exclusionMsg)
			}
		}
	}

	return nil
}

// formatTokenCount formats token count with comma separators for readability
func formatTokenCount(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	// Add comma separators for thousands
	str := fmt.Sprintf("%d", tokens)
	var result strings.Builder
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// getProjectDisplayName returns the proper display name for a project directory.
// It resolves "." to the actual directory name.
func getProjectDisplayName(dirPath string) string {
	// If dirPath is "." or relative, resolve to absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		// Fallback to original if we can't resolve
		return filepath.Base(dirPath)
	}
	return filepath.Base(absPath)
}

type cliDeps struct {
	stdout         io.Writer
	stderr         io.Writer
	usage          func()
	checkForUpdate func(string) (bool, string, error)
	updater        func(string, bool) error
	notifyUpdate   func(string)
	newInitializer initializerFactory
	processorRun   processorFunc
	absPath        func(string) (string, error)
}

func defaultCLIDeps() cliDeps {
	return cliDeps{
		stdout:         os.Stdout,
		stderr:         os.Stderr,
		usage:          customUsage,
		checkForUpdate: update.CheckForUpdate,
		updater:        update.Update,
		notifyUpdate:   update.CheckAndNotifyUpdate,
		newInitializer: func(root string, force bool, quiet bool) initializerRunner {
			return initializer.NewInitializer(root, force, quiet)
		},
		processorRun: runWithLibrary, // Use library instead of processor.Run
		absPath:      filepath.Abs,
	}
}

func run(args []string, deps cliDeps) int {
	if deps.stdout == nil {
		deps.stdout = os.Stdout
	}
	if deps.stderr == nil {
		deps.stderr = os.Stderr
	}
	if deps.usage == nil {
		deps.usage = customUsage
	}
	if deps.checkForUpdate == nil {
		deps.checkForUpdate = update.CheckForUpdate
	}
	if deps.updater == nil {
		deps.updater = update.Update
	}
	if deps.notifyUpdate == nil {
		deps.notifyUpdate = update.CheckAndNotifyUpdate
	}
	if deps.newInitializer == nil {
		deps.newInitializer = func(root string, force bool, quiet bool) initializerRunner {
			return initializer.NewInitializer(root, force, quiet)
		}
	}
	if deps.processorRun == nil {
		deps.processorRun = processor.Run
	}
	if deps.absPath == nil {
		deps.absPath = filepath.Abs
	}

	flagSet := pflag.NewFlagSet("promptext", pflag.ContinueOnError)
	flagSet.SetOutput(deps.stderr)
	flagSet.Usage = deps.usage

	help := flagSet.BoolP("help", "h", false, "Show this help message")
	showVersion := flagSet.BoolP("version", "v", false, "Show version information and exit")

	checkUpdate := flagSet.Bool("check-update", false, "Check if a new version is available")
	doUpdate := flagSet.Bool("update", false, "Update to the latest version from GitHub")

	initConfig := flagSet.Bool("init", false, "Initialize a new .promptext.yml config file with smart defaults")
	forceInit := flagSet.Bool("force", false, "Force overwrite of existing config (use with --init)")

	dirPath := flagSet.StringP("directory", "d", ".", "Directory to process (default: current directory)")
	extension := flagSet.StringP("extension", "e", "", "File extensions to include (comma-separated, e.g., .go,.js,.py)")
	gitignore := flagSet.BoolP("gitignore", "g", true, "Use .gitignore patterns for filtering")
	useDefaultRules := flagSet.BoolP("use-default-rules", "u", true, "Use built-in filtering rules for common files")

	exclude := flagSet.StringP("exclude", "x", "", "Patterns to exclude (comma-separated, e.g., vendor/,*.test.go)")

	format := flagSet.StringP("format", "f", "ptx", "Output format: ptx, toon, jsonl, toon-strict, markdown, md, or xml (default: ptx)")
	outFile := flagSet.StringP("output", "o", "", "Write output to file instead of clipboard")
	noCopy := flagSet.BoolP("no-copy", "n", false, "Don't copy output to clipboard")
	infoOnly := flagSet.BoolP("info", "i", false, "Show only project summary without file contents")
	verbose := flagSet.Bool("verbose", false, "Display full content in terminal while processing")

	dryRun := flagSet.Bool("dry-run", false, "Preview files that would be processed without reading content")
	quiet := flagSet.BoolP("quiet", "q", false, "Suppress non-essential output for scripting")

	relevant := flagSet.StringP("relevant", "r", "", "Keywords to prioritize files (comma or space separated, multi-factor scoring)")
	maxTokens := flagSet.Int("max-tokens", 0, "Maximum token budget for output (excludes lower-priority files when exceeded)")
	explainSelection := flagSet.Bool("explain-selection", false, "Show detailed priority scoring breakdown for file selection")

	debug := flagSet.BoolP("debug", "D", false, "Enable debug logging and timing information")

	if err := flagSet.Parse(args); err != nil {
		if errors.Is(err, pflag.ErrHelp) {
			deps.usage()
			return 0
		}
		return 2
	}

	if *help {
		deps.usage()
		return 0
	}
	if *showVersion {
		fmt.Fprintf(deps.stdout, "promptext version %s (%s)\n", version, date)
		return 0
	}

	if *checkUpdate {
		available, latestVersion, err := deps.checkForUpdate(version)
		if err != nil {
			fmt.Fprintf(deps.stderr, "Error checking for updates: %v\n", err)
			return 1
		}
		if available {
			fmt.Fprintf(deps.stdout, "A new version is available: %s (current: %s)\n", latestVersion, version)
			fmt.Fprintln(deps.stdout, "Run 'promptext --update' to update to the latest version")
		} else {
			fmt.Fprintf(deps.stdout, "You are running the latest version (%s)\n", version)
		}
		return 0
	}

	if *doUpdate {
		if err := deps.updater(version, true); err != nil {
			fmt.Fprintf(deps.stderr, "Error updating: %v\n", err)
			return 1
		}
		return 0
	}

	if *initConfig {
		absPath, err := deps.absPath(*dirPath)
		if err != nil {
			fmt.Fprintf(deps.stderr, "Error resolving directory path: %v\n", err)
			return 1
		}

		init := deps.newInitializer(absPath, *forceInit, *quiet)
		if err := init.Run(); err != nil {
			fmt.Fprintf(deps.stderr, "Error initializing config: %v\n", err)
			return 1
		}
		return 0
	}

	if deps.notifyUpdate != nil {
		go deps.notifyUpdate(version)
	}

	positional := flagSet.Args()
	if len(positional) > 0 {
		*dirPath = positional[0]
	}

	if *outFile != "" {
		ext := strings.ToLower(filepath.Ext(*outFile))
		detectedFormat := ""
		switch ext {
		case ".ptx":
			detectedFormat = "ptx"
		case ".toon":
			detectedFormat = "toon"
		case ".md", ".markdown":
			detectedFormat = "markdown"
		case ".xml":
			detectedFormat = "xml"
		}

		if detectedFormat != "" && *format != detectedFormat {
			formatFlag := flagSet.Lookup("format")
			if formatFlag != nil && formatFlag.Changed {
				fmt.Fprintf(deps.stderr, "⚠️  Warning: format flag '%s' conflicts with output extension '%s' - using '%s' (flag takes precedence)\n", *format, ext, *format)
			} else {
				*format = detectedFormat
			}
		}
	}

	if err := deps.processorRun(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore, *useDefaultRules, *dryRun, *quiet, *relevant, *maxTokens, *explainSelection); err != nil {
		fmt.Fprintf(deps.stderr, "%v\n", err)
		return 1
	}
	return 0
}

func main() {
	os.Exit(run(os.Args[1:], defaultCLIDeps()))
}
