package main

import (
	"golang.org/x/tools/go/analysis/singlechecker"

	"github.com/polyfloyd/go-errorlint/errorlint"
)

func main() {
	singlechecker.Main(errorlint.NewAnalyzer())
}
