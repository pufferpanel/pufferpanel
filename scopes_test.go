package pufferpanel

import (
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
	"strings"
	"testing"
)

func TestAllKnownScopes(t *testing.T) {

	t.Run("all scopes added", func(t *testing.T) {
		fs := token.NewFileSet()
		f, err := parser.ParseFile(fs, "scopes.go", nil, 0)
		if err != nil {
			panic(err)
		}

		pkg := types.NewPackage("", "pufferpanel")
		info := &types.Info{}
		checker := types.NewChecker(&types.Config{}, fs, pkg, info)
		files := []*ast.File{f}
		err = checker.Files(files)
		if err != nil {
			panic(err)
		}

		allScopes := AllKnownScopes()

		scopeFile := files[0]
		for k, v := range scopeFile.Scope.Objects {
			if k == "Scope" || !strings.HasPrefix(k, "Scope") {
				continue
			}

			spec := v.Decl.(*ast.ValueSpec)
			callExpr := spec.Values[0].(*ast.CallExpr)
			arg := callExpr.Args[0].(*ast.BasicLit)
			scopeValue := Scope(strings.TrimPrefix(strings.TrimSuffix(arg.Value, "\""), "\""))

			matches := false

			for _, s := range allScopes {
				if s == scopeValue {
					matches = true
				}
			}

			if !matches {
				t.Errorf("Scope [%s] is not in all known scopes", k)
			}
		}
	})
}
