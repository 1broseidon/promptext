# Test Coverage Implementation Checklist

## Week 1: Critical Path Testing (44.7% → 60%)

### internal/processor Package
- [ ] TestParseCommaSeparated - Input parsing
- [ ] TestPreviewDirectory - Dry-run functionality  
- [ ] TestPrioritizeFiles - File prioritization logic
- [ ] TestFilterDirectoryTree - Tree filtering
- [ ] TestFormatTokenCount - Number formatting
- [ ] TestFormatSize - Size formatting
- [ ] TestDetectEntryPoints - Entry point detection
- [ ] TestBuildProjectHeader - Header generation
- [ ] TestBuildFileAnalysis - Statistics analysis
- [ ] TestGetMetadataSummary - Summary generation
- [ ] TestFormatDryRunOutput - Dry-run output
- [ ] TestHandleDryRun - Dry-run handler
- [ ] TestHandleInfoOnly - Info mode handler
- [ ] TestHandleOutput - Output handling
- [ ] TestRun - Main entry point (integration)
- [ ] TestRunWithRelevance - Relevance filtering
- [ ] TestRunWithTokenBudget - Budget enforcement
- [ ] TestValidateFilePath - Error paths
- [ ] TestCheckFilePermissions - Edge cases
- [ ] TestProcessFile - Validation errors

**Target:** processor 19.5% → 60%

### internal/update Package
- [ ] TestUpdate - Full update flow (mocked HTTP)
- [ ] TestGetExecutablePath - Path resolution
- [ ] TestFindReleaseAssets - Asset discovery
- [ ] TestDownloadFile - HTTP download
- [ ] TestVerifyChecksum - Checksum validation
- [ ] TestVerifyChecksumMismatch - Bad checksum
- [ ] TestExtractTarGz - .tar.gz extraction
- [ ] TestExtractZip - .zip extraction
- [ ] TestCopyFile - File operations
- [ ] TestReplaceBinary - Binary replacement
- [ ] TestReplaceBinaryRollback - Rollback on failure
- [ ] TestLoadUpdateCache - Cache loading
- [ ] TestSaveUpdateCache - Cache saving
- [ ] TestGetCacheDirLinux - Linux cache path
- [ ] TestGetCacheDirDarwin - macOS cache path
- [ ] TestGetCacheDirWindows - Windows cache path

**Target:** update 18.4% → 70%

---

## Week 2: Filter & Info Testing (60% → 70%)

### internal/filter Package
- [ ] TestParseGitIgnore - .gitignore parsing
- [ ] TestParseGitIgnoreComments - Comment handling
- [ ] TestParseGitIgnoreEmpty - Empty file
- [ ] TestGetFileType - Comprehensive type detection
- [ ] TestGetFileTypeGo - Go files
- [ ] TestGetFileTypeJS - JavaScript files
- [ ] TestGetFileTypePython - Python files
- [ ] TestGetFileTypeRust - Rust files
- [ ] TestIsTestFile - Test detection patterns
- [ ] TestIsEntryPoint - Entry point detection
- [ ] TestGetConfigType - Config categorization
- [ ] TestGetDocType - Doc categorization
- [ ] TestGetSourceType - Source categorization
- [ ] TestGetDependencyType - Dependency categorization

**Target:** filter 37.5% → 80%

### internal/info Package
- [ ] TestGetPythonVersion - Python version extraction
- [ ] TestGetPythonVersionPoetry - Poetry version
- [ ] TestGetRustVersion - Rust version extraction
- [ ] TestGetJavaVersion - Java version extraction
- [ ] TestGetPipDependencies - pip requirements.txt
- [ ] TestGetPoetryDependencies - pyproject.toml parsing
- [ ] TestGetPoetryLockDependencies - poetry.lock parsing
- [ ] TestGetVenvDependencies - Virtual env deps
- [ ] TestGetRustDependencies - Cargo.toml parsing
- [ ] TestGetJavaMavenDependencies - pom.xml parsing
- [ ] TestGetJavaGradleDependencies - build.gradle parsing
- [ ] TestGetGitInfo - Git info extraction
- [ ] TestGetGitInfoNoRepo - Non-git directory
- [ ] TestGetGitInfoDetachedHead - Detached HEAD state
- [ ] TestCheckForTestFiles - Test file discovery
- [ ] TestDetectLanguage - Language detection
- [ ] TestGetLanguageVersion - Version extraction
- [ ] TestGetDependencies - Dependency routing
- [ ] TestAnalyzeProjectHealth - Health metrics
- [ ] TestIsCoreFile - Core file detection

**Target:** info 43.0% → 80%

---

## Week 3: Format & Token Testing (70% → 80%)

### internal/format Package
- [ ] TestXMLFormatter - XML output generation
- [ ] TestXMLFormatterEmpty - Empty project
- [ ] TestJSONLFormatter - JSONL output
- [ ] TestJSONLFormatterParsing - Validate JSON per line
- [ ] TestTOONStrictFormatter - TOON v1.3 strict
- [ ] TestEscapeForTOON - String escaping
- [ ] TestEscapeForTOONEdgeCases - Special characters
- [ ] TestGetFormatter - Formatter selection
- [ ] TestGetFormatterInvalid - Invalid format
- [ ] TestWriteDirectoryNode - Directory tree output
- [ ] TestFormatGitInfo - Git info formatting
- [ ] TestFormatDependencies - Dependency formatting

**Target:** format 53.5% → 80%

### internal/token Package
- [ ] TestTokenCounterFallback - Fallback mode
- [ ] TestTokenCounterFallbackAccuracy - Approximation
- [ ] TestDebugTokenCount - Debug output
- [ ] TestDebugTokenCountLogging - Log capture
- [ ] TestApproximationCode - Code token estimation
- [ ] TestApproximationProse - Prose token estimation
- [ ] TestIsLikelyCode - Code detection
- [ ] TestNewTokenCounterError - tiktoken error handling

**Target:** token 62.9% → 80%

---

## Week 4: Integration & Edge Cases (80% → 85%+)

### Integration Tests
- [ ] TestEndToEndWorkflow - Complete workflow
- [ ] TestEndToEndWithGit - With git integration
- [ ] TestRelevanceWithTokenBudget - Combined features
- [ ] TestMultipleFormatters - All output formats
- [ ] TestLargeProject - Performance with 100+ files
- [ ] TestErrorRecovery - Graceful error handling
- [ ] TestCrossFilesystem - Cross-filesystem operations
- [ ] TestSymlinks - Symlink handling
- [ ] TestSpecialCharacters - Filenames with special chars
- [ ] TestUnicode - Unicode content handling

**Target:** Overall 80% → 85%+

---

## Coverage Verification Commands

```bash
# Check current coverage
go test -cover ./...

# Generate coverage report
go test -coverprofile=/tmp/coverage.out ./...
go tool cover -func=/tmp/coverage.out

# View HTML report
go tool cover -html=/tmp/coverage.out

# Check specific package
go test -cover ./internal/processor/...

# Run with race detection
go test -race ./...
```

---

## Progress Tracking

### Week 1 Progress
- [ ] Day 1: Processor core functions (5 tests)
- [ ] Day 2: Processor formatting functions (5 tests)
- [ ] Day 3: Processor Run + integration (5 tests)
- [ ] Day 4: Update download/verify (8 tests)
- [ ] Day 5: Update extraction/cache (8 tests)

### Week 2 Progress
- [ ] Day 1: Filter type detection (7 tests)
- [ ] Day 2: Filter gitignore parsing (7 tests)
- [ ] Day 3: Info Python ecosystem (6 tests)
- [ ] Day 4: Info Rust/Java ecosystem (7 tests)
- [ ] Day 5: Info git and health (7 tests)

### Week 3 Progress
- [ ] Day 1: Format XML/JSONL (6 tests)
- [ ] Day 2: Format TOON-strict (6 tests)
- [ ] Day 3: Token fallback mode (4 tests)
- [ ] Day 4: Token approximation (4 tests)
- [ ] Day 5: Buffer day / cleanup

### Week 4 Progress
- [ ] Day 1: Integration tests (3 scenarios)
- [ ] Day 2: Edge case tests (3 scenarios)
- [ ] Day 3: Performance tests (2 scenarios)
- [ ] Day 4: Error recovery (2 scenarios)
- [ ] Day 5: Final verification + cleanup

---

## Success Criteria

✅ **Must Have (Week 1-2)**
- [ ] processor coverage ≥ 60%
- [ ] update coverage ≥ 70%
- [ ] filter coverage ≥ 80%
- [ ] info coverage ≥ 80%
- [ ] Overall coverage ≥ 70%

✅ **Should Have (Week 3)**
- [ ] format coverage ≥ 80%
- [ ] token coverage ≥ 80%
- [ ] Overall coverage ≥ 80%

✅ **Nice to Have (Week 4)**
- [ ] Integration tests complete
- [ ] Edge cases covered
- [ ] Overall coverage ≥ 85%
- [ ] CI/CD green
- [ ] No flaky tests

---

## Notes

- Focus on critical paths first (processor, update)
- Use table-driven tests for comprehensive coverage
- Mock external dependencies (HTTP, git commands)
- Always cleanup temporary resources
- Test error paths explicitly
- Document any intentionally untested code

---

**Last Updated:** 2025-11-05
**Status:** Ready to start
