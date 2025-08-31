---
title: Configuration
description: Configure promptext using YAML files and command-line flags
---

## Config File

Create `.promptext.yml` in your project root:

```yaml
extensions:
  - .go
  - .js
  - .ts
excludes:
  - node_modules/
  - vendor/
  - "*.test.*"
format: markdown
gitignore: true
no-copy: false
```

## Options

| Setting | Description | Default |
|---------|-------------|---------|
| `extensions` | File types to include | Auto-detect |
| `excludes` | Patterns to skip | Common build dirs |
| `format` | Output format | `markdown` |
| `gitignore` | Respect .gitignore | `true` |
| `no-copy` | Skip clipboard copy | `false` |
| `verbose` | Show full output | `false` |
| `debug` | Enable timing logs | `false` |

## Command Flags

Override config file with command-line flags:

```bash
promptext -e .go,.js -x vendor/ -f xml
```

## Priority

1. **Command flags** (highest)
2. **Config file**
3. **Defaults** (lowest)

Example with mixed configuration:

```bash
# Config file sets extensions to [.go, .js]
# Command overrides to include only .go files
promptext -e .go
```
