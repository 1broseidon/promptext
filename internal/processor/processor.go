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

func ProcessDirectory(config Config) (string, error) {
	var builder strings.Builder

	// Get project information
	projectInfo, err := info.GetProjectInfo(config.DirPath)
	if err != nil {
		return "", fmt.Errorf("error getting project info: %w", err)
	}

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

		if !filter.ShouldProcessFile(path, config.Extensions, config.Excludes, gitIgnore) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("error reading file %s: %w", path, err)
		}

		// Add file header
		builder.WriteString(fmt.Sprintf("\n### File: %s\n", path))
		builder.WriteString("```\n")
		builder.Write(content)
		builder.WriteString("\n```\n")

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
	
	if projectInfo.Metadata != nil {
		summary.WriteString(fmt.Sprintf("   Language: %s %s\n", projectInfo.Metadata.Language, projectInfo.Metadata.Version))
	}
	
	if projectInfo.GitInfo != nil {
		summary.WriteString(fmt.Sprintf("   Branch: %s (%s)\n", projectInfo.GitInfo.Branch, projectInfo.GitInfo.CommitHash))
	}

	if len(config.Extensions) > 0 {
		summary.WriteString(fmt.Sprintf("   Filtering: %s\n", strings.Join(config.Extensions, ", ")))
	} else {
		summary.WriteString("\n")
	}

	return summary.String(), nil
}
