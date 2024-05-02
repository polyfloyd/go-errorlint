package errorlint_test

import (
	"testing"

	"github.com/polyfloyd/go-errorlint/errorlint"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestOption(t *testing.T) {
	testCases := []struct {
		desc    string
		opt     errorlint.Option
		pattern string
	}{
		{
			desc: "WithAllowedErrors",
			opt: errorlint.WithAllowedErrors([]errorlint.AllowPair{
				{Err: "io.EOF", Fun: "example.com/pkg.Read"},
			}),
			pattern: "options/withAllowedErrors",
		},
		{
			desc: "WithAllowedWildcard",
			opt: errorlint.WithAllowedWildcard([]errorlint.AllowPair{
				{Err: "example.com/pkg.ErrMagic", Fun: "example.com/pkg.Magic"},
			}),
			pattern: "options/withAllowedWildcard",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			analyzer := errorlint.NewAnalyzer(tt.opt)
			analysistest.Run(t, analysistest.TestData(), analyzer, tt.pattern)
		})
	}
}
