package main

import (
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"regexp"
)

type Lint struct {
	Message string
	Pos     token.Position
}

type ByPosition []Lint

func (l ByPosition) Len() int           { return len(l) }
func (l ByPosition) Less(i, j int) bool { return l[i].Pos.Offset < l[j].Pos.Offset }
func (l ByPosition) Swap(i, j int)      { l[i], l[j] = l[j], l[i] }

func lintFmtErrorfCalls(fset *token.FileSet, info types.Info) []Lint {
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
					Pos:     fset.Position(arg.Pos()),
				})
			}
		}
	}
	return lints
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

func lintErrorComparisons(fset *token.FileSet, info types.Info) []Lint {
	lints := []Lint{}

	for expr := range info.Types {
		// Find == and != operations.
		binExpr, ok := expr.(*ast.BinaryExpr)
		if !ok {
			continue
		}
		if binExpr.Op != token.EQL && binExpr.Op != token.NEQ {
			continue
		}
		// Comparing errors with nil is okay.
		if isNilComparison(binExpr) {
			continue
		}
		// Find comparisons of which one side is a of type error.
		if !isErrorComparison(info, binExpr) {
			continue
		}

		lints = append(lints, Lint{
			Message: fmt.Sprintf("comparing with %s will fail on wrapped errors. Use errors.Is to check for a specific error", binExpr.Op),
			Pos:     fset.Position(binExpr.Pos()),
		})
	}

	for scope := range info.Scopes {
		// Find value switch blocks.
		switchStmt, ok := scope.(*ast.SwitchStmt)
		if !ok {
			continue
		}
		// Check whether the switch operates on an error type.
		if switchStmt.Tag == nil {
			continue
		}
		tagType := info.Types[switchStmt.Tag]
		if tagType.Type.String() != "error" {
			continue
		}

		lints = append(lints, Lint{
			Message: "switch on an error will fail on wrapped errors. Use errors.Is to check for specific errors",
			Pos:     fset.Position(switchStmt.Pos()),
		})
	}

	return lints
}

func isNilComparison(binExpr *ast.BinaryExpr) bool {
	if ident, ok := binExpr.X.(*ast.Ident); ok && ident.Name == "nil" {
		return true
	}
	if ident, ok := binExpr.Y.(*ast.Ident); ok && ident.Name == "nil" {
		return true
	}
	return false
}

func isErrorComparison(info types.Info, binExpr *ast.BinaryExpr) bool {
	tx := info.Types[binExpr.X]
	ty := info.Types[binExpr.Y]
	return tx.Type.String() == "error" || ty.Type.String() == "error"
}
