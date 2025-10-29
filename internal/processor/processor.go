package processor

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/1broseidon/promptext/internal/config"
	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/filter/rules"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/info"
	"github.com/1broseidon/promptext/internal/log"
	"github.com/1broseidon/promptext/internal/relevance"
	"github.com/1broseidon/promptext/internal/token"
	"github.com/atotto/clipboard"
	"github.com/jedib0t/go-pretty/v6/text"
)

type Config struct {
	DirPath           string
	Extensions        []string
	Excludes          []string
	GitIgnore         bool
	Filter            *filter.Filter
	RelevanceKeywords string // Keywords for relevance filtering
	MaxTokens         int    // Maximum token budget (0 = unlimited)
	ExplainSelection  bool   // Show priority scoring breakdown
}

func ParseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}

// ExcludedFileInfo contains information about an excluded file
type ExcludedFileInfo struct {
	Path   string
	Tokens int
}

// FilePriorityInfo contains information about a file's priority for explain-selection
type FilePriorityInfo struct {
	Path      string
	Tokens    int
	Score     float64
	IsEntry   bool
	IsTest    bool
	IsConfig  bool
	Depth     int
	Included  bool
}

// ProcessResult contains both display and clipboard content
type ProcessResult struct {
	ProjectOutput    *format.ProjectOutput
	DisplayContent   string
	ClipboardContent string
	TokenCount       int                 // Token count for included files
	TotalTokens      int                 // Total tokens if all files were included
	ProjectInfo      *info.ProjectInfo
	ExcludedFiles    int                 // Number of files excluded due to token budget
	ExcludedFileList []ExcludedFileInfo  // Details of excluded files
	PriorityList     []FilePriorityInfo  // Priority breakdown for explain-selection
}

// DryRunResult contains dry-run preview information
type DryRunResult struct {
	FilePaths       []string
	EstimatedTokens int
	ConfigSummary   *ConfigSummary
	ProjectInfo     *info.ProjectInfo
}

// ConfigSummary contains effective configuration information
type ConfigSummary struct {
	Extensions      []string
	Excludes        []string
	UseGitIgnore    bool
	UseDefaultRules bool
	Format          string
	OutputFile      string
}

// validateFilePath validates and gets the relative path for a file
func validateFilePath(path string, config Config) (string, error) {
	rel, err := filepath.Rel(config.DirPath, path)
	if err != nil {
		return "", fmt.Errorf("error getting relative path for %s: %w", path, err)
	}

	if !config.Filter.ShouldProcess(rel) {
		return "", nil
	}

	// Skip .DS_Store files immediately
	if filepath.Base(path) == ".DS_Store" {
		return "", nil
	}

	return rel, nil
}

// checkFilePermissions validates file type and permissions
func checkFilePermissions(path string) error {
	// Get file info first to check if it's a directory or has read permissions
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	// Skip directories
	if fileInfo.IsDir() {
		return fmt.Errorf("is directory")
	}

	// Check read permissions
	if fileInfo.Mode().Perm()&0444 == 0 {
		return fmt.Errorf("no read permissions")
	}

	// Check if file is binary using BinaryRule
	binaryRule := rules.NewBinaryRule()
	if binaryRule.Match(path) {
		return fmt.Errorf("binary file")
	}

	return nil
}

// readFileContent reads and returns file content as string
func readFileContent(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// processFile handles the processing of a single file
func processFile(path string, config Config) (*format.FileInfo, error) {
	rel, err := validateFilePath(path, config)
	if err != nil {
		return nil, err
	}
	if rel == "" {
		return nil, nil // File should be skipped
	}

	if err := checkFilePermissions(path); err != nil {
		return nil, nil // File should be skipped
	}

	content, err := readFileContent(path)
	if err != nil {
		return nil, nil // File should be skipped
	}

	return &format.FileInfo{
		Path:    rel,
		Content: content,
	}, nil
}

// populateProjectInfo adds project information to the output
func populateProjectInfo(projectOutput *format.ProjectOutput, projectInfo *info.ProjectInfo) {
	projectOutput.DirectoryTree = projectInfo.DirectoryTree

	if projectInfo.GitInfo != nil {
		projectOutput.GitInfo = &format.GitInfo{
			Branch:        projectInfo.GitInfo.Branch,
			CommitHash:    projectInfo.GitInfo.CommitHash,
			CommitMessage: projectInfo.GitInfo.CommitMessage,
		}
	}

	if projectInfo.Metadata != nil {
		projectOutput.Metadata = &format.Metadata{
			Language:     projectInfo.Metadata.Language,
			Version:      projectInfo.Metadata.Version,
			Dependencies: projectInfo.Metadata.Dependencies,
		}
	}
}

// PreviewDirectory performs a dry-run preview of what files would be processed
func PreviewDirectory(config Config) (*DryRunResult, error) {
	log.StartTimer("Dry Run Preview")
	defer log.EndTimer("Dry Run Preview")

	// Initialize result
	result := &DryRunResult{
		FilePaths: []string{},
		ConfigSummary: &ConfigSummary{
			Extensions:      config.Extensions,
			Excludes:        config.Excludes,
			UseGitIgnore:    config.GitIgnore,
			UseDefaultRules: config.Filter != nil,
		},
	}

	// Collect files that would be processed
	tokenCounter := token.NewTokenCounter()
	var estimatedTokens int

	log.Debug("=== Dry Run: Analyzing Files ===")

	err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Get relative path for filtering
		relPath, err := filepath.Rel(config.DirPath, path)
		if err != nil {
			return err
		}

		// For directories
		if d.IsDir() {
			if config.Filter.IsExcluded(relPath) {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip excluded files silently
		if config.Filter.IsExcluded(relPath) {
			return nil
		}

		// Check if file would pass validation and filtering
		rel, err := validateFilePath(path, config)
		if err != nil {
			return nil // Skip files that would fail validation
		}
		if rel == "" {
			return nil // File should be skipped due to filtering
		}

		// Check permissions and file type without reading content
		if err := checkFilePermissions(path); err != nil {
			return nil // Skip files that would fail permission check
		}

		// Add to result
		result.FilePaths = append(result.FilePaths, rel)

		// Estimate tokens based on file size (rough approximation: 4 chars per token)
		if fileInfo, err := os.Stat(path); err == nil {
			estimatedFileTokens := int(fileInfo.Size() / 4)
			estimatedTokens += estimatedFileTokens
			log.Debug("Would process: %s (estimated %d tokens)", rel, estimatedFileTokens)
		}

		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error during dry-run preview: %w", err)
	}

	result.EstimatedTokens = estimatedTokens

	// Get project info for dry-run
	if projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter); err == nil {
		result.ProjectInfo = projectInfo

		// Add estimated tokens for metadata (rough approximation)
		if projectInfo.DirectoryTree != nil {
			directoryTreeString := projectInfo.DirectoryTree.ToMarkdown(0)
			result.EstimatedTokens += tokenCounter.EstimateTokens(directoryTreeString)
		}
	}

	log.Debug("Dry run complete: %d files, ~%d tokens", len(result.FilePaths), result.EstimatedTokens)
	return result, nil
}

// processFileInWalk handles individual file processing during directory walk
func processFileInWalk(path string, d fs.DirEntry, config Config, tokenCounter *token.TokenCounter, processedFiles *[]format.FileInfo, totalTokens *int, verbose bool) error {
	if d.IsDir() {
		// Get relative path for filtering
		relPath, err := filepath.Rel(config.DirPath, path)
		if err != nil {
			return err
		}
		if config.Filter.IsExcluded(relPath) {
			return filepath.SkipDir
		}
		return nil
	}

	// Get relative path for filtering
	relPath, err := filepath.Rel(config.DirPath, path)
	if err != nil {
		return err
	}

	// Skip excluded files silently
	if config.Filter.IsExcluded(relPath) {
		return nil
	}

	// Process file
	fileInfo, err := processFile(path, config)
	if err != nil {
		log.Debug("Error processing file %s: %v", path, err)
		return nil // Continue processing other files
	}

	if fileInfo != nil {
		*processedFiles = append(*processedFiles, *fileInfo)

		// Count tokens and log immediately
		fileTokens := tokenCounter.EstimateTokens(fileInfo.Content)
		*totalTokens += fileTokens
		log.Debug("Processing: %s (%d tokens)", relPath, fileTokens)

		if verbose && !log.IsDebugEnabled() {
			fmt.Printf("\n### File: %s\n```\n%s\n```\n", path, fileInfo.Content)
		}
	}

	return nil
}

// filterDirectoryTree removes files from the tree that aren't in the included set
func filterDirectoryTree(node *format.DirectoryNode, includedFiles map[string]bool, currentPath string) *format.DirectoryNode {
	if node == nil {
		return nil
	}

	// Create a filtered node
	filtered := &format.DirectoryNode{
		Name:     node.Name,
		Type:     node.Type,
		Children: []*format.DirectoryNode{},
	}

	// Process children
	for _, child := range node.Children {
		childPath := currentPath
		if childPath != "" {
			childPath = filepath.Join(childPath, child.Name)
		} else {
			childPath = child.Name
		}

		if child.Type == "file" {
			// Include file only if it's in the included set
			if includedFiles[childPath] {
				filtered.Children = append(filtered.Children, child)
			}
		} else if child.Type == "dir" {
			// Recursively filter subdirectory
			filteredChild := filterDirectoryTree(child, includedFiles, childPath)
			// Only include directory if it has children after filtering
			if filteredChild != nil && len(filteredChild.Children) > 0 {
				filtered.Children = append(filtered.Children, filteredChild)
			}
		}
	}

	return filtered
}

// filePriority calculates priority score for sorting files
// Higher scores should be processed first
type filePriority struct {
	file     format.FileInfo
	score    float64
	isEntry  bool
	isTest   bool
	isConfig bool
	depth    int
}

// prioritizeFiles sorts files by priority based on relevance and file characteristics
func prioritizeFiles(files []format.FileInfo, scorer *relevance.Scorer, entryPoints map[string]bool) []format.FileInfo {
	if len(files) == 0 {
		return files
	}

	// Build priority list
	priorities := make([]filePriority, len(files))
	for i, file := range files {
		// Calculate path depth
		depth := strings.Count(file.Path, string(filepath.Separator))

		// Check file characteristics
		isEntry := entryPoints[file.Path]
		isTest := strings.Contains(file.Path, "test") || strings.HasSuffix(file.Path, "_test.go")
		isConfig := strings.Contains(strings.ToLower(filepath.Base(file.Path)), "config") ||
		            strings.HasSuffix(file.Path, ".yml") || strings.HasSuffix(file.Path, ".yaml") ||
		            strings.HasSuffix(file.Path, ".json") || strings.HasSuffix(file.Path, ".toml")

		// Calculate relevance score
		relevanceScore := scorer.ScoreFile(file.Path, file.Content)

		priorities[i] = filePriority{
			file:     file,
			score:    relevanceScore,
			isEntry:  isEntry,
			isTest:   isTest,
			isConfig: isConfig,
			depth:    depth,
		}
	}

	// Sort by priority (higher first)
	sort.Slice(priorities, func(i, j int) bool {
		pi, pj := priorities[i], priorities[j]

		// 1. Entry points with high relevance come first
		if pi.isEntry != pj.isEntry {
			return pi.isEntry
		}

		// 2. High relevance scores (above threshold)
		threshold := relevance.GetRelevanceThreshold()
		piHighRelevance := pi.score >= threshold
		pjHighRelevance := pj.score >= threshold
		if piHighRelevance != pjHighRelevance {
			return piHighRelevance
		}

		// 3. Within same relevance tier, prefer shallower files
		if piHighRelevance && pjHighRelevance {
			if pi.depth != pj.depth {
				return pi.depth < pj.depth
			}
			// If same depth, higher score wins
			if pi.score != pj.score {
				return pi.score > pj.score
			}
		}

		// 4. Config files before other non-relevant files
		if !piHighRelevance && !pjHighRelevance {
			if pi.isConfig != pj.isConfig {
				return pi.isConfig
			}
		}

		// 5. Tests come last
		if pi.isTest != pj.isTest {
			return !pi.isTest
		}

		// 6. Finally, prefer shallower paths
		if pi.depth != pj.depth {
			return pi.depth < pj.depth
		}

		// 7. Tie-breaker: alphabetical
		return pi.file.Path < pj.file.Path
	})

	// Extract sorted files
	sorted := make([]format.FileInfo, len(priorities))
	for i, p := range priorities {
		sorted[i] = p.file
	}

	return sorted
}

func ProcessDirectory(config Config, verbose bool) (*ProcessResult, error) {
	log.StartTimer("Project Processing")
	defer log.EndTimer("Project Processing")

	// Initialize project output
	projectOutput := &format.ProjectOutput{}

	// Combined file processing and token analysis
	log.StartTimer("Processing Files")
	tokenCounter := token.NewTokenCounter()
	log.Debug("=== Processing Files & Counting Tokens ===")
	var totalTokens int

	// Process all files first
	var processedFiles []format.FileInfo
	err := filepath.WalkDir(config.DirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		return processFileInWalk(path, d, config, tokenCounter, &processedFiles, &totalTokens, verbose)
	})

	if err != nil {
		return nil, fmt.Errorf("error processing files: %w", err)
	}
	log.EndTimer("Processing Files")

	// Get project info early for entry point detection
	log.StartTimer("Project Analysis")
	projectInfo, err := info.GetProjectInfo(config.DirPath, config.Filter)
	if err != nil {
		return &ProcessResult{}, fmt.Errorf("error getting project info: %w", err)
	}
	log.EndTimer("Project Analysis")

	// Apply relevance scoring and prioritization if keywords provided
	var excludedFileCount int
	var excludedFileList []ExcludedFileInfo
	scorer := relevance.NewScorer(config.RelevanceKeywords)
	if scorer.HasKeywords() || config.MaxTokens > 0 {
		log.Debug("=== Applying Relevance & Token Budget ===")

		// Build entry points map (detect common entry point patterns)
		entryPoints := make(map[string]bool)
		for _, file := range processedFiles {
			basename := filepath.Base(file.Path)
			// Detect common entry point file names
			if basename == "main.go" || basename == "index.js" || basename == "index.ts" ||
			   basename == "app.js" || basename == "app.ts" || basename == "main.py" ||
			   basename == "__init__.py" || basename == "index.html" {
				entryPoints[file.Path] = true
			}
		}

		// Prioritize files
		processedFiles = prioritizeFiles(processedFiles, scorer, entryPoints)
		log.Debug("Files sorted by priority")

		// Apply token budget if specified
		if config.MaxTokens > 0 {
			// Calculate overhead tokens (tree, git, metadata)
			overheadTokens := 0
			formatter, _ := format.GetFormatter("markdown")
			if formatter != nil {
				// Temporarily populate projectOutput for overhead calculation
				tempOutput := &format.ProjectOutput{}
				populateProjectInfo(tempOutput, projectInfo)

				if treeOut, err := formatter.Format(&format.ProjectOutput{DirectoryTree: tempOutput.DirectoryTree}); err == nil {
					overheadTokens += tokenCounter.EstimateTokens(treeOut)
				}
				if gitOut, err := formatter.Format(&format.ProjectOutput{GitInfo: tempOutput.GitInfo}); err == nil {
					overheadTokens += tokenCounter.EstimateTokens(gitOut)
				}
				if metaOut, err := formatter.Format(&format.ProjectOutput{Metadata: tempOutput.Metadata}); err == nil {
					overheadTokens += tokenCounter.EstimateTokens(metaOut)
				}
			}

			availableTokens := config.MaxTokens - overheadTokens
			log.Debug("Token budget: %d (available for files: %d)", config.MaxTokens, availableTokens)

			// Include files until budget is reached
			var filteredFiles []format.FileInfo
			cumulativeTokens := 0

			for _, file := range processedFiles {
				fileTokens := tokenCounter.EstimateTokens(file.Content)
				if cumulativeTokens+fileTokens <= availableTokens {
					filteredFiles = append(filteredFiles, file)
					cumulativeTokens += fileTokens
					log.Debug("Including: %s (%d tokens, cumulative: %d)", file.Path, fileTokens, cumulativeTokens)
				} else {
					excludedFileCount++
					excludedFileList = append(excludedFileList, ExcludedFileInfo{
						Path:   file.Path,
						Tokens: fileTokens,
					})
					log.Debug("Excluding: %s (%d tokens would exceed budget)", file.Path, fileTokens)
				}
			}

			processedFiles = filteredFiles
			log.Debug("Included %d files, excluded %d files due to token budget", len(processedFiles), excludedFileCount)
		}
	}

	// Store processed files
	projectOutput.Files = processedFiles

	// Calculate file statistics
	totalLines := 0
	packages := make(map[string]bool)

	for _, file := range processedFiles {
		totalLines += strings.Count(file.Content, "\n") + 1

		// Extract package directory for Go projects
		dir := filepath.Dir(file.Path)
		if dir != "." && dir != "" {
			packages[dir] = true
		}
	}

	projectOutput.FileStats = &format.FileStatistics{
		TotalFiles:   len(processedFiles),
		TotalLines:   totalLines,
		PackageCount: len(packages),
	}

	// Populate project information (projectInfo already retrieved earlier)
	populateProjectInfo(projectOutput, projectInfo)

	// Filter directory tree if files were excluded due to token budget or relevance
	if excludedFileCount > 0 || scorer.HasKeywords() {
		// Build set of included file paths
		includedFiles := make(map[string]bool)
		for _, file := range processedFiles {
			includedFiles[file.Path] = true
		}

		// Filter the directory tree to only show included files
		if projectOutput.DirectoryTree != nil {
			projectOutput.DirectoryTree = filterDirectoryTree(projectOutput.DirectoryTree, includedFiles, "")
		}
		log.Debug("Filtered directory tree to show only %d included files", len(processedFiles))
	}

	// Get formatter for output
	formatter, err := format.GetFormatter("markdown") // Default to markdown for token counting
	if err != nil {
		return nil, fmt.Errorf("error creating formatter: %w", err)
	}

	// Count tokens for directory tree
	treeOutput, _ := formatter.Format(&format.ProjectOutput{DirectoryTree: projectOutput.DirectoryTree})
	treeTokens := tokenCounter.EstimateTokens(treeOutput)
	totalTokens += treeTokens
	log.Debug("Directory structure: %d tokens", treeTokens)

	// Count tokens for git info
	gitTokens := 0
	if projectOutput.GitInfo != nil {
		gitOutput, _ := formatter.Format(&format.ProjectOutput{GitInfo: projectOutput.GitInfo})
		gitTokens = tokenCounter.EstimateTokens(gitOutput)
		totalTokens += gitTokens
	}
	log.Debug("Git information: %d tokens", gitTokens)

	// Count tokens for metadata
	metaTokens := 0
	if projectOutput.Metadata != nil {
		metaOutput, _ := formatter.Format(&format.ProjectOutput{Metadata: projectOutput.Metadata})
		metaTokens = tokenCounter.EstimateTokens(metaOutput)
		totalTokens += metaTokens
	}
	log.Debug("Project metadata: %d tokens", metaTokens)

	// Calculate source file tokens
	sourceTokens := totalTokens - treeTokens - gitTokens - metaTokens
	log.Debug("Source files: %d tokens", sourceTokens)
	log.Debug("Total token count: %d", totalTokens)

	// Add timing summary
	log.Debug("=== Performance ===")
	log.Debug("Total processing time: %.2fms", float64(time.Since(log.GetPhaseStart()).Microseconds())/1000.0)

	// Format the full output
	formattedOutput, err := formatter.Format(projectOutput)
	if err != nil {
		return nil, fmt.Errorf("error formatting output: %w", err)
	}

	var displayContent string
	if verbose {
		displayContent = formattedOutput
	}

	// Calculate total tokens (included + excluded)
	totalProjectTokens := totalTokens
	for _, excluded := range excludedFileList {
		totalProjectTokens += excluded.Tokens
	}

	return &ProcessResult{
		ProjectOutput:    projectOutput,
		DisplayContent:   displayContent,
		ClipboardContent: formattedOutput,
		TokenCount:       totalTokens,
		TotalTokens:      totalProjectTokens,
		ProjectInfo:      projectInfo,
		ExcludedFiles:    excludedFileCount,
		ExcludedFileList: excludedFileList,
	}, nil
}

// buildProjectHeader constructs the project name and basic info
func buildProjectHeader(config Config, result *ProcessResult, infoOnly bool) string {
	var content strings.Builder

	// Project name and language (always shown)
	if result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Name != "" {
		content.WriteString("📦 " + result.ProjectInfo.Metadata.Name)
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			content.WriteString("📦 " + filepath.Base(absPath))
		}
	}

	if result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Language != "" {
		content.WriteString(fmt.Sprintf(" (%s", result.ProjectInfo.Metadata.Language))
		if result.ProjectInfo.Metadata.Version != "" && infoOnly {
			content.WriteString(fmt.Sprintf(" %s", result.ProjectInfo.Metadata.Version))
		}
		content.WriteString(")")
	}

	// Basic file and token count (always shown)
	fileCount := len(result.ProjectOutput.Files)
	totalFileCount := fileCount + result.ExcludedFiles

	if result.ExcludedFiles > 0 {
		// Show "included / total" format when files were excluded
		content.WriteString(fmt.Sprintf("\n   Included: %d/%d files • ~%s tokens",
			fileCount, totalFileCount, formatTokenCount(result.TokenCount)))
		if result.TotalTokens > result.TokenCount {
			content.WriteString(fmt.Sprintf("\n   Full project: %d files • ~%s tokens",
				totalFileCount, formatTokenCount(result.TotalTokens)))
		}
	} else {
		// Normal display when no files were excluded
		content.WriteString(fmt.Sprintf("\n   Files: %d", fileCount))
		if result.TokenCount > 0 {
			content.WriteString(fmt.Sprintf(" • Tokens: ~%s", formatTokenCount(result.TokenCount)))
		}
	}

	return content.String()
}

// formatTokenCount formats token count with comma separators for readability
func formatTokenCount(tokens int) string {
	if tokens < 1000 {
		return fmt.Sprintf("%d", tokens)
	}
	// Add comma separators for thousands
	str := fmt.Sprintf("%d", tokens)
	var result strings.Builder
	for i, c := range str {
		if i > 0 && (len(str)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}

// analyzeFileStatistics collects file type and size statistics
func analyzeFileStatistics(files []format.FileInfo, config Config) (map[string]int, int64, []string) {
	fileTypes := make(map[string]int)
	var totalSize int64
	var entryPoints []string

	for _, file := range files {
		typeInfo := filter.GetFileType(file.Path, config.Filter)
		fileTypes[typeInfo.Type]++
		totalSize += typeInfo.Size

		// Track entry points
		if typeInfo.IsEntryPoint {
			entryPoints = append(entryPoints, file.Path)
		}
	}

	return fileTypes, totalSize, entryPoints
}

// buildFileAnalysis creates the file analysis section
func buildFileAnalysis(fileTypes map[string]int, totalSize int64, entryPoints []string) string {
	var content strings.Builder

	// Display File Distribution
	content.WriteString("\n   Types: ")
	first := true
	for typ, count := range fileTypes {
		if !first {
			content.WriteString(" • ")
		}
		content.WriteString(fmt.Sprintf("%s: %d", typ, count))
		first = false
	}
	content.WriteString("\n")

	// Display Size Information
	if totalSize > 0 {
		content.WriteString(fmt.Sprintf("   Total Size: %s\n", formatSize(totalSize)))
	}

	// Display Entry Points
	if len(entryPoints) > 0 {
		content.WriteString("\n🚪 Entry Points\n")
		for _, entry := range entryPoints {
			content.WriteString(fmt.Sprintf("   • %s\n", entry))
		}
	}

	return content.String()
}

// buildDependenciesSection creates the dependencies section
func buildDependenciesSection(result *ProcessResult) string {
	if result.ProjectOutput.Metadata == nil || len(result.ProjectOutput.Metadata.Dependencies) == 0 {
		return ""
	}

	var content strings.Builder
	content.WriteString("\n📚 Dependencies\n")
	for _, dep := range result.ProjectOutput.Metadata.Dependencies {
		content.WriteString(fmt.Sprintf("   • %s\n", dep))
	}
	return content.String()
}

// buildHealthSection creates the project health section
func buildHealthSection(result *ProcessResult) string {
	if result.ProjectInfo.Metadata == nil || result.ProjectInfo.Metadata.Health == nil {
		return ""
	}

	health := result.ProjectInfo.Metadata.Health
	var content strings.Builder
	content.WriteString("\n🏥 Project Health\n")

	// Documentation
	content.WriteString(fmt.Sprintf("   • README: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasReadme]))
	content.WriteString(fmt.Sprintf("   • LICENSE: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasLicense]))

	// Testing
	content.WriteString(fmt.Sprintf("   • Tests: %s\n", map[bool]string{true: "✓", false: "✗"}[health.HasTests]))

	// CI/CD
	if health.HasCI {
		content.WriteString(fmt.Sprintf("   • CI/CD: ✓ (%s)\n", health.CISystem))
	} else {
		content.WriteString("   • CI/CD: ✗\n")
	}

	return content.String()
}

// buildGitSection creates the git information section
func buildGitSection(result *ProcessResult) string {
	if result.ProjectOutput.GitInfo == nil {
		return ""
	}

	var content strings.Builder
	content.WriteString("\n🔄 Git Status\n")
	content.WriteString(fmt.Sprintf("   Branch: %s\n", result.ProjectOutput.GitInfo.Branch))

	if result.ProjectOutput.GitInfo.CommitHash != "" {
		shortHash := result.ProjectOutput.GitInfo.CommitHash
		if len(shortHash) > 7 {
			shortHash = shortHash[:7]
		}
		content.WriteString(fmt.Sprintf("   Latest: %s", shortHash))
		if result.ProjectOutput.GitInfo.CommitMessage != "" {
			content.WriteString(fmt.Sprintf(" - %s", result.ProjectOutput.GitInfo.CommitMessage))
		}
		content.WriteString("\n")
	}

	return content.String()
}

// formatBoxedOutput creates a boxed output with borders
func formatBoxedOutput(content string) string {
	contentLines := strings.Split(strings.TrimRight(content, "\n"), "\n")
	maxWidth := 0
	for _, line := range contentLines {
		width := text.RuneCount(line)
		if width > maxWidth {
			maxWidth = width
		}
	}

	// Add padding to max width
	maxWidth += 4 // 2 spaces on each side

	var summary strings.Builder
	summary.WriteString("\033[32m") // Start green color

	// Top border
	summary.WriteString("╭" + strings.Repeat("─", maxWidth) + "╮\n")

	// Content lines
	for _, line := range contentLines {
		// Calculate padding needed
		lineWidth := text.RuneCount(line)
		padding := maxWidth - lineWidth

		// Write line with padding
		summary.WriteString("│ " + line + strings.Repeat(" ", padding-2) + " │\n")
	}

	// Bottom border
	summary.WriteString("╰" + strings.Repeat("─", maxWidth) + "╯")
	summary.WriteString("\033[0m") // Reset color

	return summary.String()
}

// GetMetadataSummary returns a summary of project metadata and analysis
// If infoOnly is true, returns a rich summary with all details
// Otherwise returns a minimal summary with basic project info
func GetMetadataSummary(config Config, result *ProcessResult, infoOnly bool) (string, error) {
	var content strings.Builder

	// Build basic project header
	content.WriteString(buildProjectHeader(config, result, infoOnly))

	// Only show detailed analysis if infoOnly is true
	if infoOnly {
		content.WriteString("\n")

		// Analyze file statistics
		fileTypes, totalSize, entryPoints := analyzeFileStatistics(result.ProjectOutput.Files, config)

		// Build file analysis section
		content.WriteString(buildFileAnalysis(fileTypes, totalSize, entryPoints))

		// Build dependencies section
		content.WriteString(buildDependenciesSection(result))

		// Build health section
		content.WriteString(buildHealthSection(result))

		// Build git section
		content.WriteString(buildGitSection(result))
	}

	return formatBoxedOutput(content.String()), nil
}

// formatSize converts bytes to human readable string
func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

// FormatDryRunOutput formats the dry-run preview results for display
func FormatDryRunOutput(result *DryRunResult, config Config) string {
	var content strings.Builder

	// Project header
	if result.ProjectInfo != nil && result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Name != "" {
		content.WriteString("📦 " + result.ProjectInfo.Metadata.Name)
	} else {
		if absPath, err := filepath.Abs(config.DirPath); err == nil {
			content.WriteString("📦 " + filepath.Base(absPath))
		}
	}

	if result.ProjectInfo != nil && result.ProjectInfo.Metadata != nil && result.ProjectInfo.Metadata.Language != "" {
		content.WriteString(fmt.Sprintf(" (%s)", result.ProjectInfo.Metadata.Language))
	}

	// Dry-run summary
	content.WriteString(fmt.Sprintf("\n🔍 DRY RUN PREVIEW\n"))
	content.WriteString(fmt.Sprintf("   Would process: %d files\n", len(result.FilePaths)))
	content.WriteString(fmt.Sprintf("   Estimated tokens: ~%d\n", result.EstimatedTokens))

	// Configuration summary
	content.WriteString("\n⚙️ Effective Configuration\n")
	if len(result.ConfigSummary.Extensions) > 0 {
		content.WriteString(fmt.Sprintf("   Extensions: %s\n", strings.Join(result.ConfigSummary.Extensions, ", ")))
	} else {
		content.WriteString("   Extensions: all supported types\n")
	}

	if len(result.ConfigSummary.Excludes) > 0 {
		content.WriteString(fmt.Sprintf("   Excludes: %s\n", strings.Join(result.ConfigSummary.Excludes, ", ")))
	}

	content.WriteString(fmt.Sprintf("   Git ignore: %v\n", result.ConfigSummary.UseGitIgnore))
	content.WriteString(fmt.Sprintf("   Default rules: %v\n", result.ConfigSummary.UseDefaultRules))

	// Files to be processed (showing first 10 with more indicator)
	content.WriteString("\n📂 Files to Process\n")
	maxDisplay := 10
	for i, filePath := range result.FilePaths {
		if i >= maxDisplay {
			remaining := len(result.FilePaths) - maxDisplay
			content.WriteString(fmt.Sprintf("   ... and %d more files\n", remaining))
			break
		}
		content.WriteString(fmt.Sprintf("   • %s\n", filePath))
	}

	if len(result.FilePaths) == 0 {
		content.WriteString("   ⚠️  No files would be processed\n")
	}

	return formatBoxedOutput(content.String())
}

// Helper functions to reduce cyclomatic complexity

func setupLogging(debug, quiet bool) {
	if debug {
		log.Enable()
		log.SetColorEnabled(true)
	}
	if quiet {
		log.SetQuiet(true)
	}
}

func loadConfigurations(absPath string) (*config.FileConfig, *config.FileConfig) {
	// Load global configuration
	globalConfig, err := config.LoadGlobalConfig()
	if err != nil {
		log.Info("Warning: Failed to load global config: %v", err)
		globalConfig = &config.FileConfig{}
	}

	// Load project config file from the specified directory
	projectConfig, err := config.LoadConfig(absPath)
	if err != nil {
		log.Info("Warning: Failed to load .promptext.yml from %s: %v", absPath, err)
		projectConfig = &config.FileConfig{}
	}

	return globalConfig, projectConfig
}

func handleDryRun(procConfig Config, outputFormat, outFile string, quiet bool) error {
	dryRunResult, err := PreviewDirectory(procConfig)
	if err != nil {
		return fmt.Errorf("error during dry-run preview: %v", err)
	}

	// Update config summary with additional info
	dryRunResult.ConfigSummary.Format = outputFormat
	dryRunResult.ConfigSummary.OutputFile = outFile

	// Display dry-run results
	preview := FormatDryRunOutput(dryRunResult, procConfig)
	if quiet {
		fmt.Printf("files=%d tokens=%d\n", len(dryRunResult.FilePaths), dryRunResult.EstimatedTokens)
	} else {
		fmt.Printf("\033[32m%s\033[0m\n", preview)
	}
	return nil
}

func handleInfoOnly(procConfig Config, result *ProcessResult, infoOnly, quiet bool) (string, error) {
	info, err := GetMetadataSummary(procConfig, result, infoOnly)
	if err != nil {
		return "", fmt.Errorf("error getting project info: %v", err)
	}

	if infoOnly {
		if quiet {
			fileCount := len(result.ProjectOutput.Files)
			fmt.Printf("files=%d tokens=%d\n", fileCount, result.TokenCount)
		} else {
			fmt.Printf("\033[32m%s\033[0m\n", info)
		}
	}
	return info, nil
}

func handleOutput(formattedOutput, outputFormat, outFile, info string, result *ProcessResult, noCopy, quiet bool) error {
	// Build exclusion message if files were excluded
	exclusionMsg := ""
	if result.ExcludedFiles > 0 {
		if quiet {
			exclusionMsg = fmt.Sprintf(" excluded=%d", result.ExcludedFiles)
		} else {
			// Build detailed exclusion summary
			var summary strings.Builder
			summary.WriteString(fmt.Sprintf("\n⚠️  Excluded %d files due to token budget:\n", result.ExcludedFiles))

			// Show first 5 excluded files with token counts
			displayCount := 5
			if len(result.ExcludedFileList) < displayCount {
				displayCount = len(result.ExcludedFileList)
			}

			totalExcludedTokens := 0
			for i := 0; i < displayCount; i++ {
				excluded := result.ExcludedFileList[i]
				summary.WriteString(fmt.Sprintf("    • %s (~%d tokens)\n", excluded.Path, excluded.Tokens))
				totalExcludedTokens += excluded.Tokens
			}

			// Add summary for remaining files
			if len(result.ExcludedFileList) > displayCount {
				remaining := len(result.ExcludedFileList) - displayCount
				remainingTokens := 0
				for i := displayCount; i < len(result.ExcludedFileList); i++ {
					remainingTokens += result.ExcludedFileList[i].Tokens
				}
				totalExcludedTokens += remainingTokens
				summary.WriteString(fmt.Sprintf("    ... and %d more files (~%d tokens)\n", remaining, remainingTokens))
			}

			summary.WriteString(fmt.Sprintf("    Total excluded: ~%d tokens", totalExcludedTokens))
			exclusionMsg = summary.String()
		}
	}

	if outFile != "" {
		if err := os.WriteFile(outFile, []byte(formattedOutput), 0644); err != nil {
			return fmt.Errorf("error writing to output file: %w", err)
		}
		if quiet {
			fmt.Printf("written=%s format=%s files=%d tokens=%d%s\n", outFile, outputFormat, len(result.ProjectOutput.Files), result.TokenCount, exclusionMsg)
		} else {
			fmt.Printf("\033[32m%s\n✓ code context written to %s (%s format)%s\033[0m\n", info, outFile, outputFormat, exclusionMsg)
		}
	} else if !noCopy {
		if err := clipboard.WriteAll(formattedOutput); err != nil {
			if !quiet {
				log.Info("Warning: Failed to copy to clipboard: %v", err)
			}
			if quiet {
				return fmt.Errorf("clipboard copy failed")
			}
		} else {
			if quiet {
				fmt.Printf("clipboard=ok format=%s files=%d tokens=%d%s\n", outputFormat, len(result.ProjectOutput.Files), result.TokenCount, exclusionMsg)
			} else {
				fmt.Printf("\033[32m%s\n✓ code context copied to clipboard (%s format)%s\033[0m\n", info, outputFormat, exclusionMsg)
			}
		}
	}
	return nil
}

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool, useDefaultRules bool, dryRun bool, quiet bool, relevanceKeywords string, maxTokens int, explainSelection bool) error {
	// Enable debug logging if flag is set
	if debug {
		log.Enable()
		log.SetColorEnabled(true)
	}

	// Set quiet mode
	if quiet {
		log.SetQuiet(true)
	}

	log.Debug("=== Promptext Initialization ===")
	log.Debug("Directory: %s", dirPath)
	// Handle "md" as alias for "markdown"
	if outputFormat == "md" {
		outputFormat = "markdown"
	}

	// Validate format
	formatter, err := format.GetFormatter(outputFormat)
	if err != nil {
		return fmt.Errorf("invalid format (must be markdown or xml): %w", err)
	}
	// Convert dirPath to absolute path
	absPath, err := filepath.Abs(dirPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Load configurations
	globalConfig, projectConfig := loadConfigurations(absPath)

	// Merge global, project, and flag configurations with proper precedence
	extensions, excludes, verboseFlag, _, useGitIgnore, useDefaultRules := config.MergeConfigs(globalConfig, projectConfig, extension, exclude, verbose, debug, &gitignore, &useDefaultRules)
	log.Debug("Configuration:")
	log.Debug("  • Extensions: %v", extensions)
	log.Debug("  • Excludes: %#v", excludes)
	log.Debug("  • Git Ignore: %v", useGitIgnore)

	// Create filter options
	filterOpts := filter.Options{
		Includes:        extensions,
		Excludes:        excludes,
		UseDefaultRules: useDefaultRules,
		UseGitIgnore:    useGitIgnore,
	}

	// Create the filter once and reuse it
	f := filter.New(filterOpts)

	// Create processor configuration with filter
	procConfig := Config{
		DirPath:           absPath,
		Extensions:        extensions,
		Excludes:          excludes,
		GitIgnore:         useGitIgnore,
		Filter:            f,
		RelevanceKeywords: relevanceKeywords,
		MaxTokens:         maxTokens,
	}

	// Handle dry-run mode
	if dryRun {
		return handleDryRun(procConfig, outputFormat, outFile, quiet)
	}

	// Process directory once and reuse results
	result, err := ProcessDirectory(procConfig, verboseFlag)
	if err != nil {
		return fmt.Errorf("error processing directory: %v", err)
	}

	// Handle info-only mode
	info, err := handleInfoOnly(procConfig, result, infoOnly, quiet)
	if err != nil {
		return err
	}
	if infoOnly {
		return nil
	}

	// Format output for the selected format
	formattedOutput, err := formatter.Format(result.ProjectOutput)
	if err != nil {
		return fmt.Errorf("error formatting output: %w", err)
	}

	// Handle output
	return handleOutput(formattedOutput, outputFormat, outFile, info, result, noCopy, quiet)
}
