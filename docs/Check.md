Summary
-------
`./godelw check` runs a set of static analysis checks on the project Go files.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`

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
➜ SRC='package main

import (
	"fmt"
	"os"
	"strings"

	"PROJECT_PATH/echo"
)

func main() {
	echoer := echo.NewEchoer()
	fmt.Println(echoer.Echo(strings.Join(os.Args[1:], " ")))
}' && SRC=${SRC//PROJECT_PATH/$PROJECT_PATH} && echo "$SRC" > main.go
```

These files are formatted correctly and form a fully functioning program. Run `./godelw check` to run static code checks
on the project:

```
➜ ./godelw check
[deadcode]      Running deadcode...
[extimport]     Running extimport...
[compiles]      Running compiles...
[errcheck]      Running errcheck...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        echo/echo.go:9:1: receiver name should not be an underscore, omit the name if it is unused
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[compiles]      Finished compiles
[unconvert]     Running unconvert...
[errcheck]      Finished errcheck
[deadcode]      Finished deadcode
[varcheck]      Running varcheck...
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Check(s) produced output: [golint]
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
[errcheck]      Running errcheck...
[extimport]     Running extimport...
[compiles]      Running compiles...
[deadcode]      Running deadcode...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[errcheck]      Finished errcheck
[compiles]      Finished compiles
[varcheck]      Running varcheck...
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
```

Commit the changes to the repository:

```
➜ git add main.go echo
➜ git commit -m "Add echoer interface"
[master 626c0ec] Add echoer interface
 3 files changed, 14 insertions(+), 2 deletions(-)
 create mode 100644 echo/echoer.go
```

Refer to the "More" sections below for examples of configuring the checks in different ways.

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go` and `echo/echoer.go`

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
[errcheck]      Running errcheck...
[extimport]     Running extimport...
[deadcode]      Running deadcode...
[compiles]      Running compiles...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        echo/echoer.go:3:1: comment on exported type Echoer should be of the form "Echoer ..." (with optional leading article)
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[compiles]      Finished compiles
[unconvert]     Running unconvert...
[deadcode]      Finished deadcode
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[varcheck]      Finished varcheck
[unconvert]     Finished unconvert
Check(s) produced output: [golint]
```

Although this is a valid check performed by `go lint`, not all projects conform exactly with the Go style for comments.
In some cases, it makes sense to disable specific checks like this. This can be done by updating the
`godel/config/check.yml` file to configure the `check` command to ignore all output from the `golint` check that
contains `comment on exported type \w should be of the form` in its message.

The default configuration for `godel/config/check-plugin.yml` is as follows:

```
➜ cat godel/config/check-plugin.yml
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
      - value: "comment on exported type [[:word:]]+ should be of the form"' > godel/config/check-plugin.yml
```

Re-run `./godelw check` with the updated configuration to verify that lines that match this output are no longer
reported:

```
➜ ./godelw check
[compiles]      Running compiles...
[errcheck]      Running errcheck...
[deadcode]      Running deadcode...
[extimport]     Running extimport...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[compiles]      Finished compiles
[unconvert]     Running unconvert...
[deadcode]      Finished deadcode
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[varcheck]      Finished varcheck
[unconvert]     Finished unconvert
```

Revert the local changes by running the following:

```
➜ git checkout -- echo godel
```

Filters have a `type` and a `value`. When `type` is not specified (as in the examples above), it defaults to `message`,
which means that the value is matched against the message of the output. Currently, `message` is the only filter.

### Exclude specific file names or paths from a check
Checks have an optional `exclude` field that can be used to specify names or paths of files that should be excluded from
the check.

For example, the following configuration will ignore all issues reported by `errcheck` for `main.go`:

```yaml
checks:
  errcheck:
    exclude:
      paths:
        - main.go
```

Because the exclude type is `path`, this configuration would ignore `errcheck` issues in `./main.go`. However, issues in
other files named `main.go` in the project (for example, `./subproject/main.go`) would still be reported. Setting the
exclude type to `names` would change the behavior so that issues in all files named `main.go` would be ignored. Both
`names` and `paths` can be specified in the same `exclude` configuration.

The name match values use Go regular expressions to perform matches. For example, the following configuration ignores
all `golint` issues reported for any files that have the extension `.pb.go`:

```yaml
checks:
  golint:
    exclude:
      names:
        - ".*.pb.go"
```

The `checks` configuration also supports specifying `exclude` as a top-level value that applies to all checks (rather
than just an individual one).

### Disable checks
Checks can be disabled completely for the entire project by setting the `skip` field to `true`.

For example, the following configuration will disable the `golint` check for the project:

```
➜ echo 'checks:
  golint:
    skip: true' > godel/config/check-plugin.yml
```

Run `./godelw check` with the updated configuration to verify that the `golint` check is no longer run:

```
➜ ./godelw check
[compiles]      Running compiles...
[extimport]     Running extimport...
[deadcode]      Running deadcode...
[errcheck]      Running errcheck...
[extimport]     Finished extimport
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[compiles]      Finished compiles
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
```

Revert the local changes by running the following:

```
➜ git checkout -- godel
```

### Configure checks
Many checks offer customizable parameters for the checks. Such parameters are specified in the `config` field of the
check's configuration. Refer to the documentation for the asset that provides the check for details on configuring a
check.

### Run individual checks
Individual checks can be run in isolation by specifying the name of the check as an argument to `check`. This can be
useful when iterating on code in an attempt to fix an issue flagged by a specific check.

For example, the following runs only `deadcode` and `govet`:

```
➜ ./godelw check deadcode govet
[govet]    Running govet...
[deadcode] Running deadcode...
[govet]    Finished govet
[deadcode] Finished deadcode
```

### Run an underlying check directly
Running an individual check using `check` runs just that check, but it still runs it through the `check` task, which
uses the logic of the plugin and the asset to determine the arguments passed to the underlying check. However, sometimes
we may want to run the underlying check directly -- possibly to run it on a specific file or package or to specify flags
that are not available through configuration.

The `run-check` task can be used to do this. Running `./godelw run-check [check] [flags] [args]` calls the "underlying"
check directly. It is up to an asset to determine what this means, but most assets wrap a standalone check implemented
as its own CLI, and "running" the check means invoking the CLI. If the underlying check accepts flags, it is safest to
place a `--` after the check so that all of the flags and arguments are passed directly to the underlying check.

For example, `errcheck` can be invoked directly as follows:

```
➜ ./godelw run-check errcheck -- --help
Usage of /root/.godel/assets/com.palantir.godel-okgo-asset-errcheck-errcheck-asset-1.1.1:
  -abspath
    	print absolute paths to files
  -asserts
    	if true, check for ignored type assertion results
  -blank
    	if true, check for errors assigned to blank identifier
  -exclude string
    	Path to a file containing a list of functions to exclude from checking
  -ignore value
    	[deprecated] comma-separated list of pairs of the form pkg:regex
    	            the regex is used to ignore names within pkg. (default "fmt:.*")
  -ignorepkg string
    	comma-separated list of package paths to ignore
  -ignoretests
    	if true, checking of _test.go files is disabled
  -tags value
    	space-separated list of build tags to include
  -verbose
    	produce more verbose logging
```

The help output printed here is that of `errcheck`, and we could have supplied the errcheck flags and arguments directly
in place of `--help`. Also note the use of `--` to indicate that all of the flags/arguments should be passed to the
underlying command. In this instance, if `--` was not used, the `--help` flag would have shown the output for
`run-check errcheck` instead:

```
➜ ./godelw run-check errcheck --help
Usage:
  okgo run-check errcheck [flags]

Flags:
  -h, --help   help for errcheck

Global Flags:
      --assets stringSlice    path(s) to the plugin asset(s)
      --config string         path to the plugin configuration file
      --debug                 run in debug mode
      --godel-config string   path to the godel.yml configuration file
      --project-dir string    path to project directory
```

### Run checks sequentially
Checks can be run sequentially by running with the `--parallel=false` flag:

```
➜ ./godelw check --parallel=false
Running compiles...
Finished compiles
Running deadcode...
Finished deadcode
Running errcheck...
Finished errcheck
Running extimport...
Finished extimport
Running golint...
Finished golint
Running govet...
Finished govet
Running importalias...
Finished importalias
Running ineffassign...
Finished ineffassign
Running novendor...
Finished novendor
Running outparamcheck...
Finished outparamcheck
Running unconvert...
Finished unconvert
Running varcheck...
Finished varcheck
```

### Constituent checks
gödel includes the following checks by default:

* [`compiles`](https://github.com/palantir/go-compiles) verifies that all of the Go code in the project
  compiles, including code in test files (which is not checked by `go build`)
* [`deadcode`](https://github.com/tsenart/deadcode) finds unused code
* [`errcheck`](https://github.com/kisielk/errcheck) ensures that returned errors are checked
* [`extimport`](https://github.com/palantir/go-extimport) verifies that all non-standard library
  packages that are imported by the project are present in a vendor directory within the project
* `govet` runs [`go vet`](https://golang.org/cmd/vet/)
* [`importalias`](https://github.com/palantir/go-importalias) ensures that, if an import path in the package is imported
  using an alias, then all imports in the project that assign an alias for that path use the same alias
* [`ineffassign`](https://github.com/gordonklaus/ineffassign) flags ineffectual assignment statements
* [`novendor`](https://github.com/palantir/go-novendor) flags projects that exist in the `vendor` directory but are not
  used by the project
* [`outparamcheck`](https://github.com/palantir/outparamcheck) checks that functions that are meant to take an output
  parameter defined as an `interface{}` are passed pointers to an object rather than a concrete object
* [`unconvert`](https://github.com/mdempsky/unconvert) flags unnecessary conversions
* [`varcheck`](https://github.com/opennota/check) checks for unused global variables and constants

### Add or remove checks
The checks that are available to the `check` task are determined by the assets provided to the `okgo` plugin. Checks can
be added or removed by modifying the asset configuration for the `okgo` plugin.

For example, we can add the [nobadfuncs check](https://github.com/palantir/go-nobadfuncs) by adding the
[godel-okgo-asset-nobadfuncs asset](https://github.com/palantir/godel-okgo-asset-nobadfuncs). Modify the
`godel/config/godel.yml` file as follows:

```
➜ echo 'default-tasks:
  resolvers:
    - https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz
  tasks:
    com.palantir.okgo:check-plugin:
      assets:
        - locator:
            id: "com.palantir.godel-okgo-asset-nobadfuncs:nobadfuncs-asset:1.0.0"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
  paths:
    - "godel"' > godel/config/godel.yml
```

This adds the asset, which makes it available as a check:

```
➜ ./godelw check nobadfuncs
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-nobadfuncs/nobadfuncs-asset/1.0.0/nobadfuncs-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.80 MiB    0.00% 670.23 KiB / 3.80 MiB   17.23% 1.15 MiB / 3.80 MiB   30.38% 1.50 MiB / 3.80 MiB   39.49% 1.69 MiB / 3.80 MiB   44.60% 1s 1.75 MiB / 3.80 MiB   46.07% 1s 2.03 MiB / 3.80 MiB   53.46% 1s 2.98 MiB / 3.80 MiB   78.53% 3.80 MiB / 3.80 MiB  100.00% 1s
Running nobadfuncs...
Finished nobadfuncs
```

Although the `skip` configuration makes it easy to disable a check, it is also possible to remove the check entirely.
This can be done by updating the `exclude-default-assets` configuration. For example, the following configuration
removes the `novendor` check entirely:

```
➜ echo 'default-tasks:
  resolvers:
    - https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz
  tasks:
    com.palantir.okgo:check-plugin:
      exclude-default-assets:
        - "com.palantir.godel-okgo-asset-novendor:novendor-asset"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
  paths:
    - "godel"' > godel/config/godel.yml
```

Verify that `novendor` is no longer present as a check:

```
➜ ./godelw check novendor
Error: provided checker type(s) [novendor] not valid: valid values are [compiles deadcode errcheck extimport golint govet importalias ineffassign outparamcheck unconvert varcheck]
```

Revert the local changes by running the following:

```
➜ git checkout -- godel
```
