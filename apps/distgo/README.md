distgo
======
distgo is a library and CLI for running, building, distributing and publishing products in a Go project.

Documentation
-------------
Documentation for distgo is provided in the Go code and as part of the application itself.

* Run `distgo --help` to get an overview of the commands and flags
* distgo is configured using a YML or JSON configuration file. Refer to the documentation in
  `apps/distgo/config/config.go` for information on the configuration parameters that are available.
* Refer to `apps/distgo/config/example_test.go` for sample configuration files

Development
-----------
Use the following commands for development. All paths in the example commands assume that they are run from the root
project directory of godel -- if the current working directory is `apps/distgo`, use `../../godelw` instead.

* Run `./godelw verify` to apply formatting, perform linting checks and run the gödel tests
* Run `./godelw test --tags=distgo` to run the distgo-specific tests (not included by default in the tests run by `./godlew verify`)
* Run `./godelw build` to build the `distgo` binary in `apps/distgo/build`
