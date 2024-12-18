---
sidebar_position: 4
---

# File Filtering

## Default Filters

promptext comes with intelligent default filters to exclude common non-source files.

### Ignored Directories and Files

- System files: `.DS_Store`

- Version Control:
  - `.git/`, `.git*`
  - `.svn/`
  - `.hg/`

- Dependencies and Packages:
  - `node_modules/`
  - `vendor/`
  - `bower_components/`
  - `jspm_packages/`

- IDE and Editor:
  - `.idea/`
  - `.vscode/`
  - `.vs/`
  - `*.sublime-*`

- Build and Output:
  - `dist/`
  - `build/`
  - `out/`
  - `bin/`
  - `target/`

- Cache Directories:
  - `__pycache__/`
  - `.pytest_cache/`
  - `.sass-cache/`
  - `.npm/`
  - `.yarn/`

- Test Coverage:
  - `coverage/`
  - `.nyc_output/`

- Infrastructure:
  - `.terraform/`
  - `.vagrant/`

- Logs and Temp:
  - `logs/`
  - `*.log`
  - `tmp/`
  - `temp/`

### Ignored File Extensions

- Images:
  - `.jpg`, `.jpeg`, `.png`, `.gif`, `.bmp`
  - `.tiff`, `.webp`, `.ico`, `.icns`
  - `.svg`, `.eps`, `.raw`, `.cr2`, `.nef`

- Binary Files:
  - `.exe`, `.dll`, `.so`, `.dylib`
  - `.bin`, `.obj`, `.class`
  - `.pyc`, `.pyo`, `.pyd`
  - `.o`, `.a`

- Archives:
  - `.zip`, `.tar`, `.gz`, `.7z`, `.rar`, `.iso`

- Documents:
  - `.pdf`, `.doc`, `.docx`
  - `.xls`, `.xlsx`
  - `.ppt`, `.pptx`

- Database:
  - `.db`, `.db-shm`, `.db-wal`

## Custom Filtering

### Pattern Matching Options

promptext supports several pattern matching options in both configuration files and command-line arguments:

1. **Directory Patterns**: Ending with `/` matches directories and their contents
   ```yaml
   excludes:
     - test/        # Excludes test directory and all contents
     - internal/tmp/ # Excludes tmp directory under internal
   ```

2. **Wildcard Patterns**: Using `*` for flexible matching
   ```yaml
   excludes:
     - "*.generated.go" # Excludes all generated Go files
     - ".aider*"        # Excludes all files starting with .aider
   ```

3. **Exact Matches**: For specific files or paths
   ```yaml
   excludes:
     - "config.json"     # Excludes config.json in any directory
     - "src/constants.go" # Excludes specific file in specific path
   ```

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
