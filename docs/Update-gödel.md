Summary
-------
The version of gödel used by a project can be updated by updating the `godel/config/godel.properties` file, running
`./godelw update` and checking in the updated files.

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
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0

([Link](https://github.com/nmiyake/echgo/tree/25d27eb1763e55f228282594691798ca0c2bbe28))

Update gödel
------------

The version of gödel used by a project can be updated by updating the `godel/config/godel.properties` file and running
the `./godelw update` command.

Updating gödel requires knowing the distribution URL for the target version. Although it is optional, it is recommended
to have the SHA-256 checksum of the distribution as well to ensure the integrity of the update.

Information for the latest release can be found on gödel's [Bintray page](https://bintray.com/palantir/releases/godel).
This page also displays the SHA-256 checksum for the distribution:

![SHA checksum](images/tutorial/sha_checksum.png)

The tutorial used version 0.26.0 of gödel, but version 0.27.0 also exists. Update the version of gödel used by the
project by running the following:

```
➜ echo 'distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/0.27.0/godel-0.27.0.tgz
distributionSHA256=0869fc0fb10b4cdd179185c0e59e28a3568c447a2a7ab3d379d6037900a96bf3' > godel/config/godel.properties
➜ ./godelw update
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/0.27.0/godel-0.27.0.tgz...
 10.74 MB / 10.74 MB [======================================================================================] 100.00% 4s
```

The `./godelw update` operation compares the project's current version of gödel with the version specified in
`godel/config/godel.properties`. If the versions do not match, it downloads the version specified in the properties file
and updates the files in the project as necessary.

Updating gödel will typically update the `godelw` wrapper file and may add new configuration files to the `godel/config`
directory. All existing configuration will remain unmodified.

Check in the update by running the following:

```
➜ git add godelw godel
➜ git commit -m "Update godel to 0.27.0"
[master 3825f36] Update godel to 0.27.0
 2 files changed, 5 insertions(+), 5 deletions(-)
```

The `godelw` wrapper ensures that the required version of gödel is present on the system (downloading it if necessary)
on every invocation, so when other developers (or the CI system) checks out a version of the project that updates the
version of gödel that is used, the first `./godelw` command they invoke will ensure that the updated version is
available. This ensures that all gödel operations are always using the correct version of the program.

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
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to 0.27.0

([Link](https://github.com/nmiyake/echgo/tree/3825f36b06ee50703ad10e01068ceb13e7719acd))

Tutorial next step
------------------
[Other commands](https://github.com/palantir/godel/wiki/Other-commands)
