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

package generator

import (
	"io/ioutil"
	"path"
	"regexp"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

type inputFileWithParsedContent struct {
	FileInfo      inputFile
	ParsedContent parsedTemplateFile
}

func readTemplateFile(inputDir string, f inputFile) (inputFileWithParsedContent, error) {
	fBytes, err := ioutil.ReadFile(path.Join(inputDir, f.TemplateFileName))
	if err != nil {
		return inputFileWithParsedContent{}, errors.Wrapf(err, "failed to read file")
	}
	parsedContent, err := parseTemplateFile(fBytes)
	if err != nil {
		return inputFileWithParsedContent{}, errors.Wrapf(err, "failed to parse template file")
	}
	return inputFileWithParsedContent{
		FileInfo:      f,
		ParsedContent: parsedContent,
	}, nil
}

type parsedTemplateFile struct {
	FullContent       []byte
	TutorialCodeParts []tutorialCodePart
}

type tutorialCodePart struct {
	Code     string
	WantFail bool
}

func (f *parsedTemplateFile) Render(renderedCode []string) ([]byte, error) {
	idx := 0
	var rErr error
	rendered := tutorialCodeRegexp.ReplaceAllFunc(f.FullContent, func([]byte) []byte {
		if idx >= len(renderedCode) {
			rErr = errors.Errorf("index %d is >= than number of rendered code parts %d", len(renderedCode))
			return nil
		}
		curr := strings.Join([]string{
			tutorialCodeStartLineLiteral,
			renderedCode[idx],
			tutorialCodeEndLineLiteral,
		}, "\n")
		idx++
		return []byte(curr)
	})
	if rErr != nil {
		return nil, rErr
	}
	if idx != len(f.TutorialCodeParts) {
		return nil, errors.Errorf("only found %d tutorial code parts in content, but expected %d", idx, len(f.TutorialCodeParts))
	}

	// merge adjacent tutorial code sections
	rendered = adjacentTutorialCodeRegexp.ReplaceAll(rendered, nil)

	// replace start and end tutorial code lines with regular escapes
	rendered = tutorialCodeStartLineRegexp.ReplaceAll(rendered, []byte(codeEscape))
	rendered = tutorialCodeEndLineRegexp.ReplaceAll(rendered, []byte(codeEscape))

	return rendered, nil
}

const (
	tutorialCodeStartLineLiteral = "```START_TUTORIAL_CODE"
	tutorialCodeEndLineLiteral   = "```END_TUTORIAL_CODE"

	codeEscape = "```"
)

var (
	// start line can be followed by optional vertical bar followed by options
	tutorialCodeStartLine = regexp.QuoteMeta(tutorialCodeStartLineLiteral) + "(" + regexp.QuoteMeta("|") + `[^\n]*)?`

	tutorialCodeRegexp          = regexp.MustCompile(`(?sm)^` + tutorialCodeStartLine + `\n(.*?)\n` + regexp.QuoteMeta(tutorialCodeEndLineLiteral) + `$`)
	adjacentTutorialCodeRegexp  = regexp.MustCompile(`(?sm)^` + regexp.QuoteMeta(tutorialCodeEndLineLiteral) + `\n` + tutorialCodeStartLine + `$\n`)
	tutorialCodeStartLineRegexp = regexp.MustCompile(`(?m)^` + tutorialCodeStartLine + `$`)
	tutorialCodeEndLineRegexp   = regexp.MustCompile(`(?m)^` + regexp.QuoteMeta(tutorialCodeEndLineLiteral) + `$`)
)

func parseTemplateFile(fBytes []byte) (parsedTemplateFile, error) {
	out := parsedTemplateFile{
		FullContent: fBytes,
	}
	for _, subMatch := range tutorialCodeRegexp.FindAllSubmatch(fBytes, -1) {
		currPart := tutorialCodePart{
			Code: string(subMatch[2]),
		}
		if optionsMatch := subMatch[1]; len(optionsMatch) > 0 {
			// trim leading "|"
			if options := string(optionsMatch[1:]); len(options) > 0 {
				if err := applyOptions(&currPart, options); err != nil {
					return parsedTemplateFile{}, errors.Wrapf(err, "failed to apply options")
				}
			}
		}
		out.TutorialCodeParts = append(out.TutorialCodeParts, currPart)
	}
	return out, nil
}

func applyOptions(codePart *tutorialCodePart, optionsStr string) error {
	optionsParts := strings.Split(optionsStr, ",")
	for _, currOption := range optionsParts {
		handled := false
		if strings.HasPrefix(currOption, "fail=") {
			failVal := strings.TrimPrefix(currOption, "fail=")
			failValBool, err := strconv.ParseBool(failVal)
			if err != nil {
				return errors.Wrapf(err, "failed to parse option %q", currOption)
			}
			codePart.WantFail = failValBool
			handled = true
		}
		if !handled {
			return errors.Errorf("unknown option: %q", currOption)
		}
	}
	return nil
}
