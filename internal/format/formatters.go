package format

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}

func (m *MarkdownFormatter) Format(project *ProjectOutput) (string, error) {
	var sb strings.Builder

	// Add directory tree
	sb.WriteString("### Project Structure:\n```\n")
	sb.WriteString(project.DirectoryTree.ToMarkdown(0))
	sb.WriteString("```\n")

	// Add git information if available
	if project.GitInfo != nil {
		sb.WriteString("\n### Git Information:\n")
		sb.WriteString(fmt.Sprintf("Branch: %s\n", project.GitInfo.Branch))
		sb.WriteString(fmt.Sprintf("Commit: %s\n", project.GitInfo.CommitHash))
		sb.WriteString(fmt.Sprintf("Message: %s\n", project.GitInfo.CommitMessage))
	}

	// Add detailed metadata section if available
	if project.Metadata != nil && len(project.Metadata.Dependencies) > 0 {
		sb.WriteString("\n## Dependencies\n")
		sb.WriteString("```\n")
		for _, dep := range project.Metadata.Dependencies {
			sb.WriteString(fmt.Sprintf("%s\n", dep))
		}
		sb.WriteString("```\n")
	}

	// Add files section with better formatting
	if len(project.Files) > 0 {
		sb.WriteString("\n## Source Files\n")
		for _, file := range project.Files {
			// Extract file extension for syntax highlighting
			ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
			if ext == "" {
				ext = "text"
			}

			sb.WriteString(fmt.Sprintf("\n### %s\n", file.Path))
			sb.WriteString(fmt.Sprintf("```%s\n", ext))
			sb.WriteString(file.Content)
			sb.WriteString("\n```\n")
		}
	}

	return sb.String(), nil
}

// Helper function to write directory nodes as XML
func writeDirectoryNode(node *DirectoryNode, b *strings.Builder, indent int) {
	if node == nil {
		return
	}

	indentStr := strings.Repeat(" ", indent)

	if node.Type != "" { // Skip root node
		b.WriteString(fmt.Sprintf("%s<node name=\"%s\" type=\"%s\"", indentStr, node.Name, node.Type))
		if len(node.Children) == 0 {
			b.WriteString("/>\n")
			return
		}
		b.WriteString(">\n")
	}

	for _, child := range node.Children {
		writeDirectoryNode(child, b, indent+2)
	}

	if node.Type != "" {
		b.WriteString(fmt.Sprintf("%s</node>\n", indentStr))
	}
}

func (x *XMLFormatter) Format(project *ProjectOutput) (string, error) {
	// Create a custom encoder that uses indentation
	var b strings.Builder
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")

	// Start with XML header
	b.WriteString(xml.Header)

	// Start the project element
	b.WriteString("<project>\n")

	// Write directory tree as structured XML
	b.WriteString("  <directoryTree>\n")
	writeDirectoryNode(project.DirectoryTree, &b, 4)
	b.WriteString("  </directoryTree>\n")

	// Write git info if available
	if project.GitInfo != nil {
		b.WriteString("  <gitInfo>\n")
		b.WriteString(fmt.Sprintf("    <branch>%s</branch>\n", project.GitInfo.Branch))
		b.WriteString(fmt.Sprintf("    <commitHash>%s</commitHash>\n", project.GitInfo.CommitHash))
		b.WriteString("    <commitMessage><![CDATA[")
		b.WriteString(project.GitInfo.CommitMessage)
		b.WriteString("]]></commitMessage>\n")
		b.WriteString("  </gitInfo>\n")
	}

	// Write metadata if available
	if project.Metadata != nil {
		b.WriteString("  <metadata>\n")
		b.WriteString(fmt.Sprintf("    <language>%s</language>\n", project.Metadata.Language))
		b.WriteString(fmt.Sprintf("    <version>%s</version>\n", project.Metadata.Version))
		if len(project.Metadata.Dependencies) > 0 {
			b.WriteString("    <dependencies>\n")
			for _, dep := range project.Metadata.Dependencies {
				b.WriteString(fmt.Sprintf("      <dependency>%s</dependency>\n", dep))
			}
			b.WriteString("    </dependencies>\n")
		}
		b.WriteString("  </metadata>\n")
	}

	// Write files if available
	if len(project.Files) > 0 {
		b.WriteString("  <files>\n")
		for _, file := range project.Files {
			b.WriteString(fmt.Sprintf("    <file path=\"%s\">\n", file.Path))
			b.WriteString("      <content><![CDATA[")
			b.WriteString(file.Content)
			b.WriteString("]]></content>\n")
			b.WriteString("    </file>\n")
		}
		b.WriteString("  </files>\n")
	}

	b.WriteString("</project>")

	return b.String(), nil
}
