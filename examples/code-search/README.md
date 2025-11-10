# Code Search - Semantic Code Discovery

A natural language code search tool that finds relevant code across your codebase using promptext's relevance filtering.

## What This Does

Instead of using traditional grep/regex, this tool lets you search your codebase using **natural language queries**:

- "Where is user authentication handled?"
- "How does database connection pooling work?"
- "Find all API endpoint definitions"
- "Show me the payment processing logic"

The tool:
1. ✅ Extracts keywords from your query
2. ✅ Uses promptext's relevance scoring to find matching files
3. ✅ Limits results to a reasonable token budget
4. ✅ Saves extracted context ready for AI analysis

## Quick Start

```bash
cd examples/code-search

# Search the parent project
go run main.go "Where is user authentication handled?"

# Search a specific directory
SEARCH_DIR=/path/to/project go run main.go "How does caching work?"
```

## Example Output

```
🔍 Searching for: Where is user authentication handled?

📋 Keywords extracted: [user authentication handled]

📂 Searching in: promptext

============================================================
✨ Found 3 relevant files (1,234 tokens)
   ℹ️  7 additional files excluded due to token budget
============================================================

📄 Relevant Files:
   1. internal/auth/handler.go (456 tokens)
   2. pkg/api/auth.go (389 tokens)
   3. cmd/server/middleware.go (389 tokens)

💡 Next Steps:
   The extracted code context is ready to send to an AI assistant.
   You can paste the output below into ChatGPT/Claude to get answers:

   Example prompt:
   "Based on this code: Where is user authentication handled?"

💾 Full context saved to: search-results.ptx
   (5,432 characters, 1,234 tokens)

============================================================
🎯 Search complete!
```

## How It Works

### 1. Keyword Extraction

The tool extracts meaningful keywords from your natural language query by:
- Removing stop words ("where", "is", "the", etc.)
- Filtering out very short words (< 3 characters)
- Removing duplicates

For production use, consider:
- Using NLP libraries (spacy, nltk)
- Calling an AI API to extract semantic keywords
- Building domain-specific keyword dictionaries

### 2. Relevance Filtering

Promptext scores files based on keyword matches:
- **Filename matches**: 10x weight
- **Directory path**: 5x weight
- **Import statements**: 3x weight
- **File content**: 1x weight

### 3. Token Budget

Results are limited to 5,000 tokens by default, which:
- Keeps the context focused and relevant
- Fits well within AI model context windows
- Prevents information overload

### 4. Output Format

Results are saved in PTX format, optimized for:
- AI consumption (25-30% fewer tokens than markdown)
- Easy parsing and processing
- Preserving code structure

## Use Cases

### 1. Onboarding New Developers
```bash
go run main.go "How does the build system work?"
go run main.go "Where are the main entry points?"
```

### 2. Bug Investigation
```bash
go run main.go "Find all error handling code"
go run main.go "Where is logging configured?"
```

### 3. Feature Planning
```bash
go run main.go "Show me the current API structure"
go run main.go "How is configuration managed?"
```

### 4. Code Review Prep
```bash
go run main.go "Find all database queries"
go run main.go "Show authentication and authorization code"
```

## Configuration

### Search Directory

```bash
# Search current directory (default)
go run main.go "your query"

# Search specific directory
SEARCH_DIR=/path/to/project go run main.go "your query"
```

### File Extensions

Edit `main.go` to customize which file types to search:

```go
promptext.WithExtensions(".go", ".js", ".ts", ".py", ".java"),
```

### Token Budget

Adjust the token limit based on your needs:

```go
// Smaller budget = more focused results
promptext.WithTokenBudget(3000),

// Larger budget = more comprehensive context
promptext.WithTokenBudget(10000),
```

## Integration Ideas

### CLI Tool

Build a standalone tool:
```bash
# Install as code-search command
go build -o code-search main.go
sudo mv code-search /usr/local/bin/

# Use anywhere
cd ~/my-project
code-search "Where is error handling?"
```

### VS Code Extension

- Bind to keyboard shortcut
- Show results in sidebar
- Jump to relevant files

### Terminal Integration

Add to your shell:
```bash
# .bashrc or .zshrc
alias csearch='go run /path/to/code-search/main.go'
```

### AI Integration

Send results directly to AI APIs:

```go
// After extracting context
response := callClaudeAPI(query, result.FormattedOutput)
fmt.Println(response)
```

## Limitations

### Keyword Extraction
- Simple word splitting (no NLP)
- May miss synonyms or related terms
- Works best with explicit technical terms

**Solutions:**
- Use AI to extract semantic keywords
- Build domain-specific dictionaries
- Allow manual keyword specification

### Search Scope
- Only searches file content and paths
- Doesn't understand code semantics
- May miss relevant code with different terminology

**Solutions:**
- Combine with traditional grep for completeness
- Use AI to refine search results
- Build index of code symbols

## Advanced Usage

### Combine with Git

Search only changed files:
```bash
# Get recently modified files
SEARCH_DIR=$(git diff --name-only HEAD~10 | head -1 | xargs dirname)
go run main.go "what changed?"
```

### Pre-filter by Directory

Search specific modules:
```bash
cd internal/api
go run ../../examples/code-search/main.go "authentication"
```

### Batch Queries

Create a query file:
```bash
cat queries.txt | while read query; do
  echo "=== $query ==="
  go run main.go "$query"
done
```

## Tips for Better Results

1. **Be specific**: "JWT authentication" > "auth"
2. **Use technical terms**: "connection pool" > "database stuff"
3. **Include context**: "user login validation" > "validation"
4. **Try variations**: If no results, rephrase your query

## Next Steps

- ✅ Add AI integration to answer queries automatically
- ✅ Build index for faster searching
- ✅ Support multiple programming languages
- ✅ Add search history and bookmarks
- ✅ Create web UI for team sharing

## Related Examples

- [CI Code Review](../ci-code-review/) - Automated PR analysis
- [Doc Generator](../doc-generator/) - Keep docs in sync
- [Migration Assistant](../migration-assistant/) - Modernize legacy code
