func analyzeProject(rootPath string) *ProjectAnalysis {
    analysis := &ProjectAnalysis{
        EntryPoints:    make(map[string]string),
        ConfigFiles:    make(map[string]string),
        CoreFiles:      make(map[string]string),
        TestFiles:      make(map[string]string),
        Documentation:  make(map[string]string),
    }

    unifiedFilter := filter.NewUnifiedFilter(nil, nil, nil)

    filepath.WalkDir(rootPath, func(path string, d fs.DirEntry, err error) error {
        if err != nil || d.IsDir() {
            return nil
        }

        rel, _ := filepath.Rel(rootPath, path)
        fileType := unifiedFilter.GetFileType(rel)

        switch {
        case strings.HasPrefix(fileType, "entry:"):
            analysis.EntryPoints[rel] = "Project entry point"
        case fileType == "config":
            analysis.ConfigFiles[rel] = getConfigDescription(rel)
        case fileType == "test":
            analysis.TestFiles[rel] = "Test suite"
        case fileType == "doc":
            analysis.Documentation[rel] = getDocDescription(rel)
        case isCoreFile(rel):
            analysis.CoreFiles[rel] = getCoreDescription(rel)
        }

        return nil
    })

    return analysis
}

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
