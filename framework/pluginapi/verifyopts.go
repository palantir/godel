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
	"github.com/palantir/godel/framework/godellauncher"
	"github.com/palantir/godel/framework/verifyorder"
)

// VerifyOptions is a JSON-serializable interface that can be translated into a godellauncher.VerifyOptions. Refer to
// that struct for field documentation.
type VerifyOptions interface {
	VerifyTaskFlags() []VerifyFlag
	Ordering() *int
	ApplyTrueArgs() []string
	ApplyFalseArgs() []string

	toGodelVerifyOptions() godellauncher.VerifyOptions
}

// verifyOptionsImpl is a concrete implementation of VerifyOptions. Note that the functions are defined on non-pointer
// receivers to reduce bugs in calling functions in closures.
type verifyOptionsImpl struct {
	VerifyTaskFlagsVar []verifyFlagImpl `json:"verifyTaskFlags"`
	OrderingVar        *int             `json:"ordering"`
	ApplyTrueArgsVar   []string         `json:"applyTrueArgs"`
	ApplyFalseArgsVar  []string         `json:"applyFalseArgs"`
}

type VerifyOptionsParam interface {
	apply(*verifyOptionsImpl)
}

type verifyOptsFunc func(*verifyOptionsImpl)

func (f verifyOptsFunc) apply(impl *verifyOptionsImpl) {
	f(impl)
}

func VerifyOptionsTaskFlags(flags ...VerifyFlag) VerifyOptionsParam {
	return verifyOptsFunc(func(impl *verifyOptionsImpl) {
		for _, f := range flags {
			impl.VerifyTaskFlagsVar = append(impl.VerifyTaskFlagsVar, verifyFlagImpl{
				NameVar:        f.Name(),
				DescriptionVar: f.Description(),
				TypeVar:        f.Type(),
			})
		}
	})
}

func VerifyOptionsOrdering(ordering *int) VerifyOptionsParam {
	return verifyOptsFunc(func(impl *verifyOptionsImpl) {
		impl.OrderingVar = ordering
	})
}

func VerifyOptionsApplyTrueArgs(args ...string) VerifyOptionsParam {
	return verifyOptsFunc(func(impl *verifyOptionsImpl) {
		impl.ApplyTrueArgsVar = args
	})
}

func VerifyOptionsApplyFalseArgs(args ...string) VerifyOptionsParam {
	return verifyOptsFunc(func(impl *verifyOptionsImpl) {
		impl.ApplyFalseArgsVar = args
	})
}

func newVerifyOptionsImpl(params ...VerifyOptionsParam) verifyOptionsImpl {
	vOpts := verifyOptionsImpl{}
	for _, p := range params {
		if p == nil {
			continue
		}
		p.apply(&vOpts)
	}
	return vOpts
}

func (vo verifyOptionsImpl) VerifyTaskFlags() []VerifyFlag {
	var flags []VerifyFlag
	for _, flag := range vo.VerifyTaskFlagsVar {
		flags = append(flags, flag)
	}
	return flags
}

func (vo verifyOptionsImpl) Ordering() *int {
	return vo.OrderingVar
}

func (vo verifyOptionsImpl) ApplyTrueArgs() []string {
	return vo.ApplyTrueArgsVar
}

func (vo verifyOptionsImpl) ApplyFalseArgs() []string {
	return vo.ApplyFalseArgsVar
}

func (vo verifyOptionsImpl) toGodelVerifyOptions() godellauncher.VerifyOptions {
	var flags []godellauncher.VerifyFlag
	for _, f := range vo.VerifyTaskFlagsVar {
		flags = append(flags, f.toGodelVerifyFlag())
	}

	ordering := verifyorder.Default
	if vo.OrderingVar != nil {
		ordering = *vo.OrderingVar
	}
	return godellauncher.VerifyOptions{
		VerifyTaskFlags: flags,
		Ordering:        ordering,
		ApplyTrueArgs:   vo.ApplyTrueArgsVar,
		ApplyFalseArgs:  vo.ApplyFalseArgsVar,
	}
}

// VerifyFlag is a JSON-serializable interface that can be translated into a godellauncher.VerifyFlag. Refer to that
// struct for field documentation.
type VerifyFlag interface {
	Name() string
	Description() string
	Type() godellauncher.FlagType

	toGodelVerifyFlag() godellauncher.VerifyFlag
}

// verifyFlagImpl is a concrete implementation of VerifyFlag. Note that the functions are defined on non-pointer
// receivers to reduce bugs in calling functions in closures.
type verifyFlagImpl struct {
	NameVar        string                 `json:"name"`
	DescriptionVar string                 `json:"description"`
	TypeVar        godellauncher.FlagType `json:"type"`
}

func NewVerifyFlag(name, description string, typ godellauncher.FlagType) VerifyFlag {
	return verifyFlagImpl{
		NameVar:        name,
		DescriptionVar: description,
		TypeVar:        typ,
	}
}

func (vf verifyFlagImpl) Name() string {
	return vf.NameVar
}

func (vf verifyFlagImpl) Description() string {
	return vf.DescriptionVar
}

func (vf verifyFlagImpl) Type() godellauncher.FlagType {
	return vf.TypeVar
}

func (vf verifyFlagImpl) toGodelVerifyFlag() godellauncher.VerifyFlag {
	return godellauncher.VerifyFlag{
		Name:        vf.NameVar,
		Description: vf.DescriptionVar,
		Type:        vf.TypeVar,
	}
}
