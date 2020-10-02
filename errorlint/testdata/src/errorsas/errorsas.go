package errorsas

import (
	"errors"
	"fmt"
)

type MyError struct{}

func (*MyError) Error() string {
	return "foo"
}

func doAnotherThing() error {
	return &MyError{}
}

func TypeCheckGood() {
	err := doAnotherThing()
	var me *MyError
	if errors.As(err, &me) {
		fmt.Println("MyError")
	}
}

func TypeAssertion() {
	err := doAnotherThing()
	_, ok := err.(*MyError) // want `type assertion on error will fail on wrapped errors. Use errors.As to check for specific errors`
	if ok {
		fmt.Println("MyError")
	}
}

func TypeSwitch() {
	err := doAnotherThing()
	switch err.(type) { // want `type switch on error will fail on wrapped errors. Use errors.As to check for specific errors`
	case *MyError:
		fmt.Println("MyError")
	}
}

func TypeSwitchInline() {
	switch doAnotherThing().(type) { // want `type switch on error will fail on wrapped errors. Use errors.As to check for specific errors`
	case *MyError:
		fmt.Println("MyError")
	}
}

func TypeSwitchAssign() {
	err := doAnotherThing()
	switch t := err.(type) { // want `type switch on error will fail on wrapped errors. Use errors.As to check for specific errors`
	case *MyError:
		fmt.Println("MyError", t)
	}
}

func TypeSwitchAssignInline() {
	switch t := doAnotherThing().(type) { // want `type switch on error will fail on wrapped errors. Use errors.As to check for specific errors`
	case *MyError:
		fmt.Println("MyError", t)
	}
}
