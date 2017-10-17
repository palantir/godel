Summary
-------
The `godel/config/exclude.yml` configuration specifies patterns for paths and names that should be ignored.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function

([Link](https://github.com/nmiyake/echgo/tree/08752b2ae998c14dd5abb789cebc8f5848f7cf4e))

Exclude names/paths
-------------------

In the previous step of the tutorial, we updated the `godel/config/exclude.yml` configuration file to specify that a
generated source file should be excluded from checks. This portion of the tutorial examines the file in more detail.

Examine the current state of `godel/config/exclude.yml`:

```
➜ cat godel/config/exclude.yml
names:
  - "\\..+"
  - "vendor"
paths:
  - "godel"
  - "echo/type_string.go"
```

As seen in the previous section, specifying `echo/type_string.go` as an exclude path caused this file to be ignored for
all checks and operations. However, listing individual paths like this can be cumbersome and fragile -- for example, if
we wanted to use `stringer` for other types or in other packages, this approach would require us to manually add each
new entry.

We can make this approach more general by excluding all files that end in `_string.go` instead. Run the following to
update the configuration:

```
➜ echo 'names:
  - "\\\\..+"
  - "vendor"
  - ".+_string.go"
paths:
  - "godel"' > godel/config/exclude.yml
```

Verify that the `golint` check still succeeds:

```
➜ ./godelw check golint
Running golint...
```

Check in the changes:

```
➜ git add godel/config/exclude.yml
➜ git commit -m "Update exclude.yml"
[master 1982133] Update exclude.yml
 1 file changed, 1 insertion(+), 1 deletion(-)
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
* `godel/config/dist.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* Go files have license headers
* `godel/config/generate.yml` is configured to generate string function
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files

([Link](https://github.com/nmiyake/echgo/tree/1982133dbe7c811f1e2d71f4dcc25ff20f84146a))

Tutorial next step
------------------

[Write integration tests](https://github.com/palantir/godel/wiki/Integration-tests)

More
----

### Names, paths and default configuration

`godel/config/exclude.yml` starts with the following defaults:

```yml
names:
  - "\\..+"
  - "vendor"
paths:
  - "godel"
```

The `names` configuration specifies a list of names that should be ignored. The values can be literals or regular
expressions. In either case, the value must fully match the name to be a match -- for regular expressions, this means
that a name is considered a match only if the entire name matches the regular expression. The `names` configuration is
useful for ignoring classes of files -- for example, `.*\\.pb\\.go` can be used to ignore all files that end in
`.pb.go`.

The default configuration for `names` ignores all hidden files (names that start with `.`) and all vendor directories.
`vendor` is specified as a name rather than a path so that vendor directories are ignored in all directories.

The `paths` configuration specifies a list of paths that should be ignored. Paths can be specified as literal values or
as globs. The `paths` configuration is useful for ignoring specific files or directories that should be ignored.

The default configuration for `paths` ignores the `godel` directory in the root level of the project.

In the case of both `names` and `paths`, if the configuration matches a file, that file is ignored, while if it matches
a directory, then that directory (and all of its contents) is ignored.
