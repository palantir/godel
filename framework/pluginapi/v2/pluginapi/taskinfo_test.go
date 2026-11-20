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

package pluginapi

import (
	"bytes"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/palantir/godel/v2/framework/godellauncher"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTaskRunPluginExitError(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("test uses a shell script plugin")
	}

	for _, tc := range []struct {
		name    string
		script  string
		wantErr string
	}{
		{
			name:    "non-zero exit returns empty error",
			script:  "#!/bin/sh\necho 'plugin error output'\nexit 1\n",
			wantErr: "",
		},
		{
			name:    "termination by signal returns error describing the signal",
			script:  "#!/bin/sh\nkill -9 $$\n",
			wantErr: "plugin process for task test-task was terminated: signal: killed",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			pluginPath := filepath.Join(t.TempDir(), "plugin.sh")
			require.NoError(t, os.WriteFile(pluginPath, []byte(tc.script), 0755))

			task := taskInfoImpl{NameVar: "test-task"}.toTask(pluginPath, "", nil)
			buf := &bytes.Buffer{}
			err := task.Run(godellauncher.GlobalConfig{}, buf)
			require.Error(t, err)
			assert.Equal(t, tc.wantErr, err.Error())
		})
	}
}
