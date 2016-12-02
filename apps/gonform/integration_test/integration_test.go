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

package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/pkg/products"
)

func TestRun(t *testing.T) {
	cli, err := products.Bin("gonform")
	require.NoError(t, err)

	wd, err := os.Getwd()
	require.NoError(t, err)

	tmpDir, cleanup, err := dirs.TempDir(wd, "")
	defer cleanup()
	require.NoError(t, err)

	for i, currCase := range []struct {
		filesToCreate []gofiles.GoFileSpec
		filesToCheck  []string
		config        string
		want          []gofiles.GoFileSpec
		wantStdout    string
	}{
		{
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: unindent(`package foo

							import "fmt"

					func main() {
									fmt.Println("indented way too far")
					}
					`),
				},
				{
					RelPath: "vendor/foo/foo.go",
					Src: unindent(`package foo

							import "fmt"

					func main() {
									fmt.Println("indented way too far")
					}
					`),
				},
			},
			config: unindent(`
					exclude:
					  paths:
					    - "vendor"
					`),
			want: []gofiles.GoFileSpec{
				{
					RelPath: "foo.go",
					Src: unindent(`package foo

					import (
						"fmt"
					)

					func main() {
						fmt.Println("indented way too far")
					}
					`),
				},
				{
					RelPath: "vendor/foo/foo.go",
					Src: unindent(`package foo

							import "fmt"

					func main() {
									fmt.Println("indented way too far")
					}
					`),
				},
			},
			wantStdout: "Running ptimports...\n",
		},
		{
			filesToCreate: []gofiles.GoFileSpec{
				{
					RelPath: "vendor/foo/foo.go",
					Src: unindent(`package foo

							import "fmt"

					func main() {
									fmt.Println("indented way too far")
					}
					`),
				},
			},
			config: unindent(`
					exclude:
					  paths:
					    - "vendor"
					`),
			want: []gofiles.GoFileSpec{
				{
					RelPath: "vendor/foo/foo.go",
					Src: unindent(`package foo

							import "fmt"

					func main() {
									fmt.Println("indented way too far")
					}
					`),
				},
			},
		},
	} {
		currCaseTmpDir, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		_, err = gofiles.Write(currCaseTmpDir, currCase.filesToCreate)
		require.NoError(t, err, "Case %d", i)

		configFile := path.Join(currCaseTmpDir, "config.yml")
		err = ioutil.WriteFile(configFile, []byte(currCase.config), 0644)
		require.NoError(t, err)

		err = os.Chdir(currCaseTmpDir)
		require.NoError(t, err)

		cmd := exec.Command(cli, "--config", configFile)
		err = cmd.Run()
		require.NoError(t, err, "Case %d", i)

		verifyFileContent(t, i, currCaseTmpDir, currCase.want)
	}
}

func verifyFileContent(t *testing.T, caseNum int, rootDir string, expected []gofiles.GoFileSpec) {
	for _, spec := range expected {
		bytes, err := ioutil.ReadFile(path.Join(rootDir, spec.RelPath))
		require.NoError(t, err, "Case %d", caseNum)

		actualContent := string(bytes)
		assert.Equal(t, spec.Src, actualContent, "Case %d", caseNum)
	}
}

func unindent(input string) string {
	return strings.Replace(input, "\n\t\t\t\t\t", "\n", -1)
}
