package log

import (
	"log"
	"os"
)

var (
	debugMode bool
	logger    *log.Logger
)

func init() {
	logger = log.New(os.Stderr, "", log.Ltime)
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
