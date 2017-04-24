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
	"fmt"
	"io"
	"os"
	"path"
	"sort"

	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"github.com/termie/go-shutil"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

type osArchsBinDister params.OSArchsBinDistInfo

func (o *osArchsBinDister) NumArtifacts() int {
	return 1
}

func (o *osArchsBinDister) ArtifactPathsInOutputDir(buildSpec params.ProductBuildSpec) []string {
	var outPaths []string
	for _, osArch := range o.OSArchs {
		outPaths = append(outPaths, fmt.Sprintf("%s-%s-%s.tgz", buildSpec.ProductName, buildSpec.ProductVersion, osArch.String()))
	}
	return outPaths
}

func (o *osArchsBinDister) Dist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, spec specdir.LayoutSpec, values specdir.TemplateValues, stdout io.Writer) (Packager, error) {
	buildSpec := buildSpecWithDeps.Spec
	for _, osArch := range o.OSArchs {
		if err := verifyDistTargetSupported(osArch, buildSpecWithDeps); err != nil {
			return nil, err
		}
	}

	// each index holds all of the files required for the OS/Arch at the corresponding index
	outputPathsForOSArchs := make([][]string, len(o.OSArchs))

	for i, osArch := range o.OSArchs {
		// copy executable for current product
		dst, err := copyArtifactForOSArch(outputProductDir, buildSpec, osArch)
		if err != nil {
			return nil, err
		}
		outputPathsForOSArchs[i] = append(outputPathsForOSArchs[i], dst)
	}

	for i, osArch := range o.OSArchs {
		// copy executables for dependent products
		for _, currDepSpec := range buildSpecWithDeps.Deps {
			dst, err := copyArtifactForOSArch(outputProductDir, currDepSpec, osArch)
			if err != nil {
				return nil, err
			}
			outputPathsForOSArchs[i] = append(outputPathsForOSArchs[i], dst)
		}
	}

	outputArtifactPaths := FullArtifactsPaths(o, buildSpec, distCfg)
	artifactToInputPaths := make(map[string][]string)
	for i, currPaths := range outputPathsForOSArchs {
		artifactToInputPaths[outputArtifactPaths[i]] = currPaths
	}
	return tgzPackager(outputArtifactPaths, artifactToInputPaths), nil
}

func (o *osArchsBinDister) DistPackageType() string {
	return "tgz"
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
	dst := path.Join(outputProductDir, osArch.String(), build.ExecutableName(buildSpec.ProductName, osArch.OS))

	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return "", errors.Wrapf(err, "failed to create output directory for artifact")
	}
	if _, err := shutil.Copy(artifactPath, dst, false); err != nil {
		return "", errors.Wrapf(err, "failed to copy build artifact from %s to %s", artifactPath, dst)
	}
	return dst, nil
}
