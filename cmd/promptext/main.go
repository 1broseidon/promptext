package main

import (
	"flag"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("dir", ".", "Directory path to process")
	extension := flag.String("ext", "", "File extension to filter (e.g., .go,.js)")
	exclude := flag.String("exclude", "", "Patterns to exclude (comma-separated)")
	noCopy := flag.Bool("no-copy", false, "Disable automatic copying to clipboard")
	infoOnly := flag.Bool("info", false, "Only display project summary")
	verbose := flag.Bool("verbose", false, "Show full code content in terminal")
	format := flag.String("format", "markdown", "Output format (markdown, xml, json)")
	outFile := flag.String("out", "", "Output file path (if specified, output will be written to file instead of clipboard)")
	debug := flag.Bool("debug", false, "Enable debug logging")

	flag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug); err != nil {
		log.Fatal(err)
	}
}
