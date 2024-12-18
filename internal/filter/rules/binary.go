package rules

import (
	"bufio"
	"bytes"
	"github.com/1broseidon/promptext/internal/filter/types"
	"os"
)

type BinaryRule struct {
	types.BaseRule
}

func NewBinaryRule() types.Rule {
	return &BinaryRule{
		BaseRule: types.NewBaseRule("", types.Exclude),
	}
}

func (r *BinaryRule) Match(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes
	buf := make([]byte, 512)
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
