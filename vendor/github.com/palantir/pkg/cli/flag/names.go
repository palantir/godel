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

package flag

import (
	"strings"
)

func WithPrefix(name string) string {
	if len(name) == 1 {
		return "-" + name
	}
	return "--" + name
}

func placeholderOrDefault(maybePlaceholder, name string) string {
	if maybePlaceholder != "" {
		return maybePlaceholder
	}
	return defaultPlaceholder(name)
}

func defaultPlaceholder(name string) string {
	name = strings.ToUpper(name)
	name = strings.Replace(name, "-", "_", -1)
	return name
}
