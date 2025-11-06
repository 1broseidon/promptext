# promptext

Convert your codebase into AI-ready prompts - a fast, token-efficient alternative to code2prompt for Claude, ChatGPT, and other LLMs.

[![Go Reference](https://pkg.go.dev/badge/github.com/1broseidon/promptext.svg)](https://pkg.go.dev/github.com/1broseidon/promptext)
[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?prx=v0.4.5)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![codecov](https://codecov.io/gh/1broseidon/promptext/branch/main/graph/badge.svg)](https://codecov.io/gh/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-astro-blue)](https://1broseidon.github.io/promptext/)

## Problem

AI assistants need code context. Sending entire repositories exceeds token limits. Manual file selection wastes time. Standard formats (JSON, XML) are verbose.

## Solution

promptext filters files, ranks by relevance, and serializes to token-efficient formats within specified budgets.

## Why promptext?

Unlike other tools like code2prompt, codebase-digest, or manual copy-pasting:
- **Faster**: Written in Go, processes large codebases in seconds
- **Smarter**: Relevance scoring automatically finds the most important files
- **Token-aware**: Built-in tiktoken counting prevents LLM context overflow
- **Format-flexible**: PTX, TOON, Markdown, or XML output for any AI assistant
- **Budget-conscious**: Enforce token limits before sending to expensive API calls

## Features

- **PTX format**: promptext's hybrid TOON format - 25-30% token reduction with explicit paths and multiline code blocks
- **Token budgeting**: Hard limits with relevance-based file selection
- **Relevance scoring**: Keyword matching in paths (10×), directories (5×), imports (3×), content (1×)
- **Standard exclusions**: `.gitignore` patterns, `node_modules/`, lock files, binaries
- **Accurate counting**: tiktoken cl100k_base tokenizer (GPT-3.5/4, Claude compatible)
- **Format options**: PTX (default), TOON-strict, Markdown, XML
- **LLM-optimized**: Works with ChatGPT, Claude, GPT-4, Gemini, and any AI assistant
- **Context window aware**: Respect token limits for Claude Haiku/Sonnet/Opus, GPT-3.5/4
- **AI-friendly formatting**: Structured output for better AI code comprehension

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

### Updating

**Check for updates:**
```bash
prx --check-update
```

**Update to latest version:**
```bash
prx --update
```

promptext automatically checks for new releases once per day and notifies you when updates are available. Network failures are silently ignored to avoid disrupting normal operation.

## Use Cases

- **AI Code Review**: Feed entire projects to Claude/ChatGPT for comprehensive code analysis
- **Context Engineering**: Build optimized prompts within LLM token limits for better AI responses
- **AI Pair Programming**: Provide full codebase context to AI assistants like GitHub Copilot, Cursor, or Windsurf
- **Documentation Generation**: Help AI understand your complete project structure for accurate docs
- **Code Migration**: Give LLMs full legacy codebase context for refactoring suggestions
- **Prompt Engineering**: Create consistent, repeatable AI prompts from code for development workflows
- **Bug Investigation**: Let AI analyze related files together with proper context
- **API Integration**: Generate structured code context for AI-powered development tools

## Usage

### Smart Context Building (The Power Features)

```bash
# Find authentication-related files within token budget
prx -r "auth login OAuth session" --max-tokens 10000

# Get database layer for Claude Haiku (8K limit)
prx -r "database SQL postgres migration" --max-tokens 8000 -o db-context.ptx

# API routes for GPT-4 analysis
prx -r "api routes handlers middleware" --max-tokens 15000

# Bug investigation: error handling code only
prx -r "error exception handler logging" --max-tokens 5000 -e .go,.js
```

**How relevance scoring works:**
- Filename match: 10 points
- Directory path match: 5 points
- Import statement match: 3 points
- Content match: 1 point

### Quick Commands

```bash
# Current directory to clipboard
prx

# Specific directory with extension filter
prx /path/to/project -e .go,.js,.ts

# Output to file (format auto-detected)
prx -o context.ptx      # PTX (default)
prx -o context.md       # Markdown
prx -o project.xml      # XML

# Summary only (file list, token counts)
prx -i

# Preview file selection
prx --dry-run -r "auth"
```

### Token Budget Output

When `--max-tokens` is set and exceeded, promptext shows exactly what was included and excluded:

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

**About PTX**: PTX is a hybrid TOON format specifically created for promptext. It balances the extreme compression of TOON-strict with human readability by using explicit file paths as keys and preserving multiline code blocks. This gives you ~25-30% token savings without sacrificing clarity - perfect for AI assistants that need both efficiency and accurate file path context.

| Format | Token Efficiency | File Path Clarity | Code Preservation | Use Case |
|--------|-----------------|-------------------|-------------------|----------|
| **PTX** (default) | 25-30% reduction | ✅ Explicit quoted paths | Multiline blocks preserved | Code analysis, debugging |
| **TOON-strict** | 30-60% reduction | ✅ Path in array | Escaped to single line | Maximum compression |
| **Markdown** | Baseline (0%) | ✅ In headings | Full fidelity | Human review, documentation |
| **XML** | -20% (more verbose) | ✅ path attribute | Structured elements | Tool integration, parsing |

### PTX Format Example

PTX uses explicit file paths as keys for zero ambiguity:

```yaml
code:
  "internal/config.go": |
    package config

    type Config struct {
        Port int
    }
  "cmd/server/main.go": |
    package main

    func main() {
        // ...
    }
files[2]{path,ext,lines}:
  internal/config.go,go,67
  cmd/server/main.go,go,45
```

**Benefits:**
- **Zero Ambiguity**: AI can instantly map code blocks to exact file paths
- **Token Efficient**: Still uses `|` multiline blocks (~30% savings vs markdown)
- **Human Readable**: No mental translation needed between sanitized keys and actual paths

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