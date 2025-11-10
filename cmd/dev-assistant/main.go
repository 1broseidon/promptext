package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/1broseidon/promptext/pkg/promptext"
)

type CheckResult struct {
	NeedsChangelog  bool
	ChangelogReason string
	DocsToUpdate    []string
	ExamplesToUpdate []string
	TestSuggestions  []string
}

func main() {
	autoFix := flag.Bool("auto-fix", false, "Automatically apply suggested fixes")
	dryRun := flag.Bool("dry-run", false, "Show what would be done without making changes")
	flag.Parse()

	fmt.Println("🔍 Analyzing staged changes...")
	fmt.Println()

	// Get staged files
	stagedFiles, err := getStagedFiles()
	if err != nil {
		log.Fatalf("Failed to get staged files: %v", err)
	}

	if len(stagedFiles) == 0 {
		fmt.Println("⚠️  No staged changes detected")
		return
	}

	fmt.Printf("   Found %d staged files\n", len(stagedFiles))

	// Extract context with promptext
	result, err := extractStagedContext(stagedFiles)
	if err != nil {
		log.Fatalf("Failed to extract context: %v", err)
	}

	fmt.Printf("   Extracted ~%d tokens from %d files\n", result.TokenCount, len(result.ProjectOutput.Files))
	fmt.Println()

	// Analyze changes
	checks := analyzeChanges(stagedFiles, result)

	// Display results
	displayChecks(checks)

	// Handle actions
	if !*dryRun && !*autoFix {
		handleInteractive(checks)
	} else if *autoFix {
		fmt.Println("\n🔧 Auto-fix mode enabled")
		applyFixes(checks)
	}
}

// getStagedFiles returns list of staged files
func getStagedFiles() ([]string, error) {
	cmd := exec.Command("git", "diff", "--cached", "--name-only")
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

// extractStagedContext uses promptext to extract context
func extractStagedContext(stagedFiles []string) (*promptext.Result, error) {
	// Focus on relevant file types
	relevantExts := []string{".go", ".md", ".yml", ".yaml"}

	return promptext.Extract(".",
		promptext.WithExtensions(relevantExts...),
		promptext.WithTokenBudget(8000),
	)
}

// analyzeChanges analyzes staged changes for issues
func analyzeChanges(stagedFiles []string, result *promptext.Result) CheckResult {
	checks := CheckResult{
		DocsToUpdate:     []string{},
		ExamplesToUpdate: []string{},
		TestSuggestions:  []string{},
	}

	// Check what types of files changed
	hasLibraryChanges := false
	hasInternalChanges := false
	hasChangelogUpdate := false
	hasDocChanges := false
	hasExampleChanges := false
	hasTestChanges := false

	for _, file := range stagedFiles {
		if strings.HasPrefix(file, "pkg/promptext/") {
			hasLibraryChanges = true
		}
		if strings.HasPrefix(file, "internal/") {
			hasInternalChanges = true
		}
		if file == "CHANGELOG.md" || strings.Contains(file, "changelog.md") {
			hasChangelogUpdate = true
		}
		if strings.HasPrefix(file, "docs-astro/") || strings.HasSuffix(file, ".md") {
			hasDocChanges = true
		}
		if strings.HasPrefix(file, "examples/") {
			hasExampleChanges = true
		}
		if strings.HasSuffix(file, "_test.go") {
			hasTestChanges = true
		}
	}

	// Check if changelog needed
	if (hasLibraryChanges || hasInternalChanges) && !hasChangelogUpdate {
		checks.NeedsChangelog = true
		if hasLibraryChanges {
			checks.ChangelogReason = "Public API changes detected in pkg/promptext/"
		} else {
			checks.ChangelogReason = "Internal changes detected"
		}
	}

	// Check if docs needed
	if hasLibraryChanges && !hasDocChanges {
		checks.DocsToUpdate = append(checks.DocsToUpdate,
			"docs-astro/src/content/docs/library-usage.md - Library changes detected")
	}

	// Check if examples needed
	if hasLibraryChanges && !hasExampleChanges {
		checks.ExamplesToUpdate = append(checks.ExamplesToUpdate,
			"examples/ - New library features should have examples")
	}

	// Check if tests needed
	if (hasLibraryChanges || hasInternalChanges) && !hasTestChanges {
		checks.TestSuggestions = append(checks.TestSuggestions,
			"Consider adding tests for new functionality")
	}

	return checks
}

// displayChecks shows the analysis results
func displayChecks(checks CheckResult) {
	fmt.Println("📋 Analysis Results:")
	fmt.Println()

	hasIssues := false

	if checks.NeedsChangelog {
		hasIssues = true
		fmt.Println("  ⚠️  Changelog entry recommended")
		fmt.Printf("      Reason: %s\n", checks.ChangelogReason)
		fmt.Println()
	}

	if len(checks.DocsToUpdate) > 0 {
		hasIssues = true
		fmt.Println("  ⚠️  Documentation updates needed")
		for _, doc := range checks.DocsToUpdate {
			fmt.Printf("      - %s\n", doc)
		}
		fmt.Println()
	}

	if len(checks.ExamplesToUpdate) > 0 {
		hasIssues = true
		fmt.Println("  ⚠️  Example updates suggested")
		for _, example := range checks.ExamplesToUpdate {
			fmt.Printf("      - %s\n", example)
		}
		fmt.Println()
	}

	if len(checks.TestSuggestions) > 0 {
		hasIssues = true
		fmt.Println("  💡 Test suggestions")
		for _, suggestion := range checks.TestSuggestions {
			fmt.Printf("      - %s\n", suggestion)
		}
		fmt.Println()
	}

	if !hasIssues {
		fmt.Println("  ✅ Everything looks good!")
		fmt.Println()
	}
}

// handleInteractive prompts user for actions
func handleInteractive(checks CheckResult) {
	reader := bufio.NewReader(os.Stdin)

	if checks.NeedsChangelog {
		fmt.Print("\nAdd changelog entry? [y/N]: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			fmt.Println("\n💡 Add your changelog entry to CHANGELOG.md")
			fmt.Println("   Then stage the file: git add CHANGELOG.md")
		}
	}

	if len(checks.DocsToUpdate) > 0 {
		fmt.Print("\nOpen documentation files? [y/N]: ")
		response, _ := reader.ReadString('\n')
		response = strings.TrimSpace(strings.ToLower(response))

		if response == "y" || response == "yes" {
			for _, doc := range checks.DocsToUpdate {
				parts := strings.SplitN(doc, " - ", 2)
				if len(parts) > 0 {
					file := parts[0]
					fmt.Printf("\n📝 Opening %s\n", file)
					fmt.Printf("   Edit the file, then stage: git add %s\n", file)
				}
			}
		}
	}

	fmt.Println()
	fmt.Println("✨ Done! Remember to review your changes before committing.")
}

// applyFixes automatically applies fixes
func applyFixes(checks CheckResult) {
	if checks.NeedsChangelog {
		fmt.Println("  📝 Changelog needs manual update")
		fmt.Println("     Please add an entry to CHANGELOG.md")
	}

	if len(checks.DocsToUpdate) > 0 {
		fmt.Println("  📚 Documentation needs manual update")
		for _, doc := range checks.DocsToUpdate {
			fmt.Printf("     - %s\n", doc)
		}
	}

	fmt.Println()
	fmt.Println("  💡 Auto-fix cannot fully automate these changes")
	fmt.Println("     Please review and update manually")
}

// confirm asks for user confirmation
func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("%s [y/N]: ", prompt)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}
