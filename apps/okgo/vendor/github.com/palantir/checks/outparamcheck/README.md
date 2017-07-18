outparamcheck
=============

`outparamcheck` is a static code checker for Go based on [errcheck](https://github.com/kisielk/errcheck). It verifies
that functions that take output parameters defined as `interface{}` types are passed pointers to an object rather than
a concrete object.

A canonical example of this is the `json.Unmarshal` function, which has the following definition:

```go
func Unmarshal(data []byte, v interface{}) error
```

As noted in the godoc ("Unmarshal parses the JSON-encoded data and stores the result in the value pointed to by v"), the
`v` must be a pointer so that the results of the operation are available to the caller. However, because `v` is declared
as an `interface{}`, the compiler allows non-pointer values to be passed to the function and the failure is not detected
until runtime.

`outparamcheck` allows these classes of checks to be performed using static analysis. By default, this tool checks the
calls to `encoding/json.Unmarshal`, `encoding/safejson.Unmarshal` and `gopkg.in/yaml.v2.Unmarshal`. It is possible to
use a configuration file to add to the set of functions that are checked.

Install
=======

```
go get -u github.com/palantir/checks/outparamcheck
```

Usage
=====

Run `outparamcheck` with the default checks on all packages within the current directory:

```
./outparamcheck ./...
```

Configuration
=============

Additional checks can be configured using JSON. The JSON can be provided to the check directly as a parameter or by
specifying the path to a file that contains the configuration. The tool accepts a single JSON map where the keys are the
name of the function to be checked and the values are an array that specifies the parameter indices of the "out"
parameter (the parameter that must be a pointer). For example, in order to check that the first (index 0) parameter of
the `github.com/palantir/example/config.Load` function is an output parameter, the JSON would be the following:

```json
{
    "github.com/palantir/example/config.Load": [0]
}
```

The configuration is provided to the tool using the `-config` flag. The value for the flag is treated as a literal JSON
string unless it starts with the `@` character, in which case it is interpreted as the path to a JSON file. The checks
that are specified in the configuration are run in addition to the built-in checks. It is not possible to override or
ignore the built-in checks.

Example invocation configured using JSON directly:

```
./outparamcheck -config '{"github.com/palantir/example/config.Load":[0]}' ./...
```

Example invocation using JSON specified in the file `config.json`:

```
./outparamcheck -config @config.json ./...
```
