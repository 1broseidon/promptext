package config

import (
	"reflect"
	"testing"
)

func TestMergeWithFlags(t *testing.T) {
	tests := []struct {
		name           string
		config         *FileConfig
		flagExt        string
		flagExclude    string
		flagVerbose    bool
		wantExtensions []string
		wantExcludes   []string
		wantVerbose    bool
	}{
		{
			name: "flags override config",
			config: &FileConfig{
				Extensions: []string{".go"},
				Excludes:   []string{"vendor"},
				Verbose:    false,
			},
			flagExt:        ".js,.ts",
			flagExclude:    "node_modules",
			flagVerbose:    true,
			wantExtensions: []string{".js", ".ts"},
			wantExcludes:   []string{"vendor", "node_modules"},
			wantVerbose:    true,
		},
		{
			name: "empty flags use config",
			config: &FileConfig{
				Extensions: []string{".go", ".py"},
				Excludes:   []string{"test"},
				Verbose:    true,
			},
			flagExt:        "",
			flagExclude:    "",
			flagVerbose:    false,
			wantExtensions: []string{".go", ".py"},
			wantExcludes:   []string{"test"},
			wantVerbose:    true,
		},
		{
			name:           "empty config and flags",
			config:         &FileConfig{},
			flagExt:        "",
			flagExclude:    "",
			flagVerbose:    false,
			wantExtensions: nil,
			wantExcludes:   []string{},
			wantVerbose:    false,
		},
		{
			name: "config verbose true with flag false",
			config: &FileConfig{
				Verbose: true,
			},
			flagVerbose:    false,
			wantExtensions: nil,
			wantExcludes:   []string{},
			wantVerbose:    true,
		},
		{
			name: "config verbose false with flag true",
			config: &FileConfig{
				Verbose: false,
			},
			flagVerbose:    true,
			wantExtensions: nil,
			wantExcludes:   []string{},
			wantVerbose:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotExt, gotExc, gotVerb, _, _ := tt.config.MergeWithFlags(tt.flagExt, tt.flagExclude, tt.flagVerbose, false, nil)

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
