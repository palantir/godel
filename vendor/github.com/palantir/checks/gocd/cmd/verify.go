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
	"os"
	"path"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/pkg/errors"

	"github.com/palantir/checks/gocd/gocd"
)

func DoVerify(dirs []string) error {
	var failedDirs []string
	errs := make(map[string]error)

	for _, dir := range dirs {
		if err := verify(dir); err != nil {
			failedDirs = append(failedDirs, dir)
			errs[dir] = err
		}
	}

	if len(failedDirs) > 0 {
		// create compound error
		msgParts := []string{fmt.Sprintf("%s out of date for %d %s:", importsFileName, len(failedDirs), plural(len(failedDirs)))}
		for _, dir := range failedDirs {
			msgParts = append(msgParts, fmt.Sprintf("\t%s: %s", dir, errs[dir].Error()))
		}
		return errors.New(strings.Join(msgParts, "\n"))
	}

	return nil
}

func plural(n int) string {
	if n == 1 {
		return "directory"
	}
	return "directories"
}

func verify(rootDir string) error {
	rootDir, err := filepath.Abs(rootDir)
	if err != nil {
		return err
	}

	inputFile := path.Join(rootDir, importsFileName)
	if _, err := os.Stat(inputFile); os.IsNotExist(err) {
		return fmt.Errorf("%s does not exist", importsFileName)
	}

	bytes, err := ioutil.ReadFile(inputFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", inputFile)
	}

	gotReport := gocd.ImportReport{}
	if err := json.Unmarshal(bytes, &gotReport); err != nil {
		return errors.Wrapf(err, "failed to unmarshal report")
	}

	wantReport, err := gocd.CreateImportReport(rootDir)
	if err != nil {
		return err
	}

	if !reflect.DeepEqual(wantReport, gotReport) {
		return errors.Errorf("%s is out of date", importsFileName)
	}

	return nil
}
