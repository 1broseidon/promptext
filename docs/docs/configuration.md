---
sidebar_position: 3
---

# Configuration

## Configuration File

promptext can be configured using a `.promptext.yml` file in your project root directory.

### Example Configuration

```yaml
extensions:
  - .go
  - .js
  - .py
excludes:
  - vendor/
  - node_modules/
  - "*.test.go"
verbose: false
format: markdown
debug: false
gitignore: true
use-default-rules: true # Use default filtering rules
no-copy: false # Disable automatic clipboard copying
```

### Configuration Options

- `extensions`: List of file extensions to include
- `excludes`: List of patterns to exclude
- `verbose`: Enable verbose terminal output showing full file contents
- `format`: Output format (markdown/xml)
- `debug`: Enable debug logging with performance timing
- `gitignore`: Use .gitignore patterns
- `use-default-rules`: Enable default filtering rules (default: true)
- `no-copy`: Disable automatic clipboard copying (default: false)

## Command Line Flags

All configuration options can be overridden using command line flags:

```bash
promptext -extension .go,.js -exclude vendor/ -format xml -u=false
```

### Available Flags

- `-d, --directory`: Directory path to process (default: current directory)
- `-e, --extension`: File extensions to include (comma-separated)
- `-x, --exclude`: Patterns to exclude (comma-separated)
- `-f, --format`: Output format (markdown/xml)
- `-o, --output`: Output file path
- `--no-copy`: Disable automatic clipboard copying
- `-i, --info`: Show only project summary with token counts
- `-g, --gitignore`: Use .gitignore patterns (default: true)
- `-u, --use-default-rules`: Use default filtering rules (default: true)
- `-v, --verbose`: Show full code content in terminal
- `-D, --debug`: Enable debug logging with performance timing
- `--version`: Show version information
- `-h, --help`: Show help message

### Priority Order

1. Command line flags (highest priority)
2. .promptext.yml file
3. Default settings (lowest priority)
