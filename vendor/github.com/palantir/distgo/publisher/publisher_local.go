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

package publisher

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

const LocalPublishTypeName = "local" // publishes output artifacts to a location in the local filesystem

type LocalPublishConfig struct {
	// BaseDir is the base directory to which the artifacts are published.
	BaseDir string `yaml:"base-dir"`
}

func NewLocalPublisherCreator() Creator {
	return NewCreator(LocalPublishTypeName, func() distgo.Publisher {
		return &localPublisher{}
	})
}

type localPublisher struct{}

func (p *localPublisher) TypeName() (string, error) {
	return LocalPublishTypeName, nil
}

var (
	localPublisherBaseDirFlag = distgo.PublisherFlag{
		Name:        "base-dir",
		Description: "base output directory for the local publish",
		Type:        distgo.StringFlag,
	}
)

func (p *localPublisher) Flags() ([]distgo.PublisherFlag, error) {
	return []distgo.PublisherFlag{
		localPublisherBaseDirFlag,
		GroupIDFlag,
	}, nil
}

func (p *localPublisher) RunPublish(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfgYML []byte, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	var cfg LocalPublishConfig
	if err := yaml.Unmarshal(cfgYML, &cfg); err != nil {
		return errors.Wrapf(err, "failed to unmarshal configuration")
	}
	groupID, err := GetRequiredGroupID(flagVals, productTaskOutputInfo)
	if err != nil {
		return err
	}
	if err := SetConfigValue(flagVals, localPublisherBaseDirFlag, &cfg.BaseDir); err != nil {
		return err
	}

	groupPath := strings.Replace(groupID, ".", "/", -1)
	productPath := path.Join(cfg.BaseDir, groupPath, string(productTaskOutputInfo.Product.ID), productTaskOutputInfo.Project.Version)
	if !dryRun {
		if err := os.MkdirAll(productPath, 0755); err != nil {
			return errors.Wrapf(err, "failed to create %s", productPath)
		}
	}

	pomName, pomContent, err := productTaskOutputInfo.POM(groupID)
	if err != nil {
		return err
	}

	pomPath := path.Join(productPath, pomName)
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Writing POM to %s", pomPath), dryRun)
	if !dryRun {
		if err := ioutil.WriteFile(pomPath, []byte(pomContent), 0644); err != nil {
			return errors.Wrapf(err, "failed to write POM")
		}
	}

	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		for _, currArtifactPath := range productTaskOutputInfo.ProductDistArtifactPaths()[currDistID] {
			if _, err := copyArtifact(currArtifactPath, productPath, dryRun, stdout); err != nil {
				return errors.Wrapf(err, "failed to copy artifact")
			}
		}
	}
	return nil
}

func copyArtifact(src, dstDir string, dryRun bool, stdout io.Writer) (string, error) {
	dst := path.Join(dstDir, path.Base(src))
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Copying artifact from %s to %s", src, dst), dryRun)
	if !dryRun {
		if err := shutil.CopyFile(src, dst, false); err != nil {
			return "", errors.Wrapf(err, "failed to copy %s to %s", src, dst)
		}
	}
	return dst, nil
}
