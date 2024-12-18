package main

import (
	"flag"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("d", ".", "Directory path to process (-d, --dir)")
	flag.StringVar(dirPath, "dir", ".", "Directory path to process (-d, --dir)")
	
	extension := flag.String("e", "", "File extension to filter (e.g., .go,.js) (-e, --ext)")
	flag.StringVar(extension, "ext", "", "File extension to filter (e.g., .go,.js) (-e, --ext)")
	
	exclude := flag.String("x", "", "Patterns to exclude (comma-separated) (-x, --exclude)")
	flag.StringVar(exclude, "exclude", "", "Patterns to exclude (comma-separated) (-x, --exclude)")
	
	noCopy := flag.Bool("n", false, "Disable automatic copying to clipboard (-n, --no-copy)")
	flag.BoolVar(noCopy, "no-copy", false, "Disable automatic copying to clipboard (-n, --no-copy)")
	
	infoOnly := flag.Bool("i", false, "Only display project summary (-i, --info)")
	flag.BoolVar(infoOnly, "info", false, "Only display project summary (-i, --info)")
	
	verbose := flag.Bool("v", false, "Show full code content in terminal (-v, --verbose)")
	flag.BoolVar(verbose, "verbose", false, "Show full code content in terminal (-v, --verbose)")
	
	format := flag.String("f", "markdown", "Output format (markdown, xml, json) (-f, --format)")
	flag.StringVar(format, "format", "markdown", "Output format (markdown, xml, json) (-f, --format)")
	
	outFile := flag.String("o", "", "Output file path (-o, --out)")
	flag.StringVar(outFile, "out", "", "Output file path (-o, --out)")
	
	debug := flag.Bool("D", false, "Enable debug logging (-D, --debug)")
	flag.BoolVar(debug, "debug", false, "Enable debug logging (-D, --debug)")
	
	gitignore := flag.Bool("g", true, "Use .gitignore patterns (-g, --gitignore)")
	flag.BoolVar(gitignore, "gitignore", true, "Use .gitignore patterns (-g, --gitignore)")

	flag.Parse()

	if err := processor.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose, *format, *outFile, *debug, *gitignore); err != nil {
		log.Fatal(err)
	}
}
