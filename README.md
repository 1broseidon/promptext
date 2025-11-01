# promptext

Code context for AI assistants. No bullshit.

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-astro-blue)](https://1broseidon.github.io/promptext/)

Send your codebase to Claude/GPT without the ceremony. Filters noise, counts tokens, fits your budget.

## Why this exists

You're talking to an AI about your code. You need context. Copy-pasting files is tedious. Dumping everything hits token limits. Manually filtering is a waste of time.

promptext does what you'd do manually, but faster: finds relevant files, respects your budget, formats for AI consumption.

## What it does

- **Smaller payloads** - PTX format (25-30% less tokens than JSON) - TOON-inspired hybrid targeting code analysis
- **Stay under budget** - Set max tokens, get the most relevant files that fit
- **Find what matters** - Score files by keywords in paths, imports, content
- **Respect .gitignore** - Plus sane defaults (no node_modules, no lock files, no binaries)
- **Accurate counts** - tiktoken (cl100k_base), same as GPT-3.5/4/Claude
- **Multiple formats** - PTX (default), TOON-strict, Markdown, XML (auto-detect from extension)

> **Format note**: PTX combines TOON v1.3 metadata efficiency with readable multiline code blocks. Perfect for code analysis where you need both token savings and debuggable output. Need maximum compression? Use `toon-strict` mode.

Format inspiration: [johannschopplich/toon](https://github.com/johannschopplich/toon)

## Install

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

More options in the [docs](https://1broseidon.github.io/promptext/).

## Usage

```bash
# Current directory → clipboard (PTX format)
prx

# Specific path
prx /path/to/project

# File types
prx -e .go,.js,.ts

# Just the summary
prx -i

# Format from extension
prx -o context.ptx      # PTX (default - readable code)
prx -o context.toon     # PTX (backward compatibility)
prx -o context.md       # Markdown
prx -o project.xml      # XML

# Or specify explicitly
prx -f ptx -o context.txt        # PTX hybrid (TOON-based, multiline code)
prx -f toon-strict -o small.txt  # TOON v1.3 strict (maximum compression)
prx -f markdown -o context.md    # Human-readable
prx -f xml -o project.xml        # Machine-parseable

# Exclude patterns
prx -x "test/,vendor/" --verbose

# Preview without processing
prx --dry-run -e .go

# Quiet (for scripts)
prx -q -o output.ptx
```

## Power moves

### Relevance filtering

Weight files by keyword matches:

```bash
# Auth code
prx --relevant "auth login OAuth"

# Database stuff
prx -r "database SQL postgres"
```

Scoring:
- Filename: 10×
- Directory: 5×
- Imports: 3×
- Content: 1×

### Token budgets

Stay under model context limits:

```bash
# Claude Haiku budget
prx --max-tokens 8000

# Combined with relevance
prx -r "api routes handlers" --max-tokens 5000

# Cost-optimized
prx --max-tokens 3000 -o quick-context.ptx
```

Exceeded budget? You'll see what got included, what got cut, and why:

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

## Format options

Pick your trade-off:

| Format | Token Savings | Code Readability | Use When |
|--------|--------------|------------------|----------|
| **PTX** (default) | 25-30% | ⭐⭐⭐⭐⭐ | Code analysis, debugging, pair programming |
| **TOON-strict** | 30-60% | ⭐⭐ | Token-limited models, cost optimization |
| **Markdown** | 0% | ⭐⭐⭐⭐⭐ | Human reading, documentation |
| **XML** | -20% | ⭐⭐⭐ | Tool integration, CI/CD |

```bash
# PTX: Readable code, good compression (default)
prx

# TOON-strict: Maximum compression, escaped strings
prx -f toon-strict

# Markdown: Human-friendly, no compression
prx -f markdown

# XML: Machine processing
prx -f xml
```

## Configuration

Precedence: CLI flags → project config → global config

**Project** (`.promptext.yml`):
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

**Global** (`~/.config/promptext/config.yml`):
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

## Docs

[Full documentation](https://1broseidon.github.io/promptext/) covers:

- Getting started
- Configuration
- Filtering rules
- Relevance scoring
- Token budgets
- Output formats (PTX, TOON-strict, Markdown, XML)
- Performance tips

## License

MIT
