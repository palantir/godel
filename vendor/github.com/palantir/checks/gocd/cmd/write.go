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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/checks/gocd"
)

const importsFileName = "gocd_imports.json"

func DoWriteImportsJSON(dirs []string) error {
	var failedDirs []string
	errs := make(map[string]error)

	for _, dir := range dirs {
		if err := writeImportsJSON(dir); err != nil {
			failedDirs = append(failedDirs, dir)
			errs[dir] = err
		}
	}

	if len(failedDirs) > 0 {
		// if there is only one error, wrap it and return
		if len(failedDirs) == 1 {
			return errors.Wrapf(errs[failedDirs[0]], "failed to write imports for %s", failedDirs[0])
		}
		// otherwise, create compound error
		msgParts := []string{fmt.Sprintf("failed to write imports for %d directories:", len(failedDirs))}
		for _, dir := range failedDirs {
			msgParts = append(msgParts, fmt.Sprintf("\t%s: %s", dir, errs[dir].Error()))
		}
		return errors.New(strings.Join(msgParts, "\n"))
	}

	return nil
}

func writeImportsJSON(rootDir string) error {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}

	report, err := gocd.CreateImportReport(rootDir)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(report, "", "    ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(path.Join(rootDir, importsFileName), bytes, 0644); err != nil {
		return err
	}

	return nil
}
