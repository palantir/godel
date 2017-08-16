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
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/google/go-github/github"
	"github.com/jtacoma/uritemplates"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/git"
)

type GitHubConnectionInfo struct {
	APIURL string
	User   string
	Token  string

	Owner      string
	Repository string
}

func (g *GitHubConnectionInfo) Publish(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) ([]string, error) {
	if version := buildSpec.VersionInfo.Version; version == git.Unspecified || git.IsSnapshotVersion(version) || strings.HasSuffix(version, ".dirty") {
		return nil, errors.Errorf("cannot perform publish on repository with version %s: GitHub publish task requires repository to be on a clean tagged commit", version)
	}

	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)))

	// if base URL does not end in "/", append it (trailing slash is required)
	if !strings.HasSuffix(g.APIURL, "/") {
		g.APIURL += "/"
	}
	// set base URL (should be of the form "https://api.github.com/")
	apiURL, err := url.Parse(g.APIURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s as URL for API calls", g.APIURL)
	}
	client.BaseURL = apiURL

	if g.Owner == "" {
		g.Owner = g.User
	}

	fmt.Fprintf(stdout, "Creating GitHub release %s for %s/%s...", buildSpec.ProductVersion, g.Owner, g.Repository)
	releaseRes, _, err := client.Repositories.CreateRelease(context.Background(), g.Owner, g.Repository, &github.RepositoryRelease{
		TagName: github.String(buildSpec.ProductVersion),
	})
	if err != nil {
		// newline to complement "..." output
		fmt.Fprintln(stdout)

		if ghErr, ok := err.(*github.ErrorResponse); ok && len(ghErr.Errors) > 0 {
			if ghErr.Errors[0].Code == "already_exists" {
				return nil, errors.Errorf("GitHub release %s already exists for %s/%s", buildSpec.ProductVersion, g.Owner, g.Repository)
			}
		}
		return nil, errors.Wrapf(err, "failed to create GitHub release %s for %s/%s...", buildSpec.ProductVersion, g.Owner, g.Repository)
	}
	fmt.Fprintln(stdout, "done")

	var uploadURLs []string
	for _, currPath := range paths.artifactPaths {
		uploadURL, err := g.uploadFileAtPath(client, releaseRes, currPath, stdout)
		if err != nil {
			return nil, err
		}
		uploadURLs = append(uploadURLs, uploadURL)
	}
	return uploadURLs, nil
}

func (g *GitHubConnectionInfo) uploadFileAtPath(client *github.Client, release *github.RepositoryRelease, filePath string, stdout io.Writer) (string, error) {
	uploadURI, err := uploadURIForProduct(release.GetUploadURL(), path.Base(filePath))
	if err != nil {
		return "", err
	}

	f, err := os.Open(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open artifact %s for upload", filePath)
	}
	uploadRes, _, err := githubUploadReleaseAssetWithProgress(context.Background(), client, uploadURI, f, stdout)
	if err != nil {
		return "", errors.Wrapf(err, "failed to upload artifact %s", filePath)
	}
	return uploadRes.GetBrowserDownloadURL(), nil
}

// uploadURIForProduct returns an asset upload URI using the provided upload template from the release creation
// response. See https://developer.github.com/v3/repos/releases/#response for the specifics of the API.
func uploadURIForProduct(githubUploadURLTemplate, name string) (string, error) {
	const nameTemplate = "name"

	t, err := uritemplates.Parse(githubUploadURLTemplate)
	if err != nil {
		return "", errors.Wrapf(err, "failed to parse upload URI template %q", githubUploadURLTemplate)
	}
	uploadURI, err := t.Expand(map[string]interface{}{
		nameTemplate: name,
	})
	if err != nil {
		return "", errors.Wrapf(err, "failed to expand URI template %q with %q = %q", githubUploadURLTemplate, nameTemplate, name)
	}
	return uploadURI, nil
}

// Based on github.Repositories.UploadReleaseAsset. Adds support for progress reporting.
func githubUploadReleaseAssetWithProgress(ctx context.Context, client *github.Client, uploadURI string, file *os.File, stdout io.Writer) (*github.ReleaseAsset, *github.Response, error) {
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	if stat.IsDir() {
		return nil, nil, errors.New("the asset to upload can't be a directory")
	}

	fmt.Fprintf(stdout, "Uploading %s to %s\n", file.Name(), uploadURI)
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.Output = stdout
	bar.SetMaxWidth(120)
	bar.Start()
	defer bar.Finish()
	reader := bar.NewProxyReader(file)

	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	req, err := client.NewUploadRequest(uploadURI, reader, stat.Size(), mediaType)
	if err != nil {
		return nil, nil, err
	}

	asset := new(github.ReleaseAsset)
	resp, err := client.Do(ctx, req, asset)
	if err != nil {
		return nil, resp, err
	}
	return asset, resp, nil
}
