package promptext

import (
	"github.com/1broseidon/promptext/internal/processor"
)

// Run executes the promptext tool with the given configuration
func Run(dirPath string, extension string, exclude string, noCopy bool, infoOnly bool, verbose bool) error {
	return processor.Run(dirPath, extension, exclude, noCopy, infoOnly, verbose)
}
