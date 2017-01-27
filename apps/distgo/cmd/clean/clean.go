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

package clean

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

const (
	defaultBuildOutputDir = "build"
	defaultDistOutputDir  = "dist"

	enotfound = "no such file or directory"

	versionRegex = `v?\d+\.\d+\.\d+(-\d+-g[0-9a-e]+)?(\.dirty)?`
)

var (
	versionPattern = regexp.MustCompile(versionRegex)
)

type cleanableDirectory struct {
	path    string
	pattern *regexp.Regexp
}

func Clean(products []string, cfg params.Project, force bool, stdout io.Writer) error {
	var dirsToRemove []cleanableDirectory

	// Only remove default output directories if no products specified
	if len(products) == 0 {
		if cfg.BuildOutputDir == "" {
			cfg.BuildOutputDir = defaultBuildOutputDir
		}
		dirsToRemove = append(dirsToRemove, cleanableDirectory{path: cfg.BuildOutputDir, pattern: versionPattern})

		if cfg.DistOutputDir == "" {
			cfg.DistOutputDir = defaultDistOutputDir
		}
		// cannot provide a better pattern here for distributions since we don't know which products use it
		dirsToRemove = append(dirsToRemove, cleanableDirectory{path: cfg.DistOutputDir, pattern: versionPattern})
	}

	dirsToRemove = append(dirsToRemove, productsDirsToRemove(products, cfg)...)
	err := removeDirs(unique(dirsToRemove), force, stdout)
	if err != nil {
		return errors.Wrap(err, "failed to clean products")
	}

	return nil
}

func productsDirsToRemove(products []string, cfg params.Project) []cleanableDirectory {
	var dirsToRemove []cleanableDirectory

	// remove all products
	if len(products) == 0 {
		for name := range cfg.Products {
			products = append(products, name)
		}
	}

	for _, product := range products {
		p := cfg.Products[product]
		buildPattern := regexp.MustCompile(versionRegex)
		distPattern := regexp.MustCompile(product + "-" + versionRegex)

		if p.Build.OutputDir != "" {
			dirsToRemove = append(dirsToRemove, cleanableDirectory{path: p.Build.OutputDir, pattern: buildPattern})
		}

		for _, dist := range p.Dist {
			if dist.OutputDir != "" {
				dirsToRemove = append(dirsToRemove, cleanableDirectory{path: dist.OutputDir, pattern: distPattern})
			}
		}
	}

	return dirsToRemove
}

func removeDirs(dirs []cleanableDirectory, force bool, stdout io.Writer) error {
	for _, dir := range dirs {
		files, err := ioutil.ReadDir(dir.path)
		if err != nil {
			if strings.Contains(err.Error(), enotfound) {
				return nil
			}
			return errors.WithStack(err)
		}

		fmt.Fprintf(stdout, "Removing %s\n", dir.path)
		for _, file := range files {
			if matchIdx := dir.pattern.FindStringIndex(file.Name()); matchIdx == nil {
				if force {
					fmt.Fprintf(stdout, "removing non-build file %s due to force\n", file.Name())
				} else {
					fmt.Fprintf(stdout, "ignoring file %s as it does not appear to be produced by build\n", file.Name())
					continue
				}
			}

			if err := os.RemoveAll(path.Join(dir.path, file.Name())); err != nil {
				return errors.Wrapf(err, "failed to remove file %s", file.Name())
			}
		}

	}
	return nil
}

func unique(dirs []cleanableDirectory) []cleanableDirectory {
	uniqueDirectories := []cleanableDirectory{}
	for _, dir := range dirs {
		if !contains(uniqueDirectories, dir) {
			uniqueDirectories = append(uniqueDirectories, dir)
		}
	}
	return uniqueDirectories
}

func contains(dirs []cleanableDirectory, dir cleanableDirectory) bool {
	for _, d := range dirs {
		if d.path == dir.path {
			return true
		}
	}
	return false
}
