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

package build

import (
	"os"
	"path"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/pkg/imports"
)

// RequiresBuild returns a pointer to a distgo.ProductParam that contains only the OS/arch parameters for the outputs
// that require building. A product is considered to require building if its output executable does not exist or if the
// output executable's modification date is older than any of the Go files required to build the product. Returns nil if
// all of the outputs exist and are up-to-date.
func RequiresBuild(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam) (*distgo.ProductParam, error) {
	if productParam.Build == nil {
		return nil, nil
	}
	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compute output information for %s", productParam.ID)
	}

	pathsMap := productTaskOutputInfo.ProductBuildArtifactPaths()
	var requiresBuildOSArchs []osarch.OSArch
	for _, currOSArch := range productParam.Build.OSArchs {
		if fi, err := os.Stat(pathsMap[currOSArch]); err == nil {
			if goFiles, err := imports.AllFiles(path.Join(projectInfo.ProjectDir, productParam.Build.MainPkg)); err == nil {
				if newerThan, err := goFiles.NewerThan(fi); err == nil && !newerThan {
					// if the build artifact for the product already exists and none of the source files for the
					// product are newer than the build artifact, consider spec up-to-date
					continue
				}
			}
		}
		requiresBuildOSArchs = append(requiresBuildOSArchs, currOSArch)
	}

	if len(requiresBuildOSArchs) == 0 {
		return nil, nil
	}
	productParam.Build.OSArchs = requiresBuildOSArchs
	return &productParam, nil
}
