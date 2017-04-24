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
	"reflect"
	"strings"

	"github.com/google/go-github/github"
	"github.com/google/go-querystring/query"
	"github.com/pkg/errors"
	"golang.org/x/oauth2"
	"gopkg.in/cheggaaa/pb.v1"

	"github.com/palantir/godel/apps/distgo/params"
)

type GitHubConnectionInfo struct {
	APIURL    string
	UploadURL string
	User      string
	Token     string

	Owner      string
	Repository string
}

func (g *GitHubConnectionInfo) Publish(buildSpec params.ProductBuildSpec, paths ProductPaths, stdout io.Writer) ([]string, error) {
	client := github.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: g.Token},
	)))

	// set base URL (should be of the form "https://api.github.com/")
	apiURL, err := url.Parse(g.APIURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s as URL for API calls", g.APIURL)
	}
	client.BaseURL = apiURL

	if g.UploadURL == "" {
		// if upload URL is not specified, derive from API URL
		g.UploadURL = strings.Replace(g.APIURL, "api.", "uploads.", 1)
	}

	// set upload URL (should be of the form "https://uploads.github.com/")
	uploadURL, err := url.Parse(g.UploadURL)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse %s as URL for upload calls", g.APIURL)
	}
	client.UploadURL = uploadURL

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
	f, err := os.Open(filePath)
	if err != nil {
		return "", errors.Wrapf(err, "failed to open artifact %s for upload", filePath)
	}
	uploadRes, _, err := githubUploadReleaseAssetWithProgress(context.Background(), client, g.Owner, g.Repository, release.GetID(), &github.UploadOptions{
		Name: path.Base(filePath),
	}, f, stdout)
	if err != nil {
		return "", errors.Wrapf(err, "failed to upload artifact %s", filePath)
	}
	return uploadRes.GetBrowserDownloadURL(), nil
}

// the following comes from github.Repositories.UploadReleaseAsset. Implementation is copied so that support for upload
// progress can be added.
func githubUploadReleaseAssetWithProgress(ctx context.Context, client *github.Client, owner, repo string, id int, opt *github.UploadOptions, file *os.File, stdout io.Writer) (*github.ReleaseAsset, *github.Response, error) {
	u := fmt.Sprintf("repos/%s/%s/releases/%d/assets", owner, repo, id)
	u, err := githubAddOptions(u, opt)
	if err != nil {
		return nil, nil, err
	}

	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	if stat.IsDir() {
		return nil, nil, errors.New("the asset to upload can't be a directory")
	}

	// new code for progress
	fmt.Fprintf(stdout, "Uploading %s to %s\n", file.Name(), client.UploadURL.String()+u)
	bar := pb.New(int(stat.Size())).SetUnits(pb.U_BYTES)
	bar.Output = stdout
	bar.SetMaxWidth(120)
	bar.Start()
	defer bar.Finish()
	reader := bar.NewProxyReader(file)
	// done

	mediaType := mime.TypeByExtension(filepath.Ext(file.Name()))
	req, err := client.NewUploadRequest(u, reader, stat.Size(), mediaType)
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

// Taken from the github package.
// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func githubAddOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}
