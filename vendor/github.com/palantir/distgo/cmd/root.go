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

package cmd

import (
	"os"
	"path"
	"time"

	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/spf13/cobra"

	"github.com/palantir/distgo/assetapi"
	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/dockerbuilder"
	"github.com/palantir/distgo/publisher"
)

var (
	projectDirFlagVal       string
	distgoConfigFileFlagVal string
	godelConfigFileFlagVal  string
	assetsFlagVal           []string

	cliDisterFactory        distgo.DisterFactory
	cliDefaultDisterCfg     distgo.DisterConfig
	cliDockerBuilderFactory distgo.DockerBuilderFactory
)

var RootCmd = &cobra.Command{
	Use: "distgo",
}

func InitAssetCmds(args []string) error {
	// parse the flags to retrieve the value of the "--assets" flag. Ignore any errors that occur in flag parsing so
	// that, if provided flags are invalid, the regular logic handles the error printing.
	_ = RootCmd.ParseFlags(args)
	allAssets, err := assetapi.LoadAssets(assetsFlagVal)
	if err != nil {
		return err
	}

	// load publisher assets
	assetPublishers, err := publisher.AssetPublisherCreators(allAssets[assetapi.Publisher]...)
	if err != nil {
		return err
	}
	if err := publisher.SetPublishers(assetPublishers); err != nil {
		return err
	}

	// add publish commands based on assets
	addPublishSubcommands()

	return nil
}

func init() {
	pluginapi.AddProjectDirPFlagPtr(RootCmd.PersistentFlags(), &projectDirFlagVal)
	pluginapi.AddConfigPFlagPtr(RootCmd.PersistentFlags(), &distgoConfigFileFlagVal)
	pluginapi.AddGodelConfigPFlagPtr(RootCmd.PersistentFlags(), &godelConfigFileFlagVal)
	pluginapi.AddAssetsPFlagPtr(RootCmd.PersistentFlags(), &assetsFlagVal)

	RootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		allAssets, err := assetapi.LoadAssets(assetsFlagVal)
		if err != nil {
			return err
		}

		assetDisters, err := dister.AssetDisterCreators(allAssets[assetapi.Dister]...)
		if err != nil {
			return err
		}
		cliDisterFactory, err = dister.NewDisterFactory(assetDisters...)
		if err != nil {
			return err
		}

		cliDefaultDisterCfg, err = dister.DefaultConfig()
		if err != nil {
			return err
		}

		assetDockerBuilders, err := dockerbuilder.AssetDockerBuilderCreators(allAssets[assetapi.DockerBuilder]...)
		if err != nil {
			return err
		}
		cliDockerBuilderFactory, err = dockerbuilder.NewDockerBuilderFactory(assetDockerBuilders...)
		if err != nil {
			return err
		}

		return nil
	}
}

func distgoProjectParamFromFlags() (distgo.ProjectInfo, distgo.ProjectParam, error) {
	return distgoProjectParamFromVals(projectDirFlagVal, distgoConfigFileFlagVal, godelConfigFileFlagVal, cliDisterFactory, cliDefaultDisterCfg, cliDockerBuilderFactory)
}

func distgoConfigModTime() *time.Time {
	if distgoConfigFileFlagVal == "" {
		return nil
	}
	fi, err := os.Stat(distgoConfigFileFlagVal)
	if err != nil {
		return nil
	}
	modTime := fi.ModTime()
	return &modTime
}

func distgoProjectParamFromVals(projectDir, distgoConfigFile, godelConfigFile string, disterFactory distgo.DisterFactory, defaultDisterCfg distgo.DisterConfig, dockerBuilderFactory distgo.DockerBuilderFactory) (distgo.ProjectInfo, distgo.ProjectParam, error) {
	var distgoCfg distgo.ProjectConfig
	if distgoConfigFile != "" {
		cfg, err := distgo.LoadConfigFromFile(distgoConfigFile)
		if err != nil {
			return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
		}
		distgoCfg = cfg
	}
	if godelConfigFile != "" {
		cfg, err := godellauncher.ReadGodelConfig(path.Dir(godelConfigFile))
		if err != nil {
			return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
		}
		distgoCfg.Exclude.Add(cfg.Exclude)
	}
	projectParam, err := distgoCfg.ToParam(projectDir, disterFactory, defaultDisterCfg, dockerBuilderFactory)
	if err != nil {
		return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
	}
	projectInfo, err := projectParam.ProjectInfo(projectDirFlagVal)
	if err != nil {
		return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
	}
	return projectInfo, projectParam, nil
}
