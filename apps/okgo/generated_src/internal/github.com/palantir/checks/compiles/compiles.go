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
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/errorstringer"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
	"golang.org/x/tools/go/loader"
)

func AmalgomatedMain() {
	const pkgsFlagName = "pkgs"
	app := cli.NewApp(cli.DebugHandler(errorstringer.SingleStack))
	app.Flags = append(app.Flags, flag.StringSlice{
		Name:	pkgsFlagName,
		Usage:	"paths to the pacakges to check",
	})
	app.Action = func(ctx cli.Context) error {
		wd, err := dirs.GetwdEvalSymLinks()
		if err != nil {
			return errors.Wrapf(err, "Failed to get working directory")
		}
		return doCompiles(wd, ctx.Slice(pkgsFlagName), ctx.App.Stdout)
	}
	os.Exit(app.Run(os.Args))
}

func doCompiles(projectDir string, pkgPaths []string, w io.Writer) error {
	if !path.IsAbs(projectDir) {
		return fmt.Errorf("projectDir must be an absolute path: %v", projectDir)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		return fmt.Errorf("GOPATH environment variable must be set")
	}

	if relPath, err := filepath.Rel(path.Join(gopath, "src"), projectDir); err != nil || strings.HasPrefix(relPath, "../") {
		return fmt.Errorf("Project directory %v must be a subdirectory of $GOPATH/src (%v)", projectDir, path.Join(gopath, "src"))
	}

	if len(pkgPaths) == 0 {
		pkgs, err := pkgpath.PackagesInDir(projectDir, pkgpath.DefaultGoPkgExcludeMatcher())
		if err != nil {
			return fmt.Errorf("Failed to list packages: %v", err)
		}

		pkgPaths, err = pkgs.Paths(pkgpath.GoPathSrcRelative)
		if err != nil {
			return fmt.Errorf("Failed to convert package paths: %v", err)
		}
	}

	cfg := loader.Config{}
	for _, currPkgPath := range pkgPaths {
		cfg.ImportWithTests(currPkgPath)
	}
	cfg.TypeChecker.Error = func(e error) {
		fmt.Fprintln(w, e)
	}

	if _, err := cfg.Load(); err != nil {
		// return blank error if any errors were encountered during load. Load function prints errors to writer
		// in proper format as they are encountered so no need to create any other output.
		return fmt.Errorf("")
	}

	return nil
}
