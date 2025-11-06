package info

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/1broseidon/promptext/internal/filter"
	"github.com/stretchr/testify/assert"
)

func TestGetProjectInfo(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "project-info-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"go.mod":          "module test\n\ngo 1.17\n\nrequire github.com/stretchr/testify v1.8.0",
		"main.go":         "package main\n\nfunc main() {}\n",
		"README.md":       "# Test Project",
		".gitignore":      "*.tmp\n",
		"internal/foo.go": "package internal\n",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	// Initialize test filter
	f := filter.New(filter.Options{
		Includes: []string{".go"},
		Excludes: []string{},
	})

	// Test GetProjectInfo
	t.Run("basic project structure", func(t *testing.T) {
		info, err := GetProjectInfo(tmpDir, f)
		assert.NoError(t, err)
		assert.NotNil(t, info)
		assert.NotNil(t, info.DirectoryTree)
		// Verify basic structure instead of specific temp dir name
		assert.NotEmpty(t, info.DirectoryTree.Name)
		assert.Equal(t, "dir", info.DirectoryTree.Type)
	})
}

func TestGenerateDirectoryTree(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "directory-tree-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directory structure
	files := []string{
		"main.go",
		"internal/pkg1/file1.go",
		"internal/pkg2/file2.go",
		"docs/README.md",
	}

	for _, file := range files {
		path := filepath.Join(tmpDir, file)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte("test content"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	f := filter.New(filter.Options{
		Includes: []string{".go", ".md"},
		Excludes: []string{},
	})

	t.Run("directory tree generation", func(t *testing.T) {
		tree, err := generateDirectoryTree(tmpDir, f)
		assert.NoError(t, err)
		assert.NotNil(t, tree)

		// Verify root node
		assert.Equal(t, filepath.Base(tmpDir), tree.Name)
		assert.Equal(t, "dir", tree.Type)

		// Verify directory structure
		foundMain := false
		foundInternal := false
		foundDocs := false

		for _, child := range tree.Children {
			switch child.Name {
			case "main.go":
				foundMain = true
			case "internal":
				foundInternal = true
			case "docs":
				foundDocs = true
			}
		}

		assert.True(t, foundMain, "main.go not found")
		assert.True(t, foundInternal, "internal/ not found")
		assert.True(t, foundDocs, "docs/ not found")
	})
}

func TestGetProjectMetadata(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "metadata-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	t.Run("Go project", func(t *testing.T) {
		// Create go.mod file
		goMod := `module test
go 1.17
require github.com/stretchr/testify v1.8.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644)
		assert.NoError(t, err)

		metadata, err := getProjectMetadata(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "Go", metadata.Language)
		assert.Equal(t, "1.17", metadata.Version)
		assert.Contains(t, metadata.Dependencies, "github.com/stretchr/testify")
	})

	t.Run("Node.js project", func(t *testing.T) {
		// Clean up previous files
		os.RemoveAll(tmpDir)
		os.MkdirAll(tmpDir, 0755)

		// Create package.json file
		packageJSON := `{
			"name": "test",
			"version": "1.0.0",
			"dependencies": {
				"express": "^4.17.1"
			}
		}`
		err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageJSON), 0644)
		assert.NoError(t, err)

		metadata, err := getProjectMetadata(tmpDir)
		assert.NoError(t, err)
		assert.Equal(t, "JavaScript/Node.js", metadata.Language)
		assert.Contains(t, metadata.Dependencies, "express")
	})
}

func TestAnalyzeProject(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "analysis-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"main.go":           "package main\n\nfunc main() {}\n",
		"internal/core.go":  "package internal\n",
		"config.yaml":       "key: value\n",
		"README.md":         "# Test Project",
		"test/main_test.go": "package test\n",
		".gitignore":        "*.tmp\n",
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.MkdirAll(filepath.Dir(path), 0755)
		if err != nil {
			t.Fatal(err)
		}
		err = os.WriteFile(path, []byte(content), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	t.Run("project analysis", func(t *testing.T) {
		// Initialize test filter
		f := filter.New(filter.Options{
			UseDefaultRules: true,
			UseGitIgnore:    false,
		})
		analysis := AnalyzeProject(tmpDir, f)
		assert.NotNil(t, analysis)

		// Check entry points
		assert.Contains(t, analysis.EntryPoints, "main.go")

		// Check core files
		assert.Contains(t, analysis.CoreFiles, "internal/core.go")

		// Check config files
		assert.Contains(t, analysis.ConfigFiles, "config.yaml")

		// Check documentation
		assert.Contains(t, analysis.Documentation, "README.md")

		// Check test files
		assert.Contains(t, analysis.TestFiles, "test/main_test.go")
	})
}

func TestGetPythonVersion(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("from pyproject.toml", func(t *testing.T) {
		pyprojectContent := `[tool.poetry]
name = "test-project"

[tool.poetry.dependencies]
python = "^3.9"
requests = "^2.28.0"
`
		err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyprojectContent), 0644)
		assert.NoError(t, err)

		version := getPythonVersion(tmpDir)
		// Function strips "^" character
		assert.Equal(t, "3.9", version)
	})

	t.Run("no python version", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		version := getPythonVersion(tmpDir2)
		assert.Empty(t, version)
	})
}

func TestGetRustVersion(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("from Cargo.toml", func(t *testing.T) {
		cargoContent := `[package]
name = "test-project"
version = "0.1.0"
edition = "2021"
`
		err := os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoContent), 0644)
		assert.NoError(t, err)

		version := getRustVersion(tmpDir)
		assert.Equal(t, "0.1.0", version)
	})

	t.Run("no Cargo.toml", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		version := getRustVersion(tmpDir2)
		assert.Empty(t, version)
	})
}

func TestGetPipDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse requirements.txt", func(t *testing.T) {
		requirementsContent := `requests==2.28.0
pytest==7.2.0
# This is a comment
flask>=2.0.0

django==4.1.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(requirementsContent), 0644)
		assert.NoError(t, err)

		depsMap := make(map[string]bool)
		getPipDependencies(tmpDir, depsMap)

		assert.True(t, depsMap["requests"])
		assert.True(t, depsMap["pytest"])
		// Function splits on "==", so "flask>=2.0.0" (no "==") stays as-is
		assert.True(t, depsMap["flask>=2.0.0"])
		assert.True(t, depsMap["django"])
		assert.Equal(t, 4, len(depsMap))
	})

	t.Run("no requirements.txt", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		depsMap := make(map[string]bool)
		getPipDependencies(tmpDir2, depsMap)
		assert.Equal(t, 0, len(depsMap))
	})
}

func TestGetPoetryDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse pyproject.toml", func(t *testing.T) {
		pyprojectContent := `[tool.poetry.dependencies]
python = "^3.9"
requests = "^2.28.0"
flask = ">=2.0.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.2.0"
black = "^22.0.0"
`
		err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyprojectContent), 0644)
		assert.NoError(t, err)

		depsMap := make(map[string]bool)
		getPoetryDependencies(tmpDir, depsMap)

		assert.True(t, depsMap["requests"])
		assert.True(t, depsMap["flask"])
		assert.True(t, depsMap["[dev] pytest"])
		assert.True(t, depsMap["[dev] black"])
		assert.False(t, depsMap["python"]) // Python version should be excluded
	})

	t.Run("no pyproject.toml", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		depsMap := make(map[string]bool)
		getPoetryDependencies(tmpDir2, depsMap)
		assert.Equal(t, 0, len(depsMap))
	})
}

func TestGetPoetryLockDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse poetry.lock", func(t *testing.T) {
		poetryLockContent := `[[package]]
name = "certifi"
version = "2023.5.7"

[[package]]
name = "charset-normalizer"
version = "3.1.0"

[[package]]
name = "idna"
version = "3.4"
`
		err := os.WriteFile(filepath.Join(tmpDir, "poetry.lock"), []byte(poetryLockContent), 0644)
		assert.NoError(t, err)

		depsMap := make(map[string]bool)
		getPoetryLockDependencies(tmpDir, depsMap)

		assert.True(t, depsMap["certifi"])
		assert.True(t, depsMap["charset-normalizer"])
		assert.True(t, depsMap["idna"])
		assert.Equal(t, 3, len(depsMap))
	})

	t.Run("no poetry.lock", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		depsMap := make(map[string]bool)
		getPoetryLockDependencies(tmpDir2, depsMap)
		assert.Equal(t, 0, len(depsMap))
	})
}

func TestGetPythonDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("combined from all sources", func(t *testing.T) {
		// Create requirements.txt
		requirementsContent := `requests==2.28.0
pytest==7.2.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(requirementsContent), 0644)
		assert.NoError(t, err)

		// Create pyproject.toml
		pyprojectContent := `[tool.poetry.dependencies]
python = "^3.9"
flask = ">=2.0.0"
`
		err = os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyprojectContent), 0644)
		assert.NoError(t, err)

		deps := getPythonDependencies(tmpDir)

		// Should contain deps from both sources
		assert.Contains(t, deps, "requests")
		assert.Contains(t, deps, "pytest")
		assert.Contains(t, deps, "flask")
		assert.GreaterOrEqual(t, len(deps), 3)
	})

	t.Run("no dependency files", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		deps := getPythonDependencies(tmpDir2)
		assert.Equal(t, 0, len(deps))
	})
}

func TestGetRustDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse Cargo.toml", func(t *testing.T) {
		cargoContent := `[package]
name = "test-project"
version = "0.1.0"

[dependencies]
serde = "1.0"
tokio = { version = "1.0", features = ["full"] }
reqwest = "0.11"
`
		err := os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoContent), 0644)
		assert.NoError(t, err)

		deps := getRustDependencies(tmpDir)

		assert.Contains(t, deps, "serde")
		assert.Contains(t, deps, "tokio")
		assert.Contains(t, deps, "reqwest")
		assert.GreaterOrEqual(t, len(deps), 3)
	})

	t.Run("no Cargo.toml", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		deps := getRustDependencies(tmpDir2)
		assert.Nil(t, deps)
	})
}

func TestGetJavaMavenDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse pom.xml", func(t *testing.T) {
		pomContent := `<?xml version="1.0" encoding="UTF-8"?>
<project>
    <dependencies>
        <dependency>
            <groupId>org.springframework.boot</groupId>
            <artifactId>spring-boot-starter-web</artifactId>
        </dependency>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
        </dependency>
    </dependencies>
</project>
`
		err := os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(pomContent), 0644)
		assert.NoError(t, err)

		deps := getJavaMavenDependencies(tmpDir)

		assert.Contains(t, deps, "spring-boot-starter-web")
		assert.Contains(t, deps, "junit")
		assert.GreaterOrEqual(t, len(deps), 2)
	})

	t.Run("no pom.xml", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		deps := getJavaMavenDependencies(tmpDir2)
		assert.Nil(t, deps)
	})
}

func TestGetJavaGradleDependencies(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("parse build.gradle", func(t *testing.T) {
		gradleContent := `plugins {
    id 'java'
}

dependencies {
    implementation 'org.springframework.boot:spring-boot-starter-web:2.7.0'
    implementation 'com.google.guava:guava:31.0-jre'
    testImplementation 'junit:junit:4.13.2'
    runtimeOnly 'com.h2database:h2:2.1.214'
}
`
		err := os.WriteFile(filepath.Join(tmpDir, "build.gradle"), []byte(gradleContent), 0644)
		assert.NoError(t, err)

		deps := getJavaGradleDependencies(tmpDir)

		// Function only parses "implementation" lines, returns full dependency string
		assert.Contains(t, deps, "org.springframework.boot:spring-boot-starter-web:2.7.0")
		assert.Contains(t, deps, "com.google.guava:guava:31.0-jre")
		// testImplementation and runtimeOnly are not parsed
		assert.NotContains(t, deps, "junit:junit:4.13.2")
		assert.GreaterOrEqual(t, len(deps), 2)
	})

	t.Run("no build.gradle", func(t *testing.T) {
		tmpDir2 := t.TempDir()
		deps := getJavaGradleDependencies(tmpDir2)
		assert.Nil(t, deps)
	})
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		filename string
		expected string
	}{
		{"go.mod", "Go"},
		{"package.json", "JavaScript/Node.js"},
		{"requirements.txt", "Python"},
		{"pyproject.toml", "Python"},
		{"poetry.lock", "Python"},
		{"Cargo.toml", "Rust"},
		{"pom.xml", "Java (Maven)"},
		{"build.gradle", "Java (Gradle)"},
		{"unknown.txt", ""},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			lang := detectLanguage(tt.filename)
			assert.Equal(t, tt.expected, lang)
		})
	}
}

func TestGetLanguageVersion(t *testing.T) {
	t.Run("Go version from go.mod", func(t *testing.T) {
		tmpDir := t.TempDir()
		goModContent := `module test

go 1.21

require github.com/test/dep v1.0.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
		assert.NoError(t, err)

		version := getLanguageVersion(tmpDir, "Go")
		assert.Equal(t, "1.21", version)
	})

	t.Run("Python version from pyproject.toml", func(t *testing.T) {
		tmpDir := t.TempDir()
		pyprojectContent := `[tool.poetry.dependencies]
python = "^3.10"
`
		err := os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyprojectContent), 0644)
		assert.NoError(t, err)

		// detectLanguage returns "Python", getLanguageVersion expects exact language string
		version := getLanguageVersion(tmpDir, "Python")
		// Function strips "^" character
		assert.Equal(t, "3.10", version)
	})

	t.Run("JavaScript version from package.json", func(t *testing.T) {
		tmpDir := t.TempDir()
		packageContent := `{
  "name": "test",
  "version": "1.2.3",
  "engines": {
    "node": ">=14.0.0"
  }
}
`
		err := os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte(packageContent), 0644)
		assert.NoError(t, err)

		version := getLanguageVersion(tmpDir, "JavaScript/Node.js")
		// Function looks for "node" field and returns "requires Node X.Y.Z"
		assert.Equal(t, "requires Node >=14.0.0", version)
	})

	t.Run("Rust version from Cargo.toml", func(t *testing.T) {
		tmpDir := t.TempDir()
		cargoContent := `[package]
name = "test"
version = "0.2.5"
`
		err := os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoContent), 0644)
		assert.NoError(t, err)

		version := getLanguageVersion(tmpDir, "Rust")
		assert.Equal(t, "0.2.5", version)
	})

	t.Run("unknown language", func(t *testing.T) {
		tmpDir := t.TempDir()
		version := getLanguageVersion(tmpDir, "Unknown")
		assert.Empty(t, version)
	})
}

func TestGetDependencies(t *testing.T) {
	t.Run("Go dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		goModContent := `module test

go 1.21

require (
	github.com/stretchr/testify v1.8.0
	github.com/gorilla/mux v1.8.0
)
`
		err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goModContent), 0644)
		assert.NoError(t, err)

		// getDependencies takes (root, filename) not (root, language)
		deps := getDependencies(tmpDir, "go.mod")
		assert.Contains(t, deps, "github.com/stretchr/testify")
		assert.Contains(t, deps, "github.com/gorilla/mux")
	})

	t.Run("Python dependencies", func(t *testing.T) {
		tmpDir := t.TempDir()
		requirementsContent := `requests==2.28.0
flask==2.3.0
`
		err := os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(requirementsContent), 0644)
		assert.NoError(t, err)

		deps := getDependencies(tmpDir, "requirements.txt")
		assert.Contains(t, deps, "requests")
		assert.Contains(t, deps, "flask")
	})

	t.Run("unknown filename", func(t *testing.T) {
		tmpDir := t.TempDir()
		deps := getDependencies(tmpDir, "unknown.txt")
		assert.Nil(t, deps)
	})
}
