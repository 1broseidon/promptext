package main

import (
	"flag"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("d", ".", "Directory path to process (-d or --dir)")
	flag.StringVar(dirPath, "dir", ".", "Directory path to process (-d or --dir)")
	
	extension := flag.String("e", "", "File extension to filter, e.g., .go,.js (-e or --ext)")
	flag.StringVar(extension, "ext", "", "File extension to filter, e.g., .go,.js (-e or --ext)")
	
	exclude := flag.String("x", "", "Patterns to exclude, comma-separated (-x or --exclude)")
	flag.StringVar(exclude, "exclude", "", "Patterns to exclude, comma-separated (-x or --exclude)")
	
	noCopy := flag.Bool("n", false, "Disable automatic copying to clipboard (-n or --no-copy)")
	flag.BoolVar(noCopy, "no-copy", false, "Disable automatic copying to clipboard (-n or --no-copy)")
	
	infoOnly := flag.Bool("i", false, "Only display project summary (-i or --info)")
	flag.BoolVar(infoOnly, "info", false, "Only display project summary (-i or --info)")
	
	verbose := flag.Bool("v", false, "Show full code content in terminal (-v or --verbose)")
	flag.BoolVar(verbose, "verbose", false, "Show full code content in terminal (-v or --verbose)")
	
	format := flag.String("f", "markdown", "Output format: markdown, xml, json (-f or --format)")
	flag.StringVar(format, "format", "markdown", "Output format: markdown, xml, json (-f or --format)")
	
	outFile := flag.String("o", "", "Output file path (-o or --out)")
	flag.StringVar(outFile, "out", "", "Output file path (-o or --out)")
	
	debug := flag.Bool("D", false, "Enable debug logging (-D or --debug)")
	flag.BoolVar(debug, "debug", false, "Enable debug logging (-D or --debug)")
	
	gitignore := flag.Bool("g", true, "Use .gitignore patterns (-g or --gitignore)")
	flag.BoolVar(gitignore, "gitignore", true, "Use .gitignore patterns (-g or --gitignore)")

	flag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
