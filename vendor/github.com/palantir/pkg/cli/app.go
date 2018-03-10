// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package cli

import (
	"context"
	"io"
	"os"

	"github.com/palantir/pkg/cli/completion"
)

type App struct {
	Command
	Before         func(ctx Context) error
	ErrorHandler   func(ctx Context, err error) int
	Version        string
	Stdout         io.Writer
	Stderr         io.Writer
	Completion     map[string]completion.Provider
	Manpage        *Manpage
	Backcompat     []Backcompat
	OnExit         OnExit
	ContextConfig  func(Context, context.Context) context.Context
	ContextOptions []ContextOption
}

type Option func(*App)

type ContextOption func(*Context)

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
