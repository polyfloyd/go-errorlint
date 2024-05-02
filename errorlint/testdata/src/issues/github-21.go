package issues

// Regression test for: https://github.com/polyfloyd/go-errorlint/issues/21

import (
	"errors"
	"fmt"
)

func Test1() error {
	return fmt.Errorf("%[1]v %[1]v: %w", "value", errors.New("abc"))
}

func Test2() error {
	return fmt.Errorf("%[1]v: %[1]w", errors.New("abc"))
}
