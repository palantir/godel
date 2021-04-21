Summary
-------
`./godelw generate` runs "go generate" tasks in a project based on configuration. This task is not a built-in task, so
we will add it as a plugin task.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository and Go module
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
* Go files have license headers

Add the `go generate` plugin
----------------------------
The [gödel generate plugin](https://github.com/palantir/godel-generate-plugin) is a gödel plugin that allows
`go generate` tasks to be defined, run and verified. The plugin identifier is
"com.palantir.godel-generate-plugin:generate-plugin:1.2.0", and it is available on GitHub.

Add the plugin definition to `godel/config/godel.yml`:

```
➜ echo 'plugins:
  resolvers:
    - "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
  plugins:
    - locator:
        id: "com.palantir.godel-generate-plugin:generate-plugin:1.2.0"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
  paths:
    - "godel"' > godel/config/godel.yml
```

Any `./godelw` invocation resolves all plugins and assets, downloading them if needed. Running `./godelw` will download
the plugin:

```
➜ ./godelw
Getting package from https://github.com/palantir/godel-generate-plugin/releases/download/v1.2.0/generate-plugin-1.2.0-linux-amd64.tgz...

188.29 KiB / 3.65 MiB [--->____________________________________________________________] 5.03% ? p/s
637.83 KiB / 3.65 MiB [---------->____________________________________________________] 17.05% ? p/s
2.10 MiB / 3.65 MiB [------------------------------------->___________________________] 57.54% ? p/s
3.65 MiB / 3.65 MiB [---------------------------------------------------------] 100.00% 6.44 MiB p/s
Usage:
  godel [command]

Available Commands:
  artifacts       Print the artifacts for products
  build           Build the executables for products
  check           Run checks (runs all checks if none are specified)
  clean           Remove the build and dist outputs for products
  dist            Create distributions for products
  docker          Create or push Docker images for products
  exec            Executes given shell command using godel
  format          Format files
  generate        Run generate task
  git-hooks       Install git commit hooks that verify that Go files are properly formatted before commit
  github-wiki     Push contents of a documents directory to a GitHub Wiki repository
  goland          GoLand project commands
  idea            Create IntelliJ project files for this project
  info            Print information regarding gödel
  install         Install gödel from a local tgz file
  license         Run license task
  packages        Lists all of the packages in the project except those excluded by configuration
  products        Print the IDs of the products in this project
  project-version Print the version of the project
  publish         Publish products
  run             Run product
  run-check       Runs a specific check
  tasks-config    Prints the full YAML configuration used to load tasks and assets
  test            Test packages
  test-tags       Print the test packages that match the provided test tags
  update          Update gödel for project
  upgrade-config  Upgrade configuration
  verify          Run verify tasks for project
  version         Print godel version

Flags:
      --debug            run in debug mode (print full stack traces on failures and include other debugging output)
  -h, --help             help for godel
      --version          print godel version
      --wrapper string   path to the wrapper script for this invocation

Use "godel [command] --help" for more information about a command.
```

All projects share the same plugin and asset files, so once a plugin or asset has been resolved, subsequent runs will
not re-download them.

Define `go generate` tasks
--------------------------
We will extend `echgo` by creating some different echo implementations. The different types will be defined as enums,
and we will use `go generate` to invoke `stringer` to create the string representation of these enum values.

Run the following to update the `echo` implementation:

```
➜ echo '// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package echo

import (
	"fmt"
	"strings"
)

type Type int

func (t Type) String() string {
	switch t {
	case Simple:
		return "Simple"
	case Reverse:
		return "Reverse"
	default:
		panic(fmt.Sprintf("unrecognized type: %d", t))
	}
}

const (
	Simple Type = iota
	Reverse
	end
)

var echoers = []Echoer{
	Simple:  &simpleEchoer{},
	Reverse: &reverseEchoer{},
}

func NewEchoer(typ Type) Echoer {
	return echoers[typ]
}

func TypeFrom(typ string) (Type, error) {
	for curr := Simple; curr < end; curr++ {
		if strings.ToLower(typ) == strings.ToLower(curr.String()) {
			return curr, nil
		}
	}
	return end, fmt.Errorf("unrecognized type: %s", typ)
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}

type reverseEchoer struct{}

func (e *reverseEchoer) Echo(in string) string {
	out := make([]byte, len(in))
	for i := 0; i < len(out); i++ {
		out[i] = in[len(in)-1-i]
	}
	return string(out)
}' > echo/echo.go
➜ SRC='// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package echo_test

import (
	"testing"

	"PROJECT_PATH/echo"
)

func TestEcho(t *testing.T) {
	echoer := echo.NewEchoer(echo.Simple)
	for i, tc := range []struct {
		in   string
		want string
	}{
		{"foo", "foo"},
		{"foo bar", "foo bar"},
	} {
		if got := echoer.Echo(tc.in); got != tc.want {
			t.Errorf("case %d failed: want %q, got %q", i, tc.want, got)
		}
	}
}' && SRC=${SRC//PROJECT_PATH/$PROJECT_PATH} && echo "$SRC" > echo/echo_test.go
➜ SRC='// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"flag"
	"fmt"
	"strings"

	"PROJECT_PATH/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	typeVar := flag.String("type", echo.Simple.String(), "type of echo")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo2 version:", version)
		return
	}
	typ, err := echo.TypeFrom(*typeVar)
	if err != nil {
		fmt.Println("invalid echo type:", *typeVar)
		return
	}
	echoer := echo.NewEchoer(typ)
	fmt.Println(echoer.Echo(strings.Join(flag.Args(), " ")))
}' && SRC=${SRC//PROJECT_PATH/$PROJECT_PATH} && echo "$SRC" > main.go
```

At a high level, this code introduces a new type named `Type` that represents the different types of echo
implementations. The code maintains a mapping from the types to the implementations and provides a function that returns
the Type for a given string. Run the code to verify that this works for the "simple" and "reverse" types that were
defined:

```
➜ go run main.go -type simple foo
foo
➜ go run main.go -type reverse foo
oof
```

The code relies on `Type` having a `String` function that returns its string representation. The current implementation
works, but it is a bit redundant since the string value is always the name of the constant. It is also a maintenance
burden: whenever a new type is added or an existing type is renamed, the `String` function must also be updated.
Furthermore, because the string definitions are simply part of the switch statement, if someone forgets to add or update
the definitions, this will not be caught at compile-time, so it's also a likely source of future bugs.

We can address this by using `go generate` and the `stringer` tool to generate this code automatically.

Run the following to ensure that you have the `stringer` tool:

```
➜ go install golang.org/x/tools/cmd/stringer@latest
go: downloading golang.org/x/tools v0.1.0
go: downloading golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4
go: downloading golang.org/x/xerrors v0.0.0-20200804184101-5ec99f83aff1
go: downloading golang.org/x/mod v0.3.0
```

Now, update `echo.go` to have a `go generate` line that invokes `stringer`:

```
➜ echo '// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate stringer -type=Type

package echo

import (
	"fmt"
	"strings"
)

type Type int

func (t Type) String() string {
	switch t {
	case Simple:
		return "Simple"
	case Reverse:
		return "Reverse"
	default:
		panic(fmt.Sprintf("unrecognized type: %v", t))
	}
}

const (
	Simple Type = iota
	Reverse
	end
)

var echoers = []Echoer{
	Simple:  &simpleEchoer{},
	Reverse: &reverseEchoer{},
}

func NewEchoer(typ Type) Echoer {
	return echoers[typ]
}

func TypeFrom(typ string) (Type, error) {
	for curr := Simple; curr < end; curr++ {
		if strings.ToLower(typ) == strings.ToLower(curr.String()) {
			return curr, nil
		}
	}
	return end, fmt.Errorf("unrecognized type: %s", typ)
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}

type reverseEchoer struct{}

func (e *reverseEchoer) Echo(in string) string {
	out := make([]byte, len(in))
	for i := 0; i < len(out); i++ {
		out[i] = in[len(in)-1-i]
	}
	return string(out)
}' > echo/echo.go
```

Now that the `//go:generate` directive exists, the standard Go approach would be to run `go generate` to run the
generation task. However, this approach depends on developers knowing/remembering to run `go generate` when they update
the definitions. Projects typically address this my noting it in their documentation or in comments, but this is
obviously quite fragile, and correctly calling all of the required generators can be especially challenging for larger
projects that may have several `go generate` tasks.

We can address this by defining the `generate` tasks as part of the declarative configuration for our project. Define a
"generate" task in `godel/config/generate-plugin.yml` by running the following:

```
➜ echo 'generators:
  stringer:
    go-generate-dir: echo
    gen-paths:
      paths:
        - echo/type_string.go' > godel/config/generate-plugin.yml
```

This specifies that we have a generator task named "stringer" (this name is specified by the user and can be anything).
The `go-generate-dir` specifies the directory (relative to the project root) in which `go generate` should be run. The
`gen-paths` parameter specifies paths to the files or directories that are generated or modified by the `go generate`
task.

Run the generator task and verify that it generates the expected code:

```
➜ ./godelw generate
➜ cat ./echo/type_string.go
// Code generated by "stringer -type=Type"; DO NOT EDIT.

package echo

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Simple-0]
	_ = x[Reverse-1]
	_ = x[end-2]
}

const _Type_name = "SimpleReverseend"

var _Type_index = [...]uint8{0, 6, 13, 16}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return "Type(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
```

We can see that `echo/type_string.go` was generated and provides an implementation of the `String` function for `Type`.
Now that this exists, we can remove the one we wrote manually in `echo.go`:

```
➜ echo '// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate stringer -type=Type

package echo

import (
	"fmt"
	"strings"
)

type Type int

const (
	Simple Type = iota
	Reverse
	end
)

var echoers = []Echoer{
	Simple:  &simpleEchoer{},
	Reverse: &reverseEchoer{},
}

func NewEchoer(typ Type) Echoer {
	return echoers[typ]
}

func TypeFrom(typ string) (Type, error) {
	for curr := Simple; curr < end; curr++ {
		if strings.ToLower(typ) == strings.ToLower(curr.String()) {
			return curr, nil
		}
	}
	return end, fmt.Errorf("unrecognized type: %s", typ)
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}

type reverseEchoer struct{}

func (e *reverseEchoer) Echo(in string) string {
	out := make([]byte, len(in))
	for i := 0; i < len(out); i++ {
		out[i] = in[len(in)-1-i]
	}
	return string(out)
}' > echo/echo.go
```

With this setup, `./godelw generate` can be called on a project to invoke all of its `generate` tasks.

We will now attempt to commit these changes. If you have followed the tutorial up to this point, the git hook that
enforces formatting for files will reject the commit:

```
➜ git add echo godel main.go
➜ git commit -m "Add support for echo types"
Unformatted files exist -- run ./godelw format to format these files:
  echo/type_string.go
```

This is because the generated Go file does not match the formatting enforced by `ptimports`. However, because this code
is generated, we do not want to modify it after the fact. In general, we want to simply exclude generated code from all
gödel tasks -- we don't want to add license headers to it, format it, run linting checks on it, etc. We will update the
exclude block of `godel/config/godel.yml` to reflect this and specify that the file should be ignored:

```
➜ echo 'plugins:
  resolvers:
    - "https://github.com/{{index GroupParts 1}}/{{index GroupParts 2}}/releases/download/v{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
  plugins:
    - locator:
        id: "com.palantir.godel-generate-plugin:generate-plugin:1.2.0"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
  paths:
    - "godel"
    - "echo/type_string.go"' > godel/config/godel.yml
```

We will go through this file in more detail in the next portion of the tutorial, but for now it is sufficient to know
that this excludes the `echo/type_string.go` file from checks and other tasks (we will make this more generic later).
We can now commit the changes:

```
➜ git add echo godel main.go
➜ git commit -m "Add support for echo types"
[master 232fd7a] Add support for echo types
 6 files changed, 87 insertions(+), 4 deletions(-)
 create mode 100644 echo/type_string.go
 create mode 100644 godel/config/generate-plugin.yml
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository and Go module
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function

Tutorial next step
------------------
[Define excludes](https://github.com/palantir/godel/wiki/Exclude)

More
----
### Verification
The `generate` task also supports a verification mode that ensures that the code that is generated by the
`./godelw generate` task does not change the contents of the generated target paths. This is useful for use in CI to
verify that developers properly ran `generate`.

To demonstrate this, update `echo/echo.go` to add another echo type:

```
➜ echo '// Copyright (c) 2021 Author Name. All rights reserved.
// Use of this source code is governed by the Apache License, Version 2.0
// that can be found in the LICENSE file.

//go:generate stringer -type=Type

package echo

import (
	"fmt"
	"math/rand"
	"strings"
)

type Type int

const (
	Simple Type = iota
	Reverse
	Random
	end
)

var echoers = []Echoer{
	Simple:  &simpleEchoer{},
	Reverse: &reverseEchoer{},
	Random:  &randomEchoer{},
}

func NewEchoer(typ Type) Echoer {
	return echoers[typ]
}

func TypeFrom(typ string) (Type, error) {
	for curr := Simple; curr < end; curr++ {
		if strings.ToLower(typ) == strings.ToLower(curr.String()) {
			return curr, nil
		}
	}
	return end, fmt.Errorf("unrecognized type: %s", typ)
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}

type reverseEchoer struct{}

func (e *reverseEchoer) Echo(in string) string {
	out := make([]byte, len(in))
	for i := 0; i < len(out); i++ {
		out[i] = in[len(in)-1-i]
	}
	return string(out)
}

type randomEchoer struct{}

func (e *randomEchoer) Echo(in string) string {
	inBytes := []byte(in)
	out := make([]byte, len(in))
	for i := 0; i < len(out); i++ {
		randIdx := rand.Intn(len(inBytes))
		out[i] = inBytes[randIdx]
		inBytes = append(inBytes[:randIdx], inBytes[randIdx+1:]...)
	}
	return string(out)
}' > echo/echo.go
```

Now, run the generate task with the `--verify` flag:

```
➜ ./godelw generate --verify
Generators produced output that differed from what already exists: [stringer]
  stringer:
    echo/type_string.go: previously had checksum 12aa70bee990397e6a8f08cdd4b5524f86d1f50cc5c33294e14e0aba9b5208f7, now has checksum 648a980df503433b06c37b23b0e4e4b2a2c7b8686beeade8baaddde335e2d01f
```

As you can see, the task determined that the `generate` task changed a file that was specified in the `gen-paths`
configuration for the task and prints a warning that specifies this.

The `gen-paths` configuration consists of a `paths` and `names` list that specify regular expressions that match the
paths or names of files created by the generator. When a `generate` task is run in `--verify` mode, the task determines
all of the paths in the project that match the `gen-paths` configuration, computes the checksums of all of the files,
runs the generation task, computes all of the paths that match after the task is run, and then compares both the file
list and the checksums. If either the file list or checksums differ, the differences are echoed and the verification
fails. Note that the `--verify` mode still runs the `generate` task -- because `go generate` itself doesn't have a
notion of a dry run or verification, the only way to determine the effects of a `generate` task is to run it and to
compare the state before and after the run. Although this same kind of check can be approximated by something like a
`git status`, this configuration mechanism is more explicit/declarative and more robust.

Revert the changes by running the following:

```
➜ git checkout -- echo
```

### Configuring multiple generate tasks
The `generate` configuration supports configuring generate tasks organized in any manner. However, for projects that
have multiple different kinds of `go generate` tasks, we have found that the following can be an effective way to
organize the generators:

* Have a `generators` directory at the top level of the project
* Have a subdirectory for each generator type within the `generators` directory
* The subdirectory contains a `generate.go` file with `go generate` directive and a `vendor` directory that vendors the
  program being called by the generator

Here is an example of such a directory structure:

```
➜ tree generators
generators [error opening dir]

0 directories, 0 files
```

### Example generators
The following are examples of tasks that can be set up and configured to run as a `go generate` task:

* [`protoc`](https://github.com/golang/protobuf) for generating protobufs
* [`mockery`](https://github.com/vektra/mockery) for generating mocks for testing
* [`stringer`](https://godoc.org/golang.org/x/tools/cmd/stringer) for generating String() functions for types
* [`gopherjs`](https://github.com/gopherjs/gopherjs/) to generate Javascript from Go code
* [`go-bindata`](https://github.com/jteeuwen/go-bindata) to generate Go code that embeds static resources
* [`amalgomate`](https://github.com/palantir/amalgomate) for creating libraries out of Go `main` packages

Generators that leverage all of these tools (as well as other custom/proprietary internal tools) have been successfully
configured and integrated with gödel projects.

Generators work most effectively if there is a guarantee that, for a given input, the generated outputs is always the
same (this is required if output files are verified).

### Consider writing plugins for common generators
Generators can be very helpful, and their flexibility makes it appealing to use them to generate output for projects,
especially on a one-off basis. However, if you find that you are commonly using the same generator across projects,
consider creating a gödel plugin for the same task instead -- although generators offer nice flexibility and the
`generate` plugin offers a way for such tasks to be integrated into a gödel project, gödel plugins are much easier to
work with. As an example, implementing a plugin for `stringer` would be easy to do, and using that plugin would prevent
having to vendor the `stringer` program.
