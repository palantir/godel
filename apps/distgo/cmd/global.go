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

package cmd

import (
	"fmt"
	"strings"

	"github.com/palantir/pkg/cli/flag"

	"github.com/palantir/godel/apps/distgo/pkg/osarch"
)

const (
	ProductsParamName = "products"
	OSArchFlagName    = "os-arch"
)

var (
	ProductsParam = flag.StringSlice{
		Name:     ProductsParamName,
		Usage:    "Products for which action should be performed",
		Optional: true,
	}
	OSArchFlag = flag.StringFlag{
		Name:  OSArchFlagName,
		Usage: "GOOS-GOARCH for the command (comma-separate for multiple values)",
	}
)

type OSArchFilter []osarch.OSArch

// Matches returns true if the provided osArch is in the filter list or if the filter list is empty.
func (f OSArchFilter) Matches(osArch osarch.OSArch) bool {
	if len(f) == 0 {
		return true
	}
	for _, curr := range f {
		if curr == osArch {
			return true
		}
	}
	return false
}

func NewOSArchFilter(osArchs string) (OSArchFilter, error) {
	if osArchs == "" {
		return nil, nil
	}

	var filterArchs []osarch.OSArch
	var invalidValues []string
	// if value was provided for flag, parse
	for _, osArchStr := range strings.Split(osArchs, ",") {
		if osArch, err := osarch.New(osArchStr); err == nil {
			filterArchs = append(filterArchs, osArch)
		} else {
			invalidValues = append(invalidValues, osArchStr)
		}
	}
	if len(invalidValues) > 0 {
		return nil, fmt.Errorf("invalid os-arch values: %v", invalidValues)
	}

	return OSArchFilter(filterArchs), nil
}
