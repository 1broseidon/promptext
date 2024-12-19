---
sidebar_position: 4
---

# File Filtering

## Default Filters

promptext comes with intelligent default filters to exclude common non-source files and automatically detects binary files. These default filters can be controlled using the `UseDefaultRules` option in configuration.

### Configuring Default Rules

By default, promptext applies a set of standard filtering rules that exclude common non-source files and directories. You can control this behavior in two ways:

1. Via configuration file:

```yaml
# .promptext.yml
use-default-rules: true  # Enable default filtering rules (default)
use-default-rules: false # Process all files except explicitly excluded ones
```

2. Via command line:

```bash
promptext -u=false  # Disable default rules, only use explicit excludes
promptext --use-default-rules=false  # Same as above, long form
```

When default rules are disabled:

- Common directories (node_modules, vendor, etc.) will not be automatically excluded
- Binary file detection remains active for safety
- Your explicit excludes (via -x flag or config file) still apply
- GitIgnore patterns (if enabled) still apply

This is useful when you need to:

- Process files in typically excluded directories
- Analyze dependencies or vendor code
- Create a custom set of exclusion rules from scratch

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

### Binary File Detection

promptext automatically detects and excludes binary files using a sophisticated detection mechanism:

1. **Null Byte Detection**: Files are scanned for null bytes (0x00) in their first 512 bytes, which typically indicates binary content.
2. **UTF-8 Validation**: Files are validated for proper UTF-8 encoding. If a file cannot be read as valid UTF-8 text, it's considered binary.

This approach means you don't need to manually specify binary file extensions. The tool will automatically exclude:

- Executables and libraries (e.g., .exe, .dll, .so)
- Object files and compiled code
- Images and media files
- Archives and compressed files
- Database files
- And any other non-text content

## Custom Filtering

### Pattern Matching Options

promptext supports several pattern matching options in both configuration files and command-line arguments:

1. **Directory Patterns**: Patterns ending with `/` match directories and their contents. These patterns are matched in two ways:

   ```yaml
   excludes:
     - test/ # Matches both 'test/' and 'path/to/test/'
     - internal/tmp/ # Matches 'internal/tmp/' and 'path/internal/tmp/'
   ```

2. **Wildcard Patterns**: Using `*` for flexible matching. The wildcard is matched against the base filename:

   ```yaml
   excludes:
     - "*.generated.go" # Excludes files ending with .generated.go
     - ".aider*" # Excludes files starting with .aider
   ```

3. **Exact Matches**: For specific files or paths. These can be matched in three ways:
   ```yaml
   excludes:
     - "config.json" # Matches 'config.json' in any directory
     - "src/constants.go" # Matches exact path
     - "internal/" # Matches directory and all contents
   ```

### Via Configuration File

Create a `.promptext.yml` file in your project root:

```yaml
excludes:
  - test/
  - "*.generated.go"
  - "internal/tmp/"
  - "docs/private/"
```

### Via Command Line

Use the `-x` or `-exclude` flag with comma-separated patterns:

```bash
promptext -x "test/,*.generated.go,internal/tmp/"
```

## GitIgnore Integration

promptext automatically respects your project's `.gitignore` patterns. This means any files or directories listed in your `.gitignore` file will also be excluded from processing. This feature can be disabled with:

```bash
promptext -gitignore=false
```

The tool merges patterns from multiple sources in this order:

1. Default exclude patterns
2. GitIgnore patterns (if enabled)
3. Custom exclude patterns from config file or command line

All patterns are deduplicated to ensure efficient processing.
