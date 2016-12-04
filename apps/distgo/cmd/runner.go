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

package cmd

import (
	"fmt"
	"io"

	"github.com/palantir/godel/apps/distgo/params"
)

type ProcessFunc func(f func(buildSpec params.ProductBuildSpecWithDeps, stdout io.Writer) error) BuildFunc
type BuildFunc func(buildSpecWithDeps []params.ProductBuildSpecWithDeps, stdout io.Writer) error

// ProcessSerially returns a BuildFunc that processes each of the provided specs in order using the provided function.
// If the function returns an error for any of the specifications, the function immediately returns that error.
func ProcessSerially(f func(buildSpec params.ProductBuildSpecWithDeps, stdout io.Writer) error) BuildFunc {
	return func(buildSpec []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		for _, currSpec := range buildSpec {
			if err := f(currSpec, stdout); err != nil {
				return err
			}
		}
		return nil
	}
}

type SpecErrors struct {
	Errors map[string]error
}

func (e *SpecErrors) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

// ProcessSeriallyBatchErrors returns a BuildFunc that processes each of the provided specs in order using the provided
// function. If the function returns an error for any of the specifications, it is stored, but the function will
// will continue processing the provided specifications. The function return nil if no errors occurred; otherwise, it
// returns a SpecErrors error that contains the individual errors.
func ProcessSeriallyBatchErrors(f func(buildSpec params.ProductBuildSpecWithDeps, stdout io.Writer) error) BuildFunc {
	return func(buildSpec []params.ProductBuildSpecWithDeps, stdout io.Writer) error {
		errors := make(map[string]error)
		for _, currSpec := range buildSpec {
			if err := f(currSpec, stdout); err != nil {
				errors[currSpec.Spec.ProductName] = err
			}
		}
		if len(errors) > 0 {
			return &SpecErrors{Errors: errors}
		}
		return nil
	}
}
