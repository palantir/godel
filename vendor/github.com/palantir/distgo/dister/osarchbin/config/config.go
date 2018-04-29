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

package config

import (
	"github.com/palantir/godel/pkg/osarch"

	"github.com/palantir/distgo/dister/osarchbin"
	"github.com/palantir/distgo/dister/osarchbin/config/internal/v0"
	"github.com/palantir/distgo/distgo"
)

type OSArchBin v0.Config

func (cfg *OSArchBin) ToDister() distgo.Dister {
	osArchs := cfg.OSArchs
	if len(osArchs) == 0 {
		osArchs = []osarch.OSArch{osarch.Current()}
	}
	return &osarchbin.Dister{
		OSArchs: osArchs,
	}
}
