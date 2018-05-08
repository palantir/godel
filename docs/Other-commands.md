Summary
-------
Commands such as `./godelw project-version`, `./godelw packages` and `./godelw artifacts` can be used as inputs to other
commands.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
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
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to the latest version

Other commands
--------------
The declarative configuration provided by gödel can be used as inputs to other programs or scripts for convenience. For
example, you may want to run a `go install` on all of the project packages so that the object files are available in
`$GOPATH/pkg` and any `main` packages are installed. `go install ./...` will often not work on products that vendor
their dependencies.

The `./godelw packages` command can be used to print all of the packages in a project based on the project
configuration:

```
➜ ./godelw packages
./.
./echo
./generator
./integration_test
```

This command excludes any packages that would be excluded based on the `exclude` configuration. This command can be
combined with `go install` (or any program that takes a list of packages as inputs):

```
➜ go install $(./godelw packages)
```

`./godelw project-version` can be used to print the version of the project as determined by gödel (this is also the
string used as the version for the `dist` and `publish` tasks):

```
➜ ./godelw project-version
1.0.0
```

`./godelw artifacts` can be invoked with `build` or `dist` (and optionally a list of products) to output the location in
which the build and distribution artifacts will be generated:

```
➜ ./godelw artifacts build echgo2
out/build/echgo2/1.0.0/darwin-amd64/echgo2
out/build/echgo2/1.0.0/linux-amd64/echgo2
➜ ./godelw artifacts dist
out/dist/echgo2/1.0.0/os-arch-bin/echgo2-1.0.0-darwin-amd64.tgz
out/dist/echgo2/1.0.0/os-arch-bin/echgo2-1.0.0-linux-amd64.tgz
```

These commands can be useful to use as inputs to other programs or as part of scripts/CI tasks for a project.

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
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
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to the latest version

Tutorial next step
------------------
[Conclusion](https://github.com/palantir/godel/wiki/Tutorial-conclusion)
