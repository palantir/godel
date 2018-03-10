// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOnExit(t *testing.T) {
	ran := false
	f := func() {
		ran = true
	}

	onExit := newOnExit()
	onExit.Register(f)
	onExit.run()

	assert.True(t, ran)
}

func TestOnExitUnregister(t *testing.T) {
	ran := false
	f := func() {
		ran = true
	}

	onExit := newOnExit()
	id := onExit.Register(f)
	priorityID := onExit.register(f, highPriority)
	onExit.Unregister(id)
	onExit.Unregister(priorityID)
	onExit.run()

	assert.False(t, ran)
}

func TestOnExitWithPanic(t *testing.T) {
	ran := false
	f1 := func() {
		panic("f1 panic")
	}
	f2 := func() {
		ran = true
	}

	onExit := newOnExit()
	onExit.Register(f1)
	onExit.Register(f2)
	onExit.run()

	assert.True(t, ran)
}

func TestRunInOrder(t *testing.T) {
	var got []int

	addToGot := func(n int) func() {
		return func() {
			got = append(got, n)
		}
	}

	onExit := newOnExit()
	onExit.Register(addToGot(0))
	oneID := onExit.Register(addToGot(1))
	onExit.Register(addToGot(2))
	onExit.Register(addToGot(3))
	onExit.Register(addToGot(4))
	onExit.Register(addToGot(5))
	onExit.register(addToGot(10), highPriority)
	onExit.Unregister(oneID)
	onExit.register(addToGot(11), highPriority)

	onExit.run()

	assert.Equal(t, []int{11, 10, 5, 4, 3, 2, 0}, got)
}
