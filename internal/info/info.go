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
    Language    string
    Version     string
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
        "go.mod",         // Go
        "package.json",   // Node.js
        "requirements.txt", // Python
        "Cargo.toml",     // Rust
        "pom.xml",        // Java (Maven)
        "build.gradle",   // Java (Gradle)
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
