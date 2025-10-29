package rules

import (
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter/types"
	"github.com/1broseidon/promptext/internal/log"
)

// EcosystemRule detects and excludes lock files based on detected package managers
type EcosystemRule struct {
	types.BaseRule
	detectedEcosystems map[string]bool
	lockFileMap        map[string][]string
	rootPath           string
}

// NewEcosystemRule creates an ecosystem-aware lock file detector
func NewEcosystemRule(rootPath string) types.Rule {
	er := &EcosystemRule{
		BaseRule:           types.NewBaseRule("ecosystem", types.Exclude),
		detectedEcosystems: make(map[string]bool),
		rootPath:           rootPath,
		lockFileMap: map[string][]string{
			"node": {
				"package-lock.json",
				"yarn.lock",
				"pnpm-lock.yaml",
				"bun.lockb",
				".pnp.cjs",        // Yarn PnP
				".pnp.loader.mjs", // Yarn PnP loader
			},
			"php": {
				"composer.lock",
			},
			"python": {
				"poetry.lock",
				"Pipfile.lock",
				"pdm.lock",
			},
			"ruby": {
				"Gemfile.lock",
			},
			"rust": {
				"Cargo.lock",
			},
			"go": {
				"go.sum",
			},
			"dotnet": {
				"packages.lock.json",
				"project.assets.json", // MSBuild generated
				"*.nuget.props",
				"*.nuget.targets",
			},
			"java": {
				"gradle.lockfile",
			},
		},
	}

	// Detect ecosystems by scanning for manifest files
	er.detectEcosystems()
	return er
}

// detectEcosystems scans for package manager manifest files
func (er *EcosystemRule) detectEcosystems() {
	manifestIndicators := map[string]string{
		"package.json":     "node",
		"composer.json":    "php",
		"pyproject.toml":   "python",
		"Pipfile":          "python",
		"requirements.txt": "python",
		"Gemfile":          "ruby",
		"Cargo.toml":       "rust",
		"go.mod":           "go",
		"*.csproj":         "dotnet",
		"*.fsproj":         "dotnet",
		"*.vbproj":         "dotnet",
		"build.gradle":     "java",
		"build.gradle.kts": "java",
		"pom.xml":          "java",
	}

	filepath.WalkDir(er.rootPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}

		basename := filepath.Base(path)
		for manifest, ecosystem := range manifestIndicators {
			if matched, _ := filepath.Match(manifest, basename); matched {
				if !er.detectedEcosystems[ecosystem] {
					er.detectedEcosystems[ecosystem] = true
					log.Debug("Detected %s ecosystem via %s", ecosystem, basename)
				}
			}
		}
		return nil
	})

	if len(er.detectedEcosystems) > 0 {
		ecosystems := make([]string, 0, len(er.detectedEcosystems))
		for eco := range er.detectedEcosystems {
			ecosystems = append(ecosystems, eco)
		}
		log.Debug("Active ecosystems: %s", strings.Join(ecosystems, ", "))
	}
}

// Match checks if a file is a lock file for any detected ecosystem
func (er *EcosystemRule) Match(path string) bool {
	basename := filepath.Base(path)

	// Check if this file is a lock file for any detected ecosystem
	for ecosystem := range er.detectedEcosystems {
		if lockFiles, exists := er.lockFileMap[ecosystem]; exists {
			for _, lockFile := range lockFiles {
				// Handle both exact matches and glob patterns
				if matched, _ := filepath.Match(lockFile, basename); matched {
					log.Debug("Excluding %s lock file (ecosystem-aware): %s", ecosystem, path)
					return true
				}
			}
		}
	}

	return false
}
