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
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	const (
		src = `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
		testSrc = `package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		importsYML = `root-dirs:
  - .`
		licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	)
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args []string
		want string
	}{
		{want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\timports --verify\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\ttest`},
		{args: []string{"--skip-imports"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tformat -v -l\n\timports --verify\n\tlicense --verify\n\tcheck`},
	} {
		cmd := exec.Command("./godelw", append([]string{"verify", "--apply=false"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)
	}
}

func TestVerifyApply(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	const (
		src = `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`
		testSrc = `package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		formattedTestSrc = `package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`
		licensedTestSrc = `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`
		licensedAndFormattedTestSrc = `/*
Copyright 2016 Palantir Technologies, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`
		importsYML = `root-dirs:
  - .`
		importsJSON = `{
    "imports": [],
    "mainOnlyImports": [],
    "testOnlyImports": []
}`

		licenseYML = `header: |
  /*
  Copyright 2016 Palantir Technologies, Inc.

  Licensed under the Apache License, Version 2.0 (the "License");
  you may not use this file except in compliance with the License.
  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

  Unless required by applicable law or agreed to in writing, software
  distributed under the License is distributed on an "AS IS" BASIS,
  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
  See the License for the specific language governing permissions and
  limitations under the License.
  */
`
	)
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "imports.yml"), []byte(importsYML), 0644)
	require.NoError(t, err)
	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args            []string
		want            string
		wantTestSrc     string
		wantImportsJSON string
	}{
		{want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-imports"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: formattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tcheck`, wantTestSrc: licensedAndFormattedTestSrc, wantImportsJSON: importsJSON},
	} {
		err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
		require.NoError(t, err, "Case %d", i)

		cmd := exec.Command("./godelw", append([]string{"verify"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)

		bytes, err := ioutil.ReadFile(path.Join(testProjectDir, "main_test.go"))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantTestSrc, string(bytes), "Case %d", i)

		importsJSONPath := path.Join(testProjectDir, "gocd_imports.json")
		if currCase.wantImportsJSON == "" {
			_, err = os.Stat(importsJSONPath)
			assert.True(t, os.IsNotExist(err), "Case %d: gocd_imports.json should not exist", i)
		} else {
			bytes, err = ioutil.ReadFile(importsJSONPath)
			require.NoError(t, err, "Case %d", i)
			assert.Equal(t, currCase.wantImportsJSON, string(bytes), "Case %d", i)
			err = os.Remove(importsJSONPath)
			require.NoError(t, err, "Case %d", i)
		}
	}
}

func TestVerifyWithJUnitOutput(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	src := `package main
	import "fmt"
	func main() {
		fmt.Println("hello, world!")
	}`
	err := ioutil.WriteFile(path.Join(testProjectDir, "main.go"), []byte(src), 0644)
	require.NoError(t, err)
	testSrc := `package main_test
	import "testing"
	func TestFoo(t *testing.T) {
	}`
	err = ioutil.WriteFile(path.Join(testProjectDir, "main_test.go"), []byte(testSrc), 0644)
	require.NoError(t, err)

	junitOutputFile := "test-output.xml"
	cmd := exec.Command("./godelw", "verify", "--apply=false", "--junit-output", junitOutputFile)
	cmd.Dir = testProjectDir
	err = cmd.Run()
	require.Error(t, err)

	fi, err := os.Stat(path.Join(testProjectDir, junitOutputFile))
	require.NoError(t, err)

	assert.False(t, fi.IsDir())
}
