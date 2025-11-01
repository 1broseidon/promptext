# promptext

Converts codebases to token-efficient formats for AI context windows.

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?prx=v0.4.3)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-astro-blue)](https://1broseidon.github.io/promptext/)

## Problem

AI assistants need code context. Sending entire repositories exceeds token limits. Manual file selection wastes time. Standard formats (JSON, XML) are verbose.

## Solution

promptext filters files, ranks by relevance, and serializes to token-efficient formats within specified budgets.

## Features

- **PTX format**: 25-30% token reduction vs JSON (TOON v1.3-based hybrid with multiline code blocks)
- **Token budgeting**: Hard limits with relevance-based file selection
- **Relevance scoring**: Keyword matching in paths (10×), directories (5×), imports (3×), content (1×)
- **Standard exclusions**: `.gitignore` patterns, `node_modules/`, lock files, binaries
- **Accurate counting**: tiktoken cl100k_base tokenizer (GPT-3.5/4, Claude compatible)
- **Format options**: PTX (default), TOON-strict, Markdown, XML

Format reference: [johannschopplich/toon](https://github.com/johannschopplich/toon)

## Installation

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash
```

**Windows:**
```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex
```

**Go:**
```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

See [installation docs](https://1broseidon.github.io/promptext/) for additional methods.

## Basic Usage

```bash
# Current directory to clipboard (PTX format)
prx

# Specific directory
prx /path/to/project

# Filter by extensions
prx -e .go,.js,.ts

# Summary only (file list, token counts)
prx -i

# Output to file (format auto-detected from extension)
prx -o context.ptx      # PTX format
prx -o context.toon     # PTX format (backward compatibility)
prx -o context.md       # Markdown
prx -o project.xml      # XML

# Explicit format specification
prx -f ptx -o context.txt        # PTX: readable code blocks
prx -f toon-strict -o small.txt  # TOON v1.3: maximum compression
prx -f markdown -o context.md    # Standard Markdown
prx -f xml -o project.xml        # XML structure

# Exclude patterns (comma-separated)
prx -x "test/,vendor/" --verbose

# Preview file selection without processing
prx --dry-run -e .go

# Suppress output (useful in scripts)
prx -q -o output.ptx
```

## Advanced Usage

### Relevance Filtering

Rank files by keyword frequency:

```bash
# Authentication-related files
prx --relevant "auth login OAuth session"

# Database layer
prx -r "database SQL postgres migration"

# API endpoints
prx -r "api routes handlers middleware"
```

**Scoring algorithm:**
- Filename match: 10 points per occurrence
- Directory path match: 5 points per occurrence
- Import statement match: 3 points per occurrence
- Content match: 1 point per occurrence

Files ranked by total score. Ties broken by file size (smaller first).

### Token Budget Control

Enforce context window limits:

```bash
# Claude 3 Haiku limit
prx --max-tokens 8000

# Combined relevance + budget
prx -r "api routes handlers" --max-tokens 5000

# Cost optimization for iterative queries
prx --max-tokens 3000 -o quick-context.ptx
```

When budget exceeded, output shows inclusion/exclusion breakdown:

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

Files included in priority order until budget exhausted.

## Output Formats

| Format | Token Efficiency | Code Preservation | Use Case |
|--------|-----------------|-------------------|----------|
| **PTX** (default) | 25-30% reduction | Multiline blocks preserved | Code analysis, debugging |
| **TOON-strict** | 30-60% reduction | Escaped to single line | Maximum compression |
| **Markdown** | Baseline (0%) | Full fidelity | Human review, documentation |
| **XML** | -20% (more verbose) | Structured elements | Tool integration, parsing |

```bash
# PTX: Default, balances compression and readability
prx

# TOON-strict: Aggressive compression, all code escaped
prx -f toon-strict

# Markdown: No compression, standard formatting
prx -f markdown

# XML: Structured output for programmatic consumption
prx -f xml
```

## Configuration

Configuration hierarchy (later overrides earlier):
1. Global config: `~/.config/promptext/config.yml`
2. Project config: `.promptext.yml`
3. CLI flags

**Project config (`.promptext.yml`):**
```yaml
extensions:
  - .go
  - .js
  - .ts
excludes:
  - vendor/
  - node_modules/
  - "*.test.go"
format: ptx        # Options: ptx, toon-strict, markdown, xml
verbose: false
```

**Global config (`~/.config/promptext/config.yml`):**
```yaml
extensions:
  - .go
  - .py
  - .js
excludes:
  - vendor/
  - __pycache__/
format: ptx
```

## Default Exclusions

Always excluded:
- `.git/`, `.hg/`, `.svn/`
- `node_modules/`, `vendor/`, `__pycache__/`
- Lock files: `*-lock.json`, `*.lock`, `Gemfile.lock`, etc.
- Binary files (detected by content)
- Files matching `.gitignore` patterns

Override with `-x` flag or config file `excludes` list.

## Documentation

[Full documentation](https://1broseidon.github.io/promptext/):

- Configuration reference
- Filtering rules and precedence
- Relevance scoring algorithm
- Token counting methodology
- Format specifications (PTX, TOON-strict, Markdown, XML)
- Performance characteristics

## Requirements

- Go 1.21+ (for building from source)
- Git (for `.gitignore` pattern support)

## License

MIT