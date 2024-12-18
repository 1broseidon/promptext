package main

import (
	"flag"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("d", ".", "Directory path to process")
	flag.StringVar(dirPath, "--dir", ".", "Directory path to process")
	
	extension := flag.String("e", "", "File extension to filter (e.g., .go,.js)")
	flag.StringVar(extension, "--ext", "", "File extension to filter (e.g., .go,.js)")
	
	exclude := flag.String("x", "", "Patterns to exclude (comma-separated)")
	flag.StringVar(exclude, "--exclude", "", "Patterns to exclude (comma-separated)")
	
	noCopy := flag.Bool("n", false, "Disable automatic copying to clipboard")
	flag.BoolVar(noCopy, "--no-copy", false, "Disable automatic copying to clipboard")
	
	infoOnly := flag.Bool("i", false, "Only display project summary")
	flag.BoolVar(infoOnly, "--info", false, "Only display project summary")
	
	verbose := flag.Bool("v", false, "Show full code content in terminal")
	flag.BoolVar(verbose, "--verbose", false, "Show full code content in terminal")
	
	format := flag.String("f", "markdown", "Output format (markdown, xml, json)")
	flag.StringVar(format, "--format", "markdown", "Output format (markdown, xml, json)")
	
	outFile := flag.String("o", "", "Output file path (if specified, output will be written to file instead of clipboard)")
	flag.StringVar(outFile, "--out", "", "Output file path (if specified, output will be written to file instead of clipboard)")
	
	debug := flag.Bool("D", false, "Enable debug logging")
	flag.BoolVar(debug, "--debug", false, "Enable debug logging")
	
	gitignore := flag.Bool("g", true, "Use .gitignore patterns (default: true)")
	flag.BoolVar(gitignore, "--gitignore", true, "Use .gitignore patterns (default: true)")

	flag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
