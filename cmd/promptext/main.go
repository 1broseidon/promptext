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

func main() {
	// Add help and version flags
	help := pflag.BoolP("help", "h", false, "Show help message")
	showVersion := pflag.BoolP("version", "v", false, "Show version information")

	// Define command line flags
	dirPath := pflag.StringP("directory", "d", ".", "Directory path to process")
	extension := pflag.StringP("extension", "e", "", "File extension to filter, e.g., .go,.js")
	exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude, comma-separated")
	infoOnly := pflag.BoolP("info", "i", false, "Only display project summary")
	verbose := pflag.BoolP("verbose", "V", false, "Show full code content in terminal") // Changed to V to avoid conflict with version
	format := pflag.StringP("format", "f", "markdown", "Output format: markdown (or md), xml")
	outFile := pflag.StringP("output", "o", "", "Output file path")
	debug := pflag.BoolP("debug", "D", false, "Enable debug logging")
	gitignore := pflag.BoolP("gitignore", "g", true, "Use .gitignore patterns")
	useDefaultRules := pflag.BoolP("use-default-rules", "u", true, "Use default filtering rules")

	pflag.Parse()

	// Handle help and version flags manually
	if *help {
		fmt.Printf("promptext version %s (%s)\n\n", version, date)
		pflag.Usage()
		os.Exit(0)
	}
	if *showVersion {
		fmt.Printf("promptext version %s (%s)\n", version, date)
		os.Exit(0)
	}

	if err := processor.Run(*dirPath, *extension, *exclude, false, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore, *useDefaultRules); err != nil {
		log.Fatal(err)
	}
}
