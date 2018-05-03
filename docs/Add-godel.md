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
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz...
 0 B / 9.32 MiB    0.00% 768.00 KiB / 9.32 MiB    8.04% 2s 1.31 MiB / 9.32 MiB   14.07% 2s 1.86 MiB / 9.32 MiB   19.94% 2s 2.56 MiB / 9.32 MiB   27.44% 2s 3.33 MiB / 9.32 MiB   35.73% 1s 4.23 MiB / 9.32 MiB   45.35% 1s 4.58 MiB / 9.32 MiB   49.17% 1s 5.28 MiB / 9.32 MiB   56.64% 1s 6.40 MiB / 9.32 MiB   68.69% 7.15 MiB / 9.32 MiB   76.74% 7.62 MiB / 9.32 MiB   81.77% 8.46 MiB / 9.32 MiB   90.74% 9.25 MiB / 9.32 MiB   99.22% 9.32 MiB / 9.32 MiB  100.00% 2s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.0.0-rc15/dist-plugin-1.0.0-rc15-linux-amd64.tgz...
 0 B / 4.73 MiB    0.00% 768.00 KiB / 4.73 MiB   15.84% 1s 1.19 MiB / 4.73 MiB   25.04% 1s 1.86 MiB / 4.73 MiB   39.24% 2.63 MiB / 4.73 MiB   55.58% 2.95 MiB / 4.73 MiB   62.31% 3.83 MiB / 4.73 MiB   80.93% 4.46 MiB / 4.73 MiB   94.14% 4.73 MiB / 4.73 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0-rc7/format-plugin-1.0.0-rc7-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 168.00 KiB / 3.32 MiB    4.94% 3s 407.15 KiB / 3.32 MiB   11.98% 2s 629.84 KiB / 3.32 MiB   18.53% 2s 812.86 KiB / 3.32 MiB   23.91% 2s 1.07 MiB / 3.32 MiB   32.10% 2s 1.31 MiB / 3.32 MiB   39.47% 1s 1.54 MiB / 3.32 MiB   46.37% 1s 1.72 MiB / 3.32 MiB   51.75% 1s 2.03 MiB / 3.32 MiB   61.11% 1s 2.16 MiB / 3.32 MiB   65.20% 1s 2.32 MiB / 3.32 MiB   69.77% 2.45 MiB / 3.32 MiB   73.86% 2.76 MiB / 3.32 MiB   83.21% 2.91 MiB / 3.32 MiB   87.78% 3.21 MiB / 3.32 MiB   96.79% 3.32 MiB / 3.32 MiB  100.00% 3s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0-rc6/ptimports-asset-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 615.77 KiB / 3.60 MiB   16.69% 1s 1.17 MiB / 3.60 MiB   32.52% 1.50 MiB / 3.60 MiB   41.64% 2.00 MiB / 3.60 MiB   55.52% 2.62 MiB / 3.60 MiB   72.84% 3.60 MiB / 3.60 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0-rc2/goland-plugin-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 679.33 KiB / 3.09 MiB   21.48% 1.50 MiB / 3.09 MiB   48.57% 2.36 MiB / 3.09 MiB   76.29% 2.81 MiB / 3.09 MiB   90.96% 3.09 MiB / 3.09 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0-rc6/check-plugin-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.53 MiB    0.00% 852.70 KiB / 3.53 MiB   23.62% 1.52 MiB / 3.53 MiB   42.98% 1.88 MiB / 3.53 MiB   53.25% 2.72 MiB / 3.53 MiB   77.02% 3.53 MiB / 3.53 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0-rc3/compiles-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 726.45 KiB / 3.71 MiB   19.13% 1.46 MiB / 3.71 MiB   39.39% 2.29 MiB / 3.71 MiB   61.85% 3.12 MiB / 3.71 MiB   84.23% 3.71 MiB / 3.71 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0-rc2/deadcode-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 831.11 KiB / 3.73 MiB   21.76% 1.67 MiB / 3.73 MiB   44.75% 2.48 MiB / 3.73 MiB   66.51% 3.27 MiB / 3.73 MiB   87.53% 3.73 MiB / 3.73 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0-rc3/errcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 1.03 MiB / 3.81 MiB   27.03% 1.75 MiB / 3.81 MiB   45.88% 2.52 MiB / 3.81 MiB   66.07% 3.32 MiB / 3.81 MiB   87.04% 3.81 MiB / 3.81 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0-rc2/extimport-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 599.77 KiB / 3.36 MiB   17.41% 1.10 MiB / 3.36 MiB   32.59% 2.34 MiB / 3.36 MiB   69.64% 3.18 MiB / 3.36 MiB   94.66% 3.36 MiB / 3.36 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0-rc3/golint-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 748.70 KiB / 3.88 MiB   18.85% 1.62 MiB / 3.88 MiB   41.85% 2.34 MiB / 3.88 MiB   60.20% 3.20 MiB / 3.88 MiB   82.51% 3.88 MiB / 3.88 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0-rc3/govet-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.17 MiB    0.00% 912.85 KiB / 3.17 MiB   28.17% 1.86 MiB / 3.17 MiB   58.70% 2.59 MiB / 3.17 MiB   81.92% 3.17 MiB / 3.17 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0-rc2/importalias-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 901.14 KiB / 3.38 MiB   26.06% 1.27 MiB / 3.38 MiB   37.47% 1.98 MiB / 3.38 MiB   58.66% 2.93 MiB / 3.38 MiB   86.90% 3.38 MiB / 3.38 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0-rc2/ineffassign-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 764.48 KiB / 3.35 MiB   22.28% 1.25 MiB / 3.35 MiB   37.30% 1.64 MiB / 3.35 MiB   48.90% 2.36 MiB / 3.35 MiB   70.37% 3.18 MiB / 3.35 MiB   95.04% 3.35 MiB / 3.35 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0-rc4/novendor-asset-1.0.0-rc4-linux-amd64.tgz...
 0 B / 3.42 MiB    0.00% 375.36 KiB / 3.42 MiB   10.73% 1s 625.88 KiB / 3.42 MiB   17.90% 1s 820.73 KiB / 3.42 MiB   23.47% 1s 904.24 KiB / 3.42 MiB   25.86% 2s 987.75 KiB / 3.42 MiB   28.25% 2s 1.70 MiB / 3.42 MiB   49.74% 1s 3.18 MiB / 3.42 MiB   93.18% 3.42 MiB / 3.42 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0-rc3/outparamcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 768.00 KiB / 3.84 MiB   19.54% 1.50 MiB / 3.84 MiB   39.00% 2.28 MiB / 3.84 MiB   59.44% 3.14 MiB / 3.84 MiB   81.81% 3.84 MiB / 3.84 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0-rc3/unconvert-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 911.33 KiB / 3.91 MiB   22.77% 1.55 MiB / 3.91 MiB   39.77% 2.42 MiB / 3.91 MiB   61.90% 3.21 MiB / 3.91 MiB   82.08% 3.91 MiB / 3.91 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0-rc2/varcheck-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 567.77 KiB / 3.75 MiB   14.80% 1s 1.25 MiB / 3.75 MiB   33.26% 1.95 MiB / 3.75 MiB   52.04% 2.83 MiB / 3.75 MiB   75.44% 3.62 MiB / 3.75 MiB   96.69% 3.75 MiB / 3.75 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0-rc1/license-plugin-1.0.0-rc1-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 717.14 KiB / 3.30 MiB   21.19% 1.13 MiB / 3.30 MiB   34.18% 1.93 MiB / 3.30 MiB   58.33% 2.75 MiB / 3.30 MiB   83.23% 3.30 MiB / 3.30 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0-rc5/test-plugin-1.0.0-rc5-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 831.11 KiB / 3.60 MiB   22.55% 1.25 MiB / 3.60 MiB   34.73% 1.73 MiB / 3.60 MiB   48.10% 2.43 MiB / 3.60 MiB   67.65% 3.22 MiB / 3.60 MiB   89.34% 3.60 MiB / 3.60 MiB  100.00% 1s
godel version 2.0.0-rc12
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
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
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
distributionSHA256=6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master fb58bd6] Add godel to project
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
  9 9546k    9  889k    0     0  2228k      0  0:00:04 --:--:--  0:00:04 2228k100 9546k  100 9546k    0     0  7066k      0  0:00:01  0:00:01 --:--:-- 9103k
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
godel version 2.0.0-rc12
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc12/godel-2.0.0-rc12.tgz
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
godel version 2.0.0-rc12
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
  6 9546k    6  583k    0     0  1138k      0  0:00:08 --:--:--  0:00:08 1138k 55 9546k   55 5303k    0     0  3369k      0  0:00:02  0:00:01  0:00:01 4448k100 9546k  100 9546k    0     0  3985k      0  0:00:02  0:00:02 --:--:-- 4762k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.0.0-rc12.tgz)= 6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116  download/godel-2.0.0-rc12.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
