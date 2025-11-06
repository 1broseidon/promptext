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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset":
			w.Write(archiveBytes)
		case "/checksums.txt":
			fmt.Fprintf(w, "%s  promptext.tar.gz\n", checksum)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	path, err := downloadAndVerifyBinary(server.URL+"/asset", server.URL+"/checksums.txt", "promptext.tar.gz", false)
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset":
			w.Write(archiveBytes)
		case "/checksums.txt":
			fmt.Fprintf(w, "%s  promptext.zip\n", checksum)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	path, err := downloadAndVerifyBinary(server.URL+"/asset", server.URL+"/checksums.txt", "promptext.zip", false)
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/asset":
			w.Write(archiveBytes)
		case "/checksums.txt":
			fmt.Fprintf(w, "%s  promptext.tar.gz\n", strings.Repeat("0", 64))
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()

	if _, err := downloadAndVerifyBinary(server.URL+"/asset", server.URL+"/checksums.txt", "promptext.tar.gz", false); err == nil {
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
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	tmp := filepath.Join(t.TempDir(), "file")
	if err := downloadFile(tmp, server.URL); err == nil {
		t.Fatalf("expected download error")
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
