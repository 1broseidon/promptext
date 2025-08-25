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
	Extensions      []string `yaml:"extensions"`
	Excludes        []string `yaml:"excludes"`
	Verbose         bool     `yaml:"verbose"`
	Format          string   `yaml:"format"`            // Add format field: "markdown", "xml", "json"
	Debug           bool     `yaml:"debug"`             // Add debug field
	GitIgnore       bool     `yaml:"gitignore"`         // Use .gitignore patterns
	UseDefaultRules bool     `yaml:"use-default-rules"` // Use default filtering rules (true by default)
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

// mergeExtensions handles extension merging logic
func (fc *FileConfig) mergeExtensions(flagExt string) []string {
	if flagExt != "" {
		return parseCommaSeparated(flagExt)
	}
	if len(fc.Extensions) > 0 {
		return fc.Extensions
	}
	return nil
}

// mergeExcludes handles exclude pattern merging logic
func (fc *FileConfig) mergeExcludes(flagExclude string) []string {
	excludes := append([]string{}, fc.Excludes...)
	if flagExclude != "" {
		excludes = append(excludes, parseCommaSeparated(flagExclude)...)
	}
	return excludes
}

// mergeBooleanFlags handles boolean flag merging with proper precedence
func (fc *FileConfig) mergeBooleanFlags(flagVerbose, flagDebug bool, flagGitIgnore, flagUseDefaultRules *bool) (verbose, debug, useGitIgnore, useDefaultRules bool) {
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

	// Handle use-default-rules - default to true unless explicitly disabled
	useDefaultRules = true
	if flagUseDefaultRules != nil {
		useDefaultRules = *flagUseDefaultRules
	} else if !fc.UseDefaultRules {
		useDefaultRules = false
	}

	return verbose, debug, useGitIgnore, useDefaultRules
}

// MergeWithFlags merges the file config with command line flags
// Command line flags take precedence over file config
func (fc *FileConfig) MergeWithFlags(flagExt, flagExclude string, flagVerbose bool, flagDebug bool, flagGitIgnore *bool, flagUseDefaultRules *bool) (extensions []string, excludes []string, verbose bool, debug bool, useGitIgnore bool, useDefaultRules bool) {
	extensions = fc.mergeExtensions(flagExt)
	excludes = fc.mergeExcludes(flagExclude)
	verbose, debug, useGitIgnore, useDefaultRules = fc.mergeBooleanFlags(flagVerbose, flagDebug, flagGitIgnore, flagUseDefaultRules)

	return extensions, excludes, verbose, debug, useGitIgnore, useDefaultRules
}

func parseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}
