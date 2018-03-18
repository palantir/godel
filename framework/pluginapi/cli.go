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

package pluginapi

import (
	"io"
)

func InfoCmd(osArgs []string, stdout io.Writer, info PluginInfo) bool {
	if len(osArgs) < 2 || osArgs[1] != PluginInfoCommandName {
		return false
	}
	_ = infoAction(info, stdout)
	return true
}
