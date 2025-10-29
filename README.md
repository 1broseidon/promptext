# promptext

Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?prx=v0.2.6)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-docusaurus-blue)](https://1broseidon.github.io/promptext/)

promptext is a code context extraction tool designed for AI assistant interactions. It analyzes codebases, filters relevant files, estimates token usage using tiktoken (GPT-3.5/4 compatible), and provides formatted output suitable for AI prompts.

## Key Features

- **TOON Format Output** - Default token-optimized format (30-60% smaller than JSON/Markdown), inspired by [johannschopplich/toon](https://github.com/johannschopplich/toon)
- **Smart Relevance Filtering** - Multi-factor scoring prioritizes files by keywords (filename, directory, imports, content)
- **Token Budget Management** - Limit output to specific token count, automatically excluding lower-priority files
- **Format Auto-Detection** - Automatically detects output format from file extension (.toon, .md, .xml)
- **Smart File Filtering** - .gitignore support and intelligent defaults
- **Accurate Token Counting** - Using tiktoken (GPT-3.5/4 compatible)
- **Comprehensive Project Analysis** - Entry points, configs, core files, tests, docs
- **Multiple Output Formats** - TOON (default), Markdown, XML
- **Flexible Configuration** - CLI flags or configuration files
- **Project Metadata Extraction** - Language, version, dependencies
- **Git Repository Information** - Branch, commit, message
- **Performance Monitoring** - Debug logging and timing analysis

## Install

Linux/macOS:

```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash
```

Windows (PowerShell):

```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex
```

Go install:

```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

See our [documentation](https://1broseidon.github.io/promptext/) for more installation options.

## Basic Usage

```bash
# Process current directory (TOON format copied to clipboard)
prx

# Process specific directory with positional argument
prx /path/to/project

# Process specific file types
prx -e .go,.js,.ts

# Show project summary only
prx -i

# Auto-detect format from file extension
prx -o context.toon     # TOON format
prx -o context.md       # Markdown format
prx -o project.xml      # XML format

# Explicit format specification
prx -f markdown -o context.md

# Process with custom exclusions and view output in terminal
prx -x "test/,vendor/" --verbose

# Dry run to preview files without processing
prx --dry-run -e .go

# Quiet mode for scripting
prx -q -o output.md
```

## Advanced Features

### Relevance Filtering

Prioritize files matching specific keywords using multi-factor scoring:

```bash
# Prioritize authentication-related files
prx --relevant "auth login OAuth"
prx -r "database SQL postgres"

# Multi-factor scoring weights:
# - Filename matches: 10x
# - Directory matches: 5x
# - Import matches: 3x
# - Content matches: 1x
```

### Token Budget Management

Limit output to stay within token limits for AI models:

```bash
# Limit to 8000 tokens (fits Claude Haiku context)
prx --max-tokens 8000

# Combine with relevance to prioritize important files
prx -r "api routes handlers" --max-tokens 5000

# Cost-optimized queries
prx --max-tokens 3000 -o quick-context.toon
```

When the budget is exceeded, promptext:
- Shows which files were included vs excluded
- Displays token breakdown for excluded files
- Filters directory tree to show only included files

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

## Configuration

Configuration is loaded with the following precedence (highest to lowest):
1. CLI flags
2. Project configuration (`.promptext.yml` in project root)
3. Global configuration (`~/.config/promptext/config.yml`)

### Project Configuration

Create `.promptext.yml` in your project root:

```yaml
extensions:
  - .go
  - .js
  - .ts
excludes:
  - vendor/
  - node_modules/
  - "*.test.go"
format: toon        # Options: toon, markdown, xml
verbose: false
```

### Global Configuration

Create `~/.config/promptext/config.yml` for system-wide defaults:

```yaml
extensions:
  - .go
  - .py
  - .js
excludes:
  - vendor/
  - __pycache__/
format: toon
```

## Documentation

Visit our [documentation site](https://1broseidon.github.io/promptext/) for comprehensive guides on:

- Getting Started Guide
- Configuration Options
- File Filtering Rules
- **Relevance Filtering** - Smart file prioritization
- **Token Budget Management** - Optimize for AI model context windows
- Token Analysis & Counting
- Project Analysis Features
- Output Format Specifications (TOON, Markdown, XML)
- Performance Tips

## License

MIT License - see [LICENSE](LICENSE) for details.
