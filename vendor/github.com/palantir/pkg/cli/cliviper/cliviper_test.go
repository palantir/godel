// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cliviper_test

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/cliviper"
	"github.com/palantir/pkg/cli/flag"
)

func TestCLIViperApp(t *testing.T) {
	var resultVal string

	const msgFlag = "message"
	const content = "messageContent"

	// set cliviper.App() option
	app := cli.NewApp(cliviper.App())
	app.Flags = []flag.Flag{
		flag.StringFlag{Name: msgFlag},
	}
	app.Action = func(ctx cli.Context) error {
		// all flags for a context should be bound by viper
		resultVal = viper.GetString(msgFlag)
		return nil
	}
	app.Run([]string{"app", "--message", content})

	assert.Equal(t, content, resultVal)
}
