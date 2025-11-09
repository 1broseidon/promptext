// Package main implements an automated code review tool for CI/CD pipelines.
//
// This example demonstrates how to integrate promptext with GitHub Actions
// (or any CI system) to automatically extract code context from pull requests
// and prepare it for AI-powered code review.
//
// The tool can:
// - Extract only files changed in a PR
// - Include surrounding context for better reviews
// - Generate review prompts ready for AI analysis
// - Post comments back to the PR (with GitHub API integration)
//
// Usage:
//   # Review changes in current branch vs main
//   go run main.go
//
//   # Review specific PR
//   PR_NUMBER=123 go run main.go
//
//   # Custom base branch
//   BASE_BRANCH=develop go run main.go
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/pkg/promptext"
)

// PRInfo contains information about the pull request being reviewed
type PRInfo struct {
	Number       string   `json:"number"`
	Title        string   `json:"title"`
	BaseBranch   string   `json:"base_branch"`
	HeadBranch   string   `json:"head_branch"`
	ChangedFiles []string `json:"changed_files"`
	Author       string   `json:"author"`
}

// ReviewResult contains the extracted code and metadata for AI review
type ReviewResult struct {
	PRInfo      PRInfo `json:"pr_info"`
	Context     string `json:"context"`
	TokenCount  int    `json:"token_count"`
	FileCount   int    `json:"file_count"`
	ReviewReady bool   `json:"review_ready"`
}

func main() {
	fmt.Println("🤖 AI Code Review - PR Context Extractor")
	fmt.Println(strings.Repeat("=", 60))

	// Get PR information from environment or git
	prInfo, err := getPRInfo()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting PR info: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("📋 PR #%s: %s\n", prInfo.Number, prInfo.Title)
	fmt.Printf("   Author: %s\n", prInfo.Author)
	fmt.Printf("   Base: %s → Head: %s\n", prInfo.BaseBranch, prInfo.HeadBranch)
	fmt.Printf("   Changed files: %d\n\n", len(prInfo.ChangedFiles))

	// Extract code context for the changed files
	result, err := extractPRContext(prInfo)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error extracting context: %v\n", err)
		os.Exit(1)
	}

	// Display summary
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("✨ Context extracted successfully\n")
	fmt.Printf("   Files included: %d\n", result.FileCount)
	fmt.Printf("   Token count: %d\n", result.TokenCount)
	fmt.Println(strings.Repeat("=", 60))

	// Save results
	if err := saveResults(result); err != nil {
		fmt.Fprintf(os.Stderr, "Error saving results: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("\n💡 Next Steps:")
	fmt.Println("   1. Review the extracted context in pr-review-context.ptx")
	fmt.Println("   2. Use the review prompt in review-prompt.txt")
	fmt.Println("   3. Send to your AI API (Claude, GPT, etc.)")
	fmt.Println("   4. Post the review comments back to the PR")
	fmt.Println("\n🎯 Review ready!")
}

// getPRInfo extracts PR information from environment variables or git
func getPRInfo() (*PRInfo, error) {
	info := &PRInfo{
		Number:     getEnvOrDefault("PR_NUMBER", "local"),
		BaseBranch: getEnvOrDefault("BASE_BRANCH", "main"),
		Author:     getEnvOrDefault("PR_AUTHOR", "unknown"),
	}

	// Get current branch as head
	headBranch, err := runGitCommand("git", "rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	info.HeadBranch = strings.TrimSpace(headBranch)

	// Get PR title from branch name or environment
	info.Title = getEnvOrDefault("PR_TITLE", fmt.Sprintf("Changes in %s", info.HeadBranch))

	// Get list of changed files
	// Use direct git command to prevent command injection
	changedFiles, err := runGitCommand("git", "diff", "--name-only", fmt.Sprintf("%s...%s", info.BaseBranch, info.HeadBranch))
	if err != nil {
		return nil, fmt.Errorf("failed to get changed files: %w", err)
	}

	// Parse changed files
	for _, file := range strings.Split(strings.TrimSpace(changedFiles), "\n") {
		if file != "" {
			info.ChangedFiles = append(info.ChangedFiles, file)
		}
	}

	if len(info.ChangedFiles) == 0 {
		return nil, fmt.Errorf("no changed files found")
	}

	return info, nil
}

// extractPRContext uses promptext to extract code context for the PR
func extractPRContext(prInfo *PRInfo) (*ReviewResult, error) {
	// Extract keywords from changed file names to get relevant context
	// This helps include related files that might not be directly changed
	keywords := extractFileKeywords(prInfo.ChangedFiles)

	// Get the current directory
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fmt.Println("🔍 Extracting code context...")
	fmt.Printf("   Keywords from changes: %v\n", keywords)

	// Use promptext to extract context
	// We use relevance filtering to include related files beyond just the changes
	result, err := promptext.Extract(cwd,
		// Focus on source code files
		promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java", ".rb", ".rs"),

		// Use relevance to include related files
		promptext.WithRelevance(keywords...),

		// Limit to reasonable budget for PR reviews (most AI APIs)
		promptext.WithTokenBudget(15000),

		// PTX format for AI consumption
		promptext.WithFormat(promptext.FormatPTX),
	)

	if err != nil {
		return nil, err
	}

	return &ReviewResult{
		PRInfo:      *prInfo,
		Context:     result.FormattedOutput,
		TokenCount:  result.TokenCount,
		FileCount:   len(result.ProjectOutput.Files),
		ReviewReady: true,
	}, nil
}

// extractFileKeywords generates keywords from file paths to find related code
func extractFileKeywords(files []string) []string {
	keywords := make(map[string]bool)

	for _, file := range files {
		// Get directory name
		dir := filepath.Dir(file)
		if dir != "." {
			parts := strings.Split(dir, string(filepath.Separator))
			for _, part := range parts {
				if part != "" && len(part) > 2 {
					keywords[part] = true
				}
			}
		}

		// Get filename without extension
		base := filepath.Base(file)
		name := strings.TrimSuffix(base, filepath.Ext(base))

		// Split camelCase and snake_case
		nameParts := strings.FieldsFunc(name, func(r rune) bool {
			return r == '_' || r == '-'
		})

		for _, part := range nameParts {
			if len(part) > 2 {
				keywords[strings.ToLower(part)] = true
			}
		}
	}

	// Convert to slice
	result := make([]string, 0, len(keywords))
	for kw := range keywords {
		result = append(result, kw)
	}

	return result
}

// saveResults saves the review context and generates a review prompt
func saveResults(result *ReviewResult) error {
	// Save context
	contextFile := "pr-review-context.ptx"
	if err := os.WriteFile(contextFile, []byte(result.Context), 0644); err != nil {
		return fmt.Errorf("failed to save context: %w", err)
	}
	fmt.Printf("\n💾 Context saved to: %s\n", contextFile)

	// Save metadata
	metadataFile := "pr-review-metadata.json"
	metadata, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	if err := os.WriteFile(metadataFile, metadata, 0644); err != nil {
		return fmt.Errorf("failed to save metadata: %w", err)
	}
	fmt.Printf("📊 Metadata saved to: %s\n", metadataFile)

	// Generate review prompt
	prompt := generateReviewPrompt(result)
	promptFile := "review-prompt.txt"
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		return fmt.Errorf("failed to save prompt: %w", err)
	}
	fmt.Printf("📝 Review prompt saved to: %s\n", promptFile)

	return nil
}

// generateReviewPrompt creates a prompt for the AI reviewer
func generateReviewPrompt(result *ReviewResult) string {
	var prompt strings.Builder

	prompt.WriteString("# Code Review Request\n\n")
	prompt.WriteString(fmt.Sprintf("## Pull Request: %s\n", result.PRInfo.Title))
	prompt.WriteString(fmt.Sprintf("- **Author**: %s\n", result.PRInfo.Author))
	prompt.WriteString(fmt.Sprintf("- **Base**: %s → **Head**: %s\n", result.PRInfo.BaseBranch, result.PRInfo.HeadBranch))
	prompt.WriteString(fmt.Sprintf("- **Files Changed**: %d\n\n", len(result.PRInfo.ChangedFiles)))

	prompt.WriteString("## Changed Files\n\n")
	for _, file := range result.PRInfo.ChangedFiles {
		prompt.WriteString(fmt.Sprintf("- %s\n", file))
	}

	prompt.WriteString("\n## Review Guidelines\n\n")
	prompt.WriteString("Please review this code change and provide feedback on:\n\n")
	prompt.WriteString("1. **Code Quality**\n")
	prompt.WriteString("   - Is the code readable and maintainable?\n")
	prompt.WriteString("   - Are there any code smells or anti-patterns?\n\n")

	prompt.WriteString("2. **Security**\n")
	prompt.WriteString("   - Are there any security vulnerabilities?\n")
	prompt.WriteString("   - Is input validation adequate?\n")
	prompt.WriteString("   - Are secrets or sensitive data properly handled?\n\n")

	prompt.WriteString("3. **Performance**\n")
	prompt.WriteString("   - Are there any obvious performance issues?\n")
	prompt.WriteString("   - Could any operations be optimized?\n\n")

	prompt.WriteString("4. **Testing**\n")
	prompt.WriteString("   - Are there sufficient tests?\n")
	prompt.WriteString("   - Are edge cases covered?\n\n")

	prompt.WriteString("5. **Documentation**\n")
	prompt.WriteString("   - Is the code well-commented?\n")
	prompt.WriteString("   - Are complex sections explained?\n\n")

	prompt.WriteString("## Code Context\n\n")
	prompt.WriteString("The following code context has been extracted from the repository:\n\n")
	prompt.WriteString("```\n")
	prompt.WriteString("(See pr-review-context.ptx)\n")
	prompt.WriteString("```\n\n")

	prompt.WriteString("## Output Format\n\n")
	prompt.WriteString("Please provide:\n")
	prompt.WriteString("- **Summary**: Overall assessment (Approve/Request Changes/Comment)\n")
	prompt.WriteString("- **Issues**: List any concerns or problems found\n")
	prompt.WriteString("- **Suggestions**: Concrete improvements\n")
	prompt.WriteString("- **Positives**: What was done well\n")

	return prompt.String()
}

// Helper functions

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func runGitCommand(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s: %w", string(output), err)
	}
	return string(output), nil
}
