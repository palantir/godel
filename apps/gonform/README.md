gonform
=======
gonform is a library and CLI for running formatting operations on Go files. It uses
[gofmt](https://golang.org/cmd/gofmt/) and [ptimports](https://github.com/palantir/checks/tree/develop/ptimports) as the
libraries that perform the formatting operations.

Documentation
-------------
Documentation for gonform is provided in the Go code and as part of the application itself.

* Run `gonform --help` to get an overview of the commands and flags
* gonform is configured using a YML or JSON configuration file. Refer to the documentation in
  `apps/gonform/config/config.go` for information on the configuration parameters that are available.
* Refer to `apps/gonform/config/example_test.go` for sample configuration files

Development
-----------
Use the following commands for development. All paths in the example commands assume that they are run from the root
project directory of godel -- if the current working directory is `apps/gonform`, use `../../godelw` instead.

* Run `./godelw verify` to apply formatting, perform linting checks and run the gödel tests
* Run `./godelw test --tags=gonform` to run the gonform-specific tests (not included by default in the tests run by `./godlew verify`)
* Run `./godelw build` to build the gonform binary in `apps/gonform/build`

### Add a new formatter
* In order for a formatter to be added, it must be a Go program that has a main package
* The formatters are managed and packaged by [amalgomate](https://github.com/palantir/amalgomate)
* Add the code required for the new check (the main package and any supporting code) to the vendor directory
* Edit `apps/gonform/formatters.yml` and add an entry for the new formatter
  * Add an entry to `packages` where the key is the name of the formatter (no whitespace) and the value has a key
    named `main` and the value is the import path to the main package for the formatter. For more details on the
    config file format, refer to the documentation for amalgomate.
* Run `go generate` in the root directory of the project to re-generate the files in `generated_src`

### Generate
Run `go generate` in the `apps/gonform` directory to create or update the `generated_src` directory and the source files
within it. The `go generate` task for this project requires the amalgomate command to run. The version of amalgomate
used to build the distribution is included as a vendored dependency.
