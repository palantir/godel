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

package golicense

import (
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"

	"github.com/palantir/pkg/matcher"
	"github.com/pkg/errors"
)

type LicenseParams struct {
	Header        string
	CustomHeaders []CustomLicenseParam
	Exclude       matcher.Matcher
}

func (p *LicenseParams) validate() error {
	var emptyNameParams []CustomLicenseParam
	nameToParams := make(map[string][]CustomLicenseParam)

	for _, v := range p.CustomHeaders {
		if v.Name == "" {
			emptyNameParams = append(emptyNameParams, v)
		}
		nameToParams[v.Name] = append(nameToParams[v.Name], v)
	}

	if len(emptyNameParams) > 0 {
		return errors.Errorf("custom header entries have blank names: %+v", emptyNameParams)
	}

	var nameCollisionMsgs []string
	for k, v := range nameToParams {
		if len(v) > 1 {
			nameCollisionMsgs = append(nameCollisionMsgs, fmt.Sprintf("\t%s: %+v", k, v))
		}
	}
	if len(nameCollisionMsgs) > 0 {
		return errors.Errorf(strings.Join(append([]string{"multiple custom header entries have the same name:"}, nameCollisionMsgs...), "\n"))
	}
	return nil
}

type CustomLicenseParam struct {
	Name    string
	Header  string
	Include matcher.Matcher
}

func LicenseFiles(files []string, params LicenseParams, modify bool) ([]string, error) {
	return processFiles(files, params, modify, applyLicenseToFiles)
}

func UnlicenseFiles(files []string, params LicenseParams, modify bool) ([]string, error) {
	return processFiles(files, params, modify, removeLicenseFromFiles)
}

func processFiles(files []string, params LicenseParams, modify bool, f func(files []string, header string, modify bool) ([]string, error)) ([]string, error) {
	if err := params.validate(); err != nil {
		return nil, errors.Wrapf(err, "license parameters invalid")
	}

	goFileMatcher := matcher.Name(`.*\.go`)
	var goFiles []string
	for _, f := range files {
		if goFileMatcher.Match(f) && (params.Exclude == nil || !params.Exclude.Match(f)) {
			goFiles = append(goFiles, f)
		}
	}

	m := make(map[string][]string)
	for _, v := range params.CustomHeaders {
		for _, f := range goFiles {
			if v.Include != nil && v.Include.Match(f) {
				m[v.Name] = append(m[v.Name], f)
			}
		}
	}
	if err := keysWithCommonValuesError(m); err != nil {
		return nil, err
	}

	// all files that were processed (considered by a matcher)
	processedFiles := make(map[string]struct{})
	// all files that were modified (or would have been modified)
	var modified []string

	// process custom matchers
	for _, v := range params.CustomHeaders {
		currModified, err := f(m[v.Name], v.Header, modify)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to process headers for matcher %s", v.Name)
		}
		modified = append(modified, currModified...)
		for _, f := range m[v.Name] {
			processedFiles[f] = struct{}{}
		}
	}

	// process all "*.go" files not matched by custom matchers
	var unprocessedGoFiles []string
	for _, f := range goFiles {
		if _, ok := processedFiles[f]; !ok {
			unprocessedGoFiles = append(unprocessedGoFiles, f)
		}
	}
	currModified, err := f(unprocessedGoFiles, params.Header, modify)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to process headers for default *.go matcher")
	}
	modified = append(modified, currModified...)
	for _, f := range currModified {
		processedFiles[f] = struct{}{}
	}

	sort.Strings(modified)
	return modified, nil
}

func keysWithCommonValuesError(in map[string][]string) error {
	// create map from k -> set of values
	m := make(map[string]map[string]struct{}, len(in))
	for k, v := range in {
		m[k] = make(map[string]struct{}, len(v))
		for _, vv := range v {
			m[k][vv] = struct{}{}
		}
	}

	sortedKeys := make([]string, 0, len(in))
	for k := range in {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	var commonFileMessages []string
	for i := range sortedKeys {
		for j := i + 1; j < len(sortedKeys); j++ {
			common := intersection(m[sortedKeys[i]], m[sortedKeys[j]])
			if len(common) > 0 {
				commonFileMessages = append(commonFileMessages, fmt.Sprintf("%s and %s both match files: %v", sortedKeys[i], sortedKeys[j], common))
			}
		}
	}

	if len(commonFileMessages) > 0 {
		return errors.Errorf(strings.Join(append([]string{"overlap exists between custom matchers"}, commonFileMessages...), "\n"))
	}
	return nil
}

func intersection(a, b map[string]struct{}) []string {
	smallerMap := a
	largerMap := b
	if len(b) < len(a) {
		smallerMap = b
		largerMap = a
	}
	var intersection []string
	for k := range smallerMap {
		if _, ok := largerMap[k]; ok {
			intersection = append(intersection, k)
		}
	}
	sort.Strings(intersection)
	return intersection
}

func applyLicenseToFiles(files []string, header string, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if !strings.HasPrefix(content, header+"\n") {
			if modify {
				content = header + "\n" + content
				if err := ioutil.WriteFile(path, []byte(content), fi.Mode()); err != nil {
					return false, errors.Wrapf(err, "failed to write file %s with new license", path)
				}
			}
			return true, nil
		}
		return false, nil
	})
}

func removeLicenseFromFiles(files []string, header string, modify bool) ([]string, error) {
	return visitFiles(files, func(path string, fi os.FileInfo, content string) (bool, error) {
		if strings.HasPrefix(content, header+"\n") {
			if modify {
				content = strings.TrimPrefix(content, header+"\n")
				if err := ioutil.WriteFile(path, []byte(content), fi.Mode()); err != nil {
					return false, errors.Wrapf(err, "failed to write file %s with license removed", path)
				}
			}
			return true, nil
		}
		return false, nil
	})
}

func visitFiles(files []string, visitor func(path string, fi os.FileInfo, content string) (bool, error)) ([]string, error) {
	var modified []string

	for _, f := range files {
		fi, err := os.Stat(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to stat %s", f)
		}
		bytes, err := ioutil.ReadFile(f)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to read %s", f)
		}
		content := string(bytes)
		if changed, err := visitor(f, fi, content); err != nil {
			return nil, errors.WithStack(err)
		} else if changed {
			modified = append(modified, f)
		}
	}

	return modified, nil
}
