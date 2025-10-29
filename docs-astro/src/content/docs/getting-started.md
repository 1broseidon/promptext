---
title: Getting Started
description: Quick start guide for promptext installation and basic usage
---

## Installation

Choose your preferred installation method:

### Go Install (Recommended)

```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

### Script Install

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash
```

**Windows PowerShell:**
```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex
```

### Manual Install

Download binaries from [GitHub Releases](https://github.com/1broseidon/promptext/releases) and add to your PATH.

## Basic Usage

### Simple Commands

```bash
# Process current directory (TOON format to clipboard)
promptext

# Use alias for convenience
prx

# Process specific directory
promptext -d /path/to/project

# Show project overview only
promptext -i

# Export to file (format auto-detected from extension)
promptext -o context.toon
promptext -o context.md
promptext -o project.xml
```

### Common Options

| Flag | Description |
|------|-------------|
| `-d` | Directory to process |
| `-e` | File extensions (`.go,.js`) |
| `-x` | Exclude patterns |
| `-f` | Format (`toon`, `markdown`, `xml`) |
| `-o` | Output file (auto-detects format) |
| `-i` | Info mode only |
| `-r` | Relevant keywords for prioritization |
| `--max-tokens` | Token budget limit |
| `-v` | Verbose output |
| `-q` | Quiet mode for scripting |

### Examples

**Filter by file type:**
```bash
promptext -e .go,.js,.ts
```

**Exclude directories:**
```bash
promptext -x "node_modules/,vendor/,test/"
```

**Generate reports:**
```bash
# TOON format (default, token-optimized)
promptext -o context.toon

# Markdown format
promptext -f markdown -o context.md

# XML format for automation
promptext -f xml -o report.xml
```

**Prioritize relevant files:**
```bash
# Focus on authentication code
promptext -r "auth login OAuth"

# Database-related files
promptext -r "database SQL migration"
```

**Stay within token budgets:**
```bash
# Limit to 8000 tokens (Claude Haiku)
promptext --max-tokens 8000

# Combine with relevance for smart selection
promptext -r "api routes" --max-tokens 5000
```

## Quick Workflows

### For AI Queries

```bash
# Quick context (3k tokens)
prx -r "auth" --max-tokens 3000

# Standard context (8k tokens)
prx -r "api database" --max-tokens 8000

# Full codebase (within limits)
prx --max-tokens 50000
```

### For Documentation

```bash
# Export project overview
prx -i -o overview.md

# Export full context in markdown
prx -f markdown -o full-context.md
```

### For CI/CD

```bash
# Machine-readable XML
prx -f xml -o build/context.xml

# Quiet mode for scripting
prx -q -o context.toon
```

## Next Steps

- [Output Formats](output-formats) - TOON, Markdown, and XML formats
- [Relevance Filtering](relevance-filtering) - Smart file prioritization
- [Configuration](configuration) - Customize behavior
- [File Filtering](file-filtering) - Advanced filtering rules
