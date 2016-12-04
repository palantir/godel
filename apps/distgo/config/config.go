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

package config

import (
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"github.com/palantir/godel/apps/distgo/params"
	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

type RawProjectConfig struct {
	Products          map[string]RawProductConfig `yaml:"products" json:"products"`
	BuildOutputDir    string                      `yaml:"build-output-dir" json:"build-output-dir"`
	DistOutputDir     string                      `yaml:"dist-output-dir" json:"dist-output-dir"`
	GroupID           string                      `yaml:"group-id" json:"group-id"`
	DistScriptInclude string                      `yaml:"dist-script-include" json:"dist-script-include"`
	Exclude           matcher.NamesPathsCfg       `yaml:"exclude" json:"exclude"`
}

func (cfg *RawProjectConfig) ToParams() (params.Project, error) {
	products := make(map[string]params.Product, len(cfg.Products))
	for k, v := range cfg.Products {
		productParam, err := v.ToParam()
		if err != nil {
			return params.Project{}, err
		}
		products[k] = productParam
	}
	return params.Project{
		Products:          products,
		BuildOutputDir:    cfg.BuildOutputDir,
		DistOutputDir:     cfg.DistOutputDir,
		DistScriptInclude: cfg.DistScriptInclude,
		GroupID:           cfg.GroupID,
		Exclude:           cfg.Exclude.Matcher(),
	}, nil
}

// RawProductConfig represents user-specified configuration on how to build a specific product.
type RawProductConfig struct {
	Build          RawBuildConfig   `yaml:"build" json:"build"`
	Run            RawRunConfig     `yaml:"run" json:"run"`
	Dist           RawDistConfigs   `yaml:"dist" json:"dist"`
	DefaultPublish RawPublishConfig `yaml:"publish" json:"publish"`
}

func (cfg *RawProductConfig) ToParam() (params.Product, error) {
	var dists []params.Dist
	for _, rawDistCfg := range cfg.Dist {
		dist, err := rawDistCfg.ToParam()
		if err != nil {
			return params.Product{}, err
		}
		dists = append(dists, dist)
	}

	return params.Product{
		Build:          cfg.Build.ToParam(),
		Run:            cfg.Run.ToParam(),
		Dist:           dists,
		DefaultPublish: cfg.DefaultPublish.ToParams(),
	}, nil
}

type RawBuildConfig struct {
	Script          string            `yaml:"script" json:"script"`
	MainPkg         string            `yaml:"main-pkg" json:"main-pkg"`
	OutputDir       string            `yaml:"output-dir" json:"output-dir"`
	BuildArgsScript string            `yaml:"build-args-script" json:"build-args-script"`
	VersionVar      string            `yaml:"version-var" json:"version-var"`
	Environment     map[string]string `yaml:"environment" json:"environment"`
	OSArchs         []osarch.OSArch   `yaml:"os-archs" json:"os-archs"`
}

func (cfg *RawBuildConfig) ToParam() params.Build {
	return params.Build{
		Script:          cfg.Script,
		MainPkg:         cfg.MainPkg,
		OutputDir:       cfg.OutputDir,
		BuildArgsScript: cfg.BuildArgsScript,
		VersionVar:      cfg.VersionVar,
		Environment:     cfg.Environment,
		OSArchs:         cfg.OSArchs,
	}
}

type RawDistConfigs []RawDistConfig

func (out *RawDistConfigs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multiple []RawDistConfig
	if err := unmarshal(&multiple); err == nil {
		if len(multiple) == 0 {
			return errors.New("if `dist` key is specified, there must be at least one dist")
		}
		*out = multiple
		return nil
	}

	var single RawDistConfig
	if err := unmarshal(&single); err != nil {
		// return the error from a single DistConfig if neither one works
		return err
	}
	*out = []RawDistConfig{single}
	return nil
}

type RawDistConfig struct {
	OutputDir     string            `yaml:"output-dir" json:"output-dir"`
	InputDir      string            `yaml:"input-dir" json:"input-dir"`
	InputProducts []string          `yaml:"input-products" json:"input-products"`
	Script        string            `yaml:"script" json:"script"`
	DistType      RawDistInfoConfig `yaml:"dist-type" json:"dist-type"`
	Publish       RawPublishConfig  `yaml:"publish" json:"publish"`
}

func (cfg *RawDistConfig) ToParam() (params.Dist, error) {
	info, err := cfg.DistType.ToParam()
	if err != nil {
		return params.Dist{}, err
	}
	return params.Dist{
		OutputDir:     cfg.OutputDir,
		InputDir:      cfg.InputDir,
		InputProducts: cfg.InputProducts,
		Script:        cfg.Script,
		Info:          info,
		Publish:       cfg.Publish.ToParams(),
	}, nil
}

type RawRunConfig struct {
	Args []string `yaml:"args" json:"args"`
}

func (cfg *RawRunConfig) ToParam() params.Run {
	return params.Run{
		Args: cfg.Args,
	}
}

type RawDistInfoConfig struct {
	Type string      `yaml:"type" json:"type"`
	Info interface{} `yaml:"info" json:"info"`
}

func (cfg *RawDistInfoConfig) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// unmarshal type alias (uses deafult unmarshal strategy)
	type rawDistInfoConfigAlias RawDistInfoConfig
	var rawAliasConfig rawDistInfoConfigAlias
	if err := unmarshal(&rawAliasConfig); err != nil {
		return err
	}

	rawDistInfoConfig := RawDistInfoConfig(rawAliasConfig)
	switch params.DistInfoType(rawDistInfoConfig.Type) {
	case params.SLSDistType:
		type typedRawConfig struct {
			Type string
			Info RawSLSDistConfig
		}
		var rawSLS typedRawConfig
		if err := unmarshal(&rawSLS); err != nil {
			return err
		}
		rawDistInfoConfig.Info = rawSLS.Info
	case params.BinDistType:
		type typedRawConfig struct {
			Type string
			Info RawBinDistConfig
		}
		var rawBin typedRawConfig
		if err := unmarshal(&rawBin); err != nil {
			return err
		}
		rawDistInfoConfig.Info = rawBin.Info
	case params.RPMDistType:
		type typedRawConfig struct {
			Type string
			Info RawRPMDistConfig
		}
		var rawRPM typedRawConfig
		if err := unmarshal(&rawRPM); err != nil {
			return err
		}
		rawDistInfoConfig.Info = rawRPM.Info
	}
	*cfg = rawDistInfoConfig
	return nil
}

func (cfg *RawDistInfoConfig) ToParam() (params.DistInfo, error) {
	var distInfo params.DistInfo
	if cfg.Info != nil {
		convertMapKeysToCamelCase(cfg.Info)
		var decodeErr error
		switch params.DistInfoType(cfg.Type) {
		case params.SLSDistType:
			val := RawSLSDistConfig{}
			decodeErr = mapstructure.Decode(cfg.Info, &val)
			distInfo = &params.SLSDistInfo{
				InitShTemplateFile:   val.InitShTemplateFile,
				ManifestTemplateFile: val.ManifestTemplateFile,
				ServiceArgs:          val.ServiceArgs,
				ProductType:          val.ProductType,
				ManifestExtensions:   val.ManifestExtensions,
				YMLValidationExclude: val.YMLValidationExclude.Matcher(),
			}
		case params.BinDistType:
			val := RawBinDistConfig{}
			decodeErr = mapstructure.Decode(cfg.Info, &val)
			distInfo = &params.BinDistInfo{
				OmitInitSh:         val.OmitInitSh,
				InitShTemplateFile: val.InitShTemplateFile,
			}
		case params.RPMDistType:
			val := RawRPMDistConfig{}
			decodeErr = mapstructure.Decode(cfg.Info, &val)
			distInfo = &params.RPMDistInfo{
				Release:             val.Release,
				ConfigFiles:         val.ConfigFiles,
				BeforeInstallScript: val.BeforeInstallScript,
				AfterInstallScript:  val.AfterInstallScript,
				AfterRemoveScript:   val.AfterRemoveScript,
			}
		default:
			return nil, errors.Errorf("No unmarshaller found for type %s for %v", cfg.Type, *cfg)
		}
		if decodeErr != nil {
			return nil, errors.Wrapf(decodeErr, "failed to unmarshal DistTypeCfg.Info for %v", *cfg)
		}
	}
	return distInfo, nil
}

func convertMapKeysToCamelCase(input interface{}) {
	if inputMap, ok := input.(map[interface{}]interface{}); ok {
		for k, v := range inputMap {
			if str, ok := k.(string); ok {
				newStr := ""
				for _, currPart := range strings.Split(str, "-") {
					newStr += strings.ToUpper(currPart[0:1]) + currPart[1:]
				}
				delete(inputMap, k)
				inputMap[newStr] = v
			}
		}
	}
}

type RawBinDistConfig struct {
	OmitInitSh         bool   `yaml:"omit-init-sh" json:"omit-init-sh"`
	InitShTemplateFile string `yaml:"init-sh-template-file" json:"init-sh-template-file"`
}

func (cfg *RawBinDistConfig) ToParams() params.BinDistInfo {
	return params.BinDistInfo{
		OmitInitSh:         cfg.OmitInitSh,
		InitShTemplateFile: cfg.InitShTemplateFile,
	}
}

type RawSLSDistConfig struct {
	InitShTemplateFile   string                 `yaml:"init-sh-template-file" json:"init-sh-template-file"`
	ManifestTemplateFile string                 `yaml:"manifest-template-file" json:"manifest-template-file"`
	ServiceArgs          string                 `yaml:"service-args" json:"service-args"`
	ProductType          string                 `yaml:"product-type" json:"product-type"`
	ManifestExtensions   map[string]interface{} `yaml:"manifest-extensions" json:"manifest-extensions"`
	YMLValidationExclude matcher.NamesPathsCfg  `yaml:"yml-validation-exclude" json:"yml-validation-exclude"`
}

func (cfg *RawSLSDistConfig) ToParams() params.SLSDistInfo {
	return params.SLSDistInfo{
		InitShTemplateFile:   cfg.InitShTemplateFile,
		ManifestTemplateFile: cfg.ManifestTemplateFile,
		ServiceArgs:          cfg.ServiceArgs,
		ProductType:          cfg.ProductType,
		ManifestExtensions:   cfg.ManifestExtensions,
		YMLValidationExclude: cfg.YMLValidationExclude.Matcher(),
	}
}

type RawRPMDistConfig struct {
	Release             string   `yaml:"release" json:"release"`
	ConfigFiles         []string `yaml:"config-files" json:"config-files"`
	BeforeInstallScript string   `yaml:"before-install-script" json:"before-install-script"`
	AfterInstallScript  string   `yaml:"after-install-script" json:"after-install-script"`
	AfterRemoveScript   string   `yaml:"after-remove-script" json:"after-remove-script"`
}

func (cfg *RawRPMDistConfig) ToParams() params.RPMDistInfo {
	return params.RPMDistInfo{
		Release:             cfg.Release,
		ConfigFiles:         cfg.ConfigFiles,
		BeforeInstallScript: cfg.BeforeInstallScript,
		AfterInstallScript:  cfg.AfterInstallScript,
		AfterRemoveScript:   cfg.AfterRemoveScript,
	}
}

type RawPublishConfig struct {
	GroupID string           `yaml:"group-id" json:"group-id"`
	Almanac RawAlmanacConfig `yaml:"almanac" json:"almanac"`
}

func (cfg *RawPublishConfig) ToParams() params.Publish {
	return params.Publish{
		GroupID: cfg.GroupID,
		Almanac: cfg.Almanac.ToParams(),
	}
}

type RawAlmanacConfig struct {
	Metadata map[string]string `yaml:"metadata" json:"metadata"`
	Tags     []string          `yaml:"tags" json:"tags"`
}

func (cfg *RawAlmanacConfig) ToParams() params.Almanac {
	return params.Almanac{
		Metadata: cfg.Metadata,
		Tags:     cfg.Tags,
	}
}

func Load(cfgPath, jsonContent string) (params.Project, error) {
	var cfgYML string
	if cfgPath != "" {
		file, err := ioutil.ReadFile(cfgPath)
		if err != nil {
			return params.Project{}, errors.Wrapf(err, "failed to read file %s", cfgPath)
		}
		cfgYML = string(file)
	}
	cfg, err := Read(cfgYML, jsonContent)
	if err != nil {
		return params.Project{}, err
	}
	return cfg.ToParams()
}

func Read(cfgYML, jsonContent string) (RawProjectConfig, error) {
	cfg := RawProjectConfig{}
	if cfgYML != "" {
		if err := yaml.Unmarshal([]byte(cfgYML), &cfg); err != nil {
			return RawProjectConfig{}, errors.Wrapf(err, "failed to unmarshal yml %s", cfgYML)
		}
	}
	if jsonContent != "" {
		jsonCfg := RawProjectConfig{}
		if err := json.Unmarshal([]byte(jsonContent), &jsonCfg); err != nil {
			return RawProjectConfig{}, err
		}
		cfg.Exclude.Add(jsonCfg.Exclude)
	}
	return cfg, nil
}
