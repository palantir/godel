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

package config_test

import (
	"fmt"

	"github.com/palantir/godel/apps/okgo/config"
)

func Example() {
	yml := `
release-tag: go1.7
checks:
  errcheck:
    args:
      - "-ignore"
      - "github.com/seelog:(Info|Warn|Error|Critical)f?"
    filters:
      - type: "message"
        value: "\\w+"
`
	cfg, err := config.LoadRawConfig(yml, "")
	if err != nil {
		panic(err)
	}
	if _, err := cfg.ToParams(); err != nil {
		panic(err)
	}
	fmt.Printf("%q", fmt.Sprintf("%+v", cfg))
	// Output: "{ReleaseTag:go1.7 Checks:map[errcheck:{Skip:false Args:[-ignore github.com/seelog:(Info|Warn|Error|Critical)f?] Filters:[{Type:message Value:\\w+}]}] Exclude:{Names:[] Paths:[]}}"
}
