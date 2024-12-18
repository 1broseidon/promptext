# Output Formats

promptext supports multiple output formats for flexibility.

## Markdown Format

Default format, ideal for documentation and human readability.

```bash
promptext -format markdown
```

Example output:
```markdown
# Project Analysis

## Files
- main.go (150 tokens)
- utils/helper.go (80 tokens)

## Statistics
Total Files: 2
Total Tokens: 230
```

## XML Format

Structured format for automated processing.

```bash
promptext -format xml
```

Example output:
```xml
<project>
  <files>
    <file name="main.go" tokens="150"/>
    <file name="utils/helper.go" tokens="80"/>
  </files>
  <statistics>
    <totalFiles>2</totalFiles>
    <totalTokens>230</totalTokens>
  </statistics>
</project>
```

## Format Selection

- Use `-format` flag
- Configure in .promptext.yml
- Supports output to file with `-output`
