package initializer

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Priority constants for project types
const (
	PriorityFrameworkSpecific = 100 // Framework-specific (Next.js, Django, Laravel, Angular)
	PriorityBuildTool         = 90  // Build tools and framework configs (Vite, Nuxt, Flask)
	PriorityLanguage          = 80  // Language-specific (Go, Rust, Java, .NET)
	PriorityGeneric           = 70  // Generic language (Python, PHP)
	PriorityBasic             = 60  // Basic/generic (Node.js)
)

// ProjectType represents a detected project framework or language
type ProjectType struct {
	Name        string
	Description string
	Priority    int // Higher priority types are listed first
}

// DetectionResult contains all detected project types
type DetectionResult struct {
	ProjectTypes []ProjectType
	RootPath     string
}

// Detector interface for project type detection
type Detector interface {
	Detect(rootPath string) ([]ProjectType, error)
}

// FileDetector detects project types based on file presence
type FileDetector struct{}

// NewFileDetector creates a new file-based detector
func NewFileDetector() *FileDetector {
	return &FileDetector{}
}

// Detect scans the directory for known project indicators
func (d *FileDetector) Detect(rootPath string) ([]ProjectType, error) {
	var detected []ProjectType

	// Define detection rules: file -> project type
	detectionRules := []struct {
		files       []string // Any of these files indicates this project type
		projectType ProjectType
	}{
		// JavaScript/TypeScript frameworks
		{
			files: []string{"next.config.js", "next.config.mjs", "next.config.ts"},
			projectType: ProjectType{
				Name:        "nextjs",
				Description: "Next.js",
				Priority:    PriorityFrameworkSpecific,
			},
		},
		{
			files: []string{"nuxt.config.js", "nuxt.config.ts"},
			projectType: ProjectType{
				Name:        "nuxt",
				Description: "Nuxt.js",
				Priority:    PriorityFrameworkSpecific,
			},
		},
		{
			files: []string{"vite.config.js", "vite.config.ts"},
			projectType: ProjectType{
				Name:        "vite",
				Description: "Vite",
				Priority:    PriorityBuildTool,
			},
		},
		{
			files: []string{"vue.config.js"},
			projectType: ProjectType{
				Name:        "vue",
				Description: "Vue.js",
				Priority:    PriorityBuildTool,
			},
		},
		{
			files: []string{"angular.json"},
			projectType: ProjectType{
				Name:        "angular",
				Description: "Angular",
				Priority:    PriorityFrameworkSpecific,
			},
		},
		{
			files: []string{"svelte.config.js"},
			projectType: ProjectType{
				Name:        "svelte",
				Description: "Svelte",
				Priority:    PriorityBuildTool,
			},
		},

		// Go
		{
			files: []string{"go.mod"},
			projectType: ProjectType{
				Name:        "go",
				Description: "Go",
				Priority:    PriorityLanguage,
			},
		},

		// Python frameworks
		{
			files: []string{"manage.py", "django"},
			projectType: ProjectType{
				Name:        "django",
				Description: "Django",
				Priority:    PriorityFrameworkSpecific,
			},
		},
		{
			files: []string{"app.py", "wsgi.py"},
			projectType: ProjectType{
				Name:        "flask",
				Description: "Flask",
				Priority:    PriorityBuildTool,
			},
		},
		{
			files: []string{"pyproject.toml", "setup.py", "requirements.txt"},
			projectType: ProjectType{
				Name:        "python",
				Description: "Python",
				Priority:    PriorityGeneric,
			},
		},

		// Rust
		{
			files: []string{"Cargo.toml"},
			projectType: ProjectType{
				Name:        "rust",
				Description: "Rust",
				Priority:    PriorityLanguage,
			},
		},

		// Java/Kotlin
		{
			files: []string{"pom.xml"},
			projectType: ProjectType{
				Name:        "maven",
				Description: "Maven (Java)",
				Priority:    PriorityLanguage,
			},
		},
		{
			files: []string{"build.gradle", "build.gradle.kts"},
			projectType: ProjectType{
				Name:        "gradle",
				Description: "Gradle (Java/Kotlin)",
				Priority:    PriorityLanguage,
			},
		},

		// Ruby
		{
			files: []string{"Gemfile", "config.ru"},
			projectType: ProjectType{
				Name:        "ruby",
				Description: "Ruby/Rails",
				Priority:    PriorityLanguage,
			},
		},

		// PHP
		{
			files: []string{"composer.json"},
			projectType: ProjectType{
				Name:        "php",
				Description: "PHP",
				Priority:    PriorityGeneric,
			},
		},
		{
			files: []string{"artisan"},
			projectType: ProjectType{
				Name:        "laravel",
				Description: "Laravel",
				Priority:    PriorityBuildTool,
			},
		},

		// .NET
		{
			files: []string{"*.csproj", "*.fsproj", "*.vbproj"},
			projectType: ProjectType{
				Name:        "dotnet",
				Description: ".NET",
				Priority:    PriorityLanguage,
			},
		},

		// Node.js (generic - lowest priority)
		{
			files: []string{"package.json"},
			projectType: ProjectType{
				Name:        "node",
				Description: "Node.js",
				Priority:    PriorityBasic,
			},
		},
	}

	// Check each detection rule
	for _, rule := range detectionRules {
		for _, file := range rule.files {
			// Safety check for empty strings
			if len(file) == 0 {
				continue
			}

			// Check if pattern contains wildcards
			if strings.Contains(file, "*") {
				// Use glob matching for wildcard patterns
				matches, err := filepath.Glob(filepath.Join(rootPath, file))
				if err == nil && len(matches) > 0 {
					detected = append(detected, rule.projectType)
					break
				}
			} else {
				// Regular file existence check
				filePath := filepath.Join(rootPath, file)
				if _, err := os.Stat(filePath); err == nil {
					detected = append(detected, rule.projectType)
					break
				}
			}
		}
	}

	// Sort by priority (highest first) using sort.Slice
	sort.Slice(detected, func(i, j int) bool {
		return detected[i].Priority > detected[j].Priority
	})

	// Deduplicate
	seen := make(map[string]bool)
	var unique []ProjectType
	for _, pt := range detected {
		if !seen[pt.Name] {
			seen[pt.Name] = true
			unique = append(unique, pt)
		}
	}

	return unique, nil
}
