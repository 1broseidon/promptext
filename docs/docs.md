# promptext Documentation

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [File Filtering](#file-filtering)
- [Output Formats](#output-formats)
- [Project Analysis](#project-analysis)

## Overview
promptext is a command-line tool designed to extract and analyze code context for AI assistants. It intelligently processes your project files, respecting various filtering rules and generating structured output suitable for AI interactions.

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
- `-dir string`: Directory to process (default ".")
- `-ext string`: File extensions to include (comma-separated, e.g., ".go,.js")
- `-exclude string`: Patterns to exclude (comma-separated)
- `-format string`: Output format (markdown/xml/json)
- `-out string`: Output file path
- `-info`: Show only project summary
- `-verbose`: Show full file contents
- `-no-copy`: Disable clipboard copy

### Examples

1. Process specific file types:
```bash
promptext -ext .go,.js
```

2. Export as XML:
```bash
promptext -format xml -out project.xml
```

3. Show project overview:
```bash
promptext -info
```

4. Process with exclusions:
```bash
promptext -exclude "test/,vendor/"
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
  - "*.test.go"
format: markdown
verbose: false
```

### Configuration Priority
1. Command-line flags (highest priority)
2. .promptext.yml file
3. Default settings

## File Filtering

### Default Ignored Extensions
- Images: .jpg, .jpeg, .png, .gif, .bmp, etc.
- Binaries: .exe, .dll, .so, .dylib, etc.
- Archives: .zip, .tar, .gz, .7z, etc.
- Other: .pdf, .doc, .class, .pyc, etc.

### Default Ignored Directories
- Version Control: .git
- Dependencies: node_modules, vendor
- IDE: .idea, .vscode
- Build: dist, build, coverage
- Cache: __pycache__, .pytest_cache

### GitIgnore Integration
promptext respects your project's .gitignore patterns for consistent filtering.

## Output Formats

### Markdown (Default)
- Clean, readable format
- Suitable for documentation
- GitHub-compatible

### XML
- Structured data format
- Good for automated processing
- Includes detailed metadata

### JSON
- Machine-readable format
- Ideal for API integration
- Compact representation

## Project Analysis

### File Classification
promptext automatically categorizes files into:
- Entry Points
- Configuration Files
- Core Implementation
- Test Files
- Documentation

### Language Detection
Supports automatic detection of:
- Go
- Python
- Node.js
- Rust
- Java (Maven/Gradle)

### Dependency Analysis
- Extracts dependencies from package managers
- Identifies dev vs. production dependencies
- Supports multiple package formats:
  - go.mod
  - package.json
  - requirements.txt
  - Cargo.toml
  - pom.xml
  - build.gradle