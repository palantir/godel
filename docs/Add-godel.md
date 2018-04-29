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
Getting package from https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc11/godel-2.0.0-rc11.tgz...
 0 B / 9.32 MiB    0.00% 159.43 KiB / 9.32 MiB    1.67% 11s 354.94 KiB / 9.32 MiB    3.72% 10s 642.05 KiB / 9.32 MiB    6.73% 8s 914.12 KiB / 9.32 MiB    9.58% 7s 1.20 MiB / 9.32 MiB   12.86% 6s 1.41 MiB / 9.32 MiB   15.08% 6s 1.66 MiB / 9.32 MiB   17.77% 6s 1.90 MiB / 9.32 MiB   20.40% 6s 2.19 MiB / 9.32 MiB   23.54% 5s 2.44 MiB / 9.32 MiB   26.15% 5s 2.68 MiB / 9.32 MiB   28.78% 5s 2.88 MiB / 9.32 MiB   30.84% 5s 3.13 MiB / 9.32 MiB   33.58% 5s 3.39 MiB / 9.32 MiB   36.40% 4s 3.58 MiB / 9.32 MiB   38.36% 4s 3.90 MiB / 9.32 MiB   41.89% 4s 4.21 MiB / 9.32 MiB   45.17% 4s 4.49 MiB / 9.32 MiB   48.22% 3s 4.82 MiB / 9.32 MiB   51.70% 3s 5.15 MiB / 9.32 MiB   55.25% 3s 5.44 MiB / 9.32 MiB   58.33% 3s 5.77 MiB / 9.32 MiB   61.85% 2s 6.06 MiB / 9.32 MiB   65.04% 2s 6.35 MiB / 9.32 MiB   68.07% 2s 6.60 MiB / 9.32 MiB   70.83% 2s 6.87 MiB / 9.32 MiB   73.69% 1s 7.15 MiB / 9.32 MiB   76.70% 1s 7.47 MiB / 9.32 MiB   80.13% 1s 7.78 MiB / 9.32 MiB   83.41% 1s 8.02 MiB / 9.32 MiB   86.05% 8.36 MiB / 9.32 MiB   89.69% 8.67 MiB / 9.32 MiB   92.96% 8.88 MiB / 9.32 MiB   95.22% 9.15 MiB / 9.32 MiB   98.10% 9.31 MiB / 9.32 MiB   99.82% 9.32 MiB / 9.32 MiB  100.00% 7s
```

Run `./godelw version` to verify that gödel was installed correctly. If this is the first run, this invocation will
download all of the plugins and assets:

```
➜ ./godelw version
Getting package from https://palantir.bintray.com/releases/com/palantir/distgo/dist-plugin/1.0.0-rc14/dist-plugin-1.0.0-rc14-linux-amd64.tgz...
 0 B / 4.73 MiB    0.00% 149.84 KiB / 4.73 MiB    3.09% 6s 427.38 KiB / 4.73 MiB    8.82% 4s 703.55 KiB / 4.73 MiB   14.51% 3s 1.02 MiB / 4.73 MiB   21.62% 2s 1.40 MiB / 4.73 MiB   29.49% 2s 1.81 MiB / 4.73 MiB   38.15% 1s 2.16 MiB / 4.73 MiB   45.60% 1s 2.51 MiB / 4.73 MiB   53.10% 1s 2.86 MiB / 4.73 MiB   60.40% 1s 3.11 MiB / 4.73 MiB   65.71% 1s 3.37 MiB / 4.73 MiB   71.09% 3.57 MiB / 4.73 MiB   75.47% 3.77 MiB / 4.73 MiB   79.64% 3.97 MiB / 4.73 MiB   83.90% 4.10 MiB / 4.73 MiB   86.61% 4.31 MiB / 4.73 MiB   91.09% 4.49 MiB / 4.73 MiB   94.96% 4.72 MiB / 4.73 MiB   99.78% 4.73 MiB / 4.73 MiB  100.00% 3s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-plugin/format-plugin/1.0.0-rc7/format-plugin-1.0.0-rc7-linux-amd64.tgz...
 0 B / 3.32 MiB    0.00% 514.62 KiB / 3.32 MiB   15.14% 1s 971.83 KiB / 3.32 MiB   28.59% 1s 1.40 MiB / 3.32 MiB   42.16% 1.85 MiB / 3.32 MiB   55.61% 2.20 MiB / 3.32 MiB   66.25% 2.50 MiB / 3.32 MiB   75.26% 2.80 MiB / 3.32 MiB   84.27% 3.12 MiB / 3.32 MiB   94.09% 3.32 MiB / 3.32 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-format-asset-ptimports/ptimports-asset/1.0.0-rc6/ptimports-asset-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 432.36 KiB / 3.60 MiB   11.72% 1s 1.02 MiB / 3.60 MiB   28.32% 1s 1.54 MiB / 3.60 MiB   42.66% 2.04 MiB / 3.60 MiB   56.56% 2.68 MiB / 3.60 MiB   74.35% 3.18 MiB / 3.60 MiB   88.25% 3.60 MiB / 3.60 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-goland-plugin/goland-plugin/1.0.0-rc2/goland-plugin-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.09 MiB    0.00% 611.46 KiB / 3.09 MiB   19.34% 1.25 MiB / 3.09 MiB   40.46% 1.89 MiB / 3.09 MiB   61.21% 2.81 MiB / 3.09 MiB   91.14% 3.09 MiB / 3.09 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/okgo/check-plugin/1.0.0-rc6/check-plugin-1.0.0-rc6-linux-amd64.tgz...
 0 B / 3.53 MiB    0.00% 653.97 KiB / 3.53 MiB   18.11% 1.29 MiB / 3.53 MiB   36.62% 2.04 MiB / 3.53 MiB   57.76% 2.76 MiB / 3.53 MiB   78.25% 3.49 MiB / 3.53 MiB   99.07% 3.53 MiB / 3.53 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-compiles/compiles-asset/1.0.0-rc3/compiles-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.71 MiB    0.00% 721.29 KiB / 3.71 MiB   19.00% 1.51 MiB / 3.71 MiB   40.68% 2.17 MiB / 3.71 MiB   58.58% 2.82 MiB / 3.71 MiB   76.18% 3.66 MiB / 3.71 MiB   98.59% 3.71 MiB / 3.71 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-deadcode/deadcode-asset/1.0.0-rc2/deadcode-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.73 MiB    0.00% 889.21 KiB / 3.73 MiB   23.28% 1.77 MiB / 3.73 MiB   47.33% 2.39 MiB / 3.73 MiB   64.09% 3.21 MiB / 3.73 MiB   85.95% 3.73 MiB / 3.73 MiB  100.00% 0s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-errcheck/errcheck-asset/1.0.0-rc3/errcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.81 MiB    0.00% 722.35 KiB / 3.81 MiB   18.50% 1.41 MiB / 3.81 MiB   37.03% 1.96 MiB / 3.81 MiB   51.28% 2.57 MiB / 3.81 MiB   67.37% 3.22 MiB / 3.81 MiB   84.48% 3.79 MiB / 3.81 MiB   99.44% 3.81 MiB / 3.81 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-extimport/extimport-asset/1.0.0-rc2/extimport-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.36 MiB    0.00% 776.94 KiB / 3.36 MiB   22.55% 1.41 MiB / 3.36 MiB   41.94% 2.06 MiB / 3.36 MiB   61.33% 2.60 MiB / 3.36 MiB   77.15% 2.96 MiB / 3.36 MiB   88.00% 3.36 MiB / 3.36 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-golint/golint-asset/1.0.0-rc3/golint-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.88 MiB    0.00% 470.54 KiB / 3.88 MiB   11.85% 1s 959.75 KiB / 3.88 MiB   24.16% 1s 1.43 MiB / 3.88 MiB   36.78% 1s 1.85 MiB / 3.88 MiB   47.59% 2.42 MiB / 3.88 MiB   62.31% 2.99 MiB / 3.88 MiB   77.02% 3.33 MiB / 3.88 MiB   85.84% 3.55 MiB / 3.88 MiB   91.44% 3.69 MiB / 3.88 MiB   95.24% 3.88 MiB / 3.88 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-govet/govet-asset/1.0.0-rc3/govet-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.17 MiB    0.00% 360.02 KiB / 3.17 MiB   11.11% 1s 737.88 KiB / 3.17 MiB   22.77% 1s 1016.24 KiB / 3.17 MiB   31.36% 1s 1.28 MiB / 3.17 MiB   40.31% 1s 1.66 MiB / 3.17 MiB   52.33% 2.01 MiB / 3.17 MiB   63.50% 2.31 MiB / 3.17 MiB   72.95% 2.65 MiB / 3.17 MiB   83.75% 3.15 MiB / 3.17 MiB   99.57% 3.17 MiB / 3.17 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-importalias/importalias-asset/1.0.0-rc2/importalias-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.38 MiB    0.00% 627.33 KiB / 3.38 MiB   18.14% 1.10 MiB / 3.38 MiB   32.63% 1.67 MiB / 3.38 MiB   49.53% 2.06 MiB / 3.38 MiB   61.14% 2.54 MiB / 3.38 MiB   75.28% 3.06 MiB / 3.38 MiB   90.57% 3.38 MiB / 3.38 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-ineffassign/ineffassign-asset/1.0.0-rc2/ineffassign-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.35 MiB    0.00% 542.72 KiB / 3.35 MiB   15.82% 1s 1.25 MiB / 3.35 MiB   37.25% 1.89 MiB / 3.35 MiB   56.37% 2.50 MiB / 3.35 MiB   74.56% 3.23 MiB / 3.35 MiB   96.47% 3.35 MiB / 3.35 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-novendor/novendor-asset/1.0.0-rc4/novendor-asset-1.0.0-rc4-linux-amd64.tgz...
 0 B / 3.42 MiB    0.00% 542.47 KiB / 3.42 MiB   15.51% 1s 1.15 MiB / 3.42 MiB   33.82% 1.74 MiB / 3.42 MiB   50.88% 2.42 MiB / 3.42 MiB   70.78% 3.23 MiB / 3.42 MiB   94.66% 3.42 MiB / 3.42 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-outparamcheck/outparamcheck-asset/1.0.0-rc3/outparamcheck-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.84 MiB    0.00% 610.99 KiB / 3.84 MiB   15.55% 1s 1.21 MiB / 3.84 MiB   31.54% 1.77 MiB / 3.84 MiB   46.01% 2.34 MiB / 3.84 MiB   60.89% 2.91 MiB / 3.84 MiB   75.76% 3.49 MiB / 3.84 MiB   91.05% 3.84 MiB / 3.84 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-unconvert/unconvert-asset/1.0.0-rc3/unconvert-asset-1.0.0-rc3-linux-amd64.tgz...
 0 B / 3.91 MiB    0.00% 499.35 KiB / 3.91 MiB   12.48% 1s 821.55 KiB / 3.91 MiB   20.52% 1s 1.03 MiB / 3.91 MiB   26.38% 1s 1.21 MiB / 3.91 MiB   30.96% 1s 1.41 MiB / 3.91 MiB   36.12% 1s 1.56 MiB / 3.91 MiB   40.00% 1s 1.85 MiB / 3.91 MiB   47.25% 1s 2.11 MiB / 3.91 MiB   53.91% 1s 2.44 MiB / 3.91 MiB   62.55% 1s 2.77 MiB / 3.91 MiB   70.89% 3.03 MiB / 3.91 MiB   77.55% 3.34 MiB / 3.91 MiB   85.50% 3.70 MiB / 3.91 MiB   94.54% 3.91 MiB / 3.91 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-okgo-asset-varcheck/varcheck-asset/1.0.0-rc2/varcheck-asset-1.0.0-rc2-linux-amd64.tgz...
 0 B / 3.75 MiB    0.00% 431.55 KiB / 3.75 MiB   11.25% 1s 904.76 KiB / 3.75 MiB   23.58% 1s 1.44 MiB / 3.75 MiB   38.39% 1.94 MiB / 3.75 MiB   51.87% 2.50 MiB / 3.75 MiB   66.68% 3.02 MiB / 3.75 MiB   80.47% 3.40 MiB / 3.75 MiB   90.62% 3.75 MiB / 3.75 MiB  100.00% 1s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-license-plugin/license-plugin/1.0.0-rc1/license-plugin-1.0.0-rc1-linux-amd64.tgz...
 0 B / 3.30 MiB    0.00% 308.43 KiB / 3.30 MiB    9.12% 1s 542.22 KiB / 3.30 MiB   16.03% 2s 804.72 KiB / 3.30 MiB   23.78% 1s 1.02 MiB / 3.30 MiB   31.02% 1s 1.30 MiB / 3.30 MiB   39.26% 1s 1.55 MiB / 3.30 MiB   46.94% 1s 1.79 MiB / 3.30 MiB   54.09% 1s 2.08 MiB / 3.30 MiB   63.06% 2.34 MiB / 3.30 MiB   70.86% 2.65 MiB / 3.30 MiB   80.31% 2.97 MiB / 3.30 MiB   90.01% 3.25 MiB / 3.30 MiB   98.50% 3.30 MiB / 3.30 MiB  100.00% 2s
Getting package from https://palantir.bintray.com/releases/com/palantir/godel-test-plugin/test-plugin/1.0.0-rc5/test-plugin-1.0.0-rc5-linux-amd64.tgz...
 0 B / 3.60 MiB    0.00% 498.55 KiB / 3.60 MiB   13.53% 1s 916.09 KiB / 3.60 MiB   24.85% 1s 1.40 MiB / 3.60 MiB   38.88% 1.74 MiB / 3.60 MiB   48.27% 2.20 MiB / 3.60 MiB   61.11% 2.68 MiB / 3.60 MiB   74.38% 3.08 MiB / 3.60 MiB   85.71% 3.60 MiB / 3.60 MiB  100.00% 1s
godel version 2.0.0-rc11
```

Technically, this is sufficient and we have a working gödel installation. However, because the installation was
performed by downloading a distribution, the gödel installation itself does not have a checksum set (no verification was
performed on the package itself):

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc11/godel-2.0.0-rc11.tgz
distributionSHA256=
```

If we are concerned about integrity, we can specify the expected checksum for the package when using godelinit. The
expected checksums for gödel releases can be found on its [Bintray](https://bintray.com/palantir/releases/godel) page.
The checksum for 2.0.0-rc5 is 08d9ed3e33e69006a9c58ec65cef0ad9bd17af4c73b5c1d1aa116e813a954314. Install the distribution
using this checksum:

```
➜ godelinit --checksum ${GODEL_CHECKSUM}
```

If the installation succeeds with a checksum specified, it is set in the properties:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc11/godel-2.0.0-rc11.tgz
distributionSHA256=35aa494446fdfa3a506e45b03abfc1127640b2c7eefaf5a8b76c8774efc13136
```

Commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master cb3fc86] Add godel to project
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0 9546k    0   857    0     0    949      0  2:51:40 --:--:--  2:51:40   949 11 9546k   11 1057k    0     0   578k      0  0:00:16  0:00:01  0:00:15  578k 23 9546k   23 2215k    0     0   789k      0  0:00:12  0:00:02  0:00:10  789k 35 9546k   35 3406k    0     0   894k      0  0:00:10  0:00:03  0:00:07  894k 49 9546k   49 4770k    0     0   976k      0  0:00:09  0:00:04  0:00:05  976k 62 9546k   62 6013k    0     0  1030k      0  0:00:09  0:00:05  0:00:04 1218k 76 9546k   76 7349k    0     0  1079k      0  0:00:08  0:00:06  0:00:02 1263k 90 9546k   90 8641k    0     0  1103k      0  0:00:08  0:00:07  0:00:01 1279k100 9546k  100 9546k    0     0  1123k      0  0:00:08  0:00:08 --:--:-- 1309k
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
godel version 2.0.0-rc11
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded manually
do not have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/2.0.0-rc11/godel-2.0.0-rc11.tgz
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
godel version 2.0.0-rc11
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
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0  0 9546k    0   857    0     0    956      0  2:50:25 --:--:--  2:50:25   956  9 9546k    9  893k    0     0   501k      0  0:00:19  0:00:01  0:00:18  501k 18 9546k   18 1776k    0     0   639k      0  0:00:14  0:00:02  0:00:12  639k 28 9546k   28 2733k    0     0   722k      0  0:00:13  0:00:03  0:00:10  722k 40 9546k   40 3820k    0     0   799k      0  0:00:11  0:00:04  0:00:07  799k 53 9546k   53 5093k    0     0   882k      0  0:00:10  0:00:05  0:00:05 1043k 65 9546k   65 6295k    0     0   927k      0  0:00:10  0:00:06  0:00:04 1078k 77 9546k   77 7391k    0     0   949k      0  0:00:10  0:00:07  0:00:03 1122k 89 9546k   89 8571k    0     0   974k      0  0:00:09  0:00:08  0:00:01 1165k100 9546k  100 9546k    0     0   990k      0  0:00:09  0:00:09 --:--:-- 1177k
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"${GODEL_VERSION}".tgz
SHA256(download/godel-2.0.0-rc11.tgz)= 35aa494446fdfa3a506e45b03abfc1127640b2c7eefaf5a8b76c8774efc13136
➜ shasum -a 256 download/godel-"${GODEL_VERSION}".tgz
35aa494446fdfa3a506e45b03abfc1127640b2c7eefaf5a8b76c8774efc13136  download/godel-2.0.0-rc11.tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
```
