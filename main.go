package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/constant"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"regexp"
)

func main() {
	flag.Parse()
	sourceFiles := flag.Args()

	if len(sourceFiles) == 0 {
		log.Fatal("no source files")
	}

	lints, err := lint(sourceFiles...)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(100)
	}
	for _, lint := range lints {
		fmt.Fprintf(os.Stderr, "%s: %s\n", lint.Pos, lint.Message)
	}
	if len(lints) > 0 {
		os.Exit(1)
	}
}

type Lint struct {
	Message string
	Pos     token.Position
}

func lint(sourceFiles ...string) ([]Lint, error) {
	fset := token.NewFileSet()
	astFiles := make([]*ast.File, len(sourceFiles))
	for i, filename := range sourceFiles {
		f, err := parser.ParseFile(fset, filename, nil, 0)
		if err != nil {
			return nil, err
		}
		astFiles[i] = f
	}

	info := types.Info{
		Types: make(map[ast.Expr]types.TypeAndValue),
		Defs:  make(map[*ast.Ident]types.Object),
		Uses:  make(map[*ast.Ident]types.Object),
	}
	conf := types.Config{
		Importer: importer.Default(),
	}
	_, err := conf.Check("test", fset, astFiles, &info)
	if err != nil {
		return nil, err
	}

	lints := []Lint{}
	for expr, t := range info.Types {
		// Search for error expressions that are the result of fmt.Errorf
		// invocations.
		if t.Type.String() != "error" {
			continue
		}
		call, ok := isFmtErrorfCallExpr(info, expr)
		if !ok {
			continue
		}

		// Find all % fields in the format string.
		formatVerbs, ok := printfFormatStringVerbs(info, call)
		if !ok {
			continue
		}
		// For all arguments that are errors, check whether the wrapping verb
		// is used.
		for i, arg := range call.Args[1:] {
			if info.Types[arg].Type.String() != "error" {
				continue
			}
			if len(formatVerbs) >= i && formatVerbs[i] != "%w" {
				lints = append(lints, Lint{
					Message: "non-wrapping format verb for fmt.Errorf. Use `%w` to format errors",
					Pos:     fset.Position(expr.Pos()),
				})
			}
		}
	}
	return lints, nil
}

func printfFormatStringVerbs(info types.Info, call *ast.CallExpr) ([]string, bool) {
	if len(call.Args) <= 1 {
		return nil, false
	}
	strLit, ok := call.Args[0].(*ast.BasicLit)
	if !ok {
		// Ignore format strings that are not literals.
		return nil, false
	}
	formatString := constant.StringVal(info.Types[strLit].Value)

	// Naive format string argument verb. This does not take modifiers such as
	// padding into account...
	re := regexp.MustCompile(`%[^%]`)
	return re.FindAllString(formatString, -1), true
}

func isFmtErrorfCallExpr(info types.Info, expr ast.Expr) (*ast.CallExpr, bool) {
	call, ok := expr.(*ast.CallExpr)
	if !ok {
		return nil, false
	}
	fn, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		// TODO: Support fmt.Errorf variable aliases?
		return nil, false
	}
	obj := info.Uses[fn.Sel]

	pkg := obj.Pkg()
	if pkg != nil && pkg.Name() == "fmt" && obj.Name() == "Errorf" {
		return call, true
	}
	return nil, false
}
