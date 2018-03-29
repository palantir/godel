#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)

GODEL_DIR=$(git rev-parse --show-toplevel)
${GODEL_DIR}/godelw dist
GODEL_ARTIFACT_PATH=$(${GODEL_DIR}/godelw artifacts dist godel)
cp "${GODEL_DIR}/${GODEL_ARTIFACT_PATH}" "${SCRIPT_DIR}/"

# Example of adding artifacts for plugin. Uncomment the code to create the distribution for the "distgo" plugin (if it
# is present in ${GOPATH}/src/github.com/palantir/distgo) and publish it into the "repository" directory in this
# directory so it will be included as part of the Docker image. The example can be modified for any other plugin.
#
# Note that this will make the plugin file present in the generated image, but using it requires either modifying the
# core gödel code to use a custom resolver and use the plugin or configuring a custom resolver in the godel.yml of the
# test project.
#
# Specifically, the following needs to be added to the "DefaultResolvers" slice of the "defaultPluginsConfig" variable:
#
# framework/godellauncher/defaulttasks/defaulttasks.go:
#   `/m2/repository/{{GroupPath}}/{{Product}}/{{Version}}/{{Product}}-{{Version}}-{{OS}}-{{Arch}}.tgz`

#################
#PLUGIN_PROJECT_DIR=${GOPATH}/src/github.com/palantir/distgo
#${PLUGIN_PROJECT_DIR}/godelw publish local --path "${SCRIPT_DIR}/repository"
#################

docker build --build-arg GODEL_ARTIFACT_NAME=$(basename ${GODEL_ARTIFACT_PATH}) -t godeltutorial:Add-godel ${SCRIPT_DIR}
