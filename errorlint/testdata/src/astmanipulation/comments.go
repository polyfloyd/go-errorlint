package astmanipulation

import (
	"os"
	"syscall"
)

// Test inline comment preservation
func InlineComments() {
	var err error

	if err == syscall.ENOENT { // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		// This block comment should be preserved
		println("File not found")
	}
}

// Test else-if comment preservation (main PR issue)
func ElseIfComments() {
	var err error

	if err != nil {
		println("error occurred")
	} else if e, ok := err.(*os.PathError); ok && e.Err == syscall.ESRCH { // want "type assertion and error comparison will fail on wrapped errors. Use errors.As and errors.Is to check for specific errors"
		// If the process exits while reading its /proc/$PID/maps, the kernel will
		// return ESRCH. Handle it as if the process did not exist.
		println("Process not found")
	}
}

// Test block comments
func BlockComments() {
	var err error

	/* Pre-condition block comment */
	if e, ok := err.(*os.PathError); ok { // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
		/* Inline block comment */ println("Path error:", e.Path)
	}
}
