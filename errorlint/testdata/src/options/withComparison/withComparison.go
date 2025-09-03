package testdata

import (
	"errors"
	"fmt"
)

var ErrSentinel = errors.New("sentinel error")

func doSomething() error {
	return ErrSentinel
}

// This should be flagged when comparison checking is enabled
func CompareWithEquals() {
	err := doSomething()
	if err == ErrSentinel { // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		fmt.Println("sentinel error")
	}
}