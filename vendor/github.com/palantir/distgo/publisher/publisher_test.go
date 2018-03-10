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

package publisher_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/dist"
	"github.com/palantir/distgo/distgo/publish"
	"github.com/palantir/distgo/dockerbuilder"
)

type publisherTestCase struct {
	name            string
	projectCfg      distgo.ProjectConfig
	wantOutput      func(projectDir string) string
	wantErrorRegexp string
}

func runPublisherTests(t *testing.T, publisherImpl distgo.Publisher, dryRun bool, testCases ...publisherTestCase) {
	tmp, cleanup, err := dirs.TempDir("", "")
	defer cleanup()
	require.NoError(t, err)

	for i, tc := range testCases {
		projectDir, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		err = os.MkdirAll(path.Join(projectDir, "foo"), 0755)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		err = ioutil.WriteFile(path.Join(projectDir, "foo", "main.go"), []byte("package main; func main(){}"), 0644)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		err = ioutil.WriteFile(path.Join(projectDir, ".gitignore"), []byte("/out\n"), 0644)
		gittest.CommitAllFiles(t, projectDir, "Initial commit")
		gittest.CreateGitTag(t, projectDir, "1.0.0")

		disterFactory, err := dister.NewDisterFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		defaultDistCfg, err := dister.DefaultConfig()
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		dockerBuilderFactory, err := dockerbuilder.NewDockerBuilderFactory()
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectParam, err := tc.projectCfg.ToParam(projectDir, disterFactory, defaultDistCfg, dockerBuilderFactory)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		projectInfo, err := projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		// run "dist" to ensure that dist outputs exist
		output := &bytes.Buffer{}
		err = dist.Products(projectInfo, projectParam, nil, false, output)
		require.NoError(t, err, "Case %d: %s\nOutput: %s", i, tc.name, output.String())

		output = &bytes.Buffer{}
		err = publish.Products(projectInfo, projectParam, nil, publisherImpl, nil, dryRun, output)
		if tc.wantErrorRegexp == "" {
			require.NoError(t, err, "Case %d: %s", i, tc.name)
			assert.Equal(t, tc.wantOutput(projectDir), output.String(), "Case %d: %s", i, tc.name)
		} else {
			require.Error(t, err, fmt.Sprintf("Case %d: %s", i, tc.name))
			assert.Regexp(t, regexp.MustCompile(tc.wantErrorRegexp), err.Error(), "Case %d: %s", i, tc.name)
		}
	}
}

func stringPtr(in string) *string {
	return &in
}

func mustMapSlicePtr(in interface{}) *yaml.MapSlice {
	out, err := distgo.ToMapSlice(in)
	if err != nil {
		panic(err)
	}
	return &out
}

func mustOSArch(in string) osarch.OSArch {
	osArch, err := osarch.New(in)
	if err != nil {
		panic(err)
	}
	return osArch
}
