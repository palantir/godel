// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signals

import (
	"context"
	"io"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// ContextWithShutdown returns a context that is cancelled when the SIGTERM or SIGINT signal is received.
func ContextWithShutdown(ctx context.Context) (context.Context, context.CancelFunc) {
	return CancelOnSignalsContext(ctx, syscall.SIGTERM, syscall.SIGINT)
}

// CancelOnSignalsContext returns a context that is cancelled when any of the provided signals are received.
func CancelOnSignalsContext(ctx context.Context, sig ...os.Signal) (context.Context, context.CancelFunc) {
	newCtx, cancel := context.WithCancel(ctx)

	signals := NewSignalReceiver(sig...)
	go func() {
		<-signals
		cancel()
	}()

	return newCtx, cancel
}

// RegisterStackTraceWriter starts a goroutine that listens for the SIGQUIT (kill -3) signal and writes a
// pprof-formatted snapshot of all running goroutines when the signal is received. If writing to out returns an
// error, that error is provided to the errHandler function if one is provided. Returns a function that unregisters the
// listener when called.
func RegisterStackTraceWriter(out io.Writer, errHandler func(error)) (unregister func()) {
	return RegisterStackTraceWriterOnSignals(out, errHandler, syscall.SIGQUIT)
}

// RegisterStackTraceWriterOnSignals starts a goroutine that listens for the specified signals and writes a pprof-formatted
// snapshot of all running goroutines to out when any of the provided signals are received. If writing to out returns an
// error, that error is provided to the errHandler function if one is provided. Returns a function that unregisters the
// listener when called.
func RegisterStackTraceWriterOnSignals(out io.Writer, errHandler func(error), sig ...os.Signal) (unregister func()) {
	cancel := make(chan bool, 1)
	unregister = func() {
		cancel <- true
	}

	signals := NewSignalReceiver(sig...)
	go func() {
		for {
			select {
			case <-signals:
				err := pprof.Lookup("goroutine").WriteTo(out, 2)
				if err != nil && errHandler != nil {
					errHandler(err)
				}
			case <-cancel:
				return
			}
		}
	}()
	return unregister
}

// NewSignalReceiver returns a buffered channel that is registered to receive the provided signals.
func NewSignalReceiver(sig ...os.Signal) <-chan os.Signal {
	// Use a buffer of 1 in case we are not ready when the signal arrives
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sig...)
	return signals
}
