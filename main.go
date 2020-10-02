package main

import (
	"github.com/polyfloyd/go-errorlint/errorlint"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(errorlint.NewAnalyzer())
}
