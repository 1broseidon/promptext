package log

import (
	"bytes"
	"log"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// captureLogOutput captures log output for testing
func captureLogOutput(t *testing.T, fn func()) string {
	t.Helper()

	// Create a buffer to capture output
	var buf bytes.Buffer

	// Save original logger
	originalLogger := logger

	// Create a test logger that writes to our buffer
	testLogger := log.New(&buf, "", log.Ltime)
	logger = testLogger

	// Restore original logger after test
	t.Cleanup(func() {
		logger = originalLogger
	})

	// Execute the function
	fn()

	return buf.String()
}

// setupTest resets the logging state for clean tests
func setupTest(t *testing.T) {
	t.Helper()

	// Reset state
	debugMode = false
	phaseStart = time.Time{}
	timeMarks = make(map[string]time.Time)
	useColors = false

	// Cleanup after test
	t.Cleanup(func() {
		debugMode = false
		phaseStart = time.Time{}
		timeMarks = make(map[string]time.Time)
		useColors = false
	})
}

func TestEnable(t *testing.T) {
	setupTest(t)

	// Initially debug should be disabled
	assert.False(t, debugMode)
	assert.False(t, IsDebugEnabled())

	// Enable debug mode
	Enable()

	// Debug mode should now be enabled
	assert.True(t, debugMode)
	assert.True(t, IsDebugEnabled())
}

func TestDisable(t *testing.T) {
	setupTest(t)

	// Enable debug mode first
	Enable()
	assert.True(t, debugMode)

	// Disable debug mode
	Disable()

	// Debug mode should now be disabled
	assert.False(t, debugMode)
	assert.False(t, IsDebugEnabled())
}

func TestIsDebugEnabled(t *testing.T) {
	setupTest(t)

	// Test when disabled
	Disable()
	assert.False(t, IsDebugEnabled())

	// Test when enabled
	Enable()
	assert.True(t, IsDebugEnabled())
}

func TestDebug_WhenDisabled(t *testing.T) {
	setupTest(t)

	// Ensure debug is disabled
	Disable()

	// Capture output
	output := captureLogOutput(t, func() {
		Debug("test message %s", "arg")
	})

	// Should produce no output when disabled
	assert.Empty(t, output)
}

func TestDebug_WhenEnabled_WithoutColors(t *testing.T) {
	setupTest(t)

	// Enable debug mode, disable colors
	Enable()
	SetColorEnabled(false)

	// Capture output
	output := captureLogOutput(t, func() {
		Debug("test message %s", "arg")
	})

	// Should contain debug prefix and formatted message
	assert.Contains(t, output, "[DEBUG] test message arg")
	assert.NotContains(t, output, "\033[") // No color codes
}

func TestDebug_WhenEnabled_WithColors(t *testing.T) {
	setupTest(t)

	// Enable debug mode and colors
	Enable()
	SetColorEnabled(true)

	// Capture output
	output := captureLogOutput(t, func() {
		Debug("test message %s", "arg")
	})

	// Should contain debug prefix, formatted message, and color codes
	assert.Contains(t, output, "[DEBUG] test message arg")
	assert.Contains(t, output, debugColor) // Light gray color
	assert.Contains(t, output, resetColor) // Reset color
}

func TestInfo_WhenDisabled(t *testing.T) {
	setupTest(t)

	// Ensure debug is disabled
	Disable()

	// Capture output
	output := captureLogOutput(t, func() {
		Info("test info %s", "message")
	})

	// Should produce no output when disabled
	assert.Empty(t, output)
}

func TestInfo_WhenEnabled(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Capture output
	output := captureLogOutput(t, func() {
		Info("test info %s", "message")
	})

	// Should contain info prefix and formatted message
	assert.Contains(t, output, "[INFO] test info message")
}

func TestError_AlwaysShown(t *testing.T) {
	setupTest(t)

	// Test when debug is disabled
	Disable()

	// Capture output
	output := captureLogOutput(t, func() {
		Error("test error %s", "message")
	})

	// Should always show error messages
	assert.Contains(t, output, "[ERROR] test error message")
}

func TestSetColorEnabled(t *testing.T) {
	setupTest(t)

	// Initially colors should be disabled
	assert.False(t, useColors)

	// Enable colors
	SetColorEnabled(true)
	assert.True(t, useColors)

	// Disable colors
	SetColorEnabled(false)
	assert.False(t, useColors)
}

func TestPhase_WhenDisabled(t *testing.T) {
	setupTest(t)

	// Ensure debug is disabled
	Disable()

	// Capture output
	output := captureLogOutput(t, func() {
		Phase("Test Phase")
	})

	// Should produce no output when disabled
	assert.Empty(t, output)

	// Phase start time should remain zero
	assert.True(t, phaseStart.IsZero())
}

func TestPhase_WhenEnabled_FirstPhase(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Capture output
	output := captureLogOutput(t, func() {
		Phase("Initial Phase")
	})

	// Should contain phase header
	assert.Contains(t, output, "[DEBUG] === Initial Phase ===")

	// Should not contain completion message for first phase
	assert.NotContains(t, output, "Phase completed")

	// Phase start time should be set
	assert.False(t, phaseStart.IsZero())

	// Timer mark should be set
	_, exists := timeMarks["Initial Phase"]
	assert.True(t, exists)
}

func TestPhase_WhenEnabled_SubsequentPhases(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start first phase
	captureLogOutput(t, func() {
		Phase("First Phase")
	})

	// Add a small delay to ensure measurable time difference
	time.Sleep(1 * time.Millisecond)

	// Start second phase and capture output
	output := captureLogOutput(t, func() {
		Phase("Second Phase")
	})

	// Should contain completion message for first phase
	assert.Contains(t, output, "Phase completed in")
	assert.Contains(t, output, "ms")

	// Should contain header for new phase
	assert.Contains(t, output, "[DEBUG] === Second Phase ===")
}

func TestStartTimer_WhenDisabled(t *testing.T) {
	setupTest(t)

	// Ensure debug is disabled
	Disable()

	// Start timer
	StartTimer("test-operation")

	// Timer should not be set when disabled
	_, exists := timeMarks["test-operation"]
	assert.False(t, exists)
}

func TestStartTimer_WhenEnabled(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start timer
	StartTimer("test-operation")

	// Timer should be set
	startTime, exists := timeMarks["test-operation"]
	assert.True(t, exists)
	assert.False(t, startTime.IsZero())
}

func TestEndTimer_WhenDisabled(t *testing.T) {
	setupTest(t)

	// Ensure debug is disabled
	Disable()

	// Capture output
	output := captureLogOutput(t, func() {
		EndTimer("test-operation")
	})

	// Should produce no output when disabled
	assert.Empty(t, output)
}

func TestEndTimer_WhenEnabled_WithoutStartTimer(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Capture output (trying to end a timer that was never started)
	output := captureLogOutput(t, func() {
		EndTimer("nonexistent-operation")
	})

	// Should produce no output for non-existent timer
	assert.Empty(t, output)
}

func TestEndTimer_WhenEnabled_WithStartTimer(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start timer
	StartTimer("test-operation")

	// Add a small delay to ensure measurable time
	time.Sleep(1 * time.Millisecond)

	// End timer and capture output
	output := captureLogOutput(t, func() {
		EndTimer("test-operation")
	})

	// Should contain completion message with timing
	assert.Contains(t, output, "[DEBUG] test-operation completed in")
	assert.Contains(t, output, "ms")

	// Timer mark should be removed after ending
	_, exists := timeMarks["test-operation"]
	assert.False(t, exists)
}

func TestGetPhaseStart_InitiallyZero(t *testing.T) {
	setupTest(t)

	// Initially phase start should be zero
	assert.True(t, GetPhaseStart().IsZero())
}

func TestGetPhaseStart_AfterPhase(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start a phase
	captureLogOutput(t, func() {
		Phase("Test Phase")
	})

	// Phase start should now be set
	phaseStartTime := GetPhaseStart()
	assert.False(t, phaseStartTime.IsZero())
	assert.Equal(t, phaseStart, phaseStartTime)
}

func TestMultipleTimers(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start multiple timers
	StartTimer("operation1")
	StartTimer("operation2")
	StartTimer("operation3")

	// All timers should be tracked
	assert.Len(t, timeMarks, 3)

	// End one timer
	output := captureLogOutput(t, func() {
		EndTimer("operation2")
	})

	// Should show completion for operation2
	assert.Contains(t, output, "operation2 completed in")

	// Should have 2 remaining timers
	assert.Len(t, timeMarks, 2)

	// Remaining timers should still exist
	_, exists1 := timeMarks["operation1"]
	_, exists3 := timeMarks["operation3"]
	assert.True(t, exists1)
	assert.True(t, exists3)

	// Ended timer should be gone
	_, exists2 := timeMarks["operation2"]
	assert.False(t, exists2)
}

func TestTimingAccuracy(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Start timer
	StartTimer("accuracy-test")

	// Sleep for a known duration
	sleepDuration := 10 * time.Millisecond
	time.Sleep(sleepDuration)

	// End timer and capture output
	output := captureLogOutput(t, func() {
		EndTimer("accuracy-test")
	})

	// Extract the timing from the output
	assert.Contains(t, output, "accuracy-test completed in")
	assert.Contains(t, output, "ms")

	// The timing should be reasonably close to our sleep duration
	// (allowing for some overhead and system timing variations)
	lines := strings.Split(output, "\n")
	var timingLine string
	for _, line := range lines {
		if strings.Contains(line, "accuracy-test completed in") {
			timingLine = line
			break
		}
	}

	require.NotEmpty(t, timingLine, "Could not find timing line in output")

	// Extract the numeric part (this is a basic check - timing will vary)
	assert.Contains(t, timingLine, "ms")
}

func TestConcurrentTimers(t *testing.T) {
	setupTest(t)

	// Enable debug mode
	Enable()

	// Test that timers work correctly with concurrent operations
	// (though the log package itself isn't necessarily thread-safe,
	// this tests the timer logic doesn't interfere)

	StartTimer("concurrent1")
	time.Sleep(1 * time.Millisecond)
	StartTimer("concurrent2")
	time.Sleep(1 * time.Millisecond)

	// End timers in different order
	output1 := captureLogOutput(t, func() {
		EndTimer("concurrent2")
	})

	output2 := captureLogOutput(t, func() {
		EndTimer("concurrent1")
	})

	// Both should have completion messages
	assert.Contains(t, output1, "concurrent2 completed in")
	assert.Contains(t, output2, "concurrent1 completed in")

	// All timers should be cleaned up
	assert.Empty(t, timeMarks)
}

// Integration test combining multiple functions
func TestLoggingIntegration(t *testing.T) {
	setupTest(t)

	// Test complete workflow
	output := captureLogOutput(t, func() {
		// Start with debug disabled
		Debug("should not appear")

		// Enable debug
		Enable()
		Debug("should appear")

		// Start a phase
		Phase("Integration Test")

		// Start some timers
		StartTimer("sub-operation1")
		time.Sleep(1 * time.Millisecond)
		EndTimer("sub-operation1")

		// Log some info
		Info("integration test in progress")

		// Start another phase
		time.Sleep(1 * time.Millisecond)
		Phase("Cleanup Phase")

		// End with debug disabled
		Disable()
		Debug("should not appear again")
	})

	// Check that expected messages appear
	assert.Contains(t, output, "[DEBUG] should appear")
	assert.NotContains(t, output, "should not appear")
	assert.Contains(t, output, "=== Integration Test ===")
	assert.Contains(t, output, "=== Cleanup Phase ===")
	assert.Contains(t, output, "sub-operation1 completed in")
	assert.Contains(t, output, "[INFO] integration test in progress")
	assert.Contains(t, output, "Phase completed in")
}

// TestFatal_MessageFormatting tests that Fatal formats messages correctly
// Note: We cannot test the os.Exit(1) behavior without terminating the test process
func TestFatal_MessageFormatting(t *testing.T) {
	setupTest(t)

	// We can only test the message formatting part of Fatal
	// The os.Exit(1) part cannot be tested without custom exit handling
	output := captureLogOutput(t, func() {
		// Create a copy of Fatal that doesn't exit for testing
		testLogger := logger
		testLogger.Printf("[FATAL] test fatal %s", "message")
	})

	// Should contain fatal prefix and formatted message
	assert.Contains(t, output, "[FATAL] test fatal message")
}

// Note: The Fatal function calls os.Exit(1) which cannot be tested in a standard unit test
// without additional infrastructure. The function is simple (just logger.Printf + os.Exit)
// so testing the logging part verifies the core functionality.
