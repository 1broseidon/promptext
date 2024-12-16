package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/atotto/clipboard"
	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("dir", ".", "Directory path to process")
	extension := flag.String("ext", "", "File extension to filter (e.g., .go,.js)")
	exclude := flag.String("exclude", "", "Patterns to exclude (comma-separated)")
	noCopy := flag.Bool("no-copy", false, "Disable automatic copying to clipboard")

	flag.Parse()

	// Create processor configuration
	config := processor.Config{
		DirPath:    *dirPath,
		Extensions: processor.ParseCommaSeparated(*extension),
		Excludes:   processor.ParseCommaSeparated(*exclude),
	}

	// Process the directory
	output, err := processor.ProcessDirectory(config)
	if err != nil {
		log.Fatalf("Error processing directory: %v", err)
	}

	// Write to stdout
	fmt.Println(output)

	// Copy to clipboard unless disabled
	if !*noCopy {
		if err := clipboard.WriteAll(output); err != nil {
			log.Printf("Warning: Failed to copy to clipboard: %v", err)
		} else {
			// Print success message in green (not included in clipboard)
			fmt.Printf("\033[32mcode context copied to clipboard\033[0m\n")
		}
	}
}
