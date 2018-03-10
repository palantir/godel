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

package printproducts

import (
	"fmt"
	"io"
	"sort"

	"github.com/palantir/distgo/distgo"
)

func Run(projectParam distgo.ProjectParam, stdout io.Writer) error {
	var productIDs []distgo.ProductID
	for currProductID := range projectParam.Products {
		productIDs = append(productIDs, currProductID)
	}
	sort.Sort(distgo.ByProductID(productIDs))
	for _, currProductID := range productIDs {
		fmt.Fprintln(stdout, currProductID)
	}
	return nil
}
