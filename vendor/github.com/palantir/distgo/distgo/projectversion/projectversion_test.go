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

package projectversion_test

import (
	"bytes"
	"io/ioutil"
	"path"
	"regexp"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/projectversion"
	"github.com/palantir/distgo/projectversioner/git"
	"github.com/palantir/distgo/projectversioner/script"
)

func TestProjectVersionDefaultParam(t *testing.T) {
	rootDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name  string
		setup func(testDir string)
		want  string
	}{
		{
			"version of project with no tags is 'unspecified'",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
			},
			"^unspecified\n$",
		},
		{
			"version of project tagged with 1.0.0 is 1.0.0",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
				gittest.CreateGitTag(t, testDir, "1.0.0")
			},
			`^` + regexp.QuoteMeta("1.0.0") + `\n$`,
		},
		{
			"version of project with tagged commit with uncommited files ends in .dirty",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
				gittest.CreateGitTag(t, testDir, "1.0.0")
				err := ioutil.WriteFile(path.Join(testDir, "random.txt"), []byte(""), 0644)
				require.NoError(t, err)
			},
			`^` + regexp.QuoteMeta("1.0.0.dirty") + `\n$`,
		},
		{
			"non-tagged commit output",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
				gittest.CreateGitTag(t, testDir, "1.0.0")
				gittest.CommitRandomFile(t, testDir, "Test commit message")
				require.NoError(t, err)
			},
			`^` + regexp.QuoteMeta("1.0.0-1-g") + `[a-f0-9]{7}\n$`,
		},
		{
			"non-tagged commit dirty output",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
				gittest.CreateGitTag(t, testDir, "1.0.0")
				gittest.CommitRandomFile(t, testDir, "Test commit message")
				err := ioutil.WriteFile(path.Join(testDir, "random.txt"), []byte(""), 0644)
				require.NoError(t, err)
			},
			`^` + regexp.QuoteMeta("1.0.0-1-g") + `[a-f0-9]{7}` + regexp.QuoteMeta(`.dirty`) + `\n$`,
		},
	} {
		projectDir, err := ioutil.TempDir(rootDir, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		tc.setup(projectDir)

		projectParam := distgo.ProjectParam{
			ProjectVersionerParam: distgo.ProjectVersionerParam{
				ProjectVersioner: git.New(),
			},
		}
		projectInfo, err := projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		buf := &bytes.Buffer{}
		err = projectversion.Run(projectInfo, buf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Regexp(t, tc.want, buf.String(), "Case %d: %s", i, tc.name)
	}
}

func TestProjectVersionScriptParam(t *testing.T) {
	rootDir, cleanup, err := dirs.TempDir("", "")
	require.NoError(t, err)
	defer cleanup()

	for i, tc := range []struct {
		name         string
		setup        func(testDir string)
		projectParam distgo.ProjectParam
		want         string
	}{
		{
			"project version uses script versioner param if specified",
			func(testDir string) {
				gittest.CommitRandomFile(t, testDir, "Initial commit")
			},
			distgo.ProjectParam{
				ProjectVersionerParam: distgo.ProjectVersionerParam{
					ProjectVersioner: script.New(`#!/usr/bin/env bash
echo "3.2.1"
`,
					),
				},
			},
			`^` + regexp.QuoteMeta("3.2.1") + `\n$`,
		},
	} {
		projectDir, err := ioutil.TempDir(rootDir, "")
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		gittest.InitGitDir(t, projectDir)
		tc.setup(projectDir)

		projectInfo, err := tc.projectParam.ProjectInfo(projectDir)
		require.NoError(t, err, "Case %d: %s", i, tc.name)

		buf := &bytes.Buffer{}
		err = projectversion.Run(projectInfo, buf)
		require.NoError(t, err, "Case %d: %s", i, tc.name)
		assert.Regexp(t, tc.want, buf.String(), "Case %d: %s", i, tc.name)
	}
}
