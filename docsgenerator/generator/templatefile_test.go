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
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseAndRenderTemplateFile(t *testing.T) {
	for i, tc := range []struct {
		name                  string
		in                    string
		wantTutorialCodeParts []tutorialCodePart
		wantRendered          string
	}{
		{
			"Adjacent dividers are merged",
			renderLiteral(`Hello, world!

Here's an example:

{{START_DIVIDER}}
mkdir "testDir"
{{END_DIVIDER}}
{{START_DIVIDER}}
ls -la
{{END_DIVIDER}}

Another line.

{{START_DIVIDER}}
echo 'multi-line
content' > foo.txt
{{END_DIVIDER}}
`),
			[]tutorialCodePart{
				{
					Code: `mkdir "testDir"`,
				},
				{
					Code: `ls -la`,
				},
				{
					Code: `echo 'multi-line
content' > foo.txt`,
				},
			},
			`Hello, world!

Here's an example:

` + "```" + `
mkdir "testDir"
ls -la
` + "```" + `

Another line.

` + "```" + `
echo 'multi-line
content' > foo.txt
` + "```" + `
`,
		},
		{
			"Options after header are parsed properly",
			`Hello, world!

` + "```START_TUTORIAL_CODE|fail=true" + `
mkdir "testDir"
` + "```END_TUTORIAL_CODE" + `
`,
			[]tutorialCodePart{
				{
					Code:     `mkdir "testDir"`,
					WantFail: true,
				},
			},
			`Hello, world!

` + "```" + `
mkdir "testDir"
` + "```" + `
`,
		},
	} {
		got, err := parseTemplateFile([]byte(tc.in))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.wantTutorialCodeParts, got.TutorialCodeParts, "Case %d", i)

		var rawCode []string
		for _, currPart := range got.TutorialCodeParts {
			rawCode = append(rawCode, currPart.Code)
		}
		gotRendered, err := got.Render(rawCode)
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, tc.wantRendered, string(gotRendered), "Case %d\nOutput:\n%s", i, string(gotRendered))
	}
}

func renderLiteral(in string) string {
	out := in
	out = strings.Replace(out, "{{START_DIVIDER}}", tutorialCodeStartLineLiteral, -1)
	out = strings.Replace(out, "{{END_DIVIDER}}", tutorialCodeEndLineLiteral, -1)
	return out
}
