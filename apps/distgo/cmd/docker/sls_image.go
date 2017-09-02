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
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/script"
)

const (
	ManifestLabel         = "com.palantir.sls.manifest"
	ConfigurationLabel    = "com.palantir.sls.configuration"
	ConfigurationFileName = "configuration.yml"
)

type slsImageBuilder struct {
	image *params.DockerImage
	info  *params.SLSDockerImageInfo
}

func (sib *slsImageBuilder) build(buildSpec params.ProductBuildSpecWithDeps, stdout io.Writer) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, sib.image.ContextDir)
	configFile := path.Join(contextDir, ConfigurationFileName)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", sib.image.Repository, sib.image.Tag))
	manifest, err := params.GetManifest(sib.info.GroupID, buildSpec.Spec.ProductName, buildSpec.Spec.ProductVersion, sib.info.ProuductType, sib.info.Extensions)
	if err != nil {
		return errors.Wrap(err, "Failed to get manifest for the image")
	}
	args = append(args, "--label", fmt.Sprintf("%s=%s", ManifestLabel, base64.StdEncoding.EncodeToString([]byte(manifest))))
	if _, err := os.Stat(configFile); err == nil {
		content, err := ioutil.ReadFile(configFile)
		if err != nil {
			return errors.Wrapf(err, "Failed to read the file %s", configFile)
		}
		args = append(args, "--label", fmt.Sprintf("%s=%s", ConfigurationLabel, base64.StdEncoding.EncodeToString(content)))
	}
	buildArgs, err := script.GetBuildArgs(buildSpec.Spec, sib.image.BuildArgsScript)
	if err != nil {
		return err
	}
	args = append(args, buildArgs...)
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	bufWriter := &bytes.Buffer{}
	buildCmd.Stdout = io.MultiWriter(stdout, bufWriter)
	buildCmd.Stderr = bufWriter
	if err := buildCmd.Run(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", bufWriter.String()))
	}
	return nil
}
