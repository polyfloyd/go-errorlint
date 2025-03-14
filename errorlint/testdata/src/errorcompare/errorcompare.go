package errorcompare

import (
	"errors"
	"fmt"
	"io"
)

var ErrSentinel = errors.New("sentinel error")

func doSomething() error {
	return ErrSentinel
}

func doSomethingElse() error {
	return fmt.Errorf("wrapped: %w", ErrSentinel)
}

func doAnotherThing() error {
	return io.EOF
}

// This should be flagged - direct comparison with ==
func CompareWithEquals() {
	err := doSomething()
	if err == ErrSentinel { // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
		fmt.Println("sentinel error")
	}
}

// This should be flagged - direct comparison with !=
func CompareWithNotEquals() {
	err := doSomething()
	if err != ErrSentinel { // want "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
		fmt.Println("not sentinel error")
	}
}

// This should be flagged - direct comparison with == in binary expression
func CompareInBinaryExpr() {
	err := doSomething()
	cond1 := err == ErrSentinel // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
	cond2 := err == io.EOF      // want "comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error"
	if cond1 || cond2 {
		fmt.Println("known error")
	}
}

// This should be flagged - direct comparison with != in binary expression
func CompareInBinaryExprNotEquals() {
	err := doSomething()
	cond1 := err != ErrSentinel // want "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
	cond2 := err != io.EOF      // want "comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error"
	if cond1 && cond2 {
		fmt.Println("unknown error")
	}
}

// This should NOT be flagged - using errors.Is
func CompareWithErrorsIs() {
	err := doSomethingElse()
	if errors.Is(err, ErrSentinel) {
		fmt.Println("sentinel error (wrapped)")
	}
}

// This should NOT be flagged - nil comparison
func CompareWithNil() {
	err := doSomething()
	if err == nil {
		fmt.Println("no error")
	}
}

// This should NOT be flagged - nil comparison with !=
func CompareWithNotNil() {
	err := doSomething()
	if err != nil {
		fmt.Println("has error")
	}
}

// This should be flagged - direct comparison in switch statement
func CompareInSwitch() {
	err := doSomething()
	switch err {
	case ErrSentinel: // want "switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors"
		fmt.Println("sentinel error")
	case io.EOF:
		fmt.Println("EOF error")
	default:
		fmt.Println("unknown error")
	}
}

// This should be flagged - direct comparison in switch with assignment
func CompareInSwitchWithAssignment() {
	err := doSomething()
	switch e := err; e {
	case ErrSentinel: // want "switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors"
		fmt.Println("sentinel error:", e)
	case io.EOF:
		fmt.Println("EOF error:", e)
	default:
		fmt.Println("unknown error:", e)
	}
}

// This should NOT be flagged - using proper errors.Is in each case
func CompareInSwitchWithErrorsIs() {
	err := doSomething()
	switch {
	case errors.Is(err, ErrSentinel):
		fmt.Println("sentinel error")
	case errors.Is(err, io.EOF):
		fmt.Println("EOF error")
	default:
		fmt.Println("unknown error")
	}
}

// This should NOT be flagged - switch on non-error value
func SwitchOnNonErrorValue() {
	code := 404
	switch code {
	case 200:
		fmt.Println("OK")
	case 404:
		fmt.Println("Not Found")
	default:
		fmt.Println("Unknown status")
	}
}
