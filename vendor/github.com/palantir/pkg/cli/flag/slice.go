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

type StringSlice struct {
	Name     string
	Usage    string
	Optional bool
}

func (f StringSlice) MainName() string {
	return f.Name
}

func (f StringSlice) FullNames() []string {
	return []string{f.Name}
}

func (f StringSlice) IsRequired() bool {
	return !f.Optional
}

func (f StringSlice) DeprecationStr() string {
	return ""
}

func (f StringSlice) HasLeader() bool {
	return false
}

func (f StringSlice) Default() interface{} {
	return []string{}
}

func (f StringSlice) Parse(str string) (interface{}, error) {
	return str, nil
}

func (f StringSlice) PlaceholderStr() string {
	return defaultPlaceholder(f.Name)
}

func (f StringSlice) DefaultStr() string {
	return ""
}

func (f StringSlice) EnvVarStr() string {
	return ""
}

func (f StringSlice) UsageStr() string {
	return f.Usage
}
