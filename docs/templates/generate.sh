#!/usr/bin/env bash
set -euo pipefail

if [ "$1" = --debug ]; then
  ./debugimage/build.sh
  go run ../../docsgenerator/main.go --base-image "godeltutorial:setup" --tag-prefix godeltutorial --input-dir . --output-dir ../ --start-step=3
  exit 0
fi

docker build ./baseimage
go run ../../docsgenerator/main.go --base-image "godeltutorial:setup" --tag-prefix godeltutorial --input-dir . --output-dir ../
