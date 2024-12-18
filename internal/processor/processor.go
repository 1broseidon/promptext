package processor

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1broseidon/promptext/internal/config"
	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/info"
	"github.com/1broseidon/promptext/internal/log"
	"github.com/1broseidon/promptext/internal/token"
	"github.com/atotto/clipboard"
)

type Config struct {
	DirPath    string
	Extensions []string
	Excludes   []string
	GitIgnore  bool
	Filter     *filter.Filter
}

func ParseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}

// ProcessResult contains both display and clipboard content
type ProcessResult struct {
	ProjectOutput    *format.ProjectOutput
	DisplayContent   string
	ClipboardContent string
	TokenCount       int
}

// initializeProjectOutput sets up the initial project output structure
func initializeProjectOutput(dirPath string, f *filter.Filter) (*format.ProjectOutput, error) {
	projectOutput := &format.ProjectOutput{}

	// Get project analysis using shared filter
	analysis := info.AnalyzeProject(dirPath, f)
	projectOutput.Analysis = &format.ProjectAnalysis{
		EntryPoints:   analysis.EntryPoints,
		ConfigFiles:   analysis.ConfigFiles,
		CoreFiles:     analysis.CoreFiles,
		TestFiles:     analysis.TestFiles,
		Documentation: analysis.Documentation,
	}

	return projectOutput, nil
}

// populateProjectInfo adds project information to the output
func populateProjectInfo(projectOutput *format.ProjectOutput, projectInfo *info.ProjectInfo) {
	projectOutput.DirectoryTree = projectInfo.DirectoryTree

	if projectInfo.GitInfo != nil {
		projectOutput.GitInfo = &format.GitInfo{
			Branch:        projectInfo.GitInfo.Branch,
			CommitHash:    projectInfo.GitInfo.CommitHash,
			CommitMessage: projectInfo.GitInfo.CommitMessage,
		}
	}

	if projectInfo.Metadata != nil {
		projectOutput.Metadata = &format.Metadata{
			Language:     projectInfo.Metadata.Language,
			Version:      projectInfo.Metadata.Version,
			Dependencies: projectInfo.Metadata.Dependencies,
		}
	}
}

// buildVerboseDisplay creates the verbose display string
func buildVerboseDisplay(projectOutput *format.ProjectOutput) string {
	var displayBuilder strings.Builder

	displayBuilder.WriteString(projectOutput.DirectoryTree.ToMarkdown(0))

	if projectOutput.GitInfo != nil {
		displayBuilder.WriteString(fmt.Sprintf("\nBranch: %s\nCommit: %s\nMessage: %s\n",
			projectOutput.GitInfo.Branch,
			projectOutput.GitInfo.CommitHash,
			projectOutput.GitInfo.CommitMessage))
	}

	if projectOutput.Metadata != nil {
		displayBuilder.WriteString(fmt.Sprintf("\nLanguage: %s\nVersion: %s\n",
			projectOutput.Metadata.Language,
			projectOutput.Metadata.Version))
		if len(projectOutput.Metadata.Dependencies) > 0 {
			displayBuilder.WriteString("Dependencies:\n")
			for _, dep := range projectOutput.Metadata.Dependencies {
				displayBuilder.WriteString(fmt.Sprintf("  - %s\n", dep))
			}
		}
	}

	return displayBuilder.String()
}

// processFile handles the processing of a single file
func processFile(path string, config Config) (*format.FileInfo, error) {
	if !config.Filter.ShouldProcess(path) {
		return nil, nil
	}

	// Skip .DS_Store files immediately
	if filepath.Base(path) == ".DS_Store" {
		return nil, nil
	}

	// Get file info first to check if it's a directory or has read permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Debug("Warning: Cannot stat file %s: %v", path, err)
		return nil, nil
	}

	// Skip directories
	if fileInfo.IsDir() {
		return nil, nil
	}

	// Check read permissions
	if fileInfo.Mode().Perm()&0444 == 0 {
		log.Debug("Warning: No read permission for file %s", path)
		return nil, nil
	}

	content, err := os.ReadFile(path)
	if err != nil {
		log.Debug("Warning: Cannot read file %s: %v", path, err)
		return nil, nil
	}

	// Check if file appears to be binary
	if len(content) > 0 {
		// Check first 1024 bytes for null bytes
		if bytes.IndexByte(content[:min(1024, len(content))], 0) != -1 {
			return nil, nil
		}

		// Check file extension for common binary types
		ext := strings.ToLower(filepath.Ext(path))
		binaryExts := map[string]bool{
			".exe": true, ".dll": true, ".so": true, ".dylib": true,
			".bin": true, ".obj": true, ".o": true,
			".zip": true, ".tar": true, ".gz": true, ".7z": true,
			".pdf": true, ".doc": true, ".docx": true,
			".xls": true, ".xlsx": true, ".ppt": true,
			".db": true, ".sqlite": true, ".sqlite3": true,
		}
		if binaryExts[ext] {
			return nil, nil
			return nil, nil
		}
	}

	rel, err := filepath.Rel(config.DirPath, path)
	if err != nil {
		return nil, fmt.Errorf("error getting relative path for %s: %w", path, err)
	}

	return &format.FileInfo{
		Path:    rel,
		Content: string(content),
	}, nil
}

func ProcessDirectory(config Config, verbose bool) (*ProcessResult, error) {
	log.StartTimer("Project Processing")
	defer log.EndTimer("Project Processing")

	// Initialize project output and get project info using shared filter
	log.StartTimer("Project Analysis")
	projectOutput, err := initializeProjectOutput(config.DirPath, config.Filter)
	if err != nil {
		return nil, err
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}
	log.EndTimer("Project Analysis")

	// Populate project information
	populateProjectInfo(projectOutput, projectInfo)

	// Token analysis
	log.StartTimer("Token Analysis")
	tokenCounter := token.NewTokenCounter()
	log.Debug("=== Token Analysis ===")
	var totalTokens int

	log.Debug("Processing project files:")
	err = filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path for filtering
		relPath, err := filepath.Rel(config.DirPath, path)
		if err != nil {
			return err
		}

		// For directories
		if d.IsDir() {
			if config.Filter.IsExcluded(relPath) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip excluded files silently
		if config.Filter.IsExcluded(relPath) {
			return nil
		}

		var displayContent string
		// Only log files that will be processed
		if config.Filter.ShouldProcess(relPath) {
			log.Debug("  Processing: %s", relPath)
		}

		fileInfo, err := processFile(path, config)
		if err != nil {
			log.Debug("Error processing file %s: %v", path, err)
			return nil // Continue processing other files
		}

		if fileInfo != nil {
			projectOutput.Files = append(projectOutput.Files, *fileInfo)

			// Count tokens for this file
			fileTokens := tokenCounter.EstimateTokens(fileInfo.Content)
			totalTokens += fileTokens

			if verbose {
				displayContent += fmt.Sprintf("\n### File: %s\n```\n%s\n```\n",
					path, fileInfo.Content)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error processing files: %w", err)
	}

	// Get formatter for output
	formatter, err := format.GetFormatter("markdown") // Default to markdown for token counting
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	// Count tokens for directory tree
	treeOutput, _ := formatter.Format(&format.ProjectOutput{DirectoryTree: projectOutput.DirectoryTree})
	treeTokens := tokenCounter.EstimateTokens(treeOutput)
	totalTokens += treeTokens
	log.Debug("Directory structure: %d tokens", treeTokens)

	// Count tokens for git info
	gitTokens := 0
	if projectOutput.GitInfo != nil {
		gitOutput, _ := formatter.Format(&format.ProjectOutput{GitInfo: projectOutput.GitInfo})
		gitTokens = tokenCounter.EstimateTokens(gitOutput)
		totalTokens += gitTokens
	}
	log.Debug("Git information: %d tokens", gitTokens)

	// Count tokens for metadata
	metaTokens := 0
	if projectOutput.Metadata != nil {
		metaOutput, _ := formatter.Format(&format.ProjectOutput{Metadata: projectOutput.Metadata})
		metaTokens = tokenCounter.EstimateTokens(metaOutput)
		totalTokens += metaTokens
	}
	log.Debug("Project metadata: %d tokens", metaTokens)

	// Calculate source file tokens
	sourceTokens := totalTokens - treeTokens - gitTokens - metaTokens
	log.Debug("Source files: %d tokens", sourceTokens)
	log.Debug("Total token count: %d", totalTokens)

	// Add timing summary
	log.Debug("=== Performance ===")
	log.Debug("Total processing time: %.2fms", float64(time.Since(log.GetPhaseStart()).Microseconds())/1000.0)
	log.EndTimer("Token Counting")

	// Format the full output
	formattedOutput, err := formatter.Format(projectOutput)
	if err != nil {
		return nil, fmt.Errorf("error formatting output: %w", err)
	}

	displayContent = ""
	if verbose {
		displayContent = formattedOutput
	}

	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error walking directory: %w", err)
	}

	return &ProcessResult{
		ProjectOutput:    projectOutput,
		DisplayContent:   displayContent,
		ClipboardContent: formattedOutput,
		TokenCount:       tokenCounter.EstimateTokens(formattedOutput),
	}, nil
}

// buildProjectSummary creates the project name summary
func buildProjectSummary(projectInfo *info.ProjectInfo, config Config) string {
	var summary strings.Builder
	summary.WriteString("📦 Project Summary:\n")

	if projectInfo.Metadata != nil && projectInfo.Metadata.Name != "" {
		summary.WriteString(fmt.Sprintf("   Project: %s\n", projectInfo.Metadata.Name))
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			summary.WriteString(fmt.Sprintf("   Project: %s\n", filepath.Base(absPath)))
		}
	}
	return summary.String()
}

// buildLanguageInfo creates the language and dependencies summary
func buildLanguageInfo(metadata *info.ProjectMetadata) string {
	if metadata == nil {
		return ""
	}

	var info strings.Builder
	info.WriteString(fmt.Sprintf("   Language: %s", metadata.Language))
	if metadata.Version != "" {
		info.WriteString(fmt.Sprintf(" %s", metadata.Version))
	}
	info.WriteString("\n")

	if len(metadata.Dependencies) > 0 {
		mainDeps, devDeps := countDependencies(metadata.Dependencies)
		if devDeps > 0 {
			info.WriteString(fmt.Sprintf("   Dependencies: %d packages (%d main, %d dev)\n",
				len(metadata.Dependencies), mainDeps, devDeps))
		} else {
			info.WriteString(fmt.Sprintf("   Dependencies: %d packages\n", mainDeps))
		}
	}

	return info.String()
}

// countDependencies counts main and dev dependencies
func countDependencies(deps []string) (main, dev int) {
	for _, dep := range deps {
		if strings.HasPrefix(dep, "[dev] ") {
			dev++
		} else {
			main++
		}
	}
	return main, dev
}

// countIncludedFiles counts files that match the filter criteria
func countIncludedFiles(config Config) (int, error) {
	fileCount := 0
	err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(config.DirPath, path)
		if config.Filter.ShouldProcess(rel) {
			fileCount++
		}
		return nil
	})
	return fileCount, err
}

// GetMetadataSummary returns a concise summary of project metadata
func GetMetadataSummary(config Config, tokenCount int) (string, error) {
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return "", err
	}

	var content strings.Builder

	// Project name
	if projectInfo.Metadata != nil && projectInfo.Metadata.Name != "" {
		content.WriteString("📦 " + projectInfo.Metadata.Name)
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			content.WriteString("📦 " + filepath.Base(absPath))
		}
	}

	// Language if detected
	if projectInfo.Metadata != nil && projectInfo.Metadata.Language != "" {
		content.WriteString(fmt.Sprintf(" (%s)", projectInfo.Metadata.Language))
	}

	content.WriteString("\n")

	// File count
	fileCount, err := countIncludedFiles(config)
	if err != nil {
		return "", fmt.Errorf("error counting files: %w", err)
	}
	content.WriteString(fmt.Sprintf("   Files: %d", fileCount))

	// Token count
	if tokenCount > 0 {
		content.WriteString(fmt.Sprintf(" • Tokens: ~%d", tokenCount))
	}
	content.WriteString("\n")

	// Create bordered output
	contentLines := strings.Split(strings.TrimRight(content.String(), "\n"), "\n")
	maxWidth := 0
	for _, line := range contentLines {
		if len(line) > maxWidth {
			maxWidth = len(line)
		}
	}

	var summary strings.Builder
	summary.WriteString("\033[32m") // Start green color

	// Top border
	summary.WriteString("╭" + strings.Repeat("─", maxWidth+4) + "╮\n")

	// Content lines with borders
	for _, line := range contentLines {
		paddedLine := line + strings.Repeat(" ", maxWidth-len(line))
		summary.WriteString("│ " + paddedLine + "     │\n")
	}

	// Bottom border
	summary.WriteString("╰" + strings.Repeat("─", maxWidth+4) + "╯")
	summary.WriteString("\033[0m") // Reset color

	return summary.String(), nil
}

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool) error {
	// Enable debug logging if flag is set
	if debug {
		log.Enable()
		log.SetColorEnabled(true)
	}

	log.Debug("=== Promptext Initialization ===")
	log.Debug("Directory: %s", dirPath)
	// Handle "md" as alias for "markdown"
	if outputFormat == "md" {
		outputFormat = "markdown"
	}

	// Validate format
	formatter, err := format.GetFormatter(outputFormat)
	if err != nil {
		return fmt.Errorf("invalid format (must be markdown or xml): %w", err)
	}
	// Convert dirPath to absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Load config file from the specified directory
	fileConfig, err := config.LoadConfig(absPath)
	if err != nil {
		log.Info("Warning: Failed to load .promptext.yml from %s: %v", absPath, err)
		fileConfig = &config.FileConfig{}
	}

	// Merge file config with command line flags
	extensions, excludes, verboseFlag, _, useGitIgnore := fileConfig.MergeWithFlags(extension, exclude, verbose, debug, &gitignore)
	log.Debug("Configuration:")
	log.Debug("  • Extensions: %v", extensions)
	log.Debug("  • Excludes: %#v", excludes)
	log.Debug("  • Git Ignore: %v", useGitIgnore)

	// Create filter options
	filterOpts := filter.Options{
		Includes:      extensions,
		Excludes:      excludes,
		IgnoreDefault: true,
		UseGitIgnore:  useGitIgnore,
	}

	// Create the filter once and reuse it
	f := filter.New(filterOpts)

	// Create processor configuration with filter
	procConfig := Config{
		DirPath:    absPath,
		Extensions: extensions,
		Excludes:   excludes,
		GitIgnore:  useGitIgnore,
		Filter:     f,
	}

	if infoOnly {
		// Process directory just to get token count
		result, err := ProcessDirectory(procConfig, false)
		if err != nil {
			return fmt.Errorf("error processing directory: %v", err)
		}

		// Display project summary with token count
		if info, err := GetMetadataSummary(procConfig, result.TokenCount); err == nil {
			fmt.Printf("\033[32m%s\033[0m\n", info)
		} else {
			return fmt.Errorf("error getting project info: %v", err)
		}
		return nil
	}

	// Process the directory
	result, err := ProcessDirectory(procConfig, verboseFlag)
	if err != nil {
		return fmt.Errorf("error processing directory: %v", err)
	}

	// Write display content to stdout
	if verbose {
		fmt.Println(result.DisplayContent)
	}

	// Format output
	formattedOutput, err := formatter.Format(result.ProjectOutput)
	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	if outFile != "" {
		// Write to file if -out is specified
		if err := os.WriteFile(outFile, []byte(formattedOutput), 0644); err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
		// Always print metadata summary and success message in green
		if info, err := GetMetadataSummary(procConfig, result.TokenCount); err == nil {
			fmt.Printf("\033[32m%s\n✓ code context written to %s (%s format)\033[0m\n",
				info, outFile, outputFormat)
		}
	} else if !noCopy {
		// Copy to clipboard if no output file is specified and clipboard is not disabled
		if err := clipboard.WriteAll(formattedOutput); err != nil {
			log.Info("Warning: Failed to copy to clipboard: %v", err)
		} else {
			// Always print metadata summary and success message in green
			if info, err := GetMetadataSummary(procConfig, result.TokenCount); err == nil {
				fmt.Printf("%s\n✓ code context copied to clipboard (%s format)\033[0m\n",
					info, outputFormat)
			}
		}
	}

	return nil
}
