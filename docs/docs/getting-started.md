---
sidebar_position: 2
---

# Getting Started

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
- `-use-default-rules, -u`: Use default filtering rules (default true)
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

5. Process all files including dependencies:

```bash
promptext -u=false -exclude "test/" # Disable default rules but keep test/ excluded
```
