// Package main implements a legacy code migration assistant using promptext.
//
// This tool helps modernize legacy codebases by extracting relevant code,
// analyzing dependencies, and preparing context for AI-assisted migration.
//
// Common migration scenarios:
// - Python 2 → Python 3
// - JavaScript (CommonJS) → TypeScript (ES Modules)
// - Legacy frameworks → Modern alternatives
// - Monolith → Microservices
// - Deprecated APIs → Current standards
//
// Usage:
//   # Analyze legacy code
//   go run main.go --analyze authentication
//
//   # Plan migration
//   go run main.go --plan --component auth --target modern-auth
//
//   # Generate migration steps
//   go run main.go --migrate --from legacy-db --to orm
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/1broseidon/promptext/pkg/promptext"
)

// MigrationPhase represents the current migration phase
type MigrationPhase string

const (
	PhaseAnalyze MigrationPhase = "analyze"
	PhasePlan    MigrationPhase = "plan"
	PhaseMigrate MigrationPhase = "migrate"
	PhaseVerify  MigrationPhase = "verify"
)

// MigrationConfig holds migration configuration
type MigrationConfig struct {
	Phase         MigrationPhase
	Component     string   // Component to migrate (e.g., "auth", "database")
	SourceDir     string   // Legacy code directory
	TargetPattern string   // Target pattern/framework
	Keywords      []string // Keywords to identify relevant code
	OutputDir     string   // Where to save migration artifacts
}

// MigrationReport contains analysis and recommendations
type MigrationReport struct {
	Timestamp    time.Time              `json:"timestamp"`
	Component    string                 `json:"component"`
	FilesAnalyzed int                   `json:"files_analyzed"`
	TokenCount   int                    `json:"token_count"`
	Dependencies []string               `json:"dependencies"`
	Issues       []MigrationIssue       `json:"issues"`
	Steps        []MigrationStep        `json:"steps"`
	Metadata     map[string]interface{} `json:"metadata"`
}

// MigrationIssue represents a problem found in legacy code
type MigrationIssue struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"`
	Description string `json:"description"`
	File        string `json:"file,omitempty"`
	Line        int    `json:"line,omitempty"`
}

// MigrationStep represents a step in the migration plan
type MigrationStep struct {
	Order       int      `json:"order"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Files       []string `json:"files"`
	Estimated   string   `json:"estimated_effort"`
}

func main() {
	config := parseFlags()

	fmt.Println("🔄 Migration Assistant")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Phase: %s\n", config.Phase)
	fmt.Printf("Component: %s\n", config.Component)
	fmt.Println(strings.Repeat("=", 60))

	switch config.Phase {
	case PhaseAnalyze:
		if err := analyzeLegacyCode(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error analyzing code: %v\n", err)
			os.Exit(1)
		}

	case PhasePlan:
		if err := planMigration(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error planning migration: %v\n", err)
			os.Exit(1)
		}

	case PhaseMigrate:
		if err := generateMigrationSteps(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error generating migration: %v\n", err)
			os.Exit(1)
		}

	case PhaseVerify:
		if err := verifyMigration(config); err != nil {
			fmt.Fprintf(os.Stderr, "Error verifying migration: %v\n", err)
			os.Exit(1)
		}
	}

	fmt.Println("\n🎯 Migration assistant complete!")
}

func parseFlags() *MigrationConfig {
	config := &MigrationConfig{
		Keywords: make([]string, 0),
	}

	analyze := flag.Bool("analyze", false, "Analyze legacy code")
	plan := flag.Bool("plan", false, "Create migration plan")
	migrate := flag.Bool("migrate", false, "Generate migration steps")
	verify := flag.Bool("verify", false, "Verify migration")

	flag.StringVar(&config.Component, "component", "", "Component to migrate (required)")
	flag.StringVar(&config.SourceDir, "source", ".", "Source directory")
	flag.StringVar(&config.TargetPattern, "target", "", "Target pattern/framework")
	flag.StringVar(&config.OutputDir, "output", "migration-output", "Output directory")

	flag.Parse()

	// Determine phase
	switch {
	case *analyze:
		config.Phase = PhaseAnalyze
	case *plan:
		config.Phase = PhasePlan
	case *migrate:
		config.Phase = PhaseMigrate
	case *verify:
		config.Phase = PhaseVerify
	default:
		config.Phase = PhaseAnalyze // Default to analyze
	}

	// Use remaining args as keywords
	config.Keywords = append(config.Keywords, flag.Args()...)

	// Add component as keyword
	if config.Component != "" {
		config.Keywords = append(config.Keywords, config.Component)
	}

	return config
}

func analyzeLegacyCode(config *MigrationConfig) error {
	fmt.Println("\n🔍 Analyzing legacy code...")

	// Extract legacy code using promptext
	result, err := promptext.Extract(config.SourceDir,
		// Include common source code extensions
		promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java", ".rb", ".php"),

		// Use relevance filtering to find component-specific code
		promptext.WithRelevance(config.Keywords...),

		// Generous budget for comprehensive analysis
		promptext.WithTokenBudget(25000),

		// Markdown for human-readable analysis
		promptext.WithFormat(promptext.FormatMarkdown),
	)

	if err != nil {
		return fmt.Errorf("failed to extract code: %w", err)
	}

	fmt.Printf("✅ Extracted %d files (%d tokens)\n",
		len(result.ProjectOutput.Files), result.TokenCount)

	// Create migration report
	report := &MigrationReport{
		Timestamp:    time.Now(),
		Component:    config.Component,
		FilesAnalyzed: len(result.ProjectOutput.Files),
		TokenCount:   result.TokenCount,
		Dependencies: extractDependencies(result),
		Issues:       identifyIssues(result, config),
		Metadata:     make(map[string]interface{}),
	}

	// Save analysis results
	if err := saveAnalysisResults(config, result, report); err != nil {
		return err
	}

	// Display summary
	displayAnalysisSummary(report)

	return nil
}

func planMigration(config *MigrationConfig) error {
	fmt.Println("\n📋 Creating migration plan...")

	// Extract code for planning
	result, err := promptext.Extract(config.SourceDir,
		promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java"),
		promptext.WithRelevance(config.Keywords...),
		promptext.WithTokenBudget(20000),
		promptext.WithFormat(promptext.FormatPTX),
	)

	if err != nil {
		return fmt.Errorf("failed to extract code: %w", err)
	}

	// Create migration plan
	plan := createMigrationPlan(config, result)

	// Save plan
	if err := saveMigrationPlan(config, plan, result); err != nil {
		return err
	}

	// Display plan
	displayMigrationPlan(plan)

	return nil
}

func generateMigrationSteps(config *MigrationConfig) error {
	fmt.Println("\n⚙️ Generating migration steps...")

	// Load existing plan if available
	planPath := filepath.Join(config.OutputDir, "migration-plan.json")
	_, err := os.Stat(planPath)
	if os.IsNotExist(err) {
		fmt.Println("ℹ️  No existing plan found. Run --plan first to create a migration plan.")
		fmt.Println("   Proceeding with default analysis...")
	}

	// Extract code
	result, err := promptext.Extract(config.SourceDir,
		promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java"),
		promptext.WithRelevance(config.Keywords...),
		promptext.WithTokenBudget(30000),
		promptext.WithFormat(promptext.FormatPTX),
	)

	if err != nil {
		return fmt.Errorf("failed to extract code: %w", err)
	}

	// Generate migration artifacts
	if err := generateMigrationArtifacts(config, result); err != nil {
		return err
	}

	fmt.Println("\n💾 Migration artifacts generated successfully!")
	fmt.Printf("   Output directory: %s\n", config.OutputDir)

	return nil
}

func verifyMigration(config *MigrationConfig) error {
	fmt.Println("\n✓ Verifying migration...")

	// This would typically compare old vs new code
	// For now, we provide a framework for verification

	fmt.Println("   Verification checklist:")
	fmt.Println("   [ ] All functionality preserved")
	fmt.Println("   [ ] Tests passing")
	fmt.Println("   [ ] Performance acceptable")
	fmt.Println("   [ ] Security requirements met")
	fmt.Println("   [ ] Documentation updated")

	return nil
}

func extractDependencies(result *promptext.Result) []string {
	// Extract dependencies from metadata
	deps := make([]string, 0)

	if result.ProjectOutput.Metadata != nil {
		deps = result.ProjectOutput.Metadata.Dependencies
	}

	return deps
}

func identifyIssues(result *promptext.Result, config *MigrationConfig) []MigrationIssue {
	// In a real implementation, this would analyze the code for:
	// - Deprecated API usage
	// - Security vulnerabilities
	// - Performance bottlenecks
	// - Code smells

	issues := []MigrationIssue{
		{
			Type:        "deprecated_api",
			Severity:    "high",
			Description: "Legacy authentication pattern detected",
		},
		{
			Type:        "security",
			Severity:    "medium",
			Description: "Potential SQL injection vulnerability",
		},
	}

	return issues
}

func createMigrationPlan(config *MigrationConfig, result *promptext.Result) *MigrationReport {
	report := &MigrationReport{
		Timestamp:    time.Now(),
		Component:    config.Component,
		FilesAnalyzed: len(result.ProjectOutput.Files),
		TokenCount:   result.TokenCount,
		Steps: []MigrationStep{
			{
				Order:       1,
				Title:       "Analyze Dependencies",
				Description: "Review all dependencies and identify outdated packages",
				Estimated:   "2 hours",
			},
			{
				Order:       2,
				Title:       "Update Core Patterns",
				Description: fmt.Sprintf("Migrate %s to %s patterns", config.Component, config.TargetPattern),
				Estimated:   "1 day",
			},
			{
				Order:       3,
				Title:       "Update Tests",
				Description: "Adapt test suite for new implementation",
				Estimated:   "4 hours",
			},
			{
				Order:       4,
				Title:       "Verify Migration",
				Description: "Run full test suite and verify functionality",
				Estimated:   "2 hours",
			},
		},
	}

	return report
}

func saveAnalysisResults(config *MigrationConfig, result *promptext.Result, report *MigrationReport) error {
	// Create output directory
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return err
	}

	// Save extracted code context
	contextFile := filepath.Join(config.OutputDir, "legacy-code-context.md")
	if err := os.WriteFile(contextFile, []byte(result.FormattedOutput), 0644); err != nil {
		return err
	}

	// Save analysis report
	reportFile := filepath.Join(config.OutputDir, "analysis-report.json")
	reportJSON, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(reportFile, reportJSON, 0644); err != nil {
		return err
	}

	// Generate analysis prompt for AI
	prompt := generateAnalysisPrompt(config, report)
	promptFile := filepath.Join(config.OutputDir, "analysis-prompt.txt")
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		return err
	}

	fmt.Printf("\n💾 Analysis results saved:\n")
	fmt.Printf("   - Context: %s\n", contextFile)
	fmt.Printf("   - Report: %s\n", reportFile)
	fmt.Printf("   - Prompt: %s\n", promptFile)

	return nil
}

func saveMigrationPlan(config *MigrationConfig, plan *MigrationReport, result *promptext.Result) error {
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return err
	}

	// Save plan
	planFile := filepath.Join(config.OutputDir, "migration-plan.json")
	planJSON, err := json.MarshalIndent(plan, "", "  ")
	if err != nil {
		return err
	}
	if err := os.WriteFile(planFile, planJSON, 0644); err != nil {
		return err
	}

	// Save context
	contextFile := filepath.Join(config.OutputDir, "migration-context.ptx")
	if err := os.WriteFile(contextFile, []byte(result.FormattedOutput), 0644); err != nil {
		return err
	}

	// Generate planning prompt
	prompt := generatePlanningPrompt(config, plan)
	promptFile := filepath.Join(config.OutputDir, "migration-prompt.txt")
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		return err
	}

	fmt.Printf("\n💾 Migration plan saved:\n")
	fmt.Printf("   - Plan: %s\n", planFile)
	fmt.Printf("   - Context: %s\n", contextFile)
	fmt.Printf("   - Prompt: %s\n", promptFile)

	return nil
}

func generateMigrationArtifacts(config *MigrationConfig, result *promptext.Result) error {
	if err := os.MkdirAll(config.OutputDir, 0755); err != nil {
		return err
	}

	// Save migration context
	contextFile := filepath.Join(config.OutputDir, "migration-code.ptx")
	if err := os.WriteFile(contextFile, []byte(result.FormattedOutput), 0644); err != nil {
		return err
	}

	// Generate migration prompt
	prompt := generateMigrationPrompt(config)
	promptFile := filepath.Join(config.OutputDir, "migration-instructions.txt")
	if err := os.WriteFile(promptFile, []byte(prompt), 0644); err != nil {
		return err
	}

	fmt.Printf("\n💾 Migration artifacts saved:\n")
	fmt.Printf("   - Code context: %s\n", contextFile)
	fmt.Printf("   - Instructions: %s\n", promptFile)

	return nil
}

func generateAnalysisPrompt(config *MigrationConfig, report *MigrationReport) string {
	var prompt strings.Builder

	prompt.WriteString("# Legacy Code Analysis Request\n\n")
	prompt.WriteString(fmt.Sprintf("## Component: %s\n\n", config.Component))
	prompt.WriteString("## Task\n")
	prompt.WriteString("Analyze the legacy code and provide a comprehensive assessment for migration planning.\n\n")
	prompt.WriteString("## Analysis Requirements\n\n")
	prompt.WriteString("1. **Architecture Review**\n")
	prompt.WriteString("   - Identify current architectural patterns\n")
	prompt.WriteString("   - Assess code organization and structure\n")
	prompt.WriteString("   - Note coupling and dependencies\n\n")
	prompt.WriteString("2. **Security Analysis**\n")
	prompt.WriteString("   - Identify security vulnerabilities\n")
	prompt.WriteString("   - Check for deprecated security practices\n")
	prompt.WriteString("   - Assess authentication/authorization patterns\n\n")
	prompt.WriteString("3. **Technical Debt**\n")
	prompt.WriteString("   - Code smells and anti-patterns\n")
	prompt.WriteString("   - Deprecated API usage\n")
	prompt.WriteString("   - Performance issues\n\n")
	prompt.WriteString("4. **Migration Risks**\n")
	prompt.WriteString("   - Breaking changes\n")
	prompt.WriteString("   - Data migration concerns\n")
	prompt.WriteString("   - Compatibility issues\n\n")
	prompt.WriteString("## Code Context\n\n")
	prompt.WriteString(fmt.Sprintf("- **Files Analyzed**: %d\n", report.FilesAnalyzed))
	prompt.WriteString(fmt.Sprintf("- **Token Count**: %d\n", report.TokenCount))
	prompt.WriteString(fmt.Sprintf("- **Dependencies**: %d\n\n", len(report.Dependencies)))
	prompt.WriteString("_See accompanying context file for full code details._\n")

	return prompt.String()
}

func generatePlanningPrompt(config *MigrationConfig, plan *MigrationReport) string {
	var prompt strings.Builder

	prompt.WriteString("# Migration Planning Request\n\n")
	prompt.WriteString(fmt.Sprintf("## Component: %s\n", config.Component))
	prompt.WriteString(fmt.Sprintf("## Target: %s\n\n", config.TargetPattern))
	prompt.WriteString("## Task\n")
	prompt.WriteString("Create a detailed, step-by-step migration plan.\n\n")
	prompt.WriteString("## Plan Requirements\n\n")
	prompt.WriteString("1. **Prioritized Steps**: Logical order of operations\n")
	prompt.WriteString("2. **Effort Estimates**: Time estimates for each step\n")
	prompt.WriteString("3. **Risk Assessment**: Potential issues and mitigations\n")
	prompt.WriteString("4. **Testing Strategy**: How to verify each step\n")
	prompt.WriteString("5. **Rollback Plan**: How to revert if needed\n\n")

	return prompt.String()
}

func generateMigrationPrompt(config *MigrationConfig) string {
	var prompt strings.Builder

	prompt.WriteString("# Code Migration Request\n\n")
	prompt.WriteString(fmt.Sprintf("## Component: %s\n", config.Component))
	prompt.WriteString(fmt.Sprintf("## Target Pattern: %s\n\n", config.TargetPattern))
	prompt.WriteString("## Task\n")
	prompt.WriteString("Provide detailed migration steps and code examples.\n\n")
	prompt.WriteString("## Deliverables\n\n")
	prompt.WriteString("1. **Step-by-Step Guide**: Detailed migration instructions\n")
	prompt.WriteString("2. **Code Examples**: Before/after code snippets\n")
	prompt.WriteString("3. **Test Cases**: How to verify the migration\n")
	prompt.WriteString("4. **Common Pitfalls**: Issues to watch out for\n")
	prompt.WriteString("5. **Best Practices**: Recommended patterns for the target\n\n")

	return prompt.String()
}

func displayAnalysisSummary(report *MigrationReport) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📊 Analysis Summary")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Printf("Files analyzed: %d\n", report.FilesAnalyzed)
	fmt.Printf("Token count: %d\n", report.TokenCount)
	fmt.Printf("Dependencies: %d\n", len(report.Dependencies))
	fmt.Printf("Issues found: %d\n", len(report.Issues))

	if len(report.Issues) > 0 {
		fmt.Println("\n⚠️  Issues Detected:")
		for _, issue := range report.Issues {
			fmt.Printf("   [%s] %s: %s\n", issue.Severity, issue.Type, issue.Description)
		}
	}
}

func displayMigrationPlan(plan *MigrationReport) {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📋 Migration Plan")
	fmt.Println(strings.Repeat("=", 60))

	for _, step := range plan.Steps {
		fmt.Printf("\n%d. %s\n", step.Order, step.Title)
		fmt.Printf("   Description: %s\n", step.Description)
		fmt.Printf("   Estimated: %s\n", step.Estimated)
	}
}
