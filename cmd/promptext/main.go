package main

import (
	"flag"
	"log"

	"promptext"
)

func main() {
	// Define command line flags
	dirPath := flag.String("dir", ".", "Directory path to process")
	extension := flag.String("ext", "", "File extension to filter (e.g., .go,.js)")
	exclude := flag.String("exclude", "", "Patterns to exclude (comma-separated)")
	noCopy := flag.Bool("no-copy", false, "Disable automatic copying to clipboard")
	infoOnly := flag.Bool("info", false, "Only display project summary")
	verbose := flag.Bool("verbose", false, "Show full code content in terminal")

	flag.Parse()

	if err := promptext.Run(*dirPath, *extension, *exclude, *noCopy, *infoOnly, *verbose); err != nil {
		log.Fatal(err)
	}
}
