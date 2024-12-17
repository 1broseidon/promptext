package token

import "testing"

func TestTokenCounter_EstimateTokens(t *testing.T) {
	tests := []struct {
		name     string
		text     string
		expected int
	}{
		{
			name:     "Empty string",
			text:     "",
			expected: 0,
		},
		{
			name:     "Simple text",
			text:     "This is a simple sentence.",
			expected: 6,
		},
		{
			name:     "Text with numbers",
			text:     "There are 123 numbers in 456 this text.",
			expected: 8,
		},
		{
			name:     "Text with symbols",
			text:     "Symbols: !@#$%^&*()_+=-`~[]\\{}|;':\",./<>?",
			expected: 32,
		},
		{
			name:     "Text with mixed content",
			text:     "Mixed content 123 with symbols! and words.",
			expected: 9,
		},
		{
			name:     "Single line code block",
			text:     "```go\nfunc main() {}\n```",
			expected: 6,
		},
		{
			name: "Multi-line code block",
			text: "```go\nfunc main() {\n    println(\"Hello, world!\")\n}\n```",
			expected: 11,
		},
		
		{
			name:     "Code block with symbols",
			text:     "```python\nprint(1 + 2 * 3)\n```",
			expected: 10,
		},
		{
			name:     "Markdown with code block",
			text:     "This is some text.\n```go\nfunc main() {}\n```\nMore text.",
			expected: 10,
		},
		{
			name:     "Markdown with links",
			text:     "See [link](https://example.com) and [another](internal/file.md).",
			expected: 8,
		},
		{
			name:     "Markdown with bold and italic",
			text:     "This is **bold** and *italic* text.",
			expected: 8,
		},
		{
			name: "Code with mixed symbols and words",
			text: "```go\nvar x = a + b * (c - d) / e\n```",
			expected: 15,
		},
		{
			name: "Code with mixed symbols and words and numbers",
			text: "```go\nvar x = a + 123 * (c - 456) / e\n```",
			expected: 17,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings",
			text: "```go\nvar x = a + 123 * (c - 456) / e; println(\"hello\")\n```",
			expected: 20,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings and comments",
			text: "```go\nvar x = a + 123 * (c - 456) / e; // comment\nprintln(\"hello\")\n```",
			expected: 22,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings and comments and newlines",
			text: `
```go
var x = a + 123 * (c - 456) / e; // comment
println("hello")
````,
			expected: 23,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings and comments and newlines and tabs",
			text: `
```go
var x = a + 123 * (c - 456) / e;     // comment
println("hello")
````,
			expected: 23,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings and comments and newlines and tabs and multiple lines",
			text: `
```go
var x = a + 123 * (c - 456) / e;     // comment
println("hello")
var y = 1
````,
			expected: 26,
		},
		{
			name: "Code with mixed symbols and words and numbers and strings and comments and newlines and tabs and multiple lines and empty lines",
			text: `
```go
var x = a + 123 * (c - 456) / e;     // comment

println("hello")
var y = 1
````,
			expected: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tc := NewTokenCounter()
			actual := tc.EstimateTokens(tt.text)
			if actual != tt.expected {
				t.Errorf("EstimateTokens(%q) = %v, want %v", tt.text, actual, tt.expected)
			}
		})
	}
}