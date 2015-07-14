package main

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//writeToDisk write the new package in the right gopath/src directory.
// it will create the package folder if needed and overwrite files
func writeToDisk(pack *ast.Package, fset *token.FileSet, packagePath string) {
	fullPath := filepath.Join(build.Default.GOPATH, "src", packagePath)
	errorDir := os.Mkdir(fullPath, 0775)
	if errorDir != nil && !os.IsExist(errorDir) {
		fmt.Println(errorDir)
	}

	filePath := filepath.Join(fullPath, ".generated")
	err := ioutil.WriteFile(filePath, []byte{}, 0666)
	if err != nil {
		panic(err)
	}
	for k, f := range pack.Files {
		var buf bytes.Buffer
		printer.Fprint(&buf, fset, f)
		filePath := filepath.Join(fullPath, filepath.Base(k))

		err := ioutil.WriteFile(filePath, buf.Bytes(), 0666)
		if err != nil {
			panic(err)
		}
	}
}

// systemImport return true if the path is a root go package, false otherwise.
func systemImport(path string) bool {
	return strings.Contains(path, build.Default.GOROOT)
}

// parsePackage parse the comments of a package with the function parser.ParseDir.
// it return an error if the package is a root package, or if parser.ParseDir failed
func parsePackage(fset *token.FileSet, packageName string) (pkgs map[string]*ast.Package,
	goroot bool, first error) {
	importDir, err := build.Import(packageName, "", 0)
	if err != nil {
		return nil, false, err
	}

	packMap, errParse := parser.ParseDir(fset, importDir.Dir,
		nil, parser.ParseComments)
	return packMap, importDir.Goroot, errParse
}
