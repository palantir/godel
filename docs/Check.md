Summary
-------
`./godelw check` runs a set of static analysis checks on the project Go files.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`

([Link](https://github.com/nmiyake/echgo/tree/24f63f727542c7189c82f04f7e2a4aa38c090137))

Run checks
----------

When writing Go code, it can be useful to check code for errors and consistency issues using static code analysis.

The current echo program simply echoes the user's input exactly. We will extend the program to allow different types of
echoes to be generated. As a first step to doing this, we will define an `Echoer` interface that defines an `Echo`
function and refactor the current echo functionality to be a simple echoer that implements this interface.

Run the following to update the program to perform this refactor:

```
➜ echo 'package echo

type Echoer interface {
	Echo(in string) string
}' > echo/echoer.go
➜ echo 'package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (_ *simpleEchoer) Echo(in string) string {
	return in
}' > echo/echo.go
➜ echo 'package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nmiyake/echgo/echo"
)

func main() {
	echoer := echo.NewEchoer()
	fmt.Println(echoer.Echo(strings.Join(os.Args[1:], " ")))
}' > main.go
```

These files are formatted correctly and form a fully functioning program. Run `./godelw check` to run static code checks
on the project:

```
➜ ./godelw check
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running golint...
echo/echo.go:9:1: receiver name should not be an underscore
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
Running outparamcheck...
Running unconvert...
Running varcheck...
Checks produced output: [golint]
```

The output indicates that there was an issue identified by the `golint` check. Fix the issue by updating the receiver
name:

```
➜ echo 'package echo

func NewEchoer() Echoer {
	return &simpleEchoer{}
}

type simpleEchoer struct{}

func (e *simpleEchoer) Echo(in string) string {
	return in
}' > echo/echo.go
```

Run `./godelw check` again to verify that the issue has been resolved:

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

Commit the changes to the repository:

```
➜ git add main.go echo
➜ git commit -m "Add echoer interface"
[master 0a64992] Add echoer interface
 3 files changed, 15 insertions(+), 3 deletions(-)
 create mode 100644 echo/echoer.go
➜ git status
On branch master
nothing to commit, working directory clean
```

Refer to the "More" sections below for examples of configuring the checks in different ways.

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go` and `echo/echoer.go`

([Link](https://github.com/nmiyake/echgo/tree/0a649925e317b7896e537ef23a4885062a3ec9fb))

Tutorial next step
------------------

[Run tests](https://github.com/palantir/godel/wiki/Test)

More
----

### Suppress check issues based on output

In some instances, it may be desirable to suppress certain issues flagged by checks. As an example, modify
`echo/echoer.go` as follows:

```
➜ echo 'package echo

// Echoes the input.
type Echoer interface {
	Echo(in string) string
}' > echo/echoer.go
```

Running `./godelw check` flags the following:

```
➜ ./godelw check
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running golint...
echo/echoer.go:3:1: comment on exported type Echoer should be of the form "Echoer ..." (with optional leading article)
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
Running outparamcheck...
Running unconvert...
Running varcheck...
Checks produced output: [golint]
```

Although this is a valid check performed by `go lint`, not all projects conform exactly with the Go style for comments.
In some cases, it makes sense to disable specific checks like this. This can be done by updating the
`godel/config/check.yml` file to configure the `check` command to ignore all output from the `golint` check that
contains `comment on exported type \w should be of the form` in its message.

The default configuration for `godel/config/check.yml` is as follows:

```
➜ cat godel/config/check.yml
checks:
  golint:
    filters:
      - value: "should have comment or be unexported"
      - value: "or a comment on this block"
```

Add the line `- value: "comment on exported type [[:word:]]+ should be of the form"` to this configuration:

```
➜ echo 'checks:
  golint:
    filters:
      - value: "should have comment or be unexported"
      - value: "or a comment on this block"
      - value: "comment on exported type [[:word:]]+ should be of the form"' > godel/config/check.yml
```

Re-run `./godelw check` with the updated configuration to verify that lines that match this output are no longer
reported:

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

Revert the local changes by running the following:

```
➜ git checkout -- echo godel
➜ git status
On branch master
nothing to commit, working directory clean
```

Filters have a `type` and a `value`. When `type` is not specified (as in the examples above), it defaults to `message`,
which means that the value is matched against the message of the output. The `type` field for filters can also be `name`
or `path`. `name` matches files based on their name, while `path` matches based on an exact relative path.

For example, the following configuration will ignore all issues reported by `errcheck` for `main.go`:

```yaml
checks:
  errcheck:
    filters:
      - type: "path"
        value: "main.go"
```

Because the `type` above is `path`, this configuration would ignore `errcheck` issues in `./main.go`. However, issues in
other files named `main.go` in the project (for example, `./subproject/main.go`) would still be reported. Setting the
`type` to `name` would change the behavior so that issues in all files named `main.go` would be ignored.

The match values use Go regular expressions to perform matches. For example, the following configuration ignores all
`golint` issues reported for any files that have the extension `.pb.go`:

```yaml
checks:
  golint:
    filters:
      - type: "name"
        value: ".*.pb.go"
```

### Disable checks

Checks can be disabled completely for the entire project by setting the `skip` field to `true`.

For example, the following configuration will disable the `golint` check for the project:

```
➜ echo 'checks:
  golint:
    skip: true' > godel/config/check.yml
```

Run `./godelw check` with the updated configuration to verify that the `golint` check is no longer run:

```
➜ ./godelw check
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
Running outparamcheck...
Running unconvert...
Running varcheck...
```

Revert the local changes by running the following:

```
➜ git checkout -- godel
➜ git status
On branch master
nothing to commit, working directory clean
```

### Configure check arguments

Many of the tools used by `check` accept command-line arguments. The arguments that are passed to the tool can be
specified using the `args` parameter. The elements of the `args` parameter are provided to the underlying check. For
example, the following configuration configures the `errcheck` check to be run with the arguments
`-ignore 'io/ioutil:ReadFile'`:

```yaml
checks:
  errcheck:
    args:
      - "-ignore"
      - "io/ioutil:ReadFile"
```

### Run individual checks

Individual checks can be run in isolation by specifying the name of the check as an argument to `check`. This can be
useful when iterating on code in an attempt to fix an issue flagged by a specific check.

For example, the following runs only `govet`:

```
➜ ./godelw check govet
Running govet...
```

### Constituent checks

The following checks are run as part of `check`:

* [`compiles`](https://github.com/palantir/checks/tree/master/compiles) verifies that all of the Go code in the project
  compiles, including code in test files (which is not checked by `go build`)
* [`deadcode`](https://github.com/tsenart/deadcode) finds unused code
* [`errcheck`](https://github.com/kisielk/errcheck) ensures that returned errors are checked
* [`extimport`](https://github.com/palantir/checks/tree/master/extimport) verifies that all non-standard library
  packages that are imported by the project are present in a vendor directory within the project
* [`govet`](https://github.com/nmiyake/govet) runs [`go vet`](https://golang.org/cmd/vet/)
* [`importalias`](https://github.com/palantir/checks/tree/master/importalias) ensures that, if an import path in the
  package is imported using an alias, then all imports in the project that assign an alias for that path use the same
  alias
* [`ineffassign`](https://github.com/gordonklaus/ineffassign) flags ineffectual assignment statements
* [`nobadfuncs`](https://github.com/palantir/checks/tree/master/nobadfuncs) allows a project to blacklist specific
  functions (for example, `fmt.Println`) and flags all uses of the blacklisted functions unless the use is specifically
  whitelisted (see project documentation for details)
* [`novendor`](https://github.com/palantir/checks/tree/master/novendor) flags projects that exist in the `vendor`
  directory but are not used by the project
* [`outparamcheck`](https://github.com/palantir/checks/tree/master/outparamcheck) checks that functions that are meant
  to take an output parameter defined as an `interface{}` are passed pointers to an object rather than a concrete object
* [`unconvert`](https://github.com/mdempsky/unconvert) flags unnecessary conversions
* [`varcheck`](https://github.com/opennota/check) checks for unused global variables and constants

One of the core principles of gödel is reproducibility, so the set of available checks (and their specific
implementations) are hard-coded in gödel itself. This means that adding a new check or upgrading the version of a check
is a change that must be made in gödel itself.

If there is a check that you would like to see added to gödel (or believe that the version of an existing check should
be updated), please file an issue on the project.

Currently, the set of checks that are run are built into gödel itself. There are plans to make the checks that are run
pluggable so that they can be customized based on the needs/desires of specific projects.
