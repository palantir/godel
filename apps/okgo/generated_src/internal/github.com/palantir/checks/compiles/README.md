compiles
========

`compiles` verifies that all of the go packages that are part of a project compile properly. This is similar to the
check done by `go build ./...`, but goes further by also verifying that test files (both those that are part of a
package and those that are part of a `_test` package) also compile and build without errors.

Usage
=====
`compiles` uses its current working directory as the project root. If no arguments are provided, it is invoked on all
of the go packages it can find in the current working directory and its subdirectories. If arguments are provided, they
are interpreted as packages relative to the working directory, and only the specified packages will be checked.
