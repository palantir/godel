Summary
-------
`./godelw dist` builds distributions for the products in the project based on the dist configuration.

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

Dist
----

Now that we have created a product and defined a build configuration for it, we can move to defining how the
distribution for the product is created. At the bare minimum, most hosting services typically require a product to be
packaged as a `tgz` (or some other archive format). Additionally, some products may want the distribution to contain
artifacts other than the binary (such as documentation or resources).

Start by running `./godelw dist` in the project to observe the default behavior:

```
➜ ./godelw dist
Creating distribution for echgo at /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-darwin-amd64.tgz, /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-linux-amd64.tgz
Finished creating distribution for echgo
```

The default dist settings creates a tgz distribution for each `bin` output for its declared OS/architecture pairs (or
the OS/architecture of the host platform if none are specified) that contains only the binary for the target platform.
In this case, because the build configuration specified that the product should be built for `darwin-amd64` and
`linux-amd64`, the task created distribution artifacts for those two targets. The example above performed only the
distribution task because the binaries for the products were already present and up-to-date. If this were not the case,
the `dist` task would build the required binaries before running.

Similarly to the `build` command, `dist` writes its output to the `dist` directory by default (the output directory can
be configured using the `output-dir` property). The `./godelw clean` command will remove any outputs created by the
`dist` task. To ensure that distribution contents are not tracked in git, add `/dist/` to the `.gitignore` file:

```
➜ echo '/dist/' >> .gitignore
➜ git add .gitignore
➜ git commit -m "Update .gitignore to ignore dist directory"
[master 0b66d9a] Update .gitignore to ignore dist directory
 1 file changed, 1 insertion(+)
➜ git status
On branch master
nothing to commit, working directory clean
```

For this product, the default distribution mechanism is suitable. However, if we were to want to either choose a
different type of distribution or customize the parameters of this distribution, we would specify the `dist`
configuration for the product explicitly. Update the `dist.yml` to explicitly configure the dist parameters of the
product:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      dist-type:
        type: os-arch-bin' > godel/config/dist.yml
```

Commit this update:

```
➜ git add godel/config/dist.yml
➜ git commit -m "Specify dist configuration"
[master 55182ff] Specify dist configuration
 1 file changed, 3 insertions(+)
```

On its own, this functionality may not seem very spectacular. However, these distribution artifacts can be used as
inputs to other tasks such as `publish` and `docker`. Furthermore, for more complicated distributions, it can be useful
to have the logic for creating distributions centrally managed in the configuration.

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

([Link](https://github.com/nmiyake/echgo/tree/55182ff79dd28048782fb240920d6f2d90b453da))

Tutorial next step
------------------

[Publish](https://github.com/palantir/godel/wiki/Publish)

More
----

### Create dists for specific products

By default, `./godelw dist` will create distributions for all of the products defined for a project. Product names can
be specified as arguments to create distributions for only those products. For example, if the `echgo` project defined
multiple products, you could specify that you want to only create the distribution for `echgo` by running the following:

```
➜ ./godelw dist echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff/darwin-amd64/echgo
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff/linux-amd64/echgo
Finished building echgo for linux-amd64 (0.373s)
Finished building echgo for darwin-amd64 (0.375s)
Creating distribution for echgo at /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-2-g55182ff-darwin-amd64.tgz, /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-2-g55182ff-linux-amd64.tgz
Finished creating distribution for echgo
```

### Use a dist type that packages all binaries together

The `os-arch-bin` distribution type creates a separate distribution for each OS/architecture combination. However, in
some instances, it may be desirable to have a single distribution that contains the binaries for all of the target
platforms. The `bin` type can be used to do this.

Update the configuration as follows:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      dist-type:
        type: bin' > godel/config/dist.yml
```

Clean the previous output and run `./godelw dist` to generate a distribution using the new configuration:

```
➜ ./godelw clean
➜ ./godelw dist
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/linux-amd64/echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/darwin-amd64/echgo
Finished building echgo for linux-amd64 (0.366s)
Finished building echgo for darwin-amd64 (0.371s)
Creating distribution for echgo at /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-2-g55182ff.dirty.tgz
Finished creating distribution for echgo
```

Examine the contents of the `dist` directory:

```
➜ tree dist
dist
├── echgo-0.0.1-2-g55182ff.dirty
│   └── bin
│       ├── darwin-amd64
│       │   └── echgo
│       └── linux-amd64
│           └── echgo
└── echgo-0.0.1-2-g55182ff.dirty.tgz

4 directories, 3 files
```

The distribution consists of a product directory that contains a `bin` directory that has a directory for each target
OS/architecture that contains the executable that was built for that target. The `tgz` file is an archive that contains
the top-level directory.

Revert these changes by running the following:

```
➜ git checkout -- godel/config/dist.yml
```

### Specify `input-dir` to copy the contents of a local directory into the distribution

Distributions may want to include static resources such as documentation, scripts or other resources as part of their
distribution. This can be done by having the resources in a directory in the project and specifying that directory as an
input directory for the distribution.

Run the following command to create a `resources` directory that contains a README to include with the binaries:

```
➜ mkdir -p resources
➜ echo 'echgo is a program that echoes its arguments.' > resources/README.md
```

Run the following to update the dist configuration to copy the contents of `resources` into its distribution directory:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      input-dir: resources
      dist-type:
        type: bin' > godel/config/dist.yml
```

The `input-dir: resources` line configures the task to copy all of the contents of the `resources` directory into the
distribution directory (the value of the `input-dir` parameter is the path the to the input directory relative to the
base directory of the project).

Run the following to clean the `dist` directory and run the `dist` task to generate the distribution outputs:

```
➜ ./godelw clean
➜ ./godelw dist
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/linux-amd64/echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/darwin-amd64/echgo
Finished building echgo for darwin-amd64 (0.352s)
Finished building echgo for linux-amd64 (0.355s)
Creating distribution for echgo at /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-2-g55182ff.dirty.tgz
Finished creating distribution for echgo
```

Verify that the distribution directory contains `README.md`:

```
➜ tree dist
dist
├── echgo-0.0.1-2-g55182ff.dirty
│   ├── README.md
│   └── bin
│       ├── darwin-amd64
│       │   └── echgo
│       └── linux-amd64
│           └── echgo
└── echgo-0.0.1-2-g55182ff.dirty.tgz

4 directories, 4 files
```

The input directory can contain any content (including directories or nested directories).

Revert these changes by running the following:

```
➜ rm -rf resources
➜ git checkout -- godel/config/dist.yml
```

### Specify a script to run arbitrary actions during the distribution step

Distributions may need to perform various actions as part of their distribution process that go beyond requiring static
files. For example, a distribution step may require downloading a file, moving files or directories to specific
locations, computing checksums and writing them to a file, etc. In order to support such scenarios, the `dist` block
allows a distribution script to be specified. The distribution script is run after the basic `dist` actions have run but
before the distribution output directory has been archived.

Run the following to update the dist configuration:

```
➜ echo 'products:
  echgo:
    build:
      main-pkg: .
      version-var: main.version
      os-archs:
        - os: darwin
          arch: amd64
        - os: linux
          arch: amd64
    dist:
      script: |
              echo "Distribution created at $(date)" > $DIST_DIR/timestamp.txt
              mv $DIST_DIR/bin $DIST_DIR/binaries
      dist-type:
        type: bin' > godel/config/dist.yml
```

The specified script writes a file called "timestamp.txt" that contains the output of running `date` and also moves the
renames the `bin` directory to be `binaries` instead. `$DIST_DIR` is an environment variable that is set before the
script is run that contains the absolute path to the directory created for the distribution. Refer to the
[dist config documentation](https://godoc.org/github.com/palantir/godel/apps/distgo/config#Dist) for a full description
of the environment variables that are available to the script.

Remove any previous distribution output and run the `dist` command:

```
➜ ./godelw clean
➜ ./godelw dist
Building echgo for linux-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/linux-amd64/echgo
Building echgo for darwin-amd64 at /Volumes/git/go/src/github.com/nmiyake/echgo/build/0.0.1-2-g55182ff.dirty/darwin-amd64/echgo
Finished building echgo for darwin-amd64 (0.374s)
Finished building echgo for linux-amd64 (0.375s)
Creating distribution for echgo at /Volumes/git/go/src/github.com/nmiyake/echgo/dist/echgo-0.0.1-2-g55182ff.dirty.tgz
Finished creating distribution for echgo
```

Verify that `timestamp.txt` was created and that the distribution directory contains a `binaries` directory rather than
a `bin` directory:

```
➜ tree dist
dist
├── echgo-0.0.1-2-g55182ff.dirty
│   ├── binaries
│   │   ├── darwin-amd64
│   │   │   └── echgo
│   │   └── linux-amd64
│   │       └── echgo
│   └── timestamp.txt
└── echgo-0.0.1-2-g55182ff.dirty.tgz

4 directories, 4 files
```

You can also verify that the contents of `timestamp.txt` is correct:

```
➜ cat dist/echgo-"$(./godelw project-version)"/timestamp.txt
Distribution created at Mon Oct 16 15:40:44 PDT 2017
```

Revert these changes by running the following:

```
➜ git checkout -- godel/config/dist.yml
```
