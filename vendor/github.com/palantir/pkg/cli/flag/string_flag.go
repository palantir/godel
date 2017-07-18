// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"os"
)

type StringFlag struct {
	Name        string
	Alias       string
	Value       string
	Usage       string
	Placeholder string
	EnvVar      string
	Required    bool
	Deprecated  string
}

func (f StringFlag) MainName() string {
	return f.Name
}

func (f StringFlag) FullNames() []string {
	if f.Alias == "" {
		return []string{WithPrefix(f.Name)}
	}
	return []string{WithPrefix(f.Name), WithPrefix(f.Alias)}
}

func (f StringFlag) IsRequired() bool {
	return f.Required
}

func (f StringFlag) DeprecationStr() string {
	return f.Deprecated
}

func (f StringFlag) HasLeader() bool {
	return true
}

func (f StringFlag) Default() interface{} {
	if f.EnvVar == "" {
		return f.Value
	}
	v := os.Getenv(f.EnvVar)
	if v == "" {
		return f.Value
	}
	return v
}

func (f StringFlag) Parse(str string) (interface{}, error) {
	return str, nil
}

func (f StringFlag) PlaceholderStr() string {
	return placeholderOrDefault(f.Placeholder, f.Name)
}

func (f StringFlag) DefaultStr() string {
	return f.Value
}

func (f StringFlag) EnvVarStr() string {
	return f.EnvVar
}

func (f StringFlag) UsageStr() string {
	return f.Usage
}
