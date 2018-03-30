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

package installupdate

import (
	"fmt"
	"regexp"
	"strconv"
)

type Type int

var unknown = Type(-1)

const (
	ReleaseCandidate Type = iota
	ReleaseCandidateSnapshot
	Release
	ReleaseSnapshot
	NonOrderable
)

func (t Type) String() string {
	switch t {
	case ReleaseCandidate:
		return "ReleaseCandidate"
	case ReleaseCandidateSnapshot:
		return "ReleaseCandidateSnapshot"
	case Release:
		return "Release"
	case ReleaseSnapshot:
		return "ReleaseSnapshot"
	case NonOrderable:
		return "NonOrderable"
	}
	return strconv.Itoa(int(t))
}

var releaseRegexps = []*regexp.Regexp{
	ReleaseCandidate:         regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9])+-rc([0-9]+)$`),
	ReleaseCandidateSnapshot: regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9])+-rc([0-9]+)-([0-9]+)-g[a-f0-9]+$`),
	Release:                  regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9])+$`),
	ReleaseSnapshot:          regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9])+-([0-9]+)-g[a-f0-9]+$`),
	NonOrderable:             regexp.MustCompile(`^([0-9]+)\.([0-9]+)\.([0-9])+(-[a-z0-9-]+)?(\.dirty)?$`),
}

type godelVersion struct {
	version string

	// computed once on construction and stored
	typ                      Type
	majorVersionNum          int
	minorVersionNum          int
	patchVersionNum          int
	firstSequenceVersionNum  *int
	secondSequenceVersionNum *int
}

func (v godelVersion) String() string {
	return v.version
}

func (v godelVersion) Type() Type {
	return getType(v.version)
}

func (v godelVersion) Orderable() bool {
	typ := v.Type()
	return typ >= ReleaseCandidate && typ < NonOrderable
}

func (v godelVersion) Value() string {
	return v.version
}

func (v godelVersion) MajorVersionNum() int {
	return v.majorVersionNum
}

func (v godelVersion) MinorVersionNum() int {
	return v.minorVersionNum
}

func (v godelVersion) PatchVersionNum() int {
	return v.patchVersionNum
}

func (v godelVersion) FirstSequenceVersionNum() *int {
	return v.firstSequenceVersionNum
}

func (v godelVersion) SecondSequenceVersionNum() *int {
	return v.secondSequenceVersionNum
}

func newGodelVersion(v string) (godelVersion, error) {
	typ := getType(v)
	if typ == unknown {
		return godelVersion{}, fmt.Errorf("%s is not a valid SLS version", v)
	}

	matches := releaseRegexps[typ].FindStringSubmatch(v)

	var firstSequenceVersionNum *int
	if typ != NonOrderable && len(matches) > 4 {
		n := mustAtoI(matches[4])
		firstSequenceVersionNum = &n
	}
	var secondSequenceVersionNum *int
	if typ != NonOrderable && len(matches) > 5 {
		n := mustAtoI(matches[5])
		secondSequenceVersionNum = &n
	}

	return godelVersion{
		version:                  v,
		typ:                      typ,
		majorVersionNum:          mustAtoI(matches[1]),
		minorVersionNum:          mustAtoI(matches[2]),
		patchVersionNum:          mustAtoI(matches[3]),
		firstSequenceVersionNum:  firstSequenceVersionNum,
		secondSequenceVersionNum: secondSequenceVersionNum,
	}, nil
}

func getType(v string) Type {
	for i, regExp := range releaseRegexps {
		if regExp.MatchString(v) {
			return Type(i)
		}
	}
	return unknown
}

// CompareTo compares the receiver to the provided version. If either version is not orderable (as defined by the spec),
// always returns -1 and false. If both versions are orderable, then returns -1 if the receiver is less than the
// argument, 0 if they are equal and 1 if the receiver is greater than the argument. If both versions are orderable, the
// second return value is always true.
func (v godelVersion) CompareTo(o godelVersion) (int, bool) {
	// if either input is not orderable, always return -1 and false
	if !v.Orderable() || !o.Orderable() {
		return -1, false
	}

	switch {
	case v.MajorVersionNum() != o.MajorVersionNum():
		return compareInts(v.MajorVersionNum(), o.MajorVersionNum()), true
	case v.MinorVersionNum() != o.MinorVersionNum():
		return compareInts(v.MinorVersionNum(), o.MinorVersionNum()), true
	case v.PatchVersionNum() != o.PatchVersionNum():
		return compareInts(v.PatchVersionNum(), o.PatchVersionNum()), true
	case (v.Type() == ReleaseSnapshot && o.Type() == Release) || (v.Type() == Release && o.Type() == ReleaseSnapshot):
		if v.Type() == ReleaseSnapshot {
			return 1, true
		}
		return -1, true
	case (v.Type() == Release && o.Type() == ReleaseCandidate) || (v.Type() == ReleaseCandidate && o.Type() == Release):
		if v.Type() == Release {
			return 1, true
		}
		return -1, true
	case v.Type() == ReleaseSnapshot && o.Type() == ReleaseSnapshot:
		return compareInts(*v.FirstSequenceVersionNum(), *o.FirstSequenceVersionNum()), true
	case (v.Type() == ReleaseCandidate || v.Type() == ReleaseCandidateSnapshot) && (o.Type() == ReleaseCandidate || o.Type() == ReleaseCandidateSnapshot):
		if cmp := compareInts(*v.FirstSequenceVersionNum(), *o.FirstSequenceVersionNum()); cmp != 0 {
			return cmp, true
		}

		if v.Type() != o.Type() {
			if v.Type() == ReleaseCandidateSnapshot {
				return 1, true
			}
			return -1, true
		}

		if v.Type() == ReleaseCandidateSnapshot {
			return compareInts(*v.SecondSequenceVersionNum(), *o.SecondSequenceVersionNum()), true
		}
	}
	return 0, true
}

func compareInts(val, other int) int {
	switch {
	default:
		return 0
	case val < other:
		return -1
	case val > other:
		return 1
	}
}

func mustAtoI(in string) int {
	out, err := strconv.Atoi(in)
	if err != nil {
		panic(fmt.Sprintf("invalid version string: %s must be parsable as an int", in))
	}
	return out
}
