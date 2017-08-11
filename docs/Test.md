Summary
-------
`./godelw test` runs the Go tests in the project.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go` and `echo/echoer.go`

([Link](https://github.com/nmiyake/echgo/tree/0a649925e317b7896e537ef23a4885062a3ec9fb))

Run tests
---------

We will now add some tests to our program. Run the following to add tests for the `echo` package:

```
➜ echo 'package echo_test

import (
	"testing"

	"github.com/nmiyake/echgo/echo"
)

func TestEcho(t *testing.T) {
	echoer := echo.NewEchoer()
	for i, tc := range []struct {
		in   string
		want string
	}{
		{"foo", "foo"},
		{"foo bar", "foo bar"},
	} {
		if got := echoer.Echo(tc.in); got != tc.want {
			t.Errorf("case %d failed: want %q, got %q", i, tc.want, got)
		}
	}
}' > echo/echo_test.go
```

Run `./godelw test` to run all of the Go tests in the project:

```
➜ ./godelw test
ok  	github.com/nmiyake/echgo     	0.090s [no tests to run]
ok  	github.com/nmiyake/echgo/echo	0.132s
```

Commit the test to the repository:

```
➜ git add echo
➜ git commit -m "Add tests for echo package"
[master 404c745] Add tests for echo package
 1 file changed, 22 insertions(+)
 create mode 100644 echo/echo_test.go
➜ git status
On branch master
nothing to commit, working directory clean
```

Tutorial end state
------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores IDEA files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`

([Link](https://github.com/nmiyake/echgo/tree/404c745e6bc0f70f4d4b58b60502e5b9620a00a7))

Tutorial next step
------------------

[Build](https://github.com/palantir/godel/wiki/Build)

More
----

### Differences between `./godelw test` and `go test ./...`

`./godelw test` has the following advantages over running `go test ./...`:

* Aligns the output so that all of the test times line up
* Generates placeholder files in packages that do not contain tests (important for coverage purposes)
* Only runs tests for files that are part of the project (does not run tests in vendor directories)
  * This has been fixed in the Go tool itself as of Go 1.9

To demonstrate this, check out the project `github.com/palantir/checks`:

```
➜ mkdir -p $GOPATH/src/github.com/palantir && cd $_
➜ git clone https://github.com/palantir/checks.git
Cloning into 'checks'...
remote: Counting objects: 2800, done.
remote: Compressing objects: 100% (9/9), done.
remote: Total 2800 (delta 1), reused 3 (delta 0), pack-reused 2791
Receiving objects: 100% (2800/2800), 6.45 MiB | 4.21 MiB/s, done.
Resolving deltas: 100% (660/660), done.
Checking connectivity... done.
➜ cd checks
```

In Go 1.8 and earlier, running `go test ./...` fails immediately because there is code in the vendor directory that
cannot be built:

```
➜ go test ./...
vendor/github.com/palantir/godel/apps/distgo/cmd/publish/github_publish.go:27:2: cannot find package "github.com/google/go-github/github" in any of:
	/Volumes/git/go2/src/github.com/palantir/checks/vendor/github.com/google/go-github/github (vendor tree)
...
```

Go 1.9 fixes this issue by having `./...` no longer match paths in vendor directories. In Go 1.8 and earlier, the
equivalent can be achieved by running `go test $(go list ./... | grep -v /vendor/)`. Doing so produces the following:

```
➜ go test $(go list ./... | grep -v /vendor/)
ok  	github.com/palantir/checks/compiles	4.867s
ok  	github.com/palantir/checks/extimport	1.296s
?   	github.com/palantir/checks/gocd	[no test files]
?   	github.com/palantir/checks/gocd/cmd	[no test files]
?   	github.com/palantir/checks/gocd/cmd/gocd	[no test files]
ok  	github.com/palantir/checks/gocd/config	1.416s
ok  	github.com/palantir/checks/gocd/gocd	0.449s
?   	github.com/palantir/checks/gogenerate	[no test files]
?   	github.com/palantir/checks/gogenerate/cmd	[no test files]
?   	github.com/palantir/checks/gogenerate/cmd/gogenerate	[no test files]
ok  	github.com/palantir/checks/gogenerate/config	1.416s
ok  	github.com/palantir/checks/gogenerate/gogenerate	8.249s
?   	github.com/palantir/checks/golicense	[no test files]
?   	github.com/palantir/checks/golicense/cmd	[no test files]
?   	github.com/palantir/checks/golicense/cmd/golicense	[no test files]
ok  	github.com/palantir/checks/golicense/config	1.183s
ok  	github.com/palantir/checks/golicense/golicense	0.984s
ok  	github.com/palantir/checks/importalias	2.157s
?   	github.com/palantir/checks/nobadfuncs	[no test files]
ok  	github.com/palantir/checks/nobadfuncs/integration_test	13.334s
ok  	github.com/palantir/checks/nobadfuncs/nobadfuncs	13.895s
ok  	github.com/palantir/checks/novendor	2.132s
?   	github.com/palantir/checks/outparamcheck	[no test files]
ok  	github.com/palantir/checks/outparamcheck/exprs	1.905s
ok  	github.com/palantir/checks/outparamcheck/outparamcheck	1.729s
?   	github.com/palantir/checks/ptimports	[no test files]
?   	github.com/palantir/checks/ptimports/ptimports	[no test files]
```

Compare this to running `./godelw test`:

```
➜ ./godelw test
ok  	github.com/palantir/checks/compiles                   	4.242s
ok  	github.com/palantir/checks/extimport                  	0.928s
ok  	github.com/palantir/checks/gocd                       	0.444s [no tests to run]
ok  	github.com/palantir/checks/gocd/cmd                   	1.693s [no tests to run]
ok  	github.com/palantir/checks/gocd/cmd/gocd              	1.600s [no tests to run]
ok  	github.com/palantir/checks/gocd/config                	1.547s
ok  	github.com/palantir/checks/gocd/gocd                  	0.962s
ok  	github.com/palantir/checks/gogenerate                 	0.626s [no tests to run]
ok  	github.com/palantir/checks/gogenerate/cmd             	0.759s [no tests to run]
ok  	github.com/palantir/checks/gogenerate/cmd/gogenerate  	0.948s [no tests to run]
ok  	github.com/palantir/checks/gogenerate/config          	1.207s
ok  	github.com/palantir/checks/gogenerate/gogenerate      	5.901s
ok  	github.com/palantir/checks/golicense                  	0.929s [no tests to run]
ok  	github.com/palantir/checks/golicense/cmd              	0.818s [no tests to run]
ok  	github.com/palantir/checks/golicense/cmd/golicense    	0.986s [no tests to run]
ok  	github.com/palantir/checks/golicense/config           	0.857s
ok  	github.com/palantir/checks/golicense/golicense        	0.929s
ok  	github.com/palantir/checks/importalias                	1.123s
ok  	github.com/palantir/checks/nobadfuncs                 	0.360s [no tests to run]
ok  	github.com/palantir/checks/nobadfuncs/nobadfuncs      	9.438s
ok  	github.com/palantir/checks/novendor                   	1.059s
ok  	github.com/palantir/checks/outparamcheck              	0.564s [no tests to run]
ok  	github.com/palantir/checks/outparamcheck/exprs        	0.775s
ok  	github.com/palantir/checks/outparamcheck/outparamcheck	0.807s
ok  	github.com/palantir/checks/ptimports                  	0.168s [no tests to run]
ok  	github.com/palantir/checks/ptimports/ptimports        	0.847s [no tests to run]
```

This output is much easier to read and has the advantage of ignoring all excluded directories.

Restore the working directory to be the original directory:

```
➜ cd $GOPATH/src/github.com/nmiyake/echgo
```

### Generate coverage reports

A combined coverage report for the project can be generated by running the command
`./godelw test cover --coverage-output=<file>`:

```
➜ ./godelw test cover --coverage-output=cover.out
ok  	github.com/nmiyake/echgo     	0.108s	coverage: 0.0% of statements [no tests to run]
ok  	github.com/nmiyake/echgo/echo	0.161s	coverage: 100.0% of statements
➜ cat cover.out
mode: count
github.com/nmiyake/echgo/main.go:11.13,14.2 2 0
github.com/nmiyake/echgo/echo/echo.go:3.25,5.2 1 1
github.com/nmiyake/echgo/echo/echo.go:9.47,11.2 1 2
```

This command runs the Go tests all of the packages in the project in coverage mode using `-covermode=count` and
generates a single combined coverage file for all of the packages in the project.

In contrast, running `go test --cover ./...` produces the following:

```
➜ go test --cover ./...
?   	github.com/nmiyake/echgo	[no test files]
ok  	github.com/nmiyake/echgo/echo	0.090s	coverage: 100.0% of statements
```

Go cover does not count packages that do not contain tests towards the coverage count by default, and also does not
provide a single command that allows the use of profile flags with multiple packages.

Remove the output by running the following:

```
➜ rm cover.out
```

### Generate JUnit reports

The `./godelw test --junit-output=<file>` command can be used to generate a JUnit-style output XML file that summarizes
the results of running the tests:

```
➜ ./godelw test --junit-output=output.xml
testing: warning: no tests to run
PASS
ok  	github.com/nmiyake/echgo     	0.121s [no tests to run]
=== RUN   TestEcho
--- PASS: TestEcho (0.00s)
PASS
ok  	github.com/nmiyake/echgo/echo	0.112s
➜ cat output.xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite tests="1" failures="0" time="0.112" name="github.com/nmiyake/echgo/echo">
		<properties>
			<property name="go.version" value="go1.9"></property>
		</properties>
		<testcase classname="echo" name="TestEcho" time="0.000"></testcase>
	</testsuite>
</testsuites>
```

Remove the output by running the following:

```
➜ rm output.xml
```
