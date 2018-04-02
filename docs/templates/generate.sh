#!/usr/bin/env bash
set -euo pipefail

ARG_0=${1:-}
if [ "$ARG_0" = --debug ]; then
  ./debugimage/build.sh
  go run ../../docsgenerator/main.go --base-image "godeltutorial:setup" --tag-prefix godeltutorial --input-dir . --output-dir ../ --start-step=3
  exit 0
fi

./baseimage/build.sh
go run ../../docsgenerator/main.go --base-image "godeltutorial:setup" --tag-prefix godeltutorial --input-dir . --output-dir ../
