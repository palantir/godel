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
	"time"

	"github.com/nmiyake/pkg/gofiles"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	licenseYML = `
header: |
  // Copyright (c) {{YEAR}} Palantir Technologies Inc. All rights reserved.
  // Use of this source code is governed by the Apache License, Version 2.0
  // that can be found in the LICENSE file.
`
)

func TestVerify(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "main.go",
			Src: `package main
import "fmt"

func main() {
	fmt.Println("hello, world!")
}`,
		},
		{
			RelPath: "main_test.go",
			Src: `package main_test
import "testing"

func TestFoo(t *testing.T) {
	t=t
	t.Fail()
}`,
		},
	}
	_, err := gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license-plugin.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args []string
		want string
	}{
		{want: `(?s).+Failed tasks:\n\tformat --verify\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\tlicense --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\tformat --verify\n\tlicense --verify\n\ttest`},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tformat --verify\n\tcheck\n\ttest`},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tformat --verify\n\tlicense --verify\n\tcheck`},
	} {
		err = os.MkdirAll(path.Join(testProjectDir, "gen"), 0755)
		require.NoError(t, err)
		err = ioutil.WriteFile(path.Join(testProjectDir, "gen", "output.txt"), []byte("bar-output"), 0644)
		require.NoError(t, err)

		cmd := exec.Command("./godelw", append([]string{"verify", "--apply=false"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err)
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)
	}
}

func TestVerifyApply(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)

	specs := []gofiles.GoFileSpec{
		{
			RelPath: "main.go",
			Src: `package main
	import "fmt"

	func main() {
		fmt.Println("hello, world!")
	}`,
		},
		{
			RelPath: "main_test.go",
			Src: `package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`,
		},
	}

	const (
		formattedTestSrc = `package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`
	)

	var (
		licensedTestSrc = fmt.Sprintf(`// Copyright (c) %d Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main_test
	import "testing"

	func TestFoo(t *testing.T) {
		t=t
		t.Fail()
	}`, time.Now().Year())
		licensedAndFormattedTestSrc = fmt.Sprintf(`// Copyright (c) %d Palantir Technologies Inc. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {
	t = t
	t.Fail()
}
`, time.Now().Year())
	)
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "license-plugin.yml"), []byte(licenseYML), 0644)
	require.NoError(t, err)

	for i, currCase := range []struct {
		args        []string
		want        string
		wantTestSrc string
	}{
		{want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc},
		{args: []string{"--skip-format"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: licensedTestSrc},
		{args: []string{"--skip-check"}, want: `(?s).+Failed tasks:\n\ttest`, wantTestSrc: licensedAndFormattedTestSrc},
		{args: []string{"--skip-license"}, want: `(?s).+Failed tasks:\n\tcheck\n\ttest`, wantTestSrc: formattedTestSrc},
		{args: []string{"--skip-test"}, want: `(?s).+Failed tasks:\n\tcheck`, wantTestSrc: licensedAndFormattedTestSrc},
	} {
		_, err := gofiles.Write(testProjectDir, specs)
		require.NoError(t, err)

		cmd := exec.Command("./godelw", append([]string{"verify"}, currCase.args...)...)
		cmd.Dir = testProjectDir
		output, err := cmd.CombinedOutput()
		require.Error(t, err, fmt.Sprintf("Case %d", i))
		assert.Regexp(t, regexp.MustCompile(currCase.want), string(output), "Case %d", i)

		bytes, err := ioutil.ReadFile(path.Join(testProjectDir, "main_test.go"))
		require.NoError(t, err, "Case %d", i)
		assert.Equal(t, currCase.wantTestSrc, string(bytes), "Case %d", i)
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

func TestVerifyTestTags(t *testing.T) {
	testProjectDir := setUpGödelTestAndDownload(t, testRootDir, gödelTGZ, version)
	specs := []gofiles.GoFileSpec{
		{
			RelPath: "main.go",
			Src: `package main

func main() {}
`,
		},
		{
			RelPath: "main_test.go",
			Src: `package main_test

import (
	"testing"
)

func TestFoo(t *testing.T) {}
`,
		},
		{
			RelPath: "integration_tests/integration_test.go",
			Src: `package main_test

import (
	"testing"
)

func TestFooIntegration(t *testing.T) {}
`,
		},
	}
	files, err := gofiles.Write(testProjectDir, specs)
	require.NoError(t, err)

	err = ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "test-plugin.yml"), []byte(`tags:
  integration:
    names:
      - "integration_tests"
`), 0644)
	require.NoError(t, err)

	// run verify with "none" tags. Should include output for main package but not for integration_test package.
	cmd := exec.Command("./godelw", "verify", "--apply=false", "--tags=none")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	outputStr := string(output)
	require.NoError(t, err, "Command %v failed with error %v. Output: %q", cmd.Args, err, outputStr)
	assert.Regexp(t, fmt.Sprintf(`(?s).+%s\s+[0-9.]+s.+`, files["main.go"].ImportPath), outputStr)
	assert.NotRegexp(t, fmt.Sprintf(`(?s).+%s\s+[0-9.]+s.+`, files["integration_tests/integration_test.go"].ImportPath), outputStr)

	// run verify with "all" tags. Should include output for integration_test package but not for main package.
	cmd = exec.Command("./godelw", "verify", "--apply=false", "--tags=all")
	cmd.Dir = testProjectDir
	output, err = cmd.CombinedOutput()
	outputStr = string(output)
	require.NoError(t, err, "Command %v failed with error %v. Output: %q", cmd.Args, err, outputStr)
	assert.Regexp(t, fmt.Sprintf(`(?s).+%s\s+[0-9.]+s.+`, files["integration_tests/integration_test.go"].ImportPath), outputStr)
	assert.NotRegexp(t, fmt.Sprintf(`(?s).+%s\s+[0-9.]+s.+`, files["main.go"].ImportPath), outputStr)

	// run verify in regular mode. Should include output for all tests.
	cmd = exec.Command("./godelw", "verify", "--apply=false")
	cmd.Dir = testProjectDir
	output, err = cmd.CombinedOutput()
	outputStr = string(output)
	require.NoError(t, err, "Command %v failed with error %v. Output: %q", cmd.Args, err, outputStr)
	assert.Regexp(t, fmt.Sprintf(`(?s).+%s\s+(\(cached\)|[0-9.]+s).+`, files["main.go"].ImportPath), outputStr)
	assert.Regexp(t, fmt.Sprintf(`(?s).+%s\s+(\(cached\)|[0-9.]+s).+`, files["integration_tests/integration_test.go"].ImportPath), outputStr)
}
