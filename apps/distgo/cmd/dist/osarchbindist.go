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
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/termie/go-shutil"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

func osArchBinDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string) (Packager, error) {
	buildSpec := buildSpecWithDeps.Spec
	osArchBinDistInfo, ok := distCfg.Info.(*params.OSArchBinDistInfo)
	if !ok {
		osArchBinDistInfo = &params.OSArchBinDistInfo{}
		distCfg.Info = osArchBinDistInfo
	}

	osArch := osArchBinDistInfo.OSArch
	if err := verifyDistTargetSupported(osArch, buildSpecWithDeps); err != nil {
		return nil, err
	}

	var outputPaths []string
	// copy executable for current product
	dst, err := copyArtifactForOSArch(outputProductDir, buildSpec, osArch)
	if err != nil {
		return nil, err
	}
	outputPaths = append(outputPaths, dst)

	// copy executables for dependent products
	for _, currDepSpec := range buildSpecWithDeps.Deps {
		dst, err := copyArtifactForOSArch(outputProductDir, currDepSpec, osArch)
		if err != nil {
			return nil, err
		}
		outputPaths = append(outputPaths, dst)
	}

	return tgzPackager(buildSpec, distCfg, outputPaths...), nil
}

func verifyDistTargetSupported(osArch osarch.OSArch, buildSpecWithDeps params.ProductBuildSpecWithDeps) error {
	spec := buildSpecWithDeps.Spec
	if !osArchInBuildSpec(osArch, spec) {
		return errors.Errorf("The OS/Arch specified for the distribution of a product must be specified as a build target for the product, "+
			"but product %s does not specify %s as one of its build targets. Current build targets: %v", spec.ProductName, osArch, spec.Build.OSArchs)
	}

	var keys []string
	for k := range buildSpecWithDeps.Deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, currKey := range keys {
		currSpec := buildSpecWithDeps.Deps[currKey]
		if !osArchInBuildSpec(osArch, spec) {
			return errors.Errorf("The OS/Arch specified for the distribution of a product must be specified as a build target for the product, "+
				"but product %s (which is a dependent product of %s) does not specify %s as one of its build targets. Current build targets: %v", currSpec.ProductName, buildSpecWithDeps.Spec.ProductName, osArch, currSpec.Build.OSArchs)
		}
	}
	return nil
}

func osArchInBuildSpec(osArch osarch.OSArch, spec params.ProductBuildSpec) bool {
	found := false
	for _, currBuildOSArch := range spec.Build.OSArchs {
		if currBuildOSArch == osArch {
			found = true
			break
		}
	}
	return found
}

func copyArtifactForOSArch(outputProductDir string, buildSpec params.ProductBuildSpec, osArch osarch.OSArch) (string, error) {
	artifactPath, ok := build.ArtifactPaths(buildSpec)[osArch]
	if !ok {
		return "", errors.Errorf("no build artifacts exist for %s", osArch)
	}
	dst := path.Join(outputProductDir, build.ExecutableName(buildSpec.ProductName, osArch.OS))
	if _, err := shutil.Copy(artifactPath, dst, false); err != nil {
		return "", errors.Wrapf(err, "failed to copy build artifact from %s to %s", artifactPath, dst)
	}
	return dst, nil
}
