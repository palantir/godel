// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/cli/flag"
)

var stringCases = []struct {
	flag flag.Flag
	str  string
}{{
	flag: flag.StringFlag{Name: "path", Alias: "p", Usage: "where it is"},
	str:  "   --path, -p\twhere it is",
}, {
	flag: flag.StringFlag{Name: "path", Usage: "where it is", Required: true},
	str:  " * --path\twhere it is",
}, {
	flag: flag.StringFlag{Name: "path", Value: "/", Usage: "where it is"},
	str:  "   --path\twhere it is (default=\"/\")",
}, {
	flag: flag.StringFlag{Name: "path", EnvVar: "THE_PATH", Usage: "where it is"},
	str:  "   --path\twhere it is (default=$THE_PATH)",
}, {
	flag: flag.StringFlag{Name: "path", Value: "/", EnvVar: "THE_PATH", Usage: "where it is"},
	str:  "   --path\twhere it is (default=$THE_PATH,\"/\")",
}, {
	flag: flag.StringParam{Name: "path", Usage: "where it is"},
	str:  " * <path>\twhere it is",
}, {
	flag: flag.BoolFlag{Name: "force", Alias: "f", Usage: "forcefully"},
	str:  "   --force, -f\tforcefully",
}}

func TestString(t *testing.T) {
	for _, c := range stringCases {
		actual := flagHelp(c.flag)
		assert.Equal(t, c.str, actual)
	}
}
