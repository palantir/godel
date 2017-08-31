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

package docker

import (
	"fmt"
	"io"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
)

func Publish(products []string, cfg params.Project, wd string, baseRepo string, stdout io.Writer) error {
	// find all products with docker images and tag them with correct version and push
	_, productsToPublishImage, err := productsToDistAndBuildImage(products, cfg)
	if err != nil {
		return err
	}
	buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, productsToPublishImage, wd)
	if err != nil {
		return err
	}
	for _, specWithDeps := range buildSpecsWithDeps {
		versionTag := specWithDeps.Spec.ProductVersion
		for _, image := range specWithDeps.Spec.DockerImages {
			repo := image.Repository
			if baseRepo != "" {
				repo = path.Join(baseRepo, repo)
			}
			buildTag := fmt.Sprintf("%s:%s", repo, image.Tag)
			publishTag := fmt.Sprintf("%s:%s", repo, versionTag)
			if err := tagImage(buildTag, publishTag); err != nil {
				return err
			}
			if err := pushImage(publishTag); err != nil {
				return err
			}
		}
	}
	return nil
}

func tagImage(original, new string) error {
	var args []string
	args = append(args, "tag")
	args = append(args, original)
	args = append(args, new)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker tag failed with error:\n%s\n Make sure to run docker build before docker publish.\n", string(output)))
	}
	return nil
}

func pushImage(tag string) error {
	var args []string
	args = append(args, "push")
	args = append(args, tag)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker push failed with error:\n%s\n", string(output)))
	}
	return nil
}
