Summary
-------
To add gödel to a project, obtain the gödel distribution and copy the `godelw` file and `godel` directory to the project
directory.

Tutorial start state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory

([Link](https://github.com/nmiyake/echgo/tree/54a23e62a4f9983d60939fdc8ed8dd59f81ddf7c))

Add gödel to a project
----------------------

Add gödel to the project by downloading the distribution and copying the `godelw` script and `godel` directory into it.
This tutorial uses version 0.26.0, but the following steps are applicable for any version.

Download the distribution into a temporary directory and expand it:

```
➜ export GODEL_VERSION=0.26.0
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/$GODEL_VERSION/godel-$GODEL_VERSION.tgz" -o download/godel-"$GODEL_VERSION".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 10.7M  100 10.7M    0     0  3127k      0  0:00:03  0:00:03 --:--:-- 4433k
➜ tar -xf download/godel-"$GODEL_VERSION".tgz -C download
```

Copy the contents of the `wrapper` directory (which contains `godelw` and `godel`) to the project directory:

```
➜ cp -r download/godel-"$GODEL_VERSION"/wrapper/* .
```

Run `./godelw version` to verify that gödel was installed correctly. This command will download the distribution if the
distribution has not previously been installed locally:

```
➜ ./godelw version
Downloading https://palantir.bintray.com/releases/com/palantir/godel/godel/0.26.0/godel-0.26.0.tgz to /Users/nmiyake/.godel/downloads/godel-0.26.0.tgz...
/Users/nmiyake/.godel/downloads/godel-0.26 100%[========================================================================================>]  10.74M  4.38MB/s    in 2.5s
godel version 0.26.0
```

Technically, this is sufficient and we have a working gödel install. However, distributions that are downloaded do not
have a checksum set in `godel/config/godel.properties`:

```
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/0.26.0/godel-0.26.0.tgz
distributionSHA256=
```

For completeness, set the checksum:

```
➜ echo 'distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/0.26.0/godel-0.26.0.tgz
distributionSHA256=c8d086da372e01ab671c57b16ab988fc2cca471906bdefa0bd1711d87883b32e' > godel/config/godel.properties
```

Now that gödel has been added to the project, remove the temporary directory and unset the version variable:

```
➜ rm -rf download
➜ unset GODEL_VERSION
```

Finally, commit the changes to the repository:

```
➜ git add godel godelw
➜ git commit -m "Add godel to project"
[master 6a73370] Add godel to project
 10 files changed, 246 insertions(+)
 create mode 100644 godel/config/check.yml
 create mode 100644 godel/config/dist.yml
 create mode 100644 godel/config/exclude.yml
 create mode 100644 godel/config/format.yml
 create mode 100644 godel/config/generate.yml
 create mode 100644 godel/config/godel.properties
 create mode 100644 godel/config/imports.yml
 create mode 100644 godel/config/license.yml
 create mode 100644 godel/config/test.yml
 create mode 100755 godelw
```

gödel has now been added to the project and is ready to use.

Tutorial end state
--------------------

* `$GOPATH/src/github.com/nmiyake/echgo` exists and is the working directory
* Project contains `godel` and `godelw`

([Link](https://github.com/nmiyake/echgo/tree/6a73370d5b9c8c32ce1a5218938c922f1218db30))

Tutorial next step
------------------

[Add Git hooks to enforce formatting](https://github.com/palantir/godel/wiki/Add-git-hooks)

More
----

### Copying gödel from an existing project

If you have local projects that already use gödel, you can add gödel to a another project by copying the `godelw` and
`godel/config/godel.properties` files from the project that already has gödel.

For example, assume that `$GOPATH/src/github.com/nmiyake/echgo` exists in the current state of the tutorial and we want
to create another project at `$GOPATH/src/github.com/nmiyake/sample` that also uses gödel. This can be done as follows:

```
➜ mkdir -p $GOPATH/src/github.com/nmiyake/sample && cd $_
➜ cp $GOPATH/src/github.com/nmiyake/echgo/godelw .
➜ mkdir -p godel/config
➜ cp $GOPATH/src/github.com/nmiyake/echgo/godel/config/godel.properties godel/config/
➜ ./godelw update
```

Verify that invoking `./godelw` works and that the `godel/config` directory has been populated with the default
configuration files:

```
➜ ./godelw version
godel version 0.26.0
➜ ls godel/config
check.yml        exclude.yml      generate.yml     imports.yml      test.yml
dist.yml         format.yml       godel.properties license.yml
```

Restore the workspace to the original state by setting the working directory back to the `echgo` project and removing
the sample project:

```
➜ cd $GOPATH/src/github.com/nmiyake/echgo
➜ rm -rf $GOPATH/src/github.com/nmiyake/sample
```

The steps above take the approach of copying the `godelw` file and `godel/config/godel.properties` file and running
`update` to ensure that all of the configuration is in its default state -- if the entire `godel/config` directory was
copied from the other project, then the current project would copy the other project's configuration as well, which is
most likely not the correct thing to do.

The `update` task ensures that the version of gödelw in the project matches the one specified by
`godel/config/godel.properties`. If the distribution is already downloaded locally in the local global gödel directory
(`~/.godel` by default) and matches the checksum specified in `godel.properties`, the required files are copied from the
local cache. If the distribution does not exist in the local cache or a checksum is not specified in `godel.properties`,
then the distribution is downloaded from the distribution URL.

### Manually install gödel by moving it to the expected location

In the main tutorial above, we downloaded the gödel distribution using `curl`. However, when we ran `./godelw version`,
that task also downloaded the distribution. In most scenarios this is probably fine, but the download is not technically
required. The second download can be avoided by moving the downloaded distribution to the proper location in cache.

The `godelw` script expects the gödel distribution directory to exist at `$GODEL_HOME/dists/godel-$VERSION` (if
`GODEL_HOME` is undefined, then `~/.godel` is used). If the distribution directory is not present at that location, the
script downloads the distribution from the URL specified in `godel/config/godel.properties` and expands it.

If you have already downloaded the gödel distribution and do not want the `godelw` script to download it again, you can
manually expand the distribution directory to the expected location.

Repeat the initial steps to download the distribution:

```
➜ export GODEL_VERSION=0.26.0
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/$GODEL_VERSION/godel-$GODEL_VERSION.tgz" -o download/godel-"$GODEL_VERSION".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 10.7M  100 10.7M    0     0  12.3M      0 --:--:-- --:--:-- --:--:-- 18.5M
```

Since we have already downloaded and expanded the distribution, run the following to remove the existing state:

```
➜ rm -rf ~/.godel/dists/godel-"$GODEL_VERSION" ~/.godel/downloads/godel-"$GODEL_VERSION".tgz
```

Now, expand the distribution to its expected destination location:

```
➜ mkdir -p ~/.godel/dists
➜ tar -xf download/godel-"$GODEL_VERSION".tgz -C ~/.godel/dists
```

Run `./godelw version` and verify that the distribution is not re-downloaded:

```
➜ ./godelw version
godel version 0.26.0
```

The above is sufficient to ensure that executing `./godelw` for the given version will not re-download the distribution.
If you want to ensure that the `update` command does not re-download the distribution unnecessarily, then move the
distribution `.tgz` file to its expected location in the `~/.godel/downloads` directory:

```
➜ mv download/godel-"$GODEL_VERSION".tgz ~/.godel/downloads/godel-"$GODEL_VERSION".tgz
```

Run the following to clean up our state:

```
➜ rm -rf download
➜ unset GODEL_VERSION
```

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
➜ export GODEL_VERSION=0.26.0
➜ mkdir -p download
➜ curl -L "https://palantir.bintray.com/releases/com/palantir/godel/godel/$GODEL_VERSION/godel-$GODEL_VERSION.tgz" -o download/godel-"$GODEL_VERSION".tgz
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
  0     0    0     0    0     0      0      0 --:--:-- --:--:-- --:--:--     0
100 10.7M  100 10.7M    0     0  9623k      0  0:00:01  0:00:01 --:--:-- 20.9M
```

The checksum can be computed using `openssl` or `shasum` as follows:

```
➜ openssl dgst -sha256 download/godel-"$GODEL_VERSION".tgz
SHA256(download/godel-0.26.0.tgz)= c8d086da372e01ab671c57b16ab988fc2cca471906bdefa0bd1711d87883b32e
➜ shasum -a 256 download/godel-"$GODEL_VERSION".tgz
c8d086da372e01ab671c57b16ab988fc2cca471906bdefa0bd1711d87883b32e  download/godel-0.26.0.tgz
```

Once the digest value is obtained, it should be specified as the value of the `distributionSHA256=` key in
`godel/config/godel.properties`:

```
➜ echo 'distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/0.26.0/godel-0.26.0.tgz
distributionSHA256=c8d086da372e01ab671c57b16ab988fc2cca471906bdefa0bd1711d87883b32e' > godel/config/godel.properties
➜ cat godel/config/godel.properties
distributionURL=https://palantir.bintray.com/releases/com/palantir/godel/godel/0.26.0/godel-0.26.0.tgz
distributionSHA256=c8d086da372e01ab671c57b16ab988fc2cca471906bdefa0bd1711d87883b32e
```

Run the following to clean up our state:

```
➜ rm -rf download
➜ unset GODEL_VERSION
```
