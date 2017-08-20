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

package publish

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/godel/apps/distgo/cmd/artifacts"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/config"
	"github.com/palantir/godel/apps/distgo/pkg/git/gittest"
)

const testMain = "package main; func main(){}"

func TestPublishBatchErrors(t *testing.T) {
	var handlerFunc func(w http.ResponseWriter, r *http.Request)
	handlerFuncPtr := &handlerFunc
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localHandlerFunc := *handlerFuncPtr
		localHandlerFunc(w, r)
	}))
	defer ts.Close()

	tmp, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	wd, err := os.Getwd()
	defer func() {
		if err := os.Chdir(wd); err != nil {
			fmt.Printf("Failed to restore working directory to %s: %v\n", wd, err)
		}
	}()
	require.NoError(t, err)

	for i, currCase := range []struct {
		cfg                 string
		mainFiles           []string
		publishProducts     []string
		handler             func(w http.ResponseWriter, r *http.Request)
		failFast            bool
		wantOutputRegexp    string
		notWantOutputRegexp string
		wantErrorRegexps    []string
	}{
		// if failFast is false, all products should attempt publish
		{
			cfg: `
products:
  test-bar:
    build:
      main-pkg: ./bar
  test-baz:
    build:
      main-pkg: ./baz
  test-foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go", "bar/main.go", "baz/main.go"},
			publishProducts: []string{"test-foo", "test-bar", "test-baz"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				status := http.StatusOK
				if strings.Contains(r.URL.String(), "test-bar") || strings.Contains(r.URL.String(), "test-baz") {
					// fail all test-bar and test-baz publishes
					status = http.StatusNotFound
				}
				w.WriteHeader(status)
			},
			wantOutputRegexp:    `(?s).+Uploading dist/test-foo-unspecified.pom to .+`,
			notWantOutputRegexp: `(?s).+Uploading dist/test-bar-unspecified.pom to .+`,
			wantErrorRegexps:    []string{`Publish failed for test-bar: uploading .+ to .+ resulted in response "404 Not Found"`, `Publish failed for test-baz: uploading .+ to .+ resulted in response "404 Not Found"`},
		},
		// if failFast is true, first fail should terminate publish
		{
			cfg: `
products:
  test-bar:
    build:
      main-pkg: ./bar
  test-baz:
    build:
      main-pkg: ./baz
  test-foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go", "bar/main.go", "baz/main.go"},
			publishProducts: []string{"test-foo", "test-bar", "test-baz"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				status := http.StatusOK
				if strings.Contains(r.URL.String(), "test-bar") || strings.Contains(r.URL.String(), "test-baz") {
					// fail all test-bar and test-baz publishes
					status = http.StatusNotFound
				}
				w.WriteHeader(status)
			},
			failFast:            true,
			notWantOutputRegexp: `(?s).+Uploading dist/test-bar-unspecified.tgz to .+`,
			wantErrorRegexps:    []string{`^Publish failed for test-bar: uploading .+ to .+ resulted in response "404 Not Found"$`},
		},
	} {
		err = os.Chdir(wd)
		require.NoError(t, err)

		handlerFunc = currCase.handler

		currTmp, err := ioutil.TempDir(tmp, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmp)

		for _, currMain := range currCase.mainFiles {
			err = os.MkdirAll(path.Dir(path.Join(currTmp, currMain)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(currTmp, currMain), []byte(testMain), 0644)
			require.NoError(t, err)
		}

		err = ioutil.WriteFile(path.Join(currTmp, "dist.yml"), []byte(currCase.cfg), 0644)
		require.NoError(t, err)

		cfgcli.ConfigPath = "dist.yml"

		err = os.Chdir(currTmp)
		require.NoError(t, err)

		p := &ArtifactoryConnectionInfo{
			BasicConnectionInfo: BasicConnectionInfo{
				URL:      ts.URL,
				Username: "username",
				Password: "password",
			},
			Repository: "repo",
		}

		buf := &bytes.Buffer{}
		err = publishAction(p, currCase.publishProducts, nil, currCase.failFast, buf, ".")
		assert.Regexp(t, regexp.MustCompile(currCase.wantOutputRegexp), buf.String(), "Case %d", i)
		assert.NotRegexp(t, regexp.MustCompile(currCase.notWantOutputRegexp), buf.String(), "Case %d", i)
		for _, currWantRegexp := range currCase.wantErrorRegexps {
			assert.Regexp(t, regexp.MustCompile(currWantRegexp), err.Error(), "Case %d", i)
		}
	}
}

func TestArtifactoryPublishChecksums(t *testing.T) {
	var handlerFunc func(w http.ResponseWriter, r *http.Request)
	handlerFuncPtr := &handlerFunc
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localHandlerFunc := *handlerFuncPtr
		localHandlerFunc(w, r)
	}))
	defer ts.Close()

	tmpDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	wd, err := os.Getwd()
	defer func() {
		if err := os.Chdir(wd); err != nil {
			fmt.Printf("Failed to restore working directory to %s: %v\n", wd, err)
		}
	}()
	require.NoError(t, err)

	for i, currCase := range []struct {
		cfg             string
		mainFiles       []string
		publishProducts []string
		handler         func(productToArtifact map[string]string) func(w http.ResponseWriter, r *http.Request)
		wantRegexp      func(productToArtifact map[string]string) string
		notWantRegexp   func(productToArtifact map[string]string) string
	}{
		// upload for products with matching checksums are skipped
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(productToArtifact map[string]string) func(w http.ResponseWriter, r *http.Request) {
				fooArtifactPath := productToArtifact["foo"]
				fooArtifactName := path.Base(fooArtifactPath)
				return func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.String(), fooArtifactName) {
						fileInfo, err := newFileInfo(fooArtifactPath)
						require.NoError(t, err)
						bytes, err := json.Marshal(map[string]checksums{
							"checksums": fileInfo.checksums,
						})
						require.NoError(t, err)
						_, err = w.Write(bytes)
						require.NoError(t, err)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			},
			wantRegexp: func(productToArtifact map[string]string) string {
				fooArtifactName := path.Base(productToArtifact["foo"])
				return fmt.Sprintf("File dist/%s already exists at .+, skipping upload", fooArtifactName)
			},
		},
		// upload for products is skipped even if only single checksum matches
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(productToArtifact map[string]string) func(w http.ResponseWriter, r *http.Request) {
				fooArtifactPath := productToArtifact["foo"]
				fooArtifactName := path.Base(fooArtifactPath)
				return func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.String(), fooArtifactName) {
						fileInfo, err := newFileInfo(fooArtifactPath)
						require.NoError(t, err)
						hashes := fileInfo.checksums
						// set 2 hashes to blank, but keep one matching hash
						hashes.SHA256 = ""
						hashes.MD5 = ""
						bytes, err := json.Marshal(map[string]checksums{
							"checksums": hashes,
						})
						require.NoError(t, err)
						_, err = w.Write(bytes)
						require.NoError(t, err)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			},
			wantRegexp: func(productToArtifact map[string]string) string {
				fooArtifactName := path.Base(productToArtifact["foo"])
				return fmt.Sprintf("File dist/%s already exists at .+, skipping upload", fooArtifactName)
			},
		},
		// product is uploaded if checksum does not match
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(productToArtifact map[string]string) func(w http.ResponseWriter, r *http.Request) {
				fooArtifactPath := productToArtifact["foo"]
				fooArtifactName := path.Base(fooArtifactPath)
				return func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.String(), fooArtifactName) {
						fileInfo, err := newFileInfo(fooArtifactPath)
						require.NoError(t, err)
						hashes := fileInfo.checksums
						// 2 hashes match, but one does not
						hashes.SHA1 = "invalid"
						bytes, err := json.Marshal(map[string]checksums{
							"checksums": hashes,
						})
						require.NoError(t, err)
						_, err = w.Write(bytes)
						require.NoError(t, err)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			},
			notWantRegexp: func(productToArtifact map[string]string) string {
				return "File .+ already exists at .+, skipping upload"
			},
		},
		// product is uploaded if response does not contain checksums
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(productToArtifact map[string]string) func(w http.ResponseWriter, r *http.Request) {
				fooArtifactPath := productToArtifact["foo"]
				fooArtifactName := path.Base(fooArtifactPath)
				return func(w http.ResponseWriter, r *http.Request) {
					if strings.Contains(r.URL.String(), fooArtifactName) {
						bytes, err := json.Marshal(map[string]string{
							"no-checksum": "placeholder",
						})
						require.NoError(t, err)
						_, err = w.Write(bytes)
						require.NoError(t, err)
					} else {
						w.WriteHeader(http.StatusOK)
					}
				}
			},
			notWantRegexp: func(productToArtifact map[string]string) string {
				return "File .+ already exists at .+, skipping upload"
			},
		},
	} {
		err = os.Chdir(wd)
		require.NoError(t, err)

		currTmp, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		currTmp, err = filepath.Abs(currTmp)
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmp)

		for _, currMain := range currCase.mainFiles {
			err = os.MkdirAll(path.Dir(path.Join(currTmp, currMain)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(currTmp, currMain), []byte(testMain), 0644)
			require.NoError(t, err)
		}

		err = ioutil.WriteFile(path.Join(currTmp, "dist.yml"), []byte(currCase.cfg), 0644)
		require.NoError(t, err)

		cfgcli.ConfigPath = "dist.yml"

		err = os.Chdir(currTmp)
		require.NoError(t, err)

		cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
		require.NoError(t, err)

		buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, currCase.publishProducts, currTmp)
		require.NoError(t, err)

		artifacts, err := artifacts.DistArtifacts(buildSpecsWithDeps, true)
		require.NoError(t, err)

		productToArtifact := make(map[string]string)
		for k, v := range artifacts {
			require.Equal(t, 1, len(v.Keys()))
			productToArtifact[k] = v.Get(v.Keys()[0])[0]
		}

		handlerFunc = currCase.handler(productToArtifact)

		p := &ArtifactoryConnectionInfo{
			BasicConnectionInfo: BasicConnectionInfo{
				URL:      ts.URL,
				Username: "username",
				Password: "password",
			},
			Repository: "repo",
		}

		buf := &bytes.Buffer{}

		err = publishAction(p, currCase.publishProducts, nil, true, buf, ".")
		require.NoError(t, err, "Case %d", i)

		if currCase.wantRegexp != nil {
			assert.Regexp(t, regexp.MustCompile(currCase.wantRegexp(productToArtifact)), buf.String(), "Case %d", i)
		}

		if currCase.notWantRegexp != nil {
			assert.NotRegexp(t, regexp.MustCompile(currCase.notWantRegexp(productToArtifact)), buf.String(), "Case %d", i)
		}
	}
}

func TestAlmanacPublishCheckURL(t *testing.T) {
	var handlerFunc func(w http.ResponseWriter, r *http.Request)
	handlerFuncPtr := &handlerFunc
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		localHandlerFunc := *handlerFuncPtr
		localHandlerFunc(w, r)
	}))
	defer ts.Close()

	tmpDir, cleanup, err := dirs.TempDir(".", "")
	defer cleanup()
	require.NoError(t, err)

	wd, err := os.Getwd()
	defer func() {
		if err := os.Chdir(wd); err != nil {
			fmt.Printf("Failed to restore working directory to %s: %v\n", wd, err)
		}
	}()
	require.NoError(t, err)

	for i, currCase := range []struct {
		cfg             string
		mainFiles       []string
		publishProducts []string
		handler         func(w http.ResponseWriter, r *http.Request)
		wantRegexp      string
		notWantRegexp   string
		wantErrorRegexp string
	}{
		// Almanac publish for products with matching URLs are skipped
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      dist-type:
        type: sls
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.String() == "/v1/units/foo/unspecified/1" {
					urlMap := map[string]string{
						"url": ts.URL + "/artifactory/repo/com/palantir/distgo-cmd-test/foo/unspecified/foo-unspecified.sls.tgz",
					}
					bytes, err := json.Marshal(urlMap)
					require.NoError(t, err)
					_, err = w.Write(bytes)
					require.NoError(t, err)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			},
			wantRegexp: fmt.Sprintf(`(?s).+Unit for product foo branch unspecified revision 1 with URL %s already exists; skipping publish.+`, ts.URL+"/artifactory/repo/com/palantir/distgo-cmd-test/foo/unspecified/foo-unspecified.sls.tgz"),
		},
		// Almanac publish for product that do not exist in Almanac succeeds
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      dist-type:
        type: sls
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				status := http.StatusOK
				if r.URL.String() == "/v1/units/foo/unspecified/1" {
					// return error for unit
					status = http.StatusBadRequest
				}
				w.WriteHeader(status)
			},
			notWantRegexp: `(?s).+Unit for product .+ branch .+ revision .+ with the URL .+ already exists; skipping publish.+`,
		},
		// Almanac publish for product that exists in Almanac (but has no URL) fails
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      dist-type:
        type: sls
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			},
			notWantRegexp:   `(?s).+Unit for product .+ branch .+ revision .+ with the URL .+ already exists; skipping publish.+`,
			wantErrorRegexp: `^Almanac publish failed for foo: unit for product foo branch unspecified revision 1 already exists; not overwriting it$`,
		},
		// Almanac publish for product that exists in Almanac (but has different URL) fails
		{
			cfg: `
products:
  foo:
    build:
      main-pkg: ./foo
    dist:
      dist-type:
        type: sls
group-id: com.palantir.distgo-cmd-test`,
			mainFiles:       []string{"foo/main.go"},
			publishProducts: []string{"foo"},
			handler: func(w http.ResponseWriter, r *http.Request) {
				if r.URL.String() == "/v1/units/foo/unspecified/1" {
					urlMap := map[string]string{
						"url": "nonMatchingURL",
					}
					bytes, err := json.Marshal(urlMap)
					require.NoError(t, err)
					_, err = w.Write(bytes)
					require.NoError(t, err)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			},
			notWantRegexp:   `(?s).+Unit for product .+ branch .+ revision .+ with the URL .+ already exists; skipping publish.+`,
			wantErrorRegexp: `^Almanac publish failed for foo: unit for product foo branch unspecified revision 1 already exists; not overwriting it$`,
		},
	} {
		err = os.Chdir(wd)
		require.NoError(t, err)

		currTmp, err := ioutil.TempDir(tmpDir, "")
		require.NoError(t, err)

		gittest.InitGitDir(t, currTmp)

		for _, currMain := range currCase.mainFiles {
			err = os.MkdirAll(path.Dir(path.Join(currTmp, currMain)), 0755)
			require.NoError(t, err)
			err = ioutil.WriteFile(path.Join(currTmp, currMain), []byte(testMain), 0644)
			require.NoError(t, err)
		}

		err = ioutil.WriteFile(path.Join(currTmp, "dist.yml"), []byte(currCase.cfg), 0644)
		require.NoError(t, err)

		cfgcli.ConfigPath = "dist.yml"

		err = os.Chdir(currTmp)
		require.NoError(t, err)

		handlerFunc = currCase.handler

		p := &ArtifactoryConnectionInfo{
			BasicConnectionInfo: BasicConnectionInfo{
				URL:      ts.URL,
				Username: "username",
				Password: "password",
			},
			Repository: "repo",
		}

		a := &AlmanacInfo{
			URL:      ts.URL,
			AccessID: "username",
			Secret:   "password",
		}

		buf := &bytes.Buffer{}

		err = publishAction(p, currCase.publishProducts, a, true, buf, ".")

		if currCase.wantErrorRegexp != "" {
			assert.Regexp(t, regexp.MustCompile(currCase.wantErrorRegexp), err.Error(), "Case %d", i)
		} else {
			require.NoError(t, err, "Case %d", i)
		}

		if currCase.wantRegexp != "" {
			assert.Regexp(t, regexp.MustCompile(currCase.wantRegexp), buf.String(), "Case %d", i)
		}

		if currCase.notWantRegexp != "" {
			assert.NotRegexp(t, regexp.MustCompile(currCase.notWantRegexp), buf.String(), "Case %d", i)
		}
	}
}
