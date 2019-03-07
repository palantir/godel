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

package artifactresolver

import (
	"fmt"

	"github.com/palantir/godel/pkg/osarch"
)

type LocatorWithResolverParam struct {
	LocatorWithChecksums LocatorParam
	Resolver             Resolver
}

type LocatorParam struct {
	Locator
	Checksums map[osarch.OSArch]string
}

type Locator struct {
	Group   string
	Product string
	Version string
}

func (l Locator) String() string {
	return fmt.Sprintf("%s:%s:%s", l.Group, l.Product, l.Version)
}
