package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// FileConfig represents the structure of .promptext.yml
type FileConfig struct {
	Extensions []string `yaml:"extensions"`
	Excludes   []string `yaml:"excludes"`
	Verbose    bool     `yaml:"verbose"`
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
func (fc *FileConfig) MergeWithFlags(flagExt, flagExclude string, flagVerbose bool) ([]string, []string, bool) {
	extensions := fc.Extensions
	if flagExt != "" {
		// Override with flag extensions
		extensions = parseCommaSeparated(flagExt)
	}

	excludes := fc.Excludes
	if flagExclude != "" {
		// Override with flag excludes
		excludes = parseCommaSeparated(flagExclude)
	}

	verbose := fc.Verbose
	if flagVerbose {
		// Override with flag verbose
		verbose = flagVerbose
	}

	return extensions, excludes, verbose
}

func parseCommaSeparated(input string) []string {
	if input == "" {
		return nil
	}
	return strings.Split(input, ",")
}
