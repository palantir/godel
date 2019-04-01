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

package godellauncher

import (
	"os"
	"path"
	"path/filepath"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

// ListProjectPaths lists all of the paths in the provided project directory that matches the provided include matcher
// and does not match the provided exclude matcher. The paths are relative to the current working directory.
func ListProjectPaths(projectDir string, include, exclude matcher.Matcher) ([]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine working directory")
	}
	if !filepath.IsAbs(projectDir) {
		projectDir = path.Join(wd, projectDir)
	}
	relPathPrefix, err := filepath.Rel(wd, projectDir)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to determine relative path")
	}

	files, err := matcher.ListFiles(projectDir, include, exclude)
	if err != nil {
		return nil, err
	}
	if relPathPrefix != "" {
		for i, file := range files {
			files[i] = path.Join(relPathPrefix, file)
		}
	}
	return files, nil
}
