# AI PR Assistant

Automatically analyzes pull requests using promptext library to provide intelligent suggestions for changelog entries, documentation updates, and potential issues.

## Features

- **Changelog Detection**: Identifies when changelog needs updating
- **Documentation Suggestions**: Detects when docs need updates based on code changes
- **Breaking Change Detection**: Identifies potential breaking changes
- **Example Updates**: Suggests when examples should be updated
- **Statistics**: Provides analysis metrics

## How It Works

1. **Extract Changes**: Gets list of changed files in PR
2. **Context Analysis**: Uses promptext to extract code context (token-aware)
3. **Pattern Detection**: Analyzes patterns in changes
4. **Smart Suggestions**: Generates actionable recommendations
5. **Comment Generation**: Creates formatted PR comment

## Usage

### In GitHub Actions

Automatically runs on every PR via `.github/workflows/pr-assistant.yml`:

```yaml
on:
  pull_request:
    types: [opened, synchronize, reopened]
```

### Local Testing

```bash
# Test with current branch vs main
PR_NUMBER=test go run cmd/pr-assistant/main.go

# Test with specific commits
BASE_SHA=abc123 HEAD_SHA=def456 PR_NUMBER=1 go run cmd/pr-assistant/main.go
```

## Example Output

```markdown
## 🤖 AI PR Assistant Analysis

### 📝 Changelog Update Needed

This PR modifies library code but doesn't update CHANGELOG.md

**Suggested entry:**
```markdown
### Added
- New WithCustomFormatter option for library
```

### 📚 Documentation Updates Suggested

- [ ] Update docs-astro/src/content/docs/library-usage.md
- [ ] Add example in examples/ directory

### 📊 Analysis Statistics

- **Files analyzed**: 5
- **Context extracted**: ~4,234 tokens
- **Project language**: Go

---
*Powered by [promptext](https://github.com/1broseidon/promptext) 🚀*
```

## Demonstrates

- ✅ Automated code review with AI
- ✅ Real-time documentation drift detection
- ✅ Library usage in CI/CD pipeline
- ✅ Token-aware context extraction
- ✅ Pattern recognition in code changes

## Future Enhancements

- [ ] Integrate GitHub API to post comments
- [ ] Detect API surface changes
- [ ] Suggest specific doc sections to update
- [ ] Test coverage analysis
- [ ] Performance impact detection
