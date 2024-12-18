package token

import (
	"log"
	"os"
	"path/filepath"

	"github.com/pkoukk/tiktoken-go"
)

func init() {
	// Set default cache directory if TIKTOKEN_CACHE_DIR is not set
	if os.Getenv("TIKTOKEN_CACHE_DIR") == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Printf("Warning: Could not get user home directory: %v", err)
			return
		}

		cacheDir := filepath.Join(homeDir, ".promptext", "cache")
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			log.Printf("Warning: Could not create cache directory: %v", err)
			return
		}

		os.Setenv("TIKTOKEN_CACHE_DIR", cacheDir)
	}
}

type TokenCounter struct {
	encoding *tiktoken.Tiktoken
}

func NewTokenCounter() *TokenCounter {
	enc, err := tiktoken.GetEncoding("cl100k_base") // Using cl100k_base as it's used by GPT-3.5/4
	if err != nil {
		log.Printf("Warning: Failed to initialize tiktoken encoder: %v", err)
		return &TokenCounter{}
	}
	return &TokenCounter{
		encoding: enc,
	}
}

func (tc *TokenCounter) EstimateTokens(text string) int {
	if tc.encoding == nil || text == "" {
		return 0
	}

	tokens := tc.encoding.Encode(text, nil, nil)
	count := len(tokens)
	log.Debug("Token count for text (%d chars): %d tokens", len(text), count)
	return count
}
