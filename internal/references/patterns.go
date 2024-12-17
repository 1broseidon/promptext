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
        // Go imports
        regexp.MustCompile(`(?m)^\s*import\s+(?:"([^"]+)"|([^"'\s]+))`),
        // Python imports
        regexp.MustCompile(`(?m)^\s*(?:from\s+([^\s;]+)\s+import|import\s+([^\s;]+))`),
        // JavaScript/TypeScript imports
        regexp.MustCompile(`(?m)(?:import\s+.*?from\s+['"]([^'"]+)['"]|require\s*\(['"]([^'"]+)['"]\))`),
        // Markdown links - only file links, not URLs
        regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),
        // Local file references in comments
        regexp.MustCompile(`(?m)(?:\/\/|#)\s*(?:see|ref|reference):\s*([^\s]+)`),
    }
)
