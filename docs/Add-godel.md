Summary
-------
Install and run `godelinit` to add gödel to a project.

Tutorial start state
--------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository

Add gödel to a project
----------------------
The godelinit program can be used to add gödel to a new project. The godelinit program determines the latest release of
gödel and installs it into the current directory, downloading it if needed.

Install the `godelinit` program using `go get` and run the program to install gödel:

```
➜ go get github.com/palantir/godel/godelinit
➜ godelinit
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0/godel-2.0.0.tgz...
 0 B / 9.33 MiB    0.00% 416.50 KiB / 9.33 MiB    4.36% 4s 1.11 MiB / 9.33 MiB   11.88% 2s 1.79 MiB / 9.33 MiB   19.18% 2s 2.40 MiB / 9.33 MiB   25.77% 2s 3.16 MiB / 9.33 MiB   33.92% 1s 3.94 MiB / 9.33 MiB   42.24% 1s 4.67 MiB / 9.33 MiB   50.11% 1s 5.41 MiB / 9.33 MiB   57.98% 1s 6.01 MiB / 9.33 MiB   64.39% 6.67 MiB / 9.33 MiB   71.51% 7.79 MiB / 9.33 MiB   83.45% 8.60 MiB / 9.33 MiB   92.19% 9.19 MiB / 9.33 MiB   98.48% 9.33 MiB / 9.33 MiB  100.00% 2s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.0.0/dist-plugin-1.0.0-linux-amd64.tgz...
 0 B / 4.73 MiB    0.00% 499.93 KiB / 4.73 MiB   10.31% 1s 512.00 KiB / 4.73 MiB   10.56% 3s 512.00 KiB / 4.73 MiB   10.56% 3s 512.00 KiB / 4.73 MiB   10.56% 2s 702.02 KiB / 4.73 MiB   14.48% 5s 1.47 MiB / 4.73 MiB   31.06% 2s 2.32 MiB / 4.73 MiB   48.94% 1s 2.93 MiB / 4.73 MiB   61.82% 3.54 MiB / 4.73 MiB   74.70% 4.23 MiB / 4.73 MiB   89.30% 4.73 MiB / 4.73 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0/format-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 462.67 KiB / 3.32 MiB   13.61% 1s 1.14 MiB / 3.32 MiB   34.31% 1.50 MiB / 3.32 MiB   45.19% 2.15 MiB / 3.32 MiB   64.90% 2.90 MiB / 3.32 MiB   87.27% 3.32 MiB / 3.32 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0/ptimports-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 512.00 KiB / 3.60 MiB   13.88% 1s 846.89 KiB / 3.60 MiB   22.96% 1s 1.31 MiB / 3.60 MiB   36.42% 1s 1.87 MiB / 3.60 MiB   52.02% 2.55 MiB / 3.60 MiB   70.92% 3.47 MiB / 3.60 MiB   96.46% 3.60 MiB / 3.60 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0/goland-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 680.00 KiB / 3.09 MiB   21.51% 1.25 MiB / 3.09 MiB   40.49% 1.54 MiB / 3.09 MiB   49.87% 1.97 MiB / 3.09 MiB   63.84% 2.97 MiB / 3.09 MiB   96.07% 3.09 MiB / 3.09 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0/check-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.52 MiB    0.00% 682.94 KiB / 3.52 MiB   18.92% 1.44 MiB / 3.52 MiB   40.95% 1.82 MiB / 3.52 MiB   51.62% 2.74 MiB / 3.52 MiB   77.75% 3.52 MiB / 3.52 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0/compiles-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 654.02 KiB / 3.71 MiB   17.23% 1.00 MiB / 3.71 MiB   26.98% 1s 1.46 MiB / 3.71 MiB   39.35% 2.00 MiB / 3.71 MiB   53.95% 2.86 MiB / 3.71 MiB   77.04% 3.63 MiB / 3.71 MiB   97.91% 3.71 MiB / 3.71 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0/deadcode-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 366.45 KiB / 3.73 MiB    9.60% 1s 1.00 MiB / 3.73 MiB   26.81% 1s 1.52 MiB / 3.73 MiB   40.63% 2.30 MiB / 3.73 MiB   61.57% 3.14 MiB / 3.73 MiB   84.19% 3.73 MiB / 3.73 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0/errcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 319.11 KiB / 3.81 MiB    8.17% 2s 1.08 MiB / 3.81 MiB   28.36% 1s 1.87 MiB / 3.81 MiB   49.03% 2.55 MiB / 3.81 MiB   66.85% 3.37 MiB / 3.81 MiB   88.24% 3.81 MiB / 3.81 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0/extimport-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 638.23 KiB / 3.36 MiB   18.53% 1.00 MiB / 3.36 MiB   29.73% 1.69 MiB / 3.36 MiB   50.13% 2.27 MiB / 3.36 MiB   67.35% 3.12 MiB / 3.36 MiB   92.85% 3.36 MiB / 3.36 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0/golint-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 512.00 KiB / 3.88 MiB   12.90% 1s 1.00 MiB / 3.88 MiB   25.79% 1s 1.50 MiB / 3.88 MiB   38.69% 2.04 MiB / 3.88 MiB   52.59% 2.95 MiB / 3.88 MiB   76.10% 3.79 MiB / 3.88 MiB   97.63% 3.88 MiB / 3.88 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0/govet-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.16 MiB    0.00% 407.77 KiB / 3.16 MiB   12.59% 1s 1.30 MiB / 3.16 MiB   41.08% 2.02 MiB / 3.16 MiB   63.92% 2.84 MiB / 3.16 MiB   89.69% 3.16 MiB / 3.16 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0/importalias-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 678.45 KiB / 3.38 MiB   19.62% 1.25 MiB / 3.38 MiB   37.01% 1.90 MiB / 3.38 MiB   56.21% 2.73 MiB / 3.38 MiB   80.87% 3.38 MiB / 3.38 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0/ineffassign-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 646.45 KiB / 3.35 MiB   18.84% 1.12 MiB / 3.35 MiB   33.48% 1.84 MiB / 3.35 MiB   54.99% 2.61 MiB / 3.35 MiB   77.83% 3.35 MiB / 3.35 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0/novendor-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.41 MiB    0.00% 256.00 KiB / 3.41 MiB    7.32% 2s 256.00 KiB / 3.41 MiB    7.32% 2s 382.23 KiB / 3.41 MiB   10.93% 4s 1.14 MiB / 3.41 MiB   33.35% 1s 2.08 MiB / 3.41 MiB   60.83% 2.90 MiB / 3.41 MiB   85.05% 3.41 MiB / 3.41 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0/outparamcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 519.35 KiB / 3.84 MiB   13.22% 1s 1.25 MiB / 3.84 MiB   32.65% 2.06 MiB / 3.84 MiB   53.60% 3.09 MiB / 3.84 MiB   80.52% 3.84 MiB / 3.84 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0/unconvert-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 698.61 KiB / 3.91 MiB   17.46% 1.25 MiB / 3.91 MiB   31.99% 1.75 MiB / 3.91 MiB   44.79% 2.25 MiB / 3.91 MiB   57.59% 2.87 MiB / 3.91 MiB   73.54% 3.37 MiB / 3.91 MiB   86.15% 3.58 MiB / 3.91 MiB   91.55% 3.71 MiB / 3.91 MiB   94.98% 3.79 MiB / 3.91 MiB   97.07% 3.91 MiB / 3.91 MiB   99.96% 3.91 MiB / 3.91 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0/varcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 96.00 KiB / 3.75 MiB    2.50% 7s 240.00 KiB / 3.75 MiB    6.26% 6s 319.33 KiB / 3.75 MiB    8.32% 6s 438.75 KiB / 3.75 MiB   11.44% 6s 605.76 KiB / 3.75 MiB   15.79% 5s 1.14 MiB / 3.75 MiB   30.30% 2s 1.73 MiB / 3.75 MiB   46.26% 1s 2.22 MiB / 3.75 MiB   59.32% 1s 2.74 MiB / 3.75 MiB   73.11% 3.54 MiB / 3.75 MiB   94.46% 3.75 MiB / 3.75 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0/license-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 494.23 KiB / 3.30 MiB   14.61% 1s 512.00 KiB / 3.30 MiB   15.13% 2s 1.16 MiB / 3.30 MiB   35.07% 1s 1.39 MiB / 3.30 MiB   42.00% 1s 1.69 MiB / 3.30 MiB   51.05% 2.11 MiB / 3.30 MiB   63.86% 2.45 MiB / 3.30 MiB   74.09% 2.83 MiB / 3.30 MiB   85.60% 3.30 MiB / 3.30 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0/test-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 391.77 KiB / 3.60 MiB   10.63% 1s 813.50 KiB / 3.60 MiB   22.08% 1s 1.30 MiB / 3.60 MiB   36.11% 1s 1.80 MiB / 3.60 MiB   50.03% 2.41 MiB / 3.60 MiB   67.09% 3.00 MiB / 3.60 MiB   83.27% 3.60 MiB / 3.60 MiB  100.00% 1s
godel version 2.0.0
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0/godel-2.0.0.tgz
distributionSHA256=
```

If we are concerned about integrity, we can specify the expected checksum for the package when using godelinit. The
expected checksums for gödel releases can be found on its [Bintray](https://bintray.com/palantir/releases/godel) page.
Install the distribution using the checksum:

```
➜ godelinit --checksum ${GODEL_CHECKSUM}
```

If the installation succeeds with a checksum specified, it is set in the properties:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0/godel-2.0.0.tgz
distributionSHA256=282bdfd9a650c12d46670460ff58db765958a64e465fbee7ec459305f1325f13
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master 34e2e5b] Add godel to project
 8 files changed, 243 insertions(+)
 create mode 100644 godel/config/check-plugin.yml
 create mode 100644 godel/config/dist-plugin.yml
 create mode 100644 godel/config/format-plugin.yml
 create mode 100644 godel/config/godel.properties
 create mode 100644 godel/config/godel.yml
 create mode 100644 godel/config/license-plugin.yml
 create mode 100644 godel/config/test-plugin.yml
 create mode 100755 godelw
```

gödel has now been added to the project and is ready to use.

Tutorial end state
------------------
* `${GOPATH}/src/${PROJECT_PATH}` exists, is the working directory and is initialized as a Git repository
* Project contains `godel` and `godelw`

Tutorial next step
------------------
[Add Git hooks to enforce formatting](https://github.com/palantir/godel/wiki/Add-git-hooks)

More
----
### Download and install gödel manually
It is possible to add gödel to a project without using godelinit by downloading and installing it manually.

Download the distribution into a temporary directory and expand it:

```
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
  0 9552k    0 55851    0     0  72436      0  0:02:15 --:--:--  0:02:15 72436 48 9552k   48 4619k    0     0  2620k      0  0:00:03  0:00:01  0:00:02 4597k 96 9552k   96 9256k    0     0  3349k      0  0:00:02  0:00:02 --:--:-- 4617k100 9552k  100 9552k    0     0  3385k      0  0:00:02  0:00:02 --:--:-- 4633k
➜ tar -xf download/godel-"${GODEL_VERSION}".tgz -C download
```

Copy the contents of the wrapper directory (which contains godelw and godel) to the project directory:

```
➜ cp -r download/godel-"${GODEL_VERSION}"/wrapper/* .
```

Run `./godelw version` to verify that gödel was installed correctly. This command will download the distribution if the
distribution has not previously been installed locally. However, if you have followed the tutorial steps, this version
has already been downloaded so it will not re-download:

```
➜ ./godelw version
godel version 2.0.0
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0/godel-2.0.0.tgz
distributionSHA256=
```

For completeness, set the checksum. The expected checksums for releases are listed on the download page at
https://bintray.com/palantir/releases/godel.

```
➜ DIST_URL='distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/GODEL_VERSION/godel-GODEL_VERSION.tgz
distributionSHA256=GODEL_CHECKSUM'; DIST_URL="${DIST_URL//GODEL_VERSION/$GODEL_VERSION}"; DIST_URL="${DIST_URL//GODEL_CHECKSUM/$GODEL_CHECKSUM}"
➜ echo "${DIST_URL}" > godel/config/godel.properties
➜ unset DIST_URL
```

Now that gödel has been added to the project, remove the temporary directory and unset the version variable:

```
➜ rm -rf download
```

### Copying gödel from an existing project
If you have local projects that already use gödel, you can add gödel to a another project by copying the `godelw` and
`godel/config/godel.properties` files from the project that already has gödel.

For example, assume that `{$GOPATH}/src/${PROJECT_PATH}` exists in the current state of the tutorial and we want
to create another project at `${GOPATH}/src/github.com/nmiyake/sample` that also uses gödel. This can be done as follows:

```
➜ mkdir -p ${GOPATH}/src/github.com/nmiyake/sample && cd $_
➜ cp ${GOPATH}/src/${PROJECT_PATH}/godelw .
➜ mkdir -p godel/config
➜ cp ${GOPATH}/src/${PROJECT_PATH}/godel/config/godel.properties godel/config/
➜ ./godelw update --sync
```

Verify that invoking `./godelw` works and that the `godel/config` directory has been populated with the default
configuration files:

```
➜ ./godelw version
godel version 2.0.0
➜ ls godel/config
check-plugin.yml
dist-plugin.yml
format-plugin.yml
godel.properties
godel.yml
license-plugin.yml
test-plugin.yml
```

Restore the workspace to the original state by setting the working directory back to the `echgo2` project and removing
the sample project:

```
➜ cd ${GOPATH}/src/${PROJECT_PATH}
➜ rm -rf ${GOPATH}/src/github.com/nmiyake/sample
```

The steps above take the approach of copying the `godelw` file and `godel/config/godel.properties` file and running
`update` to ensure that all of the configuration is in its default state -- if the entire `godel/config` directory was
copied from the other project, then the current project would copy the other project's configuration as well, which is
most likely not the correct thing to do.

The `update` task updates gödel to the latest release version. Flags can be used to configure the behavior to perform
different operations (for example, to update gödel to a specific version or to the version specified in
`godel/config/godel.properties`). If the distribution is already downloaded locally in the local global gödel directory
(`~/.godel` by default) and matches the checksum specified in `godel.properties`, the required files are copied from the
local cache. If the distribution does not exist in the local cache or a checksum is not specified in `godel.properties`,
then the distribution is downloaded from the distribution URL.

### Specify the checksum in `godel/config/godel.properties`
If the `godelw` script cannot locate the `godel` binary, it downloads it from the URL specified in
`godel/godel.properties`. If a checksum is specified in `godel/config/godel.properties`, it is used to verify the
integrity and validity of the distribution downloaded from the URL. The checksum is also used by the `update` task -- if
the `update` task for a given version can find a download of the distribution in the downloads directory with a matching
checksum, it will be used to update the version of gödel (otherwise, the `update` task will always download the
distribution from the specified distribution URL).

The checksum is the SHA-256 checksum of the `.tgz` archive of the product. The checksum can be obtained from the
Bintray page for the distribution:

![SHA checksum](images/tutorial/sha_checksum.png)

If a user has obtained a trusted distribution locally, it is also possible to manually compute the checksum.

Download the distribution by running the following:

```
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/${GODEL_VERSION}/godel-${GODEL_VERSION}.tgz" -o download/godel-"${GODEL_VERSION}".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
 14 9552k   14 1362k    0     0  1469k      0  0:00:06 --:--:--  0:00:06 1469k 48 9552k   48 4607k    0     0  2377k      0  0:00:04  0:00:01  0:00:03 3212k 48 9552k   48 4647k    0     0  1577k      0  0:00:06  0:00:02  0:00:04 1628k 48 9552k   48 4675k    0     0  1156k      0  0:00:08  0:00:04  0:00:04 1063k 49 9552k   49 4691k    0     0   904k      0  0:00:10  0:00:05  0:00:05  781k 49 9552k   49 4703k    0     0   766k      0  0:00:12  0:00:06  0:00:06  641k 60 9552k   60 5804k    0     0   838k      0  0:00:11  0:00:06  0:00:05  240k 99 9552k   99 9534k    0     0  1202k      0  0:00:07  0:00:07 --:--:--  980k100 9552k  100 9552k    0     0  1203k      0  0:00:07  0:00:07 --:--:-- 1251k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.0.0.tgz)= 282bdfd9a650c12d46670460ff58db765958a64e465fbee7ec459305f1325f13
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
282bdfd9a650c12d46670460ff58db765958a64e465fbee7ec459305f1325f13  download/godel-2.0.0.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
