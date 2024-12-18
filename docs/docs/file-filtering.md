---
sidebar_position: 4
---

# File Filtering

## Default Filters

promptext comes with intelligent default filters to exclude common non-source files.

### Ignored Extensions

- Binary files: `.exe`, `.dll`, `.so`, `.dylib`
- Images: `.jpg`, `.png`, `.gif`, `.svg`
- Archives: `.zip`, `.tar.gz`, `.rar`
- Build artifacts: `.o`, `.obj`, `.class`

### Ignored Directories

- Dependencies: `node_modules/`, `vendor/`
- Build output: `dist/`, `build/`, `target/`
- Version control: `.git/`, `.svn/`
- IDE files: `.idea/`, `.vscode/`

## Custom Filtering

### Via Configuration File

```yaml
excludes:
  - test/
  - "*.generated.go"
  - "internal/tmp/"
```

### Via Command Line

```bash
promptext -exclude "test/,*.generated.go"
```

## GitIgnore Integration

promptext automatically respects your project's `.gitignore` patterns. This can be disabled with:

```bash
promptext -gitignore=false
```
