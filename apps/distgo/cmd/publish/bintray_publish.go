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
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

type BintrayConnectionInfo struct {
	BasicConnectionInfo
	Subject       string
	Repository    string
	Release       bool
	DownloadsList bool
}

func (b *BintrayConnectionInfo) Publish(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) ([]string, error) {
	baseURL := strings.Join([]string{b.URL, "content", b.Subject, b.Repository, buildSpec.ProductName, buildSpec.ProductVersion, paths.productPath}, "/")
	artifactURLs, err := b.uploadArtifacts(baseURL, paths, nil, stdout)
	if err != nil {
		return artifactURLs, err
	}
	if b.Release {
		if err := b.release(buildSpec, stdout); err != nil {
			fmt.Fprintln(stdout, "Uploading artifacts succeeded, but publish of uploaded artifacts failed:", err)
		}
	}
	if b.DownloadsList {
		if err := b.addToDownloadsList(buildSpec, paths, stdout); err != nil {
			fmt.Fprintln(stdout, "Uploading artifacts succeeded, but addings artifact to downloads list failed:", err)
		}
	}
	return artifactURLs, err
}

func (b *BintrayConnectionInfo) release(buildSpec params.ProductBuildSpec, stdout io.Writer) error {
	publishURLString := strings.Join([]string{b.URL, "content", b.Subject, b.Repository, buildSpec.ProductName, buildSpec.ProductVersion, "publish"}, "/")
	return b.runBintrayCommand(publishURLString, http.MethodPost, `{"publish_wait_for_secs":-1}`, "running Bintray publish for uploaded artifacts", stdout)
}

func (b *BintrayConnectionInfo) addToDownloadsList(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) error {
	for _, currArtifactPath := range paths.artifactPaths {
		downloadsListURLString := strings.Join([]string{b.URL, "file_metadata", b.Subject, b.Repository, paths.productPath, path.Base(currArtifactPath)}, "/")
		if err := b.runBintrayCommand(downloadsListURLString, http.MethodPut, `{"list_in_downloads":true}`, "adding artifact to Bintray downloads list for package", stdout); err != nil {
			return err
		}
	}
	return nil
}

func (b *BintrayConnectionInfo) runBintrayCommand(urlString, httpMethod, jsonContent, cmdMsg string, stdout io.Writer) (rErr error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s as URL", urlString)
	}

	capitalizedMsg := cmdMsg
	if len(cmdMsg) > 0 {
		capitalizedMsg = strings.ToUpper(string(cmdMsg[0])) + cmdMsg[1:]
	}

	fmt.Fprintf(stdout, "%s...", capitalizedMsg)
	defer func() {
		fmt.Fprintln(stdout)
	}()

	reader := strings.NewReader(jsonContent)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	req := http.Request{
		Method:        httpMethod,
		URL:           url,
		Header:        header,
		Body:          ioutil.NopCloser(reader),
		ContentLength: int64(len([]byte(jsonContent))),
	}
	req.SetBasicAuth(b.Username, b.Password)

	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return errors.Wrapf(err, "%s", cmdMsg)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close response body for URL %s", urlString)
		}
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("%s resulted in response: %s", cmdMsg, resp.Status)
	}

	fmt.Fprint(stdout, "done")

	return nil
}
