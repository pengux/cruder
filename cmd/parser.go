package cmd

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// getPkgAndType parses src (directories/files) and return the *types.Package,
// *types.Struct for the passed in structName and the directory for the sources
func getPkgAndType(structName string, src ...string) (*types.Package, *types.Struct, string, error) {
	pkg, dir, err := parsePkg(src...)
	if err != nil {
		return nil, nil, "", fmt.Errorf("parsing package from provided sources: %s", err)
	}

	// Check that struct exists in package
	o := pkg.Scope().Lookup(structName)
	if o == nil {
		return nil, nil, "", fmt.Errorf("the struct %s doesn't seem to exists in package %s", structName, pkg.Name())
	}
	// Check that it really is of type struct
	t, ok := o.Type().Underlying().(*types.Struct)
	if !ok {
		return nil, nil, "", fmt.Errorf("the type %s is not a struct", structName)
	}

	return pkg, t, dir, nil
}

// parsePkg parses the directory or files for Go code and
// do type-checks on it. It will return the types.Package and
// directory location if successful
func parsePkg(src ...string) (*types.Package, string, error) {
	var (
		dir       string
		fileNames []string
	)
	if len(src) == 1 && isDirectory(src[0]) {
		dir = src[0]
		pkg, err := build.Default.ImportDir(dir, 0)
		if err != nil {
			return nil, dir, fmt.Errorf("cannot process directory %s: %s", dir, err)
		}
		fileNames = append(fileNames, pkg.GoFiles...)
		fileNames = append(fileNames, pkg.CgoFiles...)
		fileNames = append(fileNames, pkg.SFiles...)
		fileNames = prefixDirectory(dir, fileNames)
	} else {
		dir = filepath.Dir(src[0])
		fileNames = src
	}

	var (
		astFiles []*ast.File
		pkgName  string
	)
	fset := token.NewFileSet()
	for _, fileName := range fileNames {
		if !strings.HasSuffix(fileName, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fset, fileName, nil, 0)
		if err != nil {
			return nil, dir, fmt.Errorf("parsing package: %s: %s", fileName, err)
		}
		astFiles = append(astFiles, parsedFile)
	}
	if len(astFiles) == 0 {
		return nil, dir, fmt.Errorf("%s: no buildable Go files", dir)
	}

	pkgName = astFiles[0].Name.Name
	conf := types.Config{Importer: importer.Default()}
	pkg, err := conf.Check(pkgName, fset, astFiles, nil)
	if err != nil {
		return nil, dir, fmt.Errorf("type-checking package %s: %s", pkgName, err)
	}

	return pkg, dir, nil
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// prefixDirectory joins the directoy with each filename
func prefixDirectory(directory string, fileNames []string) []string {
	if directory == "." {
		return fileNames
	}
	ret := make([]string, len(fileNames))
	for i, fileName := range fileNames {
		ret[i] = filepath.Join(directory, fileName)
	}
	return ret
}
