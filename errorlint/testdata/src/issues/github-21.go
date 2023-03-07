package issues

// Regression test for: https://github.com/polyfloyd/go-errorlint/issues/21

import (
	"errors"
	"fmt"
)

func Test1() error {
	return fmt.Errorf("%[1]v %[1]v: %w", "value", errors.New("abc")) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}

func Test2() error {
	return fmt.Errorf("%[1]v: %[1]w", errors.New("abc")) // want "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors"
}
