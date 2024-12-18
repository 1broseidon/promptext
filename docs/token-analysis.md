# Token Analysis

## Overview

promptext uses OpenAI's tiktoken library to provide accurate token counting for GPT models.

## Token Counting Features

- Accurate GPT-3.5/4 compatible counting
- Per-file token breakdowns
- Directory structure tokens
- Metadata tokens
- Git information tokens

## Token Usage Optimization

promptext helps optimize token usage by:

1. Smart file filtering
2. Intelligent content truncation
3. Token usage warnings
4. Optimization suggestions

## Example Output

```
Project Token Analysis:
----------------------
Directory Structure:  150 tokens
Git Information:      80 tokens
Source Files:         2,450 tokens
Documentation:        320 tokens
------------------------
Total:                3,000 tokens
```

## Token Caching

- Caches token counts for improved performance
- Automatic cache invalidation on file changes
- Configurable cache location
