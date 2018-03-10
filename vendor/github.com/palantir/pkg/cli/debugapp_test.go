// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli_test

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/palantir/pkg/cli"
)

type testExitCoder struct {
	error
	exitCode int
}

func (t testExitCoder) ExitCode() int {
	return t.exitCode
}

func staticErrorStringer(val string) cli.ErrorStringer {
	return cli.ErrorStringer(func(err error) string {
		return val
	})
}

func TestNewDebugApp(t *testing.T) {
	testErrorStringer := staticErrorStringer("error-stringer")

	for i, currCase := range []struct {
		err            error
		wantExitCode   int
		errorStringer  cli.ErrorStringer
		wantDebugFalse string
		wantDebugTrue  string
	}{
		// output error message if debug is false, output result of ErrorStringer if debug is true
		{err: fmt.Errorf("default-error"), wantExitCode: 1, errorStringer: testErrorStringer, wantDebugFalse: "^default-error\n$", wantDebugTrue: "^error-stringer\n$"},
		// output is empty (no newline) if debug is false and error message is empty, output result of ErrorStringer if debug is true
		{err: fmt.Errorf(""), wantExitCode: 1, errorStringer: testErrorStringer, wantDebugFalse: "^$", wantDebugTrue: "^error-stringer\n$"},
		// output only message if ErrorStringer is nil
		{err: fmt.Errorf("foo"), wantExitCode: 1, wantDebugFalse: "^foo\n$", wantDebugTrue: "^foo\n$"},
		// output is empty (no newline) if debut is true and output of ErrorStringer is empty
		{err: fmt.Errorf("foo"), wantExitCode: 1, errorStringer: staticErrorStringer(""), wantDebugFalse: "^foo\n$", wantDebugTrue: "^$"},
		// exit code is the value returned by ExitCode if error implements ExitCoder
		{err: testExitCoder{error: fmt.Errorf("foo"), exitCode: 2}, wantExitCode: 2, errorStringer: testErrorStringer, wantDebugFalse: "^foo\n$", wantDebugTrue: "^error-stringer\n$"},
	} {
		app := cli.NewApp(cli.DebugHandler(currCase.errorStringer))
		app.Action = func(ctx cli.Context) error {
			return currCase.err
		}
		stderr := &bytes.Buffer{}

		app.Stderr = stderr
		exitCode := app.Run([]string{"app", "--debug=false"})
		assert.Equal(t, currCase.wantExitCode, exitCode, "Case %d", i)
		assert.Regexp(t, currCase.wantDebugFalse, stderr.String(), "Case %d, debug=false", i)

		stderr.Reset()
		exitCode = app.Run([]string{"app", "--debug=true"})
		assert.Equal(t, currCase.wantExitCode, exitCode, "Case %d", i)
		assert.Regexp(t, currCase.wantDebugTrue, stderr.String(), "Case %d, debug=true", i)
	}
}
