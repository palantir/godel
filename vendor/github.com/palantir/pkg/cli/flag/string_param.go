// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
