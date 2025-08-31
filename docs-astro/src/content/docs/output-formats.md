---
title: Output Formats
description: Markdown and XML output formats for AI assistant integration
---

## Markdown (Default)

Human-readable format optimized for AI assistants:

```bash
promptext -f markdown
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

## XML Format

Structured format for automated processing:

```bash
promptext -f xml -o project.xml
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

## Usage

**Command line:**
```bash
promptext -f xml -o report.xml
```

**Config file:**
```yaml
format: xml
```

**Output to file:**
```bash
promptext -o output.md  # Auto-detects .md format
promptext -o data.xml   # Auto-detects .xml format
```
