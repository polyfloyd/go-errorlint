package errorlint

import (
	"testing"
)

func Test_isAllowedErrAndFunc(t *testing.T) {
	testCases := []struct {
		desc   string
		fun    string
		err    string
		expect bool
	}{
		{
			desc:   "allowedErrors: (io.Reader).Read",
			err:    "io.EOF",
			fun:    "(io.Reader).Read",
			expect: true,
		},
		{
			desc:   "allowedErrorWildcards: golang.org/x/sys/unix.",
			err:    "golang.org/x/sys/unix.EINVAL",
			fun:    "golang.org/x/sys/unix.EpollCreate",
			expect: true,
		},
		{
			desc:   "errorlint.Reader Read",
			err:    "io.EOF",
			fun:    "(errorlint.Reader).Read",
			expect: false,
		},
		// {desc: "fail test", expect: true},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			result := isAllowedErrAndFunc(tt.err, tt.fun)
			if tt.expect != result {
				t.Errorf("Not equal:\n\texpect: %t,\n\tactual: %t\n", tt.expect, result)
			}
		})
	}
}
