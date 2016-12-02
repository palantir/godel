// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strings"
	"time"
)

type Context struct {
	App        *App
	Command    *Command
	Path       []string
	IsTerminal func() bool

	defaults  map[string]interface{}
	specified map[string]interface{}
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

func (ctx *Context) Duration(name string) time.Duration {
	return ctx.get(name).(time.Duration)
}

func (ctx *Context) Slice(name string) []string {
	return ctx.get(name).([]string)
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
