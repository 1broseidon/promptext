<div align="center">

# promptext

📝 Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)

</div>

promptext is an intelligent code context extraction tool designed specifically for AI assistant interactions. It analyzes your codebase, filters relevant files, estimates token usage, and provides formatted output suitable for AI prompts.

## Key Features

- 🔍 Smart file filtering with .gitignore support
- 📊 Automatic token counting for AI context limits
- 🗂️ Intelligent project structure analysis
- 📝 Multiple output formats (Markdown, XML)
- 🔧 Configurable via CLI flags or .promptext.yml
- 📈 Project statistics and metadata extraction

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
```

See our [full documentation](docs/docs.md) for:
- Advanced configuration options
- Output format details
- File filtering rules
- Project analysis features
- Token counting methodology
- And more!

## License

MIT License - see [LICENSE](LICENSE) for details.
