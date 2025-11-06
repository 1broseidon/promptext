package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

type fakeInitializer struct {
	runErr error
	called bool
}

func (f *fakeInitializer) Run() error {
	f.called = true
	return f.runErr
}

func newTestDeps() (cliDeps, *bytes.Buffer, *bytes.Buffer) {
	stdout := &bytes.Buffer{}
	stderr := &bytes.Buffer{}
	deps := cliDeps{
		stdout: stdout,
		stderr: stderr,
		usage:  func() {},
		checkForUpdate: func(string) (bool, string, error) {
			return false, "", nil
		},
		updater: func(string, bool) error {
			return nil
		},
		notifyUpdate: func(string) {},
		processorRun: func(string, string, string, bool, bool, bool, string, string, bool, bool, bool, bool, bool, string, int, bool) error {
			return nil
		},
		absPath: func(p string) (string, error) {
			return "/abs/" + p, nil
		},
	}
	return deps, stdout, stderr
}

func TestRunHelp(t *testing.T) {
	deps, _, _ := newTestDeps()
	usageCalled := 0
	deps.usage = func() {
		usageCalled++
	}
	deps.processorRun = func(string, string, string, bool, bool, bool, string, string, bool, bool, bool, bool, bool, string, int, bool) error {
		t.Fatalf("processor should not run when showing help")
		return nil
	}

	if code := run([]string{"--help"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if usageCalled != 1 {
		t.Fatalf("expected usage to be called once, got %d", usageCalled)
	}
}

func TestRunVersion(t *testing.T) {
	deps, stdout, _ := newTestDeps()
	deps.notifyUpdate = nil
	originalVersion, originalDate := version, date
	version = "test-version"
	date = "2024-01-01"
	t.Cleanup(func() {
		version = originalVersion
		date = originalDate
	})

	if code := run([]string{"--version"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if got := stdout.String(); got != "promptext version test-version (2024-01-01)\n" {
		t.Fatalf("unexpected version output: %q", got)
	}
}

func TestRunCheckUpdateSuccess(t *testing.T) {
	deps, stdout, _ := newTestDeps()
	deps.notifyUpdate = nil
	deps.checkForUpdate = func(current string) (bool, string, error) {
		if current != version {
			t.Fatalf("unexpected version: %s", current)
		}
		return true, "v2.0.0", nil
	}

	if code := run([]string{"--check-update"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if got := stdout.String(); got != "A new version is available: v2.0.0 (current: dev)\nRun 'promptext --update' to update to the latest version\n" {
		t.Fatalf("unexpected stdout: %q", got)
	}
}

func TestRunCheckUpdateError(t *testing.T) {
	deps, _, stderr := newTestDeps()
	deps.notifyUpdate = nil
	deps.checkForUpdate = func(string) (bool, string, error) {
		return false, "", errors.New("boom")
	}

	if code := run([]string{"--check-update"}, deps); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if got := stderr.String(); got != "Error checking for updates: boom\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestRunUpdateError(t *testing.T) {
	deps, _, stderr := newTestDeps()
	deps.notifyUpdate = nil
	deps.updater = func(string, bool) error {
		return errors.New("update failed")
	}

	if code := run([]string{"--update"}, deps); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if got := stderr.String(); got != "Error updating: update failed\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestRunInitSuccess(t *testing.T) {
	deps, _, _ := newTestDeps()
	fakeInit := &fakeInitializer{}
	deps.newInitializer = func(root string, force bool, quiet bool) initializerRunner {
		if root != "/abs/project" {
			t.Fatalf("unexpected root: %s", root)
		}
		if !force {
			t.Fatalf("expected force to be true")
		}
		if quiet {
			t.Fatalf("expected quiet to be false")
		}
		return fakeInit
	}

	if code := run([]string{"--init", "--force", "-d", "project"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !fakeInit.called {
		t.Fatalf("expected initializer to be invoked")
	}
}

func TestRunInitError(t *testing.T) {
	deps, _, stderr := newTestDeps()
	fakeInit := &fakeInitializer{runErr: errors.New("init failed")}
	deps.newInitializer = func(string, bool, bool) initializerRunner {
		return fakeInit
	}

	if code := run([]string{"--init"}, deps); code != 1 {
		t.Fatalf("expected exit code 1, got %d", code)
	}
	if got := stderr.String(); got != "Error initializing config: init failed\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}
}

func TestRunFormatWarning(t *testing.T) {
	deps, _, stderr := newTestDeps()
	formatArg := ""
	deps.processorRun = func(_ string, _ string, _ string, _ bool, _ bool, _ bool, outputFormat string, _ string, _ bool, _ bool, _ bool, _ bool, _ bool, _ string, _ int, _ bool) error {
		formatArg = outputFormat
		return nil
	}

	if code := run([]string{"--format", "jsonl", "--output", "context.ptx"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if got := stderr.String(); got != "⚠️  Warning: format flag 'jsonl' conflicts with output extension '.ptx' - using 'jsonl' (flag takes precedence)\n" {
		t.Fatalf("unexpected stderr: %q", got)
	}
	if formatArg != "jsonl" {
		t.Fatalf("expected processor to receive explicit format, got %s", formatArg)
	}
}

func TestRunFormatAutoDetection(t *testing.T) {
	deps, _, _ := newTestDeps()
	var formatArg string
	deps.processorRun = func(_ string, _ string, _ string, _ bool, _ bool, _ bool, outputFormat string, _ string, _ bool, _ bool, _ bool, _ bool, _ bool, _ string, _ int, _ bool) error {
		formatArg = outputFormat
		return nil
	}

	if code := run([]string{"--output", "context.md"}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if formatArg != "markdown" {
		t.Fatalf("expected markdown format, got %s", formatArg)
	}
}

func TestRunProcessorInvocation(t *testing.T) {
	deps, _, _ := newTestDeps()
	called := false
	deps.processorRun = func(dir string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool, outputFormat string, outFile string, debug bool, gitignore bool, useDefaultRules bool, dryRun bool, quiet bool, relevance string, maxTokens int, explainSelection bool) error {
		called = true
		if dir != "./other" {
			t.Fatalf("unexpected dir: %s", dir)
		}
		if extension != ".go" {
			t.Fatalf("unexpected extension: %s", extension)
		}
		if !noCopy {
			t.Fatalf("expected noCopy true")
		}
		if !infoOnly {
			t.Fatalf("expected infoOnly true")
		}
		if !verbose {
			t.Fatalf("expected verbose true")
		}
		if outputFormat != "ptx" {
			t.Fatalf("unexpected format: %s", outputFormat)
		}
		if outFile != "out.ptx" {
			t.Fatalf("unexpected outFile: %s", outFile)
		}
		if !debug {
			t.Fatalf("expected debug true")
		}
		if gitignore {
			t.Fatalf("expected gitignore false")
		}
		if useDefaultRules {
			t.Fatalf("expected useDefaultRules false")
		}
		if !dryRun {
			t.Fatalf("expected dryRun true")
		}
		if quiet {
			t.Fatalf("expected quiet false")
		}
		if relevance != "foo" {
			t.Fatalf("unexpected relevance: %s", relevance)
		}
		if maxTokens != 123 {
			t.Fatalf("unexpected maxTokens: %d", maxTokens)
		}
		if !explainSelection {
			t.Fatalf("expected explainSelection true")
		}
		return nil
	}

	args := []string{"-d", "./other", "--extension", ".go", "--exclude", "vendor", "--no-copy", "--info", "--verbose", "--output", "out.ptx", "--debug", "--gitignore=false", "--use-default-rules=false", "--dry-run", "--relevant", "foo", "--max-tokens", "123", "--explain-selection"}
	if code := run(args, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	if !called {
		t.Fatalf("processor was not invoked")
	}
}

func TestRunNotifiesUpdate(t *testing.T) {
	deps, _, _ := newTestDeps()
	var wg sync.WaitGroup
	wg.Add(1)
	called := 0
	deps.notifyUpdate = func(string) {
		called++
		wg.Done()
	}

	if code := run([]string{"--directory", "."}, deps); code != 0 {
		t.Fatalf("expected exit code 0, got %d", code)
	}
	wg.Wait()
	if called != 1 {
		t.Fatalf("expected notifier to be called once, got %d", called)
	}
}

func TestRunParseError(t *testing.T) {
	deps, _, _ := newTestDeps()
	if code := run([]string{"--unknown"}, deps); code != 2 {
		t.Fatalf("expected exit code 2 for parse error, got %d", code)
	}
}

func TestRunInitializesNilDependencies(t *testing.T) {
	deps := cliDeps{
		processorRun: func(string, string, string, bool, bool, bool, string, string, bool, bool, bool, bool, bool, string, int, bool) error {
			t.Fatalf("processor should not execute in help mode")
			return nil
		},
		newInitializer: func(string, bool, bool) initializerRunner {
			return &fakeInitializer{}
		},
		absPath: func(p string) (string, error) { return p, nil },
	}
	if code := run([]string{"--help"}, deps); code != 0 {
		t.Fatalf("expected success exit code, got %d", code)
	}
}

func TestRunPropagatesProcessorError(t *testing.T) {
	deps, _, stderr := newTestDeps()
	deps.processorRun = func(string, string, string, bool, bool, bool, string, string, bool, bool, bool, bool, bool, string, int, bool) error {
		return errors.New("boom")
	}

	if code := run([]string{}, deps); code != 1 {
		t.Fatalf("expected failure exit code, got %d", code)
	}
	if !strings.Contains(stderr.String(), "boom") {
		t.Fatalf("expected processor error in stderr, got %q", stderr.String())
	}
}

func TestCustomUsageWithWriter(t *testing.T) {
	var buf bytes.Buffer
	customUsageWithWriter(&buf)
	usage := buf.String()
	if !strings.Contains(usage, "USAGE:") {
		t.Fatalf("expected usage to mention USAGE section")
	}
	if !strings.Contains(usage, "UPDATE OPTIONS") {
		t.Fatalf("expected usage to mention update options")
	}
}

func TestCustomUsageWritesToStdout(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("pipe: %v", err)
	}
	os.Stdout = w
	customUsage()
	w.Close()
	os.Stdout = orig

	data, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected customUsage to write output")
	}
}

func TestDefaultCLIDepsProvidesConcreteDeps(t *testing.T) {
	deps := defaultCLIDeps()
	if deps.stdout == nil || deps.stderr == nil {
		t.Fatalf("expected default deps to configure stdio")
	}
	if deps.usage == nil {
		t.Fatalf("expected usage function to be set")
	}
	if deps.checkForUpdate == nil || deps.updater == nil {
		t.Fatalf("expected update functions to be set")
	}
	if deps.processorRun == nil {
		t.Fatalf("expected processor run function")
	}
	if deps.newInitializer == nil {
		t.Fatalf("expected initializer factory")
	}
	if init := deps.newInitializer(".", false, true); init == nil {
		t.Fatalf("expected initializer factory to return non-nil")
	}
}
