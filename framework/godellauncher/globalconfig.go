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

package godellauncher

import (
	"path/filepath"

	"github.com/pkg/errors"
)

// GlobalConfig stores the configuration provided to the initial invocation of gödel.
type GlobalConfig struct {
	// Path to the gödel executable
	Executable string
	// The value of the "--wrapper" flag provided to the gödel invocation.
	Wrapper string
	// True if the "--debug" flag was provided to the gödel invocation.
	Debug bool
	// True if the "--version" flag was provided to the gödel invocation.
	Version bool
	// True if the "--help" or "-h" flag was provided to the gödel invocation.
	Help bool
	// The first non-flag argument provided to the gödel executable. This is the task name.
	Task string
	// All of the arguments following the "Task" argument that was provided to the gödel executable.
	TaskArgs []string
}

// ProjectDir returns the project directory for the global configuration. Returns an error if the wrapper path was not
// specified.
func (g GlobalConfig) ProjectDir() (string, error) {
	if g.Wrapper == "" {
		return "", errors.Errorf("wrapper must be specified to determine project directory")
	}
	return filepath.Dir(g.Wrapper), nil
}
