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

package integration

import (
	"io/ioutil"
	"path"
	"regexp"
	"testing"

	"github.com/palantir/godel/framework/pluginapitester"
	"github.com/palantir/godel/pkg/products"
	"github.com/palantir/pkg/gittest"
	"github.com/stretchr/testify/require"

	"github.com/palantir/distgo/projectversioner/projectversiontester"
)

func TestGitProjectVersioner(t *testing.T) {
	const godelYML = `exclude:
  names:
    - "\\..+"
    - "vendor"
  paths:
    - "godel"
`

	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	projectversiontester.RunAssetProjectVersionTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]projectversiontester.TestCase{
			{
				Name: "version of project with no tags is 'unspecified'",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Second commit")
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile("^unspecified\n$")
				},
			},
			{
				Name: "version of project tagged with 1.0.0 is 1.0.0",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Second commit")
					gittest.CreateGitTag(t, testDir, "1.0.0")
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile(`^` + regexp.QuoteMeta("1.0.0") + `\n$`)
				},
			},
			{
				Name: "version of project tagged with v1.0.0 is 1.0.0",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Second commit")
					gittest.CreateGitTag(t, testDir, "v1.0.0")
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile(`^` + regexp.QuoteMeta("1.0.0") + `\n$`)
				},
			},
			{
				Name: "version of project with tagged commit with uncommited files ends in .dirty",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Initial commit")
					gittest.CreateGitTag(t, testDir, "1.0.0")
					err := ioutil.WriteFile(path.Join(testDir, "random.txt"), []byte(""), 0644)
					require.NoError(t, err)
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile(`^` + regexp.QuoteMeta("1.0.0.dirty") + `\n$`)
				},
			},
			{
				Name: "non-tagged commit output",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Initial commit")
					gittest.CreateGitTag(t, testDir, "1.0.0")
					gittest.CommitRandomFile(t, testDir, "Test commit message")
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile(`^` + regexp.QuoteMeta("1.0.0") + `-1-g[a-f0-9]{7}\n$`)
				},
			},
			{
				Name: "non-tagged commit dirty output",
				ConfigFiles: map[string]string{
					"godel/config/godel.yml": godelYML,
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
`,
				},
				Setup: func(testDir string) {
					gittest.CommitRandomFile(t, testDir, "Initial commit")
					gittest.CreateGitTag(t, testDir, "1.0.0")
					gittest.CommitRandomFile(t, testDir, "Test commit message")
					err := ioutil.WriteFile(path.Join(testDir, "random.txt"), []byte(""), 0644)
					require.NoError(t, err)
				},
				WantOutput: func(projectDir string) *regexp.Regexp {
					return regexp.MustCompile(`^` + regexp.QuoteMeta("1.0.0") + `-1-g[a-f0-9]{7}` + regexp.QuoteMeta(".dirty") + `\n$`)
				},
			},
		},
	)
}

func TestGitUpgradeConfig(t *testing.T) {
	pluginPath, err := products.Bin("dist-plugin")
	require.NoError(t, err)

	pluginapitester.RunUpgradeConfigTest(t,
		pluginapitester.NewPluginProvider(pluginPath),
		nil,
		[]pluginapitester.UpgradeConfigTestCase{
			{
				Name: `valid v0 config works`,
				ConfigFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
  config:
    # comment
`,
				},
				WantOutput: ``,
				WantFiles: map[string]string{
					"godel/config/dist-plugin.yml": `
project-versioner:
  type: git
  config:
    # comment
`,
				},
			},
		},
	)
}
