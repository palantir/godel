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

package packages

import (
	"github.com/palantir/pkg/matcher"
	"github.com/palantir/pkg/pkgpath"
	"github.com/pkg/errors"
)

func List(exclude matcher.Matcher, wd string) ([]string, error) {
	pkgs, err := pkgpath.PackagesInDir(wd, exclude)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to list packages")
	}

	pkgPaths, err := pkgs.Paths(pkgpath.Relative)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get package paths")
	}

	return pkgPaths, nil
}
