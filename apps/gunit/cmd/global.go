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

package cmd

import (
	"fmt"
	"strings"

	"github.com/palantir/amalgomate/amalgomated"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/gunit/generated_src"
	"github.com/palantir/godel/apps/gunit/params"
)

var Library = amalgomated.NewCmdLibrary(amalgomatedtesters.Instance())

const (
	junitOutputPathFlagName = "junit-output"
	raceFlagName            = "race"
	tagsFlagName            = "tags"
	verboseFlagName         = "verbose"
	verboseFlagAlias        = "v"
)

var (
	GlobalFlags = []flag.Flag{
		flag.StringFlag{
			Name:  tagsFlagName,
			Usage: "Run tests that are part of the provided tags (use commas to separate multiple tags)",
		},
		flag.BoolFlag{
			Name:  verboseFlagName,
			Alias: verboseFlagAlias,
			Usage: "Enable verbose output for tests",
		},
		flag.BoolFlag{
			Name:  raceFlagName,
			Usage: "Enable race detector for tests",
		},
		flag.StringFlag{
			Name:  junitOutputPathFlagName,
			Usage: "Path to JUnit XML output (if provided, verbose flag is set to true)",
		},
	}
)

func Tags(ctx cli.Context) []string {
	if !ctx.Has(tagsFlagName) {
		return nil
	}
	return strings.Split(strings.ToLower(ctx.String(tagsFlagName)), ",")
}

func Verbose(ctx cli.Context) bool {
	return ctx.Bool(verboseFlagName)
}

func Race(ctx cli.Context) bool {
	return ctx.Bool(raceFlagName)
}

func JUnitOutputPath(ctx cli.Context) string {
	return ctx.String(junitOutputPathFlagName)
}

type trueMatcher struct{}

func (t *trueMatcher) Match(relPath string) bool {
	return true
}

// TagsMatcher returns a Matcher that matches the provided tags. Returns nil if the provided slice of tags is empty or
// if the provided tags do not match any of the tags specified in the configuration. Returns an error if any of the
// provided tags are not specified in the configuration.
func TagsMatcher(tags []string, cfg params.GUnit) (matcher.Matcher, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	for _, tag := range tags {
		if tag == params.AllTagName {
			// contains "all" tag
			if len(tags) == 1 {
				// if "all" is the only tag specified, return a matcher that matches all paths
				return &trueMatcher{}, nil
			}
			return nil, errors.Errorf(`if "all" tag is specified, it must be the only tag specified`)
		}
	}

	var tagMatchers []matcher.Matcher
	var missingTags []string
	for _, tag := range tags {
		if include, ok := cfg.Tags[tag]; ok {
			tagMatchers = append(tagMatchers, include)
		} else {
			missingTags = append(missingTags, fmt.Sprintf("%q", tag))
		}
	}

	if len(missingTags) > 0 {
		return nil, fmt.Errorf("invalid tags: %v", strings.Join(missingTags, ", "))
	}

	if len(tagMatchers) == 0 {
		return nil, nil
	}

	return matcher.Any(tagMatchers...), nil
}

// AllTagsMatcher returns a matcher that matches paths that are part of any of the tags defined in the provided
// configuration.
func AllTagsMatcher(cfg params.GUnit) matcher.Matcher {
	tags := make([]string, 0, len(cfg.Tags))
	for tag := range cfg.Tags {
		tags = append(tags, strings.ToLower(tag))
	}
	// error cannot occur because tags are known to exist
	m, err := TagsMatcher(tags, cfg)
	if err != nil {
		panic(err)
	}
	if m == nil {
		return nil
	}
	return m
}

// PkgPaths returns a slice that contains the relative package paths for the packages "pkgPaths" relative to the
// project directory "wd" excluding any of the paths that match the provided "exclude" Matcher. If "pkgPaths" is an
// empty slice, then all of the packages in "wd" (except those that match the "exclude" matcher) are returned.
func PkgPaths(pkgPaths []string, wd string, exclude matcher.Matcher) ([]string, error) {
	var pkgs pkgpath.Packages
	var err error
	if len(pkgPaths) == 0 {
		// if input slice is empty, return all matching packages
		pkgs, err = pkgpath.PackagesInDir(wd, exclude)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to list packages in %s", wd)
		}
	} else {
		// otherwise, filter provided packages and return those that are not excluded
		var nonExcludedPkgs []string
		for _, currPkg := range pkgPaths {
			if !exclude.Match(currPkg) {
				nonExcludedPkgs = append(nonExcludedPkgs, currPkg)
			}
		}
		pkgs, err = pkgpath.PackagesFromPaths(wd, nonExcludedPkgs)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse %v as packages", pkgPaths)
		}
	}
	resultPkgPaths, err := pkgs.Paths(pkgpath.Relative)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get relative paths for packages %v", pkgs)
	}
	return resultPkgPaths, nil
}
