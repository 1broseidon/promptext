<div align="center">

# promptext

📝 Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)

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

# Using Go
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

## Basic Usage

```bash
# Process current directory
promptext

# Process specific file types
promptext -extension .go,.js

# Show project summary with token counts
promptext -info

# Export as XML with debug logging
promptext -format xml -output project.xml -debug

# Process with custom exclusions
promptext -exclude "test/,vendor/" -verbose
```

## Configuration

Create a `.promptext.yml` in your project root:

```yaml
extensions:
  - .go
  - .js
excludes:
  - vendor/
  - '*.test.go'
format: markdown
verbose: false
```

See our [full documentation](docs/docs.md) for:

- Detailed configuration options
- Output format specifications
- File filtering rules and defaults
- Project analysis features
- Token counting methodology
- Performance optimization tips
- And more!

## License

MIT License - see [LICENSE](LICENSE) for details.
