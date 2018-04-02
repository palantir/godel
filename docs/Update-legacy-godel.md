Summary
-------
Projects that use a pre-2.0 version of gödel can be updated using `godelinit`.

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
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0

Update projects using pre-2.0 version of gödel
----------------------------------------------
Projects that use a pre-2.0 version of gödel can be updated to use the latest version of gödel using the `godelinit`
program. As an exmaple, we will update the version of gödel used by `github.com/nmiyake/echgo` (which was the tutorial
project for the 1.x release of gödel) to use the latest version.

Start by getting the repository and switching to it:

```
➜ go get github.com/nmiyake/echgo
➜ cd ${GOPATH}/src/github.com/nmiyake/echgo
```

Run `./godelw version` to verify that the version of gödel used by the project is pre-1.0:

```
➜ ./godelw version
Downloading https://palantir.bintray.com/releases/com/palantir/godel/godel/0.27.0/godel-0.27.0.tgz to /root/.godel/downloads/godel-0.27.0.tgz...

     0K .......... .......... .......... .......... ..........  0% 1.75M 6s
    50K .......... .......... .......... .......... ..........  0% 1.75M 6s
   100K .......... .......... .......... .......... ..........  1% 6.46M 5s
   150K .......... .......... .......... .......... ..........  1% 2.06M 5s
   200K .......... .......... .......... .......... ..........  2% 7.12M 4s
   250K .......... .......... .......... .......... ..........  2% 3.89M 4s
   300K .......... .......... .......... .......... ..........  3% 6.14M 3s
   350K .......... .......... .......... .......... ..........  3% 3.54M 3s
   400K .......... .......... .......... .......... ..........  4% 4.66M 3s
   450K .......... .......... .......... .......... ..........  4% 12.6M 3s
   500K .......... .......... .......... .......... ..........  4% 1.21M 3s
   550K .......... .......... .......... .......... ..........  5% 2.29M 4s
   600K .......... .......... .......... .......... ..........  5% 39.7M 3s
   650K .......... .......... .......... .......... ..........  6% 38.4M 3s
   700K .......... .......... .......... .......... ..........  6% 37.5M 3s
   750K .......... .......... .......... .......... ..........  7% 71.5M 3s
   800K .......... .......... .......... .......... ..........  7% 7.61M 3s
   850K .......... .......... .......... .......... ..........  8% 6.79M 2s
   900K .......... .......... .......... .......... ..........  8% 3.40M 3s
   950K .......... .......... .......... .......... ..........  9% 5.98M 2s
  1000K .......... .......... .......... .......... ..........  9% 3.26M 2s
  1050K .......... .......... .......... .......... ..........  9% 6.30M 2s
  1100K .......... .......... .......... .......... .......... 10% 20.4M 2s
  1150K .......... .......... .......... .......... .......... 10% 1.33M 3s
  1200K .......... .......... .......... .......... .......... 11% 14.4M 2s
  1250K .......... .......... .......... .......... .......... 11% 23.4M 2s
  1300K .......... .......... .......... .......... .......... 12% 36.6M 2s
  1350K .......... .......... .......... .......... .......... 12% 42.6M 2s
  1400K .......... .......... .......... .......... .......... 13% 9.60M 2s
  1450K .......... .......... .......... .......... .......... 13% 4.12M 2s
  1500K .......... .......... .......... .......... .......... 14% 5.01M 2s
  1550K .......... .......... .......... .......... .......... 14% 26.7M 2s
  1600K .......... .......... .......... .......... .......... 14% 1.68M 2s
  1650K .......... .......... .......... .......... .......... 15% 6.71M 2s
  1700K .......... .......... .......... .......... .......... 15% 7.58M 2s
  1750K .......... .......... .......... .......... .......... 16% 4.26M 2s
  1800K .......... .......... .......... .......... .......... 16% 5.38M 2s
  1850K .......... .......... .......... .......... .......... 17% 5.59M 2s
  1900K .......... .......... .......... .......... .......... 17% 6.56M 2s
  1950K .......... .......... .......... .......... .......... 18% 6.46M 2s
  2000K .......... .......... .......... .......... .......... 18% 6.60M 2s
  2050K .......... .......... .......... .......... .......... 19% 8.86M 2s
  2100K .......... .......... .......... .......... .......... 19% 5.25M 2s
  2150K .......... .......... .......... .......... .......... 19% 4.26M 2s
  2200K .......... .......... .......... .......... .......... 20% 5.24M 2s
  2250K .......... .......... .......... .......... .......... 20% 8.69M 2s
  2300K .......... .......... .......... .......... .......... 21% 10.7M 2s
  2350K .......... .......... .......... .......... .......... 21% 4.47M 2s
  2400K .......... .......... .......... .......... .......... 22% 3.73M 2s
  2450K .......... .......... .......... .......... .......... 22% 8.11M 2s
  2500K .......... .......... .......... .......... .......... 23% 6.03M 2s
  2550K .......... .......... .......... .......... .......... 23% 2.97M 2s
  2600K .......... .......... .......... .......... .......... 24% 3.38M 2s
  2650K .......... .......... .......... .......... .......... 24% 11.3M 2s
  2700K .......... .......... .......... .......... .......... 24% 6.42M 2s
  2750K .......... .......... .......... .......... .......... 25% 5.85M 2s
  2800K .......... .......... .......... .......... .......... 25% 4.18M 2s
  2850K .......... .......... .......... .......... .......... 26% 8.41M 2s
  2900K .......... .......... .......... .......... .......... 26% 6.16M 2s
  2950K .......... .......... .......... .......... .......... 27% 7.04M 2s
  3000K .......... .......... .......... .......... .......... 27% 5.71M 2s
  3050K .......... .......... .......... .......... .......... 28% 3.66M 2s
  3100K .......... .......... .......... .......... .......... 28% 9.58M 2s
  3150K .......... .......... .......... .......... .......... 29% 6.37M 2s
  3200K .......... .......... .......... .......... .......... 29% 9.16M 2s
  3250K .......... .......... .......... .......... .......... 29% 3.88M 2s
  3300K .......... .......... .......... .......... .......... 30% 5.89M 2s
  3350K .......... .......... .......... .......... .......... 30% 6.41M 2s
  3400K .......... .......... .......... .......... .......... 31% 5.33M 1s
  3450K .......... .......... .......... .......... .......... 31% 8.69M 1s
  3500K .......... .......... .......... .......... .......... 32% 4.74M 1s
  3550K .......... .......... .......... .......... .......... 32% 7.38M 1s
  3600K .......... .......... .......... .......... .......... 33% 4.10M 1s
  3650K .......... .......... .......... .......... .......... 33% 8.90M 1s
  3700K .......... .......... .......... .......... .......... 34% 8.19M 1s
  3750K .......... .......... .......... .......... .......... 34% 5.93M 1s
  3800K .......... .......... .......... .......... .......... 34% 4.20M 1s
  3850K .......... .......... .......... .......... .......... 35% 4.06M 1s
  3900K .......... .......... .......... .......... .......... 35% 19.3M 1s
  3950K .......... .......... .......... .......... .......... 36% 5.59M 1s
  4000K .......... .......... .......... .......... .......... 36% 4.16M 1s
  4050K .......... .......... .......... .......... .......... 37% 4.06M 1s
  4100K .......... .......... .......... .......... .......... 37% 5.51M 1s
  4150K .......... .......... .......... .......... .......... 38% 8.98M 1s
  4200K .......... .......... .......... .......... .......... 38% 6.08M 1s
  4250K .......... .......... .......... .......... .......... 39% 3.29M 1s
  4300K .......... .......... .......... .......... .......... 39% 24.2M 1s
  4350K .......... .......... .......... .......... .......... 39% 4.21M 1s
  4400K .......... .......... .......... .......... .......... 40% 10.8M 1s
  4450K .......... .......... .......... .......... .......... 40% 4.19M 1s
  4500K .......... .......... .......... .......... .......... 41% 4.74M 1s
  4550K .......... .......... .......... .......... .......... 41% 7.64M 1s
  4600K .......... .......... .......... .......... .......... 42% 4.32M 1s
  4650K .......... .......... .......... .......... .......... 42% 3.02M 1s
  4700K .......... .......... .......... .......... .......... 43% 4.70M 1s
  4750K .......... .......... .......... .......... .......... 43% 8.95M 1s
  4800K .......... .......... .......... .......... .......... 44% 4.41M 1s
  4850K .......... .......... .......... .......... .......... 44% 7.63M 1s
  4900K .......... .......... .......... .......... .......... 44% 6.67M 1s
  4950K .......... .......... .......... .......... .......... 45% 5.95M 1s
  5000K .......... .......... .......... .......... .......... 45% 2.12M 1s
  5050K .......... .......... .......... .......... .......... 46% 23.9M 1s
  5100K .......... .......... .......... .......... .......... 46% 13.7M 1s
  5150K .......... .......... .......... .......... .......... 47% 2.24M 1s
  5200K .......... .......... .......... .......... .......... 47% 3.32M 1s
  5250K .......... .......... .......... .......... .......... 48% 6.34M 1s
  5300K .......... .......... .......... .......... .......... 48% 16.6M 1s
  5350K .......... .......... .......... .......... .......... 49% 4.78M 1s
  5400K .......... .......... .......... .......... .......... 49% 3.08M 1s
  5450K .......... .......... .......... .......... .......... 49% 11.2M 1s
  5500K .......... .......... .......... .......... .......... 50% 6.15M 1s
  5550K .......... .......... .......... .......... .......... 50% 3.17M 1s
  5600K .......... .......... .......... .......... .......... 51% 3.93M 1s
  5650K .......... .......... .......... .......... .......... 51% 10.6M 1s
  5700K .......... .......... .......... .......... .......... 52% 4.37M 1s
  5750K .......... .......... .......... .......... .......... 52% 6.01M 1s
  5800K .......... .......... .......... .......... .......... 53% 4.22M 1s
  5850K .......... .......... .......... .......... .......... 53% 5.63M 1s
  5900K .......... .......... .......... .......... .......... 54% 4.76M 1s
  5950K .......... .......... .......... .......... .......... 54% 2.00M 1s
  6000K .......... .......... .......... .......... .......... 54% 7.20M 1s
  6050K .......... .......... .......... .......... .......... 55% 4.36M 1s
  6100K .......... .......... .......... .......... .......... 55% 6.66M 1s
  6150K .......... .......... .......... .......... .......... 56% 5.60M 1s
  6200K .......... .......... .......... .......... .......... 56% 5.11M 1s
  6250K .......... .......... .......... .......... .......... 57% 8.58M 1s
  6300K .......... .......... .......... .......... .......... 57% 3.57M 1s
  6350K .......... .......... .......... .......... .......... 58% 3.78M 1s
  6400K .......... .......... .......... .......... .......... 58% 10.5M 1s
  6450K .......... .......... .......... .......... .......... 59% 7.44M 1s
  6500K .......... .......... .......... .......... .......... 59% 2.49M 1s
  6550K .......... .......... .......... .......... .......... 59% 5.91M 1s
  6600K .......... .......... .......... .......... .......... 60% 7.06M 1s
  6650K .......... .......... .......... .......... .......... 60% 6.96M 1s
  6700K .......... .......... .......... .......... .......... 61% 2.39M 1s
  6750K .......... .......... .......... .......... .......... 61% 27.2M 1s
  6800K .......... .......... .......... .......... .......... 62% 9.74M 1s
  6850K .......... .......... .......... .......... .......... 62% 2.79M 1s
  6900K .......... .......... .......... .......... .......... 63% 5.23M 1s
  6950K .......... .......... .......... .......... .......... 63% 4.67M 1s
  7000K .......... .......... .......... .......... .......... 64% 9.01M 1s
  7050K .......... .......... .......... .......... .......... 64% 4.47M 1s
  7100K .......... .......... .......... .......... .......... 64% 8.22M 1s
  7150K .......... .......... .......... .......... .......... 65% 5.30M 1s
  7200K .......... .......... .......... .......... .......... 65% 6.73M 1s
  7250K .......... .......... .......... .......... .......... 66% 5.80M 1s
  7300K .......... .......... .......... .......... .......... 66% 4.04M 1s
  7350K .......... .......... .......... .......... .......... 67% 3.71M 1s
  7400K .......... .......... .......... .......... .......... 67% 11.3M 1s
  7450K .......... .......... .......... .......... .......... 68% 6.30M 1s
  7500K .......... .......... .......... .......... .......... 68% 3.65M 1s
  7550K .......... .......... .......... .......... .......... 69% 5.76M 1s
  7600K .......... .......... .......... .......... .......... 69% 5.22M 1s
  7650K .......... .......... .......... .......... .......... 69% 7.14M 1s
  7700K .......... .......... .......... .......... .......... 70% 7.91M 1s
  7750K .......... .......... .......... .......... .......... 70% 2.51M 1s
  7800K .......... .......... .......... .......... .......... 71% 16.3M 1s
  7850K .......... .......... .......... .......... .......... 71% 8.70M 1s
  7900K .......... .......... .......... .......... .......... 72% 2.83M 1s
  7950K .......... .......... .......... .......... .......... 72% 9.59M 1s
  8000K .......... .......... .......... .......... .......... 73% 4.39M 1s
  8050K .......... .......... .......... .......... .......... 73% 4.79M 1s
  8100K .......... .......... .......... .......... .......... 74% 5.59M 1s
  8150K .......... .......... .......... .......... .......... 74% 5.76M 1s
  8200K .......... .......... .......... .......... .......... 74% 6.60M 1s
  8250K .......... .......... .......... .......... .......... 75% 3.23M 1s
  8300K .......... .......... .......... .......... .......... 75% 11.4M 1s
  8350K .......... .......... .......... .......... .......... 76% 5.72M 1s
  8400K .......... .......... .......... .......... .......... 76% 3.21M 0s
  8450K .......... .......... .......... .......... .......... 77% 7.14M 0s
  8500K .......... .......... .......... .......... .......... 77% 5.72M 0s
  8550K .......... .......... .......... .......... .......... 78% 5.02M 0s
  8600K .......... .......... .......... .......... .......... 78% 6.90M 0s
  8650K .......... .......... .......... .......... .......... 79% 5.67M 0s
  8700K .......... .......... .......... .......... .......... 79% 6.02M 0s
  8750K .......... .......... .......... .......... .......... 79% 6.55M 0s
  8800K .......... .......... .......... .......... .......... 80% 6.19M 0s
  8850K .......... .......... .......... .......... .......... 80% 7.25M 0s
  8900K .......... .......... .......... .......... .......... 81% 4.90M 0s
  8950K .......... .......... .......... .......... .......... 81% 7.53M 0s
  9000K .......... .......... .......... .......... .......... 82% 5.55M 0s
  9050K .......... .......... .......... .......... .......... 82% 3.86M 0s
  9100K .......... .......... .......... .......... .......... 83% 5.22M 0s
  9150K .......... .......... .......... .......... .......... 83% 14.9M 0s
  9200K .......... .......... .......... .......... .......... 84% 6.23M 0s
  9250K .......... .......... .......... .......... .......... 84% 3.37M 0s
  9300K .......... .......... .......... .......... .......... 84% 10.1M 0s
  9350K .......... .......... .......... .......... .......... 85% 5.26M 0s
  9400K .......... .......... .......... .......... .......... 85% 4.50M 0s
  9450K .......... .......... .......... .......... .......... 86% 2.05M 0s
  9500K .......... .......... .......... .......... .......... 86% 8.52M 0s
  9550K .......... .......... .......... .......... .......... 87% 9.72M 0s
  9600K .......... .......... .......... .......... .......... 87% 4.90M 0s
  9650K .......... .......... .......... .......... .......... 88% 5.00M 0s
  9700K .......... .......... .......... .......... .......... 88% 6.93M 0s
  9750K .......... .......... .......... .......... .......... 89% 4.16M 0s
  9800K .......... .......... .......... .......... .......... 89% 10.6M 0s
  9850K .......... .......... .......... .......... .......... 89% 4.71M 0s
  9900K .......... .......... .......... .......... .......... 90% 10.5M 0s
  9950K .......... .......... .......... .......... .......... 90% 2.52M 0s
 10000K .......... .......... .......... .......... .......... 91% 5.57M 0s
 10050K .......... .......... .......... .......... .......... 91% 9.21M 0s
 10100K .......... .......... .......... .......... .......... 92% 5.59M 0s
 10150K .......... .......... .......... .......... .......... 92% 6.93M 0s
 10200K .......... .......... .......... .......... .......... 93% 5.46M 0s
 10250K .......... .......... .......... .......... .......... 93% 4.24M 0s
 10300K .......... .......... .......... .......... .......... 94% 12.3M 0s
 10350K .......... .......... .......... .......... .......... 94% 5.71M 0s
 10400K .......... .......... .......... .......... .......... 94% 5.13M 0s
 10450K .......... .......... .......... .......... .......... 95% 4.66M 0s
 10500K .......... .......... .......... .......... .......... 95% 1.43M 0s
 10550K .......... .......... .......... .......... .......... 96% 38.2M 0s
 10600K .......... .......... .......... .......... .......... 96% 2.94M 0s
 10650K .......... .......... .......... .......... .......... 97% 28.2M 0s
 10700K .......... .......... .......... .......... .......... 97% 1.67M 0s
 10750K .......... .......... .......... .......... .......... 98% 3.32M 0s
 10800K .......... .......... .......... .......... .......... 98% 4.75M 0s
 10850K .......... .......... .......... .......... .......... 99% 13.4M 0s
 10900K .......... .......... .......... .......... .......... 99% 6.65M 0s
 10950K .......... .......... .......... .......... .......... 99% 3.38M 0s
 11000K ..                                                    100% 24.1M=2.1sgodel version 0.27.0
```

The godelinit program can be used to install or update gödel. Although it is technically possible to use the
`./godelw update` mechanism to get the new binary, pre-2.0 versions of gödel do not know about configuration upgrades,
so running `./godelw update` will only update the gödel version. Using godelinit has the advantage that it knows to
invoke the commands to upgrade the configuration after gödel itself has been updated.

First, ensure that the latest version of godelinit is installed:

```
➜ go get -u github.com/palantir/godel/godelinit
```

Then, run `godelinit`. Running the program with no arguments determines the latest released version of gödel and either
installs it or upgrades the version of gödel in the current project to that version. In this case, because the newest
version of gödel is >=2 and the version of gödel in the current project is <2, the legacy configuration will be
upgraded:

```
➜ godelinit
Upgraded configuration for godel.yml
Upgraded configuration for dist-plugin.yml
Upgraded configuration for license-plugin.yml
Upgraded configuration for format-plugin.yml
Upgraded configuration for test-plugin.yml
Upgraded configuration for check-plugin.yml
WARNING: The following configuration file(s) were non-empty and had no known upgraders for legacy configuration: [generate.yml]
         If these configuration file(s) are for plugins, add the plugins to the configuration and rerun the upgrade task.
```

As indicated by the output, most of the configuration is automatically upgraded. Verify that this is the case:

```
➜ git status
On branch master
Your branch is up-to-date with 'origin/master'.
Changes not staged for commit:
  (use "git add/rm <file>..." to update what will be committed)
  (use "git checkout -- <file>..." to discard changes in working directory)

	deleted:    godel/config/check.yml
	deleted:    godel/config/dist.yml
	deleted:    godel/config/exclude.yml
	deleted:    godel/config/format.yml
	modified:   godel/config/godel.properties
	deleted:    godel/config/imports.yml
	deleted:    godel/config/license.yml
	deleted:    godel/config/test.yml
	modified:   godelw

Untracked files:
  (use "git add <file>..." to include in what will be committed)

	godel/config/check-plugin.yml
	godel/config/dist-plugin.yml
	godel/config/format-plugin.yml
	godel/config/godel.yml
	godel/config/license-plugin.yml
	godel/config/test-plugin.yml

no changes added to commit (use "git add" and/or "git commit -a")
```

However, note that the output said that `generate.yml` was non-empty but there were no known upgraders for it. This is
because, although `generate` was a builtin task prior to gödel 2.0, starting with gödel 2.0 this task is provided by a
plugin. We can upgrade the configuration by adding the plugin and explicitly running the `upgrade-config` task with the
`--legacy` flag set. Start by updating the `godel.yml` configuration to add the plugin:

```
➜ echo 'plugins:
  resolvers:
    - "https://palantir.bintray.com/releases/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz"
  plugins:
    - locator:
        id: "com.palantir.go-generate:generate-plugin:1.0.0-rc3"
exclude:
  names:
    - "\\\\..+"
    - "vendor"
    - ".+_string.go"
  paths:
    - "godel"' > godel/config/godel.yml
```

Now that the plugin has been added to the configuration, run the `upgrade-config` task in legacy mode:

```
➜ ./godelw upgrade-config --legacy
Upgraded configuration for generate-plugin.yml
```

The message indicates that the configuration was upgraded. You can verify this by observing that a `generate-plugin.yml`
file is now present in the configuration directory:

```
➜ cat godel/config/generate-plugin.yml
generators:
  stringer:
    go-generate-dir: generator
    gen-paths:
      paths:
      - echo/type_string.go
```

Run the verify task to verify that the upgraded configuration is working correctly:

```
➜ ./godelw verify --apply=false
Running format...
Running generate...
Running license...
Running check...
[compiles]      Running compiles...
[extimport]     Running extimport...
[errcheck]      Running errcheck...
[deadcode]      Running deadcode...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[errcheck]      Finished errcheck
[unconvert]     Running unconvert...
[compiles]      Finished compiles
[varcheck]      Running varcheck...
[deadcode]      Finished deadcode
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	0.010s
?   	github.com/nmiyake/echgo/generator       	[no test files]
--- FAIL: TestInvalidType (0.09s)
panic: command [/go/src/github.com/nmiyake/echgo/godelw artifacts build --absolute --os-arch=linux-amd64 --requires-build echgo] failed with output:
Error: unknown flag: --os-arch
Usage:
  distgo artifacts build [flags] [product-build-ids]

Flags:
      --absolute         print the absolute path for artifacts
  -h, --help             help for build
      --requires-build   only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)

Global Flags:
      --assets stringSlice    path(s) to the plugin asset(s)
      --config string         path to the plugin configuration file
      --debug                 run in debug mode
      --godel-config string   path to the godel.yml configuration file
      --project-dir string    path to project directory
Error: exit status 1 [recovered]
	panic: command [/go/src/github.com/nmiyake/echgo/godelw artifacts build --absolute --os-arch=linux-amd64 --requires-build echgo] failed with output:
Error: unknown flag: --os-arch
Usage:
  distgo artifacts build [flags] [product-build-ids]

Flags:
      --absolute         print the absolute path for artifacts
  -h, --help             help for build
      --requires-build   only prints the artifacts that require building (omits artifacts that are already built and are up-to-date)

Global Flags:
      --assets stringSlice    path(s) to the plugin asset(s)
      --config string         path to the plugin configuration file
      --debug                 run in debug mode
      --godel-config string   path to the godel.yml configuration file
      --project-dir string    path to project directory
Error: exit status 1

goroutine 18 [running]:
testing.tRunner.func1(0xc42007c1e0)
	/usr/local/go/src/testing/testing.go:742 +0x29d
panic(0x516640, 0xc42007a560)
	/usr/local/go/src/runtime/panic.go:505 +0x229
github.com/nmiyake/echgo/integration_test_test.TestInvalidType(0xc42007c1e0)
	/go/src/github.com/nmiyake/echgo/integration_test/integration_test.go:27 +0x49e
testing.tRunner(0xc42007c1e0, 0x548608)
	/usr/local/go/src/testing/testing.go:777 +0xd0
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:824 +0x2e0
FAIL	github.com/nmiyake/echgo/integration_test	0.097s
Error: 1 package(s) had failing tests:
	github.com/nmiyake/echgo/integration_test
Failed tasks:
	test
```

Although all of the checks pass, one of the tests now fails. This is because the integration test uses the
`github.com/palantir/godel/pkg/products`. This package had an API-breaking change: projects that use version 2 or later
of gödel that use this package must use the v2 version of the package instead. In order to fix this, vendor the v2
package and update the test to use it:

```
➜ cp -r ${GOPATH}/src/github.com/palantir/godel/pkg/products/v2/ vendor/github.com/palantir/godel/pkg/products/
➜ sed -i 's:github.com/palantir/godel/pkg/products:github.com/palantir/godel/pkg/products/v2/products:g' integration_test/integration_test.go
```

Confirm that the "verify" task now succeeds:

```
➜ ./godelw verify --apply=false
Running format...
Running generate...
Running license...
Running check...
[deadcode]      Running deadcode...
[compiles]      Running compiles...
[extimport]     Running extimport...
[errcheck]      Running errcheck...
[extimport]     Finished extimport
[golint]        Running golint...
[golint]        Finished golint
[govet]         Running govet...
[govet]         Finished govet
[importalias]   Running importalias...
[importalias]   Finished importalias
[ineffassign]   Running ineffassign...
[ineffassign]   Finished ineffassign
[novendor]      Running novendor...
[novendor]      Finished novendor
[outparamcheck] Running outparamcheck...
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[errcheck]      Finished errcheck
[varcheck]      Running varcheck...
[outparamcheck] Finished outparamcheck
[compiles]      Finished compiles
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	(cached)
?   	github.com/nmiyake/echgo/generator       	[no test files]
ok  	github.com/nmiyake/echgo/integration_test	0.527s
```

Finally, clean up the state:

```
➜ cd ${GOPATH}/src/${PROJECT_PATH}
➜ rm -rf ${GOPATH}/src/github.com/nmiyake/echgo
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
* Go files have license headers
* `godel/config/godel.yml` is configured to add the go-generate plugin
* `godel/config/generate-plugin.yml` is configured to generate string function
* `godel/config/godel.yml` is configured to ignore all `.+_string.go` files
* `integration_test` contains integration tests
* `godel/config/test-plugin.yml` is configured to specify the "integration" tag
* `docs` contains documentation
* `.circleci/config.yml` exists
* Project is tagged as 1.0.0
* `godelw` is updated to the latest version

Tutorial next step
------------------
[Other commands](https://github.com/palantir/godel/wiki/Other-commands)

More
----
### Update to a specific version
The `--version` flag can be used to upgrade gödel to a specific version.

### Provide a checksum
The `--checksum` flag can be used to specify the expected checksum for the update.

### Sync installation to contents of `godel/config/godel.properties`
The `--sync` flag can be used to specify that the update operation should update gödel to match the values specified
in the `godel/config/godel.properties` file. The file should be edited to match the desired state first.

Updating gödel in this manner requires knowing the distribution URL for the target version. Although it is optional, it
is recommended to have the SHA-256 checksum of the distribution as well to ensure the integrity of the update.
