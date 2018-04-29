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

package bintray

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/publisher"
	"github.com/palantir/distgo/publisher/bintray/config"
)

const TypeName = "bintray"

type bintrayPublisher struct{}

func PublisherCreator() publisher.Creator {
	return publisher.NewCreator(TypeName, func() distgo.Publisher {
		return &bintrayPublisher{}
	})
}

func (p *bintrayPublisher) TypeName() (string, error) {
	return TypeName, nil
}

var (
	bintrayPublisherSubjectFlag = distgo.PublisherFlag{
		Name:        "subject",
		Description: "subject that is the destination for the publish",
		Type:        distgo.StringFlag,
	}
	bintrayPublisherRepositoryFlag = distgo.PublisherFlag{
		Name:        "repository",
		Description: "repository that is the destination for the publish",
		Type:        distgo.StringFlag,
	}
	bintrayPublisherProductFlag = distgo.PublisherFlag{
		Name:        "product",
		Description: "Bintray product for which publish should occur (if blank, ProductID is used)",
		Type:        distgo.StringFlag,
	}
	bintrayPublisherPublishFlag = distgo.PublisherFlag{
		Name:        "publish",
		Description: "perform a Bintray publish for the uploaded content",
		Type:        distgo.BoolFlag,
	}
	bintrayPublisherDownloadsListFlag = distgo.PublisherFlag{
		Name:        "downloads-list",
		Description: "add uploaded artifact to downloads list for package",
		Type:        distgo.BoolFlag,
	}
)

func (p *bintrayPublisher) Flags() ([]distgo.PublisherFlag, error) {
	return append(
		publisher.BasicConnectionInfoFlags(),
		bintrayPublisherSubjectFlag,
		bintrayPublisherRepositoryFlag,
		bintrayPublisherProductFlag,
		bintrayPublisherPublishFlag,
		bintrayPublisherDownloadsListFlag,
		publisher.GroupIDFlag,
	), nil
}

func (p *bintrayPublisher) RunPublish(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfgYML []byte, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	var cfg config.Bintray
	if err := yaml.Unmarshal(cfgYML, &cfg); err != nil {
		return errors.Wrapf(err, "failed to unmarshal configuration")
	}
	groupID, err := publisher.GetRequiredGroupID(flagVals, productTaskOutputInfo)
	if err != nil {
		return err
	}
	if err := cfg.BasicConnectionInfo.SetValuesFromFlags(flagVals); err != nil {
		return err
	}
	if err := publisher.SetRequiredStringConfigValues(flagVals,
		bintrayPublisherSubjectFlag, &cfg.Subject,
		bintrayPublisherRepositoryFlag, &cfg.Repository,
	); err != nil {
		return err
	}

	if err := publisher.SetConfigValue(flagVals, bintrayPublisherProductFlag, &cfg.Product); err != nil {
		return err
	}
	if cfg.Product == "" {
		cfg.Product = string(productTaskOutputInfo.Product.ID)
	}

	if err := publisher.SetConfigValues(flagVals,
		bintrayPublisherPublishFlag, &cfg.Publish,
		bintrayPublisherDownloadsListFlag, &cfg.DownloadsList,
	); err != nil {
		return err
	}

	mavenProductPath := publisher.MavenProductPath(productTaskOutputInfo, groupID)
	baseURL := strings.Join([]string{cfg.URL, "content", cfg.Subject, cfg.Repository, cfg.Product, productTaskOutputInfo.Project.Version, mavenProductPath}, "/")
	if _, _, err := cfg.BasicConnectionInfo.UploadDistArtifacts(productTaskOutputInfo, baseURL, nil, dryRun, stdout); err != nil {
		return err
	}

	if cfg.Publish {
		if err := p.publish(productTaskOutputInfo, cfg, dryRun, stdout); err != nil {
			fmt.Fprintln(stdout, "Uploading artifacts succeeded, but publish of uploaded artifacts failed:", err)
		}
	}
	if cfg.DownloadsList {
		if err := p.addToDownloadsList(productTaskOutputInfo, cfg, mavenProductPath, dryRun, stdout); err != nil {
			fmt.Fprintln(stdout, "Uploading artifacts succeeded, but addings artifact to downloads list failed:", err)
		}
	}
	return nil
}

func (p *bintrayPublisher) publish(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfg config.Bintray, dryRun bool, stdout io.Writer) error {
	publishURLString := strings.Join([]string{cfg.URL, "content", cfg.Subject, cfg.Repository, cfg.Product, productTaskOutputInfo.Project.Version, "publish"}, "/")
	return p.runBintrayCommand(publishURLString, http.MethodPost, cfg.Username, cfg.Password, `{"publish_wait_for_secs":-1}`, "running Bintray publish for uploaded artifacts", dryRun, stdout)
}

func (p *bintrayPublisher) addToDownloadsList(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfg config.Bintray, mavenProductPath string, dryRun bool, stdout io.Writer) error {
	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		for _, currArtifactPath := range productTaskOutputInfo.ProductDistArtifactPaths()[currDistID] {
			downloadsListURLString := strings.Join([]string{cfg.URL, "file_metadata", cfg.Subject, cfg.Repository, mavenProductPath, path.Base(currArtifactPath)}, "/")
			if err := p.runBintrayCommand(downloadsListURLString, http.MethodPut, cfg.Username, cfg.Password, `{"list_in_downloads":true}`, "adding artifact to Bintray downloads list for package", dryRun, stdout); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *bintrayPublisher) runBintrayCommand(urlString, httpMethod, username, password, jsonContent, cmdMsg string, dryRun bool, stdout io.Writer) (rErr error) {
	url, err := url.Parse(urlString)
	if err != nil {
		return errors.Wrapf(err, "failed to parse %s as URL", urlString)
	}

	capitalizedMsg := cmdMsg
	if len(cmdMsg) > 0 {
		capitalizedMsg = strings.ToUpper(string(cmdMsg[0])) + cmdMsg[1:]
	}

	distgo.PrintOrDryRunPrint(stdout, fmt.Sprintf("%s...", capitalizedMsg), dryRun)
	defer func() {
		// not wrapped in dry run because that has already been handled at the beginning of the line
		fmt.Fprintln(stdout)
	}()

	if !dryRun {
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
		req.SetBasicAuth(username, password)

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
	}

	// not wrapped in dry run because that has already been handled at the beginning of the line
	fmt.Fprint(stdout, "done")
	return nil
}
