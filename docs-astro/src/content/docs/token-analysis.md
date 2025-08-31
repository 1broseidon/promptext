---
title: Token Analysis
description: Accurate tiktoken-based token counting for GPT models
---

## Accurate Counting

Uses OpenAI's tiktoken library for precise GPT-3.5/4 compatible token counting.

## Breakdown

Token counts include:

- **Source Files** — Code content with syntax highlighting
- **Directory Structure** — Project layout and organization
- **Metadata** — Language detection, dependencies, entry points
- **Git Information** — Repository details and branch info

## Example Output

```
Token Analysis
--------------
Directory:     120 tokens
Metadata:      85 tokens  
Source Files:  2,340 tokens
Git Info:      45 tokens
--------------
Total:         2,590 tokens
```

## Optimization

promptext optimizes token usage through:

- **Smart Filtering** — Exclude irrelevant files automatically  
- **Binary Detection** — Skip non-text content
- **Efficient Structure** — Minimal overhead for metadata

Use `-i` flag for token analysis without full content:

```bash
promptext -i  # Show only token counts and project summary
```

## Performance

- **Caching** — Token counts cached for repeated runs
- **Incremental** — Only recalculates changed files  
- **Fast** — tiktoken-go provides efficient counting
