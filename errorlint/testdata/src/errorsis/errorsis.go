package errorsis

import (
	"errors"
	"fmt"
)

var ErrFoo = errors.New("foo")

func doThing() error {
	return ErrFoo
}

func CompareGood() {
	err := doThing()
	if errors.Is(err, ErrFoo) {
		fmt.Println("ErrFoo")
	}
}

func CompareOperatorNilGood() {
	err := doThing()
	if err == nil {
		fmt.Println("nil")
	}
}

func CompareOperatorNotNilGood() {
	err := doThing()
	if err != nil {
		fmt.Println("nil")
	}
}

func CompareOperatorNilYodaGood() {
	err := doThing()
	if nil == err {
		fmt.Println("nil")
	}
}

func CompareOperatorNotNilYodaGood() {
	err := doThing()
	if nil != err {
		fmt.Println("nil")
	}
}

func EqualOperator() {
	err := doThing()
	if err == ErrFoo { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("ErrFoo")
	}
}

func NotEqualOperator() {
	err := doThing()
	if err != ErrFoo { // want `comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("not ErrFoo")
	}
}

func EqualOperatorYoda() {
	err := doThing()
	if ErrFoo == err { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("ErrFoo")
	}
}

func NotEqualOperatorYoda() {
	err := doThing()
	if ErrFoo != err { // want `comparing with != will fail on wrapped errors. Use errors.Is to check for a specific error`
		fmt.Println("not ErrFoo")
	}
}

func CompareSwitch() {
	err := doThing()
	switch err {
	case nil:
		fmt.Println("nil")
	case ErrFoo: // want `switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors`
		fmt.Println("ErrFoo")
	}
}

func CompareSwitchSafe() {
	err := doThing()
	switch err {
	case nil:
		fmt.Println("success")
	default:
		fmt.Println("failure")
	}
}

func CompareSwitchInline() {
	switch doThing() {
	case ErrFoo: // want `switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors`
		fmt.Println("ErrFoo")
	}
}

func CompareSwitchNonError() {
	s := "foo"
	switch s {
	case "bar":
		fmt.Println("bar")
	}
}

func CompareSwitchOnNothing() {
	switch {
	case true:
		fmt.Println("foo")
	}
}
