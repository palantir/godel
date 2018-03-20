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

package legacyplugins

type LegacyConfigUpgrader struct {
	LegacyConfigFileName string
}

var LegacyConfigUpgraders = map[string]LegacyConfigUpgrader{
	"com.palantir.okgo:check-plugin": {
		LegacyConfigFileName: "check.yml",
	},
	"com.palantir.distgo:dist-plugin": {
		LegacyConfigFileName: "dist.yml",
	},
	"com.palantir.godel-format-plugin:format-plugin": {
		LegacyConfigFileName: "format.yml",
	},
	"com.palantir.generate:generate-plugin": {
		LegacyConfigFileName: "generate.yml",
	},
	"com.palantir.go-license:license-plugin": {
		LegacyConfigFileName: "license.yml",
	},
	"com.palantir.godel-test-plugin:test-plugin": {
		LegacyConfigFileName: "test.yml",
	},
}

func ReservedConfigFileNames() map[string]struct{} {
	reservedNames := make(map[string]struct{})
	for _, v := range LegacyConfigUpgraders {
		reservedNames[v.LegacyConfigFileName] = struct{}{}
	}
	return reservedNames
}
