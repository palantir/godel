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

package bin

import (
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/mholt/archiver"
	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"
	"github.com/termie/go-shutil"

	"github.com/palantir/distgo/distgo"
)

const TypeName = "bin" // distribution that consists of the binaries in a "bin" directory

type Dister struct{}

func New() distgo.Dister {
	return &Dister{}
}

func (d *Dister) TypeName() (string, error) {
	return TypeName, nil
}

func (d *Dister) Artifacts(renderedName string) ([]string, error) {
	return []string{fmt.Sprintf("%s.tgz", renderedName)}, nil
}

func (d *Dister) RunDist(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo) ([]byte, error) {
	for _, osArch := range productTaskOutputInfo.Product.BuildOutputInfo.OSArchs {
		if err := verifyDistTargetSupported(osArch, productTaskOutputInfo); err != nil {
			return nil, err
		}
	}
	distWorkDir := productTaskOutputInfo.ProductDistWorkDirs()[distID]
	distWorkDirBinDir := path.Join(distWorkDir, "bin")
	if err := os.Mkdir(distWorkDirBinDir, 0755); err != nil {
		return nil, errors.Wrapf(err, "failed to create bin directory")
	}

	for _, osArch := range productTaskOutputInfo.Product.BuildOutputInfo.OSArchs {
		for _, currProductOutputInfo := range productTaskOutputInfo.AllProductOutputInfos() {
			// copy executable for current product
			if _, err := copyArtifactForOSArch(distWorkDirBinDir, productTaskOutputInfo.Project, currProductOutputInfo, osArch); err != nil {
				return nil, err
			}
		}
	}
	return nil, nil
}

func (d *Dister) GenerateDistArtifacts(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo, runDistResult []byte) error {
	distWorkDir := productTaskOutputInfo.ProductDistWorkDirs()[distID]
	dstPath := productTaskOutputInfo.ProductDistArtifactPaths()[distID][0]
	if err := archiver.TarGz.Make(dstPath, []string{distWorkDir}); err != nil {
		return errors.Wrapf(err, "failed to create TGZ archive")
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
