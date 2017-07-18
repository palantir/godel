// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cfgcli

import (
	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

const (
	configFlag = "config"
	jsonFlag   = "json"
)

var (
	// ConfigPath is a global variable used to store the value of configFlag parsed from the flags for the cli.App
	// returned by NewApp.
	ConfigPath string
	// ConfigJSON is a global variable used to store the value of jsonFlag parsed from the flags for the cli.App
	// returned by NewApp.
	ConfigJSON string
)

// NewApp returns a new cli.App configured using Handler.
func NewApp() *cli.App {
	return cli.NewApp(Handler())
}

// Handler returns a cli.Option that configures a cli.App as a config CLI application that is configured with flags for
// a configuration file and configuration JSON. The application is configured with a "Before" hook that parses the
// values of the flag and stores it in the "ConfigPath" and "ConfigJSON" variables (and executes and "Before" hook that
// was previously defined).
func Handler() cli.Option {
	return func(app *cli.App) {
		// store app.Before previously set on App
		before := app.Before
		// add a Before hook that sets value of shared global variables based on global flags
		app.Before = func(ctx cli.Context) error {
			ConfigPath = ctx.String(configFlag)
			ConfigJSON = ctx.String(jsonFlag)

			// if app.Before was previously defined, use it
			if before != nil {
				return before(ctx)
			}
			return nil
		}
		app.Flags = append(app.Flags,
			flag.StringFlag{
				Name:  configFlag,
				Usage: "Path to configuration file",
			},
			flag.StringFlag{
				Name:  jsonFlag,
				Usage: "JSON configuration (provide as literal JSON)",
			},
		)
	}
}
