package initializer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Initializer handles config file initialization
type Initializer struct {
	detector  Detector
	generator *TemplateGenerator
	rootPath  string
	force     bool
	quiet     bool
}

// NewInitializer creates a new initializer
func NewInitializer(rootPath string, force bool, quiet bool) *Initializer {
	return &Initializer{
		detector:  NewFileDetector(),
		generator: NewTemplateGenerator(),
		rootPath:  rootPath,
		force:     force,
		quiet:     quiet,
	}
}

// Run executes the initialization process
func (i *Initializer) Run() error {
	// Validate that rootPath exists and is a directory
	info, err := os.Stat(i.rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", i.rootPath)
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", i.rootPath)
	}

	// Check if config already exists
	configPath := filepath.Join(i.rootPath, ".promptext.yml")
	if _, err := os.Stat(configPath); err == nil && !i.force {
		fmt.Println("⚠️  Configuration file already exists: .promptext.yml")
		fmt.Println()

		if !i.promptConfirm("Do you want to overwrite it?") {
			fmt.Println("❌ Initialization cancelled.")
			return nil
		}
		fmt.Println()
	}

	// Detect project types
	if !i.quiet {
		fmt.Println("🔍 Detecting project type...")
	}
	projectTypes, err := i.detector.Detect(i.rootPath)
	if err != nil {
		return fmt.Errorf("failed to detect project type: %w", err)
	}

	// Display detected types
	if len(projectTypes) == 0 {
		if !i.quiet {
			fmt.Println("📦 No specific framework detected. Using generic configuration.")
			fmt.Println()
		}
	} else {
		if !i.quiet {
			fmt.Println("✅ Detected project type(s):")
			for _, pt := range projectTypes {
				fmt.Printf("   • %s\n", pt.Description)
			}
			fmt.Println()
		}
	}

	// Ask about test files
	includeTests := false
	if !i.quiet {
		includeTests = i.promptConfirm("Do you want to include test files in your output?")
		fmt.Println()
	}

	// Generate template
	template := i.generator.Generate(projectTypes, includeTests)
	yamlContent := i.generator.GenerateYAML(template)

	// Write to file
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	// Success message
	if !i.quiet {
		fmt.Println("✨ Configuration file created successfully!")
		fmt.Println()
		fmt.Printf("📄 Location: %s\n", configPath)
		fmt.Println()
		fmt.Println("📝 Next steps:")
		fmt.Println("   1. Review and customize .promptext.yml to fit your needs")
		fmt.Println("   2. Run 'promptext' to generate your project context")
		fmt.Println()

		if len(template.Extensions) > 0 {
			fmt.Println("📌 Included file extensions:")
			extList := strings.Join(template.Extensions, ", ")
			fmt.Printf("   %s\n", extList)
			fmt.Println()
		}

		fmt.Println("💡 Tips:")
		fmt.Println("   • Use 'promptext --help' to see all available options")
		fmt.Println("   • Edit .promptext.yml to add custom exclusions or extensions")
		fmt.Println("   • The config respects your .gitignore by default")
		fmt.Println()
	}

	return nil
}

// promptConfirm asks a yes/no question and returns the answer
func (i *Initializer) promptConfirm(question string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s (y/n): ", question)
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.TrimSpace(strings.ToLower(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no":
			return false
		default:
			fmt.Println("Please answer 'y' or 'n'")
		}
	}
}

// RunQuick runs initialization with default options (no prompts)
func (i *Initializer) RunQuick() error {
	// Validate that rootPath exists and is a directory
	info, err := os.Stat(i.rootPath)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("directory does not exist: %s", i.rootPath)
		}
		return fmt.Errorf("failed to access directory: %w", err)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", i.rootPath)
	}

	// Check if config already exists
	configPath := filepath.Join(i.rootPath, ".promptext.yml")
	if _, err := os.Stat(configPath); err == nil && !i.force {
		return fmt.Errorf("configuration file already exists: .promptext.yml (use --force to overwrite)")
	}

	// Detect project types
	projectTypes, err := i.detector.Detect(i.rootPath)
	if err != nil {
		return fmt.Errorf("failed to detect project type: %w", err)
	}

	// Generate template (exclude tests by default in quick mode)
	template := i.generator.Generate(projectTypes, false)
	yamlContent := i.generator.GenerateYAML(template)

	// Write to file
	if err := os.WriteFile(configPath, []byte(yamlContent), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	if !i.quiet {
		fmt.Printf("✨ Created .promptext.yml")
		if len(projectTypes) > 0 {
			fmt.Printf(" (detected: %s)", projectTypes[0].Description)
		}
		fmt.Println()
	}

	return nil
}
