Summary
-------
The `verify` task can be used to perform all of the primary operations and checks for a project.

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
* `godel/config/exclude.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test.yml` is configured to specify the "integration" tag
* `docs` contains documentation

([Link](https://github.com/nmiyake/echgo/tree/17c7406291096306e92c6f82da2df09388766693))

Verify a project
----------------

Over the course of this tutorial, we have configured and used many different tasks -- `format` to format code, `check`
to run static checks, `generate` to run "go generate" tasks, `license` to apply license headers and `test` to run tests,
among others. Remembering to run all of these tasks separately can be a challenge -- in most cases, we simply want to
make sure that all of the tasks are run properly and that our code is working.

The `./godelw verify` task can be used to do exactly this -- it runs the `format`, `generate`, `license`, `check` and
`test` tasks. This single task is typically sufficient to ensure that all of the code for a project meets the
declarative specifications and that all of the tests in the project pass.

Run `./godelw verify` to verify that it runs all the tasks and succeeds:

```
➜ ./godelw verify
Running gofmt...
Running ptimports...
Running gogenerate...
Running gocd...
Running golicense...
Running compiles...
Running deadcode...
Running errcheck...
Running extimport...
Running golint...
Running govet...
Running importalias...
Running ineffassign...
Running nobadfuncs...
Running novendor...
Running outparamcheck...
Running unconvert...
Running varcheck...
ok  	github.com/nmiyake/echgo                 	0.160s [no tests to run]
ok  	github.com/nmiyake/echgo/echo            	0.185s
ok  	github.com/nmiyake/echgo/generator       	0.166s [no tests to run]
ok  	github.com/nmiyake/echgo/integration_test	1.279s
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
* `integration_test` contains integration tests
* `godel/config/test.yml` is configured to specify the "integration" tag
* `docs` contains documentation

([Link](https://github.com/nmiyake/echgo/tree/17c7406291096306e92c6f82da2df09388766693))

Tutorial next step
------------------

[Set up CI to run tasks](https://github.com/palantir/godel/wiki/CI-setup)

More
----

### Verify that a project complies with checks without applying changes

The `--apply=false` flag can be used to run `verify` in a mode that verifies that the project passes checks without
actually modifying it. Specifically, it runs the `format`, `generate` and `license` tasks in a mode that verifies that
the project complies with the settings without modifying the project. The task exits with a non-0 exit code if any of
the verifications fail. This task is suitable to run in a CI environment where one wants to verify that a project passes
all of its checks without actually modifying it.

The one exception to this principle is the `generate` task. Because `go generate` can run arbitrary tasks, there is no
general way in which verification can be performed besides running the `generate` tasks and comparing the state of the
impacted files before and after the task was run. Running `generate` with the `--apply=false` flag will print
information about how the state after the run differs from the state before the run for the files specified in the
configuration and will cause the verification to fail if differences exist, but the modification will have already been
made and will persist after the task.

### Skip specific tasks

In some cases, you may want to run `verify` but skip specific aspects of it -- for example, if the `generate` task takes
a long time to run and you want to run all verification tasks except for `generate`, you can use the `--skip-generate`
flag to skip the generation step. Run `./godelw verify --help` for a full list of the skip flags.
