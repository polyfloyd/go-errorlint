package errorlint

import (
	"log"
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"
)

func TestErrorsAs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(), "errorsas")
}

func TestErrorsIs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(), "errorsis")
}

func TestFmtErrorf(t *testing.T) {
	analyzer := NewAnalyzer()
	if err := analyzer.Flags.Set("errorf", "true"); err != nil {
		log.Fatal(err)
	}
	analysistest.RunWithSuggestedFixes(t, analysistest.TestData(), analyzer, "fmterrorf")
}

func TestFmtErrorfMultiple(t *testing.T) {
	analyzer := NewAnalyzer()
	if err := analyzer.Flags.Set("errorf", "true"); err != nil {
		log.Fatal(err)
	}
	if err := analyzer.Flags.Set("errorf-multi", "false"); err != nil {
		log.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), analyzer, "fmterrorf-go1.19")
}

func TestAllowedComparisons(t *testing.T) {
	analyzer := NewAnalyzer()
	analysistest.Run(t, analysistest.TestData(), analyzer, "allowed")
}

func TestOptions(t *testing.T) {
	testCases := []struct {
		desc    string
		opt     Option
		pattern string
	}{
		{
			desc: "WithAllowedErrors",
			opt: WithAllowedErrors([]AllowPair{
				{err: "io.EOF", fun: "example.com/pkg.Read"},
			}),
			pattern: "options/withAllowedErrors",
		},
		{
			desc: "WithAllowedWildcard",
			opt: WithAllowedWildcard([]AllowPair{
				{err: "example.com/pkg.ErrMagic", fun: "example.com/pkg.Magic"},
			}),
			pattern: "options/withAllowedWildcard",
		},
	}
	for _, tt := range testCases {
		t.Run(tt.desc, func(t *testing.T) {
			analyzer := NewAnalyzer(tt.opt)
			analysistest.Run(t, analysistest.TestData(), analyzer, tt.pattern)
		})
	}
}

func TestIssueRegressions(t *testing.T) {
	analyzer := NewAnalyzer()
	analysistest.Run(t, analysistest.TestData(), analyzer, "issues")
}
