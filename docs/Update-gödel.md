Updating the version of gödel used by a project involves updating the `distributionURL` property of `godel.properties`
and running the `update` task to update the `godelw` script and any other files.

# Get the distribution URL
Get the distribution URL for the latest version of gödel. This can be obtained from Bintray.

# Get the SHA-256 checksum (optional)
The SHA-256 checksum can be optionally specified to verify the integrity of the downloaded distribution. This can be
obtained from the "Checksums" section of the Bintray page for the distribution:

![SHA checksum](images/add_to_project/sha_checksum.png)

# Update the properties in godel.properties
Open the `godel/godel.properties` file and update the value of `distributionURL` to be the new distribution URL. If a
checksum is going to be used, update the value of `distributionSHA256` to be the checksum for the distribution. If a
checksum is not going to be used and was previously specified, remove the value (the value can be set to blank or the
entire key/value entry can be removed).

# Run the update task
Run `./godelw update`. This will download the new distribution if necessary and will update the `godelw` script in the
project and add any new files to the `godel` directory that are included in the distribution.

# Check the changes into version control
Commit the updated `godelw` script and `godel` configuration files and push to version control.
