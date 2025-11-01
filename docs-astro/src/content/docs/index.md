---
title: Introduction
description: Smart code context extractor for AI assistants
---

# promptext

Smart code context extraction for AI assistants. Process codebases efficiently with accurate token counting, intelligent filtering, and token-optimized output.

## Features

- **Multi-Layered Lock File Detection** — Automatically excludes package lock files (99% signature-based, 95% ecosystem-aware, 85% heuristic) — saves 50K-100K+ tokens per project
- **PTX Format** — Default hybrid format combining TOON v1.3 metadata efficiency with readable multiline code blocks (25-30% smaller than JSON)
- **TOON-Strict Mode** — Full TOON v1.3 compliance for maximum compression (30-60% smaller than JSON), based on [johannschopplich/toon](https://github.com/johannschopplich/toon)
- **Relevance Filtering** — Multi-factor scoring prioritizes files by keywords
- **Token Budget Management** — Limit output to stay within AI model context windows
- **Smart Filtering** — .gitignore integration, intelligent defaults, and generated file detection
- **Token Analysis** — Accurate counting using tiktoken `cl100k_base` encoding (GPT-4/GPT-3.5-turbo) with intelligent fallback
- **Project Detection** — Language and framework identification with ecosystem-aware filtering
- **Multiple Formats** — PTX (default), TOON-strict, Markdown, and XML with auto-detection
- **Performance Focused** — Optimized for large codebases

## Quick Start

```bash
# Install
go install github.com/1broseidon/promptext/cmd/promptext@latest

# Extract current directory (PTX format to clipboard)
promptext

# Or use the convenient alias
prx

# Prioritize authentication files, limit to 8k tokens
prx -r "auth login" --max-tokens 8000

# Auto-detect format from extension
prx -o context.ptx    # PTX format (default, readable code)
prx -o context.toon   # PTX format (backward compatibility)
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
