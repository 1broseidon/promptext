package info

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter"
	"github.com/1broseidon/promptext/internal/format"
	"github.com/1broseidon/promptext/internal/log"
)

// Config holds directory processing configuration
type Config struct {
	Extensions []string
	Excludes   []string
}

// ProjectInfo holds all discoverable information about the project
type ProjectInfo struct {
	DirectoryTree *format.DirectoryNode
	GitInfo       *GitInfo
	Metadata      *ProjectMetadata
}

// GitInfo holds git repository information
type GitInfo struct {
	Branch        string
	CommitHash    string
	CommitMessage string
}

// ProjectMetadata holds project-specific information
type ProjectMetadata struct {
	Name         string
	Language     string
	Version      string
	Dependencies []string
}

// GetProjectInfo gathers all available information about the project
func GetProjectInfo(rootPath string, f *filter.Filter) (*ProjectInfo, error) {
	info := &ProjectInfo{}

	// Get git info if available
	log.StartTimer("Git Info Collection")
	gitInfo, err := getGitInfo(rootPath)
	if err == nil {
		info.GitInfo = gitInfo
	}
	log.EndTimer("Git Info Collection")

	// Try to get project metadata if available
	metadata, err := getProjectMetadata(rootPath)
	if err == nil {
		info.Metadata = metadata
	}

	// Generate directory tree
	tree, err := generateDirectoryTree(rootPath, f)
	if err != nil {
		return nil, fmt.Errorf("error generating directory tree: %w", err)
	}
	info.DirectoryTree = tree

	return info, nil
}

func generateDirectoryTree(root string, f *filter.Filter) (*format.DirectoryNode, error) {
	rootNode := &format.DirectoryNode{
		Name: filepath.Base(root),
		Type: "dir",
	}

	// Map to track directories
	dirMap := make(map[string]*format.DirectoryNode)
	dirMap["."] = rootNode

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}

		// Skip root directory
		if rel == "." {
			return nil
		}

		// Check if path should be excluded
		if f.IsExcluded(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// For files, only include those that pass the filter
		if !d.IsDir() && !f.ShouldProcess(rel) {
			return nil
		}

		// Split path into components
		parts := strings.Split(rel, string(filepath.Separator))
		currentPath := ""
		currentNode := rootNode

		// Create nodes for each part of the path
		for i, part := range parts {
			if currentPath == "" {
				currentPath = part
			} else {
				currentPath = filepath.Join(currentPath, part)
			}

			isLast := i == len(parts)-1
			isDir := d.IsDir() || !isLast

			if isDir {
				if _, exists := dirMap[currentPath]; !exists {
					newNode := &format.DirectoryNode{
						Name: part,
						Type: "dir",
					}
					dirMap[currentPath] = newNode
					currentNode.Children = append(currentNode.Children, newNode)
				}
				currentNode = dirMap[currentPath]
			} else {
				fileNode := &format.DirectoryNode{
					Name: part,
					Type: "file",
				}
				currentNode.Children = append(currentNode.Children, fileNode)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}
	return rootNode, nil
}

func getGitInfo(root string) (*GitInfo, error) {
	// Check if it's a git repository
	if _, err := os.Stat(filepath.Join(root, ".git")); os.IsNotExist(err) {
		return nil, fmt.Errorf("not a git repository")
	}

	info := &GitInfo{}

	// Get current branch
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		info.Branch = strings.TrimSpace(string(out))
	}

	// Get latest commit hash
	cmd = exec.Command("git", "rev-parse", "--short", "HEAD")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		info.CommitHash = strings.TrimSpace(string(out))
	}

	// Get latest commit message
	cmd = exec.Command("git", "log", "-1", "--pretty=%B")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		info.CommitMessage = strings.TrimSpace(string(out))
	}

	return info, nil
}

func getProjectMetadata(root string) (*ProjectMetadata, error) {
	metadata := &ProjectMetadata{}

	// Check for different project files
	files := []string{
		"pyproject.toml",   // Python (Poetry)
		"poetry.lock",      // Python (Poetry)
		"go.mod",           // Go
		"package.json",     // Node.js
		"requirements.txt", // Python
		"Cargo.toml",       // Rust
		"pom.xml",          // Java (Maven)
		"build.gradle",     // Java (Gradle)
	}

	for _, file := range files {
		if info, err := os.Stat(filepath.Join(root, file)); err == nil && !info.IsDir() {
			metadata.Language = detectLanguage(file)
			metadata.Version = getLanguageVersion(root, metadata.Language)
			metadata.Dependencies = getDependencies(root, file)
			break
		}
	}

	if metadata.Language == "" {
		return nil, fmt.Errorf("no recognized project files found")
	}

	return metadata, nil
}

func detectLanguage(filename string) string {
	switch filename {
	case "go.mod":
		return "Go"
	case "package.json":
		return "JavaScript/Node.js"
	case "requirements.txt", "pyproject.toml", "poetry.lock":
		return "Python"
	case "Cargo.toml":
		return "Rust"
	case "pom.xml":
		return "Java (Maven)"
	case "build.gradle":
		return "Java (Gradle)"
	default:
		return ""
	}
}

func getLanguageVersion(root, language string) string {
	switch language {
	case "Go":
		return getGoVersion(root)
	case "JavaScript/Node.js":
		return getNodeVersion(root)
	case "Python":
		return getPythonVersion(root)
	case "Rust":
		return getRustVersion(root)
	case "Java (Maven)", "Java (Gradle)":
		return getJavaVersion(root)
	default:
		return ""
	}
}

func getDependencies(root, filename string) []string {
	switch filename {
	case "go.mod":
		return getGoDependencies(root)
	case "package.json":
		return getNodeDependencies(root)
	case "requirements.txt":
		return getPythonDependencies(root)
	case "Cargo.toml":
		return getRustDependencies(root)
	case "pom.xml":
		return getJavaMavenDependencies(root)
	case "build.gradle":
		return getJavaGradleDependencies(root)
	default:
		return nil
	}
}

func getGoVersion(root string) string {
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "go ") {
			return strings.TrimSpace(strings.TrimPrefix(line, "go"))
		}
	}
	return ""
}

func getNodeVersion(root string) string {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(content), "\n")
	var nodeVersion string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "\"node\"") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				nodeVersion = strings.TrimSpace(parts[3])
				return fmt.Sprintf("requires Node %s", nodeVersion)
			}
		}
	}
	return ""
}

func getPythonVersion(root string) string {
	// Try pyproject.toml
	if content, err := os.ReadFile(filepath.Join(root, "pyproject.toml")); err == nil {
		lines := strings.Split(string(content), "\n")
		inToolPoetry := false
		inDependencies := false
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "[tool.poetry]" {
				inToolPoetry = true
				continue
			}
			if line == "[tool.poetry.dependencies]" {
				inDependencies = true
				continue
			}
			if (inToolPoetry || inDependencies) && strings.HasPrefix(line, "[") {
				inToolPoetry = false
				inDependencies = false
				continue
			}
			if inDependencies && strings.HasPrefix(line, "python = ") {
				version := strings.Trim(strings.TrimPrefix(line, "python = "), "\"'^")
				return version
			}
		}
	}
	return ""
}

func getRustVersion(root string) string {
	content, err := os.ReadFile(filepath.Join(root, "Cargo.toml"))
	if err != nil {
		return ""
	}
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.Contains(line, "version = ") {
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				return strings.Trim(parts[1], "\"'")
			}
		}
	}
	return ""
}

func getJavaVersion(root string) string {
	cmd := exec.Command("java", "--version")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		return strings.Split(strings.TrimSpace(string(out)), "\n")[0]
	}
	return ""
}

func getGoDependencies(root string) []string {
	content, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "require ") {
			parts := strings.Fields(line)
			if len(parts) > 1 {
				deps = append(deps, parts[1])
			}
		} else if strings.Contains(line, " ") && !strings.HasPrefix(line, "go ") && !strings.HasPrefix(line, "module ") {
			parts := strings.Fields(line)
			if len(parts) > 0 {
				deps = append(deps, parts[0])
			}
		}
	}
	return deps
}

func getNodeDependencies(root string) []string {
	content, err := os.ReadFile(filepath.Join(root, "package.json"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	inDeps := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "\"dependencies\"") || strings.Contains(line, "\"devDependencies\"") {
			inDeps = true
			continue
		}
		if inDeps && strings.Contains(line, "}") {
			inDeps = false
			continue
		}
		if inDeps && strings.Contains(line, "\":") {
			dep := strings.Split(line, "\"")[1]
			deps = append(deps, dep)
		}
	}
	return deps
}

// getPythonDependencies returns all Python dependencies from various sources
func getPythonDependencies(root string) []string {
	depsMap := make(map[string]bool)

	// Collect dependencies from each source
	getPipDependencies(root, depsMap)
	getPoetryDependencies(root, depsMap)
	getPoetryLockDependencies(root, depsMap)
	getVenvDependencies(root, depsMap)

	// Convert map to slice
	var allDeps []string
	for dep := range depsMap {
		allDeps = append(allDeps, dep)
	}
	return allDeps
}

// getPipDependencies reads dependencies from requirements.txt
func getPipDependencies(root string, depsMap map[string]bool) {
	content, err := os.ReadFile(filepath.Join(root, "requirements.txt"))
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			dep := strings.Split(line, "==")[0]
			depsMap[dep] = true
		}
	}
}

// getPoetryDependencies reads dependencies from pyproject.toml
func getPoetryDependencies(root string, depsMap map[string]bool) {
	content, err := os.ReadFile(filepath.Join(root, "pyproject.toml"))
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	inMainDeps := false
	inDevDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		switch line {
		case "[tool.poetry.dependencies]":
			inMainDeps = true
			inDevDeps = false
			continue
		case "[tool.poetry.group.dev.dependencies]":
			inMainDeps = false
			inDevDeps = true
			continue
		}

		if (inMainDeps || inDevDeps) && strings.HasPrefix(line, "[") {
			inMainDeps = false
			inDevDeps = false
			continue
		}

		if (inMainDeps || inDevDeps) && line != "" && !strings.HasPrefix(line, "#") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) > 0 {
				dep := strings.TrimSpace(parts[0])
				if dep != "python" { // Skip python version constraint
					if inDevDeps {
						depsMap["[dev] "+dep] = true
					} else {
						depsMap[dep] = true
					}
				}
			}
		}
	}
}

// getPoetryLockDependencies reads dependencies from poetry.lock
func getPoetryLockDependencies(root string, depsMap map[string]bool) {
	content, err := os.ReadFile(filepath.Join(root, "poetry.lock"))
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	inPackage := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if strings.HasPrefix(line, "[[package]]") {
			inPackage = true
			continue
		}

		if inPackage && strings.HasPrefix(line, "name = ") {
			name := strings.Trim(strings.TrimPrefix(line, "name = "), "\"")
			depsMap[name] = true
			inPackage = false
		}
	}
}

// getVenvDependencies reads dependencies from virtual environment
func getVenvDependencies(root string, depsMap map[string]bool) {
	venvDirs := []string{".venv", "venv"}

	for _, venvDir := range venvDirs {
		sitePackages := filepath.Join(root, venvDir, "lib", "python3.*", "site-packages")
		matches, err := filepath.Glob(sitePackages)
		if err != nil || len(matches) == 0 {
			continue
		}

		entries, err := os.ReadDir(matches[0])
		if err != nil {
			continue
		}

		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				depsMap[entry.Name()] = true
			}
		}
	}
}

func getRustDependencies(root string) []string {
	content, err := os.ReadFile(filepath.Join(root, "Cargo.toml"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	inDeps := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "[dependencies]") {
			inDeps = true
			continue
		}
		if inDeps && strings.HasPrefix(line, "[") {
			inDeps = false
			continue
		}
		if inDeps && strings.Contains(line, "=") {
			dep := strings.Split(line, "=")[0]
			deps = append(deps, strings.TrimSpace(dep))
		}
	}
	return deps
}

func getJavaMavenDependencies(root string) []string {
	// This is a simplified version. For a full implementation,
	// you'd want to use an XML parser
	content, err := os.ReadFile(filepath.Join(root, "pom.xml"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	inDep := false
	var currentDep string
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "<dependency>") {
			inDep = true
			currentDep = ""
			continue
		}
		if strings.Contains(line, "</dependency>") {
			if currentDep != "" {
				deps = append(deps, currentDep)
			}
			inDep = false
			continue
		}
		if inDep && strings.Contains(line, "<artifactId>") {
			currentDep = strings.TrimSuffix(strings.TrimPrefix(line, "<artifactId>"), "</artifactId>")
		}
	}
	return deps
}

func getJavaGradleDependencies(root string) []string {
	content, err := os.ReadFile(filepath.Join(root, "build.gradle"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	inDeps := false
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.Contains(line, "dependencies {") {
			inDeps = true
			continue
		}
		if inDeps && strings.Contains(line, "}") {
			inDeps = false
			continue
		}
		if inDeps && strings.Contains(line, "implementation") {
			parts := strings.Split(line, "'")
			if len(parts) > 1 {
				deps = append(deps, parts[1])
			}
		}
	}
	return deps
}

// ProjectAnalysis contains categorized project files and their descriptions
type ProjectAnalysis struct {
	EntryPoints   map[string]string // Entry points by language pattern
	ConfigFiles   map[string]string // Config files with descriptions
	CoreFiles     map[string]string // Core implementation files
	TestFiles     map[string]string // Test files
	Documentation map[string]string // Documentation files
}

// Helper function to compare path slices
func getConfigDescription(path string) string {
	switch filepath.Base(path) {
	case ".promptext.yml":
		return "Tool configuration file"
	case "go.mod":
		return "Go module definition"
	case ".gitignore":
		return "Git ignore patterns"
	default:
		return "Configuration file"
	}
}

func getDocDescription(path string) string {
	base := filepath.Base(path)
	if strings.HasPrefix(strings.ToUpper(base), "README") {
		return "Project documentation"
	}
	if strings.HasPrefix(strings.ToUpper(base), "LICENSE") {
		return "License information"
	}
	return "Documentation"
}

func isCoreFile(path string) bool {
	// Skip node_modules
	if strings.Contains(path, "node_modules/") {
		return false
	}

	// Convert path separators to forward slashes for consistent matching
	normalizedPath := filepath.ToSlash(path)

	// Only consider files in these special directories as core
	return strings.Contains(normalizedPath, "internal/") ||
		strings.Contains(normalizedPath, "pkg/") ||
		strings.Contains(normalizedPath, "src/") ||
		strings.Contains(normalizedPath, "lib/") ||
		strings.Contains(normalizedPath, "core/")
}

func getCoreDescription(_ string) string {
	return "Core implementation"
}

func AnalyzeProject(rootPath string, f *filter.Filter) *ProjectAnalysis {
	analysis := &ProjectAnalysis{
		EntryPoints:   make(map[string]string),
		ConfigFiles:   make(map[string]string),
		CoreFiles:     make(map[string]string),
		TestFiles:     make(map[string]string),
		Documentation: make(map[string]string),
	}

	filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(rootPath, path)

		// Skip excluded paths entirely
		if f.IsExcluded(rel) {
			return nil
		}

		fileType := filter.GetFileType(rel, f)

		switch {
		case strings.HasPrefix(fileType, "entry:"):
			analysis.EntryPoints[rel] = "Project entry point"
		case fileType == "config":
			analysis.ConfigFiles[rel] = getConfigDescription(rel)
		case fileType == "test":
			analysis.TestFiles[rel] = "Test suite"
		case fileType == "doc":
			analysis.Documentation[rel] = getDocDescription(rel)
		case isCoreFile(rel):
			analysis.CoreFiles[rel] = getCoreDescription(rel)
		}

		return nil
	})

	return analysis
}
