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
