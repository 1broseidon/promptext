---
title: Introduction
description: Convert your codebase into AI-ready prompts
---

# Welcome to promptext

A fast, token-efficient tool that transforms your code into optimized context for Claude, ChatGPT, and other LLMs.

## What is promptext?

`promptext` intelligently filters your codebase, ranks files by relevance, and packages them into token-efficient formats—all within your specified budget. It's designed to help you work more effectively with AI assistants by providing the right code context without exceeding token limits.

## Key Features

- 🚀 **Fast** — Written in Go, processes large codebases in seconds
- 🧠 **Smart** — Relevance scoring automatically prioritizes important files
- 💰 **Budget-Aware** — Enforces token limits to prevent context overflow
- 📦 **Token-Efficient Formats** — PTX (25-30% savings), TOON-strict (30-60% savings), Markdown, or XML
- 🎯 **Accurate Counting** — Uses `cl100k_base` tokenizer (GPT-4, GPT-3.5, Claude compatible)
- ⚙️ **Highly Configurable** — Project-level `.promptext.yml` and global settings support
- 🔒 **Multi-Layered Lock File Detection** — Automatically excludes lock files, saving 50K-100K+ tokens per project
- 🎨 **Multiple Output Formats** — PTX (default), TOON-strict, Markdown, and XML with auto-detection

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
