// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package clitest

import (
	"github.com/palantir/pkg/cli/flag"
)

type dummyFlag struct {
	Name  string
	Value interface{}
}

func (f dummyFlag) MainName() string {
	return f.Name
}

func (f dummyFlag) FullNames() []string {
	return []string{flag.WithPrefix(f.Name)}
}

func (f dummyFlag) IsRequired() bool {
	return false
}

func (f dummyFlag) DeprecationStr() string {
	return ""
}

func (f dummyFlag) HasLeader() bool {
	return true
}

func (f dummyFlag) Default() interface{} {
	return f.Value
}

func (f dummyFlag) Parse(string) (interface{}, error) {
	return f.Value, nil
}

func (f dummyFlag) PlaceholderStr() string {
	panic("not implemented")
}

func (f dummyFlag) DefaultStr() string {
	panic("not implemented")
}

func (f dummyFlag) EnvVarStr() string {
	panic("not implemented")
}

func (f dummyFlag) UsageStr() string {
	panic("not implemented")
}
