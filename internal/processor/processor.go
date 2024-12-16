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

func ProcessDirectory(config Config, verbose bool) (string, error) {
	var builder strings.Builder

	// Get project information
	projectInfo, err := info.GetProjectInfo(config.DirPath)
	if err != nil {
		return "", fmt.Errorf("error getting project info: %w", err)
	}

	if verbose {
		// Add directory tree
		builder.WriteString(projectInfo.DirectoryTree)

		// Add git information if available
		if projectInfo.GitInfo != nil {
			builder.WriteString("\n### Git Information:\n")
			builder.WriteString(fmt.Sprintf("Branch: %s\n", projectInfo.GitInfo.Branch))
			builder.WriteString(fmt.Sprintf("Commit: %s\n", projectInfo.GitInfo.CommitHash))
			builder.WriteString(fmt.Sprintf("Message: %s\n", projectInfo.GitInfo.CommitMessage))
		}

		// Add project metadata if available
		if projectInfo.Metadata != nil {
			builder.WriteString("\n### Project Metadata:\n")
			builder.WriteString(fmt.Sprintf("Language: %s\n", projectInfo.Metadata.Language))
			builder.WriteString(fmt.Sprintf("Version: %s\n", projectInfo.Metadata.Version))
			if len(projectInfo.Metadata.Dependencies) > 0 {
				builder.WriteString("Dependencies:\n")
				for _, dep := range projectInfo.Metadata.Dependencies {
					builder.WriteString(fmt.Sprintf("  - %s\n", dep))
				}
			}
		}
	}

	// Initialize gitignore
	gitIgnore, err := gitignore.New(filepath.Join(config.DirPath, ".gitignore"))
	if err != nil {
		return "", fmt.Errorf("error reading .gitignore: %w", err)
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

		// Only process content in verbose mode
		if verbose {
			content, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("error reading file %s: %w", path, err)
			}
			
			builder.WriteString(fmt.Sprintf("\n### File: %s\n", path))
			builder.WriteString("```\n")
			builder.Write(content)
			builder.WriteString("\n```\n")
		}

		return nil
	})

	if err != nil {
		return "", fmt.Errorf("error walking directory: %w", err)
	}

	return builder.String(), nil
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
