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

package amalgomated

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

const (
	pkgsFlagName	= "pkgs"
	verboseFlagName	= "verbose"
)

var (
	pkgsFlag	= flag.StringSlice{
		Name:		pkgsFlagName,
		Usage:		"paths to the packages to check",
		Optional:	true,
	}
	verboseFlag	= flag.BoolFlag{
		Name:	verboseFlagName,
		Usage:	"print verbose analysis of all imports that have multiple aliases",
		Alias:	"v",
	}
)

func AmalgomatedMain() {
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack))
	app.Flags = append(app.Flags,
		pkgsFlag,
		verboseFlag,
	)
	app.Action = func(ctx cli.Context) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return errors.Wrapf(err, "Failed to get working directory")
		}
		return doImportAlias(wd, ctx.Slice(pkgsFlagName), ctx.Bool(verboseFlagName), ctx.App.Stdout)
	}
	os.Exit(app.Run(os.Args))
}

func doImportAlias(projectDir string, pkgPaths []string, verbose bool, w io.Writer) error {
	if !path.IsAbs(projectDir) {
		return errors.Errorf("projectDir %s must be an absolute path", projectDir)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return errors.Errorf("GOPATH environment variable must be set")
	}

	if relPath, err := filepath.Rel(path.Join(gopath, "src"), projectDir); err != nil || strings.HasPrefix(relPath, "../") {
		return errors.Wrapf(err, "Project directory %s must be a subdirectory of $GOPATH/src (%s)", projectDir, path.Join(gopath, "src"))
	}

	if len(pkgPaths) == 0 {
		pkgs, err := pkgpath.PackagesInDir(projectDir, pkgpath.DefaultGoPkgExcludeMatcher())
		if err != nil {
			return errors.Wrapf(err, "Failed to list packages")
		}

		pkgPaths, err = pkgs.Paths(pkgpath.Relative)
		if err != nil {
			return errors.Wrapf(err, "Failed to convert package paths")
		}
	}

	projectImportInfo := NewProjectImportInfo()
	for _, pkgPath := range pkgPaths {
		currPath := path.Join(projectDir, pkgPath)
		fis, err := ioutil.ReadDir(currPath)
		if err != nil {
			return errors.Wrapf(err, "Failed to list contents of directory %s", currPath)
		}
		for _, fi := range fis {
			if !fi.IsDir() && strings.HasSuffix(fi.Name(), ".go") {
				currFile := path.Join(currPath, fi.Name())
				if err := projectImportInfo.AddImportAliasesFromFile(currFile); err != nil {
					return errors.Wrapf(err, "failed to determine imports in file %s", currFile)
				}
			}
		}
	}

	importsToAliases := projectImportInfo.ImportsToAliases()
	var pkgsWithMultipleAliases []string
	pkgsWithMultipleAliasesMap := make(map[string]struct{})
	for k, v := range importsToAliases {
		if len(v) > 1 {
			// package is imported using more than 1 alias
			pkgsWithMultipleAliases = append(pkgsWithMultipleAliases, k)
			pkgsWithMultipleAliasesMap[k] = struct{}{}
		}
	}
	sort.Strings(pkgsWithMultipleAliases)
	if len(pkgsWithMultipleAliases) > 0 {
		var output []string
		if verbose {
			for _, k := range pkgsWithMultipleAliases {
				output = append(output, fmt.Sprintf("%s is imported using multiple different aliases:", k))
				for _, currAliasInfo := range importsToAliases[k] {
					var files []string
					for k, v := range currAliasInfo.Occurrences {
						relPkgPath, err := pkgpath.NewAbsPkgPath(k).Rel(projectDir)
						if err != nil {
							return errors.Wrapf(err, "failed to get package path")
						}
						relPkgPath = strings.TrimLeft(relPkgPath, "./")
						files = append(files, fmt.Sprintf("%s:%d:%d", relPkgPath, v.Line, v.Column))
					}
					sort.Strings(files)

					var numFilesMsg string
					if len(currAliasInfo.Occurrences) == 1 {
						numFilesMsg = "(1 file)"
					} else {
						numFilesMsg = fmt.Sprintf("(%d files)", len(currAliasInfo.Occurrences))
					}
					output = append(output, fmt.Sprintf("\t%s %s:\n\t\t%s", currAliasInfo.Alias, numFilesMsg, strings.Join(files, "\n\t\t")))
				}
			}
		} else {
			filesToAliases := projectImportInfo.FilesToImportAliases()

			var relPkgPaths []string
			relPkgPathToFile := make(map[string]string)
			for file := range filesToAliases {
				relPkgPath, err := pkgpath.NewAbsPkgPath(file).GoPathSrcRel()
				if err != nil {
					return errors.Wrapf(err, "failed to get package path")
				}
				relPkgPaths = append(relPkgPaths, relPkgPath)
				relPkgPathToFile[relPkgPath] = file
			}
			sort.Strings(relPkgPaths)

			for _, relPkgPath := range relPkgPaths {
				file := relPkgPathToFile[relPkgPath]
				for _, alias := range filesToAliases[file] {
					if _, ok := pkgsWithMultipleAliasesMap[alias.ImportPath]; !ok {
						continue
					}
					status := projectImportInfo.GetAliasStatus(alias.Alias, alias.ImportPath)
					if status.OK {
						continue
					}

					relPkgPath, err := pkgpath.NewAbsPkgPath(file).Rel(projectDir)
					if err != nil {
						return errors.Wrapf(err, "failed to get package path")
					}
					relPkgPath = strings.TrimLeft(relPkgPath, "./")
					msg := fmt.Sprintf("%s:%d:%d: uses alias %q to import package %s. %s.", relPkgPath, alias.Pos.Line, alias.Pos.Column, alias.Alias, alias.ImportPath, status.Recommendation)
					output = append(output, msg)
				}
			}
		}
		return errors.New(strings.Join(output, "\n"))
	}
	return nil
}
