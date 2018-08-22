Summary
-------
`./godelw test` runs the Go tests in the project.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go` and `echo/echoer.go`

Run tests
---------
We will now add some tests to our program. Run the following to add tests for the `echo` package:

```
➜ SRC='package echo_test

import (
	"testing"

	"PROJECT_PATH/echo"
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
}' && SRC=${SRC//PROJECT_PATH/$PROJECT_PATH} && echo "$SRC" > echo/echo_test.go
```

Run `./godelw test` to run all of the Go tests in the project:

```
➜ ./godelw test
?   	github.com/nmiyake/echgo2     	[no test files]
ok  	github.com/nmiyake/echgo2/echo	0.002s
```

Commit the test to the repository:

```
➜ git add echo
➜ git commit -m "Add tests for echo package"
[master 362ea06] Add tests for echo package
 1 file changed, 22 insertions(+)
 create mode 100644 echo/echo_test.go
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`

Tutorial next step
------------------
[Build](https://github.com/palantir/godel/wiki/Build)

More
----
### Differences between `./godelw test` and `go test ./...`
`./godelw test` has the following advantages over running `go test ./...`:

* Aligns the output so that all of the test times line up
* Only runs tests for files that are part of the project (uses the exclude parameters)

The test task also allows tags to be specified, which allows tests such as integration tests to be treated separately.
The [integration tests](https://github.com/palantir/godel/wiki/Integration-Tests) section of the tutorial covers this
in more detail.

### Generate JUnit reports
The `./godelw test --junit-output=<file>` command can be used to generate a JUnit-style output XML file that summarizes
the results of running the tests (implemented using [go-junit-report](https://github.com/jstemmer/go-junit-report)):

```
➜ ./godelw test --junit-output=output.xml
?   	github.com/nmiyake/echgo2     	[no test files]
=== RUN   TestEcho
--- PASS: TestEcho (0.00s)
PASS
ok  	github.com/nmiyake/echgo2/echo	0.002s
```

Verify that this operation wrote a JUnit report:

```
➜ cat output.xml
<?xml version="1.0" encoding="UTF-8"?>
<testsuites>
	<testsuite tests="1" failures="0" time="0.002" name="github.com/nmiyake/echgo2/echo">
		<properties>
			<property name="go.version" value="go1.10.1"></property>
		</properties>
		<testcase classname="echo" name="TestEcho" time="0.000"></testcase>
	</testsuite>
</testsuites>
```

Remove the output by running the following:

```
➜ rm output.xml
```

### Run tests with flags
In some instances, we may want to specify flags for the `go test` operation -- for example, in the previous section, we
wanted to pass `-count=1` to force the tests to run without using cache. Other common test flags include `-timeout` to
specify a timeout, `-p` to specify the number of test binaries that can be run in parallel, `-json` to print the output
as JSON, etc.

The flags provided after the `./godelw test` command are passed directly to the underlying `go test` invocation. The
`--` separator should be used to ensure that the flags are not interpreted.

For example, the following command prints the output as JSON:

```
➜ ./godelw test -- -json
{"Time":"2018-08-22T16:25:14.6287415Z","Action":"output","Package":"github.com/nmiyake/echgo2","Output":"?   \tgithub.com/nmiyake/echgo2\t[no test files]\n"}
{"Time":"2018-08-22T16:25:14.6295891Z","Action":"skip","Package":"github.com/nmiyake/echgo2","Elapsed":0.001}
{"Time":"2018-08-22T16:25:14.6315847Z","Action":"run","Package":"github.com/nmiyake/echgo2/echo","Test":"TestEcho"}
{"Time":"2018-08-22T16:25:14.6316225Z","Action":"output","Package":"github.com/nmiyake/echgo2/echo","Test":"TestEcho","Output":"=== RUN   TestEcho\n"}
{"Time":"2018-08-22T16:25:14.6316415Z","Action":"output","Package":"github.com/nmiyake/echgo2/echo","Test":"TestEcho","Output":"--- PASS: TestEcho (0.00s)\n"}
{"Time":"2018-08-22T16:25:14.6316551Z","Action":"pass","Package":"github.com/nmiyake/echgo2/echo","Test":"TestEcho","Elapsed":0}
{"Time":"2018-08-22T16:25:14.631668Z","Action":"output","Package":"github.com/nmiyake/echgo2/echo","Output":"PASS\n"}
{"Time":"2018-08-22T16:25:14.6316795Z","Action":"output","Package":"github.com/nmiyake/echgo2/echo","Output":"ok  \tgithub.com/nmiyake/echgo2/echo\t(cached)\n"}
{"Time":"2018-08-22T16:25:14.6317011Z","Action":"pass","Package":"github.com/nmiyake/echgo2/echo","Elapsed":0}
```
