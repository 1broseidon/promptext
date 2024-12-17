package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/filter"
	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of .promptext.yml
type FileConfig struct {
	Extensions []string `yaml:"extensions"`
	Excludes   []string `yaml:"excludes"`
	Verbose    bool     `yaml:"verbose"`
	Format     string   `yaml:"format"` // Add format field: "markdown", "xml", "json"
}

// LoadConfig attempts to load and parse the .promptext.yml file
func LoadConfig(dirPath string) (*FileConfig, error) {
	configPath := filepath.Join(dirPath, ".promptext.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &FileConfig{}, nil
		}
		return nil, err
	}

	var config FileConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// MergeWithFlags merges the file config with command line flags
// Command line flags take precedence over file config
func (fc *FileConfig) MergeWithFlags(flagExt, flagExclude string, flagVerbose bool) (extensions []string, excludes []string, verbose bool) {
	extensions = fc.Extensions
	if flagExt != "" {
		// Override with flag extensions
		extensions = parseCommaSeparated(flagExt)
	}

	// Start with default excludes from filter package
	excludes = append([]string{}, filter.DefaultIgnoreDirs...)

	// Add config file excludes
	excludes = append(excludes, fc.Excludes...)

	// Add flag excludes (highest priority)
	if flagExclude != "" {
		excludes = append(excludes, parseCommaSeparated(flagExclude)...)
	}

	// Flag verbose overrides config verbose only if set
	verbose = fc.Verbose || flagVerbose

	return extensions, excludes, verbose
}

func parseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}
