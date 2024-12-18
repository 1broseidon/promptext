# promptext Documentation

## Table of Contents

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [File Filtering](#file-filtering)
- [Token Analysis](#token-analysis)
- [Project Analysis](#project-analysis)
- [Output Formats](#output-formats)
- [Performance](#performance)

## Overview

promptext is a command-line tool designed to extract and analyze code context for AI assistants. It intelligently processes your project files, using tiktoken for accurate GPT token counting, and generates structured output suitable for AI interactions.

### Key Features

- Smart file filtering with .gitignore integration
- Accurate token counting using tiktoken (GPT-3.5/4 compatible)
- Comprehensive project analysis
- Multiple output formats (Markdown, XML)
- Project metadata extraction
- Git repository information
- Performance monitoring and debug logging

## Installation

### Prerequisites

- Go 1.22 or higher
- Git (for version control features)

### Installation Methods

1. Quick Install (Linux/macOS):

```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/install.sh | bash
```

2. Using Go Install:

```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

3. Manual Installation:

- Download the appropriate binary from the [releases page](https://github.com/1broseidon/promptext/releases)
- Add it to your PATH

## Usage

### Basic Command Structure

```bash
promptext [flags]
```

### Available Flags

- `-directory, -d string`: Directory to process (default ".")
- `-extension, -e string`: File extensions to include (comma-separated, e.g., ".go,.js")
- `-exclude, -x string`: Patterns to exclude (comma-separated)
- `-format, -f string`: Output format (markdown/xml)
- `-output, -o string`: Output file path
- `-info, -i`: Show only project summary with token counts
- `-verbose, -v`: Show full file contents
- `-debug, -D`: Enable debug logging
- `-gitignore, -g`: Use .gitignore patterns (default true)
- `-help, -h`: Show help message

### Examples

1. Process specific file types:

```bash
promptext -extension .go,.js
```

2. Export as XML with debug info:

```bash
promptext -format xml -output project.xml -debug
```

3. Show project overview with token counts:

```bash
promptext -info
```

4. Process with exclusions:

```bash
promptext -exclude "test/,vendor/" -verbose
```

## Configuration

### Configuration File (.promptext.yml)

Place a `.promptext.yml` file in your project root:

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

### Configuration Priority

1. Command-line flags (highest priority)
2. .promptext.yml file
3. Default settings

## File Filtering

### Default Ignored Extensions

- Images: .jpg, .jpeg, .png, .gif, .bmp, .tiff, .webp, etc.
- Binaries: .exe, .dll, .so, .dylib, .bin, .obj, etc.
- Archives: .zip, .tar, .gz, .7z, .rar, .iso, etc.
- Documents: .pdf, .doc, .docx, .xls, .xlsx, .ppt, etc.
- Other: .class, .pyc, .pyo, .pyd, .o, .a, .db, etc.

### Default Ignored Directories

- Version Control: .git/, .svn/, .hg/
- Dependencies: node_modules/, vendor/, bower_components/
- IDE/Editor: .idea/, .vscode/, .vs/
- Build/Output: dist/, build/, out/, bin/, target/
- Cache: **pycache**/, .pytest_cache/, .sass-cache/
- Test Coverage: coverage/, .nyc_output/
- Infrastructure: .terraform/, .vagrant/
- Logs/Temp: logs/, tmp/, temp/

### GitIgnore Integration

promptext respects your project's .gitignore patterns by default. This can be disabled with `-gitignore=false`.

## Token Analysis

### Token Counting

- Uses tiktoken library for accurate GPT-3.5/4 token counting
- Separate token counts for:
  - Directory structure
  - Git repository information
  - Project metadata
  - Source code content
- Real-time token estimation
- Token usage optimization suggestions

### Token Cache

- Caches tiktoken encodings in ~/.promptext/cache
- Improves performance for repeated operations
- Automatically manages cache directory

## Project Analysis

### File Classification

Automatically categorizes files into:

- Entry Points: Main application entry points
- Configuration Files: Project settings and configs
- Core Implementation: Core business logic
- Test Files: Unit and integration tests
- Documentation: README, docs, and comments

### Git Information

Extracts:

- Current branch
- Latest commit hash
- Commit message
- Repository status

### Project Metadata

Detects:

- Primary language
- Project version
- Dependencies (with dev vs. prod distinction)

### Language Support

Automatic detection for:

- Go (go.mod)
- Node.js (package.json)
- Python (requirements.txt, pyproject.toml)
- Rust (Cargo.toml)
- Java (pom.xml, build.gradle)

### Dependency Analysis

Supports multiple package managers:

- go.mod
- package.json
- requirements.txt
- Cargo.toml
- pom.xml
- build.gradle

## Output Formats

### Markdown (Default)

- Clean, readable format
- GitHub-compatible
- Suitable for documentation

### XML

- Structured data format
- Good for automated processing
- Includes detailed metadata

## Performance

### Optimization Features

- Efficient file traversal with early filtering
- Concurrent processing for large codebases
- Memory-efficient token counting
- Smart binary file detection
- Caching of tiktoken encodings

### Debug Mode

Enable with `-debug` flag to see:

- Processing timings
- Token analysis breakdown
- File filtering decisions
- Memory usage statistics
- Performance metrics

### Progress Tracking

- Real-time file processing status
- Token counting progress
- Performance statistics
- Error reporting and handling
