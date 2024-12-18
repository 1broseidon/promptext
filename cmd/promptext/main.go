package main

import (
	"flag"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("directory", ".", "Directory path to process (-directory, -d)")
	flag.StringVar(dirPath, "d", ".", "Directory path to process (-directory, -d)")
	
	extension := flag.String("extension", "", "File extension to filter, e.g., .go,.js (-extension, -e)")
	flag.StringVar(extension, "e", "", "File extension to filter, e.g., .go,.js (-extension, -e)")
	
	exclude := flag.String("exclude", "", "Patterns to exclude, comma-separated (-exclude, -x)")
	flag.StringVar(exclude, "x", "", "Patterns to exclude, comma-separated (-exclude, -x)")
	
	noCopy := flag.Bool("nocopy", false, "Disable automatic copying to clipboard (-nocopy, -n)")
	flag.BoolVar(noCopy, "n", false, "Disable automatic copying to clipboard (-nocopy, -n)")
	
	infoOnly := flag.Bool("info", false, "Only display project summary (-info, -i)")
	flag.BoolVar(infoOnly, "i", false, "Only display project summary (-info, -i)")
	
	verbose := flag.Bool("verbose", false, "Show full code content in terminal (-verbose, -v)")
	flag.BoolVar(verbose, "v", false, "Show full code content in terminal (-verbose, -v)")
	
	format := flag.String("format", "markdown", "Output format: markdown, xml, json (-format, -f)")
	flag.StringVar(format, "f", "markdown", "Output format: markdown, xml, json (-format, -f)")
	
	outFile := flag.String("output", "", "Output file path (-output, -o)")
	flag.StringVar(outFile, "o", "", "Output file path (-output, -o)")
	
	debug := flag.Bool("debug", false, "Enable debug logging (-debug, -D)")
	flag.BoolVar(debug, "D", false, "Enable debug logging (-debug, -D)")
	
	gitignore := flag.Bool("gitignore", true, "Use .gitignore patterns (-gitignore, -g)")
	flag.BoolVar(gitignore, "g", true, "Use .gitignore patterns (-gitignore, -g)")

	flag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
