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
)

func (ctx *Context) Printf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ctx.App.Stdout, format, a...)
}

func (ctx *Context) Println(a ...interface{}) {
	_, _ = fmt.Fprintln(ctx.App.Stdout, a...)
}

func (ctx *Context) Print(a ...interface{}) {
	_, _ = fmt.Fprint(ctx.App.Stdout, a...)
}

func (ctx *Context) Errorf(format string, a ...interface{}) {
	_, _ = fmt.Fprintf(ctx.App.Stderr, format, a...)
}

func (ctx *Context) Errorln(a ...interface{}) {
	_, _ = fmt.Fprintln(ctx.App.Stderr, a...)
}
