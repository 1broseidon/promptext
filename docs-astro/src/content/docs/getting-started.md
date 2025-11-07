---
title: Getting Started
description: Quick start guide for promptext installation and basic usage
---

## Installation

**macOS/Linux:**
```bash
curl -sSL promptext.sh/install | bash
```

**Windows:**
```powershell
irm promptext.sh/install.ps1 | iex
```

**Go Install (requires Go 1.19+):**
```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

**Manual Install:**
Download binaries from [GitHub Releases](https://github.com/1broseidon/promptext/releases) and add to your PATH.

## Uninstalling

**Linux/macOS:**
```bash
curl -sSL promptext.sh/uninstall | bash
```

**Manual Removal:**
```bash
# Remove binary (default location)
rm ~/.local/bin/promptext

# Remove aliases from shell configs
sed -i '/alias prx=/d' ~/.bashrc ~/.zshrc
```

## Basic Usage

### Simple Commands

```bash
# Process current directory (PTX format to clipboard)
promptext

# Use alias for convenience
prx

# Process specific directory
prx /path/to/project

# Show project overview only
prx -i

# Export to file (format auto-detected from extension)
prx -o context.ptx   # PTX format (default)
prx -o context.md    # Markdown format
prx -o project.xml   # XML format
```

### Common Options

| Flag | Description |
|------|-------------|
| (path) | Directory to process (e.g., `prx /path/to/project`) |
| `-e` | File extensions (`.go,.js`) |
| `-x` | Exclude patterns |
| `-f` | Format (`ptx`, `toon-strict`, `markdown`, `xml`) |
| `-o` | Output file (auto-detects format from extension) |
| `-i` | Info mode only |
| `-r` | Relevant keywords for prioritization |
| `--max-tokens` | Token budget limit |
| `-v` | Verbose output |
| `-D` | Debug mode with timing |

### Examples

**Filter by file type:**
```bash
prx -e .go,.js,.ts
```

**Exclude directories:**
```bash
prx -x "node_modules/,vendor/,test/"
```

**Generate reports:**
```bash
# PTX format (default, token-optimized)
prx -o context.ptx

# Markdown format
prx -f markdown -o context.md

# XML format for automation
prx -f xml -o report.xml
```

**Prioritize relevant files:**
```bash
# Focus on authentication code
prx -r "auth login OAuth"

# Database-related files
prx -r "database SQL migration"
```

**Stay within token budgets:**
```bash
# Limit to 8000 tokens (Claude Haiku)
prx --max-tokens 8000

# Combine with relevance for smart selection
prx -r "api routes" --max-tokens 5000
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

# PTX format for programmatic use
prx -o context.ptx
```

## Next Steps

- [Output Formats](output-formats) - PTX, TOON-strict, Markdown, and XML formats
- [Relevance Filtering](relevance-filtering) - Smart file prioritization
- [Configuration](configuration) - Customize behavior
- [File Filtering](file-filtering) - Advanced filtering rules
