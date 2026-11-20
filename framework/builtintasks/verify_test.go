// Copyright 2026 Palantir Technologies, Inc.
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

package builtintasks_test

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/palantir/godel/v2/framework/builtintasks"
	"github.com/palantir/godel/v2/framework/godel/config"
	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVerifyPrintsTaskErrors(t *testing.T) {
	origArgs := os.Args
	defer func() {
		os.Args = origArgs
	}()

	newVerifyTask := func(name string, ordering int, runErr error) godellauncher.Task {
		return godellauncher.Task{
			Name: name,
			Verify: &godellauncher.VerifyOptions{
				Ordering:       ordering,
				ApplyFalseArgs: []string{"--verify"},
			},
			RunImpl: func(t *godellauncher.Task, global godellauncher.GlobalConfig, stdout io.Writer) error {
				return runErr
			},
		}
	}

	verifyTask := builtintasks.VerifyTask([]godellauncher.Task{
		newVerifyTask("error-task", 0, fmt.Errorf("plugin process for task error-task was terminated: signal: killed")),
		newVerifyTask("empty-error-task", 1, fmt.Errorf("")),
		newVerifyTask("success-task", 2, nil),
	}, config.VerifyTasksConfig{})

	buf := &bytes.Buffer{}
	err := verifyTask.Run(godellauncher.GlobalConfig{
		Executable: "godel",
		Task:       "verify",
		TaskArgs:   []string{"--apply=false"},
	}, buf)
	require.Error(t, err)

	assert.Equal(t, `Running error-task...
Error: plugin process for task error-task was terminated: signal: killed
Running empty-error-task...
Running success-task...
Failed tasks:
	error-task --verify
	empty-error-task --verify
`, buf.String())
}
