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

package dist

import (
	"github.com/palantir/godel/apps/distgo/params"
)

var (
	distTypeOrdering = []params.DistInfoType{params.BinDistType,
		params.SLSDistType,
		params.RPMDistType,
		params.DockerDistType}
)

func OrderProductDists(dists []params.Dist) []params.Dist {
	var orderedDists []params.Dist
	for _, currDistType := range distTypeOrdering {
		for _, currDist := range dists {
			if currDist.Info.Type() != currDistType {
				continue
			}
			orderedDists = append(orderedDists, currDist)
		}
	}
	return orderedDists
}
