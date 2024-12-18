package processor

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

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
func initializeProjectOutput(config Config) (*format.ProjectOutput, error) {
	projectOutput := &format.ProjectOutput{}

	// Get project analysis
	analysis := info.AnalyzeProject(config.DirPath)
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

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", path, err)
	}

	// Check if file appears to be binary
	if len(content) > 0 {
		// Check first 1024 bytes for null bytes
		if bytes.IndexByte(content[:min(1024, len(content))], 0) != -1 {
			log.Debug("  Skipping binary file (null bytes): %s", path)
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
			log.Debug("  Skipping binary file (extension): %s", path)
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
	// Initialize project output and display content using shared filter
	projectOutput, err := initializeProjectOutput(config)
	if err != nil {
		return nil, err
	}
	
	var displayContent string

	// Get project information using shared filter
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}

	// Populate project information
	populateProjectInfo(projectOutput, projectInfo)

	// Create token counter
	tokenCounter := token.NewTokenCounter()

	// Initialize token counting
	log.Debug("\nToken counting breakdown:")
	var totalTokens int
	
	log.Debug("\nProcessing project files:")
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
			// Check exclusion before any processing
			if config.Filter.IsExcluded(relPath) {
				log.Debug("  Skipping excluded directory: %s", relPath)
				return filepath.SkipDir
			}
			
			// Only log root level directories that aren't excluded
			if filepath.Dir(path) == config.DirPath {
				log.Debug("  Scanning directory: %s", relPath)
			}
			return nil
		}

		// For files, check exclusion
		if config.Filter.IsExcluded(relPath) {
			return nil
		}

		fileInfo, err := processFile(path, config)
		if err != nil {
			return err
		}

		if fileInfo != nil {
			projectOutput.Files = append(projectOutput.Files, *fileInfo)
			
			// Count tokens for this file
			fileTokens := tokenCounter.EstimateTokens(fileInfo.Content)
			totalTokens += fileTokens
			log.Debug("  File %s: %d tokens", fileInfo.Path, fileTokens)

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
	log.Debug("  Directory tree: %d tokens", treeTokens)
	
	// Count tokens for git info
	if projectOutput.GitInfo != nil {
		gitOutput, _ := formatter.Format(&format.ProjectOutput{GitInfo: projectOutput.GitInfo})
		gitTokens := tokenCounter.EstimateTokens(gitOutput)
		totalTokens += gitTokens
		log.Debug("  Git info: %d tokens", gitTokens)
	}
	
	// Count tokens for metadata
	if projectOutput.Metadata != nil {
		metaOutput, _ := formatter.Format(&format.ProjectOutput{Metadata: projectOutput.Metadata})
		metaTokens := tokenCounter.EstimateTokens(metaOutput)
		totalTokens += metaTokens
		log.Debug("  Metadata: %d tokens", metaTokens)
	}

	log.Debug("Total tokens: %d", totalTokens)

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
	infoConfig := &info.Config{
		Extensions: config.Extensions,
		Excludes:   config.Excludes,
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return "", err
	}

	var summary strings.Builder

	// Build project summary
	summary.WriteString(buildProjectSummary(projectInfo, config))

	// Add language and dependencies info
	if projectInfo.Metadata != nil {
		summary.WriteString(buildLanguageInfo(projectInfo.Metadata))
	}

	// Add git info
	if projectInfo.GitInfo != nil {
		summary.WriteString(fmt.Sprintf("   Branch: %s (%s)\n",
			projectInfo.GitInfo.Branch, projectInfo.GitInfo.CommitHash))
	}

	// Count and add included files
	fileCount, err := countIncludedFiles(config)
	if err != nil {
		return "", fmt.Errorf("error counting files: %w", err)
	}
	summary.WriteString(fmt.Sprintf("   Files: %d included\n", fileCount))

	// Add filtering info if specified
	if len(config.Extensions) > 0 {
		summary.WriteString(fmt.Sprintf("   Filtering: %s\n",
			strings.Join(config.Extensions, ", ")))
	}

	// Add token count to summary
	if tokenCount > 0 {
		summary.WriteString(fmt.Sprintf("   Tokens: ~%d\n", tokenCount))
	}

	return summary.String(), nil
}

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool) error {
	// Enable debug logging if flag is set
	if debug {
		log.Enable()
	}

	log.Debug("Starting promptext with dir: %s", dirPath)
	// Validate format
	formatter, err := format.GetFormatter(outputFormat)
	if err != nil {
		return fmt.Errorf("invalid format: %w", err)
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
	log.Debug("Using extensions: %v", extensions)
	log.Debug("Using excludes: %#v", excludes)

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
		// Only display project summary
		if info, err := GetMetadataSummary(procConfig, 0); err == nil {
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
			fmt.Printf("\033[32m%s✓ code context written to %s (%s format)\033[0m\n",
				info, outFile, outputFormat)
		}
	} else if !noCopy {
		// Copy to clipboard if no output file is specified and clipboard is not disabled
		if err := clipboard.WriteAll(formattedOutput); err != nil {
			log.Info("Warning: Failed to copy to clipboard: %v", err)
		} else {
			// Always print metadata summary and success message in green
			if info, err := GetMetadataSummary(procConfig, result.TokenCount); err == nil {
				fmt.Printf("\033[32m%s✓ code context copied to clipboard (%s format)\033[0m\n",
					info, outputFormat)
			}
		}
	}

	return nil
}
