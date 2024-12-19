<div align="center">

# promptext

📝 Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?prx=v0.2.3)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-docusaurus-blue)](https://1broseidon.github.io/promptext/)

</div>

promptext is an intelligent code context extraction tool designed specifically for AI assistant interactions. It analyzes your codebase, filters relevant files, estimates token usage using tiktoken (GPT-3.5/4 compatible), and provides formatted output suitable for AI prompts.

## Key Features

- 🔍 Smart file filtering with .gitignore support and intelligent defaults
- 📊 Accurate token counting using tiktoken (GPT-3.5/4 compatible)
- 🗂️ Comprehensive project analysis (entry points, configs, core files, tests, docs)
- 📝 Multiple output formats (Markdown, XML)
- 🔧 Configurable via CLI flags or .promptext.yml
- 📈 Project metadata extraction (language, version, dependencies)
- 🔄 Git repository information extraction
- ⚡ Performance monitoring and debug logging

## Quick Install

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/install.sh | bash

# Windows (PowerShell)
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex

# Using Go
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

### Windows Installation Options

- **System-wide installation (requires admin)**: Run PowerShell as Administrator and use the command above
- **User installation**: Add `-UserInstall` flag to install for current user only
- **Manual installation**: Download the Windows binary from the [releases page](https://github.com/1broseidon/promptext/releases)

## Basic Usage

```bash
# Process current directory
prx

# Process specific file types
prx -e .go,.js

# Show project summary
prx -i

# Export as XML with debug logging
prx -f xml -o project.xml -D

# Process with custom exclusions and verbose output
prx -x "test/,vendor/" -v
```

## Configuration

Create a `.promptext.yml` in your project root:

```yaml
extensions:
  - .go
  - .js
excludes:
  - vendor/
  - "*.test.go"
format: markdown
verbose: false
```

## Documentation

Visit our [documentation site](https://1broseidon.github.io/promptext/) for:

- 📚 Getting Started Guide
- ⚙️ Configuration Options
- 🔍 File Filtering Rules
- 📊 Token Analysis
- 🔬 Project Analysis Features
- 📝 Output Format Specifications
- ⚡ Performance Tips

## License

MIT License - see [LICENSE](LICENSE) for details.
