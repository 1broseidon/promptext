# promptext

A command-line tool for extracting and formatting code context from projects, designed to help prepare context for AI coding assistants.

## Features

- 🔍 Smart file filtering based on extensions and patterns
- 📁 Respects .gitignore patterns
- 🌲 Generates formatted directory tree structure
- 📋 Automatic clipboard copying
- ⚙️ Configurable via YAML or command-line flags
- 📊 Project metadata detection (language, version, dependencies)
- 🎨 Colored terminal output

## Installation

### Download Binary (Recommended)

Download the latest binary for your platform from the [releases page](https://github.com/1broseidon/promptext/releases).

#### macOS
```bash
# Download and install to /usr/local/bin
curl -L https://github.com/1broseidon/promptext/releases/latest/download/promptext-darwin-amd64 -o /usr/local/bin/promptext && chmod +x /usr/local/bin/promptext
```

#### Linux
```bash
# Download and install to /usr/local/bin
curl -L https://github.com/1broseidon/promptext/releases/latest/download/promptext-linux-amd64 -o /usr/local/bin/promptext && chmod +x /usr/local/bin/promptext
```

### Build from Source

Alternatively, if you have Go installed:
```bash
git clone https://github.com/1broseidon/promptext.git
cd promptext
go build -o promptext cmd/promptext/main.go
```

## Usage

Basic usage:
```bash
promptext [flags]
```

### Command Line Flags

- `-dir string`: Directory path to process (default ".")
- `-ext string`: File extensions to filter (e.g., ".go,.js")
- `-exclude string`: Patterns to exclude (comma-separated)
- `-no-copy`: Disable automatic copying to clipboard
- `-info`: Only display project summary
- `-verbose`: Show full code content in terminal

### Configuration File

You can create a `.promptext.yml` file in your project root to set default options:

```yaml
extensions:
  - .go
  - .md
excludes:
  - vendor/
  - '*.test.go'
  - tmp/
verbose: false
```

Command-line flags take precedence over configuration file settings.

### Examples

Display project summary:
```bash
promptext -info
```

Process only Go files:
```bash
promptext -ext .go
```

Process multiple file types:
```bash
promptext -ext ".go,.js,.py"
```

Exclude specific patterns:
```bash
promptext -exclude "test/,vendor/,*.test.go"
```

Show full content in terminal:
```bash
promptext -verbose
```

## Output Format

The tool generates a formatted output containing:

1. Project Structure (directory tree)
2. Git Information (if available)
   - Branch
   - Latest commit hash
   - Commit message
3. Project Metadata
   - Language and version
   - Dependencies count
4. File Contents (when not using -info flag)
   - Path to each file
   - File contents in markdown code blocks

## Default Ignored Directories

The following directories are ignored by default:
- .git
- node_modules
- vendor
- .idea
- .vscode
- __pycache__
- .pytest_cache
- dist
- build
- coverage
- bin
- .terraform

## License

MIT License

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.
