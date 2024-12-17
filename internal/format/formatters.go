package format

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}

func (m *MarkdownFormatter) formatOverview(sb *strings.Builder, overview *ProjectOverview) {
	if overview == nil {
		return
	}
	sb.WriteString("# Project Overview\n\n")
	sb.WriteString(fmt.Sprintf("%s\n\n", overview.Description))
	
	if len(overview.Features) > 0 {
		sb.WriteString("## Key Features\n")
		for _, feature := range overview.Features {
			sb.WriteString(fmt.Sprintf("- %s\n", feature))
		}
		sb.WriteString("\n")
	}
}

func (m *MarkdownFormatter) formatAnalysis(sb *strings.Builder, analysis *ProjectAnalysis) {
	if analysis == nil {
		return
	}

	sections := []struct {
		title string
		items map[string]string
	}{
		{"Entry Points", analysis.EntryPoints},
		{"Configuration Files", analysis.ConfigFiles},
		{"Core Components", analysis.CoreFiles},
		{"Tests", analysis.TestFiles},
		{"Documentation", analysis.Documentation},
	}

	for _, section := range sections {
		if len(section.items) > 0 {
			sb.WriteString(fmt.Sprintf("### %s\n", section.title))
			for path, desc := range section.items {
				sb.WriteString(fmt.Sprintf("- %s: %s\n", path, desc))
			}
			sb.WriteString("\n")
		}
	}
}

func (m *MarkdownFormatter) formatFileStats(sb *strings.Builder, stats *FileStatistics) {
	if stats == nil {
		return
	}
	sb.WriteString(fmt.Sprintf("- Total Files: %d\n", stats.TotalFiles))
	sb.WriteString(fmt.Sprintf("- Total Lines: %d\n", stats.TotalLines))
	sb.WriteString(fmt.Sprintf("- Packages: %d\n", stats.PackageCount))
	
	sb.WriteString("\nFile Types:\n")
	for ext, count := range stats.FilesByType {
		sb.WriteString(fmt.Sprintf("- %s: %d files\n", ext, count))
	}
	sb.WriteString("\n")
}

func (m *MarkdownFormatter) formatGitInfo(sb *strings.Builder, gitInfo *GitInfo) {
	if gitInfo == nil {
		return
	}
	sb.WriteString("\n## Git Information\n")
	sb.WriteString(fmt.Sprintf("- Branch: %s\n", gitInfo.Branch))
	sb.WriteString(fmt.Sprintf("- Commit: %s\n", gitInfo.CommitHash))
	sb.WriteString(fmt.Sprintf("- Message: %s\n", gitInfo.CommitMessage))
}

func (m *MarkdownFormatter) formatDependencies(sb *strings.Builder, deps *DependencyInfo) {
	if deps == nil {
		return
	}
	sb.WriteString("\n## Package Dependencies\n")
	for pkg := range deps.Imports {
		sb.WriteString(fmt.Sprintf("### %s\n", pkg))
		for _, imp := range deps.Imports[pkg] {
			sb.WriteString(fmt.Sprintf("- %s\n", imp))
		}
	}

	if len(deps.CoreFiles) > 0 {
		sb.WriteString("\n### Core Components\n")
		for _, file := range deps.CoreFiles {
			sb.WriteString(fmt.Sprintf("- %s\n", file))
		}
	}
}

func (m *MarkdownFormatter) formatSourceFiles(sb *strings.Builder, files []FileInfo) {
	if len(files) == 0 {
		return
	}
	sb.WriteString("\n## Source Files\n")
	for _, file := range files {
		ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
		if ext == "" {
			ext = "text"
		}

		lineCount := strings.Count(file.Content, "\n") + 1
		sb.WriteString(fmt.Sprintf("\n### %s (%d lines)\n", file.Path, lineCount))
		sb.WriteString(fmt.Sprintf("```%s\n", ext))
		sb.WriteString(file.Content)
		sb.WriteString("\n```\n")

		m.formatFileReferences(sb, file.References)
	}
}

func (m *MarkdownFormatter) formatFileReferences(sb *strings.Builder, refs *references.ReferenceMap) {
	if refs == nil || (len(refs.Internal) == 0 && len(refs.External) == 0) {
		return
	}
	
	sb.WriteString("\n**References:**\n")
	if len(refs.Internal) > 0 {
		sb.WriteString("Internal:\n")
		for dir, refs := range refs.Internal {
			for _, ref := range refs {
				sb.WriteString(fmt.Sprintf("- `%s` references `%s`\n", dir, ref))
			}
		}
	}
	if len(refs.External) > 0 {
		sb.WriteString("External:\n")
		for dir, refs := range refs.External {
			for _, ref := range refs {
				sb.WriteString(fmt.Sprintf("- `%s` references `%s`\n", dir, ref))
			}
		}
	}
	sb.WriteString("\n")
}

func (m *MarkdownFormatter) Format(project *ProjectOutput) (string, error) {
	var sb strings.Builder

	m.formatOverview(&sb, project.Overview)
	
	sb.WriteString("## Quick Reference\n\n")
	m.formatAnalysis(&sb, project.Analysis)
	m.formatFileStats(&sb, project.FileStats)

	// Directory Structure
	sb.WriteString("## Project Structure\n```\n")
	sb.WriteString(project.DirectoryTree.ToMarkdown(0))
	sb.WriteString("```\n")

	m.formatGitInfo(&sb, project.GitInfo)
	m.formatDependencies(&sb, project.Dependencies)
	m.formatSourceFiles(&sb, project.Files)

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

func (x *XMLFormatter) formatOverview(b *strings.Builder, overview *ProjectOverview) {
	if overview == nil {
		return
	}
	b.WriteString("  <overview>\n")
	b.WriteString(fmt.Sprintf("    <description><![CDATA[%s]]></description>\n", overview.Description))
	b.WriteString(fmt.Sprintf("    <purpose><![CDATA[%s]]></purpose>\n", overview.Purpose))
	if len(overview.Features) > 0 {
		b.WriteString("    <features>\n")
		for _, feature := range overview.Features {
			b.WriteString(fmt.Sprintf("      <feature>%s</feature>\n", feature))
		}
		b.WriteString("    </features>\n")
	}
	b.WriteString("  </overview>\n")
}

func (x *XMLFormatter) formatFileStats(b *strings.Builder, stats *FileStatistics) {
	if stats == nil {
		return
	}
	b.WriteString("  <fileStats>\n")
	b.WriteString(fmt.Sprintf("    <totalFiles>%d</totalFiles>\n", stats.TotalFiles))
	b.WriteString(fmt.Sprintf("    <totalLines>%d</totalLines>\n", stats.TotalLines))
	b.WriteString(fmt.Sprintf("    <packageCount>%d</packageCount>\n", stats.PackageCount))
	if len(stats.FilesByType) > 0 {
		b.WriteString("    <fileTypes>\n")
		for ext, count := range stats.FilesByType {
			b.WriteString(fmt.Sprintf("      <type ext=\"%s\">%d</type>\n", ext, count))
		}
		b.WriteString("    </fileTypes>\n")
	}
	b.WriteString("  </fileStats>\n")
}

func (x *XMLFormatter) formatGitInfo(b *strings.Builder, gitInfo *GitInfo) {
	if gitInfo == nil {
		return
	}
	b.WriteString("  <gitInfo>\n")
	b.WriteString(fmt.Sprintf("    <branch>%s</branch>\n", gitInfo.Branch))
	b.WriteString(fmt.Sprintf("    <commitHash>%s</commitHash>\n", gitInfo.CommitHash))
	b.WriteString("    <commitMessage><![CDATA[")
	b.WriteString(gitInfo.CommitMessage)
	b.WriteString("]]></commitMessage>\n")
	b.WriteString("  </gitInfo>\n")
}

func (x *XMLFormatter) formatDependencies(b *strings.Builder, deps *DependencyInfo) {
	if deps == nil {
		return
	}
	b.WriteString("  <dependencies>\n")
	if len(deps.Imports) > 0 {
		b.WriteString("    <imports>\n")
		for file, imports := range deps.Imports {
			b.WriteString(fmt.Sprintf("      <file path=\"%s\">\n", file))
			for _, imp := range imports {
				b.WriteString(fmt.Sprintf("        <import>%s</import>\n", imp))
			}
			b.WriteString("      </file>\n")
		}
		b.WriteString("    </imports>\n")
	}
	if len(deps.CoreFiles) > 0 {
		b.WriteString("    <coreFiles>\n")
		for _, file := range deps.CoreFiles {
			b.WriteString(fmt.Sprintf("      <file>%s</file>\n", file))
		}
		b.WriteString("    </coreFiles>\n")
	}
	b.WriteString("  </dependencies>\n")
}

func (x *XMLFormatter) formatFiles(b *strings.Builder, files []FileInfo) {
	if len(files) == 0 {
		return
	}
	b.WriteString("  <files>\n")
	for _, file := range files {
		lineCount := strings.Count(file.Content, "\n") + 1
		b.WriteString(fmt.Sprintf("    <file path=\"%s\" lines=\"%d\">\n", file.Path, lineCount))
		b.WriteString("      <content><![CDATA[")
		b.WriteString(file.Content)
		b.WriteString("]]></content>\n")
		b.WriteString("    </file>\n")
	}
	b.WriteString("  </files>\n")
}

func (x *XMLFormatter) Format(project *ProjectOutput) (string, error) {
	var b strings.Builder
	enc := xml.NewEncoder(&b)
	enc.Indent("", "  ")

	b.WriteString(xml.Header)
	b.WriteString("<project>\n")

	x.formatOverview(&b, project.Overview)
	x.formatFileStats(&b, project.FileStats)

	// Directory Tree
	b.WriteString("  <directoryTree>\n")
	writeDirectoryNode(project.DirectoryTree, &b, 4)
	b.WriteString("  </directoryTree>\n")

	x.formatGitInfo(&b, project.GitInfo)
	x.formatDependencies(&b, project.Dependencies)
	x.formatFiles(&b, project.Files)

	b.WriteString("</project>")
	return b.String(), nil
}
