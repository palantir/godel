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
	"path"

	"github.com/palantir/pkg/specdir"

	"github.com/palantir/godel/apps/distgo/params"
)

type Dister interface {
	NumArtifacts() int
	ArtifactPathsInOutputDir(buildSpec params.ProductBuildSpec) []string
	Dist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, spec specdir.LayoutSpec, values specdir.TemplateValues, stdout io.Writer) (Packager, error)
	DistPackageType() string
}

func ToDister(info params.DistInfo) Dister {
	switch info.Type() {
	default:
		panic(fmt.Errorf("unrecognized type: %v", info.Type()))
	case params.SLSDistType:
		return (*slsDister)(info.(*params.SLSDistInfo))
	case params.BinDistType:
		return (*binDister)(info.(*params.BinDistInfo))
	case params.ManualDistType:
		return (*manualDister)(info.(*params.ManualDistInfo))
	case params.OSArchBinDistType:
		return (*osArchsBinDister)(info.(*params.OSArchsBinDistInfo))
	case params.RPMDistType:
		return (*rpmDister)(info.(*params.RPMDistInfo))
	}
}

func FullArtifactsPaths(dister Dister, buildSpec params.ProductBuildSpec, distCfg params.Dist) []string {
	var outPaths []string
	for _, currPath := range dister.ArtifactPathsInOutputDir(buildSpec) {
		outPaths = append(outPaths, path.Join(buildSpec.ProjectDir, distCfg.OutputDir, currPath))
	}
	return outPaths
}
