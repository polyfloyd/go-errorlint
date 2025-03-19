package errorassert

import (
	"errors"
	"fmt"
)

type MyError struct{}

func (*MyError) Error() string {
	return "my custom error"
}

type AnotherError struct{}

func (*AnotherError) Error() string {
	return "another error"
}

// ValueError implements error with a value receiver, not a pointer receiver
type ValueError struct{}

func (ValueError) Error() string {
	return "value error"
}

func doSomething() error {
	return &MyError{}
}

func doSomethingWrapped() error {
	return fmt.Errorf("wrapped: %w", &MyError{})
}

func doSomethingValueError() error {
	return ValueError{}
}

func doSomethingValueWrapped() error {
	return fmt.Errorf("wrapped value: %w", ValueError{})
}

// This should be flagged - direct type assertion
func TypeAssertionDirect() {
	err := doSomething()
	me, ok := err.(*MyError) // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
	if ok {
		fmt.Println("got my error:", me)
	}
}

// This should be flagged - type assertion in a simple assignment
func TypeAssertionAssignment() {
	err := doSomething()
	myErr, ok := err.(*MyError) // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
	fmt.Println(myErr, ok)
}

// This should be flagged - type assertion in if statement
func TypeAssertionInIf() {
	err := doSomething()
	if me, ok := err.(*MyError); ok { // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
		fmt.Println("got my error:", me)
	}
}

// This should be flagged - value type assertion in if statement
func ValueTypeAssertionInIf() {
	err := doSomethingValueError()
	if ve, ok := err.(ValueError); ok { // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
		fmt.Println("got value error:", ve)
	}
}

// This should be flagged - value type assertion in assignment
func ValueTypeAssertionAssignment() {
	err := doSomethingValueError()
	ve, ok := err.(ValueError) // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
	fmt.Println(ve, ok)
}

// This should be flagged - direct value type assertion
func ValueTypeAssertionDirect() {
	err := doSomethingValueError()
	_ = err.(ValueError) // want "type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors"
}

// This should be flagged - type switch
func TypeSwitchStatement() {
	err := doSomething()
	switch err.(type) { // want "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *MyError:
		fmt.Println("my error")
	case *AnotherError:
		fmt.Println("another error")
	default:
		fmt.Println("unknown error")
	}
}

// This should be flagged - type switch with value type
func TypeSwitchWithValueType() {
	err := doSomethingValueError()
	switch err.(type) { // want "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case ValueError:
		fmt.Println("value error")
	case *MyError:
		fmt.Println("my error")
	default:
		fmt.Println("unknown error")
	}
}

// This should be flagged - type switch with assignment
func TypeSwitchWithAssignment() {
	err := doSomething()
	switch e := err.(type) { // want "type switch on error will fail on wrapped errors. Use errors.As to check for specific errors"
	case *MyError:
		fmt.Println("my error:", e)
	case *AnotherError:
		fmt.Println("another error:", e)
	default:
		fmt.Println("unknown error:", e)
	}
}

// This should NOT be flagged - using errors.As
func UsingErrorsAs() {
	err := doSomethingWrapped()
	var me *MyError
	if errors.As(err, &me) {
		fmt.Println("got my error:", me)
	}
}

// This should NOT be flagged - using errors.As with value type
func UsingErrorsAsWithValueType() {
	err := doSomethingValueWrapped()
	var ve ValueError
	if errors.As(err, &ve) {
		fmt.Println("got value error:", ve)
	}
}

// This should NOT be flagged - non-error type assertion
func NonErrorTypeAssertion() {
	var i interface{} = "hello"
	if s, ok := i.(string); ok {
		fmt.Println(s)
	}
}
