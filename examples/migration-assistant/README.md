# Migration Assistant

A powerful tool for modernizing legacy codebases using AI-powered analysis and planning. The Migration Assistant helps you understand old code patterns, identify issues, and create systematic migration plans.

## Overview

The Migration Assistant provides a structured 4-phase approach to code migration:

1. **Analyze**: Deep dive into legacy code to understand patterns, dependencies, and potential issues
2. **Plan**: Create a detailed migration strategy with prioritized steps
3. **Migrate**: Execute the migration with focused context for each component
4. **Verify**: Validate changes and ensure nothing was broken

## Installation

```bash
go get github.com/1broseidon/promptext/pkg/promptext
```

## Usage

### Basic Usage

```bash
# Analyze a legacy component
go run main.go analyze /path/to/legacy/code user-auth

# Create migration plan
go run main.go plan /path/to/legacy/code user-auth

# Get migration context for implementation
go run main.go migrate /path/to/legacy/code user-auth

# Verify after migration
go run main.go verify /path/to/legacy/code user-auth
```

### Phase 1: Analyze

Extract comprehensive context about legacy code:

```bash
go run main.go analyze ./old-app authentication
```

**What it does**:
- Extracts all files related to the component
- Identifies dependencies and imports
- Detects common issues (globals, deprecated APIs, security concerns)
- Provides detailed analysis for AI review

**Example output**:
```
=== Migration Analysis: authentication ===

Files analyzed: 12
Token count: 8,543
Dependencies: crypto/md5, database/sql, net/http

Issues detected:
  - security: Using deprecated MD5 for password hashing (auth/hash.go:45)
  - globals: Global database connection variable (auth/db.go:12)
  - deprecated: Using ioutil instead of io/os (auth/file.go:8)

Analysis saved to: migration-analysis-authentication-20250109.txt
```

### Phase 2: Plan

Create a systematic migration strategy:

```bash
go run main.go plan ./old-app authentication
```

**What it does**:
- Extracts relevant code for migration planning
- Focuses on entry points and high-priority files
- Provides context for creating migration steps
- Suggests prioritization based on dependencies

**Use with AI**:
```bash
# Get the plan context
go run main.go plan ./old-app authentication > context.txt

# Send to Claude with your planning prompt
cat context.txt | claude "Based on this legacy authentication code, create a detailed migration plan to modern Go best practices. Include: 1) Priority order 2) Specific refactoring steps 3) Testing strategy 4) Rollback plan"
```

### Phase 3: Migrate

Get focused context for implementing changes:

```bash
go run main.go migrate ./old-app authentication
```

**What it does**:
- Extracts code with highest relevance to component
- Limits token count for focused AI context
- Includes related test files
- Provides implementation-ready context

**Integration with AI workflow**:
```go
// Your migration automation script
result, _ := promptext.Extract("./old-app",
    promptext.WithRelevance("authentication"),
    promptext.WithTokenBudget(15000),
)

// Send to AI for code generation
modernCode := aiProvider.Generate(
    "Modernize this authentication code: " + result.FormattedOutput,
)

// Write modernized code
ioutil.WriteFile("auth_v2.go", []byte(modernCode), 0644)
```

### Phase 4: Verify

Validate the migration results:

```bash
go run main.go verify ./old-app authentication
```

**What it does**:
- Extracts comprehensive context of migrated code
- Includes test files for validation
- Provides context for review and testing
- Helps identify missed edge cases

## Real-World Migration Workflow

### Example: Migrating Legacy API Server

```bash
# 1. Initial analysis
go run main.go analyze ./legacy-api server > analysis.txt

# Review with AI to understand scope
cat analysis.txt | claude "Analyze this legacy server code. What are the main risks and challenges for migration?"

# 2. Create migration plan
go run main.go plan ./legacy-api server > plan-context.txt

cat plan-context.txt | claude "Create a phased migration plan with the following priorities: 1) Security issues 2) Deprecated APIs 3) Code organization 4) Performance"

# 3. Migrate component by component
go run main.go migrate ./legacy-api "http handlers" > handlers-context.txt

cat handlers-context.txt | claude "Modernize these HTTP handlers using the standard library's net/http patterns, proper error handling, and context propagation"

# 4. Verify each migration
go run main.go verify ./new-api server > verify-context.txt

cat verify-context.txt | claude "Review this migrated code. Check for: 1) Maintained functionality 2) Proper error handling 3) Missing edge cases 4) Test coverage"
```

## Integration Examples

### 1. Automated Migration Pipeline

```go
package main

import (
    "fmt"
    "os"
    "github.com/1broseidon/promptext/pkg/promptext"
)

func migrateComponent(legacyPath, component string) error {
    // Phase 1: Analyze
    fmt.Printf("Analyzing %s...\n", component)
    analysis, err := promptext.Extract(legacyPath,
        promptext.WithRelevance(component),
        promptext.WithTokenBudget(20000),
        promptext.WithFormat(promptext.FormatPTX),
    )
    if err != nil {
        return err
    }

    issues := detectIssues(analysis)
    fmt.Printf("Found %d issues\n", len(issues))

    // Phase 2: Plan
    fmt.Println("Creating migration plan...")
    planContext, _ := promptext.Extract(legacyPath,
        promptext.WithRelevance(component, "main", "init"),
        promptext.WithTokenBudget(10000),
    )

    plan := aiProvider.CreatePlan(planContext.FormattedOutput)

    // Phase 3: Migrate
    for _, step := range plan.Steps {
        fmt.Printf("Migrating: %s\n", step.Description)

        migrateContext, _ := promptext.Extract(legacyPath,
            promptext.WithRelevance(step.Keywords...),
            promptext.WithTokenBudget(15000),
        )

        newCode := aiProvider.Modernize(migrateContext.FormattedOutput, step)
        writeNewCode(step.TargetFile, newCode)
    }

    // Phase 4: Verify
    fmt.Println("Verifying migration...")
    verifyContext, _ := promptext.Extract("./new-code",
        promptext.WithRelevance(component),
        promptext.WithExtensions(".go", "_test.go"),
    )

    results := runTests(verifyContext)
    return reportResults(results)
}
```

### 2. Interactive Migration Assistant

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "github.com/1broseidon/promptext/pkg/promptext"
)

func interactiveMigration() {
    scanner := bufio.NewScanner(os.Stdin)

    fmt.Println("=== Interactive Migration Assistant ===")
    fmt.Print("Legacy codebase path: ")
    scanner.Scan()
    legacyPath := scanner.Text()

    fmt.Print("Component to migrate: ")
    scanner.Scan()
    component := scanner.Text()

    // Start with analysis
    fmt.Println("\n1. Analyzing component...")
    result, _ := promptext.Extract(legacyPath,
        promptext.WithRelevance(component),
        promptext.WithTokenBudget(20000),
    )

    fmt.Printf("\nFound %d relevant files (%d tokens)\n",
        len(result.ProjectOutput.Files), result.TokenCount)

    fmt.Println("\nFiles to migrate:")
    for i, file := range result.ProjectOutput.Files {
        fmt.Printf("  %d. %s (%d tokens)\n", i+1, file.Path, file.Tokens)
    }

    fmt.Print("\nProceed to AI analysis? (y/n): ")
    scanner.Scan()
    if strings.ToLower(scanner.Text()) == "y" {
        // Send to AI
        analysis := aiProvider.Analyze(result.FormattedOutput)
        fmt.Println("\n=== AI Analysis ===")
        fmt.Println(analysis)

        // Continue with planning, migration, verification...
    }
}
```

### 3. CI/CD Migration Pipeline

```yaml
# .github/workflows/migration.yml
name: Incremental Migration

on:
  workflow_dispatch:
    inputs:
      component:
        description: 'Component to migrate'
        required: true
      phase:
        description: 'Migration phase (analyze/plan/migrate/verify)'
        required: true
        default: 'analyze'

jobs:
  migrate:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'

      - name: Install promptext
        run: go install github.com/1broseidon/promptext/cmd/promptext@latest

      - name: Run Migration Phase
        run: |
          cd examples/migration-assistant
          go run main.go ${{ github.event.inputs.phase }} \
            ../../legacy-code \
            "${{ github.event.inputs.component }}" > output.txt

      - name: Upload Context
        uses: actions/upload-artifact@v3
        with:
          name: migration-context
          path: output.txt

      # Optional: Send to AI service for automated migration
      - name: AI-Assisted Migration
        if: github.event.inputs.phase == 'migrate'
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          # Send context to Claude for modernization
          curl https://api.anthropic.com/v1/messages \
            -H "x-api-key: $ANTHROPIC_API_KEY" \
            -H "content-type: application/json" \
            -d @migration-request.json
```

## AI Provider Integration

### Claude (Anthropic)

```bash
# Analyze phase
go run main.go analyze ./legacy-app auth | \
  claude --model claude-3-opus-20240229 \
  "Analyze this legacy authentication code. Identify security issues, deprecated patterns, and modernization opportunities."

# Plan phase
go run main.go plan ./legacy-app auth | \
  claude "Create a detailed migration plan with prioritized steps, testing strategy, and rollback procedures."

# Migrate phase
go run main.go migrate ./legacy-app auth | \
  claude "Modernize this code using current Go best practices. Maintain functionality while improving security and maintainability."
```

### OpenAI GPT

```python
import openai
import subprocess

# Get migration context
result = subprocess.run(
    ["go", "run", "main.go", "analyze", "./legacy-app", "auth"],
    capture_output=True,
    text=True
)

# Send to GPT-4
response = openai.ChatCompletion.create(
    model="gpt-4-turbo-preview",
    messages=[{
        "role": "user",
        "content": f"Analyze this legacy code and create a migration plan:\n\n{result.stdout}"
    }]
)

print(response.choices[0].message.content)
```

### Local LLM (Ollama)

```bash
# Use local LLM for privacy-sensitive migrations
go run main.go migrate ./legacy-app auth | \
  ollama run codellama "Refactor this code to use modern patterns while maintaining backward compatibility."
```

## Best Practices

### 1. Start Small

Begin with low-risk, isolated components:
```bash
# Good first migration targets
go run main.go analyze ./legacy utils
go run main.go analyze ./legacy helpers
go run main.go analyze ./legacy config
```

### 2. Use Token Budgets Wisely

Different phases need different context sizes:
- **Analyze**: Large budget (20k+) for comprehensive view
- **Plan**: Medium budget (10k) focusing on entry points
- **Migrate**: Focused budget (15k) for specific components
- **Verify**: Large budget (25k+) including tests

### 3. Validate Each Step

Always run tests after migration:
```bash
# After migrating each component
go test ./new-code/...
go run main.go verify ./new-code auth
```

### 4. Track Dependencies

Pay attention to the dependency list in analysis:
```bash
go run main.go analyze ./legacy auth | grep "Dependencies:"
# Plan migration order based on dependency graph
```

### 5. Maintain Parallel Versions

Keep old code until verification is complete:
```
project/
├── legacy/        # Original code
├── v2/            # Migrated code
└── tests/         # Shared test suite
```

## Common Migration Patterns

### Pattern 1: Security Updates

```bash
# Find insecure patterns
go run main.go analyze ./legacy auth | grep "security:"

# Focus migration on security issues
go run main.go migrate ./legacy auth > context.txt
cat context.txt | claude "Update this code to fix all security issues: use bcrypt for passwords, prevent SQL injection, add input validation"
```

### Pattern 2: Modernize Error Handling

```bash
go run main.go migrate ./legacy errors | \
  claude "Convert this code to use modern Go error handling: wrap errors with context, use errors.Is/As, add proper error types"
```

### Pattern 3: Add Context Support

```bash
go run main.go migrate ./legacy http | \
  claude "Refactor these HTTP handlers to accept and propagate context.Context for cancellation and timeouts"
```

### Pattern 4: Remove Global State

```bash
go run main.go analyze ./legacy globals | \
  claude "Identify all global variables and suggest dependency injection patterns to eliminate them"
```

## Tips for Better Migrations

1. **Use specific component names**: Instead of "code", use "authentication", "database", "api-handlers"

2. **Combine with code search**: Find all usages before migrating
   ```bash
   grep -r "oldFunction" ./legacy
   go run main.go migrate ./legacy oldFunction
   ```

3. **Review AI suggestions**: Always validate AI-generated code against your requirements

4. **Keep a migration log**: Track what's been migrated and what remains

5. **Test incrementally**: Don't migrate everything before testing

6. **Use version control**: Create branches for each major component migration

## Troubleshooting

### Issue: Too many files included

**Solution**: Be more specific with component names
```bash
# Instead of:
go run main.go analyze ./legacy api

# Use:
go run main.go analyze ./legacy "user api endpoints"
```

### Issue: Missing important files

**Solution**: Adjust relevance keywords
```bash
# Add related terms
go run main.go migrate ./legacy "auth authentication security login"
```

### Issue: Token limit exceeded

**Solution**: Migrate smaller chunks
```bash
# Split into sub-components
go run main.go migrate ./legacy "auth handlers"
go run main.go migrate ./legacy "auth models"
go run main.go migrate ./legacy "auth middleware"
```

## Advanced Usage

### Custom Migration Reports

Modify the tool to generate custom reports for your organization's needs:
```go
// Add custom issue detectors
func detectCustomIssues(files []File) []Issue {
    // Your organization's specific patterns
    // Check for deprecated internal libraries
    // Verify compliance with coding standards
    // Detect performance anti-patterns
}
```

### Integration with Project Management

Generate migration tasks automatically:
```go
// Export to JIRA, GitHub Issues, etc.
analysis := runAnalysis("./legacy", "auth")
for _, issue := range analysis.Issues {
    createJiraTicket(issue)
}
```

## Example: Full Migration Workflow

```bash
# 1. Initial assessment
go run main.go analyze ./legacy-monolith all > assessment.txt
cat assessment.txt  # Review scope

# 2. Prioritize components
components=("config" "utils" "models" "api" "auth")

# 3. Migrate each component
for component in "${components[@]}"; do
    echo "=== Migrating $component ==="

    # Analyze
    go run main.go analyze ./legacy-monolith "$component" > "analysis-$component.txt"

    # Plan
    go run main.go plan ./legacy-monolith "$component" | \
        claude "Create migration plan" > "plan-$component.md"

    # Migrate
    go run main.go migrate ./legacy-monolith "$component" | \
        claude "Modernize this code" > "../new-monolith/${component}.go"

    # Verify
    go test ../new-monolith/...
    go run main.go verify ../new-monolith "$component"
done

# 4. Final verification
go test ../new-monolith/... -cover
```

## Contributing

Found a useful migration pattern? Share it! This tool improves with real-world usage examples.

## License

MIT License - see main promptext repository for details.
