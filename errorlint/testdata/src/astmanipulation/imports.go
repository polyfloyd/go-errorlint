package astmanipulation

import (
	"io"
	"os"
)

// Test missing errors import with simple comparison
func MissingImportSimple() {
	var err error
	var n int

	if n == 0 || (err != nil && err != io.EOF) { // want "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
		return
	}
}

// Test missing errors import with type assertion
func MissingImportTypeAssertion() {
	var err error

	if e, ok := err.(*os.PathError); ok { // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
		println("Path error:", e.Path)
	}
}