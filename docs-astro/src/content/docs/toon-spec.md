---
title: PTX Format Specification
description: Technical specification of promptext's PTX and TOON output formats
---

## Overview

promptext uses PTX (Promptext Context Format) as its default output format. PTX is a hybrid format that combines TOON v1.3 metadata efficiency with readable multiline code blocks, achieving 25-30% token reduction compared to JSON while maintaining code readability.

For maximum compression, promptext also supports TOON-strict mode, which follows the official TOON v1.3 specification exactly, achieving 30-60% token reduction.

## PTX Format (Default)

### Design Philosophy

The PTX format is optimized for:
1. **Balanced efficiency** - 25-30% token reduction while preserving readability
2. **Code debugging** - Multiline blocks maintain formatting and indentation
3. **LLM comprehension** - Structure that language models parse naturally
4. **TOON v1.3 metadata** - Uses official spec for metadata sections
5. **Mixed content** - Handles both structured metadata and readable code

## Format Structure

### Basic Syntax

TOON uses indentation-based hierarchy with key-value pairs:

```toon
key: value
nested:
  subkey: subvalue
  another: value
```

### Data Types

**Strings:**
```toon
simple: no quotes needed
quoted: "when special chars: colons, newlines"
multiline: |
  Line 1
  Line 2
  Preserves formatting
```

**Numbers and Booleans:**
```toon
count: 42
ratio: 3.14
active: true
disabled: false
```

**Arrays:**

Arrays use different formats based on content type:

```toon
# Primitive array (inline format)
tags[3]: go,cli,ai

# Array with count marker
dependencies[6]: pkg1,pkg2,pkg3,pkg4,pkg5,pkg6

# Uniform objects (tabular format - gotoon style)
files[2]{path,ext,lines}:
  main.go,go,100
  utils.go,go,50

# Complex/non-uniform arrays (dash-list format)
items[3]:
  - path: main.go
    complex: true
  - path: utils.go
    complex: false
```

**Objects:**
```toon
metadata:
  name: promptext
  version: 0.4.1
  language: Go
```

## promptext Output Schema

### Top-Level Structure

Every promptext PTX output contains up to 7 top-level sections:

1. **metadata** - Project identification and summary
2. **git** - Git repository information
3. **stats** - File statistics and metrics
4. **structure** - Directory tree as map
5. **files** - Array of file metadata
6. **code** - Map of file contents
7. **analysis** - Optional analysis results

### Section Details

#### metadata

Project-level metadata including language, version, and dependencies:

```toon
metadata:
  dependencies[6]: github.com/spf13/pflag,github.com/atotto/clipboard,github.com/pkoukk/tiktoken-go,github.com/jedib0t/go-pretty/v6,gopkg.in/yaml.v3,github.com/stretchr/testify
  language: Go
  total_files: 36
  total_lines: 5432
  version: 1.22.4
```

**Fields:**
- `dependencies[N]` (inline array) - Package dependencies with count
- `language` (string) - Primary programming language
- `total_files` (int) - Number of included files
- `total_lines` (int) - Total lines of code
- `version` (string, optional) - Project version if detected

**Note:** Primitive arrays use inline format with count marker for token efficiency.

#### git

Git repository status:

```toon
git:
  branch: main
  commit: 209b6f7
  message: Release v0.4.1 - Multi-layered lock file detection
```

**Fields:**
- `branch` (string) - Current branch name
- `commit` (string) - Short commit hash (7 chars)
- `message` (string) - Latest commit message

#### stats

Detailed file statistics:

```toon
stats:
  totalFiles: 36
  totalLines: 5432
  packages: 8
  fileTypes:
    - type: go
      count: 28
    - type: md
      count: 5
    - type: yaml
      count: 3
```

**Fields:**
- `totalFiles` (int) - Total file count
- `totalLines` (int) - Total line count
- `packages` (int) - Number of packages/modules
- `fileTypes` (array) - Breakdown by file extension

#### structure

Directory tree represented as map (path → files):

```toon
structure:
  .:
    - go.mod
    - go.sum
    - README.md
  cmd/promptext:
    - main.go
    - main_test.go
  internal/config:
    - config.go
    - config_test.go
  internal/filter:
    - filter.go
    - filter_test.go
```

**Format:**
- Keys are directory paths (`.` for root)
- Values are arrays of filenames (not full paths)
- Only includes directories with files

#### files

Array of file metadata with path, extension, and line count:

```toon
files:
  - path: cmd/promptext/main.go
    ext: go
    lines: 225
  - path: internal/config/config.go
    ext: go
    lines: 189
  - path: README.md
    ext: md
    lines: 142
```

**Fields per file:**
- `path` (string) - Relative file path
- `ext` (string) - File extension without dot
- `lines` (int) - Number of lines in file

#### code

Map of file contents with sanitized keys:

```toon
code:
  cmd_promptext_main_go: |
    package main

    import (
      "fmt"
      "log"
      "os"
    )

    func main() {
      // Application entry point
      ...
    }
  internal_config_config_go: |
    package config

    import (
      "gopkg.in/yaml.v3"
    )
    ...
```

**Format:**
- Keys are sanitized paths (slashes and dots → underscores)
- Values use YAML multiline `|` syntax
- Preserves original indentation and formatting
- Empty lines maintained

**Key sanitization:**
```
cmd/promptext/main.go → cmd_promptext_main_go
internal/config.go    → internal_config_go
.github/workflows/ci.yml → _github_workflows_ci_yml
```

## Token Efficiency Techniques

### 1. Minimal Syntax
- No braces (`{}`) or brackets (`[]`) except in arrays
- Indentation instead of delimiters
- Unquoted strings where possible

### 2. Separated Metadata and Content
```toon
# Metadata once (compact)
files:
  - path: main.go
    lines: 100

# Content separately (preserves formatting)
code:
  main_go: |
    package main
    ...
```

This avoids repeating metadata within code blocks.

### 3. Directory Structure as Map
```toon
# Instead of full tree:
structure:
  cmd/promptext:
    - main.go
    - main_test.go
```

This is 40-50% more compact than ASCII tree or nested objects.

### 4. Short Field Names (where clear)
- `ext` instead of `extension`
- `msg` instead of `message` (in some contexts)
- But keeps clarity: `dependencies` not `deps`

## Implementation Notes

### PTX Formatter (Default)

Located in `internal/format/formatters.go`, the `PTXFormatter`:
- Uses TOON v1.3 syntax for metadata sections
- Implements YAML-style multiline blocks for code
- Sanitizes file paths for use as keys
- Preserves code formatting and indentation
- Achieves 25-30% token reduction vs JSON

### TOON-Strict Formatter

Located in `internal/format/formatters.go`, the `TOONStrictFormatter`:
- Follows TOON v1.3 specification exactly
- Uses tabular arrays for file metadata
- Escapes all strings (newlines, quotes)
- Maximum token compression (30-60% reduction)
- Best for token-limited contexts

## Comparison with Other Formats

### vs. JSON
**Token reduction: ~50%**
- No syntax overhead (`{}`, `""`, `,`)
- Multiline strings without escaping
- Structure via indentation

### vs. Markdown
**Token reduction: ~30%**
- No repeated headers (`##`, `###`)
- No code fence markers (` ``` `)
- Compact metadata representation

### vs. YAML
**Token reduction: ~15%**
- Similar structure, fewer quotes
- Optimized field names
- Separated code blocks

### PTX vs. TOON-strict

| Feature | PTX (Default) | TOON-strict |
|---------|--------------|-------------|
| Token Reduction | 25-30% | 30-60% |
| Code Readability | Excellent | Poor (escaped) |
| TOON v1.3 Compliance | Partial (metadata only) | Full |
| Multiline Code | Yes (YAML-style) | No (escaped) |
| Best For | Debugging, code review | Maximum compression |

## Future Considerations

### Potential Optimizations

1. **Tabular file metadata** (using gotoon):
   ```toon
   files[36]{path,ext,lines}:
     cmd/promptext/main.go,go,225
     internal/config/config.go,go,189
   ```
   Could save additional 20-30% on metadata.

2. **Compressed imports**:
   ```toon
   imports:
     github.com/spf13:
       - pflag
     github.com/jedib0t:
       - go-pretty/v6
   ```
   Group by org/domain.

3. **Optional fields**:
   Allow users to exclude sections (stats, structure) for max compression.

## Examples

### Minimal Example (Single File)
```toon
metadata:
  language: Go
  total_files: 1
  total_lines: 50

files:
  - path: main.go
    ext: go
    lines: 50

code:
  main_go: |
    package main

    func main() {
      println("Hello, World!")
    }
```

### Full Example (Multi-Package Project)
See [output-formats.md](/output-formats/) for complete example with all sections.

## Validation

PTX format validation:
- Unit tests in `internal/format/formatters_test.go`
- Integration tests with real projects
- LLM compatibility testing (Claude, GPT-4, GPT-3.5)

TOON-strict format validation:
- Follows official TOON v1.3 specification
- Tests ensure proper string escaping
- Compatible with gotoon library parsers

## References

- [PTX v1.0 Specification](/docs/PTX_SPEC_V1.0.md) - PTX format specification
- [Original TOON v1.3 spec](https://github.com/johannschopplich/toon) - TOON-strict mode
- [gotoon library](https://github.com/alpkeskin/gotoon) - Go implementation of TOON v1.3
- [YAML 1.2](https://yaml.org/spec/1.2.2/) - Multiline string syntax for PTX
