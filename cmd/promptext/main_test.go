package main

import (
	"os"
	"testing"

	"github.com/spf13/pflag"
)

func TestMainFlags(t *testing.T) {
	// Save original os.Args and pflag.CommandLine
	oldArgs := os.Args
	oldFlagCommandLine := pflag.CommandLine
	defer func() {
		// Restore original values after test
		os.Args = oldArgs
		pflag.CommandLine = oldFlagCommandLine
	}()

	tests := []struct {
		name     string
		args     []string
		expected struct {
			dir      string
			ext      string
			exclude  string
			infoOnly bool
			verbose  bool
			format   string
			outFile  string
		}
	}{
		{
			name: "default values",
			args: []string{"promptext"},
			expected: struct {
				dir      string
				ext      string
				exclude  string
				infoOnly bool
				verbose  bool
				format   string
				outFile  string
			}{
				dir:      ".",
				ext:      "",
				exclude:  "",
				infoOnly: false,
				verbose:  false,
				format:   "markdown",
				outFile:  "",
			},
		},
		{
			name: "all flags set",
			args: []string{
				"promptext",
				"--directory", "/test/path",
				"--extension", ".go,.js",
				"--exclude", "vendor,node_modules",
				"--no-copy",
				"--info",
				"--verbose",
				"--format", "xml",
				"--output", "output.xml",
			},
			expected: struct {
				dir      string
				ext      string
				exclude  string
				noCopy   bool
				infoOnly bool
				verbose  bool
				format   string
				outFile  string
			}{
				dir:      "/test/path",
				ext:      ".go,.js",
				exclude:  "vendor,node_modules",
				infoOnly: true,
				verbose:  true,
				format:   "xml",
				outFile:  "output.xml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset pflag.CommandLine for each test
			pflag.CommandLine = pflag.NewFlagSet(tt.args[0], pflag.ExitOnError)
			os.Args = tt.args

			// Define flags again (since we're using a new FlagSet)
			dirPath := pflag.StringP("directory", "d", ".", "Directory path to process")
			extension := pflag.StringP("extension", "e", "", "File extension to filter, e.g., .go,.js")
			exclude := pflag.StringP("exclude", "x", "", "Patterns to exclude, comma-separated")
			infoOnly := pflag.BoolP("info", "i", false, "Only display project summary")
			verbose := pflag.BoolP("verbose", "v", false, "Show full code content in terminal")
			format := pflag.StringP("format", "f", "markdown", "Output format: markdown, xml, json")
			outFile := pflag.StringP("output", "o", "", "Output file path")

			// Parse flags
			pflag.Parse()

			// Check if parsed values match expected values
			if *dirPath != tt.expected.dir {
				t.Errorf("dirPath = %v, want %v", *dirPath, tt.expected.dir)
			}
			if *extension != tt.expected.ext {
				t.Errorf("extension = %v, want %v", *extension, tt.expected.ext)
			}
			if *exclude != tt.expected.exclude {
				t.Errorf("exclude = %v, want %v", *exclude, tt.expected.exclude)
			}
			if *infoOnly != tt.expected.infoOnly {
				t.Errorf("infoOnly = %v, want %v", *infoOnly, tt.expected.infoOnly)
			}
			if *verbose != tt.expected.verbose {
				t.Errorf("verbose = %v, want %v", *verbose, tt.expected.verbose)
			}
			if *format != tt.expected.format {
				t.Errorf("format = %v, want %v", *format, tt.expected.format)
			}
			if *outFile != tt.expected.outFile {
				t.Errorf("outFile = %v, want %v", *outFile, tt.expected.outFile)
			}
		})
	}
}
