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

	// Project Overview Section
	if project.Overview != nil {
		sb.WriteString("# Project Overview\n\n")
		sb.WriteString(fmt.Sprintf("%s\n\n", project.Overview.Description))
		
		if len(project.Overview.Features) > 0 {
			sb.WriteString("## Key Features\n")
			for _, feature := range project.Overview.Features {
				sb.WriteString(fmt.Sprintf("- %s\n", feature))
			}
			sb.WriteString("\n")
		}
	}

	// Quick Reference Section
	sb.WriteString("## Quick Reference\n\n")
	if project.Analysis != nil {
		// Entry Points
		if len(project.Analysis.EntryPoints) > 0 {
			sb.WriteString("### Entry Points\n")
			for path, desc := range project.Analysis.EntryPoints {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}

		// Config Files
		if len(project.Analysis.ConfigFiles) > 0 {
			sb.WriteString("### Configuration Files\n")
			for path, desc := range project.Analysis.ConfigFiles {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}

		// Core Files
		if len(project.Analysis.CoreFiles) > 0 {
			sb.WriteString("### Core Components\n")
			for path, desc := range project.Analysis.CoreFiles {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}

		// Test Files
		if len(project.Analysis.TestFiles) > 0 {
			sb.WriteString("### Tests\n")
			for path, desc := range project.Analysis.TestFiles {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}

		// Documentation Files
		if len(project.Analysis.Documentation) > 0 {
			sb.WriteString("### Documentation\n")
			for path, desc := range project.Analysis.Documentation {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}
	}

	if project.FileStats != nil {
		sb.WriteString(fmt.Sprintf("- Total Files: %d\n", project.FileStats.TotalFiles))
		sb.WriteString(fmt.Sprintf("- Total Lines: %d\n", project.FileStats.TotalLines))
		sb.WriteString(fmt.Sprintf("- Packages: %d\n", project.FileStats.PackageCount))
		
		sb.WriteString("\nFile Types:\n")
		for ext, count := range project.FileStats.FilesByType {
			sb.WriteString(fmt.Sprintf("- %s: %d files\n", ext, count))
		}
		sb.WriteString("\n")
	}

	// Directory Structure with annotations
	sb.WriteString("## Project Structure\n```\n")
	sb.WriteString(project.DirectoryTree.ToMarkdown(0))
	sb.WriteString("```\n")

	// Git Information
	if project.GitInfo != nil {
		sb.WriteString("\n## Git Information\n")
		sb.WriteString(fmt.Sprintf("- Branch: %s\n", project.GitInfo.Branch))
		sb.WriteString(fmt.Sprintf("- Commit: %s\n", project.GitInfo.CommitHash))
		sb.WriteString(fmt.Sprintf("- Message: %s\n", project.GitInfo.CommitMessage))
	}

	// Dependencies and Relationships
	if project.Dependencies != nil {
		sb.WriteString("\n## Package Dependencies\n")
		for pkg := range project.Dependencies.Imports {
			sb.WriteString(fmt.Sprintf("### %s\n", pkg))
			for _, imp := range project.Dependencies.Imports[pkg] {
				sb.WriteString(fmt.Sprintf("- %s\n", imp))
			}
		}

		if len(project.Dependencies.CoreFiles) > 0 {
			sb.WriteString("\n### Core Components\n")
			for _, file := range project.Dependencies.CoreFiles {
				sb.WriteString(fmt.Sprintf("- %s\n", file))
			}
		}
	}

	// Source Files with line counts
	if len(project.Files) > 0 {
		sb.WriteString("\n## Source Files\n")
		for _, file := range project.Files {
			ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
			if ext == "" {
				ext = "text"
			}

			lineCount := strings.Count(file.Content, "\n") + 1
			sb.WriteString(fmt.Sprintf("\n### %s (%d lines)\n", file.Path, lineCount))
			sb.WriteString(fmt.Sprintf("```%s\n", ext))
			sb.WriteString(file.Content)
			sb.WriteString("\n```\n")

			// Add references section if any exist
			if file.References != nil && (len(file.References.Internal) > 0 || len(file.References.External) > 0) {
				sb.WriteString("\n**References:**\n")
				if len(file.References.Internal) > 0 {
					sb.WriteString("Internal:\n")
					for dir, refs := range file.References.Internal {
						for _, ref := range refs {
							sb.WriteString(fmt.Sprintf("- `%s` references `%s`\n", dir, ref))
						}
					}
				}
				if len(file.References.External) > 0 {
					sb.WriteString("External:\n")
					for dir, refs := range file.References.External {
						for _, ref := range refs {
							sb.WriteString(fmt.Sprintf("- `%s` references `%s`\n", dir, ref))
						}
					}
				}
				sb.WriteString("\n")
			}
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
	var b strings.Builder
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")

	b.WriteString(xml.Header)
	b.WriteString("<project>\n")

	// Project Overview
	if project.Overview != nil {
		b.WriteString("  <overview>\n")
		b.WriteString(fmt.Sprintf("    <description><![CDATA[%s]]></description>\n", project.Overview.Description))
		b.WriteString(fmt.Sprintf("    <purpose><![CDATA[%s]]></purpose>\n", project.Overview.Purpose))
		if len(project.Overview.Features) > 0 {
			b.WriteString("    <features>\n")
			for _, feature := range project.Overview.Features {
				b.WriteString(fmt.Sprintf("      <feature>%s</feature>\n", feature))
			}
			b.WriteString("    </features>\n")
		}
		b.WriteString("  </overview>\n")
	}

	// File Statistics
	if project.FileStats != nil {
		b.WriteString("  <fileStats>\n")
		b.WriteString(fmt.Sprintf("    <totalFiles>%d</totalFiles>\n", project.FileStats.TotalFiles))
		b.WriteString(fmt.Sprintf("    <totalLines>%d</totalLines>\n", project.FileStats.TotalLines))
		b.WriteString(fmt.Sprintf("    <packageCount>%d</packageCount>\n", project.FileStats.PackageCount))
		if len(project.FileStats.FilesByType) > 0 {
			b.WriteString("    <fileTypes>\n")
			for ext, count := range project.FileStats.FilesByType {
				b.WriteString(fmt.Sprintf("      <type ext=\"%s\">%d</type>\n", ext, count))
			}
			b.WriteString("    </fileTypes>\n")
		}
		b.WriteString("  </fileStats>\n")
	}

	// Directory Tree
	b.WriteString("  <directoryTree>\n")
	writeDirectoryNode(project.DirectoryTree, &b, 4)
	b.WriteString("  </directoryTree>\n")

	// Git Info
	if project.GitInfo != nil {
		b.WriteString("  <gitInfo>\n")
		b.WriteString(fmt.Sprintf("    <branch>%s</branch>\n", project.GitInfo.Branch))
		b.WriteString(fmt.Sprintf("    <commitHash>%s</commitHash>\n", project.GitInfo.CommitHash))
		b.WriteString("    <commitMessage><![CDATA[")
		b.WriteString(project.GitInfo.CommitMessage)
		b.WriteString("]]></commitMessage>\n")
		b.WriteString("  </gitInfo>\n")
	}

	// Dependencies and Relationships
	if project.Dependencies != nil {
		b.WriteString("  <dependencies>\n")
		if len(project.Dependencies.Imports) > 0 {
			b.WriteString("    <imports>\n")
			for file, imports := range project.Dependencies.Imports {
				b.WriteString(fmt.Sprintf("      <file path=\"%s\">\n", file))
				for _, imp := range imports {
					b.WriteString(fmt.Sprintf("        <import>%s</import>\n", imp))
				}
				b.WriteString("      </file>\n")
			}
			b.WriteString("    </imports>\n")
		}
		if len(project.Dependencies.CoreFiles) > 0 {
			b.WriteString("    <coreFiles>\n")
			for _, file := range project.Dependencies.CoreFiles {
				b.WriteString(fmt.Sprintf("      <file>%s</file>\n", file))
			}
			b.WriteString("    </coreFiles>\n")
		}
		b.WriteString("  </dependencies>\n")
	}

	// Files with metadata
	if len(project.Files) > 0 {
		b.WriteString("  <files>\n")
		for _, file := range project.Files {
			lineCount := strings.Count(file.Content, "\n") + 1
			b.WriteString(fmt.Sprintf("    <file path=\"%s\" lines=\"%d\">\n", file.Path, lineCount))
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
