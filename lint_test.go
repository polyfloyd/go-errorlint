package main

import (
	"sort"
	"testing"

	"golang.org/x/tools/go/packages"
)

func TestLintFmtErrorfCalls(t *testing.T) {
	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, "./testdata/fmterrorf.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := lintFmtErrorfCalls(pkg.Fset, *pkg.TypesInfo)
	sort.Sort(ByPosition(lints))

	expectPositions := []struct {
		Line   int
		Column int
	}{
		{Line: 15, Column: 33}, // NonWrappingVerb
		{Line: 20, Column: 29}, // DoubleNonWrappingVerb
		{Line: 20, Column: 34},
		{Line: 25, Column: 29}, // MixedGoodAndBad
	}
	for i, l := range lints {
		exp := expectPositions[i]
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
	pkgs, err := packages.Load(cfg, "./testdata/errorsis.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := lintErrorComparisons(pkg.Fset, *pkg.TypesInfo)
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
	for i, l := range lints {
		exp := expectPositions[i]
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
	pkgs, err := packages.Load(cfg, "./testdata/errorsas.go")
	if err != nil {
		t.Fatal(err)
	}

	pkg := pkgs[0]
	lints := lintErrorTypeAssertions(pkg.Fset, *pkg.TypesInfo)
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
	for i, l := range lints {
		exp := expectPositions[i]
		if exp.Line != l.Pos.Line {
			t.Errorf("Unexpected line at index %d: exp %v, got %v", i, exp.Line, l.Pos.Line)
		}
		if exp.Column != l.Pos.Column {
			t.Errorf("Unexpected column at index %d: exp %v, got %v", i, exp.Column, l.Pos.Column)
		}
	}
}
