package log

import (
	"log"
	"os"
	"time"
)

var (
	debugMode  bool
	logger     *log.Logger
	phaseStart time.Time
	timeMarks  map[string]time.Time
	headerColor = "\033[1;36m" // Cyan
	debugColor  = "\033[0;37m" // Light gray
	resetColor  = "\033[0m"
	useColors   = false
)

func init() {
	logger = log.New(os.Stderr, "", log.Ltime)
	timeMarks = make(map[string]time.Time)
}

// SetColorEnabled enables or disables color output
func SetColorEnabled(enabled bool) {
	useColors = enabled
}

// Phase starts a new logging phase with a header
func Phase(name string) {
	if debugMode {
		if !phaseStart.IsZero() {
			duration := time.Since(phaseStart)
			logger.Printf("[DEBUG] Phase completed in %.2fms\n", float64(duration.Microseconds())/1000.0)
		}
		logger.Printf("[DEBUG] %s%s%s\n", "=== ", name, " ===")
		phaseStart = time.Now()
		timeMarks[name] = phaseStart
	}
}

// StartTimer starts timing an operation
func StartTimer(operation string) {
	if debugMode {
		timeMarks[operation] = time.Now()
	}
}

// EndTimer ends timing an operation and logs the duration
func EndTimer(operation string) {
	if debugMode {
		if start, ok := timeMarks[operation]; ok {
			duration := time.Since(start)
			logger.Printf("[DEBUG] %s completed in %.2fms\n", operation, float64(duration.Microseconds())/1000.0)
			delete(timeMarks, operation)
		}
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

// Header prints a section header in debug mode
func Header(format string, v ...interface{}) {
	if debugMode {
		msg := fmt.Sprintf(format, v...)
		if useColors {
			logger.Printf("%s=== %s ===%s", headerColor, msg, resetColor)
		} else {
			logger.Printf("=== %s ===", msg)
		}
	}
}

// Debug logs a debug message if debug mode is enabled
func Debug(format string, v ...interface{}) {
	if debugMode {
		if useColors {
			logger.Printf("%s[DEBUG] %s%s", debugColor, fmt.Sprintf(format, v...), resetColor)
		} else {
			logger.Printf("[DEBUG] "+format, v...)
		}
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
