package errorlint

import (
	"sort"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestLintFmtErrorfCalls(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, "../testdata/fmterrorf.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := LintFmtErrorfCalls(pkg.Fset, *pkg.TypesInfo)
	sort.Sort(ByPosition(lints))

	expectPositions := []struct {
		Line   int
		Column int
	}{
		{Line: 15, Column: 33}, // NonWrappingVerb
		{Line: 20, Column: 29}, // DoubleNonWrappingVerb 1
		{Line: 20, Column: 34}, // DoubleNonWrappingVerb 2
		{Line: 25, Column: 29}, // MixedGoodAndBad
		{Line: 30, Column: 33}, // ErrorStringFormat
		{Line: 35, Column: 33}, // ErrorStringFormatCustomError
	}
	for i, exp := range expectPositions {
		l := lints[i]
		if exp.Line != l.Pos.Line {
			t.Errorf("Unexpected line at index %d: exp %v, got %v", i, exp.Line, l.Pos.Line)
		}
		if exp.Column != l.Pos.Column {
			t.Errorf("Unexpected column at index %d: exp %v, got %v", i, exp.Column, l.Pos.Column)
		}
	}
}

func TestLintErrorComparisons(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, "../testdata/errorsis.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := LintErrorComparisons(pkg.Fset, *pkg.TypesInfo)
	sort.Sort(ByPosition(lints))

	expectPositions := []struct {
		Line   int
		Column int
	}{
		{Line: 51, Column: 5}, // EqualOperator
		{Line: 58, Column: 5}, // NotEqualOperator
		{Line: 65, Column: 5}, // EqualOperatorYoda
		{Line: 72, Column: 5}, // NotEqualOperatorYoda
		{Line: 79, Column: 2}, // CompareSwitch
		{Line: 86, Column: 2}, // CompareSwitchInline
	}
	for i, exp := range expectPositions {
		l := lints[i]
		if exp.Line != l.Pos.Line {
			t.Errorf("Unexpected line at index %d: exp %v, got %v", i, exp.Line, l.Pos.Line)
		}
		if exp.Column != l.Pos.Column {
			t.Errorf("Unexpected column at index %d: exp %v, got %v", i, exp.Column, l.Pos.Column)
		}
	}
}

func TestLintErrorTypeAssertions(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, "../testdata/errorsas.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := LintErrorTypeAssertions(pkg.Fset, *pkg.TypesInfo)
	sort.Sort(ByPosition(lints))

	expectPositions := []struct {
		Line   int
		Column int
	}{
		{Line: 28, Column: 11}, // TypeAssertion
		{Line: 36, Column: 9},  // TypeSwitch
		{Line: 43, Column: 9},  // TypeSwitchInline
		{Line: 51, Column: 14}, // TypeSwitchAssign
		{Line: 58, Column: 14}, // TypeSwitchAssignInline
	}
	for i, exp := range expectPositions {
		l := lints[i]
		if exp.Line != l.Pos.Line {
			t.Errorf("Unexpected line at index %d: exp %v, got %v", i, exp.Line, l.Pos.Line)
		}
		if exp.Column != l.Pos.Column {
			t.Errorf("Unexpected column at index %d: exp %v, got %v", i, exp.Column, l.Pos.Column)
		}
	}
}
