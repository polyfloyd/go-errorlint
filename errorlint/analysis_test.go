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

func TestIssueRegressions(t *testing.T) {
	analyzer := NewAnalyzer()
	analysistest.Run(t, analysistest.TestData(), analyzer, "issues")
}
