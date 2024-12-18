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
```

### Configuration Options

- `extensions`: List of file extensions to include
- `excludes`: List of patterns to exclude
- `verbose`: Enable verbose output
- `format`: Output format (markdown/xml)
- `debug`: Enable debug logging
- `gitignore`: Use .gitignore patterns

## Command Line Flags

All configuration options can be overridden using command line flags:

```bash
promptext -extension .go,.js -exclude vendor/ -format xml
```

### Priority Order

1. Command line flags (highest priority)
2. .promptext.yml file
3. Default settings (lowest priority)
