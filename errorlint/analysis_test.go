package errorlint

import (
	"golang.org/x/tools/go/analysis/analysistest"
	"log"
	"testing"
)

func TestErrorsAs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(), "errorsas")
}

func TestErrorsIs(t *testing.T) {
	analysistest.Run(t, analysistest.TestData(), NewAnalyzer(), "errorsis")
}

func TestFmtErrorf(t *testing.T) {
	analyzer := NewAnalyzer()
	err := analyzer.Flags.Set("errorf", "true")
	if err != nil {
		log.Fatal(err)
	}
	analysistest.Run(t, analysistest.TestData(), analyzer, "fmterrorf")
}
