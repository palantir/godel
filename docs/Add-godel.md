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
 0 B / 9.33 MiB    0.00% 33.58 KiB / 9.33 MiB    0.35% 56s 118.58 KiB / 9.33 MiB    1.24% 31s 203.58 KiB / 9.33 MiB    2.13% 27s 407.58 KiB / 9.33 MiB    4.27% 18s 628.58 KiB / 9.33 MiB    6.58% 14s 981.58 KiB / 9.33 MiB   10.28% 10s 1.44 MiB / 9.33 MiB   15.48% 7s 2.02 MiB / 9.33 MiB   21.60% 5s 2.57 MiB / 9.33 MiB   27.58% 4s 3.25 MiB / 9.33 MiB   34.79% 3s 4.02 MiB / 9.33 MiB   43.06% 2s 4.81 MiB / 9.33 MiB   51.52% 2s 5.81 MiB / 9.33 MiB   62.28% 1s 6.52 MiB / 9.33 MiB   69.92% 1s 7.55 MiB / 9.33 MiB   80.97% 8.42 MiB / 9.33 MiB   90.21% 9.19 MiB / 9.33 MiB   98.47% 9.33 MiB / 9.33 MiB  100.00% 3s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.2.0/dist-plugin-1.2.0-linux-amd64.tgz...
 0 B / 4.74 MiB    0.00% 575.77 KiB / 4.74 MiB   11.87% 1s 1.53 MiB / 4.74 MiB   32.33% 2.21 MiB / 4.74 MiB   46.64% 2.90 MiB / 4.74 MiB   61.19% 3.67 MiB / 4.74 MiB   77.53% 4.35 MiB / 4.74 MiB   91.88% 4.74 MiB / 4.74 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0/format-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 750.40 KiB / 3.32 MiB   22.08% 1.63 MiB / 3.32 MiB   49.11% 2.06 MiB / 3.32 MiB   62.21% 2.84 MiB / 3.32 MiB   85.61% 3.32 MiB / 3.32 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0/ptimports-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 749.13 KiB / 3.60 MiB   20.31% 1.44 MiB / 3.60 MiB   39.94% 2.11 MiB / 3.60 MiB   58.49% 2.85 MiB / 3.60 MiB   79.19% 3.53 MiB / 3.60 MiB   98.06% 3.60 MiB / 3.60 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0/goland-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 805.09 KiB / 3.09 MiB   25.47% 1.48 MiB / 3.09 MiB   47.99% 2.17 MiB / 3.09 MiB   70.38% 2.85 MiB / 3.09 MiB   92.39% 3.09 MiB / 3.09 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0/check-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.52 MiB    0.00% 666.94 KiB / 3.52 MiB   18.48% 1.13 MiB / 3.52 MiB   32.03% 1.70 MiB / 3.52 MiB   48.22% 2.39 MiB / 3.52 MiB   67.83% 2.92 MiB / 3.52 MiB   82.93% 3.45 MiB / 3.52 MiB   97.91% 3.52 MiB / 3.52 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0/compiles-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 387.94 KiB / 3.71 MiB   10.22% 1s 916.83 KiB / 3.71 MiB   24.15% 1s 1.49 MiB / 3.71 MiB   40.29% 2.01 MiB / 3.71 MiB   54.22% 2.61 MiB / 3.71 MiB   70.35% 3.17 MiB / 3.71 MiB   85.44% 3.63 MiB / 3.71 MiB   97.91% 3.71 MiB / 3.71 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0/deadcode-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 625.96 KiB / 3.73 MiB   16.39% 1s 1.08 MiB / 3.73 MiB   29.09% 1.85 MiB / 3.73 MiB   49.50% 2.61 MiB / 3.73 MiB   69.91% 3.29 MiB / 3.73 MiB   88.13% 3.73 MiB / 3.73 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0/errcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 805.33 KiB / 3.81 MiB   20.62% 1.49 MiB / 3.81 MiB   39.16% 2.19 MiB / 3.81 MiB   57.39% 2.98 MiB / 3.81 MiB   78.06% 3.70 MiB / 3.81 MiB   96.89% 3.81 MiB / 3.81 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0/extimport-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 1.00 MiB / 3.36 MiB   29.87% 1.89 MiB / 3.36 MiB   56.19% 2.74 MiB / 3.36 MiB   81.59% 3.36 MiB / 3.36 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0/golint-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 1.06 MiB / 3.88 MiB   27.29% 1.94 MiB / 3.88 MiB   50.13% 2.93 MiB / 3.88 MiB   75.66% 3.88 MiB / 3.88 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0/govet-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.16 MiB    0.00% 543.78 KiB / 3.16 MiB   16.78% 861.82 KiB / 3.16 MiB   26.60% 1s 1.16 MiB / 3.16 MiB   36.55% 1s 1.51 MiB / 3.16 MiB   47.72% 1.81 MiB / 3.16 MiB   57.17% 2.19 MiB / 3.16 MiB   69.19% 2.62 MiB / 3.16 MiB   82.94% 3.02 MiB / 3.16 MiB   95.34% 3.16 MiB / 3.16 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0/importalias-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 443.90 KiB / 3.38 MiB   12.84% 1s 917.11 KiB / 3.38 MiB   26.52% 1s 1.33 MiB / 3.38 MiB   39.40% 1.84 MiB / 3.38 MiB   54.35% 2.22 MiB / 3.38 MiB   65.62% 2.68 MiB / 3.38 MiB   79.30% 3.14 MiB / 3.38 MiB   92.99% 3.38 MiB / 3.38 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0/ineffassign-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 442.59 KiB / 3.35 MiB   12.90% 1s 876.13 KiB / 3.35 MiB   25.54% 1s 1.29 MiB / 3.35 MiB   38.52% 1.81 MiB / 3.35 MiB   53.93% 2.30 MiB / 3.35 MiB   68.54% 2.70 MiB / 3.35 MiB   80.71% 3.19 MiB / 3.35 MiB   95.31% 3.35 MiB / 3.35 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0/novendor-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.41 MiB    0.00% 526.02 KiB / 3.41 MiB   15.04% 1s 1.05 MiB / 3.41 MiB   30.63% 1.55 MiB / 3.41 MiB   45.30% 2.09 MiB / 3.41 MiB   61.22% 2.51 MiB / 3.41 MiB   73.62% 2.99 MiB / 3.41 MiB   87.49% 3.41 MiB / 3.41 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0/outparamcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 666.39 KiB / 3.84 MiB   16.96% 1.26 MiB / 3.84 MiB   32.95% 1.85 MiB / 3.84 MiB   48.13% 2.43 MiB / 3.84 MiB   63.42% 2.95 MiB / 3.84 MiB   76.88% 3.42 MiB / 3.84 MiB   89.22% 3.84 MiB / 3.84 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0/unconvert-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 682.61 KiB / 3.91 MiB   17.06% 1.24 MiB / 3.91 MiB   31.67% 1.82 MiB / 3.91 MiB   46.58% 2.53 MiB / 3.91 MiB   64.67% 3.15 MiB / 3.91 MiB   80.67% 3.68 MiB / 3.91 MiB   94.29% 3.91 MiB / 3.91 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0/varcheck-asset-1.0.0-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 665.32 KiB / 3.75 MiB   17.34% 1.30 MiB / 3.75 MiB   34.75% 1.89 MiB / 3.75 MiB   50.41% 2.46 MiB / 3.75 MiB   65.64% 3.14 MiB / 3.75 MiB   83.78% 3.75 MiB / 3.75 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0/license-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 498.64 KiB / 3.30 MiB   14.74% 1s 999.69 KiB / 3.30 MiB   29.54% 1.67 MiB / 3.30 MiB   50.58% 2.23 MiB / 3.30 MiB   67.39% 2.79 MiB / 3.30 MiB   84.31% 3.30 MiB / 3.30 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0/test-plugin-1.0.0-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 916.89 KiB / 3.60 MiB   24.88% 1.63 MiB / 3.60 MiB   45.28% 2.31 MiB / 3.60 MiB   64.17% 3.25 MiB / 3.60 MiB   90.29% 3.60 MiB / 3.60 MiB  100.00% 0s
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
[master 276e30a] Add godel to project
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
  0 9552k    0     0    0     0      0      0 --:--:--  0:00:01 --:--:--     0  7 9552k    7  730k    0     0   369k      0  0:00:25  0:00:01  0:00:24  824k 26 9552k   26 2496k    0     0   841k      0  0:00:11  0:00:02  0:00:09 1330k 56 9552k   56 5388k    0     0  1356k      0  0:00:07  0:00:03  0:00:04 1871k100 9552k  100 9552k    0     0  1959k      0  0:00:04  0:00:04 --:--:-- 2525k
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
  5 9552k    5  543k    0     0   511k      0  0:00:18  0:00:01  0:00:17  511k 50 9552k   50 4823k    0     0  2347k      0  0:00:04  0:00:02  0:00:02 4313k100 9552k  100 9552k    0     0  3148k      0  0:00:03  0:00:03 --:--:-- 4572k
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
