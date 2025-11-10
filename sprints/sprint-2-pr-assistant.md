# Sprint 2: AI-Powered PR Assistant

## Goal
Create a GitHub Actions workflow that analyzes PRs using promptext, provides AI-powered suggestions for changelog entries, documentation updates, and potential breaking changes.

## What It Does
- Triggers on pull request creation/update
- Uses promptext to extract changed files
- Analyzes changes with AI for:
  - Suggested changelog entry
  - Documentation that needs updating
  - Breaking change detection
  - Test coverage suggestions
- Posts suggestions as PR comment
- Updates on each PR push

## Implementation

### 1. Create `.github/workflows/pr-assistant.yml`
```yaml
name: PR Assistant

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  analyze:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0

      - uses: actions/setup-go@v5
        with:
          go-version: stable

      - name: Run PR Assistant
        run: go run cmd/pr-assistant/main.go
        env:
          GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
          PR_NUMBER: ${{ github.event.pull_request.number }}
```

### 2. Create `cmd/pr-assistant/main.go`
```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/1broseidon/promptext/pkg/promptext"
    "github.com/google/go-github/v57/github"
)

func main() {
    // Get PR changed files
    changedFiles := getPRChangedFiles()

    // Extract context with promptext
    result, err := promptext.Extract(".",
        promptext.WithExtensions(".go", ".md"),
        promptext.WithTokenBudget(8000),
    )

    // Analyze for suggestions
    suggestions := generatePRSuggestions(result, changedFiles)

    // Post comment to PR
    postPRComment(suggestions)
}
```

### 3. Features to Implement
- [x] Detect PR changed files
- [x] Extract context with promptext
- [x] Generate changelog suggestions
- [x] Detect documentation needs
- [x] Identify breaking changes
- [x] Post formatted comment to PR
- [x] Update comment on new pushes (not duplicate)

### 4. Comment Format
```markdown
## 🤖 PR Assistant Analysis

### 📝 Suggested Changelog Entry
```
### Added
- New feature X that does Y
```

### 📚 Documentation Updates Needed
- [ ] Update library-usage.md with new option
- [ ] Add example for feature X

### ⚠️ Potential Issues
- Breaking change detected in function signature
- Consider adding migration guide

---
*Powered by [promptext](https://github.com/1broseidon/promptext)*
```

### 5. Integration Points
- Runs on every PR
- Uses promptext library to analyze changes
- Provides actionable feedback
- Helps maintain documentation quality

## Success Criteria
- [x] Automatically analyzes PRs
- [x] Provides helpful suggestions
- [x] Uses promptext library effectively
- [x] Demonstrates AI-powered code review
- [x] Improves documentation quality

## Advanced Features (Optional)
- [ ] Detect if examples need updating
- [ ] Check if version needs bumping
- [ ] Suggest test cases based on changes
- [ ] Link to related issues/PRs
