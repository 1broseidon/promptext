---
title: File Filtering
description: Smart file filtering with gitignore support and custom patterns
---

## Smart Defaults

promptext automatically excludes common non-source files:

- **Dependencies**: `node_modules/`, `vendor/`, `bower_components/`
- **Build Output**: `dist/`, `build/`, `out/`, `target/`
- **Version Control**: `.git/`, `.svn/`, `.hg/`
- **IDE Files**: `.idea/`, `.vscode/`, `*.sublime-*`
- **Cache**: `__pycache__/`, `.sass-cache/`, `.npm/`
- **Logs**: `logs/`, `*.log`, `tmp/`

Disable default rules:
```bash
promptext -u=false  # Only use explicit excludes
```

## Binary Detection

Automatically skips binary files by detecting:
- Null bytes in first 512 bytes
- Invalid UTF-8 encoding

No need to specify binary extensions — images, executables, and archives are automatically excluded.

## Custom Patterns

### Pattern Types

**Directory matching:**
```yaml
excludes:
  - test/           # Any 'test' directory
  - internal/tmp/   # Specific path
```

**Wildcards:**
```yaml
excludes:
  - "*.test.go"     # Test files
  - ".aider*"       # Generated files
```

**Exact matches:**
```yaml
excludes:
  - config.json     # Specific file
  - src/constants.go # Exact path
```

### Configuration

**Config file:**
```yaml
# .promptext.yml
excludes:
  - test/
  - "*.generated.*"
  - docs/private/
```

**Command line:**
```bash
promptext -x "test/,*.generated.*,docs/private/"
```

## GitIgnore Integration

Automatically respects `.gitignore` patterns. Disable with:

```bash
promptext -g=false
```

## Filter Priority

1. **Default patterns** (if enabled)
2. **GitIgnore patterns** (if enabled)  
3. **Custom excludes** (config + command line)

All patterns are combined and deduplicated for optimal performance.
