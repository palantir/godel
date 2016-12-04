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

package params

type Publish struct {
	// GroupID is the product-specific configuration equivalent to the global GroupID configuration.
	GroupID string
	// Almanac contains the parameters for Almanac publish operations. Optional.
	Almanac Almanac
}

type Almanac struct {
	// Metadata contains the metadata provided to the Almanac publish task.
	Metadata map[string]string
	// Tags contains the tags provided to the Almanac publish task.
	Tags []string
}

func (a *Almanac) empty() bool {
	return len(a.Metadata) == 0 && len(a.Tags) == 0
}

func (pub *Publish) empty() bool {
	return pub.GroupID == "" && pub.Almanac.empty()
}
