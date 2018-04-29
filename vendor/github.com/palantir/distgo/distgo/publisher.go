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

package distgo

import (
	"io"

	"github.com/pkg/errors"
	"github.com/spf13/pflag"
)

type Publisher interface {
	// TypeName returns the type of this publisher.
	TypeName() (string, error)

	// Flags returns the flags provided by this Publisher.
	Flags() ([]PublisherFlag, error)

	// RunPublish runs the publish task. When this function is called, the distribution artifacts for the product should
	// already exist. If dryRun is true, then the task should print the operations that would occur without actually
	// executing them.
	RunPublish(productTaskOutputInfo ProductTaskOutputInfo, cfgYML []byte, flagVals map[PublisherFlagName]interface{}, dryRun bool, stdout io.Writer) error
}

type PublisherFactory interface {
	Types() []string
	NewPublisher(typeName string) (Publisher, error)
	ConfigUpgrader(typeName string) (ConfigUpgrader, error)
}

type PublisherFlagName string

type PublisherFlag struct {
	Name        PublisherFlagName
	Description string
	Type        FlagType
}

// AddFlag adds the flag represented by PublisherFlag to the provided pflag.FlagSet. Returns the pointer to the value
// that can be used to retrieve the value.
func (f PublisherFlag) AddFlag(fset *pflag.FlagSet) (interface{}, error) {
	switch f.Type {
	case StringFlag:
		return fset.String(string(f.Name), "", f.Description), nil
	case BoolFlag:
		return fset.Bool(string(f.Name), false, f.Description), nil
	default:
		return nil, errors.Errorf("unrecognized flag type: %v", f.Type)
	}
}

func (f PublisherFlag) GetFlagValue(fset *pflag.FlagSet) (interface{}, error) {
	var val interface{}
	var err error
	switch f.Type {
	case StringFlag:
		val, err = fset.GetString(string(f.Name))
	case BoolFlag:
		val, err = fset.GetBool(string(f.Name))
	default:
		return nil, errors.Errorf("unrecognized flag type: %v", f.Type)
	}
	return val, errors.Wrapf(err, "failed to get flag value")
}

// ToFlagArgs takes the input parameter (which should be the value returned by calling AddFlag for the receiver) and
// returns a string slice that reconstructs the flag arguments for the given flag.
func (f PublisherFlag) ToFlagArgs(flagVal interface{}) ([]string, error) {
	switch f.Type {
	case StringFlag:
		flagValStr := flagVal.(*string)
		if flagValStr == nil || len(*flagValStr) == 0 {
			return nil, nil
		}
		return []string{"--" + string(f.Name), *flagValStr}, nil
	case BoolFlag:
		flagValBool := flagVal.(*bool)
		if flagValBool == nil || *flagValBool == false {
			return nil, nil
		}
		return []string{"--" + string(f.Name)}, nil
	default:
		return nil, errors.Errorf("unrecognized flag type: %v", f.Type)
	}
}

// FlagType represents the type of a flag (string, boolean, etc.). Currently only string flags are supported.
type FlagType int

const (
	StringFlag FlagType = iota
	BoolFlag
)
