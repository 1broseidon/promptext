# Test Coverage Improvement Plan: 44.7% → 80%+

**Analysis Date:** 2025-11-05
**Current Overall Coverage:** 44.7%
**Target Coverage:** 80%+

---

## Executive Summary

This plan prioritizes test coverage improvements based on:
1. **Impact**: Critical business logic with low coverage
2. **Risk**: Code that handles file operations, updates, and user data
3. **ROI**: Maximum coverage gain per test effort

**Priority Order:**
1. **internal/processor** (19.5% → 75%): +55.5% coverage gain
2. **internal/update** (18.4% → 70%): +51.6% coverage gain
3. **internal/filter** (37.5% → 80%): +42.5% coverage gain
4. **internal/info** (43.0% → 80%): +37% coverage gain
5. **internal/format** (53.5% → 80%): +26.5% coverage gain
6. **internal/token** (62.9% → 80%): +17.1% coverage gain

**Estimated Total Coverage Gain:** ~45% → **~90% overall coverage**

---

## Current State Analysis

### Critical Packages (High Priority)

#### 1. internal/processor (19.5% coverage) - **HIGHEST PRIORITY**
**Why Critical:** Core business logic, file processing orchestration, token budget management

**Untested Functions (0% coverage):**
- `ParseCommaSeparated()` - Input parsing
- `PreviewDirectory()` - Dry-run functionality
- `filterDirectoryTree()` - Tree manipulation
- `prioritizeFiles()` - File prioritization logic
- `buildProjectHeader()` - Display formatting
- `formatTokenCount()` - Number formatting
- `analyzeFileStatistics()` - Statistics gathering
- `buildFileAnalysis()` - Analysis output
- `buildDependenciesSection()` - Dependency display
- `buildHealthSection()` - Health metrics display
- `buildGitSection()` - Git info display
- `formatBoxedOutput()` - Terminal UI formatting
- `GetMetadataSummary()` - Summary generation
- `formatSize()` - Size formatting
- `FormatDryRunOutput()` - Dry-run output
- `loadConfigurations()` - Config loading
- `handleDryRun()` - Dry-run handler
- `handleInfoOnly()` - Info mode handler
- `handleOutput()` - Output handling
- `Run()` - Main entry point
- `detectEntryPoints()` - Entry point detection

**Partially Tested Functions:**
- `validateFilePath()` (62.5%) - Missing error paths
- `checkFilePermissions()` (63.6%) - Missing edge cases
- `readFileContent()` (75%) - Missing error handling
- `processFile()` (63.6%) - Missing validation errors
- `populateProjectInfo()` (80%) - Missing nil checks
- `processFileInWalk()` (72%) - Missing error branches
- `ProcessDirectory()` (48.9%) - Main processing logic incomplete

**Test Strategy:**
- Create comprehensive integration tests for `Run()` and `ProcessDirectory()`
- Test file prioritization with various relevance scores
- Test token budget enforcement
- Test dry-run mode thoroughly
- Test all output formatting functions
- Test error paths in file operations

---

#### 2. internal/update (18.4% coverage) - **HIGH PRIORITY**
**Why Critical:** Self-update mechanism, network operations, file system manipulation

**Untested Functions (0% coverage):**
- `Update()` - Main update flow
- `getExecutablePath()` - Path resolution
- `findReleaseAssets()` - Asset discovery
- `downloadAndVerifyBinary()` - Download + verification
- `copyFile()` - File operations
- `replaceBinary()` - Binary replacement
- `downloadFile()` - HTTP download
- `verifyChecksum()` - Security validation
- `extractBinary()` - Archive extraction
- `extractTarGz()` - .tar.gz extraction
- `extractZip()` - .zip extraction
- `checkForUpdateWithTimeout()` - Timeout handling
- `loadUpdateCache()` - Cache loading
- `saveUpdateCache()` - Cache saving

**Partially Tested Functions:**
- `CheckForUpdate()` (81.8%) - Missing error cases
- `fetchLatestRelease()` (75%) - Missing HTTP errors
- `getPlatformAssetName()` (57.1%) - Missing platform combinations
- `CheckAndNotifyUpdate()` (12.5%) - Missing cache scenarios
- `getCacheDir()` (57.1%) - Missing platform-specific paths

**Test Strategy:**
- Mock HTTP calls using httptest
- Test archive extraction with real test archives
- Test checksum verification with known-good checksums
- Test cache behavior (loading, saving, expiry)
- Test error recovery (rollback on failure)
- Test cross-platform binary detection

---

#### 3. internal/filter (37.5% coverage) - **MEDIUM-HIGH PRIORITY**
**Why Critical:** File filtering is core to correct operation

**Untested Functions (0% coverage):**
- `ParseGitIgnore()` - .gitignore parsing
- `isTestFile()` - Test file detection
- `isEntryPoint()` - Entry point detection
- `getConfigType()` - Config file categorization
- `getDocType()` - Doc file categorization
- `getSourceType()` - Source file categorization
- `getDependencyType()` - Dependency file categorization
- `GetFileType()` - File type detection system

**Partially Tested Functions:**
- `New()` (88.9%) - Missing edge cases
- `ShouldProcess()` (73.3%) - Missing filter combinations

**Test Strategy:**
- Test .gitignore parsing with complex patterns
- Test file type detection for all supported languages
- Test filter rule precedence
- Test edge cases (symlinks, special chars in filenames)
- Integration tests for filter combinations

---

#### 4. internal/info (43.0% coverage) - **MEDIUM PRIORITY**
**Why Critical:** Project analysis and metadata extraction

**Untested Functions (0% coverage):**
- `getPythonVersion()` - Python version detection
- `getRustVersion()` - Rust version detection
- `getJavaVersion()` - Java version detection
- `getPythonDependencies()` - Python deps
- `getPipDependencies()` - pip deps
- `getPoetryDependencies()` - Poetry deps
- `getPoetryLockDependencies()` - Poetry lock parsing
- `getVenvDependencies()` - Virtual env deps
- `getRustDependencies()` - Cargo deps
- `getJavaMavenDependencies()` - Maven deps
- `getJavaGradleDependencies()` - Gradle deps

**Partially Tested Functions:**
- `getGitInfo()` (12.5%) - Missing error cases, git commands
- `detectLanguage()` (37.5%) - Missing language types
- `getLanguageVersion()` (42.9%) - Missing version parsers
- `getDependencies()` (37.5%) - Missing dependency types
- `isTestFile()` (54.5%) - Missing test patterns
- `checkForTestFiles()` (72.7%) - Missing edge cases
- `getGoVersion()` (75%) - Missing malformed go.mod
- `getNodeVersion()` (61.5%) - Missing package.json variations

**Test Strategy:**
- Create fixtures for each language ecosystem
- Test version extraction for Python/Rust/Java
- Test dependency parsing for all package managers
- Test git info extraction with various repo states
- Test project health analysis

---

#### 5. internal/format (53.5% coverage) - **MEDIUM PRIORITY**
**Why Critical:** Output generation for multiple formats

**Untested Functions (0% coverage):**
- `XMLFormatter.Format()` - XML output generation
- `treeToDirectoryMap()` - Tree conversion
- `buildDirectoryMap()` - Directory mapping
- `mapToList()` - List conversion
- `JSONLFormatter.Format()` - JSONL output
- `escapeForTOON()` - TOON string escaping
- `TOONStrictFormatter.Format()` - TOON strict output
- `encodeToJSON()` - JSON encoding

**Partially Tested Functions:**
- `GetFormatter()` (57.1%) - Missing format types
- `formatSourceFiles()` (91.7%) - Missing edge cases
- `writeDirectoryNode()` (15.4%) - Missing tree structures
- `formatGitInfo()` (22.2%) - Missing git scenarios
- `formatDependencies()` (11.8%) - Missing dependency cases

**Test Strategy:**
- Test each formatter (XML, JSONL, TOON-strict, PTX)
- Test edge cases (empty files, special characters)
- Test large file sets
- Validate output parsability
- Test format-specific escaping

---

#### 6. internal/token (62.9% coverage) - **LOWER PRIORITY**
**Why Lower:** Already has decent coverage, fewer critical paths

**Partially Tested Functions:**
- `init()` (63.6%) - Cache directory creation
- `NewTokenCounter()` (57.1%) - Fallback mode
- `EstimateTokens()` (83.3%) - Edge cases
- `DebugTokenCount()` (16.7%) - Debug output

**Test Strategy:**
- Test fallback mode when tiktoken unavailable
- Test approximation accuracy
- Test debug output formatting
- Test cache directory handling

---

## Detailed Implementation Plan

### Phase 1: Critical Path Testing (Week 1)
**Goal:** Cover 0% functions in processor and update packages
**Target:** 44.7% → 60% overall

#### Step 1.1: Processor Core Functions
**File:** `/home/george/Projects/personal/promptext/internal/processor/processor_test.go`

**New Test Cases:**
```go
// TestParseCommaSeparated tests input parsing
func TestParseCommaSeparated(t *testing.T) {
    tests := []struct{
        name     string
        input    string
        expected []string
    }{
        {"empty string", "", nil},
        {"single item", "foo", []string{"foo"}},
        {"multiple items", "foo,bar,baz", []string{"foo", "bar", "baz"}},
        {"with spaces", "foo, bar, baz", []string{"foo", " bar", " baz"}},
    }
    // Implementation...
}

// TestPreviewDirectory tests dry-run functionality
func TestPreviewDirectory(t *testing.T) {
    tmpDir := setupTestProject(t, map[string]string{
        "main.go": "package main",
        "README.md": "# Test",
    })
    defer os.RemoveAll(tmpDir)

    config := Config{
        DirPath: tmpDir,
        Filter: filter.New(filter.Options{UseDefaultRules: true}),
    }

    result, err := PreviewDirectory(config)
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Greater(t, len(result.FilePaths), 0)
    assert.Greater(t, result.EstimatedTokens, 0)
}

// TestPrioritizeFiles tests file prioritization logic
func TestPrioritizeFiles(t *testing.T) {
    files := []format.FileInfo{
        {Path: "main.go", Content: "package main"},
        {Path: "internal/auth.go", Content: "package auth"},
        {Path: "config.yaml", Content: "port: 8080"},
        {Path: "test_helper_test.go", Content: "package main_test"},
    }

    scorer := relevance.NewScorer("auth login")
    entryPoints := map[string]bool{"main.go": true}

    sorted := prioritizeFiles(files, scorer, entryPoints)

    // Entry points should come first
    assert.Equal(t, "main.go", sorted[0].Path)
    // High relevance files next
    assert.Contains(t, sorted[1].Path, "auth")
    // Tests should be last
    assert.Contains(t, sorted[len(sorted)-1].Path, "test")
}

// TestFilterDirectoryTree tests tree filtering
func TestFilterDirectoryTree(t *testing.T) {
    tree := &format.DirectoryNode{
        Name: "root",
        Type: "dir",
        Children: []*format.DirectoryNode{
            {Name: "included.go", Type: "file"},
            {Name: "excluded.go", Type: "file"},
            {Name: "subdir", Type: "dir", Children: []*format.DirectoryNode{
                {Name: "nested.go", Type: "file"},
            }},
        },
    }

    includedFiles := map[string]bool{
        "included.go": true,
        "subdir/nested.go": true,
    }

    filtered := filterDirectoryTree(tree, includedFiles, "")

    assert.Equal(t, 2, len(filtered.Children))
    // Verify excluded.go is not present
}

// TestFormatTokenCount tests number formatting
func TestFormatTokenCount(t *testing.T) {
    tests := []struct{
        tokens   int
        expected string
    }{
        {0, "0"},
        {999, "999"},
        {1000, "1,000"},
        {1234567, "1,234,567"},
    }

    for _, tt := range tests {
        result := formatTokenCount(tt.tokens)
        assert.Equal(t, tt.expected, result)
    }
}

// TestFormatSize tests size formatting
func TestFormatSize(t *testing.T) {
    tests := []struct{
        bytes    int64
        expected string
    }{
        {0, "0 B"},
        {1023, "1023 B"},
        {1024, "1.0 KB"},
        {1048576, "1.0 MB"},
        {1073741824, "1.0 GB"},
    }

    for _, tt := range tests {
        result := formatSize(tt.bytes)
        assert.Equal(t, tt.expected, result)
    }
}

// TestRun tests the main entry point
func TestRun(t *testing.T) {
    tmpDir := setupTestProject(t, map[string]string{
        "main.go": "package main\nfunc main() {}",
        ".gitignore": "*.log",
    })
    defer os.RemoveAll(tmpDir)

    // Test basic run
    err := Run(tmpDir, ".go", "", true, false, false, "markdown", "", false, true, true, false, false, "", 0, false)
    assert.NoError(t, err)

    // Test with relevance keywords
    err = Run(tmpDir, ".go", "", true, false, false, "markdown", "", false, true, true, false, false, "main", 0, false)
    assert.NoError(t, err)

    // Test with max tokens
    err = Run(tmpDir, ".go", "", true, false, false, "markdown", "", false, true, true, false, false, "", 1000, false)
    assert.NoError(t, err)
}
```

**Estimated Coverage Gain:** processor: 19.5% → 60% (+40.5%)

---

#### Step 1.2: Update Package Functions
**File:** `/home/george/Projects/personal/promptext/internal/update/update_test.go`

**New Test Cases:**
```go
// TestUpdate tests the full update flow with mocked HTTP
func TestUpdate(t *testing.T) {
    // Create mock HTTP server
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if strings.Contains(r.URL.Path, "releases/latest") {
            json.NewEncoder(w).Encode(ReleaseInfo{
                TagName: "v1.0.0",
                Assets: []struct{
                    Name string `json:"name"`
                    BrowserDownloadURL string `json:"browser_download_url"`
                }{
                    {Name: "promptext_Linux_x86_64.tar.gz", BrowserDownloadURL: "http://example.com/binary.tar.gz"},
                    {Name: "checksums.txt", BrowserDownloadURL: "http://example.com/checksums.txt"},
                },
            })
        }
    }))
    defer server.Close()

    // Test with mock server
    // ...
}

// TestDownloadFile tests file download
func TestDownloadFile(t *testing.T) {
    // Create test server
    content := "test binary content"
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(content))
    }))
    defer server.Close()

    tmpFile := filepath.Join(t.TempDir(), "download.bin")
    err := downloadFile(tmpFile, server.URL)
    assert.NoError(t, err)

    data, _ := os.ReadFile(tmpFile)
    assert.Equal(t, content, string(data))
}

// TestVerifyChecksum tests checksum verification
func TestVerifyChecksum(t *testing.T) {
    tmpDir := t.TempDir()

    // Create test file
    testContent := "test content"
    testFile := filepath.Join(tmpDir, "test.tar.gz")
    os.WriteFile(testFile, []byte(testContent), 0644)

    // Calculate checksum
    hash := sha256.Sum256([]byte(testContent))
    expectedChecksum := hex.EncodeToString(hash[:])

    // Create checksums file
    checksumFile := filepath.Join(tmpDir, "checksums.txt")
    checksumContent := fmt.Sprintf("%s  test.tar.gz\n", expectedChecksum)
    os.WriteFile(checksumFile, []byte(checksumContent), 0644)

    // Test verification
    err := verifyChecksum(testFile, checksumFile, "test.tar.gz")
    assert.NoError(t, err)

    // Test with wrong checksum
    os.WriteFile(checksumFile, []byte("wrongchecksum  test.tar.gz\n"), 0644)
    err = verifyChecksum(testFile, checksumFile, "test.tar.gz")
    assert.Error(t, err)
}

// TestExtractTarGz tests tar.gz extraction
func TestExtractTarGz(t *testing.T) {
    tmpDir := t.TempDir()

    // Create a test tar.gz archive
    archivePath := filepath.Join(tmpDir, "test.tar.gz")
    // ... create archive with test binary

    extractDir := filepath.Join(tmpDir, "extract")
    os.MkdirAll(extractDir, 0755)

    binaryPath, err := extractTarGz(archivePath, extractDir)
    assert.NoError(t, err)
    assert.FileExists(t, binaryPath)
}

// TestExtractZip tests zip extraction
func TestExtractZip(t *testing.T) {
    tmpDir := t.TempDir()

    // Create test zip archive
    archivePath := filepath.Join(tmpDir, "test.zip")
    // ... create archive

    extractDir := filepath.Join(tmpDir, "extract")
    os.MkdirAll(extractDir, 0755)

    binaryPath, err := extractZip(archivePath, extractDir)
    assert.NoError(t, err)
    assert.FileExists(t, binaryPath)
}

// TestCopyFile tests file copying
func TestCopyFile(t *testing.T) {
    tmpDir := t.TempDir()

    srcFile := filepath.Join(tmpDir, "source.txt")
    content := "test content"
    os.WriteFile(srcFile, []byte(content), 0644)

    dstFile := filepath.Join(tmpDir, "dest.txt")
    err := copyFile(srcFile, dstFile)
    assert.NoError(t, err)

    data, _ := os.ReadFile(dstFile)
    assert.Equal(t, content, string(data))
}

// TestReplaceBinary tests binary replacement
func TestReplaceBinary(t *testing.T) {
    tmpDir := t.TempDir()

    // Create fake current binary
    execPath := filepath.Join(tmpDir, "promptext")
    os.WriteFile(execPath, []byte("old version"), 0755)

    // Create new binary
    newBinary := filepath.Join(tmpDir, "promptext-new")
    os.WriteFile(newBinary, []byte("new version"), 0755)

    err := replaceBinary(execPath, newBinary, false)
    assert.NoError(t, err)

    // Verify replacement
    data, _ := os.ReadFile(execPath)
    assert.Equal(t, "new version", string(data))
}

// TestUpdateCache tests cache loading and saving
func TestLoadSaveUpdateCache(t *testing.T) {
    cache := UpdateCheckCache{
        LastCheck: time.Now(),
        LatestVersion: "v1.0.0",
        UpdateAvailable: true,
    }

    err := saveUpdateCache(cache)
    assert.NoError(t, err)

    loaded, err := loadUpdateCache()
    assert.NoError(t, err)
    assert.Equal(t, cache.LatestVersion, loaded.LatestVersion)
    assert.Equal(t, cache.UpdateAvailable, loaded.UpdateAvailable)
}
```

**Estimated Coverage Gain:** update: 18.4% → 70% (+51.6%)

---

### Phase 2: Filter and Info Testing (Week 2)
**Goal:** Complete filter and info package coverage
**Target:** 60% → 70% overall

#### Step 2.1: Filter Package
**File:** `/home/george/Projects/personal/promptext/internal/filter/filter_test.go`

**New Test Cases:**
```go
// TestParseGitIgnore tests .gitignore parsing
func TestParseGitIgnore(t *testing.T) {
    tmpDir := t.TempDir()

    gitignore := `# Comment
*.log
node_modules/
/build
!important.log
`
    gitignorePath := filepath.Join(tmpDir, ".gitignore")
    os.WriteFile(gitignorePath, []byte(gitignore), 0644)

    patterns, err := ParseGitIgnore(tmpDir)
    assert.NoError(t, err)
    assert.Contains(t, patterns, "*.log")
    assert.Contains(t, patterns, "node_modules/")
    assert.Contains(t, patterns, "/build")
    assert.Contains(t, patterns, "!important.log")
    assert.NotContains(t, patterns, "# Comment")
}

// TestGetFileType tests comprehensive file type detection
func TestGetFileType(t *testing.T) {
    tests := []struct{
        path         string
        expectedType string
        expectedCat  string
        isTest       bool
        isEntry      bool
    }{
        {"main.go", "source", "entry:go", false, true},
        {"auth_test.go", "test", "test:go", true, false},
        {"config.yaml", "config", "config:yaml", false, false},
        {"README.md", "doc", "doc:markdown", false, false},
        {"index.js", "source", "entry:js", false, true},
        {"package.json", "dependency", "dep:node", false, false},
        {"Cargo.toml", "dependency", "dep:rust", false, false},
    }

    filter := New(Options{UseDefaultRules: true})

    for _, tt := range tests {
        info := GetFileType(tt.path, filter)
        assert.Equal(t, tt.expectedType, info.Type, "path: %s", tt.path)
        assert.Equal(t, tt.expectedCat, info.Category, "path: %s", tt.path)
        assert.Equal(t, tt.isTest, info.IsTest, "path: %s", tt.path)
        assert.Equal(t, tt.isEntry, info.IsEntryPoint, "path: %s", tt.path)
    }
}

// TestFilterHelpers tests helper functions
func TestIsTestFile(t *testing.T) {
    tests := []struct{
        path     string
        base     string
        expected bool
    }{
        {"auth_test.go", "auth_test.go", true},
        {"test_utils.py", "test_utils.py", true},
        {"app.test.js", "app.test.js", true},
        {"main.go", "main.go", false},
    }

    for _, tt := range tests {
        result := isTestFile(tt.path, tt.base)
        assert.Equal(t, tt.expected, result, "path: %s", tt.path)
    }
}
```

**Estimated Coverage Gain:** filter: 37.5% → 80% (+42.5%)

---

#### Step 2.2: Info Package
**File:** `/home/george/Projects/personal/promptext/internal/info/info_test.go`

**New Test Cases:**
```go
// TestPythonVersionDetection tests Python version extraction
func TestGetPythonVersion(t *testing.T) {
    tmpDir := t.TempDir()

    pyproject := `[tool.poetry]
name = "myproject"

[tool.poetry.dependencies]
python = "^3.9"
`
    os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644)

    version := getPythonVersion(tmpDir)
    assert.Contains(t, version, "3.9")
}

// TestRustVersionDetection tests Rust version extraction
func TestGetRustVersion(t *testing.T) {
    tmpDir := t.TempDir()

    cargoToml := `[package]
name = "myproject"
version = "0.1.0"
`
    os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoToml), 0644)

    version := getRustVersion(tmpDir)
    assert.Equal(t, "0.1.0", version)
}

// TestPythonDependencies tests Python dependency parsing
func TestGetPythonDependencies(t *testing.T) {
    tmpDir := t.TempDir()

    // Test requirements.txt
    requirements := `pytest==7.3.1
requests>=2.31.0
black
`
    os.WriteFile(filepath.Join(tmpDir, "requirements.txt"), []byte(requirements), 0644)

    deps := getPythonDependencies(tmpDir)
    assert.Contains(t, deps, "pytest")
    assert.Contains(t, deps, "requests")
    assert.Contains(t, deps, "black")
}

// TestPoetryDependencies tests Poetry dependency parsing
func TestGetPoetryDependencies(t *testing.T) {
    tmpDir := t.TempDir()

    pyproject := `[tool.poetry.dependencies]
python = "^3.9"
flask = "^2.0.0"
requests = "^2.28.0"

[tool.poetry.group.dev.dependencies]
pytest = "^7.0.0"
black = "^22.0.0"
`
    os.WriteFile(filepath.Join(tmpDir, "pyproject.toml"), []byte(pyproject), 0644)

    depsMap := make(map[string]bool)
    getPoetryDependencies(tmpDir, depsMap)

    assert.True(t, depsMap["flask"])
    assert.True(t, depsMap["requests"])
    assert.True(t, depsMap["[dev] pytest"])
    assert.True(t, depsMap["[dev] black"])
    assert.False(t, depsMap["python"]) // Should skip python version
}

// TestRustDependencies tests Cargo dependency parsing
func TestGetRustDependencies(t *testing.T) {
    tmpDir := t.TempDir()

    cargoToml := `[dependencies]
serde = "1.0"
tokio = { version = "1.0", features = ["full"] }
`
    os.WriteFile(filepath.Join(tmpDir, "Cargo.toml"), []byte(cargoToml), 0644)

    deps := getRustDependencies(tmpDir)
    assert.Contains(t, deps, "serde")
    assert.Contains(t, deps, "tokio")
}

// TestJavaDependencies tests Maven dependency parsing
func TestGetJavaMavenDependencies(t *testing.T) {
    tmpDir := t.TempDir()

    pomXml := `<project>
    <dependencies>
        <dependency>
            <groupId>org.springframework</groupId>
            <artifactId>spring-core</artifactId>
        </dependency>
        <dependency>
            <groupId>junit</groupId>
            <artifactId>junit</artifactId>
        </dependency>
    </dependencies>
</project>`
    os.WriteFile(filepath.Join(tmpDir, "pom.xml"), []byte(pomXml), 0644)

    deps := getJavaMavenDependencies(tmpDir)
    assert.Contains(t, deps, "spring-core")
    assert.Contains(t, deps, "junit")
}

// TestGitInfoExtraction tests git information gathering
func TestGetGitInfo(t *testing.T) {
    tmpDir := t.TempDir()

    // Initialize git repo
    exec.Command("git", "init", tmpDir).Run()
    exec.Command("git", "-C", tmpDir, "config", "user.name", "Test User").Run()
    exec.Command("git", "-C", tmpDir, "config", "user.email", "test@example.com").Run()

    testFile := filepath.Join(tmpDir, "test.txt")
    os.WriteFile(testFile, []byte("test"), 0644)

    exec.Command("git", "-C", tmpDir, "add", ".").Run()
    exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial commit").Run()

    gitInfo, err := getGitInfo(tmpDir)
    assert.NoError(t, err)
    assert.NotEmpty(t, gitInfo.Branch)
    assert.NotEmpty(t, gitInfo.CommitHash)
    assert.Contains(t, gitInfo.CommitMessage, "Initial commit")
}
```

**Estimated Coverage Gain:** info: 43.0% → 80% (+37%)

---

### Phase 3: Format and Token Testing (Week 3)
**Goal:** Complete format and token package coverage
**Target:** 70% → 80% overall

#### Step 3.1: Format Package
**File:** `/home/george/Projects/personal/promptext/internal/format/format_test.go`

**New Test Cases:**
```go
// TestXMLFormatter tests XML output generation
func TestXMLFormatter(t *testing.T) {
    formatter := &XMLFormatter{}

    project := &ProjectOutput{
        Files: []FileInfo{
            {Path: "main.go", Content: "package main"},
        },
        GitInfo: &GitInfo{
            Branch: "main",
            CommitHash: "abc123",
        },
    }

    output, err := formatter.Format(project)
    assert.NoError(t, err)
    assert.Contains(t, output, "<?xml")
    assert.Contains(t, output, "main.go")
    assert.Contains(t, output, "package main")
}

// TestJSONLFormatter tests JSONL output generation
func TestJSONLFormatter(t *testing.T) {
    formatter := &JSONLFormatter{}

    project := &ProjectOutput{
        Files: []FileInfo{
            {Path: "main.go", Content: "package main", Tokens: 10},
            {Path: "util.go", Content: "package util", Tokens: 8},
        },
    }

    output, err := formatter.Format(project)
    assert.NoError(t, err)

    // Each line should be valid JSON
    lines := strings.Split(strings.TrimSpace(output), "\n")
    assert.Equal(t, 2, len(lines))

    for _, line := range lines {
        var obj map[string]interface{}
        err := json.Unmarshal([]byte(line), &obj)
        assert.NoError(t, err)
        assert.Contains(t, obj, "path")
        assert.Contains(t, obj, "content")
    }
}

// TestTOONStrictFormatter tests TOON v1.3 strict output
func TestTOONStrictFormatter(t *testing.T) {
    formatter := &TOONStrictFormatter{}

    project := &ProjectOutput{
        Files: []FileInfo{
            {Path: "main.go", Content: "package main\nfunc main() {}"},
        },
    }

    output, err := formatter.Format(project)
    assert.NoError(t, err)

    // TOON strict should escape newlines
    assert.Contains(t, output, "\\n")
    assert.NotContains(t, output, "\nfunc") // Should be escaped
}

// TestEscapeForTOON tests TOON escaping
func TestEscapeForTOON(t *testing.T) {
    tests := []struct{
        input    string
        expected string
    }{
        {"simple", "simple"},
        {"with\nnewline", "with\\nnewline"},
        {"with\ttab", "with\\ttab"},
        {"with\"quote", "with\\\"quote"},
        {"with\\backslash", "with\\\\backslash"},
    }

    for _, tt := range tests {
        result := escapeForTOON(tt.input)
        assert.Equal(t, tt.expected, result)
    }
}

// TestGetFormatter tests formatter selection
func TestGetFormatter(t *testing.T) {
    tests := []struct{
        format      string
        expectError bool
        formatType  string
    }{
        {"markdown", false, "MarkdownFormatter"},
        {"md", false, "MarkdownFormatter"},
        {"xml", false, "XMLFormatter"},
        {"ptx", false, "PTXFormatter"},
        {"toon", false, "PTXFormatter"},
        {"jsonl", false, "JSONLFormatter"},
        {"toon-strict", false, "TOONStrictFormatter"},
        {"invalid", true, ""},
    }

    for _, tt := range tests {
        formatter, err := GetFormatter(tt.format)
        if tt.expectError {
            assert.Error(t, err)
        } else {
            assert.NoError(t, err)
            assert.NotNil(t, formatter)
        }
    }
}
```

**Estimated Coverage Gain:** format: 53.5% → 80% (+26.5%)

---

#### Step 3.2: Token Package
**File:** `/home/george/Projects/personal/promptext/internal/token/tiktoken_test.go`

**New Test Cases:**
```go
// TestTokenCounterFallback tests fallback mode
func TestTokenCounterFallback(t *testing.T) {
    // Temporarily break tiktoken to force fallback
    oldCache := os.Getenv("TIKTOKEN_CACHE_DIR")
    os.Setenv("TIKTOKEN_CACHE_DIR", "/nonexistent/path")
    defer os.Setenv("TIKTOKEN_CACHE_DIR", oldCache)

    tc := NewTokenCounter()

    // Should use fallback mode
    assert.True(t, tc.IsFallbackMode())
    assert.Equal(t, "approximation", tc.GetEncodingName())

    // Should still count tokens (using approximation)
    tokens := tc.EstimateTokens("Hello, world!")
    assert.Greater(t, tokens, 0)
}

// TestDebugTokenCount tests debug output
func TestDebugTokenCount(t *testing.T) {
    tc := NewTokenCounter()

    // Enable debug logging
    log.Enable()
    defer log.Disable()

    text := "This is a test string with multiple words and punctuation!"
    tokens := tc.DebugTokenCount(text, "test string")

    assert.Greater(t, tokens, 0)
    // Debug output should be logged (verify via log capture if needed)
}

// TestApproximationAccuracy tests approximation vs tiktoken
func TestApproximationAccuracy(t *testing.T) {
    tc := NewTokenCounter()

    testTexts := []string{
        "Simple text",
        "Code with {braces} and (parentheses)",
        `Multi
line
text`,
        "function main() { return 42; }",
    }

    for _, text := range testTexts {
        tokens := tc.EstimateTokens(text)
        assert.Greater(t, tokens, 0)

        // Approximation should be reasonable (within 2x of char/4)
        charEstimate := len(text) / 4
        assert.Less(t, tokens, charEstimate*2)
    }
}
```

**Estimated Coverage Gain:** token: 62.9% → 80% (+17.1%)

---

### Phase 4: Edge Cases and Integration (Week 4)
**Goal:** Reach 80%+ overall coverage with edge case testing
**Target:** 80% → 85%+

#### Integration Tests

**File:** `/home/george/Projects/personal/promptext/internal/processor/integration_test.go`

```go
// TestEndToEndWorkflow tests complete workflow
func TestEndToEndWorkflow(t *testing.T) {
    tmpDir := setupTestProject(t, map[string]string{
        "go.mod": "module example.com/test\ngo 1.19",
        "main.go": "package main\nfunc main() { println(\"Hello\") }",
        "internal/auth/auth.go": "package auth\nfunc Login() {}",
        "internal/auth/auth_test.go": "package auth\nfunc TestLogin(t *testing.T) {}",
        "config.yaml": "port: 8080",
        "README.md": "# Test Project",
        ".gitignore": "*.log\n/build",
    })
    defer os.RemoveAll(tmpDir)

    // Initialize git repo
    exec.Command("git", "init", tmpDir).Run()
    exec.Command("git", "-C", tmpDir, "config", "user.name", "Test").Run()
    exec.Command("git", "-C", tmpDir, "config", "user.email", "test@test.com").Run()
    exec.Command("git", "-C", tmpDir, "add", ".").Run()
    exec.Command("git", "-C", tmpDir, "commit", "-m", "Initial").Run()

    // Test with relevance filtering
    outFile := filepath.Join(tmpDir, "output.ptx")
    err := Run(tmpDir, ".go", "", true, false, false, "ptx", outFile, false, true, true, false, false, "auth", 5000, false)
    assert.NoError(t, err)

    // Verify output file
    assert.FileExists(t, outFile)
    content, _ := os.ReadFile(outFile)
    assert.Contains(t, string(content), "auth.go")
    assert.NotContains(t, string(content), "auth_test.go") // Tests should be deprioritized
}

// TestRelevanceWithTokenBudget tests combined relevance and budget
func TestRelevanceWithTokenBudget(t *testing.T) {
    // Create project with many files
    files := map[string]string{}
    for i := 0; i < 100; i++ {
        files[fmt.Sprintf("file%d.go", i)] = fmt.Sprintf("package main\n// File %d", i)
    }
    files["auth.go"] = "package main\n// Authentication logic with auth and login keywords"

    tmpDir := setupTestProject(t, files)
    defer os.RemoveAll(tmpDir)

    filter := filter.New(filter.Options{UseDefaultRules: true})
    config := Config{
        DirPath: tmpDir,
        Filter: filter,
        RelevanceKeywords: "auth login",
        MaxTokens: 1000,
    }

    result, err := ProcessDirectory(config, false)
    assert.NoError(t, err)

    // Should include auth.go due to high relevance
    hasAuth := false
    for _, file := range result.ProjectOutput.Files {
        if strings.Contains(file.Path, "auth") {
            hasAuth = true
            break
        }
    }
    assert.True(t, hasAuth)

    // Should respect token budget
    assert.LessOrEqual(t, result.TokenCount, 1000)
}
```

---

## Testing Best Practices

### 1. Table-Driven Tests
Use table-driven tests for comprehensive input coverage:
```go
tests := []struct{
    name     string
    input    X
    expected Y
    wantErr  bool
}{
    // Test cases...
}
```

### 2. Test Helpers
Create reusable test helpers:
```go
func setupTestProject(t *testing.T, files map[string]string) string
func setupGitRepo(t *testing.T, dir string)
func createTestArchive(t *testing.T, format string) string
```

### 3. Cleanup
Always cleanup temporary resources:
```go
tmpDir := t.TempDir() // Automatic cleanup
// or
defer os.RemoveAll(tmpDir)
```

### 4. Mock External Dependencies
- Use `httptest.NewServer()` for HTTP calls
- Mock file system when possible
- Mock git commands for predictable tests

### 5. Test Error Paths
Explicitly test error conditions:
```go
// Test with invalid input
result, err := Function(invalidInput)
assert.Error(t, err)
assert.Nil(t, result)
```

---

## Expected Coverage After Each Phase

| Phase | Package | Before | After | Gain |
|-------|---------|--------|-------|------|
| 1 | processor | 19.5% | 60% | +40.5% |
| 1 | update | 18.4% | 70% | +51.6% |
| 2 | filter | 37.5% | 80% | +42.5% |
| 2 | info | 43.0% | 80% | +37% |
| 3 | format | 53.5% | 80% | +26.5% |
| 3 | token | 62.9% | 80% | +17.1% |
| 4 | Integration | - | - | +5% |

**Overall Coverage Projection:**
- Phase 1: 44.7% → 60%
- Phase 2: 60% → 70%
- Phase 3: 70% → 80%
- Phase 4: 80% → 85%+

---

## Priority Quick Reference

### Must Test (Critical - Week 1)
1. `processor.Run()` - Main entry point
2. `processor.ProcessDirectory()` - Core processing
3. `processor.prioritizeFiles()` - File selection logic
4. `update.Update()` - Update mechanism
5. `update.downloadAndVerifyBinary()` - Download + security

### Should Test (High Priority - Week 2)
6. `filter.GetFileType()` - Type detection
7. `filter.ParseGitIgnore()` - .gitignore parsing
8. `info` language-specific functions - Metadata extraction
9. `processor.PreviewDirectory()` - Dry-run mode

### Nice to Have (Medium Priority - Week 3)
10. All formatter `Format()` methods - Output generation
11. `token.DebugTokenCount()` - Debug output
12. Remaining helper functions - Utilities

---

## Success Criteria

- [ ] Overall coverage ≥ 80%
- [ ] processor coverage ≥ 75%
- [ ] update coverage ≥ 70%
- [ ] filter coverage ≥ 80%
- [ ] info coverage ≥ 80%
- [ ] format coverage ≥ 80%
- [ ] All critical paths have error case tests
- [ ] Integration tests cover end-to-end workflows
- [ ] CI/CD passes all tests consistently

---

## Maintenance Notes

After reaching 80% coverage:
1. Require tests for all new functions
2. Maintain coverage in CI (fail if coverage drops below 75%)
3. Regular coverage reviews during code reviews
4. Add coverage badge to README
5. Document untested edge cases with TODO comments

---

**Document Version:** 1.0
**Last Updated:** 2025-11-05
**Author:** QA Engineer (Claude Code)
