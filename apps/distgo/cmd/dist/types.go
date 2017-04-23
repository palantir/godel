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
	ArtifactPathInOutputDir(buildSpec params.ProductBuildSpec, distCfg params.Dist) string
	Dist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, spec specdir.LayoutSpec, values specdir.TemplateValues, stdout io.Writer) (Packager, error)
	DistPackageType() string
}

func DisterForType(typ params.DistInfoType) Dister {
	dister, ok := distInfoTypeDataMap[typ]
	if !ok {
		panic(fmt.Sprintf("unrecognized type: %v", typ))
	}
	return dister
}

func FullArtifactPath(distType params.DistInfoType, buildSpec params.ProductBuildSpec, distCfg params.Dist) string {
	dister := DisterForType(distType)
	return path.Join(buildSpec.ProjectDir, distCfg.OutputDir, dister.ArtifactPathInOutputDir(buildSpec, distCfg))
}

var distInfoTypeDataMap = map[params.DistInfoType]Dister{
	params.SLSDistType:       &slsDistStruct{},
	params.BinDistType:       &binDistStruct{},
	params.OSArchBinDistType: &osArchBinDistStruct{},
	params.RPMDistType:       &rpmDistStruct{},
}
