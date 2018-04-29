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

package mavenlocal

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/publisher"
	"github.com/palantir/distgo/publisher/mavenlocal/config"
)

const TypeName = "maven-local" // publishes output artifacts to a location in the local filesystem

func PublisherCreator() publisher.Creator {
	return publisher.NewCreator(TypeName, func() distgo.Publisher {
		return &mavenLocalPublisher{}
	})
}

type mavenLocalPublisher struct{}

func (p *mavenLocalPublisher) TypeName() (string, error) {
	return TypeName, nil
}

var (
	mavenLocalPublisherBaseDirFlag = distgo.PublisherFlag{
		Name:        "base-dir",
		Description: "base output directory for the local publish (if blank, defaults to ${HOME}/.m2/repository)",
		Type:        distgo.StringFlag,
	}
)

func (p *mavenLocalPublisher) Flags() ([]distgo.PublisherFlag, error) {
	return []distgo.PublisherFlag{
		mavenLocalPublisherBaseDirFlag,
		publisher.GroupIDFlag,
	}, nil
}

func (p *mavenLocalPublisher) RunPublish(productTaskOutputInfo distgo.ProductTaskOutputInfo, cfgYML []byte, flagVals map[distgo.PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error {
	var cfg config.MavenLocal
	if err := yaml.Unmarshal(cfgYML, &cfg); err != nil {
		return errors.Wrapf(err, "failed to unmarshal configuration")
	}
	groupID, err := publisher.GetRequiredGroupID(flagVals, productTaskOutputInfo)
	if err != nil {
		return err
	}
	if err := publisher.SetConfigValue(flagVals, mavenLocalPublisherBaseDirFlag, &cfg.BaseDir); err != nil {
		return err
	}

	baseDir := cfg.BaseDir
	if baseDir == "" {
		baseDir = path.Join(os.Getenv("HOME"), ".m2", "repository")
	}

	groupPath := strings.Replace(groupID, ".", "/", -1)
	productPath := path.Join(baseDir, groupPath, string(productTaskOutputInfo.Product.ID), productTaskOutputInfo.Project.Version)
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

	// if error is non-nil, wd will be empty
	wd, _ := os.Getwd()
	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		for _, currArtifactPath := range productTaskOutputInfo.ProductDistArtifactPaths()[currDistID] {
			if _, err := copyArtifact(currArtifactPath, productPath, wd, dryRun, stdout); err != nil {
				return errors.Wrapf(err, "failed to copy artifact")
			}
		}
	}
	return nil
}

func copyArtifact(src, dstDir, wd string, dryRun bool, stdout io.Writer) (string, error) {
	dst := path.Join(dstDir, path.Base(src))
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Copying artifact from %s to %s", toRelPath(src, wd), dst), dryRun)
	if !dryRun {
		if err := shutil.CopyFile(src, dst, false); err != nil {
			return "", errors.Wrapf(err, "failed to copy %s to %s", src, dst)
		}
	}
	return dst, nil
}

func toRelPath(path, wd string) string {
	if !filepath.IsAbs(path) || wd == "" {
		return path
	}
	relPath, err := filepath.Rel(wd, path)
	if err != nil {
		return path
	}
	return relPath
}
