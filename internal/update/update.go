package update

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

const (
	githubAPIURL     = "https://api.github.com/repos/1broseidon/promptext/releases/latest"
	githubReleaseURL = "https://github.com/1broseidon/promptext/releases/download"
	downloadTimeout  = 5 * time.Minute
	checkInterval    = 24 * time.Hour // Check for updates once per day
)

var (
	fetchLatestReleaseFn   = fetchLatestRelease
	downloadFileFn         = downloadFile
	verifyChecksumFn       = verifyChecksum
	extractBinaryFn        = extractBinary
	getExecutablePathFn    = getExecutablePath
	getPlatformAssetNameFn = getPlatformAssetName
	findReleaseAssetsFn    = findReleaseAssets
	replaceBinaryFn        = replaceBinary
	checkForUpdateFn       = CheckForUpdate
	loadUpdateCacheFn      = loadUpdateCache
	saveUpdateCacheFn      = saveUpdateCache
	copyFileFn             = copyFile
	getCacheDirFn          = getCacheDir
)

// ReleaseInfo represents GitHub release metadata
type ReleaseInfo struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpdateCheckCache stores the last update check information
type UpdateCheckCache struct {
	LastCheck       time.Time `json:"last_check"`
	LatestVersion   string    `json:"latest_version"`
	UpdateAvailable bool      `json:"update_available"`
}

// CheckForUpdate queries GitHub API for the latest release and compares versions
func CheckForUpdate(currentVersion string) (available bool, latestVersion string, err error) {
	// Clean current version (remove 'v' prefix if present)
	current := strings.TrimPrefix(currentVersion, "v")
	if current == "dev" || current == "unknown" {
		return false, "", fmt.Errorf("cannot check updates for development build (version: %s)", currentVersion)
	}

	// Fetch latest release info
	release, err := fetchLatestReleaseFn()
	if err != nil {
		return false, "", fmt.Errorf("failed to fetch latest release: %w", err)
	}

	latest := strings.TrimPrefix(release.TagName, "v")

	// Compare versions
	isNewer, err := isNewerVersion(latest, current)
	if err != nil {
		return false, "", fmt.Errorf("failed to compare versions: %w", err)
	}

	return isNewer, release.TagName, nil
}

// Update downloads and installs the latest version
func Update(currentVersion string, verbose bool) error {
	// Check if update is available
	available, latestVersion, err := checkForUpdateFn(currentVersion)
	if err != nil {
		return err
	}

	if !available {
		if verbose {
			fmt.Printf("Already running the latest version (%s)\n", currentVersion)
		}
		return nil
	}

	if verbose {
		fmt.Printf("Updating from %s to %s...\n", currentVersion, latestVersion)
	}

	// Get current executable path
	execPath, err := getExecutablePathFn()
	if err != nil {
		return err
	}

	// Fetch release info and find assets
	release, err := fetchLatestReleaseFn()
	if err != nil {
		return fmt.Errorf("failed to fetch release info: %w", err)
	}

	assetName, err := getPlatformAssetNameFn()
	if err != nil {
		return fmt.Errorf("unsupported platform: %w", err)
	}

	downloadURL, checksumURL, err := findReleaseAssetsFn(release, assetName)
	if err != nil {
		return err
	}

	// Download, verify, and extract binary
	binaryPath, err := downloadAndVerifyBinary(downloadURL, checksumURL, assetName, verbose)
	if err != nil {
		return err
	}

	// Replace current binary with new one
	if err := replaceBinaryFn(execPath, binaryPath, verbose); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("✓ Successfully updated to %s\n", latestVersion)
		fmt.Println("\nRestart your terminal or run the command again to use the new version.")
	}

	return nil
}

// getExecutablePath returns the resolved path to the current executable
func getExecutablePath() (string, error) {
	execPath, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("failed to get executable path: %w", err)
	}
	execPath, err = filepath.EvalSymlinks(execPath)
	if err != nil {
		return "", fmt.Errorf("failed to resolve executable path: %w", err)
	}
	return execPath, nil
}

// findReleaseAssets finds download and checksum URLs from release assets
func findReleaseAssets(release *ReleaseInfo, assetName string) (downloadURL, checksumURL string, err error) {
	for _, asset := range release.Assets {
		if asset.Name == assetName {
			downloadURL = asset.BrowserDownloadURL
		}
		if asset.Name == "checksums.txt" {
			checksumURL = asset.BrowserDownloadURL
		}
	}

	if downloadURL == "" {
		return "", "", fmt.Errorf("no release asset found for %s", assetName)
	}

	return downloadURL, checksumURL, nil
}

// downloadAndVerifyBinary downloads, verifies, and extracts the binary
func downloadAndVerifyBinary(downloadURL, checksumURL, assetName string, verbose bool) (string, error) {
	// Create temporary directory for download
	tempDir, err := os.MkdirTemp("", "promptext-update-")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Download archive
	archivePath := filepath.Join(tempDir, assetName)
	if verbose {
		fmt.Printf("Downloading %s...\n", assetName)
	}
	if err := downloadFileFn(archivePath, downloadURL); err != nil {
		return "", fmt.Errorf("failed to download update: %w", err)
	}

	// Download and verify checksum if available
	if checksumURL != "" {
		if verbose {
			fmt.Println("Verifying checksum...")
		}
		checksumPath := filepath.Join(tempDir, "checksums.txt")
		if err := downloadFileFn(checksumPath, checksumURL); err != nil {
			return "", fmt.Errorf("failed to download checksums: %w", err)
		}
		if err := verifyChecksumFn(archivePath, checksumPath, assetName); err != nil {
			return "", fmt.Errorf("checksum verification failed: %w", err)
		}
		if verbose {
			fmt.Println("✓ Checksum verified")
		}
	}

	// Extract archive
	if verbose {
		fmt.Println("Extracting archive...")
	}
	binaryPath, err := extractBinaryFn(archivePath, tempDir)
	if err != nil {
		return "", fmt.Errorf("failed to extract binary: %w", err)
	}

	// Make binary executable (Unix-like systems)
	if runtime.GOOS != "windows" {
		if err := os.Chmod(binaryPath, 0755); err != nil {
			return "", fmt.Errorf("failed to make binary executable: %w", err)
		}
	}

	// Copy to system temp directory (outside our update temp dir)
	// This survives the defer cleanup of tempDir
	sysTempDir := os.TempDir()
	permanentPath := filepath.Join(sysTempDir, "promptext-new-binary")
	if err := copyFileFn(binaryPath, permanentPath); err != nil {
		return "", fmt.Errorf("failed to copy binary: %w", err)
	}

	return permanentPath, nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// replaceBinary atomically replaces the current binary with the new one
func replaceBinary(execPath, binaryPath string, verbose bool) error {
	backupPath := execPath + ".old"
	if verbose {
		fmt.Println("Installing new version...")
	}

	// Remove old backup if it exists
	os.Remove(backupPath)

	// Rename current binary to backup
	if err := os.Rename(execPath, backupPath); err != nil {
		return fmt.Errorf("failed to backup current binary: %w", err)
	}

	// Copy new binary to executable path
	// We use copy instead of rename because binaryPath might be on a different filesystem
	if err := copyFileFn(binaryPath, execPath); err != nil {
		// Rollback: restore backup
		os.Rename(backupPath, execPath)
		return fmt.Errorf("failed to install new binary: %w", err)
	}

	// Clean up temporary binary file
	os.Remove(binaryPath)

	// Remove backup on success
	os.Remove(backupPath)
	return nil
}

// fetchLatestRelease queries GitHub API for latest release information
func fetchLatestRelease() (*ReleaseInfo, error) {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	req, err := http.NewRequest("GET", githubAPIURL, nil)
	if err != nil {
		return nil, err
	}

	// Set User-Agent (GitHub API requires it)
	req.Header.Set("User-Agent", "promptext-updater")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release ReleaseInfo
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// getPlatformAssetName returns the asset name for the current platform
func getPlatformAssetName() (string, error) {
	var osName, archName string

	// Map GOOS to release asset OS name
	switch runtime.GOOS {
	case "darwin":
		osName = "Darwin"
	case "linux":
		osName = "Linux"
	case "windows":
		osName = "Windows"
	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	// Map GOARCH to release asset architecture name
	switch runtime.GOARCH {
	case "amd64":
		archName = "x86_64"
	case "arm64":
		archName = "arm64"
	default:
		return "", fmt.Errorf("unsupported architecture: %s", runtime.GOARCH)
	}

	// Determine file extension
	ext := ".tar.gz"
	if runtime.GOOS == "windows" {
		ext = ".zip"
	}

	return fmt.Sprintf("promptext_%s_%s%s", osName, archName, ext), nil
}

// downloadFile downloads a file from URL to destination path
func downloadFile(destPath, url string) error {
	client := &http.Client{
		Timeout: downloadTimeout,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies the SHA256 checksum of the downloaded file
func verifyChecksum(filePath, checksumPath, assetName string) error {
	// Read checksums file
	checksumFile, err := os.Open(checksumPath)
	if err != nil {
		return err
	}
	defer checksumFile.Close()

	// Find checksum for our asset
	var expectedChecksum string
	scanner := bufio.NewScanner(checksumFile)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[1] == assetName {
			expectedChecksum = parts[0]
			break
		}
	}

	if expectedChecksum == "" {
		return fmt.Errorf("checksum not found for %s", assetName)
	}

	// Calculate actual checksum
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))

	// Compare checksums
	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// extractBinary extracts the binary from the downloaded archive
func extractBinary(archivePath, destDir string) (string, error) {
	if strings.HasSuffix(archivePath, ".tar.gz") {
		return extractTarGz(archivePath, destDir)
	} else if strings.HasSuffix(archivePath, ".zip") {
		return extractZip(archivePath, destDir)
	}
	return "", fmt.Errorf("unsupported archive format: %s", archivePath)
}

// extractTarGz extracts a .tar.gz archive and returns the path to the binary
func extractTarGz(archivePath, destDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	var binaryPath string
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the binary (named "promptext" or "prx")
		if header.Typeflag == tar.TypeReg {
			baseName := filepath.Base(header.Name)
			if baseName == "promptext" || baseName == "prx" {
				targetPath := filepath.Join(destDir, baseName)
				outFile, err := os.Create(targetPath)
				if err != nil {
					return "", err
				}
				if _, err := io.Copy(outFile, tr); err != nil {
					outFile.Close()
					return "", err
				}
				outFile.Close()
				binaryPath = targetPath
				break
			}
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

// extractZip extracts a .zip archive and returns the path to the binary
func extractZip(archivePath, destDir string) (string, error) {
	r, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", err
	}
	defer r.Close()

	var binaryPath string
	for _, f := range r.File {
		baseName := filepath.Base(f.Name)
		if baseName == "promptext.exe" || baseName == "prx.exe" {
			rc, err := f.Open()
			if err != nil {
				return "", err
			}

			targetPath := filepath.Join(destDir, baseName)
			outFile, err := os.Create(targetPath)
			if err != nil {
				rc.Close()
				return "", err
			}

			_, err = io.Copy(outFile, rc)
			rc.Close()
			outFile.Close()

			if err != nil {
				return "", err
			}

			binaryPath = targetPath
			break
		}
	}

	if binaryPath == "" {
		return "", fmt.Errorf("binary not found in archive")
	}

	return binaryPath, nil
}

// isNewerVersion compares two semantic versions (without 'v' prefix)
// Returns true if v1 is newer than v2
func isNewerVersion(v1, v2 string) (bool, error) {
	v1Parts, err := parseVersion(v1)
	if err != nil {
		return false, fmt.Errorf("invalid version format: %s", v1)
	}

	v2Parts, err := parseVersion(v2)
	if err != nil {
		return false, fmt.Errorf("invalid version format: %s", v2)
	}

	// Compare major, minor, patch
	for i := 0; i < 3; i++ {
		if v1Parts[i] > v2Parts[i] {
			return true, nil
		}
		if v1Parts[i] < v2Parts[i] {
			return false, nil
		}
	}

	return false, nil // versions are equal
}

// parseVersion parses a semantic version string (e.g., "0.4.3") into [major, minor, patch]
func parseVersion(version string) ([3]int, error) {
	var parts [3]int
	_, err := fmt.Sscanf(version, "%d.%d.%d", &parts[0], &parts[1], &parts[2])
	if err != nil {
		return parts, err
	}
	return parts, nil
}

// CheckAndNotifyUpdate performs a non-blocking update check and notifies user if update available
// This is called automatically during normal CLI usage to keep users informed
// Network failures are silently ignored to avoid disrupting normal operation
func CheckAndNotifyUpdate(currentVersion string) {
	// Skip for development builds
	if currentVersion == "dev" || currentVersion == "unknown" {
		return
	}

	// Check cache first to avoid excessive API calls
	cache, err := loadUpdateCacheFn()
	if err == nil && time.Since(cache.LastCheck) < checkInterval {
		// Recent check exists, use cached result
		if cache.UpdateAvailable {
			fmt.Fprintf(os.Stderr, "\n💡 Update available: %s (current: %s)\n", cache.LatestVersion, currentVersion)
			fmt.Fprintf(os.Stderr, "   Run 'promptext --update' to upgrade\n\n")
		}
		return
	}

	// Perform update check with short timeout (non-blocking)
	available, latestVersion, err := checkForUpdateWithTimeout(currentVersion, 3*time.Second)
	if err != nil {
		// Silently ignore errors (network issues, API limits, etc.)
		return
	}

	// Save result to cache
	newCache := UpdateCheckCache{
		LastCheck:       time.Now(),
		LatestVersion:   latestVersion,
		UpdateAvailable: available,
	}
	saveUpdateCacheFn(newCache) // Ignore errors

	// Notify user if update is available
	if available {
		fmt.Fprintf(os.Stderr, "\n💡 Update available: %s (current: %s)\n", latestVersion, currentVersion)
		fmt.Fprintf(os.Stderr, "   Run 'promptext --update' to upgrade\n\n")
	}
}

// checkForUpdateWithTimeout performs update check with configurable timeout
func checkForUpdateWithTimeout(currentVersion string, timeout time.Duration) (bool, string, error) {
	// Create channel for result
	type result struct {
		available bool
		version   string
		err       error
	}
	ch := make(chan result, 1)

	// Run check in goroutine
	go func() {
		available, version, err := checkForUpdateFn(currentVersion)
		ch <- result{available, version, err}
	}()

	// Wait for result or timeout
	select {
	case res := <-ch:
		return res.available, res.version, res.err
	case <-time.After(timeout):
		return false, "", fmt.Errorf("update check timed out")
	}
}

// getCacheDir returns the directory for storing update check cache
func getCacheDir() (string, error) {
	// Try to get user cache directory
	userHome := ""
	if u, err := user.Current(); err == nil {
		userHome = u.HomeDir
	}
	if userHome == "" {
		userHome = os.Getenv("HOME")
	}
	if userHome == "" {
		return "", fmt.Errorf("could not determine home directory")
	}

	// Use platform-specific cache directory
	var cacheDir string
	switch runtime.GOOS {
	case "darwin":
		cacheDir = filepath.Join(userHome, "Library", "Caches", "promptext")
	case "windows":
		appData := os.Getenv("LOCALAPPDATA")
		if appData == "" {
			appData = filepath.Join(userHome, "AppData", "Local")
		}
		cacheDir = filepath.Join(appData, "promptext", "cache")
	default: // linux and others
		xdgCache := os.Getenv("XDG_CACHE_HOME")
		if xdgCache != "" {
			cacheDir = filepath.Join(xdgCache, "promptext")
		} else {
			cacheDir = filepath.Join(userHome, ".cache", "promptext")
		}
	}

	// Create cache directory if it doesn't exist
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return cacheDir, nil
}

// loadUpdateCache loads the cached update check information
func loadUpdateCache() (*UpdateCheckCache, error) {
	cacheDir, err := getCacheDirFn()
	if err != nil {
		return nil, err
	}

	cachePath := filepath.Join(cacheDir, "update_check.json")
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, err
	}

	var cache UpdateCheckCache
	if err := json.Unmarshal(data, &cache); err != nil {
		return nil, err
	}

	return &cache, nil
}

// saveUpdateCache saves the update check cache
func saveUpdateCache(cache UpdateCheckCache) error {
	cacheDir, err := getCacheDirFn()
	if err != nil {
		return err
	}

	cachePath := filepath.Join(cacheDir, "update_check.json")
	data, err := json.Marshal(cache)
	if err != nil {
		return err
	}

	return os.WriteFile(cachePath, data, 0644)
}
