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
	"github.com/termie/go-shutil"

	"github.com/palantir/godel/apps/distgo/params"
)

type manualDister params.ManualDistInfo

func (m *manualDister) NumArtifacts() int {
	return 1
}

func (m *manualDister) ArtifactPathsInOutputDir(buildSpec params.ProductBuildSpec) []string {
	extension := ""
	if m.Extension != "" {
		extension = fmt.Sprintf(".%s", m.Extension)
	}
	return []string{fmt.Sprintf("%s-%s%s", buildSpec.ProductName, buildSpec.ProductVersion, extension)}
}

func (m *manualDister) Dist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, spec specdir.LayoutSpec, values specdir.TemplateValues, stdout io.Writer) (Packager, error) {
	// output of manual dist is already packaged and in the expected location
	return packager(func() error {
		srcPath := path.Join(outputProductDir, m.ArtifactPathsInOutputDir(buildSpecWithDeps.Spec)[0])
		dstPath := FullArtifactsPaths(m, buildSpecWithDeps.Spec, distCfg)[0]
		_, err := shutil.Copy(srcPath, dstPath, false)
		return err
	}), nil
}

func (m *manualDister) DistPackageType() string {
	return m.Extension
}
