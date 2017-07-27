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
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/nmiyake/archiver"
	"github.com/palantir/pkg/specdir"
	"github.com/pkg/errors"
	"github.com/termie/go-shutil"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/apps/distgo/pkg/script"
	"github.com/palantir/godel/apps/distgo/pkg/slsspec"
)

func RequiresDist(products []string, cfg params.Project, wd string) ([]string, error) {
	buildSpecsWithDeps, err := build.SpecsWithDepsForArgs(cfg, products, wd)
	if err != nil {
		return nil, err
	}
	var productsSlice []string
	productsMap := make(map[string]struct{})
	for _, curSpecWithDeps := range buildSpecsWithDeps {
		for _, curDistCfg := range curSpecWithDeps.Spec.Dist {
			artifactPaths := FullArtifactsPaths(ToDister(curDistCfg.Info), curSpecWithDeps.Spec, curDistCfg)
			for _, artifactPath := range artifactPaths {
				if _, err := os.Stat(artifactPath); err == nil {
					buildInfo := build.RequiresBuild(curSpecWithDeps, nil)
					if len(buildInfo.Specs()) == 0 {
						continue
					}
				}
				productsMap[curSpecWithDeps.Spec.ProductName] = struct{}{}
			}
		}
	}
	for product := range productsMap {
		productsSlice = append(productsSlice, product)
	}
	return productsSlice, nil
}

func Products(products []string, cfg params.Project, forceBuild bool, wd string, stdout io.Writer) error {
	return build.RunBuildFunc(func(buildSpecWithDeps []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		var specsToBuild []params.ProductBuildSpec
		for _, currSpecWithDeps := range buildSpecWithDeps {
			if forceBuild {
				specsToBuild = append(specsToBuild, currSpecWithDeps.AllSpecs()...)
			} else {
				specsToBuild = append(specsToBuild, build.RequiresBuild(currSpecWithDeps, nil).Specs()...)
			}
		}
		if len(specsToBuild) > 0 {
			if err := build.Run(specsToBuild, nil, build.DefaultContext(), stdout); err != nil {
				return errors.Wrapf(err, "Failed to build products required for dist")
			}
		}
		return cmd.ProcessSerially(Run)(buildSpecWithDeps, stdout)
	}, cfg, products, wd, stdout)
}

// Run produces a directory and artifacts (such as a .tgz or .rpm) for the specified product using the specified build
// specification. The binaries for the distribution must already exist in the expected locations. The distribution
// directory and artifact are written to the directory specified by "buildSpecWithDeps.Spec.DistCfgs.*.OutputDir".
func Run(buildSpecWithDeps params.ProductBuildSpecWithDeps, stdout io.Writer) error {
	// verify that required build outputs exist
	missingBinaries := build.RequiresBuild(buildSpecWithDeps, nil).Specs()
	if len(missingBinaries) > 0 {
		missingProducts := make([]string, len(missingBinaries))
		for i, currSpec := range missingBinaries {
			missingProducts[i] = currSpec.ProductName
		}
		return errors.Errorf("required output not present for build specs: %v", missingProducts)
	}

	buildSpec := buildSpecWithDeps.Spec
	for _, currDistCfg := range buildSpec.Dist {
		if currDistCfg.Info.Type() == params.RPMDistType {
			osArchs := buildSpec.Build.OSArchs
			expected := osarch.OSArch{OS: "linux", Arch: "amd64"}
			if len(osArchs) != 1 || osArchs[0] != expected {
				return fmt.Errorf("RPM is only supported for %v", expected)
			}
			if err := checkRPMDependencies(); err != nil {
				return err
			}
		}

		outputDir := path.Join(buildSpec.ProjectDir, currDistCfg.OutputDir)
		artifactPaths := FullArtifactsPaths(ToDister(currDistCfg.Info), buildSpec, currDistCfg)
		fmt.Fprintf(stdout, "Creating distribution for %s at %v\n", buildSpec.ProductName, strings.Join(artifactPaths, ", "))

		spec := slsspec.New()
		values := slsspec.TemplateValues(buildSpec.ProductName, buildSpec.ProductVersion)

		// remove output directory if it already exists
		outputProductDir := path.Join(outputDir, spec.RootDirName(values))
		if err := os.RemoveAll(outputProductDir); err != nil {
			return errors.Wrapf(err, "Failed to remove directory %v", outputProductDir)
		}

		// create output root directory
		if err := os.MkdirAll(outputProductDir, 0755); err != nil {
			return errors.Wrapf(err, "failed to create directories for %v", outputProductDir)
		}

		// if input directory is specified, copy its contents
		if currDistCfg.InputDir != "" {
			inputDir := path.Join(buildSpec.ProjectDir, currDistCfg.InputDir)

			fileInfos, err := ioutil.ReadDir(inputDir)
			if err != nil {
				return errors.Wrapf(err, "failed to list files in directory %v", inputDir)
			}

			for _, currFileInfo := range fileInfos {
				currFileName := currFileInfo.Name()
				srcPath := path.Join(inputDir, currFileName)
				dstPath := path.Join(outputProductDir, currFileName)

				if currFileInfo.IsDir() {
					if err := shutil.CopyTree(srcPath, dstPath, &shutil.CopyTreeOptions{
						CopyFunction: shutil.Copy,
						// do not copy ".gitkeep" files
						Ignore: func(dir string, files []os.FileInfo) []string {
							return []string{".gitkeep"}
						},
					}); err != nil {
						return errors.Wrapf(err, "failed to copy directory %v", currFileName)
					}
				} else if currFileName != ".gitkeep" {
					if _, err := shutil.Copy(srcPath, dstPath, false); err != nil {
						return errors.Wrapf(err, "failed to copy directory %v", currFileName)
					}
				}
			}
		}

		packager, err := ToDister(currDistCfg.Info).Dist(buildSpecWithDeps, currDistCfg, outputProductDir, spec, values, stdout)
		if err != nil {
			return err
		}

		// execute dist script
		distEnvVars := cmd.ScriptEnvVariables(buildSpec, outputProductDir)
		if err := script.WriteAndExecute(buildSpec, currDistCfg.Script, stdout, os.Stderr, distEnvVars); err != nil {
			return errors.Wrapf(err, "failed to execute dist script for %v", buildSpec.ProductName)
		}

		// create artifact for distribution
		if err := packager.Package(); err != nil {
			return errors.Wrapf(err, "failed to create artifact for %v from path %v", buildSpec.ProductName, outputProductDir)
		}

		fmt.Fprintf(stdout, "Finished creating distribution for %v\n", buildSpec.ProductName)
	}

	return nil
}

func singlePathTGZPackager(dstArtifactPath, pathToPackage string) packager {
	return tgzPackager([]string{dstArtifactPath}, map[string][]string{
		dstArtifactPath: {pathToPackage},
	})
}

func tgzPackager(dstArtifactPaths []string, dstToContentPaths map[string][]string) packager {
	return packager(func() error {
		for _, currDstPath := range dstArtifactPaths {
			if err := archiver.TarGz(currDstPath, dstToContentPaths[currDstPath]); err != nil {
				return err
			}
		}
		return nil
	})
}

func copyBuildArtifactsToBinDir(buildSpecWithDeps params.ProductBuildSpecWithDeps, binSpecDir specdir.SpecDir) error {
	buildSpec := buildSpecWithDeps.Spec

	// copy build artifacts for primary product
	if err := copyBuildArtifacts(buildSpec, binSpecDir); err != nil {
		return errors.Wrapf(err, "failed to copy build artifacts for %v", buildSpec.ProductName)
	}

	// copy build artifacts for dependent products
	for _, currDepSpec := range buildSpecWithDeps.Deps {
		if err := copyBuildArtifacts(currDepSpec, binSpecDir); err != nil {
			return errors.Wrapf(err, "failed to copy build artifacts for %v", currDepSpec.ProductName)
		}
	}

	return nil
}

func copyBuildArtifacts(buildSpec params.ProductBuildSpec, binSpecDir specdir.SpecDir) error {
	artifactPaths := build.ArtifactPaths(buildSpec)
	for _, currOSArch := range buildSpec.Build.OSArchs {
		currBuildArtifact, ok := artifactPaths[currOSArch]
		if !ok {
			return fmt.Errorf("could not determine artifact path for %s for %s", buildSpec.ProductName, currOSArch.String())
		}
		if binOSArchDir := binSpecDir.Path(currOSArch.String()); binOSArchDir != "" {
			dst := path.Join(binOSArchDir, build.ExecutableName(buildSpec.ProductName, currOSArch.OS))
			if _, err := shutil.Copy(currBuildArtifact, dst, false); err != nil {
				return errors.Wrapf(err, "failed to copy build artifact from %v to %v", currBuildArtifact, dst)
			}
		}
	}
	return nil
}
