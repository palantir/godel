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
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path"
	"strings"
	"testing"

	"github.com/palantir/godel/v2/framework/builtintasks/installupdate/layout"
	"github.com/stretchr/testify/require"
)

func setUpGodelTestAndDownload(t *testing.T, testRootDir, godelTGZ string, version string) string {
	testProjectDir, server := setUpGodelTest(t, testRootDir, godelTGZ, version)
	defer server.Close()

	cmd := exec.Command("./godelw", "--version")
	cmd.Dir = testProjectDir
	output, err := cmd.CombinedOutput()
	require.NoError(t, err, "Command %v failed. Output:\n%v", cmd.Args, string(output))

	return testProjectDir
}

func setUpGodelTest(t *testing.T, testRootDir, godelTGZ, version string) (string, *httptest.Server) {
	testProjectDir, err := ioutil.TempDir(testRootDir, "")
	require.NoError(t, err)

	installGodel(t, testProjectDir, godelTGZ, version)
	server := createTGZServer(t, godelTGZ)
	updateGodelProperties(t, testProjectDir, server.URL)

	return testProjectDir, server
}

func createTGZServer(t *testing.T, godelTGZ string) *httptest.Server {
	_, err := os.Stat(godelTGZ)
	require.NoError(t, err)

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bytes, err := ioutil.ReadFile(godelTGZ)
		require.NoError(t, err)
		_, err = w.Write(bytes)
		require.NoError(t, err)
	}))
	return ts
}

func installGodel(t *testing.T, testProjectDir, godelTGZ, version string) {
	specDir, err := layout.AppSpecDir(strings.TrimSuffix(godelTGZ, ".tgz"), version)
	require.NoError(t, err)

	err = layout.CopyFile(specDir.Path(layout.WrapperScriptFile), path.Join(testProjectDir, "godelw"))
	require.NoError(t, err)
	err = layout.CopyDir(specDir.Path(layout.WrapperAppDir), path.Join(testProjectDir, "godel"))
	require.NoError(t, err)
}

func updateGodelProperties(t *testing.T, testProjectDir, url string) {
	contents := fmt.Sprintf("distributionURL=%v\n", url)
	err := ioutil.WriteFile(path.Join(testProjectDir, "godel", "config", "godel.properties"), []byte(contents), 0644)
	require.NoError(t, err)
}
