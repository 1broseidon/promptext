package initializer

import (
	"fmt"
	"strings"
)

// ConfigTemplate represents a configuration template for a project type
type ConfigTemplate struct {
	Extensions []string
	Excludes   []string
	Comments   map[string]string // Key -> comment explaining the setting
}

// TemplateGenerator generates configuration templates based on project types
type TemplateGenerator struct{}

// NewTemplateGenerator creates a new template generator
func NewTemplateGenerator() *TemplateGenerator {
	return &TemplateGenerator{}
}

// Generate creates a configuration template based on detected project types
func (g *TemplateGenerator) Generate(projectTypes []ProjectType, includeTests bool) *ConfigTemplate {
	template := &ConfigTemplate{
		Extensions: []string{},
		Excludes:   []string{},
		Comments:   make(map[string]string),
	}

	// Track what we've added to avoid duplicates
	extSet := make(map[string]bool)
	excSet := make(map[string]bool)

	// Add base exclusions that are common to all projects
	baseExcludes := []string{
		"**/.git/**",
		"**/.svn/**",
		"**/.hg/**",
		"**/.DS_Store",
	}
	for _, exc := range baseExcludes {
		if !excSet[exc] {
			template.Excludes = append(template.Excludes, exc)
			excSet[exc] = true
		}
	}

	// Process each detected project type
	for _, pt := range projectTypes {
		switch pt.Name {
		case "nextjs":
			g.addNextJS(template, extSet, excSet, includeTests)
		case "nuxt":
			g.addNuxt(template, extSet, excSet, includeTests)
		case "vite":
			g.addVite(template, extSet, excSet, includeTests)
		case "vue":
			g.addVue(template, extSet, excSet, includeTests)
		case "angular":
			g.addAngular(template, extSet, excSet, includeTests)
		case "svelte":
			g.addSvelte(template, extSet, excSet, includeTests)
		case "node":
			g.addNode(template, extSet, excSet, includeTests)
		case "go":
			g.addGo(template, extSet, excSet, includeTests)
		case "django":
			g.addDjango(template, extSet, excSet, includeTests)
		case "flask":
			g.addFlask(template, extSet, excSet, includeTests)
		case "python":
			g.addPython(template, extSet, excSet, includeTests)
		case "rust":
			g.addRust(template, extSet, excSet, includeTests)
		case "maven", "gradle":
			g.addJava(template, extSet, excSet, includeTests)
		case "ruby":
			g.addRuby(template, extSet, excSet, includeTests)
		case "php", "laravel":
			g.addPHP(template, extSet, excSet, includeTests)
		case "dotnet":
			g.addDotNet(template, extSet, excSet, includeTests)
		}
	}

	// Add comments
	template.Comments["header"] = "Promptext Configuration File"
	template.Comments["extensions"] = "File extensions to include when processing the project"
	template.Comments["excludes"] = "Patterns to exclude (supports glob patterns)"

	return template
}

// Helper functions for each framework

func (g *TemplateGenerator) addNextJS(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	// Extensions
	exts := []string{".js", ".jsx", ".ts", ".tsx", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	// Excludes
	excludes := []string{
		"**/node_modules/**",
		"**/.next/**",
		"**/out/**",
		"**/dist/**",
		"**/build/**",
		"**/.vercel/**",
		"**/.turbo/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.ts", "**/*.test.tsx", "**/*.test.js", "**/*.test.jsx", "**/*.spec.ts", "**/*.spec.tsx", "**/*.spec.js", "**/*.spec.jsx", "**/__tests__/**", "**/__mocks__/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addNuxt(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".js", ".vue", ".ts", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/.nuxt/**",
		"**/.output/**",
		"**/dist/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.js", "**/*.test.ts", "**/*.spec.js", "**/*.spec.ts")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addVite(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".js", ".jsx", ".ts", ".tsx", ".vue", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/dist/**",
		"**/build/**",
		"**/.vite/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.*", "**/*.spec.*")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addVue(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".js", ".vue", ".ts", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/dist/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.js", "**/*.spec.js")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addAngular(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".ts", ".html", ".css", ".scss", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/dist/**",
		"**/.angular/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.spec.ts")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addSvelte(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".js", ".ts", ".svelte", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/.svelte-kit/**",
		"**/build/**",
		"**/dist/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.*", "**/*.spec.*")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addNode(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".js", ".ts", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/node_modules/**",
		"**/dist/**",
		"**/build/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/*.test.js", "**/*.test.ts", "**/*.spec.js", "**/*.spec.ts")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addGo(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".go", ".mod", ".sum", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/vendor/**",
		"**/bin/**",
		"**/dist/**",
		"**/*.exe",
		"**/*.dll",
		"**/*.so",
		"**/*.dylib",
	}
	if !includeTests {
		excludes = append(excludes, "**/*_test.go", "**/testdata/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addDjango(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".py", ".html", ".css", ".js", ".json", ".md", ".txt"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/__pycache__/**",
		"**/*.pyc",
		"**/.pytest_cache/**",
		"**/venv/**",
		"**/env/**",
		"**/.venv/**",
		"**/staticfiles/**",
		"**/media/**",
		"**/db.sqlite3",
		"**/*.log",
	}
	if !includeTests {
		excludes = append(excludes, "**/test_*.py", "**/*_test.py", "**/tests/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addFlask(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".py", ".html", ".css", ".js", ".json", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/__pycache__/**",
		"**/*.pyc",
		"**/.pytest_cache/**",
		"**/venv/**",
		"**/env/**",
		"**/.venv/**",
		"**/instance/**",
		"**/*.log",
	}
	if !includeTests {
		excludes = append(excludes, "**/test_*.py", "**/*_test.py", "**/tests/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addPython(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".py", ".pyi", ".md", ".txt"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/__pycache__/**",
		"**/*.pyc",
		"**/.pytest_cache/**",
		"**/venv/**",
		"**/env/**",
		"**/.venv/**",
		"**/.tox/**",
		"**/dist/**",
		"**/build/**",
		"**/*.egg-info/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/test_*.py", "**/*_test.py", "**/tests/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addRust(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".rs", ".toml", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/target/**",
		"**/Cargo.lock",
	}
	// Rust tests are typically in the same files with #[test] attributes
	// or in tests/ directory
	if !includeTests {
		excludes = append(excludes, "**/tests/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addJava(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".java", ".kt", ".kts", ".xml", ".properties", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/target/**",
		"**/build/**",
		"**/.gradle/**",
		"**/bin/**",
		"**/out/**",
		"**/*.class",
		"**/*.jar",
		"**/*.war",
	}
	if !includeTests {
		excludes = append(excludes, "**/src/test/**", "**/*Test.java", "**/*Test.kt")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addRuby(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".rb", ".erb", ".rake", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/vendor/**",
		"**/tmp/**",
		"**/log/**",
		"**/.bundle/**",
		"**/coverage/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/spec/**", "**/test/**", "**/*_spec.rb", "**/*_test.rb")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addPHP(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".php", ".blade.php", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/vendor/**",
		"**/storage/**",
		"**/bootstrap/cache/**",
		"**/node_modules/**",
		"**/public/build/**",
		"**/public/hot/**",
	}
	if !includeTests {
		excludes = append(excludes, "**/tests/**", "**/*Test.php")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

func (g *TemplateGenerator) addDotNet(t *ConfigTemplate, extSet, excSet map[string]bool, includeTests bool) {
	exts := []string{".cs", ".fs", ".vb", ".csproj", ".fsproj", ".vbproj", ".sln", ".md"}
	for _, ext := range exts {
		if !extSet[ext] {
			t.Extensions = append(t.Extensions, ext)
			extSet[ext] = true
		}
	}

	excludes := []string{
		"**/bin/**",
		"**/obj/**",
		"**/packages/**",
		"**/.vs/**",
		"**/*.dll",
		"**/*.exe",
		"**/*.pdb",
	}
	if !includeTests {
		excludes = append(excludes, "**/*Tests/**", "**/*Test.cs", "**/*.Tests/**")
	}

	for _, exc := range excludes {
		if !excSet[exc] {
			t.Excludes = append(t.Excludes, exc)
			excSet[exc] = true
		}
	}
}

// GenerateYAML creates a YAML string from the template
func (g *TemplateGenerator) GenerateYAML(template *ConfigTemplate) string {
	var sb strings.Builder

	// Header comment
	sb.WriteString("# " + template.Comments["header"] + "\n")
	sb.WriteString("# Auto-generated by: promptext --init\n")
	sb.WriteString("# Learn more: https://github.com/1broseidon/promptext\n\n")

	// Extensions
	if len(template.Extensions) > 0 {
		sb.WriteString("# " + template.Comments["extensions"] + "\n")
		sb.WriteString("extensions:\n")
		for _, ext := range template.Extensions {
			sb.WriteString(fmt.Sprintf("  - %s\n", ext))
		}
		sb.WriteString("\n")
	}

	// Excludes
	if len(template.Excludes) > 0 {
		sb.WriteString("# " + template.Comments["excludes"] + "\n")
		sb.WriteString("excludes:\n")
		for _, exc := range template.Excludes {
			sb.WriteString(fmt.Sprintf("  - \"%s\"\n", exc))
		}
		sb.WriteString("\n")
	}

	// Other settings
	sb.WriteString("# Use .gitignore patterns for additional filtering\n")
	sb.WriteString("gitignore: true\n\n")

	sb.WriteString("# Use built-in filtering rules for common files (node_modules, etc.)\n")
	sb.WriteString("use-default-rules: true\n\n")

	sb.WriteString("# Output format: ptx, markdown, xml, jsonl, or toon\n")
	sb.WriteString("format: ptx\n\n")

	sb.WriteString("# Enable verbose output\n")
	sb.WriteString("verbose: false\n\n")

	sb.WriteString("# Enable debug mode\n")
	sb.WriteString("debug: false\n")

	return sb.String()
}
