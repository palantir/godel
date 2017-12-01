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

package dirchecksum

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// mustDefer takes a function that contains cleanup logic (a function meant to be provided to a defer call) and ensures
// that the function is called even if the program is killed using a SIGTERM or SIGINT signal. Callers should defer the
// returned function to ensure that the signal listener registered by this function will be unregistered if the cleanup
// function is called normally in a defer. Example usage:
//
//   fn := cleanup.MustDefer(func() {
//   	_ = os.Remove(tmpFile)
//   })
//   defer fn()
//
// Note that the "cleanup.MustDefer" call must occur before the defer -- if the call is inlined as part of the defer,
// then "MustDefer" will not run until the defer action is run, which will fail to properly register the signal
// listener.
//
// This function starts a goroutine that listens for the SIGTERM and SIGINT signals. This goroutine will be terminated
// when the returned function is called. The returned function guarantees that the provided cleanup function will only
// be run once. The signal listener terminates the program using "os.Exit(1)" after running the cleanup function, so
// if other signal listeners are registered they are not guaranteed to be run.
func mustDefer(fn func()) func() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	// use sync.Once to ensure that function being deferred is only called once. Guards against the case where an
	// interrupt signal is received during a call in the defer -- in such a case, if the defer call is already executing
	// the function, the signal handler should not invoke the function again.
	var once sync.Once

	// start goroutine that listens for signal and runs cleanup function is signal is received
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		select {
		case <-signals:
			once.Do(fn)
			os.Exit(1)
		case <-ctx.Done():
			return
		}
	}()

	return func() {
		cancel()
		once.Do(fn)
	}
}
