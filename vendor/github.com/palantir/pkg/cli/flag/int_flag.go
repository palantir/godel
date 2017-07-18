// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"os"
	"strconv"
)

type IntFlag struct {
	Name       string
	Alias      string
	Value      int
	Usage      string
	EnvVar     string
	Required   bool
	Deprecated string
}

func (f IntFlag) MainName() string {
	return f.Name
}

func (f IntFlag) FullNames() []string {
	if f.Alias == "" {
		return []string{WithPrefix(f.Name)}
	}
	return []string{WithPrefix(f.Name), WithPrefix(f.Alias)}
}

func (f IntFlag) IsRequired() bool {
	return f.Required
}

func (f IntFlag) DeprecationStr() string {
	return f.Deprecated
}

func (f IntFlag) HasLeader() bool {
	return true
}

func (f IntFlag) Default() interface{} {
	// if environment variable is not defined, return value
	if f.EnvVar == "" {
		return f.Value
	}
	v := os.Getenv(f.EnvVar)
	if v == "" {
		return f.Value
	}
	i, err := f.Parse(v)
	if err != nil {
		// if environment variable is defined but cannot be parsed as an int, panic
		panic(err)
	}
	// return value parsed from environment variable
	return i
}

func (f IntFlag) Parse(str string) (interface{}, error) {
	i, err := strconv.ParseInt(str, 10, 0)
	if err != nil {
		return nil, err
	}
	return int(i), nil
}

func (f IntFlag) PlaceholderStr() string {
	return defaultPlaceholder(f.Name)
}

func (f IntFlag) DefaultStr() string {
	return strconv.Itoa(f.Value)
}

func (f IntFlag) EnvVarStr() string {
	return f.EnvVar
}

func (f IntFlag) UsageStr() string {
	return f.Usage + "; value must be convertable to an integer."
}
