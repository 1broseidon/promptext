---
title: Getting Started
description: Quick start guide for promptext installation and basic usage
---

## Installation

Choose your preferred installation method:

### Go Install (Recommended)

```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

### Script Install

**Linux/macOS:**
```bash
curl -sSL https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.sh | bash
```

**Windows PowerShell:**
```powershell
irm https://raw.githubusercontent.com/1broseidon/promptext/main/scripts/install.ps1 | iex
```

### Manual Install

Download binaries from [GitHub Releases](https://github.com/1broseidon/promptext/releases) and add to your PATH.

## Basic Usage

### Simple Commands

```bash
# Process current directory
promptext

# Process specific directory
promptext -d /path/to/project

# Show project overview only
promptext -info

# Export to file
promptext -o output.md
```

### Common Options

| Flag | Description |
|------|-------------|
| `-d` | Directory to process |
| `-e` | File extensions (`.go,.js`) |
| `-x` | Exclude patterns |
| `-f` | Format (`markdown`, `xml`) |
| `-o` | Output file |
| `-i` | Info mode only |
| `-v` | Verbose output |

### Examples

**Filter by file type:**
```bash
promptext -e .go,.js,.ts
```

**Exclude directories:**
```bash
promptext -x "node_modules/,vendor/,test/"
```

**Generate XML report:**
```bash
promptext -f xml -o report.xml
```

Continue with [Configuration](configuration) to customize behavior.
