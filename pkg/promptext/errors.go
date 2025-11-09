package promptext

import (
	"errors"
	"fmt"
)

// Sentinel errors for common failure cases.
// Use errors.Is() to check for these errors in your code.
var (
	// ErrInvalidDirectory is returned when the provided directory path is invalid or inaccessible.
	ErrInvalidDirectory = errors.New("invalid or inaccessible directory")

	// ErrNoFilesMatched is returned when no files match the specified criteria.
	ErrNoFilesMatched = errors.New("no files matched the specified criteria")

	// ErrTokenBudgetTooLow is returned when the token budget is too low to include any files.
	ErrTokenBudgetTooLow = errors.New("token budget too low to include any files")

	// ErrInvalidFormat is returned when an unsupported output format is requested.
	ErrInvalidFormat = errors.New("invalid or unsupported output format")
)

// DirectoryError wraps directory-related errors with additional context.
type DirectoryError struct {
	Path string
	Err  error
}

func (e *DirectoryError) Error() string {
	return fmt.Sprintf("directory error for '%s': %v", e.Path, e.Err)
}

func (e *DirectoryError) Unwrap() error {
	return e.Err
}

// FilterError wraps filtering-related errors with additional context.
type FilterError struct {
	Pattern string
	Err     error
}

func (e *FilterError) Error() string {
	return fmt.Sprintf("filter error for pattern '%s': %v", e.Pattern, e.Err)
}

func (e *FilterError) Unwrap() error {
	return e.Err
}

// FormatError wraps formatting-related errors with additional context.
type FormatError struct {
	Format string
	Err    error
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("format error for '%s': %v", e.Format, e.Err)
}

func (e *FormatError) Unwrap() error {
	return e.Err
}
