// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Context struct {
	App        *App
	Command    *Command
	Path       []string
	IsTerminal func() bool

	context   context.Context
	cancel    context.CancelFunc
	defaults  map[string]interface{}
	specified map[string]interface{}
	// stores all flag values in the order they were encountered
	allVals map[string][]interface{}
}

func (ctx *Context) Context() context.Context {
	return ctx.context
}

func (ctx *Context) Has(name string) bool {
	return ctx.specified[name] != nil
}

func (ctx *Context) Bool(name string) bool {
	return ctx.get(name).(bool)
}

func (ctx *Context) String(name string) string {
	return ctx.get(name).(string)
}

func (ctx *Context) StringSlice(name string) []string {
	s := ctx.getSlice(name)
	out := make([]string, len(s))
	for i, v := range s {
		out[i] = v.(string)
	}
	return out
}

func (ctx *Context) Duration(name string) time.Duration {
	return ctx.get(name).(time.Duration)
}

func (ctx *Context) DurationSlice(name string) []time.Duration {
	s := ctx.getSlice(name)
	out := make([]time.Duration, len(s))
	for i, v := range s {
		out[i] = v.(time.Duration)
	}
	return out
}

// Slice is specifically for flag.StringSlice (whereas the other "*Slice" functions are for retrieving all of the values
// specified for an individual flag).
func (ctx *Context) Slice(name string) []string {
	return ctx.get(name).([]string)
}

func (ctx *Context) Int(name string) int {
	return ctx.get(name).(int)
}

func (ctx *Context) IntSlice(name string) []int {
	s := ctx.getSlice(name)
	out := make([]int, len(s))
	for i, v := range s {
		out[i] = v.(int)
	}
	return out
}

func (ctx *Context) get(name string) interface{} {
	if v, ok := ctx.specified[name]; ok {
		return v
	}
	if v, ok := ctx.defaults[name]; ok {
		return v
	}
	panic(fmt.Errorf("command %q does not have a flag named %q", strings.Join(ctx.Path, " "), name))
}

func (ctx *Context) getSlice(name string) []interface{} {
	if _, ok := ctx.specified[name]; ok {
		// guaranteed that at least one value exists: return slice
		return ctx.allVals[name]
	}
	if v, ok := ctx.defaults[name]; ok {
		// no values exist, but default value is specified: return slice with just default value
		return []interface{}{v}
	}
	panic(fmt.Errorf("command %q does not have a flag named %q", strings.Join(ctx.Path, " "), name))
}
