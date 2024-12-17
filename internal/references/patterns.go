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
        regexp.MustCompile(`(?m)^\s*import\s+(?:"([^"]+)"|'([^']+)')`),
        // Python imports
        regexp.MustCompile(`(?m)(?:from|import)\s+([^\s;]+)`),
        // JavaScript/TypeScript imports
        regexp.MustCompile(`(?m)(?:import|require)\s*\(?['"]([^'"]+)['"]\)?`),
        // C/C++ includes
        regexp.MustCompile(`(?m)#include\s+["<]([^">]+)[">]`),
        // Ruby/PHP requires
        regexp.MustCompile(`(?m)(?:require|include)\s+['"]([^'"]+)['"]`),
        // CSS/SCSS imports
        regexp.MustCompile(`(?m)@import\s+['"]([^'"]+)['"]`),
        // Markdown links
        regexp.MustCompile(`(?m)\[[^\]]*\]\(([^)]+)\)`),
        // HTML includes (src/href)
        regexp.MustCompile(`(?m)(?:src|href)=["']([^'"]+)["']`),
    }
)
