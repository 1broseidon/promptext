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

### Installation

#### Quick Install (Linux and macOS)

```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/install.sh | bash
```

This will automatically download and install the latest version of promptext for your system.

#### Build from Source

1. Install Go 1.22 or later
2. Clone and build:
```bash
git clone https://github.com/1broseidon/promptext.git
cd promptext
go install ./cmd/promptext
```

The binary will be installed to your `$GOPATH/bin` directory.

#### Manual Installation

Download the latest release for your platform from the [releases page](https://github.com/1broseidon/promptext/releases).

#### Development Build

For development:
```bash
# Clone repository
git clone https://github.com/1broseidon/promptext.git
cd promptext

# Build locally
go build -o promptext cmd/promptext/main.go

# Create a release with GoReleaser (requires GoReleaser installed)
goreleaser release --snapshot --clean
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
