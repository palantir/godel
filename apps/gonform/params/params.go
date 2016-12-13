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

import (
	"github.com/palantir/pkg/matcher"
)

type Formatters struct {
	// Formatters specifies the configuration used by the formatters. The key is the name of the formatter and the
	// value is the custom configuration for that formatter.
	Formatters map[string]Formatter

	// Exclude specifies the files that should be excluded from formatting.
	Exclude matcher.Matcher
}

type Formatter struct {
	// Args specifies the command-line arguments that are provided to the formatter.
	Args []string
}
