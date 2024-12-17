package processor

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/config"
	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/gitignore"
	"github.com/1broseidon/promptext/internal/info"
	"github.com/1broseidon/promptext/internal/references"
	"github.com/atotto/clipboard"
)

func getAllFiles(root string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		files = append(files, rel)
		return nil
	})
	return files, err
}

type Config struct {
	DirPath    string
	Extensions []string
	Excludes   []string
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
}

func ProcessDirectory(config Config, verbose bool) (*ProcessResult, error) {
	var displayBuilder strings.Builder
	projectOutput := &format.ProjectOutput{}

	// Get all files first for reference resolution
	allFiles, err := getAllFiles(config.DirPath)
	if err != nil {
		return nil, fmt.Errorf("could not list all files: %w", err)
	}

	// Get project analysis and convert to format.ProjectAnalysis
	analysis := info.AnalyzeProject(config.DirPath)
	projectOutput.Analysis = &format.ProjectAnalysis{
		EntryPoints:   analysis.EntryPoints,
		ConfigFiles:   analysis.ConfigFiles,
		CoreFiles:     analysis.CoreFiles,
		TestFiles:     analysis.TestFiles,
		Documentation: analysis.Documentation,
	}

	// Initialize gitignore once
	gi, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error reading .gitignore: %w", err)
	}

	// Get project information with filtering config
	infoConfig := &info.Config{
		Extensions: config.Extensions,
		Excludes:   config.Excludes,
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, infoConfig, gi)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}

	// Populate ProjectOutput
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

	// Only add to display if verbose
	if verbose {
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
	}

	err = filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Create unified filter once
		unifiedFilter := filter.NewUnifiedFilter(gi, config.Extensions, config.Excludes)

		// Skip if file doesn't match our filters
		if !unifiedFilter.ShouldProcess(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Add file to ProjectOutput
		rel, err := filepath.Rel(config.DirPath, path)
		if err != nil {
			return fmt.Errorf("error getting relative path for %s: %w", path, err)
		}

		refs := references.ExtractFileReferences(string(content), filepath.Dir(rel), config.DirPath, allFiles)

		projectOutput.Files = append(projectOutput.Files, format.FileInfo{
			Path:       rel,
			Content:    string(content),
			References: refs,
		})

		// Only add to display if verbose
		if verbose {
			displayBuilder.WriteString(fmt.Sprintf("\n### File: %s\n```\n%s\n```\n", path, content))
		}

		return nil
	})

	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error walking directory: %w", err)
	}

	return &ProcessResult{
		ProjectOutput:    projectOutput,
		DisplayContent:   displayBuilder.String(),
		ClipboardContent: "", // Will be set in Run() based on format
	}, nil
}

// GetMetadataSummary returns a concise summary of project metadata
func GetMetadataSummary(config Config) (string, error) {
	gi, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return "", err
	}

	infoConfig := &info.Config{
		Extensions: config.Extensions,
		Excludes:   config.Excludes,
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, infoConfig, gi)
	if err != nil {
		return "", err
	}

	var summary strings.Builder
	summary.WriteString("📦 Project Summary:\n")

	// Add project name from metadata if available, otherwise use folder name
	if projectInfo.Metadata != nil && projectInfo.Metadata.Name != "" {
		summary.WriteString(fmt.Sprintf("   Project: %s\n", projectInfo.Metadata.Name))
	} else {
		absPath, err := filepath.Abs(config.DirPath)
		if err == nil {
			summary.WriteString(fmt.Sprintf("   Project: %s\n", filepath.Base(absPath)))
		}
	}

	if projectInfo.Metadata != nil {
		summary.WriteString(fmt.Sprintf("   Language: %s", projectInfo.Metadata.Language))
		if projectInfo.Metadata.Version != "" {
			summary.WriteString(fmt.Sprintf(" %s", projectInfo.Metadata.Version))
		}
		summary.WriteString("\n")
		if len(projectInfo.Metadata.Dependencies) > 0 {
			mainDeps := 0
			devDeps := 0
			for _, dep := range projectInfo.Metadata.Dependencies {
				if strings.HasPrefix(dep, "[dev] ") {
					devDeps++
				} else {
					mainDeps++
				}
			}
			if devDeps > 0 {
				summary.WriteString(fmt.Sprintf("   Dependencies: %d packages (%d main, %d dev)\n",
					len(projectInfo.Metadata.Dependencies), mainDeps, devDeps))
			} else {
				summary.WriteString(fmt.Sprintf("   Dependencies: %d packages\n", mainDeps))
			}
		}
	}

	if projectInfo.GitInfo != nil {
		summary.WriteString(fmt.Sprintf("   Branch: %s (%s)\n", projectInfo.GitInfo.Branch, projectInfo.GitInfo.CommitHash))
	}

	// Count included files
	fileCount := 0
	filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(config.DirPath, path)
		if filter.ShouldProcessFile(rel, config.Extensions, config.Excludes, gi) {
			fileCount++
		}
		return nil
	})
	summary.WriteString(fmt.Sprintf("   Files: %d included\n", fileCount))

	if len(config.Extensions) > 0 {
		summary.WriteString(fmt.Sprintf("   Filtering: %s\n", strings.Join(config.Extensions, ", ")))
	}

	return summary.String(), nil
}

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string) error {
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
		log.Printf("Warning: Failed to load .promptext.yml from %s: %v", absPath, err)
		fileConfig = &config.FileConfig{}
	}

	// Merge file config with command line flags
	extensions, excludes, verboseFlag := fileConfig.MergeWithFlags(extension, exclude, verbose)

	// Create processor configuration
	procConfig := Config{
		DirPath:    absPath,
		Extensions: extensions,
		Excludes:   excludes,
	}

	if infoOnly {
		// Only display project summary
		if info, err := GetMetadataSummary(procConfig); err == nil {
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
		if info, err := GetMetadataSummary(procConfig); err == nil {
			fmt.Printf("\033[32m%s   ✓ code context written to %s (%s format)\033[0m\n",
				info, outFile, outputFormat)
		}
	} else if !noCopy {
		// Copy to clipboard if no output file is specified and clipboard is not disabled
		if err := clipboard.WriteAll(formattedOutput); err != nil {
			log.Printf("Warning: Failed to copy to clipboard: %v", err)
		} else {
			// Always print metadata summary and success message in green
			if info, err := GetMetadataSummary(procConfig); err == nil {
				fmt.Printf("\033[32m%s   ✓ code context copied to clipboard (%s format)\033[0m\n",
					info, outputFormat)
			}
		}
	}

	return nil
}
