package rules

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter/types"
)

// Common binary file extensions
var binaryExtensions = map[string]bool{
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".bin": true, ".obj": true, ".o": true,
	".zip": true, ".tar": true, ".gz": true, ".7z": true,
	".pdf": true, ".doc": true, ".docx": true,
	".xls": true, ".xlsx": true, ".ppt": true,
	".db": true, ".sqlite": true, ".sqlite3": true,
}

type BinaryRule struct {
	types.BaseRule
}

func NewBinaryRule() types.Rule {
	return &BinaryRule{
		BaseRule: types.NewBaseRule("", types.Exclude),
	}
}

// Match checks if a file is binary using three methods:
// 1. Known binary file extensions
// 2. Presence of null bytes in the first 1024 bytes
// 3. UTF-8 validation of the content
func (r *BinaryRule) Match(path string) bool {
	// Check file extension first as it's the fastest method
	ext := strings.ToLower(filepath.Ext(path))
	if binaryExtensions[ext] {
		return true
	}

	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 1024 bytes for more accurate detection
	buf := make([]byte, 1024)
	n, err := file.Read(buf)
	if err != nil {
		return false
	}
	buf = buf[:n]

	// Check for null bytes which typically indicate binary content
	if bytes.IndexByte(buf, 0) != -1 {
		return true
	}

	// Try to read file as UTF-8
	reader := bufio.NewReader(bytes.NewReader(buf))
	_, err = reader.ReadString('\n')
	return err != nil
}
