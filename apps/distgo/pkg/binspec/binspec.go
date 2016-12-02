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

package binspec

import (
	"github.com/palantir/pkg/specdir"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

// New returns a LayoutSpec that is rooted at a "bin" directory and contains directories of the form "{{os}}-{{arch}}"
// that contain a file "execName" for each of the OSArch entries provided.
func New(targets []osarch.OSArch, execName string) specdir.LayoutSpec {
	providers := make([]specdir.FileNodeProvider, len(targets))
	for i, currOSArch := range targets {
		providers[i] = specdir.Dir(specdir.LiteralName(currOSArch.String()), currOSArch.String(),
			specdir.File(specdir.LiteralName(execName), ""))
	}

	return specdir.NewLayoutSpec(
		specdir.Dir(specdir.LiteralName("bin"), "",
			providers...,
		),
		true,
	)
}
