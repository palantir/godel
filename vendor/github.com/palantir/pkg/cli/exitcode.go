// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

type ExitCoder interface {
	error
	ExitCode() int
}

type exitCodeError struct {
	error
	exitCode int
}

func (e *exitCodeError) ExitCode() int {
	return e.exitCode
}

func WithExitCode(exitCode int, err error) ExitCoder {
	return &exitCodeError{
		exitCode: exitCode,
		error:    err,
	}
}
