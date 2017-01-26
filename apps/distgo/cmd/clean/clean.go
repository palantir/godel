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

package clean

import (
	"fmt"
	"os"

	"github.com/pkg/errors"

	"github.com/palantir/godel/apps/distgo/params"
)

const (
	defaultBuildOutputDir = "build"
	defaultDistOutputDir  = "dist"
)

func Clean(products []string, cfg params.Project) error {
	// Only remove default output directories if no products specified
	if len(products) == 0 {
		if cfg.BuildOutputDir == "" {
			cfg.BuildOutputDir = defaultBuildOutputDir
		}
		if err := remove(cfg.BuildOutputDir); err != nil {
			return errors.WithStack(err)
		}

		if cfg.DistOutputDir == "" {
			cfg.DistOutputDir = defaultDistOutputDir
		}
		if err := remove(cfg.DistOutputDir); err != nil {
			return errors.WithStack(err)
		}
	}

	err := cleanProducts(products, cfg)
	if err != nil {
		return errors.Wrap(err, "failed to clean products")
	}

	return nil
}

func cleanProducts(products []string, cfg params.Project) error {
	// remove all products
	if len(products) == 0 {
		for name := range cfg.Products {
			products = append(products, name)
		}
	}

	for _, product := range products {
		p := cfg.Products[product]
		if err := remove(p.Build.OutputDir); err != nil {
			return errors.WithStack(err)
		}
		for _, dist := range p.Dist {
			if err := remove(dist.OutputDir); err != nil {
				return errors.WithStack(err)
			}
		}
	}

	return nil
}

func remove(dir string) error {
	if dir != "" {
		fmt.Printf("Removing %s\n", dir)
		if err := os.RemoveAll(dir); err != nil {
			return errors.Wrapf(err, "failed to remove dir %s", dir)
		}
	}
	return nil
}
