package main

import (
	"log"

	"github.com/spf13/pflag"
	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := pflag.StringP("directory", "d", ".", "Directory path to process")
	extension := pflag.StringP("extension", "e", "", "File extension to filter, e.g., .go,.js")
	exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude, comma-separated")
	noCopy := pflag.BoolP("no-copy", "n", false, "Disable automatic copying to clipboard")
	infoOnly := pflag.BoolP("info", "i", false, "Only display project summary")
	verbose := pflag.BoolP("verbose", "v", false, "Show full code content in terminal")
	format := pflag.StringP("format", "f", "markdown", "Output format: markdown, xml, json")
	outFile := pflag.StringP("output", "o", "", "Output file path")
	debug := pflag.BoolP("debug", "D", false, "Enable debug logging")
	gitignore := pflag.BoolP("gitignore", "g", true, "Use .gitignore patterns")

	pflag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
