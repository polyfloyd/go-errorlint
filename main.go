package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"golang.org/x/tools/go/packages"

	"github.com/polyfloyd/go-errorlint/errorlint"
)

func main() {
	checkErrorf := flag.Bool("errorf", false, "Check whether fmt.Errorf uses the %w verb for formatting errors. See the readme for caveats")
	flag.Parse()

	cfg := &packages.Config{
		Mode: packages.NeedTypes | packages.NeedTypesInfo,
	}
	pkgs, err := packages.Load(cfg, flag.Args()...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load: %v\n", err)
		os.Exit(100)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(100)
	}

	lints := []errorlint.Lint{}
	for _, pkg := range pkgs {
		if *checkErrorf {
			l := errorlint.LintFmtErrorfCalls(pkg.Fset, *pkg.TypesInfo)
			lints = append(lints, l...)
		}
		l := errorlint.LintErrorComparisons(pkg.Fset, *pkg.TypesInfo)
		lints = append(lints, l...)
		l = errorlint.LintErrorTypeAssertions(pkg.Fset, *pkg.TypesInfo)
		lints = append(lints, l...)
	}
	sort.Sort(errorlint.ByPosition(lints))

	for _, lint := range lints {
		fmt.Fprintf(os.Stderr, "%s: %s\n", lint.Pos, lint.Message)
	}
	if len(lints) > 0 {
		os.Exit(1)
	}
}
