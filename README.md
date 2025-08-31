# promptext

Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?prx=v0.2.6)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-docusaurus-blue)](https://1broseidon.github.io/promptext/)

promptext is a code context extraction tool designed for AI assistant interactions. It analyzes codebases, filters relevant files, estimates token usage using tiktoken (GPT-3.5/4 compatible), and provides formatted output suitable for AI prompts.

## Key Features

- Smart file filtering with .gitignore support and intelligent defaults
- Accurate token counting using tiktoken (GPT-3.5/4 compatible)
- Comprehensive project analysis (entry points, configs, core files, tests, docs)
- Multiple output formats (Markdown, XML)
- Configurable via CLI flags or configuration files
- Project metadata extraction (language, version, dependencies)
- Git repository information extraction
- Performance monitoring and debug logging

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
# Process current directory (output copied to clipboard)
prx

# Process specific directory with positional argument
prx /path/to/project

# Process specific file types
prx -e .go,.js,.ts

# Show project summary only
prx -i

# Export as XML to file
prx -f xml -o project.xml

# Process with custom exclusions and view output in terminal
prx -x "test/,vendor/" --verbose

# Dry run to preview files without processing
prx --dry-run -e .go

# Quiet mode for scripting
prx -q -o output.md
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
format: markdown
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
format: markdown
```

## Documentation

Visit our [documentation site](https://1broseidon.github.io/promptext/) for comprehensive guides on:

- Getting Started Guide
- Configuration Options
- File Filtering Rules  
- Token Analysis
- Project Analysis Features
- Output Format Specifications
- Performance Tips

## License

MIT License - see [LICENSE](LICENSE) for details.
