Summary
-------
`./godelw dist` builds distributions for the products in the project based on the dist configuration.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`
* Project contains `main.go`
* Project contains `.gitignore` that ignores GoLand files
* Project contains `echo/echo.go`, `echo/echo_test.go` and `echo/echoer.go`
* `godel/config/dist.yml` is configured to build `echgo2`
* Project is tagged as 0.0.1

Dist
----
Now that we have created a product and defined a build configuration for it, we can move to defining how the
distribution for the product is created. At the bare minimum, most hosting services typically require a product to be
packaged as a `tgz` (or some other archive format). Additionally, some products may want the distribution to contain
artifacts other than the binary (such as documentation or resources).

Observe the default behavior by removing the configuration in the `godel/config/dist.yml` file and running
`./godelw dist`:

```
➜ echo '' > godel/config/dist-plugin.yml
➜ ./godelw dist
Building echgo2 for linux-amd64 at out/build/echgo2/0.0.1.dirty/linux-amd64/echgo2
Finished building echgo2 for linux-amd64 (0.203s)
Creating distribution for echgo2 at out/dist/echgo2/0.0.1.dirty/os-arch-bin/echgo2-0.0.1.dirty-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
```

The default dist settings creates a tgz distribution for each `bin` output for the OS/architecture of the host platform.
Note that, because the build output for the new version was not present, the build task was run as well. If the required
build output was already present, only the distribution task would have been run.

Similarly to the `build` command, `dist` writes its output to the `out/dist` directory by default (the output directory
can be configured using the `output-dir` property). The `./godelw clean` command will remove any outputs created by the
`dist` task.

Update the `dist-plugin.yml` to explicitly configure the dist parameters of the product:

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
              arch: amd64' > godel/config/dist-plugin.yml
```

Run `./godelw dist` to verify that the distributions are built:

```
➜ ./godelw dist
Building echgo2 for darwin-amd64 at out/build/echgo2/0.0.1.dirty/darwin-amd64/echgo2
Finished building echgo2 for darwin-amd64 (0.207s)
Creating distribution for echgo2 at out/dist/echgo2/0.0.1.dirty/os-arch-bin/echgo2-0.0.1.dirty-darwin-amd64.tgz, out/dist/echgo2/0.0.1.dirty/os-arch-bin/echgo2-0.0.1.dirty-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
```

The `dist` operation will run the `build` operation for inputs that need to be built. In the this run, the `dist`
operation only built the output for `linux-amd64` because the previous step in the tutorial (in which we ran `dist` with
an empty `dist.yml` to observe the default behavior) generated the `darwin-amd64` binary, and that output is still
considered up-to-date.

Commit this update:

```
➜ git add godel/config/dist-plugin.yml
➜ git commit -m "Specify dist configuration"
[master 2dcd00b] Specify dist configuration
 1 file changed, 9 insertions(+)
```

On its own, this functionality may not seem very spectacular. However, these distribution artifacts can be used as
inputs to other tasks such as `publish` and `docker`. Furthermore, for more complicated distributions, it can be useful
to have the logic for creating distributions centrally managed in the configuration.

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

Tutorial next step
------------------
[Publish](https://github.com/palantir/godel/wiki/Publish)

More
----
### Force dist generation
By default, a dist output will only be generated if it is considered out of date. A dist output is considered out of
date if any of the following is true:
  * Any of the dist output paths do not exist
  * Any build output for the product or its dependencies is newer than the modification date of the oldest dist output
  * The "godel/config/dist.yml" file was modified at or after the modification date of the oldest dist output

Run `./godelw dist` to generate the dist artifacts. This run will build and dist because the commit is new:

```
➜ ./godelw dist
Building echgo2 for darwin-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b/darwin-amd64/echgo2
Building echgo2 for linux-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b/linux-amd64/echgo2
Finished building echgo2 for darwin-amd64 (0.251s)
Finished building echgo2 for linux-amd64 (0.254s)
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-darwin-amd64.tgz, out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
```

Running this same operation again will not do anything because all of the outputs are up-to-date:

```
➜ ./godelw dist
```

The `--force` flag can be used to specify that the dist artifacts should be generated even if they are not considered
out of date:

```
➜ ./godelw dist --force
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-darwin-amd64.tgz, out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
```

### Create specific distributions
By default, `./godelw dist` will create all of the distributions for all of the products defined for a project. However,
a project can define multiple products, and a product may have multiple distribution outputs. It is possible to specify
that only specific distributions should be built.

First, start by modifying `dist-plugin.yml` to add another distribution type for `echgo2`:

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
        os-arch-bin:
          type: os-arch-bin
          config:
            os-archs:
              - os: darwin
                arch: amd64
              - os: linux
                arch: amd64
        bin:
          type: bin' > godel/config/dist-plugin.yml
```

Verify that running `./godelw dist --force` generates both distributions:

```
➜ ./godelw dist --force
Building echgo2 for darwin-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b.dirty/darwin-amd64/echgo2
Building echgo2 for linux-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b.dirty/linux-amd64/echgo2
Finished building echgo2 for darwin-amd64 (0.241s)
Finished building echgo2 for linux-amd64 (0.249s)
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b.dirty/bin/echgo2-0.0.1-1-g2dcd00b.dirty.tgz
Finished creating bin distribution for echgo2
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b.dirty/os-arch-bin/echgo2-0.0.1-1-g2dcd00b.dirty-darwin-amd64.tgz, out/dist/echgo2/0.0.1-1-g2dcd00b.dirty/os-arch-bin/echgo2-0.0.1-1-g2dcd00b.dirty-linux-amd64.tgz
Finished creating os-arch-bin distribution for echgo2
```

Because there is only one product, only the dist outputs for that product are generated. If there were multiple
products, then running `./godelw dist echgo2` would generate all of the dist outputs for the `echgo2` product.

A specific dister for a product can be run using the `<product>.<name>` syntax. For example, to generate just the
"bin" dist, run `./godelw dist --force echgo2.bin`:

```
➜ ./godelw dist --force echgo2.bin
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b.dirty/bin/echgo2-0.0.1-1-g2dcd00b.dirty.tgz
Finished creating bin distribution for echgo2
```

Revert these changes by running the following:

```
➜ ./godelw clean
```

### Specify a script to run arbitrary actions during the distribution step
Distributions may need to perform various actions as part of their distribution process that go beyond requiring static
files. For example, a distribution step may require downloading a file, moving files or directories to specific
locations, computing checksums and writing them to a file, etc. In order to support such scenarios, the `dist` block
allows a distribution script to be specified. The distribution script is run after the dister's dist actions has been
run, but before the dister's archive action is run.

Run the following to create a configuration that writes a `timestamp.txt` file to each dist output directory:

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
        type: bin
        script: |
                #!/usr/bin/env bash
                set -euo pipefail
                creation_date=$(date)
                echo "Distribution created at $(date)" > "$DIST_WORK_DIR/timestamp.txt"' > godel/config/dist-plugin.yml
```

The specified script writes a file called "timestamp.txt" that contains the output of running `date` and writes it to
the distribution output directory. The environment variables such as `DIST_WORK_DIR` are injected by `distgo`. Refer to
the dist config documentation for a full description of the environment variables that are available to the script.

Run the `dist` command:

```
➜ ./godelw dist --force
Building echgo2 for darwin-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b.dirty/darwin-amd64/echgo2
Building echgo2 for linux-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b.dirty/linux-amd64/echgo2
Finished building echgo2 for linux-amd64 (0.239s)
Finished building echgo2 for darwin-amd64 (0.240s)
Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b.dirty/bin/echgo2-0.0.1-1-g2dcd00b.dirty.tgz
Finished creating bin distribution for echgo2
```

Verify that `timestamp.txt` was created in the distribution directory:

```
➜ tree out/dist
out/dist
`-- echgo2
    `-- 0.0.1-1-g2dcd00b.dirty
        `-- bin
            |-- echgo2-0.0.1-1-g2dcd00b.dirty
            |   |-- bin
            |   |   |-- darwin-amd64
            |   |   |   `-- echgo2
            |   |   `-- linux-amd64
            |   |       `-- echgo2
            |   `-- timestamp.txt
            `-- echgo2-0.0.1-1-g2dcd00b.dirty.tgz

7 directories, 4 files
```

Revert these changes by running the following:

```
➜ ./godelw clean
➜ git checkout -- godel/config/dist-plugin.yml
```

### Dry run
The `--dry-run` flag can be used to preview the operations that would be performed by `./godelw dist` without actually
performing them:

```
➜ ./godelw dist --force --dry-run
[DRY RUN] Building echgo2 for linux-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b/linux-amd64/echgo2
[DRY RUN] Building echgo2 for darwin-amd64 at out/build/echgo2/0.0.1-1-g2dcd00b/darwin-amd64/echgo2
[DRY RUN] Run: /usr/local/go/bin/go build -o /go/src/github.com/nmiyake/echgo2/out/build/echgo2/0.0.1-1-g2dcd00b/darwin-amd64/echgo2 -ldflags -X main.version=0.0.1-1-g2dcd00b ./. with additional environment variables [GOOS=darwin GOARCH=amd64]
[DRY RUN] Run: /usr/local/go/bin/go build -o /go/src/github.com/nmiyake/echgo2/out/build/echgo2/0.0.1-1-g2dcd00b/linux-amd64/echgo2 -ldflags -X main.version=0.0.1-1-g2dcd00b ./. with additional environment variables [GOOS=linux GOARCH=amd64]
[DRY RUN] Finished building echgo2 for linux-amd64 (0.001s)
[DRY RUN] Finished building echgo2 for darwin-amd64 (0.001s)
[DRY RUN] Creating distribution for echgo2 at out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-darwin-amd64.tgz, out/dist/echgo2/0.0.1-1-g2dcd00b/os-arch-bin/echgo2-0.0.1-1-g2dcd00b-linux-amd64.tgz
[DRY RUN] Finished creating os-arch-bin distribution for echgo2
```

### Add disters
The `os-arch-bin`, `bin` and `manual` dister types are built-in as part of the distgo plugin. However, it is possible to
define and add custom disters as assets.

For example, consider a fictional dister asset that generates RPM distributions with the locator
"com.palantir.godel-distgo-asset-dist-rpm:dist-rpm-asset:1.0.0". The following configuration in `godel/config/godel.yml`
would add this dister:

```yaml
default-tasks:
  resolvers:
    - https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz
  tasks:
    com.palantir.distgo:dist-plugin:
      assets:
        - locator:
            id: "com.palantir.godel-distgo-asset-dist-rpm:dist-rpm-asset:1.0.0"
```
