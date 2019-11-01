package testdata

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
	var me MyError
	if errors.As(err, &me) {
		fmt.Println("MyError")
	}
}

func TypeAssertion() {
	err := doAnotherThing()
	_, ok := err.(*MyError) // bad
	if ok {
		fmt.Println("MyError")
	}
}

func TypeSwitch() {
	err := doAnotherThing()
	switch err.(type) { // bad
	case *MyError:
		fmt.Println("MyError")
	}
}

func TypeSwitchInline() {
	switch doAnotherThing().(type) { // bad
	case *MyError:
		fmt.Println("MyError")
	}
}

func TypeSwitchAssign() {
	err := doAnotherThing()
	switch t := err.(type) { // bad
	case *MyError:
		fmt.Println("MyError", t)
	}
}

func TypeSwitchAssignInline() {
	switch t := doAnotherThing().(type) { // bad
	case *MyError:
		fmt.Println("MyError", t)
	}
}
