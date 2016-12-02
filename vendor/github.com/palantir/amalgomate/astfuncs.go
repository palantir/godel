// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"go/ast"
	"go/build"
	"go/parser"
	"go/printer"
	"go/token"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
	"golang.org/x/tools/go/ast/astutil"
)

const (
	amalgomatedPackage = "amalgomated"
	amalgomatedMain    = "AmalgomatedMain"
	internalDir        = "internal"
)

// repackage repackages the main package specified in the provided configuration and re-writes them into the provided
// output directory. The repackaged files are placed into a directory called "vendor" that is created in the provided
// directory. This function assumes and verifies that the provided "outputDir" is a directory that exists. The provided
// configuration is processed based on the natural ordering of the name of the commands. If multiple commands specify
// the same package, it is only processed for the first command.
func repackage(config Config, outputDir string) error {
	if outputDirInfo, err := os.Stat(outputDir); err != nil {
		return errors.Wrapf(err, "failed to stat output directory: %s", outputDir)
	} else if !outputDirInfo.IsDir() {
		return errors.Wrapf(err, "not a directory: %s", outputDir)
	}

	vendorDir := path.Join(outputDir, internalDir)
	// remove output directory if it already exists
	if err := os.RemoveAll(vendorDir); err != nil {
		return errors.Wrapf(err, "failed to remove directory: %s", vendorDir)
	}

	if err := os.Mkdir(vendorDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create %s directory at %s", internalDir, vendorDir)
	}

	processedPkgs := make(map[SrcPkg]bool, len(config.Pkgs))
	for _, currName := range sortedKeys(config.Pkgs) {
		currPkg := config.Pkgs[currName]

		// if multiple keys specify the exact same source package, only process once
		if processedPkgs[currPkg] {
			continue
		}

		mainPkg, err := build.Import(currPkg.MainPkg, outputDir, build.FindOnly)
		if err != nil {
			return errors.Wrapf(err, "failed to get information for package %s for output directory %s", currPkg.MainPkg, outputDir)
		}

		// get location of main package on disk
		mainDir := mainPkg.Dir

		// get project import path and location of project package directory
		projectRootDir := mainDir
		projectImportPath := currPkg.MainPkg
		for i := 0; i < currPkg.DistanceToProjectPkg; i++ {
			projectRootDir = path.Dir(projectRootDir)
			projectImportPath = path.Dir(projectImportPath)
		}

		// copy project package into vendor directory in output dir if it does not already exist
		projectDestDir := path.Join(vendorDir, projectImportPath)

		if _, err := os.Stat(projectDestDir); os.IsNotExist(err) {
			if err := shutil.CopyTree(projectRootDir, projectDestDir, nil); err != nil {
				return errors.Wrapf(err, "failed to copy directory %s to %s", projectRootDir, projectDestDir)
			}
		} else if err != nil {
			return errors.Wrapf(err, "failed to stat %s", projectDestDir)
		}

		projectDestDirImport, err := build.ImportDir(projectDestDir, build.FindOnly)
		if err != nil {
			return errors.Wrapf(err, "unable to import project destination directory %s", projectDestDir)
		}
		projectDestDirImportPath := projectDestDirImport.ImportPath

		// rewrite imports for all files in copied directory
		fileSet := token.NewFileSet()
		foundMain := false
		goFiles := make(map[string]*ast.File)

		flagPkgImported := false
		if err := filepath.Walk(projectDestDir, func(currPath string, currInfo os.FileInfo, err error) error {
			if !currInfo.IsDir() && strings.HasSuffix(currInfo.Name(), ".go") {
				fileNode, err := parser.ParseFile(fileSet, currPath, nil, parser.ParseComments)
				if err != nil {
					return errors.Wrapf(err, "failed to parse file %s", currPath)
				}
				goFiles[currPath] = fileNode

				for _, currImport := range fileNode.Imports {
					currImportPathUnquoted, err := strconv.Unquote(currImport.Path.Value)
					if err != nil {
						return errors.Wrapf(err, "unable to unquote import %s", currImport.Path.Value)
					}

					updatedImport := ""
					if currImportPathUnquoted == "flag" {
						flagPkgImported = true
						updatedImport = path.Join(projectDestDirImportPath, "amalgomated_flag")
					} else if strings.HasPrefix(currImportPathUnquoted, projectImportPath) {
						updatedImport = strings.Replace(currImportPathUnquoted, projectImportPath, projectDestDirImportPath, -1)
					}

					if updatedImport != "" {
						if !astutil.RewriteImport(fileSet, fileNode, currImportPathUnquoted, updatedImport) {
							return errors.Errorf("failed to rewrite import from %s to %s", currImportPathUnquoted, updatedImport)
						}
					}
				}

				removeImportPathChecking(fileNode)

				// change package name for main packages
				if fileNode.Name.Name == "main" {
					fileNode.Name = ast.NewIdent(amalgomatedPackage)

					// find the main function
					mainFunc := findFunction(fileNode, "main")
					if mainFunc != nil {
						err = renameFunction(fileNode, "main", amalgomatedMain)
						if err != nil {
							return errors.Wrapf(err, "failed to rename function in file %s", currPath)
						}
						foundMain = true
					}
				}
			}
			return nil
		}); err != nil {
			return err
		}

		if !foundMain {
			return errors.Errorf("main method not found in package %s", currPkg.MainPkg)
		}

		if flagPkgImported {
			// if "flag" package is imported, add "flag" as a rewritten vendored dependency. This is done
			// because flag.CommandLine is a global variable that is often used by programs and problems can
			// arise if multiple amalgomated programs use it. A custom rewritten import is used rather than
			// vendoring so that the amalgomated program can itself be vendored.
			goRoot, err := dirs.GoRoot()
			if err != nil {
				return errors.WithStack(err)
			}
			fmtSrcDir := path.Join(goRoot, "src", "flag")
			fmtDstDir := path.Join(projectDestDir, "amalgomated_flag")
			if err := shutil.CopyTree(fmtSrcDir, fmtDstDir, vendorCopyOptions()); err != nil {
				return errors.Wrapf(err, "failed to copy directory %s to %s", projectRootDir, projectDestDir)
			}
		}

		for currGoFile, currNode := range goFiles {
			if err = writeAstToFile(currGoFile, currNode, fileSet); err != nil {
				return errors.Wrapf(err, "failed to write rewritten file %s", config)
			}
		}

		processedPkgs[currPkg] = true
	}
	return nil
}

func vendorCopyOptions() *shutil.CopyTreeOptions {
	return &shutil.CopyTreeOptions{
		Ignore: func(dir string, infos []os.FileInfo) []string {
			// ignore non-go files, go test files and testdata directories
			var ignore []string
			for _, currInfo := range infos {
				isTestDataDir := currInfo.IsDir() && currInfo.Name() == "testdata"
				notGoFile := !currInfo.IsDir() && !strings.HasSuffix(currInfo.Name(), ".go")
				goTestFile := !currInfo.IsDir() && strings.HasSuffix(currInfo.Name(), "_test.go")
				if isTestDataDir || notGoFile || goTestFile {
					ignore = append(ignore, currInfo.Name())
				}
			}
			return ignore
		},
		CopyFunction: shutil.Copy,
	}
}

func removeImportPathChecking(fileNode *ast.File) {
	var newCgList []*ast.CommentGroup
	for _, cg := range fileNode.Comments {
		var newCommentList []*ast.Comment
		for _, cc := range cg.List {
			// assume that any comment that starts with "// import" or "/* import" are import path checking
			// comments and don't add them to the new slice. This may omit some comments that are not
			// actually import checks, but downside is limited (it will just omit comment from repacked file).
			if !(strings.HasPrefix(cc.Text, "// import") || strings.HasPrefix(cc.Text, "/* import")) {
				newCommentList = append(newCommentList, cc)
			}
		}
		cg.List = newCommentList

		// CommentGroup assumes that len(List) > 0, so if logic above causes group to be empty, omit
		if len(cg.List) != 0 {
			newCgList = append(newCgList, cg)
		}
	}
	fileNode.Comments = newCgList
}

func addImports(file *ast.File, fileSet *token.FileSet, amalgomatedOutputDir string, config Config) error {
	processedPkgs := make(map[SrcPkg]bool, len(config.Pkgs))
	for _, name := range sortedKeys(config.Pkgs) {
		progPkg := config.Pkgs[name]

		// if package has already been imported, skip (can't have multiple imports for the same package)
		if processedPkgs[progPkg] {
			continue
		}

		repackagedDirPath := path.Join(amalgomatedOutputDir, progPkg.MainPkg)
		repackagedDirImportResult, err := build.ImportDir(repackagedDirPath, build.FindOnly)
		if err != nil {
			return errors.Wrapf(err, "failed to import directory %s", repackagedDirPath)
		}

		repackagedImportPath := repackagedDirImportResult.ImportPath
		added := astutil.AddNamedImport(fileSet, file, name, repackagedImportPath)
		if !added {
			return errors.Errorf("failed to add import %s", repackagedImportPath)
		}
		processedPkgs[progPkg] = true
	}
	return nil
}

func sortImports(file *ast.File) {
	for _, decl := range file.Decls {
		if gen, ok := decl.(*ast.GenDecl); ok && gen.Tok == token.IMPORT {
			sort.Sort(importSlice(gen.Specs))
			break
		}
	}
}

func getFirstToken(file *ast.File, t token.Token) *ast.GenDecl {
	for _, currDecl := range file.Decls {
		switch currDecl.(type) {
		case *ast.GenDecl:
			genDecl := currDecl.(*ast.GenDecl)
			if genDecl.Tok == t {
				return genDecl
			}
		}
	}
	return nil
}

func setVarCompositeLiteralElements(file *ast.File, constName string, elems []ast.Expr) error {
	decl := getFirstToken(file, token.VAR)
	if decl == nil {
		return errors.Errorf("could not find token of type VAR in %s", file.Name)
	}

	var constExpr ast.Expr
	for _, currSpec := range decl.Specs {
		// declaration is already known to be of type const, so all of the specs are ValueSpec
		valueSpec := currSpec.(*ast.ValueSpec)
		for i, valueSpecName := range valueSpec.Names {
			if valueSpecName.Name == constName {
				constExpr = valueSpec.Values[i]
				break
			}
		}
	}

	if constExpr == nil {
		return errors.Errorf("could not find variable with name %s in given declaration", constName)
	}

	var compLit *ast.CompositeLit
	var ok bool
	if compLit, ok = constExpr.(*ast.CompositeLit); !ok {
		return errors.Errorf("variable %s did not have a composite literal value", constName)
	}

	compLit.Elts = elems

	return nil
}

func createMapLiteralEntries(pkgs map[string]SrcPkg) []ast.Expr {
	// if multiple commands refer to the same package, the command that is lexicographically first is the one that
	// is used for the named import. Create a map that stores the mapping from the package to the import name.
	pkgToFirstCmdMap := make(map[SrcPkg]string, len(pkgs))
	for _, name := range sortedKeys(pkgs) {
		if _, ok := pkgToFirstCmdMap[pkgs[name]]; !ok {
			// add mapping only if it does not already exist
			pkgToFirstCmdMap[pkgs[name]] = name
		}
	}

	var entries []ast.Expr
	for _, name := range sortedKeys(pkgs) {
		entries = append(entries, createMapKeyValueExpression(name, pkgToFirstCmdMap[pkgs[name]]))
	}
	return entries
}

// createMapKeyValueExpression creates a new map key value function expression of the form "{{name}}": func() { {{namedImport}}.{{amalgomatedMain}}() }.
// In most cases "name" and "namedImport" will be the same, but if multiple commands refer to the same package, then the
// commands that are lexicographically later should refer to the named import of the first command.
func createMapKeyValueExpression(name, namedImport string) *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf(`"%v"`, name),
		},
		Value: &ast.FuncLit{
			Type: &ast.FuncType{
				Params: &ast.FieldList{},
			},
			Body: &ast.BlockStmt{
				List: []ast.Stmt{
					&ast.ExprStmt{
						X: &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent(namedImport),
								Sel: ast.NewIdent(amalgomatedMain),
							},
						},
					},
				},
			},
		},
	}
}

func renameFunction(fileNode *ast.File, originalName, newName string) error {
	originalFunc := findFunction(fileNode, originalName)
	if originalFunc == nil {
		return errors.Errorf("function %s does not exist", originalName)
	}

	if findFunction(fileNode, newName) != nil {
		return errors.Errorf("cannot rename function %s to %s because a function with the new name already exists", originalName, newName)
	}

	originalFunc.Name = ast.NewIdent(newName)
	return nil
}

func findFunction(fileNode *ast.File, funcName string) *ast.FuncDecl {
	for _, currDecl := range fileNode.Decls {
		switch t := currDecl.(type) {
		case *ast.FuncDecl:
			if t.Name.Name == funcName {
				return currDecl.(*ast.FuncDecl)
			}
		}
	}
	return nil
}

func writeAstToFile(path string, fileNode *ast.File, fileSet *token.FileSet) (writeErr error) {
	outputFile, err := os.Create(path)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", path)
	}
	defer func() {
		if err := outputFile.Close(); err != nil {
			writeErr = errors.Errorf("failed to close file %s", path)
		}
	}()
	if err := printer.Fprint(outputFile, fileSet, fileNode); err != nil {
		return errors.Wrapf(err, "failed to write to file %s", path)
	}
	return nil
}

func sortedKeys(pkgs map[string]SrcPkg) []string {
	sortedKeys := make([]string, 0, len(pkgs))
	for currKey := range pkgs {
		sortedKeys = append(sortedKeys, currKey)
	}
	sort.Strings(sortedKeys)
	return sortedKeys
}
