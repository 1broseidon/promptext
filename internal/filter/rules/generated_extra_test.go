package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHasGeneratedMarkers(t *testing.T) {
	dir := t.TempDir()
	plain := filepath.Join(dir, "plain.txt")
	marked := filepath.Join(dir, "marked.txt")

	if err := os.WriteFile(plain, []byte("regular content"), 0644); err != nil {
		t.Fatalf("write plain: %v", err)
	}
	if err := os.WriteFile(marked, []byte("// DO NOT EDIT - generated file"), 0644); err != nil {
		t.Fatalf("write marked: %v", err)
	}

	if hasGeneratedMarkers(plain) {
		t.Fatalf("expected plain file to have no generated markers")
	}
	if !hasGeneratedMarkers(marked) {
		t.Fatalf("expected marker file to be detected as generated")
	}
}

func TestHasLowEntropy(t *testing.T) {
	dir := t.TempDir()
	highEntropy := filepath.Join(dir, "high.txt")
	lowEntropy := filepath.Join(dir, "low.txt")

	var repetitive strings.Builder
	for i := 0; i < 60; i++ {
		repetitive.WriteString("value = 123\n")
	}

	if err := os.WriteFile(highEntropy, []byte("distinct\nlines\nthroughout\n"), 0644); err != nil {
		t.Fatalf("write high entropy: %v", err)
	}
	if err := os.WriteFile(lowEntropy, []byte(repetitive.String()), 0644); err != nil {
		t.Fatalf("write low entropy: %v", err)
	}

	if hasLowEntropy(highEntropy) {
		t.Fatalf("expected varied file to have high entropy")
	}
	if !hasLowEntropy(lowEntropy) {
		t.Fatalf("expected repetitive file to be flagged as generated")
	}
}

func TestNormalizeLineToPattern(t *testing.T) {
	line := "package version 1.2.3 with hash deadbeefdeadbeefdeadbeefdeadbeef and \"long quoted string\""
	normalized := normalizeLineToPattern(line)
	if strings.Contains(normalized, "1.2.3") {
		t.Fatalf("expected version numbers to be normalized")
	}
	if strings.Contains(normalized, "deadbeef") {
		t.Fatalf("expected hashes to be normalized")
	}
	if strings.Contains(normalized, "long quoted string") {
		t.Fatalf("expected quoted strings to be normalized")
	}
}

func TestHasLockFileSignatures(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "package-lock.json")
	content := "lockfileVersion\nresolved\nintegrity"
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatalf("write lock file: %v", err)
	}

	if !hasLockFileSignatures(path, []string{"lockfileVersion", "resolved", "integrity"}) {
		t.Fatalf("expected signatures to be detected")
	}
	if hasLockFileSignatures(path, []string{"missing", "entries"}) {
		t.Fatalf("unexpected match for unrelated signatures")
	}
}

func TestLooksLikeLockFile(t *testing.T) {
	dir := t.TempDir()
	lock := filepath.Join(dir, "yarn.lock")
	noLock := filepath.Join(dir, "README.md")

	if err := os.WriteFile(lock, []byte("# THIS IS A LOCKFILE\nlockfile v1\n"), 0644); err != nil {
		t.Fatalf("write lock: %v", err)
	}
	if err := os.WriteFile(noLock, []byte("regular documentation"), 0644); err != nil {
		t.Fatalf("write doc: %v", err)
	}

	if !looksLikeLockFile(lock) {
		t.Fatalf("expected lockfile heuristics to match")
	}
	if looksLikeLockFile(noLock) {
		t.Fatalf("did not expect regular file to look like lock file")
	}
}

func TestEcosystemRuleDetectsLockFiles(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"package.json":     `{}`,
		"go.mod":           `module example.com/test`,
		"MyProject.csproj": `<Project></Project>`,
		"pyproject.toml":   `[tool.poetry]`,
		"requirements.txt": "flask",
	}
	for name, content := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(content), 0644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}

	rule := NewEcosystemRule(dir)

	tests := []struct {
		path string
		want bool
	}{
		{filepath.Join(dir, "package-lock.json"), true},
		{filepath.Join(dir, "yarn.lock"), true},
		{filepath.Join(dir, "pnpm-lock.yaml"), true},
		{filepath.Join(dir, "go.sum"), true},
		{filepath.Join(dir, "example.nuget.props"), true},
		{filepath.Join(dir, "example.nuget.targets"), true},
		{filepath.Join(dir, "poetry.lock"), true},
		{filepath.Join(dir, "composer.lock"), false},
	}

	for _, tt := range tests {
		t.Run(filepath.Base(tt.path), func(t *testing.T) {
			if got := rule.Match(tt.path); got != tt.want {
				t.Fatalf("Match(%s) = %v, want %v", tt.path, got, tt.want)
			}
		})
	}
}
