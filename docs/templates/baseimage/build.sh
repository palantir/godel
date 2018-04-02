#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd "$(dirname "$0")" && pwd)
docker build -t godeltutorial:setup ${SCRIPT_DIR}
