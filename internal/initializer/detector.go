package initializer

import (
	"os"
	"path/filepath"
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
				Priority:    100,
			},
		},
		{
			files: []string{"nuxt.config.js", "nuxt.config.ts"},
			projectType: ProjectType{
				Name:        "nuxt",
				Description: "Nuxt.js",
				Priority:    100,
			},
		},
		{
			files: []string{"vite.config.js", "vite.config.ts"},
			projectType: ProjectType{
				Name:        "vite",
				Description: "Vite",
				Priority:    90,
			},
		},
		{
			files: []string{"vue.config.js"},
			projectType: ProjectType{
				Name:        "vue",
				Description: "Vue.js",
				Priority:    90,
			},
		},
		{
			files: []string{"angular.json"},
			projectType: ProjectType{
				Name:        "angular",
				Description: "Angular",
				Priority:    100,
			},
		},
		{
			files: []string{"svelte.config.js"},
			projectType: ProjectType{
				Name:        "svelte",
				Description: "Svelte",
				Priority:    90,
			},
		},

		// Go
		{
			files: []string{"go.mod"},
			projectType: ProjectType{
				Name:        "go",
				Description: "Go",
				Priority:    80,
			},
		},

		// Python frameworks
		{
			files: []string{"manage.py", "django"},
			projectType: ProjectType{
				Name:        "django",
				Description: "Django",
				Priority:    100,
			},
		},
		{
			files: []string{"app.py", "wsgi.py"},
			projectType: ProjectType{
				Name:        "flask",
				Description: "Flask",
				Priority:    90,
			},
		},
		{
			files: []string{"pyproject.toml", "setup.py", "requirements.txt"},
			projectType: ProjectType{
				Name:        "python",
				Description: "Python",
				Priority:    70,
			},
		},

		// Rust
		{
			files: []string{"Cargo.toml"},
			projectType: ProjectType{
				Name:        "rust",
				Description: "Rust",
				Priority:    80,
			},
		},

		// Java/Kotlin
		{
			files: []string{"pom.xml"},
			projectType: ProjectType{
				Name:        "maven",
				Description: "Maven (Java)",
				Priority:    80,
			},
		},
		{
			files: []string{"build.gradle", "build.gradle.kts"},
			projectType: ProjectType{
				Name:        "gradle",
				Description: "Gradle (Java/Kotlin)",
				Priority:    80,
			},
		},

		// Ruby
		{
			files: []string{"Gemfile", "config.ru"},
			projectType: ProjectType{
				Name:        "ruby",
				Description: "Ruby/Rails",
				Priority:    80,
			},
		},

		// PHP
		{
			files: []string{"composer.json"},
			projectType: ProjectType{
				Name:        "php",
				Description: "PHP",
				Priority:    70,
			},
		},
		{
			files: []string{"artisan"},
			projectType: ProjectType{
				Name:        "laravel",
				Description: "Laravel",
				Priority:    90,
			},
		},

		// .NET
		{
			files: []string{"*.csproj", "*.fsproj", "*.vbproj"},
			projectType: ProjectType{
				Name:        "dotnet",
				Description: ".NET",
				Priority:    80,
			},
		},

		// Node.js (generic - lowest priority)
		{
			files: []string{"package.json"},
			projectType: ProjectType{
				Name:        "node",
				Description: "Node.js",
				Priority:    60,
			},
		},
	}

	// Check each detection rule
	for _, rule := range detectionRules {
		for _, file := range rule.files {
			// Handle wildcards
			if filepath.Base(file) != file && (file[0] == '*' || file[len(file)-1] == '*') {
				matches, err := filepath.Glob(filepath.Join(rootPath, file))
				if err == nil && len(matches) > 0 {
					detected = append(detected, rule.projectType)
					break
				}
			} else {
				// Regular file check
				filePath := filepath.Join(rootPath, file)
				if _, err := os.Stat(filePath); err == nil {
					detected = append(detected, rule.projectType)
					break
				}
			}
		}
	}

	// Sort by priority (highest first)
	for i := 0; i < len(detected); i++ {
		for j := i + 1; j < len(detected); j++ {
			if detected[j].Priority > detected[i].Priority {
				detected[i], detected[j] = detected[j], detected[i]
			}
		}
	}

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
