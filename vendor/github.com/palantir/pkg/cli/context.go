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
