---
title: Introduction
description: Smart code context extractor for AI assistants
---

# promptext

Smart code context extraction for AI assistants. Process codebases efficiently with accurate token counting, intelligent filtering, and token-optimized output.

## Features

- **TOON Format** — Token-optimized output (30-60% smaller than JSON/Markdown), inspired by [johannschopplich/toon](https://github.com/johannschopplich/toon)
- **Relevance Filtering** — Multi-factor scoring prioritizes files by keywords
- **Token Budget Management** — Limit output to stay within AI model context windows
- **Smart Filtering** — .gitignore integration with intelligent defaults
- **Token Analysis** — Accurate GPT-compatible counting with tiktoken
- **Project Detection** — Language and framework identification
- **Multiple Formats** — TOON (default), Markdown, and XML with auto-detection
- **Performance Focused** — Optimized for large codebases

## Quick Start

```bash
# Install
go install github.com/1broseidon/promptext/cmd/promptext@latest

# Extract current directory (TOON format to clipboard)
promptext

# Or use the convenient alias
prx

# Prioritize authentication files, limit to 8k tokens
prx -r "auth login" --max-tokens 8000

# Auto-detect format from extension
prx -o context.toon
```

## Example Output

When using token budgets, promptext shows exactly what was included:

```
╭───────────────────────────────────────────────╮
│ 📦 promptext (Go)                             │
│    Included: 7/18 files • ~4,847 tokens       │
│    Full project: 18 files • ~19,512 tokens    │
╰───────────────────────────────────────────────╯

⚠️  Excluded 11 files due to token budget:
    • internal/cli/commands.go (~784 tokens)
    • internal/app/app.go (~60 tokens)
    ... and 9 more files (~8,453 tokens)
    Total excluded: ~9,297 tokens
```

Continue with [Getting Started](getting-started) for detailed setup and usage.
