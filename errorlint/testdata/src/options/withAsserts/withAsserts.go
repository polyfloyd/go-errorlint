package testdata

import (
	"fmt"
)

type MyError struct{}

func (*MyError) Error() string {
	return "my custom error"
}

func doSomething() error {
	return &MyError{}
}

// This should be flagged when assert checking is enabled
func TypeAssertionDirect() {
	err := doSomething()
	me, ok := err.(*MyError) // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
	if ok {
		fmt.Println("got my error:", me)
	}
}