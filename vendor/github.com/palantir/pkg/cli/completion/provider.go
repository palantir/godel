// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

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
