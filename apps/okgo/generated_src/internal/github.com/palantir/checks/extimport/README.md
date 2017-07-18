extimport
=========

`extimport` is a program that verifies that there are no imports that reference external packages in a go project. If
external package imports are found, the program prints them and exits with an exit code of 1. An exit code of 0
indicates that there are no external package imports in the project, which means that it should be possible to build the
project in a default `$GOPATH`.

A package is considered external if it is not in the standard Go library and not resolvable within the project
directory itself (the package is not in the project or vendored in the project). An import is considered to be directly
external if it imports an external package. An import is considered transitively external if the imported package itself
is not external, but one of its dependent packages is external.

Given a package, `extimports` checks the imports of all of the go files and go test files in that package. However, when
checking transitive external package dependencies, only non-test go files are considered (that is, the check will not
fail if a test file of an imported package has an external dependency).

Usage
=====
`extimport` uses its current working directory as the project root. If no arguments are provided, it is invoked on all
of the go packages it can find in the current working directory and its subdirectories. If arguments are provided, they
are interpreted as packages relative to the working directory, and only the specified packages will be checked.
