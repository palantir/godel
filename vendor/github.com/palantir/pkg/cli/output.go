// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
