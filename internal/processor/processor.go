package processor

import (
	"fmt"
	"log"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"promptext/internal/config"
	"promptext/internal/filter"
	"promptext/internal/gitignore"
	"promptext/internal/info"
	"github.com/atotto/clipboard"
)

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
	DisplayContent   string
	ClipboardContent string
}

func ProcessDirectory(config Config, verbose bool) (*ProcessResult, error) {
	var displayBuilder, clipBuilder strings.Builder

	// Initialize gitignore once
	gitIgnore, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error reading .gitignore: %w", err)
	}

	// Get project information with filtering config
	infoConfig := &info.Config{
		Extensions: config.Extensions,
		Excludes:   config.Excludes,
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, infoConfig, gitIgnore)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}

	// Always add full content to clipboard
	clipBuilder.WriteString(projectInfo.DirectoryTree)

	// Add git information if available
	if projectInfo.GitInfo != nil {
		clipBuilder.WriteString("\n### Git Information:\n")
		clipBuilder.WriteString(fmt.Sprintf("Branch: %s\n", projectInfo.GitInfo.Branch))
		clipBuilder.WriteString(fmt.Sprintf("Commit: %s\n", projectInfo.GitInfo.CommitHash))
		clipBuilder.WriteString(fmt.Sprintf("Message: %s\n", projectInfo.GitInfo.CommitMessage))
	}

	// Add project metadata if available
	if projectInfo.Metadata != nil {
		clipBuilder.WriteString("\n### Project Metadata:\n")
		clipBuilder.WriteString(fmt.Sprintf("Language: %s\n", projectInfo.Metadata.Language))
		clipBuilder.WriteString(fmt.Sprintf("Version: %s\n", projectInfo.Metadata.Version))
		if len(projectInfo.Metadata.Dependencies) > 0 {
			clipBuilder.WriteString("Dependencies:\n")
			for _, dep := range projectInfo.Metadata.Dependencies {
				clipBuilder.WriteString(fmt.Sprintf("  - %s\n", dep))
			}
		}
	}

	// Only add to display if verbose
	if verbose {
		displayBuilder.WriteString(clipBuilder.String())
	}

	err = filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Create unified filter once
		unifiedFilter := filter.NewUnifiedFilter(gitIgnore, config.Extensions, config.Excludes)

		// Skip if file doesn't match our filters
		if !unifiedFilter.ShouldProcess(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Always add to clipboard content
		clipBuilder.WriteString(fmt.Sprintf("\n### File: %s\n", path))
		clipBuilder.WriteString("```\n")
		clipBuilder.Write(content)
		clipBuilder.WriteString("\n```\n")

		// Only add to display if verbose
		if verbose {
			displayBuilder.WriteString(fmt.Sprintf("\n### File: %s\n", path))
			displayBuilder.WriteString("```\n")
			displayBuilder.Write(content)
			displayBuilder.WriteString("\n```\n")
		}

		return nil
	})

	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error walking directory: %w", err)
	}

	return &ProcessResult{
		DisplayContent:   displayBuilder.String(),
		ClipboardContent: clipBuilder.String(),
	}, nil
}

// GetMetadataSummary returns a concise summary of project metadata
func GetMetadataSummary(config Config) (string, error) {
	gitIgnore, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return "", err
	}

	infoConfig := &info.Config{
		Extensions: config.Extensions,
		Excludes:   config.Excludes,
	}
	projectInfo, err := info.GetProjectInfo(config.DirPath, infoConfig, gitIgnore)
	if err != nil {
		return "", err
	}

	var summary strings.Builder
	summary.WriteString("📦 Project Summary:\n")

	// Add root folder name
	absPath, err := filepath.Abs(config.DirPath)
	if err == nil {
		summary.WriteString(fmt.Sprintf("   Project: %s\n", filepath.Base(absPath)))
	}

	if projectInfo.Metadata != nil {
		summary.WriteString(fmt.Sprintf("   Language: %s %s\n", projectInfo.Metadata.Language, projectInfo.Metadata.Version))
		if len(projectInfo.Metadata.Dependencies) > 0 {
			summary.WriteString(fmt.Sprintf("   Dependencies: %d packages\n", len(projectInfo.Metadata.Dependencies)))
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
		if filter.ShouldProcessFile(rel, config.Extensions, config.Excludes, gitIgnore) {
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
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool) error {
	// Load config file
	fileConfig, err := config.LoadConfig(dirPath)
	if err != nil {
		log.Printf("Warning: Failed to load .promptext.yml: %v", err)
		fileConfig = &config.FileConfig{}
	}

	// Merge file config with command line flags
	extensions, excludes, verboseFlag := fileConfig.MergeWithFlags(extension, exclude, verbose)

	// Create processor configuration
	procConfig := Config{
		DirPath:    dirPath,
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

	// Copy to clipboard unless disabled
	if !noCopy {
		if err := clipboard.WriteAll(result.ClipboardContent); err != nil {
			log.Printf("Warning: Failed to copy to clipboard: %v", err)
		}
		// Always print metadata summary and success message in green
		if info, err := GetMetadataSummary(procConfig); err == nil {
			fmt.Printf("\033[32m%s   ✓ code context copied to clipboard\033[0m\n", info)
		}
	}

	return nil
}
