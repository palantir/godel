// Copyright 2016 Palantir Technologies. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package cli is a library for writing CLI applications. It supports subcommands, flags, command completion,
before and after hooks, exit code control and more.

Here is an example that creates prints "Hello, world!\n" when invoked:

	func main() {
		app := cli.NewApp()
		app.Action = func(ctx cli.Context) error {
			ctx.Printf("Hello, world!\n")
			return nil
		}
		os.Exit(app.Run(os.Args))
	}
*/
package cli
