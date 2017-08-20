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
	"encoding/json"
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

type ArtifactoryConnectionInfo struct {
	BasicConnectionInfo
	Repository string
}

func (a *ArtifactoryConnectionInfo) Publish(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) (rURLs []string, rErr error) {
	artifactoryURL := strings.Join([]string{a.URL, "artifactory"}, "/")
	baseURL := strings.Join([]string{artifactoryURL, a.Repository, paths.productPath}, "/")

	artifactExists := func(fi fileInfo, dstFileName, username, password string) bool {
		rawCheckArtifactURL := strings.Join([]string{artifactoryURL, "api", "storage", a.Repository, paths.productPath, dstFileName}, "/")
		checkArtifactURL, err := url.Parse(rawCheckArtifactURL)
		if err != nil {
			return false
		}

		header := http.Header{}
		req := http.Request{
			Method: http.MethodGet,
			URL:    checkArtifactURL,
			Header: header,
		}
		req.SetBasicAuth(username, password)

		if resp, err := http.DefaultClient.Do(&req); err == nil {
			defer func() {
				if err := resp.Body.Close(); err != nil && rErr == nil {
					rErr = errors.Wrapf(err, "failed to close response body for URL %s", rawCheckArtifactURL)
				}
			}()

			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				var jsonMap map[string]*json.RawMessage
				if err := json.Unmarshal(bytes, &jsonMap); err == nil {
					if checksumJSON, ok := jsonMap["checksums"]; ok && checksumJSON != nil {
						var dstChecksums checksums
						if err := json.Unmarshal(*checksumJSON, &dstChecksums); err == nil {
							return fi.checksums.match(dstChecksums)
						}
					}
				}
			}
			return false
		}
		return false
	}

	artifactURLs, err := a.uploadArtifacts(baseURL, paths, artifactExists, stdout)
	if err != nil {
		return artifactURLs, err
	}

	// compute SHA-256 checksums for artifacts
	if err := computeArtifactChecksums(artifactoryURL, a.Repository, a.Username, a.Password, paths, stdout); err != nil {
		// if triggering checksum computation fails, print message but don't throw error
		fmt.Fprintln(stdout, "Uploading artifacts succeeded, but failed to trigger computation of SHA-256 checksums:", err)
	}
	return artifactURLs, err
}

func computeArtifactChecksums(artifactoryURL, repoKey, username, password string, paths ProductPaths, stdout io.Writer) error {
	for _, currArtifactPath := range paths.artifactPaths {
		currArtifactURL := strings.Join([]string{paths.productPath, path.Base(currArtifactPath)}, "/")
		if err := artifactorySetSHA256Checksum(artifactoryURL, repoKey, currArtifactURL, username, password); err != nil {
			return errors.Wrapf(err, "")
		}
	}
	pomPath := strings.Join([]string{paths.productPath, path.Base(paths.pomFilePath)}, "/")
	if err := artifactorySetSHA256Checksum(artifactoryURL, repoKey, pomPath, username, password); err != nil {
		return errors.Wrapf(err, "")
	}
	return nil
}

func artifactorySetSHA256Checksum(baseURLString, repoKey, filePath, username, password string) (rErr error) {
	apiURLString := baseURLString + "/api/checksum/sha256"
	uploadURL, err := url.Parse(apiURLString)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s as URL", apiURLString)
	}

	jsonContent := fmt.Sprintf(`{"repoKey":"%s","path":"%s"}`, repoKey, filePath)
	reader := strings.NewReader(jsonContent)

	header := http.Header{}
	header.Set("Content-Type", "application/json")
	req := http.Request{
		Method:        http.MethodPost,
		URL:           uploadURL,
		Header:        header,
		Body:          ioutil.NopCloser(reader),
		ContentLength: int64(len([]byte(jsonContent))),
	}
	req.SetBasicAuth(username, password)

	resp, err := http.DefaultClient.Do(&req)
	if err != nil {
		return errors.Wrapf(err, "failed to trigger computation of SHA-256 checksum for %s", filePath)
	}
	defer func() {
		if err := resp.Body.Close(); err != nil && rErr == nil {
			rErr = errors.Wrapf(err, "failed to close response body for URL %s", apiURLString)
		}
	}()

	if resp.StatusCode >= http.StatusBadRequest {
		return errors.Errorf("triggering computation of SHA-256 checksum for %s resulted in response: %s", filePath, resp.Status)
	}
	return nil
}
