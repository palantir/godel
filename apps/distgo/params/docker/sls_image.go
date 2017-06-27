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
	"encoding/base64"
	"fmt"
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

type SLSDockerImage struct {
	DefaultDockerImage
	GroupID      string
	ProuductType string
	Extensions   map[string]interface{}
}

func (sdi *SLSDockerImage) Build(buildSpec params.ProductBuildSpecWithDeps) error {
	contextDir := path.Join(buildSpec.Spec.ProjectDir, sdi.ContextDir)
	configFile := path.Join(contextDir, ConfigurationFileName)
	var args []string
	args = append(args, "build")
	args = append(args, "--tag", fmt.Sprintf("%s:%s", sdi.Repository, sdi.Tag))
	manifest, err := params.GetManifest(sdi.GroupID, buildSpec.Spec.ProductName, buildSpec.Spec.ProductVersion, sdi.ProuductType, sdi.Extensions)
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
	buildArgs, err := script.GetBuildArgs(buildSpec.Spec, sdi.BuildArgsScript)
	if err != nil {
		return errors.Wrap(err, "")
	}
	args = append(args, buildArgs...)
	args = append(args, contextDir)

	buildCmd := exec.Command("docker", args...)
	if output, err := buildCmd.CombinedOutput(); err != nil {
		return errors.Wrap(err, fmt.Sprintf("docker build failed with error:\n%s\n", string(output)))
	}
	return nil
}
