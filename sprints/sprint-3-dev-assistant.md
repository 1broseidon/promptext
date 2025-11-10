# Sprint 3: Local Dev Assistant Script

## Goal
Create a local development assistant that developers can run before committing to check if changelog, documentation, or examples need updates based on their staged changes.

## What It Does
- Analyzes git staged changes
- Uses promptext to extract context of modified files
- Checks if:
  - Changelog needs an entry
  - Documentation needs updates
  - Examples need modifications
  - Tests are missing
- Provides interactive prompts
- Optionally auto-generates updates

## Implementation

### 1. Create `scripts/dev-assistant.sh`
```bash
#!/bin/bash
# Interactive dev assistant

echo "🤖 Promptext Dev Assistant"
echo ""

go run cmd/dev-assistant/main.go "$@"
```

### 2. Create `cmd/dev-assistant/main.go`
```go
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    fmt.Println("🔍 Analyzing your staged changes...")

    // Get staged files
    stagedFiles := getStagedFiles()
    if len(stagedFiles) == 0 {
        fmt.Println("No staged changes detected.")
        return
    }

    // Extract context with promptext
    result, err := extractStagedContext(stagedFiles)
    if err != nil {
        log.Fatal(err)
    }

    // Analyze what needs updating
    checks := analyzeChanges(result, stagedFiles)

    // Display results
    displayChecks(checks)

    // Prompt for actions
    if checks.NeedsChangelog {
        if confirm("Add changelog entry?") {
            addChangelogEntry()
        }
    }

    if len(checks.DocsToUpdate) > 0 {
        fmt.Println("\n📚 Documentation updates suggested:")
        for _, doc := range checks.DocsToUpdate {
            fmt.Printf("  - %s\n", doc)
        }
    }
}
```

### 3. Features to Implement
- [x] Detect staged changes
- [x] Extract context with promptext
- [x] Check if changelog entry needed
- [x] Detect documentation drift
- [x] Suggest example updates
- [x] Interactive prompts
- [x] Optional auto-apply changes

### 4. Usage Scenarios

**Basic check:**
```bash
./scripts/dev-assistant.sh
```

**Before commit:**
```bash
./scripts/dev-assistant.sh --pre-commit
```

**Auto-fix mode:**
```bash
./scripts/dev-assistant.sh --auto-fix
```

**Dry run:**
```bash
./scripts/dev-assistant.sh --dry-run
```

### 5. Interactive Flow
```
🤖 Promptext Dev Assistant

🔍 Analyzing staged changes...
  ✓ Found 3 modified files
  ✓ Extracted context (2,456 tokens)

📋 Analysis Results:
  ⚠️  Changelog entry recommended
      Reason: Added new public API method WithCustomFormatter()

  ⚠️  Documentation updates needed
      Files: docs-astro/src/content/docs/library-usage.md
      Reason: New API method not documented

  ✓ Examples are up to date

  ✓ Tests look good

Would you like to:
  1. Generate changelog entry
  2. Update documentation
  3. Skip for now

Choice [1-3]:
```

### 6. Integration Options

**Git pre-commit hook:**
```bash
#!/bin/bash
# .git/hooks/pre-commit
./scripts/dev-assistant.sh --pre-commit --auto-fix
```

**Manual workflow:**
```bash
# Your normal workflow
git add .
./scripts/dev-assistant.sh  # Check before commit
git commit -m "..."
```

## Success Criteria
- [x] Analyzes staged changes effectively
- [x] Uses promptext library for context extraction
- [x] Provides helpful suggestions
- [x] Interactive and user-friendly
- [x] Improves developer productivity
- [x] Can be integrated into git workflow

## Advanced Features (Optional)
- [ ] Suggest commit message based on changes
- [ ] Check for breaking changes
- [ ] Verify version bump needed
- [ ] Check if tests cover changes
- [ ] Integration with editor (VS Code extension)
