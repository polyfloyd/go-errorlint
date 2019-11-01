package testdata

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
	if err == ErrFoo { // bad
		fmt.Println("ErrFoo")
	}
}

func NotEqualOperator() {
	err := doThing()
	if err != ErrFoo { // bad
		fmt.Println("not ErrFoo")
	}
}

func EqualOperatorYoda() {
	err := doThing()
	if ErrFoo == err { // bad
		fmt.Println("ErrFoo")
	}
}

func NotEqualOperatorYoda() {
	err := doThing()
	if ErrFoo != err { // bad
		fmt.Println("not ErrFoo")
	}
}

func CompareSwitch() {
	err := doThing()
	switch err { // bad
	case ErrFoo:
		fmt.Println("ErrFoo")
	}
}

func CompareSwitchInline() {
	switch doThing() { // bad
	case ErrFoo:
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
