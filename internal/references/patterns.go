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
        // Go imports - handle both quoted and parenthesized imports
        regexp.MustCompile(`(?m)^\s*import\s*\(?\s*"([^"]+)"\)?`),
        // Python imports - handle both import and from...import
        regexp.MustCompile(`(?m)^\s*(?:from\s+([\w\._]+(?:\s+import\s+[\w\s,]+)?)|import\s+([\w\._]+))`),
        // JavaScript/TypeScript imports
        regexp.MustCompile(`(?m)import\s+(?:{[^}]*}\s+from\s+)?['"]([^'"]+)['"]|require\(['"]([^'"]+)['"]\)`),
        // Markdown links
        regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),
        // Local file references in comments
        regexp.MustCompile(`(?m)(?:\/\/|#)\s*(?:see|ref|reference):\s*([^\s;]+)`),
    }
)
