gödel
=====

[![Bintray](https://img.shields.io/bintray/v/palantir/releases/godel.svg)](https://bintray.com/palantir/releases/godel/_latestVersion)
[![CircleCI](https://circleci.com/gh/palantir/godel.svg?style=shield)](https://circleci.com/gh/palantir/godel)

gödel is a Go build tool that provides tasks for configuring, formatting, checking, testing, building and publishing Go
projects in a declarative, consistent and reproducible manner across different platforms and environments. gödel can be
used in both local development environments and for verifying the correctness of project in CI environments. gödel uses
declarative configuration to define the parameters for a project and provides an executable that orchestrates build
tasks using standard Go commands. It centralizes project configuration and eliminates the need for custom build scripts
that conflate configuration with logic. gödel is designed to be portable, fast and lightweight -- adding it to a project
consists of copying a single file and directory into the project and adds less than 50kb of version-controlled material.

Features
--------
* Add to a project by copying the `godelw` script and the `godel` configuration directory into a project
* `./godelw git-hooks` installs Git commit hook that formats files on commit
* `./godelw idea` creates and configures an IntelliJ IDEA project for the project
* Supports configuring directories and files that should be excluded by the tool (vendor directory is excluded by default)
* `./godelw format` formats all code in a project
* `./godelw check` runs a variety of code linting checks on all the code in a project
* `./godelw license` applies a specified license header to all Go files in a project
  * Supports configuring custom license headers for specific directories or files
* `./godelw generate` runs `go generate` tasks for a project
* `./godelw test` runs the tests in the project
  * Configuration can be used to define test sets (such as integration tests) and run specific test sets
  * Supports outputting the test results in a JUnit XML format
* `./godelw build` builds executables for `main` packages in the project
  * Supports cross-platform compilation
  * Supports configuration of `ldflag` for version and other variables
  * Installs packages by default to speed up repeated builds
* `./godelw dist` creates distribution files for products
  * Supports creating `tgz` and `rpm` distributions
  * Supports customizing creation of distribution using scripts
* `./godelw publish` publishes artifacts to Bintray or Artifactory
* `palantir/godel/pkg/products` package provides a mechanism to easily write integration tests for gödel projects
  * Provides a function that builds the product executable or distribution and provides a path to invoke it
* `./godelw update` updates gödel to the version specified in `godel/config/godel.properties`
* `./godelw github-wiki` mirrors a documents directory to a GitHub Wiki repository
* `./godelw verify` runs the `format`, `import`, `license`, `check` and `test` tasks
  * Can be used locally as a single command to apply changes and run checks
  * Can be used in CI to verify that a project is in the proper state without applying changes

This list is not exhaustive -- run `./godelw --help` for a list of all of the available commands.

Documentation
-------------
Documentation for this project is in the `docs` directory and the [GitHub Wiki](https://github.com/palantir/godel/wiki)
(the GitHub Wiki mirrors the contents of the `docs` directory).

Development
-----------
The code for the tasks provided by gödel is in the `cmd` directory. gödel tasks fall into 2 categories: those whose
functionality are implemented directly in gödel packages and those whose functionality is implemented by a sub-program
that exposes its tasks as library functions that are directly callable.

The `app.go` file in the top-level `godel` package registers the top-level tasks available in the CLI. It registers the
tasks whose functionality is implemented directly in gödel directly. The tasks provided by sub-programs are defined in
`cmd/clicmds/cfgcli.go`.

After making changes to the code, run `./godelw verify` to format the code, apply the proper license headers, update
dependency information, run code linting checks and all unit tests.

gödel also defines integration tests in the `test/integration` directory. The tests in this file create a distribution
of gödel and run tests against it. Run `./godelw test --tags=integration` to run the integration tests (the integration
tests are not run by `./godelw verify` or `./godelw test` by default).

### Sub-applications in the apps directory
The functionality for some gödel tasks are provided by sub-applications. `distgo`, `gonform`, `gunit` and `okgo` are
such sub-applications and are located in the `apps` directory. These sub-applications can be compiled and run as self-
contained CLI applications, but also expose functionality as libraries.

Changes that are made in these sub-applications directly impact gödel, and many of the gödel integration tests test the
functionality provided by these sub-applications. The sub-applications also have their own set of tests. Each of the
sub-applications in the `apps` directory have their own test suite that is the name of the sub-application and can be
invoked using its tag -- for example, to run the tests for `distgo`, run `./godelw test --tags=distgo`.

Refer to the README files in the sub-applications for more information on application-specific development.

### Sub-applications outside of the repository
Some tasks such as `imports` and `license` use sub-applications that are defined outside of the repository (in this
case, `gocd` and `golicense`, respectively) and vendor the sub-applications to use their functionality. Changing these
tasks is akin to changing a vendored library -- locate the original repository for the library, make changes there and
then update the vendored library in gödel.

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
