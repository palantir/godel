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

package artifacts

import (
	"github.com/nmiyake/pkg/dirs"
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cfgcli"
	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/cmd"
	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/config"
	"github.com/palantir/godel/apps/distgo/params"
)

const (
	absPathFlagName       = "absolute"
	requiresBuildFlagName = "requires-build"
)

var (
	absPathFlag = flag.BoolFlag{
		Name:  absPathFlagName,
		Usage: "Print the absolute path for artifacts",
	}
	requiresBuildFlag = flag.BoolFlag{
		Name:  requiresBuildFlagName,
		Usage: "If true, only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)",
	}
)

func Command() cli.Command {
	return cli.Command{
		Name:  "artifacts",
		Usage: "Print the artifacts for products",
		Subcommands: []cli.Command{
			buildArtifactsCommand("build", "Print the paths to the build artifacts for products"),
			artifactsCommand("dist", "Print the paths to the distribution artifacts for products", distArtifactsAction),
			artifactsCommand("docker", "Print the labels for the Docker tags for products", dockerArtifactsAction),
		},
	}
}

func artifactsCommand(name, usage string, action artifactsAction) cli.Command {
	return cli.Command{
		Name:  name,
		Usage: usage,
		Flags: []flag.Flag{
			cmd.ProductsParam,
			absPathFlag,
		},
		Action: func(ctx cli.Context) error {
			cfg, err := config.Load(cfgcli.ConfigPath, cfgcli.ConfigJSON)
			if err != nil {
				return err
			}
			wd, err := dirs.GetwdEvalSymLinks()
			if err != nil {
				return err
			}

			specs, err := build.SpecsWithDepsForArgs(cfg, ctx.Slice(cmd.ProductsParamName), wd)
			if err != nil {
				return err
			}

			artifacts, err := action(ctx, specs, ctx.Bool(absPathFlagName))
			if err != nil {
				return err
			}

			for _, spec := range specs {
				if v, ok := artifacts[spec.Spec.ProductName]; ok {
					for _, k := range v.Keys() {
						for _, currPath := range v.Get(k) {
							ctx.Println(currPath)
						}
					}
				}
			}
			return nil
		},
	}
}

// buildArtifactsCommand returns the result of calling artifactsCommand for the "build" command and adding flags to the
// command that are only relevant for the "build" action.
func buildArtifactsCommand(name, usage string) cli.Command {
	buildCmd := artifactsCommand(name, usage, buildArtifactsAction)
	buildCmd.Flags = append(buildCmd.Flags, cmd.OSArchFlag, requiresBuildFlag)
	return buildCmd
}

type artifactsAction func(ctx cli.Context, specs []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringSliceMap, error)

func buildArtifactsAction(ctx cli.Context, specs []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringSliceMap, error) {
	osArchsFilter, err := cmd.NewOSArchFilter(ctx.String(cmd.OSArchFlagName))
	if err != nil {
		return nil, err
	}
	return BuildArtifacts(specs, BuildArtifactsParams{
		AbsPath:       absPath,
		RequiresBuild: ctx.Bool(requiresBuildFlagName),
		OSArchs:       osArchsFilter,
	})
}

func distArtifactsAction(ctx cli.Context, specs []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringSliceMap, error) {
	return DistArtifacts(specs, absPath)
}

func dockerArtifactsAction(ctx cli.Context, specs []params.ProductBuildSpecWithDeps, absPath bool) (map[string]OrderedStringSliceMap, error) {
	artifacts := make(map[string]OrderedStringSliceMap)
	for k, v := range DockerArtifacts(specs) {
		ordered := newOrderedStringSliceMap()
		ordered.PutValues(k, v)
		artifacts[k] = ordered
	}
	return artifacts, nil
}
