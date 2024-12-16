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
	infoOnly := flag.Bool("info", false, "Only display project summary")
	verbose := flag.Bool("verbose", false, "Show full code content in terminal")

	flag.Parse()

	// Create processor configuration
	config := processor.Config{
		DirPath:    *dirPath,
		Extensions: processor.ParseCommaSeparated(*extension),
		Excludes:   processor.ParseCommaSeparated(*exclude),
	}

	if *infoOnly {
		// Only display project summary
		if info, err := processor.GetMetadataSummary(config); err == nil {
			fmt.Printf("\033[32m%s\033[0m\n", info)
		} else {
			log.Fatalf("Error getting project info: %v", err)
		}
	} else {
		// Process the directory
		output, err := processor.ProcessDirectory(config, *verbose)
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
				// Print metadata summary and success message in green
				if info, err := processor.GetMetadataSummary(config); err == nil {
					fmt.Printf("\033[32m%s   ✓ code context copied to clipboard\033[0m\n", info)
				}
			}
		}
	}
}
