package config

import (
	"reflect"
	"testing"
)

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
				Verbose:         false,
				UseDefaultRules: false,
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
				Verbose:         true,
				UseDefaultRules: false,
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
				Verbose:         true,
				UseDefaultRules: true,
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
				Verbose:         false,
				UseDefaultRules: false,
			},
			flagVerbose:         true,
			wantExtensions:      nil,
			wantExcludes:        []string{},
			wantVerbose:         true,
			wantUseDefaultRules: true, // Default to true when not specified
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var useDefaultRules bool = true
			gotExt, gotExc, gotVerb, _, _, gotUseDefaultRules := tt.config.MergeWithFlags(tt.flagExt, tt.flagExclude, tt.flagVerbose, false, nil, &useDefaultRules)

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
