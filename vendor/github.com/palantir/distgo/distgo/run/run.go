// Copyright 2016 Palantir Technologies, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package run

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

func Product(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam, runArgs []string, stdout, stderr io.Writer) error {
	if productParam.Build == nil {
		return errors.Errorf("product %s has no build configuration defined", productParam.ID)
	}

	mainPkgDir := path.Join(projectInfo.ProjectDir, productParam.Build.MainPkg)
	mainPkgGoFiles, err := mainPkgGoFiles(mainPkgDir)
	if err != nil {
		return errors.Wrapf(err, "failed to find Go files for main package")
	}

	cmd := exec.Command("go")
	args := []string{cmd.Path, "run"}

	// add build arguments for product
	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return errors.Wrapf(err, "failed to compute output info")
	}
	buildArgs, err := productParam.Build.BuildArgs(productTaskOutputInfo)
	if err != nil {
		return err
	}
	args = append(args, buildArgs...)

	for _, goFiles := range mainPkgGoFiles {
		args = append(args, path.Join(mainPkgDir, goFiles))
	}
	if productParam.Run != nil {
		args = append(args, productParam.Run.Args...)
	}
	args = append(args, runArgs...)
	cmd.Args = args

	cmd.Stdout = stdout
	cmd.Stderr = stderr
	cmd.Stdin = os.Stdin

	fmt.Fprintln(stdout, strings.Join(args, " "))
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "go run failed")
	}
	return nil
}

// mainPkgGoFiles returns the names of all of the non-test Go files in the "main" pkg of the specified directory.
// Returns an error if there are no files in the "main" package that declares a "main" function (or if there are
// multiple such files).
func mainPkgGoFiles(mainPkgDir string) ([]string, error) {
	fileInfos, err := ioutil.ReadDir(mainPkgDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list files in directory %v", mainPkgDir)
	}
	var mainPkgFileNames []string
	var mainFuncFileNames []string
	for _, currFile := range fileInfos {
		currFilePath := path.Join(mainPkgDir, currFile.Name())
		if !currFile.IsDir() && strings.HasSuffix(currFile.Name(), ".go") && !strings.HasSuffix(currFile.Name(), "_test.go") {
			fset := token.NewFileSet()
			fnode, err := parser.ParseFile(fset, currFilePath, nil, parser.ParseComments)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to parse file %v", currFilePath)
			}

			// find main package
			if fnode.Name.Name == "main" {
				mainPkgFileNames = append(mainPkgFileNames, currFile.Name())
				if hasMainFunc(fnode) {
					mainFuncFileNames = append(mainFuncFileNames, currFile.Name())
				}
			}
		}
	}

	switch len(mainFuncFileNames) {
	case 0:
		return nil, errors.Errorf("no go file with main package and main function exists in directory %v", mainPkgDir)
	case 1:
		return mainPkgFileNames, nil
	default:
		return nil, errors.Errorf("directory %v contain multiple files that have main package and main function: %v", mainPkgDir, mainFuncFileNames)
	}
}

func hasMainFunc(node *ast.File) bool {
	for _, currDecl := range node.Decls {
		switch t := currDecl.(type) {
		case *ast.FuncDecl:
			if t.Name.Name == "main" {
				return true
			}
		}
	}
	return false
}
