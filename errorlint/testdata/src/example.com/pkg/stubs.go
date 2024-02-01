package pkg

import (
	"errors"
	"io"
)

func Read(io.Reader) error {
	return io.EOF
}

var (
	ErrMagicOne = errors.New("magic")
)

func MagicOne() error {
	return ErrMagicOne
}
