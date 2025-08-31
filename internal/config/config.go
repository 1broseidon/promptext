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
	Verbose         *bool    `yaml:"verbose"`           // Use pointer to distinguish nil (unset) from false
	Format          string   `yaml:"format"`            // Add format field: "markdown", "xml", "json"
	Debug           *bool    `yaml:"debug"`             // Use pointer to distinguish nil (unset) from false
	GitIgnore       *bool    `yaml:"gitignore"`         // Use .gitignore patterns
	UseDefaultRules *bool    `yaml:"use-default-rules"` // Use default filtering rules (true by default)
}

// getGlobalConfigPaths returns potential global config file paths in order of preference
// Follows XDG Base Directory Specification with fallbacks
func getGlobalConfigPaths() []string {
	var paths []string
	
	// XDG_CONFIG_HOME or ~/.config/promptext/config.yml
	configHome := os.Getenv("XDG_CONFIG_HOME")
	if configHome == "" {
		if homeDir, err := os.UserHomeDir(); err == nil {
			configHome = filepath.Join(homeDir, ".config")
		}
	}
	if configHome != "" {
		paths = append(paths, filepath.Join(configHome, "promptext", "config.yml"))
	}
	
	// ~/.promptext.yml (traditional dotfile)
	if homeDir, err := os.UserHomeDir(); err == nil {
		paths = append(paths, filepath.Join(homeDir, ".promptext.yml"))
	}
	
	return paths
}

// LoadGlobalConfig attempts to load global configuration from standard locations
func LoadGlobalConfig() (*FileConfig, error) {
	configPaths := getGlobalConfigPaths()
	
	for _, configPath := range configPaths {
		data, err := os.ReadFile(configPath)
		if err != nil {
			if os.IsNotExist(err) {
				continue // Try next path
			}
			return nil, err
		}
		
		log.Debug("Found and loaded global config from %s", configPath)
		
		var config FileConfig
		if err := yaml.Unmarshal(data, &config); err != nil {
			return nil, err
		}
		
		return &config, nil
	}
	
	log.Debug("No global config found in any of: %v", configPaths)
	return &FileConfig{}, nil
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

// mergeConfigs merges global, project, and flag configurations with proper precedence
// Precedence: CLI flags > Project config > Global config
func MergeConfigs(globalConfig, projectConfig *FileConfig, flagExt, flagExclude string, flagVerbose bool, flagDebug bool, flagGitIgnore *bool, flagUseDefaultRules *bool) (extensions []string, excludes []string, verbose bool, debug bool, useGitIgnore bool, useDefaultRules bool) {
	// Start with global config as base
	merged := &FileConfig{
		Extensions:      append([]string{}, globalConfig.Extensions...),
		Excludes:        append([]string{}, globalConfig.Excludes...),
		Verbose:         globalConfig.Verbose,
		Format:          globalConfig.Format,
		Debug:           globalConfig.Debug,
		GitIgnore:       globalConfig.GitIgnore,
		UseDefaultRules: globalConfig.UseDefaultRules,
	}
	
	// Override with project config where explicitly set
	if len(projectConfig.Extensions) > 0 {
		merged.Extensions = projectConfig.Extensions
	}
	// For excludes, we want to merge (append) rather than replace
	if len(projectConfig.Excludes) > 0 {
		merged.Excludes = append(merged.Excludes, projectConfig.Excludes...)
	}
	if projectConfig.Verbose != nil {
		merged.Verbose = projectConfig.Verbose
	}
	if projectConfig.Format != "" {
		merged.Format = projectConfig.Format
	}
	if projectConfig.Debug != nil {
		merged.Debug = projectConfig.Debug
	}
	if projectConfig.GitIgnore != nil {
		merged.GitIgnore = projectConfig.GitIgnore
	}
	if projectConfig.UseDefaultRules != nil {
		merged.UseDefaultRules = projectConfig.UseDefaultRules
	}
	
	// Finally merge with CLI flags (highest precedence)
	extensions, excludes, verbose, debug, useGitIgnore, useDefaultRules = merged.MergeWithFlags(flagExt, flagExclude, flagVerbose, flagDebug, flagGitIgnore, flagUseDefaultRules)
	
	return extensions, excludes, verbose, debug, useGitIgnore, useDefaultRules
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
	if fc.Verbose != nil {
		verbose = *fc.Verbose || flagVerbose
	} else {
		verbose = flagVerbose
	}

	// Handle debug flag
	if fc.Debug != nil {
		debug = *fc.Debug || flagDebug
	} else {
		debug = flagDebug
	}

	// Handle gitignore - default to true unless explicitly disabled
	useGitIgnore = true
	if flagGitIgnore != nil {
		useGitIgnore = *flagGitIgnore
	} else if fc.GitIgnore != nil {
		useGitIgnore = *fc.GitIgnore
	}

	// Handle use-default-rules - default to true unless explicitly disabled
	useDefaultRules = true
	if flagUseDefaultRules != nil {
		useDefaultRules = *flagUseDefaultRules
	} else if fc.UseDefaultRules != nil {
		useDefaultRules = *fc.UseDefaultRules
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
