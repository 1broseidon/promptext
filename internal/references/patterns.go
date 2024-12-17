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
		// Markdown links
		regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),

		// Local file references in comments
		regexp.MustCompile(`(?m)(?:\/\/|#)\s*(?:see|ref|reference):\s*([^\s;]+)`),
	}
)
