package log

import (
	"log"
	"os"
	"time"
)

var (
	debugMode bool
	logger    *log.Logger
	phaseStart time.Time
)

func init() {
	logger = log.New(os.Stderr, "", log.Ltime)
}

// Phase starts a new logging phase with a header
func Phase(name string) {
	if debugMode {
		if !phaseStart.IsZero() {
			duration := time.Since(phaseStart)
			logger.Printf("[DEBUG] Phase completed in %.2fs\n", duration.Seconds())
		}
		logger.Printf("[DEBUG] %s%s%s\n", "=== ", name, " ===")
		phaseStart = time.Now()
	}
}

// Enable turns on debug logging
func Enable() {
	debugMode = true
}

// Disable turns off debug logging
func Disable() {
	debugMode = false
}

// Debug logs a debug message if debug mode is enabled
func Debug(format string, v ...interface{}) {
	if debugMode {
		logger.Printf("[DEBUG] "+format, v...)
	}
}

// Info logs an info message if debug mode is enabled
func Info(format string, v ...interface{}) {
	if debugMode {
		logger.Printf("[INFO] "+format, v...)
	}
}

// Error logs an error message (always shown)
func Error(format string, v ...interface{}) {
	logger.Printf("[ERROR] "+format, v...)
}

// Fatal logs an error message and exits
func Fatal(format string, v ...interface{}) {
	logger.Printf("[FATAL] "+format, v...)
	os.Exit(1)
}
