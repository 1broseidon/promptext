package rules

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter/types"
)

// Common binary file extensions - comprehensive list for fast extension-based detection
var binaryExtensions = map[string]bool{
	// Executables and libraries
	".exe": true, ".dll": true, ".so": true, ".dylib": true,
	".bin": true, ".obj": true, ".o": true, ".a": true,
	".lib": true, ".class": true, ".jar": true, ".war": true,
	
	// Archives and compressed files
	".zip": true, ".tar": true, ".gz": true, ".7z": true,
	".rar": true, ".bz2": true, ".xz": true, ".tgz": true,
	".tbz": true, ".tbz2": true, ".lz": true, ".lzma": true,
	
	// Documents
	".pdf": true, ".doc": true, ".docx": true, ".odt": true,
	".xls": true, ".xlsx": true, ".ods": true,
	".ppt": true, ".pptx": true, ".odp": true,
	".rtf": true, ".pages": true, ".numbers": true, ".key": true,
	
	// Images
	".jpg": true, ".jpeg": true, ".png": true, ".gif": true,
	".bmp": true, ".tiff": true, ".tif": true, ".webp": true,
	".ico": true, ".svg": true, ".psd": true, ".ai": true,
	".eps": true, ".raw": true, ".cr2": true, ".nef": true,
	
	// Audio and Video
	".mp3": true, ".wav": true, ".flac": true, ".aac": true,
	".ogg": true, ".wma": true, ".m4a": true,
	".mp4": true, ".avi": true, ".mkv": true, ".mov": true,
	".wmv": true, ".flv": true, ".webm": true, ".m4v": true,
	
	// Databases
	".db": true, ".sqlite": true, ".sqlite3": true, ".mdb": true,
	".accdb": true, ".dbf": true,
	
	// Fonts
	".ttf": true, ".otf": true, ".woff": true, ".woff2": true,
	".eot": true,
	
	// Other binary formats
	".iso": true, ".dmg": true, ".img": true, ".deb": true,
	".rpm": true, ".msi": true, ".pkg": true, ".app": true,
	".pyc": true, ".pyo": true, ".pyd": true,
}

type BinaryRule struct {
	types.BaseRule
}

func NewBinaryRule() types.Rule {
	return &BinaryRule{
		BaseRule: types.NewBaseRule("", types.Exclude),
	}
}

// Match checks if a file is binary using a three-stage approach for optimal performance:
// 1. Extension check (fastest - O(1) map lookup, no I/O)
// 2. File size check (fast - single stat call, no content read)  
// 3. Content analysis (slowest - reads file content as last resort)
func (r *BinaryRule) Match(path string) bool {
	// Stage 1: Check file extension first - fastest method with no I/O
	ext := strings.ToLower(filepath.Ext(path))
	if binaryExtensions[ext] {
		return true
	}

	// Stage 2: Check file size - very large files are likely binary
	// This avoids reading content for obviously binary files like large media/archives
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}
	
	// Files larger than 10MB are likely binary (videos, archives, etc.)
	// This threshold catches most binary files while allowing large text files
	if fileInfo.Size() > 10*1024*1024 {
		return true
	}
	
	// Empty files are not binary
	if fileInfo.Size() == 0 {
		return false
	}

	// Stage 3: Content analysis - only for files that passed previous checks
	// This is the expensive operation we want to minimize
	return r.isBinaryContent(path)
}

// isBinaryContent performs content-based binary detection
func (r *BinaryRule) isBinaryContent(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	// Read first 512 bytes (reduced from 1024) for faster detection
	// Most binary signatures appear in the first few bytes
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

	// Check for high ratio of non-printable characters
	// This is more reliable than UTF-8 validation alone
	nonPrintable := 0
	for _, b := range buf {
		// Count characters outside typical text range
		if b < 7 || (b > 13 && b < 32) || b > 126 {
			nonPrintable++
		}
	}
	
	// If more than 30% of characters are non-printable, likely binary
	if len(buf) > 0 && float64(nonPrintable)/float64(len(buf)) > 0.3 {
		return true
	}

	return false
}
