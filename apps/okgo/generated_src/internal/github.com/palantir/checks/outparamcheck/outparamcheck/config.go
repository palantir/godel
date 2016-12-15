// Copyright 2016 Palantir Technologies, Inc. All rights reserved.
// Licensed under the MIT License. See LICENSE in the project root
// for license information.

package outparamcheck

// Config stores a map from function name to the argument indices which are output parameters.
type Config map[string][]int

var defaultCfg = Config(
	map[string][]int{
		"encoding/json.Unmarshal":	{1},
		"encoding/safejson.Unmarshal":	{1},
		"gopkg.in/yaml.v2.Unmarshal":	{1},
	},
)
