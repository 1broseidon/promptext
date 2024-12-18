package config

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/1broseidon/promptext/internal/log"
	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of .promptext.yml
type FileConfig struct {
	Extensions []string `yaml:"extensions"`
	Excludes   []string `yaml:"excludes"`
	Verbose    bool     `yaml:"verbose"`
	Format     string   `yaml:"format"` // Add format field: "markdown", "xml", "json"
	Debug      bool     `yaml:"debug"`  // Add debug field
	GitIgnore  bool     `yaml:"gitignore"` // Use .gitignore patterns
	GitIgnore  bool     `yaml:"gitignore"` // Add gitignore support
}

// LoadConfig attempts to load and parse the .promptext.yml file
func LoadConfig(dirPath string) (*FileConfig, error) {
	configPath := filepath.Join(dirPath, ".promptext.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Debug("No .promptext.yml found in %s", dirPath)
			return &FileConfig{}, nil
		}
		return nil, err
	}
	log.Debug("Found and loaded .promptext.yml from %s", dirPath)

	var config FileConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// MergeWithFlags merges the file config with command line flags
// Command line flags take precedence over file config
func (fc *FileConfig) MergeWithFlags(flagExt, flagExclude string, flagVerbose bool, flagDebug bool, flagGitIgnore *bool) (extensions []string, excludes []string, verbose bool, debug bool, useGitIgnore bool) {
	// Handle extensions
	if flagExt != "" {
		extensions = parseCommaSeparated(flagExt)
	} else if len(fc.Extensions) > 0 {
		extensions = fc.Extensions
	} else {
		extensions = nil
	}

	// Handle excludes
	excludes = append([]string{}, fc.Excludes...)
	if flagExclude != "" {
		excludes = append(excludes, parseCommaSeparated(flagExclude)...)
	}

	// Flag verbose overrides config verbose only if set
	verbose = fc.Verbose || flagVerbose

	// Handle debug flag
	debug = fc.Debug || flagDebug

	// Handle gitignore - default to true unless explicitly disabled
	useGitIgnore = true
	if flagGitIgnore != nil {
		useGitIgnore = *flagGitIgnore
	} else if !fc.GitIgnore {
		useGitIgnore = false
	}

	return extensions, excludes, verbose, debug, useGitIgnore
}

func parseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}
