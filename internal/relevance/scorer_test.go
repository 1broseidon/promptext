package relevance

import (
	"testing"
)

func TestNewScorer_KeywordParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Comma separated",
			input:    "auth,login,oauth",
			expected: []string{"auth", "login", "oauth"},
		},
		{
			name:     "Space separated",
			input:    "auth login oauth",
			expected: []string{"auth", "login", "oauth"},
		},
		{
			name:     "Mixed comma and space",
			input:    "auth, login oauth",
			expected: []string{"auth", "login", "oauth"},
		},
		{
			name:     "Case normalization",
			input:    "Auth,LOGIN,OAuth",
			expected: []string{"auth", "login", "oauth"},
		},
		{
			name:     "Extra whitespace",
			input:    "  auth  ,  login  ,  oauth  ",
			expected: []string{"auth", "login", "oauth"},
		},
		{
			name:     "Empty string",
			input:    "",
			expected: []string{},
		},
		{
			name:     "Single keyword",
			input:    "authentication",
			expected: []string{"authentication"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scorer := NewScorer(tt.input)
			if len(scorer.keywords) != len(tt.expected) {
				t.Errorf("Expected %d keywords, got %d", len(tt.expected), len(scorer.keywords))
			}
			for i, kw := range tt.expected {
				if i >= len(scorer.keywords) || scorer.keywords[i] != kw {
					t.Errorf("Expected keyword[%d] = %q, got %q", i, kw, scorer.keywords[i])
				}
			}
		})
	}
}

func TestScorer_HasKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"With keywords", "auth login", true},
		{"Empty string", "", false},
		{"Whitespace only", "   ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scorer := NewScorer(tt.input)
			if scorer.HasKeywords() != tt.expected {
				t.Errorf("Expected HasKeywords() = %v, got %v", tt.expected, scorer.HasKeywords())
			}
		})
	}
}

func TestScorer_FilenameMatching(t *testing.T) {
	scorer := NewScorer("auth")

	tests := []struct {
		name     string
		path     string
		content  string
		expected float64
	}{
		{
			name:     "Directory and content match",
			path:     "internal/auth/handler.go",
			content:  "package auth",
			expected: DirectoryWeight + ContentWeight, // directory + content (filename doesn't match)
		},
		{
			name:     "Filename exact match",
			path:     "auth.go",
			content:  "",
			expected: FilenameWeight,
		},
		{
			name:     "Case insensitive filename match",
			path:     "AuthHandler.go",
			content:  "",
			expected: FilenameWeight,
		},
		{
			name:     "No filename match",
			path:     "handler.go",
			content:  "",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.ScoreFile(tt.path, tt.content)
			if score != tt.expected {
				t.Errorf("Expected score %.1f, got %.1f", tt.expected, score)
			}
		})
	}
}

func TestScorer_DirectoryMatching(t *testing.T) {
	scorer := NewScorer("database")

	tests := []struct {
		name     string
		path     string
		expected float64
	}{
		{
			name:     "Directory contains keyword",
			path:     "internal/database/conn.go",
			expected: DirectoryWeight,
		},
		{
			name:     "Nested directory match",
			path:     "pkg/database/mysql/client.go",
			expected: DirectoryWeight,
		},
		{
			name:     "No directory match",
			path:     "internal/auth/handler.go",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.ScoreFile(tt.path, "")
			if score != tt.expected {
				t.Errorf("Expected score %.1f, got %.1f", tt.expected, score)
			}
		})
	}
}

func TestScorer_ImportMatching(t *testing.T) {
	scorer := NewScorer("database")

	tests := []struct {
		name     string
		content  string
		expected float64
	}{
		{
			name: "Single import line",
			content: `package main
import "database/sql"`,
			expected: ImportWeight + ContentWeight, // Import + content match
		},
		{
			name: "Import block with match",
			content: `package main
import (
	"fmt"
	"database/sql"
	"net/http"
)`,
			expected: ImportWeight + ContentWeight, // Import + content match
		},
		{
			name: "Multiple import matches",
			content: `package main
import (
	"database/sql"
	"github.com/user/database"
)`,
			expected: (ImportWeight * 2) + (ContentWeight * 2), // 2 imports + 2 content matches
		},
		{
			name: "No import match",
			content: `package main
import (
	"fmt"
	"net/http"
)`,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.ScoreFile("test.go", tt.content)
			if score != tt.expected {
				t.Errorf("Expected score %.1f, got %.1f", tt.expected, score)
			}
		})
	}
}

func TestScorer_ContentMatching(t *testing.T) {
	scorer := NewScorer("authentication")

	tests := []struct {
		name     string
		content  string
		expected float64
	}{
		{
			name: "Single occurrence",
			content: `func validateAuthentication() error {
	return nil
}`,
			expected: ContentWeight,
		},
		{
			name: "Multiple occurrences",
			content: `// authentication package handles authentication
func authentication() error {
	// Check authentication
	return nil
}`,
			expected: ContentWeight * 4, // 4 occurrences (in comment twice, function name, second comment)
		},
		{
			name: "Capped at 10 occurrences",
			content: `authentication authentication authentication authentication authentication
authentication authentication authentication authentication authentication
authentication authentication authentication authentication authentication`,
			expected: ContentWeight * 10, // Capped at 10
		},
		{
			name:     "No match",
			content:  "func handler() error { return nil }",
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := scorer.ScoreFile("test.go", tt.content)
			if score != tt.expected {
				t.Errorf("Expected score %.1f, got %.1f", tt.expected, score)
			}
		})
	}
}

func TestScorer_MultiFactorScoring(t *testing.T) {
	scorer := NewScorer("auth")

	content := `package auth
import (
	"github.com/user/auth"
	"net/http"
)

// auth handles authentication
func HandleAuth() error {
	// Perform auth check
	return nil
}`

	// Should match:
	// - Filename: auth_handler.go (10 points)
	// - Directory: internal/auth (5 points)
	// - Import: github.com/user/auth (3 points)
	// - Content: 6 occurrences (6 points) - "auth" appears in package, import, comment, authentication, HandleAuth, and second comment
	// Total: 24 points
	expectedScore := FilenameWeight + DirectoryWeight + ImportWeight + (ContentWeight * 6)

	score := scorer.ScoreFile("internal/auth/auth_handler.go", content)
	if score != expectedScore {
		t.Errorf("Expected multi-factor score %.1f, got %.1f", expectedScore, score)
	}
}

func TestScorer_MultipleKeywords(t *testing.T) {
	scorer := NewScorer("auth,login")

	// File matches both keywords in filename
	score := scorer.ScoreFile("auth_login.go", "")
	expected := FilenameWeight * 2 // Both keywords match filename
	if score != expected {
		t.Errorf("Expected score %.1f for multiple keyword matches, got %.1f", expected, score)
	}
}

func TestScorer_ScoreFiles(t *testing.T) {
	scorer := NewScorer("database")

	files := []FileContent{
		{Path: "internal/api/handler.go", Content: "package api"},
		{Path: "internal/database/conn.go", Content: "package database"},
		{Path: "pkg/util/helper.go", Content: "import \"database/sql\""},
		{Path: "database.go", Content: "package main"},
	}

	scored := scorer.ScoreFiles(files)

	// Verify sorting - highest scores first
	if len(scored) != 4 {
		t.Fatalf("Expected 4 scored files, got %d", len(scored))
	}

	// database.go should be first (filename match = 10 points)
	if scored[0].Path != "database.go" {
		t.Errorf("Expected database.go first, got %s", scored[0].Path)
	}

	// internal/database/conn.go should be second (directory + content = 6 points)
	if scored[1].Path != "internal/database/conn.go" {
		t.Errorf("Expected internal/database/conn.go second, got %s", scored[1].Path)
	}

	// pkg/util/helper.go should be third (import = 3 points)
	if scored[2].Path != "pkg/util/helper.go" {
		t.Errorf("Expected pkg/util/helper.go third, got %s", scored[2].Path)
	}

	// internal/api/handler.go should be last (no match = 0 points)
	if scored[3].Path != "internal/api/handler.go" {
		t.Errorf("Expected internal/api/handler.go last, got %s", scored[3].Path)
	}
}

func TestScorer_NoKeywords(t *testing.T) {
	scorer := NewScorer("")

	score := scorer.ScoreFile("auth.go", "authentication code")
	if score != 0 {
		t.Errorf("Expected score 0 with no keywords, got %.1f", score)
	}

	files := []FileContent{
		{Path: "a.go", Content: "content"},
		{Path: "b.go", Content: "content"},
	}

	scored := scorer.ScoreFiles(files)
	for i, sf := range scored {
		if sf.Score != 0 {
			t.Errorf("Expected score 0 for file %d with no keywords, got %.1f", i, sf.Score)
		}
	}
}

func TestGetRelevanceThreshold(t *testing.T) {
	threshold := GetRelevanceThreshold()
	expected := FilenameWeight * 0.5

	if threshold != expected {
		t.Errorf("Expected threshold %.1f, got %.1f", expected, threshold)
	}

	// Verify threshold logic
	if threshold > FilenameWeight {
		t.Error("Threshold should be less than a full filename match")
	}
	if threshold < DirectoryWeight {
		t.Error("Threshold should be greater than a directory match")
	}
}
