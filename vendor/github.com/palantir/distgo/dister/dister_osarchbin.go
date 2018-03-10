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

package dister

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/mholt/archiver"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"
	"github.com/termie/go-shutil"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/distgo"
)

const OSArchBinDistTypeName = "os-arch-bin" // distribution that consists of the binaries for a specific OS/Architecture

type OSArchBinDistConfig struct {
	// OSArchs specifies the GOOS and GOARCH pairs for which TGZ distributions are created. If blank, defaults to
	// the GOOS and GOARCH of the host system at runtime.
	OSArchs []osarch.OSArch `yaml:"os-archs"`
}

func (cfg *OSArchBinDistConfig) ToDister() distgo.Dister {
	osArchs := cfg.OSArchs
	if len(osArchs) == 0 {
		osArchs = []osarch.OSArch{osarch.Current()}
	}
	return &osArchsBinDister{
		OSArchs: osArchs,
	}
}

type osArchsBinDister struct {
	OSArchs []osarch.OSArch
}

func NewOSArchBinDister(osArchs ...osarch.OSArch) distgo.Dister {
	return &osArchsBinDister{
		OSArchs: osArchs,
	}
}

func NewOSArchBinDisterFromConfig(cfgYML []byte) (distgo.Dister, error) {
	var disterCfg OSArchBinDistConfig
	if err := yaml.Unmarshal(cfgYML, &disterCfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal YAML")
	}
	return disterCfg.ToDister(), nil
}

func (d *osArchsBinDister) TypeName() (string, error) {
	return OSArchBinDistTypeName, nil
}

func (d *osArchsBinDister) Artifacts(renderedName string) ([]string, error) {
	var outPaths []string
	for _, osArch := range d.OSArchs {
		outPaths = append(outPaths, fmt.Sprintf("%s-%s.tgz", renderedName, osArch.String()))
	}
	return outPaths, nil
}

func (d *osArchsBinDister) osArchFromArtifactPath(distID distgo.DistID, artifactPath string, productTaskOutputInfo distgo.ProductTaskOutputInfo) (osarch.OSArch, error) {
	for _, osArch := range d.OSArchs {
		if strings.HasSuffix(artifactPath, fmt.Sprintf("%s-%s.tgz", productTaskOutputInfo.Product.DistOutputInfos.DistInfos[distID].DistNameTemplateRendered, osArch.String())) {
			return osArch, nil
		}
	}
	return osarch.OSArch{}, errors.Errorf("failed to determine OS/Arch for artifact with Path %s", artifactPath)
}

func (d *osArchsBinDister) RunDist(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo) ([]byte, error) {
	for _, osArch := range d.OSArchs {
		if err := verifyDistTargetSupported(osArch, productTaskOutputInfo); err != nil {
			return nil, err
		}
	}
	distWorkDir := productTaskOutputInfo.ProductDistWorkDirs()[distID]
	outputPathsForOSArchs := make(map[string][]string)
	for _, osArch := range d.OSArchs {
		for _, currProductOutputInfo := range productTaskOutputInfo.AllProductOutputInfos() {
			// copy executable for current product
			dst, err := copyArtifactForOSArch(distWorkDir, productTaskOutputInfo.Project, currProductOutputInfo, osArch)
			if err != nil {
				return nil, err
			}
			outputPathsForOSArchs[osArch.String()] = append(outputPathsForOSArchs[osArch.String()], dst)
		}
	}
	jsonBytes, err := json.Marshal(outputPathsForOSArchs)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to marshal outputPathsForOSArchs as JSON")
	}
	return jsonBytes, nil
}

func (d *osArchsBinDister) GenerateDistArtifacts(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo, runDistResult []byte) error {
	var outputPathsForOSArchs map[string][]string
	if err := json.Unmarshal(runDistResult, &outputPathsForOSArchs); err != nil {
		return errors.Wrapf(err, "failed to unmarshal runDistResult JSON %s", string(runDistResult))
	}
	artifactToInputPaths := make(map[string][]string)
	outputArtifactPaths := productTaskOutputInfo.ProductDistArtifactPaths()[distID]
	for _, artifactPath := range outputArtifactPaths {
		currOSArch, err := d.osArchFromArtifactPath(distID, artifactPath, productTaskOutputInfo)
		if err != nil {
			return err
		}
		artifactToInputPaths[artifactPath] = outputPathsForOSArchs[currOSArch.String()]
	}
	if err := createTgz(outputArtifactPaths, artifactToInputPaths); err != nil {
		return err
	}
	return nil
}

func createTgz(dstArtifactPaths []string, dstToContentPaths map[string][]string) error {
	for _, currDstPath := range dstArtifactPaths {
		if err := archiver.TarGz.Make(currDstPath, dstToContentPaths[currDstPath]); err != nil {
			return errors.Wrapf(err, "failed to create TGZ archive")
		}
	}
	return nil
}

func verifyDistTargetSupported(osArch osarch.OSArch, productTaskOutputInfo distgo.ProductTaskOutputInfo) error {
	if err := verifySingleProduct(osArch, productTaskOutputInfo.Product); err != nil {
		return err
	}
	var keys []distgo.ProductID
	for k := range productTaskOutputInfo.Deps {
		keys = append(keys, k)
	}
	sort.Sort(distgo.ByProductID(keys))
	for _, currKey := range keys {
		currSpec := productTaskOutputInfo.Deps[currKey]
		if err := verifySingleProduct(osArch, currSpec); err != nil {
			return err
		}
	}
	return nil
}

func verifySingleProduct(osArch osarch.OSArch, productOutputInfo distgo.ProductOutputInfo) error {
	if !osArchInBuildSpec(osArch, productOutputInfo) {
		buildOSArchs := "[none]"
		if productOutputInfo.BuildOutputInfo != nil {
			buildOSArchs = fmt.Sprint(productOutputInfo.BuildOutputInfo.OSArchs)
		}
		return errors.Errorf("the OS/Arch specified for the distribution of a product must be specified as a build target for the product, "+
			"but product %s does not specify %s as one of its build targets (current build targets: %s)", productOutputInfo.ID, osArch, buildOSArchs)
	}
	return nil
}

func osArchInBuildSpec(osArch osarch.OSArch, productOutputInfo distgo.ProductOutputInfo) bool {
	if productOutputInfo.BuildOutputInfo == nil {
		return false
	}
	found := false
	for _, currBuildOSArch := range productOutputInfo.BuildOutputInfo.OSArchs {
		if currBuildOSArch == osArch {
			found = true
			break
		}
	}
	return found
}

func copyArtifactForOSArch(outputDir string, projectInfo distgo.ProjectInfo, productInfo distgo.ProductOutputInfo, osArch osarch.OSArch) (string, error) {
	artifactPath, ok := distgo.ProductBuildArtifactPaths(projectInfo, productInfo)[osArch]
	if !ok {
		return "", errors.Errorf("no build artifacts exist for %s", osArch)
	}

	dst := path.Join(outputDir, osArch.String(), distgo.ExecutableName(productInfo.BuildOutputInfo.BuildNameTemplateRendered, osArch.OS))
	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return "", errors.Wrapf(err, "failed to create output directory for artifact")
	}
	if _, err := shutil.Copy(artifactPath, dst, false); err != nil {
		return "", errors.Wrapf(err, "failed to copy build artifact from %s to %s", artifactPath, dst)
	}
	return dst, nil
}
