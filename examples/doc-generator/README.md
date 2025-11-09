# Doc Generator - Automated Documentation

Keep your documentation in sync with code changes using AI-powered doc generation.

## What This Does

Automatically generates and updates documentation by:

1. ✅ Extracting relevant code using promptext
2. ✅ Creating structured prompts for different doc types
3. ✅ Generating ready-to-use documentation templates
4. ✅ Maintaining consistency across all docs

## Documentation Types

### API Reference
Complete API documentation with:
- Type definitions
- Function signatures
- Parameters and return values
- Usage examples

### README
Project overview including:
- Feature list
- Installation instructions
- Quick start guide
- Usage examples

### User Guide
Comprehensive tutorials covering:
- Core concepts
- Common tasks
- Advanced features
- Best practices

### Examples
Practical, runnable code showing:
- Basic usage
- Real-world scenarios
- Integration patterns

## Quick Start

```bash
cd examples/doc-generator

# Generate API reference
go run main.go --type api --output docs/api-reference.md

# Update README
go run main.go --type readme --output README.md

# Generate user guide
go run main.go --type guide --output docs/user-guide.md

# Generate all documentation
go run main.go --all
```

## Example Output

```
📚 Documentation Generator
============================================================
Generating api documentation...

✅ Extracted 12 files (8,456 tokens)

💾 Files saved:
   - Prompt: docs/api-prompt.txt
   - Context: docs/api-context.ptx
   - Output template: docs/api.md

💡 Next Steps:
   1. Review the generated prompt and context
   2. Send to your AI assistant (Claude, GPT, etc.)
   3. AI will generate documentation based on the code
   4. Review and commit the generated documentation

🎯 Documentation generation complete!
```

## How It Works

### 1. Code Extraction

Different doc types need different code:

**API Reference:**
```go
// Extract all public packages, exclude tests
promptext.Extract(".",
    promptext.WithExtensions(".go"),
    promptext.WithExcludes("*_test.go", "internal/"),
    promptext.WithTokenBudget(20000),
)
```

**README:**
```go
// Focus on main entry points and examples
promptext.Extract(".",
    promptext.WithExtensions(".go", ".md"),
    promptext.WithRelevance("main", "cmd", "example"),
    promptext.WithTokenBudget(10000),
)
```

### 2. Prompt Generation

Creates structured prompts with:
- Clear task description
- Specific requirements
- Output format guidelines
- Project context

### 3. Template Creation

Generates markdown templates with:
- Standard sections
- Placeholders for AI content
- Consistent formatting
- Navigation aids

### 4. AI Integration

The tool prepares everything needed:
- `*-prompt.txt` - What to ask the AI
- `*-context.ptx` - Code to analyze
- `*.md` - Template to fill in

## Automated Workflows

### Nightly Doc Updates

`.github/workflows/update-docs.yml`:

```yaml
name: Update Documentation

on:
  schedule:
    - cron: '0 2 * * *'  # Run at 2 AM daily
  workflow_dispatch:     # Allow manual triggers

jobs:
  update-docs:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Generate Documentation
        run: |
          cd examples/doc-generator
          go run main.go --all

      - name: Create Pull Request
        uses: peter-evans/create-pull-request@v5
        with:
          commit-message: 'docs: auto-update documentation'
          title: 'Automated Documentation Update'
          body: |
            Auto-generated documentation update based on latest code.

            Files updated:
            - docs/api.md
            - docs/guide.md
            - docs/example.md
```

### Pre-Commit Hook

Keep docs in sync on every commit:

`.git/hooks/pre-commit`:

```bash
#!/bin/bash

# Generate docs before commit
cd examples/doc-generator
go run main.go --type readme

# Check if README changed
if git diff --quiet README.md; then
  echo "✅ README is up to date"
else
  echo "📝 README updated, please review"
  git add README.md
fi
```

## Configuration

### Token Budgets

Adjust based on project size:

```go
// Small project
promptext.WithTokenBudget(5000)

// Medium project
promptext.WithTokenBudget(15000)

// Large project
promptext.WithTokenBudget(30000)
```

### File Filters

Customize what to include:

```go
// Just Go code
promptext.WithExtensions(".go")

// Multiple languages
promptext.WithExtensions(".go", ".js", ".ts", ".py")

// Exclude patterns
promptext.WithExcludes("vendor/", "*_test.go", "internal/")
```

### Source Directory

Generate docs for specific modules:

```bash
# Document entire project
go run main.go --source . --type api

# Document specific package
go run main.go --source pkg/mylib --type api

# Document examples only
go run main.go --source examples/ --type example
```

## AI Provider Integration

### Send to Claude

```bash
# Generate documentation
go run main.go --type api

# Send to Claude (requires anthropic CLI)
claude \
  --prompt "$(cat docs/api-prompt.txt)" \
  --file docs/api-context.ptx \
  > docs/api-reference.md
```

### Send to GPT

```bash
# Generate documentation
go run main.go --type readme

# Send to GPT (using openai CLI)
openai api chat.completions.create \
  -m gpt-4-turbo-preview \
  --messages "[{\"role\": \"user\", \"content\": \"$(cat README-prompt.txt)\n\n$(cat README-context.ptx)\"}]" \
  | jq -r '.choices[0].message.content' \
  > README.md
```

## Use Cases

### 1. Open Source Projects

Keep README and docs current:
```bash
# Before each release
go run main.go --type readme
go run main.go --type api
# Review and commit
git add docs/ README.md
git commit -m "docs: update for v1.2.0"
```

### 2. Internal Libraries

Maintain internal docs:
```bash
# Daily cron job
0 9 * * * cd /workspace && \
  go run examples/doc-generator/main.go --all && \
  git add docs/ && \
  git commit -m "docs: auto-update"
```

### 3. API Documentation

Generate API docs on deployment:
```bash
# In CI/CD pipeline
go run main.go --type api --output api-docs/reference.md
# Deploy to docs site
cp api-docs/* /var/www/docs/
```

### 4. Migration Documentation

Document changes during migrations:
```bash
# Before migration
go run main.go --type api --output docs/api-v1.md

# After migration
go run main.go --type api --output docs/api-v2.md

# Generate migration guide
diff docs/api-v1.md docs/api-v2.md > docs/migration-guide.md
```

## Output Files

### Generated Prompts (`*-prompt.txt`)

Contains:
- Task description
- Specific requirements
- Output format
- Project context

Example:
```markdown
# Documentation Generation Request: api

## Task
Generate comprehensive API reference documentation...

## Requirements
1. Document all public types...
2. Include parameter descriptions...
...
```

### Extracted Context (`*-context.ptx`)

Contains:
- Relevant source code
- Project structure
- Dependencies
- Metadata

Format: PTX (optimized for AI consumption)

### Documentation Templates (`*.md`)

Contains:
- Standard markdown structure
- Section placeholders
- Navigation
- Metadata

## Best Practices

### 1. Regular Updates

Keep docs fresh:
- Daily automated updates for active projects
- Weekly for stable projects
- On every release

### 2. Human Review

Always review AI-generated docs:
- Check technical accuracy
- Verify examples work
- Ensure tone is appropriate
- Fix formatting issues

### 3. Version Control

Track doc changes:
```bash
# Commit separately from code
git add docs/
git commit -m "docs: update API reference"
```

### 4. Incremental Updates

Update specific sections:
```bash
# Just update examples
go run main.go --type example

# Just update README
go run main.go --type readme
```

## Tips for Better Documentation

### 1. Use Clear Code Comments

AI generates better docs from commented code:
```go
// Good: Detailed comment
// ExtractContext analyzes the codebase and extracts
// relevant files based on the provided options.
// It returns a Result containing both structured data
// and formatted output ready for AI consumption.
func ExtractContext(dir string, opts ...Option) (*Result, error)

// Bad: Minimal comment
// ExtractContext extracts context
func ExtractContext(dir string, opts ...Option) (*Result, error)
```

### 2. Include Examples in Code

Example functions get documented:
```go
func ExampleExtract() {
    result, _ := Extract(".")
    fmt.Println(result.TokenCount)
    // Output: 1234
}
```

### 3. Maintain CHANGELOG

AI can incorporate recent changes:
```bash
# Include CHANGELOG in context
go run main.go --source . --type readme

# AI will reference recent changes
```

### 4. Organize Code Well

Clear structure = better docs:
```
pkg/
  api/           # API surface
  internal/      # Implementation details
  examples/      # Usage examples
docs/
  api.md         # Generated API docs
  guide.md       # Generated guide
```

## Advanced Usage

### Custom Prompts

Edit generated prompts before sending to AI:

```bash
# Generate prompt
go run main.go --type api

# Edit prompt to add specific instructions
vim docs/api-prompt.txt

# Send customized prompt to AI
cat docs/api-prompt.txt docs/api-context.ptx | claude
```

### Multiple Languages

Generate docs for polyglot projects:

```bash
# Go documentation
go run main.go --source go-code/ --type api -o docs/go-api.md

# JavaScript documentation
go run main.go --source js-code/ --type api -o docs/js-api.md
```

### Versioned Documentation

Keep docs for multiple versions:

```bash
# Current version
go run main.go --type api -o docs/v2/api.md

# Older version (switch branch first)
git checkout v1.0
go run main.go --type api -o docs/v1/api.md
git checkout main
```

## Troubleshooting

### Token Budget Exceeded

If project is too large:
```bash
# Document in parts
go run main.go --source pkg/core --type api
go run main.go --source pkg/utils --type api
```

### Missing Context

If docs seem incomplete:
```bash
# Increase token budget
# Edit main.go:
promptext.WithTokenBudget(30000)
```

### Outdated Docs

Set up automation:
```yaml
# GitHub Actions
on:
  push:
    branches: [main]
  schedule:
    - cron: '0 0 * * 0'  # Weekly
```

## Related Examples

- [CI Code Review](../ci-code-review/) - Automated PR analysis
- [Code Search](../code-search/) - Find code with natural language
- [Migration Assistant](../migration-assistant/) - Modernize legacy code
