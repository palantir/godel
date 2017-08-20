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
	"sort"
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

// TagsMatcher returns a Matcher that matches all packages that are matched by the provided tags. If no tags are
// provided, returns nil. If the tags consist of a single tag named "all", the returned matcher matches the union of all
// known tags. If the tags consist of a single tag named "none", the returned matcher matches everything except the
// union of all known tags (untagged tests).
func TagsMatcher(tags []string, cfg params.GUnit) (matcher.Matcher, error) {
	if len(tags) == 0 {
		// if no tags were provided, does not match anything
		return nil, nil
	}

	if len(tags) == 1 {
		var allMatchers []matcher.Matcher
		for _, matcher := range cfg.Tags {
			allMatchers = append(allMatchers, matcher)
		}
		anyTagMatcher := matcher.Any(allMatchers...)
		switch tags[0] {
		case params.AllTagName:
			// if tags contains only a single tag that is the "all" tag, return matcher that matches union of all tags
			return anyTagMatcher, nil
		case params.NoneTagName:
			// if tags contains only a single tag that is the "none" tag, return matcher that matches not of union of all tags
			return matcher.Not(anyTagMatcher), nil
		}
	}

	// due to previous check, if "all" or "none" tag exists at this point it means that it was one of multiple tags
	for _, tag := range tags {
		switch tag {
		case params.AllTagName, params.NoneTagName:
			return nil, errors.Errorf("if %q tag is specified, it must be the only tag specified", tag)
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
		var allTags []string
		for tag := range cfg.Tags {
			allTags = append(allTags, fmt.Sprintf("%q", tag))
		}
		sort.Strings(allTags)
		validTagsOutput := fmt.Sprintf("Valid tags: %v", strings.Join(allTags, ", "))
		if len(allTags) == 0 {
			validTagsOutput = "No tags are defined."
		}
		return nil, fmt.Errorf("Tags %v not defined in configuration. %s", strings.Join(missingTags, ", "), validTagsOutput)
	}

	// not possible: if initial tags were empty then should have already returned, if specified tags did not match then
	// missing block should have executed and returned, so at this point matchers must exist
	if len(tagMatchers) == 0 {
		panic("no matching tags found")
	}

	// OR of tags
	return matcher.Any(tagMatchers...), nil
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
