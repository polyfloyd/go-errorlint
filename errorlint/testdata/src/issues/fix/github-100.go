package issues

import (
	"os"
	"syscall"
)

// Test case for https://github.com/polyfloyd/go-errorlint/issues/100
// This reproduces the conflicting edits issue where both type assertion
// and error comparison rules apply to the same line of code.
func ConflictingEdits() {
	var err error = &os.PathError{Err: syscall.ESRCH}
	
	// This line should trigger both:
	// 1. Type assertion linting (err.(*os.PathError))
	// 2. Error comparison linting (e.Err == syscall.ESRCH)
	// Leading to a combined fix instead of conflicting edits
	if e, ok := err.(*os.PathError); ok && e.Err == syscall.ESRCH { // want "type assertion and error comparison will fail on wrapped errors. Use errors.As and errors.Is to check for specific errors"
		println("Found PathError with ESRCH")
	}
	
	// Additional test cases with similar patterns
	if pathErr, ok := err.(*os.PathError); ok && pathErr.Err == syscall.ENOENT { // want "type assertion and error comparison will fail on wrapped errors. Use errors.As and errors.Is to check for specific errors"
		println("Found PathError with ENOENT")
	}
	
	// Test the exact pattern from the original issue #100
	// This is the actual problematic line that caused conflicting edits
	if err != nil {
		println("error occurred")
	} else if e, ok := err.(*os.PathError); ok && e.Err == syscall.ESRCH { // want "type assertion and error comparison will fail on wrapped errors. Use errors.As and errors.Is to check for specific errors"
		// If the process exits while reading its /proc/$PID/maps, the kernel will
		// return ESRCH. Handle it as if the process did not exist.
		println("Found PathError with ESRCH in else if - from original issue")
	}
}