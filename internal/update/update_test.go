package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

func writeTarGz(t *testing.T, fileName string, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gw)
	hdr := &tar.Header{
		Name: fileName,
		Mode: 0644,
		Size: int64(len(data)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("write header: %v", err)
	}
	if _, err := tw.Write(data); err != nil {
		t.Fatalf("write data: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("close tar: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("close gzip: %v", err)
	}
	return buf.Bytes()
}

func writeZip(t *testing.T, fileName string, data []byte) []byte {
	t.Helper()
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	w, err := zw.Create(fileName)
	if err != nil {
		t.Fatalf("create zip file: %v", err)
	}
	if _, err := w.Write(data); err != nil {
		t.Fatalf("write zip data: %v", err)
	}
	if err := zw.Close(); err != nil {
		t.Fatalf("close zip: %v", err)
	}
	return buf.Bytes()
}

func sha256Hex(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

func captureStderr(t *testing.T, fn func()) string {
	t.Helper()
	orig := os.Stderr
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stderr = w
	fn()
	w.Close()
	os.Stderr = orig
	var buf bytes.Buffer
	if _, err := io.Copy(&buf, r); err != nil {
		t.Fatalf("copy stderr: %v", err)
	}
	return buf.String()
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func newTestServer(t *testing.T, handler http.Handler) *httptest.Server {
	t.Helper()
	ln, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := httptest.NewUnstartedServer(handler)
	server.Listener = ln
	server.Start()
	t.Cleanup(server.Close)
	return server
}

func TestFindReleaseAssets(t *testing.T) {
	release := &ReleaseInfo{
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{Name: "promptext.tar.gz", BrowserDownloadURL: "http://example.com/asset"},
			{Name: "checksums.txt", BrowserDownloadURL: "http://example.com/checksums"},
		},
	}

	downloadURL, checksumURL, err := findReleaseAssets(release, "promptext.tar.gz")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if downloadURL != "http://example.com/asset" || checksumURL != "http://example.com/checksums" {
		t.Fatalf("unexpected urls: %s %s", downloadURL, checksumURL)
	}

	if _, _, err := findReleaseAssets(release, "missing"); err == nil {
		t.Fatalf("expected error for missing asset")
	}
}

func TestDownloadAndVerifyBinaryTarGz(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("tar.gz extraction not used on windows")
	}
	binaryData := []byte("test binary")
	archiveBytes := writeTarGz(t, "promptext", binaryData)
	checksum := sha256Hex(archiveBytes)

	originalDownload := downloadFileFn
	defer func() { downloadFileFn = originalDownload }()

	downloadFileFn = func(destPath, url string) error {
		switch url {
		case "asset":
			return os.WriteFile(destPath, archiveBytes, 0644)
		case "checksum":
			return os.WriteFile(destPath, []byte(fmt.Sprintf("%s  promptext.tar.gz\n", checksum)), 0644)
		default:
			t.Fatalf("unexpected download url: %s", url)
		}
		return nil
	}

	path, err := downloadAndVerifyBinary("asset", "checksum", "promptext.tar.gz", false)
	if err != nil {
		t.Fatalf("downloadAndVerifyBinary failed: %v", err)
	}
	t.Cleanup(func() {
		os.Remove(path)
	})
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read binary: %v", err)
	}
	if !bytes.Equal(got, binaryData) {
		t.Fatalf("unexpected binary contents: %q", got)
	}
}

func TestDownloadAndVerifyBinaryZip(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("zip extraction primarily for windows")
	}
	binaryData := []byte("zip binary")
	archiveBytes := writeZip(t, "promptext.exe", binaryData)
	checksum := sha256Hex(archiveBytes)

	originalDownload := downloadFileFn
	defer func() { downloadFileFn = originalDownload }()

	downloadFileFn = func(destPath, url string) error {
		switch url {
		case "asset":
			return os.WriteFile(destPath, archiveBytes, 0644)
		case "checksum":
			return os.WriteFile(destPath, []byte(fmt.Sprintf("%s  promptext.zip\n", checksum)), 0644)
		default:
			t.Fatalf("unexpected download url: %s", url)
		}
		return nil
	}

	path, err := downloadAndVerifyBinary("asset", "checksum", "promptext.zip", false)
	if err != nil {
		t.Fatalf("downloadAndVerifyBinary failed: %v", err)
	}
	t.Cleanup(func() {
		os.Remove(path)
	})
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read binary: %v", err)
	}
	if !bytes.Equal(got, binaryData) {
		t.Fatalf("unexpected binary contents: %q", got)
	}
}

func TestDownloadAndVerifyBinaryChecksumMismatch(t *testing.T) {
	binaryData := []byte("test binary")
	archiveBytes := writeTarGz(t, "promptext", binaryData)

	originalDownload := downloadFileFn
	defer func() { downloadFileFn = originalDownload }()

	downloadFileFn = func(destPath, url string) error {
		switch url {
		case "asset":
			return os.WriteFile(destPath, archiveBytes, 0644)
		case "checksum":
			return os.WriteFile(destPath, []byte(fmt.Sprintf("%s  promptext.tar.gz\n", strings.Repeat("0", 64))), 0644)
		default:
			t.Fatalf("unexpected download url: %s", url)
		}
		return nil
	}

	if _, err := downloadAndVerifyBinary("asset", "checksum", "promptext.tar.gz", false); err == nil {
		t.Fatalf("expected checksum error")
	}
}

func TestExtractBinaryUnsupported(t *testing.T) {
	if _, err := extractBinary("file.bin", t.TempDir()); err == nil {
		t.Fatalf("expected unsupported format error")
	}
}

func TestCheckForUpdateWithTimeoutSuccess(t *testing.T) {
	original := checkForUpdateFn
	checkForUpdateFn = func(string) (bool, string, error) {
		return true, "v1.2.3", nil
	}
	t.Cleanup(func() {
		checkForUpdateFn = original
	})

	available, version, err := checkForUpdateWithTimeout("v1.0.0", time.Second)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !available || version != "v1.2.3" {
		t.Fatalf("unexpected result: %v %s", available, version)
	}
}

func TestCheckForUpdateWithTimeoutTimeout(t *testing.T) {
	original := checkForUpdateFn
	checkForUpdateFn = func(string) (bool, string, error) {
		time.Sleep(20 * time.Millisecond)
		return false, "", nil
	}
	t.Cleanup(func() {
		checkForUpdateFn = original
	})

	if _, _, err := checkForUpdateWithTimeout("v1.0.0", time.Millisecond); err == nil {
		t.Fatalf("expected timeout error")
	}
}

func TestCheckForUpdateWithTimeoutError(t *testing.T) {
	original := checkForUpdateFn
	checkForUpdateFn = func(string) (bool, string, error) {
		return false, "", errors.New("boom")
	}
	t.Cleanup(func() {
		checkForUpdateFn = original
	})

	if _, _, err := checkForUpdateWithTimeout("v1.0.0", time.Second); err == nil {
		t.Fatalf("expected propagated error")
	}
}

func TestCheckAndNotifyUpdateUsesCache(t *testing.T) {
	originalLoad := loadUpdateCacheFn
	originalSave := saveUpdateCacheFn
	originalCheck := checkForUpdateFn
	defer func() {
		loadUpdateCacheFn = originalLoad
		saveUpdateCacheFn = originalSave
		checkForUpdateFn = originalCheck
	}()

	loadUpdateCacheFn = func() (*UpdateCheckCache, error) {
		return &UpdateCheckCache{
			LastCheck:       time.Now(),
			LatestVersion:   "v2.0.0",
			UpdateAvailable: true,
		}, nil
	}
	checkCalled := false
	checkForUpdateFn = func(string) (bool, string, error) {
		checkCalled = true
		return false, "", nil
	}
	saveUpdateCacheFn = func(UpdateCheckCache) error {
		return nil
	}

	output := captureStderr(t, func() {
		CheckAndNotifyUpdate("v1.0.0")
	})

	if checkCalled {
		t.Fatalf("expected cached result to avoid network call")
	}
	if !strings.Contains(output, "Update available: v2.0.0") {
		t.Fatalf("expected notification, got %q", output)
	}
}

func TestCheckAndNotifyUpdatePerformsCheckWhenStale(t *testing.T) {
	originalLoad := loadUpdateCacheFn
	originalSave := saveUpdateCacheFn
	originalCheck := checkForUpdateFn
	defer func() {
		loadUpdateCacheFn = originalLoad
		saveUpdateCacheFn = originalSave
		checkForUpdateFn = originalCheck
	}()

	loadUpdateCacheFn = func() (*UpdateCheckCache, error) {
		return &UpdateCheckCache{
			LastCheck:       time.Now().Add(-2 * checkInterval),
			LatestVersion:   "v1.0.0",
			UpdateAvailable: false,
		}, nil
	}
	var saved UpdateCheckCache
	saveCalled := false
	saveUpdateCacheFn = func(cache UpdateCheckCache) error {
		saveCalled = true
		saved = cache
		return nil
	}
	checkForUpdateFn = func(string) (bool, string, error) {
		return true, "v3.0.0", nil
	}

	output := captureStderr(t, func() {
		CheckAndNotifyUpdate("v1.0.0")
	})

	if !saveCalled {
		t.Fatalf("expected cache to be saved")
	}
	if saved.LatestVersion != "v3.0.0" || !saved.UpdateAvailable {
		t.Fatalf("unexpected cache saved: %+v", saved)
	}
	if !strings.Contains(output, "Update available: v3.0.0") {
		t.Fatalf("expected notification, got %q", output)
	}
}

func TestReplaceBinaryRollback(t *testing.T) {
	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "promptext")
	backupPath := execPath + ".old"
	if err := os.WriteFile(execPath, []byte("old"), 0755); err != nil {
		t.Fatalf("write exec: %v", err)
	}
	newPath := filepath.Join(tmpDir, "new")
	if err := os.WriteFile(newPath, []byte("new"), 0755); err != nil {
		t.Fatalf("write new: %v", err)
	}

	originalCopy := copyFileFn
	copyFileFn = func(string, string) error {
		return errors.New("copy failed")
	}
	t.Cleanup(func() {
		copyFileFn = originalCopy
	})

	if err := replaceBinary(execPath, newPath, false); err == nil {
		t.Fatalf("expected replaceBinary to fail")
	}
	data, err := os.ReadFile(execPath)
	if err != nil {
		t.Fatalf("read exec: %v", err)
	}
	if string(data) != "old" {
		t.Fatalf("expected original binary restored, got %q", data)
	}
	if _, err := os.Stat(backupPath); err == nil {
		t.Fatalf("expected backup to be cleaned up")
	}
}

func TestReplaceBinarySuccess(t *testing.T) {
	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "promptext")
	if err := os.WriteFile(execPath, []byte("old"), 0755); err != nil {
		t.Fatalf("write exec: %v", err)
	}
	newPath := filepath.Join(tmpDir, "new")
	if err := os.WriteFile(newPath, []byte("new"), 0755); err != nil {
		t.Fatalf("write new: %v", err)
	}

	if err := replaceBinary(execPath, newPath, false); err != nil {
		t.Fatalf("replaceBinary failed: %v", err)
	}
	data, err := os.ReadFile(execPath)
	if err != nil {
		t.Fatalf("read exec: %v", err)
	}
	if string(data) != "new" {
		t.Fatalf("expected new binary installed, got %q", data)
	}
	if _, err := os.Stat(newPath); !os.IsNotExist(err) {
		t.Fatalf("expected temporary binary to be removed")
	}
}

func TestDownloadFileHTTPError(t *testing.T) {
	tmp := filepath.Join(t.TempDir(), "file")
	if err := downloadFile(tmp, ":bad-url"); err == nil {
		t.Fatalf("expected download error")
	}
}

func TestGetExecutablePathReturnsPath(t *testing.T) {
	path, err := getExecutablePath()
	if err != nil {
		t.Fatalf("getExecutablePath error: %v", err)
	}
	if path == "" {
		t.Fatalf("expected non-empty executable path")
	}
}

func TestFetchLatestReleaseSuccess(t *testing.T) {
	originalURL := githubAPIURL
	originalClient := httpClient
	defer func() {
		githubAPIURL = originalURL
		httpClient = originalClient
	}()

	githubAPIURL = "https://example.com/latest"
	httpClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			if req.URL.String() != githubAPIURL {
				t.Fatalf("unexpected URL: %s", req.URL)
			}
			if ua := req.Header.Get("User-Agent"); !strings.Contains(ua, "promptext") {
				t.Fatalf("expected user agent to contain promptext, got %s", ua)
			}
			body := `{"tag_name":"v2.0.0","name":"Release"}`
			return &http.Response{
				StatusCode: http.StatusOK,
				Body:       io.NopCloser(strings.NewReader(body)),
				Header:     make(http.Header),
			}, nil
		}),
	}

	release, err := fetchLatestRelease()
	if err != nil {
		t.Fatalf("fetchLatestRelease error: %v", err)
	}
	if release.TagName != "v2.0.0" || release.Name != "Release" {
		t.Fatalf("unexpected release: %+v", release)
	}
}

func TestFetchLatestReleaseHTTPError(t *testing.T) {
	originalURL := githubAPIURL
	originalClient := httpClient
	defer func() {
		githubAPIURL = originalURL
		httpClient = originalClient
	}()

	githubAPIURL = "https://example.com/latest"
	httpClient = &http.Client{
		Transport: roundTripperFunc(func(req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(strings.NewReader("")),
				Header:     make(http.Header),
			}, nil
		}),
	}

	if _, err := fetchLatestRelease(); err == nil {
		t.Fatalf("expected error for non-200 response")
	}
}

func TestVerifyChecksumMissingEntry(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "archive")
	if err := os.WriteFile(archive, []byte("data"), 0644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	checksumPath := filepath.Join(dir, "checksums.txt")
	if err := os.WriteFile(checksumPath, []byte("deadbeef  other.tar.gz\n"), 0644); err != nil {
		t.Fatalf("write checksum: %v", err)
	}

	if err := verifyChecksum(archive, checksumPath, "promptext.tar.gz"); err == nil {
		t.Fatalf("expected missing checksum error")
	}
}

func TestVerifyChecksumMismatch(t *testing.T) {
	dir := t.TempDir()
	archive := filepath.Join(dir, "archive")
	if err := os.WriteFile(archive, []byte("data"), 0644); err != nil {
		t.Fatalf("write archive: %v", err)
	}
	checksumPath := filepath.Join(dir, "checksums.txt")
	if err := os.WriteFile(checksumPath, []byte("deadbeef  archive\n"), 0644); err != nil {
		t.Fatalf("write checksum: %v", err)
	}

	if err := verifyChecksum(archive, checksumPath, "archive"); err == nil {
		t.Fatalf("expected mismatch error")
	}
}

func TestLoadAndSaveUpdateCache(t *testing.T) {
	dir := t.TempDir()
	original := getCacheDirFn
	getCacheDirFn = func() (string, error) {
		return dir, nil
	}
	t.Cleanup(func() {
		getCacheDirFn = original
	})

	cache := UpdateCheckCache{LatestVersion: "v1.0.0", UpdateAvailable: true, LastCheck: time.Unix(10, 0)}
	if err := saveUpdateCache(cache); err != nil {
		t.Fatalf("save cache: %v", err)
	}

	loaded, err := loadUpdateCache()
	if err != nil {
		t.Fatalf("load cache: %v", err)
	}
	if loaded.LatestVersion != cache.LatestVersion || loaded.UpdateAvailable != cache.UpdateAvailable {
		t.Fatalf("unexpected cache contents: %+v", loaded)
	}
}

func TestCheckForUpdateDevelopmentBuild(t *testing.T) {
	if _, _, err := CheckForUpdate("dev"); err == nil {
		t.Fatalf("expected error for development version")
	}
	if _, _, err := CheckForUpdate("unknown"); err == nil {
		t.Fatalf("expected error for unknown version")
	}
}

func TestCheckForUpdateNewerVersion(t *testing.T) {
	originalFetch := fetchLatestReleaseFn
	defer func() { fetchLatestReleaseFn = originalFetch }()

	fetchLatestReleaseFn = func() (*ReleaseInfo, error) {
		return &ReleaseInfo{TagName: "v1.2.3"}, nil
	}

	available, latest, err := CheckForUpdate("v1.0.0")
	if err != nil {
		t.Fatalf("CheckForUpdate error: %v", err)
	}
	if !available {
		t.Fatalf("expected update to be available")
	}
	if latest != "v1.2.3" {
		t.Fatalf("unexpected latest version: %s", latest)
	}
}

func TestCheckForUpdateSameVersion(t *testing.T) {
	originalFetch := fetchLatestReleaseFn
	defer func() { fetchLatestReleaseFn = originalFetch }()

	fetchLatestReleaseFn = func() (*ReleaseInfo, error) {
		return &ReleaseInfo{TagName: "v1.0.0"}, nil
	}

	available, latest, err := CheckForUpdate("v1.0.0")
	if err != nil {
		t.Fatalf("CheckForUpdate error: %v", err)
	}
	if available {
		t.Fatalf("expected no update for same version")
	}
	if latest != "v1.0.0" {
		t.Fatalf("expected latest to match release tag, got %s", latest)
	}
}

func TestUpdateNoUpdateAvailable(t *testing.T) {
	originalCheck := checkForUpdateFn
	defer func() { checkForUpdateFn = originalCheck }()

	checkForUpdateFn = func(string) (bool, string, error) {
		return false, "v1.0.0", nil
	}

	if err := Update("v1.0.0", true); err != nil {
		t.Fatalf("expected no error when update unavailable, got %v", err)
	}
}

func TestUpdateSuccess(t *testing.T) {
	tmpDir := t.TempDir()
	execPath := filepath.Join(tmpDir, "promptext")
	if err := os.WriteFile(execPath, []byte("old"), 0755); err != nil {
		t.Fatalf("write exec: %v", err)
	}

	originalCheck := checkForUpdateFn
	originalFetch := fetchLatestReleaseFn
	originalGetExec := getExecutablePathFn
	originalGetAsset := getPlatformAssetNameFn
	originalFindAssets := findReleaseAssetsFn
	originalDownload := downloadFileFn
	originalVerify := verifyChecksumFn
	originalExtract := extractBinaryFn
	originalCopy := copyFileFn
	originalReplace := replaceBinaryFn
	defer func() {
		checkForUpdateFn = originalCheck
		fetchLatestReleaseFn = originalFetch
		getExecutablePathFn = originalGetExec
		getPlatformAssetNameFn = originalGetAsset
		findReleaseAssetsFn = originalFindAssets
		downloadFileFn = originalDownload
		verifyChecksumFn = originalVerify
		extractBinaryFn = originalExtract
		copyFileFn = originalCopy
		replaceBinaryFn = originalReplace
	}()

	checkForUpdateFn = func(string) (bool, string, error) {
		return true, "v2.0.0", nil
	}
	fetchLatestReleaseFn = func() (*ReleaseInfo, error) {
		return &ReleaseInfo{Assets: nil}, nil
	}
	getExecutablePathFn = func() (string, error) {
		return execPath, nil
	}
	getPlatformAssetNameFn = func() (string, error) {
		return "promptext_Linux_x86_64.tar.gz", nil
	}
	findReleaseAssetsFn = func(*ReleaseInfo, string) (string, string, error) {
		return "mock-asset", "mock-checksum", nil
	}

	var archivePath, checksumPath string
	downloadFileFn = func(destPath, url string) error {
		switch url {
		case "mock-asset":
			archivePath = destPath
			return os.WriteFile(destPath, []byte("archive"), 0644)
		case "mock-checksum":
			checksumPath = destPath
			return os.WriteFile(destPath, []byte("checksum"), 0644)
		default:
			t.Fatalf("unexpected download url: %s", url)
		}
		return nil
	}
	verifyChecksumFn = func(filePath, checksumFile, assetName string) error {
		if filePath != archivePath || checksumFile != checksumPath || assetName == "" {
			t.Fatalf("unexpected checksum args: %s %s %s", filePath, checksumFile, assetName)
		}
		return nil
	}
	newBinary := filepath.Join(tmpDir, "new-binary")
	extractBinaryFn = func(path, destDir string) (string, error) {
		if path != archivePath {
			t.Fatalf("unexpected extract path: %s", path)
		}
		if err := os.WriteFile(newBinary, []byte("new"), 0755); err != nil {
			return "", err
		}
		return newBinary, nil
	}
	copyFileFn = func(src, dst string) error {
		if src != newBinary {
			t.Fatalf("unexpected copy src: %s", src)
		}
		return os.WriteFile(dst, []byte("new"), 0755)
	}
	var replaced bool
	replaceBinaryFn = func(oldPath, newPath string, verbose bool) error {
		replaced = true
		if oldPath != execPath || newPath == "" {
			t.Fatalf("unexpected replace args: %s %s", oldPath, newPath)
		}
		return nil
	}

	if err := Update("v1.0.0", true); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if !replaced {
		t.Fatalf("expected replaceBinary to be called")
	}
}

func TestGetPlatformAssetName(t *testing.T) {
	name, err := getPlatformAssetName()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if name == "" {
		t.Fatalf("expected non-empty asset name")
	}

	if runtime.GOOS == "windows" && !strings.HasSuffix(name, ".zip") {
		t.Fatalf("expected windows asset to be zip, got %s", name)
	}
	if runtime.GOOS != "windows" && !strings.HasSuffix(name, ".tar.gz") {
		t.Fatalf("expected non-windows asset to be tar.gz, got %s", name)
	}
}

func TestExtractZipBinary(t *testing.T) {
	zipBytes := writeZip(t, "promptext.exe", []byte("binary"))
	archive := filepath.Join(t.TempDir(), "archive.zip")
	if err := os.WriteFile(archive, zipBytes, 0644); err != nil {
		t.Fatalf("write archive: %v", err)
	}

	path, err := extractZip(archive, t.TempDir())
	if err != nil {
		t.Fatalf("extractZip error: %v", err)
	}
	if filepath.Base(path) != "promptext.exe" {
		t.Fatalf("unexpected binary path: %s", path)
	}
}

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		v1, v2   string
		expected bool
	}{
		{"major", "2.0.0", "1.9.9", true},
		{"minor", "1.2.0", "1.1.5", true},
		{"patch", "1.0.1", "1.0.0", true},
		{"equal", "1.0.0", "1.0.0", false},
		{"older", "1.0.0", "1.1.0", false},
	}

	for _, tt := range tests {
		res, err := isNewerVersion(tt.v1, tt.v2)
		if err != nil {
			t.Fatalf("%s: unexpected error: %v", tt.name, err)
		}
		if res != tt.expected {
			t.Fatalf("%s: expected %v, got %v", tt.name, tt.expected, res)
		}
	}
}

func TestParseVersionInvalid(t *testing.T) {
	if _, err := parseVersion("invalid"); err == nil {
		t.Fatalf("expected error for invalid version")
	}
}

func TestGetCacheDirHonorsXDG(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("XDG cache home not used on Windows")
	}
	oldXDG := os.Getenv("XDG_CACHE_HOME")
	defer os.Setenv("XDG_CACHE_HOME", oldXDG)

	tmp := t.TempDir()
	if err := os.Setenv("XDG_CACHE_HOME", tmp); err != nil {
		t.Fatalf("set env: %v", err)
	}

	cacheDir, err := getCacheDir()
	if err != nil {
		t.Fatalf("getCacheDir error: %v", err)
	}
	if !strings.HasPrefix(cacheDir, tmp) {
		t.Fatalf("expected cache dir inside temp, got %s", cacheDir)
	}
	if fi, err := os.Stat(cacheDir); err != nil || !fi.IsDir() {
		t.Fatalf("expected cache dir to exist: %v", err)
	}
}

func TestCheckForUpdateVersionCompareError(t *testing.T) {
	originalFetch := fetchLatestReleaseFn
	defer func() { fetchLatestReleaseFn = originalFetch }()

	fetchLatestReleaseFn = func() (*ReleaseInfo, error) {
		return &ReleaseInfo{TagName: "bad.version"}, nil
	}

	if _, _, err := CheckForUpdate("v1.0.0"); err == nil {
		t.Fatalf("expected error when release tag is invalid")
	}
}

func TestUpdateUnsupportedPlatform(t *testing.T) {
	originalCheck := checkForUpdateFn
	originalFetch := fetchLatestReleaseFn
	originalGetAsset := getPlatformAssetNameFn
	originalGetExec := getExecutablePathFn
	defer func() {
		checkForUpdateFn = originalCheck
		fetchLatestReleaseFn = originalFetch
		getPlatformAssetNameFn = originalGetAsset
		getExecutablePathFn = originalGetExec
	}()

	checkForUpdateFn = func(string) (bool, string, error) {
		return true, "v2.0.0", nil
	}
	fetchLatestReleaseFn = func() (*ReleaseInfo, error) {
		return &ReleaseInfo{}, nil
	}
	getExecutablePathFn = func() (string, error) {
		return filepath.Join(t.TempDir(), "promptext"), nil
	}
	getPlatformAssetNameFn = func() (string, error) {
		return "", fmt.Errorf("unsupported platform")
	}

	if err := Update("v1.0.0", false); err == nil {
		t.Fatalf("expected Update to fail when platform unsupported")
	}
}

func TestDownloadAndVerifyBinaryDownloadError(t *testing.T) {
	originalDownload := downloadFileFn
	defer func() { downloadFileFn = originalDownload }()

	downloadFileFn = func(string, string) error {
		return fmt.Errorf("download failed")
	}

	if _, err := downloadAndVerifyBinary("asset", "", "promptext.tar.gz", false); err == nil {
		t.Fatalf("expected error when download fails")
	}
}

func TestDownloadAndVerifyBinaryChecksumDownloadError(t *testing.T) {
	originalDownload := downloadFileFn
	defer func() { downloadFileFn = originalDownload }()

	called := 0
	downloadFileFn = func(destPath, url string) error {
		called++
		if url == "asset" {
			return os.WriteFile(destPath, []byte("archive"), 0644)
		}
		return fmt.Errorf("checksum download failed")
	}

	if _, err := downloadAndVerifyBinary("asset", "checksum", "promptext.tar.gz", false); err == nil {
		t.Fatalf("expected error when checksum download fails")
	}
	if called < 2 {
		t.Fatalf("expected download to be attempted twice")
	}
}

func TestDownloadAndVerifyBinaryExtractError(t *testing.T) {
	originalDownload := downloadFileFn
	originalExtract := extractBinaryFn
	defer func() {
		downloadFileFn = originalDownload
		extractBinaryFn = originalExtract
	}()

	archiveBytes := writeTarGz(t, "promptext", []byte("binary"))
	downloadFileFn = func(destPath, url string) error {
		if url == "asset" {
			return os.WriteFile(destPath, archiveBytes, 0644)
		}
		return nil
	}
	extractBinaryFn = func(string, string) (string, error) {
		return "", fmt.Errorf("extract failed")
	}

	if _, err := downloadAndVerifyBinary("asset", "", "promptext.tar.gz", true); err == nil {
		t.Fatalf("expected error when extraction fails")
	}
}

func TestDownloadAndVerifyBinaryCopyError(t *testing.T) {
	originalDownload := downloadFileFn
	originalCopy := copyFileFn
	defer func() {
		downloadFileFn = originalDownload
		copyFileFn = originalCopy
	}()

	archiveBytes := writeTarGz(t, "promptext", []byte("binary"))
	downloadFileFn = func(destPath, url string) error {
		if url == "asset" {
			return os.WriteFile(destPath, archiveBytes, 0644)
		}
		return nil
	}
	copyFileFn = func(string, string) error {
		return fmt.Errorf("copy failed")
	}

	if _, err := downloadAndVerifyBinary("asset", "", "promptext.tar.gz", true); err == nil {
		t.Fatalf("expected error when copying fails")
	}
}
