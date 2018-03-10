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
	"fmt"
	"io"
)

type ProcessFunc func(f func(projectInfo ProjectInfo, productParam ProductParam, stdout io.Writer) error) BuildFunc
type BuildFunc func(projectInfo ProjectInfo, productParams []ProductParam, stdout io.Writer) error

// ProcessSerially returns a BuildFunc that processes each of the provided specs in order using the provided function.
// If the function returns an error for any of the specifications, the function immediately returns that error.
func ProcessSerially(f func(projectInfo ProjectInfo, productParam ProductParam, stdout io.Writer) error) BuildFunc {
	return func(projectInfo ProjectInfo, productParams []ProductParam, stdout io.Writer) error {
		for _, currParam := range productParams {
			if err := f(projectInfo, currParam, stdout); err != nil {
				return err
			}
		}
		return nil
	}
}

type ProductErrors struct {
	Errors map[ProductID]error
}

func (e *ProductErrors) Error() string {
	return fmt.Sprintf("%v", e.Errors)
}

// ProcessSeriallyBatchErrors returns a BuildFunc that processes each of the provided specs in order using the provided
// function. If the function returns an error for any of the specifications, it is stored, but the function will
// will continue processing the provided specifications. The function return nil if no errors occurred; otherwise, it
// returns a ProductErrors error that contains the individual errors.
func ProcessSeriallyBatchErrors(f func(projectInfo ProjectInfo, productParam ProductParam, stdout io.Writer) error) BuildFunc {
	return func(projectInfo ProjectInfo, productParams []ProductParam, stdout io.Writer) error {
		errors := make(map[ProductID]error)
		for _, currParam := range productParams {
			if err := f(projectInfo, currParam, stdout); err != nil {
				errors[currParam.ID] = err
			}
		}
		if len(errors) == 0 {
			return nil
		}
		return &ProductErrors{Errors: errors}
	}
}
