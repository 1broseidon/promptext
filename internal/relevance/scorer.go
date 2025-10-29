package relevance

import (
	"path/filepath"
	"strings"
)

// Scoring weights for different match types
const (
	FilenameWeight   = 10.0 // Matches in filename are most important
	DirectoryWeight  = 5.0  // Matches in directory/package name
	ImportWeight     = 3.0  // Matches in import statements
	ContentWeight    = 1.0  // Matches in file content
)

// ScoredFile represents a file with its relevance score
type ScoredFile struct {
	Path  string
	Score float64
}

// Scorer handles relevance scoring for files based on keywords
type Scorer struct {
	keywords []string
}

// NewScorer creates a new scorer with parsed keywords
func NewScorer(keywordString string) *Scorer {
	if keywordString == "" {
		return &Scorer{keywords: []string{}}
	}

	// Parse keywords - support both comma and space separation
	keywordString = strings.ReplaceAll(keywordString, ",", " ")
	parts := strings.Fields(keywordString)

	// Normalize keywords to lowercase for case-insensitive matching
	keywords := make([]string, 0, len(parts))
	for _, kw := range parts {
		if normalized := strings.ToLower(strings.TrimSpace(kw)); normalized != "" {
			keywords = append(keywords, normalized)
		}
	}

	return &Scorer{keywords: keywords}
}

// HasKeywords returns true if scorer has any keywords configured
func (s *Scorer) HasKeywords() bool {
	return len(s.keywords) > 0
}

// ScoreFile calculates relevance score for a single file
// Returns 0 if no keywords are configured
func (s *Scorer) ScoreFile(path, content string) float64 {
	if !s.HasKeywords() {
		return 0
	}

	score := 0.0

	// Extract components for scoring
	filename := filepath.Base(path)
	dir := filepath.Dir(path)
	contentLower := strings.ToLower(content)

	// Score each keyword
	for _, keyword := range s.keywords {
		// 1. Filename matches (highest weight)
		if strings.Contains(strings.ToLower(filename), keyword) {
			score += FilenameWeight
		}

		// 2. Directory/package name matches
		if strings.Contains(strings.ToLower(dir), keyword) {
			score += DirectoryWeight
		}

		// 3. Import statement matches
		importScore := s.scoreImports(content, keyword)
		score += float64(importScore) * ImportWeight

		// 4. Content matches (lowest weight)
		// Count occurrences but cap at 10 to prevent single keyword spam from dominating
		contentMatches := strings.Count(contentLower, keyword)
		if contentMatches > 10 {
			contentMatches = 10
		}
		score += float64(contentMatches) * ContentWeight
	}

	return score
}

// scoreImports counts keyword matches in import statements
func (s *Scorer) scoreImports(content, keyword string) int {
	matches := 0
	lines := strings.Split(content, "\n")

	inImportBlock := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Detect import blocks
		if strings.HasPrefix(trimmed, "import (") {
			inImportBlock = true
			continue
		}
		if inImportBlock && trimmed == ")" {
			inImportBlock = false
			continue
		}

		// Check single-line imports and lines within import blocks
		if strings.HasPrefix(trimmed, "import ") || inImportBlock {
			if strings.Contains(strings.ToLower(trimmed), keyword) {
				matches++
			}
		}
	}

	return matches
}

// ScoreFiles scores multiple files and returns them sorted by relevance
func (s *Scorer) ScoreFiles(files []FileContent) []ScoredFile {
	if !s.HasKeywords() {
		// Return all files with zero score if no keywords
		result := make([]ScoredFile, len(files))
		for i, file := range files {
			result[i] = ScoredFile{Path: file.Path, Score: 0}
		}
		return result
	}

	scored := make([]ScoredFile, len(files))
	for i, file := range files {
		scored[i] = ScoredFile{
			Path:  file.Path,
			Score: s.ScoreFile(file.Path, file.Content),
		}
	}

	// Sort by score (highest first)
	for i := 0; i < len(scored)-1; i++ {
		for j := i + 1; j < len(scored); j++ {
			if scored[j].Score > scored[i].Score {
				scored[i], scored[j] = scored[j], scored[i]
			}
		}
	}

	return scored
}

// FileContent represents a file with its content for scoring
type FileContent struct {
	Path    string
	Content string
}

// GetRelevanceThreshold returns suggested minimum score for high-priority files
// Files with scores above this should be prioritized when token budget is limited
func GetRelevanceThreshold() float64 {
	// A file with 1 filename match or 2 directory matches should be considered relevant
	return FilenameWeight * 0.5 // Half of a filename match
}
