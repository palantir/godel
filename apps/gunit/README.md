gunit
=====
gunit is a library and CLI for running Go tests, generating coverage reports and generating JUnit-style XML output. It
uses [go-junit-report](https://github.com/jstemmer/go-junit-report) as the library for formatting test reports and
[gt](https://godoc.org/rsc.io/gt) for running cached tests.

Documentation
-------------
Documentation for `gunit` is provided in the Go code and as part of the application itself.

* Run `gunit --help` to get an overview of the commands and flags
* gunit is configured using a YML or JSON configuration file. Refer to the documentation in
  `apps/gunit/config/config.go` for information on the configuration parameters that are available.
* Refer to `apps/gunit/config/example_test.go` for sample configuration files

Development
-----------
Use the following commands for development. All paths in the example commands assume that they are run from the root
project directory of godel -- if the current working directory is `apps/gunit`, use `../../godelw` instead.

* Run `./godelw verify` to apply formatting, perform linting checks and run the gödel tests
* Run `./godelw test --tags=gunit` to run the gunit-specific tests (not included by default in the tests run by `./godlew verify`)
* Run `./godelw build` to build the gunit binary in `apps/gunit/build`

### Add a new test utility
* In order for a test utility to be added, it must be a Go program that has a main package
* The test utilities are managed and packaged by [amalgomate](https://github.com/palantir/amalgomate)
* Add the code required for the new test utility (the main package and any supporting code) to the vendor directory
* Edit `apps/gunit/testers.yml` and add an entry for the new test utility
  * Add an entry to `packages` where the key is the name of the test utility (no whitespace) and the value has a key
    named `main` and the value is the import path to the `main` package for the test utility. For more details on the
    config file format, refer to the documentation for amalgomate.
* Run `go generate` in the root directory of the project to re-generate the files in `generated_src`

### Generate
Run `go generate` in the `apps/gonform` directory to create or update the `generated_src` directory and the source files
within it. The `go generate` task for this project requires the amalgomate command to run. The version of amalgomate
used to build the distribution is included as a vendored dependency.
