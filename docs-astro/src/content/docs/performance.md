---
title: Performance
description: Performance optimization and benchmarking for large codebases
---

## Built for Speed

promptext is optimized for large codebases with:

- **Concurrent Processing** — Parallel file handling
- **Smart Caching** — Reduced redundant operations  
- **Binary Detection** — Fast exclusion of non-text files
- **Efficient Filtering** — Skip irrelevant files early

## Benchmarks

| Project Size | Processing Time |
|--------------|-----------------|
| Small (&lt;100 files) | &lt;1 second |
| Medium (100-1K files) | 1-3 seconds |
| Large (1K+ files) | 3-10 seconds |

## Optimization Tips

**Filter early:**
```bash
promptext -e .go,.js  # Process only specific types
```

**Exclude heavy directories:**
```bash
promptext -x "node_modules/,vendor/,dist/"
```

**Use info mode for quick analysis:**
```bash
promptext -i  # Skip full content processing
```

**Leverage gitignore:**
```bash
promptext -g=true  # Respect .gitignore patterns (default)
```

## Debug Mode

Monitor performance with debug logging:

```bash
promptext -D  # Show timing details
```

**Example output:**
```
[DEBUG] File filtering: 45ms
[DEBUG] Binary detection: 120ms  
[DEBUG] Token counting: 230ms
[DEBUG] Total processing: 395ms
```

## Memory Usage

Efficient memory management:
- **Stream Processing** — No large file buffers
- **Incremental Parsing** — Process files as needed
- **Cache Management** — Automatic cleanup of old entries

For memory-constrained environments, use info mode to minimize memory usage.
