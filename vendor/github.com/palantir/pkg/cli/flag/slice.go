// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
