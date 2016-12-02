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
