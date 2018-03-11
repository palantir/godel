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
	"sort"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBashScript(t *testing.T) {
	for i, tc := range []struct {
		codeParts []tutorialCodePart
		want      string
	}{
		{
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
			`#!/usr/bin/env bash
print_then_run () {
    echo "BASH_RUN:-------------"
    echo "$1"
    echo "----------------------"

    echo "OUTPUT:---------------"
    eval "$1"
    echo "----------------------"
}

set +e
read -d '' ACTION <<"EOF"
mkdir "testDir"
EOF
set -e
print_then_run "$ACTION"

set +e
read -d '' ACTION <<"EOF"
ls -la
EOF
set -e
print_then_run "$ACTION"

set +e
read -d '' ACTION <<"EOF"
echo 'multi-line
content' > foo.txt
EOF
set -e
print_then_run "$ACTION"
`,
		},
	} {
		got := bashScript(tc.codeParts)
		assert.Equal(t, tc.want, got, "Case %d\nOutput:\n%s", i, got)
	}
}

func TestDockerfile(t *testing.T) {
	for i, tc := range []struct {
		fromImage, scriptFileName string
		want                      string
	}{
		{
			"tutorial:step1",
			"step1.sh",
			`FROM tutorial:step1

ADD step1.sh /scripts/
RUN /scripts/step1.sh 2>&1
`,
		},
	} {
		got := dockerFile(tc.fromImage, tc.scriptFileName)
		assert.Equal(t, tc.want, got, "Case %d\nOutput:\n%s", i, got)
	}
}

func TestWriteOutputFiles(t *testing.T) {
	tmpDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		inFile    inputFileWithParsedContent
		fromImage string
		wantFiles map[string]string
	}{
		{
			inputFileWithParsedContent{
				FileInfo: mustNewInputFile("1_add.md.tmpl"),
				ParsedContent: mustParseTemplateFile([]byte(renderLiteral(`Hello, world!

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
`))),
			},
			"tutorial:step1",
			map[string]string{
				"1_add/run-add.sh": `#!/usr/bin/env bash
print_then_run () {
    echo "BASH_RUN:-------------"
    echo "$1"
    echo "----------------------"

    echo "OUTPUT:---------------"
    eval "$1"
    echo "----------------------"
}

set +e
read -d '' ACTION <<"EOF"
mkdir "testDir"
EOF
set -e
print_then_run "$ACTION"

set +e
read -d '' ACTION <<"EOF"
ls -la
EOF
set -e
print_then_run "$ACTION"

set +e
read -d '' ACTION <<"EOF"
echo 'multi-line
content' > foo.txt
EOF
set -e
print_then_run "$ACTION"
`,
				"1_add/Dockerfile": `FROM tutorial:step1

ADD run-add.sh /scripts/
RUN /scripts/run-add.sh 2>&1
`,
			},
		},
		{
			inputFileWithParsedContent{
				FileInfo: mustNewInputFile("1_add.md.tmpl"),
				ParsedContent: mustParseTemplateFile([]byte(`Hello, world!

` + "```START_TUTORIAL_CODE|fail=true" + `
mkdir "testDir"
` + "```END_TUTORIAL_CODE" + `
`)),
			},
			"tutorial:step1",
			map[string]string{
				"1_add/run-add.sh": `#!/usr/bin/env bash
print_then_run () {
    echo "BASH_RUN:-------------"
    echo "$1"
    echo "----------------------"

    echo "OUTPUT:---------------"
    eval "$1"
    echo "----------------------"
}

set +e
read -d '' ACTION <<"EOF"
mkdir "testDir" || true
EOF
set -e
print_then_run "$ACTION"
`,
				"1_add/Dockerfile": `FROM tutorial:step1

ADD run-add.sh /scripts/
RUN /scripts/run-add.sh 2>&1
`,
			},
		},
	} {
		currDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err, "Case %d", i)

		_, err = writeOutputFiles(currDir, tc.inFile, tc.fromImage)
		require.NoError(t, err, "Case %d", i)

		var sortedKeys []string
		for k := range tc.wantFiles {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)

		for _, k := range sortedKeys {
			want := tc.wantFiles[k]

			filePath := path.Join(currDir, k)
			gotBytes, err := ioutil.ReadFile(filePath)
			require.NoError(t, err, "Case %d\nPath: %s", i, filePath)
			got := string(gotBytes)

			assert.Equal(t, want, got, "Case %d\nPath: %s\nOutput:\n%s", i, filePath, got)
		}
	}
}

func TestParseBashRunCmdFromOutput(t *testing.T) {
	for i, tc := range []struct {
		in   string
		want []bashRunCmd
	}{
		{
			`BASH_RUN:-------------
mkdir "testDir"
----------------------
OUTPUT:---------------
----------------------
`,
			[]bashRunCmd{
				{
					cmd:    `mkdir "testDir"`,
					output: ``,
				},
			},
		},
		{
			`BASH_RUN:-------------
ls -la
----------------------
OUTPUT:---------------
total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir
----------------------
`,
			[]bashRunCmd{
				{
					cmd: `ls -la`,
					output: `total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir`,
				},
			},
		},
		{
			`BASH_RUN:-------------
mkdir "testDir"
----------------------
OUTPUT:---------------
----------------------
BASH_RUN:-------------
ls -la
----------------------
OUTPUT:---------------
total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir
----------------------
`,
			[]bashRunCmd{
				{
					cmd:    `mkdir "testDir"`,
					output: ``,
				},
				{
					cmd: `ls -la`,
					output: `total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir`,
				},
			},
		},
		{
			`BASH_RUN:-------------
mkdir "testDir"
----------------------
OUTPUT:---------------
----------------------
BASH_RUN:-------------
ls -la
----------------------
OUTPUT:---------------
total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir
----------------------
BASH_RUN:-------------
echo 'multi-line
content' > foo.txt
----------------------
OUTPUT:---------------
----------------------
`,
			[]bashRunCmd{
				{
					cmd:    `mkdir "testDir"`,
					output: ``,
				},
				{
					cmd: `ls -la`,
					output: `total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir`,
				},
				{
					cmd: `echo 'multi-line
content' > foo.txt`,
					output: ``,
				},
			},
		},
	} {
		got := parseBashRunCmdFromOutput(tc.in)
		assert.Equal(t, tc.want, got, "Case %d\nOutput:\n%s", i, got)
	}
}

func TestBashRunCmdString(t *testing.T) {
	for i, tc := range []struct {
		cmd  bashRunCmd
		want string
	}{
		{
			cmd: bashRunCmd{
				cmd:    `mkdir "testDir"`,
				output: ``,
			},
			want: `➜ mkdir "testDir"`,
		},
		{
			cmd: bashRunCmd{
				cmd: `echo 'multi-line
content' > foo.txt`,
				output: ``,
			},
			want: `➜ echo 'multi-line
content' > foo.txt`,
		},
		{
			cmd: bashRunCmd{
				cmd: `ls -la`,
				output: `total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir`,
			},
			want: `➜ ls -la
total 16
drwxr-xr-x  5 test  170 Mar 28 13:45 .
drwx------  3 test  102 Mar 28 13:45 ..
-rw-r--r--  1 test  62 Mar 28 13:45 Dockerfile
-rwxr-xr-x  1 test  476 Mar 28 13:45 add.sh
drwxr-xr-x  2 test  68 Mar 28 13:45 testDir`,
		},
	} {
		got := tc.cmd.String()
		assert.Equal(t, tc.want, got, "Case %d\nOutput:\n%s", i, got)
	}
}

func mustNewInputFile(fileName string) inputFile {
	out, err := newInputFile(fileName)
	if err != nil {
		panic(err)
	}
	return out
}

func mustParseTemplateFile(in []byte) parsedTemplateFile {
	out, err := parseTemplateFile(in)
	if err != nil {
		panic(err)
	}
	return out
}
