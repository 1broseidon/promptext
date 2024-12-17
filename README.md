<div align="center">

# promptext

📝 Smart code context extractor for AI assistants

[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)

</div>

promptext helps you extract relevant code context from your projects for AI assistants. It intelligently filters files, respects `.gitignore`, and provides clean, formatted output.

## Quick Install

```bash
# Linux/macOS
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/install.sh | bash

# Using Go
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

Or download from our [releases page](https://github.com/1broseidon/promptext/releases).

## Basic Usage

```bash
# Process current directory
promptext

# Process specific file types
promptext -ext .go,.js

# Show project summary
promptext -info

# Export as XML
promptext -format xml -out project.xml
```

See our [full documentation](docs/docs.md) for:
- Advanced configuration
- Output formats
- File filtering rules
- Project analysis features
- And more!

## License

MIT License - see [LICENSE](LICENSE) for details.
