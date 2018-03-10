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
	"os"
	"time"

	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

// RequiresDist returns a pointer to a distgo.ProductParam that contains only the Dister parameters for the output dist
// artifacts that require generation. A product is considered to require generating dist artifacts if any of the
// following is true:
//   * Any of the dist artifact output paths do not exist
//   * The product has dependencies and any of the dependent build or dist artifacts are newer (have a later
//     modification date) than any of the dist artifacts for the provided product
//   * The product does not define a dist configuration
//
// Returns nil if all of the outputs exist and are up-to-date.
func RequiresDist(projectInfo distgo.ProjectInfo, productParam distgo.ProductParam) (*distgo.ProductParam, error) {
	if productParam.Dist == nil {
		return nil, nil
	}
	productTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, productParam)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to compute output information for %s", productParam.ID)
	}

	requiresDistIDs := make(map[distgo.DistID]struct{})
	for _, currDistID := range productTaskOutputInfo.Product.DistOutputInfos.DistIDs {
		if !disterRequiresDist(currDistID, productTaskOutputInfo) {
			continue
		}
		requiresDistIDs[currDistID] = struct{}{}
	}

	if len(requiresDistIDs) == 0 {
		return nil, nil
	}
	requiresDistParams := make(map[distgo.DistID]distgo.DisterParam)
	for distID, distParam := range productParam.Dist.DistParams {
		if _, ok := requiresDistIDs[distID]; !ok {
			continue
		}
		requiresDistParams[distID] = distParam
	}
	productParam.Dist.DistParams = requiresDistParams
	return &productParam, nil
}

func disterRequiresDist(distID distgo.DistID, productTaskOutputInfo distgo.ProductTaskOutputInfo) bool {
	// determine oldest dist artifact for current Dister. If any artifact is missing, Dister needs to be run.
	oldestDistTime := time.Now()
	for _, currArtifactPath := range productTaskOutputInfo.ProductDistWorkDirsAndArtifactPaths()[distID] {
		fi, err := os.Stat(currArtifactPath)
		if err != nil {
			return true
		}
		if fiModTime := fi.ModTime(); fiModTime.Before(oldestDistTime) {
			oldestDistTime = fiModTime
		}
	}

	// if any dependent artifact (build or dist) is newer than the oldest dist artifact, consider dist artifact out-of-date
	for _, depProductOutputInfo := range productTaskOutputInfo.Deps {
		newestDependencyTime := newestDistArtifactModTime(productTaskOutputInfo.Project, distID, depProductOutputInfo)
		if newestDependencyTime != nil && newestDependencyTime.After(oldestDistTime) {
			return true
		}
	}
	return false
}

func newestDistArtifactModTime(projectInfo distgo.ProjectInfo, distID distgo.DistID, productInfo distgo.ProductOutputInfo) *time.Time {
	var newestModTime *time.Time
	newestModTimeFn := func(currPath string) {
		fi, err := os.Stat(currPath)
		if err != nil {
			return
		}
		if fiModTime := fi.ModTime(); newestModTime == nil || fiModTime.After(*newestModTime) {
			newestModTime = &fiModTime
		}
	}
	if productInfo.BuildOutputInfo != nil {
		for _, v := range distgo.ProductBuildArtifactPaths(projectInfo, productInfo) {
			newestModTimeFn(v)
		}
	}
	if productInfo.DistOutputInfos != nil {
		for _, v := range distgo.ProductDistWorkDirsAndArtifactPaths(projectInfo, productInfo)[distID] {
			newestModTimeFn(v)
		}
	}
	return newestModTime
}
