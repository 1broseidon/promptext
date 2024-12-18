package main

import (
	"log"
	"os"

	"github.com/spf13/pflag"
	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Add help flag
	help := pflag.BoolP("help", "h", false, "Show help message")
	// Define command line flags
	dirPath := pflag.StringP("directory", "d", ".", "Directory path to process")
	extension := pflag.StringP("extension", "e", "", "File extension to filter, e.g., .go,.js")
	exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude, comma-separated")
	infoOnly := pflag.BoolP("info", "i", false, "Only display project summary")
	verbose := pflag.BoolP("verbose", "v", false, "Show full code content in terminal")
	format := pflag.StringP("format", "f", "markdown", "Output format: markdown (or md), xml")
	outFile := pflag.StringP("output", "o", "", "Output file path")
	debug := pflag.BoolP("debug", "D", false, "Enable debug logging")
	gitignore := pflag.BoolP("gitignore", "g", true, "Use .gitignore patterns")

	pflag.Parse()

	// Handle help flag manually
	if *help {
		pflag.Usage()
		os.Exit(0)
	}

	if err := processor.Run(*dirPath, *extension, *exclude, false, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
