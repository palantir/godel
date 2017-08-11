Summary
-------
`./godelw run` can be used to run a product from source.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1

([Link](https://github.com/nmiyake/echgo/tree/7799802bb82db52e99dda67edf9c98333b28fca3))

Run
---

`echgo` is now defined as a product and can be built using `build`. Its functionality can be tested using tests or by
invoking the executable that was built with `build`.

Although the program can be run by building the product and invoking the executable (or by running `go install` and
running the executable), this can be cumbersome for quick iteration. `./godelw run` can be used to quickly build and run
a product from source for faster iteration.

Use `./godelw run` to invoke `echgo`:

```
➜ ./godelw run foo
/usr/local/go/bin/go run /Volumes/git/go/src/github.com/nmiyake/echgo/main.go foo
foo
```

This uses `go run` to run the product. The above works because our project only defines a single product. If a project
defines multiple products, then the `--product` flag must be used to specify the product that should be run:

```
➜ ./godelw run --product=echgo foo
/usr/local/go/bin/go run /Volumes/git/go/src/github.com/nmiyake/echgo/main.go foo
foo
```

Because the `run` task uses `go run`, it does not build the product using the build parameters specified in the
configuration. This can be verified by running the command with the `-version` flag and verifying the output. Flags that
are meant to be processed by the program invoked by `run` must be prepended with `flag:` to distinguish them from flags
that for the `run` task itself. Run the following to run the equivalent of `echgo -version`:

```
➜ ./godelw run flag:-version
/usr/local/go/bin/go run /Volumes/git/go/src/github.com/nmiyake/echgo/main.go -version
echgo version: none
```

The output demonstrates that the code was run from source (if `echgo` was built using `./godelw build` and that
executable was run, the version output would be the output of `git describe` since we configured the build parameters to
set the version variable on build).

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo`
* Project is tagged as 0.0.1

([Link](https://github.com/nmiyake/echgo/tree/7799802bb82db52e99dda67edf9c98333b28fca3))

Tutorial next step
------------------

[Dist](https://github.com/palantir/godel/wiki/Dist)
