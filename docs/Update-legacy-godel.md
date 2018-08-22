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
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
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

     0K .......... .......... .......... .......... ..........  0%  659K 17s
    50K .......... .......... .......... .......... ..........  0% 1.60M 12s
   100K .......... .......... .......... .......... ..........  1% 6.65M 8s
   150K .......... .......... .......... .......... ..........  1% 1.19M 8s
   200K .......... .......... .......... .......... ..........  2% 5.76M 7s
   250K .......... .......... .......... .......... ..........  2% 8.43M 6s
   300K .......... .......... .......... .......... ..........  3% 3.13M 6s
   350K .......... .......... .......... .......... ..........  3% 1.65M 6s
   400K .......... .......... .......... .......... ..........  4% 4.18M 5s
   450K .......... .......... .......... .......... ..........  4% 2.03M 5s
   500K .......... .......... .......... .......... ..........  4% 3.38M 5s
   550K .......... .......... .......... .......... ..........  5% 5.74M 5s
   600K .......... .......... .......... .......... ..........  5% 5.31M 4s
   650K .......... .......... .......... .......... ..........  6% 8.72M 4s
   700K .......... .......... .......... .......... ..........  6% 4.36M 4s
   750K .......... .......... .......... .......... ..........  7% 2.70M 4s
   800K .......... .......... .......... .......... ..........  7% 2.42M 4s
   850K .......... .......... .......... .......... ..........  8% 16.7M 4s
   900K .......... .......... .......... .......... ..........  8% 7.21M 4s
   950K .......... .......... .......... .......... ..........  9% 2.44M 4s
  1000K .......... .......... .......... .......... ..........  9% 5.63M 4s
  1050K .......... .......... .......... .......... ..........  9% 6.30M 3s
  1100K .......... .......... .......... .......... .......... 10% 7.96M 3s
  1150K .......... .......... .......... .......... .......... 10% 1.13M 4s
  1200K .......... .......... .......... .......... .......... 11% 8.14M 3s
  1250K .......... .......... .......... .......... .......... 11% 1.05M 4s
  1300K .......... .......... .......... .......... .......... 12% 3.82M 4s
  1350K .......... .......... .......... .......... .......... 12% 7.43M 3s
  1400K .......... .......... .......... .......... .......... 13% 1.68M 4s
  1450K .......... .......... .......... .......... .......... 13% 6.08M 3s
  1500K .......... .......... .......... .......... .......... 14% 3.25M 3s
  1550K .......... .......... .......... .......... .......... 14% 3.51M 3s
  1600K .......... .......... .......... .......... .......... 14% 3.82M 3s
  1650K .......... .......... .......... .......... .......... 15% 3.26M 3s
  1700K .......... .......... .......... .......... .......... 15% 3.67M 3s
  1750K .......... .......... .......... .......... .......... 16% 2.44M 3s
  1800K .......... .......... .......... .......... .......... 16% 2.19M 3s
  1850K .......... .......... .......... .......... .......... 17% 11.1M 3s
  1900K .......... .......... .......... .......... .......... 17% 6.41M 3s
  1950K .......... .......... .......... .......... .......... 18% 2.37M 3s
  2000K .......... .......... .......... .......... .......... 18% 2.77M 3s
  2050K .......... .......... .......... .......... .......... 19% 13.8M 3s
  2100K .......... .......... .......... .......... .......... 19% 3.81M 3s
  2150K .......... .......... .......... .......... .......... 19% 3.08M 3s
  2200K .......... .......... .......... .......... .......... 20% 6.13M 3s
  2250K .......... .......... .......... .......... .......... 20% 7.17M 3s
  2300K .......... .......... .......... .......... .......... 21% 3.29M 3s
  2350K .......... .......... .......... .......... .......... 21% 3.57M 3s
  2400K .......... .......... .......... .......... .......... 22% 6.72M 3s
  2450K .......... .......... .......... .......... .......... 22% 6.56M 3s
  2500K .......... .......... .......... .......... .......... 23% 4.35M 3s
  2550K .......... .......... .......... .......... .......... 23% 5.40M 3s
  2600K .......... .......... .......... .......... .......... 24% 7.17M 3s
  2650K .......... .......... .......... .......... .......... 24% 3.55M 3s
  2700K .......... .......... .......... .......... .......... 24% 2.15M 3s
  2750K .......... .......... .......... .......... .......... 25% 6.93M 3s
  2800K .......... .......... .......... .......... .......... 25% 7.11M 2s
  2850K .......... .......... .......... .......... .......... 26% 3.67M 2s
  2900K .......... .......... .......... .......... .......... 26% 10.8M 2s
  2950K .......... .......... .......... .......... .......... 27% 5.71M 2s
  3000K .......... .......... .......... .......... .......... 27% 4.12M 2s
  3050K .......... .......... .......... .......... .......... 28% 5.94M 2s
  3100K .......... .......... .......... .......... .......... 28% 5.52M 2s
  3150K .......... .......... .......... .......... .......... 29% 4.67M 2s
  3200K .......... .......... .......... .......... .......... 29% 4.00M 2s
  3250K .......... .......... .......... .......... .......... 29% 3.14M 2s
  3300K .......... .......... .......... .......... .......... 30% 5.33M 2s
  3350K .......... .......... .......... .......... .......... 30% 5.81M 2s
  3400K .......... .......... .......... .......... .......... 31% 9.09M 2s
  3450K .......... .......... .......... .......... .......... 31% 2.17M 2s
  3500K .......... .......... .......... .......... .......... 32% 4.53M 2s
  3550K .......... .......... .......... .......... .......... 32% 4.94M 2s
  3600K .......... .......... .......... .......... .......... 33% 3.07M 2s
  3650K .......... .......... .......... .......... .......... 33% 4.79M 2s
  3700K .......... .......... .......... .......... .......... 34% 4.85M 2s
  3750K .......... .......... .......... .......... .......... 34% 4.56M 2s
  3800K .......... .......... .......... .......... .......... 34% 6.25M 2s
  3850K .......... .......... .......... .......... .......... 35% 2.73M 2s
  3900K .......... .......... .......... .......... .......... 35% 5.69M 2s
  3950K .......... .......... .......... .......... .......... 36% 6.57M 2s
  4000K .......... .......... .......... .......... .......... 36% 6.18M 2s
  4050K .......... .......... .......... .......... .......... 37% 4.25M 2s
  4100K .......... .......... .......... .......... .......... 37% 6.09M 2s
  4150K .......... .......... .......... .......... .......... 38% 3.20M 2s
  4200K .......... .......... .......... .......... .......... 38% 3.47M 2s
  4250K .......... .......... .......... .......... .......... 39% 4.74M 2s
  4300K .......... .......... .......... .......... .......... 39% 6.87M 2s
  4350K .......... .......... .......... .......... .......... 39% 7.94M 2s
  4400K .......... .......... .......... .......... .......... 40% 3.23M 2s
  4450K .......... .......... .......... .......... .......... 40% 3.02M 2s
  4500K .......... .......... .......... .......... .......... 41% 5.67M 2s
  4550K .......... .......... .......... .......... .......... 41% 6.93M 2s
  4600K .......... .......... .......... .......... .......... 42% 6.42M 2s
  4650K .......... .......... .......... .......... .......... 42% 2.30M 2s
  4700K .......... .......... .......... .......... .......... 43% 6.72M 2s
  4750K .......... .......... .......... .......... .......... 43% 2.86M 2s
  4800K .......... .......... .......... .......... .......... 44% 7.28M 2s
  4850K .......... .......... .......... .......... .......... 44% 6.43M 2s
  4900K .......... .......... .......... .......... .......... 44% 4.70M 2s
  4950K .......... .......... .......... .......... .......... 45% 3.58M 2s
  5000K .......... .......... .......... .......... .......... 45% 2.19M 2s
  5050K .......... .......... .......... .......... .......... 46% 6.02M 2s
  5100K .......... .......... .......... .......... .......... 46% 8.46M 2s
  5150K .......... .......... .......... .......... .......... 47% 4.45M 2s
  5200K .......... .......... .......... .......... .......... 47% 2.71M 2s
  5250K .......... .......... .......... .......... .......... 48% 6.51M 2s
  5300K .......... .......... .......... .......... .......... 48% 6.12M 1s
  5350K .......... .......... .......... .......... .......... 49% 6.75M 1s
  5400K .......... .......... .......... .......... .......... 49% 6.71M 1s
  5450K .......... .......... .......... .......... .......... 49% 3.73M 1s
  5500K .......... .......... .......... .......... .......... 50% 5.00M 1s
  5550K .......... .......... .......... .......... .......... 50% 5.27M 1s
  5600K .......... .......... .......... .......... .......... 51% 5.23M 1s
  5650K .......... .......... .......... .......... .......... 51% 3.23M 1s
  5700K .......... .......... .......... .......... .......... 52% 2.73M 1s
  5750K .......... .......... .......... .......... .......... 52% 5.94M 1s
  5800K .......... .......... .......... .......... .......... 53% 5.89M 1s
  5850K .......... .......... .......... .......... .......... 53% 9.11M 1s
  5900K .......... .......... .......... .......... .......... 54% 6.34M 1s
  5950K .......... .......... .......... .......... .......... 54% 3.46M 1s
  6000K .......... .......... .......... .......... .......... 54% 3.82M 1s
  6050K .......... .......... .......... .......... .......... 55% 6.59M 1s
  6100K .......... .......... .......... .......... .......... 55% 5.84M 1s
  6150K .......... .......... .......... .......... .......... 56% 3.44M 1s
  6200K .......... .......... .......... .......... .......... 56% 4.13M 1s
  6250K .......... .......... .......... .......... .......... 57% 3.88M 1s
  6300K .......... .......... .......... .......... .......... 57% 7.54M 1s
  6350K .......... .......... .......... .......... .......... 58% 5.77M 1s
  6400K .......... .......... .......... .......... .......... 58% 3.54M 1s
  6450K .......... .......... .......... .......... .......... 59% 4.00M 1s
  6500K .......... .......... .......... .......... .......... 59% 4.33M 1s
  6550K .......... .......... .......... .......... .......... 59% 5.69M 1s
  6600K .......... .......... .......... .......... .......... 60% 7.68M 1s
  6650K .......... .......... .......... .......... .......... 60% 3.17M 1s
  6700K .......... .......... .......... .......... .......... 61% 4.53M 1s
  6750K .......... .......... .......... .......... .......... 61% 3.17M 1s
  6800K .......... .......... .......... .......... .......... 62% 3.99M 1s
  6850K .......... .......... .......... .......... .......... 62% 3.80M 1s
  6900K .......... .......... .......... .......... .......... 63% 9.39M 1s
  6950K .......... .......... .......... .......... .......... 63% 3.46M 1s
  7000K .......... .......... .......... .......... .......... 64% 5.58M 1s
  7050K .......... .......... .......... .......... .......... 64% 5.30M 1s
  7100K .......... .......... .......... .......... .......... 64% 3.47M 1s
  7150K .......... .......... .......... .......... .......... 65% 5.90M 1s
  7200K .......... .......... .......... .......... .......... 65% 7.23M 1s
  7250K .......... .......... .......... .......... .......... 66% 2.45M 1s
  7300K .......... .......... .......... .......... .......... 66% 12.5M 1s
  7350K .......... .......... .......... .......... .......... 67% 2.84M 1s
  7400K .......... .......... .......... .......... .......... 67% 4.47M 1s
  7450K .......... .......... .......... .......... .......... 68% 4.06M 1s
  7500K .......... .......... .......... .......... .......... 68% 6.11M 1s
  7550K .......... .......... .......... .......... .......... 69% 5.50M 1s
  7600K .......... .......... .......... .......... .......... 69% 3.64M 1s
  7650K .......... .......... .......... .......... .......... 69% 7.19M 1s
  7700K .......... .......... .......... .......... .......... 70% 4.52M 1s
  7750K .......... .......... .......... .......... .......... 70% 4.88M 1s
  7800K .......... .......... .......... .......... .......... 71% 1.68M 1s
  7850K .......... .......... .......... .......... .......... 71% 5.64M 1s
  7900K .......... .......... .......... .......... .......... 72% 3.52M 1s
  7950K .......... .......... .......... .......... .......... 72% 11.2M 1s
  8000K .......... .......... .......... .......... .......... 73% 1.67M 1s
  8050K .......... .......... .......... .......... .......... 73% 5.81M 1s
  8100K .......... .......... .......... .......... .......... 74% 6.18M 1s
  8150K .......... .......... .......... .......... .......... 74% 7.64M 1s
  8200K .......... .......... .......... .......... .......... 74% 3.59M 1s
  8250K .......... .......... .......... .......... .......... 75% 7.24M 1s
  8300K .......... .......... .......... .......... .......... 75% 3.24M 1s
  8350K .......... .......... .......... .......... .......... 76% 2.80M 1s
  8400K .......... .......... .......... .......... .......... 76% 4.43M 1s
  8450K .......... .......... .......... .......... .......... 77% 8.66M 1s
  8500K .......... .......... .......... .......... .......... 77% 1.87M 1s
  8550K .......... .......... .......... .......... .......... 78% 5.54M 1s
  8600K .......... .......... .......... .......... .......... 78% 8.92M 1s
  8650K .......... .......... .......... .......... .......... 79% 4.75M 1s
  8700K .......... .......... .......... .......... .......... 79% 3.51M 1s
  8750K .......... .......... .......... .......... .......... 79% 6.34M 1s
  8800K .......... .......... .......... .......... .......... 80% 5.55M 1s
  8850K .......... .......... .......... .......... .......... 80% 4.28M 1s
  8900K .......... .......... .......... .......... .......... 81% 7.33M 1s
  8950K .......... .......... .......... .......... .......... 81% 5.61M 0s
  9000K .......... .......... .......... .......... .......... 82% 5.47M 0s
  9050K .......... .......... .......... .......... .......... 82% 1.84M 0s
  9100K .......... .......... .......... .......... .......... 83% 5.92M 0s
  9150K .......... .......... .......... .......... .......... 83% 15.7M 0s
  9200K .......... .......... .......... .......... .......... 84% 2.30M 0s
  9250K .......... .......... .......... .......... .......... 84% 3.15M 0s
  9300K .......... .......... .......... .......... .......... 84% 2.41M 0s
  9350K .......... .......... .......... .......... .......... 85% 6.90M 0s
  9400K .......... .......... .......... .......... .......... 85% 2.50M 0s
  9450K .......... .......... .......... .......... .......... 86% 5.99M 0s
  9500K .......... .......... .......... .......... .......... 86% 4.89M 0s
  9550K .......... .......... .......... .......... .......... 87% 2.96M 0s
  9600K .......... .......... .......... .......... .......... 87% 4.75M 0s
  9650K .......... .......... .......... .......... .......... 88% 5.22M 0s
  9700K .......... .......... .......... .......... .......... 88% 5.54M 0s
  9750K .......... .......... .......... .......... .......... 89% 3.50M 0s
  9800K .......... .......... .......... .......... .......... 89% 4.32M 0s
  9850K .......... .......... .......... .......... .......... 89% 8.45M 0s
  9900K .......... .......... .......... .......... .......... 90% 2.06M 0s
  9950K .......... .......... .......... .......... .......... 90% 5.17M 0s
 10000K .......... .......... .......... .......... .......... 91% 6.83M 0s
 10050K .......... .......... .......... .......... .......... 91% 4.06M 0s
 10100K .......... .......... .......... .......... .......... 92% 5.81M 0s
 10150K .......... .......... .......... .......... .......... 92% 5.64M 0s
 10200K .......... .......... .......... .......... .......... 93% 4.34M 0s
 10250K .......... .......... .......... .......... .......... 93% 5.61M 0s
 10300K .......... .......... .......... .......... .......... 94% 5.41M 0s
 10350K .......... .......... .......... .......... .......... 94% 8.15M 0s
 10400K .......... .......... .......... .......... .......... 94% 3.66M 0s
 10450K .......... .......... .......... .......... .......... 95% 5.64M 0s
 10500K .......... .......... .......... .......... .......... 95% 7.63M 0s
 10550K .......... .......... .......... .......... .......... 96% 1.11M 0s
 10600K .......... .......... .......... .......... .......... 96% 5.05M 0s
 10650K .......... .......... .......... .......... .......... 97% 6.27M 0s
 10700K .......... .......... .......... .......... .......... 97% 4.83M 0s
 10750K .......... .......... .......... .......... .......... 98% 5.03M 0s
 10800K .......... .......... .......... .......... .......... 98% 5.74M 0s
 10850K .......... .......... .......... .......... .......... 99% 4.02M 0s
 10900K .......... .......... .......... .......... .......... 99% 13.4M 0s
 10950K .......... .......... .......... .......... .......... 99% 5.27M 0s
 11000K ..                                                    100% 4848G=2.7sgodel version 0.27.0
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
        id: "com.palantir.godel-generate-plugin:generate-plugin:1.0.0"
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
[extimport]     Running extimport...
[compiles]      Running compiles...
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
[compiles]      Finished compiles
[deadcode]      Finished deadcode
[unconvert]     Running unconvert...
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	0.002s
?   	github.com/nmiyake/echgo/generator       	[no test files]
--- FAIL: TestInvalidType (0.71s)
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
testing.tRunner.func1(0xc4200781e0)
	/usr/local/go/src/testing/testing.go:742 +0x29d
panic(0x5153c0, 0xc42004c090)
	/usr/local/go/src/runtime/panic.go:502 +0x229
github.com/nmiyake/echgo/integration_test_test.TestInvalidType(0xc4200781e0)
	/go/src/github.com/nmiyake/echgo/integration_test/integration_test.go:27 +0x49e
testing.tRunner(0xc4200781e0, 0x5472e0)
	/usr/local/go/src/testing/testing.go:777 +0xd0
created by testing.(*T).Run
	/usr/local/go/src/testing/testing.go:824 +0x2e0
FAIL	github.com/nmiyake/echgo/integration_test	0.715s
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
[extimport]     Running extimport...
[deadcode]      Running deadcode...
[compiles]      Running compiles...
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
[compiles]      Finished compiles
[unconvert]     Running unconvert...
[deadcode]      Finished deadcode
[varcheck]      Running varcheck...
[errcheck]      Finished errcheck
[outparamcheck] Finished outparamcheck
[unconvert]     Finished unconvert
[varcheck]      Finished varcheck
Running test...
?   	github.com/nmiyake/echgo                 	[no test files]
ok  	github.com/nmiyake/echgo/echo            	(cached)
?   	github.com/nmiyake/echgo/generator       	[no test files]
ok  	github.com/nmiyake/echgo/integration_test	2.313s
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
* `dockerctx` directory exists and `godel/config/dist-plugin.yml` is configured to build Docker images for the product
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
