---
title: Output Formats
description: PTX, TOON-strict, Markdown, and XML output formats for AI assistant integration
---

## PTX Format (Default)

PTX (Promptext Context Format) is a hybrid format that combines TOON v1.3 metadata efficiency with readable multiline code blocks. This default format provides 25-30% token reduction compared to JSON while maintaining code readability for debugging and analysis.

```bash
promptext               # Uses PTX by default
promptext -f ptx        # Explicit PTX format
promptext -o code.ptx   # Auto-detected from extension
promptext -o code.toon  # Backward compatibility (maps to PTX)
```

**Structure:**

PTX uses TOON v1.3 compliant syntax for metadata sections with multiline code blocks for readability:

```ptx
metadata:
  dependencies[6]: github.com/spf13/pflag,github.com/jedib0t/go-pretty/v6,github.com/pkoukk/tiktoken-go,github.com/atotto/clipboard,gopkg.in/yaml.v3,github.com/stretchr/testify
  language: Go
  total_files: 18
  total_lines: 2390
  version: 1.24.1

git:
  branch: main
  commit: f8fbf27
  message: docs: Update changelog

stats:
  totalFiles: 18
  totalLines: 2390
  packages: 8
  fileTypes:
    - type: go
      count: 14
    - type: md
      count: 3

structure:
  .:
    - go.mod
    - go.sum
    - README.md
  cmd/promptext:
    - main.go
  internal/config:
    - config.go
    - config_test.go

files:
  - path: cmd/promptext/main.go
    ext: go
    lines: 89
  - path: internal/config/config.go
    ext: go
    lines: 156

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
      "os"
      "path/filepath"
    )
    ...
```

**Format Design:**

The PTX format prioritizes:
- **TOON v1.3 metadata** - Token-efficient metadata using inline arrays and tabular formats
- **Multiline code blocks** - YAML-style `|` syntax preserves code formatting and indentation
- **Hybrid approach** - Combines TOON efficiency for metadata with readability for code
- **Count markers** - Array lengths included (`[N]`) for quick sizing
- **Sanitized keys** - File paths converted to valid keys (e.g., `cmd/main.go` → `cmd_main_go`)

**Benefits:**
- **25-30% token reduction** - Significant savings vs JSON while maintaining readability
- **Debugging-friendly** - Code remains scannable with preserved formatting
- **LLM-optimized** - Structure designed for AI assistant comprehension
- **No escaping needed** - Multiline blocks avoid escape sequences in code
- **Multiline code support** - Preserves formatting with `|` syntax
- **Compact metadata** - Nested objects avoid repetition

**When to use:**
- Default for all AI assistant interactions
- Code-heavy projects requiring debugging
- When readability matters alongside token efficiency
- Large codebases where you need to review output

## TOON-Strict Format

Full TOON v1.3 specification compliance for maximum token compression (30-60% reduction vs JSON):

```bash
promptext -f toon-strict      # TOON v1.3 strict mode
promptext -f toon-v1.3        # Alias for toon-strict
```

**Structure:**

TOON-strict follows the official [TOON v1.3 specification](https://github.com/johannschopplich/toon) with escaped strings:

```toon
metadata:
  dependencies[6]: github.com/spf13/pflag,github.com/jedib0t/go-pretty/v6,github.com/pkoukk/tiktoken-go,github.com/atotto/clipboard,gopkg.in/yaml.v3,github.com/stretchr/testify
  language: Go

files[2]{path,ext,lines}:
  cmd/main.go,go,89
  internal/config.go,go,156

code[2]{path,content}:
  "cmd/main.go","package main\n\nimport (\n  \"fmt\"\n  \"log\"\n  \"os\"\n)\n\nfunc main() {\n  // Application entry point\n  ...\n}"
  "internal/config.go","package config\n\nimport \"os\"\n\n..."
```

**Benefits:**
- **Maximum compression** - 30-60% token reduction
- **TOON v1.3 compliant** - Follows official specification
- **Token-optimized** - Best for very limited contexts

**When to use:**
- Token-limited models with strict budgets
- Metadata-heavy projects with less code
- When maximum compression is required
- API integrations with token costs

## Markdown

Human-readable format with rich formatting:

```bash
promptext -f markdown
promptext -o context.md  # Auto-detected from extension
```

**Structure:**
- Project metadata and analysis
- Directory structure with file tree
- Source code with syntax highlighting
- Token counts and statistics

**Benefits:**
- Direct copy-paste to AI chats
- Readable in any text editor
- Preserves code formatting
- Includes contextual information

**When to use:**
- Documentation generation
- Human review of extracted context
- Integration with Markdown-based tools
- When token efficiency is not a concern

## XML Format

Structured format for automated processing:

```bash
promptext -f xml
promptext -o project.xml  # Auto-detected from extension
```

**Structure:**
```xml
<project>
  <metadata language="go" tokens="2590"/>
  <files>
    <file path="main.go" tokens="150">...</file>
    <file path="pkg/utils.go" tokens="85">...</file>
  </files>
  <dependencies>
    <dependency name="github.com/spf13/pflag" version="v1.0.5"/>
  </dependencies>
</project>
```

**Benefits:**
- Machine parseable
- Preserves hierarchical structure
- Suitable for build tools
- Schema validation ready

**When to use:**
- Integration with CI/CD pipelines
- Automated code analysis tools
- Systems requiring strict schema validation
- Legacy systems expecting XML input

## Format Auto-Detection

Promptext automatically detects the output format from file extensions:

```bash
promptext -o context.ptx    # → PTX format
promptext -o context.toon   # → PTX format (backward compatibility)
promptext -o context.md     # → Markdown format
promptext -o project.xml    # → XML format
```

**Conflict handling:**
```bash
# If both flag and extension are specified, flag takes precedence
promptext -f xml -o output.md
# ⚠️  Warning: format flag 'xml' conflicts with output extension '.md'
#     Using 'xml' (flag takes precedence)
```

## Format Selection Guide

| Use Case | Recommended Format | Reason |
|----------|-------------------|---------|
| AI assistant queries | PTX | 25-30% reduction + readable code |
| Code debugging sessions | PTX | Preserves formatting and indentation |
| Token-limited models | TOON-strict | 30-60% maximum compression |
| Smaller context windows | TOON-strict | Extreme token efficiency |
| Large context windows | PTX or Markdown | Either works, PTX saves tokens |
| Human code review | Markdown | Better readability |
| CI/CD integration | XML | Machine-parseable structure |
| Documentation | Markdown | Rich formatting |
| Cost optimization | TOON-strict | Maximum token reduction |

## Configuration

**Command line:**
```bash
promptext -f toon -o report.toon
promptext -f markdown -o context.md
promptext -f xml -o data.xml
```

**Config file (.promptext.yml):**
```yaml
format: toon  # Options: toon, markdown, xml
```

**Global config (~/.config/promptext/config.yml):**
```yaml
format: toon  # Set default for all projects
```
