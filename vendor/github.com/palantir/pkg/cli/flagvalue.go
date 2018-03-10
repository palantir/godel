// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/palantir/pkg/cli/flag"
)

// FlagValue matches the interface specified by the viper package.
type FlagValue interface {
	HasChanged() bool
	Name() string
	ValueString() string
	ValueType() string
}

type flagVal struct {
	name          string
	typ           string
	ctx           *Context
	valStringFunc func() string
}

func (f *flagVal) HasChanged() bool {
	return f.ctx.Has(f.name)
}

func (f *flagVal) Name() string {
	return f.name
}

func (f *flagVal) ValueString() string {
	return f.valStringFunc()
}

func (f *flagVal) ValueType() string {
	return f.typ
}

func (ctx *Context) FlagValue(name string) FlagValue {
	for _, currFlag := range ctx.Command.Flags {
		currFlagName := currFlag.MainName()
		if name != currFlagName {
			continue
		}

		var typ string
		var valStringFunc func() string

		switch currFlag.(type) {
		default:
			typ = "string"
			valStringFunc = func() string {
				stringVal, ok := ctx.specified[currFlagName]
				if !ok {
					return ""
				}
				return stringVal.(string)
			}
		case flag.BoolFlag:
			typ = "bool"
			valStringFunc = func() string {
				boolVal, ok := ctx.specified[currFlagName]
				if !ok {
					return ""
				}
				return strconv.FormatBool(boolVal.(bool))
			}
		case flag.IntFlag:
			typ = "int"
			valStringFunc = func() string {
				intVal, ok := ctx.specified[currFlagName]
				if !ok {
					return ""
				}
				return strconv.FormatInt(int64(intVal.(int)), 10)
			}
		case flag.StringSlice:
			typ = "stringSlice"
			valStringFunc = func() string {
				strSliceVal, ok := ctx.specified[currFlagName]
				if !ok {
					return ""
				}
				return fmt.Sprintf("[%s]", strings.Join(strSliceVal.([]string), ","))
			}
		case flag.DurationFlag:
			typ = "string"
			valStringFunc = func() string {
				durationVal, ok := ctx.specified[currFlagName]
				if !ok {
					return ""
				}
				return fmt.Sprint(durationVal.(time.Duration))
			}
		}
		return &flagVal{
			name:          currFlagName,
			typ:           typ,
			ctx:           ctx,
			valStringFunc: valStringFunc,
		}
	}
	return nil
}
