package main

import (
	"flag"
	"fmt"
	"os"
	"sort"

	"golang.org/x/tools/go/packages"
)

func main() {
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

	lints := []Lint{}
	for _, pkg := range pkgs {
		l := lintFmtErrorfCalls(pkg.Fset, *pkg.TypesInfo)
		lints = append(lints, l...)
		l = lintErrorComparisons(pkg.Fset, *pkg.TypesInfo)
		lints = append(lints, l...)
		l = lintErrorTypeAssertions(pkg.Fset, *pkg.TypesInfo)
		lints = append(lints, l...)
	}
	sort.Sort(ByPosition(lints))

	for _, lint := range lints {
		fmt.Fprintf(os.Stderr, "%s: %s\n", lint.Pos, lint.Message)
	}
	if len(lints) > 0 {
		os.Exit(1)
	}
}
