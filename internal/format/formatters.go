package format

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}
type PTXFormatter struct{}        // PTX v1.0 - TOON-based with multiline code
type TOONStrictFormatter struct{} // TOON v1.3 strict compliance

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

// PTXFormatter formats project data in PTX v1.0 format (TOON-based with multiline code)
func (t *PTXFormatter) Format(project *ProjectOutput) (string, error) {
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

		// Create content section with literal file paths as keys
		// File paths will be quoted by TOON encoder (e.g., "internal/config.go")
		// This provides zero ambiguity while maintaining token efficiency
		contents := make(map[string]interface{})
		for _, file := range project.Files {
			// Use literal file path as key (PTX v2.0)
			// TOON encoder will automatically quote paths with special chars
			contents[file.Path] = file.Content
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
func (t *PTXFormatter) treeToDirectoryMap(node *DirectoryNode) map[string]interface{} {
	if node == nil {
		return nil
	}

	structure := make(map[string]interface{})
	t.buildDirectoryMap(node, "", structure)
	return structure
}

// Recursive helper to build directory map
func (t *PTXFormatter) buildDirectoryMap(node *DirectoryNode, currentPath string, structure map[string]interface{}) {
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
func (t *PTXFormatter) mapToList(m map[string]string) []map[string]interface{} {
	var result []map[string]interface{}
	for k, v := range m {
		result = append(result, map[string]interface{}{
			"path": k,
			"desc": v,
		})
	}
	return result
}

// TOONStrictFormatter implements TOON v1.3 strict compliance
// This formatter follows the official TOON specification exactly,
// using escaped strings for code content instead of multiline blocks.
func (t *TOONStrictFormatter) Format(project *ProjectOutput) (string, error) {
	// For now, we'll implement a simplified version that converts
	// code to escaped strings. In production, we'd use gotoon library.

	// Build structured data similar to PTX but with escaped strings
	data := make(map[string]interface{})

	// Project metadata (same as PTX)
	if project.Metadata != nil {
		metadata := make(map[string]interface{})
		metadata["language"] = project.Metadata.Language
		if project.Metadata.Version != "" {
			metadata["version"] = project.Metadata.Version
		}
		if len(project.Metadata.Dependencies) > 0 {
			metadata["dependencies"] = project.Metadata.Dependencies
		}

		// Add project size stats
		if project.FileStats != nil {
			metadata["total_files"] = project.FileStats.TotalFiles
			metadata["total_lines"] = project.FileStats.TotalLines
		}

		data["metadata"] = metadata
	}

	// Git information (same as PTX)
	if project.GitInfo != nil {
		gitInfo := make(map[string]interface{})
		gitInfo["branch"] = project.GitInfo.Branch
		gitInfo["commit"] = project.GitInfo.CommitHash
		if project.GitInfo.CommitMessage != "" {
			// Escape newlines in commit message for TOON v1.3
			gitInfo["message"] = escapeForTOON(project.GitInfo.CommitMessage)
		}
		data["git"] = gitInfo
	}

	// File statistics (same as PTX)
	if project.FileStats != nil {
		stats := make(map[string]interface{})
		stats["totalFiles"] = project.FileStats.TotalFiles
		stats["totalLines"] = project.FileStats.TotalLines
		stats["packages"] = project.FileStats.PackageCount

		if len(project.FileStats.FilesByType) > 0 {
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

	// Directory tree (same as PTX)
	if project.DirectoryTree != nil {
		ptxFormatter := &PTXFormatter{}
		structure := ptxFormatter.treeToDirectoryMap(project.DirectoryTree)
		if len(structure) > 0 {
			data["structure"] = structure
		}
	}

	// Files - create tabular arrays for TOON v1.3 compliance
	if len(project.Files) > 0 {
		// File metadata array (tabular format)
		var fileMetadata []map[string]interface{}
		// Code content array (tabular format with escaped strings)
		var codeContent []map[string]interface{}

		for _, file := range project.Files {
			lineCount := strings.Count(file.Content, "\n") + 1
			ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
			if ext == "" {
				ext = "txt"
			}

			// Add to file metadata (tabular)
			fileMetadata = append(fileMetadata, map[string]interface{}{
				"path":  file.Path,
				"ext":   ext,
				"lines": lineCount,
			})

			// Add to code content (tabular with escaped content)
			codeContent = append(codeContent, map[string]interface{}{
				"path":    file.Path,
				"content": escapeForTOON(file.Content),
			})
		}

		data["files"] = fileMetadata
		data["code"] = codeContent
	}

	// Use custom TOON encoder (we'd use gotoon in production)
	encoder := NewTOONEncoder()
	result, err := encoder.Encode(data)
	if err != nil {
		return "", fmt.Errorf("TOON v1.3 encoding error: %w", err)
	}

	return result, nil
}

// escapeForTOON escapes a string according to TOON v1.3 specification
func escapeForTOON(s string) string {
	// TOON v1.3 requires escaping: \, ", \n, \r, \t
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	s = strings.ReplaceAll(s, "\r", "\\r")
	s = strings.ReplaceAll(s, "\t", "\\t")
	return s
}
