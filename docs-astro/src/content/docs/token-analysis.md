---
title: Token Analysis
description: Accurate tiktoken-based token counting for GPT models
---

## Accurate Counting

promptext uses OpenAI's tiktoken library with the `cl100k_base` encoding for precise token counting compatible with:
- **GPT-4** (all variants)
- **GPT-3.5-turbo** (all variants)
- **Claude** (similar tokenization patterns)

### Fallback Mode

If tiktoken is unavailable or fails to initialize, promptext automatically falls back to an intelligent approximation system:

- **Hybrid Estimation** — Combines word-based (1.3 tokens/word) and character-based estimates
- **Code vs Prose Detection** — Adjusts estimates based on content type
  - Code: ~3.5 characters per token
  - Prose: ~4.0 characters per token
- **User Notification** — Displays "Token counting using approximation (tiktoken unavailable)" when in fallback mode

The fallback mode provides reasonable accuracy (±10-15%) for most use cases.

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
