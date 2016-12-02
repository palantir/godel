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

package cli

import (
	"fmt"
	"sort"
	"sync/atomic"
)

type OnExitID struct {
	id int32
}

type priority int

const (
	lowPriority  priority = 0
	highPriority          = 1000
)

type OnExit interface {
	// Registers a function to be executed on exit of the cli application. If multiple functions are registered
	// they are executed in LIFO order.
	Register(f func()) OnExitID
	Unregister(id OnExitID)
	run()
	register(f func(), p priority) OnExitID
}

func newOnExit() OnExit {
	return &onExit{
		registrations: make(map[OnExitID]registration),
	}
}

type registration struct {
	id       int32
	priority priority
	f        func()
}

type byPriorityLIFO []registration

func (a byPriorityLIFO) Len() int           { return len(a) }
func (a byPriorityLIFO) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriorityLIFO) Less(i, j int) bool { return a[i].priority > a[j].priority || a[i].id > a[j].id }

type onExit struct {
	counter       OnExitID
	registrations map[OnExitID]registration
}

func (m *onExit) Register(f func()) OnExitID {
	return m.register(f, lowPriority)
}

func (m *onExit) register(f func(), p priority) OnExitID {
	id := OnExitID{
		id: atomic.AddInt32(&m.counter.id, 1),
	}
	m.registrations[id] = registration{
		id:       id.id,
		priority: p,
		f:        f,
	}
	return id
}

func (m *onExit) Unregister(id OnExitID) {
	delete(m.registrations, id)
}

func (m *onExit) run() {
	// sort registrations first by priority and then in reverse order of registration (LIFO)
	sortable := make([]registration, 0, len(m.registrations))
	for id := range m.registrations {
		sortable = append(sortable, m.registrations[id])
	}
	sort.Sort(byPriorityLIFO(sortable))

	for _, registration := range sortable {
		// invoke all registered onExit functions. Invoke in a manner that recovers from panics so that a
		// function panicking does not prevent the other functions from running.
		invokeAndRecover(registration.f)
	}
}

// Invokes the provided function. If the function panics, recovers and prints to console.
func invokeAndRecover(f func()) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in invokeAndRecover:", r)
		}
	}()
	f()
}
