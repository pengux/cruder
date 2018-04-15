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

// parsePkg parses the directory or files for Go code and
// do type-checks on it. It will return the types.Package and
// directory location if successful
func parsePkg(sources []string) (*types.Package, string, error) {
	var (
		dir       string
		fileNames []string
	)
	if len(sources) == 1 && isDirectory(sources[0]) {
		dir = sources[0]
		pkg, err := build.Default.ImportDir(dir, 0)
		if err != nil {
			return nil, dir, fmt.Errorf("cannot process directory %s: %s", dir, err)
		}
		fileNames = append(fileNames, pkg.GoFiles...)
		fileNames = append(fileNames, pkg.CgoFiles...)
		fileNames = append(fileNames, pkg.SFiles...)
		fileNames = prefixDirectory(dir, fileNames)
	} else {
		dir = filepath.Dir(sources[0])
		fileNames = sources
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
