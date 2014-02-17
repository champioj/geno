package main

import (
	"errors"
	"go/ast"
	"go/parser"
	"path"
	"strings"
)

// generic type is the computed form of a typepath
type genericType struct {
	genImport string
	name      string
	expr      *ast.Expr
}

// splitDotNotation return split s in two part,  the import path and the type
// ex: container/list.List -> container/list, list.List
// ex: int -> "", int
func splitDotNotation(s string) (string, string, error) {
	dot := strings.Split(s, ".")
	base := path.Base(s)
	if len(dot) == 1 {
		return "", dot[0], nil
	} else if len(dot) == 2 {
		return dot[0], base, nil
	} else {
		return "", "", errors.New("path format not valid: " + s)
	}
}

// parseExpr parse all typePaths and return the computed type.
// typePath not parseable are put in errs
func parseExpr(fullPath typePaths) (types []genericType, errs []error) {
	genTypes := make([]genericType, 0)
	for _, path := range fullPath {
		importPath, typeName, errPath := splitDotNotation(path)
		if errPath != nil {
			errs = append(errs, errPath)
			continue
		}
		xpr, errExpr := parser.ParseExpr(typeName)
		if errExpr != nil {
			errs = append(errs, errExpr)
			continue
		}
		genTypes = append(genTypes, genericType{importPath, typeName, &xpr})
	}
	return genTypes, errs
}
