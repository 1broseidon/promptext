package format

import (
	"encoding/xml"
	"fmt"
	"path/filepath"
	"sort"
	"strings"
)

type MarkdownFormatter struct{}
type XMLFormatter struct{}
type PTXFormatter struct{}        // PTX v2.0 - TOON-based with multiline code and enhanced manifest
type TOONStrictFormatter struct{} // TOON v1.3 strict compliance
type JSONLFormatter struct{}      // JSONL - Machine-friendly sidecar format (one JSON object per line)

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

// PTXFormatter formats project data in PTX v2.0 format (TOON-based with multiline code and enhanced manifest)
func (t *PTXFormatter) Format(project *ProjectOutput) (string, error) {
	// Build a structured map for TOON encoding
	data := make(map[string]interface{})

	// PTX schema version and manifest
	promptext := make(map[string]interface{})
	promptext["schema"] = "ptx/v2.0"
	data["promptext"] = promptext

	// Project metadata with enhanced fields
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

	// Budget information (token estimation and truncation tracking)
	if project.Budget != nil {
		budget := make(map[string]interface{})
		budget["max_tokens"] = project.Budget.MaxTokens
		budget["est_tokens"] = project.Budget.EstimatedTokens
		if project.Budget.FileTruncations > 0 {
			budget["file_truncations"] = project.Budget.FileTruncations
		}
		data["budget"] = budget
	}

	// Filter configuration used to generate this output
	if project.FilterConfig != nil {
		filters := make(map[string]interface{})
		if len(project.FilterConfig.Includes) > 0 {
			filters["includes"] = project.FilterConfig.Includes
		}
		if len(project.FilterConfig.Excludes) > 0 {
			filters["excludes"] = project.FilterConfig.Excludes
		}
		data["filters"] = filters
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

	// Files - enhanced manifest with per-file metadata including token counts and truncation info
	if len(project.Files) > 0 {
		// Sort files by path for deterministic output (PTX v2.0 requirement)
		sortedFiles := make([]FileInfo, len(project.Files))
		copy(sortedFiles, project.Files)
		sort.Slice(sortedFiles, func(i, j int) bool {
			return sortedFiles[i].Path < sortedFiles[j].Path
		})

		// Create tabular array with comprehensive file metadata
		var fileMetadata []map[string]interface{}
		for _, file := range sortedFiles {
			lineCount := strings.Count(file.Content, "\n") + 1
			ext := strings.TrimPrefix(filepath.Ext(file.Path), ".")
			if ext == "" {
				ext = "txt"
			}

			fileEntry := map[string]interface{}{
				"path":  file.Path,
				"lines": lineCount,
			}

			// Add token count if available
			if file.Tokens > 0 {
				fileEntry["tokens"] = file.Tokens
			}

			// Add truncation info if file was truncated
			if file.Truncation != nil {
				truncInfo := make(map[string]interface{})
				truncInfo["mode"] = file.Truncation.Mode
				truncInfo["original_tokens"] = file.Truncation.OriginalTokens
				fileEntry["truncation"] = truncInfo
			}

			fileMetadata = append(fileMetadata, fileEntry)
		}
		data["files"] = fileMetadata

		// Create content section with literal file paths as keys
		// File paths will be quoted by TOON encoder (e.g., "internal/config.go")
		// This provides zero ambiguity while maintaining token efficiency
		contents := make(map[string]interface{})
		for _, file := range sortedFiles {
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

// JSONLFormatter implements machine-friendly JSONL output (one JSON object per line)
// This format is ideal for programmatic processing, streaming, and pipeline integration
func (j *JSONLFormatter) Format(project *ProjectOutput) (string, error) {
	var sb strings.Builder
	encoder := NewTOONEncoder() // We'll use JSON encoding from TOON encoder

	// Line 1: Metadata header
	metadataLine := make(map[string]interface{})
	metadataLine["type"] = "metadata"
	if project.Metadata != nil {
		metadataLine["language"] = project.Metadata.Language
		if project.Metadata.Version != "" {
			metadataLine["version"] = project.Metadata.Version
		}
		if len(project.Metadata.Dependencies) > 0 {
			metadataLine["dependencies"] = project.Metadata.Dependencies
		}
		if project.FileStats != nil {
			metadataLine["total_files"] = project.FileStats.TotalFiles
			metadataLine["total_lines"] = project.FileStats.TotalLines
		}
	}
	if metadataJSON, err := encoder.encodeToJSON(metadataLine); err == nil {
		sb.WriteString(metadataJSON)
		sb.WriteString("\n")
	}

	// Line 2: Git info
	if project.GitInfo != nil {
		gitLine := map[string]interface{}{
			"type":   "git",
			"branch": project.GitInfo.Branch,
			"commit": project.GitInfo.CommitHash,
		}
		if project.GitInfo.CommitMessage != "" {
			gitLine["message"] = project.GitInfo.CommitMessage
		}
		if gitJSON, err := encoder.encodeToJSON(gitLine); err == nil {
			sb.WriteString(gitJSON)
			sb.WriteString("\n")
		}
	}

	// Line 3: Budget info (if present)
	if project.Budget != nil {
		budgetLine := map[string]interface{}{
			"type":       "budget",
			"max_tokens": project.Budget.MaxTokens,
			"est_tokens": project.Budget.EstimatedTokens,
		}
		if project.Budget.FileTruncations > 0 {
			budgetLine["file_truncations"] = project.Budget.FileTruncations
		}
		if budgetJSON, err := encoder.encodeToJSON(budgetLine); err == nil {
			sb.WriteString(budgetJSON)
			sb.WriteString("\n")
		}
	}

	// Line 4: Filter config (if present)
	if project.FilterConfig != nil {
		filterLine := map[string]interface{}{
			"type": "filters",
		}
		if len(project.FilterConfig.Includes) > 0 {
			filterLine["includes"] = project.FilterConfig.Includes
		}
		if len(project.FilterConfig.Excludes) > 0 {
			filterLine["excludes"] = project.FilterConfig.Excludes
		}
		if filterJSON, err := encoder.encodeToJSON(filterLine); err == nil {
			sb.WriteString(filterJSON)
			sb.WriteString("\n")
		}
	}

	// Sort files by path for deterministic output
	sortedFiles := make([]FileInfo, len(project.Files))
	copy(sortedFiles, project.Files)
	sort.Slice(sortedFiles, func(i, j int) bool {
		return sortedFiles[i].Path < sortedFiles[j].Path
	})

	// Lines N: One line per file with metadata and content
	for _, file := range sortedFiles {
		lineCount := strings.Count(file.Content, "\n") + 1
		fileLine := map[string]interface{}{
			"type":    "file",
			"path":    file.Path,
			"lines":   lineCount,
			"content": file.Content,
		}

		if file.Tokens > 0 {
			fileLine["tokens"] = file.Tokens
		}

		if file.Truncation != nil {
			fileLine["truncation"] = map[string]interface{}{
				"mode":            file.Truncation.Mode,
				"original_tokens": file.Truncation.OriginalTokens,
			}
		}

		if fileJSON, err := encoder.encodeToJSON(fileLine); err == nil {
			sb.WriteString(fileJSON)
			sb.WriteString("\n")
		}
	}

	return sb.String(), nil
}
