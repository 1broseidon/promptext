package format

import (
    "encoding/json"
    "encoding/xml"
    "fmt"
    "strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}
type JSONFormatter struct{}

func (m *MarkdownFormatter) Format(project *ProjectOutput) (string, error) {
    var sb strings.Builder

    // Add directory tree
    sb.WriteString(project.DirectoryTree)

    // Add git information if available
    if project.GitInfo != nil {
        sb.WriteString("\n### Git Information:\n")
        sb.WriteString(fmt.Sprintf("Branch: %s\n", project.GitInfo.Branch))
        sb.WriteString(fmt.Sprintf("Commit: %s\n", project.GitInfo.CommitHash))
        sb.WriteString(fmt.Sprintf("Message: %s\n", project.GitInfo.CommitMessage))
    }

    // Add metadata if available
    if project.Metadata != nil {
        sb.WriteString("\n### Project Metadata:\n")
        sb.WriteString(fmt.Sprintf("Language: %s\n", project.Metadata.Language))
        sb.WriteString(fmt.Sprintf("Version: %s\n", project.Metadata.Version))
        if len(project.Metadata.Dependencies) > 0 {
            sb.WriteString("Dependencies:\n")
            for _, dep := range project.Metadata.Dependencies {
                sb.WriteString(fmt.Sprintf("  - %s\n", dep))
            }
        }
    }

    // Add files if available
    if len(project.Files) > 0 {
        for _, file := range project.Files {
            sb.WriteString(fmt.Sprintf("\n### File: %s\n", file.Path))
            sb.WriteString("```\n")
            sb.WriteString(file.Content)
            sb.WriteString("\n```\n")
        }
    }

    return sb.String(), nil
}

func (x *XMLFormatter) Format(project *ProjectOutput) (string, error) {
    output, err := xml.MarshalIndent(project, "", "  ")
    if err != nil {
        return "", err
    }
    return string(output), nil
}

func (j *JSONFormatter) Format(project *ProjectOutput) (string, error) {
    output, err := json.MarshalIndent(project, "", "  ")
    if err != nil {
        return "", err
    }
    return string(output), nil
}
