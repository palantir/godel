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

// +build !windows

package checkpath

import (
	"os"
	"syscall"
)

func pathsOnSameDevice(p1, p2 string) bool {
	id1, ok := getDeviceID(p1)
	if !ok {
		return false
	}
	id2, ok := getDeviceID(p2)
	if !ok {
		return false
	}
	return id1 == id2
}

func getDeviceID(p string) (interface{}, bool) {
	fi, err := os.Stat(p)
	if err != nil {
		return 0, false
	}
	s := fi.Sys()
	switch s := s.(type) {
	case *syscall.Stat_t:
		return s.Dev, true
	}
	return 0, false
}
