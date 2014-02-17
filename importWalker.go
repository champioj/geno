package main

import (
	"fmt"
	"go/ast"
	"go/token"
	"path"
	"regexp"
	"strings"
)

var importRegExp *regexp.Regexp

func init() {
	importRegExp = regexp.MustCompile(`<gen:(.*)>`)
}

// findImportComment search for comment matching the regexp importRegExp in
// the comment group
func findImportComment(commentGroup *ast.CommentGroup) (typePaths, bool) {
	if commentGroup == nil {
		return nil, false
	}
	if len(commentGroup.List) == 1 {
		response := importRegExp.FindStringSubmatch(commentGroup.List[0].Text)
		if len(response) == 2 {
			return strings.Split(response[1], ","), response[1] != ""
		}
	}
	return nil, false
}

// walkImport walk import walk import recursively and return every genericDef
// it found. It does avoid root package
func walkImport(pack string) (packages []genericDef) {
	fset := token.NewFileSet()
	packs, root, err := parsePackage(fset, pack)
	if root {
		return packages
	}
	if err != nil {
		fmt.Println(err)
		return packages
	}

	for _, vp := range packs {
		for _, vf := range vp.Files {
			for _, vi := range vf.Imports {
				typesString, found := findImportComment(vi.Comment)
				importPath := strings.Replace(vi.Path.Value, "\"", "", 2)
				if found {
					packages = append(packages, genericDef{importPath,
						typesString})
					importPath, _ = path.Split(importPath)
				}
				packages = append(packages, walkImport(importPath)...)
			}
		}
	}
	return packages
}
