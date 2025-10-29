---
title: Relevance Filtering & Token Budgets
description: Smart file prioritization and token budget management for AI models
---

## Overview

Relevance filtering intelligently prioritizes files based on keyword matching, allowing you to focus on the most important code for your task. Combined with token budgets, this ensures your context stays within AI model limits while including the most relevant information.

## Relevance Filtering

### Basic Usage

Prioritize files matching specific keywords:

```bash
# Prioritize authentication-related files
promptext --relevant "auth login OAuth"

# Short form
promptext -r "database SQL postgres"

# Multiple keywords (comma or space separated)
promptext -r "api,routes,handlers"
```

### Multi-Factor Scoring

Promptext uses multi-factor scoring to determine file relevance:

| Match Location | Weight | Example |
|---------------|--------|---------|
| Filename | 10x | `auth.go` matches "auth" |
| Directory | 5x | `internal/auth/` matches "auth" |
| Imports | 3x | `import "auth"` matches "auth" |
| Content | 1x | `func authenticate()` matches "auth" |

**Scoring Example:**
```
File: internal/auth/handler.go
Keywords: "auth"

Matches:
- Filename "handler.go" ❌ (0 points)
- Directory "internal/auth" ✓ (5 points)
- Import "github.com/pkg/auth" ✓ (3 points)
- Content "authentication" appears 4x ✓ (4 points)

Total Score: 12 points
```

### Prioritization Strategy

Files are sorted by priority:

1. **Entry points with high relevance** - `main.go`, `index.js`, etc. matching keywords
2. **High relevance files** - Score above threshold (5 points)
3. **Shallow files first** - Prefer root-level over deeply nested
4. **Config files** - Configuration before implementation
5. **Tests last** - Test files have lowest priority

```bash
# With --relevant flag, files are prioritized:
promptext -r "database" -e .go

# Priority order:
# 1. cmd/main.go (entry point, mentions database)
# 2. internal/database/conn.go (filename + directory match)
# 3. internal/database/queries.go (directory match)
# 4. internal/config/database.go (config + filename match)
# 5. pkg/models/user.go (imports database package)
# 6. internal/database/conn_test.go (test file)
```

## Token Budget Management

### Basic Usage

Limit output to stay within token constraints:

```bash
# Limit to 8000 tokens (Claude Haiku sweet spot)
promptext --max-tokens 8000

# Combine with relevance for smart prioritization
promptext -r "api routes" --max-tokens 5000

# Quick context for simple queries
promptext --max-tokens 3000
```

### Budget Allocation

When a token budget is set:

1. **Calculate overhead** - Directory tree, git info, metadata (~500-2000 tokens)
2. **Available budget** - Total budget minus overhead
3. **Include files** - Add highest-priority files until budget exhausted
4. **Report exclusions** - Show what was excluded and why

```
╭───────────────────────────────────────────────╮
│ 📦 promptext (Go)                             │
│    Included: 7/18 files • ~4,847 tokens       │
│    Full project: 18 files • ~19,512 tokens    │
╰───────────────────────────────────────────────╯

⚠️  Excluded 11 files due to token budget:
    • internal/cli/commands.go (~784 tokens)
    • internal/app/app.go (~60 tokens)
    • internal/ollama/client.go (~211 tokens)
    • PROJECT_PLAN.md (~202 tokens)
    • README.md (~279 tokens)
    ... and 6 more files (~7,761 tokens)
    Total excluded: ~9,297 tokens
```

### Filtered Directory Tree

The directory structure automatically adjusts to show only included files:

**Without token budget:**
```yaml
structure:
  cmd/chain[1]: main.go
  internal/cli[5]: commands.go,flags.go,help.go,output.go,version.go
  internal/app[1]: app.go
```

**With token budget (5000 tokens):**
```yaml
structure:
  cmd/chain[1]: main.go
  internal/cli[2]: commands.go,help.go
```

## Use Cases

### AI Model Context Windows

Match your budget to AI model capabilities:

```bash
# Claude Haiku (100k context, but efficient at 8k)
promptext --max-tokens 8000

# Claude Sonnet/Opus (200k context)
promptext --max-tokens 100000

# GPT-3.5-Turbo (16k context)
promptext --max-tokens 12000

# GPT-4 (8k context, older models)
promptext --max-tokens 6000
```

### Cost Optimization

Fewer tokens = lower API costs:

```bash
# Quick question about architecture
promptext -r "architecture design" --max-tokens 3000

# Detailed API review
promptext -r "api endpoints routes" --max-tokens 10000

# Full codebase analysis (cost-aware)
promptext --max-tokens 50000
```

### Focused Context

Zero in on specific functionality:

```bash
# Authentication system only
promptext -r "auth login session OAuth" --max-tokens 8000

# Database layer only
promptext -r "database SQL migration schema" --max-tokens 6000

# API routes and handlers
promptext -r "routes handlers endpoints middleware" --max-tokens 7000
```

## Advanced Combinations

### Relevance + Budget + Format

```bash
# Token-optimized TOON format with relevance filtering
promptext -r "testing unittest" --max-tokens 5000 -o tests.toon

# Focused markdown for documentation
promptext -r "api public" --max-tokens 8000 -o api-docs.md

# XML for CI/CD with budget control
promptext -r "config deployment" --max-tokens 10000 -o deploy.xml
```

### File Type + Relevance + Budget

```bash
# Go files only, auth-related, within 6k tokens
promptext -e .go -r "auth security" --max-tokens 6000

# TypeScript API files, limited context
promptext -e .ts,.tsx -r "api fetch axios" --max-tokens 4000
```

### Exclude + Relevance + Budget

```bash
# Exclude tests, focus on core logic, limit tokens
promptext -x "test/,spec/" -r "business logic core" --max-tokens 8000

# Exclude vendors, focus on database, optimize for Claude Haiku
promptext -x "vendor/,node_modules/" -r "database" --max-tokens 8000
```

## Best Practices

### Keyword Selection

**Good keywords:**
- Specific: "authentication", "database", "api"
- Domain-relevant: "OAuth", "JWT", "PostgreSQL"
- Functional: "routes", "handlers", "middleware"

**Avoid:**
- Too broad: "code", "function", "class"
- Too specific: "line42", "temporary"
- Common words: "get", "set", "data"

### Budget Guidelines

| Task Type | Recommended Budget | Rationale |
|-----------|-------------------|-----------|
| Quick questions | 2,000-3,000 | Architecture docs + entry points |
| Feature review | 5,000-8,000 | Relevant subsystem + context |
| Full codebase | 15,000-50,000 | Complete understanding |
| Debug session | 6,000-10,000 | Error context + related code |

### Workflow Tips

**1. Start narrow:**
```bash
# Get focused context first
promptext -r "payment checkout" --max-tokens 5000
```

**2. Expand if needed:**
```bash
# Increase budget or broaden keywords
promptext -r "payment checkout stripe" --max-tokens 10000
```

**3. Use dry-run to preview:**
```bash
# See what would be included
promptext -r "auth" --max-tokens 5000 --dry-run
```

**4. Check exclusions:**
```bash
# Review what was excluded
promptext -r "api" --max-tokens 8000
# Read the exclusion summary to see if critical files were missed
```

## Configuration

Set defaults in `.promptext.yml`:

```yaml
# Not recommended: relevance keywords should be task-specific
# Better to specify on command line

# But token budget can be a good default:
# (commented out by default)
# max-tokens: 8000  # Default budget for this project
```

## Troubleshooting

**Missing important files?**
- Increase `--max-tokens` budget
- Add more relevant keywords with `-r`
- Check exclusion summary for clues

**Too many irrelevant files?**
- Use more specific keywords
- Reduce `--max-tokens` to force prioritization
- Add exclusion patterns with `-x`

**Context still too large?**
- Narrow keyword focus
- Decrease token budget
- Exclude test files: `-x "test/,spec/"`
- Use specific file extensions: `-e .go`
