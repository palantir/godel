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

package dist

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path"
	"text/template"

	"github.com/pkg/errors"

	"github.com/ghodss/yaml"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/templating"
)

const (
	slsLabelPrefix = "com.palantir.sls."
)

func dockerDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist) (Packager, error) {
	var dockerDistInfo params.DockerDistInfo
	if info, ok := distCfg.Info.(*params.DockerDistInfo); ok {
		dockerDistInfo = *info
	} else {
		dockerDistInfo = params.DockerDistInfo{}
		distCfg.Info = &dockerDistInfo
	}
	if dockerDistInfo.Tag == "" {
		dockerDistInfo.Tag = buildSpecWithDeps.Spec.ProductVersion
	}
	if dockerDistInfo.Repository == "" {
		dockerDistInfo.Repository = buildSpecWithDeps.Spec.ProductName
	}

	completeTag := fmt.Sprintf("%s:%s", dockerDistInfo.Repository, dockerDistInfo.Tag)

	return packager(func() error {
		if err := buildWithCmd(completeTag, dockerDistInfo, distCfg, buildSpecWithDeps.Spec); err != nil {
			return errors.WithStack(err)
		}
		return nil
	}), nil
}

func buildWithCmd(tag string, distInfo params.DockerDistInfo, distCfg params.Dist, buildSpec params.ProductBuildSpec) error {
	labels := map[string]string{}
	if manifest, err := getManifestLabel(distInfo, distCfg, buildSpec); err == nil {
		labels[slsLabelPrefix+"manifest"] = manifest
	}
	if configuration, err := getConfigurationLabel(distInfo.ConfigurationFile); err == nil {
		labels[slsLabelPrefix+"configuration"] = configuration
	}
	for labelKey, labelValue := range distInfo.Labels {
		labels[labelKey] = labelValue
	}

	var args []string
	args = append(args, "build")
	args = append(args, "--tag", tag)
	for k, v := range labels {
		args = append(args, "--label", fmt.Sprintf("%s=%s", k, v))
	}
	args = append(args, distCfg.InputDir)

	dockerBuild := exec.Command("docker", args...)
	if output, err := dockerBuild.CombinedOutput(); err != nil {
		fmt.Printf("docker build failed with error:\n%s", string(output))
		return errors.Wrap(err, "failed to run")
	}
	return nil
}

func getManifestLabel(distInfo params.DockerDistInfo, distCfg params.Dist, buildSpec params.ProductBuildSpec) (string, error) {
	if distInfo.ManifestTemplateFile != "" {
		manifestTemplateFilePath := path.Join(buildSpec.ProjectDir, distInfo.ManifestTemplateFile)
		manifestBytes, err := ioutil.ReadFile(manifestTemplateFilePath)
		if err != nil {
			return "", errors.Wrap(err, "unable to read configuration file")
		}
		t := template.Must(template.New("manifest").Parse(string(manifestBytes)))
		manifestBuf := bytes.Buffer{}
		if err := t.Execute(&manifestBuf, templating.ConvertSpec(buildSpec, distCfg)); err != nil {
			return "", errors.Wrap(err, "unable to read configuration file")
		}
		manifestJSON, err := yaml.YAMLToJSON(manifestBuf.Bytes())
		if err != nil {
			return "", errors.Wrap(err, "failed to convert manifest YAML to JSON")
		}
		return string(manifestJSON), nil
	}

	manifest := slsManifest{
		ManifestVersion: "1.0",
		ProductGroup:    distCfg.Publish.GroupID,
		ProductName:     buildSpec.ProductName,
		ProductVersion:  buildSpec.ProductVersion,
		ProductType:     distInfo.ProductType,
		Extensions:      distInfo.ManifestExtensions,
	}

	manifestBytes, err := json.Marshal(&manifest)
	if err != nil {
		return "", errors.Wrap(err, "unable to marshal manifest")
	}
	return string(manifestBytes), nil
}

func getConfigurationLabel(configurationFile string) (string, error) {
	conf, err := ioutil.ReadFile(configurationFile)
	if err != nil {
		return "", errors.Wrap(err, "unable to read configuration file")
	}

	confJSONBytes, err := yaml.YAMLToJSON(conf)
	if err != nil {
		return "", errors.Wrap(err, "unable to convert configuration file to JSON")
	}

	return string(confJSONBytes), nil
}
