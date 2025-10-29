package format

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}
type TOONFormatter struct{}

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

	}
}

func (m *MarkdownFormatter) Format(project *ProjectOutput) (string, error) {
	var sb strings.Builder

	// Start with language and metadata
	if project.Metadata != nil {
		sb.WriteString(fmt.Sprintf("Language: %s\n", project.Metadata.Language))
		if project.Metadata.Version != "" {
			sb.WriteString(fmt.Sprintf("Version: %s\n", project.Metadata.Version))
		}
		if len(project.Metadata.Dependencies) > 0 {
			sb.WriteString("Dependencies:\n")
			for _, dep := range project.Metadata.Dependencies {
				sb.WriteString(fmt.Sprintf("  - %s\n", dep))
			}
			sb.WriteString("\n")
		}
	}

	// Add directory tree right after metadata
	if project.DirectoryTree != nil {
		sb.WriteString("Project Structure:\n")
		// Skip the root node name but process its children
		for _, child := range project.DirectoryTree.Children {
			sb.WriteString(child.ToMarkdown(1))
		}
		sb.WriteString("\n")
	}

	// Add source files
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

// TOONFormatter formats project data in TOON format for optimal token efficiency
func (t *TOONFormatter) Format(project *ProjectOutput) (string, error) {
	// Build a structured map for TOON encoding
	data := make(map[string]interface{})

	// Project metadata
	if project.Metadata != nil {
		metadata := make(map[string]interface{})
		metadata["language"] = project.Metadata.Language
		if project.Metadata.Version != "" {
			metadata["version"] = project.Metadata.Version
		}
		if len(project.Metadata.Dependencies) > 0 {
			metadata["dependencies"] = project.Metadata.Dependencies
		}

		// Add project size stats for instant intuition
		if project.FileStats != nil {
			metadata["total_files"] = project.FileStats.TotalFiles
			metadata["total_lines"] = project.FileStats.TotalLines
		}

		data["metadata"] = metadata
	}

	// Git information
	if project.GitInfo != nil {
		gitInfo := make(map[string]interface{})
		gitInfo["branch"] = project.GitInfo.Branch
		gitInfo["commit"] = project.GitInfo.CommitHash
		if project.GitInfo.CommitMessage != "" {
			gitInfo["message"] = project.GitInfo.CommitMessage
		}
		data["git"] = gitInfo
	}

	// File statistics
	if project.FileStats != nil {
		stats := make(map[string]interface{})
		stats["totalFiles"] = project.FileStats.TotalFiles
		stats["totalLines"] = project.FileStats.TotalLines
		stats["packages"] = project.FileStats.PackageCount

		if len(project.FileStats.FilesByType) > 0 {
			// Convert map to tabular array for token efficiency
			var fileTypes []map[string]interface{}
			for ext, count := range project.FileStats.FilesByType {
				fileTypes = append(fileTypes, map[string]interface{}{
					"type":  ext,
					"count": count,
				})
			}
			stats["fileTypes"] = fileTypes
		}

		data["stats"] = stats
	}

	// Directory tree (convert to compact map representation)
	if project.DirectoryTree != nil {
		structure := t.treeToDirectoryMap(project.DirectoryTree)
		if len(structure) > 0 {
			data["structure"] = structure
		}
	}

	// Analysis sections
	if project.Analysis != nil {
		analysis := make(map[string]interface{})
		if len(project.Analysis.EntryPoints) > 0 {
			analysis["entryPoints"] = t.mapToList(project.Analysis.EntryPoints)
		}
		if len(project.Analysis.ConfigFiles) > 0 {
			analysis["configFiles"] = t.mapToList(project.Analysis.ConfigFiles)
		}
		if len(project.Analysis.CoreFiles) > 0 {
			analysis["coreFiles"] = t.mapToList(project.Analysis.CoreFiles)
		}
		if len(project.Analysis.TestFiles) > 0 {
			analysis["testFiles"] = t.mapToList(project.Analysis.TestFiles)
		}
		if len(project.Analysis.Documentation) > 0 {
			analysis["documentation"] = t.mapToList(project.Analysis.Documentation)
		}
		if len(analysis) > 0 {
			data["analysis"] = analysis
		}
	}

	// Dependencies
	if project.Dependencies != nil {
		deps := make(map[string]interface{})
		if len(project.Dependencies.Packages) > 0 {
			deps["packages"] = project.Dependencies.Packages
		}
		if len(project.Dependencies.CoreFiles) > 0 {
			deps["coreFiles"] = project.Dependencies.CoreFiles
		}
		if len(deps) > 0 {
			data["dependencies"] = deps
		}
	}

	// Files - use tabular format for metadata only
	if len(project.Files) > 0 {
		// Create tabular array with file metadata
		var fileMetadata []map[string]interface{}
		for _, file := range project.Files {
			lineCount := strings.Count(file.Content, "\n") + 1
			ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
			if ext == "" {
				ext = "txt"
			}

			fileMetadata = append(fileMetadata, map[string]interface{}{
				"path":  file.Path,
				"ext":   ext,
				"lines": lineCount,
			})
		}
		data["files"] = fileMetadata

		// Create content section as nested objects with multiline strings
		// This is more token-efficient than tabular format for large code blocks
		contents := make(map[string]interface{})
		for _, file := range project.Files {
			// Use a sanitized version of the path as the key
			key := strings.ReplaceAll(file.Path, "/", "_")
			key = strings.ReplaceAll(key, ".", "_")
			contents[key] = file.Content
		}
		data["code"] = contents
	}

	// Use TOON encoder
	encoder := NewTOONEncoder()
	result, err := encoder.Encode(data)
	if err != nil {
		return "", fmt.Errorf("TOON encoding error: %w", err)
	}

	return result, nil
}

// Helper function to convert directory tree to map structure
// Returns map[directory_path][]filenames for compact representation
func (t *TOONFormatter) treeToDirectoryMap(node *DirectoryNode) map[string]interface{} {
	if node == nil {
		return nil
	}

	structure := make(map[string]interface{})
	t.buildDirectoryMap(node, "", structure)
	return structure
}

// Recursive helper to build directory map
func (t *TOONFormatter) buildDirectoryMap(node *DirectoryNode, currentPath string, structure map[string]interface{}) {
	if node == nil {
		return
	}

	// Collect files and subdirectories at this level
	var files []string
	var subdirs []*DirectoryNode

	for _, child := range node.Children {
		if child.Type == "file" {
			files = append(files, child.Name)
		} else if child.Type == "dir" {
			subdirs = append(subdirs, child)
		}
	}

	// Add files for this directory if any
	if len(files) > 0 {
		structure[currentPath] = files
	}

	// Recursively process subdirectories
	for _, subdir := range subdirs {
		newPath := currentPath
		if newPath == "" {
			newPath = subdir.Name
		} else {
			newPath = currentPath + "/" + subdir.Name
		}
		t.buildDirectoryMap(subdir, newPath, structure)
	}
}

// Helper to convert map[string]string to list of maps for tabular format
func (t *TOONFormatter) mapToList(m map[string]string) []map[string]interface{} {
	var result []map[string]interface{}
	for k, v := range m {
		result = append(result, map[string]interface{}{
			"path": k,
			"desc": v,
		})
	}
	return result
}
