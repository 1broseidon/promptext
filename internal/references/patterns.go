package references

import "regexp"

var (
	// Common prefixes that indicate non-local references
	nonLocalPrefixes = []string{
		"http://", "https://", "mailto:", "tel:", "ftp://",
		"github.com/", "golang.org/", "gopkg.in/",
		"@", // npm packages
		"~", // home directory
	}

	// Common file extensions to try when resolving references
	commonExtensions = []string{
		// Source files
		".go", ".py", ".js", ".ts", ".jsx", ".tsx",
		".rb", ".php", ".java", ".cpp", ".c", ".h",
		// Documentation
		".md", ".rst", ".txt",
		// Config files
		".yml", ".yaml", ".json", ".toml",
	}

	referencePatterns = []*regexp.Regexp{
		// Go imports - simpler pattern for single-line imports
		regexp.MustCompile(`(?m)^\s*import\s+(?:"([^"]+)"|([A-Za-z0-9_/\.-]+))`),

		// Python imports - separate patterns for "import" and "from ... import ..."
		regexp.MustCompile(`(?m)^\s*import\s+([\w\.]+)`),
		regexp.MustCompile(`(?m)^\s*from\s+([\w\.]+)`),

		// JavaScript/TypeScript imports - separate patterns for import and require
		regexp.MustCompile(`(?m)import\s+(?:{[^}]*}\s+from\s+)?['"]([^'"]+)['"]`),
		regexp.MustCompile(`(?m)require\(['"]([^'"]+)['"]\)`),

		// Markdown links
		regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),

		// Local file references in comments
		regexp.MustCompile(`(?m)(?:\/\/|#)\s*(?:see|ref|reference):\s*([^\s;]+)`),
	}
)
