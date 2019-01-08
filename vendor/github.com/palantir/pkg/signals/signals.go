// Copyright (c) 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package signals

import (
	"bytes"
	"context"
	"io"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

// ContextWithShutdown returns a context that is cancelled when the SIGTERM or SIGINT signal is received. This can be
// useful to use for cleanup tasks that would typically be deferred and which you want to run even if the program is
// terminated or cancelled. The following is a common usage pattern:
//
//   cleanupCtx, cancel := signals.ContextWithShutdown(context.Background())
//   cleanupDone := make(chan struct{})
//   defer func() {
//   	cancel()
//   	<-cleanupDone
//   }()
//   go func() {
//   	select {
//   	case <-cleanupCtx.Done():
//   		// perform cleanup action
//   	}
//   	cleanupDone <- struct{}{}
//   }()
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
	ctx, cancel := context.WithCancel(context.Background())
	handler := func(stackTraceOutput []byte) error {
		_, err := out.Write(stackTraceOutput)
		return err
	}
	RegisterStackTraceHandlerOnSignals(ctx, handler, errHandler, sig...)
	return cancel
}

// RegisterStackTraceHandlerOnSignals starts a goroutine that listens for the specified signals and calls stackTraceHandler with a
// pprof-formatted snapshot of all running goroutines when any of the provided signals are received. If stackTraceHandler returns
// an error, that error is provided to the errHandler function if one is provided. No goroutine is created if stackTraceHandler is nil.
// If no signals are provided, the handler will receive notifications for all signals (matching the os/signal.Notify API).
// The handler will exit when ctx is cancelled.
func RegisterStackTraceHandlerOnSignals(ctx context.Context, stackTraceHandler func(stackTraceOutput []byte) error, errHandler func(error), sig ...os.Signal) {
	if stackTraceHandler == nil {
		return
	}
	signals := NewSignalReceiver(sig...)
	go func() {
		for {
			select {
			case <-signals:
				var buf bytes.Buffer
				_ = pprof.Lookup("goroutine").WriteTo(&buf, 2) // bytes.Buffer's Write never returns an error, so we swallow it
				if err := stackTraceHandler(buf.Bytes()); err != nil && errHandler != nil {
					errHandler(err)
				}
			case <-ctx.Done():
				return
			}
		}
	}()
}

// NewSignalReceiver returns a buffered channel that is registered to receive the provided signals.
func NewSignalReceiver(sig ...os.Signal) <-chan os.Signal {
	// Use a buffer of 1 in case we are not ready when the signal arrives
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, sig...)
	return signals
}
