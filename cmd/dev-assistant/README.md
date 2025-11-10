# Dev Assistant

Interactive local development tool that analyzes staged changes before committing to ensure changelog, documentation, and examples are up to date.

## Features

- **Changelog Detection**: Checks if staged changes need changelog entry
- **Documentation Checks**: Identifies when docs should be updated
- **Example Validation**: Suggests when examples need updates
- **Test Coverage**: Reminds about test requirements
- **Interactive Mode**: Prompts for actions
- **Token-Aware**: Uses promptext for intelligent analysis

## Usage

### Basic Check

```bash
# Stage your changes first
git add .

# Run the assistant
./scripts/dev-assistant.sh
```

### Dry Run Mode

```bash
# See what would be checked without prompting
./scripts/dev-assistant.sh --dry-run
```

### Auto-Fix Mode

```bash
# Attempt to automatically apply fixes
./scripts/dev-assistant.sh --auto-fix
```

## Example Session

```
🤖 Promptext Dev Assistant

🔍 Analyzing staged changes...

   Found 3 staged files
   Extracted ~2,456 tokens from 3 files

📋 Analysis Results:

  ⚠️  Changelog entry recommended
      Reason: Public API changes detected in pkg/promptext/

  ⚠️  Documentation updates needed
      - docs-astro/src/content/docs/library-usage.md - Library changes detected

  ✅ Examples are up to date
  ✅ Tests look good

Add changelog entry? [y/N]: y

💡 Add your changelog entry to CHANGELOG.md
   Then stage the file: git add CHANGELOG.md

✨ Done! Remember to review your changes before committing.
```

## Integration with Git

### Pre-Commit Hook

Add to `.git/hooks/pre-commit`:

```bash
#!/bin/bash
./scripts/dev-assistant.sh --dry-run
```

### Alias

Add to your `.bashrc` or `.zshrc`:

```bash
alias check='./scripts/dev-assistant.sh'
```

Then use:

```bash
git add .
check
git commit -m "..."
```

## How It Works

1. **Detect Staged Files**: Uses `git diff --cached` to get staged changes
2. **Extract Context**: Uses promptext library with 8K token budget
3. **Pattern Analysis**: Checks for library changes, doc updates, etc.
4. **Interactive Prompts**: Asks for confirmation before actions
5. **Helpful Suggestions**: Provides specific file paths and reasons

## Demonstrates

- ✅ Local development productivity
- ✅ Interactive CLI tools with promptext
- ✅ Git workflow integration
- ✅ Pre-commit validation
- ✅ Token-aware context extraction

## Checks Performed

- **Changelog**: Detects if CHANGELOG.md needs an entry
- **Documentation**: Identifies doc files needing updates
- **Examples**: Checks if examples/ should be updated
- **Tests**: Reminds about test coverage
- **Breaking Changes**: Warns about API changes

## Future Enhancements

- [ ] Auto-generate changelog entries
- [ ] Suggest specific doc sections to update
- [ ] Check test coverage percentage
- [ ] Integrate with editor (VS Code extension)
- [ ] AI-powered commit message suggestions
