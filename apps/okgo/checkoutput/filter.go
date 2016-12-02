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

package checkoutput

import (
	"regexp"

	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

type Filterer interface {
	Filter(line Issue) (bool, error)
}

func ApplyFilters(lineInfo []Issue, filters []Filterer) ([]Issue, error) {
	output := make([]Issue, 0, len(lineInfo))
	for _, currLine := range lineInfo {
		filterOutCurrLine := false
		var err error
		for _, currLineFilter := range filters {
			if filterOutCurrLine, err = currLineFilter.Filter(currLine); err != nil {
				return nil, errors.Wrapf(err, "failed to apply filter to lineInfo %v", currLine)
			} else if filterOutCurrLine {
				break
			}
		}
		if !filterOutCurrLine {
			output = append(output, currLine)
		}
	}
	return output, nil
}

func NamePathFilter(namePattern string) Filterer {
	return MatcherFilter(matcher.Name(namePattern))
}

func RelativePathFilter(relativePathToExclude string) Filterer {
	return MatcherFilter(matcher.Path(relativePathToExclude))
}

type matcherFilter struct {
	m matcher.Matcher
}

func (f matcherFilter) Filter(line Issue) (bool, error) {
	currLineRelativePath, err := line.Path(pkgpath.Relative)
	if err != nil {
		return false, errors.Wrapf(err, "failed to convert line %v to relative path", line)
	}
	return f.m != nil && f.m.Match(currLineRelativePath), nil
}

func MatcherFilter(matcher matcher.Matcher) Filterer {
	return &matcherFilter{m: matcher}
}

type msgRegexpFilter struct {
	exp string
}

// cache for regular expressions. Used so that the msgRegexpFilter can store a string rather than a regexp, but avoids
// re-compiling the same regular expression repeatedly. Not thread-safe, but OK.
var regExpCache = make(map[string]*regexp.Regexp)

func (f msgRegexpFilter) Filter(line Issue) (bool, error) {
	exp, ok := regExpCache[f.exp]
	if !ok {
		// not cached -- compile regexp and add to cache
		exp = regexp.MustCompile(f.exp)
		regExpCache[f.exp] = exp
	}
	return exp.MatchString(line.Message()), nil
}

func MessageRegexpFilter(messagePattern string) Filterer {
	// add to cache
	regExpCache[messagePattern] = regexp.MustCompile(messagePattern)
	return msgRegexpFilter{exp: messagePattern}
}
