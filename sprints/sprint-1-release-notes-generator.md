# Sprint 1: Smart Release Notes Generator

## Goal
Automate release note generation by using promptext to analyze code changes between releases, then AI to generate comprehensive, well-formatted release notes.

## What It Does
- Analyzes git commits since last release tag
- Uses promptext library to extract changed files
- Categorizes changes (features, fixes, breaking changes)
- Generates formatted release notes
- Updates CHANGELOG.md and docs changelog automatically

## Implementation

### 1. Create `cmd/release-notes/main.go`
```go
package main

import (
    "fmt"
    "log"
    "os"
    "os/exec"
    "strings"
    "time"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    // Get last release tag
    lastTag := getLastTag()

    // Get changed files since last tag
    changedFiles := getChangedFilesSinceTag(lastTag)

    // Extract code context for changed files
    result, err := extractChangedFilesContext(changedFiles)
    if err != nil {
        log.Fatal(err)
    }

    // Generate release notes using AI context
    releaseNotes := generateReleaseNotes(result, lastTag)

    fmt.Println(releaseNotes)
}
```

### 2. Features to Implement
- [x] Get last git tag
- [x] List changed files between tags
- [x] Extract context with promptext for only changed files
- [x] Format output for AI processing
- [x] Generate categorized release notes
- [x] Update CHANGELOG.md
- [x] Update docs-astro changelog

### 3. Usage
```bash
# Generate release notes for current version
go run cmd/release-notes/main.go

# Or with specific version
go run cmd/release-notes/main.go --version v0.7.3
```

### 4. Integration Points
- Can be run manually before tagging
- Can be added to release workflow
- Can generate draft release notes for GitHub

## Success Criteria
- [x] Successfully extracts changes since last release
- [x] Uses promptext library to analyze changed code
- [x] Generates well-formatted release notes
- [x] Demonstrates real-world library usage
- [x] Can be used for future releases

## Example Output
```markdown
## [0.7.3] - 2025-11-09

### Added
- Comprehensive default exclusions for 10+ ecosystems
- 70+ new exclusion patterns

### Fixed
- Critical: Python virtual environments included in output

### Changed
- Organized exclusions by ecosystem
```
