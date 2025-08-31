package processor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1broseidon/promptext/internal/config"
	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/filter/rules"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/info"
	"github.com/1broseidon/promptext/internal/log"
	"github.com/1broseidon/promptext/internal/token"
	"github.com/atotto/clipboard"
	"github.com/jedib0t/go-pretty/v6/text"
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
	ProjectInfo      *info.ProjectInfo
}

// DryRunResult contains dry-run preview information
type DryRunResult struct {
	FilePaths       []string
	EstimatedTokens int
	ConfigSummary   *ConfigSummary
	ProjectInfo     *info.ProjectInfo
}

// ConfigSummary contains effective configuration information
type ConfigSummary struct {
	Extensions      []string
	Excludes        []string
	UseGitIgnore    bool
	UseDefaultRules bool
	Format          string
	OutputFile      string
}

// validateFilePath validates and gets the relative path for a file
func validateFilePath(path string, config Config) (string, error) {
	rel, err := filepath.Rel(config.DirPath, path)
	if err != nil {
		return "", fmt.Errorf("error getting relative path for %s: %w", path, err)
	}

	if !config.Filter.ShouldProcess(rel) {
		return "", nil
	}

	// Skip .DS_Store files immediately
	if filepath.Base(path) == ".DS_Store" {
		return "", nil
	}

	return rel, nil
}

// checkFilePermissions validates file type and permissions
func checkFilePermissions(path string) error {
	// Get file info first to check if it's a directory or has read permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Skip directories
	if fileInfo.IsDir() {
		return fmt.Errorf("is directory")
	}

	// Check read permissions
	if fileInfo.Mode().Perm()&0444 == 0 {
		return fmt.Errorf("no read permissions")
	}

	// Check if file is binary using BinaryRule
	binaryRule := rules.NewBinaryRule()
	if binaryRule.Match(path) {
		return fmt.Errorf("binary file")
	}

	return nil
}

// readFileContent reads and returns file content as string
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// processFile handles the processing of a single file
func processFile(path string, config Config) (*format.FileInfo, error) {
	rel, err := validateFilePath(path, config)
	if err != nil {
		return nil, err
	}
	if rel == "" {
		return nil, nil // File should be skipped
	}

	if err := checkFilePermissions(path); err != nil {
		return nil, nil // File should be skipped
	}

	content, err := readFileContent(path)
	if err != nil {
		return nil, nil // File should be skipped
	}

	return &format.FileInfo{
		Path:    rel,
		Content: content,
	}, nil
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

// PreviewDirectory performs a dry-run preview of what files would be processed
func PreviewDirectory(config Config) (*DryRunResult, error) {
	log.StartTimer("Dry Run Preview")
	defer log.EndTimer("Dry Run Preview")

	// Initialize result
	result := &DryRunResult{
		FilePaths: []string{},
		ConfigSummary: &ConfigSummary{
			Extensions:      config.Extensions,
			Excludes:        config.Excludes,
			UseGitIgnore:    config.GitIgnore,
			UseDefaultRules: config.Filter != nil,
		},
	}

	// Collect files that would be processed
	tokenCounter := token.NewTokenCounter()
	var estimatedTokens int

	log.Debug("=== Dry Run: Analyzing Files ===")

	err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
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

		// Check if file would pass validation and filtering
		rel, err := validateFilePath(path, config)
		if err != nil {
			return nil // Skip files that would fail validation
		}
		if rel == "" {
			return nil // File should be skipped due to filtering
		}

		// Check permissions and file type without reading content
		if err := checkFilePermissions(path); err != nil {
			return nil // Skip files that would fail permission check
		}

		// Add to result
		result.FilePaths = append(result.FilePaths, rel)

		// Estimate tokens based on file size (rough approximation: 4 chars per token)
		if fileInfo, err := os.Stat(path); err == nil {
			estimatedFileTokens := int(fileInfo.Size() / 4)
			estimatedTokens += estimatedFileTokens
			log.Debug("Would process: %s (estimated %d tokens)", rel, estimatedFileTokens)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error during dry-run preview: %w", err)
	}

	result.EstimatedTokens = estimatedTokens

	// Get project info for dry-run
	if projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter); err == nil {
		result.ProjectInfo = projectInfo
		
		// Add estimated tokens for metadata (rough approximation)
		if projectInfo.DirectoryTree != nil {
			directoryTreeString := projectInfo.DirectoryTree.ToMarkdown(0)
			result.EstimatedTokens += tokenCounter.EstimateTokens(directoryTreeString)
		}
	}

	log.Debug("Dry run complete: %d files, ~%d tokens", len(result.FilePaths), result.EstimatedTokens)
	return result, nil
}

func ProcessDirectory(config Config, verbose bool) (*ProcessResult, error) {
	log.StartTimer("Project Processing")
	defer log.EndTimer("Project Processing")

	// Initialize project output
	projectOutput := &format.ProjectOutput{}

	// Combined file processing and token analysis
	log.StartTimer("Processing Files")
	tokenCounter := token.NewTokenCounter()
	log.Debug("=== Processing Files & Counting Tokens ===")
	var totalTokens int

	// Process all files first
	var processedFiles []format.FileInfo
	err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
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

		// Process file
		fileInfo, err := processFile(path, config)
		if err != nil {
			log.Debug("Error processing file %s: %v", path, err)
			return nil // Continue processing other files
		}

		if fileInfo != nil {
			processedFiles = append(processedFiles, *fileInfo)

			// Count tokens and log immediately
			fileTokens := tokenCounter.EstimateTokens(fileInfo.Content)
			totalTokens += fileTokens
			log.Debug("Processing: %s (%d tokens)", relPath, fileTokens)

			if verbose && !log.IsDebugEnabled() {
				fmt.Printf("\n### File: %s\n```\n%s\n```\n",
					path, fileInfo.Content)
			}
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error processing files: %w", err)
	}
	log.EndTimer("Processing Files")

	// Store processed files
	projectOutput.Files = processedFiles

	// Get project info using processed files
	log.StartTimer("Project Analysis")
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}
	log.EndTimer("Project Analysis")

	// Populate project information
	populateProjectInfo(projectOutput, projectInfo)

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

	// Format the full output
	formattedOutput, err := formatter.Format(projectOutput)
	if err != nil {
		return nil, fmt.Errorf("error formatting output: %w", err)
	}

	var displayContent string
	if verbose {
		displayContent = formattedOutput
	}

	return &ProcessResult{
		ProjectOutput:    projectOutput,
		DisplayContent:   displayContent,
		ClipboardContent: formattedOutput,
		TokenCount:       totalTokens,
		ProjectInfo:      projectInfo,
	}, nil
}

// buildProjectHeader constructs the project name and basic info
func buildProjectHeader(config Config, result *ProcessResult, infoOnly bool) string {
	var content strings.Builder

	// Project name and language (always shown)
	if result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Name != "" {
		content.WriteString("📦 " + result.ProjectInfo.Metadata.Name)
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			content.WriteString("📦 " + filepath.Base(absPath))
		}
	}

	if result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Language != "" {
		content.WriteString(fmt.Sprintf(" (%s", result.ProjectInfo.Metadata.Language))
		if result.ProjectInfo.Metadata.Version != "" && infoOnly {
			content.WriteString(fmt.Sprintf(" %s", result.ProjectInfo.Metadata.Version))
		}
		content.WriteString(")")
	}

	// Basic file and token count (always shown)
	fileCount := len(result.ProjectOutput.Files)
	content.WriteString(fmt.Sprintf("\n   Files: %d", fileCount))
	if result.TokenCount > 0 {
		content.WriteString(fmt.Sprintf(" • Tokens: ~%d", result.TokenCount))
	}

	return content.String()
}

// analyzeFileStatistics collects file type and size statistics
func analyzeFileStatistics(files []format.FileInfo, config Config) (map[string]int, int64, []string) {
	fileTypes := make(map[string]int)
	var totalSize int64
	var entryPoints []string

	for _, file := range files {
		typeInfo := filter.GetFileType(file.Path, config.Filter)
		fileTypes[typeInfo.Type]++
		totalSize += typeInfo.Size

		// Track entry points
		if typeInfo.IsEntryPoint {
			entryPoints = append(entryPoints, file.Path)
		}
	}

	return fileTypes, totalSize, entryPoints
}

// buildFileAnalysis creates the file analysis section
func buildFileAnalysis(fileTypes map[string]int, totalSize int64, entryPoints []string) string {
	var content strings.Builder

	// Display File Distribution
	content.WriteString("\n   Types: ")
	first := true
	for typ, count := range fileTypes {
		if !first {
			content.WriteString(" • ")
		}
		content.WriteString(fmt.Sprintf("%s: %d", typ, count))
		first = false
	}
	content.WriteString("\n")

	// Display Size Information
	if totalSize > 0 {
		content.WriteString(fmt.Sprintf("   Total Size: %s\n", formatSize(totalSize)))
	}

	// Display Entry Points
	if len(entryPoints) > 0 {
		content.WriteString("\n🚪 Entry Points\n")
		for _, entry := range entryPoints {
			content.WriteString(fmt.Sprintf("   • %s\n", entry))
		}
	}

	return content.String()
}

// buildDependenciesSection creates the dependencies section
func buildDependenciesSection(result *ProcessResult) string {
	if result.ProjectOutput.Metadata == nil || len(result.ProjectOutput.Metadata.Dependencies) == 0 {
		return ""
	}

	var content strings.Builder
	content.WriteString("\n📚 Dependencies\n")
	for _, dep := range result.ProjectOutput.Metadata.Dependencies {
		content.WriteString(fmt.Sprintf("   • %s\n", dep))
	}
	return content.String()
}

// buildHealthSection creates the project health section
func buildHealthSection(result *ProcessResult) string {
	if result.ProjectInfo.Metadata == nil || result.ProjectInfo.Metadata.Health == nil {
		return ""
	}

	health := result.ProjectInfo.Metadata.Health
	var content strings.Builder
	content.WriteString("\n🏥 Project Health\n")

	// Documentation
	content.WriteString(fmt.Sprintf("   • README: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasReadme]))
	content.WriteString(fmt.Sprintf("   • LICENSE: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasLicense]))

	// Testing
	content.WriteString(fmt.Sprintf("   • Tests: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasTests]))

	// CI/CD
	if health.HasCI {
		content.WriteString(fmt.Sprintf("   • CI/CD: ✓ (%s)\n", health.CISystem))
	} else {
		content.WriteString("   • CI/CD: ✗\n")
	}

	return content.String()
}

// buildGitSection creates the git information section
func buildGitSection(result *ProcessResult) string {
	if result.ProjectOutput.GitInfo == nil {
		return ""
	}

	var content strings.Builder
	content.WriteString("\n🔄 Git Status\n")
	content.WriteString(fmt.Sprintf("   Branch: %s\n", result.ProjectOutput.GitInfo.Branch))

	if result.ProjectOutput.GitInfo.CommitHash != "" {
		shortHash := result.ProjectOutput.GitInfo.CommitHash
		if len(shortHash) > 7 {
			shortHash = shortHash[:7]
		}
		content.WriteString(fmt.Sprintf("   Latest: %s", shortHash))
		if result.ProjectOutput.GitInfo.CommitMessage != "" {
			content.WriteString(fmt.Sprintf(" - %s", result.ProjectOutput.GitInfo.CommitMessage))
		}
		content.WriteString("\n")
	}

	return content.String()
}

// formatBoxedOutput creates a boxed output with borders
func formatBoxedOutput(content string) string {
	contentLines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	maxWidth := 0
	for _, line := range contentLines {
		width := text.RuneCount(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Add padding to max width
	maxWidth += 4 // 2 spaces on each side

	var summary strings.Builder
	summary.WriteString("\033[32m") // Start green color

	// Top border
	summary.WriteString("╭" + strings.Repeat("─", maxWidth) + "╮\n")

	// Content lines
	for _, line := range contentLines {
		// Calculate padding needed
		lineWidth := text.RuneCount(line)
		padding := maxWidth - lineWidth

		// Write line with padding
		summary.WriteString("│ " + line + strings.Repeat(" ", padding-2) + " │\n")
	}

	// Bottom border
	summary.WriteString("╰" + strings.Repeat("─", maxWidth) + "╯")
	summary.WriteString("\033[0m") // Reset color

	return summary.String()
}

// GetMetadataSummary returns a summary of project metadata and analysis
// If infoOnly is true, returns a rich summary with all details
// Otherwise returns a minimal summary with basic project info
func GetMetadataSummary(config Config, result *ProcessResult, infoOnly bool) (string, error) {
	var content strings.Builder

	// Build basic project header
	content.WriteString(buildProjectHeader(config, result, infoOnly))

	// Only show detailed analysis if infoOnly is true
	if infoOnly {
		content.WriteString("\n")

		// Analyze file statistics
		fileTypes, totalSize, entryPoints := analyzeFileStatistics(result.ProjectOutput.Files, config)

		// Build file analysis section
		content.WriteString(buildFileAnalysis(fileTypes, totalSize, entryPoints))

		// Build dependencies section
		content.WriteString(buildDependenciesSection(result))

		// Build health section
		content.WriteString(buildHealthSection(result))

		// Build git section
		content.WriteString(buildGitSection(result))
	}

	return formatBoxedOutput(content.String()), nil
}

// formatSize converts bytes to human readable string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDryRunOutput formats the dry-run preview results for display
func FormatDryRunOutput(result *DryRunResult, config Config) string {
	var content strings.Builder

	// Project header
	if result.ProjectInfo != nil && result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Name != "" {
		content.WriteString("📦 " + result.ProjectInfo.Metadata.Name)
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			content.WriteString("📦 " + filepath.Base(absPath))
		}
	}

	if result.ProjectInfo != nil && result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Language != "" {
		content.WriteString(fmt.Sprintf(" (%s)", result.ProjectInfo.Metadata.Language))
	}

	// Dry-run summary
	content.WriteString(fmt.Sprintf("\n🔍 DRY RUN PREVIEW\n"))
	content.WriteString(fmt.Sprintf("   Would process: %d files\n", len(result.FilePaths)))
	content.WriteString(fmt.Sprintf("   Estimated tokens: ~%d\n", result.EstimatedTokens))

	// Configuration summary
	content.WriteString("\n⚙️ Effective Configuration\n")
	if len(result.ConfigSummary.Extensions) > 0 {
		content.WriteString(fmt.Sprintf("   Extensions: %s\n", strings.Join(result.ConfigSummary.Extensions, ", ")))
	} else {
		content.WriteString("   Extensions: all supported types\n")
	}

	if len(result.ConfigSummary.Excludes) > 0 {
		content.WriteString(fmt.Sprintf("   Excludes: %s\n", strings.Join(result.ConfigSummary.Excludes, ", ")))
	}

	content.WriteString(fmt.Sprintf("   Git ignore: %v\n", result.ConfigSummary.UseGitIgnore))
	content.WriteString(fmt.Sprintf("   Default rules: %v\n", result.ConfigSummary.UseDefaultRules))

	// Files to be processed (showing first 10 with more indicator)
	content.WriteString("\n📂 Files to Process\n")
	maxDisplay := 10
	for i, filePath := range result.FilePaths {
		if i >= maxDisplay {
			remaining := len(result.FilePaths) - maxDisplay
			content.WriteString(fmt.Sprintf("   ... and %d more files\n", remaining))
			break
		}
		content.WriteString(fmt.Sprintf("   • %s\n", filePath))
	}

	if len(result.FilePaths) == 0 {
		content.WriteString("   ⚠️  No files would be processed\n")
	}

	return formatBoxedOutput(content.String())
}

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool, useDefaultRules bool, dryRun bool, quiet bool) error {
	// Enable debug logging if flag is set
	if debug {
		log.Enable()
		log.SetColorEnabled(true)
	}
	
	// Set quiet mode
	if quiet {
		log.SetQuiet(true)
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

	// Load global configuration
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		log.Info("Warning: Failed to load global config: %v", err)
		globalConfig = &config.FileConfig{}
	}

	// Load project config file from the specified directory
	projectConfig, err := config.LoadConfig(absPath)
	if err != nil {
		log.Info("Warning: Failed to load .promptext.yml from %s: %v", absPath, err)
		projectConfig = &config.FileConfig{}
	}

	// Merge global, project, and flag configurations with proper precedence
	extensions, excludes, verboseFlag, _, useGitIgnore, useDefaultRules := config.MergeConfigs(globalConfig, projectConfig, extension, exclude, verbose, debug, &gitignore, &useDefaultRules)
	log.Debug("Configuration:")
	log.Debug("  • Extensions: %v", extensions)
	log.Debug("  • Excludes: %#v", excludes)
	log.Debug("  • Git Ignore: %v", useGitIgnore)

	// Create filter options
	filterOpts := filter.Options{
		Includes:        extensions,
		Excludes:        excludes,
		UseDefaultRules: useDefaultRules,
		UseGitIgnore:    useGitIgnore,
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

	// Handle dry-run mode
	if dryRun {
		dryRunResult, err := PreviewDirectory(procConfig)
		if err != nil {
			return fmt.Errorf("error during dry-run preview: %v", err)
		}

		// Update config summary with additional info
		dryRunResult.ConfigSummary.Format = outputFormat
		dryRunResult.ConfigSummary.OutputFile = outFile

		// Display dry-run results
		preview := FormatDryRunOutput(dryRunResult, procConfig)
		if quiet {
			// In quiet mode, output minimal dry-run info
			fmt.Printf("files=%d tokens=%d\n", len(dryRunResult.FilePaths), dryRunResult.EstimatedTokens)
		} else {
			fmt.Printf("\033[32m%s\033[0m\n", preview)
		}
		
		return nil
	}

	// Process directory once and reuse results
	result, err := ProcessDirectory(procConfig, verboseFlag)
	if err != nil {
		return fmt.Errorf("error processing directory: %v", err)
	}

	// Get metadata summary using the already processed result
	info, err := GetMetadataSummary(procConfig, result, infoOnly)
	if err != nil {
		return fmt.Errorf("error getting project info: %v", err)
	}

	// If info-only flag is set, just display the summary and return
	if infoOnly {
		if quiet {
			// In quiet mode, output minimal info
			fileCount := len(result.ProjectOutput.Files)
			fmt.Printf("files=%d tokens=%d\n", fileCount, result.TokenCount)
		} else {
			fmt.Printf("\033[32m%s\033[0m\n", info)
		}
		return nil
	}

	// Format output for the selected format
	formattedOutput, err := formatter.Format(result.ProjectOutput)
	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	// Handle output based on flags
	if outFile != "" {
		if err := os.WriteFile(outFile, []byte(formattedOutput), 0644); err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
		if quiet {
			// In quiet mode, output minimal success info
			fmt.Printf("written=%s format=%s files=%d tokens=%d\n", outFile, outputFormat, len(result.ProjectOutput.Files), result.TokenCount)
		} else {
			fmt.Printf("\033[32m%s\n✓ code context written to %s (%s format)\033[0m\n",
				info, outFile, outputFormat)
		}
	} else if !noCopy {
		if err := clipboard.WriteAll(formattedOutput); err != nil {
			if !quiet {
				log.Info("Warning: Failed to copy to clipboard: %v", err)
			}
			// In quiet mode, exit with error code for clipboard failure
			if quiet {
				return fmt.Errorf("clipboard copy failed")
			}
		} else {
			if quiet {
				// In quiet mode, output minimal success info
				fmt.Printf("clipboard=ok format=%s files=%d tokens=%d\n", outputFormat, len(result.ProjectOutput.Files), result.TokenCount)
			} else {
				fmt.Printf("\033[32m%s\n✓ code context copied to clipboard (%s format)\033[0m\n",
					info, outputFormat)
			}
		}
	}

	return nil
}
