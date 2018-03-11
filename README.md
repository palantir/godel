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

gödel is also highly extensible and configurable. The core functionality of gödel is provided by plugins and assets,
and it is easy to write new plugins or assets and to configure a gödel instance to use custom plugins or assets as
needed.

Features
--------
The following features are provided by a default gödel installation (either as builtin tasks or default plugin tasks):

* Add to a project by running the `godelinit` program
* `./godelw git-hooks` installs Git commit hook that formats files on commit
* `./godelw goland` creates and configures a GoLand project for the project
* Supports configuring directories and files that should be excluded by the tool
* `./godelw format` formats all code in a project
* `./godelw check` runs a variety of code linting checks on all the code in a project
  * Default configuration includes a wide variety of checks that catch common errors
  * Custom checks can be added as assets
* `./godelw license` applies a specified license header to all Go files in a project
  * Supports configuring custom license headers for specific directories or files
* `./godelw test` runs the tests in the project
  * Configuration can be used to define test sets (such as integration tests) and run specific test sets
  * Supports outputting the test results in a JUnit XML format
* `./godelw build` builds executables for `main` packages in the project
  * Supports cross-platform compilation
  * Supports configuration of `ldflag` for version and other variables
  * Installs packages by default to speed up repeated builds
* `./godelw dist` creates distribution files for products
  * Supports creating `tgz` distributions
  * Supports customizing creation of distribution using scripts
  * Supports creating custom distributions using assets
* `./godelw publish` publishes artifacts to Bintray, Artifactory or GitHub
  * Supports other forms of publishing using assets
* `palantir/godel/pkg/products` package provides a mechanism to easily write integration tests for gödel projects
  * Provides a function that builds the product executable or distribution and provides a path to invoke it
* `./godelw update` updates gödel to the version specified in `godel/config/godel.properties`
* `./godelw verify` runs all of the tasks that declare support for verification
  * Can be used locally as a single command to apply changes and run checks
  * Can be used in CI to verify that a project is in the proper state without applying changes
* `./godelw github-wiki` mirrors a documents directory to a GitHub Wiki repository

This list is not exhaustive -- run `./godelw --help` for a list of all of the available commands. Furthermore, custom

Documentation
-------------
Documentation for this project is in the `docs` directory and the [GitHub Wiki](https://github.com/palantir/godel/wiki)
(the GitHub Wiki mirrors the contents of the `docs` directory).

License
-------
This project is made available under the [Apache 2.0 License](http://www.apache.org/licenses/LICENSE-2.0).
