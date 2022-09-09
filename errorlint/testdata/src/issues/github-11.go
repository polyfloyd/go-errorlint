package issues

import (
	"errors"
	"os"
	"syscall"
)

// Regression test for: https://github.com/polyfloyd/go-errorlint/issues/11

var ErrInvalidConfig = errors.New("invalid configuration")

// invalidConfig type is needed to make any error returned from Validator
// to be ErrInvalidConfig.
type invalidConfig struct {
	err error
}

func (e *invalidConfig) Is(target error) bool {
	return target == ErrInvalidConfig
}

type Errno uintptr

func (e Errno) Is(target error) bool {
	switch target {
	case os.ErrPermission:
		return e == Errno(syscall.EACCES) || e == Errno(syscall.EPERM)
	case os.ErrExist:
		return e == Errno(syscall.EEXIST) || e == Errno(syscall.ENOTEMPTY)
	case os.ErrNotExist:
		return e == Errno(syscall.ENOENT)
	}
	return false
}
