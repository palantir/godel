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
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.1.0/godel-2.1.0.tgz...
 0 B / 9.33 MiB    0.00% 376.67 KiB / 9.33 MiB    3.94% 4s 1.05 MiB / 9.33 MiB   11.23% 3s 1.89 MiB / 9.33 MiB   20.26% 2s 2.76 MiB / 9.33 MiB   29.59% 1s 3.56 MiB / 9.33 MiB   38.20% 1s 4.26 MiB / 9.33 MiB   45.61% 1s 5.17 MiB / 9.33 MiB   55.40% 1s 6.08 MiB / 9.33 MiB   65.14% 7.14 MiB / 9.33 MiB   76.50% 7.91 MiB / 9.33 MiB   84.83% 8.69 MiB / 9.33 MiB   93.11% 9.33 MiB / 9.33 MiB  100.00% 2s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.2.0/dist-plugin-1.2.0-linux-amd64.tgz...
 0 B / 4.74 MiB    0.00% 566.07 KiB / 4.74 MiB   11.67% 1s 1.80 MiB / 4.74 MiB   38.08% 2.78 MiB / 4.74 MiB   58.75% 3.62 MiB / 4.74 MiB   76.54% 4.66 MiB / 4.74 MiB   98.36% 4.74 MiB / 4.74 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0/format-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 869.70 KiB / 3.32 MiB   25.59% 1.72 MiB / 3.32 MiB   51.80% 2.52 MiB / 3.32 MiB   76.02% 3.05 MiB / 3.32 MiB   91.93% 3.32 MiB / 3.32 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0/ptimports-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 920.21 KiB / 3.60 MiB   24.95% 1.56 MiB / 3.60 MiB   43.39% 2.26 MiB / 3.60 MiB   62.69% 3.14 MiB / 3.60 MiB   87.17% 3.60 MiB / 3.60 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0/goland-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 498.90 KiB / 3.09 MiB   15.78% 1s 1.14 MiB / 3.09 MiB   36.92% 1.47 MiB / 3.09 MiB   47.48% 1.93 MiB / 3.09 MiB   62.45% 2.68 MiB / 3.09 MiB   86.73% 3.09 MiB / 3.09 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0/check-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.52 MiB    0.00% 527.76 KiB / 3.52 MiB   14.62% 1s 1.29 MiB / 3.52 MiB   36.66% 2.06 MiB / 3.52 MiB   58.58% 2.80 MiB / 3.52 MiB   79.40% 3.52 MiB / 3.52 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0/compiles-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 817.41 KiB / 3.71 MiB   21.53% 1.52 MiB / 3.71 MiB   41.02% 2.16 MiB / 3.71 MiB   58.20% 2.91 MiB / 3.71 MiB   78.42% 3.65 MiB / 3.71 MiB   98.53% 3.71 MiB / 3.71 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0/deadcode-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 498.71 KiB / 3.73 MiB   13.06% 1s 1.18 MiB / 3.73 MiB   31.59% 1.97 MiB / 3.73 MiB   52.73% 2.97 MiB / 3.73 MiB   79.70% 3.73 MiB / 3.73 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0/errcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 777.50 KiB / 3.81 MiB   19.91% 1.45 MiB / 3.81 MiB   38.14% 2.15 MiB / 3.81 MiB   56.26% 2.79 MiB / 3.81 MiB   73.07% 3.36 MiB / 3.81 MiB   88.04% 3.81 MiB / 3.81 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0/extimport-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 746.61 KiB / 3.36 MiB   21.67% 1.50 MiB / 3.36 MiB   44.65% 2.15 MiB / 3.36 MiB   64.04% 2.67 MiB / 3.36 MiB   79.39% 3.31 MiB / 3.36 MiB   98.44% 3.36 MiB / 3.36 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0/golint-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 753.36 KiB / 3.88 MiB   18.97% 1.33 MiB / 3.88 MiB   34.40% 2.00 MiB / 3.88 MiB   51.63% 2.79 MiB / 3.88 MiB   71.96% 3.61 MiB / 3.88 MiB   92.99% 3.88 MiB / 3.88 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0/govet-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.16 MiB    0.00% 595.39 KiB / 3.16 MiB   18.38% 1.17 MiB / 3.16 MiB   36.91% 1.82 MiB / 3.16 MiB   57.53% 2.45 MiB / 3.16 MiB   77.30% 3.16 MiB / 3.16 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0/importalias-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 809.84 KiB / 3.38 MiB   23.42% 1.47 MiB / 3.38 MiB   43.54% 2.16 MiB / 3.38 MiB   64.01% 3.05 MiB / 3.38 MiB   90.23% 3.38 MiB / 3.38 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0/ineffassign-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 571.91 KiB / 3.35 MiB   16.67% 1s 1.20 MiB / 3.35 MiB   35.67% 1.75 MiB / 3.35 MiB   52.37% 2.49 MiB / 3.35 MiB   74.27% 3.13 MiB / 3.35 MiB   93.28% 3.35 MiB / 3.35 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0/novendor-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.41 MiB    0.00% 532.09 KiB / 3.41 MiB   15.22% 1s 1.07 MiB / 3.41 MiB   31.48% 1.73 MiB / 3.41 MiB   50.58% 2.30 MiB / 3.41 MiB   67.30% 2.81 MiB / 3.41 MiB   82.43% 3.30 MiB / 3.41 MiB   96.76% 3.41 MiB / 3.41 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0/outparamcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 269.34 KiB / 3.84 MiB    6.85% 2s 694.32 KiB / 3.84 MiB   17.67% 1s 1.25 MiB / 3.84 MiB   32.55% 1s 2.02 MiB / 3.84 MiB   52.69% 2.48 MiB / 3.84 MiB   64.73% 3.10 MiB / 3.84 MiB   80.72% 3.78 MiB / 3.84 MiB   98.43% 3.84 MiB / 3.84 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0/unconvert-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 698.51 KiB / 3.91 MiB   17.46% 1.28 MiB / 3.91 MiB   32.76% 1.92 MiB / 3.91 MiB   49.06% 2.67 MiB / 3.91 MiB   68.25% 3.43 MiB / 3.91 MiB   87.73% 3.91 MiB / 3.91 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0/varcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 591.74 KiB / 3.75 MiB   15.42% 1s 1.35 MiB / 3.75 MiB   36.16% 2.37 MiB / 3.75 MiB   63.31% 3.12 MiB / 3.75 MiB   83.32% 3.75 MiB / 3.75 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0/license-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 641.79 KiB / 3.30 MiB   18.97% 1.39 MiB / 3.30 MiB   42.00% 2.11 MiB / 3.30 MiB   63.86% 2.76 MiB / 3.30 MiB   83.61% 3.30 MiB / 3.30 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0/test-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 765.78 KiB / 3.60 MiB   20.78% 1.25 MiB / 3.60 MiB   34.82% 1.89 MiB / 3.60 MiB   52.51% 2.71 MiB / 3.60 MiB   75.18% 3.47 MiB / 3.60 MiB   96.33% 3.60 MiB / 3.60 MiB  100.00% 1s
godel version 2.1.0
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.1.0/godel-2.1.0.tgz
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
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.1.0/godel-2.1.0.tgz
distributionSHA256=a1c33e701f18411f72a8b81ba148ec26e2cb0ef5a18ed6d49fc7cc3149acab28
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master dea81f1] Add godel to project
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
 15 9552k   15 1502k    0     0  1439k      0  0:00:06  0:00:01  0:00:05 1439k 57 9552k   57 5538k    0     0  2709k      0  0:00:03  0:00:02  0:00:01 4036k100 9552k  100 9552k    0     0  3434k      0  0:00:02  0:00:02 --:--:-- 4632k
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
godel version 2.1.0
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.1.0/godel-2.1.0.tgz
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
godel version 2.1.0
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
 33 9552k   33 3172k    0     0  2638k      0  0:00:03  0:00:01  0:00:02 2638k 75 9552k   75 7224k    0     0  3280k      0  0:00:02  0:00:02 --:--:-- 4052k100 9552k  100 9552k    0     0  3208k      0  0:00:02  0:00:02 --:--:-- 3594k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.1.0.tgz)= a1c33e701f18411f72a8b81ba148ec26e2cb0ef5a18ed6d49fc7cc3149acab28
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
a1c33e701f18411f72a8b81ba148ec26e2cb0ef5a18ed6d49fc7cc3149acab28  download/godel-2.1.0.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
