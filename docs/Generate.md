Summary
-------
`./godelw generate` runs "go generate" tasks in a project based on configuration.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers

([Link](https://github.com/nmiyake/echgo/tree/0239b282904d05bb9eef6c3c3edfe1c28f888ad3))

Define `go generate` tasks
--------------------------

We will extend `echgo` by creating some different echo implementations. The different types will be defined as enums,
and we will use `go generate` to invoke `stringer` to create the string representation of these enum values.

Run the following to update the `echo` implementation:

```
➜ echo '// Copyright (c) 2017 Author Name
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
➜ echo '// Copyright (c) 2017 Author Name
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

package echo_test

import (
	"testing"

	"github.com/nmiyake/echgo/echo"
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
}' > echo/echo_test.go
➜ echo '// Copyright (c) 2017 Author Name
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

package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nmiyake/echgo/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	typeVar := flag.String("type", echo.Simple.String(), "type of echo")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo version:", version)
		return
	}
	typ, err := echo.TypeFrom(*typeVar)
	if err != nil {
		fmt.Println("invalid echo type:", *typeVar)
		return
	}
	echoer := echo.NewEchoer(typ)
	fmt.Println(echoer.Echo(strings.Join(flag.Args(), " ")))
}' > main.go
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
➜ go get -u golang.org/x/tools/cmd/stringer
```

Now, update `echo.go` to have a `go generate` line that invokes `stringer`:

```
➜ echo '// Copyright (c) 2017 Author Name
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
"generate" task in `godel/config/generate.yml` by running the following:

```
➜ echo 'generators:
  stringer:
    go-generate-dir: echo
    gen-paths:
      paths:
        - echo/type_string.go' > godel/config/generate.yml
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

import "fmt"

const _Type_name = "SimpleReverseend"

var _Type_index = [...]uint8{0, 6, 13, 16}

func (i Type) String() string {
	if i < 0 || i >= Type(len(_Type_index)-1) {
		return fmt.Sprintf("Type(%d)", i)
	}
	return _Type_name[_Type_index[i]:_Type_index[i+1]]
}
```

We can see that `echo/type_string.go` was generated and provides an implementation of the `String` function for `Type`.
Now that this exists, we can remove the one we wrote manually in `echo.go`:

```
➜ echo '// Copyright (c) 2017 Author Name
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
  /Volumes/git/go/src/github.com/nmiyake/echgo/echo/type_string.go
```

This is because the generated Go file does not match the formatting enforced by `ptimports`. However, because this code
is generated, we do not want to modify it after the fact. In general, we want to simply exclude generated code from all
gödel tasks -- we don't want to add license headers to it, format it, run linting checks on it, etc. We will update the
`godel/config/exclude.yml` to reflect this and specify that the file should be ignored:

```
➜ echo 'names:
  - "\\\..+"
  - "vendor"
paths:
  - "godel"
  - "echo/type_string.go"' > godel/config/exclude.yml
```

We will go through this file in more detail in the next portion of the tutorial, but for now it is sufficient to know
that this excludes the `echo/type_string.go` file from checks and other tasks (we will make this more generic later).
We can now commit the changes:

```
➜ git add echo godel main.go
➜ git commit -m "Add support for echo types"
[master 4e528d0] Add support for echo types
 6 files changed, 72 insertions(+), 4 deletions(-)
 create mode 100644 echo/type_string.go
```

Many Go projects would consider this sufficient -- they would document the requirement that developers must run
`go get golang.org/x/tools/cmd/stringer` locally to in order to run "generate" and also ensure that this same action is
performed in their CI environment. However, this introduces an external dependency on the ability to get and install
`stringer`. Furthermore, the version of `stringer` is not defined/locked in anywhere -- the `go get` action will fetch
whatever version is the latest at that time. This may not be an issue for tools that have a completely mature API, but
if there are behavior changes between versions of the tools it can lead to the generation tasks creating inconsistent
output.

For that reason, if the `generate` task is running a Go program, we have found it helpful to vendor the entire program
within the project and to run it using `go run` to ensure that the `generate` task does not have any external
dependencies. We will use this construction for this project.

Run the following to create a new directory for the generator and create a `vendor` directory within that directory:

```
➜ mkdir -p generator/vendor
```

Putting the `vendor` directory within `generator` ensures that the code we vendor will only be accessible within the
`generator` directory. We will now vendor the `stringer` program. In a real workflow, you would use the vendoring tool
of your choice to do so. For the purposes of this tutorial, we will handle our vendoring manually by copying the code we
need to the expected location:

```
➜ mkdir -p generator/vendor/golang.org/x/tools/cmd/stringer
➜ cp $(find $GOPATH/src/golang.org/x/tools/cmd/stringer -name '*.go' -not -name '*_test.go' -maxdepth 1 -type f) generator/vendor/golang.org/x/tools/cmd/stringer/
```

Note: the `find` command above performs some pruning to copy only the buildable Go files that will be used -- if you
don't care about pulling in extra unneeded files (such as tests and testdata files), you can run
`cp -r $GOPATH/src/golang.org/x/tools/cmd/stringer/* generator/vendor/golang.org/x/tools/cmd/stringer/` instead.

We will now define a generator that invokes this:

```
➜ echo '// Copyright (c) 2017 Author Name
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

//go:generate -command runstringer go run vendor/golang.org/x/tools/cmd/stringer/stringer.go vendor/golang.org/x/tools/cmd/stringer/importer18.go

//go:generate runstringer -type=Type ../echo

package generator' > generator/generate.go
```

This generator now runs `stringer` directly from the vendor directory. If we had other packages on which we wanted to
invoke `stringer`, we could simply update this file to do so.

Update the previous code to remove its generation logic:

```
➜ echo '// Copyright (c) 2017 Author Name
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

Update the `generate.yml` configuration:

```
➜ echo 'generators:
  stringer:
    go-generate-dir: generator
    gen-paths:
      paths:
        - echo/type_string.go' > godel/config/generate.yml
```

Run the `generate` task to verify that it still succeeds:

```
➜ ./godelw generate
```

Run the `check` command to verify that the project is still valid:

```
➜ ./godelw check
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running golint...
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
golang.org/x/tools
Running outparamcheck...
Running unconvert...
Running varcheck...
Checks produced output: [novendor]
```

You can see that the `novendor` check now fails. This is because the `golang.org/x/tools` package is present in the
`vendor` directory, but no packages in the project are importing its packages and the `novendor` check has identified
it as an unused vendored project. However, in this instance we know that this is valid because we call the code directly
from `go generate`. Update the `godel/config/check.yml` configuration to reflect this:

```
➜ echo 'checks:
  golint:
    filters:
      - value: "should have comment or be unexported"
      - value: "or a comment on this block"
  novendor:
    args:
      # ignore packages added for generation
      - "--ignore"
      - "./generator/vendor/golang.org/x/tools"' > godel/config/check.yml
```

Run `check` again to verify that the checks now pass:

```
➜ ./godelw check
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running golint...
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
Running outparamcheck...
Running unconvert...
Running varcheck...
```

Commit these changes by running the following:

```
➜ git add echo generator godel
➜ git commit -m "Update generator code"
[master 08752b2] Update generator code
 8 files changed, 693 insertions(+), 4 deletions(-)
 create mode 100644 generator/generate.go
 create mode 100644 generator/vendor/golang.org/x/tools/cmd/stringer/importer18.go
 create mode 100644 generator/vendor/golang.org/x/tools/cmd/stringer/importer19.go
 create mode 100644 generator/vendor/golang.org/x/tools/cmd/stringer/stringer.go
```

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function

([Link](https://github.com/nmiyake/echgo/tree/08752b2ae998c14dd5abb789cebc8f5848f7cf4e))

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
➜ echo '// Copyright (c) 2017 Author Name
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
    echo/type_string.go: previously had checksum d594017fce62ad2e2a8a98f9c7d519012d1df157c0f59088aaea2702a24f70e0, now has checksum 5b57686f254b93087a006aa6ab65753356f659015e7dd0c7e3053ea9fc2c024f
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
generators
├── mocks
│   ├── generate.go
│   └── vendor
│       └── ...
└── stringer
    ├── generate.go
    └── vendor
        └── ...
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
