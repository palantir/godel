// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"os"
	"strconv"
)

type BoolFlag struct {
	Name       string
	Alias      string
	Value      bool
	Usage      string
	EnvVar     string
	Deprecated string
}

func (f BoolFlag) MainName() string {
	return f.Name
}

func (f BoolFlag) FullNames() []string {
	if f.Alias == "" {
		return []string{WithPrefix(f.Name)}
	}
	return []string{WithPrefix(f.Name), WithPrefix(f.Alias)}
}

func (f BoolFlag) IsRequired() bool {
	return false
}

func (f BoolFlag) DeprecationStr() string {
	return f.Deprecated
}

func (f BoolFlag) HasLeader() bool {
	return true
}

func (f BoolFlag) Default() interface{} {
	// if environment variable is not defined, return value
	if f.EnvVar == "" {
		return f.Value
	}
	v := os.Getenv(f.EnvVar)
	if v == "" {
		return f.Value
	}
	b, err := f.Parse(v)
	if err != nil {
		// if environment variable is defined but cannot be parsed as a bool, return false
		return false
	}
	// return value parsed from environment variable
	return b
}

func (f BoolFlag) Parse(str string) (interface{}, error) {
	return strconv.ParseBool(str)
}

func (f BoolFlag) PlaceholderStr() string {
	panic("bool flag does not have placeholder")
}

func (f BoolFlag) DefaultStr() string {
	return ""
}

func (f BoolFlag) EnvVarStr() string {
	return f.EnvVar
}

func (f BoolFlag) UsageStr() string {
	return f.Usage
}
