// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"os"
	"strconv"
	"time"
	"unicode"

	"github.com/palantir/pkg/cli/completion"
)

type DurationFlag struct {
	Name  string
	Alias string
	// Value is string instead of time.Duration so help can format it
	// however makes sense, like 120s instead of 2m0s
	Value      string
	Usage      string
	EnvVar     string
	Required   bool
	Deprecated string
}

func (f DurationFlag) MainName() string {
	return f.Name
}

func (f DurationFlag) FullNames() []string {
	if f.Alias == "" {
		return []string{WithPrefix(f.Name)}
	}
	return []string{WithPrefix(f.Name), WithPrefix(f.Alias)}
}

func (f DurationFlag) IsRequired() bool {
	return f.Required
}

func (f DurationFlag) DeprecationStr() string {
	return f.Deprecated
}

func (f DurationFlag) HasLeader() bool {
	return true
}

func (f DurationFlag) Default() interface{} {
	dur, err := f.Parse(f.DefaultStr())
	if err != nil {
		panic(err)
	}
	return dur
}

func (f DurationFlag) Parse(str string) (interface{}, error) {
	// default to seconds if no suffix
	sec, err := strconv.ParseInt(str, 10, 64)
	if err == nil {
		return time.Duration(sec) * time.Second, nil
	}
	return time.ParseDuration(str)
}

func (f DurationFlag) PlaceholderStr() string {
	return defaultPlaceholder(f.Name)
}

func (f DurationFlag) DefaultStr() string {
	if f.EnvVar == "" {
		return f.Value
	}
	v := os.Getenv(f.EnvVar)
	if v == "" {
		return f.Value
	}
	return v
}

func (f DurationFlag) EnvVarStr() string {
	return f.EnvVar
}

func (f DurationFlag) UsageStr() string {
	return f.Usage + "; value in seconds or with suffix like 500ms, 30s, 15m, 2h"
}

func DurationProvider(ctx *completion.ProviderCtx) []string {
	if ctx.Partial == "" {
		return nil
	}
	if ctx.Partial == "0" {
		return []string{"0"} // units not required
	}

	digits := 0
	for len(ctx.Partial) > digits && unicode.IsDigit(rune(ctx.Partial[digits])) {
		digits++
	}
	if digits == 0 {
		return nil
	}

	num := ctx.Partial[:digits]
	return []string{num + "ms", num + "s", num + "m", num + "h"}
}
