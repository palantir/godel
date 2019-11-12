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
	"github.com/palantir/godel/v2/framework/godellauncher"
)

// GlobalFlagOptions is a JSON-serializable interface that can be translated into a godellauncher.GlobalFlagOptions.
// See godellauncher.GlobalFlagOptions for documentation.
type GlobalFlagOptions interface {
	DebugFlag() string
	ProjectDirFlag() string
	GodelConfigFlag() string
	ConfigFlag() string

	toGodelGlobalFlagOptions() godellauncher.GlobalFlagOptions
}

// globalFlagOptionsImpl is a concrete implementation of GlobalFlagOptions. Note that the functions are defined on
// non-pointer receivers to reduce bugs in calling functions in closures.
type globalFlagOptionsImpl struct {
	DebugFlagVar       string `json:"debugFlag"`
	ProjectDirFlagVar  string `json:"projectDirFlag"`
	GodelConfigFlagVar string `json:"godelConfigFlag"`
	ConfigFlagVar      string `json:"configFlag"`
}

type GlobalFlagOptionsParam interface {
	apply(*globalFlagOptionsImpl)
}

type globalFlagOptionsParamFunc func(*globalFlagOptionsImpl)

func (f globalFlagOptionsParamFunc) apply(impl *globalFlagOptionsImpl) {
	f(impl)
}

func GlobalFlagOptionsParamDebugFlag(debugFlag string) GlobalFlagOptionsParam {
	return globalFlagOptionsParamFunc(func(impl *globalFlagOptionsImpl) {
		impl.DebugFlagVar = debugFlag
	})
}

func GlobalFlagOptionsParamProjectDirFlag(projectDirFlag string) GlobalFlagOptionsParam {
	return globalFlagOptionsParamFunc(func(impl *globalFlagOptionsImpl) {
		impl.ProjectDirFlagVar = projectDirFlag
	})
}

func GlobalFlagOptionsParamGodelConfigFlag(godelConfigFlag string) GlobalFlagOptionsParam {
	return globalFlagOptionsParamFunc(func(impl *globalFlagOptionsImpl) {
		impl.GodelConfigFlagVar = godelConfigFlag
	})
}

func GlobalFlagOptionsParamConfigFlag(configFlag string) GlobalFlagOptionsParam {
	return globalFlagOptionsParamFunc(func(impl *globalFlagOptionsImpl) {
		impl.ConfigFlagVar = configFlag
	})
}

func newGlobalFlagOptionsImpl(params ...GlobalFlagOptionsParam) *globalFlagOptionsImpl {
	if len(params) == 0 {
		return nil
	}
	impl := &globalFlagOptionsImpl{}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(impl)
	}
	return impl
}

func (g globalFlagOptionsImpl) DebugFlag() string {
	return g.DebugFlagVar
}

func (g globalFlagOptionsImpl) ProjectDirFlag() string {
	return g.ProjectDirFlagVar
}

func (g globalFlagOptionsImpl) GodelConfigFlag() string {
	return g.GodelConfigFlagVar
}

func (g globalFlagOptionsImpl) ConfigFlag() string {
	return g.ConfigFlagVar
}

func (g globalFlagOptionsImpl) toGodelGlobalFlagOptions() godellauncher.GlobalFlagOptions {
	return godellauncher.GlobalFlagOptions{
		DebugFlag:       g.DebugFlagVar,
		ProjectDirFlag:  g.ProjectDirFlagVar,
		GodelConfigFlag: g.GodelConfigFlagVar,
		ConfigFlag:      g.ConfigFlagVar,
	}
}
