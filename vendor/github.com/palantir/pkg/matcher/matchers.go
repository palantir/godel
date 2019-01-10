// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package matcher

import (
	"fmt"
	"path"
	"path/filepath"
	"regexp"
)

type Matcher interface {
	Match(relPath string) bool
}

type allMatcher []Matcher

func (m allMatcher) Match(relPath string) bool {
	nonNilMatcherExists := false
	for _, currMatcher := range []Matcher(m) {
		if currMatcher != nil {
			nonNilMatcherExists = true
			if !currMatcher.Match(relPath) {
				return false
			}
		}
	}
	return nonNilMatcherExists
}

// All returns a compound Matcher that returns true if all of its provided non-nil Matchers return true. Returns false
// if no matchers are provided or if all of the provided matchers are nil.
func All(matchers ...Matcher) Matcher {
	return allMatcher(append([]Matcher{}, matchers...))
}

type anyMatcher []Matcher

func (m anyMatcher) Match(relPath string) bool {
	for _, currMatcher := range []Matcher(m) {
		if currMatcher != nil && currMatcher.Match(relPath) {
			return true
		}
	}
	return false
}

// Not returns a matcher that returns the negation of the provided matcher.
func Not(matcher Matcher) Matcher {
	return notMatcher{
		matcher: matcher,
	}
}

type notMatcher struct {
	matcher Matcher
}

func (m notMatcher) Match(relPath string) bool {
	return !m.matcher.Match(relPath)
}

// Any returns a compound Matcher that returns true if any of the provided Matchers return true.
func Any(matchers ...Matcher) Matcher {
	return anyMatcher(append([]Matcher{}, matchers...))
}

// Hidden returns a matcher that matches all hidden files or directories (any path that begins with `.`).
func Hidden() Matcher {
	return Name(`\..+`)
}

// Name returns a Matcher that matches the on the name of all of the components of a path using the provided
// expressions. Each part of the path (except for ".." components, which are ignored and cannot be matched) is tested
// against the expressions independently (no path separators). The name must fully match the expression to be considered
// a match.
func Name(regexps ...string) Matcher {
	compiled := make([]*regexp.Regexp, len(regexps))
	for i, curr := range regexps {
		compiled[i] = regexp.MustCompile(curr)
	}
	return nameMatcher(compiled)
}

type nameMatcher []*regexp.Regexp

func (m nameMatcher) Match(inputRelPath string) bool {
	for _, currSubpath := range allSubpaths(inputRelPath) {
		currName := path.Base(currSubpath)
		// do not match relative path components
		if currName == ".." {
			continue
		}
		for _, currRegExp := range []*regexp.Regexp(m) {
			matchLoc := currRegExp.FindStringIndex(currName)
			if len(matchLoc) > 0 && matchLoc[0] == 0 && matchLoc[1] == len(currName) {
				return true
			}
		}
	}
	return false
}

// Path returns a Matcher that matches any path that matches or is a subpath of any of the provided paths. For example,
// a value of "foo" would match the relative directory "foo" and all of its sub-paths ("foo/bar", "foo/bar.txt"), but
// not every directory named "foo" (would not match "bar/foo"). Matches are done using glob matching (same as
// filepath.Match). However, unlike filepath.Match, subpath matches will match all of the sub-paths of a given match as
// well (for example, the pattern "foo/*/bar" matches "foo/*/bar/baz").
func Path(paths ...string) Matcher {
	return &pathMatcher{paths: paths, glob: true}
}

// PathLiteral returns a Matcher that is equivalent to that returned by Paths except that matches are done using string
// equality rather than using glob matching.
func PathLiteral(paths ...string) Matcher {
	return &pathMatcher{paths: paths, glob: false}
}

type pathMatcher struct {
	paths []string
	glob  bool
}

func (m *pathMatcher) Match(inputRelPath string) bool {
	subpaths := allSubpaths(inputRelPath)
	for _, currMatcherPath := range m.paths {
		for _, currSubpath := range subpaths {
			var match bool
			if m.glob {
				var err error
				match, err = filepath.Match(currMatcherPath, currSubpath)
				if err != nil {
					// only possible error is bad pattern
					panic(fmt.Sprintf("filepath: Match(%q): %v", currMatcherPath, err))
				}
			} else {
				match = currMatcherPath == currSubpath
			}
			if match {
				return true
			}
		}
	}
	return false
}

// allSubpaths returns the provided relative path and all of its subpaths up to (but not including) ".". For example,
// "foo/bar/baz.txt" returns [foo/bar/baz.txt foo/bar foo], while "foo.txt" returns [foo.txt]. This applies for ".."
// paths as well: a path of the form "../foo/bar/baz.txt" returns [../foo/bar/baz.txt ../foo/bar ../foo ..]. Returns nil
// if the input path is not a relative path.
func allSubpaths(relPath string) []string {
	if path.IsAbs(relPath) {
		return nil
	}
	var subpaths []string
	for currRelPath := relPath; currRelPath != "."; currRelPath = path.Dir(currRelPath) {
		subpaths = append(subpaths, currRelPath)
	}
	return subpaths
}
