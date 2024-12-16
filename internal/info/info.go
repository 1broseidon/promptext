package info

import (
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ProjectInfo holds all discoverable information about the project
type ProjectInfo struct {
	DirectoryTree string
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
	Language     string
	Version      string
	Dependencies []string
}

// GetProjectInfo gathers all available information about the project
func GetProjectInfo(rootPath string) (*ProjectInfo, error) {
	info := &ProjectInfo{}

	// Get directory tree
	tree, err := generateDirectoryTree(rootPath)
	if err != nil {
		return nil, fmt.Errorf("error generating directory tree: %w", err)
	}
	info.DirectoryTree = tree

	// Try to get git info if available
	gitInfo, err := getGitInfo(rootPath)
	if err == nil {
		info.GitInfo = gitInfo
	}

	// Try to get project metadata if available
	metadata, err := getProjectMetadata(rootPath)
	if err == nil {
		info.Metadata = metadata
	}

	return info, nil
}

func generateDirectoryTree(root string) (string, error) {
	var builder strings.Builder
	builder.WriteString("### Project Structure:\n```\n")

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

		// Skip hidden files and common ignore patterns
		if shouldSkip(rel) {
			if d.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		indent := strings.Repeat("  ", strings.Count(rel, string(filepath.Separator)))
		prefix := "├──"
		if isLastItem(path) {
			prefix = "└──"
		}

		if d.IsDir() {
			builder.WriteString(fmt.Sprintf("%s%s %s/\n", indent, prefix, d.Name()))
		} else {
			builder.WriteString(fmt.Sprintf("%s%s %s\n", indent, prefix, d.Name()))
		}
		return nil
	})

	builder.WriteString("```\n")
	return builder.String(), err
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
	case "requirements.txt":
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
	cmd := exec.Command("node", "--version")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func getPythonVersion(root string) string {
	cmd := exec.Command("python", "--version")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(out))
	}
	return ""
}

func getRustVersion(root string) string {
	cmd := exec.Command("rustc", "--version")
	cmd.Dir = root
	if out, err := cmd.Output(); err == nil {
		return strings.TrimSpace(string(out))
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
		if strings.HasPrefix(line, "require ") || (strings.Contains(line, " ") && !strings.HasPrefix(line, "go ")) {
			deps = append(deps, strings.Split(line, " ")[0])
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

func getPythonDependencies(root string) []string {
	content, err := os.ReadFile(filepath.Join(root, "requirements.txt"))
	if err != nil {
		return nil
	}

	var deps []string
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" && !strings.HasPrefix(line, "#") {
			deps = append(deps, strings.Split(line, "==")[0])
		}
	}
	return deps
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

func shouldSkip(path string) bool {
	patterns := []string{
		".git/",
		"node_modules/",
		"__pycache__/",
		".env",
		".DS_Store",
		"*.pyc",
		"*.pyo",
		"*.pyd",
		"*.so",
		"*.dylib",
		"*.dll",
		"*.class",
		"target/",
		"dist/",
		"build/",
	}

	for _, pattern := range patterns {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return true
		}
		if strings.Contains(path, pattern) {
			return true
		}
	}
	return false
}

func isLastItem(path string) bool {
	// This is a simplified version. For a more accurate implementation,
	// we'd need to track parent directories and their remaining items
	return true
}
