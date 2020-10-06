package fmterrorf

import (
	"errors"
	"fmt"
)

func Good() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %w", err)
}

func NonWrappingVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %v", err) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func DoubleNonWrappingVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("%v %v", err, err) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors" "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func MixedGoodAndBad() error {
	err := errors.New("oops")
	return fmt.Errorf("%v %w", err, err) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func ErrorStringFormat() error {
	err := errors.New("oops")
	return fmt.Errorf("error: %s", err.Error()) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func ErrorStringFormatCustomError() error {
	err := MyError{}
	return fmt.Errorf("error: %s", err.Error()) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func NotAnError() error {
	err := "oops"
	return fmt.Errorf("%v", err)
}

type MyError struct{}

func (MyError) Error() string {
	return "oops"
}
