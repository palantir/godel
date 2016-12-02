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
	"bufio"
	"io"
	"regexp"
	"strconv"
	"strings"

	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

type IssueParser interface {
	// IsStartToken returns true if the provided line contains a token that signals that it represents a new Issue.
	// If the parser is strictly a line-based one (every line is a new issue), the implementation can always return
	// true.
	IsStartToken(line string) bool

	// ParseSingleIssue parses the provided string and returns the Issue that is parsed from it. The first line of
	// the input must be a line that causes "IsStartToken" to return true. Returns an error if the provided input
	// cannot be parsed as an issue.
	ParseSingleIssue(input string) (Issue, error)
}

type LineParser func(line, rootDir string) (Issue, error)

type SingleLineIssueParser struct {
	LineParser LineParser
	RootDir    string
}

func (p *SingleLineIssueParser) IsStartToken(line string) bool {
	_, err := p.LineParser(line, p.RootDir)
	return err == nil
}

func (p *SingleLineIssueParser) ParseSingleIssue(input string) (Issue, error) {
	return p.LineParser(input, p.RootDir)
}

func ParseIssues(reader io.Reader, parser IssueParser, rawLineFilter func(line string) bool) ([]Issue, error) {
	var issues []Issue

	numLinesRead := 0
	scanner := bufio.NewScanner(reader)

	// read the first line
	atEnd, firstLineOfCurrIssue, err := nextValidLine(scanner, &numLinesRead, rawLineFilter)
	if err != nil {
		return nil, errors.Wrapf(err, "failed at line %d", numLinesRead)
	}

	if atEnd {
		return issues, nil
	}

	for {
		// verify that first line is valid
		if !parser.IsStartToken(firstLineOfCurrIssue) {
			return nil, errors.Errorf("failed on line %d: line %s is not valid as the start token for an issue", numLinesRead, firstLineOfCurrIssue)
		}

		currIssueText := firstLineOfCurrIssue

		// read until first line of next issue or end
		firstLineOfNextIssue := ""
		nextIssueExists := false
		for {
			atEnd, nextLine, err := nextValidLine(scanner, &numLinesRead, rawLineFilter)
			if err != nil {
				return nil, errors.Wrapf(err, "failed at line %d", numLinesRead)
			}

			if atEnd {
				break
			}

			// if the next line is the start of a new issue, break. "firstLineOfNextIssue" contains the line.
			if parser.IsStartToken(nextLine) {
				nextIssueExists = true
				firstLineOfNextIssue = nextLine
				break
			}

			// otherwise, append the current line to the text for the current issue
			currIssueText = currIssueText + "\n" + nextLine
		}

		// parse current issue
		currIssue, err := parser.ParseSingleIssue(currIssueText)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to parse issue from text %s", currIssueText)
		}
		issues = append(issues, currIssue)

		if !nextIssueExists {
			break
		}

		// update
		firstLineOfCurrIssue = firstLineOfNextIssue
	}

	return issues, nil
}

// reads the next line that does not pass the given filter
func nextValidLine(scanner *bufio.Scanner, numLinesRead *int, rawLineFilter func(line string) bool) (bool, string, error) {
	for scanner.Scan() {
		nextLine := scanner.Text()
		(*numLinesRead)++

		// passes, return
		if rawLineFilter == nil || !rawLineFilter(nextLine) {
			return false, nextLine, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return true, "", errors.Wrapf(err, "failed reading line %d", numLinesRead)
	}

	// reached end of buffer but no error
	return true, "", nil
}

func DefaultParser(pathType pkgpath.Type) LineParser {
	return func(line, rootDir string) (Issue, error) {
		return parseStandardLine(line, pathType, rootDir, false)
	}
}

func MultiLineParser(pathType pkgpath.Type) LineParser {
	return func(line, rootDir string) (Issue, error) {
		return parseStandardLine(line, pathType, rootDir, true)
	}
}

func StartAfterFirstWhitespaceParser(pathType pkgpath.Type) LineParser {
	// some tools have output of the form "(text):(whitespace)" before the standard output, so provide a parser
	// that skips everything up to after the first chunk of whitespace
	return func(line, rootDir string) (Issue, error) {
		spaceIndex := whitespace.FindStringIndex(line)
		return parseStandardLine(line[spaceIndex[1]:], pathType, rootDir, false)
	}
}

func RawParser() LineParser {
	return func(line, rootDir string) (Issue, error) {
		return Issue{
			message: line,
			baseDir: rootDir,
		}, nil
	}
}

var whitespace = regexp.MustCompile(`\s+`)

func parseStandardLine(line string, pathType pkgpath.Type, rootDir string, strict bool) (Issue, error) {
	spaceIndex := whitespace.FindStringIndex(line)
	if spaceIndex == nil {
		return Issue{}, errors.Errorf("failed to find whitespace in line %s", line)
	}

	filePath, lineNum, columnNum, err := parseStandardLocation(line[0:spaceIndex[0]], strict)
	if err != nil {
		return Issue{}, errors.Wrapf(err, "failed to parse location from line %s", line)
	}

	pkgPather := newPkgPather(filePath, rootDir, pathType)
	if pkgPather == nil {
		return Issue{}, errors.Errorf("failed to create PkgPather for %s", filePath)
	}

	relPath, err := pkgPather.Rel(rootDir)
	if err != nil {
		return Issue{}, errors.WithStack(err)
	}
	relPath = strings.TrimPrefix(relPath, "./")

	messagePart := line[spaceIndex[1]:]
	message := strings.TrimSpace(messagePart)

	return Issue{
		path:    relPath,
		line:    lineNum,
		column:  columnNum,
		message: message,
		baseDir: rootDir,
	}, nil
}

func newPkgPather(path string, rootDir string, pathType pkgpath.Type) pkgpath.PkgPather {
	switch pathType {
	case pkgpath.Absolute:
		return pkgpath.NewAbsPkgPath(path)
	case pkgpath.GoPathSrcRelative:
		return pkgpath.NewGoPathSrcRelPkgPath(path)
	case pkgpath.Relative:
		return pkgpath.NewRelPkgPath(path, rootDir)
	default:
		return nil
	}
}

func parseStandardLocation(locationString string, strict bool) (string, int, int, error) {
	// trim final ":" so split is cleaner
	if strings.HasSuffix(locationString, ":") {
		locationString = locationString[:len(locationString)-1]
	} else if strict {
		return "", 0, 0, errors.Errorf("location input %s did not have suffix ':'", locationString)
	}

	locationParts := strings.Split(locationString, ":")
	if len(locationParts) > 3 {
		// too many parts -- max is "message:line:col", which is 3 parts
		return "", 0, 0, errors.Errorf(`splitting %q on character ':' resulted in greater than 3 parts: %v`, locationString, locationParts)
	} else if strict && len(locationString) < 2 {
		return "", 0, 0, errors.Errorf("splitting %q on character ':' resulted in fewer than 2 parts: %v", locationString, locationParts)
	}

	currIndex := 0

	filePath := locationParts[currIndex]
	currIndex++

	var lineNum int
	var err error
	if currIndex < len(locationParts) {
		lineNum, err = parsePartAsInt(locationParts, currIndex)
		if err != nil {
			return "", 0, 0, errors.Wrapf(err, "failed to parse element %d from parts %v of line %s as an integer", currIndex, locationParts, locationString)
		}
	}
	currIndex++

	var columnNum int
	if currIndex < len(locationParts) {
		columnNum, err = parsePartAsInt(locationParts, currIndex)
		if err != nil {
			return "", 0, 0, errors.Wrapf(err, "failed to parse element %d from parts %v of line %s as an integer", currIndex, locationParts, locationString)
		}
	}

	return filePath, lineNum, columnNum, nil
}

func parsePartAsInt(parts []string, index int) (int, error) {
	numStr := parts[index]
	return strconv.Atoi(strings.TrimSpace(numStr))
}
