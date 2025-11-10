# CI Code Review - Automated PR Analysis

Automatically extract code context from pull requests for AI-powered code reviews in your CI/CD pipeline.

## What This Does

This tool integrates with your CI system (GitHub Actions, GitLab CI, etc.) to:

1. ✅ Detect files changed in a pull request
2. ✅ Extract relevant code context using promptext
3. ✅ Generate review prompts for AI analysis
4. ✅ Prepare output ready to send to Claude/GPT/etc.

## Quick Start

### Local Testing

```bash
cd examples/ci-code-review

# Review current branch vs main
go run main.go

# Review specific PR
PR_NUMBER=123 PR_TITLE="Add authentication" go run main.go

# Custom base branch
BASE_BRANCH=develop go run main.go
```

### GitHub Actions Integration

See `.github/workflows/ai-review.yml` for a complete example.

## Example Output

```
🤖 AI Code Review - PR Context Extractor
============================================================
📋 PR #42: Add user authentication endpoints
   Author: alice
   Base: main → Head: feature/auth
   Changed files: 5

🔍 Extracting code context...
   Keywords from changes: [auth api handlers middleware users]

============================================================
✨ Context extracted successfully
   Files included: 8
   Token count: 12,456
============================================================

💾 Context saved to: pr-review-context.ptx
📊 Metadata saved to: pr-review-metadata.json
📝 Review prompt saved to: review-prompt.txt

💡 Next Steps:
   1. Review the extracted context in pr-review-context.ptx
   2. Use the review prompt in review-prompt.txt
   3. Send to your AI API (Claude, GPT, etc.)
   4. Post the review comments back to the PR

🎯 Review ready!
```

## How It Works

### 1. Changed File Detection

```go
// Uses git to find changes between base and head branch
git diff --name-only main...feature-branch
```

### 2. Keyword Extraction

Automatically extracts keywords from changed files:
- Directory names: `internal/auth` → `["auth", "internal"]`
- File names: `user_handler.go` → `["user", "handler"]`
- Helps find related code not directly changed

### 3. Context Extraction

Uses promptext with relevance filtering:
```go
promptext.Extract(cwd,
    promptext.WithRelevance(keywords...),
    promptext.WithTokenBudget(15000),
    promptext.WithFormat(promptext.FormatPTX),
)
```

### 4. Review Prompt Generation

Creates a structured prompt covering:
- Code quality and maintainability
- Security concerns
- Performance issues
- Test coverage
- Documentation

## GitHub Actions Setup

### 1. Create Workflow File

`.github/workflows/ai-review.yml`:

```yaml
name: AI Code Review

on:
  pull_request:
    types: [opened, synchronize]

jobs:
  review:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
        with:
          fetch-depth: 0  # Need full history for diff

      - uses: actions/setup-go@v4
        with:
          go-version: '1.22'

      - name: Extract PR Context
        env:
          PR_NUMBER: ${{ github.event.pull_request.number }}
          PR_TITLE: ${{ github.event.pull_request.title }}
          PR_AUTHOR: ${{ github.event.pull_request.user.login }}
          BASE_BRANCH: ${{ github.event.pull_request.base.ref }}
        run: |
          cd examples/ci-code-review
          go run main.go

      - name: Upload Context
        uses: actions/upload-artifact@v3
        with:
          name: review-context
          path: |
            pr-review-context.ptx
            review-prompt.txt
            pr-review-metadata.json
```

### 2. Add AI Review Step

Extend the workflow to call an AI API:

```yaml
      - name: AI Review
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          # Send context to Claude API
          curl https://api.anthropic.com/v1/messages \
            -H "x-api-key: $ANTHROPIC_API_KEY" \
            -H "content-type: application/json" \
            -d @- <<EOF
          {
            "model": "claude-3-5-sonnet-20241022",
            "max_tokens": 4096,
            "messages": [{
              "role": "user",
              "content": "$(cat review-prompt.txt)\n\n$(cat pr-review-context.ptx)"
            }]
          }
          EOF
```

### 3. Post Comment to PR

```yaml
      - name: Post Review Comment
        uses: actions/github-script@v7
        with:
          script: |
            const fs = require('fs');
            const review = fs.readFileSync('ai-review.txt', 'utf8');

            github.rest.issues.createComment({
              issue_number: context.issue.number,
              owner: context.repo.owner,
              repo: context.repo.repo,
              body: `## 🤖 AI Code Review\n\n${review}`
            });
```

## Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `PR_NUMBER` | Pull request number | `local` |
| `PR_TITLE` | Pull request title | Branch name |
| `PR_AUTHOR` | PR author username | `unknown` |
| `BASE_BRANCH` | Base branch for comparison | `main` |

### Token Budget

Adjust based on your AI model's context window:

```go
// GPT-4 Turbo (128k context)
promptext.WithTokenBudget(100000),

// Claude 3.5 (200k context)
promptext.WithTokenBudget(150000),

// GPT-3.5 (16k context)
promptext.WithTokenBudget(12000),
```

### File Extensions

Customize which files to review:

```go
promptext.WithExtensions(".go", ".js", ".ts", ".py"),
```

## AI Provider Integration

### Anthropic Claude

```go
import "github.com/anthropics/anthropic-sdk-go"

func sendToAnthropic(prompt, context string) (string, error) {
    client := anthropic.NewClient(
        option.WithAPIKey(os.Getenv("ANTHROPIC_API_KEY")),
    )

    response, err := client.Messages.New(ctx, anthropic.MessageNewParams{
        Model: anthropic.F(anthropic.ModelClaude3_5SonnetLatest),
        MaxTokens: anthropic.F(int64(4096)),
        Messages: anthropic.F([]anthropic.MessageParam{
            anthropic.NewUserMessage(anthropic.NewTextBlock(prompt + "\n\n" + context)),
        }),
    })

    return response.Content[0].Text, err
}
```

### OpenAI GPT

```go
import "github.com/sashabaranov/go-openai"

func sendToOpenAI(prompt, context string) (string, error) {
    client := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

    resp, err := client.CreateChatCompletion(
        context.Background(),
        openai.ChatCompletionRequest{
            Model: openai.GPT4Turbo,
            Messages: []openai.ChatCompletionMessage{
                {
                    Role: openai.ChatMessageRoleUser,
                    Content: prompt + "\n\n" + context,
                },
            },
        },
    )

    return resp.Choices[0].Message.Content, err
}
```

## Review Guidelines

The generated prompt asks the AI to review:

### 1. Code Quality
- Readability and maintainability
- Code smells and anti-patterns
- Naming conventions
- Code organization

### 2. Security
- Input validation
- Authentication/authorization
- Secrets management
- Common vulnerabilities (SQL injection, XSS, etc.)

### 3. Performance
- Algorithmic efficiency
- Database queries
- Network calls
- Resource usage

### 4. Testing
- Test coverage
- Edge cases
- Error scenarios
- Integration points

### 5. Documentation
- Code comments
- API documentation
- Complex logic explanations

## Output Files

### pr-review-context.ptx
Complete code context in PTX format, including:
- Changed files
- Related files (via relevance)
- Project structure
- Dependencies

### review-prompt.txt
Structured review prompt with:
- PR information
- Review guidelines
- Expected output format

### pr-review-metadata.json
Machine-readable metadata:
```json
{
  "pr_info": {
    "number": "42",
    "title": "Add authentication",
    "base_branch": "main",
    "head_branch": "feature/auth",
    "changed_files": ["..."],
    "author": "alice"
  },
  "token_count": 12456,
  "file_count": 8,
  "review_ready": true
}
```

## Tips for Better Reviews

### 1. Keep PRs Focused
Smaller PRs = better reviews:
- Easier for AI to analyze
- More focused feedback
- Faster review cycles

### 2. Use Descriptive Titles
Good titles help keyword extraction:
- ✅ "Add JWT authentication middleware"
- ❌ "Fix stuff"

### 3. Include Context
If PR needs special context, add it to the prompt:
```go
// Edit review-prompt.txt to add:
// "Note: This PR implements the authentication
//  design from RFC-123"
```

### 4. Review AI Output
Always human-review AI suggestions:
- AI can miss context
- May suggest overly complex solutions
- Domain knowledge matters

## Common Issues

### Too Many Files
If token budget exceeded:
- Increase budget for larger PRs
- Review in multiple passes
- Focus on most critical changes

### Missing Context
If AI review seems shallow:
- Check that related files were included
- Adjust relevance keywords
- Increase token budget

### False Positives
If AI reports non-issues:
- Refine review guidelines
- Add project-specific rules to prompt
- Use few-shot examples

## Advanced Usage

### Multi-Language Support

```go
// Detect languages in PR
languages := detectLanguages(prInfo.ChangedFiles)

// Custom extensions per language
for lang, files := range languages {
    result, _ := promptext.Extract(cwd,
        promptext.WithExtensions(getExtensions(lang)...),
        // ...
    )
}
```

### Incremental Reviews

Only review new commits:
```bash
# Compare against last review
BASE_COMMIT=$(cat .last-review-commit)
git diff --name-only $BASE_COMMIT...HEAD
```

### Custom Review Rules

Add project-specific rules:
```go
prompt += "\n## Project-Specific Rules\n"
prompt += "- All API handlers must have auth middleware\n"
prompt += "- Database queries must use prepared statements\n"
prompt += "- New features require integration tests\n"
```

## Cost Optimization

### Token Usage
- Default budget: ~15K tokens
- Average cost: $0.10-0.50 per review
- Reduce budget for frequent reviews

### Caching
Cache extracted context for quick re-reviews:
```go
// Save context with commit hash
contextKey := fmt.Sprintf("pr-%s-%s", prNumber, commitHash)
saveToCache(contextKey, result.Context)
```

### Selective Reviews
Only review certain file types or directories:
```go
// Only review src/ directory
promptext.Extract("src/",
    // ...
)
```

## Related Examples

- [Code Search](../code-search/) - Find code with natural language
- [Doc Generator](../doc-generator/) - Auto-generate documentation
- [Migration Assistant](../migration-assistant/) - Modernize legacy code
