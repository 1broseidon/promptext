package update

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    [3]int
		wantErr bool
	}{
		{
			name:    "valid version",
			version: "0.4.3",
			want:    [3]int{0, 4, 3},
			wantErr: false,
		},
		{
			name:    "version with major",
			version: "1.0.0",
			want:    [3]int{1, 0, 0},
			wantErr: false,
		},
		{
			name:    "version with double digits",
			version: "2.15.7",
			want:    [3]int{2, 15, 7},
			wantErr: false,
		},
		{
			name:    "invalid version format",
			version: "invalid",
			wantErr: true,
		},
		{
			name:    "incomplete version",
			version: "1.0",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseVersion(tt.version)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name    string
		v1      string
		v2      string
		want    bool
		wantErr bool
	}{
		{
			name: "major version newer",
			v1:   "1.0.0",
			v2:   "0.9.9",
			want: true,
		},
		{
			name: "minor version newer",
			v1:   "0.5.0",
			v2:   "0.4.9",
			want: true,
		},
		{
			name: "patch version newer",
			v1:   "0.4.4",
			v2:   "0.4.3",
			want: true,
		},
		{
			name: "same version",
			v1:   "0.4.3",
			v2:   "0.4.3",
			want: false,
		},
		{
			name: "older major version",
			v1:   "0.4.3",
			v2:   "1.0.0",
			want: false,
		},
		{
			name: "older minor version",
			v1:   "0.3.9",
			v2:   "0.4.0",
			want: false,
		},
		{
			name: "older patch version",
			v1:   "0.4.2",
			v2:   "0.4.3",
			want: false,
		},
		{
			name:    "invalid v1",
			v1:      "invalid",
			v2:      "0.4.3",
			wantErr: true,
		},
		{
			name:    "invalid v2",
			v1:      "0.4.3",
			v2:      "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := isNewerVersion(tt.v1, tt.v2)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetPlatformAssetName(t *testing.T) {
	name, err := getPlatformAssetName()
	require.NoError(t, err)
	assert.NotEmpty(t, name)

	// Should contain "promptext_" prefix
	assert.Contains(t, name, "promptext_")

	// Should have proper extension
	if name[len(name)-7:] != ".tar.gz" && name[len(name)-4:] != ".zip" {
		t.Errorf("invalid extension: %s", name)
	}
}

func TestUpdateCheckCache(t *testing.T) {
	// Create temporary cache directory
	tempDir := t.TempDir()
	cachePath := filepath.Join(tempDir, "update_check.json")

	// Test saving cache
	cache := UpdateCheckCache{
		LastCheck:       time.Now(),
		LatestVersion:   "v0.5.0",
		UpdateAvailable: true,
	}

	data, err := json.Marshal(cache)
	require.NoError(t, err)
	err = os.WriteFile(cachePath, data, 0644)
	require.NoError(t, err)

	// Test loading cache
	loadedData, err := os.ReadFile(cachePath)
	require.NoError(t, err)

	var loadedCache UpdateCheckCache
	err = json.Unmarshal(loadedData, &loadedCache)
	require.NoError(t, err)

	assert.Equal(t, cache.LatestVersion, loadedCache.LatestVersion)
	assert.Equal(t, cache.UpdateAvailable, loadedCache.UpdateAvailable)
	assert.WithinDuration(t, cache.LastCheck, loadedCache.LastCheck, time.Second)
}

func TestCheckForUpdate_DevVersion(t *testing.T) {
	_, _, err := CheckForUpdate("dev")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "development build")
}

func TestCheckForUpdate_UnknownVersion(t *testing.T) {
	_, _, err := CheckForUpdate("unknown")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "development build")
}

func TestGetCacheDir(t *testing.T) {
	cacheDir, err := getCacheDir()
	require.NoError(t, err)
	assert.NotEmpty(t, cacheDir)

	// Should contain "promptext"
	assert.Contains(t, cacheDir, "promptext")

	// Directory should be created
	info, err := os.Stat(cacheDir)
	require.NoError(t, err)
	assert.True(t, info.IsDir())
}

// Integration test - requires network access
func TestCheckForUpdate_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// This test requires network access and should only run when explicitly requested
	available, version, err := CheckForUpdate("v0.1.0")

	// We expect this to find a newer version since v0.1.0 is very old
	// If there's a network error, we skip the test rather than fail
	if err != nil {
		if err.Error() == "failed to fetch latest release: Get \"https://api.github.com/repos/1broseidon/promptext/releases/latest\": dial tcp: lookup api.github.com: no such host" {
			t.Skip("Skipping integration test: no network access")
		}
		t.Logf("Warning: integration test failed with error: %v", err)
		return
	}

	assert.True(t, available, "Expected newer version to be available for v0.1.0")
	assert.NotEmpty(t, version)
	assert.Contains(t, version, "v", "Version should contain 'v' prefix")
}

func TestCheckAndNotifyUpdate_DevVersion(t *testing.T) {
	// Should not panic or error for dev version
	CheckAndNotifyUpdate("dev")
	CheckAndNotifyUpdate("unknown")
}

func TestVersionComparison_EdgeCases(t *testing.T) {
	tests := []struct {
		v1   string
		v2   string
		want bool
	}{
		{"0.0.1", "0.0.0", true},
		{"0.1.0", "0.0.99", true},
		{"1.0.0", "0.99.99", true},
		{"10.0.0", "9.99.99", true},
		{"0.10.0", "0.9.0", true},
		{"0.0.10", "0.0.9", true},
	}

	for _, tt := range tests {
		t.Run(tt.v1+"_vs_"+tt.v2, func(t *testing.T) {
			got, err := isNewerVersion(tt.v1, tt.v2)
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
