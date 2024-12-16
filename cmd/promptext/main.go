package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/1broseidon/promptext/internal/processor"
)

func main() {
	// Define command line flags
	dirPath := flag.String("dir", ".", "Directory path to process")
	extension := flag.String("ext", "", "File extension to filter (e.g., .go,.js)")
	exclude := flag.String("exclude", "", "Patterns to exclude (comma-separated)")

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
}
