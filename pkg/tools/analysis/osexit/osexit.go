// Package osexit defines an Analyzer for the presence of the os.Exit function in the main package.
package osexit

import (
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/analysis"
)

const doc = "Checks for the presence of the os.Exit function in the main package"

// Analyzer implements an analysis function.
var Analyzer = &analysis.Analyzer{
	Name: "osexit",
	Doc:  doc,
	Run:  run,
}

func run(pass *analysis.Pass) (interface{}, error) {
	for _, file := range pass.Files {
		if pass.Pkg.Name() != "main" {
			continue
		}

		if len(file.Scope.Objects) != 1 {
			continue
		}

		ast.Inspect(file, func(n ast.Node) bool {
			if fun, ok := n.(*ast.SelectorExpr); ok && fun.Sel.Name == "Exit" &&
				fmt.Sprintf("%v", fun.X) == "os" {
				pass.Report(analysis.Diagnostic{
					Pos:     fun.Pos(),
					Message: "Found call of os.Exit on main package",
				})
			}
			return true
		})
	}

	return nil, nil
}
