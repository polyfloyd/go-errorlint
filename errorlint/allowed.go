package errorlint

import (
	"fmt"
	"go/ast"
	"go/types"
)

var allowedErrors = []struct {
	err string
	fun string
}{
	// pkg/archive/tar
	{err: "io.EOF", fun: "(*tar.Reader).Next"},
	{err: "io.EOF", fun: "(*tar.Reader).Read"},
	// pkg/bufio
	{err: "io.EOF", fun: "(*bufio.Reader).Read"},
	{err: "io.EOF", fun: "(*bufio.Reader).ReadByte"},
	{err: "io.EOF", fun: "(*bufio.Reader).ReadBytes"},
	{err: "io.EOF", fun: "(*bufio.Reader).ReadLine"},
	{err: "io.EOF", fun: "(*bufio.Reader).ReadSlice"},
	{err: "io.EOF", fun: "(*bufio.Reader).ReadString"},
	{err: "io.EOF", fun: "(*bufio.Scanner).Scan"},
	// pkg/bytes
	{err: "io.EOF", fun: "(*bytes.Buffer).Read"},
	{err: "io.EOF", fun: "(*bytes.Buffer).ReadByte"},
	{err: "io.EOF", fun: "(*bytes.Buffer).ReadBytes"},
	{err: "io.EOF", fun: "(*bytes.Buffer).ReadRune"},
	{err: "io.EOF", fun: "(*bytes.Buffer).ReadString"},
	{err: "io.EOF", fun: "(*bytes.Reader).Read"},
	{err: "io.EOF", fun: "(*bytes.Reader).ReadAt"},
	{err: "io.EOF", fun: "(*bytes.Reader).ReadByte"},
	{err: "io.EOF", fun: "(*bytes.Reader).ReadRune"},
	{err: "io.EOF", fun: "(*bytes.Reader).ReadString"},
	// pkg/database/sql
	{err: "sql.ErrNoRows", fun: "(*database/sql.Row).Scan"},
	// pkg/io
	{err: "io.EOF", fun: "(io.Reader).Read"},
	{err: "io.ErrClosedPipe", fun: "(*io.PipeWriter).Write"},
	{err: "io.ErrShortBuffer", fun: "io.ReadAtLeast"},
	{err: "io.ErrUnexpectedEOF", fun: "io.ReadAtLeast"},
	{err: "io.ErrUnexpectedEOF", fun: "io.ReadFull"},
	// pkg/net/http
	{err: "http.ErrServerClosed", fun: "(*net/http.Server).ListenAndServe"},
	{err: "http.ErrServerClosed", fun: "(*net/http.Server).ListenAndServeTLS"},
	{err: "http.ErrServerClosed", fun: "(*net/http.Server).Serve"},
	{err: "http.ErrServerClosed", fun: "(*net/http.Server).ServeTLS"},
	{err: "http.ErrServerClosed", fun: "http.ListenAndServe"},
	{err: "http.ErrServerClosed", fun: "http.ListenAndServeTLS"},
	{err: "http.ErrServerClosed", fun: "http.Serve"},
	{err: "http.ErrServerClosed", fun: "http.ServeTLS"},
	// pkg/os
	{err: "io.EOF", fun: "(*os.File).Read"},
	{err: "io.EOF", fun: "(*os.File).ReadAt"},
	{err: "io.EOF", fun: "(*os.File).ReadDir"},
	{err: "io.EOF", fun: "(*os.File).Readdir"},
	{err: "io.EOF", fun: "(*os.File).Readdirnames"},
	// pkg/strings
	{err: "io.EOF", fun: "(*strings.Reader).Read"},
	{err: "io.EOF", fun: "(*strings.Reader).ReadAt"},
	{err: "io.EOF", fun: "(*strings.Reader).ReadByte"},
	{err: "io.EOF", fun: "(*strings.Reader).ReadRune"},
}

func isAllowedErrAndFunc(err, fun string) bool {
	for _, allow := range allowedErrors {
		if allow.fun == fun && allow.err == err {
			return true
		}
	}
	return false
}

func isAllowedErrorComparison(info types.Info, binExpr *ast.BinaryExpr) bool {
	var errName string // `<package>.<name>`, e.g. `io.EOF`
	var callExprs []*ast.CallExpr

	// Figure out which half of the expression is the returned error and which
	// half is the presumed error declaration.
	for _, expr := range []ast.Expr{binExpr.X, binExpr.Y} {
		switch t := expr.(type) {
		case *ast.SelectorExpr:
			// A selector which we assume refers to a staticaly declared error
			// in a package.
			errName = selectorToString(t)
		case *ast.Ident:
			// Identifier, most likely to be the `err` variable or whatever
			// produces it.
			callExprs = assigningCallExprs(info, t)
		case *ast.CallExpr:
			callExprs = append(callExprs, t)
		}
	}

	// Unimplemented or not sure, disallow the expression.
	if errName == "" || len(callExprs) == 0 {
		return false
	}

	// Map call expressions to the function name format of the allow list.
	functionNames := make([]string, len(callExprs))
	for i, callExpr := range callExprs {
		functionSelector, ok := callExpr.Fun.(*ast.SelectorExpr)
		if !ok {
			// If the function is not a selector it is not an Std function that is
			// allowed.
			return false
		}
		if sel, ok := info.Selections[functionSelector]; ok {
			functionNames[i] = fmt.Sprintf("(%s).%s", sel.Recv(), sel.Obj().Name())
		} else {
			// If there is no selection, assume it is a package.
			functionNames[i] = selectorToString(callExpr.Fun.(*ast.SelectorExpr))
		}
	}

	// All assignments done must be allowed.
	for _, funcName := range functionNames {
		if !isAllowedErrAndFunc(errName, funcName) {
			return false
		}
	}
	return true
}

// assigningCallExprs finds all *ast.CallExpr nodes that are part of an
// *ast.AssignStmt that assign to the subject identifier.
func assigningCallExprs(info types.Info, subject *ast.Ident) []*ast.CallExpr {
	if subject.Obj == nil {
		return nil
	}

	// - Find object from identifier
	// - Find other identifiers that reference the object
	// - Walk through identifier parents to find assignments
	// - Find call expressions for assignments

	// Find the object that the identifier points to. We need this to find
	// identifiers that reference it.
	sobj := info.ObjectOf(subject)

	// Find other identifiers that reference this same object. Make sure to
	// exclude the subject identifier as it will cause an infinite recursion
	// and is being used in a read operation anyway.
	identifiers := []*ast.Ident{}
	for node, obj := range info.Uses {
		if obj == sobj && subject.Pos() != node.Pos() {
			identifiers = append(identifiers, node)
		}
	}
	for node, obj := range info.Defs {
		if obj == sobj && subject.Pos() != node.Pos() {
			identifiers = append(identifiers, node)
		}
	}

	// Find the scope in which the subject is declared. We need this to search
	// for parent nodes.
	var scopeNode ast.Node
	for node, scope := range info.Scopes {
		if scope == sobj.Parent() {
			scopeNode = node
			break
		}
	}
	if scopeNode == nil {
		return nil // TODO: When does this happen?
	}

	// Function scopes are mapped to only the function type, which does not
	// contain the function body needed for finding identifier parent nodes.
	// If any, remap the function type to its body.
	if funcType, ok := scopeNode.(*ast.FuncType); ok {
	outer:
		for node := range info.Scopes {
			if file, ok := node.(*ast.File); ok {
				for _, decl := range file.Decls {
					if funcDecl, ok := decl.(*ast.FuncDecl); ok {
						if funcDecl.Type == funcType {
							scopeNode = funcDecl
							break outer
						}
					}
				}
			}
		}
	}

	// Find the identifiers in the scope, but record their parent nodes. It's
	// too bad there is no parent node mapping in go/ast itself, so we have to
	// inspect and use a stack.
	parentNodes := map[*ast.Ident]ast.Node{}
	stack := []ast.Node{scopeNode}
	ast.Inspect(scopeNode, func(n ast.Node) bool {
		for _, ident := range identifiers {
			if n == ident {
				parentNodes[ident] = stack[len(stack)-1]
				break
			}
		}

		if n == nil {
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, n)
		}
		return true
	})

	// Find call expressions for assignments.
	var callExprs []*ast.CallExpr
	for _, parent := range parentNodes {
		switch declT := parent.(type) {
		case *ast.AssignStmt:
			// The identifier is LHS of an assignment.
			assignment := declT

			assigningExpr := assignment.Rhs[0]
			// If the assignment is comprised of multiple expressions, find out
			// which LHS expression we should use by finding its index in the LHS.
			if len(assignment.Rhs) > 1 {
				for i, lhs := range assignment.Lhs {
					if subject.Name == lhs.(*ast.Ident).Name {
						assigningExpr = assignment.Rhs[i]
						break
					}
				}
			}

			switch assignT := assigningExpr.(type) {
			case *ast.CallExpr:
				// Found the function call.
				callExprs = append(callExprs, assignT)
			case *ast.Ident:
				// The subject was the result of assigning from another identifier.
				callExprs = append(callExprs, assigningCallExprs(info, assignT)...)
			default:
				// TODO: inconclusive?
			}
		}
	}
	return callExprs
}

func selectorToString(selExpr *ast.SelectorExpr) string {
	if ident, ok := selExpr.X.(*ast.Ident); ok {
		return ident.Name + "." + selExpr.Sel.Name
	}
	return ""
}
