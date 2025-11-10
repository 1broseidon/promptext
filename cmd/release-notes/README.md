# Release Notes Generator

Automatically generates release notes by analyzing git changes using the promptext library.

## Features

- Analyzes git commits since last release
- Extracts code context with promptext
- Categorizes changes (feat, fix, docs, chore)
- Detects breaking changes
- Provides statistics on changes

## Usage

```bash
# Generate notes since last tag
go run cmd/release-notes/main.go

# Generate notes for specific version
go run cmd/release-notes/main.go --version v0.7.4

# Generate notes since specific tag
go run cmd/release-notes/main.go --since v0.7.2

# Save to file
go run cmd/release-notes/main.go --output release-notes.md
```

## Example Output

```markdown
## [0.7.4] - 2025-11-10

### Added
- New feature X

### Fixed
- Bug fix Y

### Statistics
- **Files changed**: 15
- **Commits**: 7
- **Context analyzed**: ~7,187 tokens
```

## How It Works

1. **Git Analysis**: Compares current HEAD with last release tag
2. **Context Extraction**: Uses promptext library to extract code context from changed files
3. **Categorization**: Parses commit messages following conventional commits
4. **Formatting**: Generates changelog-compatible markdown

## Integration

Can be integrated into release workflow:

```yaml
# .github/workflows/release.yml
- name: Generate Release Notes
  run: go run cmd/release-notes/main.go --version ${{ github.ref_name }} --output RELEASE_NOTES.md
```

## Demonstrates

- ✅ Real-world promptext library usage
- ✅ Automated documentation workflow
- ✅ Git integration
- ✅ Token budget management
