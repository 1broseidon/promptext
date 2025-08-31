---
title: Introduction
description: Smart code context extractor for AI assistants
---

# promptext

Smart code context extraction for AI assistants. Process codebases efficiently with accurate token counting and intelligent filtering.

## Features

- **Smart Filtering** — .gitignore integration with intelligent defaults
- **Token Analysis** — Accurate GPT-compatible counting with tiktoken
- **Project Detection** — Language and framework identification
- **Multiple Formats** — Markdown and XML output options
- **Performance Focused** — Optimized for large codebases

## Quick Start

```bash
# Install
go install github.com/1broseidon/promptext/cmd/promptext@latest

# Extract current directory
promptext

# Analyze specific project
promptext -d /path/to/project
```

Continue with [Getting Started](getting-started) for detailed setup and usage.
