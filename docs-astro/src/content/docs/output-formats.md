---
title: Output Formats
description: TOON, Markdown, and XML output formats for AI assistant integration
---

## TOON (Default)

Token-Oriented Object Notation - optimized format for LLM consumption with 30-60% token reduction compared to JSON/Markdown. Inspired by [johannschopplich/toon](https://github.com/johannschopplich/toon).

```bash
promptext              # Uses TOON by default
promptext -f toon      # Explicit TOON format
promptext -o code.toon # Auto-detected from extension
```

**Structure:**
```toon
metadata:
  dependencies[14]: github.com/spf13/pflag,github.com/jedib0t/go-pretty/v6,...
  language: Go
  total_files: 18
  total_lines: 2390
  version: 1.24.1

structure:
  internal/config[2]: config.go,config_test.go
  internal/filter[3]: filter.go,filter_bench_test.go,filter_test.go
  cmd/promptext[1]: main.go

git:
  branch: main
  commit: f8fbf27
  message: "docs: Update changelog"

files[18]{ext,lines,path}:
  .go,89,cmd/promptext/main.go
  .go,156,internal/config/config.go

code:
  cmd_promptext_main_go: |
    package main

    import (
      "fmt"
      ...
```

**Benefits:**
- **30-60% token reduction** - Significantly fewer tokens than JSON/Markdown
- **Scannable structure** - Easy for both humans and LLMs to parse
- **Zero redundancy** - Tabular arrays eliminate repeated field names
- **Compact metadata** - Inline arrays for dependencies and file lists
- **Multiline code support** - YAML-style code blocks preserve formatting

**When to use:**
- Default for all AI assistant interactions
- Token-constrained models (Claude Haiku, GPT-3.5)
- Cost optimization for high-volume queries
- Large codebases where token efficiency matters

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
promptext -o context.toon   # → TOON format
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
| AI assistant queries | TOON | 30-60% token reduction |
| Claude Haiku / GPT-3.5 | TOON | Token efficiency matters |
| Claude Opus / GPT-4 | TOON or Markdown | Either works, TOON saves tokens |
| Human code review | Markdown | Better readability |
| CI/CD integration | XML | Machine-parseable structure |
| Documentation | Markdown | Rich formatting |
| Cost optimization | TOON | Fewer tokens = lower cost |

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
