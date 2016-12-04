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

package products

import (
	"fmt"
	"io"

	"github.com/palantir/godel/apps/distgo/cmd/build"
	"github.com/palantir/godel/apps/distgo/params"
)

func PrintProducts(cfg params.Project, wd string, stdout io.Writer) error {
	return build.RunBuildFunc(func(buildSpec []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		for _, spec := range buildSpec {
			fmt.Fprintln(stdout, spec.Spec.ProductName)
		}
		return nil
	}, cfg, nil, wd, stdout)
}
