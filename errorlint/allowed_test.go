package errorlint

import (
	"testing"
)

func Test_isAllowedErrAndFunc(t *testing.T) {
	setDefaultAllowedErrors()

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

func Benchmark_isAllowedErrAndFunc(b *testing.B) {
	var benchCases = []struct {
		desc string
		err  string
		fun  string
	}{
		{
			desc: "empty",
			err:  "",
			fun:  "",
		},
		{
			desc: "short",
			err:  "x",
			fun:  "x",
		},
		{
			desc: "long not existed",
			// should pass strings.HasPrefix length check, 30 symbols here
			err: "xxxx_xxxx_yyyy_yyyy_zzzz_zzzz_",
			fun: "xxxx_xxxx_yyyy_yyyy_zzzz_zzzz_",
		},
		{
			desc: "existed, not wildcard",
			err:  "io.EOF",
			fun:  "(io.Reader).Read",
		},
		{
			desc: "existed, wildcard",
			err:  "golang.org/x/sys/unix.Exx",
			fun:  "golang.org/x/sys/unix.xxx",
		},
	}
	for _, bb := range benchCases {
		b.Run(bb.desc, func(b *testing.B) {
			err, fun := bb.err, bb.fun
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				isAllowedErrAndFunc(err, fun)
			}
		})
	}
}
