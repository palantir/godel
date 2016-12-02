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
	"io"
	"os"

	"github.com/palantir/pkg/cli/completion"
)

type App struct {
	Command
	Before       func(ctx Context) error
	ErrorHandler func(ctx Context, err error) int
	Version      string
	Stdout       io.Writer
	Stderr       io.Writer
	Completion   map[string]completion.Provider
	Manpage      *Manpage
	Backcompat   []Backcompat
	OnExit       OnExit
}

type Option func(*App)

func NewApp(opts ...Option) *App {
	app := &App{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		OnExit: newOnExit(),
	}
	for _, opt := range opts {
		opt(app)
	}
	return app
}
