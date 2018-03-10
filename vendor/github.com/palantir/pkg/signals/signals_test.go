// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signals_test

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/palantir/pkg/signals"
)

func TestCancelOnSignalsContext(t *testing.T) {
	ctx, _ := signals.CancelOnSignalsContext(context.Background(), syscall.SIGHUP)

	sendSignalToCurrProcess(t, syscall.SIGHUP)

	timer := time.NewTimer(time.Second * 3)
	done := false
	select {
	case <-ctx.Done():
		done = true
	case <-timer.C:
	}
	assert.True(t, done)
}

func TestRegisterStackTraceWriterOnSignals(t *testing.T) {
	out := &bytes.Buffer{}
	signals.RegisterStackTraceWriterOnSignals(out, nil, syscall.SIGHUP)

	sendSignalToCurrProcess(t, syscall.SIGHUP)

	// output stack should contain current routine
	assert.Contains(t, out.String(), "signals_test.TestRegisterStackTraceWriterOnSignals")
}

type errWriter struct{}

func (w errWriter) Write(p []byte) (n int, err error) {
	return 0, fmt.Errorf("test error")
}

func TestRegisterStackTraceWriterErrorHandler(t *testing.T) {
	out := errWriter{}
	var handlerErr error
	errHandler := func(err error) {
		handlerErr = err
	}

	signals.RegisterStackTraceWriterOnSignals(out, errHandler, syscall.SIGHUP)

	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	err = proc.Signal(syscall.SIGHUP)
	require.NoError(t, err)

	// add sleep because write to buffer happens on a separate channel
	time.Sleep(1 * time.Second)

	// handler should have been called with error
	assert.EqualError(t, handlerErr, "test error")
}

func TestUnregisterStackTraceWriterOnSignals(t *testing.T) {
	out := &bytes.Buffer{}
	unregister := signals.RegisterStackTraceWriterOnSignals(out, nil, syscall.SIGHUP)
	unregister()

	sendSignalToCurrProcess(t, syscall.SIGHUP)

	// output stack should be empty
	assert.Empty(t, out.String())
}

func TestNewSignalReceiver(t *testing.T) {
	c := signals.NewSignalReceiver(syscall.SIGHUP)

	sendSignalToCurrProcess(t, syscall.SIGHUP)

	timer := time.NewTimer(time.Second * 3)
	done := false
	select {
	case <-c:
		done = true
	case <-timer.C:
	}
	assert.True(t, done)
}

func sendSignalToCurrProcess(t *testing.T, sig os.Signal) {
	proc, err := os.FindProcess(os.Getpid())
	require.NoError(t, err)
	go func() {
		err = proc.Signal(sig)
		require.NoError(t, err)
	}()

	// add sleep because write to buffer happens on a separate channel
	time.Sleep(1 * time.Second)
}
