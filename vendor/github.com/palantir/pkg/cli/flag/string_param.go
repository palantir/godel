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

package flag

type StringParam struct {
	Name  string
	Usage string
}

func (f StringParam) MainName() string {
	return f.Name
}

func (f StringParam) FullNames() []string {
	return []string{f.Name}
}

func (f StringParam) IsRequired() bool {
	return true
}

func (f StringParam) DeprecationStr() string {
	return ""
}

func (f StringParam) HasLeader() bool {
	return false
}

func (f StringParam) Default() interface{} {
	panic("always required")
}

func (f StringParam) Parse(str string) (interface{}, error) {
	return str, nil
}

func (f StringParam) PlaceholderStr() string {
	return defaultPlaceholder(f.Name)
}

func (f StringParam) DefaultStr() string {
	panic("always required")
}

func (f StringParam) EnvVarStr() string {
	panic("always required")
}

func (f StringParam) UsageStr() string {
	return f.Usage
}
