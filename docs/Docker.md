Summary
-------
`./godelw docker` provides commands for building and publishing Docker images.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2

Build Docker Images
-------------------
Our project has been configured to build and publish its binaries. We now want to build and publish a Docker image that
contains the output of the distribution for our product. gödel provides a `docker build` command that can be used to
build Docker images with the outputs of `build` and `dist` tasks (and can use the output of other `docker` tasks to use
as a base image).

The `docker build` task requires a directory within the project that is used as the Docker build directory and a
Dockerfile used to build the Docker image (if a Dockerfile is not specified in configuration, it defaults to
"Dockerfile" within the build context directory). When run, the task ensures that all of the build and dist inputs
required by the Docker task have been created (and runs the `build` and `dist` tasks as necessary if up-to-date outputs
are not available) and then hard-links those outputs into the context directory. The paths to all of these dependencies
are provided as template functions.

Start by creating the Docker build context directory:

```
➜ mkdir dockerctx
```

Run the following to define a Dockerfile:

```
➜ echo 'FROM alpine:3.7

COPY {{InputBuildArtifact Product "linux-amd64"}} /usr/local/bin/

ENTRYPOINT [ "/usr/local/bin/echgo2" ]' > dockerctx/Dockerfile
```

For the most part, this is a standard Dockerfile. The custom portion is `{{InputBuildArtifact Product "linux-amd64"}}`.
This renders to the path of the build output of the `echgo2` product built for the `linux-amd64` OS/architecture. The
example above uses the build artifact (which is the generated Go binary). The `InputDistArtifacts` function can be used
get the paths to the distribution outputs for a product (useful if the `dist` task for a product packages resources or
other files required by the program) and the `Tags` function can be used to retrieve the Docker tags for a product
(useful if an image that is being built by a project needs to use another image created by the product as its base
image). See the "More" section below for details on these use cases.

Now that the context directory and Dockerfile exists, update `dist-plugin.yml` to specify that a Docker image should be
built:

```
➜ echo 'products:
  echgo2:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      disters:
        type: os-arch-bin
        config:
          os-archs:
            - os: darwin
              arch: amd64
            - os: linux
              arch: amd64
    docker:
      docker-builders:
        echgo2:
          type: default
          context-dir: dockerctx
          input-products-dir: inputs
          tag-templates:
            - "{{Repository}}echgo2:{{Version}}"
            - "{{Repository}}echgo2:latest"' > godel/config/dist-plugin.yml
```

The `tag-templates` configuration specifies the tags that should be applied to the Docker image that is built.
`{{Product}}`, `{{Version}}` and `{{Repository}}` are special values that can be used that will be rendered when the
`docker build` task is run. The value of `{{Repository}}` can be specified as configuration (the `repository` field in
the `docker` configuration block) or as a flag when invoking the `docker build` or `docker push` tasks.

Run the task in dry run mode to verify that it will invoke the proper command:

```
➜ ./godelw docker build --dry-run
[DRY RUN] Creating distribution for echgo2 at out/dist/echgo2/0.0.2.dirty/os-arch-bin/echgo2-0.0.2.dirty-darwin-amd64.tgz, out/dist/echgo2/0.0.2.dirty/os-arch-bin/echgo2-0.0.2.dirty-linux-amd64.tgz
[DRY RUN] Finished creating os-arch-bin distribution for echgo2
[DRY RUN] Running Docker build for configuration echgo2 of product echgo2...
[DRY RUN] Run [docker build --file /go/src/github.com/nmiyake/echgo2/dockerctx/Dockerfile -t echgo2:0.0.2.dirty -t echgo2:latest /go/src/github.com/nmiyake/echgo2/dockerctx]
```

Note that specifying the repository using the `--repository` flag updates the tag:

```
➜ ./godelw docker build --dry-run --repository=myregistryhost:5000/
[DRY RUN] Creating distribution for echgo2 at out/dist/echgo2/0.0.2.dirty/os-arch-bin/echgo2-0.0.2.dirty-darwin-amd64.tgz, out/dist/echgo2/0.0.2.dirty/os-arch-bin/echgo2-0.0.2.dirty-linux-amd64.tgz
[DRY RUN] Finished creating os-arch-bin distribution for echgo2
[DRY RUN] Running Docker build for configuration echgo2 of product echgo2...
[DRY RUN] Run [docker build --file /go/src/github.com/nmiyake/echgo2/dockerctx/Dockerfile -t myregistryhost:5000/echgo2:0.0.2.dirty -t myregistryhost:5000/echgo2:latest /go/src/github.com/nmiyake/echgo2/dockerctx]
```

The following is example output of a successful build:

```
➜ ./godelw docker build
Creating distribution for echgo2 at /Volumes/git/go/src/github.com/nmiyake/echgo2/out/dist/echgo2/0.0.2-dirty/os-arch-bin/echgo2-0.0.2-dirty-darwin-amd64.tgz, /Volumes/git/go/src/github.com/nmiyake/echgo2/out/dist/echgo2/0.0.2-dirty/os-arch-bin/echgo2-0.0.2-dirty-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
Running Docker build for configuration echgo2 of product echgo2...
```

The default behavior of the `docker build` task does not include the output of `docker build` itself. The `--verbose`
flag can be used to include all output:

```
➜ ./godelw docker build --verbose
Running Docker build for configuration echgo2 of product echgo2...
Sending build context to Docker daemon  6.021MB
Step 1/3 : FROM alpine:3.7
 ---> 3fd9065eaf02
Step 2/3 : COPY inputs/echgo2/build/linux-amd64/echgo2 /usr/local/bin/
 ---> 57767e5b5a61
Step 3/3 : ENTRYPOINT [ "/usr/local/bin/echgo2" ]
 ---> Running in 32700255ed57
Removing intermediate container 32700255ed57
 ---> 2538ee323de3
Successfully built 2538ee323de3
Successfully tagged echgo2:0.0.2-dirty
Successfully tagged echgo2:latest
```

Verify that the Docker image can be run and works as expected:

```
➜ docker run echgo2:latest 'Hello, world!'
Hello, world!
```

In order to make sure that the hard-linked artifacts are not checked in, add a `.gitignore` file:

```
➜ echo 'inputs/' > dockerctx/.gitignore
```

In this example, the `.gitignore` file is configured to ignore the `inputs/` directory because we configured the Docker
task to use that directory as the directory into which the build and dist artifacts should be hard-linked.

Commit the changes that were made to the repository:

```
➜ git add dockerctx godel/config/dist-plugin.yml
➜ git commit -m "Add Docker build configuration"
[master 4e54a51] Add Docker build configuration
 3 files changed, 15 insertions(+)
 create mode 100644 dockerctx/.gitignore
 create mode 100644 dockerctx/Dockerfile
```

Push Docker Images
------------------
The `docker push` task can be used to push Docker images built using the `docker build` task. Run the following command
to observe how the Docker image would be pushed:

```
➜ ./godelw docker push --dry-run
[DRY RUN] Running Docker push for configuration echgo2 of product echgo2...
[DRY RUN] Run [docker push echgo2:0.0.2-1-g4e54a51]
[DRY RUN] Run [docker push echgo2:latest]
```

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist-plugin.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1
* `godel/config/dist-plugin.yml` is configured to create distributions for `echgo`
* Project is tagged as 0.0.2
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product

Tutorial next step
------------------
[Generate license headers](https://github.com/palantir/godel/wiki/License-headers)

More
----
### Use build or dist artifacts of dependent products
By default, the `docker build` task only hard-links in the build and dist artifacts for the current product. However, if
the current product depends on other products, the Docker task configuration can specify that build and/or dist outputs
of any of those products should be hard-linked in to the context directory. The `input-builds` and `input-dists`
configuration values can be used to specify the builds and distributions that should be hard-linked in.
