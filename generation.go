package main

import (
	"code.google.com/p/go.tools/astutil"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
)

var interfaceRegExp *regexp.Regexp

func init() {
	interfaceRegExp = regexp.MustCompile(`<gen:([0-9]+)>`)
}

// findInterfaceComment search for comment matching the regexp interfaceRegExp in
// the comment group.
func findInterfaceComment(commentGroup *ast.CommentGroup) (int, bool, error) {
	if commentGroup == nil {
		return 0, false, nil
	}
	if len(commentGroup.List) == 1 {
		response := interfaceRegExp.FindStringSubmatch(commentGroup.List[0].Text)
		if len(response) == 2 {
			fmt.Println(response)

			number, err := strconv.Atoi(response[1])
			if err == nil {
				return number, true, nil
			} else {
				return 0, false, errors.New(
					"The format should be: <gen:x> where x is an integer")
			}
		}
	}
	return 0, false, nil
}

func addImport(pack *ast.Package, fset *token.FileSet, name string) {
	for _, f := range pack.Files {
		astutil.AddImport(fset, f, name)
	}
}

// setNewName change the name of the package
func setNewName(pack *ast.Package, name string) {
	for _, f := range pack.Files {
		newIdent := ast.Ident{f.Package, name, nil}
		f.Name = &newIdent
	}
}

// isInterfaceAlias check if the declaration is an alias of an interface and
// return it if true
// ok: type data interface{}
// not ok : type data int
// not ok : type data interface{
//	           name() string
//          }
func isInterfaceAlias(spec ast.Spec) (*ast.TypeSpec, bool) {
	t, ok := spec.(*ast.TypeSpec)
	if ok {
		switch t.Type.(type) {
		case *ast.InterfaceType:
			return t, true
		case *ast.Ident:
			return t, true
		}
	}
	return t, false
}

// generate generate from the the base package, a copy of it with the name genPackageName
// The new package will be located at the base of the base package.
// All of its interface alias marked with interfaceRegExp will be replace with
// the types in genTypes
func generate(basePackage, genPackageName string, genTypes []genericType) {
	basePackageName := path.Base(basePackage)
	genPackagePath := filepath.Join(basePackage, genPackageName)
	fset := token.NewFileSet() // positions are relative to fset
	// check that the number of types given is enough
	packMap, root, err := parsePackage(fset, basePackage)
	fmt.Println("bp:", basePackage)

	if root {
		fmt.Println("root package are invalid target")
	}
	if err != nil {
		fmt.Println(err)
		return
	}
	pack, ok := packMap[basePackageName]
	if !ok {
		fmt.Println("There is no package ", basePackageName, " in ", basePackage)
		return
	}

	types := make(map[int]*ast.TypeSpec)

	for _, v := range pack.Files { // check if import is really used
		findGenericDecl(v, fset, types)
	}
	for k, v := range types {
		if k > len(genTypes) {
			fmt.Println(basePackage, " need more than ", len(genTypes), " types")
			return
		}

		for _, file := range pack.Files { // TODO:check if import is really used
			tchange := findType(file, v.Name.Name)
			for _, ident := range tchange {
				ident.Name = genTypes[k-1].name
			}
		}
	}
	for _, v := range genTypes {
		if v.genImport != "" {
			addImport(pack, fset, v.genImport)
		}
	}
	setNewName(pack, genPackageName)
	writeToDisk(pack, fset, genPackagePath)
}

func findGenericDecl(f *ast.File, fset *token.FileSet, types map[int]*ast.TypeSpec) {
	for _, s := range f.Decls {
		t, ok := s.(*ast.GenDecl)
		if ok {
			for _, v := range t.Specs {
				t, ok := isInterfaceAlias(v)
				if ok {
					fmt.Println(t, t.Comment)
					i, found, _ := findInterfaceComment(t.Comment)
					if found {
						types[i] = t
					}
				}
			}
		}
	}
}

func findType(f *ast.File, name string) (result []*ast.Ident) {
	result = make([]*ast.Ident, 0)
	ast.Inspect(f, func(n ast.Node) bool {
		switch x := n.(type) {
		case *ast.Ident:
			s := x.Name
			if x.Obj != nil && s == name && x.Obj.Name == name && x.Obj.Pos() != x.Pos() {
				result = append(result, x)
			}
		}
		return true
	})
	return result
}
