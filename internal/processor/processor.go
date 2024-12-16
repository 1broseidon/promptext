package processor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/gitignore"
	"github.com/1broseidon/promptext/internal/info"
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
	DisplayContent  string
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
	projectInfo, err := info.GetProjectInfo(config.DirPath, &config, gitIgnore)
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

	// Initialize gitignore
	gitIgnore, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error reading .gitignore: %w", err)
	}

	err = filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Skip if file doesn't match our filters
		if !filter.ShouldProcessFile(path, config.Extensions, config.Excludes, gitIgnore) {
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
	projectInfo, err := info.GetProjectInfo(config.DirPath)
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
		if filter.ShouldProcessFile(path, config.Extensions, config.Excludes, &gitignore.GitIgnore{}) {
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
