# promptext Documentation

## Table of Contents
- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [File Filtering](#file-filtering)
- [Output Formats](#output-formats)
- [Project Analysis](#project-analysis)
- [Token Analysis](#token-analysis)
- [Performance](#performance)

## Overview
promptext is a command-line tool designed to extract and analyze code context for AI assistants. It intelligently processes your project files, respecting various filtering rules and generating structured output suitable for AI interactions. The tool focuses on smart file filtering, token counting, and providing comprehensive project analysis while maintaining efficient performance.

### Key Features
- Smart file filtering with .gitignore integration
- Automatic token counting and estimation
- Multiple output formats (Markdown, XML)
- Project structure analysis
- Git repository information
- Dependency analysis
- Performance monitoring

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

4. Process with exclusions and custom output:
```bash
promptext -exclude "test/,vendor/" -verbose -output summary.md
```

5. Disable gitignore patterns:
```bash
promptext -gitignore=false
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
- Entry Points: Main application entry points
- Configuration Files: Project settings and configs
- Core Implementation: Core business logic
- Test Files: Unit and integration tests
- Documentation: README, docs, and comments

### Token Analysis
- Automatic token counting for AI context limits using tiktoken
- Separate token counts for:
  - Directory structure and file hierarchy
  - Git repository information (branch, commits)
  - Project metadata (language, version, dependencies)
  - Source code content and documentation
- Real-time token estimation for different output formats
- Token usage optimization suggestions
- Configurable token counting strategies

### Performance
- Efficient file traversal with early filtering
- Concurrent processing for large codebases
- Memory-efficient token counting
- Progress tracking and timing information
- Debug mode for detailed performance metrics
- Configurable processing options:
  - File extension filtering
  - Directory exclusions
  - GitIgnore integration
  - Token counting optimization

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
