package main

import (
	"errors"
	"flag"
	"fmt"
	"path"
	"strings"
)

// typePaths contains a list a path fully describing a type
// example for a base type: int
// example for a type defined in a pacakge : container/list.List
type typePaths []string

func (ts *typePaths) String() string {
	return fmt.Sprint(*ts)
}

func (ts *typePaths) Set(value string) error {
	if len(*ts) > 0 {
		return errors.New("types flag already set")
	}
	for _, dt := range strings.Split(value, ",") {
		typeStr := strings.TrimSpace(dt)
		*ts = append(*ts, typeStr)
	}
	return nil
}

var packageFlag string
var typesFlag typePaths

var recursiveFlag string

// genericDef contains the string description of a package to be generated
type genericDef struct {
	packageImport string
	types         typePaths
}

func init() {
	flag.StringVar(&packageFlag, "package", "", "the package to use as template (use with type)")
	flag.Var(&typesFlag, "types", "comma-separated list of types (use with package)")

	flag.StringVar(&recursiveFlag, "recursive", "", "Recursively generate packages for the package and its import (use alone)")
}

func main() {
	var genDefs []genericDef
	flag.Parse()

	fmt.Println(packageFlag)
	fmt.Println(typesFlag)
	fmt.Println(recursiveFlag)

	if recursiveFlag != "" {
		genDefs = walkImport(recursiveFlag)
	} else if len(typesFlag) > 0 && packageFlag != "" {
		genDefs = append(genDefs, genericDef{packageFlag, typesFlag})
	} else {
		flag.Usage()
		return
	}

	for _, v := range genDefs {
		fmt.Println("Generating: ", v)
		typesExpr, errs := parseExpr(v.types)
		if errs != nil {
			fmt.Println(errs)
			continue
		}
		base, name := path.Split(v.packageImport)
		generate(base, name, typesExpr)
	}
}
