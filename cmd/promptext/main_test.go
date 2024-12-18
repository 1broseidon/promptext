package main

import (
	"flag"
	"os"
	"testing"
)

func TestMainFlags(t *testing.T) {
	// Save original os.Args and flag.CommandLine
	oldArgs := os.Args
	oldFlagCommandLine := flag.CommandLine
	defer func() {
		// Restore original values after test
		os.Args = oldArgs
		flag.CommandLine = oldFlagCommandLine
	}()

	tests := []struct {
		name     string
		args     []string
		expected struct {
			dir      string
			ext      string
			exclude  string
			noCopy   bool
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
				noCopy   bool
				infoOnly bool
				verbose  bool
				format   string
				outFile  string
			}{
				dir:      ".",
				ext:      "",
				exclude:  "",
				noCopy:   false,
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
				"-d", "/test/path",
				"-e", ".go,.js",
				"-x", "vendor,node_modules",
				"-n",
				"-i",
				"-v",
				"-f", "xml",
				"-o", "output.xml",
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
				noCopy:   true,
				infoOnly: true,
				verbose:  true,
				format:   "xml",
				outFile:  "output.xml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flag.CommandLine for each test
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)
			os.Args = tt.args

			// Define flags again (since we're using a new FlagSet)
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

			// Parse flags
			flag.Parse()

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
			if *noCopy != tt.expected.noCopy {
				t.Errorf("noCopy = %v, want %v", *noCopy, tt.expected.noCopy)
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
