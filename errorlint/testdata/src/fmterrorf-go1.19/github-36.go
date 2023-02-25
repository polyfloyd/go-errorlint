package issues

import (
	"errors"
	"fmt"
)

func Single() error {
	err1 := errors.New("oops1")
	err2 := errors.New("oops2")
	err3 := errors.New("oops3")
	return fmt.Errorf("%w, %v, %v", err1, err2, err3)
}

func Multiple() error {
	err1 := errors.New("oops1")
	err2 := errors.New("oops2")
	err3 := errors.New("oops3")
	return fmt.Errorf("%w, %w, %w", err1, err2, err3) // want "only one %w verb is permitted per format string"
}
