
// Helper functions for file analysis
func getConfigDescription(path string) string {
    switch filepath.Base(path) {
    case ".promptext.yml":
        return "Tool configuration file"
    case "go.mod":
        return "Go module definition"
    case ".gitignore":
        return "Git ignore patterns"
    default:
        return "Configuration file"
    }
}

func getDocDescription(path string) string {
    base := filepath.Base(path)
    if strings.HasPrefix(strings.ToUpper(base), "README") {
        return "Project documentation"
    }
    if strings.HasPrefix(strings.ToUpper(base), "LICENSE") {
        return "License information"
    }
    return "Documentation"
}

func isCoreFile(path string) bool {
    dir := filepath.Dir(path)
    return strings.Contains(dir, "internal/") ||
        strings.Contains(dir, "pkg/") ||
        strings.Contains(dir, "lib/")
}

func getCoreDescription(path string) string {
    return "Core implementation"
}

// Helper function to compare path slices
func pathEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
