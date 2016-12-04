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
	"os/exec"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/script"
)

const defaultRPMRelease = "1"

func checkRPMDependencies() error {
	var missing []string
	if _, err := exec.LookPath("fpm"); err != nil {
		missing = append(missing, "Missing `fpm` command required to build RPMs. Install with `gem install fpm`.")
	}
	if _, err := exec.LookPath("rpmbuild"); err != nil {
		missing = append(missing, "Missing `rpmbuild` command required to build RPMs. Install with `yum install rpm-build` or `apt-get install rpm` or `brew install rpm`.")
	}
	if len(missing) > 0 {
		return errors.New(strings.Join(missing, "\n"))
	}
	return nil
}

func rpmDist(buildSpecWithDeps params.ProductBuildSpecWithDeps, distCfg params.Dist, outputProductDir string, stdout io.Writer) (p Packager, rErr error) {
	buildSpec := buildSpecWithDeps.Spec
	rpmDistInfo, ok := distCfg.Info.(*params.RPMDistInfo)
	if !ok {
		rpmDistInfo = &params.RPMDistInfo{}
		distCfg.Info = rpmDistInfo
	}

	release := defaultRPMRelease
	if rpmDistInfo.Release != "" {
		release = rpmDistInfo.Release
	}

	// These are run after the cmd is executed.
	var cleanups []func() error
	// clean up unless everything below succeeds
	runCleanups := true
	defer func() {
		if runCleanups {
			var errs []string
			for _, cleanup := range cleanups {
				if err := cleanup(); err != nil {
					errs = append(errs, err.Error())
				}
			}
			if len(errs) > 0 {
				rErr = errors.Errorf(strings.Join(append([]string{"encountered errors during cleanup:"}, errs...), "\n"))
			}
		}
	}()

	cmd := exec.Command("fpm")
	cmd.Dir = buildSpec.ProjectDir
	cmd.Stdout = stdout
	cmd.Stderr = os.Stderr

	cmd.Args = []string{
		"fpm",
		"-t", "rpm",
		"-n", buildSpec.ProductName,
		"-v", buildSpec.ProductVersion,
		"--iteration", release,
		"-p", ArtifactPath(buildSpec, distCfg),
		"-s", "dir",
		"-C", outputProductDir,
	}

	for _, configFile := range rpmDistInfo.ConfigFiles {
		cmd.Args = append(cmd.Args, "--config-files", configFile)
	}

	scriptArg := func(name string, content string) error {
		if content == "" {
			return nil
		}
		f, cleanup, err := script.Write(buildSpec, rpmDistInfo.BeforeInstallScript)
		if err != nil {
			return errors.Wrapf(err, "failed to write %v script for %v", name, buildSpec.ProductName)
		}
		cleanups = append(cleanups, cleanup)
		cmd.Args = append(cmd.Args, "--"+name, f)
		return nil
	}
	if err := scriptArg("before-install", rpmDistInfo.BeforeInstallScript); err != nil {
		return nil, err
	}
	if err := scriptArg("after-install", rpmDistInfo.AfterInstallScript); err != nil {
		return nil, err
	}
	if err := scriptArg("after-remove", rpmDistInfo.AfterRemoveScript); err != nil {
		return nil, err
	}

	runCleanups = false
	return packager(func() error {
		err := cmd.Run()
		for _, cleanup := range cleanups {
			if err := cleanup(); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
		return err
	}), nil
}
