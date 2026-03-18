<div align="center">

<img src=".github/logo.svg" alt="promptext" width="300">

**Convert your codebase into AI-ready prompts**

A fast, token-efficient tool that transforms your code into optimized context for Claude, ChatGPT, and other LLMs.

[![GitHub Stars](https://img.shields.io/github/stars/1broseidon/promptext?style=social)](https://github.com/1broseidon/promptext/stargazers)
[![Go Reference](https://pkg.go.dev/badge/github.com/1broseidon/promptext.svg)](https://pkg.go.dev/github.com/1broseidon/promptext)
[![Go Report Card](https://goreportcard.com/badge/github.com/1broseidon/promptext?v=0.7.1)](https://goreportcard.com/report/github.com/1broseidon/promptext)
[![codecov](https://codecov.io/gh/1broseidon/promptext/branch/main/graph/badge.svg)](https://codecov.io/gh/1broseidon/promptext)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Release](https://img.shields.io/github/release/1broseidon/promptext.svg)](https://github.com/1broseidon/promptext/releases/latest)
[![Documentation](https://img.shields.io/badge/docs-chain.sh-blue)](https://chain.sh/promptext)

[Documentation](https://1broseidon.github.io/promptext/) • [Installation](#installation) • [Quick Start](#quick-start) • [Examples](#usage)

</div>

---

## The Problem

Working with AI assistants requires code context, but:
- 🔴 Entire repositories exceed token limits
- 🔴 Manual file selection is tedious and error-prone
- 🔴 Standard formats (JSON, XML) waste precious tokens
- 🔴 You never know if you'll hit the context limit until it's too late

## The Solution

`promptext` intelligently filters your codebase, ranks files by relevance, and packages them into token-efficient formats—all within your specified budget.

## Why Choose promptext?

| Challenge | Manual Approach | promptext Solution |
|-----------|----------------|-------------------|
| **Selecting relevant files** | 😓 Manually browse and choose | 🧠 Automatic relevance scoring |
| **Staying within token limits** | ❌ Trial and error, wasted API calls | ✅ Enforced budgets with preview |
| **Efficient formatting** | 📝 Verbose markdown/JSON | 📦 25-60% token reduction |
| **Token counting** | ❓ Guesswork | 🎯 Accurate tiktoken counting |
| **Processing speed** | 🐌 Copy-paste each file | ⚡ Entire codebase in seconds |

### Key Features

- 🚀 **Fast**: Written in Go—processes large codebases in seconds
- 🧠 **Smart**: Relevance scoring automatically prioritizes important files
- 💰 **Budget-Aware**: Enforces token limits to prevent context overflow and save on API costs
- 📦 **Token-Efficient Formats**: PTX (25-30% savings), TOON-strict (30-60% savings), Markdown, or XML
- 🎯 **Accurate Counting**: Uses `cl100k_base` tokenizer (GPT-4, GPT-3.5, Claude compatible)
- ⚙️ **Highly Configurable**: Project-level `.promptext.yml` and global settings support

## Installation

**macOS/Linux:**
```bash
curl -sSL promptext.sh/install | bash
```

**Windows:**
```powershell
irm promptext.sh/install.ps1 | iex
```

**Go Install (requires Go 1.19+):**
```bash
go install github.com/1broseidon/promptext/cmd/promptext@latest
```

**Manual Download:**
Download pre-built binaries from [GitHub Releases](https://github.com/1broseidon/promptext/releases/latest)

The executable is installed as `promptext` with `prx` alias.

### Updating

```bash
# Check for updates
prx --check-update

# Update to latest version
prx --update
```

**Uninstall:**
```bash
curl -sSL promptext.sh/uninstall | bash
```

> **Note:** `promptext` automatically checks for new releases once per day and notifies you when updates are available.

---

## Quick Start

Navigate to your project directory and run:

```bash
promptext
# or use the alias
prx
```

That's it! `promptext` will:
1. Analyze your project structure
2. Filter out unnecessary files (node_modules, binaries, etc.)
3. Package everything into a token-efficient format (PTX by default)
4. Copy the result to your clipboard

Now paste into ChatGPT, Claude, or your favorite LLM and start coding!

---

## Use Cases

Perfect for:

- 🔍 **AI Code Review** — Feed complete projects to Claude/ChatGPT for comprehensive analysis
- 🤖 **AI Pair Programming** — Provide full context to GitHub Copilot, Cursor, or Windsurf
- 📚 **Documentation Generation** — Help AI understand your complete project structure
- 🐛 **Bug Investigation** — Let AI analyze related files together with proper context
- 🔄 **Code Migration** — Give LLMs full legacy codebase context for refactoring
- 🎯 **Prompt Engineering** — Create consistent, repeatable AI prompts from code
- 🔌 **Library Integration** — Use the Go API to integrate code extraction into your AI/ML workflows
- 🛠️ **Build Custom Tools** — Embed promptext capabilities in your own applications

---

## Usage

### Basic Commands

```bash
# Process current directory and copy to clipboard
prx

# Process specific directory
prx /path/to/project

# Filter by file extensions
prx -e .go,.js,.ts

# Output to file (format auto-detected from extension)
prx -o context.ptx      # PTX format
prx -o context.md       # Markdown format  
prx -o project.xml      # XML format

# Show file list and token counts (no output)
prx -i

# Preview file selection without generating output
prx --dry-run
```

### Smart Context Building

Build focused prompts with relevance scoring and token budgets. Start simple and combine options as needed:

```bash
# Start simple: Find authentication-related files
prx -r "auth login session"

# Add a token budget for smaller context windows
prx -r "auth login session" --max-tokens 8000

# Narrow down by file extensions
prx -r "auth login session" --max-tokens 8000 -e .go,.js

# Save to a file for reuse
prx -r "auth login session" --max-tokens 8000 -e .go,.js -o auth-context.ptx

# Complex example: Database layer with multiple keywords
prx -r "database SQL postgres migration schema" --max-tokens 12000 -e .go,.sql -o db-layer.ptx
```

**Real-world scenarios:**

```bash
# Bug investigation: error handling code for limited context LLM
prx -r "error exception handler logging" --max-tokens 5000

# API routes for models with larger context windows
prx -r "api routes handlers middleware" --max-tokens 20000

# Quick security audit: authentication and authorization
prx -r "auth token jwt security session" --max-tokens 10000 -e .go,.js,.ts
```

#### How Relevance Scoring Works

Files are ranked by keyword matches with weighted scoring:

| Match Location | Score | Example |
|----------------|-------|---------|
| Filename | 10x | `auth.go` matches "auth" |
| Directory path | 5x | `auth/handlers/` matches "auth" |
| Import statements | 3x | `import auth` matches "auth" |
| File content | 1x | "auth" appears in code |

Files with the highest scores are included first until the token budget is exhausted.

### Understanding Token Budget Output

When `--max-tokens` is set, `promptext` shows exactly what was included:

```
╭───────────────────────────────────────────────╮
│ 📦 promptext (Go)                             │
│    Included: 7/18 files • ~4,847 tokens       │
│    Full project: 18 files • ~19,512 tokens    │
╰───────────────────────────────────────────────╯

⚠️  Excluded 11 files due to token budget:
    • internal/cli/commands.go (~784 tokens)
    • internal/app/app.go (~60 tokens)
    ... and 9 more files (~8,453 tokens)
    Total excluded: ~9,297 tokens
```

This helps you understand the trade-offs and adjust your filters or budget as needed.

---

## Using as a Library

`promptext` can be used as a Go library in your own applications, allowing you to programmatically extract code context and integrate it into AI/ML workflows.

### Installation

```bash
go get github.com/1broseidon/promptext/pkg/promptext
```

### Quick Start

```go
package main

import (
    "fmt"
    "log"

    "github.com/1broseidon/promptext/pkg/promptext"
)

func main() {
    // Simple extraction
    result, err := promptext.Extract(".")
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Extracted %d files (%d tokens)\n",
        len(result.ProjectOutput.Files),
        result.TokenCount)

    // Use the formatted output
    fmt.Println(result.FormattedOutput)
}
```

### Common Patterns

**Filter by extensions:**
```go
result, err := promptext.Extract(".",
    promptext.WithExtensions(".go", ".mod", ".sum"),
    promptext.WithExcludes("*_test.go", "vendor/"),
)
```

**AI-optimized extraction with token budget:**
```go
result, err := promptext.Extract(".",
    promptext.WithRelevance("auth", "login"),
    promptext.WithTokenBudget(8000),
    promptext.WithFormat(promptext.FormatPTX),
)

// Send to AI API
sendToAI(result.FormattedOutput)
```

**Reusable extractor:**
```go
extractor := promptext.NewExtractor(
    promptext.WithExtensions(".go"),
    promptext.WithTokenBudget(5000),
)

result1, _ := extractor.Extract("/project1")
result2, _ := extractor.Extract("/project2")
```

**Format conversion:**
```go
result, _ := promptext.Extract(".", promptext.WithFormat(promptext.FormatPTX))

// Convert to different formats
markdown, _ := result.As(promptext.FormatMarkdown)
jsonl, _ := result.As(promptext.FormatJSONL)
```

### Available Options

- `WithExtensions(extensions ...string)` - Include specific file extensions
- `WithExcludes(patterns ...string)` - Exclude files matching patterns
- `WithGitIgnore(enabled bool)` - Respect .gitignore patterns (default: true)
- `WithDefaultRules(enabled bool)` - Use built-in filtering rules (default: true)
- `WithRelevance(keywords ...string)` - Filter by keyword relevance
- `WithTokenBudget(maxTokens int)` - Limit output to token budget
- `WithFormat(format Format)` - Set output format (PTX, JSONL, Markdown, XML)
- `WithVerbose(enabled bool)` - Enable verbose logging
- `WithDebug(enabled bool)` - Enable debug logging with timing

### Output Formats

- `FormatPTX` - PTX v2.0 (recommended for AI)
- `FormatJSONL` - Machine-friendly JSONL
- `FormatMarkdown` - Human-readable markdown
- `FormatXML` - Machine-parseable XML

### Error Handling

```go
result, err := promptext.Extract("/invalid/path")
if err != nil {
    if errors.Is(err, promptext.ErrInvalidDirectory) {
        // Handle invalid directory
    }
    if errors.Is(err, promptext.ErrNoFilesMatched) {
        // Handle no matching files
    }
}
```

### Examples

See the [examples/](examples/) directory for complete working examples:
- `examples/basic/` - Simple usage patterns
- `examples/token-budget/` - AI-focused extraction with token limits

For full API documentation, see [pkg.go.dev/github.com/1broseidon/promptext/pkg/promptext](https://pkg.go.dev/github.com/1broseidon/promptext/pkg/promptext)

---

## Output Formats

`promptext` supports multiple output formats optimized for different use cases:

| Format | Token Efficiency | Best For |
|--------|-----------------|----------|
| **PTX** (default) | 25-30% reduction | General AI interactions, code analysis |
| **TOON-strict** | 30-60% reduction | Maximum compression, large codebases |
| **Markdown** | Baseline (0%) | Human readability, documentation |
| **XML** | -20% (more verbose) | Structured parsing, tool integration |

### PTX Format (Recommended)

PTX is a hybrid format created specifically for `promptext`. It balances token efficiency with readability by using explicit file paths and preserving multiline code blocks.

**Example:**

```yaml
code:
  "internal/config.go": |
    package config

    type Config struct {
        Port int
    }
  "cmd/server/main.go": |
    package main

    func main() {
        // ...
    }
files[2]{path,ext,lines}:
  internal/config.go,go,67
  cmd/server/main.go,go,45
```

**Why PTX?**
- ✅ Zero ambiguity — AI instantly maps code to exact file paths
- ✅ Token efficient — ~30% savings vs. markdown
- ✅ Human readable — No mental translation needed
- ✅ LLM-friendly — Clear structure for better AI comprehension

### Switching Formats

```bash
# PTX (default) — balanced compression and readability
prx

# TOON-strict — maximum compression
prx -f toon-strict

# Markdown — no compression, human-friendly
prx -f markdown

# XML — structured output
prx -f xml
```

> **Format Reference:** PTX and TOON-strict are based on [johannschopplich/toon](https://github.com/johannschopplich/toon)

---

## Configuration

Customize `promptext` behavior with configuration files. Settings are applied in order (later overrides earlier):

1. Global config: `~/.config/promptext/config.yml`
2. Project config: `.promptext.yml`
3. CLI flags

### Project Configuration

Generate a starter configuration file in your project:

```bash
prx --init
```

This creates a `.promptext.yml` file with sensible defaults. Customize it for your project:

```yaml
# File extensions to include
extensions:
  - .go
  - .js
  - .ts

# Patterns to exclude (supports glob patterns)
excludes:
  - "vendor/"
  - "node_modules/"
  - "*.test.go"

# Default output format
format: ptx        # Options: ptx, toon-strict, markdown, xml

# Use .gitignore patterns
gitignore: true

# Enable verbose output
verbose: false
```

### Global Configuration

Set defaults for all projects in `~/.config/promptext/config.yml`:

```yaml
extensions:
  - .go
  - .py
  - .js
  - .ts

excludes:
  - "vendor/"
  - "__pycache__/"
  
format: ptx
```

### Default Exclusions

The following are **always excluded** automatically:

- **Version control:** `.git/`, `.hg/`, `.svn/`
- **Dependencies:** `node_modules/`, `vendor/`, `__pycache__/`
- **Lock files:** `*-lock.json`, `*.lock`, `Gemfile.lock`, `poetry.lock`, etc.
- **Binary files:** Detected by content analysis
- **Gitignored files:** Respects your `.gitignore` patterns

> **Tip:** Override exclusions with the `-x` flag or `excludes` list in your config file.

---

## Documentation

For comprehensive documentation, visit [1broseidon.github.io/promptext](https://1broseidon.github.io/promptext/)

Topics covered:
- 📖 **Configuration Reference** — All options and settings
- 🎯 **Filtering Rules** — How files are selected and excluded
- 🧮 **Relevance Scoring** — Algorithm details and tuning
- 🔢 **Token Counting** — Methodology and accuracy
- 📦 **Format Specifications** — PTX, TOON-strict, Markdown, and XML
- ⚡ **Performance** — Benchmarks and optimization tips

---

## Contributing

Contributions are welcome! Whether it's bug reports, feature requests, or code contributions, we'd love your help.

### Ways to Contribute

- 🐛 **Report bugs** — [Open an issue](https://github.com/1broseidon/promptext/issues/new)
- 💡 **Suggest features** — [Start a discussion](https://github.com/1broseidon/promptext/discussions)
- 📝 **Improve docs** — Help make documentation clearer
- 🔧 **Submit PRs** — Fix bugs or add features

### Development Setup

```bash
# Clone the repository
git clone https://github.com/1broseidon/promptext.git
cd promptext

# Build the project
go build -o prx ./cmd/promptext

# Run tests
go test ./...

# Run with coverage
go test -coverprofile=coverage.out ./...
```

### Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes and add tests
4. Ensure tests pass (`go test ./...`)
5. Commit your changes (`git commit -m 'Add amazing feature'`)
6. Push to your fork (`git push origin feature/amazing-feature`)
7. Open a Pull Request

---

## License

This project is licensed under the **MIT License** — see the [LICENSE](LICENSE) file for details.

---

<div align="center">

**Built with ❤️ by the promptext community**

[⭐ Star on GitHub](https://github.com/1broseidon/promptext) • [📖 Documentation](https://1broseidon.github.io/promptext/) • [🐛 Report Bug](https://github.com/1broseidon/promptext/issues)

</div>