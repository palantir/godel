okgo
====
okgo is a library and CLI for running static linting checks on Go code. okgo performs the following checks:

* [deadcode](https://github.com/tsenart/deadcode)
* [errcheck](https://github.com/kisielk/errcheck)
* [extimport](https://github.com/palantir/checks/tree/develop/extimport)
* [go vet](https://golang.org/cmd/vet/)
* [golint](https://github.com/golang/lint/tree/master/golint)
* [ineffassign](https://github.com/gordonklaus/ineffassign)
* [outparamcheck](https://github.com/palantir/checks/tree/develop/outparamcheck)
* [unconvert](https://github.com/mdempsky/unconvert)
* [varcheck](https://github.com/opennota/check/tree/master/cmd/varcheck)

Documentation
-------------
Documentation for okgo is provided in the Go code and as part of the application itself.

* Run `okgo --help` to get an overview of the commands and flags
* okgo is configured using a YML or JSON configuration file. Refer to the documentation in
  `apps/okgo/config/config.go` for information on the configuration parameters that are available.
* Refer to `apps/okgo/config/example_test.go` for sample configuration files

Development
-----------
Use the following commands for development. All paths in the example commands assume that they are run from the root
project directory of godel -- if the current working directory is `apps/okgo`, use `../../godelw` instead.

* Run `./godelw verify` to apply formatting, perform linting checks and run the gödel tests
* Run `./godelw test --tags=okgo` to run the okgo-specific tests (not included by default in the tests run by `./godlew verify`)
* Run `./godelw build` to build the okgo binary in `apps/okgo/build`

### Add a new check
* In order for a check to be added, it must be a Go program that has a main package
* The checks are managed and packaged by [amalgomate](https://github.com/palantir/amalgomate)
* Add the code required for the new check (the main package and any supporting code) to the vendor directory
* Edit `apps/okgo/checks.yml` and add an entry for the new check
  * Add an entry to `packages` where the key is the name of the check (no whitespace) and the value has a key named
    `main` and the value is the import path to the main package for the formatter. For more details on the config file
    format, refer to the documentation for amalgomate.
* Run `go generate` in the root directory of the project to re-generate the files in `generated_src`
* Add a definition for the check in `ckecks/definition.go`

### Generate
Run `go generate` in the `apps/gonform` directory to create or update the `generated_src` directory and the source files
within it. The `go generate` task for this project requires the amalgomate command to run. The version of amalgomate
used to build the distribution is included as a vendored dependency.
