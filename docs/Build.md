Summary
-------
`./godelw build` builds the products in the project based on the build configuration.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`

([Link](https://github.com/nmiyake/echgo/tree/404c745e6bc0f70f4d4b58b60502e5b9620a00a7))

Build
-----

Now that we have a functional program with some tests, it's time to think about how to build and distribute the program.
We will start by considering how to build the program. We will be providing binary distributions for `darwin-amd64` and
`linux-amd64` systems. We also want to include a version variable in the product so that users can determine the version
of the product by invoking it with a `-version` flag.

Start by running `./godelw build` in the project to observe the default behavior:

```
➜ ./godelw build
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/unspecified/darwin-amd64/echgo
Finished building echgo for darwin-amd64 (0.442s)
```

The default build settings builds all `main` packages for the executing platform. As stated by the output, this command
built an executable for the OS/architecture of the host system in the `build` directory of the project. The executable
is fully built and can be run:

```
➜ ./build/unspecified/darwin-amd64/echgo foo
foo
```

The output of the `build` command can be removed by running `./godelw clean`. The output directory can be specified
using the `output-dir` key, but we will keep it as `build` for this product. We do not want to track build artifacts in
git, so run the following to add `/build/` to the `.gitignore` file and commit it:

```
➜ echo '/build/' >> .gitignore
➜ git add .gitignore
➜ git commit -m "Update .gitignore to ignore build directory"
[master 4e51282] Update .gitignore to ignore build directory
 1 file changed, 1 insertion(+)
```

The first step of customizing the build parameters for the product is to explicitly define the products in the project.
Builds operate on the concepts of products, where a product is a Go `main` package. A single project may have multiple
products. Define the `echgo` product for this project as follows:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .' > godel/config/dist.yml
```

The `main-pkg` key specifies that the main package for the product is located at `.` (relative to the root directory of
the project).

Define OS/architectures
-----------------------

By default, the `./godelw build` command builds an executable that matches the OS and architecture of the host system.
The `os-archs` parameter can be used to define the exact OS/architecture pairs for which a product should be built. Run
the following command to update the configuration to specify that `echgo` should be built for `darwin-amd64` and
`linux-amd64`:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64' > godel/config/dist.yml
```

`os-archs` is a list of OS and architectures for which the product should be built. For a list of available
OS/architecture combinations, refer to the [Go documentation](https://golang.org/doc/install/source#environment) or run
`go tool dist list`.

Run `./godelw build` with the updated configuration:

```
➜ ./godelw build
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/unspecified/darwin-amd64/echgo
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/unspecified/linux-amd64/echgo
Finished building echgo for darwin-amd64 (0.362s)
Finished building echgo for linux-amd64 (0.365s)
```

As indicated by the output, the product has now been built for both `darwin-amd64` and `linux-amd64`.

NOTE: if you have not cross-compiled a Go program before, the command may produce an error of the form:
      `go install <path>: mkdir <path>: permission denied`. gödel should provide a description of the error and a
      suggested fix. This error occurs because the standard Go distribution only includes the compiled object files for
      the standard library for the target platform. The gödel command attempts to build the standard library for all
      target platforms, but the Go library is often in a location that is not writable by the standard user. When
      cross-compiling, gödel expects the standard library for the target platform to be present (this saves time by
      ensuring that the standard library is not constantly re-compiled). Running `go install std` using `sudo` (or as a
      user who can write to the Go library directory) for the target OS/architecture should resolve the issue. For
      example, if you were running on a `darwin` system and encountered an error having to do with the `linux` libraries
      not being installed, run `sudo env GOOS=linux GOARCH=amd64 go install std`.

Set a version variable
----------------------

The `version-var` configuration can be used to set the value of a variable to be the version of the project determined
at build time. The project version is based on the output of `git describe`.

Update `main.go` to add support for printing a version variable:

```
➜ echo 'package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/nmiyake/echgo/echo"
)

var version = "none"

func main() {
	versionVar := flag.Bool("version", false, "print version")
	flag.Parse()
	if *versionVar {
		fmt.Println("echgo version:", version)
		return
	}
	echoer := echo.NewEchoer()
	fmt.Println(echoer.Echo(strings.Join(flag.Args(), " ")))
}' > main.go
```

Running this program with `-version` outputs the value `none`:

```
➜ go run main.go -version
echgo version: none
```

The desired behavior is to set the `version` variable to be the version of the program when it is built. Configure this
by updating the build configuration as follows:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64' > godel/config/dist.yml
```

Run `./godelw build` and invoke the build executable to see that the version variable has been set:

```
➜ ./godelw build
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/unspecified/linux-amd64/echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/unspecified/darwin-amd64/echgo
Finished building echgo for linux-amd64 (0.495s)
Finished building echgo for darwin-amd64 (0.501s)
➜ ./build/unspecified/darwin-amd64/echgo -version
echgo version: unspecified
```

The fact that the output is now "unspecified" (rather than "none" as specified in the source code) demonstrates that the
version was set by the build. The version is "unspecified" because there are no git tags present. We can fix this by
tagging a release.

Commit the files that were modified and tag a release:

```
➜ git add main.go godel/config/dist.yml
➜ git commit -m "Add version variable and define build configuration"
[master 7799802] Add version variable and define build configuration
 2 files changed, 20 insertions(+), 2 deletions(-)
➜ git tag 0.0.1
```

Now that the repository is tagged, run `./godelw build` and run version on the executable:

```
➜ ./godelw build
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1/linux-amd64/echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1/darwin-amd64/echgo
Finished building echgo for linux-amd64 (0.362s)
Finished building echgo for darwin-amd64 (0.365s)
➜ ./build/0.0.1/darwin-amd64/echgo -version
echgo version: 0.0.1
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

([Link](https://github.com/nmiyake/echgo/tree/7799802bb82db52e99dda67edf9c98333b28fca3))

Tutorial next step
------------------

[Run](https://github.com/palantir/godel/wiki/Run)

More
----

### Differences between `./godelw build` and `go build ./...`

The primary difference between `./godelw build` and `go build ./...` is that `./godelw build` uses the `dist.yml`
configuration to build the outputs, which gives much more control over the build process. It also creates binaries for
all of its products, whereas `go build` will only generate a binary if only a `main` package is built.

`./godelw build` also strives to be efficient in its operation by doing the following:

* Build operations are performed in parallel
* Builds operate in "install" mode so that intermediate object files are generated and re-used on subsequent runs
* Products will only be built if code has changed in such a way that necessitates a re-build

These optimizations make building projects that build binaries for multiple different OS/architecture combinations
significantly faster than simply using `go build` to build them serially.

### Build specific products

By default, `./godelw build` will build all of the products defined for a project. Product names can be specified as
arguments to build only those products. For example, if the `echgo` project defined multiple products, you could specify
that you want to build only `echgo` by running the following:

```
➜ ./godelw build echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1/darwin-amd64/echgo
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1/linux-amd64/echgo
Finished building echgo for linux-amd64 (0.402s)
Finished building echgo for darwin-amd64 (0.405s)
```

### Specify build environment variables

In some instances, products will want specific environment variables set during the build phase. These can be specified
using the `environment` map, where the key specifies the name of the environment variable and the value specifies the
value of the environment variable. For example, the following configuration sets `CGO_ENABLED` to `0`:

```yaml
products:
  echgo:
    build:
      main-pkg: .
      environment:
        CGO_ENABLED: "0"
```

### Specify build arguments

In some cases, it is desirable to have values that are dynamically computed passed to the `build` command. It is
possible to specify a script that will be executed to generate the arguments that are passed to the build script by
using the `build-args-script` configuration key. The content of the script is evaluated when `build` is run for the
product and each line of output is passed as a separate argument to the `build` command. For example, consider the
following configuration:

```yaml
products:
  echgo:
    build:
      main-pkg: .
      build-args-script: |
                          YEAR=$(date +%Y)
                          echo "-ldflags"
                          echo "-X"
                          echo "main.year=$YEAR"
```

This configuration evaluates `$(date +%Y)` when the `echgo` product is built and passes the arguments "-ldflags" "-X"
"main.year=$YEAR" to the build command.
