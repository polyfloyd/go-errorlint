package issues

// Regression test for: https://github.com/polyfloyd/go-errorlint/issues/19

import (
	"errors"
	"io"
)

var errChecksum = errors.New("checksum error")

type checksumReader struct {
	rc   io.ReadCloser
	hash uint32
}

func (r *checksumReader) Read(b []byte) (n int, err error) {
	n, err = r.rc.Read(b)
	if err == nil {
		return
	}
	if r.hash != 123 {
		err = errChecksum
	}
	if r.hash != 123 {
		err = errChecksum
	}
	if err == errChecksum { // want `comparing with == will fail on wrapped errors. Use errors.Is to check for a specific error`
		return 0, err
	}
	return
}
