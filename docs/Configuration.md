Configuration
=============
gödel projects are configured using YML files in the "godel/config" directory. The YML files in that directory specify
the declarative configuration for various different aspects of gödel projects.

Overview
--------

The following is a list of the supported configuration files. Refer to the struct definitions, tutorial steps or GoDocs
for more information on the configuration values.

Example configuration for these files can usually be found on the tutorial page. Many of the `config.go` files also
have `example_test.go` files in the same directory that contain example YML. Examining the equivalent files in other
projects that use gödel (including gödel itself) is also a good way to get a sense of the configuration.

| File | Tasks | Struct definitions | Tutorial step |
| ---- | ----- | ------------------ | ------------- |
| `check.yml` | `check`, `verify` | [apps/okgo/config/config.go](https://github.com/palantir/godel/blob/master/apps/okgo/config/config.go) | [Check](https://github.com/palantir/godel/wiki/Check) |
| `dist.yml`  | `artifacts`, `build`, `dist`, `products`, `publish`, `run` | [apps/distgo/config/config.go](https://github.com/palantir/godel/blob/master/apps/distgo/config/config.go) | [Dist](https://github.com/palantir/godel/wiki/Dist) |
| `exclude.yml` | All | [config/config.go](https://github.com/palantir/godel/blob/master/config/config.go) | [Exclude](https://github.com/palantir/godel/wiki/Exclude) |
| `format.yml` | `format`, `verify` | [apps/gonform/config/config.go](https://github.com/palantir/godel/blob/master/apps/gonform/config/config.go) | [Format](https://github.com/palantir/godel/wiki/Format) |
| `generate.yml` | `generate`, `verify` | [vendor/github.com/palantir/checks/gogenerate/config/config.go](https://github.com/palantir/godel/blob/master/vendor/github.com/palantir/checks/gogenerate/config/config.go) | [Generate](https://github.com/palantir/godel/wiki/Generate) |
| `imports.yml` | `imports`, `verify` | [vendor/github.com/palantir/checks/gocd/config/config.go](https://github.com/palantir/godel/blob/master/vendor/github.com/palantir/checks/gocd/config/config.go) | N/A |
| `license.yml` | `license`, `verify` | [vendor/github.com/palantir/checks/golicense/config/config.go](https://github.com/palantir/godel/blob/master/vendor/github.com/palantir/checks/golicense/config/config.go) | [License](https://github.com/palantir/godel/wiki/License-headers) |
| `test.yml` | `test`, `verify` | [apps/gunit/config/config.go](https://github.com/palantir/godel/blob/master/apps/gunit/config/config.go) | [Test](https://github.com/palantir/godel/wiki/Test) |

Go toolchain management
-----------------------

Set `manage-gotoolchain: true` in `godel/config/godel.yml` to select an exact Go
toolchain for all subprocesses that inherit the gödel environment, including
plugins and `./godelw exec -- go ...`. The default is false.

```yaml
manage-gotoolchain: true
```

Gödel reads `go.mod` beside the invoked `godelw`, regardless of the working
directory. An explicit `toolchain` directive takes precedence. If it is absent
or is `default`, gödel uses the `go` directive, adding the `go` prefix and
appending `.0` to a bare major.minor version. For example, `go 1.27` selects
`go1.27.0`. Release candidates and beta versions are preserved. Custom
toolchains with a `-suffix` are unsupported.

When enabled, gödel overrides any inherited `GOTOOLCHAIN` with the selected
name, without `+auto`. Specifying `GOTOOLCHAIN` under `environment` is a
startup error, including when its value is empty. The module file must exist,
be valid, and contain a `go` directive. A toolchain older than the `go`
requirement is passed through; Go reports the incompatibility when invoked.

Go 1.21 or later must be available to honor `GOTOOLCHAIN`. Go handles
toolchain lookup, downloading, and verification. Configure `GOPROXY`,
`GOSUMDB`, and related settings through the existing `environment` map as
needed; gödel does not download toolchains itself.

When management is false or unspecified, gödel does not read `go.mod` and
continues applying `environment` as usual, including a configured
`GOTOOLCHAIN`. Explicit false values are preserved so they remain opt-outs
if a future release enables management by default.

This setting covers one module rooted beside the wrapper. Operations spanning
multiple modules and any behavior affected by `go.work` are unsupported.
