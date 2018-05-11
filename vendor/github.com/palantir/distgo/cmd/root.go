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
	"io/ioutil"
	"os"
	"time"

	godelconfig "github.com/palantir/godel/framework/godel/config"
	"github.com/palantir/godel/framework/pluginapi"
	"github.com/palantir/pkg/cobracli"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"

	"github.com/palantir/distgo/assetapi"
	"github.com/palantir/distgo/dister"
	"github.com/palantir/distgo/dister/disterfactory"
	"github.com/palantir/distgo/distgo"
	"github.com/palantir/distgo/distgo/config"
	"github.com/palantir/distgo/dockerbuilder"
	"github.com/palantir/distgo/dockerbuilder/dockerbuilderfactory"
	"github.com/palantir/distgo/projectversioner/projectversionerfactory"
	"github.com/palantir/distgo/publisher"
	"github.com/palantir/distgo/publisher/publisherfactory"
)

var (
	debugFlagVal            bool
	projectDirFlagVal       string
	distgoConfigFileFlagVal string
	godelConfigFileFlagVal  string
	assetsFlagVal           []string

	cliProjectVersionerFactory distgo.ProjectVersionerFactory
	cliDisterFactory           distgo.DisterFactory
	cliDefaultDisterCfg        config.DisterConfig
	cliDockerBuilderFactory    distgo.DockerBuilderFactory
	cliPublisherFactory        distgo.PublisherFactory
)

var rootCmd = &cobra.Command{
	Use: "distgo",
}

func Execute() int {
	return cobracli.ExecuteWithDebugVarAndDefaultParams(rootCmd, &debugFlagVal)
}

func restoreRootFlagsFn() func() {
	origProjectDirFlagVal := projectDirFlagVal
	origDistgoConfigFileFlagVal := distgoConfigFileFlagVal
	origGodelConfigFileFlagVal := godelConfigFileFlagVal
	origAssetsFlagVal := assetsFlagVal
	return func() {
		projectDirFlagVal = origProjectDirFlagVal
		distgoConfigFileFlagVal = origDistgoConfigFileFlagVal
		godelConfigFileFlagVal = origGodelConfigFileFlagVal
		assetsFlagVal = origAssetsFlagVal
	}
}

func InitAssetCmds(args []string) error {
	restoreFn := restoreRootFlagsFn()
	// parse the flags to retrieve the value of the "--assets" flag. Ignore any errors that occur in flag parsing so
	// that, if provided flags are invalid, the regular logic handles the error printing.
	_ = rootCmd.ParseFlags(args)
	allAssets, err := assetapi.LoadAssets(assetsFlagVal)
	// restore the root flags to undo any parsing done by rootCmd.ParseFlags
	restoreFn()
	if err != nil {
		return err
	}

	// load publisher assets
	assetPublishers, upgraderPublishers, err := publisher.AssetPublisherCreators(allAssets[assetapi.Publisher]...)
	if err != nil {
		return err
	}

	cliPublisherFactory, err = publisherfactory.New(assetPublishers, upgraderPublishers)
	if err != nil {
		return err
	}

	publisherTypeNames := cliPublisherFactory.Types()
	var publishers []distgo.Publisher
	for _, typeName := range publisherTypeNames {
		publisher, err := cliPublisherFactory.NewPublisher(typeName)
		if err != nil {
			return errors.Wrapf(err, "failed to create publisher %q", typeName)
		}
		publishers = append(publishers, publisher)
	}

	// add publish commands based on assets
	addPublishSubcommands(publisherTypeNames, publishers)

	return nil
}

func init() {
	pluginapi.AddDebugPFlagPtr(rootCmd.PersistentFlags(), &debugFlagVal)
	pluginapi.AddProjectDirPFlagPtr(rootCmd.PersistentFlags(), &projectDirFlagVal)
	pluginapi.AddConfigPFlagPtr(rootCmd.PersistentFlags(), &distgoConfigFileFlagVal)
	pluginapi.AddGodelConfigPFlagPtr(rootCmd.PersistentFlags(), &godelConfigFileFlagVal)
	pluginapi.AddAssetsPFlagPtr(rootCmd.PersistentFlags(), &assetsFlagVal)

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		allAssets, err := assetapi.LoadAssets(assetsFlagVal)
		if err != nil {
			return err
		}

		// parameters will become non-nil if/when support for projectversioner assets are added
		cliProjectVersionerFactory, err = projectversionerfactory.New(nil, nil)
		if err != nil {
			return err
		}

		assetDisters, upgraderDisters, err := dister.AssetDisterCreators(allAssets[assetapi.Dister]...)
		if err != nil {
			return err
		}
		cliDisterFactory, err = disterfactory.New(assetDisters, upgraderDisters)
		if err != nil {
			return err
		}

		cliDefaultDisterCfg, err = disterfactory.DefaultConfig()
		if err != nil {
			return err
		}

		assetDockerBuilders, upgraderDockerBuilders, err := dockerbuilder.AssetDockerBuilderCreators(allAssets[assetapi.DockerBuilder]...)
		if err != nil {
			return err
		}
		cliDockerBuilderFactory, err = dockerbuilderfactory.New(assetDockerBuilders, upgraderDockerBuilders)
		if err != nil {
			return err
		}

		return nil
	}
}

func distgoProjectParamFromFlags() (distgo.ProjectInfo, distgo.ProjectParam, error) {
	return distgoProjectParamFromVals(projectDirFlagVal, distgoConfigFileFlagVal, godelConfigFileFlagVal, cliProjectVersionerFactory, cliDisterFactory, cliDefaultDisterCfg, cliDockerBuilderFactory, cliPublisherFactory)
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

func distgoProjectParamFromVals(projectDir, distgoConfigFile, godelConfigFile string, projectVersionerFactory distgo.ProjectVersionerFactory, disterFactory distgo.DisterFactory, defaultDisterCfg config.DisterConfig, dockerBuilderFactory distgo.DockerBuilderFactory, publisherFactory distgo.PublisherFactory) (distgo.ProjectInfo, distgo.ProjectParam, error) {
	var distgoCfg config.ProjectConfig
	if distgoConfigFile != "" {
		cfg, err := loadConfigFromFile(distgoConfigFile)
		if err != nil {
			return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
		}
		distgoCfg = cfg
	}
	if godelConfigFile != "" {
		cfg, err := godelconfig.ReadGodelConfigFromFile(godelConfigFile)
		if err != nil {
			return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
		}
		distgoCfg.Exclude.Add(cfg.Exclude)
	}
	projectParam, err := distgoCfg.ToParam(projectDir, projectVersionerFactory, disterFactory, defaultDisterCfg, dockerBuilderFactory, publisherFactory)
	if err != nil {
		return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
	}
	projectInfo, err := projectParam.ProjectInfo(projectDirFlagVal)
	if err != nil {
		return distgo.ProjectInfo{}, distgo.ProjectParam{}, err
	}
	return projectInfo, projectParam, nil
}

func loadConfigFromFile(cfgFile string) (config.ProjectConfig, error) {
	cfgBytes, err := ioutil.ReadFile(cfgFile)
	if os.IsNotExist(err) {
		return config.ProjectConfig{}, nil
	}
	if err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to read configuration file")
	}
	upgradedCfgBytes, err := config.UpgradeConfig(cfgBytes, cliProjectVersionerFactory, cliDisterFactory, cliDockerBuilderFactory, cliPublisherFactory)
	if err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to upgrade configuration")
	}

	var cfg config.ProjectConfig
	if err := yaml.Unmarshal(upgradedCfgBytes, &cfg); err != nil {
		return config.ProjectConfig{}, errors.Wrapf(err, "failed to unmarshal configuration")
	}
	return cfg, nil
}
