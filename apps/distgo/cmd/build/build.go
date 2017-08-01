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
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
	"github.com/palantir/godel/apps/distgo/pkg/script"
)

type buildUnit struct {
	buildSpec params.ProductBuildSpec
	osArch    osarch.OSArch
}

type Context struct {
	Parallel bool
	Install  bool
	Pkgdir   bool
}

func Products(products []string, osArchs cmd.OSArchFilter, buildCtx Context, cfg params.Project, wd string, stdout io.Writer) error {
	return RunBuildFunc(func(buildSpec []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		specs := make([]params.ProductBuildSpec, len(buildSpec))
		for i, curr := range buildSpec {
			specs[i] = curr.Spec
		}
		return Run(specs, osArchs, buildCtx, stdout)
	}, cfg, products, wd, stdout)
}

// Run builds all of the executables specified by buildSpecs using the mode specified in ctx. If ctx.Parallel is true,
// then the products will be built in parallel with N workers, where N is the number of logical processors reported by
// Go. When builds occur in parallel, each (Product, OSArch) pair is treated as an individual unit of work. Thus, it is
// possible that different products may be built in parallel. If any build process returns an error, the first error
// returned is propagated back (and any builds that have not started will not be started). If ctx.PkgDir is true, a
// custom per-OS/Arch "pkg" directory is used and the "install" command is run before build for each unit, which can
// speed up compilations on repeated runs by writing compiled packages to disk for reuse.
func Run(buildSpecs []params.ProductBuildSpec, osArchs cmd.OSArchFilter, ctx Context, stdout io.Writer) error {
	var units []buildUnit
	for _, currSpec := range distinct(buildSpecs) {
		// execute pre-build script
		distEnvVars := cmd.ScriptEnvVariables(currSpec, "")
		if err := script.WriteAndExecute(currSpec, currSpec.Build.Script, stdout, os.Stderr, distEnvVars); err != nil {
			return errors.Wrapf(err, "failed to execute build script for %v", currSpec.ProductName)
		}

		for _, currOSArch := range currSpec.Build.OSArchs {
			if osArchs.Matches(currOSArch) {
				units = append(units, buildUnit{
					buildSpec: currSpec,
					osArch:    currOSArch,
				})
			}
		}
	}

	if len(units) == 1 || !ctx.Parallel {
		// process serially
		for _, currUnit := range units {
			if err := executeBuild(stdout, currUnit.buildSpec, ctx, currUnit.osArch); err != nil {
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
			cs = append(cs, worker(stdout, buildUnitsJobs, ctx))
		}

		for err := range merge(done, cs...) {
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// ArtifactPaths returns a map that contains the paths to the executables created by the provided spec. The keys in the
// map are the OS/architecture of the executable, and the value is the output path for the executable for that
// OS/architecture.
func ArtifactPaths(buildSpec params.ProductBuildSpec) map[osarch.OSArch]string {
	paths := make(map[osarch.OSArch]string)
	for _, osArch := range buildSpec.Build.OSArchs {
		paths[osArch] = path.Join(buildSpec.ProjectDir, buildSpec.Build.OutputDir, buildSpec.VersionInfo.Version, osArch.String(), ExecutableName(buildSpec.ProductName, osArch.OS))
	}
	return paths
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

func worker(stdout io.Writer, in <-chan buildUnit, ctx Context) <-chan error {
	out := make(chan error)
	go func() {
		for unit := range in {
			out <- executeBuild(stdout, unit.buildSpec, ctx, unit.osArch)
		}
		close(out)
	}()
	return out
}

func executeBuild(stdout io.Writer, buildSpec params.ProductBuildSpec, ctx Context, osArch osarch.OSArch) error {
	name := buildSpec.ProductName

	if buildSpec.Build.Skip {
		fmt.Fprintf(stdout, "Skipping build for %s because skip configuration for product is true\n", name)
		return nil
	}

	start := time.Now()
	outputArtifactPath, ok := ArtifactPaths(buildSpec)[osArch]
	if !ok {
		return fmt.Errorf("failed to determine artifact path for %s for %s", name, osArch.String())
	}
	currOutputDir := path.Dir(outputArtifactPath)
	fmt.Fprintf(stdout, "Building %s for %s at %s\n", name, osArch.String(), path.Join(currOutputDir, name))

	if err := os.MkdirAll(currOutputDir, 0755); err != nil {
		return errors.Wrapf(err, "failed to create directories for %s", currOutputDir)
	}
	if ctx.Install {
		if err := doBuildAction(doInstall, buildSpec, "", osArch, ctx.Pkgdir); err != nil {
			return fmt.Errorf("go install failed: %v", err)
		}
	}
	if err := doBuildAction(doBuild, buildSpec, currOutputDir, osArch, ctx.Pkgdir); err != nil {
		return errors.Wrapf(err, "go build failed")
	}

	elapsed := time.Since(start)
	fmt.Fprintf(stdout, "Finished building %s for %s (%.3fs)\n", name, osArch.String(), elapsed.Seconds())

	return nil
}

type buildAction int

const (
	doBuild buildAction = iota
	doInstall
)

func doBuildAction(action buildAction, buildSpec params.ProductBuildSpec, outputDir string, osArch osarch.OSArch, pkgdir bool) error {
	cmd := exec.Command("go")
	cmd.Dir = buildSpec.ProjectDir

	var env []string
	goos := runtime.GOOS
	if osArch.OS != "" {
		env = append(env, "GOOS="+osArch.OS)
		goos = osArch.OS
	}
	goarch := runtime.GOARCH
	if osArch.Arch != "" {
		env = append(env, "GOARCH="+osArch.Arch)
		goarch = osArch.Arch
	}
	for k, v := range buildSpec.Build.Environment {
		env = append(env, fmt.Sprintf("%v=%v", k, v))
	}
	cmd.Env = append(os.Environ(), env...)

	args := []string{cmd.Path}
	switch action {
	case doBuild:
		args = append(args, "build")
		args = append(args, "-o", path.Join(outputDir, ExecutableName(buildSpec.ProductName, goos)))
	case doInstall:
		args = append(args, "install")
	default:
		return errors.Errorf("unrecognized action: %v", action)
	}

	// get build args
	buildArgs, err := script.GetBuildArgs(buildSpec, buildSpec.Build.BuildArgsScript)
	if err != nil {
		return err
	}
	args = append(args, buildArgs...)

	if buildSpec.Build.VersionVar != "" {
		args = append(args, "-ldflags", fmt.Sprintf("-X %v=%v", buildSpec.Build.VersionVar, buildSpec.ProductVersion))
	}

	if pkgdir {
		// specify custom pkgdir if isolation of packages is desired
		args = append(args, "-pkgdir", fmt.Sprintf("%v/pkg/_%v_%v", os.Getenv("GOPATH"), goos, goarch))
	}
	args = append(args, buildSpec.Build.MainPkg)
	cmd.Args = args

	if output, err := cmd.CombinedOutput(); err != nil {
		errOutput := strings.TrimSpace(string(output))
		err = fmt.Errorf("build command %v run with additional environment variables %v failed with output:\n%s", cmd.Args, env, errOutput)

		if action == doInstall && regexp.MustCompile(installPermissionDenied).MatchString(errOutput) {
			// if "install" command failed due to lack of permissions, return error that contains explanation
			return fmt.Errorf(goInstallErrorMsg(osArch, err))
		}
		return err
	}
	return nil
}

const installPermissionDenied = `^go install [a-zA-Z0-9_/]+: mkdir .+: permission denied$`

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

func distinct(buildSpecs []params.ProductBuildSpec) []params.ProductBuildSpec {
	distinctSpecs := make([]params.ProductBuildSpec, 0, len(buildSpecs))
	for _, spec := range buildSpecs {
		if contains(distinctSpecs, spec) {
			continue
		}
		distinctSpecs = append(distinctSpecs, spec)
	}
	return distinctSpecs
}

func contains(specs []params.ProductBuildSpec, spec params.ProductBuildSpec) bool {
	for _, currSpec := range specs {
		if reflect.DeepEqual(currSpec, spec) {
			return true
		}
	}
	return false
}
