// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type SrcPkg struct {
	MainPkg              string `yaml:"main"`
	DistanceToProjectPkg int    `yaml:"distance-to-project-pkg"`
}

type Config struct {
	Pkgs map[string]SrcPkg `yaml:"packages"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read file %s", configPath)
	}

	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal file %s", configPath)
	}

	if len(cfg.Pkgs) == 0 {
		return nil, errors.Errorf("configuration read from file %s with content %q was empty", configPath, string(file))
	}

	for name, pkg := range cfg.Pkgs {
		if name == "" {
			return nil, errors.Errorf("config cannot contain a blank name: %v", cfg)
		}

		if pkg.MainPkg == "" {
			return nil, errors.Errorf("config for package %s had a blank main package directory: %v", name, cfg)
		}
	}

	return &cfg, nil
}
