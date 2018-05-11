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
	"fmt"
	"io"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/palantir/godel/pkg/osarch"
	"github.com/pkg/errors"

	"github.com/palantir/distgo/distgo"
)

type buildUnit struct {
	buildParam            distgo.BuildParam
	productTaskOutputInfo distgo.ProductTaskOutputInfo
	osArch                osarch.OSArch
}

type Options struct {
	Parallel bool
	Install  bool
	DryRun   bool
}

func Products(projectInfo distgo.ProjectInfo, projectParam distgo.ProjectParam, productBuildIDs []distgo.ProductBuildID, buildOpts Options, stdout io.Writer) error {
	productParams, err := distgo.ProductParamsForBuildProductArgs(projectParam.Products, productBuildIDs...)
	if err != nil {
		return err
	}
	if err := Run(projectInfo, productParams, buildOpts, stdout); err != nil {
		return err
	}
	return nil
}

// Run builds the executables for the products specified by productParams using the options specified in buildOpts. If
// buildOpts.Parallel is true, then the products will be built in parallel with N workers, where N is the number of
// logical processors reported by Go. When builds occur in parallel, each (Product, OSArch) pair is treated as an
// individual unit of work. Thus, it is possible that different products may be built in parallel. If any build process
// returns an error, the first error returned is propagated back (and any builds that have not started will not be
// started).
func Run(projectInfo distgo.ProjectInfo, productParams []distgo.ProductParam, buildOpts Options, stdout io.Writer) error {
	var units []buildUnit
	for _, currProductParam := range productParams {
		currProductTaskOutputInfo, err := distgo.ToProductTaskOutputInfo(projectInfo, currProductParam)
		if err != nil {
			return errors.Wrapf(err, "failed to compute output information for %s", currProductParam.ID)
		}
		if currProductParam.Build == nil {
			continue
		}
		for _, currOSArch := range currProductParam.Build.OSArchs {
			units = append(units, buildUnit{
				buildParam:            *currProductParam.Build,
				productTaskOutputInfo: currProductTaskOutputInfo,
				osArch:                currOSArch,
			})
		}
	}

	if len(units) == 1 || !buildOpts.Parallel {
		// process serially
		for _, currUnit := range units {
			if err := executeBuild(currUnit, buildOpts, stdout); err != nil {
				return err
			}
		}
	} else {
		done := make(chan struct{})
		defer close(done)

		// send all jobs
		nUnits := len(units)
		buildUnitsJobs := make(chan buildUnit, nUnits)
		for _, currUnit := range units {
			buildUnitsJobs <- currUnit
		}
		close(buildUnitsJobs)

		// create workers
		nWorkers := runtime.NumCPU()
		if nUnits < nWorkers {
			nWorkers = nUnits
		}
		var cs []<-chan error
		for i := 0; i < nWorkers; i++ {
			cs = append(cs, worker(buildUnitsJobs, buildOpts, stdout))
		}

		for err := range merge(done, cs...) {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// merge handles "fanning in" the result of multiple output channels into a single output channel. If a signal is
// received on the "done" channel, output processing will stop.
func merge(done <-chan struct{}, cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	out := make(chan error)

	output := func(c <-chan error) {
		defer wg.Done()
		for err := range c {
			select {
			case out <- err:
			case <-done:
				return
			}
		}
	}

	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func worker(in <-chan buildUnit, buildOpts Options, stdout io.Writer) <-chan error {
	out := make(chan error)
	go func() {
		for unit := range in {
			out <- executeBuild(unit, buildOpts, stdout)
		}
		close(out)
	}()
	return out
}

func executeBuild(unit buildUnit, buildOpts Options, stdout io.Writer) error {
	name := unit.productTaskOutputInfo.Product.ID

	osArch := unit.osArch
	start := time.Now()
	outputArtifactPath, ok := unit.productTaskOutputInfo.ProductBuildArtifactPaths()[osArch]
	if !ok {
		return fmt.Errorf("failed to determine artifact path for %s for %s", name, osArch.String())
	}
	outputArtifactDisplayPath := outputArtifactPath
	if wd, err := os.Getwd(); err == nil {
		if relPath, err := filepath.Rel(wd, outputArtifactPath); err == nil {
			outputArtifactDisplayPath = relPath
		}
	}
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Building %s for %s at %s", name, osArch.String(), outputArtifactDisplayPath), buildOpts.DryRun)

	if !buildOpts.DryRun {
		if err := os.MkdirAll(path.Dir(outputArtifactPath), 0755); err != nil {
			return errors.Wrapf(err, "failed to create directories for %s", path.Dir(outputArtifactPath))
		}
	}
	if err := doBuildAction(unit, outputArtifactPath, buildOpts.Install, buildOpts.DryRun, stdout); err != nil {
		return errors.Wrapf(err, "go build failed")
	}

	elapsed := time.Since(start)
	distgo.PrintlnOrDryRunPrintln(stdout, fmt.Sprintf("Finished building %s for %s (%.3fs)", name, osArch.String(), elapsed.Seconds()), buildOpts.DryRun)
	return nil
}

func doBuildAction(unit buildUnit, outputArtifactPath string, doInstall, dryRun bool, stdout io.Writer) error {
	osArch := unit.osArch

	cmd := exec.Command("go")
	cmd.Dir = unit.productTaskOutputInfo.Project.ProjectDir

	var env []string
	if osArch.OS != "" {
		env = append(env, "GOOS="+osArch.OS)
	}
	if osArch.Arch != "" {
		env = append(env, "GOARCH="+osArch.Arch)
	}
	for k, v := range unit.buildParam.Environment {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Env = append(os.Environ(), env...)

	args := []string{cmd.Path}
	args = append(args, "build")
	if doInstall {
		args = append(args, "-i")
	}

	if !path.IsAbs(outputArtifactPath) {
		// if outputArtifactPath is relative, then if it starts with ProjectDir the prefix needs to be trimmed because
		// the working directory for the build command is set to the project directory
		outputArtifactPath = strings.TrimPrefix(outputArtifactPath, path.Clean(unit.productTaskOutputInfo.Project.ProjectDir)+"/")
	}
	args = append(args, "-o", outputArtifactPath)

	buildArgs, err := unit.buildParam.BuildArgs(unit.productTaskOutputInfo)
	if err != nil {
		return err
	}
	args = append(args, buildArgs...)

	mainPkg := unit.buildParam.MainPkg
	args = append(args, mainPkg)
	cmd.Args = args

	if dryRun {
		dryRunMsg := fmt.Sprintf("Run: %s", strings.Join(cmd.Args, " "))
		if len(env) > 0 {
			dryRunMsg += fmt.Sprintf(" with additional environment variables %v", env)
		}
		distgo.DryRunPrintln(stdout, dryRunMsg)
	} else {
		if output, err := cmd.CombinedOutput(); err != nil {
			errOutput := strings.TrimSpace(string(output))
			err = fmt.Errorf("build command %v run in directory %s with additional environment variables %v failed with output:\n%s", cmd.Args, cmd.Dir, env, errOutput)
			if regexp.MustCompile(installPermissionDenied).MatchString(errOutput) {
				// if "install" command failed due to lack of permissions, return error that contains explanation
				return fmt.Errorf(goInstallErrorMsg(osArch, err))
			}
			return err
		}
	}
	return nil
}

const installPermissionDenied = `(?s)^go build [a-zA-Z0-9_/]+: mkdir [^:]+: permission denied.+`

func goInstallErrorMsg(osArch osarch.OSArch, err error) string {
	goBinary := "go"
	if output, err := exec.Command("command", "-v", "go").CombinedOutput(); err == nil {
		goBinary = strings.TrimSpace(string(output))
	}
	return strings.Join([]string{
		`failed to install a Go standard library package due to insufficient permissions to create directory.`,
		`This typically means that the standard library for the OS/architecture combination have not been installed locally and the current user does not have write permissions to GOROOT/pkg.`,
		fmt.Sprintf(`Run "sudo env GOOS=%s GOARCH=%s %s install std" to install the standard packages for this combination as root and then try again.`, osArch.OS, osArch.Arch, goBinary),
		fmt.Sprintf(`Full error: %s`, err.Error()),
	}, "\n")
}
