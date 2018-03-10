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

package distgo

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"sort"
	"strings"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/pkg/git"
)

func ProjectVersion(projectDir, versionScriptContent string) (rVersion string, rErr error) {
	if versionScriptContent == "" {
		return git.ProjectVersion(projectDir)
	}

	tmpDir, err := ioutil.TempDir("", "distgo_script")
	if err != nil {
		return "", errors.Wrapf(err, "failed to create temporary directory")
	}
	defer func() {
		if err := os.RemoveAll(tmpDir); rErr == nil && err != nil {
			rErr = errors.Wrapf(err, "failed to remove temporary directory")
		}
	}()

	versionScript := path.Join(tmpDir, "version")
	if err := ioutil.WriteFile(versionScript, []byte(versionScriptContent), 0755); err != nil {
		return "", errors.Wrapf(err, "failed to write version script to %s", versionScript)
	}
	versionScriptCmd := exec.Command(versionScript)
	versionScriptCmd.Dir = projectDir
	outputBytes, err := versionScriptCmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		return "", errors.Wrapf(err, "command %v failed with output %s", versionScriptCmd.Args, output)
	}
	parts := strings.Split(output, "\n")
	if len(parts) == 0 {
		return git.Unspecified, nil
	}
	scriptOutput := strings.TrimSpace(parts[0])
	if scriptOutput == "" {
		return git.Unspecified, nil
	}
	return scriptOutput, nil
}

func ToProductTaskOutputInfo(projectInfo ProjectInfo, productParam ProductParam) (ProductTaskOutputInfo, error) {
	var deps map[ProductID]ProductOutputInfo
	if len(productParam.AllDependencies) > 0 {
		deps = make(map[ProductID]ProductOutputInfo)
		for k, v := range productParam.AllDependencies {
			productOutputInfo, err := v.ToProductOutputInfo(projectInfo.Version)
			if err != nil {
				return ProductTaskOutputInfo{}, err
			}
			deps[k] = productOutputInfo
		}
	}
	productOutputInfo, err := productParam.ToProductOutputInfo(projectInfo.Version)
	if err != nil {
		return ProductTaskOutputInfo{}, err
	}
	return ProductTaskOutputInfo{
		Project: projectInfo,
		Product: productOutputInfo,
		Deps:    deps,
	}, nil
}

type ProductTaskOutputInfo struct {
	Project ProjectInfo                     `json:"project"`
	Product ProductOutputInfo               `json:"product"`
	Deps    map[ProductID]ProductOutputInfo `json:"deps"`
}

func (p *ProductTaskOutputInfo) AllProductOutputInfos() []ProductOutputInfo {
	allProductOutputInfos := []ProductOutputInfo{p.Product}
	for _, buildInfo := range p.Deps {
		allProductOutputInfos = append(allProductOutputInfos, buildInfo)
	}
	return allProductOutputInfos
}

func (p *ProductTaskOutputInfo) AllProductOutputInfosMap() map[ProductID]ProductOutputInfo {
	allMap := make(map[ProductID]ProductOutputInfo)
	allMap[p.Product.ID] = p.Product
	for k, v := range p.Deps {
		allMap[k] = v
	}
	return allMap
}

func (p *ProductTaskOutputInfo) ProductBuildOutputDir() string {
	return ProductBuildOutputDir(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductBuildArtifactPaths() map[osarch.OSArch]string {
	return ProductBuildArtifactPaths(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductDistOutputDir(distID DistID) string {
	return ProductDistOutputDir(p.Project, p.Product, distID)
}

func (p *ProductTaskOutputInfo) ProductDistWorkDirs() map[DistID]string {
	return ProductDistWorkDirs(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductDistArtifactPaths() map[DistID][]string {
	return ProductDistArtifactPaths(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductDistWorkDirsAndArtifactPaths() map[DistID][]string {
	return ProductDistWorkDirsAndArtifactPaths(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductDockerBuildArtifactPaths() map[DockerID]map[ProductID]map[osarch.OSArch]string {
	return ProductDockerBuildArtifactPaths(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) ProductDockerDistArtifactPaths() map[DockerID]map[ProductID]map[DistID][]string {
	return ProductDockerDistArtifactPaths(p.Project, p.Product)
}

func (p *ProductTaskOutputInfo) POM(groupID string) (string, string, error) {
	pomName := fmt.Sprintf("%s-%s.pom", p.Product.ID, p.Project.Version)
	pomContent, err := generatePOMContent(p.Product.ID, p.Project.Version, groupID)
	if err != nil {
		return "", "", err
	}
	return pomName, pomContent, nil
}

func ExecutableName(productName, goos string) string {
	executableName := productName
	if goos == "windows" {
		executableName += ".exe"
	}
	return executableName
}

// ProductBuildOutputDir returns the output directory for the build outputs, which is
// "{{ProjectDir}}/{{OutputDir}}/{{ProductID}}/{{Version}}".
func ProductBuildOutputDir(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) string {
	if productOutputInfo.BuildOutputInfo == nil {
		return ""
	}
	return path.Join(projectInfo.ProjectDir, productOutputInfo.BuildOutputInfo.BuildOutputDir, string(productOutputInfo.ID), projectInfo.Version)
}

// ProductBuildArtifactPaths returns a map that contains the paths to the executables created by the provided product
// for the provided project. The keys in the map are the OS/architecture of the executable and the values are the
// executable output paths for that OS/architecture. The output paths are of the form
// "{{ProjectDir}}/{{OutputDir}}/{{ProductID}}/{{Version}}/{{OSArch}}/{{NameTemplateRendered}}" (and if the OS is
// Windows, the ".exe" extension is appended).
func ProductBuildArtifactPaths(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[osarch.OSArch]string {
	if productOutputInfo.BuildOutputInfo == nil {
		return nil
	}
	paths := make(map[osarch.OSArch]string)
	for _, osArch := range productOutputInfo.BuildOutputInfo.OSArchs {
		executableName := ExecutableName(productOutputInfo.BuildOutputInfo.BuildNameTemplateRendered, osArch.OS)
		paths[osArch] = path.Join(ProductBuildOutputDir(projectInfo, productOutputInfo), osArch.String(), executableName)
	}
	return paths
}

// ProductDistOutputDir returns the output directory for the dist outputs for the dist with the given DistID, which is
// "{{ProjectDir}}/{{OutputDir}}/{{ProductID}}/{{Version}}/{{DistID}}".
func ProductDistOutputDir(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo, distID DistID) string {
	if productOutputInfo.DistOutputInfos == nil {
		return ""
	}
	return path.Join(projectInfo.ProjectDir, productOutputInfo.DistOutputInfos.DistOutputDir, string(productOutputInfo.ID), projectInfo.Version, string(distID))
}

// ProductDistWorkDirs returns a map from DistID to the directory used to prepare the distribution for that DistID,
// which is "{{ProjectDir}}/{{OutputDir}}/{{ProductID}}/{{Version}}/{{DistID}}/{{NameTemplateRendered}}".
func ProductDistWorkDirs(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[DistID]string {
	if productOutputInfo.DistOutputInfos == nil {
		return nil
	}
	workDirs := make(map[DistID]string)
	for distID, distOutputInfo := range productOutputInfo.DistOutputInfos.DistInfos {
		workDirs[distID] = path.Join(ProductDistOutputDir(projectInfo, productOutputInfo, distID), distOutputInfo.DistNameTemplateRendered)
	}
	return workDirs
}

// ProductDistArtifactPaths returns a map from DistID to the output paths for the dist, which is
// "{{ProjectDir}}/{{OutputDir}}/{{ProductID}}/{{Version}}/{{DistID}}/{{Artifacts}}".
func ProductDistArtifactPaths(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[DistID][]string {
	if productOutputInfo.DistOutputInfos == nil {
		return nil
	}
	paths := make(map[DistID][]string)
	for distID, distOutputInfo := range productOutputInfo.DistOutputInfos.DistInfos {
		for _, currArtifactPath := range distOutputInfo.DistArtifactNames {
			paths[distID] = append(paths[distID], path.Join(ProductDistOutputDir(projectInfo, productOutputInfo, distID), currArtifactPath))
		}
		sort.Strings(paths[distID])
	}
	return paths
}

// ProductDistWorkDirsAndArtifactPaths returns a map that is the result of joining the values of the outputs of
// ProductDistWorkDirs and ProductDistArtifactPaths.
func ProductDistWorkDirsAndArtifactPaths(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[DistID][]string {
	paths := ProductDistArtifactPaths(projectInfo, productOutputInfo)
	if paths == nil {
		return nil
	}
	for k, v := range ProductDistWorkDirs(projectInfo, productOutputInfo) {
		paths[k] = append(paths[k], v)
		sort.Strings(paths[k])
	}
	return paths
}

// ProductDockerBuildArtifactPaths returns a map that contains the paths to the locations where the input build
// artifacts should be placed in the Docker context directory. The DockerID key identifies the DockerBuilder, the
// ProductID represents the input product for that DockerBuilder, and the osarch.OSArch represents the OS/Arch for the
// build. Paths are of the form "{{ProjectDir}}/{{DockerID.ContextDir}}/{{DockerID.InputProductsDir}}/{{ProductID}}/build/{{OSArch}}/{{ExecutableName}}".
func ProductDockerBuildArtifactPaths(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[DockerID]map[ProductID]map[osarch.OSArch]string {
	if productOutputInfo.DockerOutputInfos == nil {
		return nil
	}
	out := make(map[DockerID]map[ProductID]map[osarch.OSArch]string)
	for _, dockerID := range productOutputInfo.DockerOutputInfos.DockerIDs {
		out[dockerID] = make(map[ProductID]map[osarch.OSArch]string)

		dockerOutputInfo := productOutputInfo.DockerOutputInfos.DockerBuilderOutputInfos[dockerID]
		pathToInputProductsDir := path.Join(projectInfo.ProjectDir, dockerOutputInfo.ContextDir, dockerOutputInfo.InputProductsDir)
		for productID, valMap := range dockerOutputInfo.InputBuilds {
			if _, ok := out[dockerID][productID]; !ok {
				out[dockerID][productID] = make(map[osarch.OSArch]string)
			}
			for osArchID := range valMap {
				osArch, err := osarch.New(string(osArchID))
				if err != nil {
					panic(errors.Wrapf(err, "OSArchID was not in a valid state"))
				}
				artifactPath := path.Join(pathToInputProductsDir, string(productID), "build", string(osArchID), ExecutableName(productOutputInfo.BuildOutputInfo.BuildNameTemplateRendered, osArch.OS))
				out[dockerID][productID][osArch] = artifactPath
			}
		}
	}
	return out
}

// ProductDockerDistArtifactPaths returns a map that contains the paths to the locations where the input dist artifacts
// should be placed in the Docker context directory. The DockerID key identifies the DockerBuilder, the ProductID
// represents the input product for that DockerBuilder, and the DistID represents the Dister for the product. Paths are
// of the form "{{ProjectDir}}/{{DockerID.ContextDir}}/{{DockerID.InputProductsDir}}/{{ProductID}}/dist/{{DistID}}/{{Artifacts}}".
func ProductDockerDistArtifactPaths(projectInfo ProjectInfo, productOutputInfo ProductOutputInfo) map[DockerID]map[ProductID]map[DistID][]string {
	if productOutputInfo.DockerOutputInfos == nil {
		return nil
	}
	out := make(map[DockerID]map[ProductID]map[DistID][]string)
	for _, dockerID := range productOutputInfo.DockerOutputInfos.DockerIDs {
		out[dockerID] = make(map[ProductID]map[DistID][]string)

		dockerOutputInfo := productOutputInfo.DockerOutputInfos.DockerBuilderOutputInfos[dockerID]
		pathToInputProductsDir := path.Join(projectInfo.ProjectDir, dockerOutputInfo.ContextDir, dockerOutputInfo.InputProductsDir)
		for productID, valMap := range dockerOutputInfo.InputDists {
			if _, ok := out[dockerID][productID]; !ok {
				out[dockerID][productID] = make(map[DistID][]string)
			}
			productDistArtifacts := ProductDistArtifactPaths(projectInfo, productOutputInfo)
			for distID := range valMap {
				for _, origArtifactPath := range productDistArtifacts[distID] {
					artifactPath := path.Join(pathToInputProductsDir, string(productID), "dist", string(distID), path.Base(origArtifactPath))
					out[dockerID][productID][distID] = append(out[dockerID][productID][distID], artifactPath)
				}
			}
		}
	}
	return out
}

type ProjectInfo struct {
	ProjectDir string `json:"projectDir"`
	Version    string `json:"version"`
}
