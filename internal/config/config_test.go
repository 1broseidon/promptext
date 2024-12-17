package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/1broseidon/promptext/internal/filter"
)

func TestLoadConfig(t *testing.T) {
	// Create a temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name     string
		content  string
		wantErr  bool
		expected *FileConfig
	}{
		{
			name: "valid config",
			content: `
extensions:
  - .go
  - .js
excludes:
  - vendor
  - node_modules
verbose: true
format: json
`,
			wantErr: false,
			expected: &FileConfig{
				Extensions: []string{".go", ".js"},
				Excludes:   []string{"vendor", "node_modules"},
				Verbose:    true,
				Format:     "json",
			},
		},
		{
			name:     "empty config",
			content:  "",
			wantErr:  false,
			expected: &FileConfig{},
		},
		{
			name: "invalid yaml",
			content: `
extensions: [
  - not valid yaml
`,
			wantErr:  true,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create config file
			configPath := filepath.Join(tmpDir, ".promptext.yml")
			if err := os.WriteFile(configPath, []byte(tt.content), 0644); err != nil {
				t.Fatalf("Failed to write config file: %v", err)
			}

			got, err := LoadConfig(tmpDir)
			if (err != nil) != tt.wantErr {
				t.Errorf("LoadConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("LoadConfig() = %v, want %v", got, tt.expected)
			}
		})

		// Clean up config file between tests
		os.Remove(filepath.Join(tmpDir, ".promptext.yml"))
	}
}

func TestMergeWithFlags(t *testing.T) {
	tests := []struct {
		name            string
		config          *FileConfig
		flagExt         string
		flagExclude     string
		flagVerbose     bool
		wantExtensions  []string
		wantExcludes    []string
		wantVerbose     bool
	}{
		{
			name: "flags override config",
			config: &FileConfig{
				Extensions: []string{".go"},
				Excludes:   []string{"vendor"},
				Verbose:    false,
			},
			flagExt:         ".js,.ts",
			flagExclude:     "node_modules",
			flagVerbose:     true,
			wantExtensions:  []string{".js", ".ts"},
			wantExcludes:    []string{"vendor", "node_modules"},
			wantVerbose:     true,
		},
		{
			name: "empty flags use config",
			config: &FileConfig{
				Extensions: []string{".go", ".py"},
				Excludes:   []string{"test"},
				Verbose:    true,
			},
			flagExt:         "",
			flagExclude:     "",
			flagVerbose:     false,
			wantExtensions:  []string{".go", ".py"},
			wantExcludes:    append(filter.DefaultIgnoreDirs, "test"),
			wantVerbose:     true, // Changed: respect config verbose when flag not set
		},
		{
			name:            "empty config and flags",
			config:          &FileConfig{},
			flagExt:         "",
			flagExclude:     "",
			flagVerbose:     false,
			wantExtensions:  nil,
			wantExcludes:    []string{},
			wantVerbose:     false,
		},
		{
			name: "config verbose true with flag false",
			config: &FileConfig{
				Verbose: true,
			},
			flagVerbose: false,
			wantVerbose: true,
		},
		{
			name: "config verbose false with flag true",
			config: &FileConfig{
				Verbose: false,
			},
			flagVerbose: true,
			wantVerbose: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, gotExc, gotVerb := tt.config.MergeWithFlags(tt.flagExt, tt.flagExclude, tt.flagVerbose)

			if !reflect.DeepEqual(gotExt, tt.wantExtensions) {
				t.Errorf("Extensions = %v, want %v", gotExt, tt.wantExtensions)
			}

			if !reflect.DeepEqual(gotExc, tt.wantExcludes) {
				t.Errorf("Excludes = %v, want %v", gotExc, tt.wantExcludes)
			}

			if gotVerb != tt.wantVerbose {
				t.Errorf("Verbose = %v, want %v", gotVerb, tt.wantVerbose)
			}
		})
	}
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
			input: ".go",
			want:  []string{".go"},
		},
		{
			name:  "multiple values",
			input: ".go,.js,.py",
			want:  []string{".go", ".js", ".py"},
		},
		{
			name:  "with spaces",
			input: " .go, .js,.py ",
			want:  []string{" .go", " .js", ".py "},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseCommaSeparated(tt.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseCommaSeparated() = %v, want %v", got, tt.want)
			}
		})
	}
}
