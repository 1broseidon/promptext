package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/pkg/promptext"
)

type PRAnalysis struct {
	ChangelogNeeded     bool
	ChangelogSuggestion string
	DocsNeeded          []string
	BreakingChanges     []string
	ExamplesNeeded      []string
	TestSuggestions     []string
}

func main() {
	// Get environment variables
	prNumber := os.Getenv("PR_NUMBER")
	baseSHA := os.Getenv("BASE_SHA")
	headSHA := os.Getenv("HEAD_SHA")

	if prNumber == "" {
		log.Fatal("PR_NUMBER environment variable required")
	}

	fmt.Printf("🤖 PR Assistant analyzing PR #%s\n\n", prNumber)

	// Get changed files in PR
	changedFiles, err := getChangedFilesInPR(baseSHA, headSHA)
	if err != nil {
		log.Fatalf("Failed to get changed files: %v", err)
	}

	if len(changedFiles) == 0 {
		fmt.Println("⚠️  No changes detected in PR")
		return
	}

	fmt.Printf("📁 Analyzing %d changed files...\n", len(changedFiles))

	// Extract context with promptext
	result, err := extractPRContext(changedFiles)
	if err != nil {
		log.Fatalf("Failed to extract context: %v", err)
	}

	fmt.Printf("   Extracted ~%d tokens from %d files\n\n", result.TokenCount, len(result.ProjectOutput.Files))

	// Analyze PR
	analysis := analyzePR(changedFiles, result)

	// Generate comment
	comment := generatePRComment(analysis, result)

	fmt.Println(comment)

	// TODO: Post comment to GitHub (requires GitHub API integration)
	// For now, just output to stdout for testing
	fmt.Println("\n💡 To post this comment to GitHub, integrate with GitHub API")
}

// getChangedFilesInPR returns list of changed files in the PR
func getChangedFilesInPR(baseSHA, headSHA string) ([]string, error) {
	var cmd *exec.Cmd

	if baseSHA != "" && headSHA != "" {
		// Use provided SHAs
		cmd = exec.Command("git", "diff", "--name-only", baseSHA+"..."+headSHA)
	} else {
		// Fallback: compare with main branch
		cmd = exec.Command("git", "diff", "--name-only", "origin/main...HEAD")
	}

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git diff failed: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var files []string
	for _, line := range lines {
		if line != "" {
			files = append(files, line)
		}
	}
	return files, nil
}

// extractPRContext uses promptext to extract context from changed files
func extractPRContext(changedFiles []string) (*promptext.Result, error) {
	// Focus on code files
	relevantExts := []string{".go", ".md", ".yml", ".yaml", ".json"}

	// Filter relevant files
	var relevantFiles []string
	for _, file := range changedFiles {
		ext := filepath.Ext(file)
		for _, relevantExt := range relevantExts {
			if ext == relevantExt {
				relevantFiles = append(relevantFiles, file)
				break
			}
		}
	}

	// Extract with promptext
	return promptext.Extract(".",
		promptext.WithExtensions(relevantExts...),
		promptext.WithTokenBudget(10000),
		promptext.WithFormat(promptext.FormatMarkdown),
	)
}

// analyzePR analyzes the PR changes
func analyzePR(changedFiles []string, result *promptext.Result) PRAnalysis {
	analysis := PRAnalysis{
		DocsNeeded:      []string{},
		BreakingChanges: []string{},
		ExamplesNeeded:  []string{},
		TestSuggestions: []string{},
	}

	// Check for library changes
	hasLibraryChanges := false
	hasInternalChanges := false
	hasExampleChanges := false
	hasDocChanges := false

	for _, file := range changedFiles {
		if strings.HasPrefix(file, "pkg/promptext/") {
			hasLibraryChanges = true
		}
		if strings.HasPrefix(file, "internal/") {
			hasInternalChanges = true
		}
		if strings.HasPrefix(file, "examples/") {
			hasExampleChanges = true
		}
		if strings.HasPrefix(file, "docs-astro/") || strings.HasSuffix(file, ".md") {
			hasDocChanges = true
		}
	}

	// Changelog check
	hasChangelogUpdate := false
	for _, file := range changedFiles {
		if file == "CHANGELOG.md" || strings.Contains(file, "changelog.md") {
			hasChangelogUpdate = true
			break
		}
	}

	if (hasLibraryChanges || hasInternalChanges) && !hasChangelogUpdate {
		analysis.ChangelogNeeded = true
		analysis.ChangelogSuggestion = generateChangelogSuggestion(changedFiles)
	}

	// Documentation check
	if hasLibraryChanges && !hasDocChanges {
		analysis.DocsNeeded = append(analysis.DocsNeeded,
			"Library changes detected - consider updating docs-astro/src/content/docs/library-usage.md")
	}

	// Examples check
	if hasLibraryChanges && !hasExampleChanges {
		analysis.ExamplesNeeded = append(analysis.ExamplesNeeded,
			"New library features - consider adding examples/")
	}

	// Breaking changes detection (simple heuristic)
	for _, file := range result.ProjectOutput.Files {
		content := strings.ToLower(file.Content)
		if strings.Contains(content, "breaking") || strings.Contains(content, "deprecated") {
			analysis.BreakingChanges = append(analysis.BreakingChanges, file.Path)
		}
	}

	return analysis
}

// generateChangelogSuggestion creates a suggested changelog entry
func generateChangelogSuggestion(changedFiles []string) string {
	var suggestion strings.Builder

	suggestion.WriteString("```markdown\n")
	suggestion.WriteString("### Added\n")
	suggestion.WriteString("- [Describe new feature]\n\n")

	// Check file types to suggest category
	hasGoFiles := false
	hasTests := false

	for _, file := range changedFiles {
		if strings.HasSuffix(file, ".go") && !strings.HasSuffix(file, "_test.go") {
			hasGoFiles = true
		}
		if strings.HasSuffix(file, "_test.go") {
			hasTests = true
		}
	}

	if hasGoFiles {
		suggestion.WriteString("### Changed\n")
		suggestion.WriteString("- [Describe changes]\n\n")
	}

	if hasTests {
		suggestion.WriteString("### Testing\n")
		suggestion.WriteString("- Added test coverage for [feature]\n")
	}

	suggestion.WriteString("```")

	return suggestion.String()
}

// generatePRComment generates the formatted PR comment
func generatePRComment(analysis PRAnalysis, result *promptext.Result) string {
	var comment strings.Builder

	comment.WriteString("## 🤖 AI PR Assistant Analysis\n\n")

	// Changelog section
	if analysis.ChangelogNeeded {
		comment.WriteString("### 📝 Changelog Update Needed\n\n")
		comment.WriteString("This PR modifies library or internal code but doesn't update CHANGELOG.md\n\n")
		comment.WriteString("**Suggested entry:**\n")
		comment.WriteString(analysis.ChangelogSuggestion)
		comment.WriteString("\n\n")
	} else {
		comment.WriteString("### ✅ Changelog\n\n")
		comment.WriteString("Changelog appears to be updated\n\n")
	}

	// Documentation section
	if len(analysis.DocsNeeded) > 0 {
		comment.WriteString("### 📚 Documentation Updates Suggested\n\n")
		for _, doc := range analysis.DocsNeeded {
			comment.WriteString(fmt.Sprintf("- [ ] %s\n", doc))
		}
		comment.WriteString("\n")
	}

	// Examples section
	if len(analysis.ExamplesNeeded) > 0 {
		comment.WriteString("### 💡 Example Updates Suggested\n\n")
		for _, example := range analysis.ExamplesNeeded {
			comment.WriteString(fmt.Sprintf("- [ ] %s\n", example))
		}
		comment.WriteString("\n")
	}

	// Breaking changes
	if len(analysis.BreakingChanges) > 0 {
		comment.WriteString("### ⚠️ Potential Breaking Changes Detected\n\n")
		for _, file := range analysis.BreakingChanges {
			comment.WriteString(fmt.Sprintf("- `%s`\n", file))
		}
		comment.WriteString("\nConsider:\n")
		comment.WriteString("- [ ] Version bump (major version)\n")
		comment.WriteString("- [ ] Migration guide\n")
		comment.WriteString("- [ ] Deprecation warnings\n\n")
	}

	// Statistics
	comment.WriteString("### 📊 Analysis Statistics\n\n")
	comment.WriteString(fmt.Sprintf("- **Files analyzed**: %d\n", len(result.ProjectOutput.Files)))
	comment.WriteString(fmt.Sprintf("- **Context extracted**: ~%d tokens\n", result.TokenCount))
	comment.WriteString(fmt.Sprintf("- **Project language**: %s\n", result.ProjectOutput.Metadata.Language))
	comment.WriteString("\n")

	// Footer
	comment.WriteString("---\n")
	comment.WriteString("*Powered by [promptext](https://github.com/1broseidon/promptext) 🚀*\n")

	return comment.String()
}
