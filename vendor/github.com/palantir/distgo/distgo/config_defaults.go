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
	"os/exec"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

func mainPkgsProductsConfig(projectDir string, defaultDisterCfg DisterConfig, exclude matcher.Matcher) (map[ProductID]ProductConfig, error) {
	mainPkgPaths, err := mainPkgPaths(projectDir, exclude)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine paths to main packages in %s", projectDir)
	}

	if len(mainPkgPaths) == 0 {
		return nil, nil
	}

	mainPkgPathToProductID := make(map[string]ProductID)
	mainProductIDs := make(map[ProductID]struct{})
	for _, currMainPkgPath := range mainPkgPaths {
		currMainPkgProjectPath := path.Join(projectDir, currMainPkgPath)
		if currMainPkgProjectPath == "." {
			absPath, err := filepath.Abs(currMainPkgProjectPath)
			if err != nil {
				return nil, errors.Wrapf(err, "failed to convert to absolute path")
			}
			currMainPkgProjectPath = absPath
		}
		productID := ProductID(path.Base(currMainPkgProjectPath))

		mainPkgPathToProductID[currMainPkgPath] = productID
		mainProductIDs[productID] = struct{}{}
	}

	usedProductIDs := make(map[ProductID]struct{})
	productsCfg := make(map[ProductID]ProductConfig)
	for _, currMainPkgPath := range mainPkgPaths {
		// redeclare locally so address can be taken
		currMainPkgPath := currMainPkgPath
		productID := uniqueProductID(mainPkgPathToProductID[currMainPkgPath], mainProductIDs, usedProductIDs)
		productsCfg[productID] = ProductConfig{
			Build: &BuildConfig{
				MainPkg: &currMainPkgPath,
			},
			Run: &RunConfig{},
			Dist: &DistConfig{
				Disters: &DistersConfig{
					DistID(*defaultDisterCfg.Type): defaultDisterCfg,
				},
			},
			Publish: &PublishConfig{},
			Docker:  &DockerConfig{},
		}
	}
	return productsCfg, nil
}

func combineProductIDSets(one, two map[ProductID]struct{}) map[ProductID]struct{} {
	if len(one) == 0 {
		return two
	}
	if len(two) == 0 {
		return one
	}
	out := make(map[ProductID]struct{})
	for k := range one {
		out[k] = struct{}{}
	}
	for k := range two {
		out[k] = struct{}{}
	}
	return out
}

func uniqueProductID(candidate ProductID, primaryIDs, used map[ProductID]struct{}) (rVal ProductID) {
	// add returned value to "used" map
	defer func() {
		used[rVal] = struct{}{}
	}()

	if _, ok := used[candidate]; !ok {
		// name is unique
		return candidate
	}

	// current name has already been used: create a unique one
	idx := strings.LastIndex(string(candidate), "-")
	if idx == -1 {
		// existing name does not have a hyphen, so doesn't conform to naming scheme: add hyphen and number
		return ProductID(nextAvailableNumName(string(candidate)+"-", 1, combineProductIDSets(used, primaryIDs)))
	}

	lastPortion := string(candidate[idx+1:])
	lastPortionNum, err := strconv.Atoi(lastPortion)
	if err != nil {
		// the portion of the name after the hyphen cannot be parsed as a number, so doesn't conform to naming scheme: add hyphen and number
		return ProductID(nextAvailableNumName(string(candidate)+"-", 1, combineProductIDSets(used, primaryIDs)))
	}

	// existing name ends with a hyphen and number -- increment number after hyphen to next available one
	return ProductID(nextAvailableNumName(string(candidate[:idx])+"-", lastPortionNum, combineProductIDSets(used, primaryIDs)))
}

func nextAvailableNumName(nameWithHyphen string, currNum int, used map[ProductID]struct{}) string {
	var currName string
	for {
		currName = nameWithHyphen + strconv.Itoa(currNum)
		if _, ok := used[ProductID(currName)]; !ok {
			// current name is not used
			break
		}
		currNum++
	}
	return currName
}

func mainPkgPaths(projectDir string, exclude matcher.Matcher) ([]string, error) {
	projectPkgOutput, err := runGoList(projectDir, "-e")
	if err != nil {
		return nil, err
	}
	projectBasePkg := projectPkgOutput[0]

	allProjectPkgsOutput, err := runGoList(projectDir, "-f", "{{.Name}} {{.ImportPath}}", "./...")
	if err != nil {
		return nil, err
	}

	var mainPkgPaths []string
	for _, currPkgOutput := range allProjectPkgsOutput {
		firstSpaceIdx := strings.Index(currPkgOutput, " ")
		if firstSpaceIdx == -1 {
			return nil, errors.Errorf("failed to find space in output %q", currPkgOutput)
		}
		if currPkgOutput[:firstSpaceIdx] != "main" {
			continue
		}
		currPkgRelPath, err := filepath.Rel(projectBasePkg, currPkgOutput[firstSpaceIdx+1:])
		if err != nil {
			return nil, errors.Wrapf(err, "failed to convert package Path to relative Path")
		}
		if exclude != nil && exclude.Match(currPkgRelPath) {
			continue
		}
		mainPkgPaths = append(mainPkgPaths, currPkgRelPath)
	}
	sort.Strings(mainPkgPaths)
	return mainPkgPaths, nil
}

func runGoList(dir string, args ...string) ([]string, error) {
	goListCmd := exec.Command("go", append([]string{"list"}, args...)...)
	goListCmd.Dir = dir
	outputBytes, err := goListCmd.CombinedOutput()
	output := string(outputBytes)
	if err != nil {
		return nil, errors.Wrapf(err, "command %v run in directory %s failed with outputBytes %q", goListCmd.Args, dir, output)
	}
	return strings.Split(strings.TrimSpace(output), "\n"), nil
}
