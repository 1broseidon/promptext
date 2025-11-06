package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// Helper functions for creating boolean pointers
func boolPtr(b bool) *bool {
	return &b
}

func TestMergeWithFlags(t *testing.T) {
	tests := []struct {
		name                string
		config              *FileConfig
		flagExt             string
		flagExclude         string
		flagVerbose         bool
		wantExtensions      []string
		wantExcludes        []string
		wantVerbose         bool
		wantUseDefaultRules bool
	}{
		{
			name: "flags override config",
			config: &FileConfig{
				Extensions:      []string{".go"},
				Excludes:        []string{"vendor"},
				Verbose:         boolPtr(false),
				UseDefaultRules: boolPtr(false),
			},
			flagExt:             ".js,.ts",
			flagExclude:         "node_modules",
			flagVerbose:         true,
			wantExtensions:      []string{".js", ".ts"},
			wantExcludes:        []string{"vendor", "node_modules"},
			wantVerbose:         true,
			wantUseDefaultRules: true, // Default to true when flag is set
		},
		{
			name: "empty flags use config",
			config: &FileConfig{
				Extensions:      []string{".go", ".py"},
				Excludes:        []string{"test"},
				Verbose:         boolPtr(true),
				UseDefaultRules: boolPtr(false),
			},
			flagExt:             "",
			flagExclude:         "",
			flagVerbose:         false,
			wantExtensions:      []string{".go", ".py"},
			wantExcludes:        []string{"test"},
			wantVerbose:         true,
			wantUseDefaultRules: false,
		},
		{
			name:                "empty config and flags",
			config:              &FileConfig{},
			flagExt:             "",
			flagExclude:         "",
			flagVerbose:         false,
			wantExtensions:      nil,
			wantExcludes:        []string{},
			wantVerbose:         false,
			wantUseDefaultRules: true, // Default to true for empty config
		},
		{
			name: "config verbose true with flag false",
			config: &FileConfig{
				Verbose:         boolPtr(true),
				UseDefaultRules: boolPtr(true),
			},
			flagVerbose:         false,
			wantExtensions:      nil,
			wantExcludes:        []string{},
			wantVerbose:         true,
			wantUseDefaultRules: true,
		},
		{
			name: "config verbose false with flag true",
			config: &FileConfig{
				Verbose:         boolPtr(false),
				UseDefaultRules: boolPtr(false),
			},
			flagVerbose:         true,
			wantExtensions:      nil,
			wantExcludes:        []string{},
			wantVerbose:         true,
			wantUseDefaultRules: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var useDefaultRules *bool
			if tt.name == "flags override config" || tt.name == "empty config and flags" {
				trueVal := true
				useDefaultRules = &trueVal
			}
			gotExt, gotExc, gotVerb, _, _, gotUseDefaultRules := tt.config.MergeWithFlags(tt.flagExt, tt.flagExclude, tt.flagVerbose, false, nil, useDefaultRules)

			if !reflect.DeepEqual(gotExt, tt.wantExtensions) {
				t.Errorf("Extensions = %v, want %v", gotExt, tt.wantExtensions)
			}

			if !reflect.DeepEqual(gotExc, tt.wantExcludes) {
				t.Errorf("Excludes = %v, want %v", gotExc, tt.wantExcludes)
			}

			if gotVerb != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", gotVerb, tt.wantVerbose)
			}

			if gotUseDefaultRules != tt.wantUseDefaultRules {
				t.Errorf("UseDefaultRules = %v, want %v", gotUseDefaultRules, tt.wantUseDefaultRules)
			}
		})
	}
}

func TestMergeConfigs(t *testing.T) {
	tests := []struct {
		name                string
		globalConfig        *FileConfig
		projectConfig       *FileConfig
		flagExt             string
		flagExclude         string
		flagVerbose         bool
		flagDebug           bool
		wantExtensions      []string
		wantExcludes        []string
		wantVerbose         bool
		wantDebug           bool
		wantUseGitIgnore    bool
		wantUseDefaultRules bool
	}{
		{
			name: "global config only",
			globalConfig: &FileConfig{
				Extensions:      []string{".go", ".js"},
				Excludes:        []string{"node_modules"},
				Verbose:         boolPtr(true),
				Debug:           boolPtr(false),
				GitIgnore:       boolPtr(true),
				UseDefaultRules: boolPtr(true),
			},
			projectConfig:       &FileConfig{},
			wantExtensions:      []string{".go", ".js"},
			wantExcludes:        []string{"node_modules"},
			wantVerbose:         true,
			wantDebug:           false,
			wantUseGitIgnore:    true,
			wantUseDefaultRules: true,
		},
		{
			name: "project config overrides global",
			globalConfig: &FileConfig{
				Extensions:      []string{".go"},
				Excludes:        []string{"vendor"},
				Verbose:         boolPtr(false),
				Debug:           boolPtr(false),
				GitIgnore:       boolPtr(true),
				UseDefaultRules: boolPtr(true),
			},
			projectConfig: &FileConfig{
				Extensions:      []string{".js", ".ts"},
				Excludes:        []string{"node_modules"},
				Verbose:         boolPtr(true),
				Debug:           boolPtr(true),
				GitIgnore:       boolPtr(false),
				UseDefaultRules: boolPtr(false),
			},
			wantExtensions:      []string{".js", ".ts"},
			wantExcludes:        []string{"vendor", "node_modules"}, // Merged, not replaced
			wantVerbose:         true,
			wantDebug:           true,
			wantUseGitIgnore:    false,
			wantUseDefaultRules: false,
		},
		{
			name: "CLI flags override everything",
			globalConfig: &FileConfig{
				Extensions: []string{".go"},
				Excludes:   []string{"vendor"},
				Verbose:    boolPtr(false),
				Debug:      boolPtr(false),
			},
			projectConfig: &FileConfig{
				Extensions: []string{".js"},
				Verbose:    boolPtr(false),
				Debug:      boolPtr(false),
			},
			flagExt:             ".py,.rb",
			flagExclude:         "tmp",
			flagVerbose:         true,
			flagDebug:           true,
			wantExtensions:      []string{".py", ".rb"},
			wantExcludes:        []string{"vendor", "tmp"},
			wantVerbose:         true,
			wantDebug:           true,
			wantUseGitIgnore:    true, // Default
			wantUseDefaultRules: true, // Default
		},
		{
			name: "partial project config override",
			globalConfig: &FileConfig{
				Extensions:      []string{".go", ".js"},
				Excludes:        []string{"vendor", "tmp"},
				Verbose:         boolPtr(false),
				Debug:           boolPtr(false),
				GitIgnore:       boolPtr(true),
				UseDefaultRules: boolPtr(true),
			},
			projectConfig: &FileConfig{
				Verbose: boolPtr(true), // Only override verbose
			},
			wantExtensions:      []string{".go", ".js"},    // From global
			wantExcludes:        []string{"vendor", "tmp"}, // From global
			wantVerbose:         true,                      // From project
			wantDebug:           false,                     // From global
			wantUseGitIgnore:    true,                      // From global
			wantUseDefaultRules: true,                      // From global
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, gotExc, gotVerb, gotDebug, gotGitIgnore, gotDefaultRules := MergeConfigs(
				tt.globalConfig, tt.projectConfig, tt.flagExt, tt.flagExclude,
				tt.flagVerbose, tt.flagDebug, nil, nil,
			)

			if !reflect.DeepEqual(gotExt, tt.wantExtensions) {
				t.Errorf("Extensions = %v, want %v", gotExt, tt.wantExtensions)
			}

			if !reflect.DeepEqual(gotExc, tt.wantExcludes) {
				t.Errorf("Excludes = %v, want %v", gotExc, tt.wantExcludes)
			}

			if gotVerb != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", gotVerb, tt.wantVerbose)
			}

			if gotDebug != tt.wantDebug {
				t.Errorf("Debug = %v, want %v", gotDebug, tt.wantDebug)
			}

			if gotGitIgnore != tt.wantUseGitIgnore {
				t.Errorf("GitIgnore = %v, want %v", gotGitIgnore, tt.wantUseGitIgnore)
			}

			if gotDefaultRules != tt.wantUseDefaultRules {
				t.Errorf("UseDefaultRules = %v, want %v", gotDefaultRules, tt.wantUseDefaultRules)
			}
		})
	}
}

func TestGetGlobalConfigPaths(t *testing.T) {
	// Save original env vars
	originalXDGConfigHome := os.Getenv("XDG_CONFIG_HOME")

	defer func() {
		// Restore original env vars
		if originalXDGConfigHome == "" {
			os.Unsetenv("XDG_CONFIG_HOME")
		} else {
			os.Setenv("XDG_CONFIG_HOME", originalXDGConfigHome)
		}
	}()

	t.Run("with XDG_CONFIG_HOME set", func(t *testing.T) {
		testConfigHome := "/tmp/test-config"
		os.Setenv("XDG_CONFIG_HOME", testConfigHome)

		paths := getGlobalConfigPaths()

		expectedXDGPath := filepath.Join(testConfigHome, "promptext", "config.yml")
		if len(paths) < 1 || paths[0] != expectedXDGPath {
			t.Errorf("First path should be XDG config path %s, got %v", expectedXDGPath, paths)
		}

		// Should also include home dotfile as fallback
		if len(paths) < 2 {
			t.Errorf("Should have at least 2 paths, got %d", len(paths))
		}
	})

	t.Run("without XDG_CONFIG_HOME set", func(t *testing.T) {
		os.Unsetenv("XDG_CONFIG_HOME")

		paths := getGlobalConfigPaths()

		if len(paths) < 1 {
			t.Errorf("Should have at least 1 path, got %d", len(paths))
		}

		// Should include both ~/.config/promptext/config.yml and ~/.promptext.yml
		if len(paths) >= 2 {
			// Second path should be the dotfile in home directory
			homeDir, _ := os.UserHomeDir()
			expectedDotfilePath := filepath.Join(homeDir, ".promptext.yml")
			if paths[1] != expectedDotfilePath {
				t.Errorf("Second path should be %s, got %s", expectedDotfilePath, paths[1])
			}
		}
	})
}

func TestParseCommaSeparated(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "empty string",
			input: "",
			want:  nil,
		},
		{
			name:  "single value",
			input: "value1",
			want:  []string{"value1"},
		},
		{
			name:  "multiple values",
			input: "value1,value2,value3",
			want:  []string{"value1", "value2", "value3"},
		},
		{
			name:  "values with spaces",
			input: "value 1,value 2,value 3",
			want:  []string{"value 1", "value 2", "value 3"},
		},
		{
			name:  "file extensions",
			input: ".go,.js,.py",
			want:  []string{".go", ".js", ".py"},
		},
		{
			name:  "paths with commas",
			input: "path/to/file,another/path",
			want:  []string{"path/to/file", "another/path"},
		},
		{
			name:  "single character values",
			input: "a,b,c",
			want:  []string{"a", "b", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseCommaSeparated(tt.input)
			if !stringSlicesEqual(got, tt.want) {
				t.Errorf("parseCommaSeparated(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

// Helper function to compare string slices
func stringSlicesEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
