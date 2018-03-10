// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cliviper

import (
	"github.com/spf13/viper"

	"github.com/palantir/pkg/cli"
)

func FlagValueSet(ctx *cli.Context) viper.FlagValueSet {
	return (*flagValueSet)(ctx)
}

type flagValueSet cli.Context

func (p *flagValueSet) VisitAll(fn func(viper.FlagValue)) {
	for _, currFlag := range p.Command.Flags {
		fn((*cli.Context)(p).FlagValue(currFlag.MainName()))
	}
}

func App() cli.Option {
	return func(app *cli.App) {
		app.ContextOptions = append(app.ContextOptions, func(ctx *cli.Context) {
			_ = viper.BindFlagValues(FlagValueSet(ctx))
		})
	}
}
