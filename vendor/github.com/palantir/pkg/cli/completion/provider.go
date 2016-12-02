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

package completion

type Provider func(ctx *ProviderCtx) []string

type ProviderCtx struct {
	Flags    map[string]string
	Command  []string
	Partial  string
	ExitCode *uint8
}

func Filepath(ctx *ProviderCtx) []string {
	*ctx.ExitCode = FilepathCode
	return nil
}

func Directory(ctx *ProviderCtx) []string {
	*ctx.ExitCode = DirectoryCode
	return nil
}

func CustomFilepath(ctx *ProviderCtx, files []string) []string {
	*ctx.ExitCode = CustomFilepathCode
	return files
}

func (ctx *ProviderCtx) CommandIs(cmd ...string) bool {
	if len(ctx.Command) != len(cmd) {
		return false
	}
	for i := range ctx.Command {
		if ctx.Command[i] != cmd[i] {
			return false
		}
	}
	return true
}

func (ctx *ProviderCtx) CommandBeginsWith(cmd ...string) bool {
	if len(ctx.Command) < len(cmd) {
		return false
	}
	for i := range cmd {
		if cmd[i] != ctx.Command[i] {
			return false
		}
	}
	return true
}
