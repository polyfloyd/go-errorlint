package test

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
	return fmt.Errorf("error: %v", err)
}

func DoubleNonWrappingVerb() error {
	err := errors.New("oops")
	return fmt.Errorf("%v %v", err, err)
}

func MixedGoodAndBad() error {
	err := errors.New("oops")
	return fmt.Errorf("%v %w", err, err)
}

func NotAnError() error {
	err := "oops"
	return fmt.Errorf("%v", err)
}
