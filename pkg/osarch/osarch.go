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

package osarch

import (
	"fmt"
	"runtime"
	"strings"
)

type OSArch struct {
	OS   string
	Arch string
}

// String returns a string representation of the form "GOOS-GOARCH".
func (o OSArch) String() string {
	return fmt.Sprintf("%v-%v", o.OS, o.Arch)
}

// Current returns and OSArch that reflects the GOOS/GOARCH value for the current runtime.
func Current() OSArch {
	return OSArch{
		OS:   runtime.GOOS,
		Arch: runtime.GOARCH,
	}
}

// New returns an OSArch value for the provided input. Returns an error if the provided input is not of the correct
// format, which is defined as two non-empty alphanumeric strings separated by a single hyphen ("-").
func New(input string) (OSArch, error) {
	if parts := strings.Split(input, "-"); len(parts) == 2 && isAlphaNumericOnly(parts[0]) && isAlphaNumericOnly(parts[1]) {
		return OSArch{OS: parts[0], Arch: parts[1]}, nil
	}
	return OSArch{}, fmt.Errorf("not a valid OSArch value: %s", input)
}

func isAlphaNumericOnly(input string) bool {
	if input == "" {
		return false
	}
	for _, r := range input {
		if !((r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			return false
		}
	}
	return true
}
