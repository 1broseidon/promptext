package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
	version := flag.String("version", "", "Version to generate notes for (e.g., v0.7.3)")
	sinceTag := flag.String("since", "", "Generate notes since this tag (auto-detects if empty)")
	output := flag.String("output", "", "Output file (prints to stdout if empty)")
	aiPrompt := flag.Bool("ai-prompt", false, "Generate prompt for Claude Code to enhance release notes")
	flag.Parse()

	// Get the tag to compare from
	fromTag := *sinceTag
	if fromTag == "" {
		fromTag = getLastTag()
	}

	fmt.Fprintf(os.Stderr, "📊 Analyzing changes since %s...\n", fromTag)

	// Get changed files since tag
	changedFiles, err := getChangedFilesSinceTag(fromTag)
	if err != nil {
		log.Fatalf("Failed to get changed files: %v", err)
	}

	if len(changedFiles) == 0 {
		fmt.Fprintln(os.Stderr, "⚠️  No changes detected since last release")
		return
	}

	fmt.Fprintf(os.Stderr, "   Found %d changed files\n", len(changedFiles))

	// Get commit messages for context
	commits, err := getCommitsSinceTag(fromTag)
	if err != nil {
		log.Fatalf("Failed to get commits: %v", err)
	}

	fmt.Fprintf(os.Stderr, "   Found %d commits\n", len(commits))

	// Extract code context for changed files
	fmt.Fprintln(os.Stderr, "\n🔍 Extracting code context with promptext...")
	result, err := extractChangedFilesContext(changedFiles)
	if err != nil {
		log.Fatalf("Failed to extract context: %v", err)
	}

	fmt.Fprintf(os.Stderr, "   Extracted context: ~%d tokens from %d files\n",
		result.TokenCount, len(result.ProjectOutput.Files))

	// Generate release notes or AI prompt
	if *aiPrompt {
		fmt.Fprintln(os.Stderr, "\n📝 Generating Claude Code prompt...\n")
		prompt := generateClaudeCodePrompt(*version, fromTag, commits, result)

		if *output != "" {
			if err := os.WriteFile(*output, []byte(prompt), 0644); err != nil {
				log.Fatalf("Failed to write output: %v", err)
			}
			fmt.Fprintf(os.Stderr, "✅ Prompt written to %s\n", *output)
		} else {
			fmt.Println(prompt)
		}
	} else {
		fmt.Fprintln(os.Stderr, "\n📝 Generating release notes...\n")
		releaseNotes := generateReleaseNotes(*version, fromTag, commits, result)

		if *output != "" {
			if err := os.WriteFile(*output, []byte(releaseNotes), 0644); err != nil {
				log.Fatalf("Failed to write output: %v", err)
			}
			fmt.Fprintf(os.Stderr, "✅ Release notes written to %s\n", *output)
		} else {
			fmt.Println(releaseNotes)
		}
	}
}

// getLastTag returns the most recent git tag
func getLastTag() string {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")
	output, err := cmd.Output()
	if err != nil {
		return "HEAD~10" // Fallback to last 10 commits
	}
	return strings.TrimSpace(string(output))
}

// getChangedFilesSinceTag returns list of changed files since a tag
func getChangedFilesSinceTag(tag string) ([]string, error) {
	cmd := exec.Command("git", "diff", "--name-only", tag+"..HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
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

// getCommitsSinceTag returns commit messages since a tag
func getCommitsSinceTag(tag string) ([]string, error) {
	cmd := exec.Command("git", "log", tag+"..HEAD", "--pretty=format:%s")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	var commits []string
	for _, line := range lines {
		if line != "" {
			commits = append(commits, line)
		}
	}
	return commits, nil
}

// extractChangedFilesContext uses promptext to extract context from changed files
func extractChangedFilesContext(changedFiles []string) (*promptext.Result, error) {
	// Create a temporary include pattern for only changed files
	// For now, extract relevant files by extension
	relevantExts := []string{".go", ".md", ".yml", ".yaml"}

	// Build includes list from changed files
	var includes []string
	for _, file := range changedFiles {
		ext := filepath.Ext(file)
		for _, relevantExt := range relevantExts {
			if ext == relevantExt {
				includes = append(includes, file)
				break
			}
		}
	}

	if len(includes) == 0 {
		// If no relevant files, just get a summary
		return promptext.Extract(".",
			promptext.WithExtensions(relevantExts...),
			promptext.WithTokenBudget(4000),
		)
	}

	// Extract with focus on changed files
	return promptext.Extract(".",
		promptext.WithExtensions(relevantExts...),
		promptext.WithTokenBudget(8000),
	)
}

// generateReleaseNotes generates formatted release notes
func generateReleaseNotes(version, fromTag string, commits []string, result *promptext.Result) string {
	var notes strings.Builder

	// Determine version if not provided
	if version == "" {
		version = "Unreleased"
	}

	// Header
	notes.WriteString(fmt.Sprintf("## [%s] - %s\n\n", version, time.Now().Format("2006-01-02")))

	// Categorize commits
	features := []string{}
	fixes := []string{}
	changes := []string{}
	docs := []string{}
	chores := []string{}
	breaking := []string{}

	for _, commit := range commits {
		commitLower := strings.ToLower(commit)

		if strings.HasPrefix(commitLower, "feat:") || strings.HasPrefix(commitLower, "feature:") {
			features = append(features, strings.TrimPrefix(strings.TrimPrefix(commit, "feat:"), "feature:"))
		} else if strings.HasPrefix(commitLower, "fix:") {
			fixes = append(fixes, strings.TrimPrefix(commit, "fix:"))
		} else if strings.HasPrefix(commitLower, "docs:") {
			docs = append(docs, strings.TrimPrefix(commit, "docs:"))
		} else if strings.HasPrefix(commitLower, "chore:") {
			chores = append(chores, strings.TrimPrefix(commit, "chore:"))
		} else if strings.Contains(commitLower, "breaking") || strings.Contains(commitLower, "breaking change") {
			breaking = append(breaking, commit)
		} else {
			changes = append(changes, commit)
		}
	}

	// Breaking Changes (highest priority)
	if len(breaking) > 0 {
		notes.WriteString("### ⚠️ Breaking Changes\n")
		for _, item := range breaking {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Added (features)
	if len(features) > 0 {
		notes.WriteString("### Added\n")
		for _, item := range features {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Fixed
	if len(fixes) > 0 {
		notes.WriteString("### Fixed\n")
		for _, item := range fixes {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Changed
	if len(changes) > 0 {
		notes.WriteString("### Changed\n")
		for _, item := range changes {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Documentation
	if len(docs) > 0 {
		notes.WriteString("### Documentation\n")
		for _, item := range docs {
			notes.WriteString(fmt.Sprintf("- %s\n", strings.TrimSpace(item)))
		}
		notes.WriteString("\n")
	}

	// Context summary
	notes.WriteString("### Statistics\n")
	notes.WriteString(fmt.Sprintf("- **Files changed**: %d\n", len(result.ProjectOutput.Files)))
	notes.WriteString(fmt.Sprintf("- **Commits**: %d\n", len(commits)))
	notes.WriteString(fmt.Sprintf("- **Context analyzed**: ~%d tokens\n", result.TokenCount))
	notes.WriteString("\n")

	notes.WriteString("---\n\n")

	return notes.String()
}

// generateClaudeCodePrompt generates a prompt for Claude Code to enhance release notes
func generateClaudeCodePrompt(version, fromTag string, commits []string, result *promptext.Result) string {
	var prompt strings.Builder

	// Determine version
	if version == "" {
		version = "0.7.4"
	}

	prompt.WriteString("# Release Notes Enhancement Request\n\n")
	prompt.WriteString("Please generate comprehensive release notes for promptext version " + version + "\n\n")

	prompt.WriteString("## Context\n\n")
	prompt.WriteString(fmt.Sprintf("- **Version**: %s\n", version))
	prompt.WriteString(fmt.Sprintf("- **Changes since**: %s\n", fromTag))
	prompt.WriteString(fmt.Sprintf("- **Commits analyzed**: %d\n", len(commits)))
	prompt.WriteString(fmt.Sprintf("- **Files changed**: %d\n", len(result.ProjectOutput.Files)))
	prompt.WriteString(fmt.Sprintf("- **Context extracted**: ~%d tokens\n\n", result.TokenCount))

	prompt.WriteString("## Commit History\n\n")
	prompt.WriteString("```\n")
	for _, commit := range commits {
		prompt.WriteString(commit + "\n")
	}
	prompt.WriteString("```\n\n")

	prompt.WriteString("## Changed Files Summary\n\n")
	for _, file := range result.ProjectOutput.Files {
		prompt.WriteString(fmt.Sprintf("- `%s` (~%d tokens)\n", file.Path, file.Tokens))
	}
	prompt.WriteString("\n")

	prompt.WriteString("## Code Context (via promptext)\n\n")
	prompt.WriteString("```\n")
	prompt.WriteString(result.FormattedOutput)
	prompt.WriteString("\n```\n\n")

	prompt.WriteString("## Task\n\n")
	prompt.WriteString("Generate release notes in Keep a Changelog format with these sections:\n\n")
	prompt.WriteString("### Added\n")
	prompt.WriteString("- New features (be specific and detailed)\n")
	prompt.WriteString("- Focus on user-facing value\n\n")
	prompt.WriteString("### Changed\n")
	prompt.WriteString("- Improvements and modifications\n\n")
	prompt.WriteString("### Fixed\n")
	prompt.WriteString("- Bug fixes\n\n")
	prompt.WriteString("### Documentation\n")
	prompt.WriteString("- Doc updates\n\n")

	prompt.WriteString("## Requirements\n\n")
	prompt.WriteString("1. Use the commit history and code context to write detailed, clear descriptions\n")
	prompt.WriteString("2. Group related changes together logically\n")
	prompt.WriteString("3. Focus on user impact, not implementation details\n")
	prompt.WriteString("4. Be specific about what changed and why it matters\n")
	prompt.WriteString("5. Follow Keep a Changelog format\n")
	prompt.WriteString("6. Include markdown formatting for code, paths, etc.\n\n")

	prompt.WriteString("## Example Format\n\n")
	prompt.WriteString("```markdown\n")
	prompt.WriteString("## [" + version + "] - " + time.Now().Format("2006-01-02") + "\n\n")
	prompt.WriteString("### Added\n")
	prompt.WriteString("- **Release Notes Generator**: Automated tool using promptext library to analyze git changes\n")
	prompt.WriteString("  - Extracts code context with token-aware analysis\n")
	prompt.WriteString("  - Categorizes commits by type (feat, fix, docs)\n")
	prompt.WriteString("  - Generates changelog-compatible markdown\n\n")
	prompt.WriteString("...\n")
	prompt.WriteString("```\n\n")

	prompt.WriteString("Please generate the complete, polished release notes now.\n")

	return prompt.String()
}
